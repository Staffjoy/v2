// code is a web service that makes v2.staffjoy.com a custom
// go import doman, e.g. it lets `go get v2.staffjoy.com` function
package main

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"v2.staffjoy.com/account"
	"v2.staffjoy.com/auth"
	"v2.staffjoy.com/environments"
	"v2.staffjoy.com/errorpages"
	"v2.staffjoy.com/healthcheck"
	"v2.staffjoy.com/middlewares"
)

const (
	addr = 80
	// ServiceName is how this service is identified in logs
	ServiceName = "superpowers"
)

var (
	wildcard bool
	logger   *logrus.Entry
	config   environments.Config
	c        account.AccountServiceClient
)

func init() {
	var err error

	// Set the ENV environment variable to control dev/stage/prod behavior
	config, err = environments.GetConfig(os.Getenv(environments.EnvVar))
	if err != nil {
		panic("Unable to determine configuration")
	}
	logger = config.GetLogger(ServiceName)

	logger.Debugf("initialized superpowers %s environment", config.Name)
}

// NewRouter builds the mux router for the site
// (abstracted for testing purposes)
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", infoHandler).Methods("GET")
	r.HandleFunc("/", superpowerHandler).Methods("POST")
	r.HandleFunc(healthcheck.HEALTHPATH, healthcheck.Handler)
	r.NotFoundHandler = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		errorpages.NotFound(res)
	})
	return r
}

func main() {
	n := negroni.New()
	n.Use(middlewares.NewRecovery(ServiceName, config, ""))
	n.UseHandler(NewRouter())

	conn, err := grpc.Dial(account.Endpoint, grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c = account.NewAccountServiceClient(conn)

	s := &http.Server{
		Addr:           ":80",
		Handler:        n,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logger.Panicf("%s", s.ListenAndServe())
}

var supportTmpl = template.Must(template.New("main").Parse(`<!DOCTYPE html>
<html>
<head>
</head>
<body>
<h1>You already have superpowers</h1>
<p>Go make the world a better place. And remember - no capes.</p>
</body>
</html>
`))

var activateTmpl = template.Must(template.New("main").Parse(`<!DOCTYPE html>
<html>
<head>
</head>
<body>
<h1>Get Superpowers</h1>
<p>Click below and you will be given superpowers. (You'll have to log back in afterward).
<form action="/" method="post">
	<input type="submit" value="Find your power" />
</form>
</body>
</html>
`))

func infoHandler(res http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	if req.Header.Get(auth.AuthorizationHeader) == auth.AuthorizationSupportUser {
		if err := supportTmpl.Execute(&buf, nil); err != nil {
			panic("cannot render support template")
		}
	} else {
		if err := activateTmpl.Execute(&buf, nil); err != nil {
			panic("cannot render support template")
		}
	}
	res.Write(buf.Bytes())
}

// superpowerHandler grants superpowers to the current user
func superpowerHandler(res http.ResponseWriter, req *http.Request) {
	var uuid string
	var err error
	if uuid, err = auth.GetCurrentUserUUIDFromHeader(req.Header); err != nil {
		panic("Could not get user id")
	}
	md := metadata.New(map[string]string{auth.AuthorizationMetadata: auth.AuthorizationSuperpowersService})
	ctx, cancel := context.WithCancel(metadata.NewContext(context.Background(), md))
	defer cancel()
	a, err := c.Get(ctx, &account.GetAccountRequest{Uuid: uuid})
	if err != nil {
		panic(err)
	}
	a.Support = true
	_, err = c.Update(ctx, a)
	if err != nil {
		panic(err)
	}
	logger.Infof("Superpowers granted to user %s", uuid)
	auth.Logout(res)
	http.Redirect(res, req, "http://www."+config.ExternalApex+"/login/", http.StatusTemporaryRedirect)
}
