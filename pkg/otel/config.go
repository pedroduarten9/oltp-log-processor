package otel

import (
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Tracer trace.Tracer
	Meter  metric.Meter
	Logger *slog.Logger
}

func NewConfig(name string) Config {
	return Config{
		Tracer: otel.Tracer(name),
		Meter:  otel.Meter(name),
		Logger: otelslog.NewLogger(name),
	}
}
