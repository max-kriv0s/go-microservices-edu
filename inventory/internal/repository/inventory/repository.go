package inventory

import (
	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	collection := db.Collection("parts")

	return &repository{
		collection: collection,
	}
}
