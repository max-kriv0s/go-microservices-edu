package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var serviceName = "assembly_service"

var meter = otel.Meter(serviceName)

// AssemblyDuration — гистограмма длительности сборки
var AssemblyDuration metric.Float64Histogram

func InitMetrics() error {
	var err error

	AssemblyDuration, err = meter.Float64Histogram(
		serviceName+"_assembly_duration_seconds",
		metric.WithDescription("Время сборки заказов в секундах"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(
			0.001, 0.005, 0.01, 0.025, 0.05,
			0.1, 0.25, 0.5, 1.0, 2.0, 5.0,
		),
	)
	if err != nil {
		return err
	}

	return nil
}
