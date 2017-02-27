package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"v2.staffjoy.com/auth"
)

// Each permission has a public convenience checker, and a private relationship checker.
// Recall that support users have a different authorization, and will not use these functions.

// PermissionCompanyAdmin checks that the current user is an admin of the given company
func (s *companyServer) PermissionCompanyAdmin(md metadata.MD, companyUUID string) error {
	userUUID, err := auth.GetCurrentUserUUIDFromMetadata(md)
	if err != nil {
		return s.internalError(err, "failed to find current user uuid")

	}

	ok, err := s.checkCompanyAdmin(userUUID, companyUUID)
	if err != nil {
		return s.internalError(err, "failed to check company admin permissions")
	}
	if !ok {
		return grpc.Errorf(codes.PermissionDenied, "you do not have admin access to this service")
	}
	return nil
}

// PermissionTeamWorker checks whether a user is a worker of a given team in a given company, or is an admin of that company
func (s *companyServer) PermissionTeamWorker(md metadata.MD, companyUUID, teamUUID string) error {
	userUUID, err := auth.GetCurrentUserUUIDFromMetadata(md)
	if err != nil {
		return s.internalError(err, "failed to find current user uuid")
	}

	// Check if worker
	ok, err := s.checkCompanyAdmin(userUUID, companyUUID)
	if err != nil {
		return s.internalError(err, "failed to check company admin permissions")
	}
	if ok {
		// Admin - allow access
		return nil
	}

	// Not admin - check if worker
	ok, err = s.checkTeamWorker(userUUID, teamUUID)
	if err != nil {
		return s.internalError(err, "failed to check team member permissions")
	}
	if ok {
		// worker in team - allow access
		return nil
	}
	return grpc.Errorf(codes.PermissionDenied, "you do not have worker access to this team")
}

// PermissionCompanyDirectory checks whether a user exists in the directory of a company. It is the lowest level of security.
// The user may no longer be associated with a team (i.e. may be a former employee)
func (s *companyServer) PermissionCompanyDirectory(md metadata.MD, companyUUID string) error {
	userUUID, err := auth.GetCurrentUserUUIDFromMetadata(md)
	if err != nil {
		return s.internalError(err, "failed to find current user uuid")

	}

	// Admins are in directory, so this is all we have to check
	ok, err := s.checkInDirectory(userUUID, companyUUID)
	if err != nil {
		return s.internalError(err, "failed to check directory existence")
	}
	if !ok {
		return grpc.Errorf(codes.PermissionDenied, "you are not associated with this company")
	}
	return nil
}

func (s *companyServer) checkCompanyAdmin(userUUID string, companyUUID string) (ok bool, err error) {
	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM admin WHERE (company_uuid=? AND user_uuid=?))", companyUUID, userUUID).Scan(&ok)
	return
}

func (s *companyServer) checkTeamWorker(userUUID string, teamUUID string) (ok bool, err error) {
	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM worker WHERE (team_uuid=? AND user_uuid=?))", teamUUID, userUUID).Scan(&ok)
	return
}

func (s *companyServer) checkInDirectory(userUUID string, companyUUID string) (ok bool, err error) {
	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM directory WHERE (company_uuid=? AND user_uuid=?))", companyUUID, userUUID).Scan(&ok)
	return
}
