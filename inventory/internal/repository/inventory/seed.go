package inventory

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/seed"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
)

func (r *repository) Seed(ctx context.Context, count int) error {
	countDocum, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		logger.Error(ctx, "count documents error", zap.String("func", "seed"), zap.String("collection", r.collection.Name()), zap.Error(err))

		return err
	}

	// Если коллекция пустая заполним её первоначальными данными
	if countDocum > 0 {
		return nil
	}

	for i := 0; i < count; i++ {
		part := seed.GeneratePart()
		err := r.Create(ctx, part)
		if err != nil {
			logger.Error(ctx, "create error", zap.String("func", "seed"), zap.String("collection", r.collection.Name()), zap.Error(err))

			return err
		}
	}

	return nil
}
