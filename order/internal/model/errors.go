package model

import (
	"errors"
	"fmt"
)

var (
	ErrUnknownRepoOrderStatus = errors.New("unknown repo order status")
	ErrOrderNotFound          = errors.New("order not found")
	ErrInternalServer         = errors.New("unknown error")
	ErrBadRequest             = errors.New("bad request")
	ErrConflict               = errors.New("conflict")
)

func NewBadRequestError(msg string) error {
	return fmt.Errorf("%w: %s", ErrBadRequest, msg)
}

func NewConflictError(msg string) error {
	return fmt.Errorf("%w: %s", ErrConflict, msg)
}
