// Package mod contains a generic implementation of finite field arithmetic
// on integer fields with a constant modulus.
package mod

import (
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"io"
	"math/big"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/internal/marshalling"
	"go.dedis.ch/kyber/v3/util/random"
)

var one = big.NewInt(1)
var two = big.NewInt(2)
var marshalScalarID = [8]byte{'m', 'o', 'd', '.', 'i', 'n', 't', ' '}

// ByteOrder denotes the endianness of the operation.
type ByteOrder bool

const (
	// LittleEndian endianness
	LittleEndian ByteOrder = true
	// BigEndian endianness
	BigEndian ByteOrder = false
)

// Int is a generic implementation of finite field arithmetic
// on integer finite fields with a given constant modulus,
// built using Go's built-in big.Int package.
// Int satisfies the kyber.Scalar interface,
// and hence serves as a basic implementation of kyber.Scalar,
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
	V  big.Int   // Integer value from 0 through M-1
	M  *big.Int  // Modulus for finite field arithmetic
	BO ByteOrder // Endianness which will be used on input and output
}

// NewInt creaters a new Int with a given big.Int and a big.Int modulus.
func NewInt(v *big.Int, m *big.Int) *Int {
	return new(Int).Init(v, m)
}

// NewInt64 creates a new Int with a given int64 value and big.Int modulus.
func NewInt64(v int64, M *big.Int) *Int {
	return new(Int).Init64(v, M)
}

// NewIntBytes creates a new Int with a given slice of bytes and a big.Int
// modulus.
func NewIntBytes(a []byte, m *big.Int, byteOrder ByteOrder) *Int {
	return new(Int).InitBytes(a, m, byteOrder)
}

// NewIntString creates a new Int with a given string and a big.Int modulus.
// The value is set to a rational fraction n/d in a given base.
func NewIntString(n, d string, base int, m *big.Int) *Int {
	return new(Int).InitString(n, d, base, m)
}

// Init a Int with a given big.Int value and modulus pointer.
// Note that the value is copied; the modulus is not.
func (i *Int) Init(V *big.Int, m *big.Int) *Int {
	i.M = m
	i.BO = BigEndian
	i.V.Set(V).Mod(&i.V, m)
	return i
}

// Init64 creates an Int with an int64 value and big.Int modulus.
func (i *Int) Init64(v int64, m *big.Int) *Int {
	i.M = m
	i.BO = BigEndian
	i.V.SetInt64(v).Mod(&i.V, m)
	return i
}

// InitBytes init the Int to a number represented in a big-endian byte string.
func (i *Int) InitBytes(a []byte, m *big.Int, byteOrder ByteOrder) *Int {
	i.M = m
	i.BO = byteOrder
	i.SetBytes(a)
	return i
}

// InitString inits the Int to a rational fraction n/d
// specified with a pair of strings in a given base.
func (i *Int) InitString(n, d string, base int, m *big.Int) *Int {
	i.M = m
	i.BO = BigEndian
	if _, succ := i.SetString(n, d, base); !succ {
		panic("InitString: invalid fraction representation")
	}
	return i
}

// Return the Int's integer value in hexadecimal string representation.
func (i *Int) String() string {
	return hex.EncodeToString(i.V.Bytes())
}

// SetString sets the Int to a rational fraction n/d represented by a pair of strings.
// If d == "", then the denominator is taken to be 1.
// Returns (i,true) on success, or
// (nil,false) if either string fails to parse.
func (i *Int) SetString(n, d string, base int) (*Int, bool) {
	if _, succ := i.V.SetString(n, base); !succ {
		return nil, false
	}
	if d != "" {
		var di Int
		di.M = i.M
		if _, succ := di.SetString(d, "", base); !succ {
			return nil, false
		}
		i.Div(i, &di)
	}
	return i, true
}

// Cmp compares two Ints for equality or inequality
func (i *Int) Cmp(s2 kyber.Scalar) int {
	return i.V.Cmp(&s2.(*Int).V)
}

// Equal returns true if the two Ints are equal
func (i *Int) Equal(s2 kyber.Scalar) bool {
	return i.V.Cmp(&s2.(*Int).V) == 0
}

// Nonzero returns true if the integer value is nonzero.
func (i *Int) Nonzero() bool {
	return i.V.Sign() != 0
}

