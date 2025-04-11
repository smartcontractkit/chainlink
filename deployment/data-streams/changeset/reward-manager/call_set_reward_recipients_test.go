package reward_manager

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	rewardManager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

func runSetRewardRecipientsTest(t *testing.T, useMCMS bool) {
	e := testutil.NewMemoryEnv(t, true)
	chainSelector := testutil.TestChain.Selector
	chain := e.Chains[chainSelector]
	e, rewardManagerAddr, _ := DeployRewardManagerAndLink(t, e)

	var poolID [32]byte
	copy(poolID[:], []byte("poolId"))

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

	var timelocks map[uint64]*proposalutils.TimelockExecutionContracts
	if useMCMS {
		e, _, timelocks = testutil.DeployMCMS(t, e, map[uint64][]common.Address{
			chainSelector: {rewardManagerAddr},
		})
	}

	_, err := commonChangesets.Apply(
		t, e, timelocks,
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
		),
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
