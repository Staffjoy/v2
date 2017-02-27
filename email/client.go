package email

import (
	"fmt"

	"google.golang.org/grpc"
)

// NewClient returns a gRPC client for interacting with the email.
// After calling it, run a defer close on the close function
func NewClient() (EmailServiceClient, func() error, error) {
	conn, err := grpc.Dial(Endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, nil, fmt.Errorf("did not connect: %v", err)
	}
	return NewEmailServiceClient(conn), conn.Close, nil
}
