package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/csrf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"v2.staffjoy.com/account"
	"v2.staffjoy.com/auth"
	"v2.staffjoy.com/company"
	"v2.staffjoy.com/email"
	"v2.staffjoy.com/helpers"

	"golang.org/x/net/context"
)

const (
	defaultTimezone      = "UTC"
	defaultDayWeekStarts = "monday"
	defaultTeamName      = "Team"
	newCompanyTmpl       = "new_company.tmpl"
	defaultTeamColor     = "744fc6"
)

func newCompanyHandler(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get(auth.AuthorizationHeader) == auth.AuthorizationAnonymousWeb {
		http.Redirect(res, req, "/login/", http.StatusFound)
	}
	name := req.FormValue("name") // not everything sends this
	timezone := req.FormValue("timezone")
	team := req.FormValue("team")
	if name != "" {
		if timezone == "" {
			timezone = defaultTimezone
		}
		if team == "" {
			team = defaultTeamName
		}

		md := metadata.New(map[string]string{auth.AuthorizationMetadata: auth.AuthorizationWWWService})
		ctx, cancel := context.WithCancel(metadata.NewOutgoingContext(context.Background(), md))
		defer cancel()

		// fetch current user Infof
		currentUserUUID, err := auth.GetCurrentUserUUIDFromHeader(req.Header)
		if err != nil {
			panic(err)
		}
		accountClient, close, err := account.NewClient()
		if err != nil {
			panic(err)
		}
		defer close()

		currentUser, err := accountClient.Get(ctx, &account.GetAccountRequest{Uuid: currentUserUUID})

		companyClient, companyClose, err := company.NewClient()
		if err != nil {
			panic(err)
		}
		defer companyClose()

		// Make the company
		c, err := companyClient.CreateCompany(ctx, &company.CreateCompanyRequest{Name: name, DefaultTimezone: timezone, DefaultDayWeekStarts: defaultDayWeekStarts})
		if codes.InvalidArgument == grpc.Code(err) {
			// retry with default timezone
			if c, err = companyClient.CreateCompany(ctx, &company.CreateCompanyRequest{Name: name, DefaultTimezone: defaultTimezone, DefaultDayWeekStarts: defaultDayWeekStarts}); err != nil {
				panic(err)
			}
		} else if err != nil {
			panic(err)
		}

		// register current user in directory
		if _, err = companyClient.CreateDirectory(ctx, &company.NewDirectoryEntry{CompanyUuid: c.Uuid, Email: currentUser.Email}); err != nil {
			panic(err)
		}

		// create admin
		if _, err := companyClient.CreateAdmin(ctx, &company.DirectoryEntryRequest{CompanyUuid: c.Uuid, UserUuid: currentUserUUID}); err != nil {
			panic(err)
		}

		// create team
		team, err := companyClient.CreateTeam(ctx, &company.CreateTeamRequest{CompanyUuid: c.Uuid, Name: team, Color: defaultTeamColor})
		if err != nil {
			panic(err)
		}

		// register as worker
		if _, err = companyClient.CreateWorker(ctx, &company.Worker{CompanyUuid: c.Uuid, TeamUuid: team.Uuid, UserUuid: currentUserUUID}); err != nil {
			panic(err)
		}

		// redirect
		logger.Infof("new company signup - %v", c)
		url := url.URL{
			Scheme: "http",
			Host:   "app." + config.ExternalApex,
		}
		go accountClient.SyncUser(ctx, &account.SyncUserRequest{Uuid: currentUser.Uuid})
		go helpers.TrackEvent(currentUserUUID, "freetrial_created")

		if config.Name == "production" && !currentUser.Support {
			// Alert sales of a new account signup
			go func(a *account.Account, c *company.Company) {
				msg := &email.EmailRequest{
					To:       "sales@staffjoy.com",
					Name:     "",
					Subject:  fmt.Sprintf("%s from %s just joined Staffjoy", a.Name, c.Name),
					HtmlBody: fmt.Sprintf("Name: %s<br>Phone: %s<br>Email: %s<br>Company: %s<br>App: https://app.staffjoy.com/#/companies/%s/employees/", a.Name, a.Phonenumber, a.Email, c.Name, c.Uuid),
				}
				mailer, close, err := email.NewClient()
				if err != nil {
					logger.Errorf("unable to initiate email service connection - %s", err)
					return
				}
				defer close()

				ctx, cancel := context.WithCancel(metadata.NewOutgoingContext(context.Background(), md))
				defer cancel()

				if _, err = mailer.Send(ctx, msg); err != nil {
					logger.Errorf("Unable to send email - %s", err)
					return
				}
			}(currentUser, c)
		}
		http.Redirect(res, req, url.String(), http.StatusFound)

	}

	newCompanyPage.CsrfField = csrf.TemplateField(req)
	if err := tmpl.ExecuteTemplate(res, newCompanyPage.TemplateName, newCompanyPage); err != nil {
		panic(err)
	}
}
