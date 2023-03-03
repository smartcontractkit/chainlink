package curve25519

import (
	"crypto/cipher"
	"io"
	"math/big"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/internal/marshalling"
	"go.dedis.ch/kyber/v3/group/mod"
)

type projPoint struct {
	X, Y, Z mod.Int
	c       *ProjectiveCurve
}

func (P *projPoint) initXY(x, y *big.Int, c kyber.Group) {
	P.c = c.(*ProjectiveCurve)
	P.X.Init(x, &P.c.P)
	P.Y.Init(y, &P.c.P)
	P.Z.Init64(1, &P.c.P)
}

func (P *projPoint) getXY() (x, y *mod.Int) {
	P.normalize()
	return &P.X, &P.Y
}

func (P *projPoint) String() string {
	P.normalize()
	return P.c.pointString(&P.X, &P.Y)
}

func (P *projPoint) MarshalSize() int {
	return P.c.PointLen()
}

func (P *projPoint) MarshalBinary() ([]byte, error) {
	P.normalize()
	return P.c.encodePoint(&P.X, &P.Y), nil
}

func (P *projPoint) UnmarshalBinary(b []byte) error {
	P.Z.Init64(1, &P.c.P)
	return P.c.decodePoint(b, &P.X, &P.Y)
}

func (P *projPoint) MarshalTo(w io.Writer) (int, error) {
	return marshalling.PointMarshalTo(P, w)
}

func (P *projPoint) UnmarshalFrom(r io.Reader) (int, error) {
	return marshalling.PointUnmarshalFrom(P, r)
}

// Equality test for two Points on the same curve.
// We can avoid inversions here because:
//
//	(X1/Z1,Y1/Z1) == (X2/Z2,Y2/Z2)
//		iff
//	(X1*Z2,Y1*Z2) == (X2*Z1,Y2*Z1)
//
func (P *projPoint) Equal(CP2 kyber.Point) bool {
	P2 := CP2.(*projPoint)
	var t1, t2 mod.Int
	xeq := t1.Mul(&P.X, &P2.Z).Equal(t2.Mul(&P2.X, &P.Z))
	yeq := t1.Mul(&P.Y, &P2.Z).Equal(t2.Mul(&P2.Y, &P.Z))
	return xeq && yeq
}

func (P *projPoint) Set(CP2 kyber.Point) kyber.Point {
	P2 := CP2.(*projPoint)
	P.c = P2.c
	P.X.Set(&P2.X)
	P.Y.Set(&P2.Y)
	P.Z.Set(&P2.Z)
	return P
}

func (P *projPoint) Clone() kyber.Point {
	P2 := projPoint{}
	P2.c = P.c
	P2.X.Set(&P.X)
	P2.Y.Set(&P.Y)
	P2.Z.Set(&P.Z)
	return &P2
}

func (P *projPoint) Null() kyber.Point {
	P.Set(&P.c.null)
	return P
}

func (P *projPoint) Base() kyber.Point {
	P.Set(&P.c.base)
	return P
}

func (P *projPoint) EmbedLen() int {
	return P.c.embedLen()
}

// Normalize the point's representation to Z=1.
func (P *projPoint) normalize() {
	P.Z.Inv(&P.Z)
	P.X.Mul(&P.X, &P.Z)
	P.Y.Mul(&P.Y, &P.Z)
	P.Z.V.SetInt64(1)
}

func (P *projPoint) Embed(data []byte, rand cipher.Stream) kyber.Point {
	P.c.embed(P, data, rand)
	return P
}

func (P *projPoint) Pick(rand cipher.Stream) kyber.Point {
	return P.Embed(nil, rand)
}

// Extract embedded data from a point group element
func (P *projPoint) Data() ([]byte, error) {
	P.normalize()
	return P.c.data(&P.X, &P.Y)
}

