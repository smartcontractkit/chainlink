package evm

import (
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Fees(t *testing.T) {
	BaseUSDFee, err := decimal.NewFromString("0.70")
	require.NoError(t, err)
	t.Run("with token price > 1", func(t *testing.T) {
		tokenPriceInUSD := decimal.NewFromInt32(1630)
		fee := CalculateFee(tokenPriceInUSD, BaseUSDFee)
		expectedFee := big.NewInt(429447852760736)
		if fee.Cmp(expectedFee) != 0 {
			t.Errorf("Expected fee to be %v, got %v", expectedFee, fee)
		}
	})

	t.Run("with token price < 1", func(t *testing.T) {
		tokenPriceInUSD := decimal.NewFromFloat32(0.4)
		fee := CalculateFee(tokenPriceInUSD, BaseUSDFee)
		expectedFee := big.NewInt(1750000000000000000)
		if fee.Cmp(expectedFee) != 0 {
			t.Errorf("Expected fee to be %v, got %v", expectedFee, fee)
		}
	})

	t.Run("with token price == 0", func(t *testing.T) {
		tokenPriceInUSD := decimal.NewFromInt32(0)
		fee := CalculateFee(tokenPriceInUSD, BaseUSDFee)
		assert.Equal(t, big.NewInt(0), fee)
	})

	t.Run("with base fee == 0", func(t *testing.T) {
		tokenPriceInUSD := decimal.NewFromInt32(123)
		BaseUSDFee = decimal.NewFromInt32(0)
		fee := CalculateFee(tokenPriceInUSD, BaseUSDFee)
		assert.Equal(t, big.NewInt(0), fee)
	})
}
