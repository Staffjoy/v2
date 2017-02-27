// Package main implements a gRPC server that handles Staffjoy accounts.
package main

import (
	"database/sql"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"

	"github.com/Sirupsen/logrus"

	"google.golang.org/grpc"

	pb "v2.staffjoy.com/company"
	"v2.staffjoy.com/environments"

	"v2.staffjoy.com/healthcheck"
)

const (
	// ServiceName identifies this app in logs
	ServiceName = "companyserver"
)

var (
	logger         *logrus.Entry
	config         environments.Config
	serverLocation *time.Location
)

type companyServer struct {
	logger       *logrus.Entry
	db           *sql.DB
	errorClient  environments.SentryClient
	signingToken string
	dbMap        *gorp.DbMap
}

// Setup environment, logger, etc
func init() {
	// Set the ENV environment variable to control dev/stage/prod behavior
	var err error
	config, err = environments.GetConfig(os.Getenv(environments.EnvVar))
	if err != nil {
		panic("Unable to determine accountserver configuration")
	}
	logger = config.GetLogger(ServiceName)
	serverLocation, err = time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Debugf("PANIC! %s", err)
		}
	}()

	var err error

	logger.Debugf("Booting companyserver environment %s", config.Name)

	s := &companyServer{logger: logger, signingToken: os.Getenv("SIGNING_SECRET")}
	if !config.Debug {
		s.errorClient = environments.ErrorClient(&config)
	}

	s.db, err = sql.Open("mysql", os.Getenv("MYSQL_CONFIG")+"?parseTime=true")
	if err != nil {
		logger.Panicf("Cannot connect to company db - %v", err)
	}
	defer s.db.Close()

	s.dbMap = &gorp.DbMap{Db: s.db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}
	_ = s.dbMap.AddTableWithName(pb.Company{}, "company").SetKeys(false, "uuid")
	_ = s.dbMap.AddTableWithName(pb.Team{}, "team").SetKeys(false, "uuid")
	_ = s.dbMap.AddTableWithName(pb.Shift{}, "shift").SetKeys(false, "uuid")
	_ = s.dbMap.AddTableWithName(pb.Job{}, "job").SetKeys(false, "uuid")
	_ = s.dbMap.AddTableWithName(pb.DirectoryEntry{}, "directory")

	if config.Debug {
		s.dbMap.TraceOn("[gorp]", logger)
	}

	// listen for incoming conections
	lis, err := net.Listen("tcp", pb.ServerPort)
	if err != nil {
		logger.Panicf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterCompanyServiceServer(grpcServer, s)

	// set up a health check listener for kubernetes
	go func() {
		logger.Debugf("Booting companyserver health check %s", config.Name)
		http.HandleFunc(healthcheck.HEALTHPATH, healthcheck.Handler)
		http.ListenAndServe(":80", nil)
	}()

	grpcServer.Serve(lis)
}