// Add two points using optimized projective coordinate addition formulas.
// Formulas taken from:
//
//	http://eprint.iacr.org/2008/013.pdf
//	https://hyperelliptic.org/EFD/g1p/auto-twisted-projective.html
//
func (P *projPoint) Add(CP1, CP2 kyber.Point) kyber.Point {
	P1 := CP1.(*projPoint)
	P2 := CP2.(*projPoint)
	X1, Y1, Z1 := &P1.X, &P1.Y, &P1.Z
	X2, Y2, Z2 := &P2.X, &P2.Y, &P2.Z
	var A, B, C, D, E, F, G, X3, Y3, Z3 mod.Int

	A.Mul(Z1, Z2)
	B.Mul(&A, &A)
	C.Mul(X1, X2)
	D.Mul(Y1, Y2)
	E.Mul(&C, &D).Mul(&P.c.d, &E)
	F.Sub(&B, &E)
	G.Add(&B, &E)
	X3.Add(X1, Y1).Mul(&X3, Z3.Add(X2, Y2)).Sub(&X3, &C).Sub(&X3, &D).
		Mul(&F, &X3).Mul(&A, &X3)
	Y3.Mul(&P.c.a, &C).Sub(&D, &Y3).Mul(&G, &Y3).Mul(&A, &Y3)
	Z3.Mul(&F, &G)

	P.c = P1.c
	P.X.Set(&X3)
	P.Y.Set(&Y3)
	P.Z.Set(&Z3)
	return P
}

// Subtract points so that their scalars subtract homomorphically
func (P *projPoint) Sub(CP1, CP2 kyber.Point) kyber.Point {
	P1 := CP1.(*projPoint)
	P2 := CP2.(*projPoint)
	X1, Y1, Z1 := &P1.X, &P1.Y, &P1.Z
	X2, Y2, Z2 := &P2.X, &P2.Y, &P2.Z
	var A, B, C, D, E, F, G, X3, Y3, Z3 mod.Int

	A.Mul(Z1, Z2)
	B.Mul(&A, &A)
	C.Mul(X1, X2)
	D.Mul(Y1, Y2)
	E.Mul(&C, &D).Mul(&P.c.d, &E)
	F.Add(&B, &E)
	G.Sub(&B, &E)
	X3.Add(X1, Y1).Mul(&X3, Z3.Sub(Y2, X2)).Add(&X3, &C).Sub(&X3, &D).
		Mul(&F, &X3).Mul(&A, &X3)
	Y3.Mul(&P.c.a, &C).Add(&D, &Y3).Mul(&G, &Y3).Mul(&A, &Y3)
	Z3.Mul(&F, &G)

	P.c = P1.c
	P.X.Set(&X3)
	P.Y.Set(&Y3)
	P.Z.Set(&Z3)
	return P
}

// Find the negative of point A.
// For Edwards curves, the negative of (x,y) is (-x,y).
func (P *projPoint) Neg(CA kyber.Point) kyber.Point {
	A := CA.(*projPoint)
	P.c = A.c
	P.X.Neg(&A.X)
	P.Y.Set(&A.Y)
	P.Z.Set(&A.Z)
	return P
}

// Optimized point doubling for use in scalar multiplication.
func (P *projPoint) double() {
	var B, C, D, E, F, H, J mod.Int

	B.Add(&P.X, &P.Y).Mul(&B, &B)
	C.Mul(&P.X, &P.X)
	D.Mul(&P.Y, &P.Y)
	E.Mul(&P.c.a, &C)
	F.Add(&E, &D)
	H.Mul(&P.Z, &P.Z)
	J.Add(&H, &H).Sub(&F, &J)
	P.X.Sub(&B, &C).Sub(&P.X, &D).Mul(&P.X, &J)
	P.Y.Sub(&E, &D).Mul(&F, &P.Y)
	P.Z.Mul(&F, &J)
}

// Multiply point p by scalar s using the repeated doubling method.
func (P *projPoint) Mul(s kyber.Scalar, G kyber.Point) kyber.Point {
	v := s.(*mod.Int).V
	if G == nil {
		return P.Base().Mul(s, P)
	}
	T := P
	if G == P { // Must use temporary for in-place multiply
		T = &projPoint{}
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

// ProjectiveCurve implements Twisted Edwards curves
// using projective coordinate representation (X:Y:Z),
// satisfying the identities x = X/Z, y = Y/Z.
// This representation still supports all Twisted Edwards curves
// and avoids expensive modular inversions on the critical paths.
// Uses the projective arithmetic formulas in:
// http://cr.yp.to/newelliptic/newelliptic-20070906.pdf
//
type ProjectiveCurve struct {
	curve           // generic Edwards curve functionality
	null  projPoint // Constant identity/null point (0,1)
	base  projPoint // Standard base point
}

// Point creates a new Point on this curve.
func (c *ProjectiveCurve) Point() kyber.Point {
	P := new(projPoint)
	P.c = c
	//P.Set(&c.null)
	return P
}

// Init initializes the curve with given parameters.
func (c *ProjectiveCurve) Init(p *Param, fullGroup bool) *ProjectiveCurve {
	c.curve.init(c, p, fullGroup, &c.null, &c.base)
	return c
}
