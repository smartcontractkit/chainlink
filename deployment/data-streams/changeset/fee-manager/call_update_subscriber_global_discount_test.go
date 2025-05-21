package fee_manager

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

func TestUpdateSubscriberGlobalDiscount(t *testing.T) {
	res, err := DeployTestEnvironment(t, NewDefaultOptions())
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

	t.Run("VerifyMetadata", func(t *testing.T) {
		// Use View To Confirm Data
		_, outputs, err := commonChangesets.ApplyChangesetsV2(t, e,
			[]commonChangesets.ConfiguredChangeSet{
				commonChangesets.Configure(
					changeset.SaveContractViews,
					changeset.SaveContractViewsConfig{
						Chains: []uint64{testutil.TestChain.Selector},
					},
				),
			},
		)
		require.NoError(t, err)
		require.Len(t, outputs, 1)
		output := outputs[0]

		contractMetadata := testutil.MustGetContractMetaData[v0_5.FeeManagerView](t, output.DataStore, testutil.TestChain.Selector, feeManagerAddress.Hex())
		require.NotNil(t, contractMetadata)
		discountRecord, ok := contractMetadata.View.SubscriberDiscounts[subscriber.String()]["global"]
		require.True(t, ok)
		require.Equal(t, "2000", discountRecord.Link)
		require.True(t, discountRecord.IsGlobal)
	})
}
