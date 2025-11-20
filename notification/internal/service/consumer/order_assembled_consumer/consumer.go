package order_assembled_consumer

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/max-kriv0s/go-microservices-edu/notification/internal/config"
	kafkaConverter "github.com/max-kriv0s/go-microservices-edu/notification/internal/converter/kafka"
	def "github.com/max-kriv0s/go-microservices-edu/notification/internal/service"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/kafka"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
	telegramService        def.TelegramService
}

func NewService(orderAssembledConsumer kafka.Consumer, orderAssembledDecoder kafkaConverter.OrderAssembledDecoder, telegramService def.TelegramService) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
		telegramService:        telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order AssembledConsumer service")

	topicName := config.AppConfig().OrderAssembledConsumer.Topic()

	err := s.orderAssembledConsumer.Consume(ctx, s.OrderAssembledHandler)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("Consume from %s topic error", topicName), zap.Error(err))
		return err
	}

	return nil
}
