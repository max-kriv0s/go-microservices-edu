package inventory

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	inventoryRepository *mocks.InventoryRepository

	service *service
}

func (s *ServiceSuite) Ctx() context.Context {
	return context.Background()
}

func (s *ServiceSuite) SetupTest() {
	s.inventoryRepository = mocks.NewInventoryRepository(s.T())

	s.service = NewService(
		s.inventoryRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func fakePart() model.Part {
	now := time.Now()

	category := model.CategoryEngine

	metadata := make(map[string]*model.Value)

	key := gofakeit.Word()
	metadata[key] = &model.Value{String: lo.ToPtr(gofakeit.Word())}

	tagsCount := gofakeit.Number(1, 5)
	tags := make([]string, tagsCount)
	for i := 0; i < tagsCount; i++ {
		tags[i] = gofakeit.Word()
	}

	part := model.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.Sentence(),
		Price:         100,
		StockQuantity: int64(gofakeit.IntRange(0, 100)),
		Category:      category,
		Dimensions: &model.Dimensions{
			Length: gofakeit.Float64Range(1, 100),
			Width:  gofakeit.Float64Range(1, 100),
			Height: gofakeit.Float64Range(1, 100),
			Weight: gofakeit.Float64Range(0.1, 50),
		},
		Manufacturer: &model.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags:      tags,
		Metadata:  metadata,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	return part
}
