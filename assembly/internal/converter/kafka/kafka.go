package kafka

import "github.com/max-kriv0s/go-microservices-edu/assembly/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}
