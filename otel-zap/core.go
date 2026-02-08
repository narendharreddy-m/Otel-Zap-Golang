package otelzap

import (
	"context"
	"time"

	otelLog "go.opentelemetry.io/otel/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// otelCore is a zapcore.Core that converts zap entries to OpenTelemetry log records.
type otelCore struct {
	ctx    context.Context
	logger otelLog.Logger
	level  zapcore.LevelEnabler
}

func defaultLevel() zapcore.LevelEnabler {
	// default to Info if user doesn't specify a level in the public API
	return zapcore.InfoLevel
}

func buildCore(cfg Config, lp otelLog.LoggerProvider) zapcore.Core {
	otelLogger := lp.Logger(cfg.ServiceName)

	core := &otelCore{
		ctx:    context.Background(),
		logger: otelLogger,
		level:  defaultLevel(),
	}

	if !cfg.EnableStdout {
		return core
	}

	stdLogger, _ := zap.NewProduction()
	return zapcore.NewTee(stdLogger.Core(), core)
}

func (c *otelCore) Enabled(lvl zapcore.Level) bool {
	return c.level.Enabled(lvl)
}

func (c *otelCore) With(fields []zapcore.Field) zapcore.Core {
	// Minimal v1: ignore With(fields). Later we can store fields and map to OTel attributes.
	return c
}

func (c *otelCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *otelCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	var rec otelLog.Record

	// Message -> Body
	rec.SetBody(otelLog.StringValue(ent.Message))

	// Severity
	rec.SetSeverityText(ent.Level.String())

	// Timestamp (zap Entry has Time)
	if !ent.Time.IsZero() {
		rec.SetTimestamp(ent.Time)
	} else {
		rec.SetTimestamp(time.Now())
	}

	// Minimal v1: skip field mapping. We'll add zap fields -> OTel attributes next.
	c.logger.Emit(c.ctx, rec)
	return nil
}

func (c *otelCore) Sync() error {
	return nil
}
