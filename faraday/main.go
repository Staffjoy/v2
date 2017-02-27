// Faraday proxies all requests to Staffjoy
package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"v2.staffjoy.com/faraday/services"

	"v2.staffjoy.com/auth"
	"v2.staffjoy.com/environments"
	"v2.staffjoy.com/errorpages"
	"v2.staffjoy.com/healthcheck"
	"v2.staffjoy.com/middlewares"
)

type contextKey int // Used for gorilla/context

const (
	// ServiceName is how this app is identified in logs and error handlers
	ServiceName          string     = "faraday"
	userID               contextKey = iota // Used for gorilla/context
	userSudo             contextKey = iota // Used for gorilla/context
	requestAuthenticated contextKey = iota // Used for gorilla/context
	requestedService     contextKey = iota // Used for gorilla/context
)

var (
	logger       *logrus.Entry
	config       environments.Config
	signingToken = os.Getenv("SIGNING_SECRET")
	bannedUsers  = map[string]string{ // Use a map for constant time lookups. Value doesn't matter
		// Hypothetically these should be universally unique, so we don't have to limit by env
		"d7b9dbed-9719-4856-5f19-23da2d0e3dec": "hidden",
	}
)

// Setup environment, logger, etc
func init() {
	// Set the ENV environment variable to control dev/stage/prod behavior
	var err error
	config, err = environments.GetConfig(os.Getenv(environments.EnvVar))
	if err != nil {
		panic("Unable to determine configuration")
	}
	logger = config.GetLogger(ServiceName)
}

// Listen for incoming requests, then validate, sanitize, and route them.
func main() {
	logger.Infof("Initialized environment %s", config.Name)

	r := NewRouter(config, logger)
	// Set up http server
	// Note - we do this without the Negroni convenience func so that
	// we can add in TLS support in the future too.
	s := &http.Server{
		Addr:           ":80",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// TODO - add in a logging system and have it do a fatal call here
	logger.Panicf("%v", s.ListenAndServe())
}

// NewRouter returns a router composed of internal and external parts
func NewRouter(config environments.Config, logger *logrus.Entry) http.Handler {
	// Create a new router. We use Gorilla instead of stdlib because it handles
	// memory clean up for the 'context' package correctly
	externalRouter := mux.NewRouter()
	internalRouter := mux.NewRouter().PathPrefix("/").Subrouter().StrictSlash(true)

	// Make this available always, e.g. for kubernetes health checks
	externalRouter.HandleFunc(healthcheck.HEALTHPATH, healthcheck.Handler)
	externalRouter.HandleFunc(MobileConfigPath, MobileConfigHandler)

	sentryPublicDSN, err := environments.GetPublicSentryDSN(config.GetSentryDSN())
	if err != nil {
		logger.Fatalf("Cannot get sentry info - %s", err)
	}

	traceMW, err := NewTraceMiddleware(logger, config)
	if err != nil {
		logger.Fatalf("Unable to load trace middleware - %v", err)
	}

	// only apply security to the internal routes
	externalRouter.PathPrefix("/").Handler(negroni.New(
		middlewares.NewRecovery(ServiceName, config, sentryPublicDSN),
		NewSecurityMiddleware(config),
		NewServiceMiddleware(config, services.StaffjoyServices),
		traceMW,
		NewRobotstxtMiddleware(config),
		negroni.Wrap(internalRouter),
	))
	internalRouter.PathPrefix("/").HandlerFunc(proxyHandler)

	return externalRouter
}

// HTTP function that handles proxying after all of the middlewares
func proxyHandler(res http.ResponseWriter, req *http.Request) {
	service := context.Get(req, requestedService).(services.Service)
	// No security on backend right now :-(
	destination := "http://" + service.BackendDomain + req.URL.RequestURI()
	logger.Debugf("Proxying to %s", destination)
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		panic(fmt.Sprintf("Could not read request body - %s", err))
	}

	internalReq, err := http.NewRequest(req.Method, destination, bytes.NewReader(b))
	if err != nil {
		panic(fmt.Sprintf("Unable to create request - %s", err))
	}

	auth.SetInternalHeaders(req, internalReq.Header)

	currentUserUUID, err := auth.GetCurrentUserUUIDFromHeader(internalReq.Header)
	if err == nil {
		// authenticated request
		if _, isBanned := bannedUsers[currentUserUUID]; isBanned {
			logger.Warningf("Banned user accessing service - user %s", currentUserUUID)
			errorpages.Forbidden(res)
			return
		}
	}

	// Right here - check response Authorization and see if it's ok
	// with the requested service

	// Check perimeter authorization
	switch internalReq.Header.Get(auth.AuthorizationHeader) {
	case auth.AuthorizationAnonymousWeb:
		if service.Security != services.Public {
			// send to login
			scheme := "https"
			if config.Name == "development" || config.Name == "test" {
				scheme = "http"
			}
			redirectDest := &url.URL{Host: "www." + config.ExternalApex, Scheme: scheme, Path: "/login/"}

			url := req.Host + req.URL.EscapedPath()

			http.Redirect(res, req, redirectDest.String()+"?return_to="+url, http.StatusTemporaryRedirect)
			return
		}
	case auth.AuthorizationAuthenticatedUser:
		if service.Security == services.Admin {
			errorpages.Forbidden(res)
			return
		}
	case auth.AuthorizationSupportUser:
		// no restrictions
	default:
		logger.Panicf("unknown authorization header")
	}

	client := http.Client{
		// RETURN a redirect, do not FOLLOW it (which ends up causing relative redirect issues)
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	internalRes, err := client.Do(internalReq)
	if err != nil {
		logger.Warningf("Unable to query backend - %s", err)
		errorpages.GatewayTimeout(res)
		return
	}
	// Copy headers from service to user
	auth.ProxyHeaders(internalRes.Header, res.Header())

	if service.NoCacheHTML {
		if strings.Contains(strings.Join(res.Header()["Content-Type"], ""), "text/html") {
			// insert header to prevent caching
			res.Header().Set("Cache-Control", "no-cache")
		}
	}

	res.WriteHeader(internalRes.StatusCode)
	io.Copy(res, internalRes.Body)
	internalRes.Body.Close()

}
