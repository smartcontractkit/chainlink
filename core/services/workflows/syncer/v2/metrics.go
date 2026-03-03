package v2

import (
	"context"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

type metrics struct {
	handleDuration   metric.Int64Histogram
	fetchedWorkflows metric.Int64Gauge
	runningWorkflows metric.Int64Gauge
	completedSyncs   metric.Int64Counter

	// Per-source metrics for multi-source observability
	sourceHealth        metric.Int64Gauge     // 1=healthy, 0=unhealthy per source
	workflowsPerSource  metric.Int64Gauge     // workflows fetched per source
	sourceFetchDuration metric.Int64Histogram // fetch latency per source
	sourceFetchErrors   metric.Int64Counter   // error count per source

	// Per-tick reconciliation metrics
	reconcileEventsDispatched metric.Int64Histogram // events dispatched per source per tick
	reconcileDuration         metric.Int64Histogram // wall-clock ms for parallel event processing
	reconcileEventsBackoff    metric.Int64Counter   // events skipped due to backoff
}

func (m *metrics) recordHandleDuration(ctx context.Context, d time.Duration, event string, success bool) {
	m.handleDuration.Record(ctx, d.Milliseconds(), metric.WithAttributes(
		attribute.String("success", strconv.FormatBool(success)),
		attribute.String("eventType", event),
	))
}

func (m *metrics) recordFetchedWorkflows(ctx context.Context, count int) {
	m.fetchedWorkflows.Record(ctx, int64(count))
}

func (m *metrics) recordRunningWorkflows(ctx context.Context, count int) {
	m.runningWorkflows.Record(ctx, int64(count))
}

func (m *metrics) incrementCompletedSyncs(ctx context.Context) {
	m.completedSyncs.Add(ctx, 1)
}

// recordSourceFetch records metrics for a source fetch operation.
func (m *metrics) recordSourceFetch(ctx context.Context, sourceName string, workflowCount int, duration time.Duration, err error) {
	attrs := metric.WithAttributes(attribute.String("source", sourceName))

	// Record fetch duration
	m.sourceFetchDuration.Record(ctx, duration.Milliseconds(), attrs)

	// Record workflow count per source
	m.workflowsPerSource.Record(ctx, int64(workflowCount), attrs)

	// Record health status (1=healthy, 0=unhealthy)
	if err != nil {
		m.sourceHealth.Record(ctx, 0, attrs)
		m.sourceFetchErrors.Add(ctx, 1, attrs)
	} else {
		m.sourceHealth.Record(ctx, 1, attrs)
	}
}

func (m *metrics) recordReconcileBatch(ctx context.Context, source string, dispatched int, duration time.Duration) {
	attrs := metric.WithAttributes(attribute.String("source", source))
	m.reconcileEventsDispatched.Record(ctx, int64(dispatched), attrs)
	m.reconcileDuration.Record(ctx, duration.Milliseconds(), attrs)
}

func (m *metrics) recordReconcileBackoff(ctx context.Context, source string, count int) {
	m.reconcileEventsBackoff.Add(ctx, int64(count), metric.WithAttributes(
		attribute.String("source", source),
	))
}

func newMetrics() (*metrics, error) {
	handleDuration, err := beholder.GetMeter().Int64Histogram("platform_workflow_registry_syncer_handler_duration_ms")
	if err != nil {
		return nil, err
	}

	fetchedWorkflows, err := beholder.GetMeter().Int64Gauge("platform_workflow_registry_syncer_fetched_workflows")
	if err != nil {
		return nil, err
	}

	runningWorkflows, err := beholder.GetMeter().Int64Gauge("platform_workflow_registry_syncer_running_workflows")
	if err != nil {
		return nil, err
	}

	completedSyncs, err := beholder.GetMeter().Int64Counter("platform_workflow_registry_syncer_completed_syncs_total")
	if err != nil {
		return nil, err
	}

	// Per-source metrics
	sourceHealth, err := beholder.GetMeter().Int64Gauge("platform_workflow_registry_syncer_source_health")
	if err != nil {
		return nil, err
	}

	workflowsPerSource, err := beholder.GetMeter().Int64Gauge("platform_workflow_registry_syncer_workflows_per_source")
	if err != nil {
		return nil, err
	}

	sourceFetchDuration, err := beholder.GetMeter().Int64Histogram("platform_workflow_registry_syncer_source_fetch_duration_ms")
	if err != nil {
		return nil, err
	}

	sourceFetchErrors, err := beholder.GetMeter().Int64Counter("platform_workflow_registry_syncer_source_fetch_errors_total")
	if err != nil {
		return nil, err
	}

	reconcileEventsDispatched, err := beholder.GetMeter().Int64Histogram("platform_workflow_registry_syncer_reconcile_events_dispatched")
	if err != nil {
		return nil, err
	}

	reconcileDuration, err := beholder.GetMeter().Int64Histogram("platform_workflow_registry_syncer_reconcile_duration_ms")
	if err != nil {
		return nil, err
	}

	reconcileEventsBackoff, err := beholder.GetMeter().Int64Counter("platform_workflow_registry_syncer_reconcile_events_backoff_total")
	if err != nil {
		return nil, err
	}

	return &metrics{
		handleDuration:            handleDuration,
		fetchedWorkflows:          fetchedWorkflows,
		runningWorkflows:          runningWorkflows,
		completedSyncs:            completedSyncs,
		sourceHealth:              sourceHealth,
		workflowsPerSource:        workflowsPerSource,
		sourceFetchDuration:       sourceFetchDuration,
		sourceFetchErrors:         sourceFetchErrors,
		reconcileEventsDispatched: reconcileEventsDispatched,
		reconcileDuration:         reconcileDuration,
		reconcileEventsBackoff:    reconcileEventsBackoff,
	}, nil
}
