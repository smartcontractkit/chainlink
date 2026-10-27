package v2

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/shardownership"
	wfTypes "github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
)

// shardFailoverReader receives every metric emitted through beholder.GetMeter()
// for the duration of the test binary. TestMain installs it before any test
// runs, so the GlobalShardFailoverMetrics singleton binds to it regardless of
// which test emits first, and these tests assert on before/after deltas
// because the counters are process-global.
var shardFailoverReader *sdkmetric.ManualReader

func TestMain(m *testing.M) {
	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	client := beholder.NoopClientConfig{}.New()
	client.Meter = mp.Meter("shard-failover-test")
	client.MeterProvider = mp
	beholder.SetClient(client)
	shardFailoverReader = reader
	code := m.Run()
	_ = mp.Shutdown(context.Background())
	os.Exit(code)
}

type failoverTotals struct {
	triggerCached int64
	cacheExit     map[string]int64
	replay        map[string]int64
	replayCount   map[string]uint64
}

func collectFailoverTotals(t *testing.T) failoverTotals {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, shardFailoverReader.Collect(t.Context(), &rm))
	totals := failoverTotals{cacheExit: map[string]int64{}, replay: map[string]int64{}, replayCount: map[string]uint64{}}
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			switch m.Name {
			case "platform_shard_failover_trigger_cached":
				if sum, ok := m.Data.(metricdata.Sum[int64]); ok {
					for _, dp := range sum.DataPoints {
						totals.triggerCached += dp.Value
					}
				}
			case "platform_shard_failover_cache_exit":
				if sum, ok := m.Data.(metricdata.Sum[int64]); ok {
					for _, dp := range sum.DataPoints {
						totals.cacheExit[failoverOutcome(dp.Attributes)] += dp.Value
					}
				}
			case "platform_shard_failover_replay":
				if sum, ok := m.Data.(metricdata.Sum[int64]); ok {
					for _, dp := range sum.DataPoints {
						totals.replay[failoverOutcome(dp.Attributes)] += dp.Value
					}
				}
			case "platform_shard_failover_replay_duration_seconds":
				if h, ok := m.Data.(metricdata.Histogram[int64]); ok {
					for _, dp := range h.DataPoints {
						totals.replayCount[failoverOutcome(dp.Attributes)] += dp.Count
					}
				}
			}
		}
	}
	return totals
}

// failoverGauge returns the last recorded value of a failover gauge.
func failoverGauge(t *testing.T, name string) (int64, bool) {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, shardFailoverReader.Collect(t.Context(), &rm))
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name != name {
				continue
			}
			g, ok := m.Data.(metricdata.Gauge[int64])
			if !ok {
				return 0, false
			}
			for _, dp := range g.DataPoints {
				return dp.Value, true
			}
		}
	}
	return 0, false
}

func failoverOutcome(attrs attribute.Set) string {
	if v, ok := attrs.Value(attribute.Key("outcome")); ok {
		return v.AsString()
	}
	return ""
}

// newFailoverTestManager wraps a recording engine in a secondary-shard manager
// whose failover gate is open, so trigger events are cached.
func newFailoverTestManager(t *testing.T, engine *recordingEngine) *ShardFailoverManager {
	t.Helper()
	return newAdmissionTestManager(t, engine, ShardFailoverManagerConfig{
		ShardingEnabled: true,
		ShardResolver:   secondaryResolver,
		Failover:        FailoverConfig{ShardingFailoverGate: limits.NewGateLimiter(true)},
	})
}

