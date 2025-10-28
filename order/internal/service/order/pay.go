package order

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/samber/lo"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (string, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return "", model.ErrOrderNotFound
		}

		log.Printf("[service.PayOrder] internal error pay order (uuid=%s): %v", orderUUID, err)

		return "", model.ErrInternalServer
	}

	if order.Status != model.OrderStatusPendingPayment {
		return "", model.NewConflictError(fmt.Sprintf("You can't pay an order. Order status %s", order.Status))
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, order, paymentMethod)
	if err != nil {
		return "", model.ErrInternalServer
	}

	updateOrder := model.UpdateOrder{
		Status:          lo.ToPtr(model.OrderStatusPaid),
		PaymentMethod:   lo.ToPtr(paymentMethod),
		TransactionUUID: lo.ToPtr(transactionUUID),
	}

	err = s.orderRepository.Update(ctx, order.OrderUUID, updateOrder)
	if err != nil {
		log.Printf("[service.PayOrder] internal error update order (uuid=%s): %v", orderUUID, err)
		return "", model.ErrInternalServer
	}

	return transactionUUID, nil
}
