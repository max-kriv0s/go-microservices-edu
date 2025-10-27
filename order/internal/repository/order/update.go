package order

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"

	sq "github.com/Masterminds/squirrel"
	repoConverter "github.com/max-kriv0s/go-microservices-edu/order/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, uuid string, updateOrder model.UpdateOrder) error {
	var foundOrderId string
	err := r.dbPool.QueryRow(ctx, "SELECT id FROM orders WHERE orders.id = $1", uuid).Scan(&foundOrderId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.ErrOrderNotFound
		}
		return err
	}

	tx, err := r.dbPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	builderUpdate := sq.Update("orders").PlaceholderFormat(sq.Dollar).Where(sq.Eq{"id": uuid})

	hasUpdates := false
	if updateOrder.UserUUID != nil {
		builderUpdate = builderUpdate.Set("user_id", *updateOrder.UserUUID)
		hasUpdates = true
	}

	if updateOrder.TotalPrice != nil {
		builderUpdate = builderUpdate.Set("total_price", *updateOrder.TotalPrice)
		hasUpdates = true
	}

	if updateOrder.TransactionUUID != nil {
		builderUpdate = builderUpdate.Set("transaction_uuid", updateOrder.TransactionUUID)
		hasUpdates = true
	}

	if updateOrder.PaymentMethod != nil {
		builderUpdate = builderUpdate.Set("payment_method", repoConverter.PaymentMethodToRepoPaymentMethod(updateOrder.PaymentMethod))
		hasUpdates = true
	}

	if updateOrder.Status != nil {
		builderUpdate = builderUpdate.Set("status", repoConverter.StatusToRepoStatus(*updateOrder.Status))
		hasUpdates = true
	}

	if !hasUpdates {
		return errors.New("no fields to update")
	}

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if updateOrder.PartsUUIDs != nil {
		_, err := tx.Exec(ctx, "DELETE FROM order_items WHERE order_id = $1", uuid)
		if err != nil {
			return err
		}

		if len(*updateOrder.PartsUUIDs) > 0 {
			rows := make([][]interface{}, len(*updateOrder.PartsUUIDs))
			for i, partUUID := range *updateOrder.PartsUUIDs {
				rows[i] = []interface{}{uuid, partUUID}
			}

			_, err := tx.CopyFrom(
				ctx,
				pgx.Identifier{"order_items"},
				[]string{"order_id", "part_uuid"},
				pgx.CopyFromRows(rows),
			)
			if err != nil {
				return err
			}
		}

	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
