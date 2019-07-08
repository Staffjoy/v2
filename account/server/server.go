package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	intercom "gopkg.in/intercom/intercom-go.v2"

	"github.com/go-gorp/gorp"
	"github.com/sirupsen/logrus"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"github.com/golang/protobuf/ptypes/empty"
	pb "v2.staffjoy.com/account"
	"v2.staffjoy.com/auth"
	"v2.staffjoy.com/company"
	"v2.staffjoy.com/crypto"
	"v2.staffjoy.com/email"
	"v2.staffjoy.com/environments"
	"v2.staffjoy.com/helpers"
	"v2.staffjoy.com/sms"
)

const (
	minPasswordLength = 6
)

type accountServer struct {
	config       environments.Config
	logger       *logrus.Entry
	db           *sql.DB
	errorClient  environments.SentryClient
	signingToken string
	dbMap        *gorp.DbMap
	smsClient    sms.SmsServiceClient
}

// GetOrCreate is for internal use by other APIs to match a user based on their phonenumber or email.
func (s *accountServer) GetOrCreate(ctx context.Context, req *pb.GetOrCreateRequest) (*pb.Account, error) {
	// rely on downstream permissions

	var err error
	req.Email = strings.ToLower(req.Email)
	if req.Phonenumber, err = ParseAndFormatPhonenumber(req.Phonenumber); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Invalid phone number")
	}
	// check for existing user
	var existingUserUUID string
	if len(req.Email) > 0 {
		err = s.db.QueryRow("SELECT uuid FROM account WHERE email=?", req.Email).Scan(&existingUserUUID)
		if err != nil && err != sql.ErrNoRows {
			return nil, s.internalError(err, "failed to query database for existing email")
		}
	}
	if len(req.Phonenumber) > 0 && existingUserUUID == "" {
		err = s.db.QueryRow("SELECT uuid FROM account WHERE phonenumber=?", req.Phonenumber).Scan(&existingUserUUID)
		if err != nil && err != sql.ErrNoRows {
			return nil, s.internalError(err, "failed to query database for existing phonenumber")
		}

	}
	if existingUserUUID != "" {
		return s.Get(ctx, &pb.GetAccountRequest{Uuid: existingUserUUID})
	}
	return s.Create(ctx, &pb.CreateAccountRequest{Phonenumber: req.Phonenumber, Name: req.Name, Email: req.Email})
}

func (s *accountServer) GetAccountByPhonenumber(ctx context.Context, req *pb.GetAccountByPhonenumberRequest) (*pb.Account, error) {
	// rely on downstream permissions

	var err error
	if req.Phonenumber, err = ParseAndFormatPhonenumber(req.Phonenumber); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Invalid phone number")
	}
	if req.Phonenumber == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "No phonenumber provided")
	}

	var uuid string
	err = s.db.QueryRow("SELECT uuid FROM account WHERE phonenumber=?", req.Phonenumber).Scan(&uuid)
	if err == sql.ErrNoRows {
		return nil, grpc.Errorf(codes.NotFound, "")
	} else if err != nil {
		return nil, s.internalError(err, "failed to query database for existing phonenumber")
	}
	return s.Get(ctx, &pb.GetAccountRequest{Uuid: uuid})
}

