package monitoring

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestGatewayMetrics_RecordReadiness(t *testing.T) { //nolint:paralleltest // The test replaces the process-global Beholder client.
	reader := sdkmetric.NewManualReader()
	meterProvider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { require.NoError(t, meterProvider.Shutdown(context.Background())) })

	previousClient := beholder.GetClient()
	t.Cleanup(func() { beholder.SetClient(previousClient) })
	client := beholder.NoopClientConfig{Lggr: logger.Test(t)}.New()
	client.Meter = meterProvider.Meter("gateway-test")
	client.MeterProvider = meterProvider
	beholder.SetClient(client)

	metrics, err := NewGatewayMetrics()
	require.NoError(t, err)
	metrics.RecordDONConnectionState(t.Context(), "workflow_1_zone-a", 5, 7, 10)
	metrics.RecordUserReady(t.Context(), true)

	var resourceMetrics metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &resourceMetrics))
	requireGaugeValue(t, resourceMetrics, "platform_gateway_don_connected_nodes", 5, "workflow_1_zone-a")
	requireGaugeValue(t, resourceMetrics, "platform_gateway_don_required_nodes", 7, "workflow_1_zone-a")
	requireGaugeValue(t, resourceMetrics, "platform_gateway_don_configured_nodes", 10, "workflow_1_zone-a")
	requireGaugeValue(t, resourceMetrics, "platform_gateway_user_ready", 1, "")

	metrics.RecordUserReady(t.Context(), false)
	require.NoError(t, reader.Collect(t.Context(), &resourceMetrics))
	requireGaugeValue(t, resourceMetrics, "platform_gateway_user_ready", 0, "")
}

func requireGaugeValue(t *testing.T, resourceMetrics metricdata.ResourceMetrics, name string, want int64, donID string) {
	t.Helper()
	for _, scopeMetrics := range resourceMetrics.ScopeMetrics {
		for _, metric := range scopeMetrics.Metrics {
			if metric.Name != name {
				continue
			}
			gauge, ok := metric.Data.(metricdata.Gauge[int64])
			require.True(t, ok)
			require.Len(t, gauge.DataPoints, 1)
			require.Equal(t, want, gauge.DataPoints[0].Value)
			if donID != "" {
				attributeValue, ok := gauge.DataPoints[0].Attributes.Value("donID")
				require.True(t, ok)
				require.Equal(t, donID, attributeValue.AsString())
			}
			return
		}
	}
	t.Fatalf("metric %s not found", name)
}
