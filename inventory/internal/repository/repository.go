package repository

import (
	"context"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
)

type InventoryRepository interface {
	Create(ctx context.Context, part model.Part) error
	Get(ctx context.Context, uuid string) (model.Part, error)
	ListParts(ctx context.Context, filter *model.PartsFilter) ([]model.Part, error)
	Seed(ctx context.Context, count int) error
}
