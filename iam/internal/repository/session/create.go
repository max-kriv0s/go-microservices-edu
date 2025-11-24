package session

import (
	"context"
	"time"

	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
	repoConverter "github.com/max-kriv0s/go-microservices-edu/iam/internal/repository/converter"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
	"go.uber.org/zap"
)

func (r *repository) Create(ctx context.Context, uuid string, sessionInfo model.SessionInfo, ttl time.Duration) error {
	cacheKey := r.getCacheKey(uuid)

	redisView := repoConverter.SessionInfoToRedisView(sessionInfo)

	err := r.cache.HashSet(ctx, cacheKey, redisView)
	if err != nil {
		return err
	}

	err = r.addSessionToUserSet(ctx, sessionInfo.User.UUID, uuid)
	if err != nil {
		logger.Warn(ctx, "internal error add session to user set", zap.String("func", "Create"), zap.Error(err))
	}

	return r.cache.Expire(ctx, cacheKey, ttl)
}
