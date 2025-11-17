package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/max-kriv0s/go-microservices-edu/assembly/internal/converter/kafka"
	def "github.com/max-kriv0s/go-microservices-edu/assembly/internal/service"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/kafka"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderPaidRecordedConsumer kafka.Consumer
	orderPaidRecordedDecoder  kafkaConverter.OrderPaidDecoder
	orderProducerService      def.OrderProducerService
}

func NewService(orderPaidRecordedConsumer kafka.Consumer, orderPaidRecordedDecoder kafkaConverter.OrderPaidDecoder, orderProducerService def.OrderProducerService) *service {
	return &service{
		orderPaidRecordedConsumer: orderPaidRecordedConsumer,
		orderPaidRecordedDecoder:  orderPaidRecordedDecoder,
		orderProducerService:      orderProducerService,
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
