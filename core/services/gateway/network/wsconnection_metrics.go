package network

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// WSConnectionObserver observes the blocking operations of a single
// WSConnectionWrapper: how long writes wait to be accepted by the write pump,
// how long each socket write takes, how long read messages wait to be handed to
// the consumer, and how many writers are currently queued. Implementations must
// be cheap and non-blocking.
type WSConnectionObserver interface {
	// RecordWriteQueueWait records the time a Write call spent in the first
	// select, up to the moment the write pump accepted the item. It excludes
	// the subsequent wait for the write result.
	RecordWriteQueueWait(ctx context.Context, d time.Duration)
	// RecordSocketWrite records the duration of a single conn.WriteMessage
	// call in the write pump, including failed writes.
	RecordSocketWrite(ctx context.Context, d time.Duration)
	// RecordReadDispatchWait records the time a successfully read message spent
	// waiting to be delivered to the read channel consumer. ReadMessage itself
	// is excluded because it blocks for normal idle time.
	RecordReadDispatchWait(ctx context.Context, d time.Duration)
	// AddPendingWriters adjusts the count of goroutines blocked in Write's
	// first select. Every increment must be paired with a decrement on exit,
	// including cancellation and shutdown paths.
	AddPendingWriters(ctx context.Context, delta int64)
	// RecordConnectionClosed records why the wrapper closed a live connection.
	// It is called once per connection.
	RecordConnectionClosed(ctx context.Context, reason CloseReason)
	// RecordPingRTT records the round trip of a keepalive ping, from sending
	// it to receiving its pong.
	RecordPingRTT(ctx context.Context, d time.Duration)
}

// CloseReason says which mechanism closed a connection.
type CloseReason string

const (
	// CloseReasonReadTimeout means a read deadline expired, as when keepalive pongs stop arriving.
	CloseReasonReadTimeout CloseReason = "read_timeout"
	CloseReasonReadError   CloseReason = "read_error"
	// CloseReasonWriteTimeout means a write deadline expired.
	CloseReasonWriteTimeout CloseReason = "write_timeout"
	CloseReasonWriteError   CloseReason = "write_error"
	// CloseReasonShutdown means the wrapper was closed or told to disconnect.
	CloseReasonShutdown CloseReason = "shutdown"
	// CloseReasonReplaced means Reset swapped in a new connection while the old one was live.
	CloseReasonReplaced CloseReason = "replaced"
)

type noopWSConnectionObserver struct{}

func (noopWSConnectionObserver) RecordWriteQueueWait(context.Context, time.Duration)   {}
func (noopWSConnectionObserver) RecordSocketWrite(context.Context, time.Duration)      {}
func (noopWSConnectionObserver) RecordReadDispatchWait(context.Context, time.Duration) {}
func (noopWSConnectionObserver) AddPendingWriters(context.Context, int64)              {}
func (noopWSConnectionObserver) RecordConnectionClosed(context.Context, CloseReason)   {}
func (noopWSConnectionObserver) RecordPingRTT(context.Context, time.Duration)          {}

type WSConnectionMetrics struct {
	writeQueueWait   metric.Int64Histogram
	socketWrite      metric.Int64Histogram
	readDispatchWait metric.Int64Histogram
	pendingWriters   metric.Int64UpDownCounter
	connectionClosed metric.Int64Counter
	pingRTT          metric.Float64Histogram
}

