package gapi

import (
	"context"
	db "simple-bank/db/sqlc"
	"simple-bank/pb"
	"simple-bank/util"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to hash password")
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "Username already exist %s", arg.Username)
			}
		}
		return nil, status.Errorf(codes.Internal, "Failed to create user")
	}

	rsp := &pb.CreateUserResponse{
		User: converUser(user),
	}

	return rsp, nil
}