func (s *accountServer) Create(ctx context.Context, req *pb.CreateAccountRequest) (*pb.Account, error) {
	md, authz, err := getAuth(ctx)
	if err != nil {
		return nil, s.internalError(err, "Failed to authorize")
	}
	al := newAuditEntry(md, "account", "")

	switch authz {
	case auth.AuthorizationSupportUser:
	case auth.AuthorizationWWWService:
	case auth.AuthorizationCompanyService:
	default:
		return nil, grpc.Errorf(codes.PermissionDenied, "You do not have access to this service")
	}

	if (len(req.Email) + len(req.Name) + len(req.Phonenumber)) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Empty request")
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if len(req.Email) > 0 && strings.Index(req.Email, "@") == -1 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Invalid email")
	}
	req.Phonenumber, err = ParseAndFormatPhonenumber(req.Phonenumber)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Invalid phone number")
	}

	if req.Email != "" {
		// Check to see if account exists
		var existingUser string
		err = s.db.QueryRow("SELECT uuid FROM account WHERE email=?", req.Email).Scan(&existingUser)
		// We expect an sql.ErrNoRows, which means that the user doesn't exist.
		if err == nil {
			return nil, grpc.Errorf(codes.AlreadyExists, "A user with that email already exists. Try a password reset")
		} else if err != sql.ErrNoRows {
			return nil, s.internalError(err, "An unknown error occurred while searching for that email.")
		}
	}
	if req.Phonenumber != "" {
		_, err = s.GetAccountByPhonenumber(ctx, &pb.GetAccountByPhonenumberRequest{Phonenumber: req.Phonenumber})
		if err == nil {
			return nil, grpc.Errorf(codes.AlreadyExists, "A user with that phonenumber already exists. Try a password reset.")
		} else if grpc.Code(err) != codes.NotFound {
			return nil, s.internalError(err, "An unknown error occurred")
		}
	}
	// Create the account
	uuid, err := crypto.NewUUID()
	if err != nil {
		return nil, s.internalError(err, "Cannot generate a user id")
	}

	a := &pb.Account{Uuid: uuid.String(), Email: req.Email, Name: req.Name, Phonenumber: req.Phonenumber}
	a.PhotoUrl = GenerateGravatarURL(a.Email)
	a.MemberSince = time.Now()

	if err = s.dbMap.Insert(a); err != nil {
		return nil, s.internalError(err, "Could not create user account")
	}
	go s.SyncUser(ctx, &pb.SyncUserRequest{Uuid: a.Uuid})

	if len(req.Email) > 0 {
		// Email confirmation
		token, err := crypto.EmailConfirmationToken(a.Uuid, a.Email, s.signingToken)
		if err != nil {
			return nil, s.internalError(err, "Could not create token")
		}
		link := url.URL{Host: "www." + config.ExternalApex, Path: fmt.Sprintf("/activate/%s", token), Scheme: "http"}
		// Send verification email
		emailName := req.Name
		if emailName == "" {
			emailName = "there"
		}
		msg := &email.EmailRequest{
			To:       a.Email,
			Name:     a.Name,
			Subject:  "Activate your Staffjoy account",
			HtmlBody: fmt.Sprintf(activateAccountTmpl, emailName, link.String(), link.String(), link.String()),
		}
		mailer, close, err := email.NewClient()
		if err != nil {
			return nil, s.internalError(err, "unable to initiate email service connection")
		}
		defer close()

		_, err = mailer.Send(ctx, msg)
		if err != nil {
			return nil, s.internalError(err, "Unable to send email")
		}
	}

	// todo - sms onboarding (if worker??)

	al.TargetUUID = a.Uuid
	al.UpdatedContents = a
	al.Log(logger, "created account")

	return a, nil
}

func (s *accountServer) List(ctx context.Context, req *pb.GetAccountListRequest) (*pb.AccountList, error) {
	_, authz, err := getAuth(ctx)
	if err != nil {
		return nil, s.internalError(err, "Failed to authorize")
	}
	switch authz {
	case auth.AuthorizationSupportUser:
	default:
		return nil, grpc.Errorf(codes.PermissionDenied, "You do not have access to this service")
	}

	if req.Offset < 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Invalid offset - must be greater than or equal to zero")
	}
	if req.Limit <= 0 {
		// Set a default
		req.Limit = 10
	}
	res := &pb.AccountList{Limit: req.Limit, Offset: req.Offset}

	rows, err := s.db.Query("select uuid from account limit ? offset ?", req.Limit, req.Offset)
	if err != nil {
		return nil, s.internalError(err, "Unable to query database")
	}

	for rows.Next() {
		r := &pb.GetAccountRequest{}
		if err := rows.Scan(&r.Uuid); err != nil {
			return nil, s.internalError(err, "Error scanning database")
		}

		var a *pb.Account
		if a, err = s.Get(ctx, r); err != nil {
			return nil, err
		}
		res.Accounts = append(res.Accounts, *a)
	}
	return res, nil
}

