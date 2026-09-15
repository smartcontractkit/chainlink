package headreporter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func newTestBeholderHeadMetrics(t *testing.T) (*beholderHeadMetrics, *sdkmetric.ManualReader) {
	t.Helper()
	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { _ = mp.Shutdown(context.Background()) })

	m := mp.Meter("test")
	latestNumber, err := m.Int64Gauge("head_reporter_latest_block_number")
	require.NoError(t, err)
	latestTs, err := m.Int64Gauge("head_reporter_latest_block_timestamp_seconds")
	require.NoError(t, err)
	finalizedNumber, err := m.Int64Gauge("head_reporter_finalized_block_number")
	require.NoError(t, err)
	finalizedTs, err := m.Int64Gauge("head_reporter_finalized_block_timestamp_seconds")
	require.NoError(t, err)
	finalityDepth, err := m.Int64Gauge("head_reporter_finality_depth_blocks")
	require.NoError(t, err)

	return &beholderHeadMetrics{
		latestNumber:    latestNumber,
		latestTs:        latestTs,
		finalizedNumber: finalizedNumber,
		finalizedTs:     finalizedTs,
		finalityDepth:   finalityDepth,
	}, reader
}

func collectGauge(t *testing.T, rm metricdata.ResourceMetrics, name string) metricdata.Gauge[int64] {
	t.Helper()
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name == name {
				data, ok := m.Data.(metricdata.Gauge[int64])
				require.True(t, ok, "metric %s is not an int64 gauge", name)
				return data
			}
		}
	}
	t.Fatalf("metric %s not found", name)
	return metricdata.Gauge[int64]{}
}

func Test_BeholderHeadMetrics_RecordHeadReport_WithFinalized(t *testing.T) {
	t.Parallel()
	metrics, reader := newTestBeholderHeadMetrics(t)

	metrics.RecordHeadReport(t.Context(), headReport{
		chainID:       "100",
		network:       "evm",
		chainSelector: 465200170687744372,
		hasSelector:   true,
		latestNumber:  42,
		latestTs:      1000,
		finalized: &finalizedBlock{
			number: 40,
			ts:     900,
		},
	})

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))

	latest := collectGauge(t, rm, "head_reporter_latest_block_number")
	require.Len(t, latest.DataPoints, 1)
	require.Equal(t, int64(42), latest.DataPoints[0].Value)

	depth := collectGauge(t, rm, "head_reporter_finality_depth_blocks")
	require.Len(t, depth.DataPoints, 1)
	require.Equal(t, int64(2), depth.DataPoints[0].Value)

	finalizedNumber := collectGauge(t, rm, "head_reporter_finalized_block_number")
	require.Len(t, finalizedNumber.DataPoints, 1)
	require.Equal(t, int64(40), finalizedNumber.DataPoints[0].Value)
}

func Test_BeholderHeadMetrics_RecordHeadReport_NoFinalized(t *testing.T) {
	t.Parallel()
	metrics, reader := newTestBeholderHeadMetrics(t)

	metrics.RecordHeadReport(t.Context(), headReport{
		chainID:      "testchain",
		network:      "solana",
		latestNumber: 42,
		latestTs:     1000,
	})

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))

	latest := collectGauge(t, rm, "head_reporter_latest_block_number")
	require.Len(t, latest.DataPoints, 1)
	require.Equal(t, int64(42), latest.DataPoints[0].Value)

	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name == "head_reporter_finalized_block_number" || m.Name == "head_reporter_finality_depth_blocks" {
				data, ok := m.Data.(metricdata.Gauge[int64])
				require.True(t, ok)
				require.Empty(t, data.DataPoints)
			}
		}
	}
}
