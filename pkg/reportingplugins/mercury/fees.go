package mercury

import "math/big"

// PriceScalingFactor indicates the multiplier applied to token prices.
// e.g. for a 1e8 multiplier, a LINK/USD value of 7.42 will be represented as 742000000
// This is what we expect from our data source.
var PRICE_SCALING_FACTOR = big.NewInt(1e8)

// FeeScalingFactor indicates the multiplier applied to fees.
// e.g. for a 1e18 multiplier, a LINK fee of 7.42 will be represented as 7.42e18
// This is what will be baked into the report for use on-chain.
var FEE_SCALING_FACTOR = big.NewInt(1e18)

var CENTS_PER_DOLLAR = big.NewInt(100)

// CalculateFee outputs a fee in wei
func CalculateFee(tokenPriceInUSD *big.Int, baseUSDFeeCents uint32) (fee *big.Int) {
	fee = new(big.Int).Mul(big.NewInt(int64(baseUSDFeeCents)), tokenPriceInUSD)
	fee = fee.Mul(fee, FEE_SCALING_FACTOR)
	fee = fee.Div(fee, PRICE_SCALING_FACTOR)
	fee = fee.Div(fee, CENTS_PER_DOLLAR)
	return
}
