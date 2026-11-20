package v2

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
)

// syncerMetricLabeler wraps monitoring.MetricsLabeler to provide workflow specific utilities
// for monitoring resources
type syncerMetricLabeler struct {
	metrics.Labeler
	remoteRegistrySyncFailureCounter metric.Int64Counter
	launcherFailureCounter           metric.Int64Counter
	registryEntriesGauge             metric.Int64Gauge
}

func newSyncerMetricLabeler(meter metric.Meter) (*syncerMetricLabeler, error) {
	remoteRegistrySyncFailureCounter, err := meter.Int64Counter("platform_registrysyncer_sync_failures")
	if err != nil {
		return nil, fmt.Errorf("failed to register sync failure counter: %w", err)
	}

	launcherFailureCounter, err := meter.Int64Counter("platform_registrysyncer_launch_failures")
	if err != nil {
		return nil, fmt.Errorf("failed to register launcher failure counter: %w", err)
	}

	registryEntriesGauge, err := meter.Int64Gauge("platform_registrysyncer_registry_entries")
	if err != nil {
		return nil, fmt.Errorf("failed to register registry entries gauge: %w", err)
	}

	return &syncerMetricLabeler{
		remoteRegistrySyncFailureCounter: remoteRegistrySyncFailureCounter,
		launcherFailureCounter:           launcherFailureCounter,
		registryEntriesGauge:             registryEntriesGauge,
	}, nil
}

func (c *syncerMetricLabeler) with(keyValues ...string) syncerMetricLabeler {
	labeled := *c
	labeled.Labeler = c.With(keyValues...)
	return labeled
}

func (c *syncerMetricLabeler) incrementRemoteRegistryFailureCounter(ctx context.Context) {
	otelLabels := beholder.OtelAttributes(c.Labels).AsStringAttributes()
	c.remoteRegistrySyncFailureCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c *syncerMetricLabeler) incrementLauncherFailureCounter(ctx context.Context) {
	otelLabels := beholder.OtelAttributes(c.Labels).AsStringAttributes()
	c.launcherFailureCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

// recordRegistryEntries records how many entries of the given type ("capabilities", "dons", "nodes") were read on-chain.
func (c *syncerMetricLabeler) recordRegistryEntries(ctx context.Context, entryType string, count int) {
	otelLabels := beholder.OtelAttributes(c.Labels).AsStringAttributes()
	otelLabels = append(otelLabels, attribute.String("type", entryType))
	c.registryEntriesGauge.Record(ctx, int64(count), metric.WithAttributes(otelLabels...))
}
