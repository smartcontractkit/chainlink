package vault

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

type pluginMetrics struct {
	configDigest string

	queueOverflow       metric.Int64Counter
	kvOperationDuration metric.Float64Histogram
}

func newPluginMetrics(configDigest string) (*pluginMetrics, error) {
	queueOverflow, err := beholder.GetMeter().Int64Counter("platform_vault_plugin_queue_overflow")
	if err != nil {
		return nil, fmt.Errorf("failed to create queue overflow counter: %w", err)
	}

	kvOperationDuration, err := beholder.GetMeter().Float64Histogram(
		"platform_vault_plugin_kv_operation_duration_seconds",
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create kv operation duration histogram: %w", err)
	}

	return &pluginMetrics{
		configDigest:        configDigest,
		queueOverflow:       queueOverflow,
		kvOperationDuration: kvOperationDuration,
	}, nil
}

func (m *pluginMetrics) trackKVOperation(ctx context.Context, method string, durationSeconds float64) {
	m.kvOperationDuration.Record(ctx, durationSeconds, metric.WithAttributes(
		attribute.String("configDigest", m.configDigest),
		attribute.String("method", method),
	))
}

func (m *pluginMetrics) trackQueueOverflow(ctx context.Context, queueSize int, batchSize int) {
	m.queueOverflow.Add(ctx, 1, metric.WithAttributes(
		attribute.String("configDigest", m.configDigest),
		attribute.Int("queueSize", queueSize),
		attribute.Int("batchSize", batchSize),
	))
}
