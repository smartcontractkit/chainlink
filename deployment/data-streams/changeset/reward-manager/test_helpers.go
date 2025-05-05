package reward_manager

import (
	"testing"

	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	dsTypes "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

// DeployRewardManagerAndLink deploys Link Token and deploys Verifier
// and returns the updated environment and the addresses of RewardManager.
func DeployRewardManagerAndLink(
	t *testing.T,
	e deployment.Environment,
) (env deployment.Environment, rewardManagerAddr common.Address, linkState *commonState.LinkTokenState) {
	t.Helper()

	chainSelector := testutil.TestChain.Selector

	// 1) Deploy Link
	env, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			cldf.CreateLegacyChangeSet(commonChangesets.DeployLinkToken),
			[]uint64{testutil.TestChain.Selector},
		),
	)

	require.NoError(t, err)

	addresses, err := env.ExistingAddresses.AddressesForChain(testutil.TestChain.Selector)
	require.NoError(t, err)

	chain := env.Chains[testutil.TestChain.Selector]
	linkState, err = commonState.MaybeLoadLinkTokenChainState(chain, addresses)
	require.NoError(t, err)

	// 2) Deploy RewardManager
	deployRewardManagerCfg := DeployRewardManagerConfig{
		ChainsToDeploy: map[uint64]DeployRewardManager{
			testutil.TestChain.Selector: {LinkTokenAddress: linkState.LinkToken.Address()},
		},
	}

	env, err = commonChangesets.Apply(t, env, nil,
		commonChangesets.Configure(
			DeployRewardManagerChangeset,
			deployRewardManagerCfg,
		),
	)
	require.NoError(t, err, "deploying RewardManager should not fail")

	// Get the RewardManager address
	rewardManagerAddrHex, err := deployment.SearchAddressBook(env.ExistingAddresses, chainSelector, dsTypes.RewardManager)
	require.NoError(t, err, "unable to find RewardManager address in address book")
	rewardManagerAddr = common.HexToAddress(rewardManagerAddrHex)
	require.NotEqual(t, common.Address{}, rewardManagerAddr, "RewardManager should not be zero address")

	return env, rewardManagerAddr, linkState
}
