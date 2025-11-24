package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/max-kriv0s/go-microservices-edu/iam/internal/model"
	repoConverter "github.com/max-kriv0s/go-microservices-edu/iam/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, userInfo model.UserInfo, passwordHash string) (string, error) {
	repoUser, err := repoConverter.UserInfoToRepoModel(userInfo, passwordHash)
	if err != nil {
		return "", err
	}

	builderInsert := sq.
		Insert(usersTable).
		PlaceholderFormat(sq.Dollar).
		Columns(uLoginColumn, uEmailColumn, uNotificationMethodColumn, uPasswordHash).
		Values(repoUser.Login, repoUser.Email, repoUser.NotificationMethod, repoUser.PasswordHash).
		Suffix("RETURNING " + uUuidColumn)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return "", nil
	}

	var userUUID string
	err = r.dbPool.QueryRow(ctx, query, args...).Scan(&userUUID)
	if err != nil {
		return "", nil
	}

	return userUUID, nil
}
