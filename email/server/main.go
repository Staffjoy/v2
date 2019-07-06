// Package main implements a gRPC server that handles Staffjoy emails.
package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/mailgun/mailgun-go/v3"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

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
)

var (
	logger *logrus.Entry
	config environments.Config
)

type emailServer struct {
	logger      *logrus.Entry
	errorClient environments.SentryClient
	client      *mailgun.MailgunImpl
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
	s.client = mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"))
	s.client.SetAPIBase("https://api.eu.mailgun.net/v3")

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

	//templateContent := map[string]string{"body": req.HtmlBody, "title": req.Subject}

	// The message object allows you to add attachments and Bcc recipients
	message := s.client.NewMessage(from, req.Subject, req.HtmlBody, req.To)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, id, err := s.client.Send(ctx, message)
	if err != nil {
		if s.errorClient != nil {
			s.errorClient.CaptureError(err, map[string]string{
				"subject": req.Subject,
				"to":      req.To,
			})
		}
		logLine.Errorf("Unable to send email - %s %s %v", err, id, res)
		return
	}
	logLine.Infof("successfully sent - %v", res)
	return
}
