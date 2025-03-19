package v0_5_0

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestDeployVerifierProxy(t *testing.T) {
	e := testutil.NewMemoryEnv(t, false, 0)
	cc := DeployVerifierProxyConfig{
		ChainsToDeploy: []uint64{testutil.TestChain.Selector},
		Version:        *semver.MustParse("0.5.0"),
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
	verifierAddr := common.HexToAddress(verifierProxyAddrHex)
	require.NotEqual(t, common.HexToAddress("0x0000000000000000000000000000000000000000"), verifierAddr)
}
