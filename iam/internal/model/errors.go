package model

import "errors"

var (
	ErrFailedToHashPassword   = errors.New("failed to hash password")
	ErrUserNotFound           = errors.New("user not found")
	ErrLoginOrEmailExists     = errors.New("login or email exists")
	ErrInternal               = errors.New("unknown error")
	ErrInvalidLoginOrPassword = errors.New("invalid login or password")
	ErrSessionNotFound        = errors.New("session not found")
)
