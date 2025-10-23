package inventory

import (
	"context"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
	repoConverter "github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoPart, ok := r.data[uuid]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	return repoConverter.PartToModel(repoPart), nil
}
