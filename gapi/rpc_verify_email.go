package gapi

import (
	"context"
	db "simple-bank/db/sqlc"
	"simple-bank/pb"
	"simple-bank/val"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmail(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	arg := db.VerifyEmailTxParams{
		EmailId:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	}

	result, err := server.store.VerifyEmailTx(ctx, arg)
	if err != nil {

		return nil, status.Errorf(codes.Internal, "Failed to verify email")
	}
	rsp := &pb.VerifyEmailResponse{
		IsVerified: result.User.IsVerifiedEmail,
	}

	return rsp, nil
}

func validateVerifyEmail(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateEmailId(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}
	if err := val.ValidateScretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}

	return violations
}
