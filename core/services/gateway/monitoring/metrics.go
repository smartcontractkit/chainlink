package monitoring

import (
	"context"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

type GatewayMetrics struct {
	nodeMsgHandleDuration metric.Int64Histogram
	nodeMsgHandleCount    metric.Int64Counter

	userMsgHandleDuration metric.Int64Histogram
	userMsgHandleCount    metric.Int64Counter

	nodeConnectedEvents    metric.Int64Counter
	keepalivePingsSent     metric.Int64Counter
	keepalivePongsReceived metric.Int64Counter
	wsPingRoundTrip        metric.Int64Histogram
	donConnectedNodes      metric.Int64Gauge
	donRequiredNodes       metric.Int64Gauge
	donConfiguredNodes     metric.Int64Gauge
	userReady              metric.Int64Gauge
}

type HTTPServerMetrics struct {
	requestDuration metric.Int64Histogram
	requestCount    metric.Int64Counter
}

func (m *GatewayMetrics) RecordNodeMsgHandlerDuration(ctx context.Context, nodeAddress string, nodeName string, duration time.Duration, success bool) {
	m.nodeMsgHandleDuration.Record(ctx, duration.Milliseconds(), metric.WithAttributes(
		attribute.String("nodeAddress", nodeAddress),
		attribute.String("nodeName", nodeName),
		attribute.String("success", strconv.FormatBool(success)),
	))
}

func (m *GatewayMetrics) RecordNodeMsgHandlerInvocation(ctx context.Context, nodeAddress string, nodeName string, success bool) {
	m.nodeMsgHandleCount.Add(ctx, 1, metric.WithAttributes(
		attribute.String("nodeAddress", nodeAddress),
		attribute.String("nodeName", nodeName),
		attribute.String("success", strconv.FormatBool(success)),
	))
}

func (m *GatewayMetrics) RecordUserMsgHandlerDuration(ctx context.Context, method string, responseCode string, duration time.Duration) {
	m.userMsgHandleDuration.Record(ctx, duration.Milliseconds(), metric.WithAttributes(
		attribute.String("method", method),
		attribute.String("responseCode", responseCode),
	))
}

func (m *GatewayMetrics) RecordUserMsgHandlerInvocation(ctx context.Context, method string, responseCode string) {
	m.userMsgHandleCount.Add(ctx, 1, metric.WithAttributes(
		attribute.String("method", method),
		attribute.String("responseCode", responseCode),
	))
}

func (m *GatewayMetrics) RecordNodeConnectedEvent(ctx context.Context, nodeAddress string, nodeName string) {
	m.nodeConnectedEvents.Add(ctx, 1, metric.WithAttributes(
		attribute.String("nodeAddress", nodeAddress),
		attribute.String("nodeName", nodeName),
	))
}

func (m *GatewayMetrics) RecordKeepalivePingsSent(ctx context.Context, nodeAddress string, nodeName string, success bool) {
	m.keepalivePingsSent.Add(ctx, 1, metric.WithAttributes(
		attribute.String("nodeAddress", nodeAddress),
		attribute.String("nodeName", nodeName),
		attribute.String("success", strconv.FormatBool(success)),
	))
}

func (m *GatewayMetrics) RecordKeepalivePongsReceived(ctx context.Context, nodeAddress string, nodeName string) {
	m.keepalivePongsReceived.Add(ctx, 1, metric.WithAttributes(
		attribute.String("nodeAddress", nodeAddress),
		attribute.String("nodeName", nodeName),
	))
}

// RecordWSPingRoundTrip records the observed round trip of a keepalive ping:
// the time from registering the probe (immediately before the ping write) until
// the pong echoing its correlation token arrives
func (m *GatewayMetrics) RecordWSPingRoundTrip(ctx context.Context, nodeAddress string, nodeName string, duration time.Duration) {
	m.wsPingRoundTrip.Record(ctx, duration.Milliseconds(), metric.WithAttributes(
		attribute.String("nodeAddress", nodeAddress),
		attribute.String("nodeName", nodeName),
	))
}

func (m *GatewayMetrics) RecordDONConnectionState(ctx context.Context, donID string, connected, required, configured int) {
	attrs := metric.WithAttributes(attribute.String("donID", donID))
	m.donConnectedNodes.Record(ctx, int64(connected), attrs)
	m.donRequiredNodes.Record(ctx, int64(required), attrs)
	m.donConfiguredNodes.Record(ctx, int64(configured), attrs)
}

func (m *GatewayMetrics) RecordUserReady(ctx context.Context, ready bool) {
	value := int64(0)
	if ready {
		value = 1
	}
	m.userReady.Record(ctx, value)
}

func NewGatewayMetrics() (*GatewayMetrics, error) {
	return NewGatewayMetricsWithMeter(beholder.GetMeter())
}

func NewGatewayMetricsWithMeter(meter metric.Meter) (*GatewayMetrics, error) {
	nodeMsgHandleDuration, err := meter.Int64Histogram("platform_gateway_node_msg_handler_duration_ms")
	if err != nil {
		return nil, err
	}

	nodeMsgHandleCount, err := meter.Int64Counter("platform_gateway_node_msgs_handled_total")
	if err != nil {
		return nil, err
	}

	userMsgHandleDuration, err := meter.Int64Histogram("platform_gateway_user_msg_handler_duration_ms")
	if err != nil {
		return nil, err
	}

	userMsgHandleCount, err := meter.Int64Counter("platform_gateway_user_msgs_handled_total")
	if err != nil {
		return nil, err
	}

	nodeConnectedEvents, err := meter.Int64Counter("platform_gateway_node_connected_events_total")
	if err != nil {
		return nil, err
	}

	keepalivePingsSent, err := meter.Int64Counter("platform_gateway_keepalive_pings_sent_total")
	if err != nil {
		return nil, err
	}

	keepalivePongsReceived, err := meter.Int64Counter("platform_gateway_keepalive_pongs_received_total")
	if err != nil {
		return nil, err
	}

	wsPingRoundTrip, err := meter.Int64Histogram("platform_gateway_ws_ping_round_trip_ms",
		metric.WithUnit("ms"),
		metric.WithDescription("Observed websocket ping round trip in milliseconds, from ping enqueue to receipt of the pong echoing its correlation token. Includes local write-pump queue wait, transport, and peer control-frame processing; it is not a pure network RTT. Only the connection's latest probe yields a sample"),
	)
	if err != nil {
		return nil, err
	}

	donConnectedNodes, err := meter.Int64Gauge("platform_gateway_don_connected_nodes")
	if err != nil {
		return nil, err
	}

	donRequiredNodes, err := meter.Int64Gauge("platform_gateway_don_required_nodes")
	if err != nil {
		return nil, err
	}

	donConfiguredNodes, err := meter.Int64Gauge("platform_gateway_don_configured_nodes")
	if err != nil {
		return nil, err
	}

	userReady, err := meter.Int64Gauge("platform_gateway_user_ready")
	if err != nil {
		return nil, err
	}

	return &GatewayMetrics{
		nodeMsgHandleDuration:  nodeMsgHandleDuration,
		nodeMsgHandleCount:     nodeMsgHandleCount,
		userMsgHandleDuration:  userMsgHandleDuration,
		userMsgHandleCount:     userMsgHandleCount,
		nodeConnectedEvents:    nodeConnectedEvents,
		keepalivePingsSent:     keepalivePingsSent,
		keepalivePongsReceived: keepalivePongsReceived,
		wsPingRoundTrip:        wsPingRoundTrip,
		donConnectedNodes:      donConnectedNodes,
		donRequiredNodes:       donRequiredNodes,
		donConfiguredNodes:     donConfiguredNodes,
		userReady:              userReady,
	}, nil
}

func (m *HTTPServerMetrics) RecordRequestDuration(ctx context.Context, responseCode int, duration time.Duration) {
	m.requestDuration.Record(ctx, duration.Milliseconds(), metric.WithAttributes(
		attribute.Int("responseCode", responseCode),
	))
}

func (m *HTTPServerMetrics) RecordRequestCount(ctx context.Context, responseCode int) {
	m.requestCount.Add(ctx, 1, metric.WithAttributes(
		attribute.Int("responseCode", responseCode),
	))
}

func NewHTTPServerMetrics() (*HTTPServerMetrics, error) {
	requestDuration, err := beholder.GetMeter().Int64Histogram("platform_gateway_http_server_request_duration_ms")
	if err != nil {
		return nil, err
	}

	requestCount, err := beholder.GetMeter().Int64Counter("platform_gateway_http_server_requests_total")
	if err != nil {
		return nil, err
	}
	return &HTTPServerMetrics{requestDuration: requestDuration, requestCount: requestCount}, nil
}

// MetricViews returns histogram bucket definitions for this package's metrics.
// Due to the OTEL specification, all histogram buckets must be defined when the beholder client is created.
func MetricViews() []sdkmetric.View {
	return []sdkmetric.View{
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: "platform_gateway_ws_ping_round_trip_ms"},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				// 1ms up to ~32s: healthy pongs return in milliseconds, but the
				// read deadline lets a slow round trip stretch to tens of seconds.
				Boundaries: prometheus.ExponentialBuckets(1, 2, 16),
			}},
		),
	}
}
