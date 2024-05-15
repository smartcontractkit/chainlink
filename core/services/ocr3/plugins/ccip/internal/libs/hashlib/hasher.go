package hashlib

import (
	"bytes"

	"golang.org/x/crypto/sha3"
)

func Keccak256Fixed(in []byte) [32]byte {
	hash := sha3.NewLegacyKeccak256()
	// Note this Keccak256 cannot error https://github.com/golang/crypto/blob/master/sha3/sha3.go#L126
	// if we start supporting hashing algos which do, we can change this API to include an error.
	hash.Write(in)
	var h [32]byte
	copy(h[:], hash.Sum(nil))
	return h
}

// Hash contains all supported hash formats.
// Add additional hash types e.g. [20]byte as needed here.
type Hash interface {
	[32]byte
}

type Ctx[H Hash] interface {
	Hash(l []byte) H
	HashInternal(a, b H) H
	ZeroHash() H
}

type keccakCtx struct {
	InternalDomainSeparator [32]byte
}

func NewKeccakCtx() Ctx[[32]byte] {
	return keccakCtx{
		InternalDomainSeparator: [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	}
}

// Hash hashes a byte array with Keccak256
func (k keccakCtx) Hash(l []byte) [32]byte {
	// Note this Keccak256 cannot error https://github.com/golang/crypto/blob/master/sha3/sha3.go#L126
	// if we start supporting hashing algos which do, we can change this API to include an error.
	return Keccak256Fixed(l)
}

// HashInternal orders two [32]byte values and prepends them with
// a separator before hashing them.
func (k keccakCtx) HashInternal(a, b [32]byte) [32]byte {
	if bytes.Compare(a[:], b[:]) < 0 {
		return k.Hash(append(k.InternalDomainSeparator[:], append(a[:], b[:]...)...))
	}
	return k.Hash(append(k.InternalDomainSeparator[:], append(b[:], a[:]...)...))
}

// ZeroHash returns the zero hash: 0xFF..FF
// We use bytes32 0xFF..FF for zeroHash in the CCIP research spec, this needs to match.
// This value is chosen since it is unlikely to be the result of a hash, and cannot match any internal node preimage.
func (k keccakCtx) ZeroHash() [32]byte {
	var zeroes [32]byte
	for i := 0; i < 32; i++ {
		zeroes[i] = 0xFF
	}
	return zeroes
}
