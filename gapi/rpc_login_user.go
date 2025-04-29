package gapi

import (
	"context"
	"database/sql"
	"simple-bank/pb"
	"simple-bank/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			status.Errorf(codes.NotFound, "user not found")
		}
		status.Errorf(codes.Internal, err.Error())

	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		status.Errorf(codes.Unauthenticated, "Unauthorized")
	}

	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)

	if err != nil {
		status.Errorf(codes.Internal, "Failed to create access token")
	}
	rsp := &pb.LoginUserResponse{
		AccessToken: accessToken,
		User:        converUser(user),
	}
	return rsp, nil
}
