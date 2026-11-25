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
	tests := []struct {
		name             string
		callbackGasLimit uint32
		gasOverhead      uint32
	}{
		{"zero callback gas limit", 0, v2.GasProofVerification},
		{"small callback gas limit", 20_000, v2.GasProofVerification},
		{"typical callback gas limit", 100_000, v2.GasProofVerification},
		{"larger callback gas limit", 250_000, v2.GasProofVerification},
		{"large callback gas limit", 900_000, v2.GasProofVerification},
		{"max gas limit", 2_500_000, v2.GasProofVerification},
		{"job spec override", 100_000, 250_000},
		{"extreme values must not overflow", math.MaxUint32, math.MaxUint32},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			floor := v2.WorstCaseGasFloor(tt.callbackGasLimit, tt.gasOverhead)
			require.GreaterOrEqual(t, floor,
				uint64(tt.callbackGasLimit)*64/63+uint64(tt.gasOverhead))
			if tt.callbackGasLimit > 0 {
				require.Greater(t, floor, v2.WorstCaseGasFloor(tt.callbackGasLimit-1, tt.gasOverhead))
			}
		})
	}

	// The default overhead must cover the highest coordinator overhead measured
	// across chains at callbackGasLimit=0 (~165k gas), with headroom.
	require.Greater(t, v2.GasProofVerification, uint32(165_000))

	// A max-size request stays within the gas limit Start() treats as required,
	// so correctly configured nodes are unaffected.
	require.LessOrEqual(t, v2.WorstCaseGasFloor(2_500_000, v2.GasProofVerification),
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
