package otelzap

import "context"

// Shutdown flushes logs and metrics.
func (l *Logger) Shutdown(ctx context.Context) error {
	if err := l.lp.Shutdown(ctx); err != nil {
		return err
	}
	if l.metricsShutdown != nil {
		return l.metricsShutdown(ctx)
	}
	return nil
}
