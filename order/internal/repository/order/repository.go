package order

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/max-kriv0s/go-microservices-edu/order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	dbPool *pgxpool.Pool
}

func NewRepository(dbPool *pgxpool.Pool) *repository {
	return &repository{
		dbPool: dbPool,
	}
}
