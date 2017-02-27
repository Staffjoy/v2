// Package main implements a gRPC server that handles Staffjoy emails.
package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/Sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/keighl/mandrill"
	pb "v2.staffjoy.com/email"
	"v2.staffjoy.com/environments"
	"v2.staffjoy.com/healthcheck"
)

const (
	// ServiceName identifies this app in logs
	ServiceName         = "emailapi"
	fromName            = "Staffjoy"
	from                = "help@staffjoy.com"
	staffjoyEmailSuffix = "@staffjoy.com"
	mandrillTemplate    = "staffjoy-base"
)

var (
	logger *logrus.Entry
	config environments.Config
)

type emailServer struct {
	logger      *logrus.Entry
	errorClient environments.SentryClient
	client      *mandrill.Client
	clientMutex *sync.Mutex
	config      *environments.Config
}

// Setup environment, logger, etc
func init() {
	// Set the ENV environment variable to control dev/stage/prod behavior
	var err error
	config, err = environments.GetConfig(os.Getenv(environments.EnvVar))
	if err != nil {
		panic("Unable to determine emailserver configuration")
	}
	logger = config.GetLogger(ServiceName)
}

func main() {
	logger.Debugf("Booting emailserver environment %s", config.Name)
	s := &emailServer{logger: logger, config: &config, clientMutex: &sync.Mutex{}}
	if !config.Debug {
		s.errorClient = environments.ErrorClient(&config)
	}
	s.client = mandrill.ClientWithKey(os.Getenv("MANDRILL_API_KEY"))

	var err error

	lis, err := net.Listen("tcp", pb.ServerPort)
	if err != nil {
		logger.Panicf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterEmailServiceServer(grpcServer, s)

	// set up a health check listener for kubernetes
	go func() {
		logger.Debugf("Booting emailserver health check %s", config.Name)
		http.HandleFunc(healthcheck.HEALTHPATH, healthcheck.Handler)
		http.ListenAndServe(":80", nil)
	}()

	grpcServer.Serve(lis)
}

func (s *emailServer) Send(ctx context.Context, req *pb.EmailRequest) (*empty.Empty, error) {
	if len(req.To) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Please provide an email")
	}
	if len(req.Subject) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Please provide a subject")
	}
	if len(req.HtmlBody) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Please provide a valid body")
	}
	// Send asynchronously
	go s.processSend(req)
	return &empty.Empty{}, nil
}

func (s *emailServer) processSend(req *pb.EmailRequest) {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()

	logLine := s.logger.WithFields(logrus.Fields{
		"subject":   req.Subject,
		"to":        req.To,
		"html_body": req.HtmlBody,
	})

	// In development and staging - only send emails to @staffjoy.com
	if s.config.Name != "production" {
		// prepend env for sanity
		req.Subject = fmt.Sprintf("[%s] %s", s.config.Name, req.Subject)

		if !strings.HasSuffix(req.To, staffjoyEmailSuffix) {
			logLine.Warningf("Intercepted sending due to non-production environment.")
			return
		}
	}
	message := &mandrill.Message{}
	message.AddRecipient(req.To, req.Name, "to")
	message.FromEmail = from
	message.FromName = fromName
	message.Subject = req.Subject

	templateContent := map[string]string{"body": req.HtmlBody, "title": req.Subject}

	res, err := s.client.MessagesSendTemplate(message, mandrillTemplate, templateContent)
	if err != nil {
		if s.errorClient != nil {
			s.errorClient.CaptureError(err, map[string]string{
				"subject": req.Subject,
				"to":      req.To,
			})
		}
		logLine.Errorf("Unable to send email - %s %v", err, res)
		return
	}
	logLine.Infof("successfully sent - %v", res)
	return
}
