package beholderwrapper

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

type functionType string

const (
	query               functionType = "query"
	observation         functionType = "observation"
	validateObservation functionType = "validateObservation"
	observationQuorum   functionType = "observationQuorum"
	stateTransition     functionType = "stateTransition"
	committed           functionType = "committed"
	reports             functionType = "reports"
	shouldAccept        functionType = "shouldAccept"
	shouldTransmit      functionType = "shouldTransmit"
)

type pluginMetrics struct {
	plugin       string
	configDigest string

	durations        metric.Int64Histogram
	reportsGenerated metric.Int64Counter
	sizes            metric.Int64Histogram
	status           metric.Int64Gauge
}

func newPluginMetrics(plugin, configDigest string) (*pluginMetrics, error) {
	durations, err := beholder.GetMeter().Int64Histogram("platform_ocr3_1_reporting_plugin_duration_ms", metric.WithUnit("ms"))
	if err != nil {
		return nil, fmt.Errorf("failed to create duration histogram: %w", err)
	}

	reportsGenerated, err := beholder.GetMeter().Int64Counter("platform_ocr3_1_reporting_plugin_reports_processed", metric.WithUnit("1"))
	if err != nil {
		return nil, fmt.Errorf("failed to create reports counter: %w", err)
	}

	sizes, err := beholder.GetMeter().Int64Histogram("platform_ocr3_1_reporting_plugin_data_sizes", metric.WithUnit("By"))
	if err != nil {
		return nil, fmt.Errorf("failed to create sizes counter: %w", err)
	}

	status, err := beholder.GetMeter().Int64Gauge("platform_ocr3_1_reporting_plugin_status")
	if err != nil {
		return nil, fmt.Errorf("failed to create status gauge: %w", err)
	}

	return &pluginMetrics{
		plugin:           plugin,
		configDigest:     configDigest,
		durations:        durations,
		reportsGenerated: reportsGenerated,
		sizes:            sizes,
		status:           status,
	}, nil
}

func (m *pluginMetrics) recordDuration(ctx context.Context, function functionType, d time.Duration, success bool) {
	m.durations.Record(ctx, d.Milliseconds(), metric.WithAttributes(
		attribute.String("plugin", m.plugin),
		attribute.String("function", string(function)),
		attribute.String("success", strconv.FormatBool(success)),
		attribute.String("configDigest", m.configDigest),
	))
}

func (m *pluginMetrics) trackReports(ctx context.Context, function functionType, count int, success bool) {
	m.reportsGenerated.Add(ctx, int64(count), metric.WithAttributes(
		attribute.String("plugin", m.plugin),
		attribute.String("function", string(function)),
		attribute.String("success", strconv.FormatBool(success)),
		attribute.String("configDigest", m.configDigest),
	))
}

func (m *pluginMetrics) trackSize(ctx context.Context, function functionType, size int) {
	m.sizes.Record(ctx, int64(size), metric.WithAttributes(
		attribute.String("plugin", m.plugin),
		attribute.String("function", string(function)),
		attribute.String("configDigest", m.configDigest),
	))
}

func (m *pluginMetrics) updateStatus(ctx context.Context, up bool) {
	val := int64(0)
	if up {
		val = 1
	}
	m.status.Record(ctx, val, metric.WithAttributes(
		attribute.String("plugin", m.plugin),
		attribute.String("configDigest", m.configDigest),
	))
}

// Note: due to the OTEL specification, all histogram buckets
// Must be defined when the beholder client is created
func MetricViews() []sdkmetric.View {
	return []sdkmetric.View{
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: "platform_ocr3_1_reporting_plugin_duration_ms"},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				// 5, 10, 20, 40, 80, 160, 320, 640, 1280, 2560
				Boundaries: prometheus.ExponentialBuckets(5, 2, 10),
			}},
		),
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: "platform_ocr3_1_reporting_plugin_data_sizes"},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				// 512KB is the max value possible
				// 1KB, 2KB, 4KB, 8KB, 16KB, 32KB, 64KB, 128KB, 256KB, 512KB
				Boundaries: prometheus.ExponentialBuckets(1024, 2, 10),
			}},
		),
	}
}
