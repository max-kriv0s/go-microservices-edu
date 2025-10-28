package inventory

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/model"
	repoConverter "github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/converter"
	repoModel "github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/model"
)

func (r *repository) ListParts(ctx context.Context, filter *model.PartsFilter) ([]model.Part, error) {
	searchFilter := buildPartsFilter(filter)
	cursor, err := r.collection.Find(ctx, searchFilter)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close cursor: %v\n", cerr)
		}
	}()

	var repoParts []repoModel.Part
	err = cursor.All(ctx, &repoParts)
	if err != nil {
		return nil, err
	}

	parts := make([]model.Part, 0, len(repoParts))
	for _, repoPart := range repoParts {
		part := repoConverter.PartToModel(repoPart)
		parts = append(parts, part)
	}

	return parts, nil
}

func buildPartsFilter(f *model.PartsFilter) bson.M {
	filter := bson.M{}
	if f == nil {
		return filter
	}

	if len(f.Uuids) > 0 {
		filter["uuid"] = bson.M{"$in": f.Uuids}
	}
	if len(f.Names) > 0 {
		filter["name"] = bson.M{"$in": f.Names}
	}
	if len(f.Categories) > 0 {
		filter["category"] = bson.M{"$in": f.Categories}
	}
	if len(f.ManufacturerCountries) > 0 {
		filter["manufacturer.country"] = bson.M{"$in": f.ManufacturerCountries}
	}
	if len(f.Tags) > 0 {
		filter["tags"] = bson.M{"$in": f.Tags}
	}

	return filter
}
