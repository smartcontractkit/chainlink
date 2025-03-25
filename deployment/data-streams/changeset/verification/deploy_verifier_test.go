package verification

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestDeployVerifier(t *testing.T) {
	e := testutil.NewMemoryEnv(t, true)

	cc := DeployVerifierProxyConfig{
		ChainsToDeploy: map[uint64]DeployVerifierProxy{
			testutil.TestChain.Selector: {AccessControllerAddress: common.Address{}},
		},
		Version: deployment.Version0_5_0,
	}

	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			DeployVerifierProxyChangeset,
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
