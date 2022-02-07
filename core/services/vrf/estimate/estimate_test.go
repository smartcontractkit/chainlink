package estimate

import (
	"math/big"
	"testing"

	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink/core/assets"
)

func TestListener_EstimateLinkNeeded(t *testing.T) {
	callbackGasLimit := uint32(150_000)
	maxGasPriceGwei := assets.GWei(30)
	weiPerUnitLink := big.NewInt(5898160000000000)
	actual := JuelsNeeded(callbackGasLimit, maxGasPriceGwei, weiPerUnitLink)
	expected := big.NewInt(1780216203019246680)
	require.True(t, actual.Cmp(expected) == 0, "expected:", expected.String(), "actual:", actual.String())
}
