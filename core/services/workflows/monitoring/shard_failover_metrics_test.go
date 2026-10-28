package monitoring

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

// newTestShardFailoverMetrics swaps the beholder client for a manual reader
// and returns a fresh metrics instance bound to it, so each test asserts
// absolute values free of the singleton's cross-test state.
func newTestShardFailoverMetrics(t *testing.T) (*sdkmetric.ManualReader, *ShardFailoverMetrics) {
	t.Helper()
	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	prev := beholder.GetClient()
	t.Cleanup(func() { beholder.SetClient(prev) })
	client := beholder.NoopClientConfig{}.New()
	client.Meter = mp.Meter("test")
	client.MeterProvider = mp
	beholder.SetClient(client)
	fm, err := newShardFailoverMetrics()
	require.NoError(t, err)
	return reader, fm
}

// cachedEventsValue returns the current value of the node-wide cached-events
// up/down counter (a non-monotonic sum in metricdata).
func cachedEventsValue(t *testing.T, reader *sdkmetric.ManualReader) int64 {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name != "platform_shard_failover_cached_events" {
				continue
			}
			if s, ok := m.Data.(metricdata.Sum[int64]); ok {
				for _, dp := range s.DataPoints {
					return dp.Value
				}
			}
		}
	}
	return 0
}

//nolint:paralleltest // swaps the global beholder client behind the metrics singleton
func TestShardFailoverMetrics_Emissions(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	prev := beholder.GetClient()
	t.Cleanup(func() { beholder.SetClient(prev) })
	client := beholder.NoopClientConfig{}.New()
	client.Meter = mp.Meter("test")
	client.MeterProvider = mp
	beholder.SetClient(client)

	fm, err := GlobalShardFailoverMetrics()
	require.NoError(t, err)
	require.NotNil(t, fm)

	ctx := t.Context()
	fm.SetFailoverEnabled(ctx, true)
	fm.RecordTriggerCached(ctx)
	fm.RecordTriggerCached(ctx)
	fm.RecordCacheExit(ctx, ShardFailoverCacheExitExpired, 2)
	fm.RecordTriggerCached(ctx)
	fm.RecordCacheExit(ctx, ShardFailoverCacheExitClosed, 1)
	fm.RecordReplay(ctx, ShardFailoverReplaySucceeded)
	fm.RecordReplayDuration(ctx, ShardFailoverReplaySucceeded, 1500*time.Millisecond)
	fm.RecordReplayDuration(ctx, ShardFailoverReplaySucceeded, 400*time.Millisecond)

	var nilFM *ShardFailoverMetrics
	nilFM.SetFailoverEnabled(ctx, true)
	nilFM.RecordTriggerCached(ctx)
	nilFM.RecordCacheExit(ctx, ShardFailoverCacheExitDrained, 1)
	nilFM.RecordReplay(ctx, ShardFailoverReplayFailed)
	nilFM.RecordReplayDuration(ctx, ShardFailoverReplayFailed, time.Second)

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(ctx, &rm))

	gauges := make(map[string]int64)
	sums := make(map[string]map[string]int64)
	units := make(map[string]string)
	var replayHist *metricdata.HistogramDataPoint[float64]
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			units[m.Name] = m.Unit
			switch data := m.Data.(type) {
			case metricdata.Gauge[int64]:
				for _, dp := range data.DataPoints {
					gauges[m.Name] = dp.Value
				}
			case metricdata.Sum[int64]:
				for _, dp := range data.DataPoints {
					outcome := attrOutcome(dp.Attributes)
					if sums[m.Name] == nil {
						sums[m.Name] = make(map[string]int64)
					}
					sums[m.Name][outcome] += dp.Value
				}
			case metricdata.Histogram[float64]:
				for _, dp := range data.DataPoints {
					if attrOutcome(dp.Attributes) == ShardFailoverReplaySucceeded {
						replayHist = &dp
					}
				}
			}
		}
	}

	assert.Equal(t, int64(1), gauges["platform_shard_failover_enabled"])
	assert.Equal(t, int64(0), sums["platform_shard_failover_cached_events"][""], "node-wide total returns to 0 after every cached event exits")
	assert.Equal(t, int64(3), sums["platform_shard_failover_trigger_cached"][""])
	assert.Equal(t, int64(2), sums["platform_shard_failover_cache_exit"][ShardFailoverCacheExitExpired])
	assert.Equal(t, int64(1), sums["platform_shard_failover_cache_exit"][ShardFailoverCacheExitClosed])
	assert.Equal(t, int64(1), sums["platform_shard_failover_replay"][ShardFailoverReplaySucceeded])
	assert.Equal(t, "s", units["platform_shard_failover_replay_duration_seconds"])

	require.NotNil(t, replayHist, "replay duration histogram recorded")
	assert.Equal(t, uint64(2), replayHist.Count)
	assert.InDelta(t, 1.9, replayHist.Sum, 1e-9, "1500ms records as 1.5 and 400ms as 0.4, not truncated whole seconds")
	assert.Equal(t, []float64{0.25, 0.5, 1, 2, 5, 10, 30, 120, 300}, replayHist.Bounds)
	// 0.4 lands in the (0.25, 0.5] bucket (index 1) and 1.5 in the (1, 2]
	// bucket (index 3); the 10th entry is the implied (300, +inf) bucket.
	assert.Equal(t, []uint64{0, 1, 0, 1, 0, 0, 0, 0, 0, 0}, replayHist.BucketCounts)
}

