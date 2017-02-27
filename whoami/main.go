// whoami displays session information
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"v2.staffjoy.com/account"
	"v2.staffjoy.com/auth"
	"v2.staffjoy.com/company"
	"v2.staffjoy.com/crypto"
	"v2.staffjoy.com/environments"
	"v2.staffjoy.com/errorpages"
	"v2.staffjoy.com/healthcheck"
	"v2.staffjoy.com/middlewares"

	"github.com/urfave/negroni"
)

var (
	logger *logrus.Entry
	config environments.Config
)

const (
	// ServiceName is how we refer to this app in logs
	ServiceName = "whoami"
)

func init() {
	var err error

	// Set the ENV environment variable to control dev/stage/prod behavior
	config, err = environments.GetConfig(os.Getenv(environments.EnvVar))
	if err != nil {
		panic("Unable to determine configuration")
	}
	logger = config.GetLogger(ServiceName)

	logger.Debugf("Initialized whami %s environment", config.Name)
}

// NewRouter builds the mux router for the site
// (abstracted for testing purposes)
func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc(healthcheck.HEALTHPATH, healthcheck.Handler)
	r.HandleFunc("/", whoamiHandler)
	r.HandleFunc("/intercom/", intercomHandler)
	r.NotFoundHandler = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		errorpages.NotFound(res)
	})

	return r
}

type iAm struct {
	Support  bool                  `json:"support"`
	UserUUID string                `json:"user_uuid"`
	Worker   *company.WorkerOfList `json:"worker"`
	Admin    *company.AdminOfList  `json:"admin"`
}

type intercomSettings struct {
	AppID     string `json:"app_id"`
	UserUUID  string `json:"user_id"`
	UserHash  string `json:"user_hash"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
}

func whoamiHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	payload := iAm{}
	header := http.StatusOK
	var err error

	switch req.Header.Get(auth.AuthorizationHeader) {
	case auth.AuthorizationAnonymousWeb:
		header = http.StatusForbidden
	case auth.AuthorizationAuthenticatedUser:
		payload.UserUUID, err = auth.GetCurrentUserUUIDFromHeader(req.Header)
		if err != nil {
			panic(err)
		}
	case auth.AuthorizationSupportUser:
		payload.Support = true
		payload.UserUUID, err = auth.GetCurrentUserUUIDFromHeader(req.Header)
		if err != nil {
			panic(err)
		}
	default:
		logger.Panicf("unknown authorization header %v", req.Header.Get(auth.AuthorizationHeader))
	}

	if payload.UserUUID != "" {
		// Get worker stuff
		svc, close, err := company.NewClient()
		if err != nil {
			panic(err)
		}
		defer close()

		md := metadata.New(map[string]string{auth.AuthorizationMetadata: auth.AuthorizationWhoamiService})
		ctx, cancel := context.WithCancel(metadata.NewContext(context.Background(), md))
		defer cancel()
		if payload.Worker, err = svc.GetWorkerOf(ctx, &company.WorkerOfRequest{UserUuid: payload.UserUUID}); err != nil {
			panic(err)
		}
		if payload.Admin, err = svc.GetAdminOf(ctx, &company.AdminOfRequest{UserUuid: payload.UserUUID}); err != nil {
			panic(err)
		}
	}
	encodedData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	res.WriteHeader(header)
	res.Header().Set("Content-Type", "application/json")
	res.Write(encodedData)
}

func main() {
	r := NewRouter()
	n := negroni.New()

	sentryPublicDSN, err := environments.GetPublicSentryDSN(config.GetSentryDSN())
	if err != nil {
		logger.Fatalf("Cannot get sentry info - %s", err)
	}
	n.Use(middlewares.NewRecovery(ServiceName, config, sentryPublicDSN))
	n.UseHandler(r)

	s := &http.Server{
		Addr:           ":80",
		Handler:        n,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.Panicf("%s", s.ListenAndServe())
}

func intercomHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	payload := intercomSettings{}
	header := http.StatusOK
	var err error
	payload.AppID = os.Getenv("INTERCOM_APP_ID")

	switch req.Header.Get(auth.AuthorizationHeader) {
	case auth.AuthorizationAnonymousWeb:
		header = http.StatusForbidden
	case auth.AuthorizationSupportUser:
		fallthrough
	case auth.AuthorizationAuthenticatedUser:
		payload.UserUUID, err = auth.GetCurrentUserUUIDFromHeader(req.Header)
		if err != nil {
			panic(err)
		}
		payload.UserHash = crypto.ComputeHmac256(payload.UserUUID, os.Getenv("INTERCOM_SIGNING_SECRET"))
	default:
		logger.Panicf("unknown authorization header %v", req.Header.Get(auth.AuthorizationHeader))
	}

	if payload.UserUUID != "" {
		// Get user account info so we can fill in name and email

		svc, close, err := account.NewClient()
		if err != nil {
			panic(err)
		}
		defer close()

		md := metadata.New(map[string]string{auth.AuthorizationMetadata: auth.AuthorizationWhoamiService})
		ctx, cancel := context.WithCancel(metadata.NewContext(context.Background(), md))
		defer cancel()

		var u *account.Account
		if u, err = svc.Get(ctx, &account.GetAccountRequest{Uuid: payload.UserUUID}); err != nil {
			panic(err)
		}
		payload.Name = u.Name
		payload.Email = u.Email
		payload.CreatedAt = int64(u.MemberSince.Unix())

	}

	encodedData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	res.WriteHeader(header)
	res.Header().Set("Content-Type", "application/json")
	res.Write(encodedData)
}
