package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	clientMocks "github.com/max-kriv0s/go-microservices-edu/order/internal/client/grpc/mocks"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/repository/mocks"
	producerMocks "github.com/max-kriv0s/go-microservices-edu/order/internal/service/mocks"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
)

type ServiceSuite struct {
	suite.Suite

	inventoryServiceClient *clientMocks.InventoryServiceClient
	paymentServiceClient   *clientMocks.PaymentServiceClient
	orderRepository        *mocks.OrderRepository
	orderProducerService   *producerMocks.OrderProducerService

	service *service
}

func (s *ServiceSuite) Ctx() context.Context {
	return context.Background()
}

func (s *ServiceSuite) SetupTest() {
	logger.InitForTest()

	s.inventoryServiceClient = clientMocks.NewInventoryServiceClient(s.T())
	s.paymentServiceClient = clientMocks.NewPaymentServiceClient(s.T())
	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.orderProducerService = producerMocks.NewOrderProducerService(s.T())

	s.service = NewService(
		s.inventoryServiceClient, s.paymentServiceClient, s.orderRepository, s.orderProducerService,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
