package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"v2.staffjoy.com/auth"
)

func getAuth(ctx context.Context) (md metadata.MD, authz string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, "", fmt.Errorf("Context missing metadata")
	}
	if len(md[auth.AuthorizationMetadata]) == 0 {
		return nil, "", fmt.Errorf("Missing Authorization")
	}
	authz = md[auth.AuthorizationMetadata][0]
	return
}

func (s *smsServer) internalError(err error, format string, a ...interface{}) error {
	s.logger.Errorf("%s: %v", format, err)
	if s.errorClient != nil {
		s.errorClient.CaptureError(err, nil)
	}
	return grpc.Errorf(codes.Unknown, format, a...)
}
