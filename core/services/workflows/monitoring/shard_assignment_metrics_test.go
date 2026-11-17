package monitoring

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

func TestShardAssignmentConstants(t *testing.T) {
	t.Parallel()
	require.Equal(t, int64(0), ShardAssignmentModeManualOnly)
	require.Equal(t, int64(1), ShardAssignmentModeRingOCROnly)
	require.Equal(t, int64(2), ShardAssignmentModeRingOCRWithOverrides)

	require.Equal(t, "manual", ShardResolverManual)
	require.Equal(t, "ring", ShardResolverRing)
	require.Equal(t, "override", ShardResolverOverride)

	require.Equal(t, "found", ShardResolverOutcomeFound)
	require.Equal(t, "not_found", ShardResolverOutcomeNotFound)
	require.Equal(t, "error", ShardResolverOutcomeError)

	require.Equal(t, "per_owner", ShardResolverSourcePerOwner)
	require.Equal(t, "per_org", ShardResolverSourcePerOrg)
	require.Equal(t, "static_default", ShardResolverSourceStaticDefault)
	require.Equal(t, "hashed", ShardResolverSourceHashed)
	require.Equal(t, "none", ShardResolverSourceNone)

	require.Equal(t, "applied", ShardAssignmentUpdateApplied)
	require.Equal(t, "parse_error", ShardAssignmentUpdateParseError)
}

func TestGlobalShardAssignmentMetricsSingleton(t *testing.T) {
	t.Parallel()
	m, err := GlobalShardAssignmentMetrics()
	require.NoError(t, err)
	require.NotNil(t, m)

	m2, err := GlobalShardAssignmentMetrics()
	require.NoError(t, err)
	require.Same(t, m, m2)
}

func TestShardAssignmentMetricsNilReceiver(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	var m *ShardAssignmentMetrics
	m.RecordAssignmentMode(ctx, ShardAssignmentModeManualOnly)
	m.RecordResolverRequest(ctx, ShardResolverManual, ShardResolverOutcomeFound, ShardResolverSourcePerOwner, time.Second)
	m.AddResolverRequest(ctx, ShardResolverManual, ShardResolverOutcomeNotFound, "")
	m.RecordResolverDuration(ctx, ShardResolverRing, time.Second)
	m.RecordSettingsArrival("raw", time.Now())
	require.False(t, m.RecordSettingsApplied(ctx, 1, "raw"))
	require.False(t, m.RecordSettingsParseError(ctx, "raw"))
}

// TestGlobalShardAssignmentMetricsRebindsWhenClientChanges proves the global
// binding follows a beholder client installed after the first call: beholder
// starts with a noop client, and one made too early must not be cached for
// the life of the process.
func TestGlobalShardAssignmentMetricsRebindsWhenClientChanges(t *testing.T) { //nolint:paralleltest // swaps the process-global beholder client
	prevClient := beholder.GetClient()
	t.Cleanup(func() { beholder.SetClient(prevClient) })

	newClient := func(name string) (*beholder.Client, *sdkmetric.ManualReader) {
		reader := sdkmetric.NewManualReader()
		mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
		t.Cleanup(func() { _ = mp.Shutdown(context.Background()) })
		return &beholder.Client{Meter: mp.Meter(name)}, reader
	}

	beholder.SetClient(beholder.NewNoopClient())
	mNoop, err := GlobalShardAssignmentMetrics()
	require.NoError(t, err)
	require.NotNil(t, mNoop)

	client1, reader1 := newClient("rebind-test-1")
	beholder.SetClient(client1)
	m1, err := GlobalShardAssignmentMetrics()
	require.NoError(t, err)
	require.NotNil(t, m1, "the binding must be retried after the global client changes")
	require.NotSame(t, mNoop, m1)

	client2, reader2 := newClient("rebind-test-2")
	beholder.SetClient(client2)
	m2, err := GlobalShardAssignmentMetrics()
	require.NoError(t, err)
	require.NotSame(t, m1, m2)

	m2.RecordResolverRequest(t.Context(), ShardResolverManual, ShardResolverOutcomeFound, ShardResolverSourcePerOwner, time.Second)
	require.Equal(t, int64(1), counterTotal(collect(t, reader2), "platform_shard_resolver_requests",
		map[string]string{"resolver": ShardResolverManual, "outcome": ShardResolverOutcomeFound}),
		"metrics recorded after the swap must land on the new client's meter")
	require.Zero(t, counterTotal(collect(t, reader1), "platform_shard_resolver_requests", nil),
		"the abandoned binding must not receive further recordings")
}

