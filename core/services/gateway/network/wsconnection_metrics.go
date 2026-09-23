package network

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
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
}

type noopWSConnectionObserver struct{}

func (noopWSConnectionObserver) RecordWriteQueueWait(context.Context, time.Duration)   {}
func (noopWSConnectionObserver) RecordSocketWrite(context.Context, time.Duration)      {}
func (noopWSConnectionObserver) RecordReadDispatchWait(context.Context, time.Duration) {}
func (noopWSConnectionObserver) AddPendingWriters(context.Context, int64)              {}


type WSConnectionMetrics struct {
	writeQueueWait   metric.Int64Histogram
	socketWrite      metric.Int64Histogram
	readDispatchWait metric.Int64Histogram
	pendingWriters   metric.Int64UpDownCounter
}

func NewWSConnectionMetrics() (*WSConnectionMetrics, error) {
	meter := beholder.GetMeter()

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

	return &WSConnectionMetrics{
		writeQueueWait:   writeQueueWait,
		socketWrite:      socketWrite,
		readDispatchWait: readDispatchWait,
		pendingWriters:   pendingWriters,
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
	return []sdkmetric.View{
		histogramView("platform_gateway_ws_write_queue_wait_ms"),
		histogramView("platform_gateway_ws_socket_write_ms"),
		histogramView("platform_gateway_ws_read_dispatch_wait_ms"),
	}
}
