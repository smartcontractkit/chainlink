package verification

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"
	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestSetAccessController(t *testing.T) {
	e := testutil.NewMemoryEnv(t, true, 0)
	acAddress := common.HexToAddress("0x0000000000000000000000000000000000000123")
	testChain := e.AllChainSelectors()[0]
	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
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

	cfg := VerifierProxySetAccessControllerConfig{
		ConfigPerChain: map[uint64][]SetAccessControllerConfig{
			testChain: {
				{AccessControllerAddress: acAddress, ContractAddress: verifierProxyAddr},
			},
		},
	}

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			SetAccessControllerChangeset,
			cfg,
		),
	)
	require.NoError(t, err)

	client := e.Chains[testChain].Client
	verifierProxy, err := verifier_proxy_v0_5_0.NewVerifierProxy(verifierProxyAddr, client)
	require.NoError(t, err)

	accessController, err := verifierProxy.SAccessController(nil)
	require.NoError(t, err)
	require.Equal(t, acAddress, accessController)
}
