// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import "math/big"

// IsForked returns whether a fork scheduled at block s is active at the given head block.
// Note: [s] and [head] can be either a block number or a block timestamp.
func IsForked(s, head *big.Int) bool {
	if s == nil || head == nil {
		return false
	}
	return s.Cmp(head) <= 0
}

// IsForkTransition returns true if [fork] activates during the transition from [parent]
// to [current].
// Note: this works for both block number and timestamp activated forks.
func IsForkTransition(fork *big.Int, parent *big.Int, current *big.Int) bool {
	parentForked := IsForked(fork, parent)
	currentForked := IsForked(fork, current)
	return !parentForked && currentForked
}
