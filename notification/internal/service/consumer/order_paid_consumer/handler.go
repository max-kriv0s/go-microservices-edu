package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

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

	// Отправляем уведомление в Telegram
	if err := s.telegramService.SendOrderPaidNotification(ctx, event); err != nil {
		// Логируем ошибку, но не прерываем выполнение
		logger.Error(ctx, "Failed to send telegram notification for OrderPaid", zap.String("event_uuid", event.EventUUID), zap.String("order_uuid", event.OrderUUID), zap.Error(err))
	}

	return nil
}
