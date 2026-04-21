package metrics

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

// FunctionType represents the OCR plugin function being measured
type FunctionType string

const (
	Query               FunctionType = "query"
	Observation         FunctionType = "observation"
	ValidateObservation FunctionType = "validateObservation"
	// OCR3 specific
	Outcome FunctionType = "outcome"
	// OCR3.1 specific
	ObservationQuorum FunctionType = "observationQuorum"
	StateTransition   FunctionType = "stateTransition"
	Committed         FunctionType = "committed"
	// Common
	Reports        FunctionType = "reports"
	ShouldAccept   FunctionType = "shouldAccept"
	ShouldTransmit FunctionType = "shouldTransmit"
)

// PluginMetrics holds OTEL metrics for OCR plugin instrumentation
type PluginMetrics struct {
	plugin       string
	configDigest string

	durations        metric.Int64Histogram
	reportsGenerated metric.Int64Counter
	sizes            metric.Int64Histogram
	status           metric.Int64Gauge
}

// NewPluginMetrics creates metrics with the given prefix (e.g., "platform_ocr3_reporting_plugin" or "platform_ocr3_1_reporting_plugin")
func NewPluginMetrics(metricPrefix, plugin, configDigest string) (*PluginMetrics, error) {
	durations, err := beholder.GetMeter().Int64Histogram(metricPrefix+"_duration_ms", metric.WithUnit("ms"))
	if err != nil {
		return nil, fmt.Errorf("failed to create duration histogram: %w", err)
	}

	reportsGenerated, err := beholder.GetMeter().Int64Counter(metricPrefix+"_reports_processed", metric.WithUnit("1"))
	if err != nil {
		return nil, fmt.Errorf("failed to create reports counter: %w", err)
	}

	sizes, err := beholder.GetMeter().Int64Histogram(metricPrefix+"_data_sizes", metric.WithUnit("By"))
	if err != nil {
		return nil, fmt.Errorf("failed to create sizes histogram: %w", err)
	}

	status, err := beholder.GetMeter().Int64Gauge(metricPrefix + "_status")
	if err != nil {
		return nil, fmt.Errorf("failed to create status gauge: %w", err)
	}

	return &PluginMetrics{
		plugin:           plugin,
		configDigest:     configDigest,
		durations:        durations,
		reportsGenerated: reportsGenerated,
		sizes:            sizes,
		status:           status,
	}, nil
}

// RecordDuration records the duration of a function execution
func (m *PluginMetrics) RecordDuration(ctx context.Context, function FunctionType, d time.Duration, success bool) {
	m.durations.Record(ctx, d.Milliseconds(), metric.WithAttributes(
		attribute.String("plugin", m.plugin),
		attribute.String("function", string(function)),
		attribute.String("success", strconv.FormatBool(success)),
		attribute.String("configDigest", m.configDigest),
	))
}

// TrackReports increments the reports processed counter
func (m *PluginMetrics) TrackReports(ctx context.Context, function FunctionType, count int, success bool) {
	m.reportsGenerated.Add(ctx, int64(count), metric.WithAttributes(
		attribute.String("plugin", m.plugin),
		attribute.String("function", string(function)),
		attribute.String("success", strconv.FormatBool(success)),
		attribute.String("configDigest", m.configDigest),
	))
}

// TrackSize records the size of data produced
func (m *PluginMetrics) TrackSize(ctx context.Context, function FunctionType, size int) {
	m.sizes.Record(ctx, int64(size), metric.WithAttributes(
		attribute.String("plugin", m.plugin),
		attribute.String("function", string(function)),
		attribute.String("configDigest", m.configDigest),
	))
}

// UpdateStatus updates the plugin status gauge (1 = up, 0 = down)
func (m *PluginMetrics) UpdateStatus(ctx context.Context, up bool) {
	val := int64(0)
	if up {
		val = 1
	}
	m.status.Record(ctx, val, metric.WithAttributes(
		attribute.String("plugin", m.plugin),
		attribute.String("configDigest", m.configDigest),
	))
}

// MetricViews returns histogram bucket definitions for the given metric prefix.
// Note: due to the OTEL specification, all histogram buckets must be defined when the beholder client is created.
func MetricViews(metricPrefix string) []sdkmetric.View {
	return []sdkmetric.View{
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: metricPrefix + "_duration_ms"},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				// 5, 10, 20, 40, 80, 160, 320, 640, 1280, 2560, 5120, 10240, 20480, 40960
				Boundaries: prometheus.ExponentialBuckets(5, 2, 14),
			}},
		),
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: metricPrefix + "_data_sizes"},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				// 1KB, 2KB, 4KB, 8KB, 16KB, 32KB, 64KB, 128KB, 256KB, 512KB, 1024KB, 2048KB, 4096KB, 8192KB
				Boundaries: prometheus.ExponentialBuckets(1024, 2, 14),
			}},
		),
	}
}
