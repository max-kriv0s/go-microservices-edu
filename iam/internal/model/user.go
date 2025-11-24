package model

import "time"

type (
	ProviderName string
	Target       string
)

type UserInfo struct {
	Login              string
	Email              string
	NotificationMethod map[ProviderName]Target
}

type User struct {
	UUID         string
	UserInfo     UserInfo
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