func collect(t *testing.T, reader *sdkmetric.ManualReader) metricdata.ResourceMetrics {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	return rm
}

func TestShardAssignmentMetricsInstruments(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	reader, m := newTestShardAssignmentMetrics(t)

	m.RecordAssignmentMode(ctx, ShardAssignmentModeRingOCRWithOverrides)
	m.RecordResolverRequest(ctx, ShardResolverManual, ShardResolverOutcomeFound, ShardResolverSourcePerOwner, 2*time.Second)
	m.RecordResolverRequest(ctx, ShardResolverRing, ShardResolverOutcomeNotFound, "", 4*time.Millisecond)
	m.AddResolverRequest(ctx, ShardResolverOverride, ShardResolverOutcomeError, "")
	m.RecordSettingsArrival("settings-blob", time.Now().Add(-90*time.Second))
	m.RecordSettingsApplied(ctx, 4242, "settings-blob")

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(ctx, &rm))
	metrics := map[string]metricdata.Metrics{}
	for _, sm := range rm.ScopeMetrics {
		for _, met := range sm.Metrics {
			metrics[met.Name] = met
		}
	}

	mode, ok := metrics["platform_shard_assignment_mode"]
	require.True(t, ok, "mode gauge must be registered")
	gauge, ok := mode.Data.(metricdata.Gauge[int64])
	require.True(t, ok)
	require.Len(t, gauge.DataPoints, 1)
	require.Equal(t, ShardAssignmentModeRingOCRWithOverrides, gauge.DataPoints[0].Value)

	requests, ok := metrics["platform_shard_resolver_requests"]
	require.True(t, ok, "resolver requests counter must be registered")
	sum, ok := requests.Data.(metricdata.Sum[int64])
	require.True(t, ok)
	require.Equal(t, int64(3), sumDataTotal(sum))

	duration, ok := metrics["platform_shard_resolver_duration_seconds"]
	require.True(t, ok, "resolver duration histogram must be registered")
	require.Equal(t, "s", duration.Unit)
	hist, ok := duration.Data.(metricdata.Histogram[float64])
	require.True(t, ok)
	require.NotEmpty(t, hist.DataPoints)
	require.Equal(t, []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5, 30}, hist.DataPoints[0].Bounds)
	require.Equal(t, uint64(2), histogramDataCount(hist, map[string]string{}),
		"histogram counts API-call durations; AddResolverRequest does not record one")
	require.Equal(t, []uint64{0, 0, 0, 0, 0, 0, 0, 1, 0, 0}, histogramBucketCounts(hist, map[string]string{"resolver": ShardResolverManual}))
	require.Equal(t, []uint64{0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, histogramBucketCounts(hist, map[string]string{"resolver": ShardResolverRing}),
		"a 4ms resolution must record as 0.004 in the (0.001, 0.005] bucket, not truncate to 0")

	hashGauge, ok := metrics["platform_shard_assignment_settings_hash"]
	require.True(t, ok, "settings hash gauge must be registered")
	gauge, ok = hashGauge.Data.(metricdata.Gauge[int64])
	require.True(t, ok)
	require.Equal(t, int64(4242), gaugeValue(rm, "platform_shard_assignment_settings_hash"))

	updateCountGauge, ok := metrics["platform_shard_assignment_config_update_count"]
	require.True(t, ok, "config update count gauge must be registered")
	gauge, ok = updateCountGauge.Data.(metricdata.Gauge[int64])
	require.True(t, ok)
	require.Len(t, gauge.DataPoints, 1)
	require.Equal(t, int64(1), gauge.DataPoints[0].Value, "the first applied settings update count is 1")

	stalenessMet, ok := metrics["platform_shard_assignment_config_staleness_seconds"]
	require.True(t, ok, "config staleness histogram must be registered")
	require.Equal(t, "s", stalenessMet.Unit)
	stalenessHist, ok := stalenessMet.Data.(metricdata.Histogram[float64])
	require.True(t, ok)
	require.Equal(t, []float64{0.1, 0.5, 1, 2, 5, 10, 30, 60, 120, 300}, stalenessHist.DataPoints[0].Bounds)
	require.Equal(t, uint64(1), histogramDataCount(stalenessHist, map[string]string{}))
	require.Equal(t, []uint64{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0}, histogramBucketCounts(stalenessHist, map[string]string{}),
		"a 90s arrival→in-effect delay must land in the (60, 120] bucket")

	_, ok = metrics["platform_shard_assignment_updated"]
	require.True(t, ok, "assignment updated counter must be registered")
	require.Equal(t, int64(1), counterTotal(rm, "platform_shard_assignment_updated", map[string]string{"outcome": ShardAssignmentUpdateApplied}))
}

