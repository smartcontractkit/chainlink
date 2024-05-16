package slicelib

import (
	"sort"

	"github.com/smartcontractkit/ccipocr3/internal/model"
)

// BigIntSortedMiddle returns the middle number after sorting the provided numbers. nil is returned if the provided slice is empty.
// If length of the provided slice is even, the right-hand-side value of the middle 2 numbers is returned.
// The objective of this function is to always pick within the range of values reported by honest nodes when we have 2f+1 values.
func BigIntSortedMiddle(vals []model.BigInt) model.BigInt {
	if len(vals) == 0 {
		return model.BigInt{}
	}

	valsCopy := make([]model.BigInt, len(vals))
	copy(valsCopy[:], vals[:])

	sort.Slice(valsCopy, func(i, j int) bool {
		return (valsCopy[i].Int).Cmp(valsCopy[j].Int) < 0
	})
	return valsCopy[len(valsCopy)/2]
}
