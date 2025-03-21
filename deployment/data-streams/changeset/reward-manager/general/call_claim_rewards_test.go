package general

import (
	"testing"

	"github.com/stretchr/testify/require"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

func TestCallClaimRewards(t *testing.T) {
	e := testutil.NewMemoryEnv(t, true)

	e, rewardManagerAddr, _ := DeployRewardManagerAndLink(t, e)

	var poolID [32]byte
	copy(poolID[:], []byte("poolId"))

	_, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			ClaimRewardsChangeset,
			ClaimRewardsConfig{
				ConfigsByChain: map[uint64][]ClaimRewards{
					testutil.TestChain.Selector: {ClaimRewards{
						RewardManagerAddress: rewardManagerAddr,
						PoolIDs:              [][32]byte{poolID},
					}},
				},
			},
		),
	)
	// Need Configured Fee Manager For ClaimRewards Event
	require.NoError(t, err)
}
