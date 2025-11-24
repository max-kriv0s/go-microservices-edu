package v1

import (
	"context"
	"errors"

	"github.com/max-kriv0s/go-microservices-edu/iam/internal/converter"
	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
	userV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	userUUID, err := a.userService.Register(ctx, converter.RegisterUserRequestToModel(req.GetInfo()), req.Info.GetPassword())
	if err != nil {
		if errors.Is(err, model.ErrLoginOrEmailExists) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &userV1.RegisterResponse{UserUuid: userUUID}, nil
}
