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

func (r *repository) Update(ctx context.Context, orderUUUiD string, updateOrder model.UpdateOrder) error {
	builderFoundOrder := sq.Select(orderUuidColumn).From(ordersTable).PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{orderUuidColumn: orderUUUiD}).
		Limit(1)

	query, args, err := builderFoundOrder.ToSql()
	if err != nil {
		return err
	}

	var foundOrderId string
	err = r.dbPool.QueryRow(ctx, query, args...).Scan(&foundOrderId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ErrOrderNotFound
		}
		return err
	}

	tx, err := r.dbPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("tx rollback failed: %v", err)
		}
	}()

	builderUpdate := sq.Update(ordersTable).PlaceholderFormat(sq.Dollar).Where(sq.Eq{orderUuidColumn: orderUUUiD})

	hasUpdates := false
	if updateOrder.UserUUID != nil {
		builderUpdate = builderUpdate.Set(orderUserUUIDColumn, *updateOrder.UserUUID)
		hasUpdates = true
	}

	if updateOrder.TotalPrice != nil {
		builderUpdate = builderUpdate.Set(orderTotalPriceColumn, *updateOrder.TotalPrice)
		hasUpdates = true
	}

	if updateOrder.TransactionUUID != nil {
		builderUpdate = builderUpdate.Set(orderTransactionUUIDColumn, updateOrder.TransactionUUID)
		hasUpdates = true
	}

	if updateOrder.PaymentMethod != nil {
		builderUpdate = builderUpdate.Set(orderPaymentMethodColumn, repoConverter.PaymentMethodToRepoPaymentMethod(updateOrder.PaymentMethod))
		hasUpdates = true
	}

	if updateOrder.Status != nil {
		builderUpdate = builderUpdate.Set(orderStatusColumn, repoConverter.StatusToRepoStatus(*updateOrder.Status))
		hasUpdates = true
	}

	if !hasUpdates {
		return errors.New("no fields to update")
	}

	query, args, err = builderUpdate.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if updateOrder.PartsUUIDs != nil {
		builderDelete := sq.Delete(ordersTable).
			PlaceholderFormat(sq.Dollar).
			Where(sq.Eq{orderUuidColumn: orderUUUiD})

		query, args, err = builderDelete.ToSql()
		if err != nil {
			return err
		}

		_, err := tx.Exec(ctx, query, args...)
		if err != nil {
			return err
		}

		if len(*updateOrder.PartsUUIDs) > 0 {
			rows := make([][]interface{}, len(*updateOrder.PartsUUIDs))
			for i, partUUID := range *updateOrder.PartsUUIDs {
				rows[i] = []interface{}{orderUUUiD, partUUID}
			}

			_, err := tx.CopyFrom(
				ctx,
				pgx.Identifier{orderItemTable},
				[]string{itemOrderUuid, itemPartUuidColumn},
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
