package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/kafka"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
	eventsV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/events/v1"
)

type service struct {
	orderRecordedProducer kafka.Producer
}

func NewService(orderRecordedProducer kafka.Producer) *service {
	return &service{
		orderRecordedProducer: orderRecordedProducer,
	}
}

func (p *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	msg := &eventsV1.OrderPaid{
		EventUuid:       event.EventUUID,
		OrderUuid:       event.OrderUUID,
		UserUuid:        event.UserUUID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUuid: event.TransactionUUID,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal OrderPaid", zap.Error(err))
		return err
	}
	err = p.orderRecordedProducer.Send(ctx, []byte(event.EventUUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish UFORecorded", zap.Error(err))
		return err
	}

	return nil
}