func TestShardAssignmentMetricsSettingsUpdateDedup(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	reader, m := newTestShardAssignmentMetrics(t)

	require.True(t, m.RecordSettingsApplied(ctx, 111, "raw-111"))
	require.False(t, m.RecordSettingsApplied(ctx, 111, "raw-111"))
	require.True(t, m.RecordSettingsApplied(ctx, 222, "raw-222"))
	require.True(t, m.RecordSettingsApplied(ctx, 111, "raw-111"))

	require.True(t, m.RecordSettingsParseError(ctx, "raw-333"))
	require.False(t, m.RecordSettingsParseError(ctx, "raw-333"))
	require.True(t, m.RecordSettingsParseError(ctx, "raw-444"))

	// A successfully applied version resets the parse-error dedup, so the
	// same broken blob recurring after a good one counts (and logs) again.
	require.True(t, m.RecordSettingsApplied(ctx, 555, "raw-555"))
	require.True(t, m.RecordSettingsParseError(ctx, "raw-333"),
		"a bad version recurring after a good one must be counted again")

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(ctx, &rm))
	require.Equal(t, int64(4), counterTotal(rm, "platform_shard_assignment_updated", map[string]string{"outcome": ShardAssignmentUpdateApplied}))
	require.Equal(t, int64(3), counterTotal(rm, "platform_shard_assignment_updated", map[string]string{"outcome": ShardAssignmentUpdateParseError}))

	require.Equal(t, int64(555), gaugeValue(rm, "platform_shard_assignment_settings_hash"),
		"parse errors must not move the hash gauge: the last applied version stays")
}

func TestShardAssignmentMetricsConfigUpdateCount(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	reader, m := newTestShardAssignmentMetrics(t)

	require.Zero(t, gaugeValue(collect(t, reader), "platform_shard_assignment_config_update_count"),
		"no update count is recorded before anything is applied")

	require.True(t, m.RecordSettingsApplied(ctx, 111, "raw-111"))
	require.Equal(t, int64(1), gaugeValue(collect(t, reader), "platform_shard_assignment_config_update_count"))

	require.False(t, m.RecordSettingsApplied(ctx, 111, "raw-111"), "re-observing the same blob must not bump the update count")
	require.Equal(t, int64(1), gaugeValue(collect(t, reader), "platform_shard_assignment_config_update_count"))

	require.True(t, m.RecordSettingsApplied(ctx, 222, "raw-222"))
	require.Equal(t, int64(2), gaugeValue(collect(t, reader), "platform_shard_assignment_config_update_count"))

	require.True(t, m.RecordSettingsParseError(ctx, "raw-333"))
	require.Equal(t, int64(2), gaugeValue(collect(t, reader), "platform_shard_assignment_config_update_count"),
		"a parse error must not move the update count gauge: a broken version is not applied")

	require.True(t, m.RecordSettingsApplied(ctx, 111, "raw-111"), "a version returning after a different one is a change")
	require.Equal(t, int64(3), gaugeValue(collect(t, reader), "platform_shard_assignment_config_update_count"))
}

