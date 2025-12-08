package monitoring

import (
	"context"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

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

func NewGatewayMetrics() (*GatewayMetrics, error) {
	nodeMsgHandleDuration, err := beholder.GetMeter().Int64Histogram("platform_gateway_node_msg_handler_duration_ms")
	if err != nil {
		return nil, err
	}

	nodeMsgHandleCount, err := beholder.GetMeter().Int64Counter("platform_gateway_node_msgs_handled_total")
	if err != nil {
		return nil, err
	}

	userMsgHandleDuration, err := beholder.GetMeter().Int64Histogram("platform_gateway_user_msg_handler_duration_ms")
	if err != nil {
		return nil, err
	}

	userMsgHandleCount, err := beholder.GetMeter().Int64Counter("platform_gateway_user_msgs_handled_total")
	if err != nil {
		return nil, err
	}

	nodeConnectedEvents, err := beholder.GetMeter().Int64Counter("platform_gateway_node_connected_events_total")
	if err != nil {
		return nil, err
	}

	keepalivePingsSent, err := beholder.GetMeter().Int64Counter("platform_gateway_keepalive_pings_sent_total")
	if err != nil {
		return nil, err
	}

	keepalivePongsReceived, err := beholder.GetMeter().Int64Counter("platform_gateway_keepalive_pongs_received_total")
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
