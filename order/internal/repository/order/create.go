package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	repoConverter "github.com/max-kriv0s/go-microservices-edu/order/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, order model.Order) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if order.OrderUUID == "" {
		order.OrderUUID = uuid.NewString()
	}
	repoOrder, err := repoConverter.OrderToRepoModel(order)
	if err != nil {
		return "", err
	}

	r.data[order.OrderUUID] = repoOrder

	return order.OrderUUID, nil
}
