package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/max-kriv0s/go-microservices-edu/iam/internal/converter"
	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
	userV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/user/v1"
)

func (a *api) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	user, err := a.userService.GetUserById(ctx, req.GetUserUuid())
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &userV1.GetUserResponse{User: converter.UserToProto(user)}, nil
}
