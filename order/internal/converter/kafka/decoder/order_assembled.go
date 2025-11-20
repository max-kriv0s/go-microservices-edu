package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	eventsV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewOrderAssembledDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.ShipAssembledEvent, error) {
	var pb eventsV1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.ShipAssembledEvent{
		EventUUID:    pb.EventUuid,
		OrderUUID:    pb.OrderUuid,
		UserUUID:     pb.UserUuid,
		BuildTimeSec: pb.BuildTimeSec,
	}, nil
}
