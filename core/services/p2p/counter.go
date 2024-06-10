package p2p

import (
	"encoding/binary"
)

// counter is a simple abstraction that can be used to generate unique peer group IDs.
type counter struct {
	x uint64
}

// Bytes returns the counter as a 32-byte array.
func (g *counter) Bytes() [32]byte {
	var b [32]byte
	binary.BigEndian.PutUint64(b[24:], g.x)
	return b
}

// Inc increments the counter.
func (g *counter) Inc() *counter {
	g.x++
	return g
}
