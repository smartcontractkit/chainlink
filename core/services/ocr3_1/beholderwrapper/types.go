package beholderwrapper

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

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
	sizes            metric.Int64Counter
	status           metric.Int64Gauge
}

func newPluginMetrics(plugin, configDigest string) (*pluginMetrics, error) {
	durations, err := beholder.GetMeter().Int64Histogram("platform_ocr3_1_reporting_plugin_duration_ms")
	if err != nil {
		return nil, fmt.Errorf("failed to create duration histogram: %w", err)
	}

	reportsGenerated, err := beholder.GetMeter().Int64Counter("platform_ocr3_1_reporting_plugin_reports_processed")
	if err != nil {
		return nil, fmt.Errorf("failed to create reports counter: %w", err)
	}

	sizes, err := beholder.GetMeter().Int64Counter("platform_ocr3_1_reporting_plugin_data_sizes")
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

func (m *pluginMetrics) trackReports(ctx context.Context, function functionType, count int) {
	m.reportsGenerated.Add(ctx, int64(count), metric.WithAttributes(
		attribute.String("plugin", m.plugin),
		attribute.String("function", string(function)),
		attribute.String("configDigest", m.configDigest),
	))
}

func (m *pluginMetrics) trackSize(ctx context.Context, function functionType, size int) {
	m.sizes.Add(ctx, int64(size), metric.WithAttributes(
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
