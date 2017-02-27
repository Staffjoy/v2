package sms

import (
	"fmt"

	"google.golang.org/grpc"
)

// NewClient returns a gRPC client for interacting with the sms.
// After calling it, run a defer close on the close function
func NewClient() (SmsServiceClient, func() error, error) {
	conn, err := grpc.Dial(Endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, nil, fmt.Errorf("did not connect: %v", err)
	}
	return NewSmsServiceClient(conn), conn.Close, nil
}
