package otelzap

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

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

// WithContext returns a shallow copy of Logger that uses the provided context.
// If the context contains an active span, logs will be trace-correlated.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	if ctx == nil {
		return l
	}

	clone := *l

	// zapcore.NewTee returns a tee core; we must walk it
	if tee, ok := clone.Zap.Core().(interface {
		Core(i int) zapcore.Core
	}); ok {
		for i := 0; ; i++ {
			core := tee.Core(i)
			if oc, ok := core.(*otelCore); ok {
				oc.ctx = ctx
				break
			}
		}
	}

	return &clone
}
