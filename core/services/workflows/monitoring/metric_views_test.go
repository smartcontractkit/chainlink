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

func TestMetricViews_platformEngineHistogramBuckets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		instrument string
		unit       string
		wantBounds []float64
	}{
		{
			name:       "workflow_completed_time",
			instrument: "platform_engine_workflow_completed_time_seconds",
			unit:       "s",
			wantBounds: []float64{0, 10, 40, 90, 150, 300, 900},
		},
		{
			name:       "trigger_event_queue_wait",
			instrument: "platform_engine_trigger_event_queue_wait_seconds",
			unit:       "s",
			wantBounds: []float64{0, 0.001, 0.01, 0.1, 1, 10, 60},
		},
		{
			name:       "trigger_queue_to_execution_start",
			instrument: "platform_engine_trigger_queue_to_execution_start_seconds",
			unit:       "s",
			wantBounds: []float64{0, 0.01, 0.1, 1, 10, 60, 300},
		},
		{
			name:       "trigger_payload_bytes",
			instrument: "platform_engine_trigger_payload_bytes",
			unit:       "By",
			wantBounds: []float64{0, 512, 4096, 32768, 262144, 1048576, 4194304},
		},
		{
			name:       "execution_semaphore_wait",
			instrument: "platform_engine_execution_semaphore_wait_seconds",
			unit:       "s",
			wantBounds: []float64{0, 0.001, 0.01, 0.1, 1, 10, 60},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reader := sdkmetric.NewManualReader()
			mp := sdkmetric.NewMeterProvider(
				sdkmetric.WithReader(reader),
				sdkmetric.WithView(monitoring.MetricViews()...),
			)
			t.Cleanup(func() { _ = mp.Shutdown(context.Background()) })

			hist, err := mp.Meter("test").Float64Histogram(
				tt.instrument,
				metric.WithUnit(tt.unit),
			)
			require.NoError(t, err)
			hist.Record(context.Background(), 1)

			var rm metricdata.ResourceMetrics
			require.NoError(t, reader.Collect(context.Background(), &rm))

			require.Len(t, rm.ScopeMetrics, 1)
			require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
			require.Equal(t, tt.instrument, rm.ScopeMetrics[0].Metrics[0].Name)

			data, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[float64])
			require.True(t, ok)
			require.Len(t, data.DataPoints, 1)
			assert.Equal(t, tt.wantBounds, data.DataPoints[0].Bounds)
			assert.Len(t, tt.wantBounds, 7, "expected 7 boundaries (8 Prometheus buckets including +Inf)")
		})
	}
}

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
	assert.Equal(t, wantBoundaries, data.DataPoints[0].Bounds)
	assert.Len(t, data.DataPoints[0].Bounds, 7, "expected 7 boundaries (8 Prometheus buckets including +Inf)")
}
