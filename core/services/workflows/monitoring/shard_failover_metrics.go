package monitoring

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

// Outcome values for platform_shard_failover_cache_exit (low cardinality).
const (
	ShardFailoverCacheExitDrained  = "drained"
	ShardFailoverCacheExitReplayed = "replayed"
	ShardFailoverCacheExitExpired  = "expired"
)

// Outcome values for platform_shard_failover_replay (low cardinality).
const (
	ShardFailoverReplaySucceeded = "succeeded"
	ShardFailoverReplayFailed    = "failed"
	ShardFailoverReplayMiss      = "miss"
)

// ShardFailoverMetrics are the process-global failover instruments, emitted by
// the ShardFailoverManager and the syncer through the shared unlabeled meter.
type ShardFailoverMetrics struct {
	enabledGauge          metric.Int64Gauge
	triggerCachedCounter  metric.Int64Counter
	cachedEventsGauge     metric.Int64Gauge
	cacheExitCounter      metric.Int64Counter
	replayCounter         metric.Int64Counter
	replayDurationSeconds metric.Int64Histogram
}

var (
	globalShardFailoverMetrics     *ShardFailoverMetrics
	errGlobalShardFailoverMetrics  error
	globalShardFailoverMetricsOnce sync.Once
)

func GlobalShardFailoverMetrics() (*ShardFailoverMetrics, error) {
	globalShardFailoverMetricsOnce.Do(func() {
		globalShardFailoverMetrics, errGlobalShardFailoverMetrics = newShardFailoverMetrics()
	})
	return globalShardFailoverMetrics, errGlobalShardFailoverMetrics
}

func newShardFailoverMetrics() (*ShardFailoverMetrics, error) {
	m := &ShardFailoverMetrics{}
	var err error
	m.enabledGauge, err = beholder.GetMeter().Int64Gauge("platform_shard_failover_enabled",
		metric.WithDescription("Whether sharding failover is enabled on this node (1) or not (0)"))
	if err != nil {
		return nil, err
	}
	m.triggerCachedCounter, err = beholder.GetMeter().Int64Counter("platform_shard_failover_trigger_cached",
		metric.WithDescription("Trigger events cached on a secondary shard for potential failover replay"))
	if err != nil {
		return nil, err
	}
	m.cachedEventsGauge, err = beholder.GetMeter().Int64Gauge("platform_shard_failover_cached_events",
		metric.WithDescription("Trigger events currently cached for failover replay"))
	if err != nil {
		return nil, err
	}
	m.cacheExitCounter, err = beholder.GetMeter().Int64Counter("platform_shard_failover_cache_exit",
		metric.WithDescription("Cached failover trigger events leaving the cache, by outcome: drained (verdict arrived), replayed (popped for replay), expired (pruned without a verdict)"))
	if err != nil {
		return nil, err
	}
	m.replayCounter, err = beholder.GetMeter().Int64Counter("platform_shard_failover_replay",
		metric.WithDescription("Failover replays of cached trigger events, by outcome: succeeded, failed, miss (no cached event for the triggerEventID)"))
	if err != nil {
		return nil, err
	}
	m.replayDurationSeconds, err = beholder.GetMeter().Int64Histogram("platform_shard_failover_replay_duration_seconds",
		metric.WithDescription("Duration of failover replays of cached trigger events"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.5, 1, 2, 5, 10, 30, 60, 120, 300))
	if err != nil {
		return nil, err
	}
	return m, nil
}

// SetFailoverEnabled records whether the sharding failover gate is open.
func (m *ShardFailoverMetrics) SetFailoverEnabled(ctx context.Context, enabled bool) {
	if m == nil {
		return
	}
	enabledVal := int64(0)
	if enabled {
		enabledVal = 1
	}
	m.enabledGauge.Record(ctx, enabledVal)
}

// RecordTriggerCached counts one trigger event cached for failover.
func (m *ShardFailoverMetrics) RecordTriggerCached(ctx context.Context) {
	if m == nil {
		return
	}
	m.triggerCachedCounter.Add(ctx, 1)
}

// UpdateCachedEvents records the current number of cached trigger events.
func (m *ShardFailoverMetrics) UpdateCachedEvents(ctx context.Context, n int64) {
	if m == nil {
		return
	}
	m.cachedEventsGauge.Record(ctx, n)
}

// RecordCacheExit counts n cached events leaving the cache. outcome must be a
// ShardFailoverCacheExit* constant.
func (m *ShardFailoverMetrics) RecordCacheExit(ctx context.Context, outcome string, n int64) {
	if m == nil {
		return
	}
	m.cacheExitCounter.Add(ctx, n, metric.WithAttributes(attribute.String("outcome", outcome)))
}

// RecordReplay counts one failover replay attempt. outcome must be a
// ShardFailoverReplay* constant.
func (m *ShardFailoverMetrics) RecordReplay(ctx context.Context, outcome string) {
	if m == nil {
		return
	}
	m.replayCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("outcome", outcome)))
}

// RecordReplayDuration records the duration of one failover replay. outcome
// must be a ShardFailoverReplay* constant.
func (m *ShardFailoverMetrics) RecordReplayDuration(ctx context.Context, outcome string, d time.Duration) {
	if m == nil {
		return
	}
	m.replayDurationSeconds.Record(ctx, int64(d/time.Second), metric.WithAttributes(attribute.String("outcome", outcome)))
}
