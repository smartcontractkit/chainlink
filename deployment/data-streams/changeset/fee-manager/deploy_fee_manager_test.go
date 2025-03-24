<<<<<<<< HEAD:deployment/data-streams/changeset/fee-manager/deploy_fee_manager_test.go
package fee_manager
========
package v0_5_0
>>>>>>>> d6421371295beb55cfa19e666e24812b38317e63:deployment/data-streams/changeset/fee-manager/v0_5_0/deploy_fee_manager_test.go

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestDeployFeeManager(t *testing.T) {
	e := testutil.NewMemoryEnv(t, false, 0)

	// Need the Link Token
	e, err := commonchangeset.Apply(t, e, nil,
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
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
<<<<<<<< HEAD:deployment/data-streams/changeset/fee-manager/deploy_fee_manager_test.go
		ChainsToDeploy:       []uint64{testutil.TestChain.Selector},
		LinkTokenAddress:     linkState.LinkToken.Address(),
		NativeTokenAddress:   common.HexToAddress("0x3e5e9111ae8eb78fe1cc3bb8915d5d461f3ef9a9"),
		VerifierProxyAddress: common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e"),
		RewardManagerAddress: common.HexToAddress("0x0fd8b81e3d1143ec7f1ce474827ab93c43523ea2"),
	}

	resp, err := commonchangeset.Apply(t, e, nil,
		commonchangeset.Configure(
			FeeManagerDeploy{},
			cc,
		),
========
		ChainsToDeploy: map[uint64]DeployFeeManager{testutil.TestChain.Selector: {
			LinkTokenAddress:     linkState.LinkToken.Address(),
			NativeTokenAddress:   common.HexToAddress("0x3e5e9111ae8eb78fe1cc3bb8915d5d461f3ef9a9"),
			ProxyAddress:         common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e"),
			RewardManagerAddress: common.HexToAddress("0x0fd8b81e3d1143ec7f1ce474827ab93c43523ea2"),
		}},
	}

	resp, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(DeployFeeManagerChangeset, cc),
>>>>>>>> d6421371295beb55cfa19e666e24812b38317e63:deployment/data-streams/changeset/fee-manager/v0_5_0/deploy_fee_manager_test.go
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
