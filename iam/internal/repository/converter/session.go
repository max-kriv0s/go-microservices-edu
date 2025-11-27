package converter

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"

	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
	repoModel "github.com/max-kriv0s/go-microservices-edu/iam/internal/repository/model"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
)

func SessionInfoToRedisView(sessionInfo model.SessionInfo) repoModel.SessionRedisView {
	var (
		user    = sessionInfo.User
		session = sessionInfo.Session
	)

	var notificationMethod string
	if len(user.UserInfo.NotificationMethod) > 0 {
		data, err := json.Marshal(user.UserInfo.NotificationMethod)
		if err == nil {
			notificationMethod = string(data)
		}
	}

	return repoModel.SessionRedisView{
		SessionUUID:        session.Uuid,
		CreatedAtNs:        session.CreatedAt.UnixNano(),
		UpdatedAtNs:        session.UpdatedAt.UnixNano(),
		ExpiresAtNs:        session.ExpiresAt.UnixNano(),
		UserUUID:           user.UUID,
		Login:              user.UserInfo.Login,
		Email:              user.UserInfo.Email,
		NotificationMethod: notificationMethod,
		UserCreatedAtNs:    user.CreatedAt.UnixNano(),
		UserUpdatedAtNs:    user.CreatedAt.UnixNano(),
	}
}

func RedisViewToSessionInfo(ctx context.Context, redisView repoModel.SessionRedisView) model.SessionInfo {
	notificationMethod := make(map[model.ProviderName]model.Target, 0)
	if redisView.NotificationMethod != "" {
		err := json.Unmarshal([]byte(redisView.NotificationMethod), &notificationMethod)
		if err != nil {
			logger.Warn(ctx, "error unmarshal NotificationMethod", zap.String("func", "RedisViewToSessionInfo"), zap.String("user_uuid", redisView.UserUUID), zap.Error(err))
		}
	}

	return model.SessionInfo{
		Session: model.Session{
			Uuid:      redisView.SessionUUID,
			CreatedAt: time.Unix(0, redisView.CreatedAtNs),
			UpdatedAt: time.Unix(0, redisView.UpdatedAtNs),
			ExpiresAt: time.Unix(0, redisView.ExpiresAtNs),
		},
		User: model.User{
			UUID:      redisView.UserUUID,
			CreatedAt: time.Unix(0, redisView.UserCreatedAtNs),
			UpdatedAt: time.Unix(0, redisView.UserUpdatedAtNs),
			UserInfo: model.UserInfo{
				Login:              redisView.Login,
				Email:              redisView.Email,
				NotificationMethod: notificationMethod,
			},
		},
	}
}
