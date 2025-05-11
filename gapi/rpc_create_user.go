package gapi

import (
	"context"
	"log"
	"time"

	db "github.com/guncv/Simple-Bank/db/sqlc"
	"github.com/guncv/Simple-Bank/pb"
	"github.com/guncv/Simple-Bank/util"
	"github.com/guncv/Simple-Bank/worker"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Printf("create user request: username=%s, email=%s, fullName=%s", req.GetUsername(), req.GetEmail(), req.GetFullName())
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		log.Printf("failed to hash password: %s", err)
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		log.Printf("failed to create user: %s", err)
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				log.Printf("username already exists: %s", err)
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
			}
		}
		log.Printf("failed to create user: %s", err)
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	// TODO: use db transaction
	taskPayload := &worker.PayloadSendVerifyEmail{
		Username: user.Username,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.Queue(worker.QueueCritical),
		asynq.ProcessIn(3 * time.Second),
	}
	if err := server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...); err != nil {
		log.Printf("failed to distribute task: %s", err)
		return nil, status.Errorf(codes.Internal, "failed to distribute task: %s", err)
	}

	resp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	log.Printf("create user response: %v", resp)
	return resp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolations("username", err))
	}
	if err := util.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolations("password", err))
	}
	if err := util.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolations("email", err))
	}
	if err := util.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolations("full_name", err))
	}
	return violations
}
