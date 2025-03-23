package fee_manager_v0_5_0

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/stretchr/testify/require"
)

func TestUpdateSubscriberDiscount(t *testing.T) {

	res, err := NewDataStreamsEnvironment(t, NewDefaultOptions())
	require.NoError(t, err)

	linkTokenAddress := res.LinkTokenAddress
	feeManagerAddress := res.FeeManagerAddress
	e := res.Env

	chain := e.Chains[testutil.TestChain.Selector]
	require.NotNil(t, chain)

	subscriber := common.HexToAddress("0x0fd8b81e3d1143ec7f1ce474827ab93c43523ea2")

	feedId := [32]byte{1}

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			UpdateSubscriberDiscountChangeset,
			UpdateSubscriberDiscountConfig{
				ConfigPerChain: map[uint64][]UpdateSubscriberDiscount{
					testutil.TestChain.Selector: {
						{FeeManagerAddress: feeManagerAddress,
							SubscriberAddress: subscriber,
							FeedId:            feedId,
							TokenAddress:      linkTokenAddress,
							Discount:          1000,
						},
					},
				},
			},
		))

	feeManager, err := LoadFeeManagerState(e, testutil.TestChain.Selector, feeManagerAddress.String())
	require.NoError(t, err)
	require.NotNil(t, feeManager)

	actualDiscount, err := feeManager.SSubscriberDiscounts(nil, subscriber, feedId, linkTokenAddress)

	require.NoError(t, err)
	require.Equal(t, actualDiscount, big.NewInt(1000))
}
