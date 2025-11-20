package order_consumer

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/config"
	kafkaConverter "github.com/max-kriv0s/go-microservices-edu/order/internal/converter/kafka"
	def "github.com/max-kriv0s/go-microservices-edu/order/internal/service"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/kafka"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
	orderService           def.OrderService
}

func NewService(orderAssembledConsumer kafka.Consumer, orderAssembledDecoder kafkaConverter.OrderAssembledDecoder, orderService def.OrderService) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
		orderService:           orderService,
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
