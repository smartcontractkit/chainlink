package v0_5_0

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	feeManagerCs "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/fee-manager/v0_5_0"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"
)

func TestSetFeeManager(t *testing.T) {
	e := testutil.NewMemoryEnv(t, false, 0)

	testChain := testutil.TestChain.Selector
	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			DeployVerifierProxyChangeset,
			DeployVerifierProxyConfig{
				ChainsToDeploy: []uint64{testChain},
				Version:        *semver.MustParse("0.5.0"),
			},
		),
	)
	require.NoError(t, err)

	// Ensure the VerifierProxy was deployed
	verifierProxyAddrHex, err := deployment.SearchAddressBook(e.ExistingAddresses, testChain, types.VerifierProxy)
	require.NoError(t, err)
	verifierProxyAddr := common.HexToAddress(verifierProxyAddrHex)

	// Need the Link Token For FeeManager
	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			deployment.CreateLegacyChangeSet(commonChangesets.DeployLinkToken),
			[]uint64{testutil.TestChain.Selector},
		),
	)
	require.NoError(t, err)

	addresses, err := e.ExistingAddresses.AddressesForChain(testutil.TestChain.Selector)
	require.NoError(t, err)

	chain := e.Chains[testutil.TestChain.Selector]
	linkState, err := commonChangesets.MaybeLoadLinkTokenChainState(chain, addresses)
	require.NoError(t, err)

	// Deploy Fee Manager
	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			feeManagerCs.DeployFeeManagerChangeset,
			feeManagerCs.DeployFeeManagerConfig{
				ChainsToDeploy: map[uint64]feeManagerCs.DeployFeeManager{
					testutil.TestChain.Selector: {
						LinkTokenAddress:     linkState.LinkToken.Address(),
						NativeTokenAddress:   common.HexToAddress("0x3e5e9111ae8eb78fe1cc3bb8915d5d461f3ef9a9"),
						ProxyAddress:         common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e"),
						RewardManagerAddress: common.HexToAddress("0x0fd8b81e3d1143ec7f1ce474827ab93c43523ea2"),
					},
				},
			},
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
	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
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
