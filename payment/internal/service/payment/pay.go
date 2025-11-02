package payment

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/max-kriv0s/go-microservices-edu/payment/internal/model"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
)

func (s *service) PayOrder(ctx context.Context, dto model.PayOrder) (string, error) {
	newUUID := uuid.NewString()

	logger.Info(ctx, "Оплата прошла успешно", zap.String("transaction_uuid", newUUID))

	return newUUID, nil
}
