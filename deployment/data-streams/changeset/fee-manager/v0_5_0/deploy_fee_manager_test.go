package v0_5_0

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestDeployFeeManager(t *testing.T) {
	e := testutil.NewMemoryEnv(t, false, 0)

	// Need the Link Token
	e, err := commonChangesets.Apply(t, e, nil,
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

	cc := DeployFeeManagerConfig{
		ChainsToDeploy: map[uint64]DeployFeeManager{testutil.TestChain.Selector: {
			LinkTokenAddress:     linkState.LinkToken.Address(),
			NativeTokenAddress:   common.HexToAddress("0x3e5e9111ae8eb78fe1cc3bb8915d5d461f3ef9a9"),
			ProxyAddress:         common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e"),
			RewardManagerAddress: common.HexToAddress("0x0fd8b81e3d1143ec7f1ce474827ab93c43523ea2"),
		}},
	}

	resp, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(DeployFeeManagerChangeset, cc),
	)

	require.NoError(t, err)

	// Check the address book for fm existence
	chainAddresses, err := resp.ExistingAddresses.AddressesForChain(testutil.TestChain.Selector)
	require.NoError(t, err)

	var fmAddress common.Address
	for addr, tv := range chainAddresses {
		if tv.Type == types.FeeManager {
			fmAddress = common.HexToAddress(addr)
			break
		}
	}
	require.NotEqual(t, "", fmAddress)
	require.NotEqual(t, common.HexToAddress("0x0000000000000000000000000000000000000000").String(), fmAddress)

}
