package order

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
)

func (s *service) GetOrder(ctx context.Context, orderUUID string) (model.Order, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return model.Order{}, model.ErrOrderNotFound
		}

		logger.Error(ctx, "order get error", zap.String("func", "GetOrder"), zap.String("uuid", orderUUID), zap.Error(err))

		return model.Order{}, model.ErrInternalServer
	}
	return order, nil
}
