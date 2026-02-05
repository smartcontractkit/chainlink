package capregconfig

import (
	"context"
	"fmt"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

// Metrics holds all the metric instruments for OCRConfigService.
type Metrics struct {
	configUpdatesCounter        metric.Int64Counter
	parseErrorsCounter          metric.Int64Counter
	configCountGauge            metric.Int64Gauge
	capabilityConfigErrors      metric.Int64Counter
	trackerLegacyFallbackGauge  metric.Int64Gauge
	digesterLegacyFallbackGauge metric.Int64Gauge
}

// InitMetrics initializes the Beholder metrics for OCRConfigService.
func InitMetrics() (*Metrics, error) {
	m := &Metrics{}
	var err error

	m.configUpdatesCounter, err = beholder.GetMeter().Int64Counter(
		"platform_ocr_config_service_updates_total",
		metric.WithDescription("Total OCR config updates received from registry"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register config updates counter: %w", err)
	}

	m.parseErrorsCounter, err = beholder.GetMeter().Int64Counter(
		"platform_ocr_config_service_parse_errors_total",
		metric.WithDescription("Total OCR config parse errors"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register parse errors counter: %w", err)
	}

	m.configCountGauge, err = beholder.GetMeter().Int64Gauge(
		"platform_ocr_config_service_config_count",
		metric.WithDescription("Current config count per capability/DON/key"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register config count gauge: %w", err)
	}

	m.capabilityConfigErrors, err = beholder.GetMeter().Int64Counter(
		"platform_ocr_config_service_capability_config_errors_total",
		metric.WithDescription("Total capability config processing errors (unmarshal failures)"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register capability config errors counter: %w", err)
	}

	m.trackerLegacyFallbackGauge, err = beholder.GetMeter().Int64Gauge(
		"platform_ocr_config_service_tracker_legacy_fallback",
		metric.WithDescription("Indicates whether the config tracker is using legacy fallback (1) or registry config (0)"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register tracker legacy fallback gauge: %w", err)
	}

	m.digesterLegacyFallbackGauge, err = beholder.GetMeter().Int64Gauge(
		"platform_ocr_config_service_digester_legacy_fallback",
		metric.WithDescription("Indicates whether the config digester is using legacy fallback (1) or registry config (0)"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register digester legacy fallback gauge: %w", err)
	}

	return m, nil
}

func (m *Metrics) IncrementConfigUpdates(ctx context.Context, capabilityID string, donID uint32, ocrConfigKey string) {
	if m == nil {
		return
	}
	m.configUpdatesCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("capability_id", capabilityID),
		attribute.String("don_id", strconv.FormatUint(uint64(donID), 10)),
		attribute.String("ocr_config_key", ocrConfigKey),
	))
}

func (m *Metrics) IncrementParseErrors(ctx context.Context, capabilityID string, donID uint32, ocrConfigKey string) {
	if m == nil {
		return
	}
	m.parseErrorsCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("capability_id", capabilityID),
		attribute.String("don_id", strconv.FormatUint(uint64(donID), 10)),
		attribute.String("ocr_config_key", ocrConfigKey),
	))
}

func (m *Metrics) SetConfigCount(ctx context.Context, capabilityID string, donID uint32, ocrConfigKey string, count int64) {
	if m == nil {
		return
	}
	m.configCountGauge.Record(ctx, count, metric.WithAttributes(
		attribute.String("capability_id", capabilityID),
		attribute.String("don_id", strconv.FormatUint(uint64(donID), 10)),
		attribute.String("ocr_config_key", ocrConfigKey),
	))
}

func (m *Metrics) IncrementCapabilityConfigErrors(ctx context.Context, capabilityID string, donID uint32) {
	if m == nil {
		return
	}
	m.capabilityConfigErrors.Add(ctx, 1, metric.WithAttributes(
		attribute.String("capability_id", capabilityID),
		attribute.String("don_id", strconv.FormatUint(uint64(donID), 10)),
	))
}

func (m *Metrics) SetTrackerLegacyFallback(ctx context.Context, capabilityID string, ocrConfigKey string, isLegacy bool) {
	if m == nil {
		return
	}
	var value int64
	if isLegacy {
		value = 1
	}
	m.trackerLegacyFallbackGauge.Record(ctx, value, metric.WithAttributes(
		attribute.String("capability_id", capabilityID),
		attribute.String("ocr_config_key", ocrConfigKey),
	))
}

func (m *Metrics) SetDigesterLegacyFallback(ctx context.Context, capabilityID string, ocrConfigKey string, isLegacy bool) {
	if m == nil {
		return
	}
	var value int64
	if isLegacy {
		value = 1
	}
	m.digesterLegacyFallbackGauge.Record(ctx, value, metric.WithAttributes(
		attribute.String("capability_id", capabilityID),
		attribute.String("ocr_config_key", ocrConfigKey),
	))
}
