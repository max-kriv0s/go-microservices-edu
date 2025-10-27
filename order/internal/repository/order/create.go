package order

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	repoConverter "github.com/max-kriv0s/go-microservices-edu/order/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, order model.Order) (string, error) {
	repoOrder := repoConverter.OrderToRepoModel(order)

	tx, err := r.dbPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	var orderId string
	err = tx.QueryRow(ctx, "INSERT INTO orders (user_id, total_price, status) VALUES ($1, $2, $3) RETURNING id", repoOrder.UserUUID, repoOrder.TotalPrice, repoOrder.Status).Scan(&orderId)
	if err != nil {
		return "", err
	}

	if len(repoOrder.PartsUUIDs) > 0 {
		rows := make([][]interface{}, len(repoOrder.PartsUUIDs))
		for i, partUUID := range repoOrder.PartsUUIDs {
			rows[i] = []interface{}{orderId, partUUID}
		}

		_, err := tx.CopyFrom(
			ctx,
			pgx.Identifier{"order_items"},
			[]string{"order_id", "part_uuid"},
			pgx.CopyFromRows(rows),
		)
		if err != nil {
			return "", err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", err
	}

	return orderId, nil
}
