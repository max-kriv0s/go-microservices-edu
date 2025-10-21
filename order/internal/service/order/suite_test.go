package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	clientMocks "github.com/max-kriv0s/go-microservices-edu/order/internal/client/grpc/mocks"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	inventoryServiceClient *clientMocks.InventoryServiceClient
	paymentServiceClient   *clientMocks.PaymentServiceClient
	orderRepository        *mocks.OrderRepository

	service *service
}

func (s *ServiceSuite) Ctx() context.Context {
	return context.Background()
}

func (s *ServiceSuite) SetupTest() {
	s.inventoryServiceClient = clientMocks.NewInventoryServiceClient(s.T())
	s.paymentServiceClient = clientMocks.NewPaymentServiceClient(s.T())
	s.orderRepository = mocks.NewOrderRepository(s.T())

	s.service = NewService(
		s.inventoryServiceClient, s.paymentServiceClient, s.orderRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