func TestShardAssignmentMetricsConfigStaleness(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	// An arrival preceding the apply is consumed and recorded in the (30, 60]
	// bucket of the staleness histogram.
	reader, m := newTestShardAssignmentMetrics(t)
	m.RecordSettingsArrival("raw-1111", time.Now().Add(-45*time.Second))
	m.RecordSettingsApplied(ctx, 111, "raw-1111")
	hist := stalenessHistogram(t, reader)
	require.Equal(t, uint64(1), histogramDataCount(hist, nil))
	require.Equal(t, []uint64{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0}, histogramBucketCounts(hist, nil),
		"a 45s arrival→in-effect delay must land in the (30, 60] bucket")
	require.Zero(t, m.lastArrivalRawHash, "a matched arrival is consumed on apply")

	// Re-observing the same blob is not a change and records nothing.
	m.RecordSettingsApplied(ctx, 111, "raw-1111")
	require.Equal(t, uint64(1), histogramDataCount(stalenessHistogram(t, reader), nil))

	// A superseded arrival (a newer blob arrived after the one applied)
	// produces no sample and keeps its slot for a later matching apply.
	m.RecordSettingsArrival("raw-9999", time.Now().Add(-10*time.Second))
	m.RecordSettingsApplied(ctx, 222, "raw-2222")
	require.Equal(t, uint64(1), histogramDataCount(stalenessHistogram(t, reader), nil))
	require.Equal(t, settingsStringHash("raw-9999"), m.lastArrivalRawHash)

	// A change with no arrival observed at all produces no sample.
	reader2, m2 := newTestShardAssignmentMetrics(t)
	m2.RecordSettingsApplied(ctx, 333, "raw-3333")
	require.Zero(t, histogramDataCount(stalenessHistogram(t, reader2), nil))

	// Semantically identical settings (same parsed hash, new raw blob) are
	// not a change: no sample and the arrival slot is left alone.
	m2.RecordSettingsArrival("raw-5555", time.Now().Add(-10*time.Second))
	m2.RecordSettingsApplied(ctx, 333, "raw-5555")
	require.Zero(t, histogramDataCount(stalenessHistogram(t, reader2), nil))
	require.Equal(t, settingsStringHash("raw-5555"), m2.lastArrivalRawHash)
}

func stalenessHistogram(t *testing.T, reader *sdkmetric.ManualReader) metricdata.Histogram[float64] {
	t.Helper()
	for _, sm := range collect(t, reader).ScopeMetrics {
		for _, met := range sm.Metrics {
			if met.Name != "platform_shard_assignment_config_staleness_seconds" {
				continue
			}
			hist, ok := met.Data.(metricdata.Histogram[float64])
			require.True(t, ok)
			return hist
		}
	}
	return metricdata.Histogram[float64]{}
}

func newTestShardAssignmentMetrics(t *testing.T) (*sdkmetric.ManualReader, *ShardAssignmentMetrics) {
	t.Helper()
	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { require.NoError(t, mp.Shutdown(context.Background())) })
	m, err := NewShardAssignmentMetricsWithMeter(mp.Meter("test-shard-assignment"))
	require.NoError(t, err)
	return reader, m
}

func counterTotal(rm metricdata.ResourceMetrics, name string, attrs map[string]string) int64 {
	var total int64
	for _, sm := range rm.ScopeMetrics {
		for _, met := range sm.Metrics {
			if met.Name != name {
				continue
			}
			sum, ok := met.Data.(metricdata.Sum[int64])
			if !ok {
				continue
			}
			for _, dp := range sum.DataPoints {
				if hasAttrs(dp.Attributes, attrs) {
					total += dp.Value
				}
			}
		}
	}
	return total
}

func sumDataTotal(sum metricdata.Sum[int64]) int64 {
	var total int64
	for _, dp := range sum.DataPoints {
		total += dp.Value
	}
	return total
}

func gaugeValue(rm metricdata.ResourceMetrics, name string) int64 {
	for _, sm := range rm.ScopeMetrics {
		for _, met := range sm.Metrics {
			if met.Name != name {
				continue
			}
			gauge, ok := met.Data.(metricdata.Gauge[int64])
			if !ok || len(gauge.DataPoints) != 1 {
				continue
			}
			return gauge.DataPoints[0].Value
		}
	}
	return 0
}

func histogramDataCount(hist metricdata.Histogram[float64], attrs map[string]string) uint64 {
	var count uint64
	for _, dp := range hist.DataPoints {
		if hasAttrs(dp.Attributes, attrs) {
			count += dp.Count
		}
	}
	return count
}

func histogramBucketCounts(hist metricdata.Histogram[float64], attrs map[string]string) []uint64 {
	var counts []uint64
	for _, dp := range hist.DataPoints {
		if hasAttrs(dp.Attributes, attrs) {
			if counts == nil {
				counts = make([]uint64, len(dp.BucketCounts))
			}
			for i, c := range dp.BucketCounts {
				counts[i] += c
			}
		}
	}
	return counts
}

func hasAttrs(set attribute.Set, attrs map[string]string) bool {
	for k, want := range attrs {
		if got, ok := set.Value(attribute.Key(k)); !ok || got.AsString() != want {
			return false
		}
	}
	return true
}
