package inventory

import (
	"context"
	"log"

	"github.com/max-kriv0s/go-microservices-edu/inventory/internal/repository/seed"
)

func (r *repository) seed(ctx context.Context, count int) {
	for i := 0; i < count; i++ {
		part := seed.GeneratePart()
		err := r.Create(ctx, part)
		if err != nil {
			log.Println("Seed error")
			break
		}
	}
}