// Set both value and modulus to be equal to another Int.
// Since this method copies the modulus as well,
// it may be used as an alternative to Init().
func (i *Int) Set(a kyber.Scalar) kyber.Scalar {
	ai := a.(*Int)
	i.V.Set(&ai.V)
	i.M = ai.M
	return i
}

// Clone returns a separate duplicate of this Int.
func (i *Int) Clone() kyber.Scalar {
	ni := new(Int).Init(&i.V, i.M)
	ni.BO = i.BO
	return ni
}

// Zero set the Int to the value 0.  The modulus must already be initialized.
func (i *Int) Zero() kyber.Scalar {
	i.V.SetInt64(0)
	return i
}

// One sets the Int to the value 1.  The modulus must already be initialized.
func (i *Int) One() kyber.Scalar {
	i.V.SetInt64(1)
	return i
}

// SetInt64 sets the Int to an arbitrary 64-bit "small integer" value.
// The modulus must already be initialized.
func (i *Int) SetInt64(v int64) kyber.Scalar {
	i.V.SetInt64(v).Mod(&i.V, i.M)
	return i
}

// Int64 returns the int64 representation of the value.
// If the value is not representable in an int64 the result is undefined.
func (i *Int) Int64() int64 {
	return i.V.Int64()
}

// SetUint64 sets the Int to an arbitrary uint64 value.
// The modulus must already be initialized.
func (i *Int) SetUint64(v uint64) kyber.Scalar {
	i.V.SetUint64(v).Mod(&i.V, i.M)
	return i
}

// Uint64 returns the uint64 representation of the value.
// If the value is not representable in an uint64 the result is undefined.
func (i *Int) Uint64() uint64 {
	return i.V.Uint64()
}

// Add sets the target to a + b mod M, where M is a's modulus..
func (i *Int) Add(a, b kyber.Scalar) kyber.Scalar {
	ai := a.(*Int)
	bi := b.(*Int)
	i.M = ai.M
	i.V.Add(&ai.V, &bi.V).Mod(&i.V, i.M)
	return i
}

// Sub sets the target to a - b mod M.
// Target receives a's modulus.
func (i *Int) Sub(a, b kyber.Scalar) kyber.Scalar {
	ai := a.(*Int)
	bi := b.(*Int)
	i.M = ai.M
	i.V.Sub(&ai.V, &bi.V).Mod(&i.V, i.M)
	return i
}

