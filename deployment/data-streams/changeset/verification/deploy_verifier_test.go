package verification

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_v0_5_0"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
)

func TestDeployVerifier(t *testing.T) {
	testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{DeployMCMS: true})

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

	verifierProxyAddrHex, err := deployment.SearchAddressBook(e.ExistingAddresses, testutil.TestChain.Selector, types.VerifierProxy)
	require.NoError(t, err)
	verifierProxyAddr := common.HexToAddress(verifierProxyAddrHex)

	// Specifying timelock configurations will execute any resulting proposals - in this case, the ownership transfer
	e, err = commonChangesets.Apply(t, e, testEnv.Timelocks,
		commonChangesets.Configure(
			DeployVerifierChangeset,
			DeployVerifierConfig{
				ChainsToDeploy: map[uint64]DeployVerifier{
					testutil.TestChain.Selector: {VerifierProxyAddress: verifierProxyAddr},
				},
				MCMSConfig: &proposalutils.TimelockConfig{
					MinDelay: 0,
				},
			},
		),
	)

	require.NoError(t, err)

	// Confirm address exists
	verifierAddr, err := dsutil.MaybeFindEthAddress(e.ExistingAddresses, testutil.TestChain.Selector, types.Verifier)
	require.NoError(t, err)

	// Confirm transfer to MCMS was successful
	chain := e.Chains[testutil.TestChain.Selector]
	verifier, err := verifier_v0_5_0.NewVerifier(verifierAddr, chain.Client)
	require.NoError(t, err)

	owner, err := verifier.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, testEnv.Timelocks[testutil.TestChain.Selector].Timelock.Address(), owner)

}
