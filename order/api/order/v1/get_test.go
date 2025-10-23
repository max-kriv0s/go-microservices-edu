package v1

import (
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/max-kriv0s/go-microservices-edu/order/converter"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	orderV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestAPIV1OrdersOrderUUIDGetSuccess() {
	var (
		orderUUID  = uuid.New()
		userUUID   = gofakeit.UUID()
		partsUUIDs = []string{gofakeit.UUID()}

		param = orderV1.APIV1OrdersOrderUUIDGetParams{
			OrderUUID: orderUUID,
		}

		order = model.Order{
			OrderUUID:  orderUUID.String(),
			UserUUID:   userUUID,
			PartsUUIDs: partsUUIDs,
			TotalPrice: 100,
			Status:     model.OrderStatusPendingPayment,
		}

		expecedOrder, _ = converter.OrderToGetResponseDto(order)
	)

	s.orderService.On("GetOrder", s.Ctx(), orderUUID.String()).Return(order, nil)

	res, err := s.api.APIV1OrdersOrderUUIDGet(s.Ctx(), param)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expecedOrder, res)
}

func (s *APISuite) TestAPIV1OrdersOrderUUIDGetNotFoundError() {
	var (
		serviceErr   = model.ErrOrderNotFound
		expectedCode = http.StatusNotFound

		orderUUID = uuid.New()
		param     = orderV1.APIV1OrdersOrderUUIDGetParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("GetOrder", s.Ctx(), orderUUID.String()).Return(model.Order{}, serviceErr)

	res, err := s.api.APIV1OrdersOrderUUIDGet(s.Ctx(), param)
	s.Require().NoError(err)
	s.Require().NotEmpty(res)

	s.Require().IsType(&orderV1.NotFoundError{}, res)

	resErr := res.(*orderV1.NotFoundError)
	s.Require().Equal(expectedCode, resErr.Code)
}

func (s *APISuite) TestAPIV1OrdersOrderUUIDGetInternalError() {
	var (
		serviceErr   = gofakeit.Error()
		expectedCode = http.StatusInternalServerError

		orderUUID = uuid.New()
		param     = orderV1.APIV1OrdersOrderUUIDGetParams{
			OrderUUID: orderUUID,
		}
	)

	s.orderService.On("GetOrder", s.Ctx(), orderUUID.String()).Return(model.Order{}, serviceErr)

	res, err := s.api.APIV1OrdersOrderUUIDGet(s.Ctx(), param)
	s.Require().NoError(err)
	s.Require().NotEmpty(res)

	s.Require().IsType(&orderV1.InternalServerError{}, res)

	resErr := res.(*orderV1.InternalServerError)
	s.Require().Equal(expectedCode, resErr.Code)
}
