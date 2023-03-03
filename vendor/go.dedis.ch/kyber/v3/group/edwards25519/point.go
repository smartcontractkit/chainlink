// Package edwards25519 provides an optimized Go implementation of a
// Twisted Edwards curve that is isomorphic to Curve25519. For details see:
// http://ed25519.cr.yp.to/.
//
// This code is based on Adam Langley's Go port of the public domain,
// "ref10" implementation of the ed25519 signing scheme in C from SUPERCOP.
// It was generalized and extended to support full kyber.Group arithmetic
// by the DEDIS lab at Yale and EPFL.
//
// Due to the field element and group arithmetic optimizations
// described in the Ed25519 paper, this implementation generally
// performs extremely well, typically comparable to native C
// implementations.  The tradeoff is that this code is completely
// specialized to a single curve.
package edwards25519

import (
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"io"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/internal/marshalling"
)

var marshalPointID = [8]byte{'e', 'd', '.', 'p', 'o', 'i', 'n', 't'}

type point struct {
	ge      extendedGroupElement
	varTime bool
}

func (P *point) String() string {
	var b [32]byte
	P.ge.ToBytes(&b)
	return hex.EncodeToString(b[:])
}

func (P *point) MarshalSize() int {
	return 32
}

func (P *point) MarshalBinary() ([]byte, error) {
	var b [32]byte
	P.ge.ToBytes(&b)
	return b[:], nil
}

// MarshalID returns the type tag used in encoding/decoding
func (P *point) MarshalID() [8]byte {
	return marshalPointID
}

func (P *point) UnmarshalBinary(b []byte) error {
	if !P.ge.FromBytes(b) {
		return errors.New("invalid Ed25519 curve point")
	}
	return nil
}

func (P *point) MarshalTo(w io.Writer) (int, error) {
	return marshalling.PointMarshalTo(P, w)
}

func (P *point) UnmarshalFrom(r io.Reader) (int, error) {
	return marshalling.PointUnmarshalFrom(P, r)
}

// Equality test for two Points on the same curve
func (P *point) Equal(P2 kyber.Point) bool {

	var b1, b2 [32]byte
	P.ge.ToBytes(&b1)
	P2.(*point).ge.ToBytes(&b2)
	for i := range b1 {
		if b1[i] != b2[i] {
			return false
		}
	}
	return true
}

// Set point to be equal to P2.
func (P *point) Set(P2 kyber.Point) kyber.Point {
	P.ge = P2.(*point).ge
	return P
}

// Set point to be equal to P2.
func (P *point) Clone() kyber.Point {
	return &point{ge: P.ge}
}

// Set to the neutral element, which is (0,1) for twisted Edwards curves.
func (P *point) Null() kyber.Point {
	P.ge.Zero()
	return P
}

// Set to the standard base point for this curve
func (P *point) Base() kyber.Point {
	P.ge = baseext
	return P
}

func (P *point) EmbedLen() int {
	// Reserve the most-significant 8 bits for pseudo-randomness.
	// Reserve the least-significant 8 bits for embedded data length.
	// (Hopefully it's unlikely we'll need >=2048-bit curves soon.)
	return (255 - 8 - 8) / 8
}

func (P *point) Embed(data []byte, rand cipher.Stream) kyber.Point {

	// How many bytes to embed?
	dl := P.EmbedLen()
	if dl > len(data) {
		dl = len(data)
	}

	for {
		// Pick a random point, with optional embedded data
		var b [32]byte
		rand.XORKeyStream(b[:], b[:])
		if data != nil {
			b[0] = byte(dl)       // Encode length in low 8 bits
			copy(b[1:1+dl], data) // Copy in data to embed
		}
		if !P.ge.FromBytes(b[:]) { // Try to decode
			continue // invalid point, retry
		}

		// If we're using the full group,
		// we just need any point on the curve, so we're done.
		//		if c.full {
		//			return P,data[dl:]
		//		}

		// We're using the prime-order subgroup,
		// so we need to make sure the point is in that subencoding.
		// If we're not trying to embed data,
		// we can convert our point into one in the subgroup
		// simply by multiplying it by the cofactor.
		if data == nil {
			P.Mul(cofactorScalar, P) // multiply by cofactor
			if P.Equal(nullPoint) {
				continue // unlucky; try again
			}
			return P // success
		}

		// Since we need the point's y-coordinate to hold our data,
		// we must simply check if the point is in the subgroup
		// and retry point generation until it is.
		var Q point
		Q.Mul(primeOrderScalar, P)
		if Q.Equal(nullPoint) {
			return P // success
		}
		// Keep trying...
	}
}

