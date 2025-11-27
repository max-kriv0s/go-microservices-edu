package auth

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
)

func (s *service) Whoami(ctx context.Context, sessionUUID string) (model.SessionInfo, error) {
	sessionInfo, err := s.sessionRepository.Get(ctx, sessionUUID)
	if err != nil {
		if errors.Is(err, model.ErrSessionNotFound) {
			return model.SessionInfo{}, model.ErrSessionNotFound
		}

		logger.Error(ctx, "internal error sessionUUID", zap.String("func", "Whoami"), zap.String("session_uuid", sessionUUID), zap.Error(err))
		return model.SessionInfo{}, model.ErrInternal
	}

	return sessionInfo, nil
}
