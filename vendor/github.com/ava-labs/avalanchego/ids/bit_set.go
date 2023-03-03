// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ids

import (
	"fmt"
	"math/bits"
)

// BitSet64 is a set that can contain uints in the range [0, 64). All functions
// are O(1). The zero value is the empty set.
type BitSet64 uint64

// Add [i] to the set of ints
func (bs *BitSet64) Add(i uint) { *bs |= 1 << i }

// Union adds all the elements in [s] to this set
func (bs *BitSet64) Union(s BitSet64) { *bs |= s }

// Intersection takes the intersection of [s] with this set
func (bs *BitSet64) Intersection(s BitSet64) { *bs &= s }

// Difference removes all the elements in [s] from this set
func (bs *BitSet64) Difference(s BitSet64) { *bs &^= s }

// Remove [i] from the set of ints
func (bs *BitSet64) Remove(i uint) { *bs &^= 1 << i }

// Clear removes all elements from this set
func (bs *BitSet64) Clear() { *bs = 0 }

// Contains returns true if [i] was previously added to this set
func (bs BitSet64) Contains(i uint) bool { return bs&(1<<i) != 0 }

// Len returns the number of elements in this set
func (bs BitSet64) Len() int { return bits.OnesCount64(uint64(bs)) }

func (bs BitSet64) String() string { return fmt.Sprintf("%016x", uint64(bs)) }
