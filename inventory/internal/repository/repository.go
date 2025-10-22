package repository

import (
	"context"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
)

type InventoryRepository interface {
	Create(ctx context.Context, part model.Part) error
	Get(ctx context.Context, uuid string) (model.Part, error)
	FindAll(ctx context.Context, filter *model.PartsFilter) ([]model.Part, error)
}
