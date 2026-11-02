package shardownership

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/services/cresettings"
	wfmon "github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

const (
	ownerA = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
	ownerB = "0x70997970c51812dc3a010c7d01b50e0d17dc79c8"
	ownerC = "0x0000000000000000000000000000000000000003"
	ownerD = "0x0000000000000000000000000000000000000004"
)

type assignmentMetricsFixture struct {
	reader  *sdkmetric.ManualReader
	metrics *wfmon.ShardAssignmentMetrics
}

func newAssignmentMetricsFixture(t *testing.T) *assignmentMetricsFixture {
	t.Helper()
	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { require.NoError(t, mp.Shutdown(context.Background())) })
	metrics, err := wfmon.NewShardAssignmentMetricsWithMeter(mp.Meter("test-shard-assignment"))
	require.NoError(t, err)
	return &assignmentMetricsFixture{reader: reader, metrics: metrics}
}

func (f *assignmentMetricsFixture) manualResolver(settings *loop.AtomicSettings, orgResolver OrgResolver, lggr logger.Logger) *manualShardResolver {
	r := NewManualShardResolver(settings, orgResolver, nil, lggr).(*manualShardResolver)
	r.metrics = f.metrics
	return r
}

func (f *assignmentMetricsFixture) collect(t *testing.T) metricdata.ResourceMetrics {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, f.reader.Collect(t.Context(), &rm))
	return rm
}

func (f *assignmentMetricsFixture) counterTotal(t *testing.T, name string, attrs map[string]string) int64 {
	t.Helper()
	var total int64
	for _, sm := range f.collect(t).ScopeMetrics {
		for _, met := range sm.Metrics {
			if met.Name != name {
				continue
			}
			sum, ok := met.Data.(metricdata.Sum[int64])
			if !ok {
				continue
			}
			for _, dp := range sum.DataPoints {
				if hasAllAttrs(dp.Attributes, attrs) {
					total += dp.Value
				}
			}
		}
	}
	return total
}

func (f *assignmentMetricsFixture) gaugeValue(t *testing.T, name string) (int64, bool) {
	t.Helper()
	for _, sm := range f.collect(t).ScopeMetrics {
		for _, met := range sm.Metrics {
			if met.Name != name {
				continue
			}
			gauge, ok := met.Data.(metricdata.Gauge[int64])
			if !ok || len(gauge.DataPoints) != 1 {
				continue
			}
			return gauge.DataPoints[0].Value, true
		}
	}
	return 0, false
}

func (f *assignmentMetricsFixture) histogramCount(t *testing.T, name string, attrs map[string]string) uint64 {
	t.Helper()
	var count uint64
	for _, sm := range f.collect(t).ScopeMetrics {
		for _, met := range sm.Metrics {
			if met.Name != name {
				continue
			}
			hist, ok := met.Data.(metricdata.Histogram[float64])
			if !ok {
				continue
			}
			for _, dp := range hist.DataPoints {
				if hasAllAttrs(dp.Attributes, attrs) {
					count += dp.Count
				}
			}
		}
	}
	return count
}

func (f *assignmentMetricsFixture) observedSources(t *testing.T) []string {
	t.Helper()
	seen := map[string]struct{}{}
	for _, sm := range f.collect(t).ScopeMetrics {
		for _, met := range sm.Metrics {
			if met.Name != "platform_shard_resolver_requests" {
				continue
			}
			sum, ok := met.Data.(metricdata.Sum[int64])
			if !ok {
				continue
			}
			for _, dp := range sum.DataPoints {
				if v, ok := dp.Attributes.Value(attribute.Key("source")); ok {
					seen[v.AsString()] = struct{}{}
				}
			}
		}
	}
	return slices.Sorted(maps.Keys(seen))
}

func hasAllAttrs(set attribute.Set, attrs map[string]string) bool {
	for k, want := range attrs {
		if got, ok := set.Value(attribute.Key(k)); !ok || got.AsString() != want {
			return false
		}
	}
	return true
}

func (f *assignmentMetricsFixture) requireSourceWithinBoundedValues(t *testing.T) {
	t.Helper()
	allowed := []string{
		wfmon.ShardResolverSourcePerOwner,
		wfmon.ShardResolverSourcePerOrg,
		wfmon.ShardResolverSourceStaticDefault,
		wfmon.ShardResolverSourceHashed,
		wfmon.ShardResolverSourceNone,
	}
	for _, source := range f.observedSources(t) {
		require.Contains(t, allowed, source, "source label must stay within its 5 enum values")
	}
}

