package model

import "time"

type ProviderName string
type Target string

type User struct {
	UUID               *string
	Login              string
	PasswordHash       string
	Email              string
	NotificationMethod []byte // jsonb как []byte
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
