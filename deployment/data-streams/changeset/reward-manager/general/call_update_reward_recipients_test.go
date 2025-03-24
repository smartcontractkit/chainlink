package general

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	rewardManager "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"
)

func TestCallUpdateRewardRecipients(t *testing.T) {
	e := testutil.NewMemoryEnv(t, true)
	chain := e.Chains[testutil.TestChain.Selector]

	e, rewardManagerAddr, _ := DeployRewardManagerAndLink(t, e)

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
			},
		),
	)
	require.NoError(t, err)

	rm, err := rewardManager.NewRewardManager(rewardManagerAddr, chain.Client)
	require.NoError(t, err)
	recipientsUpdatedIterator, err := rm.FilterRewardRecipientsUpdated(nil, [][32]byte{poolID})
	require.NoError(t, err)
	defer recipientsUpdatedIterator.Close()
	require.NoError(t, err)
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
			Addr:   common.HexToAddress("0x3333333333333333333333333333333333333333"),
			Weight: 300000000000000000,
		},
		{
			Addr:   common.HexToAddress("0x4444444444444444444444444444444444444444"),
			Weight: 700000000000000000,
		},
	}
	var poolIDUpdated [32]byte
	copy(poolID[:], []byte("poolIdUpdated"))

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			SetRewardRecipientsChangeset,
			SetRewardRecipientsConfig{
				ConfigsByChain: map[uint64][]SetRewardRecipients{
					testutil.TestChain.Selector: {SetRewardRecipients{
						RewardManagerAddress:      rewardManagerAddr,
						PoolID:                    poolIDUpdated,
						RewardRecipientAndWeights: recipientsUpdated,
					}},
				},
			},
		),
	)
	require.NoError(t, err)

	recipientsUpdatedIterator, err = rm.FilterRewardRecipientsUpdated(nil, [][32]byte{poolIDUpdated})
	require.NoError(t, err)
	defer recipientsUpdatedIterator.Close()
	require.NoError(t, err)
	foundExpected = false

	for recipientsUpdatedIterator.Next() {
		event := recipientsUpdatedIterator.Event
		if poolIDUpdated == event.PoolId && reflect.DeepEqual(recipientsUpdated, event.NewRewardRecipients) {
			foundExpected = true
			break
		}
	}
	require.True(t, foundExpected)
}
