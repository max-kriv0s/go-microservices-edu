package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/max-kriv0s/go-microservices-edu/order/internal/api/health"
	"github.com/max-kriv0s/go-microservices-edu/order/internal/config"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/closer"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
	orderV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/openapi/order/v1"
)

type App struct {
	diContainer *diContainer
	httpServer  *orderV1.Server
	router      *chi.Mux
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	// Канал для ошибок от компонентов
	errCh := make(chan error, 2)

	// Контекст для остановки всех горутин
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Консьюмер
	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- errors.Errorf("consumer crashed: %v", err)
		}
	}()

	// HTTP сервер
	go func() {
		if err := a.runHTTPServer(ctx); err != nil {
			errCh <- errors.Errorf("http server crashed: %v", err)
		}
	}()

	// Ожидание либо ошибки, либо завершения контекста (например, сигнал SIGINT/SIGTERM)
	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))
		// Триггерим cancel, чтобы остановить второй компонент
		cancel()
		// Дождись завершения всех задач (если есть graceful shutdown внутри)
		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initHTTPServer,
		a.initRouter,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(ctx context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	orderServer, err := orderV1.NewServer(a.diContainer.OrderV1Api(ctx))
	if err != nil {
		return err
	}

	a.httpServer = orderServer

	return nil
}

func (a *App) initRouter(ctx context.Context) error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	serverTimeout := config.AppConfig().OrderHTTP.ServerTimeout()
	r.Use(middleware.Timeout(serverTimeout))

	r.Get("/health", health.Handler)

	r.Mount("/", a.httpServer)

	a.router = r

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	serverAddr := config.AppConfig().OrderHTTP.Address()

	logger.Info(ctx, fmt.Sprintf("🚀 HTTP OrderService server listening on %s", serverAddr))

	server := &http.Server{
		Addr:              serverAddr,
		Handler:           a.router,
		ReadHeaderTimeout: config.AppConfig().OrderHTTP.ReadHeaderTimeout(),
		ReadTimeout:       config.AppConfig().OrderHTTP.ReadTimeout(),
	}

	closer.AddNamed("", func(ctx context.Context) error {
		return server.Shutdown(ctx)
	})

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(ctx, fmt.Sprintf("❌ Ошибка запуска сервера: %v\n", err.Error()))

		return err
	}

	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderAssembled Kafka consumer running")

	err := a.diContainer.OrderConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
