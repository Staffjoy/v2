package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"v2.staffjoy.com/account"
	"v2.staffjoy.com/auditlog"
	"v2.staffjoy.com/auth"
	pb "v2.staffjoy.com/company"
)

var (
	daysOfWeek = []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}
)

func newAuditEntry(md metadata.MD, targetType, targetUUID, companyUUID, teamUUID string) *auditlog.Entry {
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
		CompanyUUID:     companyUUID,
		TeamUUID:        teamUUID,
	}
}

func (s *companyServer) internalError(err error, format string, a ...interface{}) error {
	s.logger.Errorf("%s: %v", format, err)
	if s.errorClient != nil {
		s.errorClient.CaptureError(err, nil)
	}
	return grpc.Errorf(codes.Unknown, format, a...)
}

func sanitizeDayOfWeek(input string) (string, error) {
	input = strings.TrimSpace(strings.ToLower(input))
	for _, day := range daysOfWeek {
		if day == input {
			return input, nil
		}
	}
	return "", fmt.Errorf("unknown day of week")
}

func validTimezone(input string) (err error) {
	_, err = time.LoadLocation(input)
	return
}

func validColor(input string) error {
	re := regexp.MustCompile(`^[0-9A-Fa-f]{3}$|^[0-9A-Fa-f]{6}$`)
	if !re.MatchString(input) {
		return fmt.Errorf("invalid color hex code")
	}

	return nil
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

func copyAccountToDirectory(a *account.Account, d *pb.DirectoryEntry) {
	d.UserUuid = a.Uuid
	d.Name = a.Name
	d.ConfirmedAndActive = a.ConfirmedAndActive
	d.Phonenumber = a.Phonenumber
	d.PhotoUrl = a.PhotoUrl
	d.Email = a.Email
	return
}

func asyncContext() context.Context {
	md := metadata.New(map[string]string{auth.AuthorizationMetadata: auth.AuthorizationCompanyService})
	return metadata.NewContext(context.Background(), md)
}
