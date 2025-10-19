package v1

import (
	"github.com/max-kriv0s/go-microservices-edu/order/internal/service"
)

type api struct {
	orderService service.OrderService
}

func NewAPI(orderService service.OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
