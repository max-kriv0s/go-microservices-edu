package model

import "time"

type SessionInfo struct {
	Uuid      string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

type SessionRedisView struct {
	SessionUUID        string `redis:"session_uuid"`
	CreatedAtNs        int64  `redis:"created_at"`
	UpdatedAtNs        int64  `redis:"updated_at"`
	ExpiresAtNs        int64  `redis:"expires_at"`
	UserUUID           string `redis:"user_uuid"`
	Login              string `redis:"login"`
	Email              string `redis:"email"`
	NotificationMethod string `redis:"notification_method"`
	UserCreatedAtNs    int64  `redis:"user_created_at"`
	UserUpdatedAtNs    int64  `redis:"user_updated_at"`
}
