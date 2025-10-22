package inventory

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
)

func (r *repository) seed(ctx context.Context, count int) {
	for i := 0; i < count; i++ {
		part := generatePart()
		err := r.Create(ctx, part)
		if err != nil {
			log.Println("Seed error")
			break
		}
	}
}

func generatePart() model.Part {
	now := time.Now()

	// Случайный Category (не UNKNOWN)
	//nolint:gosec // Используем math/rand для некритичных целей, OK
	category := model.Category(rand.Int31n(int32(model.CategoryWing)) + 1)

	// Заполняем map metadata с несколькими парами
	metadata := make(map[string]any)
	for i := 0; i < 3; i++ {
		key := gofakeit.Word()
		metadata[key] = randomValue()
	}

	// Заполняем теги
	tagsCount := gofakeit.Number(1, 5)
	tags := make([]string, tagsCount)
	for i := 0; i < tagsCount; i++ {
		tags[i] = gofakeit.Word()
	}

	part := model.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.Sentence(),
		Price:         gofakeit.Price(1, 1000),
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
		CreatedAt: lo.ToPtr(now.Add(-time.Duration(gofakeit.Number(0, 1000)) * time.Hour)),
		UpdatedAt: lo.ToPtr(now),
	}

	return part
}

func randomValue() any {
	//nolint:gosec // Используем math/rand для некритичных целей, OK
	switch rand.Intn(4) {
	case 0:
		return lo.ToPtr(gofakeit.Word())
	case 1:
		return lo.ToPtr(gofakeit.Int64())
	case 2:
		return lo.ToPtr(gofakeit.Float64())
	default:
		return lo.ToPtr(gofakeit.Bool())
	}
}
