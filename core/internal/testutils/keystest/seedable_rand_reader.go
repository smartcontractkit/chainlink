package keystest

import (
	"io"
	"math/rand"
)

type randReader struct{}

func (randReader) Read(b []byte) (n int, err error) {
	return rand.Read(b)
}

// NewRandReaderFromSeed returns a seedable random io reader, producing deterministic
// output. This is useful for deterministically producing keys for tests. This is an
// insecure source of randomness and therefor should only be used in tests.
func NewRandReaderFromSeed(seed int64) io.Reader {
	rand.Seed(seed)
	return randReader{}
}
