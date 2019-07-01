package main

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/ptypes/empty"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"v2.staffjoy.com/auth"
	"v2.staffjoy.com/bot"
	pb "v2.staffjoy.com/company"
	"v2.staffjoy.com/crypto"
	"v2.staffjoy.com/helpers"
)

var (
	maxShiftDuration = time.Duration(23 * time.Hour)
)

func (s *companyServer) CreateShift(ctx context.Context, req *pb.CreateShiftRequest) (*pb.Shift, error) {
	md, authz, err := getAuth(ctx)
	if err != nil {
		return nil, s.internalError(err, "failed to authorize")
	}
	switch authz {
	case auth.AuthorizationSupportUser:
		if err = s.PermissionCompanyAdmin(md, req.CompanyUuid); err != nil {
			return nil, err
		}
	case auth.AuthorizationAuthenticatedUser:
	default:
		return nil, grpc.Errorf(codes.PermissionDenied, "You do not have access to this service")
	}

	if _, err = s.GetTeam(ctx, &pb.GetTeamRequest{Uuid: req.TeamUuid, CompanyUuid: req.CompanyUuid}); err != nil {
		return nil, err
	}
	if req.JobUuid != "" {
		if _, err = s.GetJob(ctx, &pb.GetJobRequest{Uuid: req.JobUuid, CompanyUuid: req.CompanyUuid, TeamUuid: req.TeamUuid}); err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, "Invalid job parameter")
		}
	}

	if req.UserUuid != "" {
		if _, err = s.GetDirectoryEntry(ctx, &pb.DirectoryEntryRequest{CompanyUuid: req.CompanyUuid, UserUuid: req.UserUuid}); err != nil {
			return nil, err
		}
	}

	uuid, err := crypto.NewUUID()
	if err != nil {
		return nil, s.internalError(err, "cannot generate a uuid")
	}

	dur := req.Stop.Sub(req.Start)
	if dur <= 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "stop must be after start")
	} else if dur > maxShiftDuration {
		return nil, grpc.Errorf(codes.InvalidArgument, "Shifts exceed max %f hour duration", maxShiftDuration.Hours())
	}

	shift := &pb.Shift{Uuid: uuid.String(), CompanyUuid: req.CompanyUuid, TeamUuid: req.TeamUuid, JobUuid: req.JobUuid, Start: req.Start, Stop: req.Stop, Published: req.Published, UserUuid: req.UserUuid}
	if err = s.dbMap.Insert(shift); err != nil {
		return nil, s.internalError(err, "could not create shift")
	}
	al := newAuditEntry(md, "shift", shift.Uuid, shift.CompanyUuid, req.TeamUuid)
	al.UpdatedContents = shift
	al.Log(logger, "created shift")

	go func() {
		if shift.UserUuid != "" && shift.Published {
			botClient, close, err := bot.NewClient()
			if err != nil {
				s.internalError(err, "unable to initiate bot connection")
				return
			}
			defer close()
			if _, err := botClient.AlertNewShift(asyncContext(), &bot.AlertNewShiftRequest{UserUuid: shift.UserUuid, NewShift: shift}); err != nil {
				s.internalError(err, "failed to alert worker about new shift")
			}
		}
	}()
	go helpers.TrackEventFromMetadata(md, "shift_created")
	if req.Published {
		go helpers.TrackEventFromMetadata(md, "shift_published")
	}

	return shift, nil
}

