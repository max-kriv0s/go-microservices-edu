package session

import "context"

func (r *repository) addSessionToUserSet(ctx context.Context, userUUID string, sessionUUID string) error {
	cacheKey := r.getCacheSetKey(userUUID)
	return r.cache.SAdd(ctx, cacheKey, sessionUUID)
}
