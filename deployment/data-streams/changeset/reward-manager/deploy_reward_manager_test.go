package reward_manager

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
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

	envDatastore, err := datastore.FromDefault[
		metadata.SerializedContractMetadata,
		datastore.DefaultMetadata,
	](e.DataStore)
	require.NoError(t, err)

	// Verify Contract Is Deployed
	record, err := envDatastore.Addresses().Get(
		datastore.NewAddressRefKey(testutil.TestChain.Selector, datastore.ContractType(types.RewardManager), &deployment.Version0_5_0, ""),
	)
	require.NoError(t, err)
	require.NotNil(t, record)

	t.Run("VerifyOwnershipTransferred", func(t *testing.T) {
		chain := e.Chains[testutil.TestChain.Selector]
		owner, _, err := commonChangesets.LoadOwnableContract(common.HexToAddress(record.Address), chain.Client)
		require.NoError(t, err)
		assert.Equal(t, testEnv.Timelocks[testutil.TestChain.Selector].Timelock.Address(), owner)
	})
}
