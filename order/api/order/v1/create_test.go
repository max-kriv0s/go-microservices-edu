package v1

import (
	"net/http"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/max-kriv0s/go-microservices-edu/order/converter"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	orderV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestCreateOrderSuccess() {
	var (
		userUUID   = gofakeit.UUID()
		partsUUIDs = []string{gofakeit.UUID(), gofakeit.UUID()}

		req = &orderV1.CreateOrderRequestDto{
			UserUUID:  userUUID,
			PartUuids: partsUUIDs,
		}

		expectedModel = converter.CreateOrderRequestToModel(req)
		expectedOrder = model.Order{
			OrderUUID:  gofakeit.UUID(),
			UserUUID:   userUUID,
			PartsUUIDs: partsUUIDs,
			TotalPrice: 100,
			Status:     model.OrderStatusPendingPayment,
		}

		expectedResponseDto = converter.OrderToCreateOrderResponseDto(expectedOrder)
	)

	s.orderService.On("CreateOrder", s.Ctx(), expectedModel).Return(expectedOrder, nil)

	res, err := s.api.CreateOrder(s.Ctx(), req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedResponseDto, res)
}

func (s *APISuite) TestCreateOrderBadRequestError() {
	var (
		serviceErr = model.ErrBadRequest
		userUUID   = gofakeit.UUID()
		partsUUIDs = []string{gofakeit.UUID(), gofakeit.UUID()}

		req = &orderV1.CreateOrderRequestDto{
			UserUUID:  userUUID,
			PartUuids: partsUUIDs,
		}

		expectedModel = converter.CreateOrderRequestToModel(req)
		expectedCode  = http.StatusBadRequest
	)

	s.orderService.On("CreateOrder", s.Ctx(), expectedModel).Return(model.Order{}, serviceErr)

	res, err := s.api.CreateOrder(s.Ctx(), req)
	s.Require().NoError(err)
	s.Require().NotEmpty(res)

	s.Require().IsType(&orderV1.BadRequestError{}, res)

	resErr := res.(*orderV1.BadRequestError)
	s.Require().Equal(expectedCode, resErr.Code)
}

func (s *APISuite) TestCreateOrderInternalError() {
	var (
		serviceErr = gofakeit.Error()
		userUUID   = gofakeit.UUID()
		partsUUIDs = []string{gofakeit.UUID(), gofakeit.UUID()}

		req = &orderV1.CreateOrderRequestDto{
			UserUUID:  userUUID,
			PartUuids: partsUUIDs,
		}

		expectedModel = converter.CreateOrderRequestToModel(req)
		expectedCode  = http.StatusInternalServerError
	)

	s.orderService.On("CreateOrder", s.Ctx(), expectedModel).Return(model.Order{}, serviceErr)

	res, err := s.api.CreateOrder(s.Ctx(), req)
	s.Require().NoError(err)
	s.Require().NotEmpty(res)

	s.Require().IsType(&orderV1.InternalServerError{}, res)

	resErr := res.(*orderV1.InternalServerError)
	s.Require().Equal(expectedCode, resErr.Code)
}