func (s *accountServer) Get(ctx context.Context, req *pb.GetAccountRequest) (*pb.Account, error) {
	md, authz, err := getAuth(ctx)
	if err != nil {
		return nil, s.internalError(err, "Failed to authorize")
	}
	switch authz {
	case auth.AuthorizationWWWService:
	case auth.AuthorizationAccountService:
	case auth.AuthorizationCompanyService:
	case auth.AuthorizationWhoamiService:
	case auth.AuthorizationBotService:
	case auth.AuthorizationAuthenticatedUser:
		uuid, err := auth.GetCurrentUserUUIDFromMetadata(md)
		if err != nil {
			return nil, s.internalError(err, "failed to find current user uuid %v", md)
		}
		if uuid != req.Uuid {
			return nil, grpc.Errorf(codes.PermissionDenied, "You do not have access to this service")
		}
	case auth.AuthorizationSupportUser:
	case auth.AuthorizationSuperpowersService:
		if s.config.Name != "development" {
			s.logger.Warningf("Development service trying to connect outside development environment")
			return nil, grpc.Errorf(codes.PermissionDenied, "This service is not available outside development environments")
		}
	default:
		return nil, grpc.Errorf(codes.PermissionDenied, "You do not have access to this service")
	}

	if req.Uuid == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "uuid must be specified")
	}
	obj, err := s.dbMap.Get(pb.Account{}, req.Uuid)
	if err != nil {
		return nil, s.internalError(err, "Unable to query database")
	} else if obj == nil {
		return nil, grpc.Errorf(codes.NotFound, "User with id %s not found", req.Uuid)
	}
	return obj.(*pb.Account), nil
}

