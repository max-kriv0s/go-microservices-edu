package order

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/model"
	repoConverter "github.com/max-kriv0s/go-microservices-edu/order/internal/repository/converter"
	repoModel "github.com/max-kriv0s/go-microservices-edu/order/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, orderUUID string) (model.Order, error) {
	builderSelect := sq.Select(
		"o.id",
		"o.user_id",
		"o.total_price",
		"o.transaction_uuid",
		"o.payment_method",
		"o.status",
		"oi.part_uuid",
	).
		From("orders o").
		LeftJoin("order_items oi ON oi.order_id = o.id").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"o.id": orderUUID})

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
	first := true

	for rows.Next() {
		var (
			id              string
			userId          string
			totalPrice      float64
			transactionUuid sql.NullString
			paymentMethod   sql.NullString
			status          string
			partUUID        sql.NullString
		)

		if err := rows.Scan(&id, &userId, &totalPrice, &transactionUuid, &paymentMethod, &status, &partUUID); err != nil {
			return model.Order{}, err
		}

		if first {
			repoOrder.OrderUUID = id
			repoOrder.UserUUID = userId
			repoOrder.TotalPrice = totalPrice

			if transactionUuid.Valid {
				repoOrder.TransactionUUID = &transactionUuid.String
			}

			if paymentMethod.Valid {
				pm := repoModel.PaymentMethod(paymentMethod.String)
				repoOrder.PaymentMethod = &pm
			}

			repoOrder.Status = repoModel.OrderStatus(status)
			first = false
		}
		if partUUID.Valid {
			parts = append(parts, partUUID.String)
		}

	}

	repoOrder.PartsUUIDs = parts

	if first {
		return model.Order{}, model.ErrOrderNotFound
	}

	return repoConverter.OrderToModel(repoOrder), nil
}
