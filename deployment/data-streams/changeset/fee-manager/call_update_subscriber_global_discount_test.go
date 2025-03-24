package fee_manager

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

func TestUpdateSubscriberGlobalDiscount(t *testing.T) {
	res, err := NewDataStreamsEnvironment(t, NewDefaultOptions())
	require.NoError(t, err)

	linkTokenAddress := res.LinkTokenAddress
	feeManagerAddress := res.FeeManagerAddress
	e := res.Env

	chain := e.Chains[testutil.TestChain.Selector]
	require.NotNil(t, chain)

	subscriber := common.HexToAddress("0x0fd8b81e3d1143ec7f1ce474827ab93c43523ea2")

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			UpdateSubscriberGlobalDiscountChangeset,
			UpdateSubscriberGlobalDiscountConfig{
				ConfigPerChain: map[uint64][]UpdateSubscriberGlobalDiscount{
					testutil.TestChain.Selector: {
						{FeeManagerAddress: feeManagerAddress,
							SubscriberAddress: subscriber,
							TokenAddress:      linkTokenAddress,
							Discount:          2000,
						},
					},
				},
			},
		))
	require.NoError(t, err)

	feeManager, err := LoadFeeManagerState(e, testutil.TestChain.Selector, feeManagerAddress.String())
	require.NoError(t, err)
	require.NotNil(t, feeManager)

	actualDiscount, err := feeManager.SGlobalDiscounts(nil, subscriber, linkTokenAddress)

	require.NoError(t, err)
	require.Equal(t, actualDiscount, big.NewInt(2000))
}