func (s *accountServer) Update(ctx context.Context, req *pb.Account) (*pb.Account, error) {
	md, authz, err := getAuth(ctx)
	if err != nil {
		return nil, s.internalError(err, "Failed to authorize")
	}
	switch authz {
	case auth.AuthorizationWWWService:
	case auth.AuthorizationCompanyService:
	case auth.AuthorizationAuthenticatedUser:
		uuid, err := auth.GetCurrentUserUUIDFromMetadata(md)
		if err != nil {
			return nil, s.internalError(err, "failed to find current user uuid")

		}
		if uuid != req.Uuid {
			return nil, grpc.Errorf(codes.PermissionDenied, "You do not have access to this service")
		}
	case auth.AuthorizationSupportUser:
	case auth.AuthorizationSuperpowersService:
		if s.config.Name != "development" {
			s.logger.Warningf("Development service trying to connect outside development environment")
			return nil, grpc.Errorf(codes.PermissionDenied, "This service is not available outside development environments")
		}
	default:
		return nil, grpc.Errorf(codes.PermissionDenied, "You do not have access to this service")
	}
	al := newAuditEntry(md, "account", req.Uuid)

	existing, err := s.Get(ctx, &pb.GetAccountRequest{Uuid: req.Uuid})
	if err != nil {
		// This handles 404 and everything!
		return nil, err
	}
	al.OriginalContents = existing

	// Some validations
	if req.Phonenumber, err = ParseAndFormatPhonenumber(req.Phonenumber); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Invalid phone number")
	}
	req.Email = strings.ToLower(req.Email)
	if !req.MemberSince.Equal(existing.MemberSince) {
		return nil, grpc.Errorf(codes.PermissionDenied, "You cannot modify the member_since date")
	}
	if req.Email != "" && (req.Email != existing.Email) {
		// Check to see if account exists
		var existingUser string
		err = s.db.QueryRow("SELECT uuid FROM account WHERE email=?", req.Email).Scan(&existingUser)
		// We expect an sql.ErrNoRows, which means that the user doesn't exist.
		if err == nil {
			return nil, grpc.Errorf(codes.AlreadyExists, "A user with that email already exists. Try a password reset")
		} else if err != sql.ErrNoRows {
			return nil, s.internalError(err, "An unknown error occurred while searching for that email.")
		}
	}
	if req.Phonenumber != "" && (req.Phonenumber != existing.Phonenumber) {
		_, err = s.GetAccountByPhonenumber(ctx, &pb.GetAccountByPhonenumberRequest{Phonenumber: req.Phonenumber})
		if err == nil {
			return nil, grpc.Errorf(codes.AlreadyExists, "A user with that phonenumber already exists. Try a password reset.")
		} else if grpc.Code(err) != codes.NotFound {
			return nil, s.internalError(err, "An unknown error occurred")
		}
	}

	if authz == auth.AuthorizationAuthenticatedUser {
		if (req.ConfirmedAndActive != existing.ConfirmedAndActive) && (existing.ConfirmedAndActive == false) {
			return nil, grpc.Errorf(codes.PermissionDenied, "You cannot activate this account")
		}
		if req.Support != existing.Support {
			return nil, grpc.Errorf(codes.PermissionDenied, "You cannot change the support parameter")
		}
		if req.PhotoUrl != existing.PhotoUrl {
			return nil, grpc.Errorf(codes.PermissionDenied, "You cannot change the photo through this endpoint (see docs)")
		}
		// User can request email change - not do it :-)
		if req.Email != existing.Email {
			s.RequestEmailChange(ctx, &pb.EmailChangeRequest{Uuid: req.Uuid, Email: req.Email})
			// revert
			req.Email = existing.Email
		}
	}

	req.PhotoUrl = GenerateGravatarURL(req.Email)

	if _, err := s.dbMap.Update(req); err != nil {
		return nil, s.internalError(err, "Could not update the user account")
	}

	go s.SyncUser(ctx, &pb.SyncUserRequest{Uuid: req.Uuid})

	al.UpdatedContents = req
	al.Log(logger, "updated account")

	// If account is being activated, or if phone number is changed by current user - send text
	if req.ConfirmedAndActive && len(req.Phonenumber) > 0 && req.Phonenumber != existing.Phonenumber {
		s.sendSmsGreeting(req.Phonenumber)
	}

	go helpers.TrackEventFromMetadata(md, "account_updated")

	return req, nil
}

func (s *accountServer) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*empty.Empty, error) {
	md, authz, err := getAuth(ctx)
	if err != nil {
		return nil, s.internalError(err, "Failed to authorize")
	}
	switch authz {
	case auth.AuthorizationAuthenticatedUser:
		uuid, err := auth.GetCurrentUserUUIDFromMetadata(md)
		if err != nil {
			return nil, s.internalError(err, "failed to find current user uuid")

		}
		if uuid != req.Uuid {
			return nil, grpc.Errorf(codes.PermissionDenied, "You do not have access to this service")
		}
	case auth.AuthorizationWWWService:
	case auth.AuthorizationSupportUser:
	default:
		return nil, grpc.Errorf(codes.PermissionDenied, "You do not have access to this service")
	}
	al := newAuditEntry(md, "account", req.Uuid)

	// Verify inputs
	if req.Uuid == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Invalid UUID")
	}
	if len(req.Password) < minPasswordLength {
		return nil, grpc.Errorf(codes.InvalidArgument, "Invalid password - it must be at least %d characters long", minPasswordLength)
	}
	salt, err := crypto.NewSalt()
	if err != nil {
		return nil, s.internalError(err, "Failed to generate a salt")
	}

	pwHash, err := crypto.HashPassword(salt, []byte(req.Password))
	if err != nil {
		return nil, s.internalError(err, "Failed to hash the password")
	}

	// Run the update . . .
	// We can see whether the user exists in a single step - by just
	// running the query and seeing how many rows are affectd
	res, err := s.db.Exec("UPDATE account SET password_hash=?, password_salt=? where uuid=? limit 1", pwHash, salt, req.Uuid)
	if err != nil {
		return nil, s.internalError(err, "failed to query the database to update the password hash")
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, s.internalError(err, "Failed to read the database")
	}
	if affected != 1 {
		return nil, grpc.Errorf(codes.NotFound, "")
	}
	al.Log(logger, "updated password")
	go helpers.TrackEventFromMetadata(md, "password_updated")
	return &empty.Empty{}, nil
}

