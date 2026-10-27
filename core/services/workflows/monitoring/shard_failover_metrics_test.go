package monitoring

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

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
	fm.UpdateCachedEvents(ctx, 2)
	fm.RecordCacheExit(ctx, ShardFailoverCacheExitExpired, 2)
	fm.RecordReplay(ctx, ShardFailoverReplaySucceeded)
	fm.RecordReplayDuration(ctx, ShardFailoverReplaySucceeded, 1500*time.Millisecond)

	var nilFM *ShardFailoverMetrics
	nilFM.SetFailoverEnabled(ctx, true)
	nilFM.RecordTriggerCached(ctx)
	nilFM.UpdateCachedEvents(ctx, 1)
	nilFM.RecordCacheExit(ctx, ShardFailoverCacheExitDrained, 1)
	nilFM.RecordReplay(ctx, ShardFailoverReplayFailed)
	nilFM.RecordReplayDuration(ctx, ShardFailoverReplayFailed, time.Second)

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(ctx, &rm))

	gauges := make(map[string]int64)
	sums := make(map[string]map[string]int64)
	histCounts := make(map[string]map[string]uint64)
	units := make(map[string]string)
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
			case metricdata.Histogram[int64]:
				for _, dp := range data.DataPoints {
					outcome := attrOutcome(dp.Attributes)
					if histCounts[m.Name] == nil {
						histCounts[m.Name] = make(map[string]uint64)
					}
					histCounts[m.Name][outcome] += dp.Count
				}
			}
		}
	}

	assert.Equal(t, int64(1), gauges["platform_shard_failover_enabled"])
	assert.Equal(t, int64(2), gauges["platform_shard_failover_cached_events"])
	assert.Equal(t, int64(1), sums["platform_shard_failover_trigger_cached"][""])
	assert.Equal(t, int64(2), sums["platform_shard_failover_cache_exit"][ShardFailoverCacheExitExpired])
	assert.Equal(t, int64(1), sums["platform_shard_failover_replay"][ShardFailoverReplaySucceeded])
	assert.Equal(t, uint64(1), histCounts["platform_shard_failover_replay_duration_seconds"][ShardFailoverReplaySucceeded])
	assert.Equal(t, "s", units["platform_shard_failover_replay_duration_seconds"])
}

func attrOutcome(attrs attribute.Set) string {
	if v, ok := attrs.Value(attribute.Key("outcome")); ok {
		return v.AsString()
	}
	return ""
}
