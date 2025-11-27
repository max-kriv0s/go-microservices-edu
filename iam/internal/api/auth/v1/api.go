package v1

import (
	"github.com/max-kriv0s/go-microservices-edu/iam/internal/service"
	authV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/auth/v1"
)

type api struct {
	authV1.UnimplementedAuthServiceServer

	authService service.AuthService
}

func NewApi(authService service.AuthService) *api {
	return &api{
		authService: authService,
	}
}
