package kyber

import (
	"crypto/cipher"
	"io"
)

// An XOF is an extendable output function, which is a cryptographic
// primitive that can take arbitrary input in the same way a hash
// function does, and then create a stream of output, up to a limit
// determined by the size of the internal state of the hash function
// the underlies the XOF.
//
// When XORKeyStream is called with zeros for the source, an XOF
// also acts as a PRNG. If it is seeded with an appropriate amount
// of keying material, it is a cryptographically secure source of random
// bits.
type XOF interface {
	// Write absorbs more data into the hash's state. It panics if called
	// after Read. Use Reseed() to reset the XOF into a state where more data
	// can be absorbed via Write.
	io.Writer

	// Read reads more output from the hash. It returns io.EOF if the
	// limit of available data for reading has been reached.
	io.Reader

	// An XOF implements cipher.Stream, so that callers can use XORKeyStream
	// to encrypt/decrypt data. The key stream is read from the XOF using
	// the io.Reader interface. If Read returns an error, then XORKeyStream
	// will panic.
	cipher.Stream

	// Reseed makes an XOF writeable again after it has been read from
	// by sampling a key from it's output and initializing a fresh XOF implementation
	// with that key.
	Reseed()

	// Clone returns a copy of the XOF in its current state.
	Clone() XOF
}

// An XOFFactory is an interface that can be mixed in to local suite definitions.
type XOFFactory interface {
	// XOF creates a new XOF, feeding seed to it via it's Write method. If seed
	// is nil or []byte{}, the XOF is left unseeded, it will produce a fixed, predictable
	// stream of bits (Caution: this behavior is useful for testing but fatal for
	// production use).
	XOF(seed []byte) XOF
}
