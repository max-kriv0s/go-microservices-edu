package model

import "time"

type Login struct {
	Login    string
	Password string
}

type Session struct {
	Uuid      string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

type SessionInfo struct {
	Session Session
	User    User
}
