package order

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/max-kriv0s/go-microservices-edu/order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

const (
	ordersTable = "orders"

	orderUuidColumn            = "id"
	orderUserUUIDColumn        = "user_id"
	orderTotalPriceColumn      = "total_price"
	orderTransactionUUIDColumn = "transaction_uuid"
	orderPaymentMethodColumn   = "payment_method"
	orderStatusColumn          = "status"

	orderItemTable = "order_items"

	itemOrderUuid      = "order_id"
	itemPartUuidColumn = "part_uuid"
)

type repository struct {
	dbPool *pgxpool.Pool
}

func NewRepository(dbPool *pgxpool.Pool) *repository {
	return &repository{
		dbPool: dbPool,
	}
}

func col(alias, column string) string {
	if alias == "" {
		return column
	}
	return alias + "." + column
}
