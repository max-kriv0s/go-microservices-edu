package inventory

import (
	"context"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
	repoConverter "github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, part model.Part) error {
	repoPart := repoConverter.PartToRepoModel(part)

	_, err := r.collection.InsertOne(ctx, repoPart)
	if err != nil {
		return err
	}

	return nil
}
