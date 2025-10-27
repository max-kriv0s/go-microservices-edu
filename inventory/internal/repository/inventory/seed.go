package inventory

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/seed"
)

func (r *repository) seed(ctx context.Context, count int) {
	countDocum, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("Seed error: %v\n", err)
		return
	}

	// Если коллекция пустая заполним её первоначальными данными
	if countDocum > 0 {
		return
	}

	for i := 0; i < count; i++ {
		part := seed.GeneratePart()
		err := r.Create(ctx, part)
		if err != nil {
			log.Println("Seed error")
			break
		}
	}
}
