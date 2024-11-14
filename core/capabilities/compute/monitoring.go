package compute

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"

	localMonitoring "github.com/smartcontractkit/chainlink/v2/core/monitoring"
)

const timestampKey = "computeTimestamp"

type computeMetricsLabeler struct {
	metrics.Labeler
	computeHTTPRequestCounter metric.Int64Counter
}

func newComputeMetricsLabeler(l metrics.Labeler) (*computeMetricsLabeler, error) {
	computeHTTPRequestCounter, err := beholder.GetMeter().Int64Counter("capabilities_compute_http_request_count")
	if err != nil {
		return nil, fmt.Errorf("failed to register compute http request counter: %w", err)
	}

	return &computeMetricsLabeler{Labeler: l, computeHTTPRequestCounter: computeHTTPRequestCounter}, nil
}

func (c *computeMetricsLabeler) with(keyValues ...string) *computeMetricsLabeler {
	return &computeMetricsLabeler{c.With(keyValues...), c.computeHTTPRequestCounter}
}

func (c *computeMetricsLabeler) incrementHTTPRequestCounter(ctx context.Context) {
	otelLabels := localMonitoring.KvMapToOtelAttributes(c.Labels)
	c.computeHTTPRequestCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}
