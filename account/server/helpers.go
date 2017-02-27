package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"time"

	"golang.org/x/net/context"

	"github.com/ttacon/libphonenumber"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"v2.staffjoy.com/auditlog"
	"v2.staffjoy.com/auth"
	"v2.staffjoy.com/sms"
)

const (
	defaultRegion = "US" // USA
)

func (s *accountServer) internalError(err error, format string, a ...interface{}) error {
	s.logger.Errorf("%s: %v", format, err)
	if s.errorClient != nil {
		s.errorClient.CaptureError(err, nil)
	}
	return grpc.Errorf(codes.Unknown, format, a...)
}

// ParseAndFormatPhonenumber takes a raw input phone number,
// and exports it in E164 format. All phone numbers going to the
// database should use this.
func ParseAndFormatPhonenumber(input string) (cleanPhonenumber string, err error) {
	// If empty string input - return empy string
	if input == "" {
		return
	}
	p, err := libphonenumber.Parse(input, defaultRegion)
	if err != nil {
		return "", fmt.Errorf("Invalid phone number")
	}
	cleanPhonenumber = libphonenumber.Format(p, libphonenumber.E164)
	return
}

// GenerateGravatarURL returns the gravatar photo corresponding to the email
// (passing in a null string is fine)
func GenerateGravatarURL(email string) string {
	h := md5.New()
	io.WriteString(h, email)
	return fmt.Sprintf("https://www.gravatar.com/avatar/%x.jpg?s=400&d=identicon", h.Sum(nil))
}

func newAuditEntry(md metadata.MD, targetType, targetUUID string) *auditlog.Entry {
	u, err := auth.GetCurrentUserUUIDFromMetadata(md)
	if err != nil {
		u = ""
	}

	var a string
	aa, ok := md[auth.AuthorizationMetadata]
	if ok {
		a = aa[0]
	}

	return &auditlog.Entry{
		CurrentUserUUID: u,
		Authorization:   a,
		TargetType:      targetType,
		TargetUUID:      targetUUID,
	}
}

func getAuth(ctx context.Context) (md metadata.MD, authz string, err error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil, "", fmt.Errorf("Context missing metadata")
	}
	if len(md[auth.AuthorizationMetadata]) == 0 {
		return nil, "", fmt.Errorf("Missing Authorization")
	}
	authz = md[auth.AuthorizationMetadata][0]
	return
}

func (s *accountServer) sendSmsGreeting(phonenumber string) {
	go func() {
		md := metadata.New(map[string]string{auth.AuthorizationMetadata: auth.AuthorizationAccountService})
		ctx, cancel := context.WithTimeout(metadata.NewContext(context.Background(), md), 2*time.Second)
		defer cancel()
		_, err := s.smsClient.QueueSend(ctx, &sms.SmsRequest{To: phonenumber, Body: "Welcome to Staffjoy!"})
		if err != nil {
			s.internalError(err, "could not send welcome sms")
		}
	}()
}
