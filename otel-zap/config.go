package otelzap

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// Config controls how the zap -> OpenTelemetry logger behaves.
type Config struct {
	// OTLP gRPC endpoint, e.g. "localhost:4317"
	Endpoint string

	// ServiceName is reported as service.name in OTel resource.
	ServiceName string

	// EnableStdout tees logs to stdout using zap's production logger.
	EnableStdout bool

	// Timeout used when creating the OTLP exporter.
	Timeout time.Duration

	// Level controls which zap logs are emitted to OTel.
	// If nil, defaults to zapcore.InfoLevel in core.go via defaultLevel().
	Level zapcore.Level

	EnableMetrics bool
}
