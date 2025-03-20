package v0_5_0

import (
	"testing"

	"github.com/stretchr/testify/require"

	dsTypes "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	verifier_proxy "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/verifier-proxy/v0_5_0"
)

// DeployVerifierProxyAndVerifier deploys a VerifierProxy, deploys a Verifier,
// initializes the VerifierProxy with the Verifier, and returns the updated
// environment and the addresses of VerifierProxy and Verifier.
func DeployVerifierProxyAndVerifier(
	t *testing.T,
	e deployment.Environment,
) (env deployment.Environment, verifierProxyAddr common.Address, verifierAddr common.Address) {
	t.Helper()

	chainSelector := testutil.TestChain.Selector

	// 1) Deploy VerifierProxy
	deployProxyCfg := verifier_proxy.DeployVerifierProxyConfig{
		ChainsToDeploy: map[uint64]verifier_proxy.DeployVerifierProxy{
			chainSelector: {
				VerifierProxyAddress: common.Address{},
			},
		},
	}
	env, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			verifier_proxy.DeployVerifierProxyChangeset,
			deployProxyCfg,
		),
	)
	require.NoError(t, err, "deploying verifier proxy should not fail")

	// Get the VerifierProxy address
	verifierProxyAddrHex, err := deployment.SearchAddressBook(env.ExistingAddresses, chainSelector, dsTypes.VerifierProxy)
	require.NoError(t, err, "unable to find verifier proxy address in address book")
	verifierProxyAddr = common.HexToAddress(verifierProxyAddrHex)
	require.NotEqual(t, common.Address{}, verifierProxyAddr, "verifier proxy should not be zero address")

	// 2) Deploy Verifier
	deployVerifierCfg := DeployVerifierConfig{
		ChainsToDeploy: map[uint64]DeployVerifier{
			chainSelector: {
				VerifierProxyAddress: verifierProxyAddr,
			},
		},
	}
	env, err = commonChangesets.Apply(t, env, nil,
		commonChangesets.Configure(
			DeployVerifierChangeset,
			deployVerifierCfg,
		),
	)
	require.NoError(t, err, "deploying verifier should not fail")

	// Get the Verifier address
	verifierAddrHex, err := deployment.SearchAddressBook(env.ExistingAddresses, chainSelector, dsTypes.Verifier)
	require.NoError(t, err, "unable to find verifier address in address book")
	verifierAddr = common.HexToAddress(verifierAddrHex)
	require.NotEqual(t, common.Address{}, verifierAddr, "verifier should not be zero address")

	// 3) Initialize the VerifierProxy
	initCfg := verifier_proxy.InitializeConfig{
		ConfigsByChain: map[uint64][]verifier_proxy.Initialize{
			chainSelector: {
				{
					VerifierAddress:      verifierAddr,
					VerifierProxyAddress: verifierProxyAddr,
				},
			},
		},
	}
	env, err = commonChangesets.Apply(t, env, nil,
		commonChangesets.Configure(
			verifier_proxy.InitializeChangeset,
			initCfg,
		),
	)
	require.NoError(t, err, "initializing verifier proxy should not fail")

	t.Logf("VerifierProxy deployed at %s, Verifier deployed at %s, and successfully initialized",
		verifierProxyAddrHex, verifierAddrHex)

	return env, verifierProxyAddr, verifierAddr
}
