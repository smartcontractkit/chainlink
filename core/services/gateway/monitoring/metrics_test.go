package monitoring

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestGatewayMetrics_RecordReadiness(t *testing.T) { //nolint:paralleltest // The test replaces the process-global Beholder client.
	reader := sdkmetric.NewManualReader()
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithResource(resource.NewSchemaless(attribute.String("donID", "cre-gateway-1"))),
	)
	t.Cleanup(func() { require.NoError(t, meterProvider.Shutdown(context.Background())) })

	previousClient := beholder.GetClient()
	t.Cleanup(func() { beholder.SetClient(previousClient) })
	client := beholder.NoopClientConfig{Lggr: logger.Test(t)}.New()
	client.Meter = meterProvider.Meter("gateway-test")
	client.MeterProvider = meterProvider
	beholder.SetClient(client)

	metrics, err := NewGatewayMetrics()
	require.NoError(t, err)
	metrics.RecordDONConnectionState(t.Context(), "workflow_1_zone-a", 4, 3, 4)
	metrics.RecordDONConnectionState(t.Context(), "vault_1", 7, 5, 7)
	metrics.RecordUserReady(t.Context(), true)

	var resourceMetrics metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &resourceMetrics))
	resourceDONID, ok := resourceMetrics.Resource.Set().Value("donID")
	require.True(t, ok)
	require.Equal(t, "cre-gateway-1", resourceDONID.AsString())
	requireGaugeValue(t, resourceMetrics, "platform_gateway_don_connected_nodes", 4, "workflow_1_zone-a")
	requireGaugeValue(t, resourceMetrics, "platform_gateway_don_connected_nodes", 7, "vault_1")
	requireGaugeValue(t, resourceMetrics, "platform_gateway_don_required_nodes", 3, "workflow_1_zone-a")
	requireGaugeValue(t, resourceMetrics, "platform_gateway_don_required_nodes", 5, "vault_1")
	requireGaugeValue(t, resourceMetrics, "platform_gateway_don_configured_nodes", 4, "workflow_1_zone-a")
	requireGaugeValue(t, resourceMetrics, "platform_gateway_don_configured_nodes", 7, "vault_1")
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
			for _, dataPoint := range gauge.DataPoints {
				if donID == "" {
					require.Equal(t, want, dataPoint.Value)
					return
				}
				attributeValue, ok := dataPoint.Attributes.Value("donShardID")
				if ok && attributeValue.AsString() == donID {
					require.Equal(t, want, dataPoint.Value)
					return
				}
			}
			t.Fatalf("metric %s has no datapoint for DON shard %s", name, donID)
		}
	}
	t.Fatalf("metric %s not found", name)
}
