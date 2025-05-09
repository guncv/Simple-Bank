package gapi

import (
	"context"
	"database/sql"
	"log"
	"time"

	db "github.com/guncv/Simple-Bank/db/sqlc"
	"github.com/guncv/Simple-Bank/pb"
	"github.com/guncv/Simple-Bank/util"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	log.Printf("update user request: username=%s, email=%s, fullName=%s", req.GetUsername(), req.GetEmail(), req.GetFullName())
	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	passwordChangeAt := sql.NullTime{}
	newPassword := sql.NullString{}
	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			log.Printf("failed to hash password: %s", err)
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}
		newPassword = sql.NullString{
			String: hashedPassword,
			Valid:  req.Password != nil,
		}
		passwordChangeAt = sql.NullTime{
			Time:  time.Now(),
			Valid: req.Password != nil,
		}
	}

	arg := db.UpdateUserParams{
		Username:         req.GetUsername(),
		HashedPassword:   newPassword,
		FullName:         toNullString(req.FullName),
		Email:            toNullString(req.Email),
		PasswordChangeAt: passwordChangeAt,
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		log.Printf("failed to update user: %s", err)
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
	}

	resp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}

	log.Printf("update user response: %v", resp)
	return resp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolations("username", err))
	}
	if req.Password != nil {
		if err := util.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolations("password", err))
		}
	}
	if req.Email != nil {
		if err := util.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolations("email", err))
		}
	}
	if req.FullName != nil {
		if err := util.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, fieldViolations("full_name", err))
		}
	}
	return violations
}

func toNullString(str *string) sql.NullString {
	if str == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *str,
		Valid:  str != nil,
	}
}
