package v2

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	handlermocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
)

// handlersDON is a local alias to keep the test matrix literals readable.
type handlersDON = handlers.DON

func TestBuildShardEndpoints(t *testing.T) {
	t.Parallel()

	node := func(addr string) config.NodeConfig {
		return config.NodeConfig{Name: addr, Address: addr}
	}

	t.Run("single DON single shard (backward compatible)", func(t *testing.T) {
		t.Parallel()

		don := handlermocks.NewDON(t)
		shardedDONs := []config.ShardedDONConfig{
			{DonName: "donA", F: 1, Shards: []config.Shard{{Nodes: []config.NodeConfig{node("n1"), node("n2")}}}},
		}
		connMgrs := [][]handlersDON{{don}}

		eps, addrMap, err := buildShardEndpoints(shardedDONs, connMgrs)
		require.NoError(t, err)
		require.Len(t, eps, 1)
		require.Equal(t, "donA", eps[0].donID) // shard 0 => bare name
		require.Equal(t, 0, eps[0].shardIdx)
		require.Equal(t, 1, eps[0].f)
		require.Len(t, addrMap, 2)
		require.Same(t, eps[0], addrMap["n1"])
	})

	t.Run("full matrix: two DONs each with two shards", func(t *testing.T) {
		t.Parallel()

		dons := [][]handlersDON{
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

		eps, addrMap, err := buildShardEndpoints(shardedDONs, dons)
		require.NoError(t, err)
		require.Len(t, eps, 4)
		// donIDs: shard 0 bare, shard 1 suffixed
		require.Equal(t, "donA", eps[0].donID)
		require.Equal(t, "donA_1", eps[1].donID)
		require.Equal(t, "donB", eps[2].donID)
		require.Equal(t, "donB_1", eps[3].donID)
		require.Len(t, addrMap, 6)
		require.Same(t, eps[1], addrMap["a2"])
		require.Same(t, eps[3], addrMap["b1"])
		require.Equal(t, 2, addrMap["b1"].f)
	})

	t.Run("rejects empty matrix", func(t *testing.T) {
		t.Parallel()

		_, _, err := buildShardEndpoints(nil, nil)
		require.Error(t, err)
	})

	t.Run("rejects ragged outer dimension", func(t *testing.T) {
		t.Parallel()

		shardedDONs := []config.ShardedDONConfig{
			{DonName: "donA", Shards: []config.Shard{{Nodes: []config.NodeConfig{node("n1")}}}},
		}
		_, _, err := buildShardEndpoints(shardedDONs, [][]handlersDON{})
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
		connMgrs := [][]handlersDON{{handlermocks.NewDON(t)}} // only 1 conn mgr for 2 shards
		_, _, err := buildShardEndpoints(shardedDONs, connMgrs)
		require.Error(t, err)
	})

	t.Run("rejects nil connection manager", func(t *testing.T) {
		t.Parallel()

		shardedDONs := []config.ShardedDONConfig{
			{DonName: "donA", Shards: []config.Shard{{Nodes: []config.NodeConfig{node("n1")}}}},
		}
		connMgrs := [][]handlersDON{{nil}}
		_, _, err := buildShardEndpoints(shardedDONs, connMgrs)
		require.Error(t, err)
	})

	t.Run("rejects duplicate node address across shards", func(t *testing.T) {
		t.Parallel()

		shardedDONs := []config.ShardedDONConfig{
			{DonName: "donA", Shards: []config.Shard{
				{Nodes: []config.NodeConfig{node("n1")}},
				{Nodes: []config.NodeConfig{node("n1")}}, // duplicate
			}},
		}
		connMgrs := [][]handlersDON{{handlermocks.NewDON(t), handlermocks.NewDON(t)}}
		_, _, err := buildShardEndpoints(shardedDONs, connMgrs)
		require.Error(t, err)
		require.Contains(t, err.Error(), "disjoint")
	})

	t.Run("allMembers unions all shards", func(t *testing.T) {
		t.Parallel()

		eps := []*shardEndpoint{
			{members: []config.NodeConfig{node("n1"), node("n2")}},
			{members: []config.NodeConfig{node("n3")}},
		}
		members := allMembers(eps)
		require.Len(t, members, 3)
	})
}
