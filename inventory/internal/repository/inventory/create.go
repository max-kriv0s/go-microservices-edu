package inventory

import (
	"context"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
	repoConverter "github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/converter"
)

func (r *repository) Create(_ context.Context, part model.Part) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[part.Uuid] = repoConverter.PartToRepoModel(part)

	return nil
}
