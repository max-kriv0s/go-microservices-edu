package inventory

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
	repoConverter "github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/converter"
	repoModel "github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, partUUID string) (model.Part, error) {
	var repoPart repoModel.Part
	err := r.collection.FindOne(ctx, bson.M{"uuid": partUUID}).Decode(&repoPart)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, model.ErrPartNotFound
		}
		return model.Part{}, err
	}

	return repoConverter.PartToModel(repoPart), nil
}