func TestResolverMetrics_ManualSources(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	f := newAssignmentMetricsFixture(t)

	settings := &loop.AtomicSettings{}
	storeShardAssignment(t, settings, `
static_default_assignment = [0]
hashed_owner_assignment = ["`+ownerC+`"]

[per_owner_assignment]
  `+ownerA+` = [1]

[per_org_assignment]
  org_abc123 = [2]
`)
	orgR := &mockOrgResolver{mappings: map[string]string{ownerB: "org_abc123"}}
	r := f.manualResolver(settings, orgR, nopLogger)

	donID, found, err := r.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, uint32(1), donID)

	donID, found, err = r.ResolveShard(ctx, "wf-b", ownerB)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, uint32(2), donID)

	_, found, err = r.ResolveShard(ctx, "wf-c", ownerC)
	require.NoError(t, err)
	require.False(t, found)

	donID, found, err = r.ResolveShard(ctx, "wf-d", ownerD)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, uint32(0), donID)

	shards, found, err := r.ResolveAllShards(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, []uint32{1}, shards)

	resolved, err := r.ResolveShards(ctx, []string{"wf-a", "wf-d"}, []string{ownerA, ownerD})
	require.NoError(t, err)
	require.Equal(t, map[string]uint32{"wf-a": 1, "wf-d": 0}, resolved)

	require.Equal(t, int64(3), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverManual, "outcome": wfmon.ShardResolverOutcomeFound, "source": wfmon.ShardResolverSourcePerOwner,
	}))
	require.Equal(t, int64(1), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverManual, "outcome": wfmon.ShardResolverOutcomeFound, "source": wfmon.ShardResolverSourcePerOrg,
	}))
	require.Equal(t, int64(1), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverManual, "outcome": wfmon.ShardResolverOutcomeNotFound, "source": wfmon.ShardResolverSourceHashed,
	}))
	require.Equal(t, int64(2), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverManual, "outcome": wfmon.ShardResolverOutcomeFound, "source": wfmon.ShardResolverSourceStaticDefault,
	}))
	require.Equal(t, uint64(6), f.histogramCount(t, "platform_shard_resolver_duration_seconds", map[string]string{"resolver": wfmon.ShardResolverManual}))

	require.ElementsMatch(t, []string{
		wfmon.ShardResolverSourcePerOwner,
		wfmon.ShardResolverSourcePerOrg,
		wfmon.ShardResolverSourceStaticDefault,
		wfmon.ShardResolverSourceHashed,
	}, f.observedSources(t))
	f.requireSourceWithinBoundedValues(t)

	noDefault := &loop.AtomicSettings{}
	storeShardAssignment(t, noDefault, `
[per_owner_assignment]
  `+ownerA+` = [1]
`)
	rNoDefault := f.manualResolver(noDefault, nil, nopLogger)
	_, found, err = rNoDefault.ResolveShard(ctx, "wf-d", ownerD)
	require.NoError(t, err)
	require.False(t, found)
	require.Equal(t, int64(1), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverManual, "outcome": wfmon.ShardResolverOutcomeNotFound, "source": wfmon.ShardResolverSourceNone,
	}), "the nothing-matched-no-default fall-through must report source=none, not static_default")
	f.requireSourceWithinBoundedValues(t)
}

