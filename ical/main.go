// ical is a web service that makes ical.staffjoy.com a custom
// go import doman, e.g. it lets `go get v2.staffjoy.com` function
package main

import (
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
	"v2.staffjoy.com/auth"
	"v2.staffjoy.com/company"
	pb "v2.staffjoy.com/company"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"google.golang.org/grpc/metadata"
	"v2.staffjoy.com/environments"
	"v2.staffjoy.com/errorpages"
	"v2.staffjoy.com/healthcheck"
	"v2.staffjoy.com/middlewares"
)

const (
	addr = 80

	// ServiceName is how this service is identified in logs
	ServiceName = "ical"
)

var (
	wildcard bool
	logger   *logrus.Entry
	config   environments.Config
)

func init() {
	var err error

	// Set the ENV environment variable to control dev/stage/prod behavior
	config, err = environments.GetConfig(os.Getenv(environments.EnvVar))
	if err != nil {
		panic("Unable to determine configuration")
	}
	logger = config.GetLogger(ServiceName)

	logger.Debugf("initialized ical %s environment", config.Name)
}

func icalContext() context.Context {
	md := metadata.New(map[string]string{auth.AuthorizationMetadata: auth.AuthorizationICalService})
	return metadata.NewOutgoingContext(context.Background(), md)
}

// NewRouter builds the mux router for the site
// (abstracted for testing purposes)
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/{user_uuid}.ics", CalHandler).Methods("GET")
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

	s := &http.Server{
		Addr:           ":80",
		Handler:        n,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logger.Panicf("%s", s.ListenAndServe())
}

// CalHandler returns a basic ical for user given their UUID
// it currently has no auth
func CalHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	old := &pb.Worker{UserUuid: vars["user_uuid"]}

	ctx := icalContext()
	companyClient, close, err := company.NewClient()
	if err != nil {
		logger.Debugf("unable to initiate company connection")
	}
	defer close()

	tinfo, err := companyClient.GetWorkerTeamInfo(ctx, old)
	if err != nil {
		logger.Debugf("unable to get team info %s", err)
	}

	wsl := &pb.WorkerShiftListRequest{
		CompanyUuid:      tinfo.CompanyUuid,
		TeamUuid:         tinfo.TeamUuid,
		WorkerUuid:       tinfo.UserUuid,
		ShiftStartAfter:  time.Now().AddDate(0, -1, 0),
		ShiftStartBefore: time.Now().AddDate(0, 3, 0),
	}

	company, err := companyClient.GetCompany(ctx, &pb.GetCompanyRequest{Uuid: tinfo.CompanyUuid})
	if err != nil {
		logger.Debugf("unable to get company %s", err)
	}

	shifts, err := companyClient.ListWorkerShifts(ctx, wsl)
	if err != nil {
		logger.Debugf("unable to get worker shifts %s", err)
	}

	cal := Cal{
		Shifts:  shifts.Shifts,
		Company: company.Name,
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/calendar; charset=utf-8")
	res.Header().Set("Content-Disposition", "attachment; filename="+vars["user_uid"]+".ics")
	res.Write([]byte(cal.Build()))
	return
}
