package registrysyncer

import (
	"context"
	"fmt"

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
}

func newSyncerMetricLabeler() (*syncerMetricLabeler, error) {
	remoteRegistrySyncFailureCounter, err := beholder.GetMeter().Int64Counter("platform_registrysyncer_sync_failures")
	if err != nil {
		return nil, fmt.Errorf("failed to register sync failure counter: %w", err)
	}

	launcherFailureCounter, err := beholder.GetMeter().Int64Counter("platform_registrysyncer_launch_failures")
	if err != nil {
		return nil, fmt.Errorf("failed to register launcher failure counter: %w", err)
	}

	return &syncerMetricLabeler{remoteRegistrySyncFailureCounter: remoteRegistrySyncFailureCounter, launcherFailureCounter: launcherFailureCounter}, nil
}

func (c *syncerMetricLabeler) with(keyValues ...string) syncerMetricLabeler {
	return syncerMetricLabeler{c.With(keyValues...), c.remoteRegistrySyncFailureCounter, c.launcherFailureCounter}
}

func (c *syncerMetricLabeler) incrementRemoteRegistryFailureCounter(ctx context.Context) {
	otelLabels := beholder.OtelAttributes(c.Labels).AsStringAttributes()
	c.remoteRegistrySyncFailureCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c *syncerMetricLabeler) incrementLauncherFailureCounter(ctx context.Context) {
	otelLabels := beholder.OtelAttributes(c.Labels).AsStringAttributes()
	c.launcherFailureCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}
