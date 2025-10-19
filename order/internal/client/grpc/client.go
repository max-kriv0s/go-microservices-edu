package client

import (
	"context"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
)

type InventoryServiceClient interface {
	ListParts(ctx context.Context, partsUUIDs []string) ([]model.Part, error)
}

type PaymentServiceClient interface {
	PayOrder(ctx context.Context, order model.Order, paymentMethod model.PaymentMethod) (string, error)
}