func (s *accountServer) VerifyPassword(ctx context.Context, req *pb.VerifyPasswordRequest) (*pb.Account, error) {
	// Prep
	_, authz, err := getAuth(ctx)
	if err != nil {
		return nil, s.internalError(err, "Failed to authorize")
	}
	switch authz {
	case auth.AuthorizationWWWService:
	case auth.AuthorizationSupportUser:
	default:
		return nil, grpc.Errorf(codes.PermissionDenied, "You do not have access to this service")
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	var dbHash, salt, uuid sql.NullString
	var confirmedAndActive sql.NullBool

	err = s.db.QueryRow("SELECT uuid, password_hash, password_salt, confirmed_and_active FROM account WHERE email=?", req.Email).Scan(&uuid, &dbHash, &salt, &confirmedAndActive)
	// We expect an sql.ErrNoRows, which means that the user doesn't exist.
	switch {
	case err == sql.ErrNoRows:
		return nil, grpc.Errorf(codes.NotFound, "")
	case err != nil:
		return nil, s.internalError(err, "Unable to query database")
	default:
		if !confirmedAndActive.Bool {
			return nil, grpc.Errorf(codes.PermissionDenied, "This user has not confirmed their account")
		}
		if len(dbHash.String) == 0 {
			return nil, grpc.Errorf(codes.PermissionDenied, "This user has not set up their password ")
		}

		if err != nil {
			return nil, s.internalError(err, "Unable to determine account creation time")
		}

		if crypto.CheckPasswordHash([]byte(dbHash.String), []byte(salt.String), []byte(req.Password)) != nil {
			return nil, grpc.Errorf(codes.Unauthenticated, "Incorrect password")
		}

		a, err := s.Get(ctx, &pb.GetAccountRequest{Uuid: uuid.String})
		if err != nil {
			return nil, s.internalError(err, "Unable to query account")
		}

		// You shall pass
		return a, nil
	}
}

// RequestPasswordReset sends an email to a user with a password reset link
func (s *accountServer) RequestPasswordReset(ctx context.Context, req *pb.PasswordResetRequest) (*empty.Empty, error) {
	_, authz, err := getAuth(ctx)
	if err != nil {
		return nil, s.internalError(err, "Failed to authorize")
	}
	switch authz {
	case auth.AuthorizationWWWService:
	case auth.AuthorizationSupportUser:
	default:
		return nil, grpc.Errorf(codes.PermissionDenied, "You do not have access to this service")
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if len(req.Email) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "No UUID provided")
	}

	var existingUser string
	err = s.db.QueryRow("SELECT uuid FROM account WHERE email=?", req.Email).Scan(&existingUser)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "No user with that email exists")
	}
	a, err := s.Get(ctx, &pb.GetAccountRequest{Uuid: existingUser})

	token, err := crypto.EmailConfirmationToken(a.Uuid, a.Email, s.signingToken)
	if err != nil {
		return nil, s.internalError(err, "Could not create token")
	}

	link := url.URL{Host: "www." + config.ExternalApex, Scheme: "http"}
	message := "Reset your Staffjoy password"
	tmpl := resetPasswordTmpl
	if a.ConfirmedAndActive {
		link.Path = fmt.Sprintf("/reset/%s", token)
	} else {
		// Not actually active - make some tweaks for activate instead of password reset
		link.Path = fmt.Sprintf("/activate/%s", token)
		message = "Activate your Staffjoy account"
		tmpl = activateAccountTmpl
	}

	// Send verification email
	msg := &email.EmailRequest{
		To:       a.Email,
		Name:     a.Name,
		Subject:  message,
		HtmlBody: fmt.Sprintf(tmpl, link.String(), link.String()),
	}
	mailer, close, err := email.NewClient()
	if err != nil {
		panic(err)
	}
	defer close()

	_, err = mailer.Send(ctx, msg)
	if err != nil {
		return nil, s.internalError(err, "Unable to send email")
	}

	return &empty.Empty{}, nil
}

