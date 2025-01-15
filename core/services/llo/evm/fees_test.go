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
		baseUSDFee := decimal.NewFromInt32(0)
		fee := CalculateFee(tokenPriceInUSD, baseUSDFee)
		assert.Equal(t, big.NewInt(0), fee)
	})

	t.Run("negative fee rounds up to zero", func(t *testing.T) {
		tokenPriceInUSD := decimal.NewFromInt32(-123)
		baseUSDFee := decimal.NewFromInt32(1)
		fee := CalculateFee(tokenPriceInUSD, baseUSDFee)
		assert.Equal(t, big.NewInt(0), fee)

		tokenPriceInUSD = decimal.NewFromInt32(123)
		baseUSDFee = decimal.NewFromInt32(-1)
		fee = CalculateFee(tokenPriceInUSD, baseUSDFee)
		assert.Equal(t, big.NewInt(0), fee)

		// Multiple negative values also return a zero fee since negative
		// prices are always nonsensical
		tokenPriceInUSD = decimal.NewFromInt32(-123)
		baseUSDFee = decimal.NewFromInt32(-1)
		fee = CalculateFee(tokenPriceInUSD, baseUSDFee)
		assert.Equal(t, big.NewInt(0), fee)
	})

	t.Run("ridiculously high value rounds down fee to zero", func(t *testing.T) {
		// 20dp
		tokenPriceInUSD, err := decimal.NewFromString("12984833000000000000")
		require.NoError(t, err)
		BaseUSDFee, err = decimal.NewFromString("0.1")
		require.NoError(t, err)
		fee := CalculateFee(tokenPriceInUSD, BaseUSDFee)
		assert.Equal(t, big.NewInt(0), fee)
	})
}
