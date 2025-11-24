package session

import (
	"fmt"

	def "github.com/max-kriv0s/go-microservices-edu/iam/internal/repository"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/cache"
)

var _ def.SessonRepository = (*repository)(nil)

const (
	cacheKeyPrefix    = "iam:user_session:"
	cacheSetKeyPrefix = "iam:user_session_set:"
)

type repository struct {
	cache cache.RedisClient
}

func NewRepository(cache cache.RedisClient) *repository {
	return &repository{
		cache: cache,
	}
}

func (r *repository) getCacheKey(uuid string) string {
	return fmt.Sprintf("%s%s", cacheKeyPrefix, uuid)
}

func (r *repository) getCacheSetKey(uuid string) string {
	return fmt.Sprintf("%s%s", cacheSetKeyPrefix, uuid)
}
