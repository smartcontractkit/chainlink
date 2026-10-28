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
	ShardFailoverCacheExitClosed   = "closed"
)

// Outcome values for platform_shard_failover_replay (low cardinality).
const (
	ShardFailoverReplaySucceeded = "succeeded"
	ShardFailoverReplayFailed    = "failed"
	ShardFailoverReplayMiss      = "miss"
)

// ShardFailoverMetrics are the process-global failover instruments, emitted by
// the per-workflow ShardFailoverManager instances and the syncer through the
// shared unlabeled meter. The instruments are per-node aggregates with no
// per-workflow labels, so all managers on this node emit through this one
// singleton: cachedEvents is an up/down counter whose value is the node-wide
// number of currently cached events (OTel translates it to a Prometheus
// gauge) — every manager's +1 on cache and -n on exit sums into it, whereas a
// plain gauge would be last-writer-wins and hide the other managers' events.
type ShardFailoverMetrics struct {
	enabledGauge          metric.Int64Gauge
	triggerCachedCounter  metric.Int64Counter
	cachedEvents          metric.Int64UpDownCounter
	cacheExitCounter      metric.Int64Counter
	replayCounter         metric.Int64Counter
	replayDurationSeconds metric.Float64Histogram
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
	m.cachedEvents, err = beholder.GetMeter().Int64UpDownCounter("platform_shard_failover_cached_events",
		metric.WithDescription("Trigger events currently cached for failover replay, node-wide across all failover managers on this node"))
	if err != nil {
		return nil, err
	}
	m.cacheExitCounter, err = beholder.GetMeter().Int64Counter("platform_shard_failover_cache_exit",
		metric.WithDescription("Cached failover trigger events leaving the cache, by outcome: drained (verdict arrived), replayed (popped for replay), expired (pruned without a verdict), closed (manager shut down with events still cached)"))
	if err != nil {
		return nil, err
	}
	m.replayCounter, err = beholder.GetMeter().Int64Counter("platform_shard_failover_replay",
		metric.WithDescription("Failover replays of cached trigger events, by outcome: succeeded, failed, miss (no cached event for the triggerEventID)"))
	if err != nil {
		return nil, err
	}
	m.replayDurationSeconds, err = beholder.GetMeter().Float64Histogram("platform_shard_failover_replay_duration_seconds",
		metric.WithDescription("Duration of failover replays of cached trigger events"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.25, 0.5, 1, 2, 5, 10, 30, 120, 300))
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

// RecordTriggerCached counts one trigger event cached for failover and bumps
// the node-wide cached-events total.
func (m *ShardFailoverMetrics) RecordTriggerCached(ctx context.Context) {
	if m == nil {
		return
	}
	m.triggerCachedCounter.Add(ctx, 1)
	m.cachedEvents.Add(ctx, 1)
}

// RecordCacheExit counts n cached events leaving the cache and lowers the
// node-wide cached-events total by n. outcome must be a
// ShardFailoverCacheExit* constant.
func (m *ShardFailoverMetrics) RecordCacheExit(ctx context.Context, outcome string, n int64) {
	if m == nil {
		return
	}
	m.cacheExitCounter.Add(ctx, n, metric.WithAttributes(attribute.String("outcome", outcome)))
	m.cachedEvents.Add(ctx, -n)
}

// RecordReplay counts one failover replay attempt. outcome must be a
// ShardFailoverReplay* constant.
func (m *ShardFailoverMetrics) RecordReplay(ctx context.Context, outcome string) {
	if m == nil {
		return
	}
	m.replayCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("outcome", outcome)))
}

// RecordReplayDuration records the duration of one failover replay, in
// fractional seconds. outcome must be a ShardFailoverReplay* constant.
func (m *ShardFailoverMetrics) RecordReplayDuration(ctx context.Context, outcome string, d time.Duration) {
	if m == nil {
		return
	}
	m.replayDurationSeconds.Record(ctx, d.Seconds(), metric.WithAttributes(attribute.String("outcome", outcome)))
}
