package inventory

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
)

func (s *ServiceSuite) TestGetPartSuccess() {
	expectedPart := fakePart()

	s.inventoryRepository.On("Get", s.Ctx(), expectedPart.Uuid).Return(expectedPart, nil)

	res, err := s.service.GetPart(s.Ctx(), expectedPart.Uuid)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedPart, res)
}

func (s *ServiceSuite) TestGetPartNotFoundError() {
	var (
		repoErr = model.ErrPartNotFound
		uuid    = gofakeit.UUID()
	)

	s.inventoryRepository.On("Get", s.Ctx(), uuid).Return(model.Part{}, repoErr)

	res, err := s.service.GetPart(s.Ctx(), uuid)
	s.Require().Empty(res)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestGetPartInternalError() {
	var (
		repoErr = gofakeit.Error()
		uuid    = gofakeit.UUID()

		expectedErr = model.ErrInternalServer
	)

	s.inventoryRepository.On("Get", s.Ctx(), uuid).Return(model.Part{}, repoErr)

	res, err := s.service.GetPart(s.Ctx(), uuid)
	s.Require().Empty(res)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedErr)
}
