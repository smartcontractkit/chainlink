package targets

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/metrics"

	localMonitoring "github.com/smartcontractkit/chainlink/v2/core/monitoring"
)

type writeTargetMetricsLabeler struct {
	metrics.Labeler
	chainWriterFailureCount metric.Int64Counter
}

func newWriteTargetMetricsLabeler(labeler metrics.Labeler) (*writeTargetMetricsLabeler, error) {
	chainWriterFailureCount, err := beholder.GetMeter().Int64Counter("write_target_failures_count")
	if err != nil {
		return nil, fmt.Errorf("failed to register write target failure counter: %w", err)
	}
	return &writeTargetMetricsLabeler{
		Labeler:                 labeler,
		chainWriterFailureCount: chainWriterFailureCount,
	}, nil
}

func (l *writeTargetMetricsLabeler) with(keyValues ...string) *writeTargetMetricsLabeler {
	return &writeTargetMetricsLabeler{l.With(keyValues...), l.chainWriterFailureCount}
}

func (l *writeTargetMetricsLabeler) incrementChainWriterFailureCount(ctx context.Context) {
	otelLabels := localMonitoring.KvMapToOtelAttributes(l.Labels)
	l.chainWriterFailureCount.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}
