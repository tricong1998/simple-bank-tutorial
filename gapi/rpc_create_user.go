package gapi

import (
	"context"

	db "github.com/Sotatek-CongNguyen/simple-bank-practice/db/sqlc"
	"github.com/Sotatek-CongNguyen/simple-bank-practice/pb"
	"github.com/Sotatek-CongNguyen/simple-bank-practice/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}
	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		Email:          req.GetEmail(),
		FullName:       req.GetFullName(),
		HashedPassword: hashedPassword,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				{
					return nil, status.Errorf(codes.AlreadyExists, "Username already exists: %s", err)
				}
			}
		}
		return nil, status.Errorf(codes.AlreadyExists, "failed to create user: %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}