func TestResolverMetrics_SettingsHash(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	raw := `static_default_assignment = [0]

[per_owner_assignment]
  ` + ownerA + ` = [1]
  ` + ownerB + ` = [2]
`
	reformatted := `static_default_assignment = [0]

[per_owner_assignment]
  "0X` + ownerB[2:] + `" = [2]
  "` + ownerA + `" = [1]
`
	oneOwnerChanged := `static_default_assignment = [0]

[per_owner_assignment]
  ` + ownerA + ` = [3]
  ` + ownerB + ` = [2]
`

	cfg, err := cresettings.ParseShardAssignmentConfig(raw)
	require.NoError(t, err)
	cfgReformatted, err := cresettings.ParseShardAssignmentConfig(reformatted)
	require.NoError(t, err)
	cfgChanged, err := cresettings.ParseShardAssignmentConfig(oneOwnerChanged)
	require.NoError(t, err)

	stable := shardAssignmentSettingsHash(cfg)
	require.Equal(t, stable, shardAssignmentSettingsHash(cfgReformatted), "identical settings must hash identically")
	require.NotEqual(t, stable, shardAssignmentSettingsHash(cfgChanged), "a one-owner change must change the hash")

	f := newAssignmentMetricsFixture(t)
	settings := &loop.AtomicSettings{}
	storeShardAssignment(t, settings, raw)
	r := f.manualResolver(settings, nil, nopLogger)

	donID, found, err := r.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, uint32(1), donID)

	hash, ok := f.gaugeValue(t, "platform_shard_assignment_settings_hash")
	require.True(t, ok)
	require.Equal(t, int64(stable), hash)
	require.Equal(t, int64(1), f.counterTotal(t, "platform_shard_assignment_updated", map[string]string{"outcome": wfmon.ShardAssignmentUpdateApplied}))

	_, found, err = r.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, int64(1), f.counterTotal(t, "platform_shard_assignment_updated", map[string]string{"outcome": wfmon.ShardAssignmentUpdateApplied}),
		"identical settings must not count another applied update")

	storeShardAssignment(t, settings, raw)
	_, found, err = r.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, int64(1), f.counterTotal(t, "platform_shard_assignment_updated", map[string]string{"outcome": wfmon.ShardAssignmentUpdateApplied}))

	storeShardAssignment(t, settings, oneOwnerChanged)
	donID, found, err = r.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, uint32(3), donID)
	require.Equal(t, int64(2), f.counterTotal(t, "platform_shard_assignment_updated", map[string]string{"outcome": wfmon.ShardAssignmentUpdateApplied}))

	hash, ok = f.gaugeValue(t, "platform_shard_assignment_settings_hash")
	require.True(t, ok)
	require.Equal(t, int64(shardAssignmentSettingsHash(cfgChanged)), hash)
	require.NotEqual(t, int64(stable), hash)
}

func TestResolverMetrics_ParseErrorPropagates(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	f := newAssignmentMetricsFixture(t)

	lggr, logs := logger.TestObserved(t, zapcore.WarnLevel)
	settings := &loop.AtomicSettings{}
	storeShardAssignment(t, settings, `
static_default_assignment = [0]

[per_owner_assignment]
  `+ownerA+` = [1]
`)
	r := f.manualResolver(settings, nil, lggr)

	donID, found, err := r.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, uint32(1), donID)
	firstHash, ok := f.gaugeValue(t, "platform_shard_assignment_settings_hash")
	require.True(t, ok)

	bad := `static_default_assignment = "not-an-array"`
	storeShardAssignment(t, settings, bad)
	_, found, err = r.ResolveShard(ctx, "wf-a", ownerA)
	require.Error(t, err, "a parse error must fail resolution exactly as it did before instrumentation")
	require.False(t, found)
	require.Equal(t, int64(1), f.counterTotal(t, "platform_shard_assignment_updated", map[string]string{"outcome": wfmon.ShardAssignmentUpdateParseError}))
	require.Equal(t, int64(1), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverManual, "outcome": wfmon.ShardResolverOutcomeError, "source": "",
	}))
	require.Len(t, logs.FilterMessage("Failed to parse shard assignment settings; returning error").All(), 1, "one warning per new bad settings version")
	hash, ok := f.gaugeValue(t, "platform_shard_assignment_settings_hash")
	require.True(t, ok)
	require.Equal(t, firstHash, hash, "a parse error must not move the settings hash gauge")

	// The same broken version keeps failing resolution but is counted once.
	_, found, err = r.ResolveShard(ctx, "wf-a", ownerA)
	require.Error(t, err)
	require.False(t, found)
	require.Equal(t, int64(2), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverManual, "outcome": wfmon.ShardResolverOutcomeError, "source": "",
	}), "every failed resolve still counts one error")
	require.Equal(t, int64(1), f.counterTotal(t, "platform_shard_assignment_updated", map[string]string{"outcome": wfmon.ShardAssignmentUpdateParseError}),
		"the same broken settings version is counted once")
	require.Len(t, logs.FilterMessage("Failed to parse shard assignment settings; returning error").All(), 1)

	// The batch path propagates the same error, counting one error per call,
	// exactly like a single-shot call for the same failure.
	_, err = r.ResolveShards(ctx, []string{"wf-a"}, []string{ownerA})
	require.Error(t, err)
	require.Equal(t, int64(3), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverManual, "outcome": wfmon.ShardResolverOutcomeError, "source": "",
	}), "a failed batch call counts one error")

	// A good version applies again and resets the parse-error dedup.
	storeShardAssignment(t, settings, `
static_default_assignment = [0]

[per_owner_assignment]
  `+ownerA+` = [1]
`)
	donID, found, err = r.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, uint32(1), donID)
	require.Equal(t, int64(1), f.counterTotal(t, "platform_shard_assignment_updated", map[string]string{"outcome": wfmon.ShardAssignmentUpdateApplied}),
		"the same assignment returning after a broken version is not a new applied update")

	// The same broken version recurring after a good one counts and warns again.
	storeShardAssignment(t, settings, bad)
	_, found, err = r.ResolveShard(ctx, "wf-a", ownerA)
	require.Error(t, err)
	require.False(t, found)
	require.Equal(t, int64(2), f.counterTotal(t, "platform_shard_assignment_updated", map[string]string{"outcome": wfmon.ShardAssignmentUpdateParseError}),
		"a bad version recurring after a good one is counted again")
	require.Len(t, logs.FilterMessage("Failed to parse shard assignment settings; returning error").All(), 2)

	// A different broken version is its own parse error.
	storeShardAssignment(t, settings, `[per_owner_assignment]
  "0x1234" = [1]
`)
	_, _, err = r.ResolveShard(ctx, "wf-a", ownerA)
	require.Error(t, err)
	require.Equal(t, int64(3), f.counterTotal(t, "platform_shard_assignment_updated", map[string]string{"outcome": wfmon.ShardAssignmentUpdateParseError}),
		"a new broken settings version is counted again")
	require.Equal(t, int64(1), f.counterTotal(t, "platform_shard_assignment_updated", map[string]string{"outcome": wfmon.ShardAssignmentUpdateApplied}))
}

