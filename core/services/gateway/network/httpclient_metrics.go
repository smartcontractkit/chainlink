package network

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

// TracePhase identifies a phase of the HTTP client request lifecycle as
// observed via net/http/httptrace.
type TracePhase string

const (
	PhaseGetConn         TracePhase = "get_conn"
	PhaseDNSLookup       TracePhase = "dns_lookup"
	PhaseTCPConnect      TracePhase = "tcp_connect"
	PhaseTLSHandshake    TracePhase = "tls_handshake"
	PhaseWroteRequest    TracePhase = "wrote_request"
	PhaseTimeToFirstByte TracePhase = "time_to_first_byte"
	PhaseTotal           TracePhase = "total"
)

type httpClientMetrics struct {
	phaseDuration metric.Int64Histogram
}

func newHTTPClientMetrics() (*httpClientMetrics, error) {
	phaseDuration, err := beholder.GetMeter().Int64Histogram(
		"platform_gateway_http_client_phase_duration",
		metric.WithUnit("ms"),
		metric.WithDescription("HTTP client request phase duration observed via httptrace. The count of phase=total observations is the request count, partitioned by method, statusCode, success, and connectionReused. Success does not imply a 2xx status code"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform_gateway_http_client_phase_duration histogram: %w", err)
	}

	return &httpClientMetrics{
		phaseDuration: phaseDuration,
	}, nil
}

func (m *httpClientMetrics) recordPhase(ctx context.Context, method string, phase TracePhase, d time.Duration) {
	m.phaseDuration.Record(ctx, d.Milliseconds(), metric.WithAttributes(
		attribute.String("method", method),
		attribute.String("phase", string(phase)),
	))
}

// recordTotal records the total request lifetime with the result attributes.
// The histogram's count for phase=total doubles as the request counter.
// Success means the request returned a response, it does not imply a successful 2xx statusCode.
func (m *httpClientMetrics) recordTotal(ctx context.Context, method string, statusCode int, success, connReused bool, d time.Duration) {
	m.phaseDuration.Record(ctx, d.Milliseconds(), metric.WithAttributes(
		attribute.String("method", method),
		attribute.String("phase", string(PhaseTotal)),
		attribute.String("statusCode", strconv.Itoa(statusCode)),
		attribute.String("success", strconv.FormatBool(success)),
		attribute.String("connectionReused", strconv.FormatBool(connReused)),
	))
}

// HTTPClientMetricViews returns histogram bucket definitions for the HTTP client trace metrics.
// Due to the OTEL specification, all histogram buckets must be defined when the beholder client is created.
func HTTPClientMetricViews() []sdkmetric.View {
	return []sdkmetric.View{
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: "platform_gateway_http_client_phase_duration"},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				// 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768 (ms)
				Boundaries: prometheus.ExponentialBuckets(1, 2, 16),
			}},
		),
	}
}
