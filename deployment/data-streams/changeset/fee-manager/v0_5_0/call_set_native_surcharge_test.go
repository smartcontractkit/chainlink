package fee_manager_v0_5_0

import (
	"math/big"
	"testing"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/stretchr/testify/require"
)

func TestSetNativeSurcharge(t *testing.T) {

	res, err := NewDataStreamsEnvironment(t, NewDefaultOptions())
	require.NoError(t, err)

	feeManagerAddress := res.FeeManagerAddress
	e := res.Env

	chain := e.Chains[testutil.TestChain.Selector]
	require.NotNil(t, chain)

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			SetNativeChangeset,
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

	feeManager, err := LoadFeeManagerState(e, testutil.TestChain.Selector, feeManagerAddress.String())
	require.NoError(t, err)
	require.NotNil(t, feeManager)

	actualNativeSurcharge, err := feeManager.SNativeSurcharge(nil)
	require.NoError(t, err)
	require.Equal(t, actualNativeSurcharge, big.NewInt(5000))

}
