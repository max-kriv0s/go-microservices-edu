package session

import (
	"context"
	"errors"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
	repoConverter "github.com/max-kriv0s/go-microservices-edu/iam/internal/repository/converter"
	repoModel "github.com/max-kriv0s/go-microservices-edu/iam/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.SessionInfo, error) {
	cacheKey := r.getCacheKey(uuid)

	values, err := r.cache.HGetAll(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return model.SessionInfo{}, model.ErrSessionNotFound
		}
		return model.SessionInfo{}, err
	}

	if len(values) == 0 {
		return model.SessionInfo{}, model.ErrSessionNotFound
	}

	var sessionRedisView repoModel.SessionRedisView
	err = redigo.ScanStruct(values, &sessionRedisView)
	if err != nil {
		return model.SessionInfo{}, err
	}

	return repoConverter.RedisViewToSessionInfo(ctx, sessionRedisView), nil
}
