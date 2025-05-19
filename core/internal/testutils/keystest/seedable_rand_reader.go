package keystest

import (
	"io"
	"math/rand"
)

// NewRandReaderFromSeed returns a seedable random io reader, producing deterministic
// output. This is useful for deterministically producing keys for tests. This is an
// insecure source of randomness and therefor should only be used in tests.
func NewRandReaderFromSeed(seed int64) io.Reader {
	return rand.New(rand.NewSource(seed))
}
