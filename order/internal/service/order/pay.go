package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/tracing"
)

func (s *service) PayOrder(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (string, error) {
	// Создаем спан для вызова Payment сервиса
	ctx, span := tracing.StartSpan(ctx, "order.call_payment",
		trace.WithAttributes(
			attribute.String("order.uuid", orderUUID),
		),
	)
	defer span.End()

	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		span.RecordError(err)

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
		span.RecordError(err)
		return "", model.ErrInternalServer
	}

	updateOrder := model.UpdateOrder{
		Status:          lo.ToPtr(model.OrderStatusPaid),
		PaymentMethod:   lo.ToPtr(paymentMethod),
		TransactionUUID: lo.ToPtr(transactionUUID),
	}

	err = s.orderRepository.Update(ctx, order.OrderUUID, updateOrder)
	if err != nil {
		span.RecordError(err)
		logger.Error(ctx, "order update error", zap.String("func", "PayOrder"), zap.String("uuid", orderUUID), zap.Error(err))

		return "", model.ErrInternalServer
	}

	err = s.orderProducerService.ProduceOrderPaid(ctx, model.OrderPaidEvent{
		EventUUID:       uuid.NewString(),
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PaymentMethod:   string(paymentMethod),
		TransactionUUID: transactionUUID,
	})
	if err != nil {
		span.RecordError(err)
		return "", nil
	}

	return transactionUUID, nil
}
