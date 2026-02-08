package otelzap

import "context"

// Shutdown flushes and shuts down the OpenTelemetry logger provider.
// Call this on app exit.
func (l *Logger) Shutdown(ctx context.Context) error {
	if l == nil || l.lp == nil {
		return nil
	}
	return l.lp.Shutdown(ctx)
}
