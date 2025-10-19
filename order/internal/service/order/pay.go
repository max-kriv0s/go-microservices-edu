package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (string, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return "", model.ErrOrderNotFound
		}
		return "", model.ErrInternalServer
	}

	if order.Status != model.OrderStatusPendingPayment {
		return "", model.NewConflictError(fmt.Sprintf("You can't pay an order. Order status %s", order.Status))
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, order, paymentMethod)
	if err != nil {
		return "", model.ErrInternalServer
	}

	order.Status = model.OrderStatusPaid
	order.PaymentMethod = &paymentMethod
	order.TransactionUUID = &transactionUUID

	err = s.orderRepository.Update(ctx, order.OrderUUID, order)
	if err != nil {
		return "", model.ErrInternalServer
	}

	return transactionUUID, nil
}
