package main

import (
	"net/http"

	"google.golang.org/grpc/metadata"

	"v2.staffjoy.com/account"
	"v2.staffjoy.com/auth"

	"golang.org/x/net/context"
)

const (
	signUpPath = "/sign-up/"
	loginPath  = "/login/"
)

func signUpHandler(res http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name") // not everything sends this
	email := req.FormValue("email")
	if len(email) <= 0 {
		http.Redirect(res, req, signUpPath, http.StatusFound)
	}
	md := metadata.New(map[string]string{auth.AuthorizationMetadata: auth.AuthorizationWWWService})
	ctx, cancel := context.WithCancel(metadata.NewOutgoingContext(context.Background(), md))
	defer cancel()

	accountClient, close, err := account.NewClient()
	if err != nil {
		panic(err)
	}
	defer close()

	a, err := accountClient.Create(ctx, &account.CreateAccountRequest{Name: name, Email: email})
	if err != nil {
		// TODO - check if user exists, and send a reset link
		logger.Infof("Failed to create account - %v", err)
		http.Redirect(res, req, signUpPath, http.StatusFound)
		return
	}
	logger.Infof("New account signup - %v", a)
	if err = tmpl.ExecuteTemplate(res, confirmPage.TemplateName, confirmPage); err != nil {
		panic(err)
	}
}
