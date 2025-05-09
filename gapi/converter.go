package gapi

import (
	db "github.com/guncv/Simple-Bank/db/sqlc"
	"github.com/guncv/Simple-Bank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.Users) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangeAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
