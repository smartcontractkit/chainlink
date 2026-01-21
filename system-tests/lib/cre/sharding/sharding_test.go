package sharding

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

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
