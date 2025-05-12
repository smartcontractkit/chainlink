package reward_manager

import (
	"testing"

	"github.com/stretchr/testify/require"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

func runClaimRewardsTest(t *testing.T, useMCMS bool) {
	testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{
		ShouldDeployMCMS:      useMCMS,
		ShouldDeployLinkToken: true,
	})
	chainSelector := testutil.TestChain.Selector

	e, rewardManagerAddr := RewardManagerDeploy(t, testEnv)

	var poolID [32]byte
	copy(poolID[:], []byte("poolId"))

	_, _, err := commonChangesets.ApplyChangesetsV2(
		t, e, []commonChangesets.ConfiguredChangeSet{
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
			)},
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
