package syncer

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	monutils "github.com/smartcontractkit/chainlink/v2/core/monitoring"
)

// workflow registry metrics is to locally scope instruments to avoid
// data races in testing
type workflowRegistryMetrics struct {
	activateCounter           metric.Int64Counter
	deleteCounter             metric.Int64Counter
	forceUpdateSecretsCounter metric.Int64Counter
	pauseCounter              metric.Int64Counter
	registerCounter           metric.Int64Counter
	updateCounter             metric.Int64Counter

	totalWorkflows metric.Int64Gauge
}

func initMonitoringResources() (m *workflowRegistryMetrics, err error) {
	m = &workflowRegistryMetrics{}

	m.activateCounter, err = beholder.GetMeter().Int64Counter("platform_workflow_syncer_register")
	if err != nil {
		return nil, fmt.Errorf("error initializing activate counter: %w", err)
	}

	m.deleteCounter, err = beholder.GetMeter().Int64Counter("platform_workflow_syncer_delete")
	if err != nil {
		return nil, fmt.Errorf("error initializing delete counter: %w", err)
	}

	m.forceUpdateSecretsCounter, err = beholder.GetMeter().Int64Counter("platform_workflow_syncer_force_update_secrets")
	if err != nil {
		return nil, fmt.Errorf("error initializing force update secrets counter: %w", err)
	}

	m.pauseCounter, err = beholder.GetMeter().Int64Counter("platform_workflow_syncer_pause")
	if err != nil {
		return nil, fmt.Errorf("error initializing pause counter: %w", err)
	}

	m.registerCounter, err = beholder.GetMeter().Int64Counter("platform_workflow_syncer_register")
	if err != nil {
		return nil, fmt.Errorf("error initializing register counter: %w", err)
	}

	m.updateCounter, err = beholder.GetMeter().Int64Counter("platform_workflow_syncer_update")
	if err != nil {
		return nil, fmt.Errorf("error initializing update counter: %w", err)
	}

	m.totalWorkflows, err = beholder.GetMeter().Int64Gauge("platform_workflow_syncer_total")
	if err != nil {
		return nil, fmt.Errorf("error initializing total workflows: %w", err)
	}

	return m, nil
}

// workflowRegistryMetricsLabeler wraps m to provide utilities for
// monitoring instrumentation
type workflowRegistryMetricsLabeler struct {
	metrics.Labeler
	m *workflowRegistryMetrics
}

func newWorkflowRegistryMetricsLabeler(m *workflowRegistryMetrics) workflowRegistryMetricsLabeler {
	return workflowRegistryMetricsLabeler{
		metrics.NewLabeler(),
		m,
	}
}

func (l workflowRegistryMetricsLabeler) with(keyValues ...string) workflowRegistryMetricsLabeler {
	return workflowRegistryMetricsLabeler{l.With(keyValues...), l.m}
}

func (l workflowRegistryMetricsLabeler) incrementActivateCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(l.Labels)
	l.m.activateCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (l workflowRegistryMetricsLabeler) incrementDeleteCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(l.Labels)
	l.m.deleteCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (l workflowRegistryMetricsLabeler) incrementForceUpdateSecretsCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(l.Labels)
	l.m.forceUpdateSecretsCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (l workflowRegistryMetricsLabeler) incrementPauseCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(l.Labels)
	l.m.pauseCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (l workflowRegistryMetricsLabeler) incrementRegisterCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(l.Labels)
	l.m.registerCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (l workflowRegistryMetricsLabeler) incrementUpdateCounter(ctx context.Context) {
	otelLabels := monutils.KvMapToOtelAttributes(l.Labels)
	l.m.updateCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (l workflowRegistryMetricsLabeler) updateTotalWorkflowsGauge(ctx context.Context, val int64) {
	otelLabels := monutils.KvMapToOtelAttributes(l.Labels)
	l.m.totalWorkflows.Record(ctx, val, metric.WithAttributes(otelLabels...))
}
