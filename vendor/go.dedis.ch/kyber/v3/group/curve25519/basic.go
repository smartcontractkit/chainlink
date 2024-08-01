//go:build experimental
// +build experimental

package curve25519

import (
	"crypto/cipher"
	"io"
	"math/big"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/internal/marshalling"
	"go.dedis.ch/kyber/v3/group/mod"
)

type basicPoint struct {
	x, y mod.Int
	c    *BasicCurve
}

func (P *basicPoint) initXY(x, y *big.Int, c kyber.Group) {
	P.c = c.(*BasicCurve)
	P.x.Init(x, &P.c.P)
	P.y.Init(y, &P.c.P)
}

func (P *basicPoint) getXY() (x, y *mod.Int) {
	return &P.x, &P.y
}

func (P *basicPoint) String() string {
	return P.c.pointString(&P.x, &P.y)
}

// coord creates a new ModInt representing a coordinate on this curve,
// with a given int64 integer value for constant-initialization convenience.
func (P *basicPoint) coord(v int64) *mod.Int {
	return mod.NewInt64(v, &P.c.P)
}

func (P *basicPoint) MarshalSize() int {
	return (P.y.M.BitLen() + 7 + 1) / 8
}

// MarshalBinary encodew an Edwards curve point.
func (P *basicPoint) MarshalBinary() ([]byte, error) {
	return P.c.encodePoint(&P.x, &P.y), nil
}

// UnmarshalBinary decodes an Edwards curve point.
func (P *basicPoint) UnmarshalBinary(b []byte) error {
	return P.c.decodePoint(b, &P.x, &P.y)
}

func (P *basicPoint) MarshalTo(w io.Writer) (int, error) {
	return marshalling.PointMarshalTo(P, w)
}

func (P *basicPoint) UnmarshalFrom(r io.Reader) (int, error) {
	return marshalling.PointUnmarshalFrom(P, r)
}

// Equal tests for two Points on the same curve
func (P *basicPoint) Equal(P2 kyber.Point) bool {
	E2 := P2.(*basicPoint)
	return P.x.Equal(&E2.x) && P.y.Equal(&E2.y)
}

// Set point to be equal to P2.
func (P *basicPoint) Set(P2 kyber.Point) kyber.Point {
	E2 := P2.(*basicPoint)
	P.c = E2.c
	P.x.Set(&E2.x)
	P.y.Set(&E2.y)
	return P
}

// Clone returns the given point
func (P *basicPoint) Clone() kyber.Point {
	p2 := new(basicPoint)
	p2.Set(P)
	return p2
}

// Null sets to the neutral element, which is (0,1) for twisted Edwards curves.
func (P *basicPoint) Null() kyber.Point {
	P.Set(&P.c.null)
	return P
}

// Base sets to the standard base point for this curve
func (P *basicPoint) Base() kyber.Point {
	P.Set(&P.c.base)
	return P
}

func (P *basicPoint) EmbedLen() int {
	return P.c.embedLen()
}

func (P *basicPoint) Embed(data []byte, rand cipher.Stream) kyber.Point {
	P.c.embed(P, data, rand)
	return P
}

func (P *basicPoint) Pick(rand cipher.Stream) kyber.Point {
	return P.Embed(nil, rand)
}

// Data extracts embedded data from a point group element
func (P *basicPoint) Data() ([]byte, error) {
	return P.c.data(&P.x, &P.y)
}

// Add two points using the basic unified addition laws for Edwards curves:
//
//	x' = ((x1*y2 + x2*y1) / (1 + d*x1*x2*y1*y2))
//	y' = ((y1*y2 - a*x1*x2) / (1 - d*x1*x2*y1*y2))
//
func (P *basicPoint) Add(P1, P2 kyber.Point) kyber.Point {
	E1 := P1.(*basicPoint)
	E2 := P2.(*basicPoint)
	x1, y1 := E1.x, E1.y
	x2, y2 := E2.x, E2.y

	var t1, t2, dm, nx, dx, ny, dy mod.Int

	// Reused part of denominator: dm = d*x1*x2*y1*y2
	dm.Mul(&P.c.d, &x1).Mul(&dm, &x2).Mul(&dm, &y1).Mul(&dm, &y2)

	// x' numerator/denominator
	nx.Add(t1.Mul(&x1, &y2), t2.Mul(&x2, &y1))
	dx.Add(&P.c.one, &dm)

	// y' numerator/denominator
	ny.Sub(t1.Mul(&y1, &y2), t2.Mul(&x1, &x2).Mul(&P.c.a, &t2))
	dy.Sub(&P.c.one, &dm)

	// result point
	P.x.Div(&nx, &dx)
	P.y.Div(&ny, &dy)
	return P
}

// Point doubling, which for Edwards curves can be accomplished
// simply by adding a point to itself (no exceptions for equal input points).
func (P *basicPoint) double() kyber.Point {
	return P.Add(P, P)
}

// Subtract points so that their scalars subtract homomorphically
func (P *basicPoint) Sub(A, B kyber.Point) kyber.Point {
	var nB basicPoint
	return P.Add(A, nB.Neg(B))
}

// Find the negative of point A.
// For Edwards curves, the negative of (x,y) is (-x,y).
func (P *basicPoint) Neg(A kyber.Point) kyber.Point {
	E := A.(*basicPoint)
	P.c = E.c
	P.x.Neg(&E.x)
	P.y.Set(&E.y)
	return P
}

// Multiply point p by scalar s using the repeated doubling method.
func (P *basicPoint) Mul(s kyber.Scalar, G kyber.Point) kyber.Point {
	v := s.(*mod.Int).V
	if G == nil {
		return P.Base().Mul(s, P)
	}
	T := P
	if G == P { // Must use temporary in case G == P
		T = &basicPoint{}
	}
	T.Set(&P.c.null) // Initialize to identity element (0,1)
	for i := v.BitLen() - 1; i >= 0; i-- {
		T.double()
		if v.Bit(i) != 0 {
			T.Add(T, G)
		}
	}
	if T != P {
		P.Set(T)
	}
	return P
}

// Basic unoptimized reference implementation of Twisted Edwards curves.
// This reference implementation is mainly intended for testing, debugging,
// and instructional uses, and not for production use.
// The projective coordinates implementation (ProjectiveCurve)
// is just as general and much faster.
//
type BasicCurve struct {
	curve            // generic Edwards curve functionality
	null  basicPoint // Neutral/identity point (0,1)
	base  basicPoint // Standard base point
}

// Create a new Point on this curve.
func (c *BasicCurve) Point() kyber.Point {
	P := new(basicPoint)
	P.c = c
	P.Set(&c.null)
	return P
}

// Initialize the curve with given parameters.
func (c *BasicCurve) Init(p *Param, fullGroup bool) *BasicCurve {
	c.curve.init(c, p, fullGroup, &c.null, &c.base)
	return c
}
