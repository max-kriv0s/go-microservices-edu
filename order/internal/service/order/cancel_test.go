package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
)

func (s *ServiceSuite) TestCancelOrderSuccess() {
	var (
		orderUUID  = gofakeit.UUID()
		userUUID   = gofakeit.UUID()
		partsUUIDs = []string{gofakeit.UUID()}

		order = model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   userUUID,
			PartsUUIDs: partsUUIDs,
			TotalPrice: 100,
			Status:     model.OrderStatusPendingPayment,
		}
	)
	s.orderRepository.On("Get", s.Ctx(), orderUUID).Return(order, nil)

	updateOrder := model.UpdateOrder{
		Status: lo.ToPtr(model.OrderStatusCancelled),
	}
	s.orderRepository.On("Update", s.Ctx(), orderUUID, updateOrder).Return(nil)

	err := s.service.CancelOrder(s.Ctx(), orderUUID)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCancelOrderNotFoundErrro() {
	var (
		expectedErr = model.ErrOrderNotFound
		repoErr     = model.ErrOrderNotFound
		orderUUID   = gofakeit.UUID()
	)
	s.orderRepository.On("Get", s.Ctx(), orderUUID).Return(model.Order{}, repoErr)

	err := s.service.CancelOrder(s.Ctx(), orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}

func (s *ServiceSuite) TestCancelOrderInternalErrro() {
	var (
		expectedErr = model.ErrInternalServer
		repoErr     = gofakeit.Error()
		orderUUID   = gofakeit.UUID()
	)
	s.orderRepository.On("Get", s.Ctx(), orderUUID).Return(model.Order{}, repoErr)

	err := s.service.CancelOrder(s.Ctx(), orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}

func (s *ServiceSuite) TestCancelOrderConflictError() {
	var (
		expectedErr = model.ErrConflict

		orderUUID  = gofakeit.UUID()
		userUUID   = gofakeit.UUID()
		partsUUIDs = []string{gofakeit.UUID()}

		order = model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   userUUID,
			PartsUUIDs: partsUUIDs,
			TotalPrice: 100,
			Status:     model.OrderStatusPaid,
		}
	)
	s.orderRepository.On("Get", s.Ctx(), orderUUID).Return(order, nil)

	err := s.service.CancelOrder(s.Ctx(), orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}

func (s *ServiceSuite) TestCancelOrderupdatedError() {
	var (
		expectedErr = model.ErrInternalServer
		repoErr     = gofakeit.Error()

		orderUUID  = gofakeit.UUID()
		userUUID   = gofakeit.UUID()
		partsUUIDs = []string{gofakeit.UUID()}

		order = model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   userUUID,
			PartsUUIDs: partsUUIDs,
			TotalPrice: 100,
			Status:     model.OrderStatusPendingPayment,
		}
	)
	s.orderRepository.On("Get", s.Ctx(), orderUUID).Return(order, nil)

	updateOrder := model.UpdateOrder{
		Status: lo.ToPtr(model.OrderStatusCancelled),
	}
	s.orderRepository.On("Update", s.Ctx(), orderUUID, updateOrder).Return(repoErr)

	err := s.service.CancelOrder(s.Ctx(), orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}
