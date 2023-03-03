// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"crypto/rand"
)

// CopyBytes returns a copy of the provided byte slice. If nil is provided, nil
// will be returned.
func CopyBytes(b []byte) []byte {
	if b == nil {
		return nil
	}

	cb := make([]byte, len(b))
	copy(cb, b)
	return cb
}

// RandomBytes returns a slice of n random bytes
// Intended for use in testing
func RandomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}
