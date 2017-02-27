package main

import (
	"fmt"
	"net/url"
	"time"

	"golang.org/x/net/context"

	"github.com/golang/protobuf/ptypes/empty"
	"v2.staffjoy.com/account"
	"v2.staffjoy.com/bot"
	"v2.staffjoy.com/company"
	"v2.staffjoy.com/sms"
)

func (s *botServer) OnboardWorker(ctx context.Context, req *bot.OnboardWorkerRequest) (*empty.Empty, error) {
	botCtx := botContext()

	accountClient, close, err := account.NewClient()
	if err != nil {
		return nil, s.internalError(err, "unable to initiate account connection")
	}
	defer close()

	a, err := accountClient.Get(botCtx, &account.GetAccountRequest{Uuid: req.UserUuid})
	if err != nil {
		return nil, s.internalError(err, "cannot find user")
	}
	u := user(*a)

	companyClient, close, err := company.NewClient()
	if err != nil {
		return nil, s.internalError(err, "unable to initiate company connection")
	}
	defer close()

	c, err := companyClient.GetCompany(botCtx, &company.GetCompanyRequest{Uuid: req.CompanyUuid})
	if err != nil {
		return nil, s.internalError(err, "cannot find company")
	}

	d := u.PreferredDispatch()

	switch d {
	case dispatchSms:
		go s.smsOnboard(a, c)
	case dispatchEmail:
		s.logger.Warningf("Email dispatch not implemented - skiping user %s", req.UserUuid)
	default:
		s.logger.Infof("Unable to onboard user %s - no comm method found", req.UserUuid)

	}
	return &empty.Empty{}, nil
}

func (s *botServer) smsOnboard(a *account.Account, c *company.Company) {
	u := user(*a)
	ctx := botContext()
	icalURL := url.URL{Host: "ical." + config.ExternalApex, Path: fmt.Sprintf("/%s.ics", u.Uuid), Scheme: config.Scheme}

	var onboardingMessages = []string{
		fmt.Sprintf("%s Your manager just added you to %s on Staffjoy to share your work schedule.", u.Greet(), c.Name),
		"When your manager publishes your shifts, we'll send them to you here. (To disable Staffjoy messages, reply STOP at any time)",
		fmt.Sprintf("Click this link to sync your shifts to your phone's calendar app: %s", icalURL.String()),
		//"Reply HELP now to see what you can do with Staffjoy",
	}

	smsClient, close, err := sms.NewClient()
	if err != nil {
		s.internalError(err, "unable to initiate sms connection")
		return
	}
	defer close()

	for _, m := range onboardingMessages {
		if _, err := smsClient.QueueSend(ctx, &sms.SmsRequest{To: a.Phonenumber, Body: m}); err != nil {
			s.internalError(err, "could not send welcome sms")
		}
		time.Sleep(4 * time.Second)
	}
	// todo - check if upcoming shifts, and if there are - send them
	s.logger.Infof("onboarded worker %s (%s) for company %s (%s)", a.Uuid, a.Name, c.Uuid, c.Name)
}
