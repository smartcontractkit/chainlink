package v0_5_0

import (
	"testing"

	"github.com/stretchr/testify/require"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

func TestCallClaimRewards(t *testing.T) {
	e := testutil.NewMemoryEnv(t, true)

	e, rewardManagerAddr, _ := DeployRewardManagerAndLink(t, e)

	var poolId [32]byte
	copy(poolId[:], []byte("poolId"))

	_, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			ClaimRewardsChangeset,
			ClaimRewardsConfig{
				ConfigsByChain: map[uint64][]ClaimRewards{
					testutil.TestChain.Selector: {ClaimRewards{
						RewardManagerAddress: rewardManagerAddr,
						PoolIds:              [][32]byte{poolId},
					}},
				},
			},
		),
	)
	// Need Configured Fee Manager For ClaimRewards Event
	require.NoError(t, err)
}
