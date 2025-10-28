package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	repoConverter "github.com/max-kriv0s/go-microservices-edu/order/internal/repository/converter"
	repoModel "github.com/max-kriv0s/go-microservices-edu/order/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, orderUUID string) (model.Order, error) {
	const (
		ordersAlias     = "o"
		orderItemsAlias = "oi"
	)
	var (
		oOrderUuidColumn = col(ordersAlias, orderUuidColumn)
		joinOrderItems   = fmt.Sprintf(
			`%s  %s ON %s = %s`,
			orderItemTable, orderItemsAlias,
			col(orderItemsAlias, itemOrderUuid),
			oOrderUuidColumn)
	)

	builderSelect := sq.Select(
		oOrderUuidColumn,
		col(ordersAlias, orderUserUUIDColumn),
		col(ordersAlias, orderTotalPriceColumn),
		col(ordersAlias, orderTransactionUUIDColumn),
		col(ordersAlias, orderPaymentMethodColumn),
		col(ordersAlias, orderStatusColumn),
		col(orderItemsAlias, itemPartUuidColumn),
	).
		From(ordersTable + " " + ordersAlias).
		LeftJoin(joinOrderItems).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{oOrderUuidColumn: orderUUID})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return model.Order{}, err
	}

	rows, err := r.dbPool.Query(ctx, query, args...)
	if err != nil {
		return model.Order{}, err
	}
	defer rows.Close()

	var repoOrder repoModel.Order
	parts := []string{}

	for rows.Next() {
		var partUUID sql.NullString

		err := rows.Scan(&repoOrder.OrderUUID, &repoOrder.UserUUID, &repoOrder.TotalPrice, &repoOrder.TransactionUUID, &repoOrder.PaymentMethod, &repoOrder.Status, &partUUID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return model.Order{}, model.ErrOrderNotFound
			}
			return model.Order{}, err
		}

		if partUUID.Valid {
			parts = append(parts, partUUID.String)
		}
	}

	repoOrder.PartsUUIDs = parts

	return repoConverter.OrderToModel(repoOrder), nil
}
