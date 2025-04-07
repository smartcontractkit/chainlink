package verification

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestDeployVerifierProxy(t *testing.T) {
	testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{ShouldDeployMCMS: true})

	cc := DeployVerifierProxyConfig{
		ChainsToDeploy: map[uint64]DeployVerifierProxy{
			testutil.TestChain.Selector: {AccessControllerAddress: common.Address{}},
		},
		Ownership: types.OwnershipSettings{
			ShouldTransfer: true,
			MCMSProposalConfig: &proposalutils.TimelockConfig{
				MinDelay: 0,
			},
		},
		Version: *semver.MustParse("0.5.0"),
	}

	e, err := commonChangesets.Apply(t, testEnv.Environment, testEnv.Timelocks,
		commonChangesets.Configure(
			DeployVerifierProxyChangeset,
			cc,
		),
	)
	require.NoError(t, err)

	verifierProxyAddr, err := dsutil.MaybeFindEthAddress(e.ExistingAddresses, testutil.TestChain.Selector, types.VerifierProxy)
	require.NoError(t, err)

	chain := e.Chains[testutil.TestChain.Selector]

	owner, _, err := commonChangesets.LoadOwnableContract(verifierProxyAddr, chain.Client)
	require.NoError(t, err)
	require.Equal(t, testEnv.Timelocks[testutil.TestChain.Selector].Timelock.Address(), owner)
}
