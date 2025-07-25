package gapi

import (
	db "github.com/mahanth/simplebank/db/sqlc"
	"github.com/mahanth/simplebank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertDBUserToPBUser(user db.User) *pb.User {
	return &pb.User{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt.Time),
	}
}
