package service

import (
	"context"

	"github.com/max-kriv0s/go-microservices-edu/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, dto model.PayOrder) (string, error)
}
