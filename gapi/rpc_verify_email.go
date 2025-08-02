package gapi

import (
	"context"

	db "github.com/mahanth/simplebank/db/sqlc"
	"github.com/mahanth/simplebank/pb"
	"github.com/mahanth/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	verifyEmail, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	})
	if err != nil {
		return nil, err
	}

	response := &pb.VerifyEmailResponse{
		IsVerified: verifyEmail.VerifyEmail.IsUsed,
	}

	return response, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if req.GetEmailId() <= 0 {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "email_id",
			Description: "email_id must be a positive integer",
		})
	}

	if val.ValidateString(req.GetSecretCode(), 32, 128) != nil {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "secret_code",
			Description: "secret_code must be between 32 and 128 characters",
		})
	}

	return violations
}
