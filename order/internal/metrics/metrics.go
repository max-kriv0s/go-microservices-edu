package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var serviceName = "order_service"

var meter = otel.Meter(serviceName)

var (
	// OrdersTotal — счетчик созданных заказов
	OrderTotal metric.Int64Counter

	// OrdersRevenueTotal — счетчик суммарной выручки
	OrdersRevenueTotal metric.Float64Counter
)

func InitMetrics() error {
	var err error

	OrderTotal, err = meter.Int64Counter(serviceName+"_order_total", metric.WithDescription("Количество созданных заказов"))
	if err != nil {
		return err
	}

	OrdersRevenueTotal, err = meter.Float64Counter(
		serviceName+"_orders_revenue_total",
		metric.WithDescription("Cуммарная выручка по собранным заказам"),
		metric.WithUnit("currency"),
	)
	if err != nil {
		return err
	}

	return nil
}
