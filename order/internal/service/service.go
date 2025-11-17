package service

import (
	"context"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, data model.CreateOrderRequest) (model.Order, error)
	PayOrder(ctx context.Context, orderUUID string, paymenthMetod model.PaymentMethod) (string, error)
	GetOrder(ctx context.Context, orderUUID string) (model.Order, error)
	CancelOrder(ctx context.Context, orderUUID string) error
}

type OrderProducerService interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error
}