func NewWSConnectionMetrics(meter metric.Meter) (*WSConnectionMetrics, error) {
	writeQueueWait, err := meter.Int64Histogram(
		"platform_gateway_ws_write_queue_wait_ms",
		metric.WithUnit("ms"),
		metric.WithDescription("Time a websocket write spent waiting for the write pump to accept it (the first select in Write), in milliseconds. Excludes the subsequent wait for the write result"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform_gateway_ws_write_queue_wait_ms histogram: %w", err)
	}

	socketWrite, err := meter.Int64Histogram(
		"platform_gateway_ws_socket_write_ms",
		metric.WithUnit("ms"),
		metric.WithDescription("Duration of a blocking websocket conn.WriteMessage call in the write pump, in milliseconds, including errors. Only completed writes produce samples; a write that never returns leaves no observation"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform_gateway_ws_socket_write_ms histogram: %w", err)
	}

	readDispatchWait, err := meter.Int64Histogram(
		"platform_gateway_ws_read_dispatch_wait_ms",
		metric.WithUnit("ms"),
		metric.WithDescription("Time a successfully read websocket message spent waiting to be handed to the read channel consumer, in milliseconds. ReadMessage blocking time is excluded because it includes normal idle time"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform_gateway_ws_read_dispatch_wait_ms histogram: %w", err)
	}

	pendingWriters, err := meter.Int64UpDownCounter(
		"platform_gateway_ws_pending_writers",
		metric.WithDescription("Number of goroutines currently blocked waiting to enqueue a websocket write. The write channel is unbuffered, so this cannot be derived from its length"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform_gateway_ws_pending_writers up-down counter: %w", err)
	}

	connectionClosed, err := meter.Int64Counter(
		"platform_gateway_ws_connection_closed_total",
		metric.WithDescription("Number of live websocket connections closed by the wrapper, by reason: read_timeout, read_error, write_timeout, write_error, shutdown or replaced"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform_gateway_ws_connection_closed_total counter: %w", err)
	}

	pingRTT, err := meter.Float64Histogram(
		"platform_gateway_ws_ping_rtt_ms",
		metric.WithUnit("ms"),
		metric.WithDescription("Round trip of a websocket keepalive ping, from sending it to receiving its pong, in milliseconds. Only pings that were answered produce samples"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform_gateway_ws_ping_rtt_ms histogram: %w", err)
	}

	return &WSConnectionMetrics{
		writeQueueWait:   writeQueueWait,
		socketWrite:      socketWrite,
		readDispatchWait: readDispatchWait,
		pendingWriters:   pendingWriters,
		connectionClosed: connectionClosed,
		pingRTT:          pingRTT,
	}, nil
}

func (m *WSConnectionMetrics) Observer(attrs ...attribute.KeyValue) WSConnectionObserver {
	if m == nil {
		return noopWSConnectionObserver{}
	}
	return &wsConnectionObserver{metrics: m, attrs: metric.WithAttributes(attrs...)}
}

type wsConnectionObserver struct {
	metrics *WSConnectionMetrics
	attrs   metric.MeasurementOption
}

func (o *wsConnectionObserver) RecordWriteQueueWait(ctx context.Context, d time.Duration) {
	o.metrics.writeQueueWait.Record(ctx, d.Milliseconds(), o.attrs)
}

func (o *wsConnectionObserver) RecordSocketWrite(ctx context.Context, d time.Duration) {
	o.metrics.socketWrite.Record(ctx, d.Milliseconds(), o.attrs)
}

func (o *wsConnectionObserver) RecordReadDispatchWait(ctx context.Context, d time.Duration) {
	o.metrics.readDispatchWait.Record(ctx, d.Milliseconds(), o.attrs)
}

func (o *wsConnectionObserver) AddPendingWriters(ctx context.Context, delta int64) {
	o.metrics.pendingWriters.Add(ctx, delta, o.attrs)
}

// WSConnectionMetricViews returns histogram bucket definitions for the
// websocket connection wrapper metrics. Due to the OTEL specification, all
// histogram buckets must be defined when the beholder client is created.
func WSConnectionMetricViews() []sdkmetric.View {
	histogramView := func(name string) sdkmetric.View {
		return sdkmetric.NewView(
			sdkmetric.Instrument{Name: name},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				// 1ms up to ~32s: queue waits and socket writes are normally
				// sub-second, but degraded peers can block for many seconds.
				Boundaries: prometheus.ExponentialBuckets(1, 2, 16),
			}},
		)
	}
	pingRTTView := sdkmetric.NewView(
		sdkmetric.Instrument{Name: "platform_gateway_ws_ping_rtt_ms"},
		sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
			// Normal round trips are a few ms. The top buckets show a peer
			// degrading towards the 30s pong timeout.
			Boundaries: []float64{1, 2, 5, 10, 25, 50, 100, 250, 1000, 5000, 15000, 30000},
		}},
	)
	return []sdkmetric.View{
		pingRTTView,
		histogramView("platform_gateway_ws_write_queue_wait_ms"),
		histogramView("platform_gateway_ws_socket_write_ms"),
		histogramView("platform_gateway_ws_read_dispatch_wait_ms"),
	}
}

func (o *wsConnectionObserver) RecordConnectionClosed(ctx context.Context, reason CloseReason) {
	o.metrics.connectionClosed.Add(ctx, 1, o.attrs, metric.WithAttributes(attribute.String("reason", string(reason))))
}

func (o *wsConnectionObserver) RecordPingRTT(ctx context.Context, d time.Duration) {
	o.metrics.pingRTT.Record(ctx, float64(d)/float64(time.Millisecond), o.attrs)
}
