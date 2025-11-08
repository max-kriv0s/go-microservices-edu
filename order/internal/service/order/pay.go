package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
)

func (s *service) PayOrder(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (string, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return "", model.ErrOrderNotFound
		}

		logger.Error(ctx, "order get error", zap.String("func", "PayOrder"), zap.String("uuid", orderUUID), zap.Error(err))

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
		logger.Error(ctx, "order update error", zap.String("func", "PayOrder"), zap.String("uuid", orderUUID), zap.Error(err))

		return "", model.ErrInternalServer
	}

	return transactionUUID, nil
}
