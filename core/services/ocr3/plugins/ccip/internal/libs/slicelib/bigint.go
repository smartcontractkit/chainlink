package slicelib

import (
	"sort"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// BigIntSortedMiddle returns the middle number after sorting the provided numbers.
// nil is returned if the provided slice is empty.
// If length of the provided slice is even, the right-hand-side value of the middle 2 numbers is returned.
// The objective of this function is to always pick within the range of values reported by honest nodes
// when we have 2f+1 values.
func BigIntSortedMiddle(vals []cciptypes.BigInt) cciptypes.BigInt {
	if len(vals) == 0 {
		return cciptypes.BigInt{}
	}

	valsCopy := make([]cciptypes.BigInt, len(vals))
	copy(valsCopy[:], vals[:])

	sort.Slice(valsCopy, func(i, j int) bool {
		return (valsCopy[i].Int).Cmp(valsCopy[j].Int) < 0
	})
	return valsCopy[len(valsCopy)/2]
}
