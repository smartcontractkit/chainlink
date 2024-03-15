// Package secp256k1 is an implementation of the kyber.{Group,Point,Scalar}
// //////////////////////////////////////////////////////////////////////////////
//
//	XXX: Do not use in production until this code has been audited.
//
// //////////////////////////////////////////////////////////////////////////////
// interfaces, based on btcd/btcec and kyber/group/mod
//
// XXX: NOT CONSTANT TIME!
package secp256k1

// Arithmetic operations in the base field of secp256k1, i.e. ℤ/qℤ, where q is
// the base field characteristic.

import (
	"crypto/cipher"
	"fmt"
	"math/big"

	"go.dedis.ch/kyber/v3/util/random"
)

// q is the field characteristic (cardinality) of the secp256k1 base field. All
// arithmetic operations on the field are modulo this.
var q = s256.P

type fieldElt big.Int

// newFieldZero returns a newly allocated field element.
func newFieldZero() *fieldElt { return (*fieldElt)(big.NewInt(0)) }

// Int returns f as a big.Int
func (f *fieldElt) int() *big.Int { return (*big.Int)(f) }

// modQ reduces f's underlying big.Int modulo q, and returns it
func (f *fieldElt) modQ() *fieldElt {
	if f.int().Cmp(q) != -1 || f.int().Cmp(bigZero) == -1 {
		// f ∉ {0, ..., q-1}. Find the representative of f+qℤ in that set.
		//
		// Per Mod docstring, "Mod implements Euclidean modulus", meaning that after
		// this, f will be the smallest non-negative representative of its
		// equivalence class in ℤ/qℤ. TODO(alx): Make this faster
		f.int().Mod(f.int(), q)
	}
	return f
}

// This differs from SetInt below, in that it does not take a copy of v.
func fieldEltFromBigInt(v *big.Int) *fieldElt { return (*fieldElt)(v).modQ() }

func fieldEltFromInt(v int64) *fieldElt {
	return fieldEltFromBigInt(big.NewInt(v)).modQ()
}

var fieldZero = fieldEltFromInt(0)
var bigZero = big.NewInt(0)

// String returns the string representation of f
func (f *fieldElt) String() string {
	return fmt.Sprintf("fieldElt{%x}", f.int())
}

// Equal returns true iff f=g, i.e. the backing big.Ints satisfy f ≡ g mod q
func (f *fieldElt) Equal(g *fieldElt) bool {
	if f == (*fieldElt)(nil) && g == (*fieldElt)(nil) {
		return true
	}
	if f == (*fieldElt)(nil) { // f is nil, g is not
		return false
	}
	if g == (*fieldElt)(nil) { // g is nil, f is not
		return false
	}
	return bigZero.Cmp(newFieldZero().Sub(f, g).modQ().int()) == 0
}

// Add sets f to the sum of a and b modulo q, and returns it.
func (f *fieldElt) Add(a, b *fieldElt) *fieldElt {
	f.int().Add(a.int(), b.int())
	return f.modQ()
}

// Sub sets f to a-b mod q, and returns it.
func (f *fieldElt) Sub(a, b *fieldElt) *fieldElt {
	f.int().Sub(a.int(), b.int())
	return f.modQ()
}

// Set sets f's value to v, and returns f.
func (f *fieldElt) Set(v *fieldElt) *fieldElt {
	f.int().Set(v.int())
	return f.modQ()
}

// SetInt sets f's value to v mod q, and returns f.
func (f *fieldElt) SetInt(v *big.Int) *fieldElt {
	f.int().Set(v)
	return f.modQ()
}

// Pick samples uniformly from {0, ..., q-1}, assigns sample to f, and returns f
func (f *fieldElt) Pick(rand cipher.Stream) *fieldElt {
	return f.SetInt(random.Int(q, rand)) // random.Int safe because q≅2²⁵⁶, q<2²⁵⁶
}

// Neg sets f to the negation of g modulo q, and returns it
func (f *fieldElt) Neg(g *fieldElt) *fieldElt {
	f.int().Neg(g.int())
	return f.modQ()
}

// Clone returns a new fieldElt, backed by a clone of f
func (f *fieldElt) Clone() *fieldElt { return newFieldZero().Set(f.modQ()) }

// SetBytes sets f to the 32-byte big-endian value represented by buf, reduces
// it, and returns it.
func (f *fieldElt) SetBytes(buf [32]byte) *fieldElt {
	f.int().SetBytes(buf[:])
	return f.modQ()
}

// Bytes returns the 32-byte big-endian representation of f
func (f *fieldElt) Bytes() [32]byte {
	bytes := f.modQ().int().Bytes()
	if len(bytes) > 32 {
		panic("field element longer than 256 bits")
	}
	var rv [32]byte
	copy(rv[32-len(bytes):], bytes) // leftpad w zeros
	return rv
}

var two = big.NewInt(2)

// square returns y² mod q
func fieldSquare(y *fieldElt) *fieldElt {
	return fieldEltFromBigInt(newFieldZero().int().Exp(y.int(), two, q))
}

func i() *big.Int { return new(big.Int) }

// sqrtPower is s.t. n^sqrtPower≡sqrt(n) mod q, if n has a root at all. See
// https://math.stackexchange.com/a/1816280, for instance
var sqrtPower = i().Rsh(i().Add(q, big.NewInt(1)), 2) // (q +1)/4

// maybeSqrtInField returns a square root of v, if it has any, else nil
func maybeSqrtInField(v *fieldElt) *fieldElt {
	s := newFieldZero()
	s.int().Exp(v.int(), sqrtPower, q)
	if !fieldSquare(s).Equal(v) {
		return nil
	}
	return s
}

var three = big.NewInt(3)
var seven = fieldEltFromInt(7)

// rightHandSide returns the RHS of the secp256k1 equation, x³+7 mod q, given x
func rightHandSide(x *fieldElt) *fieldElt {
	xCubed := newFieldZero()
	xCubed.int().Exp(x.int(), three, q)
	return xCubed.Add(xCubed, seven)
}

// isEven returns true if f is even, false otherwise
func (f *fieldElt) isEven() bool {
	return big.NewInt(0).Mod(f.int(), two).Cmp(big.NewInt(0)) == 0
}
