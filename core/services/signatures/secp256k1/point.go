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

// Implementation of kyber.Point interface for elliptic-curve arithmetic
// operations on secpk256k1.
//
// This is mostly a wrapper of the functionality provided by btcec

import (
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
	"math/big"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/util/key"
	"golang.org/x/crypto/sha3"
)

// btcec's public interface uses this affine representation for points on the
// curve. This does not naturally accommodate the point at infinity. btcec
// represents it as (0, 0), which is not a point on {y²=x³+7}.
type secp256k1Point struct {
	X *fieldElt
	Y *fieldElt
}

func newPoint() *secp256k1Point {
	return &secp256k1Point{newFieldZero(), newFieldZero()}
}

// String returns a string representation of P
func (P *secp256k1Point) String() string {
	return fmt.Sprintf("Secp256k1{X: %s, Y: %s}", P.X, P.Y)
}

// Equal returns true if p and pPrime represent the same point, false otherwise.
func (P *secp256k1Point) Equal(pPrime kyber.Point) bool {
	return P.X.Equal(pPrime.(*secp256k1Point).X) &&
		P.Y.Equal(pPrime.(*secp256k1Point).Y)
}

// Null sets p to the group-identity value, and returns it.
func (P *secp256k1Point) Null() kyber.Point {
	P.X = fieldEltFromInt(0) // btcec representation of null point is (0,0)
	P.Y = fieldEltFromInt(0)
	return P
}

// Base sets p to a copy of the standard group generator, and returns it.
func (P *secp256k1Point) Base() kyber.Point {
	P.X.SetInt(s256.Gx)
	P.Y.SetInt(s256.Gy)
	return P
}

// Pick sets P to a random point sampled from rand, and returns it.
func (P *secp256k1Point) Pick(rand cipher.Stream) kyber.Point {
	for { // Keep trying X's until one fits the curve (~50% probability of
		// success each iteration
		P.X.Set(newFieldZero().Pick(rand))
		maybeRHS := rightHandSide(P.X)
		if maybeY := maybeSqrtInField(maybeRHS); maybeY != (*fieldElt)(nil) {
			P.Y.Set(maybeY)
			// Take the negative with 50% probability
			b := make([]byte, 1)
			rand.XORKeyStream(b, b)
			if b[0]&1 == 0 {
				P.Y.Neg(P.Y)
			}
			return P
		}
	}
}

// Set sets P to copies of pPrime's values, and returns it.
func (P *secp256k1Point) Set(pPrime kyber.Point) kyber.Point {
	P.X.Set(pPrime.(*secp256k1Point).X)
	P.Y.Set(pPrime.(*secp256k1Point).Y)
	return P
}

// Clone returns a copy of P.
func (P *secp256k1Point) Clone() kyber.Point {
	return &secp256k1Point{X: P.X.Clone(), Y: P.Y.Clone()}
}

// EmbedLen returns the number of bytes of data which can be embedded in a point.
func (*secp256k1Point) EmbedLen() int {
	// Reserve the most-significant 8 bits for pseudo-randomness.
	// Reserve the least-significant 8 bits for embedded data length.
	return (255 - 8 - 8) / 8
}

// Embed encodes a limited amount of specified data in the Point, using r as a
// source of cryptographically secure random data. Implementations only embed
// the first EmbedLen bytes of the given data.
func (P *secp256k1Point) Embed(data []byte, r cipher.Stream) kyber.Point {
	numEmbedBytes := P.EmbedLen()
	if len(data) > numEmbedBytes {
		panic("too much data to embed in a point")
	}
	numEmbedBytes = len(data)
	var x [32]byte
	randStart := 1 // First byte to fill with random data
	if data != nil {
		x[0] = byte(numEmbedBytes)       // Encode length in low 8 bits
		copy(x[1:1+numEmbedBytes], data) // Copy in data to embed
		randStart = 1 + numEmbedBytes
	}
	maxAttempts := 10000
	// Try random x ordinates satisfying the constraints, until one provides
	// a point on secp256k1
	for numAttempts := 0; numAttempts < maxAttempts; numAttempts++ {
		// Fill the rest of the x ordinate with random data
		r.XORKeyStream(x[randStart:], x[randStart:])
		xOrdinate := newFieldZero().SetBytes(x)
		// RHS of secp256k1 equation is x³+7 mod p. Success if square.
		// We optimistically don't use btcec.IsOnCurve, here, because we
		// hope to assign the intermediate result maybeY to P.Y
		secp256k1RHS := rightHandSide(xOrdinate)
		if maybeY := maybeSqrtInField(secp256k1RHS); maybeY != (*fieldElt)(nil) {
			P.X = xOrdinate // success: found (x,y) s.t. y²=x³+7
			P.Y = maybeY
			return P
		}
	}
	// Probability 2^{-maxAttempts}, under correct operation.
	panic("failed to find point satisfying all constraints")
}

