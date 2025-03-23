package v0_5_0

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	verifierproxy "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/verifier-proxy/v0_5_0"
)

func TestDeployVerifier(t *testing.T) {
	e := testutil.NewMemoryEnv(t, true)

	cc := verifierproxy.DeployVerifierProxyConfig{
		ChainsToDeploy: map[uint64]verifierproxy.DeployVerifierProxy{
			testutil.TestChain.Selector: {AccessControllerAddress: common.Address{}},
		},
	}

	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			verifierproxy.DeployVerifierProxyChangeset,
			cc,
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
