package gapi

import (
	"context"
	"log"

	db "github.com/guncv/Simple-Bank/db/sqlc"
	"github.com/guncv/Simple-Bank/pb"
	"github.com/guncv/Simple-Bank/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Printf("create user request: username=%s, email=%s, fullName=%s", req.GetUsername(), req.GetEmail(), req.GetFullName())
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

	resp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	log.Printf("create user response: %v", resp)
	return resp, nil
}
