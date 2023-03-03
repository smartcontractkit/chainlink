// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package hashing

// Hasher is an interface to compute a hash value.
type Hasher interface {
	// Hash takes a string and computes its hash value.
	// Values must be computed deterministically.
	Hash([]byte) uint64
}
