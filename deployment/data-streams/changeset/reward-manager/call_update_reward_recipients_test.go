package reward_manager

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	rewardManager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

func runCallUpdateRewardRecipients(t *testing.T, useMCMS bool) {
	testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{
		ShouldDeployMCMS:      useMCMS,
		ShouldDeployLinkToken: true,
	})
	e, rewardManagerAddr := RewardManagerDeploy(t, testEnv)
	chainSelector := testutil.TestChain.Selector
	chain := e.Chains[chainSelector]

	recipients := []rewardManager.CommonAddressAndWeight{
		{
			Addr:   common.HexToAddress("0x1111111111111111111111111111111111111111"),
			Weight: 500000000000000000,
		},
		{
			Addr:   common.HexToAddress("0x2222222222222222222222222222222222222222"),
			Weight: 500000000000000000,
		},
	}
	var poolID [32]byte
	copy(poolID[:], []byte("poolId"))

	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			SetRewardRecipientsChangeset,
			SetRewardRecipientsConfig{
				ConfigsByChain: map[uint64][]SetRewardRecipients{
					testutil.TestChain.Selector: {SetRewardRecipients{
						RewardManagerAddress:      rewardManagerAddr,
						PoolID:                    poolID,
						RewardRecipientAndWeights: recipients,
					}},
				},
				MCMSConfig: testutil.GetMCMSConfig(useMCMS),
			},
		),
	)
	require.NoError(t, err)

	rm, err := rewardManager.NewRewardManager(rewardManagerAddr, chain.Client)
	require.NoError(t, err)
	recipientsUpdatedIterator, err := rm.FilterRewardRecipientsUpdated(nil, [][32]byte{poolID})
	require.NoError(t, err)
	defer recipientsUpdatedIterator.Close()
	foundExpected := false

	for recipientsUpdatedIterator.Next() {
		event := recipientsUpdatedIterator.Event
		if poolID == event.PoolId && reflect.DeepEqual(recipients, event.NewRewardRecipients) {
			foundExpected = true
			break
		}
	}
	require.True(t, foundExpected)

	recipientsUpdated := []rewardManager.CommonAddressAndWeight{
		{
			Addr:   common.HexToAddress("0x1111111111111111111111111111111111111111"),
			Weight: 600000000000000000,
		},
		{
			Addr:   common.HexToAddress("0x2222222222222222222222222222222222222222"),
			Weight: 400000000000000000,
		},
	}

	e, _, err = commonChangesets.ApplyChangesetsV2(
		t, e, []commonChangesets.ConfiguredChangeSet{
			commonChangesets.Configure(
				UpdateRewardRecipientsChangeset,
				UpdateRewardRecipientsConfig{
					ConfigsByChain: map[uint64][]UpdateRewardRecipients{
						testutil.TestChain.Selector: {UpdateRewardRecipients{
							RewardManagerAddress:      rewardManagerAddr,
							PoolID:                    poolID,
							RewardRecipientAndWeights: recipientsUpdated,
						}},
					},
					MCMSConfig: testutil.GetMCMSConfig(useMCMS),
				},
			)},
	)
	require.NoError(t, err)

	recipientsUpdatedIterator, err = rm.FilterRewardRecipientsUpdated(nil, [][32]byte{poolID})
	require.NoError(t, err)
	defer recipientsUpdatedIterator.Close()
	require.NoError(t, err)
	foundExpected = false

	for recipientsUpdatedIterator.Next() {
		event := recipientsUpdatedIterator.Event
		if poolID == event.PoolId && reflect.DeepEqual(recipientsUpdated, event.NewRewardRecipients) {
			foundExpected = true
			break
		}
	}
	require.True(t, foundExpected)
}

func TestCallUpdateRewardRecipients(t *testing.T) {
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
			runCallUpdateRewardRecipients(t, tc.useMCMS)
		})
	}
}
