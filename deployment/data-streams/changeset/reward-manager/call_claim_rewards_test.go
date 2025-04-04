package reward_manager

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

func runClaimRewardsTest(t *testing.T, useMCMS bool) {
	e := testutil.NewMemoryEnv(t, true)
	chainSelector := testutil.TestChain.Selector

	e, rewardManagerAddr, _ := DeployRewardManagerAndLink(t, e)

	var poolID [32]byte
	copy(poolID[:], []byte("poolId"))

	var timelocks map[uint64]*proposalutils.TimelockExecutionContracts
	if useMCMS {
		e, _, timelocks = testutil.DeployMCMS(t, e, map[uint64][]common.Address{
			chainSelector: {rewardManagerAddr},
		})
	}

	_, err := commonChangesets.Apply(
		t, e, timelocks,
		commonChangesets.Configure(
			ClaimRewardsChangeset,
			ClaimRewardsConfig{
				ConfigsByChain: map[uint64][]ClaimRewards{
					chainSelector: {{
						RewardManagerAddress: rewardManagerAddr,
						PoolIDs:              [][32]byte{poolID},
					}},
				},
				MCMSConfig: testutil.GetMCMSConfig(useMCMS),
			},
		),
	)
	require.NoError(t, err)
}

func TestClaimRewards(t *testing.T) {
	testCases := []struct {
		name    string
		useMCMS bool
	}{
		{
			name:    "Without MCMS",
			useMCMS: false,
		},
		{
			name:    "With MCMS",
			useMCMS: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runClaimRewardsTest(t, tc.useMCMS)
		})
	}
}
