package main

import (
	"context"
	"time"

	otelzap "github.com/narendharreddy-m/Otel-Zap-Golang/otel-zap"
	"go.uber.org/zap"
)

func main() {
	l, err := otelzap.New(otelzap.Config{
		Endpoint:     "localhost:4317",
		ServiceName:  "demo-service",
		EnableStdout: true,
		Timeout:      5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer func() { _ = l.Shutdown(context.Background()) }()

	l.Zap.Info("hello from otel-zap", zap.String("env", "local"))
}