func TestResolverMetrics_RingRequests(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	f := newAssignmentMetricsFixture(t)

	r := NewRingOCRShardResolver(&stubClient{
		resp: &ringpb.GetWorkflowShardMappingResponse{Mappings: map[string]uint32{"wf-a": 1}},
	}, nopLogger).(*ringOCRShardResolver)
	r.metrics = f.metrics

	donID, found, err := r.ResolveShard(ctx, "wf-a", "")
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, uint32(1), donID)

	_, found, err = r.ResolveShard(ctx, "wf-b", "")
	require.NoError(t, err)
	require.False(t, found)

	mappings, err := r.ResolveShards(ctx, []string{"wf-a", "wf-b"}, nil)
	require.NoError(t, err)
	require.Equal(t, map[string]uint32{"wf-a": 1}, mappings)

	require.Equal(t, int64(2), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverRing, "outcome": wfmon.ShardResolverOutcomeFound, "source": "",
	}))
	require.Equal(t, int64(2), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverRing, "outcome": wfmon.ShardResolverOutcomeNotFound, "source": "",
	}))

	rErr := NewRingOCRShardResolver(&stubClient{err: errors.New("rpc down")}, nopLogger).(*ringOCRShardResolver)
	rErr.metrics = f.metrics
	_, _, err = rErr.ResolveShard(ctx, "wf-a", "")
	require.Error(t, err)
	_, err = rErr.ResolveShards(ctx, []string{"wf-a"}, nil)
	require.Error(t, err)
	require.Equal(t, int64(2), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverRing, "outcome": wfmon.ShardResolverOutcomeError, "source": "",
	}))
	require.Equal(t, uint64(5), f.histogramCount(t, "platform_shard_resolver_duration_seconds", map[string]string{"resolver": wfmon.ShardResolverRing}))
}

