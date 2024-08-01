// Package nist implements cryptographic groups and ciphersuites
// based on the NIST standards, using Go's built-in crypto library.
package nist

import (
	"crypto/cipher"
	"crypto/elliptic"
	"errors"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/tdh2/go/tdh2/internal/group"
	"github.com/smartcontractkit/tdh2/go/tdh2/internal/group/mod"
)

// streamReader implements io.Reader from cipher.Stream
type streamReader struct {
	stream cipher.Stream
}

func (s *streamReader) Read(p []byte) (int, error) {
	s.stream.XORKeyStream(p, p)
	return len(p), nil
}

type curvePoint struct {
	x, y *big.Int
	c    *curve
}

func (p *curvePoint) String() string {
	return "(" + p.x.String() + "," + p.y.String() + ")"
}

func (p *curvePoint) Equal(p2 group.Point) bool {
	cp2 := p2.(*curvePoint)

	// Make sure both coordinates are normalized.
	// Apparently Go's elliptic curve code doesn't always ensure this.
	M := p.c.p.P
	p.x.Mod(p.x, M)
	p.y.Mod(p.y, M)
	cp2.x.Mod(cp2.x, M)
	cp2.y.Mod(cp2.y, M)

	return p.x.Cmp(cp2.x) == 0 && p.y.Cmp(cp2.y) == 0
}

func (p *curvePoint) Null() group.Point {
	p.x = new(big.Int).SetInt64(0)
	p.y = new(big.Int).SetInt64(0)
	return p
}

func (p *curvePoint) Base() group.Point {
	p.x = p.c.p.Gx
	p.y = p.c.p.Gy
	return p
}

func (p *curvePoint) Valid() bool {
	// The IsOnCurve function in Go's elliptic curve package
	// doesn't consider the point-at-infinity to be "on the curve"
	return p.c.IsOnCurve(p.x, p.y) ||
		(p.x.Sign() == 0 && p.y.Sign() == 0)
}

func (p *curvePoint) Pick(rand cipher.Stream) group.Point {
	var err error
	_, p.x, p.y, err = elliptic.GenerateKey(p.c, &streamReader{rand})
	if err != nil {
		// It cannot panic since GenerateKey returns errors only on reading
		// from the randomness source which is deterministic in our case.
		panic(fmt.Sprintf("cannot generate point: %v", err))
	}
	return p
}

func (p *curvePoint) Add(a, b group.Point) group.Point {
	ca := a.(*curvePoint)
	cb := b.(*curvePoint)
	p.x, p.y = p.c.Add(ca.x, ca.y, cb.x, cb.y)
	return p
}

func (p *curvePoint) Sub(a, b group.Point) group.Point {
	ca := a.(*curvePoint)
	cb := b.(*curvePoint)

	cbn := p.c.Point().Neg(cb).(*curvePoint)
	p.x, p.y = p.c.Add(ca.x, ca.y, cbn.x, cbn.y)
	return p
}

func (p *curvePoint) Neg(a group.Point) group.Point {
	s := p.c.Scalar().One()
	s.Neg(s)
	return p.Mul(s, a).(*curvePoint)
}

func (p *curvePoint) Mul(s group.Scalar, b group.Point) group.Point {
	cs := s.(*mod.Int)
	if b != nil {
		cb := b.(*curvePoint)
		p.x, p.y = p.c.ScalarMult(cb.x, cb.y, cs.V.Bytes())
	} else {
		p.x, p.y = p.c.ScalarBaseMult(cs.V.Bytes())
	}
	return p
}

func (p *curvePoint) MarshalSize() int {
	coordlen := (p.c.Params().BitSize + 7) >> 3
	return 1 + 2*coordlen // uncompressed ANSI X9.62 representation
}

func (p *curvePoint) MarshalBinary() ([]byte, error) {
	return elliptic.Marshal(p.c, p.x, p.y), nil
}

func (p *curvePoint) UnmarshalBinary(buf []byte) error {
	if len(buf) != p.MarshalSize() {
		return errors.New("wrong buffer size")
	}
	// Check whether all bytes after first one are 0, so we
	// just return the initial point. Read everything to
	// prevent timing-leakage.
	var c byte
	for _, b := range buf[1:] {
		c |= b
	}
	if c != 0 {
		p.x, p.y = elliptic.Unmarshal(p.c, buf)
		if p.x == nil || !p.Valid() {
			return errors.New("invalid elliptic curve point")
		}
	} else {
		// All bytes are 0, so we initialize x and y
		p.x = big.NewInt(0)
		p.y = big.NewInt(0)
	}
	return nil
}

// Curve is an implementation of the group.Group interface
// for NIST elliptic curves, built on Go's native elliptic curve library.
type curve struct {
	elliptic.Curve
	p *elliptic.CurveParams
}

// Return the number of bytes in the encoding of a Scalar for this curve.
func (c *curve) ScalarLen() int { return (c.p.N.BitLen() + 7) / 8 }

// Create a Scalar associated with this curve. The scalars created by
// this package implement group.Scalar's SetBytes method, interpreting
// the bytes as a big-endian integer, so as to be compatible with the
// Go standard library's big.Int type.
func (c *curve) Scalar() group.Scalar {
	return mod.NewInt64(0, c.p.N)
}

// Number of bytes required to store one coordinate on this curve
func (c *curve) coordLen() int {
	return (c.p.BitSize + 7) / 8
}

// Return the number of bytes in the encoding of a Point for this curve.
// Currently uses uncompressed ANSI X9.62 format with both X and Y coordinates;
// this could change.
func (c *curve) PointLen() int {
	return 1 + 2*c.coordLen() // ANSI X9.62: 1 header byte plus 2 coords
}

// Create a Point associated with this curve.
func (c *curve) Point() group.Point {
	p := new(curvePoint)
	p.c = c
	return p
}

func (p *curvePoint) Set(P group.Point) group.Point {
	p.x = P.(*curvePoint).x
	p.y = P.(*curvePoint).y
	return p
}

func (p *curvePoint) Clone() group.Point {
	return &curvePoint{x: p.x, y: p.y, c: p.c}
}

// Return the order of this curve: the prime N in the curve parameters.
func (c *curve) Order() *big.Int {
	return c.p.N
}

// P256 implements the group.Group interface for the NIST P-256 elliptic curve.
type P256 struct {
	curve
}

func (curve *P256) String() string {
	return "P256"
}

// NewP256 returns a new instance of P256.
func NewP256() *P256 {
	var g P256
	g.curve.Curve = elliptic.P256()
	g.p = g.Params()
	return &g
}

// P384 implements the group.Group interface for the NIST P-384 elliptic curve.
type P384 struct {
	curve
}

func (curve *P384) String() string {
	return "P384"
}

// NewP384 returns a new instance of P384.
func NewP384() *P384 {
	var g P384
	g.curve.Curve = elliptic.P384()
	g.p = g.Params()
	return &g
}

// P521 implements the group.Group interface for the NIST P-521 elliptic curve.
type P521 struct {
	curve
}

func (curve *P521) String() string {
	return "P521"
}

// NewP521 returns a new instance of P521.
func NewP521() *P521 {
	var g P521
	g.curve.Curve = elliptic.P521()
	g.p = g.Params()
	return &g
}
