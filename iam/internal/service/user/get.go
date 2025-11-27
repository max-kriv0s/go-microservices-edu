package user

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
)

func (s *service) GetUserById(ctx context.Context, userUUID string) (model.User, error) {
	user, err := s.userRepository.GetUserById(ctx, userUUID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return model.User{}, model.ErrUserNotFound
		}
		logger.Error(ctx, "internal error getting user", zap.String("func", "GetUserById"), zap.String("user_uuid", userUUID), zap.Error(err))

		return model.User{}, model.ErrInternal
	}

	return user, nil
}

func (s *service) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	user, err := s.userRepository.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return model.User{}, model.ErrUserNotFound
		}
		logger.Error(ctx, "internal error getting user", zap.String("func", "GetUserById"), zap.String("login", login), zap.Error(err))

		return model.User{}, model.ErrInternal
	}

	return user, nil
}
