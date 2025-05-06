package reward_manager

import (
	"testing"

	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestDeployRewardManager(t *testing.T) {
	testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{ShouldDeployMCMS: true})

	e, err := commonChangesets.Apply(t, testEnv.Environment, nil,
		commonChangesets.Configure(
			cldf.CreateLegacyChangeSet(commonChangesets.DeployLinkToken),
			[]uint64{testutil.TestChain.Selector},
		),
	)

	require.NoError(t, err)

	addresses, err := e.ExistingAddresses.AddressesForChain(testutil.TestChain.Selector)
	require.NoError(t, err)

	chain := e.Chains[testutil.TestChain.Selector]
	linkState, err := commonState.MaybeLoadLinkTokenChainState(chain, addresses)
	require.NoError(t, err)

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			DeployRewardManagerChangeset,
			DeployRewardManagerConfig{
				ChainsToDeploy: map[uint64]DeployRewardManager{
					testutil.TestChain.Selector: {LinkTokenAddress: linkState.LinkToken.Address()},
				},
				Ownership: types.OwnershipSettings{
					ShouldTransfer: true,
					MCMSProposalConfig: &proposalutils.TimelockConfig{
						MinDelay: 0,
					},
				},
			},
		),
	)

	require.NoError(t, err)

	rmAddr, err := dsutil.MaybeFindEthAddress(e.ExistingAddresses, testutil.TestChain.Selector, types.RewardManager)
	require.NoError(t, err)

	owner, _, err := commonChangesets.LoadOwnableContract(rmAddr, chain.Client)
	require.NoError(t, err)
	require.Equal(t, testEnv.Timelocks[testutil.TestChain.Selector].Timelock.Address(), owner)
}