func TestResolverMetrics_OverrideRequests(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	f := newAssignmentMetricsFixture(t)

	settings := &loop.AtomicSettings{}
	storeShardAssignment(t, settings, `
static_default_assignment = [0]
hashed_default_assignment = true

[per_owner_assignment]
  `+ownerA+` = [1]
`)

	ring := NewRingOCRShardResolver(&stubClient{
		resp: &ringpb.GetWorkflowShardMappingResponse{Mappings: map[string]uint32{"wf-b": 2}},
	}, nopLogger).(*ringOCRShardResolver)
	ring.metrics = f.metrics

	o := NewOverrideShardResolver(settings, nil, nil, ring, nopLogger).(*overrideShardResolver)
	o.metrics = f.metrics
	o.manual.metrics = f.metrics

	// Single-shot: a manual hit counts once on the override resolver and must
	// not touch the ring resolver.
	donID, found, err := o.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, uint32(1), donID)
	require.Equal(t, int64(1), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverOverride, "outcome": wfmon.ShardResolverOutcomeFound, "source": "",
	}))
	require.Equal(t, int64(0), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{"resolver": wfmon.ShardResolverRing}),
		"a manual override hit must not go to the ring resolver")

	// Single-shot ring fallback: only the ring resolver counts, once.
	donID, found, err = o.ResolveShard(ctx, "wf-b", ownerB)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, uint32(2), donID)
	require.Equal(t, int64(1), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverRing, "outcome": wfmon.ShardResolverOutcomeFound, "source": "",
	}))
	require.Equal(t, int64(1), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{"resolver": wfmon.ShardResolverOverride}),
		"the ring fallback must not double-count on the override resolver")

	_, found, err = o.ResolveShard(ctx, "wf-c", ownerB)
	require.NoError(t, err)
	require.False(t, found)
	require.Equal(t, int64(1), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverRing, "outcome": wfmon.ShardResolverOutcomeNotFound, "source": "",
	}))

	// ResolveAllShards follows the same accounting on both sides.
	shards, found, err := o.ResolveAllShards(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, []uint32{1}, shards)
	shards, found, err = o.ResolveAllShards(ctx, "wf-b", ownerB)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, []uint32{2}, shards)
	require.Equal(t, int64(2), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverOverride, "outcome": wfmon.ShardResolverOutcomeFound, "source": "",
	}))

	// Batch: manual-decided workflows count under override, ring fallbacks
	// under ring, and a workflow without an owner counts as not-found — each
	// exactly once.
	resolved, err := o.ResolveShards(ctx, []string{"wf-a", "wf-b", "wf-tail"}, []string{ownerA, ownerB})
	require.NoError(t, err)
	require.Equal(t, map[string]uint32{"wf-a": 1, "wf-b": 2}, resolved)

	require.Equal(t, int64(3), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverOverride, "outcome": wfmon.ShardResolverOutcomeFound, "source": "",
	}), "one override count per manual-decided workflow, single-shot and batch agreeing")
	require.Equal(t, int64(1), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverOverride, "outcome": wfmon.ShardResolverOutcomeNotFound, "source": wfmon.ShardResolverSourceNone,
	}), "the workflow without an owner counts as not-found with source=none")
	require.Equal(t, int64(3), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverRing, "outcome": wfmon.ShardResolverOutcomeFound, "source": "",
	}), "the ring fallback inside the batch is counted on the ring resolver")
	require.Equal(t, uint64(3), f.histogramCount(t, "platform_shard_resolver_duration_seconds", map[string]string{"resolver": wfmon.ShardResolverOverride}),
		"one override duration per override-decided call")
	require.Equal(t, uint64(4), f.histogramCount(t, "platform_shard_resolver_duration_seconds", map[string]string{"resolver": wfmon.ShardResolverRing}),
		"one ring duration per ring-fallback call")

	// With no manual settings at all, the whole batch is decided (and counted)
	// by the ring resolver; the override resolver adds nothing.
	empty := &loop.AtomicSettings{}
	o2 := NewOverrideShardResolver(empty, nil, nil, ring, nopLogger).(*overrideShardResolver)
	o2.metrics = f.metrics
	o2.manual.metrics = f.metrics
	resolved, err = o2.ResolveShards(ctx, []string{"wf-b"}, []string{ownerB})
	require.NoError(t, err)
	require.Equal(t, map[string]uint32{"wf-b": 2}, resolved)
	require.Equal(t, int64(4), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverRing, "outcome": wfmon.ShardResolverOutcomeFound, "source": "",
	}), "a nil manual config delegates the whole batch to the ring resolver")
	require.Equal(t, int64(4), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{"resolver": wfmon.ShardResolverOverride}),
		"no override counts are recorded for workflows the ring decided")
	require.Equal(t, uint64(3), f.histogramCount(t, "platform_shard_resolver_duration_seconds", map[string]string{"resolver": wfmon.ShardResolverOverride}),
		"a fully delegated batch adds no override duration either")
}

