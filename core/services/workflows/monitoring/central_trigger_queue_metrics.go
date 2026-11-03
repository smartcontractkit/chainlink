package monitoring

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
)

// Drop reason values for platform_engine_central_trigger_queue_dropped_events_total (low cardinality).
const (
	CentralTriggerQueueDropReasonQueueFull = "queue_full"
)

// CentralTriggerQueueMetrics are the central trigger queue's instruments. They are separate from
// EngineMetrics because the queue is node-wide and its metrics must never carry workflow labels.
type CentralTriggerQueueMetrics struct {
	putsCounter          metric.Int64Counter
	droppedEventsCounter metric.Int64Counter
	drainedEventsCounter metric.Int64Counter
	depthGauge           metric.Int64Gauge
	capacityGauge        metric.Int64Gauge
}

// InitCentralTriggerQueueMetrics registers the central trigger queue instruments on the beholder meter.
func InitCentralTriggerQueueMetrics() (*CentralTriggerQueueMetrics, error) {
	return NewCentralTriggerQueueMetrics(beholder.GetMeter())
}

// NewCentralTriggerQueueMetrics registers the central trigger queue instruments on meter.
func NewCentralTriggerQueueMetrics(meter metric.Meter) (*CentralTriggerQueueMetrics, error) {
	var (
		qm  CentralTriggerQueueMetrics
		err error
	)
	qm.putsCounter, err = meter.Int64Counter(
		"platform_engine_central_trigger_queue_puts_total",
		metric.WithDescription("Trigger events accepted by the central trigger queue"),
	)
	if err != nil {
		return nil, fmt.Errorf("register central trigger queue puts counter: %w", err)
	}

	qm.droppedEventsCounter, err = meter.Int64Counter(
		"platform_engine_central_trigger_queue_dropped_events_total",
		metric.WithDescription("Trigger events the central trigger queue rejected, by drop_reason"),
	)
	if err != nil {
		return nil, fmt.Errorf("register central trigger queue dropped events counter: %w", err)
	}

	qm.drainedEventsCounter, err = meter.Int64Counter(
		"platform_engine_central_trigger_queue_drained_events_total",
		metric.WithDescription("Trigger events drained from the central trigger queue"),
	)
	if err != nil {
		return nil, fmt.Errorf("register central trigger queue drained events counter: %w", err)
	}

	qm.depthGauge, err = meter.Int64Gauge(
		"platform_engine_central_trigger_queue_depth",
		metric.WithDescription("Trigger events in the central trigger queue"),
	)
	if err != nil {
		return nil, fmt.Errorf("register central trigger queue depth gauge: %w", err)
	}

	qm.capacityGauge, err = meter.Int64Gauge(
		"platform_engine_central_trigger_queue_capacity",
		metric.WithDescription("Current capacity of the central trigger queue"),
	)
	if err != nil {
		return nil, fmt.Errorf("register central trigger queue capacity gauge: %w", err)
	}
	return &qm, nil
}

// CentralTriggerQueueMetricLabeler wraps metrics.Labeler to record the central trigger queue's
// instruments with node-level labels.
type CentralTriggerQueueMetricLabeler struct {
	metrics.Labeler
	qm *CentralTriggerQueueMetrics
}

// NewCentralTriggerQueueMetricLabeler returns a labeler that records qm with the labels of labeler. Pass node-level
// labels only.
func NewCentralTriggerQueueMetricLabeler(labeler metrics.Labeler, qm *CentralTriggerQueueMetrics) *CentralTriggerQueueMetricLabeler {
	return &CentralTriggerQueueMetricLabeler{labeler, qm}
}

// With returns a copy of the labeler with keyValues added to its labels.
func (c CentralTriggerQueueMetricLabeler) With(keyValues ...string) *CentralTriggerQueueMetricLabeler {
	return &CentralTriggerQueueMetricLabeler{c.Labeler.With(keyValues...), c.qm}
}

// IncrementPutsCounter counts one accepted put.
func (c CentralTriggerQueueMetricLabeler) IncrementPutsCounter(ctx context.Context) {
	c.qm.putsCounter.Add(ctx, 1, c.attributes())
}

// IncrementDroppedEventsCounter adds n dropped events. Use the CentralTriggerQueueDropReason* constants.
func (c CentralTriggerQueueMetricLabeler) IncrementDroppedEventsCounter(ctx context.Context, reason string, n int) {
	lc := c.With(platform.KeyTriggerDropReason, reason)
	lc.qm.droppedEventsCounter.Add(ctx, int64(n), lc.attributes())
}

// IncrementDrainedEventsCounter adds n drained events.
func (c CentralTriggerQueueMetricLabeler) IncrementDrainedEventsCounter(ctx context.Context, n int) {
	c.qm.drainedEventsCounter.Add(ctx, int64(n), c.attributes())
}

// RecordDepth records the number of queued events.
func (c CentralTriggerQueueMetricLabeler) RecordDepth(ctx context.Context, depth int) {
	c.qm.depthGauge.Record(ctx, int64(depth), c.attributes())
}

// RecordCapacity records the queue's current capacity.
func (c CentralTriggerQueueMetricLabeler) RecordCapacity(ctx context.Context, capacity int) {
	c.qm.capacityGauge.Record(ctx, int64(capacity), c.attributes())
}

func (c CentralTriggerQueueMetricLabeler) attributes() metric.MeasurementOption {
	return metric.WithAttributes(beholder.OtelAttributes(c.Labels).AsStringAttributes()...)
}
