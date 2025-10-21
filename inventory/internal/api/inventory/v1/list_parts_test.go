package v1

import (
	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/converter"
	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
	inventoryV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestListPartsSuccess() {
	var (
		part1 = fakePart()
		part2 = fakePart()

		parts = []model.Part{part1, part2}

		filter = &inventoryV1.PartsFilter{
			Uuids: []string{part1.Uuid, part2.Uuid},
		}

		req = &inventoryV1.ListPartsRequest{Filter: filter}

		expectedFilter = converter.PartsFilterToModel(filter)

		expectedProtoPart1 = converter.PartToProto(part1)
		expectedProtoPart2 = converter.PartToProto(part2)
		expectedProtoParts = []*inventoryV1.Part{expectedProtoPart1, expectedProtoPart2}
	)

	s.inventoryService.On("ListParts", s.Ctx(), expectedFilter).Return(parts, nil)

	res, err := s.api.ListParts(s.Ctx(), req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedProtoParts, res.GetParts())
}

func (s *APISuite) TestListPartsInternalError() {
	var (
		serviceErr = gofakeit.Error()
		part1      = fakePart()
		part2      = fakePart()

		filter = &inventoryV1.PartsFilter{
			Uuids: []string{part1.Uuid, part2.Uuid},
		}

		req            = &inventoryV1.ListPartsRequest{Filter: filter}
		expectedFilter = converter.PartsFilterToModel(filter)
	)
	s.inventoryService.On("ListParts", s.Ctx(), expectedFilter).Return(nil, serviceErr)

	res, err := s.api.ListParts(s.Ctx(), req)

	s.Require().Nil(res)
	s.Require().Error(err)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.Internal, st.Code())
}
