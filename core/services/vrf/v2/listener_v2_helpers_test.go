//go:build integration

package v2_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
)

func TestListener_EstimateFeeJuels(t *testing.T) {
	t.Parallel()
	callbackGasLimit := uint32(150_000)
	maxGasPriceGwei := assets.GWei(30).ToInt()
	weiPerUnitLink := big.NewInt(5898160000000000)
	actual, err := v2.EstimateFeeJuels(callbackGasLimit, maxGasPriceGwei, weiPerUnitLink)
	expected := big.NewInt(1780216203019246680)
	require.Equal(t, 0, actual.Cmp(expected), "expected:", expected.String(), "actual:", actual.String())
	require.NoError(t, err)

	weiPerUnitLink = big.NewInt(5898161234554321)
	actual, err = v2.EstimateFeeJuels(callbackGasLimit, maxGasPriceGwei, weiPerUnitLink)
	expected = big.NewInt(1780215830399116719)
	require.Equal(t, 0, actual.Cmp(expected), "expected:", expected.String(), "actual:", actual.String())
	require.NoError(t, err)

	actual, err = v2.EstimateFeeJuels(callbackGasLimit, maxGasPriceGwei, big.NewInt(0))
	require.Nil(t, actual)
	require.Error(t, err)
}

func TestListener_WorstCaseGasFloor(t *testing.T) {
	t.Parallel()
	// gasAfterPaymentCalculation from the Ethereum mainnet VRFCoordinatorV2_5 config.
	const mainnetGasAfterPaymentCalculation = uint32(38_900)

	tests := []struct {
		name                string
		callbackGasLimit    uint32
		gasAfterPaymentCalc uint32
	}{
		{"zero callback gas limit", 0, mainnetGasAfterPaymentCalculation},
		{"small callback gas limit", 20_000, mainnetGasAfterPaymentCalculation},
		{"typical callback gas limit", 100_000, mainnetGasAfterPaymentCalculation},
		{"larger callback gas limit", 250_000, mainnetGasAfterPaymentCalculation},
		{"large callback gas limit", 900_000, mainnetGasAfterPaymentCalculation},
		{"max gas limit", 2_500_000, mainnetGasAfterPaymentCalculation},
		{"extreme values must not overflow", math.MaxUint32, math.MaxUint32},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			floor := v2.WorstCaseGasFloor(tt.callbackGasLimit, tt.gasAfterPaymentCalc)
			require.GreaterOrEqual(t, floor,
				uint64(tt.callbackGasLimit)+uint64(v2.GasProofVerification)+uint64(tt.gasAfterPaymentCalc))
			if tt.callbackGasLimit > 0 {
				require.Greater(t, floor, v2.WorstCaseGasFloor(tt.callbackGasLimit-1, tt.gasAfterPaymentCalc))
			}
		})
	}

	// A max-size request stays within the gas limit Start() treats as required,
	// so correctly configured nodes are unaffected.
	require.LessOrEqual(t, v2.WorstCaseGasFloor(2_500_000, mainnetGasAfterPaymentCalculation),
		uint64(2_500_000)+uint64(v2.GasProofVerification)*2)
}

func TestListener_FallbackGasAfterPaymentCalculation(t *testing.T) {
	t.Parallel()
	// The fallback must exceed gasAfterPaymentCalculation on every coordinator we
	// deploy to, so the floor it produces is never undersized.
	deployedValues := map[string]uint32{
		"ethereum":  38_900,
		"base":      42_500,
		"optimism":  42_500,
		"soneium":   42_500,
		"bsc":       48_500,
		"polygon":   48_500,
		"ronin":     48_500,
		"arbitrum":  56_000,
		"avalanche": 56_000,
	}
	for chain, gasAfterPaymentCalculation := range deployedValues {
		require.Greater(t, v2.FallbackGasAfterPaymentCalculation, gasAfterPaymentCalculation,
			"fallback must exceed the value deployed on %s", chain)
		require.Greater(t,
			v2.WorstCaseGasFloor(100_000, v2.FallbackGasAfterPaymentCalculation),
			v2.WorstCaseGasFloor(100_000, gasAfterPaymentCalculation),
			"fallback floor must be conservative relative to %s", chain)
	}

	// The fallback must not push fulfillments past the node's configured limit.
	require.LessOrEqual(t, v2.WorstCaseGasFloor(2_500_000, v2.FallbackGasAfterPaymentCalculation),
		uint64(2_500_000)+uint64(v2.GasProofVerification)*2)
}

func Test_TxListDeduper(t *testing.T) {
	t.Parallel()
	tx1 := &txmgr.Tx{
		ID:      1,
		Value:   *big.NewInt(0),
		ChainID: big.NewInt(0),
	}
	tx2 := &txmgr.Tx{
		ID:      1,
		Value:   *big.NewInt(1),
		ChainID: big.NewInt(0),
	}
	tx3 := &txmgr.Tx{
		ID:      2,
		Value:   *big.NewInt(1),
		ChainID: big.NewInt(0),
	}
	txList := vrfcommon.DedupeTxList([]*txmgr.Tx{tx1, tx2, tx3})
	require.Len(t, txList, 2)
}