// Data returns data embedded in P, or an error if inconsistent with encoding
func (P *secp256k1Point) Data() ([]byte, error) {
	b := P.X.Bytes()
	dataLength := int(b[0])
	if dataLength > P.EmbedLen() {
		return nil, errors.New("point specifies too much data")
	}
	return b[1 : dataLength+1], nil
}

// Add sets P to a+b (secp256k1 group operation) and returns it.
func (P *secp256k1Point) Add(a, b kyber.Point) kyber.Point {
	X, Y := s256.Add(
		a.(*secp256k1Point).X.int(), a.(*secp256k1Point).Y.int(),
		b.(*secp256k1Point).X.int(), b.(*secp256k1Point).Y.int())
	P.X.SetInt(X)
	P.Y.SetInt(Y)
	return P
}

// Add sets P to a-b (secp256k1 group operation), and returns it.
func (P *secp256k1Point) Sub(a, b kyber.Point) kyber.Point {
	X, Y := s256.Add(
		a.(*secp256k1Point).X.int(), a.(*secp256k1Point).Y.int(),
		b.(*secp256k1Point).X.int(),
		newFieldZero().Neg(b.(*secp256k1Point).Y).int()) // -b_y
	P.X.SetInt(X)
	P.Y.SetInt(Y)
	return P
}

// Neg sets P to -a (in the secp256k1 group), and returns it.
func (P *secp256k1Point) Neg(a kyber.Point) kyber.Point {
	P.X = a.(*secp256k1Point).X.Clone()
	P.Y = newFieldZero().Neg(a.(*secp256k1Point).Y)
	return P
}

// Mul sets P to s*a (in the secp256k1 group, i.e. adding a to itself s times),
// and returns it. If a is nil, it is replaced by the secp256k1 generator.
func (P *secp256k1Point) Mul(s kyber.Scalar, a kyber.Point) kyber.Point {
	sBytes, err := s.(*secp256k1Scalar).MarshalBinary()
	if err != nil {
		panic(fmt.Errorf("failure while marshaling multiplier: %w",
			err))
	}
	var X, Y *big.Int
	if a == (*secp256k1Point)(nil) || a == nil {
		X, Y = s256.ScalarBaseMult(sBytes)
	} else {
		X, Y = s256.ScalarMult(a.(*secp256k1Point).X.int(),
			a.(*secp256k1Point).Y.int(), sBytes)
	}
	P.X.SetInt(X)
	P.Y.SetInt(Y)
	return P
}

// MarshalBinary returns the concatenated big-endian representation of the X
// ordinate and a byte which is 0 if Y is even, 1 if it's odd. Or it returns an
// error on failure.
func (P *secp256k1Point) MarshalBinary() ([]byte, error) {
	maybeSqrt := maybeSqrtInField(rightHandSide(P.X))
	if maybeSqrt == (*fieldElt)(nil) {
		return nil, errors.New("x³+7 not a square")
	}
	minusMaybeSqrt := newFieldZero().Neg(maybeSqrt)
	if !P.Y.Equal(maybeSqrt) && !P.Y.Equal(minusMaybeSqrt) {
		return nil, errors.New("y ≠ ±maybeSqrt(x³+7), so not a point on the curve")
	}
	rv := make([]byte, P.MarshalSize())
	signByte := P.MarshalSize() - 1 // Last byte contains sign of Y.
	xordinate := P.X.Bytes()
	copyLen := copy(rv[:signByte], xordinate[:])
	if copyLen != P.MarshalSize()-1 {
		return []byte{}, errors.New("marshal of x ordinate too short")
	}
	if P.Y.isEven() {
		rv[signByte] = 0
	} else {
		rv[signByte] = 1
	}
	return rv, nil
}

// MarshalSize returns the length of the byte representation of P
func (P *secp256k1Point) MarshalSize() int { return 33 }

// MarshalID returns the ID for a secp256k1 point
func (P *secp256k1Point) MarshalID() [8]byte {
	return [8]byte{'s', 'p', '2', '5', '6', '.', 'p', 'o'}
}

// UnmarshalBinary sets P to the point represented by contents of buf, or
// returns an non-nil error
func (P *secp256k1Point) UnmarshalBinary(buf []byte) error {
	var err error
	if len(buf) != P.MarshalSize() {
		err = errors.New("wrong length for marshaled point")
	}
	if err == nil && !(buf[32] == 0 || buf[32] == 1) {
		err = errors.New("bad sign byte (the last one)")
	}
	if err != nil {
		return err
	}
	var xordinate [32]byte
	copy(xordinate[:], buf[:32])
	P.X = newFieldZero().SetBytes(xordinate)
	secp256k1RHS := rightHandSide(P.X)
	maybeY := maybeSqrtInField(secp256k1RHS)
	if maybeY == (*fieldElt)(nil) {
		return errors.New("x ordinate does not correspond to a curve point")
	}
	isEven := maybeY.isEven()
	P.Y.Set(maybeY)
	if (buf[32] == 0 && !isEven) || (buf[32] == 1 && isEven) {
		P.Y.Neg(P.Y)
	} else {
		if buf[32] != 0 && buf[32] != 1 {
			return errors.New("parity byte must be 0 or 1")
		}
	}
	return nil
}

