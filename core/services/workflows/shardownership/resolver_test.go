package shardownership

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

var nopLogger = logger.Nop()

func storeShardAssignment(t *testing.T, settings *loop.AtomicSettings, toml string) {
	t.Helper()
	require.NoError(t, settings.Store(core.SettingsUpdate{Settings: toml, Hash: "test"}))
}

func TestManualShardResolver_PerOwnerAssignment(t *testing.T) {
	t.Parallel()
	settings := &loop.AtomicSettings{}
	storeShardAssignment(t, settings, `
static_default_assignment = [0]

[per_owner_assignment]
  "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266" = [1, 0]
  "0x70997970c51812dc3a010c7d01b50e0d17dc79c8" = [0, 2]
`)

	r := NewManualShardResolver(settings, nopLogger)

	shard, found, err := r.ResolveShard(context.Background(), "wf-1", "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, uint32(1), shard)

	shard, found, err = r.ResolveShard(context.Background(), "wf-2", "0x70997970c51812dc3a010c7d01b50e0d17dc79c8")
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, uint32(0), shard)
}

func TestManualShardResolver_StaticDefaultFallback(t *testing.T) {
	t.Parallel()
	settings := &loop.AtomicSettings{}
	storeShardAssignment(t, settings, `
static_default_assignment = [0]
`)

	r := NewManualShardResolver(settings, nopLogger)

	shard, found, err := r.ResolveShard(context.Background(), "wf-1", "0x0000000000000000000000000000000000000001")
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, uint32(0), shard)
}

func TestManualShardResolver_HashedOwnerDelegatesToRingOCR(t *testing.T) {
	t.Parallel()
	settings := &loop.AtomicSettings{}
	storeShardAssignment(t, settings, `
hashed_owner_assignment = ["0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"]
`)

	r := NewManualShardResolver(settings, nopLogger)

	shard, found, err := r.ResolveShard(context.Background(), "wf-1", "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")
	require.NoError(t, err)
	assert.False(t, found)
	assert.Equal(t, uint32(0), shard)
}

func TestManualShardResolver_HashedDefaultDelegatesToRingOCR(t *testing.T) {
	t.Parallel()
	settings := &loop.AtomicSettings{}
	storeShardAssignment(t, settings, `
hashed_default_assignment = true
static_default_assignment = [0]
`)

	r := NewManualShardResolver(settings, nopLogger)

	_, found, err := r.ResolveShard(context.Background(), "wf-1", "0x0000000000000000000000000000000000000001")
	require.NoError(t, err)
	assert.False(t, found)
}

func TestManualShardResolver_NilConfig(t *testing.T) {
	t.Parallel()
	settings := &loop.AtomicSettings{}
	r := NewManualShardResolver(settings, nopLogger)

	shard, found, err := r.ResolveShard(context.Background(), "wf-1", "0x0000000000000000000000000000000000000001")
	require.NoError(t, err)
	assert.False(t, found)
	assert.Equal(t, uint32(0), shard)
}

func TestManualShardResolver_ResolveShards(t *testing.T) {
	t.Parallel()
	settings := &loop.AtomicSettings{}
	storeShardAssignment(t, settings, `
static_default_assignment = [0]

[per_owner_assignment]
  "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266" = [1]
`)

	r := NewManualShardResolver(settings, nopLogger)

	wfIDs := []string{"wf-1", "wf-2"}
	owners := []string{"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266", "0x0000000000000000000000000000000000000001"}
	result, err := r.ResolveShards(context.Background(), wfIDs, owners)
	require.NoError(t, err)
	assert.Equal(t, uint32(1), result["wf-1"])
	assert.Equal(t, uint32(0), result["wf-2"])
}

func TestOverrideShardResolver_ManualWins(t *testing.T) {
	t.Parallel()
	settings := &loop.AtomicSettings{}
	storeShardAssignment(t, settings, `
static_default_assignment = [0]

[per_owner_assignment]
  "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266" = [1]
`)

	mockRing := &mockRingOCRResolver{mappings: map[string]uint32{"wf-1": 0}}
	r := NewOverrideShardResolver(settings, &shardResolverAdapter{inner: mockRing}, nopLogger)

	shard, found, err := r.ResolveShard(context.Background(), "wf-1", "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, uint32(1), shard)
}

func TestOverrideShardResolver_RingOCRWinsForHashed(t *testing.T) {
	t.Parallel()
	settings := &loop.AtomicSettings{}
	storeShardAssignment(t, settings, `
hashed_default_assignment = true
`)

	mockRing := &mockRingOCRResolver{mappings: map[string]uint32{"wf-1": 1}}
	r := NewOverrideShardResolver(settings, &shardResolverAdapter{inner: mockRing}, nopLogger)

	shard, found, err := r.ResolveShard(context.Background(), "wf-1", "0x0000000000000000000000000000000000000001")
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, uint32(1), shard)
}

func TestOverrideShardResolver_ResolveShards_Mixed(t *testing.T) {
	t.Parallel()
	settings := &loop.AtomicSettings{}
	storeShardAssignment(t, settings, `
static_default_assignment = [0]
hashed_default_assignment = true

[per_owner_assignment]
  "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266" = [1]
`)

	mockRing := &mockRingOCRResolver{mappings: map[string]uint32{"wf-2": 0}}
	r := NewOverrideShardResolver(settings, &shardResolverAdapter{inner: mockRing}, nopLogger)

	wfIDs := []string{"wf-1", "wf-2"}
	owners := []string{"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266", "0x0000000000000000000000000000000000000001"}
	result, err := r.ResolveShards(context.Background(), wfIDs, owners)
	require.NoError(t, err)
	assert.Equal(t, uint32(1), result["wf-1"])
	assert.Equal(t, uint32(0), result["wf-2"])
}

type mockRingOCRResolver struct {
	mappings map[string]uint32
}

func (m *mockRingOCRResolver) ResolveShard(_ context.Context, workflowID string, _ string) (uint32, bool, error) {
	s, ok := m.mappings[workflowID]
	return s, ok, nil
}

func (m *mockRingOCRResolver) ResolveShards(_ context.Context, workflowIDs []string, _ []string) (map[string]uint32, error) {
	result := make(map[string]uint32, len(workflowIDs))
	for _, id := range workflowIDs {
		if s, ok := m.mappings[id]; ok {
			result[id] = s
		}
	}
	return result, nil
}

type shardResolverAdapter struct {
	inner *mockRingOCRResolver
}

func (a *shardResolverAdapter) ResolveShard(ctx context.Context, workflowID string, ownerHex string) (uint32, bool, error) {
	return a.inner.ResolveShard(ctx, workflowID, ownerHex)
}

func (a *shardResolverAdapter) ResolveShards(ctx context.Context, workflowIDs []string, ownerHexes []string) (map[string]uint32, error) {
	return a.inner.ResolveShards(ctx, workflowIDs, ownerHexes)
}

var _ ShardResolver = (*shardResolverAdapter)(nil)