// Neg sets the target to -a mod M.
func (i *Int) Neg(a kyber.Scalar) kyber.Scalar {
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
func (i *Int) Mul(a, b kyber.Scalar) kyber.Scalar {
	ai := a.(*Int)
	bi := b.(*Int)
	i.M = ai.M
	i.V.Mul(&ai.V, &bi.V).Mod(&i.V, i.M)
	return i
}

// Div sets the target to a * b^-1 mod M, where b^-1 is the modular inverse of b.
func (i *Int) Div(a, b kyber.Scalar) kyber.Scalar {
	ai := a.(*Int)
	bi := b.(*Int)
	var t big.Int
	i.M = ai.M
	i.V.Mul(&ai.V, t.ModInverse(&bi.V, i.M))
	i.V.Mod(&i.V, i.M)
	return i
}

// Inv sets the target to the modular inverse of a with respect to modulus M.
func (i *Int) Inv(a kyber.Scalar) kyber.Scalar {
	ai := a.(*Int)
	i.M = ai.M
	i.V.ModInverse(&a.(*Int).V, i.M)
	return i
}

// Exp sets the target to a^e mod M,
// where e is an arbitrary big.Int exponent (not necessarily 0 <= e < M).
func (i *Int) Exp(a kyber.Scalar, e *big.Int) kyber.Scalar {
	ai := a.(*Int)
	i.M = ai.M
	// to protect against golang/go#22830
	var tmp big.Int
	tmp.Exp(&ai.V, e, i.M)
	i.V = tmp
	return i
}

// Jacobi computes the Jacobi symbol of (a/M), which indicates whether a is
// zero (0), a positive square in M (1), or a non-square in M (-1).
func (i *Int) Jacobi(as kyber.Scalar) kyber.Scalar {
	ai := as.(*Int)
	i.M = ai.M
	i.V.SetInt64(int64(big.Jacobi(&ai.V, i.M)))
	return i
}

// Sqrt computes some square root of a mod M of one exists.
// Assumes the modulus M is an odd prime.
// Returns true on success, false if input a is not a square.
func (i *Int) Sqrt(as kyber.Scalar) bool {
	ai := as.(*Int)
	out := i.V.ModSqrt(&ai.V, ai.M)
	i.M = ai.M
	return out != nil
}

// Pick a [pseudo-]random integer modulo M
// using bits from the given stream cipher.
func (i *Int) Pick(rand cipher.Stream) kyber.Scalar {
	i.V.Set(random.Int(i.M, rand))
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
// It uses i's ByteOrder to determine which byte order to output.
func (i *Int) MarshalBinary() ([]byte, error) {
	l := i.MarshalSize()
	b := i.V.Bytes() // may be shorter than l
	offset := l - len(b)

	if i.BO == LittleEndian {
		return i.LittleEndian(l, l), nil
	}

	if offset != 0 {
		nb := make([]byte, l)
		copy(nb[offset:], b)
		b = nb
	}
	return b, nil
}

// MarshalID returns a unique identifier for this type
func (i *Int) MarshalID() [8]byte {
	return marshalScalarID
}

// UnmarshalBinary tries to decode a Int from a byte-slice buffer.
// Returns an error if the buffer is not exactly Len() bytes long
// or if the contents of the buffer represents an out-of-range integer.
func (i *Int) UnmarshalBinary(buf []byte) error {
	if len(buf) != i.MarshalSize() {
		return errors.New("UnmarshalBinary: wrong size buffer")
	}
	// Still needed here because of the comparison with the modulo
	if i.BO == LittleEndian {
		buf = reverse(nil, buf)
	}
	i.V.SetBytes(buf)
	if i.V.Cmp(i.M) >= 0 {
		return errors.New("UnmarshalBinary: value out of range")
	}
	return nil
}

// MarshalTo encodes this Int to the given Writer.
func (i *Int) MarshalTo(w io.Writer) (int, error) {
	return marshalling.ScalarMarshalTo(i, w)
}

// UnmarshalFrom tries to decode an Int from the given Reader.
func (i *Int) UnmarshalFrom(r io.Reader) (int, error) {
	return marshalling.ScalarUnmarshalFrom(i, r)
}

// BigEndian encodes the value of this Int into a big-endian byte-slice
// at least min bytes but no more than max bytes long.
// Panics if max != 0 and the Int cannot be represented in max bytes.
func (i *Int) BigEndian(min, max int) []byte {
	act := i.MarshalSize()
	pad, ofs := act, 0
	if pad < min {
		pad, ofs = min, min-act
	}
	if max != 0 && pad > max {
		panic("Int not representable in max bytes")
	}
	buf := make([]byte, pad)
	copy(buf[ofs:], i.V.Bytes())
	return buf
}

// SetBytes set the value value to a number represented
// by a byte string.
// Endianness depends on the endianess set in i.
func (i *Int) SetBytes(a []byte) kyber.Scalar {
	var buff = a
	if i.BO == LittleEndian {
		buff = reverse(nil, a)
	}
	i.V.SetBytes(buff).Mod(&i.V, i.M)
	return i
}

// LittleEndian encodes the value of this Int into a little-endian byte-slice
// at least min bytes but no more than max bytes long.
// Panics if max != 0 and the Int cannot be represented in max bytes.
func (i *Int) LittleEndian(min, max int) []byte {
	act := i.MarshalSize()
	vBytes := i.V.Bytes()
	vSize := len(vBytes)
	if vSize < act {
		act = vSize
	}
	pad := act
	if pad < min {
		pad = min
	}
	if max != 0 && pad > max {
		panic("Int not representable in max bytes")
	}
	buf := make([]byte, pad)
	reverse(buf[:act], vBytes)
	return buf
}

// reverse copies src into dst in byte-reversed order and returns dst,
// such that src[0] goes into dst[len-1] and vice versa.
// dst and src may be the same slice but otherwise must not overlap.
func reverse(dst, src []byte) []byte {
	if dst == nil {
		dst = make([]byte, len(src))
	}
	l := len(dst)
	for i, j := 0, l-1; i < (l+1)/2; {
		dst[i], dst[j] = src[j], src[i]
		i++
		j--
	}
	return dst
}
