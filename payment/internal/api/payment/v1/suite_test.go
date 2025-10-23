package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/max-kriv0s/go-microservices-edu/payment/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	paymentService *mocks.PaymentService

	api *api
}

func (s *APISuite) Ctx() context.Context {
	return context.Background()
}

func (s *APISuite) SetupTest() {
	s.paymentService = mocks.NewPaymentService(s.T())

	s.api = NewAPI(
		s.paymentService,
	)
}

func (s *APISuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
