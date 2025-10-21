package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
)

func (s *ServiceSuite) TestPayOrderSuccess() {
	var (
		orderUUID     = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCard

		order = model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   gofakeit.UUID(),
			PartsUUIDs: []string{gofakeit.UUID()},
			TotalPrice: 100,
			Status:     model.OrderStatusPendingPayment,
		}

		transactionUUID = gofakeit.UUID()
	)

	s.orderRepository.On("Get", s.Ctx(), orderUUID).Return(order, nil)
	s.paymentServiceClient.On("PayOrder", s.Ctx(), order, paymentMethod).Return(transactionUUID, nil)

	paymentOrder := order
	paymentOrder.Status = model.OrderStatusPaid
	paymentOrder.PaymentMethod = &paymentMethod
	paymentOrder.TransactionUUID = &transactionUUID

	s.orderRepository.On("Update", s.Ctx(), orderUUID, paymentOrder).Return(nil)

	res, err := s.service.PayOrder(s.Ctx(), orderUUID, paymentMethod)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(transactionUUID, res)
}

func (s *ServiceSuite) TestPayOrderNotFoundError() {
	var (
		repoErr     = model.ErrOrderNotFound
		expectedErr = model.ErrOrderNotFound

		orderUUID     = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCard
	)

	s.orderRepository.On("Get", s.Ctx(), orderUUID).Return(model.Order{}, repoErr)

	res, err := s.service.PayOrder(s.Ctx(), orderUUID, paymentMethod)
	s.Require().Empty(res)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}

func (s *ServiceSuite) TestPayOrderInternalError() {
	var (
		repoErr     = gofakeit.Error()
		expectedErr = model.ErrInternalServer

		orderUUID     = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCard
	)

	s.orderRepository.On("Get", s.Ctx(), orderUUID).Return(model.Order{}, repoErr)

	res, err := s.service.PayOrder(s.Ctx(), orderUUID, paymentMethod)
	s.Require().Empty(res)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}

func (s *ServiceSuite) TestPayOrderConflictError() {
	var (
		expectedErr = model.ErrConflict

		orderUUID     = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCard

		order = model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   gofakeit.UUID(),
			PartsUUIDs: []string{gofakeit.UUID()},
			TotalPrice: 100,
			Status:     model.OrderStatusPaid,
		}
	)

	s.orderRepository.On("Get", s.Ctx(), orderUUID).Return(order, nil)

	res, err := s.service.PayOrder(s.Ctx(), orderUUID, paymentMethod)
	s.Require().Empty(res)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}

func (s *ServiceSuite) TestPayOrderPaymentServiceError() {
	var (
		expectedErr       = model.ErrInternalServer
		paymentServiceErr = gofakeit.Error()

		orderUUID     = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCard

		order = model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   gofakeit.UUID(),
			PartsUUIDs: []string{gofakeit.UUID()},
			TotalPrice: 100,
			Status:     model.OrderStatusPendingPayment,
		}
	)

	s.orderRepository.On("Get", s.Ctx(), orderUUID).Return(order, nil)
	s.paymentServiceClient.On("PayOrder", s.Ctx(), order, paymentMethod).Return("", paymentServiceErr)

	res, err := s.service.PayOrder(s.Ctx(), orderUUID, paymentMethod)
	s.Require().Empty(res)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}

func (s *ServiceSuite) TestPayOrderUpdatedError() {
	var (
		repoErr     = gofakeit.Error()
		expectedErr = model.ErrInternalServer

		orderUUID     = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCard

		order = model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   gofakeit.UUID(),
			PartsUUIDs: []string{gofakeit.UUID()},
			TotalPrice: 100,
			Status:     model.OrderStatusPendingPayment,
		}

		transactionUUID = gofakeit.UUID()
	)

	s.orderRepository.On("Get", s.Ctx(), orderUUID).Return(order, nil)
	s.paymentServiceClient.On("PayOrder", s.Ctx(), order, paymentMethod).Return(transactionUUID, nil)

	paymentOrder := order
	paymentOrder.Status = model.OrderStatusPaid
	paymentOrder.PaymentMethod = &paymentMethod
	paymentOrder.TransactionUUID = &transactionUUID

	s.orderRepository.On("Update", s.Ctx(), orderUUID, paymentOrder).Return(repoErr)

	res, err := s.service.PayOrder(s.Ctx(), orderUUID, paymentMethod)
	s.Require().Empty(res)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}
