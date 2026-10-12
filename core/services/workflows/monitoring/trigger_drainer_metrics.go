package monitoring

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
)

// Drop reason values for platform_engine_trigger_drainer_dropped_events_total (low cardinality).
const (
	TriggerDrainerDropReasonHookBusy = "hook_busy"
)

// TriggerDrainerMetrics are the trigger queue drainer's instruments.
type TriggerDrainerMetrics struct {
	droppedEventsCounter metric.Int64Counter
	hookDurationSeconds  metric.Float64Histogram
}

// InitTriggerDrainerMetrics registers the drainer instruments on the beholder meter.
func InitTriggerDrainerMetrics() (*TriggerDrainerMetrics, error) {
	return NewTriggerDrainerMetrics(beholder.GetMeter())
}

// NewTriggerDrainerMetrics registers the drainer instruments on meter.
func NewTriggerDrainerMetrics(meter metric.Meter) (*TriggerDrainerMetrics, error) {
	var (
		tm  TriggerDrainerMetrics
		err error
	)
	tm.droppedEventsCounter, err = meter.Int64Counter(
		"platform_engine_trigger_drainer_dropped_events_total",
		metric.WithDescription("Trigger events dropped by the queue drainer, by drop_reason"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register trigger drainer dropped events counter: %w", err)
	}

	tm.hookDurationSeconds, err = meter.Float64Histogram(
		"platform_engine_trigger_drainer_hook_duration_seconds",
		metric.WithDescription("Duration of a trigger drainer OnObservedEvents hook call"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register trigger drainer hook duration histogram: %w", err)
	}
	return &tm, nil
}

// TriggerDrainerMetricLabeler wraps metrics.Labeler to record the drainer's
// instruments with labels.
type TriggerDrainerMetricLabeler struct {
	metrics.Labeler
	tm *TriggerDrainerMetrics
}

func NewTriggerDrainerMetricLabeler(labeler metrics.Labeler, tm *TriggerDrainerMetrics) *TriggerDrainerMetricLabeler {
	return &TriggerDrainerMetricLabeler{labeler, tm}
}

func (c TriggerDrainerMetricLabeler) With(keyValues ...string) *TriggerDrainerMetricLabeler {
	return &TriggerDrainerMetricLabeler{c.Labeler.With(keyValues...), c.tm}
}

// IncrementDroppedEventsCounter adds n dropped events. reason must be a
// low-cardinality value; use TriggerDrainerDropReason* constants.
func (c TriggerDrainerMetricLabeler) IncrementDroppedEventsCounter(ctx context.Context, reason string, n int) {
	lc := c.With(platform.KeyTriggerDropReason, reason)
	otelLabels := beholder.OtelAttributes(lc.Labels).AsStringAttributes()
	lc.tm.droppedEventsCounter.Add(ctx, int64(n), metric.WithAttributes(otelLabels...))
}

func (c TriggerDrainerMetricLabeler) RecordHookDurationSeconds(ctx context.Context, seconds float64) {
	otelLabels := beholder.OtelAttributes(c.Labels).AsStringAttributes()
	c.tm.hookDurationSeconds.Record(ctx, seconds, metric.WithAttributes(otelLabels...))
}
