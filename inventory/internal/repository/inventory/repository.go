package inventory

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	def "github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	collection := db.Collection("parts")

	// Создаем индексы при инициализации
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		panic(err)
	}

	repo := &repository{
		collection: collection,
	}

	repo.seed(context.Background(), 10)
	return repo
}
