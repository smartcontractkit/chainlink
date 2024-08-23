package fee

import (
	"errors"
	"math/big"

	bigmath "github.com/smartcontractkit/chainlink-common/pkg/utils/big_math"
)

var (
	ErrBumpFeeExceedsLimit = errors.New("fee bump exceeds limit")
	ErrBump                = errors.New("fee bump failed")
	ErrConnectivity        = errors.New("transaction propagation issue: transactions are not being mined")
)

func IsBumpErr(err error) bool {
	return err != nil && (errors.Is(err, ErrBumpFeeExceedsLimit) || errors.Is(err, ErrBump) || errors.Is(err, ErrConnectivity))
}

// CalculateFee computes the fee price for a transaction.
// The fee price is the minimum of:
// - max fee price specified, default fee price and max fee price for the node.
func CalculateFee(
	maxFeePrice, defaultPrice, maxFeePriceConfigured *big.Int,
) *big.Int {
	maxFeePriceAllowed := bigmath.Min(maxFeePrice, maxFeePriceConfigured)
	return bigmath.Min(defaultPrice, maxFeePriceAllowed)
}

// Returns highest bumped fee price of originalFeePrice bumped by fixed units or percentage.
func MaxBumpedFee(originalFeePrice *big.Int, feeBumpPercent uint16, feeBumpUnits *big.Int) *big.Int {
	return bigmath.Max(
		AddPercentage(originalFeePrice, feeBumpPercent),
		new(big.Int).Add(originalFeePrice, feeBumpUnits),
	)
}
