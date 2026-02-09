package otelzap

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

// metricsState holds internal logger metrics.
type metricsState struct {
	logsTotal      metric.Int64Counter
	logsErrorTotal metric.Int64Counter
	droppedTotal   metric.Int64Counter
}

// initMetrics initializes OTEL metrics for the logger.
func initMetrics(
	ctx context.Context,
	cfg Config,
	res *resource.Resource,
) (*metricsState, func(context.Context) error, error) {

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	exp, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(cfg.Endpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, nil, err
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(exp),
		),
		sdkmetric.WithResource(res),
	)

	otel.SetMeterProvider(provider)

	meter := provider.Meter("otelzap")

	logsTotal, _ := meter.Int64Counter(
		"otelzap_log_records_total",
		metric.WithDescription("Total number of log records emitted"),
	)

	logsErrorTotal, _ := meter.Int64Counter(
		"otelzap_log_records_error_total",
		metric.WithDescription("Total number of error log records"),
	)

	droppedTotal, _ := meter.Int64Counter(
		"otelzap_log_records_dropped_total",
		metric.WithDescription("Total number of dropped log records"),
	)

	state := &metricsState{
		logsTotal:      logsTotal,
		logsErrorTotal: logsErrorTotal,
		droppedTotal:   droppedTotal,
	}

	return state, provider.Shutdown, nil
}
