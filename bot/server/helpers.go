package main

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"v2.staffjoy.com/auth"
	"v2.staffjoy.com/company"
)

const (
	smsStartTimeFormat = "Mon 1/2 3:04PM"
	smsStopTimeFormat  = "3:04PM"
	smsShiftFormat     = "%s - %s" // sprint in start and stop
)

func (u *user) FirstName() string {
	names := strings.Split(u.Name, " ")
	if len(names) == 0 {
		return "there"
	}
	return names[0]
}

func botContext() context.Context {
	md := metadata.New(map[string]string{auth.AuthorizationMetadata: auth.AuthorizationBotService})
	return metadata.NewOutgoingContext(context.Background(), md)
}

func (s *botServer) internalError(err error, format string, a ...interface{}) error {
	s.logger.Errorf("%s: %v", format, err)
	if s.errorClient != nil {
		s.errorClient.CaptureError(err, nil)
	}
	return grpc.Errorf(codes.Unknown, format, a...)
}

func printShiftSms(shift *company.Shift, tz string) (string, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return "", err
	}
	startTime := shift.Start.In(loc).Format(smsStartTimeFormat)
	endTime := shift.Stop.In(loc).Format(smsStopTimeFormat)
	return fmt.Sprintf(smsShiftFormat, startTime, endTime), nil
}

// JobName returns the name of a job, given its UUID
func JobName(companyUUID, teamUUID, jobUUID string) (string, error) {
	if jobUUID == "" {
		return "", nil
	}

	companyClient, close, err := company.NewClient()
	if err != nil {
		return "", err
	}
	defer close()

	j, err := companyClient.GetJob(botContext(), &company.GetJobRequest{CompanyUuid: companyUUID, TeamUuid: teamUUID, Uuid: jobUUID})
	if err != nil {
		return "", err
	}
	return j.Name, nil

}
