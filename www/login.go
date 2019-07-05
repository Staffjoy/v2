package main

import (
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/grpc/metadata"
	"v2.staffjoy.com/faraday/services"

	"v2.staffjoy.com/account"
	"v2.staffjoy.com/auth"
	"v2.staffjoy.com/helpers"
	"v2.staffjoy.com/suite"

	"github.com/gorilla/csrf"
	"github.com/sirupsen/logrus"
)

type loginPage struct {
	Title         string // Used in <title>
	CSSId         string // e.g. 'careers'
	Version       string // e.g. master-1, for cachebusting
	CsrfField     template.HTML
	Denied        bool
	Description   string
	TemplateName  string
	PreviousEmail string
	ReturnTo      string
}

// isValidSub returns true if url contains valid subdomain
func isValidSub(sub string) bool {
	u, err := url.Parse(sub)
	if err != nil {
		logger.Errorf("can't parse url %v", err)
		return false
	}

	bare := strings.Replace(u.Host, "."+config.ExternalApex, "", -1)
	for k := range services.StaffjoyServices {
		if k == bare {
			return true
		}
	}
	return false
}

func loginHandler(res http.ResponseWriter, req *http.Request) {

	// for GET
	returnTo := req.URL.Query().Get("return_to")

	p := loginPage{
		Title:        "Staffjoy Log in",
		CSSId:        "login",
		CsrfField:    csrf.TemplateField(req),
		Version:      config.GetDeployVersion(),
		TemplateName: "login.tmpl",
		Description:  "Log in to Staffjoy to start scheduling your workers. All you’ll need is your email and password.",
		ReturnTo:     returnTo,
	}

	// if logged in - go away
	if req.Header.Get(auth.AuthorizationHeader) != auth.AuthorizationAnonymousWeb {
		destination := &url.URL{Host: "myaccount." + config.ExternalApex, Scheme: "http"}
		http.Redirect(res, req, destination.String(), http.StatusFound)
	}

	if req.Method == http.MethodPost {
		email := req.FormValue("email")
		password := req.FormValue("password")
		// POST and GET are in the same handler - reset
		returnTo = req.FormValue("return_to")
		// rememberMe=True means that the session is set for a month instead of a day
		rememberMe := len(req.FormValue("remember-me")) > 0

		md := metadata.New(map[string]string{auth.AuthorizationMetadata: auth.AuthorizationWWWService})
		ctx, cancel := context.WithCancel(metadata.NewOutgoingContext(context.Background(), md))
		defer cancel()

		accountClient, close, err := account.NewClient()
		if err != nil {
			panic(err)
		}
		defer close()

		a, err := accountClient.VerifyPassword(ctx, &account.VerifyPasswordRequest{Email: email, Password: password})
		if err == nil {
			auth.LoginUser(a.Uuid, a.Support, rememberMe, res)
			go helpers.TrackEvent(a.Uuid, "login")
			go accountClient.SyncUser(ctx, &account.SyncUserRequest{Uuid: a.Uuid})

			logger.WithFields(logrus.Fields{"user_uuid": a.Uuid}).Info("Logging in user")

			scheme := "https"
			if config.Name == "development" || config.Name == "test" {
				scheme = "http"
			}

			if returnTo == "" {
				destination := &url.URL{Host: "app." + config.ExternalApex, Scheme: scheme}
				returnTo = destination.String()
			} else {
				returnTo = "http://" + returnTo

				// sanitize
				if !isValidSub(returnTo) {
					destination := &url.URL{Host: "myaccount." + config.ExternalApex, Scheme: scheme}
					returnTo = destination.String()
				}

			}

			http.Redirect(res, req, returnTo, http.StatusFound)
			return
		}

		// Check suite.staffjoy.com for an account
		if email != "" {
			suiteAccountExists, err := suite.AccountExists(email)
			if err != nil {
				// Log but don't crash
				// (which is important for dev too)
				logger.Warningf("Unable to query suite for login attempt by %v - %v", email, err)
			} else if suiteAccountExists {
				u, ok := suite.SuiteConfigs[config.Name]
				if !ok {
					panic("failed to generate suite url")
				}
				u.Path = "/auth/login"
				logger.Infof("Redirecting email %s to suite", email)
				http.Redirect(res, req, u.String(), http.StatusTemporaryRedirect)
				return
			}
		}

		logger.WithFields(logrus.Fields{"email": email}).Infof("Login attempt denied")
		p.Denied = true
		p.PreviousEmail = email
	}
	err := tmpl.ExecuteTemplate(res, p.TemplateName, p)
	if err != nil {
		panic(err)
	}

}
