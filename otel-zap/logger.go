package otelzap

import (
	"go.uber.org/zap"

	otelLog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type Logger struct {
	Zap *zap.Logger
	lp  *sdklog.LoggerProvider
}

// New creates a zap.Logger that exports logs through OpenTelemetry Logs (OTLP).
func New(cfg Config) (*Logger, error) {
	lp, err := newLoggerProvider(cfg)
	if err != nil {
		return nil, err
	}

	// sdk LoggerProvider implements the otel/log LoggerProvider interface.
	var apiLP otelLog.LoggerProvider = lp

	core := buildCore(cfg, apiLP)
	zlogger := zap.New(core)

	return &Logger{
		Zap: zlogger,
		lp:  lp,
	}, nil
}
