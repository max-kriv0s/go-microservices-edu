package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/max-kriv0s/go-microservices-edu/notification/internal/converter/kafka"
	def "github.com/max-kriv0s/go-microservices-edu/notification/internal/service"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/kafka"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderPaidRecordedConsumer kafka.Consumer
	orderPaidRecordedDecoder  kafkaConverter.OrderPaidDecoder
	telegramService           def.TelegramService
}

func NewService(orderPaidRecordedConsumer kafka.Consumer, orderPaidRecordedDecoder kafkaConverter.OrderPaidDecoder, telegramService def.TelegramService) *service {
	return &service{
		orderPaidRecordedConsumer: orderPaidRecordedConsumer,
		orderPaidRecordedDecoder:  orderPaidRecordedDecoder,
		telegramService:           telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order assemblyRecordedConsumer service")

	err := s.orderPaidRecordedConsumer.Consume(ctx, s.OrderPaidHandler)
	if err != nil {
		logger.Error(ctx, "Consume from assembly.recorded topic error", zap.Error(err))
		return err
	}

	return nil
}