func TestResolverMetrics_SettingsDiffLog(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	f := newAssignmentMetricsFixture(t)

	lggr, logs := logger.TestObserved(t, zapcore.InfoLevel)
	settings := &loop.AtomicSettings{}
	storeShardAssignment(t, settings, `
static_default_assignment = [0]

[per_owner_assignment]
  `+ownerA+` = [1]
  `+ownerB+` = [2]
  `+ownerC+` = [5]
`)
	r := f.manualResolver(settings, nil, lggr)

	_, _, err := r.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	entries := logs.FilterMessage("Shard assignment settings applied").All()
	require.Len(t, entries, 1, "one INFO line per settings update")
	require.Equal(t, true, entries[0].ContextMap()["initial"])
	require.Equal(t, int64(3), entries[0].ContextMap()["perOwnerCount"])

	_, _, err = r.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	require.Len(t, logs.FilterMessage("Shard assignment settings applied").All(), 1,
		"re-observing identical settings must not log another line")

	storeShardAssignment(t, settings, `
static_default_assignment = [0]

[per_owner_assignment]
  `+ownerA+` = [3]
  `+ownerB+` = [2]
  `+ownerD+` = [4]
`)
	_, _, err = r.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	entries = logs.FilterMessage("Shard assignment settings applied").All()
	require.Len(t, entries, 2)
	fields := entries[1].ContextMap()
	require.Equal(t, []any{ownerA[2:] + "=[1]->[3]"}, fields["perOwnerMoved"])
	require.Equal(t, []any{ownerD[2:] + "=[4]"}, fields["perOwnerAdded"])
	require.Equal(t, []any{ownerC[2:]}, fields["perOwnerRemoved"])
	require.Equal(t, int64(3), fields["perOwnerCount"])
	require.NotContains(t, fields["perOwnerMoved"], ownerB[2:])
	require.NotContains(t, fields["perOwnerAdded"], ownerB[2:])
	require.NotContains(t, fields["perOwnerRemoved"], ownerB[2:])

	_, _, err = r.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	require.Len(t, logs.FilterMessage("Shard assignment settings applied").All(), 2)
}

func TestResolverMetrics_ManualBatchCountsEveryWorkflow(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	f := newAssignmentMetricsFixture(t)

	settings := &loop.AtomicSettings{}
	storeShardAssignment(t, settings, `
static_default_assignment = [0]

[per_owner_assignment]
  `+ownerA+` = [1]
`)
	r := f.manualResolver(settings, nil, nopLogger)

	// More workflows than owners: the ownerless tail still counts as
	// not-found with source=none and stays absent from the result.
	resolved, err := r.ResolveShards(ctx, []string{"wf-a", "wf-tail-1", "wf-tail-2"}, []string{ownerA})
	require.NoError(t, err)
	require.Equal(t, map[string]uint32{"wf-a": 1}, resolved)
	require.Equal(t, int64(1), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverManual, "outcome": wfmon.ShardResolverOutcomeFound, "source": wfmon.ShardResolverSourcePerOwner,
	}))
	require.Equal(t, int64(2), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverManual, "outcome": wfmon.ShardResolverOutcomeNotFound, "source": wfmon.ShardResolverSourceNone,
	}), "workflows without an owner are counted as not-found, not dropped")

	// A batch hitting an unparseable settings version fails like single-shot
	// and counts one error for the call.
	storeShardAssignment(t, settings, `static_default_assignment = "not-an-array"`)
	_, err = r.ResolveShards(ctx, []string{"wf-a"}, []string{ownerA})
	require.Error(t, err)
	require.Equal(t, int64(1), f.counterTotal(t, "platform_shard_resolver_requests", map[string]string{
		"resolver": wfmon.ShardResolverManual, "outcome": wfmon.ShardResolverOutcomeError, "source": "",
	}), "a failed batch call counts one error, matching a single-shot call")
}