// RequestPasswordReset sends an email to a user with a password reset link
func (s *accountServer) RequestEmailChange(ctx context.Context, req *pb.EmailChangeRequest) (*empty.Empty, error) {
	md, authz, err := getAuth(ctx)
	if err != nil {
		return nil, s.internalError(err, "Failed to authorize")
	}
	switch authz {
	case auth.AuthorizationAuthenticatedUser:
		uuid, err := auth.GetCurrentUserUUIDFromMetadata(md)
		if err != nil {
			return nil, s.internalError(err, "failed to find current user uuid")

		}
		if uuid != req.Uuid {
			return nil, grpc.Errorf(codes.PermissionDenied, "You do not have access to this service")
		}
	case auth.AuthorizationSupportUser:
	default:
		return nil, grpc.Errorf(codes.PermissionDenied, "You do not have access to this service")
	}

	newEmail := strings.TrimSpace(strings.ToLower(req.Email))
	if len(req.Uuid) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "No UUID provided")
	}

	a, err := s.Get(ctx, &pb.GetAccountRequest{Uuid: req.Uuid})
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "No user with that email exists")
	}

	token, err := crypto.EmailConfirmationToken(a.Uuid, newEmail, s.signingToken)
	if err != nil {
		return nil, s.internalError(err, "Could not create token")
	}
	link := url.URL{Host: "www." + config.ExternalApex, Path: fmt.Sprintf("/activate/%s", token), Scheme: "http"}
	// Send verification email
	msg := &email.EmailRequest{
		To:       a.Email,
		Name:     a.Name,
		Subject:  "Confirm Your New Email Address",
		HtmlBody: fmt.Sprintf(confirmEmailTmpl, a.Name, link.String(), link.String(), link.String()),
	}
	mailer, close, err := email.NewClient()
	if err != nil {
		return nil, s.internalError(err, "unable to initiate email service connection")
	}
	defer close()

	if _, err = mailer.Send(ctx, msg); err != nil {
		return nil, s.internalError(err, "Unable to send email")
	}
	return &empty.Empty{}, nil
}

