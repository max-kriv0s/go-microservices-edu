package order

import (
	client "github.com/max-kriv0s/go-microservices-edu/order/internal/client/grpc"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/repository"
	def "github.com/max-kriv0s/go-microservices-edu/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	inventoryClient client.InventoryServiceClient
	paymentClient   client.PaymentServiceClient
	orderRepository repository.OrderRepository
}

func NewService(inventoryServiceClient client.InventoryServiceClient, paymentServiceClient client.PaymentServiceClient, orderRepository repository.OrderRepository) *service {
	return &service{
		inventoryClient: inventoryServiceClient,
		paymentClient:   paymentServiceClient,
		orderRepository: orderRepository,
	}
}
