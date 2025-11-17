package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/max-kriv0s/go-microservices-edu/assembly/internal/model"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/kafka"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
	eventsV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/events/v1"
)

type service struct {
	orderRecorderProducer kafka.Producer
}

func NewService(orderRecorderProducer kafka.Producer) *service {
	return &service{
		orderRecorderProducer: orderRecorderProducer,
	}
}

func (p *service) ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	msg := &eventsV1.ShipAssembled{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal ShipAssembled", zap.Error(err))
		return err
	}

	err = p.orderRecorderProducer.Send(ctx, []byte(event.EventUUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish ShipAssembled", zap.Error(err))
		return err
	}

	return nil
}
