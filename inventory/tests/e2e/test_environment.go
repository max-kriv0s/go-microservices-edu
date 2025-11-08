package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ClearSightingsCollection — удаляет все записи из коллекции parts
func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = defaultDatabaseName
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}

func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUUID := gofakeit.UUID()
	now := time.Now()

	partDoc := bson.M{
		"uuid":           partUUID,
		"name":           gofakeit.ProductName(),
		"description":    gofakeit.Sentence(),
		"price":          gofakeit.Price(1, 1000),
		"stock_quantity": int64(gofakeit.IntRange(0, 100)),
		"category":       int32(1), // CategoryEngine
		"dimensions": bson.M{
			"length": gofakeit.Float64Range(1, 100),
			"width":  gofakeit.Float64Range(1, 100),
			"height": gofakeit.Float64Range(1, 100),
			"weight": gofakeit.Float64Range(0.1, 50),
		},
		"manufacturer": bson.M{
			"name":    gofakeit.Company(),
			"country": gofakeit.Country(),
			"website": gofakeit.URL(),
		},
		"tags":      []string{gofakeit.Word()},
		"metadata":  map[string]any{"key": gofakeit.Word()},
		"createdAt": primitive.NewDateTimeFromTime(now),
		"updatedAt": primitive.NewDateTimeFromTime(now),
	}

	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = defaultDatabaseName
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}
