package main

import (
	"sync"

	"golang.org/x/net/context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"v2.staffjoy.com/auth"
	"v2.staffjoy.com/environments"
	pb "v2.staffjoy.com/sms"
)

type smsServer struct {
	logger        *logrus.Entry
	errorClient   environments.SentryClient
	config        *environments.Config
	sendingConfig *sendingConfig
	queue         chan *pb.SmsRequest
	wg            sync.WaitGroup
	twilioSid     string
	twilioToken   string
}

func (s *smsServer) QueueSend(ctx context.Context, req *pb.SmsRequest) (*empty.Empty, error) {
	_, authz, err := getAuth(ctx)
	if err != nil {
		return nil, s.internalError(err, "failed to authorize")
	}

	switch authz {
	case auth.AuthorizationCompanyService:
	case auth.AuthorizationBotService:
	case auth.AuthorizationAccountService:
	default:
		return nil, grpc.Errorf(codes.PermissionDenied, "you do not have access to this service")
	}

	if s.sendingConfig.WhitelistOnly {
		allowedToSend := false
		for _, w := range whitelist {
			if w == req.To {
				allowedToSend = true
			}
		}
		if !allowedToSend {
			s.logger.Warningf("prevented sending to number %s due to whitelist", req.To)
			return &empty.Empty{}, nil
		}
	}

	s.queue <- req
	s.logger.Debugf("enqueued new message to %v", req.To)
	return &empty.Empty{}, nil
}
