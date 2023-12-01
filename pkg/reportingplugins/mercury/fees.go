package mercury

import (
	"math/big"

	"github.com/shopspring/decimal"
)

// PriceScalingFactor indicates the multiplier applied to token prices that we expect from data source
// e.g. for a 1e8 multiplier, a LINK/USD value of 7.42 will be derived from a data source value of 742000000
var PRICE_SCALING_FACTOR = decimal.NewFromInt(1e18) //nolint:revive

// FeeScalingFactor indicates the multiplier applied to fees.
// e.g. for a 1e18 multiplier, a LINK fee of 7.42 will be represented as 7.42e18
// This is what will be baked into the report for use on-chain.
var FEE_SCALING_FACTOR = decimal.NewFromInt(1e18) //nolint:revive

// CalculateFee outputs a fee in wei according to the formula: baseUSDFee * scaleFactor / tokenPriceInUSD
func CalculateFee(tokenPriceInUSD *big.Int, baseUSDFee decimal.Decimal) *big.Int {
	if tokenPriceInUSD.Cmp(big.NewInt(0)) == 0 || baseUSDFee.IsZero() {
		// zero fee if token price or base fee is zero
		return big.NewInt(0)
	}

	// scale baseFee in USD
	baseFeeScaled := baseUSDFee.Mul(PRICE_SCALING_FACTOR)

	tokenPrice := decimal.NewFromBigInt(tokenPriceInUSD, 0)

	// fee denominated in token
	fee := baseFeeScaled.Div(tokenPrice)

	// scale fee to the expected format
	fee = fee.Mul(FEE_SCALING_FACTOR)

	// convert to big.Int
	return fee.BigInt()
}
