package evm

import (
	"math/big"

	"github.com/shopspring/decimal"
)

// FeeScalingFactor indicates the multiplier applied to fees.
// e.g. for a 1e18 multiplier, a LINK fee of 7.42 will be represented as 7.42e18
// This is what will be baked into the report for use on-chain.
var FeeScalingFactor = decimal.NewFromInt(1e18)

// NOTE: Inexact divisions will have this degree of precision
const Precision int32 = 18

// CalculateFee outputs a fee in wei according to the formula: baseUSDFee / tokenPriceInUSD
func CalculateFee(tokenPriceInUSD decimal.Decimal, baseUSDFee decimal.Decimal) *big.Int {
	if tokenPriceInUSD.IsZero() || baseUSDFee.IsZero() {
		// zero fee if token price or base fee is zero
		return big.NewInt(0)
	}

	// fee denominated in token
	fee := baseUSDFee.DivRound(tokenPriceInUSD, Precision)

	// fee scaled up
	fee = fee.Mul(FeeScalingFactor)

	return fee.BigInt()
}
