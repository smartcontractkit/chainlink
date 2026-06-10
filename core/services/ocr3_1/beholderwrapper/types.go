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

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3/beholderwrapper/metrics"
)

// MetricPrefix is the prefix for all OCR3.1 beholder metrics
const MetricPrefix = "platform_ocr3_1_reporting_plugin"

// pluginMetrics extends the shared PluginMetrics with OCR3.1-specific metrics
type pluginMetrics struct {
	*metrics.PluginMetrics

	plugin       string
	configDigest string

	// OCR3.1 specific metrics for blob and KV operations
	blobDurations metric.Int64Histogram
	kvDurations   metric.Int64Histogram
}

func newPluginMetrics(plugin, configDigest string) (*pluginMetrics, error) {
	// Create base metrics using shared package
	base, err := metrics.NewPluginMetrics(MetricPrefix, plugin, configDigest)
	if err != nil {
		return nil, err
	}

	// Create OCR3.1-specific metrics
	blobDurations, err := beholder.GetMeter().Int64Histogram(MetricPrefix+"_blob_duration_ms", metric.WithUnit("ms"))
	if err != nil {
		return nil, fmt.Errorf("failed to create blob duration histogram: %w", err)
	}

	kvDurations, err := beholder.GetMeter().Int64Histogram(MetricPrefix+"_kv_duration_ms", metric.WithUnit("ms"))
	if err != nil {
		return nil, fmt.Errorf("failed to create kv duration histogram: %w", err)
	}

	return &pluginMetrics{
		PluginMetrics: base,
		plugin:        plugin,
		configDigest:  configDigest,
		blobDurations: blobDurations,
		kvDurations:   kvDurations,
	}, nil
}

func (m *pluginMetrics) recordKVDuration(ctx context.Context, method string, d time.Duration, success bool) {
	m.kvDurations.Record(ctx, d.Milliseconds(), metric.WithAttributes(
		attribute.String("plugin", m.plugin),
		attribute.String("method", method),
		attribute.String("success", strconv.FormatBool(success)),
		attribute.String("configDigest", m.configDigest),
	))
}

func (m *pluginMetrics) recordBlobDuration(ctx context.Context, method string, d time.Duration, success bool) {
	m.blobDurations.Record(ctx, d.Milliseconds(), metric.WithAttributes(
		attribute.String("plugin", m.plugin),
		attribute.String("method", method),
		attribute.String("success", strconv.FormatBool(success)),
		attribute.String("configDigest", m.configDigest),
	))
}

// MetricViews returns histogram bucket definitions for OCR3.1 metrics
// Note: due to the OTEL specification, all histogram buckets must be defined when the beholder client is created
func MetricViews() []sdkmetric.View {
	// Get base views from shared package
	baseViews := metrics.MetricViews(MetricPrefix)

	// Add OCR3.1-specific views
	ocr31Views := []sdkmetric.View{
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: MetricPrefix + "_kv_duration_ms"},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				// 5, 10, 20, 40, 80, 160, 320, 640, 1280, 2560, 5120, 10240, 20480, 40960
				Boundaries: prometheus.ExponentialBuckets(5, 2, 14),
			}},
		),
		sdkmetric.NewView(
			sdkmetric.Instrument{Name: MetricPrefix + "_blob_duration_ms"},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				// 5, 10, 20, 40, 80, 160, 320, 640, 1280, 2560, 5120, 10240, 20480, 40960
				Boundaries: prometheus.ExponentialBuckets(5, 2, 14),
			}},
		),
	}

	return append(baseViews, ocr31Views...)
}
