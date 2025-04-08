package fee_manager

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
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

	actualNativeSurcharge, err := feeManager.SNativeSurcharge(nil)
	require.NoError(t, err)
	require.Equal(t, actualNativeSurcharge, big.NewInt(5000))
}
