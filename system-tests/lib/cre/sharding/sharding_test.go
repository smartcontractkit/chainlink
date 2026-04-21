package sharding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func makeNodes(count int) []*cre.Node {
	nodes := make([]*cre.Node, count)
	for i := range nodes {
		nodes[i] = &cre.Node{Name: "node-" + string(rune('0'+i))}
	}
	return nodes
}

func TestGetShardZeroDON(t *testing.T) {
	t.Run("returns shard zero when present", func(t *testing.T) {
		shardZero := &cre.Don{
			Name:  "shard-zero",
			ID:    1,
			Flags: []cre.CapabilityFlag{cre.ShardDON},
		}
		shardOne := &cre.Don{
			Name:  "shard-one",
			ID:    2,
			Flags: []cre.CapabilityFlag{cre.ShardDON},
		}

		dons := cre.NewDons([]*cre.Don{shardZero, shardOne}, nil)

		result, err := getShardLeaderDON(dons)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, "shard-zero", result.Name)
	})

	t.Run("returns error when no shard DONs exist", func(t *testing.T) {
		nonShardDON := &cre.Don{
			Name:  "workflow-don",
			ID:    1,
			Flags: []cre.CapabilityFlag{cre.WorkflowDON},
		}
		dons := cre.NewDons([]*cre.Don{nonShardDON}, nil)

		result, err := getShardLeaderDON(dons)
		require.Error(t, err)
		require.Nil(t, result)
		require.Contains(t, err.Error(), "no shard leader DON found")
	})
}

func TestRingContractQualifier(t *testing.T) {
	require.Equal(t, "ring", RingContractQualifier)
}

func TestShardOrchestratorAddress(t *testing.T) {
	nm := &cre.NodeMetadata{Host: "10.0.0.1"}

	t.Run("default port matches constant", func(t *testing.T) {
		addr := nm.ShardOrchestratorAddress()
		assert.Equal(t, "10.0.0.1:50051", addr)
		assert.Equal(t, uint16(50051), cre.DefaultShardOrchestratorPort)
	})

	t.Run("custom port", func(t *testing.T) {
		addr := nm.ShardOrchestratorAddressWithPort(60051)
		assert.Equal(t, "10.0.0.1:60051", addr)
	})

	t.Run("default arbiter port", func(t *testing.T) {
		assert.Equal(t, uint16(9876), cre.DefaultArbiterPort)
	})
}

func TestValidateShardTopology(t *testing.T) {
	t.Run("valid topology passes", func(t *testing.T) {
		dons := cre.NewDons([]*cre.Don{
			{Name: "shard0", Flags: []cre.CapabilityFlag{cre.ShardDON}, ShardIndex: 0, Nodes: makeNodes(4)},
			{Name: "shard1", Flags: []cre.CapabilityFlag{cre.ShardDON}, ShardIndex: 1, Nodes: makeNodes(4)},
		}, nil)
		require.NoError(t, ValidateShardTopology(dons))
	})

	t.Run("single shard DON fails", func(t *testing.T) {
		dons := cre.NewDons([]*cre.Don{
			{Name: "shard0", Flags: []cre.CapabilityFlag{cre.ShardDON}, ShardIndex: 0, Nodes: makeNodes(4)},
		}, nil)
		err := ValidateShardTopology(dons)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "at least 2 shard DONs")
	})

	t.Run("no leader fails", func(t *testing.T) {
		dons := cre.NewDons([]*cre.Don{
			{Name: "shard1", Flags: []cre.CapabilityFlag{cre.ShardDON}, ShardIndex: 1, Nodes: makeNodes(4)},
			{Name: "shard2", Flags: []cre.CapabilityFlag{cre.ShardDON}, ShardIndex: 2, Nodes: makeNodes(4)},
		}, nil)
		err := ValidateShardTopology(dons)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "shard_index=0")
	})

	t.Run("duplicate leaders fail", func(t *testing.T) {
		dons := cre.NewDons([]*cre.Don{
			{Name: "shard0a", Flags: []cre.CapabilityFlag{cre.ShardDON}, ShardIndex: 0, Nodes: makeNodes(4)},
			{Name: "shard0b", Flags: []cre.CapabilityFlag{cre.ShardDON}, ShardIndex: 0, Nodes: makeNodes(4)},
		}, nil)
		err := ValidateShardTopology(dons)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple shard DONs")
	})

	t.Run("mismatched node count fails", func(t *testing.T) {
		dons := cre.NewDons([]*cre.Don{
			{Name: "shard0", Flags: []cre.CapabilityFlag{cre.ShardDON}, ShardIndex: 0, Nodes: makeNodes(4)},
			{Name: "shard1", Flags: []cre.CapabilityFlag{cre.ShardDON}, ShardIndex: 1, Nodes: makeNodes(3)},
		}, nil)
		err := ValidateShardTopology(dons)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "same node count")
	})

	t.Run("non-shard DONs are ignored", func(t *testing.T) {
		dons := cre.NewDons([]*cre.Don{
			{Name: "workflow", Flags: []cre.CapabilityFlag{cre.WorkflowDON}, Nodes: makeNodes(2)},
			{Name: "shard0", Flags: []cre.CapabilityFlag{cre.ShardDON}, ShardIndex: 0, Nodes: makeNodes(4)},
			{Name: "shard1", Flags: []cre.CapabilityFlag{cre.ShardDON}, ShardIndex: 1, Nodes: makeNodes(4)},
		}, nil)
		require.NoError(t, ValidateShardTopology(dons))
	})
}
