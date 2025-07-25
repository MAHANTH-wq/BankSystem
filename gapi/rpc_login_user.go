package gapi

import (
	"context"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/mahanth/simplebank/db/sqlc"
	"github.com/mahanth/simplebank/pb"
	"github.com/mahanth/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := server.store.GetUser(ctx, req.Username)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user details not found %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to login user %s", err)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)

	if err != nil {

		return nil, status.Errorf(codes.Unauthenticated, "password is not correct %s", err)
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token %s", err)
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token %s", err)
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           pgtype.UUID{Bytes: refreshTokenPayload.ID, Valid: true},
		Username:     refreshTokenPayload.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    pgtype.Timestamptz{Time: refreshTokenPayload.ExpiredAt, Valid: true},
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create new session %s", err)
	}

	response := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
		User:                  convertDBUserToPBUser(user),
	}

	return response, nil
}
