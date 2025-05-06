package fee_manager

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestDeployFeeManager(t *testing.T) {
	testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{ShouldDeployMCMS: true})

	// Need the Link Token
	e, err := commonChangesets.Apply(t, testEnv.Environment, nil,
		commonChangesets.Configure(
			cldf.CreateLegacyChangeSet(commonChangesets.DeployLinkToken),
			[]uint64{testutil.TestChain.Selector},
		),
	)
	require.NoError(t, err)

	addresses, err := e.ExistingAddresses.AddressesForChain(testutil.TestChain.Selector)
	require.NoError(t, err)

	chain := e.Chains[testutil.TestChain.Selector]
	linkState, err := commonstate.MaybeLoadLinkTokenChainState(chain, addresses)
	require.NoError(t, err)

	cc := DeployFeeManagerConfig{
		ChainsToDeploy: map[uint64]DeployFeeManager{testutil.TestChain.Selector: {
			LinkTokenAddress:     linkState.LinkToken.Address(),
			NativeTokenAddress:   common.HexToAddress("0x3e5e9111ae8eb78fe1cc3bb8915d5d461f3ef9a9"),
			VerifierProxyAddress: common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e"),
			RewardManagerAddress: common.HexToAddress("0x0fd8b81e3d1143ec7f1ce474827ab93c43523ea2"),
		}},
		Ownership: types.OwnershipSettings{
			ShouldTransfer: true,
			MCMSProposalConfig: &proposalutils.TimelockConfig{
				MinDelay: 0,
			},
		},
	}

	resp, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(DeployFeeManagerChangeset, cc),
	)

	require.NoError(t, err)

	// Check the address book for fm existence
	fmAddr, err := dsutil.MaybeFindEthAddress(resp.ExistingAddresses, testutil.TestChain.Selector, types.FeeManager)
	require.NoError(t, err)

	owner, _, err := commonChangesets.LoadOwnableContract(fmAddr, chain.Client)
	require.NoError(t, err)
	require.Equal(t, testEnv.Timelocks[testutil.TestChain.Selector].Timelock.Address(), owner)
}
