package shardownership

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func makePeerID(b byte) ragetypes.PeerID {
	var id ragetypes.PeerID
	id[0] = b
	return id
}

// testDON describes a DON to seed into a test registry. members[0] (if set)
// is treated as "my" peer when localPeerID matches it.
type testDON struct {
	id               uint32
	name             string
	families         []string
	acceptsWorkflows bool
	members          []ragetypes.PeerID
}

// newTestRegistry builds a minimally valid *registry.RegistryMetadata (so
// LocalNode/NodeByPeerID don't hit ensureNotEmpty or nil-Logger panics),
// with localPeerID as the local node's peer ID.
func newTestRegistry(t *testing.T, localPeerID ragetypes.PeerID, dons ...testDON) *registry.RegistryMetadata {
	t.Helper()

	idsToDONs := make(map[registry.DonID]registry.DON, len(dons))
	for _, d := range dons {
		idsToDONs[registry.DonID(d.id)] = registry.DON{
			ID:               d.id,
			Name:             d.name,
			Families:         d.families,
			AcceptsWorkflows: d.acceptsWorkflows,
			Members:          d.members,
		}
	}

	reg := registry.NewRegistryMetadata(
		logger.Test(t),
		func() (ragetypes.PeerID, error) { return localPeerID, nil },
		idsToDONs,
		map[ragetypes.PeerID]registry.NodeInfo{localPeerID: {}},
		map[string]registry.Capability{"cron-trigger@1.0.0": {}},
	)
	return &reg
}

func TestShardIndexMapper_DonByShardIndex(t *testing.T) {
	t.Parallel()

	me := makePeerID(1)
	reg := newTestRegistry(t, me,
		testDON{id: 10, name: "workflow-1-zone-a", families: []string{"zone-a_shard-0", "zone-a"}, acceptsWorkflows: true, members: []ragetypes.PeerID{me}},
		testDON{id: 30, name: "workflow-1-zone-a-shard-2", families: []string{"zone-a_shard-2", "zone-a"}, acceptsWorkflows: true},
		testDON{id: 20, name: "workflow-1-zone-a-shard-1", families: []string{"zone-a_shard-1", "zone-a"}, acceptsWorkflows: true},
		// Different family entirely: must be excluded even though it accepts workflows.
		testDON{id: 40, name: "workflow-2-zone-b", families: []string{"zone-b_shard-0", "zone-b"}, acceptsWorkflows: true},
		// Capability-only DON in the same family: must be excluded.
		testDON{id: 99, name: "chain-capabilities-zone-a", families: []string{"zone-a"}, acceptsWorkflows: false},
	)

	idx := NewShardIndexMapper(logger.Test(t))
	require.NoError(t, idx.OnNewRegistry(t.Context(), reg))

	assert.Equal(t, uint32(10), idx.DonByShardIndex(t.Context(), 0).ID)
	assert.Equal(t, uint32(20), idx.DonByShardIndex(t.Context(), 1).ID)
	assert.Equal(t, uint32(30), idx.DonByShardIndex(t.Context(), 2).ID)
	assert.Nil(t, idx.DonByShardIndex(t.Context(), 3))
}

func TestShardIndexMapper_ExcludesSameFamilyDifferentNamePrefix(t *testing.T) {
	t.Parallel()

	me := makePeerID(1)
	reg := newTestRegistry(t, me,
		testDON{id: 10, name: "workflow-1-zone-a", families: []string{"zone-a_shard-0", "zone-a"}, acceptsWorkflows: true, members: []ragetypes.PeerID{me}},
		testDON{id: 20, name: "workflow-1-zone-a-shard-1", families: []string{"zone-a_shard-1", "zone-a"}, acceptsWorkflows: true},
		// Same family as me, but a different name prefix: a distinct shard group that happens to share the family. Must be excluded.
		testDON{id: 40, name: "workflow-2-zone-a-shard-1", families: []string{"zone-a_shard-1", "zone-a"}, acceptsWorkflows: true},
	)

	idx := NewShardIndexMapper(logger.Test(t))
	require.NoError(t, idx.OnNewRegistry(t.Context(), reg))

	assert.Equal(t, uint32(10), idx.DonByShardIndex(t.Context(), 0).ID)
	assert.Equal(t, uint32(20), idx.DonByShardIndex(t.Context(), 1).ID)
	assert.Nil(t, idx.DonByShardIndex(t.Context(), 2))
}

