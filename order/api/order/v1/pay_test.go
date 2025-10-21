package v1

import (
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/max-kriv0s/go-microservices-edu/order/converter"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	orderV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestPayOrderSuccess() {
	var (
		orderUUID = uuid.New()
		req       = &orderV1.PayOrderRequestDto{
			PaymentMethod: orderV1.PaymentMethodCARD,
		}
		params = orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		}

		expectedPaymentMethod = converter.ApiPaymentMethodToPaymentMethod(req.PaymentMethod)

		transactionUUID = gofakeit.UUID()
		expectedRes     = &orderV1.PayOrderResponseDto{TransactionUUID: transactionUUID}
	)

	s.orderService.On("PayOrder", s.Ctx(), orderUUID.String(), expectedPaymentMethod).Return(transactionUUID, nil)

	res, err := s.api.PayOrder(s.Ctx(), req, params)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedRes, res)
}

func (s *APISuite) TestPayOrderNotFoundError() {
	var (
		serviceErr   = model.ErrOrderNotFound
		expectedCode = http.StatusNotFound

		orderUUID = uuid.New()
		req       = &orderV1.PayOrderRequestDto{
			PaymentMethod: orderV1.PaymentMethodCARD,
		}
		params = orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		}

		expectedPaymentMethod = converter.ApiPaymentMethodToPaymentMethod(req.PaymentMethod)
	)

	s.orderService.On("PayOrder", s.Ctx(), orderUUID.String(), expectedPaymentMethod).Return("", serviceErr)

	res, err := s.api.PayOrder(s.Ctx(), req, params)
	s.Require().NoError(err)
	s.Require().NotEmpty(res)

	s.Require().IsType(&orderV1.NotFoundError{}, res)

	resErr := res.(*orderV1.NotFoundError)
	s.Require().Equal(expectedCode, resErr.Code)
}

func (s *APISuite) TestPayOrderConflictError() {
	var (
		serviceErr   = model.ErrConflict
		expectedCode = http.StatusConflict

		orderUUID = uuid.New()
		req       = &orderV1.PayOrderRequestDto{
			PaymentMethod: orderV1.PaymentMethodCARD,
		}
		params = orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		}

		expectedPaymentMethod = converter.ApiPaymentMethodToPaymentMethod(req.PaymentMethod)
	)

	s.orderService.On("PayOrder", s.Ctx(), orderUUID.String(), expectedPaymentMethod).Return("", serviceErr)

	res, err := s.api.PayOrder(s.Ctx(), req, params)
	s.Require().NoError(err)
	s.Require().NotEmpty(res)

	s.Require().IsType(&orderV1.ConflictError{}, res)

	resErr := res.(*orderV1.ConflictError)
	s.Require().Equal(expectedCode, resErr.Code)
}

func (s *APISuite) TestPayOrderInternalError() {
	var (
		serviceErr   = gofakeit.Error()
		expectedCode = http.StatusInternalServerError

		orderUUID = uuid.New()
		req       = &orderV1.PayOrderRequestDto{
			PaymentMethod: orderV1.PaymentMethodCARD,
		}
		params = orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		}

		expectedPaymentMethod = converter.ApiPaymentMethodToPaymentMethod(req.PaymentMethod)
	)

	s.orderService.On("PayOrder", s.Ctx(), orderUUID.String(), expectedPaymentMethod).Return("", serviceErr)

	res, err := s.api.PayOrder(s.Ctx(), req, params)
	s.Require().NoError(err)
	s.Require().NotEmpty(res)

	s.Require().IsType(&orderV1.InternalServerError{}, res)

	resErr := res.(*orderV1.InternalServerError)
	s.Require().Equal(expectedCode, resErr.Code)
}
