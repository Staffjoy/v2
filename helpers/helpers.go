package helpers

import (
	"context"

	"google.golang.org/grpc/metadata"
	"v2.staffjoy.com/account"
	"v2.staffjoy.com/auth"
)

// TrackEventFromMetadata determines the current user from gRPC context metadata
// and tracks the event if it is an authenticated user
func TrackEventFromMetadata(md metadata.MD, eventName string) (err error) {
	if len(md[auth.AuthorizationMetadata]) == 0 {
		// noop - no authentication
		return
	}
	authz := md[auth.AuthorizationMetadata][0]
	if authz != auth.AuthorizationAuthenticatedUser {
		// Not an action performed by a normal user
		// (noop - not an error)
		return
	}

	userUUID, err := auth.GetCurrentUserUUIDFromMetadata(md)
	if err != nil {
		return
	}
	err = TrackEvent(userUUID, eventName)
	return
}

// TrackEvent is a helper function for tracking user events
func TrackEvent(userUUID, eventName string) (err error) {
	var s account.AccountServiceClient
	var close func() error

	s, close, err = account.NewClient()
	if err != nil {
		return
	}
	defer close()
	ctx := context.Background()
	_, err = s.TrackEvent(ctx, &account.TrackEventRequest{Uuid: userUUID, Event: eventName})
	return
}

// SyncUser is a helper function for re-sending user info to tracking services
func SyncUser(userUUID string) (err error) {
	var s account.AccountServiceClient
	var close func() error

	s, close, err = account.NewClient()
	if err != nil {
		return
	}
	defer close()
	ctx := context.Background()
	_, err = s.SyncUser(ctx, &account.SyncUserRequest{Uuid: userUUID})
	return
}
