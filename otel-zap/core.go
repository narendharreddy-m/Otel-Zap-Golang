package otelzap

import (
	"context"
	"math"
	"time"

	otelLog "go.opentelemetry.io/otel/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// otelCore is a zapcore.Core that converts zap entries to OpenTelemetry log records.
type otelCore struct {
	ctx     context.Context
	logger  otelLog.Logger
	level   zapcore.LevelEnabler
	metrics *metricsState // nil when metrics disabled
}

func defaultLevel() zapcore.LevelEnabler {
	return zapcore.InfoLevel
}

func buildCore(cfg Config, lp otelLog.LoggerProvider, m *metricsState) zapcore.Core {
	otelLogger := lp.Logger(cfg.ServiceName)

	core := &otelCore{
		ctx:     context.Background(),
		logger:  otelLogger,
		level:   defaultLevel(),
		metrics: m,
	}

	if !cfg.EnableStdout {
		return core
	}

	stdLogger, _ := zap.NewProduction()
	return zapcore.NewTee(stdLogger.Core(), core)
}

func (c *otelCore) Enabled(lvl zapcore.Level) bool {
	enabled := c.level.Enabled(lvl)

	if !enabled && c.metrics != nil {
		c.metrics.droppedTotal.Add(c.ctx, 1)
	}

	return enabled
}

func (c *otelCore) With(fields []zapcore.Field) zapcore.Core {
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

	// Metrics: total logs
	if c.metrics != nil {
		c.metrics.logsTotal.Add(c.ctx, 1)
		if ent.Level >= zapcore.ErrorLevel {
			c.metrics.logsErrorTotal.Add(c.ctx, 1)
		}
	}

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

	// Convert zap fields -> OpenTelemetry attributes
	attrs := make([]otelLog.KeyValue, 0, len(fields))

	for _, f := range fields {
		switch f.Type {

		case zapcore.StringType:
			attrs = append(attrs, otelLog.String(f.Key, f.String))

		case zapcore.Int64Type:
			attrs = append(attrs, otelLog.Int64(f.Key, f.Integer))

		case zapcore.Int32Type:
			attrs = append(attrs, otelLog.Int64(f.Key, int64(f.Integer)))

		case zapcore.BoolType:
			attrs = append(attrs, otelLog.Bool(f.Key, f.Integer == 1))

		case zapcore.Float64Type:
			attrs = append(attrs, otelLog.Float64(
				f.Key,
				math.Float64frombits(uint64(f.Integer)),
			))

		case zapcore.ErrorType:
			if err, ok := f.Interface.(error); ok {
				attrs = append(attrs, otelLog.String(f.Key, err.Error()))
			}
		}
	}

	if len(attrs) > 0 {
		rec.AddAttributes(attrs...)
	}

	c.logger.Emit(c.ctx, rec)
	return nil
}

func (c *otelCore) Sync() error {
	return nil
}
