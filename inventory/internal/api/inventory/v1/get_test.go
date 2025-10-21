package v1

import (
	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/converter"
	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestGetPartSuccess() {
	var (
		uuid = gofakeit.UUID()

		part = fakePart()

		req = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}

		expectedProtoPart = converter.PartToProto(part)
	)

	s.inventoryService.On("GetPart", s.Ctx(), uuid).Return(part, nil)

	res, err := s.api.GetPart(s.Ctx(), req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedProtoPart, res.GetPart())
}

func (s *APISuite) TestGetPartNotFoundError() {
	var (
		serviceErr = model.ErrPartNotFound
		uuid       = gofakeit.UUID()

		req = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}
	)

	s.inventoryService.On("GetPart", s.Ctx(), uuid).Return(model.Part{}, serviceErr)

	res, err := s.api.GetPart(s.Ctx(), req)
	s.Require().Nil(res)
	s.Require().Error(err)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestGetPartInternalError() {
	var (
		serviceErr = gofakeit.Error()
		uuid       = gofakeit.UUID()

		req = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}
	)

	s.inventoryService.On("GetPart", s.Ctx(), uuid).Return(model.Part{}, serviceErr)

	res, err := s.api.GetPart(s.Ctx(), req)
	s.Require().Nil(res)
	s.Require().Error(err)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.Internal, st.Code())
}
