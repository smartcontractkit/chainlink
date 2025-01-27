package ccipcalc

import (
	"math"
	"math/big"
	"sort"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	// CurveBasedDeviationPPB is the deviation threshold when writing to ethereum's PriceRegistry and is the trigger for
	// using curve-based deviation logic.
	CurveBasedDeviationPPB = 4e9
)

// ContiguousReqs checks if seqNrs contains all numbers from min to max.
func ContiguousReqs(lggr logger.Logger, min, max uint64, seqNrs []uint64) bool {
	if int(max-min+1) != len(seqNrs) {
		return false
	}

	for i, j := min, 0; i <= max && j < len(seqNrs); i, j = i+1, j+1 {
		if seqNrs[j] != i {
			lggr.Errorw("unexpected gap in seq nums", "seqNr", i, "minSeqNr", min, "maxSeqNr", max)
			return false
		}
	}
	return true
}

// CalculateUsdPerUnitGas returns: (sourceGasPrice * usdPerFeeCoin) / 1e18
func CalculateUsdPerUnitGas(sourceGasPrice *big.Int, usdPerFeeCoin *big.Int) *big.Int {
	// (wei / gas) * (usd / eth) * (1 eth / 1e18 wei)  = usd/gas
	tmp := new(big.Int).Mul(sourceGasPrice, usdPerFeeCoin)
	return tmp.Div(tmp, big.NewInt(1e18))
}

// BigIntSortedMiddle returns the middle number after sorting the provided numbers. nil is returned if the provided slice is empty.
// If length of the provided slice is even, the right-hand-side value of the middle 2 numbers is returned.
// The objective of this function is to always pick within the range of values reported by honest nodes when we have 2f+1 values.
func BigIntSortedMiddle(vals []*big.Int) *big.Int {
	if len(vals) == 0 {
		return nil
	}

	valsCopy := make([]*big.Int, len(vals))
	copy(valsCopy, vals)
	sort.Slice(valsCopy, func(i, j int) bool {
		return valsCopy[i].Cmp(valsCopy[j]) == -1
	})
	return valsCopy[len(valsCopy)/2]
}

// Deviates checks if x1 and x2 deviates based on the provided ppb (parts per billion)
// ppb is calculated based on the smaller value of the two
// e.g, if x1 > x2, deviation_parts_per_billion = ((x1 - x2) / x2) * 1e9
func Deviates(x1, x2 *big.Int, ppb int64) bool {
	// if x1 == 0 or x2 == 0, deviates if x2 != x1, to avoid the relative division by 0 error
	if x1.BitLen() == 0 || x2.BitLen() == 0 {
		return x1.Cmp(x2) != 0
	}
	diff := big.NewInt(0).Sub(x1, x2) // diff = x1-x2
	diff.Mul(diff, big.NewInt(1e9))   // diff = diff * 1e9
	// dividing by the smaller value gives consistent ppb regardless of input order, and supports >100% deviation.
	if x1.Cmp(x2) > 0 {
		diff.Div(diff, x2)
	} else {
		diff.Div(diff, x1)
	}
	return diff.CmpAbs(big.NewInt(ppb)) > 0 // abs(diff) > ppb
}

// DeviatesOnCurve calculates a deviation threshold on the fly using xNew. For now it's only used for gas price
// deviation calculation. It's important to make sure the order of xNew and xOld is correct when passed into this
// function to get an accurate deviation threshold.
func DeviatesOnCurve(xNew, xOld, noDeviationLowerBound *big.Int, ppb int64) bool {
	// This is a temporary gating mechanism that ensures we only apply the gas curve deviation logic to eth-bound price
	// updates. If ppb from config is not equal to 4000000000, do not apply the gas curve.
	if ppb != CurveBasedDeviationPPB {
		return Deviates(xOld, xNew, ppb)
	}

	// If xNew < noDeviationLowerBound, Deviates should never be true
	if xNew.Cmp(noDeviationLowerBound) < 0 {
		return false
	}

	xNewFloat := new(big.Float).SetInt(xNew)
	xNewFloat64, _ := xNewFloat.Float64()

	// We use xNew to generate the threshold so that when going from cheap --> expensive, xNew generates a smaller
	// deviation threshold so we are more likely to update the gas price on chain. When going from expensive --> cheap,
	// xNew generates a larger deviation threshold since it's not as urgent to update the gas price on chain.
	curveThresholdPPB := calculateCurveThresholdPPB(xNewFloat64)
	return Deviates(xNew, xOld, curveThresholdPPB)
}

// calculateCurveThresholdPPB calculates the deviation threshold percentage with x using the formula:
// y = (10e11) / (x^0.665). This sliding scale curve was created by collecting several thousands of historical
// PriceRegistry gas price update samples from chains with volatile gas prices like Zircuit and Mode and then using that
// historical data to define thresholds of gas deviations that were acceptable given their USD value. The gas prices
// (X coordinates) and these new thresholds (Y coordinates) were used to fit this sliding scale curve that returns large
// thresholds at low gas prices and smaller thresholds at higher gas prices. Constructing the curve in USD terms allows
// us to more easily reason about these thresholds and also better translates to non-evm chains. For example, when the
// per unit gas price is 0.000006 USD, (6e12 USDgwei), the curve will output a threshold of around 3,000%. However, when
// the per unit gas price is 0.005 USD (5e15 USDgwei), the curve will output a threshold of only ~30%.
func calculateCurveThresholdPPB(x float64) int64 {
	const constantFactor = 10e11
	const exponent = 0.665
	xPower := math.Pow(x, exponent)
	threshold := constantFactor / xPower

	// Convert curve output percentage to PPB
	thresholdPPB := int64(threshold * 1e7)
	return thresholdPPB
}

func MergeEpochAndRound(epoch uint32, round uint8) uint64 {
	return uint64(epoch)<<8 + uint64(round)
}
