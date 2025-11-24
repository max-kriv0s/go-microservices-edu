package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) Login(ctx context.Context, login string, password string) (string, error) {
	user, err := s.userService.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return "", model.ErrInvalidLoginOrPassword
		}
		logger.Error(ctx, "internal error login", zap.String("func", "Login"), zap.String("login", login), zap.Error(err))

		return "", err
	}

	result, err := s.userService.VerifyPassword(password, user.PasswordHash)
	if err != nil || !result {
		return "", model.ErrInvalidLoginOrPassword
	}

	sessionUUID := uuid.NewString()
	now := time.Now()

	sessionInfo := model.SessionInfo{
		Session: model.Session{
			Uuid:      sessionUUID,
			CreatedAt: now,
			UpdatedAt: now,
			ExpiresAt: now.Add(s.sessionTTL),
		},
		User: user,
	}

	err = s.sessionRepository.Create(ctx, sessionUUID, sessionInfo, s.sessionTTL)
	if err != nil {
		logger.Error(ctx, "set cache error", zap.String("func", "Login"), zap.String("login", login), zap.Error(err))
		return "", model.ErrInternal
	}

	return sessionUUID, nil
}
