package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/service/mocks"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
)

type APISuite struct {
	suite.Suite

	orderService *mocks.OrderService

	api *api
}

func (s *APISuite) Ctx() context.Context {
	return context.Background()
}

func (s *APISuite) SetupTest() {
	logger.InitForTest()

	s.orderService = mocks.NewOrderService(s.T())

	s.api = NewAPI(
		s.orderService,
	)
}

func (s *APISuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
