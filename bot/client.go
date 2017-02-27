package bot

import (
	"fmt"

	"google.golang.org/grpc"
)

// NewClient returns a gRPC client for interacting with the bot.
// After calling it, run a defer close on the close function
func NewClient() (BotServiceClient, func() error, error) {
	conn, err := grpc.Dial(Endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, nil, fmt.Errorf("did not connect: %v", err)
	}
	return NewBotServiceClient(conn), conn.Close, nil
}
