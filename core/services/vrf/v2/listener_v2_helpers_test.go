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
	maxGasPriceWei := assets.GWei(30).ToInt()

	tests := []struct {
		name              string
		gasOverhead       uint32
		weiPerUnitLink    *big.Int
		premiumPercentage uint32
		expected          *big.Int
	}{
		{"no premium", v2.GasProofVerification, big.NewInt(5898160000000000), 0, big.NewInt(1780216203019246680)},
		{"default premium", v2.GasProofVerification, big.NewInt(5898160000000000), 20, big.NewInt(2136259443623096016)},
		{"max premium", v2.GasProofVerification, big.NewInt(5898160000000000), 155, big.NewInt(4539551317699079034)},
		{"no premium, different feed", v2.GasProofVerification, big.NewInt(5898161234554321), 0, big.NewInt(1780215830399116719)},
		{"default premium, different feed", v2.GasProofVerification, big.NewInt(5898161234554321), 20, big.NewInt(2136258996478940062)},
		{"max premium, different feed", v2.GasProofVerification, big.NewInt(5898161234554321), 155, big.NewInt(4539550367517747633)},
		{"gas overhead override", 250_000, big.NewInt(5898160000000000), 20, big.NewInt(2441439364140681161)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := v2.EstimateFeeJuels(callbackGasLimit, tt.gasOverhead, maxGasPriceWei, tt.weiPerUnitLink, tt.premiumPercentage)
			require.NoError(t, err)
			require.Equal(t, 0, actual.Cmp(tt.expected), "expected:", tt.expected.String(), "actual:", actual.String())
		})
	}

	actual, err := v2.EstimateFeeJuels(callbackGasLimit, v2.GasProofVerification, maxGasPriceWei, big.NewInt(0), 20)
	require.Nil(t, actual)
	require.Error(t, err)
}

func TestListener_EstimateFeeWei(t *testing.T) {
	t.Parallel()
	callbackGasLimit := uint32(150_000)
	maxGasPriceWei := assets.GWei(30).ToInt()

	tests := []struct {
		name              string
		gasOverhead       uint32
		premiumPercentage uint32
		expected          *big.Int
	}{
		{"no premium", v2.GasProofVerification, 0, big.NewInt(10500000000000000)},
		{"default premium", v2.GasProofVerification, 20, big.NewInt(12600000000000000)},
		{"max premium", v2.GasProofVerification, 155, big.NewInt(26775000000000000)},
		{"gas overhead override", 250_000, 20, big.NewInt(14400000000000000)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := v2.EstimateFeeWei(callbackGasLimit, tt.gasOverhead, maxGasPriceWei, tt.premiumPercentage)
			require.NoError(t, err)
			require.Equal(t, 0, actual.Cmp(tt.expected), "expected:", tt.expected.String(), "actual:", actual.String())
		})
	}
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
