package model

import "errors"

var (
	ErrPartNotFound   = errors.New("part not found")
	ErrInternalServer = errors.New("unknown error")
)
