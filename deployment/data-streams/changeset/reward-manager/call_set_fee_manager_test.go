package reward_manager

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	rewardManager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"
	"github.com/smartcontractkit/chainlink/deployment"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	feeManagerCs "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/fee-manager"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func runSetFeeManagerTest(t *testing.T, useMCMS bool) {
	testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{
		ShouldDeployMCMS:      useMCMS,
		ShouldDeployLinkToken: true,
	})
	e := testEnv.Environment

	chainSelector := testutil.TestChain.Selector
	chain := e.Chains[chainSelector]

	e, rewardManagerAddr := RewardManagerDeploy(t, testEnv)

	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			feeManagerCs.DeployFeeManagerChangeset,
			feeManagerCs.DeployFeeManagerConfig{
				ChainsToDeploy: map[uint64]feeManagerCs.DeployFeeManager{
					testutil.TestChain.Selector: {
						LinkTokenAddress:     testEnv.LinkTokenState.LinkToken.Address(),
						NativeTokenAddress:   common.HexToAddress("0x3e5e9111ae8eb78fe1cc3bb8915d5d461f3ef9a9"),
						VerifierProxyAddress: common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e"),
						RewardManagerAddress: rewardManagerAddr,
					},
				},
			},
		),
	)

	require.NoError(t, err)

	// Get the FeeManager address
	record, err := e.DataStore.Addresses().Get(ds.NewAddressRefKey(chainSelector, ds.ContractType(types.FeeManager), &deployment.Version0_5_0, ""))
	require.NoError(t, err)
	feeManagerAddr := common.HexToAddress(record.Address)

	e, _, err = commonChangesets.ApplyChangesetsV2(
		t, e, []commonChangesets.ConfiguredChangeSet{
			commonChangesets.Configure(
				SetFeeManagerChangeset,
				SetFeeManagerConfig{
					ConfigsByChain: map[uint64][]SetFeeManager{
						chainSelector: {{
							RewardManagerAddress: rewardManagerAddr,
							FeeManagerAddress:    feeManagerAddr,
						}},
					},
					MCMSConfig: testutil.GetMCMSConfig(useMCMS),
				},
			)},
	)
	require.NoError(t, err)

	rm, err := rewardManager.NewRewardManager(rewardManagerAddr, chain.Client)
	require.NoError(t, err)

	logIterator, err := rm.FilterFeeManagerUpdated(nil)
	require.NoError(t, err)
	defer logIterator.Close()
	require.NoError(t, err)
	foundExpected := false

	for logIterator.Next() {
		event := logIterator.Event
		if feeManagerAddr == event.NewFeeManagerAddress {
			foundExpected = true
			break
		}
	}
	require.True(t, foundExpected)
}

func TestSetFeeManager(t *testing.T) {
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
			runSetFeeManagerTest(t, tc.useMCMS)
		})
	}
}
