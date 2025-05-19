package fee_manager

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
)

func TestSetNativeSurcharge(t *testing.T) {
	res, err := DeployTestEnvironment(t, NewDefaultOptions())
	require.NoError(t, err)

	feeManagerAddress := res.FeeManagerAddress
	e := res.Env

	chain := e.Chains[testutil.TestChain.Selector]
	require.NotNil(t, chain)

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			SetNativeSurchargeChangeset,
			SetNativeSurchargeConfig{
				ConfigPerChain: map[uint64][]SetNativeSurcharge{
					testutil.TestChain.Selector: {
						{
							FeeManagerAddress: feeManagerAddress,
							Surcharge:         5000,
						},
					},
				},
			},
		))
	require.NoError(t, err)

	feeManager, err := LoadFeeManagerState(e, testutil.TestChain.Selector, feeManagerAddress.String())
	require.NoError(t, err)
	require.NotNil(t, feeManager)

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

		require.Equal(t, strconv.Itoa(5000), contractMetadata.View.NativeSurcharge)
	})
}