func (s *companyServer) ListWorkerShifts(ctx context.Context, req *pb.WorkerShiftListRequest) (*pb.ShiftList, error) {
	// Prep
	md, authz, err := getAuth(ctx)
	if err != nil {
		return nil, s.internalError(err, "Failed to authorize")
	}

	switch authz {
	case auth.AuthorizationAuthenticatedUser:
	case auth.AuthorizationSupportUser:
		if err = s.PermissionTeamWorker(md, req.CompanyUuid, req.TeamUuid); err != nil {
			return nil, err
		}
	case auth.AuthorizationBotService:
	case auth.AuthorizationICalService:
	default:
		return nil, grpc.Errorf(codes.PermissionDenied, "you do not have access to this service")
	}

	if _, err = s.GetTeam(ctx, &pb.GetTeamRequest{Uuid: req.TeamUuid, CompanyUuid: req.CompanyUuid}); err != nil {
		return nil, err
	}

	if req.ShiftStartAfter.After(req.ShiftStartBefore) {
		return nil, grpc.Errorf(codes.InvalidArgument, "shift_start_after must be before shift_start_before")
	}
	res := &pb.ShiftList{ShiftStartAfter: req.ShiftStartAfter, ShiftStartBefore: req.ShiftStartBefore}

	var dbShifts []pb.Shift
	if _, err = s.dbMap.Select(&dbShifts, "select * from shift where team_uuid=? and user_uuid=? AND start>=? AND start<? order by start asc", req.TeamUuid, req.WorkerUuid, req.ShiftStartAfter, req.ShiftStartBefore); err != nil {
		return nil, s.internalError(err, "unable to query database")
	}

	for _, shift := range dbShifts {
		shift.CompanyUuid = req.CompanyUuid
		res.Shifts = append(res.Shifts, shift)
	}
	return res, nil
}

func (s *companyServer) ListShifts(ctx context.Context, req *pb.ShiftListRequest) (*pb.ShiftList, error) {
	// Prep
	md, authz, err := getAuth(ctx)
	if err != nil {
		return nil, s.internalError(err, "Failed to authorize")
	}

	switch authz {
	case auth.AuthorizationAuthenticatedUser:
		if err = s.PermissionTeamWorker(md, req.CompanyUuid, req.TeamUuid); err != nil {
			return nil, err
		}
	case auth.AuthorizationSupportUser:
	default:
		return nil, grpc.Errorf(codes.PermissionDenied, "you do not have access to this service")
	}

	if _, err = s.GetTeam(ctx, &pb.GetTeamRequest{Uuid: req.TeamUuid, CompanyUuid: req.CompanyUuid}); err != nil {
		return nil, err
	}

	shiftStartAfter, err := time.Parse(time.RFC3339, req.ShiftStartAfter)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "shift_start_after could not be parsed")
	}

	shiftStartBefore, err := time.Parse(time.RFC3339, req.ShiftStartBefore)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "shift_start_before could not be parsed")
	}

	if shiftStartAfter.After(shiftStartBefore) {
		return nil, grpc.Errorf(codes.InvalidArgument, "shift_start_after must be before shift_start_before")
	}
	res := &pb.ShiftList{ShiftStartAfter: shiftStartAfter, ShiftStartBefore: shiftStartBefore}

	var dbShifts []pb.Shift

	if req.UserUuid != "" && req.JobUuid == "" {
		if _, err = s.dbMap.Select(&dbShifts, "select * from shift where team_uuid=? and user_uuid = ? AND start>=? AND start<?", req.TeamUuid, req.UserUuid, shiftStartAfter, shiftStartBefore); err != nil {
			return nil, s.internalError(err, "unable to query database")
		}
	}

	if req.JobUuid != "" && req.UserUuid == "" {
		if _, err = s.dbMap.Select(&dbShifts, "select * from shift where team_uuid=? and job_uuid = ? AND start>=? AND start<?", req.TeamUuid, req.JobUuid, shiftStartAfter, shiftStartBefore); err != nil {
			return nil, s.internalError(err, "unable to query database")
		}
	}

	if req.JobUuid != "" && req.UserUuid != "" {
		if _, err = s.dbMap.Select(&dbShifts, "select * from shift where team_uuid=? and user_uuid = ? and job_uuid = ? AND start>=? AND start<?", req.TeamUuid, req.UserUuid, req.JobUuid, shiftStartAfter, shiftStartBefore); err != nil {
			return nil, s.internalError(err, "unable to query database")
		}
	}

	if req.JobUuid == "" && req.UserUuid == "" {
		if _, err = s.dbMap.Select(&dbShifts, "select * from shift where team_uuid=? AND start>=? AND start<?", req.TeamUuid, shiftStartAfter, shiftStartBefore); err != nil {
			return nil, s.internalError(err, "unable to query database")
		}
	}

	for _, shift := range dbShifts {
		shift.CompanyUuid = req.CompanyUuid
		res.Shifts = append(res.Shifts, shift)
	}
	return res, nil
}

