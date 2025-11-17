package order_consumer

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/max-kriv0s/go-microservices-edu/assembly/internal/model"
	kafkaConsumer "github.com/max-kriv0s/go-microservices-edu/platform/pkg/kafka/consumer"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
)

func (s *service) OrderPaidHandler(ctx context.Context, msg kafkaConsumer.Message) error {
	event, err := s.orderPaidRecordedDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode AssemblyRecorded", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.String("payment_method", event.PaymentMethod),
		zap.String("transaction_uuid", event.TransactionUUID),
	)

	go func(event model.OrderPaidEvent) {
		n, err := rand.Int(rand.Reader, big.NewInt(10)) // 0..9
		if err != nil {
			logger.Error(ctx, "Failed to generate random delay", zap.Error(err))
			return
		}
		delay := time.Duration(n.Int64()+1) * time.Second // 1..10 секунд

		select {
		case <-time.After(delay):
			// Таймер сработал, продолжаем
		case <-ctx.Done():
			// Контекст отменён, выходим
			logger.Info(ctx, "Order processing cancelled")
			return
		}

		shipAssemblyEvent := model.ShipAssembledEvent{
			EventUUID:    uuid.NewString(),
			OrderUUID:    event.OrderUUID,
			UserUUID:     event.UserUUID,
			BuildTimeSec: int64(delay.Seconds()),
		}

		err = s.orderProducerService.ProduceShipAssembled(ctx, shipAssemblyEvent)
		if err != nil {
			logger.Error(ctx, "Failed to produce ShipAssembled", zap.Error(err))
		}
	}(event)

	return nil
}
