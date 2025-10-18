package inventory

import (
	"context"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
	repoConverter "github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/converter"
)

func (r *repository) FindAll(ctx context.Context) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parts := make([]model.Part, 0, len(r.data))
	for _, repoPart := range r.data {
		parts = append(parts, repoConverter.PartToModel(repoPart))
	}
	return parts, nil
}
