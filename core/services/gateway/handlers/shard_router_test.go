package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	handlermocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
)

func TestShardRouter_SendToNode_RoutesToCorrectShard(t *testing.T) {
	t.Parallel()

	shard0DON := handlermocks.NewDON(t)
	shard1DON := handlermocks.NewDON(t)

	shardedDONs := []config.ShardedDONConfig{{
		DonName: "myDON",
		F:       1,
		Shards: []config.Shard{
			{Nodes: []config.NodeConfig{
				{Name: "s0n0", Address: "0xaaaa"},
				{Name: "s0n1", Address: "0xbbbb"},
			}},
			{Nodes: []config.NodeConfig{
				{Name: "s1n0", Address: "0xcccc"},
				{Name: "s1n1", Address: "0xdddd"},
			}},
		},
	}}
	connMgrs := [][]handlers.DON{{shard0DON, shard1DON}}

	router, err := handlers.NewShardRouter(shardedDONs, connMgrs)
	require.NoError(t, err)

	ctx := testutils.Context(t)

	shard0DON.EXPECT().SendToNode(mock.Anything, "0xaaaa", mock.Anything).Return(nil).Once()
	require.NoError(t, router.SendToNode(ctx, "0xaaaa", nil))

	shard1DON.EXPECT().SendToNode(mock.Anything, "0xcccc", mock.Anything).Return(nil).Once()
	require.NoError(t, router.SendToNode(ctx, "0xcccc", nil))
}

func TestShardRouter_SendToNode_UnknownNodeReturnsError(t *testing.T) {
	t.Parallel()

	shardedDONs := []config.ShardedDONConfig{{
		DonName: "myDON",
		F:       0,
		Shards: []config.Shard{
			{Nodes: []config.NodeConfig{{Name: "n0", Address: "0xaaaa"}}},
		},
	}}
	connMgrs := [][]handlers.DON{{handlermocks.NewDON(t)}}

	router, err := handlers.NewShardRouter(shardedDONs, connMgrs)
	require.NoError(t, err)

	err = router.SendToNode(testutils.Context(t), "0xunknown", nil)
	require.ErrorContains(t, err, "not found in any shard")
}

func TestShardRouter_DuplicateNodeAddressAcrossShardsErrors(t *testing.T) {
	t.Parallel()

	shardedDONs := []config.ShardedDONConfig{{
		DonName: "myDON",
		F:       0,
		Shards: []config.Shard{
			{Nodes: []config.NodeConfig{{Name: "n0", Address: "0xaaaa"}}},
			{Nodes: []config.NodeConfig{{Name: "n1", Address: "0xaaaa"}}},
		},
	}}
	connMgrs := [][]handlers.DON{{handlermocks.NewDON(t), handlermocks.NewDON(t)}}

	_, err := handlers.NewShardRouter(shardedDONs, connMgrs)
	require.ErrorContains(t, err, "duplicate node address")
}

func TestShardRouter_MultipleDONs(t *testing.T) {
	t.Parallel()

	don1Shard0 := handlermocks.NewDON(t)
	don2Shard0 := handlermocks.NewDON(t)

	shardedDONs := []config.ShardedDONConfig{
		{
			DonName: "don1",
			F:       0,
			Shards:  []config.Shard{{Nodes: []config.NodeConfig{{Name: "d1n0", Address: "0x1111"}}}},
		},
		{
			DonName: "don2",
			F:       0,
			Shards:  []config.Shard{{Nodes: []config.NodeConfig{{Name: "d2n0", Address: "0x2222"}}}},
		},
	}
	connMgrs := [][]handlers.DON{{don1Shard0}, {don2Shard0}}

	router, err := handlers.NewShardRouter(shardedDONs, connMgrs)
	require.NoError(t, err)

	ctx := testutils.Context(t)

	don1Shard0.EXPECT().SendToNode(mock.Anything, "0x1111", mock.Anything).Return(nil).Once()
	require.NoError(t, router.SendToNode(ctx, "0x1111", nil))

	don2Shard0.EXPECT().SendToNode(mock.Anything, "0x2222", mock.Anything).Return(nil).Once()
	require.NoError(t, router.SendToNode(ctx, "0x2222", nil))
}

func TestShardRouter_CaseInsensitiveAddressLookup(t *testing.T) {
	t.Parallel()

	mockDON := handlermocks.NewDON(t)

	shardedDONs := []config.ShardedDONConfig{{
		DonName: "myDON",
		F:       0,
		Shards:  []config.Shard{{Nodes: []config.NodeConfig{{Name: "n0", Address: "0xAaBb"}}}},
	}}
	connMgrs := [][]handlers.DON{{mockDON}}

	router, err := handlers.NewShardRouter(shardedDONs, connMgrs)
	require.NoError(t, err)

	mockDON.EXPECT().SendToNode(mock.Anything, "0xaabb", mock.Anything).Return(nil).Once()
	require.NoError(t, router.SendToNode(testutils.Context(t), "0xaabb", nil))
}

func TestShardRouter_ShardTopology(t *testing.T) {
	t.Parallel()

	shard0DON := handlermocks.NewDON(t)
	shard1DON := handlermocks.NewDON(t)
	shard2DON := handlermocks.NewDON(t)

	shardedDONs := []config.ShardedDONConfig{
		{
			DonName: "don1",
			F:       1,
			Shards: []config.Shard{
				{Nodes: []config.NodeConfig{
					{Name: "d1s0n0", Address: "0x1"},
					{Name: "d1s0n1", Address: "0x2"},
				}},
				{Nodes: []config.NodeConfig{
					{Name: "d1s1n0", Address: "0x3"},
				}},
			},
		},
		{
			DonName: "don2",
			F:       2,
			Shards: []config.Shard{
				{Nodes: []config.NodeConfig{
					{Name: "d2s0n0", Address: "0x4"},
				}},
			},
		},
	}
	connMgrs := [][]handlers.DON{{shard0DON, shard1DON}, {shard2DON}}

	router, err := handlers.NewShardRouter(shardedDONs, connMgrs)
	require.NoError(t, err)

	require.Equal(t, 3, router.NumShards())

	shard0 := router.Shard(0)
	require.Len(t, shard0.Members, 2)
	require.Equal(t, 1, shard0.F)

	shard1 := router.Shard(1)
	require.Len(t, shard1.Members, 1)
	require.Equal(t, 1, shard1.F)

	shard2 := router.Shard(2)
	require.Len(t, shard2.Members, 1)
	require.Equal(t, 2, shard2.F)

	allMembers := router.AllMembers()
	require.Len(t, allMembers, 4)

	idx, ok := router.ShardIndexForNode("0x1")
	require.True(t, ok)
	require.Equal(t, 0, idx)

	idx, ok = router.ShardIndexForNode("0x3")
	require.True(t, ok)
	require.Equal(t, 1, idx)

	idx, ok = router.ShardIndexForNode("0x4")
	require.True(t, ok)
	require.Equal(t, 2, idx)

	_, ok = router.ShardIndexForNode("0xunknown")
	require.False(t, ok)
}
