package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/max-kriv0s/go-microservices-edu/payment/internal/config"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/closer"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/grpc/health"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/tracing"
	paymentV1 "github.com/max-kriv0s/go-microservices-edu/shared/pkg/proto/payment/v1"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
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
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.registerLoggerClose,
		a.initListener,
		a.initGRPCServer,
		a.initTracing,
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
	opts := logger.InitOptions{
		LogLevel:     config.AppConfig().Logger.Level(),
		AsJSON:       config.AppConfig().Logger.AsJson(),
		EnableStdout: config.AppConfig().Logger.EnableStdout(),
		EnableOTLP:   config.AppConfig().Logger.EnableOTLP(),
		OTLPEndpoint: config.AppConfig().OtelCollector.CollectorEndpoint(),
		ServiceName:  config.AppConfig().OtelCollector.ServiceName(),
		ServiceEnv:   config.AppConfig().OtelCollector.ServiceEnv(),
	}

	return logger.Init(ctx, opts)
}

func (a *App) initCloser(ctx context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initListener(ctx context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().PaymentGRPC.Address())
	if err != nil {
		return err
	}
	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})

	a.listener = listener

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(tracing.UnaryServerInterceptor(config.AppConfig().OtelCollector.ServiceName())),
	)

	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register((a.grpcServer))

	// Регистрируем health service для проверки работоспособности
	health.RegisterService(a.grpcServer)

	paymentV1.RegisterPaymentServiceServer(a.grpcServer, a.diContainer.PaymentV1API(ctx))

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC PaymentService server listening on %s", config.AppConfig().PaymentGRPC.Address()))

	err := a.grpcServer.Serve(a.listener)
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

func (a *App) initTracing(ctx context.Context) error {
	err := tracing.InitTracer(ctx, config.AppConfig().OtelCollector)
	if err != nil {
		return err
	}

	closer.AddNamed("tracer", tracing.ShutdownTracer)

	return nil
}
