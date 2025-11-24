package user

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
	repoConverter "github.com/max-kriv0s/go-microservices-edu/iam/internal/repository/converter"
	repoModel "github.com/max-kriv0s/go-microservices-edu/iam/internal/repository/model"
)

func (r *repository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	return r.getByColumn(ctx, uEmailColumn, email)
}

func (r *repository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	return r.getByColumn(ctx, uLoginColumn, login)
}

func (r *repository) GetUserById(ctx context.Context, userUUID string) (model.User, error) {
	return r.getByColumn(ctx, uUuidColumn, userUUID)
}

func (r *repository) getByColumn(ctx context.Context, column, value string) (model.User, error) {
	whereColumn := col(usersTable, column)

	builderSelect := sq.Select(
		col(usersTable, uUuidColumn),
		col(usersTable, uLoginColumn),
		col(usersTable, uEmailColumn),
		col(usersTable, uNotificationMethodColumn),
		col(usersTable, uPasswordHash),
		col(usersTable, uCreatedAt),
		col(usersTable, uUpdatedAt),
	).
		From(usersTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{whereColumn: value})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return model.User{}, err
	}

	var repoUser repoModel.User
	err = r.dbPool.QueryRow(ctx, query, args...).Scan(
		&repoUser.UUID,
		&repoUser.Login,
		&repoUser.Email,
		&repoUser.NotificationMethod,
		&repoUser.PasswordHash,
		&repoUser.CreatedAt,
		&repoUser.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, err
	}

	user, err := repoConverter.UserToModel(repoUser)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
