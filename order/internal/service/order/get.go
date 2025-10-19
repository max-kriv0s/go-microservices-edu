package order

import (
	"context"
	"errors"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
)

func (s *service) GetOrder(ctx context.Context, orderUUID string) (model.Order, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return model.Order{}, model.ErrOrderNotFound
		}
		return model.Order{}, model.ErrInternalServer
	}
	return order, nil
}
