package connector

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

// connectorMetrics holds Beholder-backed metrics for the gateway connector.
type connectorMetrics struct {
	gatewaysPerDon metric.Int64Gauge
}

// newConnectorMetrics constructs the connector metrics. Returns nil and a
// non-nil error if the Beholder meter is unavailable; callers should treat a
// nil *connectorMetrics as "metrics disabled" and skip recording.
func newConnectorMetrics() (*connectorMetrics, error) {
	gatewaysPerDon, err := beholder.GetMeter().Int64Gauge("platform_gateway_connector_gateways_per_don")
	if err != nil {
		return nil, err
	}
	return &connectorMetrics{gatewaysPerDon: gatewaysPerDon}, nil
}

// recordGatewaysPerDon records the number of gateways configured on this node's
// gateway connector for a given DON ID. csaKey is the node's CSA key ID (hex)
// and is used as the per-node label. When csaKey is empty the call is a no-op
// so that unlabeled series are never emitted.
func (m *connectorMetrics) recordGatewaysPerDon(ctx context.Context, csaKey, donID string, count int) {
	if m == nil || csaKey == "" {
		return
	}
	m.gatewaysPerDon.Record(ctx, int64(count), metric.WithAttributes(
		attribute.String("csa_key", csaKey),
		attribute.String("gateway_don_id", donID),
	))
}
