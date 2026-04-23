package metrics

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

// Metrics contains metrics for consensus capability
type Metrics struct {
	PendingConsensusRequests      metric.Int64Gauge
	batchCapacityExceeded         metric.Int64Counter
	batchRequestsTotal            metric.Int64Counter
	requestSizeHistogram          metric.Float64Histogram
	observationBatchSizeHistogram metric.Float64Histogram
	outcomeBatchSizeHistogram     metric.Float64Histogram
}

// NewMetrics creates a new instance of Metrics
func NewMetrics() (*Metrics, error) {
	m := &Metrics{}
	if err := m.init(); err != nil {
		return nil, err
	}
	return m, nil
}

func MetricViews() []sdkmetric.View {
	return []sdkmetric.View{
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: "consensus_capability_request_size_bytes"},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				Boundaries: []float64{0, 10, 100, 1000, 10000, 100000, 1000000},
			}},
		),
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: "consensus_capability_observation_batch_size_bytes"},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				Boundaries: []float64{0, 10, 100, 1000, 10000, 100000, 1000000},
			}},
		),
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: "consensus_capability_outcome_batch_size_bytes"},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				Boundaries: []float64{0, 10, 100, 1000, 10000, 100000, 1000000},
			}},
		),
	}
}

func (m *Metrics) init() error {
	meter := beholder.GetMeter()
	var err error

	m.PendingConsensusRequests, err = meter.Int64Gauge(
		"consensus_capability_pending_consensus_requests",
		metric.WithDescription("Number of consensus requests awaiting a response"),
	)
	if err != nil {
		return fmt.Errorf("failed to create pending consensus requests gauge: %w", err)
	}

	m.batchCapacityExceeded, err = meter.Int64Counter(
		"consensus_capability_batch_capacity_exceeded",
		metric.WithDescription("Number of times the batch capacity has been exceeded"),
	)
	if err != nil {
		return fmt.Errorf("failed to create batch capacity exceeded gauge: %w", err)
	}

	m.batchRequestsTotal, err = meter.Int64Counter(
		"consensus_capability_batch_requests_total",
		metric.WithDescription("Total number of batch requests"),
	)
	if err != nil {
		return fmt.Errorf("failed to create batch requests total counter: %w", err)
	}

	m.requestSizeHistogram, err = meter.Float64Histogram(
		"consensus_capability_request_size_bytes",
		metric.WithDescription("Histogram of consensus request sizes in bytes"),
	)
	if err != nil {
		return fmt.Errorf("failed to create request size histogram: %w", err)
	}

	m.observationBatchSizeHistogram, err = meter.Float64Histogram(
		"consensus_capability_observation_batch_size_bytes",
		metric.WithDescription("Histogram of consensus observation batch sizes in bytes"),
	)
	if err != nil {
		return fmt.Errorf("failed to create consensus observation batch size histogram: %w", err)
	}

	m.outcomeBatchSizeHistogram, err = meter.Float64Histogram(
		"consensus_capability_outcome_batch_size_bytes",
		metric.WithDescription("Histogram of consensus outcome batch sizes in bytes"),
	)
	if err != nil {
		return fmt.Errorf("failed to create consensus outcome batch size histogram: %w", err)
	}

	return nil
}

func (m *Metrics) SetPendingRequestsCount(ctx context.Context, pendingRequests int64) {
	m.PendingConsensusRequests.Record(ctx, pendingRequests)
}

func (m *Metrics) IncBatchCapacityExceeded(ctx context.Context, step string) {
	m.batchCapacityExceeded.Add(ctx, 1, metric.WithAttributes(attribute.String("step", step)))
}

func (m *Metrics) IncBatchRequestsTotal(ctx context.Context, step string) {
	m.batchRequestsTotal.Add(ctx, 1, metric.WithAttributes(attribute.String("step", step)))
}

func (m *Metrics) RecordRequestSize(ctx context.Context, size float64) {
	m.requestSizeHistogram.Record(ctx, size)
}

func (m *Metrics) RecordObservationBatchSize(ctx context.Context, size float64) {
	m.observationBatchSizeHistogram.Record(ctx, size)
}

func (m *Metrics) RecordOutcomeBatchSize(ctx context.Context, size float64) {
	m.outcomeBatchSizeHistogram.Record(ctx, size)
}
