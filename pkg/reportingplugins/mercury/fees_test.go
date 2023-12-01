package mercury

import (
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// simulate price from DP
func scalePrice(usdPrice float64) *big.Int {
	scaledPrice := new(big.Float).Mul(big.NewFloat(usdPrice), big.NewFloat(1e18))
	scaledPriceInt, _ := scaledPrice.Int(nil)
	return scaledPriceInt
}

func Test_Fees(t *testing.T) {
	BaseUSDFee, err := decimal.NewFromString("0.70")
	require.NoError(t, err)
	t.Run("with token price > 1", func(t *testing.T) {
		tokenPriceInUSD := scalePrice(1630)
		fee := CalculateFee(tokenPriceInUSD, BaseUSDFee)
		expectedFee := big.NewInt(429447852760700) // 0.0004294478527607 18 decimals
		if fee.Cmp(expectedFee) != 0 {
			t.Errorf("Expected fee to be %v, got %v", expectedFee, fee)
		}
	})

	t.Run("with token price < 1", func(t *testing.T) {
		tokenPriceInUSD := scalePrice(0.4)
		fee := CalculateFee(tokenPriceInUSD, BaseUSDFee)
		expectedFee := big.NewInt(1750000000000000000) // 1.75 18 decimals
		if fee.Cmp(expectedFee) != 0 {
			t.Errorf("Expected fee to be %v, got %v", expectedFee, fee)
		}
	})

	t.Run("with token price == 0", func(t *testing.T) {
		tokenPriceInUSD := scalePrice(0)
		fee := CalculateFee(tokenPriceInUSD, BaseUSDFee)
		assert.Equal(t, big.NewInt(0), fee)
	})

	t.Run("with base fee == 0", func(t *testing.T) {
		tokenPriceInUSD := scalePrice(123)
		BaseUSDFee = decimal.NewFromInt32(0)
		fee := CalculateFee(tokenPriceInUSD, BaseUSDFee)
		assert.Equal(t, big.NewInt(0), fee)
	})
}
