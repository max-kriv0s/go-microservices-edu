package repository

import (
	"context"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) (string, error)
	Get(ctx context.Context, uuid string) (model.Order, error)
	Update(ctx context.Context, uuid string, updateOrder model.UpdateOrder) error
}
