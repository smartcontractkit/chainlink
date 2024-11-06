package registrysyncer

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"

	localMonitoring "github.com/smartcontractkit/chainlink/v2/core/monitoring"
)

var remoteRegistrySyncFailureCounter metric.Int64Counter
var launcherFailureCounter metric.Int64Counter

func initMonitoringResources() (err error) {
	remoteRegistrySyncFailureCounter, err = beholder.GetMeter().Int64Counter("platform_registrysyncer_sync_failures")
	if err != nil {
		return fmt.Errorf("failed to register sync failure counter: %w", err)
	}

	launcherFailureCounter, err = beholder.GetMeter().Int64Counter("platform_registrysyncer_launch_failures")
	if err != nil {
		return fmt.Errorf("failed to register launcher failure counter: %w", err)
	}

	return nil
}

// syncerMetricLabeler wraps monitoring.MetricsLabeler to provide workflow specific utilities
// for monitoring resources
type syncerMetricLabeler struct {
	metrics.Labeler
}

func (c syncerMetricLabeler) with(keyValues ...string) syncerMetricLabeler {
	return syncerMetricLabeler{c.With(keyValues...)}
}

func (c syncerMetricLabeler) incrementRemoteRegistryFailureCounter(ctx context.Context) {
	otelLabels := localMonitoring.KvMapToOtelAttributes(c.Labels)
	remoteRegistrySyncFailureCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c syncerMetricLabeler) incrementLauncherFailureCounter(ctx context.Context) {
	otelLabels := localMonitoring.KvMapToOtelAttributes(c.Labels)
	launcherFailureCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}
