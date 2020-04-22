// Package cryptotest provides convenience functions for kyber-based APIs.
//
// It is separate from cltest to prevent an import cycle.
package cryptotest

import (
	"math/rand"
	"testing"
)

// RandomStream implements cipher.Stream, but with a deterministic output.
type RandomStream rand.Rand

// NewStream returns a randomStream seeded from seed, for deterministic
// randomness in tests of random outputs, and for small property-based tests.
//
// This API is deliberately awkward to prevent it from being used outside of
// tests.
//
// The testing framework runs the tests in a file in series, unless you
// explicitly request otherwise with testing.T.Parallel(). So one such stream
// per file is enough, most of the time.
func NewStream(t *testing.T, seed int64) *RandomStream {
	return (*RandomStream)(rand.New(rand.NewSource(seed)))
}

// XORKeyStream dumps the output from a math/rand PRNG on dst.
//
// It gives no consideration for the contents of src, and is named so
// misleadingly purely to satisfy the cipher.Stream interface.
func (s *RandomStream) XORKeyStream(dst, src []byte) {
	(*rand.Rand)(s).Read(dst)
}
