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

	queueOverflow                   metric.Int64Counter
	kvOperationDuration             metric.Int64Histogram
	localQueueSize                  metric.Int64Histogram
	observationPendingPackedItems   metric.Int64Histogram
	pendingQueueWrittenSize         metric.Int64Histogram
	observationPrefixCoverage       metric.Int64Histogram
	observationPrefixCoverageSpread metric.Int64Histogram
	pendingQueueStallSignals        metric.Int64Counter
	pendingQueuePurges              metric.Int64Counter
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

	observationPrefixCoverage, err := beholder.GetMeter().Int64Histogram(
		"platform_vault_plugin_observation_prefix_coverage",
		metric.WithUnit("{request}"),
		metric.WithDescription("Per-observer count of pending-queue items an oracle contributed an Ok observation for. Low outliers indicate a node withholding or truncating its observation prefix."),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create observation prefix coverage histogram: %w", err)
	}

	observationPrefixCoverageSpread, err := beholder.GetMeter().Int64Histogram(
		"platform_vault_plugin_observation_prefix_coverage_spread",
		metric.WithUnit("{request}"),
		metric.WithDescription("Max minus min per-observer Ok prefix coverage across oracles for one round. High spread indicates prefix divergence stalling head-of-queue quorum."),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create observation prefix coverage spread histogram: %w", err)
	}

	pendingQueueStallSignals, err := beholder.GetMeter().Int64Counter(
		"platform_vault_plugin_pending_queue_stall_signal",
		metric.WithUnit("{observation}"),
		metric.WithDescription("Count of Vault observations that signal a pending queue stall."),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create pending queue stall signal counter: %w", err)
	}

	pendingQueuePurges, err := beholder.GetMeter().Int64Counter(
		"platform_vault_plugin_pending_queue_purge",
		metric.WithUnit("{purge}"),
		metric.WithDescription("Count of replicated Vault pending queue purges committed after F+1 stall signals."),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create pending queue purge counter: %w", err)
	}

	return &pluginMetrics{
		configDigest:                    configDigest,
		queueOverflow:                   queueOverflow,
		kvOperationDuration:             kvOperationDuration,
		localQueueSize:                  localQueueSize,
		observationPendingPackedItems:   observationPendingPackedItems,
		pendingQueueWrittenSize:         pendingQueueWrittenSize,
		observationPrefixCoverage:       observationPrefixCoverage,
		observationPrefixCoverageSpread: observationPrefixCoverageSpread,
		pendingQueueStallSignals:        pendingQueueStallSignals,
		pendingQueuePurges:              pendingQueuePurges,
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

func (m *pluginMetrics) trackObservationPrefixCoverage(ctx context.Context, observer uint8, coverage int) {
	if m == nil {
		return
	}
	m.observationPrefixCoverage.Record(ctx, int64(coverage), metric.WithAttributes(
		attribute.String("configDigest", m.configDigest),
		attribute.Int("observer", int(observer)),
	))
}

func (m *pluginMetrics) trackObservationPrefixCoverageSpread(ctx context.Context, spread int) {
	if m == nil {
		return
	}
	m.observationPrefixCoverageSpread.Record(ctx, int64(spread), metric.WithAttributes(
		attribute.String("configDigest", m.configDigest),
	))
}

func (m *pluginMetrics) trackPendingQueueStallSignal(ctx context.Context) {
	if m == nil {
		return
	}
	m.pendingQueueStallSignals.Add(ctx, 1, metric.WithAttributes(
		attribute.String("configDigest", m.configDigest),
	))
}

func (m *pluginMetrics) trackPendingQueuePurge(ctx context.Context) {
	if m == nil {
		return
	}
	m.pendingQueuePurges.Add(ctx, 1, metric.WithAttributes(
		attribute.String("configDigest", m.configDigest),
	))
}
