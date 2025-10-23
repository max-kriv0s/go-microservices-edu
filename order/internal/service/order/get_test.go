package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
)

func (s *ServiceSuite) TestGetOrderSuccess() {
	var (
		orderUUID     = gofakeit.UUID()
		expectedOrder = model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   gofakeit.UUID(),
			PartsUUIDs: []string{gofakeit.UUID()},
			TotalPrice: 100,
			Status:     model.OrderStatusPendingPayment,
		}
	)

	s.orderRepository.On("Get", s.Ctx(), orderUUID).Return(expectedOrder, nil)

	res, err := s.service.GetOrder(s.Ctx(), orderUUID)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedOrder, res)
}

func (s *ServiceSuite) TestGetOrderNotFoundError() {
	var (
		repoErr     = model.ErrOrderNotFound
		expectedErr = model.ErrOrderNotFound

		orderUUID = gofakeit.UUID()
	)

	s.orderRepository.On("Get", s.Ctx(), orderUUID).Return(model.Order{}, repoErr)

	res, err := s.service.GetOrder(s.Ctx(), orderUUID)
	s.Require().Empty(res)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}

func (s *ServiceSuite) TestGetOrderInternalError() {
	var (
		repoErr     = gofakeit.Error()
		expectedErr = model.ErrInternalServer

		orderUUID = gofakeit.UUID()
	)

	s.orderRepository.On("Get", s.Ctx(), orderUUID).Return(model.Order{}, repoErr)

	res, err := s.service.GetOrder(s.Ctx(), orderUUID)
	s.Require().Empty(res)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}
