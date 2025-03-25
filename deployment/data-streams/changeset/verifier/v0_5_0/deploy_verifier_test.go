package v0_5_0

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	verifier_proxy "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/verifier-proxy/v0.5.0"
)

func TestDeployVerifier(t *testing.T) {
	e := testutil.NewMemoryEnv(t, true)
	chainSelector := e.AllChainSelectors()[0]

	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			verifier_proxy.DeployVerifierProxyChangeset,
			verifier_proxy.DeployVerifierProxyConfig{
				ChainsToDeploy: []uint64{chainSelector},
				Version:        *semver.MustParse("0.5.0"),
			},
		),
	)

	require.NoError(t, err)

	verifierProxyAddrHex, err := deployment.SearchAddressBook(e.ExistingAddresses, testutil.TestChain.Selector, types.VerifierProxy)
	require.NoError(t, err)
	verifierProxyAddr := common.HexToAddress(verifierProxyAddrHex)

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			DeployVerifierChangeset,
			DeployVerifierConfig{
				ChainsToDeploy: map[uint64]DeployVerifier{
					testutil.TestChain.Selector: {VerifierProxyAddress: verifierProxyAddr},
				},
			},
		),
	)

	require.NoError(t, err)

	_, err = deployment.SearchAddressBook(e.ExistingAddresses, testutil.TestChain.Selector, types.Verifier)
	require.NoError(t, err)
}
