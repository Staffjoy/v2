// Package main implements a server that handles messaging to workers.
package main

import (
	"net"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"

	"v2.staffjoy.com/account"
	"v2.staffjoy.com/bot"
	"v2.staffjoy.com/environments"
	"v2.staffjoy.com/healthcheck"
)

const (
	// ServiceName identifies this app in logs
	ServiceName = "smsserver"

	// ShiftWindow is the number of days out the bot will inform the worker
	// of their schedule
	ShiftWindow = 10
)

type user account.Account

var (
	logger *logrus.Entry
	config environments.Config
)

type botServer struct {
	logger      *logrus.Entry
	errorClient environments.SentryClient
	config      *environments.Config
}

// Setup environment, logger, etc
func init() {
	// Set the ENV environment variable to control dev/stage/prod behavior
	var err error
	config, err = environments.GetConfig(os.Getenv(environments.EnvVar))
	if err != nil {
		panic("Unable to determine botserver configuration")
	}
	logger = config.GetLogger(ServiceName)
}

func main() {
	logger.Debugf("Booting botserver environment %s", config.Name)
	s := &botServer{logger: logger, config: &config}
	if !config.Debug {
		s.errorClient = environments.ErrorClient(&config)
	}

	lis, err := net.Listen("tcp", bot.ServerPort)
	if err != nil {
		logger.Panicf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	bot.RegisterBotServiceServer(grpcServer, s)

	// set up a health check listener for kubernetes
	go func() {
		logger.Debugf("Booting botserver health check %s", config.Name)
		http.HandleFunc(healthcheck.HEALTHPATH, healthcheck.Handler)
		http.ListenAndServe(":80", nil)
	}()

	grpcServer.Serve(lis)
}
