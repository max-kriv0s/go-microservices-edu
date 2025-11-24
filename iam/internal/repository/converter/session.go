package converter

import (
	"encoding/json"
	"time"

	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
	repoModel "github.com/max-kriv0s/go-microservices-edu/iam/internal/repository/model"
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

func RedisViewToSessionInfo(redisView repoModel.SessionRedisView) model.SessionInfo {
	notificationMethod := make(map[model.ProviderName]model.Target, 0)
	if redisView.NotificationMethod != "" {
		_ = json.Unmarshal([]byte(redisView.NotificationMethod), &notificationMethod)

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
