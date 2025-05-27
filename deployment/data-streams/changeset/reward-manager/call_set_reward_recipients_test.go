package reward_manager

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5"

	rewardManager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

func runSetRewardRecipientsTest(t *testing.T, useMCMS bool) {
	testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{
		ShouldDeployMCMS:      useMCMS,
		ShouldDeployLinkToken: true,
	})
	chainSelector := testutil.TestChain.Selector
	e, rewardManagerAddr := RewardManagerDeploy(t, testEnv)
	chain := e.Chains[chainSelector]

	var poolID [32]byte
	copy(poolID[:], []byte("poolId"))

	r1 := "0x1111111111111111111111111111111111111111"
	r2 := "0x2222222222222222222222222222222222222222"

	recipients := []rewardManager.CommonAddressAndWeight{
		{
			Addr:   common.HexToAddress(r1),
			Weight: 400000000000000000,
		},
		{
			Addr:   common.HexToAddress(r2),
			Weight: 600000000000000000,
		},
	}

	_, _, err := commonChangesets.ApplyChangesetsV2(
		t, e, []commonChangesets.ConfiguredChangeSet{
			commonChangesets.Configure(
				SetRewardRecipientsChangeset,
				SetRewardRecipientsConfig{
					ConfigsByChain: map[uint64][]SetRewardRecipients{
						chainSelector: {{
							RewardManagerAddress:      rewardManagerAddr,
							PoolID:                    poolID,
							RewardRecipientAndWeights: recipients,
						}},
					},
					MCMSConfig: testutil.GetMCMSConfig(useMCMS),
				},
			)},
	)
	require.NoError(t, err)

	rm, err := rewardManager.NewRewardManager(rewardManagerAddr, chain.Client)
	require.NoError(t, err)
	it, err := rm.FilterRewardRecipientsUpdated(nil, [][32]byte{poolID})
	require.NoError(t, err)
	defer it.Close()

	foundExpected := false
	for it.Next() {
		event := it.Event
		if poolID == event.PoolId && reflect.DeepEqual(recipients, event.NewRewardRecipients) {
			foundExpected = true
			break
		}
	}
	require.True(t, foundExpected)

	t.Run("VerifyMetadata", func(t *testing.T) {
		// Use View To Confirm Data
		_, outputs, err := commonChangesets.ApplyChangesetsV2(t, e,
			[]commonChangesets.ConfiguredChangeSet{
				commonChangesets.Configure(
					changeset.SaveContractViews,
					changeset.SaveContractViewsConfig{
						Chains: []uint64{testutil.TestChain.Selector},
					},
				),
			},
		)
		require.NoError(t, err)
		require.Len(t, outputs, 1)
		output := outputs[0]

		contractMetadata := testutil.MustGetContractMetaData[v0_5.RewardManagerView](t, output.DataStore, testutil.TestChain.Selector, rewardManagerAddr.Hex())

		require.NotNil(t, contractMetadata)
		poolIDHex := dsutil.HexEncodeBytes32(poolID)
		recipientWeights := contractMetadata.View.RecipientWeights[poolIDHex]
		require.Equal(t, len(recipients), len(recipientWeights))
		for _, recipient := range recipients {
			// Compare configured (expected) recipients with the ones retrieved from the view
			switch recipient.Addr.Hex() {
			case r1:
				require.Equal(t, strconv.FormatUint(recipient.Weight, 10), recipientWeights[r1].Weight)
			case r2:
				require.Equal(t, strconv.FormatUint(recipient.Weight, 10), recipientWeights[r2].Weight)
			default:
				t.Fatalf("Unexpected recipient address: %s", recipient.Addr.Hex())
			}
		}
	})
}

func TestSetRewardRecipients(t *testing.T) {
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
			runSetRewardRecipientsTest(t, tc.useMCMS)
		})
	}
}
