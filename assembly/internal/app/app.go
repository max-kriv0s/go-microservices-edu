package app

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/max-kriv0s/go-microservices-edu/assembly/internal/config"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/closer"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
)

type App struct {
	diContainer *diContainer
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
	g, gCtx := errgroup.WithContext(ctx)
	// gCtx – контекст группы. Он автоматически отменяется, если одна из горутин в группе вернула ошибку.

	g.Go(func() error {
		logger.Info(ctx, "Starting order consumer service")
		if err := a.runConsumer(gCtx); err != nil {
			return fmt.Errorf("order consumer service error: %w", err)
		}
		return nil
	})

	return g.Wait()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.registerLoggerClose,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(ctx context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	opts := logger.InitOptions{
		LogLevel:     config.AppConfig().Logger.Level(),
		AsJSON:       config.AppConfig().Logger.AsJson(),
		EnableStdout: config.AppConfig().Logger.EnableStdout(),
		EnableOTLP:   config.AppConfig().Logger.EnableOTLP(),
		OTLPEndpoint: config.AppConfig().Logger.OTLPEndpoint(),
		ServiceName:  config.AppConfig().Logger.ServiceName(),
		ServiceEnv:   config.AppConfig().Logger.ServiceEnv(),
	}

	return logger.Init(ctx, opts)
}

func (a *App) initCloser(ctx context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 AssemblyRecorded Kafka consumer running")

	err := a.diContainer.OrderConsumerService().RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) registerLoggerClose(ctx context.Context) error {
	closer.AddNamed("logger zap", logger.Sync)   // Сбрасываем буферы zap
	closer.AddNamed("logger otlp", logger.Close) // Закрываем OTLP ресурсы

	return nil
}