// TestShardFailoverManager_CacheExitFunnel proves the cache-exit funnel
// invariant: cached = drained + replayed + expired + still-cached. Every
// cached event leaves exactly one counted exit, a verdict for an event that
// was never cached is not an exit, and replays record their outcome and
// duration.
//
//nolint:paralleltest // asserts deltas on process-global failover metrics
func TestShardFailoverManager_CacheExitFunnel(t *testing.T) {
	ctx := t.Context()
	m := newFailoverTestManager(t, &recordingEngine{})
	failing := newFailoverTestManager(t, &recordingEngine{execErr: errors.New("engine boom")})

	before := collectFailoverTotals(t)

	for _, id := range []string{"evt-drain", "evt-replay", "evt-expire", "evt-keep"} {
		require.ErrorIs(t, m.ExecuteTrigger(ctx, testCoordinatedEvent(id)), triggers.ErrAdmissionCached)
	}
	require.ErrorIs(t, failing.ExecuteTrigger(ctx, testCoordinatedEvent("evt-fail")), triggers.ErrAdmissionCached)

	// A SUCCESS verdict drains one cached event; a verdict for an event that was
	// never cached must not count as an exit.
	m.HandleExecutionStatusUpdate(&ringpb.ExecutionStatusUpdate{TriggerEventId: "evt-drain", Status: ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS})
	m.HandleExecutionStatusUpdate(&ringpb.ExecutionStatusUpdate{TriggerEventId: "evt-uncached", Status: ringpb.ExecutionStatus_EXECUTION_STATUS_USER_ERROR})

	// A SYSTEM_ERROR verdict replays the cached event, succeeding on one manager
	// and failing on the other.
	m.HandleExecutionStatusUpdate(&ringpb.ExecutionStatusUpdate{TriggerEventId: "evt-replay", Status: ringpb.ExecutionStatus_EXECUTION_STATUS_SYSTEM_ERROR})
	failing.HandleExecutionStatusUpdate(&ringpb.ExecutionStatusUpdate{TriggerEventId: "evt-fail", Status: ringpb.ExecutionStatus_EXECUTION_STATUS_SYSTEM_ERROR})

	// The 10-minute prune expires the event whose verdict never arrived.
	m.mu.Lock()
	ce := m.cache["evt-expire"]
	ce.cachedAt = time.Now().Add(-cachedExpiry - time.Minute)
	m.cache["evt-expire"] = ce
	m.mu.Unlock()
	m.pruneExpiredCachedEvents(ctx)

	after := collectFailoverTotals(t)

	cached := after.triggerCached - before.triggerCached
	drained := after.cacheExit[monitoring.ShardFailoverCacheExitDrained] - before.cacheExit[monitoring.ShardFailoverCacheExitDrained]
	replayed := after.cacheExit[monitoring.ShardFailoverCacheExitReplayed] - before.cacheExit[monitoring.ShardFailoverCacheExitReplayed]
	expired := after.cacheExit[monitoring.ShardFailoverCacheExitExpired] - before.cacheExit[monitoring.ShardFailoverCacheExitExpired]
	stillCached := int64(len(m.cacheIDs()))

	assert.Equal(t, int64(5), cached, "events cached across both managers")
	assert.Equal(t, int64(1), drained, "drained exits")
	assert.Equal(t, int64(2), replayed, "replayed exits")
	assert.Equal(t, int64(1), expired, "expired exits")
	assert.Equal(t, int64(1), stillCached, "events still cached")
	assert.Equal(t, cached, drained+replayed+expired+stillCached, "cache funnel invariant")

	assert.Equal(t, int64(1), after.replay[monitoring.ShardFailoverReplaySucceeded]-before.replay[monitoring.ShardFailoverReplaySucceeded])
	assert.Equal(t, int64(1), after.replay[monitoring.ShardFailoverReplayFailed]-before.replay[monitoring.ShardFailoverReplayFailed])
	assert.Equal(t, uint64(1), after.replayCount[monitoring.ShardFailoverReplaySucceeded]-before.replayCount[monitoring.ShardFailoverReplaySucceeded])
	assert.Equal(t, uint64(1), after.replayCount[monitoring.ShardFailoverReplayFailed]-before.replayCount[monitoring.ShardFailoverReplayFailed])

	cachedGauge, ok := failoverGauge(t, "platform_shard_failover_cached_events")
	require.True(t, ok, "cached_events gauge recorded")
	assert.Equal(t, stillCached, cachedGauge, "cached_events gauge tracks the cache size")
}

// TestShardFailoverManager_ReplayMiss covers the miss replay outcome: a
// SYSTEM_ERROR verdict for a trigger event that was never cached counts as a
// miss and leaves the cache funnel untouched.
//
//nolint:paralleltest // asserts deltas on process-global failover metrics
func TestShardFailoverManager_ReplayMiss(t *testing.T) {
	m := newFailoverTestManager(t, &recordingEngine{})

	before := collectFailoverTotals(t)

	m.HandleExecutionStatusUpdate(&ringpb.ExecutionStatusUpdate{
		TriggerEventId: "evt-never-cached",
		Status:         ringpb.ExecutionStatus_EXECUTION_STATUS_SYSTEM_ERROR,
	})

	after := collectFailoverTotals(t)

	assert.Equal(t, int64(1), after.replay[monitoring.ShardFailoverReplayMiss]-before.replay[monitoring.ShardFailoverReplayMiss], "miss replays")
	assert.Equal(t, int64(0), after.triggerCached-before.triggerCached, "nothing was cached")
	for _, outcome := range []string{
		monitoring.ShardFailoverCacheExitDrained,
		monitoring.ShardFailoverCacheExitReplayed,
		monitoring.ShardFailoverCacheExitExpired,
	} {
		assert.Equal(t, int64(0), after.cacheExit[outcome]-before.cacheExit[outcome], "no cache exit for outcome %s", outcome)
	}
	assert.Empty(t, m.cacheIDs(), "the cache stays empty")
}

