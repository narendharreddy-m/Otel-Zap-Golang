package otelzap

// Package otelzap provides a Zap logger that exports logs to OpenTelemetry.
//
// It allows applications to keep using Zap as their logging API while
// emitting OpenTelemetry-compatible log records over OTLP.
//
// Features:
//   - Zap → OpenTelemetry Logs bridge
//   - Structured zap fields exported as log attributes
//   - Optional stdout logging
//   - Optional internal metrics
//   - Optional trace correlation via context.Context
//
// This package does NOT:
//   - Initialize tracing automatically
//   - Perform automatic environment variable parsing
//   - Provide vendor-specific exporters
//
// Tracing and metrics are explicitly configured by the user to keep
// behavior predictable and controlled.