func TestShardIndexMapper_MatchesOnPartiallyOverlappingFamilies(t *testing.T) {
	t.Parallel()

	me := makePeerID(1)
	reg := newTestRegistry(t, me,
		// Local DON and its shard peer share only "zone-a"; each also has its own
		// shard-specific family that the other doesn't have. A single shared
		// family is enough to match.
		testDON{id: 10, name: "workflow-1-zone-a", families: []string{"zone-a_shard-0", "zone-a", "extra-family"}, acceptsWorkflows: true, members: []ragetypes.PeerID{me}},
		testDON{id: 20, name: "workflow-1-zone-a-shard-1", families: []string{"zone-a_shard-1", "zone-a"}, acceptsWorkflows: true},
	)

	idx := NewShardIndexMapper(logger.Test(t))
	require.NoError(t, idx.OnNewRegistry(t.Context(), reg))

	assert.Equal(t, uint32(10), idx.DonByShardIndex(t.Context(), 0).ID)
	assert.Equal(t, uint32(20), idx.DonByShardIndex(t.Context(), 1).ID)
}

func TestShardIndexMapper_ExcludesSameNamePrefixNoSharedFamily(t *testing.T) {
	t.Parallel()

	me := makePeerID(1)
	reg := newTestRegistry(t, me,
		testDON{id: 10, name: "workflow-1-zone-a", families: []string{"zone-a_shard-0", "zone-a"}, acceptsWorkflows: true, members: []ragetypes.PeerID{me}},
		// Same name prefix, but no family in common with the local DON: must be excluded.
		testDON{id: 20, name: "workflow-1-zone-a-shard-1", families: []string{"zone-b"}, acceptsWorkflows: true},
	)

	idx := NewShardIndexMapper(logger.Test(t))
	require.NoError(t, idx.OnNewRegistry(t.Context(), reg))

	assert.Equal(t, uint32(10), idx.DonByShardIndex(t.Context(), 0).ID)
	assert.Nil(t, idx.DonByShardIndex(t.Context(), 1))
}

func TestShardIndexMapper_NoShardSuffixMeansIndexZero(t *testing.T) {
	t.Parallel()

	me := makePeerID(1)
	reg := newTestRegistry(t, me,
		testDON{id: 10, name: "workflow-1-zone-a", families: []string{"zone-a_shard-0", "zone-a"}, acceptsWorkflows: true, members: []ragetypes.PeerID{me}},
	)

	idx := NewShardIndexMapper(logger.Test(t))
	require.NoError(t, idx.OnNewRegistry(t.Context(), reg))
	require.NotNil(t, idx.DonByShardIndex(t.Context(), 0))
	assert.Equal(t, uint32(10), idx.DonByShardIndex(t.Context(), 0).ID)
}

func TestShardIndexMapper_RefreshesOnNewRegistry(t *testing.T) {
	t.Parallel()

	me := makePeerID(1)
	idx := NewShardIndexMapper(logger.Test(t))

	require.NoError(t, idx.OnNewRegistry(t.Context(), newTestRegistry(t, me,
		testDON{id: 5, name: "workflow-1-zone-a", families: []string{"zone-a_shard-0", "zone-a"}, acceptsWorkflows: true, members: []ragetypes.PeerID{me}},
	)))
	require.NotNil(t, idx.DonByShardIndex(t.Context(), 0))
	assert.Nil(t, idx.DonByShardIndex(t.Context(), 1))

	// A later registry snapshot with a different shard layout replaces the stale entry.
	require.NoError(t, idx.OnNewRegistry(t.Context(), newTestRegistry(t, me,
		testDON{id: 7, name: "workflow-1-zone-a", families: []string{"zone-a_shard-0", "zone-a"}, acceptsWorkflows: true, members: []ragetypes.PeerID{me}},
		testDON{id: 8, name: "workflow-1-zone-a-shard-1", families: []string{"zone-a_shard-1", "zone-a"}, acceptsWorkflows: true},
	)))
	assert.Equal(t, uint32(7), idx.DonByShardIndex(t.Context(), 0).ID)
	assert.Equal(t, uint32(8), idx.DonByShardIndex(t.Context(), 1).ID)
}

