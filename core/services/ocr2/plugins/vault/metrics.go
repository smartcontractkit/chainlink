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

	queueOverflow metric.Int64Counter
}

func newPluginMetrics(configDigest string) (*pluginMetrics, error) {
	queueOverflow, err := beholder.GetMeter().Int64Counter("platform_vault_plugin_queue_overflow")
	if err != nil {
		return nil, fmt.Errorf("failed to create queue overflow counter: %w", err)
	}

	return &pluginMetrics{
		configDigest:  configDigest,
		queueOverflow: queueOverflow,
	}, nil
}

func (m *pluginMetrics) trackQueueOverflow(ctx context.Context, queueSize int, batchSize int) {
	m.queueOverflow.Add(ctx, 1, metric.WithAttributes(
		attribute.String("configDigest", m.configDigest),
		attribute.Int("queueSize", queueSize),
		attribute.Int("batchSize", batchSize),
	))
}
