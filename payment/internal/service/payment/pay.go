package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/max-kriv0s/go-microservices-edu/payment/internal/model"
)

func (s *service) PayOrder(ctx context.Context, dto model.PayOrder) (string, error) {
	newUUID := uuid.NewString()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", newUUID)

	return newUUID, nil
}
