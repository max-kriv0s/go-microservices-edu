package order_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/kafka/consumer"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
)

func (s *service) OrderAssembledHandler(ctx context.Context, msg consumer.Message) error {
	event, err := s.orderAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ShipAssembled", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	err = s.orderService.AssembledOrder(ctx, event.OrderUUID)
	if err != nil {
		logger.Error(ctx, "Handler error ShipAssembled", zap.Error(err))
		return err
	}

	return nil
}
