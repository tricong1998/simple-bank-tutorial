package gapi

import (
	"context"
	"database/sql"

	db "github.com/Sotatek-CongNguyen/simple-bank-practice/db/sqlc"
	"github.com/Sotatek-CongNguyen/simple-bank-practice/pb"
	"github.com/Sotatek-CongNguyen/simple-bank-practice/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
		return nil, status.Errorf(codes.Internal, "Find user error: %s", err)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Incorrect password")

	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create access token")

	}

	refreshToken, rfPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create refresh token")

	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           rfPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientID:     "",
		IsBlocked:    false,
		ExpiresAt:    rfPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create session")

	}
	rsp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(rfPayload.ExpiredAt),
	}
	return rsp, nil
}
