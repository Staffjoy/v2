package main

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"v2.staffjoy.com/apidocs"
	"v2.staffjoy.com/company"
	"v2.staffjoy.com/environments"
	"v2.staffjoy.com/healthcheck"
)

const (
	// ServiceName identifies this app in logs
	ServiceName = "company"
	// swagger is the name of the swagger file
	swagger = "company.swagger.json"
)

var (
	logger *logrus.Entry
	config environments.Config
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

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()

	mux.HandleFunc("/swagger.json", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		data, err := Asset(swagger)
		if err != nil {
			panic("Unable to load swagger")
		}
		res.Write(data)
	})

	mux.HandleFunc(healthcheck.HEALTHPATH, healthcheck.Handler)
	apidocs.Serve(mux, logger)

	// Custom runtime option to emit empty fields (like false bools)
	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := RegisterCompanyServiceHandlerFromEndpoint(ctx, gwmux, company.Endpoint, opts); err != nil {
		return err
	}
	mux.Handle("/", gwmux)

	return http.ListenAndServe(":80", mux)
}

func main() {
	logger.Debugf("Initialized companyapi environment %s", config.Name)

	if err := run(); err != nil {
		logger.Fatal(err)
	}
}