// TestShardFailoverManager_PruneExpiresCachedEvents covers the expired exit:
// events older than the 10-minute expiry are pruned with no verdict ever
// arriving, and the cached_events gauge returns to 0.
//
//nolint:paralleltest // asserts deltas on process-global failover metrics
func TestShardFailoverManager_PruneExpiresCachedEvents(t *testing.T) {
	ctx := t.Context()
	m := newFailoverTestManager(t, &recordingEngine{})

	before := collectFailoverTotals(t)

	require.ErrorIs(t, m.ExecuteTrigger(ctx, testCoordinatedEvent("evt-1")), triggers.ErrAdmissionCached)
	require.ErrorIs(t, m.ExecuteTrigger(ctx, testCoordinatedEvent("evt-2")), triggers.ErrAdmissionCached)

	m.mu.Lock()
	for id, ce := range m.cache {
		ce.cachedAt = time.Now().Add(-cachedExpiry - time.Minute)
		m.cache[id] = ce
	}
	m.mu.Unlock()
	m.pruneExpiredCachedEvents(ctx)

	after := collectFailoverTotals(t)

	assert.Equal(t, int64(2), after.cacheExit[monitoring.ShardFailoverCacheExitExpired]-before.cacheExit[monitoring.ShardFailoverCacheExitExpired], "expired exits")
	assert.Empty(t, m.cacheIDs(), "cache emptied by the prune")
	cachedGauge, ok := failoverGauge(t, "platform_shard_failover_cached_events")
	require.True(t, ok, "cached_events gauge recorded")
	assert.Equal(t, int64(0), cachedGauge, "cached_events gauge back to 0 after prune")
}

// TestShardFailoverManager_TimeoutBasedFailoverCacheExits covers the two
// timeout-based pop paths: a due event dropped because this shard owns the
// workflow again exits as drained, while a due event auto-executed by this
// shard exits as replayed.
//
//nolint:paralleltest // asserts deltas on process-global failover metrics
func TestShardFailoverManager_TimeoutBasedFailoverCacheExits(t *testing.T) {
	ctx := t.Context()
	eng := &recordingEngine{}
	resolver := &switchableShardResolver{ownerDON: 1, found: true}
	m := newTimeoutBasedFailoverManager(t, eng, &fakeGateLimiter{open: true}, &fakeWindowLimiter{window: 10 * time.Millisecond}, resolver)

	before := collectFailoverTotals(t)

	// The event is due while this shard owns the workflow again: dropped, not executed.
	require.ErrorIs(t, m.admissionCheck(ctx, timeoutBasedFailoverTestEvent("evt-owner-again")), v2.ErrAdmissionCache)
	time.Sleep(15 * time.Millisecond)
	resolver.setOwner(2)
	m.executeDueTimeoutBasedFailovers(ctx)
	assert.Empty(t, eng.executedIDs(), "cached event executed although this shard became the owner again")
	assert.Empty(t, m.cacheIDs(), "owner-again event left the cache without a replayed exit")

	// The event is due while this shard is still the secondary: auto-executed.
	resolver.setOwner(1)
	require.ErrorIs(t, m.admissionCheck(ctx, timeoutBasedFailoverTestEvent("evt-auto")), v2.ErrAdmissionCache)
	time.Sleep(15 * time.Millisecond)
	m.executeDueTimeoutBasedFailovers(ctx)
	assert.Eventually(t, func() bool {
		ids := eng.executedIDs()
		return len(ids) == 1 && ids[0] == "evt-auto"
	}, 5*time.Second, 10*time.Millisecond, "cached event was not auto-executed after the failover window elapsed")

	after := collectFailoverTotals(t)

	assert.Equal(t, int64(1), after.cacheExit[monitoring.ShardFailoverCacheExitDrained]-before.cacheExit[monitoring.ShardFailoverCacheExitDrained], "owner-again drop exits as drained")
	assert.Equal(t, int64(1), after.cacheExit[monitoring.ShardFailoverCacheExitReplayed]-before.cacheExit[monitoring.ShardFailoverCacheExitReplayed], "auto-execution exits as replayed")
	assert.Empty(t, m.cacheIDs())
}

// TestWorkflowRegistry_ShardFailoverEnabledGauge covers the failover gate
// evaluation in the shard filter: the enabled gauge follows the gate per sync
// tick.
//
//nolint:paralleltest // asserts on process-global failover metrics
func TestWorkflowRegistry_ShardFailoverEnabledGauge(t *testing.T) {
	ctx := t.Context()
	wf := wfTypes.WorkflowID([32]byte{1})
	workflows := []WorkflowMetadataView{{WorkflowID: wf, WorkflowName: "wf1"}}
	client := &mockShardMappingClient{mappings: map[string]uint32{wf.Hex(): 0}}

	newReg := func(gate limits.GateLimiter) *workflowRegistry {
		return &workflowRegistry{
			shardResolver:        shardownership.NewRingOCRShardResolver(client, logger.TestLogger(t)),
			myDonID:              1,
			shardingEnabled:      true,
			shardingFailoverGate: gate,
		}
	}

	wr := newReg(limits.NewGateLimiter(true))
	_, err := wr.filterWorkflowsByShard(ctx, workflows)
	require.NoError(t, err)
	enabled, ok := failoverGauge(t, "platform_shard_failover_enabled")
	require.True(t, ok, "enabled gauge recorded")
	assert.Equal(t, int64(1), enabled, "failover enabled")

	wr = newReg(limits.NewGateLimiter(false))
	_, err = wr.filterWorkflowsByShard(ctx, workflows)
	require.NoError(t, err)
	enabled, ok = failoverGauge(t, "platform_shard_failover_enabled")
	require.True(t, ok, "enabled gauge recorded")
	assert.Equal(t, int64(0), enabled, "failover disabled")
}
