## otel-zap-golang

`otel-zap-golang` is a lightweight Go library that lets you use **Zap** for logging while exporting logs through **OpenTelemetry Logs (OTLP)**.

In simple terms:
**you keep Zap**, and your logs automatically flow into any OpenTelemetry-compatible backend.

---

## Why this exists

Zap is fast and widely used, but on its own it doesn’t integrate with OpenTelemetry logs.

This library:

* keeps Zap as the logging API
* converts Zap logs into OpenTelemetry log records
* exports them using OTLP
* optionally still logs to stdout

This makes it easy to adopt OpenTelemetry **without rewriting existing Zap logging**.

---

## Features

* Zap → OpenTelemetry Logs bridge
* Structured zap fields are exported as OpenTelemetry log attributes
* OTLP gRPC export support
* Optional stdout logging (Zap production logger)
* Clean shutdown and flush support
* Minimal configuration

**Out of scope (by design):**
* Metrics (planned)
* Automatic tracing setup (user-controlled)

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
		Endpoint:      "localhost:4317",
		ServiceName:   "my-service",
		EnableStdout:  true,
		Timeout:       5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer logger.Shutdown(context.Background())

	logger.Zap.Info("hello from otel-zap", zap.String("env", "local"))
}
```

---
## Trace correlation (optional)

This library supports **context-based trace correlation**.

If you log using a `context.Context` that contains an active span, the OpenTelemetry Logs SDK will automatically attach:

* `trace_id`
* `span_id`

to the emitted log record.

### Example

```go
import "go.opentelemetry.io/otel"

tracer := otel.Tracer("example")

ctx, span := tracer.Start(context.Background(), "demo-span")
defer span.End()

logger.WithContext(ctx).Zap.Info(
	"processing request",
	zap.String("env", "prod"),
)
```

When a span is present, logs and traces can be correlated in observability backends.

> Note: This library does **not** initialize tracing automatically.
> Users are expected to configure tracing explicitly if needed.

---

## How it works (high level)

1. Your application logs using Zap
2. A custom Zap core converts each log entry into an OpenTelemetry log record
3. Structured zap fields become OpenTelemetry attributes
4. Logs are exported using OTLP
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

A minimal collector configuration for local testing:

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
```

Run the collector:

```bash
otelcol-contrib --config otel-collector.yaml
```

Then run your app and verify logs appear in the collector output.

---
