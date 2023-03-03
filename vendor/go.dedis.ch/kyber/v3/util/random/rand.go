// Package random provides facilities for generating
// random or pseudorandom cryptographic objects.
package random

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"math/big"

	"go.dedis.ch/kyber/v3/xof/blake2xb"
)

// Bits chooses a uniform random BigInt with a given maximum BitLen.
// If 'exact' is true, choose a BigInt with _exactly_ that BitLen, not less
func Bits(bitlen uint, exact bool, rand cipher.Stream) []byte {
	b := make([]byte, (bitlen+7)/8)
	rand.XORKeyStream(b, b)
	highbits := bitlen & 7
	if highbits != 0 {
		b[0] &= ^(0xff << highbits)
	}
	if exact {
		if highbits != 0 {
			b[0] |= 1 << (highbits - 1)
		} else {
			b[0] |= 0x80
		}
	}
	return b
}

// Int chooses a uniform random big.Int less than a given modulus
func Int(mod *big.Int, rand cipher.Stream) *big.Int {
	bitlen := uint(mod.BitLen())
	i := new(big.Int)
	for {
		i.SetBytes(Bits(bitlen, false, rand))
		if i.Sign() > 0 && i.Cmp(mod) < 0 {
			return i
		}
	}
}

// Bytes fills a slice with random bytes from rand.
func Bytes(b []byte, rand cipher.Stream) {
	rand.XORKeyStream(b, b)
}

type randstream struct {
	Readers []io.Reader
}

func (r *randstream) XORKeyStream(dst, src []byte) {

	l := len(dst)
	if len(src) != l {
		panic("XORKeyStream: mismatched buffer lengths")
	}

	// readerBytes is how many bytes we expect from each source
	readerBytes := 32

	// try to read readerBytes bytes from all readers and write them in a buffer
	var b bytes.Buffer
	var nerr int
	buff := make([]byte, readerBytes)
	for _, reader := range r.Readers {
		n, err := io.ReadFull(reader, buff)
		if err != nil {
			nerr++
		}
		b.Write(buff[:n])
	}

	// we are ok with few sources being insecure (i.e., providing less than
	// readerBytes bytes), but not all of them
	if nerr == len(r.Readers) {
		panic("all readers failed")
	}

	// create the XOF output, with hash of collected data as seed
	h := sha256.New()
	h.Write(b.Bytes())
	seed := h.Sum(nil)
	blake2 := blake2xb.New(seed)
	blake2.XORKeyStream(dst, src)
}

// New returns a new cipher.Stream that gets random data from the given
// readers. If no reader was provided, Go's crypto/rand package is used.
// Otherwise, for each source, 32 bytes are read. They are concatenated and
// then hashed, and the resulting hash is used as a seed to a PRNG.
// The resulting cipher.Stream can be used in multiple threads.
func New(readers ...io.Reader) cipher.Stream {
	if len(readers) == 0 {
		readers = []io.Reader{rand.Reader}
	}
	return &randstream{readers}
}
