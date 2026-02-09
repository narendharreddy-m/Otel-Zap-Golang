package main

import (
	"context"
	"time"

	otelzap "github.com/narendharreddy-m/Otel-Zap-Golang/otel-zap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// Initialize tracing
	shutdownTracing, err := initTracing(ctx)
	if err != nil {
		panic(err)
	}
	defer shutdownTracing(ctx)

	// --- Initialize logger (logs + metrics) ---
	l, err := otelzap.New(otelzap.Config{
		Endpoint:      "localhost:4317",
		ServiceName:   "demo-service",
		EnableStdout:  true,
		EnableMetrics: true,
		Timeout:       5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer func() { _ = l.Shutdown(ctx) }()

	// Create a span
	tracer := otel.Tracer("example-service")
	ctx, span := tracer.Start(ctx, "demo-span")
	defer span.End()

	// --- Log with context (logs + traces + metrics) ---
	l.WithContext(ctx).Zap.Info(
		"hello from otel-zap",
		zap.String("env", "local"),
	)
}

func initTracing(ctx context.Context) (func(context.Context) error, error) {
	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("localhost:4317"),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("demo-service"),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}
