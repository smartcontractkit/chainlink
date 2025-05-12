package verification

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink/deployment"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestDeployVerifierProxy(t *testing.T) {
	t.Parallel()
	testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{ShouldDeployMCMS: true})

	cc := DeployVerifierProxyConfig{
		ChainsToDeploy: map[uint64]DeployVerifierProxy{
			testutil.TestChain.Selector: {AccessControllerAddress: common.HexToAddress("0x001")},
		},
		Ownership: types.OwnershipSettings{
			ShouldTransfer: true,
			MCMSProposalConfig: &proposalutils.TimelockConfig{
				MinDelay: 0,
			},
		},
		Version: deployment.Version0_5_0,
	}

	e, _, err := commonChangesets.ApplyChangesetsV2(t, testEnv.Environment, []commonChangesets.ConfiguredChangeSet{
		commonChangesets.Configure(
			DeployVerifierProxyChangeset,
			cc,
		),
	})
	require.NoError(t, err)

	envDatastore, err := datastore.FromDefault[metadata.SerializedContractMetadata, datastore.DefaultMetadata](e.DataStore)
	require.NoError(t, err)

	// Verify Contract Is Deployed
	record, err := envDatastore.Addresses().Get(
		datastore.NewAddressRefKey(testutil.TestChain.Selector, datastore.ContractType(types.VerifierProxy), &deployment.Version0_5_0, ""),
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
