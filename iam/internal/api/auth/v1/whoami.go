package v1

import (
	"context"
	"errors"

	"github.com/max-kriv0s/go-microservices-edu/iam/internal/converter"
	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
	authV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) Whoami(ctx context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {

	sessionInfo, err := a.authService.Whoami(ctx, req.GetSessionUuid())
	if err != nil {
		if errors.Is(err, model.ErrSessionNotFound) {
			return nil, status.Error(codes.Unauthenticated, model.ErrSessionNotFound.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return converter.SessionInfoToWhoamiResponse(sessionInfo), nil
}
