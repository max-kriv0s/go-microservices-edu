package order

import (
	"context"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	repoConverter "github.com/max-kriv0s/go-microservices-edu/order/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoOrder, ok := r.data[uuid]
	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	return repoConverter.OrderToModel(repoOrder), nil
}
