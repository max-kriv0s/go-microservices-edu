package app

import (
	"context"
	"fmt"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	authV1API "github.com/max-kriv0s/go-microservices-edu/iam/internal/api/auth/v1"
	envoyV1API "github.com/max-kriv0s/go-microservices-edu/iam/internal/api/envoy/v1"
	userV1API "github.com/max-kriv0s/go-microservices-edu/iam/internal/api/user/v1"
	"github.com/max-kriv0s/go-microservices-edu/iam/internal/config"
	"github.com/max-kriv0s/go-microservices-edu/iam/internal/repository"
	sessionRepository "github.com/max-kriv0s/go-microservices-edu/iam/internal/repository/session"
	userRepository "github.com/max-kriv0s/go-microservices-edu/iam/internal/repository/user"
	"github.com/max-kriv0s/go-microservices-edu/iam/internal/service"
	authService "github.com/max-kriv0s/go-microservices-edu/iam/internal/service/auth"
	userService "github.com/max-kriv0s/go-microservices-edu/iam/internal/service/user"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/cache"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/cache/redis"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/closer"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/migrator"
	migratorPg "github.com/max-kriv0s/go-microservices-edu/platform/pkg/migrator/pg"
	authV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/auth/v1"
	userV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/user/v1"
)

type diContainer struct {
	authV1API  authV1.AuthServiceServer
	userV1API  userV1.UserServiceServer
	envoyV1API authv3.AuthorizationServer

	authService service.AuthService
	userService service.UserService

	userRepository    repository.UserRepository
	sessionRepository repository.SessonRepository

	postgresDBConn *pgx.Conn
	postgresDBPool *pgxpool.Pool

	migrator migrator.Migrator

	redisPool   *redigo.Pool
	redisClient cache.RedisClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) AuthV1API(ctx context.Context) authV1.AuthServiceServer {
	if d.authV1API == nil {
		d.authV1API = authV1API.NewApi(d.AuthService(ctx))
	}

	return d.authV1API
}

func (d *diContainer) UserV1API(ctx context.Context) userV1.UserServiceServer {
	if d.userV1API == nil {
		d.userV1API = userV1API.NewApi(d.UserService(ctx))
	}

	return d.userV1API
}

func (d *diContainer) EnvoyV1API(ctx context.Context) authv3.AuthorizationServer {
	if d.envoyV1API == nil {
		d.envoyV1API = envoyV1API.NewApi(d.AuthV1API(ctx))
	}

	return d.envoyV1API
}

func (d *diContainer) AuthService(ctx context.Context) service.AuthService {
	if d.authService == nil {
		d.authService = authService.NewService(d.UserService(ctx), d.SessionRepository(), config.AppConfig().Session.SessionTTL())
	}

	return d.authService
}

func (d *diContainer) UserService(ctx context.Context) service.UserService {
	if d.userService == nil {
		d.userService = userService.NewService(d.UserRepository(ctx))
	}

	return d.userService
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userRepository.NewRepository(d.PostgresDBPool(ctx))
	}

	return d.userRepository
}

func (d *diContainer) PostgresDBConn(ctx context.Context) *pgx.Conn {
	if d.postgresDBConn != nil {
		return d.postgresDBConn
	}

	dbURI := config.AppConfig().Postgres.URI()

	dbConn, err := pgx.Connect(ctx, dbURI)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v\n", err.Error()))
	}

	closer.AddNamed("Postgres database connect", func(ctx context.Context) error {
		return dbConn.Close(ctx)
	})

	err = dbConn.Ping(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to ping postgresDB: %v\n", err.Error()))
	}

	d.postgresDBConn = dbConn

	return d.postgresDBConn
}

func (d *diContainer) PostgresDBPool(ctx context.Context) *pgxpool.Pool {
	if d.postgresDBPool != nil {
		return d.postgresDBPool
	}

	dbURI := config.AppConfig().Postgres.URI()

	// Создаем пул соединений с базой данных
	dbPool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v\n", err.Error()))
	}

	closer.AddNamed("", func(ctx context.Context) error {
		dbPool.Close()
		return nil
	})

	// Проверяем, что соединение с базой установлено
	err = dbPool.Ping(ctx)
	if err != nil {
		panic(fmt.Sprintf("База данных недоступна: %v\n", err.Error()))
	}

	d.postgresDBPool = dbPool

	return d.postgresDBPool
}

func (d *diContainer) Migrator(ctx context.Context) migrator.Migrator {
	if d.migrator != nil {
		return d.migrator
	}
	pgxConfig := d.PostgresDBConn(ctx).Config().Copy()

	migrationsDir := config.AppConfig().Postgres.MigrationDirectory()

	d.migrator = migratorPg.NewMigrator(stdlib.OpenDB(*pgxConfig), migrationsDir)

	return d.migrator
}

func (d *diContainer) RedisPool() *redigo.Pool {
	if d.redisPool == nil {
		d.redisPool = &redigo.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}
	}

	return d.redisPool
}

func (d *diContainer) RedisClient() cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = redis.NewClient(d.RedisPool(), logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
	}

	return d.redisClient
}

func (d *diContainer) SessionRepository() repository.SessonRepository {
	if d.sessionRepository == nil {
		d.sessionRepository = sessionRepository.NewRepository(d.RedisClient())
	}

	return d.sessionRepository
}
