package reward_manager

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

func runPayRecipientsTest(t *testing.T, useMCMS bool) {
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
				PayRecipientsChangeset,
				PayRecipientsConfig{
					ConfigsByChain: map[uint64][]PayRecipients{
						chainSelector: {{
							RewardManagerAddress: rewardManagerAddr,
							PoolID:               poolID,
							Recipients:           []common.Address{},
						}},
					},
					MCMSConfig: testutil.GetMCMSConfig(useMCMS),
				},
			)},
	)
	// Need Configured Fee Manager For PayRecipients Event
	require.NoError(t, err)
}

func TestPayRecipients(t *testing.T) {
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
			runPayRecipientsTest(t, tc.useMCMS)
		})
	}
}