// quickTime is a helper method to show elapsed time in a request
func quickTime(s time.Time) int {
	endTime := time.Now()
	return int(((endTime.Sub(s)).Nanoseconds() / 1000000))
}

func (s *companyServer) BulkPublishShifts(ctx context.Context, req *pb.BulkPublishShiftsRequest) (*pb.ShiftList, error) {
	startTime := time.Now()
	logger.Infof("time so far %v", quickTime(startTime))

	orig, err := s.ListShifts(ctx, &pb.ShiftListRequest{CompanyUuid: req.CompanyUuid,
		TeamUuid: req.TeamUuid, UserUuid: req.UserUuid, JobUuid: req.JobUuid,
		ShiftStartAfter:  req.ShiftStartAfter.Format(time.RFC3339),
		ShiftStartBefore: req.ShiftStartBefore.Format(time.RFC3339)})
	if err != nil {
		return nil, err
	}

	// create a new context so we can suppress notifications in the "UpdateShift" endpoint
	newMd := metadata.New(map[string]string{"suppressnotification": "true"})
	oldMd, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, s.internalError(err, "context missing metadata")
	}
	combinedMd := metadata.Join(oldMd, newMd)
	newCtx := metadata.NewOutgoingContext(context.Background(), combinedMd)

	res := &pb.ShiftList{ShiftStartAfter: req.ShiftStartAfter, ShiftStartBefore: req.ShiftStartBefore}
	// Keep track of notifications  - user to orig shift
	notifs := make(map[string][]*pb.Shift)

	logger.Infof("before shifts update %v", quickTime(startTime))

	for _, shift := range orig.Shifts {
		// keep track of what changed for messaging purposes
		if shift.UserUuid != "" && shift.Published != req.Published && shift.Start.After(time.Now()) {
			copy := shift
			notifs[shift.UserUuid] = append(notifs[shift.UserUuid], &copy)
		}
		// do the change
		shift.Published = req.Published
		go func(sh pb.Shift) {
			_, err = s.UpdateShift(newCtx, &sh)
			if err != nil {
				s.internalError(err, "failed to alert worker about new shifts")
			}
		}(shift)
		res.Shifts = append(res.Shifts, shift)
	}
	logger.Infof("before shifts notifications %v", quickTime(startTime))

	go func() {
		s.logger.Debugf("starting bulk shift notifications %v", notifs)
		for userUUID, shifts := range notifs {
			botClient, close, err := bot.NewClient()
			if err != nil {
				s.internalError(err, "unable to initiate bot connection")
				return
			}
			defer close()
			if req.Published {
				// alert published
				if _, err = botClient.AlertNewShifts(asyncContext(), &bot.AlertNewShiftsRequest{UserUuid: userUUID, NewShifts: shifts}); err != nil {
					s.internalError(err, "failed to alert worker about new shifts")
				}
			} else {
				// alert removed
				if _, err = botClient.AlertRemovedShifts(asyncContext(), &bot.AlertRemovedShiftsRequest{UserUuid: userUUID, OldShifts: shifts}); err != nil {
					s.internalError(err, "failed to alert worker about removed shifts")
				}

			}
		}
	}()
	logger.Infof("total time %v", quickTime(startTime))

	return res, nil
}

