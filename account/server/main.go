// Package main implements a gRPC server that handles Staffjoy accounts.
package main

import (
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"

	pb "v2.staffjoy.com/account"
	"v2.staffjoy.com/environments"
	"v2.staffjoy.com/sms"

	"v2.staffjoy.com/healthcheck"
)

const (
	// ServiceName identifies this app in logs
	ServiceName = "accountserver"
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
		panic("Unable to determine accountserver configuration")
	}
	logger = config.GetLogger(ServiceName)
}

func main() {
	logger.Debugf("Booting accountserver environment %s", config.Name)

	s := &accountServer{logger: logger, signingToken: os.Getenv("SIGNING_SECRET"), config: config}
	if !config.Debug {
		s.errorClient = environments.ErrorClient(&config)
	}

	smsConn, err := grpc.Dial(sms.Endpoint, grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer smsConn.Close()
	s.smsClient = sms.NewSmsServiceClient(smsConn)

	s.db, err = sql.Open("mysql", os.Getenv("MYSQL_CONFIG")+"?parseTime=true")
	if err != nil {
		logger.Panicf("Cannot connect to account db - %v", err)
	}
	defer s.db.Close()

	s.dbMap = &gorp.DbMap{Db: s.db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}
	_ = s.dbMap.AddTableWithName(pb.Account{}, "account").SetKeys(false, "uuid")
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
	pb.RegisterAccountServiceServer(grpcServer, s)

	// set up a health check listener for kubernetes
	go func() {
		logger.Debugf("Booting accountserver health check %s", config.Name)
		http.HandleFunc(healthcheck.HEALTHPATH, healthcheck.Handler)
		http.ListenAndServe(":80", nil)
	}()

	grpcServer.Serve(lis)
}
