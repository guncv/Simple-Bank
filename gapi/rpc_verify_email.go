package gapi

import (
	"context"
	"database/sql"
	"log"

	db "github.com/guncv/Simple-Bank/db/sqlc"
	"github.com/guncv/Simple-Bank/pb"
	"github.com/guncv/Simple-Bank/util"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	log.Printf("verify email request: email_id=%d, secret_code=%s", req.GetEmailId(), req.GetSecretCode())
	violations := validateVerifyEmailRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.VerifyEmailTxParams{
		EmailID:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	}

	result, err := server.store.VerifyEmailTx(ctx, arg)
	if err != nil {
		if sql.ErrNoRows == err {
			return nil, status.Errorf(codes.NotFound, "email not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to verify email: %s", err)
	}

	return &pb.VerifyEmailResponse{
		IsVerified: result.User.IsEmailVerified,
	}, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateEmailId(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolations("email_id", err))
	}
	if err := util.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolations("secret_code", err))
	}
	return violations
}
