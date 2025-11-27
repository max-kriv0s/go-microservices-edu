package repository

import (
	"context"
	"time"

	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (model.User, error)
	GetByLogin(ctx context.Context, login string) (model.User, error)
	GetUserById(ctx context.Context, userUUID string) (model.User, error)
	Create(ctx context.Context, userInfo model.UserInfo, passwordHash string) (string, error)
}

type SessonRepository interface {
	Get(ctx context.Context, uuid string) (model.SessionInfo, error)
	Create(ctx context.Context, uuid string, sessionInfo model.SessionInfo, ttl time.Duration) error
}
