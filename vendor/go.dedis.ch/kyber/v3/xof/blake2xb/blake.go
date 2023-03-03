// Package blake2xb provides an implementation of kyber.XOF based on the
// Blake2xb construction.
package blake2xb

import (
	"go.dedis.ch/kyber/v3"
	"golang.org/x/crypto/blake2b"
)

type xof struct {
	impl blake2b.XOF
	// key is here to not make excess garbage during repeated calls
	// to XORKeyStream.
	key []byte
}

// New creates a new XOF using the Blake2b hash.
func New(seed []byte) kyber.XOF {
	seed1 := seed
	var seed2 []byte
	if len(seed) > blake2b.Size {
		seed1 = seed[0:blake2b.Size]
		seed2 = seed[blake2b.Size:]
	}
	b, err := blake2b.NewXOF(blake2b.OutputLengthUnknown, seed1)
	if err != nil {
		panic("blake2b.NewXOF should not return error: " + err.Error())
	}

	if seed2 != nil {
		_, err := b.Write(seed2)
		if err != nil {
			panic("blake2b.XOF.Write should not return error: " + err.Error())
		}
	}
	return &xof{impl: b}
}

func (x *xof) Clone() kyber.XOF {
	return &xof{impl: x.impl.Clone()}
}

func (x *xof) Read(dst []byte) (int, error) {
	return x.impl.Read(dst)
}

func (x *xof) Write(src []byte) (int, error) {
	return x.impl.Write(src)
}

func (x *xof) Reseed() {
	// Use New to create a new one seeded with output from the old one.
	if len(x.key) < 128 {
		x.key = make([]byte, 128)
	} else {
		x.key = x.key[0:128]
	}
	x.Read(x.key)
	y := New(x.key)
	// Steal the XOF implementation, and put it inside of x.
	x.impl = y.(*xof).impl
}

func (x *xof) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic("dst too short")
	}
	if len(x.key) < len(src) {
		x.key = make([]byte, len(src))
	} else {
		x.key = x.key[0:len(src)]
	}

	n, err := x.Read(x.key)
	if err != nil {
		panic("blake xof error: " + err.Error())
	}
	if n != len(src) {
		panic("short read on key")
	}

	for i := range src {
		dst[i] = src[i] ^ x.key[i]
	}
}
