package kafka

import "github.com/max-kriv0s/go-microservices-edu/notification/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}

type OrderAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembledEvent, error)
}
