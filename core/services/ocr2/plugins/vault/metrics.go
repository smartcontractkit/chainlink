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

	queueOverflow                 metric.Int64Counter
	kvOperationDuration           metric.Int64Histogram
	localQueueSize                metric.Int64Histogram
	observationPendingPackedItems metric.Int64Histogram
	pendingQueueWrittenSize       metric.Int64Histogram
}

func newPluginMetrics(configDigest string) (*pluginMetrics, error) {
	queueOverflow, err := beholder.GetMeter().Int64Counter("platform_vault_plugin_queue_overflow")
	if err != nil {
		return nil, fmt.Errorf("failed to create queue overflow counter: %w", err)
	}

	kvOperationDuration, err := beholder.GetMeter().Int64Histogram(
		"platform_vault_plugin_kv_operation_duration_ms",
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create kv operation duration histogram: %w", err)
	}

	localQueueSize, err := beholder.GetMeter().Int64Histogram(
		"platform_vault_plugin_local_queue_size",
		metric.WithUnit("{request}"),
		metric.WithDescription("Count of items in the Vault reporting plugin local request store at Observation time"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create local queue size histogram: %w", err)
	}

	observationPendingPackedItems, err := beholder.GetMeter().Int64Histogram(
		"platform_vault_plugin_observation_pending_packed_items",
		metric.WithUnit("{request}"),
		metric.WithDescription("Count of local-queue requests packed into pending-queue blobs in one Observation (after dedupe against KV pending queue)."),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create observation pending packed items histogram: %w", err)
	}

	pendingQueueWrittenSize, err := beholder.GetMeter().Int64Histogram(
		"platform_vault_plugin_pending_queue_written_size",
		metric.WithUnit("{request}"),
		metric.WithDescription("Items written to the KV pending queue after F+1 consensus aggregation."),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create pending queue written size histogram: %w", err)
	}

	return &pluginMetrics{
		configDigest:                  configDigest,
		queueOverflow:                 queueOverflow,
		kvOperationDuration:           kvOperationDuration,
		localQueueSize:                localQueueSize,
		observationPendingPackedItems: observationPendingPackedItems,
		pendingQueueWrittenSize:       pendingQueueWrittenSize,
	}, nil
}

func (m *pluginMetrics) trackKVOperation(ctx context.Context, method string, durationMs int64) {
	if m == nil {
		return
	}
	m.kvOperationDuration.Record(ctx, durationMs, metric.WithAttributes(
		attribute.String("configDigest", m.configDigest),
		attribute.String("method", method),
	))
}

func (m *pluginMetrics) trackQueueOverflow(ctx context.Context, queueSize int, batchSize int) {
	if m == nil {
		return
	}
	m.queueOverflow.Add(ctx, 1, metric.WithAttributes(
		attribute.String("configDigest", m.configDigest),
		attribute.Int("queueSize", queueSize),
		attribute.Int("batchSize", batchSize),
	))
}

func (m *pluginMetrics) trackLocalQueueSize(ctx context.Context, size int) {
	if m == nil {
		return
	}
	m.localQueueSize.Record(ctx, int64(size), metric.WithAttributes(
		attribute.String("configDigest", m.configDigest),
	))
}

func (m *pluginMetrics) trackObservationPendingPack(ctx context.Context, packedItemCount, blobHandleCount int) {
	if m == nil {
		return
	}
	m.observationPendingPackedItems.Record(ctx, int64(packedItemCount), metric.WithAttributes(
		attribute.String("configDigest", m.configDigest),
		attribute.Int("blobHandleCount", blobHandleCount),
	))
}

func (m *pluginMetrics) trackPendingQueueWrittenSize(ctx context.Context, writtenCount int) {
	if m == nil {
		return
	}
	m.pendingQueueWrittenSize.Record(ctx, int64(writtenCount), metric.WithAttributes(
		attribute.String("configDigest", m.configDigest),
	))
}
