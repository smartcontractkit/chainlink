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

var computeHTTPRequestCounter metric.Int64Counter

func initMonitoringResources() (err error) {
	computeHTTPRequestCounter, err = beholder.GetMeter().Int64Counter("capabilities_compute_http_request_count")
	if err != nil {
		return fmt.Errorf("failed to register compute http request counter: %w", err)
	}

	return nil
}

type computeMetricsLabeler struct {
	metrics.Labeler
}

func (c computeMetricsLabeler) with(keyValues ...string) computeMetricsLabeler {
	return computeMetricsLabeler{c.With(keyValues...)}
}

func (c computeMetricsLabeler) incrementHTTPRequestCounter(ctx context.Context) {
	otelLabels := localMonitoring.KvMapToOtelAttributes(c.Labels)
	computeHTTPRequestCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}
