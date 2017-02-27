package main

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"v2.staffjoy.com/account"
	"v2.staffjoy.com/bot"
	"v2.staffjoy.com/company"
	"v2.staffjoy.com/sms"

	"time"
)

func (s *botServer) AlertNewShift(ctx context.Context, req *bot.AlertNewShiftRequest) (*empty.Empty, error) {
	companyUUID := req.NewShift.CompanyUuid
	teamUUID := req.NewShift.TeamUuid

	botCtx := botContext()

	accountClient, close, err := account.NewClient()
	if err != nil {
		return nil, s.internalError(err, "unable to initiate account connection")
	}
	defer close()

	a, err := accountClient.Get(botCtx, &account.GetAccountRequest{Uuid: req.UserUuid})
	if err != nil {
		return nil, s.internalError(err, "cannot find user")
	}

	companyClient, close, err := company.NewClient()
	if err != nil {
		return nil, s.internalError(err, "unable to initiate company connection")
	}
	defer close()

	c, err := companyClient.GetCompany(botCtx, &company.GetCompanyRequest{Uuid: companyUUID})
	if err != nil {
		return nil, s.internalError(err, "cannot find company")
	}
	t, err := companyClient.GetTeam(botCtx, &company.GetTeamRequest{CompanyUuid: companyUUID, Uuid: teamUUID})
	if err != nil {
		return nil, s.internalError(err, "cannot find company")
	}

	newShift, err := printShiftSms(req.NewShift, t.Timezone)
	if err != nil {
		return nil, s.internalError(err, "cannot format start of shift")
	}

	jobName, err := JobName(companyUUID, teamUUID, req.NewShift.JobUuid)
	if err != nil {
		logger.Errorf("unable to fetch job name - %s", err)
	}
	// Format name with leading space.
	if jobName != "" {
		jobName = fmt.Sprintf(" %s", jobName)
	}

	u := user(*a)
	if u.PreferredDispatch() != dispatchSms {
		return &empty.Empty{}, nil
	}

	msg := fmt.Sprintf("%s Your %s manager just published a new%s shift for you: \n%s", u.Greet(), c.Name, jobName, newShift)

	smsClient, close, err := sms.NewClient()
	if err != nil {
		return nil, s.internalError(err, "unable to initiate sms connection")
	}
	defer close()

	if _, err := smsClient.QueueSend(ctx, &sms.SmsRequest{To: a.Phonenumber, Body: msg}); err != nil {
		s.internalError(err, "could not send welcome sms")
	}
	return &empty.Empty{}, nil
}

func (s *botServer) AlertNewShifts(ctx context.Context, req *bot.AlertNewShiftsRequest) (*empty.Empty, error) {
	shifts := req.NewShifts
	if len(shifts) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "empty shifts array in request")
	}

	companyUUID := shifts[0].CompanyUuid
	teamUUID := shifts[0].TeamUuid

	botCtx := botContext()

	accountClient, close, err := account.NewClient()
	if err != nil {
		return nil, s.internalError(err, "unable to initiate account connection")
	}
	defer close()

	a, err := accountClient.Get(botCtx, &account.GetAccountRequest{Uuid: req.UserUuid})
	if err != nil {
		return nil, s.internalError(err, "cannot find user")
	}

	companyClient, close, err := company.NewClient()
	if err != nil {
		return nil, s.internalError(err, "unable to initiate company connection")
	}
	defer close()

	c, err := companyClient.GetCompany(botCtx, &company.GetCompanyRequest{Uuid: companyUUID})
	if err != nil {
		return nil, s.internalError(err, "cannot find company")
	}
	t, err := companyClient.GetTeam(botCtx, &company.GetTeamRequest{CompanyUuid: companyUUID, Uuid: teamUUID})
	if err != nil {
		return nil, s.internalError(err, "cannot find team")
	}

	var newShifts string
	for _, shift := range shifts {
		newShift, err := printShiftSms(shift, t.Timezone)
		if err != nil {
			return nil, s.internalError(err, "cannot format start of shift")
		}

		jobName, err := JobName(companyUUID, teamUUID, shift.JobUuid)
		if err != nil {
			logger.Errorf("unable to fetch job name - %s", err)
		}
		// Format name with leading space.
		if jobName != "" {
			jobName = fmt.Sprintf(" (%s)", jobName)
		}

		newShifts += fmt.Sprintf("%s%s\n", newShift, jobName)
	}

	u := user(*a)
	if u.PreferredDispatch() != dispatchSms {
		return &empty.Empty{}, nil
	}

	msg := fmt.Sprintf("%s Your %s manager just published %d new shifts that you are working: \n%s", u.Greet(), c.Name, len(shifts), newShifts)

	smsClient, close, err := sms.NewClient()
	if err != nil {
		return nil, s.internalError(err, "unable to initiate sms connection")
	}
	defer close()

	if _, err = smsClient.QueueSend(ctx, &sms.SmsRequest{To: a.Phonenumber, Body: msg}); err != nil {
		s.internalError(err, "could not send bulk new shifts sms")
	}
	return &empty.Empty{}, nil
}

