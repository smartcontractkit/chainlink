package verification

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestDeployVerifier(t *testing.T) {
	t.Parallel()
	testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{ShouldDeployMCMS: true})

	// Step 1 deploy VerifierProxy
	cc := DeployVerifierProxyConfig{
		ChainsToDeploy: map[uint64]DeployVerifierProxy{
			testutil.TestChain.Selector: {AccessControllerAddress: common.Address{}},
		},
		Version: deployment.Version0_5_0,
	}

	e, err := commonChangesets.Apply(t, testEnv.Environment, nil,
		commonChangesets.Configure(
			DeployVerifierProxyChangeset,
			cc,
		),
	)

	require.NoError(t, err)

	record, err := e.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(testutil.TestChain.Selector, datastore.ContractType(types.VerifierProxy), &deployment.Version0_5_0, ""),
	)
	require.NoError(t, err)
	verifierProxyAddr := common.HexToAddress(record.Address)

	// Step 2 deploy Verifier
	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			DeployVerifierChangeset,
			DeployVerifierConfig{
				ChainsToDeploy: map[uint64]DeployVerifier{
					testutil.TestChain.Selector: {VerifierProxyAddress: verifierProxyAddr},
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
	record, err = envDatastore.Addresses().Get(
		datastore.NewAddressRefKey(testutil.TestChain.Selector, datastore.ContractType(types.Verifier), &deployment.Version0_5_0, ""),
	)
	require.NoError(t, err)
	require.NotNil(t, record)

	t.Run("VerifyOwnershipTransferred", func(t *testing.T) {
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

		contractMetadata := testutil.MustGetContractMetaData[v0_5.VerifierView](t, output.DataStore, testutil.TestChain.Selector, record.Address)

		assert.Equal(t, testEnv.Timelocks[testutil.TestChain.Selector].Timelock.Address().Hex(), contractMetadata.View.Owner)
	})
}
