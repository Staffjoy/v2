// Package main implements a gRPC server that handles Staffjoy emails.
package main

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"

	"v2.staffjoy.com/environments"
	"v2.staffjoy.com/healthcheck"
	pb "v2.staffjoy.com/sms"
)

const (
	// ServiceName identifies this app in logs
	ServiceName = "sms"
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
		panic("Unable to determine smsserver configuration")
	}
	logger = config.GetLogger(ServiceName)
}

func main() {
	logger.Debugf("Booting smsserver environment %s", config.Name)
	s := &smsServer{logger: logger, config: &config}
	if !config.Debug {
		s.errorClient = environments.ErrorClient(&config)
	}
	s.twilioSid = os.Getenv("TWILIO_SID")
	s.twilioToken = os.Getenv("TWILIO_TOKEN")
	if sc, ok := sendingConfigs[config.Name]; !ok {
		panic("could not determine sending config")
	} else {
		s.sendingConfig = &sc
	}
	s.queue = make(chan *pb.SmsRequest, 60) // buffered to 60 messages

	lis, err := net.Listen("tcp", pb.ServerPort)
	if err != nil {
		logger.Panicf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterSmsServiceServer(grpcServer, s)

	// set up a health check listener for kubernetes
	go func() {
		logger.Debugf("Booting smsserver health check %s", config.Name)
		http.HandleFunc(healthcheck.HEALTHPATH, healthcheck.Handler)
		http.ListenAndServe(":80", nil)
	}()

	go func() {
		grpcServer.Serve(lis)
	}()

	for i := 0; i < s.sendingConfig.Concurrency; i++ {
		s.Sender()
	}

	// graceful shutdown stuff

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// This goroutine executes a blocking receive for
	// signals. When it gets one it'll print it out
	// and then notify the program that it can finish.
	go func() {
		sig := <-sigs
		s.logger.Infof("caught shutdown signal %s - flushing queue", sig)
		close(s.queue)
		s.wg.Wait()
		s.logger.Infof("queue emptied for shutdown")
		done <- true
	}()
	<-done
	s.logger.Infof("graceful shutdown complete")

}
