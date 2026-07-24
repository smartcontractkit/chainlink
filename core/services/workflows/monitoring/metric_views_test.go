package monitoring_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

func TestMetricViews_getSecretsDurationBuckets(t *testing.T) {
	t.Parallel()

	wantBoundaries := []float64{0, 10, 50, 100, 250, 500, 1000}

	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithView(monitoring.MetricViews()...),
	)
	t.Cleanup(func() { _ = mp.Shutdown(context.Background()) })

	hist, err := mp.Meter("test").Int64Histogram(
		"platform_engine_get_secrets_duration_ms",
		metric.WithUnit("ms"),
	)
	require.NoError(t, err)
	hist.Record(context.Background(), 75)

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(context.Background(), &rm))

	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
	require.Equal(t, "platform_engine_get_secrets_duration_ms", rm.ScopeMetrics[0].Metrics[0].Name)

	data, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[int64])
	require.True(t, ok)
	require.Len(t, data.DataPoints, 1)
	// 7 boundaries produce 8 Prometheus buckets (including the implicit +Inf).
	assert.Equal(t, wantBoundaries, data.DataPoints[0].Bounds)
}
