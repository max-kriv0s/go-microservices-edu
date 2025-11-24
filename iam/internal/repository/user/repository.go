package user

import (
	"github.com/jackc/pgx/v5/pgxpool"
	def "github.com/max-kriv0s/go-microservices-edu/iam/internal/repository"
)

var _ def.UserRepository = (*repository)(nil)

const (
	usersTable = "users"

	uUuidColumn               = "id"
	uLoginColumn              = "login"
	uEmailColumn              = "email"
	uNotificationMethodColumn = "notification_method"
	uPasswordHash             = "password_hash"
	uCreatedAt                = "created_at"
	uUpdatedAt                = "updated_at"
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