func (P *point) Pick(rand cipher.Stream) kyber.Point {
	return P.Embed(nil, rand)
}

// Extract embedded data from a point group element
func (P *point) Data() ([]byte, error) {
	var b [32]byte
	P.ge.ToBytes(&b)
	dl := int(b[0]) // extract length byte
	if dl > P.EmbedLen() {
		return nil, errors.New("invalid embedded data length")
	}
	return b[1 : 1+dl], nil
}

func (P *point) Add(P1, P2 kyber.Point) kyber.Point {
	E1 := P1.(*point)
	E2 := P2.(*point)

	var t2 cachedGroupElement
	var r completedGroupElement

	E2.ge.ToCached(&t2)
	r.Add(&E1.ge, &t2)
	r.ToExtended(&P.ge)

	return P
}

func (P *point) Sub(P1, P2 kyber.Point) kyber.Point {
	E1 := P1.(*point)
	E2 := P2.(*point)

	var t2 cachedGroupElement
	var r completedGroupElement

	E2.ge.ToCached(&t2)
	r.Sub(&E1.ge, &t2)
	r.ToExtended(&P.ge)

	return P
}

// Neg finds the negative of point A.
// For Edwards curves, the negative of (x,y) is (-x,y).
func (P *point) Neg(A kyber.Point) kyber.Point {
	P.ge.Neg(&A.(*point).ge)
	return P
}

// Mul multiplies point p by scalar s using the repeated doubling method.
func (P *point) Mul(s kyber.Scalar, A kyber.Point) kyber.Point {

	a := &s.(*scalar).v

	if A == nil {
		geScalarMultBase(&P.ge, a)
	} else {
		if P.varTime {
			geScalarMultVartime(&P.ge, a, &A.(*point).ge)
		} else {
			geScalarMult(&P.ge, a, &A.(*point).ge)
		}
	}

	return P
}

// HasSmallOrder determines whether the group element has small order
//
// Provides resilience against malicious key substitution attacks (M-S-UEO)
// and message bound security (MSB) even for malicious keys
// See paper https://eprint.iacr.org/2020/823.pdf for definitions and theorems
//
// This is the same code as in
// https://github.com/jedisct1/libsodium/blob/4744636721d2e420f8bbe2d563f31b1f5e682229/src/libsodium/crypto_core/ed25519/ref10/ed25519_ref10.c#L1170
func (P *point) HasSmallOrder() bool {
	s, err := P.MarshalBinary()
	if err != nil {
		return false
	}

	var c [5]byte

	for j := 0; j < 31; j++ {
		for i := 0; i < 5; i++ {
			c[i] |= s[j] ^ weakKeys[i][j]
		}
	}
	for i := 0; i < 5; i++ {
		c[i] |= (s[31] & 0x7f) ^ weakKeys[i][31]
	}

	// Constant time verification if one or more of the c's are zero
	var k uint16
	for i := 0; i < 5; i++ {
		k |= uint16(c[i]) - 1
	}

	return (k>>8)&1 > 0
}

// IsCanonical determines whether the group element is canonical
//
// Checks whether group element s is less than p, according to RFC8032§5.1.3.1
// https://tools.ietf.org/html/rfc8032#section-5.1.3
//
// Taken from
// https://github.com/jedisct1/libsodium/blob/4744636721d2e420f8bbe2d563f31b1f5e682229/src/libsodium/crypto_core/ed25519/ref10/ed25519_ref10.c#L1113
//
// The method accepts a buffer instead of calling `MarshalBinary` on the receiver
// because that always returns a value modulo `prime`.
func (P *point) IsCanonical(s []byte) bool {
	if len(s) != 32 {
		return false
	}

	c := (s[31] & 0x7f) ^ 0x7f
	for i := 30; i > 0; i-- {
		c |= s[i] ^ 0xff
	}

	// subtraction might underflow
	c = byte((uint16(c) - 1) >> 8)
	d := byte((0xed - 1 - uint16(s[0])) >> 8)

	return 1-(c&d&1) == 1
}