func (s *companyServer) GetShift(ctx context.Context, req *pb.GetShiftRequest) (*pb.Shift, error) {
	md, authz, err := getAuth(ctx)
	if err != nil {
		return nil, s.internalError(err, "failed to authorize")
	}

	switch authz {
	case auth.AuthorizationAuthenticatedUser:
		if err = s.PermissionTeamWorker(md, req.CompanyUuid, req.TeamUuid); err != nil {
			return nil, err
		}
	case auth.AuthorizationSupportUser:
	default:
		return nil, grpc.Errorf(codes.PermissionDenied, "you do not have access to this service")
	}

	if _, err = s.GetTeam(ctx, &pb.GetTeamRequest{Uuid: req.TeamUuid, CompanyUuid: req.CompanyUuid}); err != nil {
		return nil, err
	}

	obj, err := s.dbMap.Get(pb.Shift{}, req.Uuid)
	if err != nil {
		return nil, s.internalError(err, "unable to query database")
	} else if obj == nil {
		return nil, grpc.Errorf(codes.NotFound, "shift not found")
	}
	shift := obj.(*pb.Shift)
	shift.CompanyUuid = req.CompanyUuid
	shift.TeamUuid = req.TeamUuid
	return shift, nil

}

// noChange returns true if the shifts are the same
func (s *companyServer) noChange(s1 *pb.Shift, s2 *pb.Shift) bool {
	if s1.CompanyUuid != s2.CompanyUuid {
		return false
	}

	if s1.TeamUuid != s2.TeamUuid {
		return false
	}

	if s1.Start != s2.Start {
		return false
	}

	if s1.Stop != s2.Stop {
		return false
	}

	if s1.JobUuid != s2.JobUuid {
		return false
	}

	if s1.UserUuid != s2.UserUuid {
		return false
	}

	if s1.Published != s2.Published {
		return false
	}

	return true
}

func (s *companyServer) UpdateShift(ctx context.Context, req *pb.Shift) (*pb.Shift, error) {
	md, authz, err := getAuth(ctx)
	switch authz {
	case auth.AuthorizationAuthenticatedUser:
		if err = s.PermissionCompanyAdmin(md, req.CompanyUuid); err != nil {
			return nil, err
		}
	case auth.AuthorizationSupportUser:
	default:
		return nil, grpc.Errorf(codes.PermissionDenied, "you do not have access to this service")
	}

	orig, err := s.GetShift(ctx, &pb.GetShiftRequest{Uuid: req.Uuid, TeamUuid: req.TeamUuid, CompanyUuid: req.CompanyUuid})
	if err != nil {
		return nil, grpc.Errorf(codes.NotFound, "shift not found")
	}

	if s.noChange(orig, req) {
		return req, nil
	}

	if req.JobUuid != "" {
		if _, err = s.GetJob(ctx, &pb.GetJobRequest{Uuid: req.JobUuid, CompanyUuid: req.CompanyUuid, TeamUuid: req.TeamUuid}); err != nil {
			return nil, err
		}
	}

	if req.UserUuid != "" {
		if _, err = s.GetDirectoryEntry(ctx, &pb.DirectoryEntryRequest{CompanyUuid: req.CompanyUuid, UserUuid: req.UserUuid}); err != nil {
			return nil, err
		}
	}

	dur := req.Stop.Sub(req.Start)
	if dur <= 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "stop must be after start")
	} else if dur > maxShiftDuration {
		return nil, grpc.Errorf(codes.InvalidArgument, "duration exceeds max %f hour duration", maxShiftDuration.Hours())
	}

	if _, err := s.dbMap.Update(req); err != nil {
		return nil, s.internalError(err, "could not update the shift")
	}
	al := newAuditEntry(md, "shift", req.Uuid, req.CompanyUuid, req.TeamUuid)
	al.OriginalContents = orig
	al.UpdatedContents = req
	al.Log(logger, "updated shift")
	go helpers.TrackEventFromMetadata(md, "shift_updated")
	if !orig.Published && req.Published {
		go helpers.TrackEventFromMetadata(md, "shift_published")
	}

	go func() {
		botClient, close, err := bot.NewClient()
		if err != nil {
			s.internalError(err, "unable to initiate bot connection")
			return
		}
		defer close()
		// Send bot notifications
		switch {
		case len(md["suppressnotification"]) != 0:
			// The BulkPublishShifts endpoint is suppressing this endpoint's notifications
		case orig.Published == false && req.Published == true:
			if req.Start.After(time.Now()) && req.UserUuid != "" {
				// looks like a new shift
				if _, err = botClient.AlertNewShift(asyncContext(), &bot.AlertNewShiftRequest{UserUuid: req.UserUuid, NewShift: req}); err != nil {
					s.internalError(err, "failed to alert worker about new shift")
				}
			}

		case orig.Published == true && req.Published == false:
			if orig.Start.After(time.Now()) && orig.UserUuid != "" {
				// removed a shift
				if _, err = botClient.AlertRemovedShift(asyncContext(), &bot.AlertRemovedShiftRequest{UserUuid: orig.UserUuid, OldShift: orig}); err != nil {
					s.internalError(err, "failed to alert worker about removed shift")
				}
			}
		case orig.Published == false && req.Published == false:
			// NOOP - basically return
		case orig.UserUuid == req.UserUuid:
			if orig.UserUuid != "" && req.Start.After(time.Now()) {
				if _, err = botClient.AlertChangedShift(asyncContext(), &bot.AlertChangedShiftRequest{UserUuid: orig.UserUuid, OldShift: orig, NewShift: req}); err != nil {
					s.internalError(err, "failed to alert worker about changed shift")
				}
			}
		case orig.UserUuid != req.UserUuid:
			if orig.UserUuid != "" && orig.Start.After(time.Now()) {
				if _, err = botClient.AlertRemovedShift(asyncContext(), &bot.AlertRemovedShiftRequest{UserUuid: orig.UserUuid, OldShift: orig}); err != nil {
					s.internalError(err, "failed to alert worker about removed shift")
				}
			}
			if req.UserUuid != "" && req.Start.After(time.Now()) {
				if _, err = botClient.AlertNewShift(asyncContext(), &bot.AlertNewShiftRequest{UserUuid: req.UserUuid, NewShift: req}); err != nil {
					s.internalError(err, "failed to alert worker about new shift")
				}
			}
		default:
			logger.Errorf("unable to determine updated shift messaging - orig %v new %v", orig, req)
		}
	}()

	return req, nil
}

