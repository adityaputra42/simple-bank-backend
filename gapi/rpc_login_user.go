package gapi

import (
	"context"
	"log"
	db "simple-bank/db/sqlc"
	"simple-bank/pb"
	"simple-bank/util"
	"simple-bank/val"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUser(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)

	}
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == db.ErrRecordNotFound {
			status.Errorf(codes.NotFound, "user not found")
		}
		status.Errorf(codes.Internal, "%s", err.Error())

	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		status.Errorf(codes.Unauthenticated, "Unauthorized")
	}
	mtdt := server.extractMetadata(ctx)
	log.Printf("userAgent: %v", mtdt.UserAgent)
	log.Printf("ClientIp: %v", mtdt.ClientIP)
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

func validateLoginUser(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	return violations
}
