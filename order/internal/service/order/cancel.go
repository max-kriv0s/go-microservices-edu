package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
)

func (s *service) CancelOrder(ctx context.Context, orderUUID string) error {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return model.ErrOrderNotFound
		}
		return model.ErrInternalServer
	}

	if order.Status != model.OrderStatusPendingPayment {
		return model.NewConflictError(fmt.Sprintf("You can't cancel an order. Order status %s", order.Status))
	}

	updateOrder := model.UpdateOrder{
		Status: lo.ToPtr(model.OrderStatusCancelled),
	}

	err = s.orderRepository.Update(ctx, order.OrderUUID, updateOrder)
	if err != nil {
		return model.ErrInternalServer
	}

	return nil
}
