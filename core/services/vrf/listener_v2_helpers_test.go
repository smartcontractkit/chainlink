package vrf_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
)

func TestListener_EstimateFeeJuels(t *testing.T) {
	callbackGasLimit := uint32(150_000)
	maxGasPriceGwei := assets.GWei(30)
	weiPerUnitLink := big.NewInt(5898160000000000)
	actual := vrf.EstimateFeeJuels(callbackGasLimit, maxGasPriceGwei, weiPerUnitLink)
	expected := big.NewInt(1780216203019246680)
	require.True(t, actual.Cmp(expected) == 0, "expected:", expected.String(), "actual:", actual.String())

	weiPerUnitLink = big.NewInt(5898161234554321)
	actual = vrf.EstimateFeeJuels(callbackGasLimit, maxGasPriceGwei, weiPerUnitLink)
	expected = big.NewInt(1780215830399116719)
	require.True(t, actual.Cmp(expected) == 0, "expected:", expected.String(), "actual:", actual.String())
}
