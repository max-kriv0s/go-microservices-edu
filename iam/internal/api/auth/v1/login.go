package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
	authV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/auth/v1"
)

func (a *api) Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	sessionUUID, err := a.authService.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		if errors.Is(err, model.ErrInvalidLoginOrPassword) {
			return nil, status.Error(codes.Unauthenticated, model.ErrInvalidLoginOrPassword.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &authV1.LoginResponse{SessionUuid: sessionUUID}, nil
}