func (s *companyServer) DeleteShift(ctx context.Context, req *pb.GetShiftRequest) (*empty.Empty, error) {
	md, authz, err := getAuth(ctx)
	if err != nil {
		return nil, s.internalError(err, "Failed to authorize")
	}

	switch authz {
	case auth.AuthorizationAuthenticatedUser:
		if err = s.PermissionTeamWorker(md, req.CompanyUuid, req.TeamUuid); err != nil {
			return nil, err
		}
	case auth.AuthorizationSupportUser:
	default:
		return nil, grpc.Errorf(codes.PermissionDenied, "You do not have access to this service")
	}

	orig, err := s.GetShift(ctx, &pb.GetShiftRequest{Uuid: req.Uuid, TeamUuid: req.TeamUuid, CompanyUuid: req.CompanyUuid})
	if err != nil {
		return nil, err
	}
	if _, err = s.dbMap.Delete(orig); err != nil {
		return nil, s.internalError(err, "failed to delete shift")
	}

	al := newAuditEntry(md, "shift", req.Uuid, req.CompanyUuid, req.TeamUuid)
	al.OriginalContents = orig
	al.Log(logger, "deleted shift")

	go func() {
		if orig.UserUuid != "" && orig.Published && orig.Start.After(time.Now()) {
			botClient, close, err := bot.NewClient()
			if err != nil {
				s.internalError(err, "unable to initiate bot connection")
				return
			}
			defer close()
			if _, err = botClient.AlertRemovedShift(asyncContext(), &bot.AlertRemovedShiftRequest{UserUuid: orig.UserUuid, OldShift: orig}); err != nil {
				s.internalError(err, "failed to alert worker about removed shift")
			}
		}
	}()
	go helpers.TrackEventFromMetadata(md, "shift_deleted")

	return &empty.Empty{}, nil
}
