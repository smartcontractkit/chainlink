package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	handlermocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
)

func node(addr string) config.NodeConfig {
	return config.NodeConfig{Name: addr, Address: addr}
}

func TestNewShardedDONs(t *testing.T) {
	t.Parallel()

	t.Run("rejects empty matrix", func(t *testing.T) {
		t.Parallel()

		_, err := handlers.NewShardedDONs(nil, nil)
		require.Error(t, err)
	})

	t.Run("rejects ragged outer dimension", func(t *testing.T) {
		t.Parallel()

		shardedDONs := []config.ShardedDONConfig{
			{DonName: "donA", Shards: []config.Shard{{Nodes: []config.NodeConfig{node("n1")}}}},
		}
		_, err := handlers.NewShardedDONs(shardedDONs, [][]handlers.DON{})
		require.Error(t, err)
	})

	t.Run("rejects ragged inner dimension", func(t *testing.T) {
		t.Parallel()

		shardedDONs := []config.ShardedDONConfig{
			{DonName: "donA", Shards: []config.Shard{
				{Nodes: []config.NodeConfig{node("n1")}},
				{Nodes: []config.NodeConfig{node("n2")}},
			}},
		}
		connMgrs := [][]handlers.DON{{handlermocks.NewDON(t)}} // only 1 conn mgr for 2 shards
		_, err := handlers.NewShardedDONs(shardedDONs, connMgrs)
		require.Error(t, err)
	})

	t.Run("rejects nil connection manager", func(t *testing.T) {
		t.Parallel()

		shardedDONs := []config.ShardedDONConfig{
			{DonName: "donA", Shards: []config.Shard{{Nodes: []config.NodeConfig{node("n1")}}}},
		}
		connMgrs := [][]handlers.DON{{nil}}
		_, err := handlers.NewShardedDONs(shardedDONs, connMgrs)
		require.Error(t, err)
	})
}

func TestShardedDONs_BuildShardEndpoints(t *testing.T) {
	t.Parallel()

	t.Run("single DON single shard (backward compatible)", func(t *testing.T) {
		t.Parallel()

		don := handlermocks.NewDON(t)
		shardedDONs := []config.ShardedDONConfig{
			{DonName: "donA", F: 1, Shards: []config.Shard{{Nodes: []config.NodeConfig{node("n1"), node("n2")}}}},
		}
		dons, err := handlers.NewShardedDONs(shardedDONs, [][]handlers.DON{{don}})
		require.NoError(t, err)

		eps, addrMap, err := dons.BuildShardEndpoints()
		require.NoError(t, err)
		require.Len(t, eps, 1)
		require.Equal(t, "donA", eps[0].DonID) // shard 0 => bare name
		require.Equal(t, 0, eps[0].ShardIdx)
		require.Equal(t, 1, eps[0].F)
		require.Len(t, addrMap, 2)
		require.Same(t, eps[0], addrMap["n1"])
	})

	t.Run("full matrix: two DONs each with two shards", func(t *testing.T) {
		t.Parallel()

		connMgrs := [][]handlers.DON{
			{handlermocks.NewDON(t), handlermocks.NewDON(t)},
			{handlermocks.NewDON(t), handlermocks.NewDON(t)},
		}
		shardedDONs := []config.ShardedDONConfig{
			{DonName: "donA", F: 1, Shards: []config.Shard{
				{Nodes: []config.NodeConfig{node("a0"), node("a1")}},
				{Nodes: []config.NodeConfig{node("a2"), node("a3")}},
			}},
			{DonName: "donB", F: 2, Shards: []config.Shard{
				{Nodes: []config.NodeConfig{node("b0")}},
				{Nodes: []config.NodeConfig{node("b1")}},
			}},
		}
		dons, err := handlers.NewShardedDONs(shardedDONs, connMgrs)
		require.NoError(t, err)

		eps, addrMap, err := dons.BuildShardEndpoints()
		require.NoError(t, err)
		require.Len(t, eps, 4)
		// donIDs: shard 0 bare, shard 1 suffixed
		require.Equal(t, "donA", eps[0].DonID)
		require.Equal(t, "donA_shard-1", eps[1].DonID)
		require.Equal(t, "donB", eps[2].DonID)
		require.Equal(t, "donB_shard-1", eps[3].DonID)
		require.Len(t, addrMap, 6)
		require.Same(t, eps[1], addrMap["a2"])
		require.Same(t, eps[3], addrMap["b1"])
		require.Equal(t, 2, addrMap["b1"].F)
	})

	t.Run("rejects duplicate node address across shards", func(t *testing.T) {
		t.Parallel()

		shardedDONs := []config.ShardedDONConfig{
			{DonName: "donA", Shards: []config.Shard{
				{Nodes: []config.NodeConfig{node("n1")}},
				{Nodes: []config.NodeConfig{node("n1")}}, // duplicate
			}},
		}
		connMgrs := [][]handlers.DON{{handlermocks.NewDON(t), handlermocks.NewDON(t)}}
		dons, err := handlers.NewShardedDONs(shardedDONs, connMgrs)
		require.NoError(t, err)

		_, _, err = dons.BuildShardEndpoints()
		require.Error(t, err)
		require.Contains(t, err.Error(), "disjoint")
	})

	t.Run("AllMembers unions all shards", func(t *testing.T) {
		t.Parallel()

		don0, don1 := handlermocks.NewDON(t), handlermocks.NewDON(t)
		shardedDONs := []config.ShardedDONConfig{
			{DonName: "donA", Shards: []config.Shard{
				{Nodes: []config.NodeConfig{node("n1"), node("n2")}},
				{Nodes: []config.NodeConfig{node("n3")}},
			}},
		}
		dons, err := handlers.NewShardedDONs(shardedDONs, [][]handlers.DON{{don0, don1}})
		require.NoError(t, err)

		eps, _, err := dons.BuildShardEndpoints()
		require.NoError(t, err)
		require.Len(t, handlers.AllMembers(eps), 3)
	})
}
