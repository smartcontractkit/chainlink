package bandersnatch

import (
	"fmt"
	"io"

	gnarkbandersnatch "github.com/consensys/gnark-crypto/ecc/bls12-381/bandersnatch"
	gnarkfr "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/crate-crypto/go-ipa/bandersnatch/fp"
)

var CurveParams = gnarkbandersnatch.GetEdwardsCurve()

type PointAffine = gnarkbandersnatch.PointAffine
type PointProj = gnarkbandersnatch.PointProj
type PointExtended = gnarkbandersnatch.PointExtended

var Identity = PointProj{
	X: fp.Zero(),
	Y: fp.One(),
	Z: fp.One(),
}

var IdentityExt = PointExtendedFromProj(&Identity)

// Reads an uncompressed affine point
// Point is not guaranteed to be in the prime subgroup
func ReadUncompressedPoint(r io.Reader) (PointAffine, error) {
	var xy = make([]byte, 64)
	if _, err := io.ReadAtLeast(r, xy, 64); err != nil {
		return PointAffine{}, fmt.Errorf("reading bytes: %s", err)
	}

	var x_fp, y_fp fp.Element
	x_fp.SetBytes(xy[:32])
	y_fp.SetBytes(xy[32:])

	return PointAffine{
		X: x_fp,
		Y: y_fp,
	}, nil
}

// Writes an uncompressed affine point to an io.Writer
func WriteUncompressedPoint(w io.Writer, p *PointAffine) (int, error) {
	x_bytes := p.X.Bytes()
	y_bytes := p.Y.Bytes()
	n1, err := w.Write(x_bytes[:])
	if err != nil {
		return n1, err
	}
	n2, err := w.Write(y_bytes[:])
	total_bytes_written := n1 + n2
	if err != nil {
		return total_bytes_written, err
	}
	return total_bytes_written, nil
}

func GetPointFromX(x *fp.Element, choose_largest bool) *PointAffine {
	y := computeY(x, choose_largest)
	if y == nil { // not a square
		return nil
	}
	return &PointAffine{X: *x, Y: *y}
}

// ax^2 + y^2 = 1 + dx^2y^2
// ax^2 -1 = dx^2y^2 - y^2
// ax^2 -1 = y^2(dx^2 -1)
// ax^2 - 1 / (dx^2 - 1) = y^2
func computeY(x *fp.Element, choose_largest bool) *fp.Element {
	var one, num, den, y fp.Element
	one.SetOne()
	num.Square(x)                 // x^2
	den.Mul(&num, &CurveParams.D) //dx^2
	den.Sub(&den, &one)           //dx^2 - 1

	num.Mul(&num, &CurveParams.A) // ax^2
	num.Sub(&num, &one)           // ax^2 - 1
	y.Div(&num, &den)
	sqrtY := fp.SqrtPrecomp(&y)

	// If the square root does not exist, then the Sqrt method returns nil
	// and leaves the receiver unchanged.
	// Note the fact that it leaves the receiver unchanged, means we do not return &y
	if sqrtY == nil {
		return nil
	}

	// Choose between `y` and it's negation
	is_largest := sqrtY.LexicographicallyLargest()
	if choose_largest == is_largest {
		return sqrtY
	} else {
		return sqrtY.Neg(sqrtY)
	}
}

// PointExtendedFromProj converts a point in projective coordinates to extended coordinates.
func PointExtendedFromProj(p *PointProj) PointExtended {
	var pzinv fp.Element
	pzinv.Inverse(&p.Z)
	var z fp.Element
	z.Mul(&p.X, &p.Y).Mul(&z, &pzinv)
	return PointExtended{
		X: p.X,
		Y: p.Y,
		Z: p.Z,
		T: z,
	}
}

// PointExtendedNormalized is an extended point which is normalized.
// i.e: Z=1. We store it this way to save 32 bytes per point in memory.
type PointExtendedNormalized struct {
	X, Y, T gnarkfr.Element
}

// Neg computes p = -p1
func (p *PointExtendedNormalized) Neg(p1 *PointExtendedNormalized) *PointExtendedNormalized {
	p.X.Neg(&p1.X)
	p.Y = p1.Y
	p.T.Neg(&p1.T)
	return p
}

// ExtendedAddNormalized computes p = p1 + p2.
// https://hyperelliptic.org/EFD/g1p/auto-twisted-extended.html#addition-madd-2008-hwcd
func ExtendedAddNormalized(p, p1 *PointExtended, p2 *PointExtendedNormalized) *gnarkbandersnatch.PointExtended {
	var A, B, C, D, E, F, G, H, tmp gnarkfr.Element
	A.Mul(&p1.X, &p2.X)
	B.Mul(&p1.Y, &p2.Y)
	C.Mul(&p1.T, &p2.T).Mul(&C, &CurveParams.D)
	D.Set(&p1.Z)
	tmp.Add(&p1.X, &p1.Y)
	E.Add(&p2.X, &p2.Y).
		Mul(&E, &tmp).
		Sub(&E, &A).
		Sub(&E, &B)
	F.Sub(&D, &C)
	G.Add(&D, &C)
	H.Set(&A)

	// mulBy5(&H)
	H.Neg(&H)
	gnarkfr.MulBy5(&H)

	H.Sub(&B, &H)

	p.X.Mul(&E, &F)
	p.Y.Mul(&G, &H)
	p.T.Mul(&E, &H)
	p.Z.Mul(&F, &G)

	return p
}
