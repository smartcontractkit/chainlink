package mercury

import (
	"math/big"
	"testing"
)

func Test_Fees(t *testing.T) {
	t.Run("CalculateFee", func(t *testing.T) {
		tokenPriceInUSD := big.NewInt(655000000)
		var baseUSDFeeCents uint32 = 100
		fee := CalculateFee(tokenPriceInUSD, baseUSDFeeCents)
		if fee.Cmp(big.NewInt(6.55e18)) != 0 {
			t.Errorf("Expected fee to be 6550000000000000000, got %v", fee)
		}
	})
}
