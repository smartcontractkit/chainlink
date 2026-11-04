package shardownership

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
)

// TestResolverMetrics_SettingsArrivalToApply proves the arrival wiring end to
// end: a settings Store is observed by the resolver's settings subscription,
// and the next resolve applies it, bumping the config update count and
// recording exactly one arrival→in-effect staleness sample per actual
// settings change.
func TestResolverMetrics_SettingsArrivalToApply(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	f := newAssignmentMetricsFixture(t)

	settings := &loop.AtomicSettings{}
	r := f.manualResolver(settings, nil, nopLogger)

	blobA := `
static_default_assignment = [0]

[per_owner_assignment]
  ` + ownerA + ` = [1]
`
	storeShardAssignment(t, settings, blobA)
	// The subscription goroutine records the arrival asynchronously; give it
	// a beat so the arrival is observed before the resolve applies the blob.
	time.Sleep(100 * time.Millisecond)
	donID, found, err := r.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, uint32(1), donID)
	updateCount, ok := f.gaugeValue(t, "platform_shard_assignment_config_update_count")
	require.True(t, ok)
	require.Equal(t, int64(1), updateCount)
	require.Equal(t, uint64(1), f.histogramCount(t, "platform_shard_assignment_config_staleness_seconds", nil))

	// Re-observing the same blob records nothing.
	_, _, err = r.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	require.Equal(t, uint64(1), f.histogramCount(t, "platform_shard_assignment_config_staleness_seconds", nil))

	blobB := `
static_default_assignment = [0]

[per_owner_assignment]
  ` + ownerA + ` = [2]
`
	storeShardAssignment(t, settings, blobB)
	time.Sleep(100 * time.Millisecond)
	donID, found, err = r.ResolveShard(ctx, "wf-a", ownerA)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, uint32(2), donID)
	updateCount, ok = f.gaugeValue(t, "platform_shard_assignment_config_update_count")
	require.True(t, ok)
	require.Equal(t, int64(2), updateCount, "a real settings change bumps the config update count")
	require.Equal(t, uint64(2), f.histogramCount(t, "platform_shard_assignment_config_staleness_seconds", nil),
		"one staleness sample per actual settings change")
}