// ChangeEmail sets an account to active and updates its email. It is
// used after a user clicks a confirmation link in their email.
func (s *accountServer) ChangeEmail(ctx context.Context, req *pb.EmailConfirmation) (*empty.Empty, error) {
	md, authz, err := getAuth(ctx)
	if err != nil {
		return nil, s.internalError(err, "Failed to authorize")
	}
	switch authz {
	case auth.AuthorizationWWWService:
	case auth.AuthorizationSupportUser:
	default:
		return nil, grpc.Errorf(codes.PermissionDenied, "You do not have access to this service")
	}
	al := newAuditEntry(md, "account", req.Uuid)

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	res, err := s.db.Exec("UPDATE account SET email=?, confirmed_and_active=true where uuid=? limit 1", req.Email, req.Uuid)

	if err != nil {
		return nil, s.internalError(err, "Unable to query database")
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, s.internalError(err, "Failed to read the database")
	}
	if affected != 1 {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	go s.SyncUser(ctx, &pb.SyncUserRequest{Uuid: req.Uuid})

	al.UpdatedContents = req.Email
	al.Log(logger, "changed email")
	go helpers.TrackEventFromMetadata(md, "email_updated")

	return &empty.Empty{}, nil
}

func (s *accountServer) TrackEvent(ctx context.Context, req *pb.TrackEventRequest) (*empty.Empty, error) {
	if s.config.Debug {
		s.logger.Debugf("intercom disabled in dev environment")
		return &empty.Empty{}, nil
	}
	icEvent := intercom.Event{
		UserID:    req.Uuid,
		EventName: "v2_" + req.Event,
		CreatedAt: int64(time.Now().Unix()),
	}

	ic := s.config.GetIntercomClient()
	if err := ic.Events.Save(&icEvent); err != nil {
		return nil, s.internalError(err, "failed to update intercom")
	}
	s.logger.Debugf("updated intercom")
	return &empty.Empty{}, nil
}

func (s *accountServer) SyncUser(ctx context.Context, req *pb.SyncUserRequest) (*empty.Empty, error) {
	if s.config.Debug {
		s.logger.Debugf("intercom disabled in dev environment")
		return &empty.Empty{}, nil
	}

	// Setup for communication
	md := metadata.New(map[string]string{auth.AuthorizationMetadata: auth.AuthorizationAccountService})
	newCtx := metadata.NewOutgoingContext(context.Background(), md)

	u, err := s.Get(newCtx, &pb.GetAccountRequest{Uuid: req.Uuid})
	if err != nil {
		return nil, s.internalError(err, "could not fetch user")
	}
	if u.Phonenumber == "" && u.Email == "" {
		s.logger.Infof("skipping sync for user %v because no email or phonenumber", u.Uuid)
	}

	companyClient, close, err := company.NewClient()
	if err != nil {
		return nil, s.internalError(err, "could not create company client")
	}
	defer close()

	// use a map to de-dupe
	memberships := make(map[string]*company.Company)

	workerOfList, err := companyClient.GetWorkerOf(newCtx, &company.WorkerOfRequest{UserUuid: u.Uuid})
	isWorker := len(workerOfList.Teams) > 0
	for _, t := range workerOfList.Teams {
		c, err := companyClient.GetCompany(newCtx, &company.GetCompanyRequest{Uuid: t.CompanyUuid})
		if err != nil {
			return nil, s.internalError(err, "could not fetch company from team")
		}
		memberships[c.Uuid] = c
	}

	adminOfList, err := companyClient.GetAdminOf(newCtx, &company.AdminOfRequest{UserUuid: u.Uuid})
	isAdmin := len(adminOfList.Companies) > 0
	for _, c := range adminOfList.Companies {
		memberships[c.Uuid] = &c
	}

	// process companies
	var icCompanyList intercom.CompanyList
	for _, c := range memberships {
		icCompanyList.Companies = append(icCompanyList.Companies, intercom.Company{
			CompanyID: c.Uuid,
			Name:      c.Name,
		})
	}

	icUser := intercom.User{
		UserID:        u.Uuid,
		Email:         u.Email,
		Name:          u.Name,
		SignedUpAt:    int64(u.MemberSince.Unix()),
		Avatar:        &intercom.UserAvatar{ImageURL: u.PhotoUrl},
		UpdatedAt:     int64(time.Now().Unix()),
		LastRequestAt: int64(time.Now().Unix()),
		CustomAttributes: map[string]interface{}{
			"v2":                   true,
			"phonenumber":          u.Phonenumber,
			"confirmed_and_active": u.ConfirmedAndActive,
			"is_worker":            isWorker,
			"is_admin":             isAdmin,
			"is_staffjoy_support":  u.Support,
		},
		Companies: &icCompanyList,
	}

	ic := s.config.GetIntercomClient()
	if _, err := ic.Users.Save(&icUser); err != nil {
		return nil, s.internalError(err, "failed to update intercom")
	}
	s.logger.Debugf("updated intercom")
	return &empty.Empty{}, nil
}