func (s *botServer) AlertRemovedShift(ctx context.Context, req *bot.AlertRemovedShiftRequest) (*empty.Empty, error) {
	companyUUID := req.OldShift.CompanyUuid
	teamUUID := req.OldShift.TeamUuid

	botCtx := botContext()

	accountClient, close, err := account.NewClient()
	if err != nil {
		return nil, s.internalError(err, "unable to initiate account connection")
	}
	defer close()

	a, err := accountClient.Get(botCtx, &account.GetAccountRequest{Uuid: req.UserUuid})
	if err != nil {
		return nil, s.internalError(err, "cannot find user")
	}
	companyClient, close, err := company.NewClient()
	if err != nil {
		return nil, s.internalError(err, "unable to initiate company connection")
	}
	defer close()

	c, err := companyClient.GetCompany(botCtx, &company.GetCompanyRequest{Uuid: companyUUID})
	if err != nil {
		return nil, s.internalError(err, "cannot find company")
	}
	t, err := companyClient.GetTeam(botCtx, &company.GetTeamRequest{CompanyUuid: companyUUID, Uuid: teamUUID})
	if err != nil {
		return nil, s.internalError(err, "cannot find team")
	}

	u := user(*a)
	if u.PreferredDispatch() != dispatchSms {
		return &empty.Empty{}, nil
	}

	nshifts, err := companyClient.ListWorkerShifts(botCtx, &company.WorkerShiftListRequest{CompanyUuid: companyUUID, TeamUuid: teamUUID, WorkerUuid: req.UserUuid, ShiftStartAfter: time.Now(), ShiftStartBefore: time.Now().AddDate(0, 0, ShiftWindow)})
	if err != nil {
		return nil, s.internalError(err, "cannot find team")
	}

	var newShifts string
	for _, shift := range nshifts.Shifts {
		newShift, err := printShiftSms(&shift, t.Timezone)
		if err != nil {
			return nil, s.internalError(err, "cannot format start of shift")
		}
		newShifts += fmt.Sprintf("%s\n", newShift)
	}

	msg := fmt.Sprintf("%s Your %s manager just removed you from a shift, so you are no longer working it. Here is your new schedule: \n%s", u.Greet(), c.Name, newShifts)

	logger.Infof("msg sent to user was: %v", msg)

	smsClient, close, err := sms.NewClient()
	if err != nil {
		return nil, s.internalError(err, "unable to initiate sms connection")
	}
	defer close()

	if _, err := smsClient.QueueSend(ctx, &sms.SmsRequest{To: a.Phonenumber, Body: msg}); err != nil {
		s.internalError(err, "could not send sms")
	}
	return &empty.Empty{}, nil
}

