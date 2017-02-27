package company

import (
	"fmt"

	"google.golang.org/grpc"
)

// NewClient returns a gRPC client for interacting with the company.
// After calling it, run a defer close on the close function
func NewClient() (CompanyServiceClient, func() error, error) {
	conn, err := grpc.Dial(Endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, nil, fmt.Errorf("did not connect: %v", err)
	}
	return NewCompanyServiceClient(conn), conn.Close, nil
}
