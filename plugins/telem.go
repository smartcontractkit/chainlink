package plugins

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"

	"github.com/smartcontractkit/chainlink/v2/core/static"
)

// TelemetryProviders simplifies creating and closing a [*trace.TracerProvider] and a [*metric.MeterProvider] for a
// common [resource.Resource].
type TelemetryProviders struct {
	*trace.TracerProvider
	*metric.MeterProvider
}

func NewTelemetryProviders(serviceName, serviceInstanceID string) (*TelemetryProviders, error) {
	sha, ver := static.Short()
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(fmt.Sprintf("%s@%s", ver, sha)),
		semconv.ServiceInstanceID(serviceInstanceID),
		semconv.ProcessPID(os.Getpid()),
	)
	texp, err := stdouttrace.New(
		stdouttrace.WithWriter(io.Discard), //TODO where to send traces?
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}
	mexp, err := otelprom.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	tp := trace.NewTracerProvider(trace.WithBatcher(texp), trace.WithResource(res))
	mp := metric.NewMeterProvider(metric.WithReader(mexp), metric.WithResource(res))
	return &TelemetryProviders{tp, mp}, nil
}

func (p *TelemetryProviders) Close() error {
	ctx := context.Background()
	return errors.Join(p.TracerProvider.Shutdown(ctx), p.MeterProvider.Shutdown(ctx))
}

// writerFunc is a Writer implemented by the underlying func.
type writerFunc func(p []byte) (int, error)

func (f writerFunc) Write(p []byte) (int, error) {
	return f(p)
}
