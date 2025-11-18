package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/max-kriv0s/go-microservices-edu/notification/internal/model"
	eventsV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/events/v1"
)

type orderPaidDecoder struct{}

func NewOrderPaidDecoder() *orderPaidDecoder {
	return &orderPaidDecoder{}
}

func (d *orderPaidDecoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsV1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.OrderPaidEvent{
		EventUUID:       pb.EventUuid,
		OrderUUID:       pb.OrderUuid,
		UserUUID:        pb.UserUuid,
		PaymentMethod:   pb.PaymentMethod,
		TransactionUUID: pb.TransactionUuid,
	}, nil
}