// MarshalTo writes the serialized P to w, and returns the number of bytes
// written, or an error on failure.
func (P *secp256k1Point) MarshalTo(w io.Writer) (int, error) {
	buf, err := P.MarshalBinary()
	if err != nil {
		return 0, err
	}
	return w.Write(buf)
}

// UnmarshalFrom sets P to the secp256k1 point represented by bytes read from r,
// and returns the number of bytes read, or an error on failure.
func (P *secp256k1Point) UnmarshalFrom(r io.Reader) (int, error) {
	buf := make([]byte, P.MarshalSize())
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, err
	}
	return n, P.UnmarshalBinary(buf)
}

// EthereumAddress returns the 160-bit address corresponding to p as public key.
func EthereumAddress(p kyber.Point) (rv [20]byte) {
	// The Ethereum address of P is the bottom 160 bits of keccak256(P.X‖P.Y),
	// where P.X and P.Y are represented in 32 bytes as big-endian. See equations
	// (277, 284) of Ethereum Yellow Paper version 3e36772, or go-ethereum's
	// crypto.PubkeyToAddress.
	h := sha3.NewLegacyKeccak256()
	if _, err := h.Write(LongMarshal(p)); err != nil {
		panic(err)
	}
	copy(rv[:], h.Sum(nil)[12:])
	return rv
}

// IsSecp256k1Point returns true if p is a secp256k1Point
func IsSecp256k1Point(p kyber.Point) bool {
	switch p.(type) {
	case *secp256k1Point:
		return true
	default:
		return false
	}
}

// Coordinates returns the coordinates of p
func Coordinates(p kyber.Point) (*big.Int, *big.Int) {
	return p.(*secp256k1Point).X.int(), p.(*secp256k1Point).Y.int()
}

// ValidPublicKey returns true iff p can be used in the optimized on-chain
// Schnorr-signature verification. See SchnorrSECP256K1.sol for details.
func ValidPublicKey(p kyber.Point) bool {
	if p == (*secp256k1Point)(nil) || p == nil {
		return false
	}
	P, ok := p.(*secp256k1Point)
	if !ok {
		return false
	}
	maybeY := maybeSqrtInField(rightHandSide(P.X))
	return maybeY != nil && (P.Y.Equal(maybeY) || P.Y.Equal(maybeY.Neg(maybeY)))
}

// Generate generates a public/private key pair, which can be verified cheaply
// on-chain
func Generate(random cipher.Stream) *key.Pair {
	p := key.Pair{}
	for !ValidPublicKey(p.Public) {
		p.Private = (&Secp256k1{}).Scalar().Pick(random)
		p.Public = (&Secp256k1{}).Point().Mul(p.Private, nil)
	}
	return &p
}

// LongMarshal returns the concatenated coordinates serialized as uint256's
func LongMarshal(p kyber.Point) []byte {
	xMarshal := p.(*secp256k1Point).X.Bytes()
	yMarshal := p.(*secp256k1Point).Y.Bytes()
	return append(xMarshal[:], yMarshal[:]...)
}

// LongUnmarshal returns the secp256k1 point represented by m, as a concatenated
// pair of uint256's
func LongUnmarshal(m []byte) (kyber.Point, error) {
	if len(m) != 64 {
		return nil, fmt.Errorf(
			"0x%x does not represent an uncompressed secp256k1Point. Should be length 64, but is length %d",
			m, len(m))
	}
	p := newPoint()
	p.X.SetInt(big.NewInt(0).SetBytes(m[:32]))
	p.Y.SetInt(big.NewInt(0).SetBytes(m[32:]))
	if !ValidPublicKey(p) {
		return nil, fmt.Errorf("%s is not a valid secp256k1 point", p)
	}
	return p, nil
}

// ScalarToPublicPoint returns the public secp256k1 point associated to s
func ScalarToPublicPoint(s kyber.Scalar) kyber.Point {
	publicPoint := (&Secp256k1{}).Point()
	return publicPoint.Mul(s, nil)
}

// SetCoordinates returns the point (x,y), or panics if an invalid secp256k1Point
func SetCoordinates(x, y *big.Int) kyber.Point {
	rv := newPoint()
	rv.X.SetInt(x)
	rv.Y.SetInt(y)
	if !ValidPublicKey(rv) {
		panic("point requested from invalid coordinates")
	}
	return rv
}