// TestResolverMetrics_SettingsDiffLogBaselineIsGlobal proves the applied-diff
// baseline lives with the global dedup state: a second resolver that never
// observed the previous version still logs the real diff, not "initial".
func TestResolverMetrics_SettingsDiffLogBaselineIsGlobal(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	f := newAssignmentMetricsFixture(t)

	lggr, logs := logger.TestObserved(t, zapcore.InfoLevel)
	settings := &loop.AtomicSettings{}
	storeShardAssignment(t, settings, `
static_default_assignment = [0]

[per_owner_assignment]
  `+ownerA+` = [1]
`)
	r1 := f.manualResolver(settings, nil, lggr)
	_, _, err := r1.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	require.Len(t, logs.FilterMessage("Shard assignment settings applied").All(), 1)

	storeShardAssignment(t, settings, `
static_default_assignment = [0]

[per_owner_assignment]
  `+ownerA+` = [3]
`)
	r2 := f.manualResolver(settings, nil, lggr)
	_, _, err = r2.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	entries := logs.FilterMessage("Shard assignment settings applied").All()
	require.Len(t, entries, 2)
	fields := entries[1].ContextMap()
	require.NotContains(t, fields, "initial", "the diff baseline is the globally applied config, not the new resolver's nil baseline")
	require.Equal(t, []any{ownerA[2:] + "=[1]->[3]"}, fields["perOwnerMoved"])
}

// TestResolverMetrics_SettingsDiffLogCapsOwnerLists proves each owner list in
// the applied log is capped, so a large owner map cannot produce an
// unbounded log line.
func TestResolverMetrics_SettingsDiffLogCapsOwnerLists(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	f := newAssignmentMetricsFixture(t)

	lggr, logs := logger.TestObserved(t, zapcore.InfoLevel)
	settings := &loop.AtomicSettings{}
	storeShardAssignment(t, settings, `static_default_assignment = [0]`)
	r := f.manualResolver(settings, nil, lggr)
	_, _, err := r.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)

	var cfgToml strings.Builder
	cfgToml.WriteString("static_default_assignment = [0]\n\n[per_owner_assignment]\n")
	for i := range assignmentDiffOwnerCap + 2 {
		fmt.Fprintf(&cfgToml, "  \"0x%040x\" = [1]\n", i+1)
	}
	storeShardAssignment(t, settings, cfgToml.String())
	_, _, err = r.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)

	entries := logs.FilterMessage("Shard assignment settings applied").All()
	require.Len(t, entries, 2)
	added, ok := entries[1].ContextMap()["perOwnerAdded"].([]any)
	require.True(t, ok)
	require.Len(t, added, assignmentDiffOwnerCap+1, "the added list is capped at the first 10 owners plus a marker")
	require.Equal(t, fmt.Sprintf("...and %d more", 2), added[len(added)-1])
}

// TestResolverMetrics_ConcurrentResolveWhileSettingsChange exercises the
// settings memo and metric dedup under concurrent resolves with the settings
// changing underneath (run with -race).
func TestResolverMetrics_ConcurrentResolveWhileSettingsChange(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	f := newAssignmentMetricsFixture(t)

	goodA := `
static_default_assignment = [0]

[per_owner_assignment]
  ` + ownerA + ` = [1]
`
	goodB := `
static_default_assignment = [0]

[per_owner_assignment]
  ` + ownerA + ` = [2]
`
	settings := &loop.AtomicSettings{}
	storeShardAssignment(t, settings, goodA)
	r := f.manualResolver(settings, nil, nopLogger)

	var wg sync.WaitGroup
	for range 4 {
		wg.Go(func() {
			for range 50 {
				_, _, _ = r.ResolveShard(ctx, "wf-a", ownerA)
				_, _ = r.ResolveShards(ctx, []string{"wf-a"}, []string{ownerA})
			}
		})
	}
	wg.Go(func() {
		for range 25 {
			_ = settings.Store(core.SettingsUpdate{Settings: goodB, Hash: "test"})
			_ = settings.Store(core.SettingsUpdate{Settings: goodA, Hash: "test"})
		}
	})
	wg.Wait()

	cfgA, err := cresettings.ParseShardAssignmentConfig(goodA)
	require.NoError(t, err)
	cfgB, err := cresettings.ParseShardAssignmentConfig(goodB)
	require.NoError(t, err)
	hash, ok := f.gaugeValue(t, "platform_shard_assignment_settings_hash")
	require.True(t, ok)
	require.Contains(t, []int64{int64(shardAssignmentSettingsHash(cfgA)), int64(shardAssignmentSettingsHash(cfgB))}, hash)
}
