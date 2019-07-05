package main

import (
	"html/template"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"google.golang.org/grpc/metadata"

	"v2.staffjoy.com/account"
	"v2.staffjoy.com/auth"

	recaptcha "github.com/dpapathanasiou/go-recaptcha"
	"github.com/gorilla/csrf"
	"github.com/sirupsen/logrus"
)

type resetPage struct {
	Title           string // Used in <title>
	CSSId           string // e.g. 'careers'
	Version         string // e.g. master-1, for cachebusting
	CsrfField       template.HTML
	Denied          bool
	Description     string
	TemplateName    string
	RecaptchaPublic string
}

func resetHandler(res http.ResponseWriter, req *http.Request) {
	p := resetPage{
		Title:           "Password Reset",
		CSSId:           "sign-up",
		CsrfField:       csrf.TemplateField(req),
		Version:         config.GetDeployVersion(),
		TemplateName:    "reset.tmpl",
		Description:     "Reset the password for your Staffjoy account.",
		RecaptchaPublic: os.Getenv("RECAPTCHA_PUBLIC"),
	}
	if req.Method == http.MethodPost {
		email := req.FormValue("email")
		recaptcha.Init(os.Getenv("RECAPTCHA_PRIVATE"))
		recaptchaResponse, ok := req.Form["g-recaptcha-response"]
		if !ok {
			res.Write([]byte("Recaptcha absent"))
			return
		}

		var remoteIP string
		if config.Debug {
			remoteIP = req.RemoteAddr
		} else {
			// Cloudflare proxy
			remoteIP = req.Header.Get("CF-Connecting-IP")
		}
		result, err := recaptcha.Confirm(remoteIP, recaptchaResponse[0])
		if !result {
			res.Write([]byte("Recaptcha incorrect."))
			return
		}

		md := metadata.New(map[string]string{auth.AuthorizationMetadata: auth.AuthorizationWWWService})
		ctx, cancel := context.WithCancel(metadata.NewOutgoingContext(context.Background(), md))
		defer cancel()

		accountClient, close, err := account.NewClient()
		if err != nil {
			panic(err)
		}
		defer close()

		_, err = accountClient.RequestPasswordReset(ctx, &account.PasswordResetRequest{Email: email})
		if err == nil {
			logger.WithFields(logrus.Fields{"email": email}).Infof("Initiating password reset")
		} else {
			logger.WithFields(logrus.Fields{"email": email}).Infof("Failed password reset - %v", err)
		}
		if err = tmpl.ExecuteTemplate(res, resetConfirmPage.TemplateName, resetConfirmPage); err != nil {
			panic(err)
		}
		return
	}
	err := tmpl.ExecuteTemplate(res, p.TemplateName, p)
	if err != nil {
		panic(err)
	}
}
