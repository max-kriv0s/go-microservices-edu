package kafka

import "github.com/max-kriv0s/go-microservices-edu/order/internal/model"

type OrderAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembledEvent, error)
}
