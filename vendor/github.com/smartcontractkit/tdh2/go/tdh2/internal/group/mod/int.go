// Package mod contains a generic implementation of finite field arithmetic
// on integer fields with a constant modulus.
package mod

import (
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/smartcontractkit/tdh2/go/tdh2/internal/group"
)

// Int is a generic implementation of finite field arithmetic
// on integer finite fields with a given constant modulus,
// built using Go's built-in big.Int package.
// Int satisfies the group.Scalar interface,
// and hence serves as a basic implementation of group.Scalar,
// e.g., representing discrete-log exponents of Schnorr groups
// or scalar multipliers for elliptic curves.
//
// Int offers an API similar to and compatible with big.Int,
// but "carries around" a pointer to the relevant modulus
// and automatically normalizes the value to that modulus
// after all arithmetic operations, simplifying modular arithmetic.
// Binary operations assume that the source(s)
// have the same modulus, but do not check this assumption.
// Unary and binary arithmetic operations may be performed on uninitialized
// target objects, and receive the modulus of the first operand.
// For efficiency the modulus field M is a pointer,
// whose target is assumed never to change.
type Int struct {
	V big.Int  // Integer value from 0 through M-1
	M *big.Int // Modulus for finite field arithmetic
}

// NewInt64 creates a new Int with a given int64 value and big.Int modulus.
func NewInt64(v int64, M *big.Int) *Int {
	i := &Int{M: M}
	i.V.SetInt64(v).Mod(&i.V, M)
	return i
}

// Return the Int's integer value in hexadecimal string representation.
func (i *Int) String() string {
	return hex.EncodeToString(i.V.Bytes())
}

// Equal returns true if the two Ints are equal
func (i *Int) Equal(s2 group.Scalar) bool {
	return i.V.Cmp(&s2.(*Int).V) == 0
}

// Clone returns a separate duplicate of this Int.
func (i *Int) Clone() group.Scalar {
	ni := &Int{M: i.M}
	ni.V.Set(&i.V).Mod(&i.V, i.M)
	return ni
}

// Zero set the Int to the value 0.  The modulus must already be initialized.
func (i *Int) Zero() group.Scalar {
	i.V.SetInt64(0)
	return i
}

// One sets the Int to the value 1.  The modulus must already be initialized.
func (i *Int) One() group.Scalar {
	i.V.SetInt64(1)
	return i
}

// SetInt64 sets the Int to an arbitrary 64-bit "small integer" value.
// The modulus must already be initialized.
func (i *Int) SetInt64(v int64) group.Scalar {
	i.V.SetInt64(v).Mod(&i.V, i.M)
	return i
}

// Add sets the target to a + b mod M, where M is a's modulus..
func (i *Int) Add(a, b group.Scalar) group.Scalar {
	ai := a.(*Int)
	bi := b.(*Int)
	i.M = ai.M
	i.V.Add(&ai.V, &bi.V).Mod(&i.V, i.M)
	return i
}

// Sub sets the target to a - b mod M.
// Target receives a's modulus.
func (i *Int) Sub(a, b group.Scalar) group.Scalar {
	ai := a.(*Int)
	bi := b.(*Int)
	i.M = ai.M
	i.V.Sub(&ai.V, &bi.V).Mod(&i.V, i.M)
	return i
}

// Neg sets the target to -a mod M.
func (i *Int) Neg(a group.Scalar) group.Scalar {
	ai := a.(*Int)
	i.M = ai.M
	if ai.V.Sign() > 0 {
		i.V.Sub(i.M, &ai.V)
	} else {
		i.V.SetUint64(0)
	}
	return i
}

// Mul sets the target to a * b mod M.
// Target receives a's modulus.
func (i *Int) Mul(a, b group.Scalar) group.Scalar {
	ai := a.(*Int)
	bi := b.(*Int)
	i.M = ai.M
	i.V.Mul(&ai.V, &bi.V).Mod(&i.V, i.M)
	return i
}

// Div sets the target to a * b^-1 mod M, where b^-1 is the modular inverse of b.
func (i *Int) Div(a, b group.Scalar) group.Scalar {
	ai := a.(*Int)
	bi := b.(*Int)
	var t big.Int
	i.M = ai.M
	i.V.Mul(&ai.V, t.ModInverse(&bi.V, i.M))
	i.V.Mod(&i.V, i.M)
	return i
}

// Inv sets the target to the modular inverse of a with respect to modulus M.
func (i *Int) Inv(a group.Scalar) group.Scalar {
	ai := a.(*Int)
	i.M = ai.M
	i.V.ModInverse(&a.(*Int).V, i.M)
	return i
}

// Pick a [pseudo-]random integer modulo M
// using bits from the given stream cipher.
// This code is adopted from Go's elliptic.GenerateKey()
// and the rejection sampling can lead to up to a two-fold
// slowdown, if M is not close to 2**bitSize.
func (i *Int) Pick(rand cipher.Stream) group.Scalar {
	var n *big.Int
	// This is just a bitmask with the number of ones starting at 8 then
	// incrementing by index. To account for fields with bitsizes that are not a whole
	// number of bytes, we mask off the unnecessary bits. h/t agl
	mask := []byte{0xff, 0x1, 0x3, 0x7, 0xf, 0x1f, 0x3f, 0x7f}
	bitSize := i.M.BitLen()
	byteLen := (bitSize + 7) / 8
	b := make([]byte, byteLen)

	for {
		rand.XORKeyStream(b, b)
		// We have to mask off any excess bits in the case that the size of the
		// underlying field is not a whole number of bytes.
		b[0] &= mask[bitSize%8]
		// This is because, in tests, rand will return all zeros and we don't
		// want to get the point at infinity and loop forever.
		b[1] ^= 0x42

		n = new(big.Int).SetBytes(b)
		// If the scalar is out of range, sample another random number.
		if n.Cmp(i.M) < 0 {
			break
		}
	}

	i.V.Set(n)
	return i
}

// MarshalSize returns the length in bytes of encoded integers with modulus M.
// The length of encoded Ints depends only on the size of the modulus,
// and not on the the value of the encoded integer,
// making the encoding is fixed-length for simplicity and security.
func (i *Int) MarshalSize() int {
	return (i.M.BitLen() + 7) / 8
}

// MarshalBinary encodes the value of this Int into a byte-slice exactly Len() bytes long.
// It uses big endian.
func (i *Int) MarshalBinary() ([]byte, error) {
	l := i.MarshalSize()
	b := i.V.Bytes() // may be shorter than l
	offset := l - len(b)

	if offset != 0 {
		nb := make([]byte, l)
		copy(nb[offset:], b)
		b = nb
	}
	return b, nil
}

// UnmarshalBinary tries to decode a Int from a byte-slice buffer.
// Returns an error if the buffer is not exactly Len() bytes long
// or if the contents of the buffer represents an out-of-range integer.
func (i *Int) UnmarshalBinary(buf []byte) error {
	if len(buf) != i.MarshalSize() {
		return errors.New("UnmarshalBinary: wrong size buffer")
	}

	i.V.SetBytes(buf)
	if i.V.Cmp(i.M) >= 0 {
		return errors.New("UnmarshalBinary: value out of range")
	}
	return nil
}

// SetBytes set the value value to a number represented
// by a byte string.
func (i *Int) SetBytes(a []byte) group.Scalar {
	var buff = a
	i.V.SetBytes(buff).Mod(&i.V, i.M)
	return i
}
