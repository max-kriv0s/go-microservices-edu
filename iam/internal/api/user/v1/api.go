package v1

import (
	"github.com/max-kriv0s/go-microservices-edu/iam/internal/service"
	userV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/user/v1"
)

type api struct {
	userV1.UnimplementedUserServiceServer

	userService service.UserService
}

func NewApi(userService service.UserService) *api {
	return &api{
		userService: userService,
	}
}
