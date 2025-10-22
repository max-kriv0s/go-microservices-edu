package inventory

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
)

func (s *ServiceSuite) TestListPartsReturnAllParts() {
	var (
		part1 = fakePart()
		part2 = fakePart()

		parts = []model.Part{part1, part2}

		filter = (*model.PartsFilter)(nil)
	)

	s.inventoryRepository.On("FindAll", s.Ctx(), filter).Return(parts, nil)

	res, err := s.service.ListParts(s.Ctx(), nil)

	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(parts, res)
}

func (s *ServiceSuite) TestListPartsSuccess() {
	var (
		part1 = fakePart()

		parts  = []model.Part{part1}
		filter = &model.PartsFilter{
			Uuids: []string{part1.Uuid},
		}

		expectedParts = []model.Part{part1}
	)

	s.inventoryRepository.On("FindAll", s.Ctx(), filter).Return(parts, nil)

	res, err := s.service.ListParts(s.Ctx(), filter)

	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res, 1)
	s.Require().Equal(expectedParts, res)
}

func (s *ServiceSuite) TestListPartsError() {
	var (
		repoErr     = gofakeit.Error()
		expectedErr = model.ErrInternalServer
		filter      = (*model.PartsFilter)(nil)
	)

	s.inventoryRepository.On("FindAll", s.Ctx(), filter).Return(nil, repoErr)

	res, err := s.service.ListParts(s.Ctx(), nil)
	s.Require().Nil(res)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}
