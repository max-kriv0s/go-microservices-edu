package v1

import (
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	orderV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestAPIV1OrdersOrderUUIDCancelPostSuccess() {
	var (
		uuid  = uuid.New()
		param = orderV1.APIV1OrdersOrderUUIDCancelPostParams{
			OrderUUID: uuid,
		}
	)

	s.orderService.On("CancelOrder", s.Ctx(), uuid.String()).Return(nil)
	res, err := s.api.APIV1OrdersOrderUUIDCancelPost(s.Ctx(), param)
	s.Require().NoError(err)
	s.Require().Empty(res)
}

func (s *APISuite) TestAPIV1OrdersOrderUUIDCancelPostNotFoundError() {
	var (
		serviceErr = model.ErrOrderNotFound
		uuid       = uuid.New()
		param      = orderV1.APIV1OrdersOrderUUIDCancelPostParams{
			OrderUUID: uuid,
		}
		expectedCode = http.StatusNotFound
	)

	s.orderService.On("CancelOrder", s.Ctx(), uuid.String()).Return(serviceErr)
	res, err := s.api.APIV1OrdersOrderUUIDCancelPost(s.Ctx(), param)
	s.Require().NoError(err)
	s.Require().NotEmpty(res)

	s.Require().IsType(&orderV1.NotFoundError{}, res)

	notFoundErr := res.(*orderV1.NotFoundError)
	s.Require().Equal(expectedCode, notFoundErr.Code)
}

func (s *APISuite) TestAPIV1OrdersOrderUUIDCancelPostConflictError() {
	var (
		serviceErr = model.ErrConflict
		uuid       = uuid.New()
		param      = orderV1.APIV1OrdersOrderUUIDCancelPostParams{
			OrderUUID: uuid,
		}
		expectedCode = http.StatusConflict
	)

	s.orderService.On("CancelOrder", s.Ctx(), uuid.String()).Return(serviceErr)
	res, err := s.api.APIV1OrdersOrderUUIDCancelPost(s.Ctx(), param)
	s.Require().NoError(err)
	s.Require().NotEmpty(res)

	s.Require().IsType(&orderV1.ConflictError{}, res)

	conflictErr := res.(*orderV1.ConflictError)
	s.Require().Equal(expectedCode, conflictErr.Code)
}

func (s *APISuite) TestAPIV1OrdersOrderUUIDCancelPostInternalError() {
	var (
		serviceErr = gofakeit.Error()
		uuid       = uuid.New()
		param      = orderV1.APIV1OrdersOrderUUIDCancelPostParams{
			OrderUUID: uuid,
		}
		expectedCode = http.StatusInternalServerError
	)

	s.orderService.On("CancelOrder", s.Ctx(), uuid.String()).Return(serviceErr)
	res, err := s.api.APIV1OrdersOrderUUIDCancelPost(s.Ctx(), param)
	s.Require().NoError(err)
	s.Require().NotEmpty(res)

	s.Require().IsType(&orderV1.InternalServerError{}, res)

	internalErr := res.(*orderV1.InternalServerError)
	s.Require().Equal(expectedCode, internalErr.Code)
}
