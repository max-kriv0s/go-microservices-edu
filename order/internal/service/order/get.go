package order

import (
	"context"
	"errors"
	"log"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
)

func (s *service) GetOrder(ctx context.Context, orderUUID string) (model.Order, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return model.Order{}, model.ErrOrderNotFound
		}

		log.Printf("[service.GetOrder] internal error getting order (uuid=%s): %v", orderUUID, err)

		return model.Order{}, model.ErrInternalServer
	}
	return order, nil
}