func TestShardIndexMapper_ErrorsWhenLocalDonHasNoFamilies(t *testing.T) {
	t.Parallel()

	me := makePeerID(1)
	reg := newTestRegistry(t, me,
		testDON{id: 10, name: "workflow-1-zone-a", families: nil, acceptsWorkflows: true, members: []ragetypes.PeerID{me}},
	)

	idx := NewShardIndexMapper(logger.Test(t))
	err := idx.OnNewRegistry(t.Context(), reg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no workflow DONs found")
}

func TestShardIndexMapper_ErrorsWhenLocalNodeNotInAnyWorkflowDon(t *testing.T) {
	t.Parallel()

	me := makePeerID(1)
	other := makePeerID(2)
	reg := newTestRegistry(t, me,
		testDON{id: 10, name: "workflow-1-zone-a", families: []string{"zone-a_shard-0", "zone-a"}, acceptsWorkflows: true, members: []ragetypes.PeerID{other}},
	)

	idx := NewShardIndexMapper(logger.Test(t))
	err := idx.OnNewRegistry(t.Context(), reg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not belong to a workflow DON")
}

func TestShardIndexMapper_ErrorsOnMultiDigitShardSuffix(t *testing.T) {
	t.Parallel()

	me := makePeerID(1)
	reg := newTestRegistry(t, me,
		testDON{id: 10, name: "workflow-1-zone-a", families: []string{"zone-a_shard-0", "zone-a"}, acceptsWorkflows: true, members: []ragetypes.PeerID{me}},
		testDON{id: 11, name: "workflow-1-zone-a-shard-12", families: []string{"zone-a"}, acceptsWorkflows: true},
	)

	idx := NewShardIndexMapper(logger.Test(t))
	err := idx.OnNewRegistry(t.Context(), reg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "single digit")
}

func TestShardIndexMapper_ErrorsOnDuplicateShardIndex(t *testing.T) {
	t.Parallel()

	me := makePeerID(1)
	reg := newTestRegistry(t, me,
		testDON{id: 10, name: "workflow-1-zone-a", families: []string{"zone-a_shard-0", "zone-a"}, acceptsWorkflows: true, members: []ragetypes.PeerID{me}},
		testDON{id: 11, name: "workflow-1-zone-a", families: []string{"zone-a_shard-0", "zone-a"}, acceptsWorkflows: true},
	)

	idx := NewShardIndexMapper(logger.Test(t))
	err := idx.OnNewRegistry(t.Context(), reg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "two DONs at shard index")
}

func TestShardIndexMapper_WaitReady(t *testing.T) {
	t.Parallel()

	me := makePeerID(1)
	idx := NewShardIndexMapper(logger.Test(t))

	// No registry snapshot yet: WaitReady must block until ctx is done.
	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Millisecond)
	defer cancel()
	require.ErrorIs(t, idx.WaitReady(ctx), context.DeadlineExceeded)

	// Once a snapshot arrives, WaitReady returns immediately, including for
	// callers that started waiting before the snapshot was processed.
	done := make(chan error, 1)
	go func() {
		done <- idx.WaitReady(t.Context())
	}()

	require.NoError(t, idx.OnNewRegistry(t.Context(), newTestRegistry(t, me,
		testDON{id: 1, name: "workflow-1-zone-a", families: []string{"zone-a_shard-0", "zone-a"}, acceptsWorkflows: true, members: []ragetypes.PeerID{me}},
	)))

	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("WaitReady did not unblock after OnNewRegistry")
	}

	require.NoError(t, idx.WaitReady(t.Context()))
}

func TestShardIndexFromName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		donName   string
		wantIndex uint32
		wantErr   bool
	}{
		{name: "no suffix", donName: "workflow-1-zone-a", wantIndex: 0},
		{name: "shard 1", donName: "workflow-1-zone-a-shard-1", wantIndex: 1},
		{name: "shard 9", donName: "workflow-1-zone-a-shard-9", wantIndex: 9},
		{name: "multi-digit rejected", donName: "workflow-1-zone-a-shard-12", wantErr: true},
		{name: "non-digit rejected", donName: "workflow-1-zone-a-shard-x", wantErr: true},
		{name: "trailing hyphen rejected", donName: "workflow-1-zone-a-shard-", wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := shardIndexFromName(tc.donName)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantIndex, got)
		})
	}
}

func TestShardGroupNamePrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		donName    string
		wantPrefix string
	}{
		{name: "no suffix", donName: "workflow-1-zone-a", wantPrefix: "workflow-1-zone-a"},
		{name: "hyphen shard 1", donName: "workflow-1-zone-a-shard-1", wantPrefix: "workflow-1-zone-a"},
		{name: "hyphen shard 9", donName: "workflow-1-zone-a-shard-9", wantPrefix: "workflow-1-zone-a"},
		{name: "underscore shard 1", donName: "workflow-1-zone-a_shard-1", wantPrefix: "workflow-1-zone-a"},
		{name: "malformed suffix still stripped", donName: "workflow-1-zone-a-shard-12", wantPrefix: "workflow-1-zone-a"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.wantPrefix, shardGroupNamePrefix(tc.donName))
		})
	}
}
