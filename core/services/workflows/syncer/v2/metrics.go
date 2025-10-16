package v2

import (
	"context"
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
}

func (m *metrics) recordHandleDuration(ctx context.Context, d time.Duration, event string, success bool) {
	// Beholder doesn't support non-string attributes
	successStr := "false"
	if success {
		successStr = "true"
	}
	m.handleDuration.Record(ctx, d.Milliseconds(), metric.WithAttributes(
		attribute.String("success", successStr),
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

	return &metrics{
		handleDuration:   handleDuration,
		fetchedWorkflows: fetchedWorkflows,
		runningWorkflows: runningWorkflows,
		completedSyncs:   completedSyncs,
	}, nil
}
