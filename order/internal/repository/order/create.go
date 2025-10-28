package order

import (
	"context"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
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
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("tx rollback failed: %v", err)
		}
	}()

	builderInsert := sq.Insert(ordersTable).PlaceholderFormat(sq.Dollar).Columns(orderUserUUIDColumn, orderTotalPriceColumn, orderStatusColumn).Values(repoOrder.UserUUID, repoOrder.TotalPrice, repoOrder.Status).Suffix("RETURNING " + orderUuidColumn)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return "", nil
	}

	var orderId string
	err = r.dbPool.QueryRow(ctx, query, args...).Scan(&orderId)
	if err != nil {
		return "", nil
	}

	if len(repoOrder.PartsUUIDs) > 0 {
		rows := make([][]interface{}, len(repoOrder.PartsUUIDs))
		for i, partUUID := range repoOrder.PartsUUIDs {
			rows[i] = []interface{}{orderId, partUUID}
		}

		_, err := tx.CopyFrom(
			ctx,
			pgx.Identifier{orderItemTable},
			[]string{itemOrderUuid, itemPartUuidColumn},
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
