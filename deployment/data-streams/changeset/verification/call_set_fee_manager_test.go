package verification

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	feemanager "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/fee-manager"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestSetFeeManager(t *testing.T) {
	t.Parallel()
	testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{
		ShouldDeployLinkToken: true,
	})

	e := testEnv.Environment

	testChain := testutil.TestChain.Selector
	e, verifierProxyAddr, _ := DeployVerifierProxyAndVerifier(t, e)

	// Deploy Fee Manager
	cfgFeeManager := feemanager.DeployFeeManagerConfig{
		ChainsToDeploy: map[uint64]feemanager.DeployFeeManager{testutil.TestChain.Selector: {
			LinkTokenAddress:     testEnv.LinkTokenState.LinkToken.Address(),
			NativeTokenAddress:   common.HexToAddress("0x3e5e9111ae8eb78fe1cc3bb8915d5d461f3ef9a9"),
			VerifierProxyAddress: common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e"),
			RewardManagerAddress: common.HexToAddress("0x0fd8b81e3d1143ec7f1ce474827ab93c43523ea2"),
		}},
	}

	e, err := commonchangeset.Apply(t, e, nil,
		commonchangeset.Configure(
			feemanager.DeployFeeManagerChangeset,
			cfgFeeManager,
		),
	)

	require.NoError(t, err)

	// Ensure the FeeManager was deployed
	record, err := e.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(testutil.TestChain.Selector, datastore.ContractType(types.FeeManager), &deployment.Version0_5_0, ""),
	)
	require.NoError(t, err)
	feeManagerAddr := common.HexToAddress(record.Address)

	// Set Fee Manager on Verifier Proxy
	cfg := VerifierProxySetFeeManagerConfig{
		ConfigPerChain: map[uint64][]SetFeeManagerConfig{
			testChain: {
				{FeeManagerAddress: feeManagerAddr, VerifierProxyAddress: verifierProxyAddr},
			},
		},
	}
	e, err = commonchangeset.Apply(t, e, nil,
		commonchangeset.Configure(
			SetFeeManagerChangeset,
			cfg,
		),
	)
	require.NoError(t, err)

	t.Run("VerifyMetadata", func(t *testing.T) {
		// Use View To Confirm Data
		_, outputs, err := commonchangeset.ApplyChangesetsV2(t, e,
			[]commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(
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

		contractMetadata := testutil.MustGetContractMetaData[v0_5.VerifierProxyView](t, output.DataStore, testutil.TestChain.Selector, verifierProxyAddr.Hex())
		require.NotNil(t, contractMetadata)
		require.Equal(t, feeManagerAddr.Hex(), contractMetadata.View.FeeManager)
	})
}
