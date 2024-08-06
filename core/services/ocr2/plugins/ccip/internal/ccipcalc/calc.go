package ccipcalc

import (
	"math/big"
	"sort"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
	copy(valsCopy[:], vals[:])
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

func MergeEpochAndRound(epoch uint32, round uint8) uint64 {
	return uint64(epoch)<<8 + uint64(round)
}
