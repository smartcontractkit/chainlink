package verification

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	feemanager "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/fee-manager"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestSetFeeManager(t *testing.T) {
	e := testutil.NewMemoryEnv(t, false, 0)

	testChain := testutil.TestChain.Selector
	e, err := commonchangeset.Apply(t, e, nil,
		commonchangeset.Configure(
			DeployVerifierProxyChangeset,
			DeployVerifierProxyConfig{
				ChainsToDeploy: map[uint64]DeployVerifierProxy{
					testChain: {AccessControllerAddress: common.Address{}},
				},
				Version: *semver.MustParse("0.5.0"),
			},
		),
	)
	require.NoError(t, err)

	// Ensure the VerifierProxy was deployed
	verifierProxyAddrHex, err := deployment.SearchAddressBook(e.ExistingAddresses, testChain, types.VerifierProxy)
	require.NoError(t, err)
	verifierProxyAddr := common.HexToAddress(verifierProxyAddrHex)

	// Need the Link Token For FeeManager
	e, err = commonchangeset.Apply(t, e, nil,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
			[]uint64{testutil.TestChain.Selector},
		),
	)
	require.NoError(t, err)

	addresses, err := e.ExistingAddresses.AddressesForChain(testutil.TestChain.Selector)
	require.NoError(t, err)

	chain := e.Chains[testutil.TestChain.Selector]
	linkState, err := commonstate.MaybeLoadLinkTokenChainState(chain, addresses)
	require.NoError(t, err)

	// Deploy Fee Manager
	cfgFeeManager := feemanager.DeployFeeManagerConfig{
		ChainsToDeploy: map[uint64]feemanager.DeployFeeManager{testutil.TestChain.Selector: {
			LinkTokenAddress:     linkState.LinkToken.Address(),
			NativeTokenAddress:   common.HexToAddress("0x3e5e9111ae8eb78fe1cc3bb8915d5d461f3ef9a9"),
			VerifierProxyAddress: common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e"),
			RewardManagerAddress: common.HexToAddress("0x0fd8b81e3d1143ec7f1ce474827ab93c43523ea2"),
		}},
	}

	e, err = commonchangeset.Apply(t, e, nil,
		commonchangeset.Configure(
			feemanager.DeployFeeManagerChangeset,
			cfgFeeManager,
		),
	)

	require.NoError(t, err)

	// Ensure the FeeManager was deployed
	feeManagerAddrHex, err := deployment.SearchAddressBook(e.ExistingAddresses, testChain, types.FeeManager)
	require.NoError(t, err)
	feeManagerAddr := common.HexToAddress(feeManagerAddrHex)

	// Set Fee Manager on Verifier Proxy
	cfg := VerifierProxySetFeeManagerConfig{
		ConfigPerChain: map[uint64][]SetFeeManagerConfig{
			testChain: {
				{FeeManagerAddress: feeManagerAddr, ContractAddress: verifierProxyAddr},
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

	client := e.Chains[testChain].Client
	verifierProxy, err := verifier_proxy_v0_5_0.NewVerifierProxy(verifierProxyAddr, client)
	require.NoError(t, err)

	// Check VerifierProxy has the correct FeeManager set
	feeManager, err := verifierProxy.SFeeManager(nil)
	require.NoError(t, err)
	require.Equal(t, feeManagerAddr, feeManager)
}
