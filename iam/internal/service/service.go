package service

import (
	"context"

	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
)

type AuthService interface {
	Login(ctx context.Context, login, password string) (string, error)
	Whoami(ctx context.Context, sessionUUID string) (model.SessionInfo, error)
}

type UserService interface {
	Register(ctx context.Context, userInfo model.UserInfo, password string) (string, error)
	GetUserById(ctx context.Context, userUUID string) (model.User, error)
	GetUserByLogin(ctx context.Context, login string) (model.User, error)
	VerifyPassword(password, hash string) (bool, error)
}
