package general

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestDeployVerifier(t *testing.T) {
	e := testutil.NewMemoryEnv(t, true)

	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			deployment.CreateLegacyChangeSet(commonChangesets.DeployLinkToken),
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
					testutil.TestChain.Selector: {LinkAddress: linkState.LinkToken.Address()},
				},
			},
		),
	)

	require.NoError(t, err)

	_, err = deployment.SearchAddressBook(e.ExistingAddresses, testutil.TestChain.Selector, types.RewardManager)
	require.NoError(t, err)
}
