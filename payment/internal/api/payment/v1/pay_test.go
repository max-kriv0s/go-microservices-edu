package v1

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/max-kriv0s/go-microservices-edu/payment/internal/converter"
	paymentV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *APISuite) TestPayOrderSuccess() {
	var (
		orderUUId     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		paymentMethod = paymentV1.PaymentMethod_CARD

		req = &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUId,
			UserUuid:      userUUID,
			PaymentMethod: paymentMethod,
		}

		expectdModel = converter.PayDtoToModel(req)

		transactionUuid = gofakeit.UUID()
	)

	s.paymentService.On("PayOrder", s.ctx, expectdModel).Return(transactionUuid, nil)

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(transactionUuid, res.GetTransactionUuid())
}

func (s *APISuite) TestPayOrderValidateError() {
	var (
		orderUUId     = "1"
		userUUID      = gofakeit.UUID()
		paymentMethod = paymentV1.PaymentMethod_CARD

		req = &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUId,
			UserUuid:      userUUID,
			PaymentMethod: paymentMethod,
		}
	)

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().Nil(res)
	s.Require().Error(err)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.InvalidArgument, st.Code())
}

func (s *APISuite) TestPayOrderInternalError() {
	var (
		repoErr = gofakeit.Error()

		orderUUId     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		paymentMethod = paymentV1.PaymentMethod_CARD

		req = &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUId,
			UserUuid:      userUUID,
			PaymentMethod: paymentMethod,
		}

		expectdModel = converter.PayDtoToModel(req)
	)

	s.paymentService.On("PayOrder", s.ctx, expectdModel).Return("", repoErr)

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().Nil(res)
	s.Require().Error(err)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.Internal, st.Code())
}
