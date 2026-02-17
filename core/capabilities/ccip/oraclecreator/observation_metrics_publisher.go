package oraclecreator

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// BeholderMetricsPublisher publishes observation metrics to Beholder
type BeholderMetricsPublisher struct {
	logger                     logger.Logger
	sentObservationsMetric     metric.Int64Counter
	includedObservationsMetric metric.Int64Counter
}

// NewBeholderMetricsPublisher creates a new Beholder-based metrics publisher for OCR3 observation metrics
func NewBeholderMetricsPublisher(logger logger.Logger, packageName string) (*BeholderMetricsPublisher, error) {
	bhClient := beholder.GetClient().ForPackage(packageName)

	sentObservationsMetric, err := bhClient.Meter.Int64Counter("ocr3_sent_observations_total")
	if err != nil {
		return nil, fmt.Errorf("failed to register ocr3_sent_observations_total counter: %w", err)
	}

	includedObservationsMetric, err := bhClient.Meter.Int64Counter("ocr3_included_observations_total")
	if err != nil {
		return nil, fmt.Errorf("failed to register ocr3_included_observations_total counter: %w", err)
	}

	return &BeholderMetricsPublisher{
		logger:                     logger,
		sentObservationsMetric:     sentObservationsMetric,
		includedObservationsMetric: includedObservationsMetric,
	}, nil
}

// PublishMetric publishes a metric to Beholder
func (p *BeholderMetricsPublisher) PublishMetric(ctx context.Context, metricName string, value float64, labels map[string]string) {
	// Convert labels to OTEL attributes
	attrs := make([]attribute.KeyValue, 0, len(labels))
	for k, v := range labels {
		attrs = append(attrs, attribute.String(k, v))
	}

	p.logger.Debugw("Publishing observation metric to Beholder",
		"metric", metricName,
		"value", value,
		"labels", labels,
	)

	// Publish the delta to Beholder (Beholder counters are cumulative)
	// The value parameter already contains just the delta (increment amount), not the cumulative total
	switch metricName {
	case "ocr3_sent_observations_total":
		p.sentObservationsMetric.Add(ctx, int64(value), metric.WithAttributes(attrs...))
	case "ocr3_included_observations_total":
		p.includedObservationsMetric.Add(ctx, int64(value), metric.WithAttributes(attrs...))
	default:
		p.logger.Warnw("Unknown observation metric", "metric", metricName)
	}
}
