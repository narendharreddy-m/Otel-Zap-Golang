## otel-zap-golang

`otel-zap-golang` is a lightweight Go library that lets you use **Zap** for logging while exporting logs through **OpenTelemetry Logs (OTLP)**.

In simple terms:  
**you keep Zap**, and your logs automatically flow into any OpenTelemetry-compatible backend — with optional trace correlation and internal metrics.

---

## Why this exists

Zap is fast and widely used, but on its own it doesn’t integrate with OpenTelemetry logs.

This library:

- keeps Zap as the logging API
- converts Zap logs into OpenTelemetry log records
- exports them using OTLP
- optionally still logs to stdout
- exposes internal metrics about the logging pipeline

This makes it easy to adopt OpenTelemetry **without rewriting existing Zap logging**.

---

## Features

### Logs
- Zap → OpenTelemetry Logs bridge
- Structured zap fields exported as OpenTelemetry log attributes
- OTLP gRPC export support
- Optional stdout logging (Zap production logger)

### Trace correlation (optional)
- Context-based correlation with active OpenTelemetry spans
- Logs automatically correlate with traces when a span is present
- Tracing setup remains user-controlled

### Metrics (optional)
- Internal logger health metrics via OpenTelemetry Metrics
- Zero overhead unless explicitly enabled

### Operational
- Clean shutdown and flush support
- Minimal, explicit configuration
- No global side effects unless enabled

---

## Installation

```bash
go get github.com/narendharreddy-m/otel-zap-golang
```

---

## Basic usage

```go
package main

import (
	"context"
	"time"

	"github.com/narendharreddy-m/otel-zap-golang/otelzap"
	"go.uber.org/zap"
)

func main() {
	logger, err := otelzap.New(otelzap.Config{
		Endpoint:     "localhost:4317",
		ServiceName:  "my-service",
		EnableStdout: true,
		Timeout:      5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer logger.Shutdown(context.Background())

	logger.Zap.Info("hello from otel-zap", zap.String("env", "local"))
}
```

---

## Structured log attributes

Structured zap fields are automatically converted into OpenTelemetry log attributes.

```go
logger.Zap.Info(
	"request processed",
	zap.String("env", "prod"),
	zap.Int("status", 200),
)
```

Results in OpenTelemetry logs like:

```
Body: "request processed"
Attributes:
  env = "prod"
  status = 200
```

This enables filtering, querying, and alerting in observability backends.

---

## Trace correlation (optional)

This library supports **context-based trace correlation**.

If you log using a `context.Context` that contains an active span, logs will automatically correlate with that span.

### Example

```go
import "go.opentelemetry.io/otel"

tracer := otel.Tracer("example-service")

ctx, span := tracer.Start(context.Background(), "demo-span")
defer span.End()

logger.WithContext(ctx).Zap.Info(
	"processing request",
	zap.String("env", "prod"),
)
```

> Note: This library does **not** initialize tracing automatically.
> Users are expected to configure tracing explicitly if needed.

---

## Metrics (optional)

`otel-zap-golang` can emit **internal metrics** describing the health of the logging pipeline.

### Available metrics

| Metric name                         | Description                                   |
| ----------------------------------- | --------------------------------------------- |
| `otelzap_log_records_total`         | Total number of log records emitted           |
| `otelzap_log_records_error_total`   | Number of error-level log records             |
| `otelzap_log_records_dropped_total` | Number of logs dropped due to level filtering |

### Enabling metrics

Metrics are **opt-in**:

```go
logger, err := otelzap.New(otelzap.Config{
	Endpoint:       "localhost:4317",
	ServiceName:    "demo-service",
	EnableStdout:   true,
	EnableMetrics:  true,
})
```

Metrics are exported using OTLP and can be consumed by any OpenTelemetry-compatible backend.

---

## How it works (high level)

1. Your application logs using Zap
2. A custom Zap core converts log entries into OpenTelemetry log records
3. Structured fields become OpenTelemetry attributes
4. Logs, traces, and metrics are exported using OTLP
5. An OpenTelemetry Collector or backend receives and processes them

```
Zap Logger
   ↓
Custom Zap Core
   ↓
OpenTelemetry Logs SDK
   ↓
OTLP Exporter
   ↓
Collector / Backend
```

---

## Local testing

This library has been validated end-to-end using the **OpenTelemetry Collector (contrib)**.

Minimal collector configuration for logs, traces, and metrics:

```yaml
receivers:
  otlp:
    protocols:
      grpc:

exporters:
  debug:
    verbosity: detailed

service:
  pipelines:
    logs:
      receivers: [otlp]
      exporters: [debug]

    traces:
      receivers: [otlp]
      exporters: [debug]

    metrics:
      receivers: [otlp]
      exporters: [debug]
```

Run the collector:

```bash
otelcol-contrib --config otel-collector.yaml
```

Then run your app and verify logs, traces, and metrics appear in the output.

---