// TestShardFailoverMetrics_CachedEventsAggregatesAcrossManagers proves the
// cached-events instrument is the node-wide total: two managers interleaving
// cache and exit calls leave it at the SUM of both caches, not the last
// writer's per-manager length, and the funnel invariant
// cached = drained + replayed + expired + closed + still-cached holds.
//
//nolint:paralleltest // swaps the global beholder client
func TestShardFailoverMetrics_CachedEventsAggregatesAcrossManagers(t *testing.T) {
	reader, fm := newTestShardFailoverMetrics(t)
	ctx := t.Context()

	// manager A caches one, manager B caches two, manager A caches one more
	fm.RecordTriggerCached(ctx)
	fm.RecordTriggerCached(ctx)
	fm.RecordTriggerCached(ctx)
	assert.Equal(t, int64(3), cachedEventsValue(t, reader), "up/down counter sums both managers' caches, not the last writer's length")

	// manager A replays one, manager B expires one
	fm.RecordCacheExit(ctx, ShardFailoverCacheExitReplayed, 1)
	fm.RecordCacheExit(ctx, ShardFailoverCacheExitExpired, 1)
	assert.Equal(t, int64(1), cachedEventsValue(t, reader), "up/down counter reads the sum after exits from both managers")

	// manager A is closed with its last event still cached
	fm.RecordCacheExit(ctx, ShardFailoverCacheExitClosed, 1)
	assert.Equal(t, int64(0), cachedEventsValue(t, reader), "closed exit returns the node-wide total to 0")

	stillCached := cachedEventsValue(t, reader)

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(ctx, &rm))
	sums := make(map[string]map[string]int64)
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if sum, ok := m.Data.(metricdata.Sum[int64]); ok {
				for _, dp := range sum.DataPoints {
					outcome := attrOutcome(dp.Attributes)
					if sums[m.Name] == nil {
						sums[m.Name] = make(map[string]int64)
					}
					sums[m.Name][outcome] += dp.Value
				}
			}
		}
	}
	cached := sums["platform_shard_failover_trigger_cached"][""]
	drained := sums["platform_shard_failover_cache_exit"][ShardFailoverCacheExitDrained]
	replayed := sums["platform_shard_failover_cache_exit"][ShardFailoverCacheExitReplayed]
	expired := sums["platform_shard_failover_cache_exit"][ShardFailoverCacheExitExpired]
	closed := sums["platform_shard_failover_cache_exit"][ShardFailoverCacheExitClosed]
	assert.Equal(t, int64(3), cached)
	assert.Equal(t, int64(1), replayed)
	assert.Equal(t, int64(1), expired)
	assert.Equal(t, int64(1), closed)
	assert.Equal(t, cached, drained+replayed+expired+closed+stillCached, "cache funnel invariant against the node-wide total")
}

// TestShardFailoverMetrics_ConcurrentManagers hammers the shared instruments
// from concurrent goroutines the way concurrent managers do: concurrent Add
// calls on the up/down counter must sum to the exact node-wide total (run
// with -race).
//
//nolint:paralleltest // swaps the global beholder client
func TestShardFailoverMetrics_ConcurrentManagers(t *testing.T) {
	reader, fm := newTestShardFailoverMetrics(t)
	ctx := t.Context()

	const managers = 8
	const perManager = 50
	var wg sync.WaitGroup
	for range managers {
		wg.Go(func() {
			for j := range perManager {
				fm.RecordTriggerCached(ctx)
				if j%2 == 0 {
					fm.RecordCacheExit(ctx, ShardFailoverCacheExitReplayed, 1)
				}
			}
		})
	}
	wg.Wait()

	assert.Equal(t, int64(200), cachedEventsValue(t, reader), "every manager caches 50 and exits 25: 8*(50-25)")
}

func attrOutcome(attrs attribute.Set) string {
	if v, ok := attrs.Value(attribute.Key("outcome")); ok {
		return v.AsString()
	}
	return ""
}
