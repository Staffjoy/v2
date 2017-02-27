// Package auth controls how users and services authenticate
package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"v2.staffjoy.com/environments"

	"google.golang.org/grpc/metadata"
)

const (
	cookieName = "staffjoy-faraday"
	cookie
	uuidKey       = "uuid"
	supportKey    = "support"
	expirationKey = "exp"
	// for GRPC
	currentUserMetadata = "faraday-current-user-uuid"
	// header set for internal user id
	currentUserHeader = "Grpc-Metadata-Faraday-Current-User-Uuid"

	// AuthorizationHeader is the http request header
	// key used for accessing the internal authorization.
	AuthorizationHeader = "Authorization"

	// AuthorizationMetadata is the grpce metadadata key used
	// for accessing the internal authorization
	AuthorizationMetadata = "authorization"

	// AuthorizationAnonymousWeb is set as the Authorization header to denote that
	// a request is being made bu an unauthenticated web user
	AuthorizationAnonymousWeb = "faraday-anonymous"

	// AuthorizationAuthenticatedUser is set as the Authorization header to denote that
	// a request is being made by an authenticated web user
	AuthorizationAuthenticatedUser = "faraday-authenticated"

	// AuthorizationSupportUser is set as the Authorization header to denote that
	// a request is being made by a Staffjoy team me
	AuthorizationSupportUser = "faraday-support"

	// AuthorizationWWWService is set as the Authorization header to denote that
	// a request is being made by the www login / signup system
	AuthorizationWWWService = "www-service"

	// AuthorizationCompanyService is set as the Authorization header to denote
	// that a request is being made by the company api/server
	AuthorizationCompanyService = "company-service"

	// AuthorizationSuperpowersService is set as the Authorization header to
	// denote that a request is being made by the dev-only superpowers service
	AuthorizationSuperpowersService = "superpowers-service"

	// AuthorizationWhoamiService is set as the Authorization heade to denote that
	// a request is being made by the whoami microservice
	AuthorizationWhoamiService = "whoami-service"

	// AuthorizationBotService is set as the Authorization header to denote that
	// a request is being made by the bot microservice
	AuthorizationBotService = "bot-service"

	// AuthorizationAccountService is set as the Authorization header to denote that
	// a request is being made by the account service
	AuthorizationAccountService = "account-service"

	// AuthorizationICalService is set as the Authorization header to denote that
	// a request is being made by the ical service
	AuthorizationICalService = "ical-service"
)

var (
	signingSecret string
	shortSession  = time.Duration(12 * time.Hour)
	longSession   = time.Duration(30 * 24 * time.Hour)
	config        environments.Config
)

func init() {
	signingSecret = os.Getenv("SIGNING_SECRET")

	var err error
	config, err = environments.GetConfig(os.Getenv(environments.EnvVar))
	if err != nil {
		panic("Unable to determine configuration")
	}

}

// SetInternalHeaders is used by Faraday to sanitize incoming external requests
// and convert them to internal requests with authorization information
func SetInternalHeaders(externalReq *http.Request, internalHeaders http.Header) {
	ProxyHeaders(externalReq.Header, internalHeaders)

	// default to anonymous web then prove otherwise
	authorization := AuthorizationAnonymousWeb
	uuid, support, err := getSession(externalReq)
	if err == nil {
		if support {
			authorization = AuthorizationSupportUser
		} else {
			authorization = AuthorizationAuthenticatedUser
		}
		internalHeaders.Set(currentUserHeader, uuid)
	}
	internalHeaders.Set(AuthorizationHeader, authorization)
	return
}

// ProxyHeaders copies http headers
func ProxyHeaders(from, to http.Header) {
	// Range over the headres
	for k, v := range from {
		// TODO - filter restricted headers

		// Multiple header values may exist per key
		for _, x := range v {
			to.Add(k, x)
		}

	}
}

// GetCurrentUserUUIDFromMetadata allows backend gRPC services with
// authorization methods of AuthenticatedUser or SupportUser to access
// the uuid of the user making the request
func GetCurrentUserUUIDFromMetadata(data metadata.MD) (uuid string, err error) {
	res, ok := data[currentUserMetadata]
	if !ok || len(res) == 0 {
		err = fmt.Errorf("User not authenticated")
		return
	}
	uuid = res[0]
	return
}

// GetCurrentUserUUIDFromHeader allows backend http services with
// authorization methods of AuthenticatedUser or SupportUser to access
// the uuid of the user making the request
func GetCurrentUserUUIDFromHeader(data http.Header) (uuid string, err error) {
	res, ok := data[currentUserHeader]
	if !ok || len(res) == 0 {
		err = fmt.Errorf("User not authenticated")
		return
	}
	uuid = res[0]
	return
}