func (s *botServer) AlertRemovedShifts(ctx context.Context, req *bot.AlertRemovedShiftsRequest) (*empty.Empty, error) {
	shifts := req.OldShifts
	if len(shifts) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "empty shifts array in request")
	}

	companyUUID := shifts[0].CompanyUuid
	teamUUID := shifts[0].TeamUuid

	botCtx := botContext()

	accountClient, close, err := account.NewClient()
	if err != nil {
		return nil, s.internalError(err, "unable to initiate account connection")
	}
	defer close()

	a, err := accountClient.Get(botCtx, &account.GetAccountRequest{Uuid: req.UserUuid})
	if err != nil {
		return nil, s.internalError(err, "cannot find user")
	}

	companyClient, close, err := company.NewClient()
	if err != nil {
		return nil, s.internalError(err, "unable to initiate company connection")
	}
	defer close()

	c, err := companyClient.GetCompany(botCtx, &company.GetCompanyRequest{Uuid: companyUUID})
	if err != nil {
		return nil, s.internalError(err, "cannot find company")
	}

	t, err := companyClient.GetTeam(botCtx, &company.GetTeamRequest{CompanyUuid: companyUUID, Uuid: teamUUID})
	if err != nil {
		return nil, s.internalError(err, "cannot find team")
	}

	var oldShifts string
	for _, shift := range shifts {
		oldShift, err := printShiftSms(shift, t.Timezone)
		if err != nil {
			return nil, s.internalError(err, "cannot format start of shift")
		}
		oldShifts += fmt.Sprintf("%s\n", oldShift)
	}

	u := user(*a)
	if u.PreferredDispatch() != dispatchSms {
		return &empty.Empty{}, nil
	}

	nshifts, err := companyClient.ListWorkerShifts(botCtx, &company.WorkerShiftListRequest{CompanyUuid: companyUUID, TeamUuid: teamUUID, WorkerUuid: req.UserUuid, ShiftStartAfter: time.Now(), ShiftStartBefore: time.Now().AddDate(0, 0, ShiftWindow)})
	if err != nil {
		return nil, s.internalError(err, "cannot find worker shifts")
	}

	var newShifts string
	for _, shift := range nshifts.Shifts {
		newShift, err := printShiftSms(&shift, t.Timezone)
		if err != nil {
			return nil, s.internalError(err, "cannot format start of shift")
		}
		newShifts += fmt.Sprintf("%s\n", newShift)
	}

	msg := fmt.Sprintf("%s Your %s manager just removed %d of your shifts so you are no longer working it. \n Your new shifts are: \n%s", u.Greet(), c.Name, len(shifts), newShifts)
	smsClient, close, err := sms.NewClient()
	if err != nil {
		return nil, s.internalError(err, "unable to initiate sms connection")
	}
	defer close()

	if _, err := smsClient.QueueSend(ctx, &sms.SmsRequest{To: a.Phonenumber, Body: msg}); err != nil {
		s.internalError(err, "could not send bulk new shifts sms")
	}
	return &empty.Empty{}, nil
}

func (s *botServer) AlertChangedShift(ctx context.Context, req *bot.AlertChangedShiftRequest) (*empty.Empty, error) {
	companyUUID := req.OldShift.CompanyUuid
	teamUUID := req.OldShift.TeamUuid

	botCtx := botContext()

	accountClient, close, err := account.NewClient()
	if err != nil {
		return nil, s.internalError(err, "unable to initiate account connection")
	}
	defer close()

	a, err := accountClient.Get(botCtx, &account.GetAccountRequest{Uuid: req.UserUuid})
	if err != nil {
		return nil, s.internalError(err, "cannot find user")
	}
	companyClient, close, err := company.NewClient()
	if err != nil {
		return nil, s.internalError(err, "unable to initiate company connection")
	}
	defer close()

	c, err := companyClient.GetCompany(botCtx, &company.GetCompanyRequest{Uuid: companyUUID})
	if err != nil {
		return nil, s.internalError(err, "cannot find company")
	}
	t, err := companyClient.GetTeam(botCtx, &company.GetTeamRequest{CompanyUuid: companyUUID, Uuid: teamUUID})
	if err != nil {
		return nil, s.internalError(err, "cannot find team")
	}

	oldShift, err := printShiftSms(req.OldShift, t.Timezone)
	if err != nil {
		return nil, s.internalError(err, "cannot format old shift")
	}
	oldJobName, err := JobName(companyUUID, teamUUID, req.NewShift.JobUuid)
	if err != nil {
		logger.Errorf("unable to fetch job name - %s", err)
	}
	// Format name with leading space.
	if oldJobName != "" {
		oldShift += fmt.Sprintf(" (%s)", oldJobName)
	}

	newShift, err := printShiftSms(req.NewShift, t.Timezone)
	if err != nil {
		return nil, s.internalError(err, "cannot format new shift")
	}
	newJobName, err := JobName(companyUUID, teamUUID, req.OldShift.JobUuid)
	if err != nil {
		logger.Errorf("unable to fetch job name - %s", err)
	}
	// Format name with leading space.
	if newJobName != "" {
		newShift += fmt.Sprintf(" (%s)", newJobName)
	}

	u := user(*a)
	if u.PreferredDispatch() != dispatchSms {
		return &empty.Empty{}, nil
	}

	msg := fmt.Sprintf("%s Your %s manager just changed your shift: \nOld: %s\nNew:%s", u.Greet(), c.Name, oldShift, newShift)

	smsClient, close, err := sms.NewClient()
	if err != nil {
		return nil, s.internalError(err, "unable to initiate sms connection")
	}
	defer close()

	if _, err := smsClient.QueueSend(ctx, &sms.SmsRequest{To: a.Phonenumber, Body: msg}); err != nil {
		s.internalError(err, "could not send sms")
	}
	return &empty.Empty{}, nil
}
