package banderwagon

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/crate-crypto/go-ipa/bandersnatch"
	"github.com/crate-crypto/go-ipa/bandersnatch/fp"
	"github.com/crate-crypto/go-ipa/bandersnatch/fr"
	"github.com/crate-crypto/go-ipa/common/parallel"
)

const (
	coordinateSize   = fp.Limbs * 8
	CompressedSize   = coordinateSize
	UncompressedSize = 2 * coordinateSize
)

// Fr is the scalar field underlying the group.
type Fr = fr.Element

// Generator is the generator of the group.
var Generator = Element{inner: bandersnatch.PointProj{
	X: bandersnatch.CurveParams.Base.X,
	Y: bandersnatch.CurveParams.Base.Y,
	Z: fp.One(),
}}

// Identity is the identity element of the group.
var Identity = Element{inner: bandersnatch.PointProj{
	X: fp.Zero(),
	Y: fp.One(),
	Z: fp.One(),
}}

// Element is an element of the group.
type Element struct {
	inner bandersnatch.PointProj
}

// Bytes returns the compressed serialized version of the element.
func (p Element) Bytes() [CompressedSize]byte {
	// Serialisation takes the x co-ordinate and multiplies it by the sign of y.
	affineX := p.inner.X
	affineY := p.inner.Y
	if !p.inner.Z.IsOne() {
		// Convert underlying point to affine representation.
		var affine bandersnatch.PointAffine
		affine.FromProj(&p.inner)
		affineX = affine.X
		affineY = affine.Y
	}

	if !affineY.LexicographicallyLargest() {
		affineX.Neg(&affineX)
	}
	return affineX.Bytes()
}

// BytesUncompressed returns the uncompressed serialized version of the element.
func (p Element) BytesUncompressed() [UncompressedSize]byte {
	// Convert underlying point to affine representation
	var affine bandersnatch.PointAffine
	affine.FromProj(&p.inner)

	xbytes := affine.X.Bytes()
	ybytes := affine.Y.Bytes()

	var xy [UncompressedSize]byte
	copy(xy[:], xbytes[:])
	copy(xy[coordinateSize:], ybytes[:])

	return xy
}

// BatchNormalize normalizes a slice of group elements.
func BatchNormalize(elements []*Element) error {
	// The elements slice might contain duplicate pointers,
	// dedupe them to avoid double work.
	mapDedupedElements := make(map[*Element]struct{}, len(elements))
	for _, e := range elements {
		mapDedupedElements[e] = struct{}{}
	}
	dedupedElements := make([]*Element, 0, len(mapDedupedElements))
	for e := range mapDedupedElements {
		dedupedElements = append(dedupedElements, e)
	}

	invs := make([]fp.Element, len(elements))
	accumulator := fp.One()

	// batch invert all points[].Z coordinates with Montgomery batch inversion trick
	// (stores points[].Z^-1 in result[i].X to avoid allocating a slice of fr.Elements)
	for i := 0; i < len(dedupedElements); i++ {
		if dedupedElements[i].inner.Z.IsZero() {
			return errors.New("can not normalize point at infinity")
		}
		invs[i] = accumulator
		accumulator.Mul(&accumulator, &dedupedElements[i].inner.Z)
	}

	var accInverse fp.Element
	accInverse.Inverse(&accumulator)

	for i := len(dedupedElements) - 1; i >= 0; i-- {
		invs[i].Mul(&invs[i], &accInverse)
		accInverse.Mul(&accInverse, &dedupedElements[i].inner.Z)
	}

	// batch convert to affine.
	parallel.Execute(len(dedupedElements), func(start, end int) {
		for i := start; i < end; i++ {
			dedupedElements[i].inner.X.Mul(&dedupedElements[i].inner.X, &invs[i])
			dedupedElements[i].inner.Y.Mul(&dedupedElements[i].inner.Y, &invs[i])
			dedupedElements[i].inner.Z = fp.One()
		}
	})
	return nil
}

// ElementsToBytes serialises a slice of group elements in compressed form.
func ElementsToBytes(elements ...*Element) [][CompressedSize]byte {
	// Collect all z co-ordinates
	zs := make([]fp.Element, len(elements))
	for i := 0; i < len(elements); i++ {
		zs[i] = elements[i].inner.Z
	}

	// Invert z co-ordinates
	zInvs := fp.BatchInvert(zs)

	serialised_points := make([][CompressedSize]byte, len(elements))

	// Multiply x and y by zInv
	for i := 0; i < len(elements); i++ {
		var X fp.Element
		var Y fp.Element

		element := elements[i]

		X.Mul(&element.inner.X, &zInvs[i])
		Y.Mul(&element.inner.Y, &zInvs[i])

		// Serialisation takes the x co-ordinate and multiplies it by the sign of y
		if !Y.LexicographicallyLargest() {
			X.Neg(&X)
		}

		serialised_points[i] = X.Bytes()
	}

	return serialised_points
}

// BatchToBytesUncompressed serialises a slice of group elements in uncompressed form.
func BatchToBytesUncompressed(elements ...*Element) [][UncompressedSize]byte {
	// Collect all z co-ordinates
	zs := make([]fp.Element, len(elements))
	for i := 0; i < len(elements); i++ {
		zs[i] = elements[i].inner.Z
	}

	// Invert z co-ordinates
	zInvs := fp.BatchInvert(zs)

	uncompressedPoints := make([][UncompressedSize]byte, len(elements))

	// Multiply x and y by zInv
	for i := 0; i < len(elements); i++ {
		var X fp.Element
		var Y fp.Element

		element := elements[i]

		X.Mul(&element.inner.X, &zInvs[i])
		Y.Mul(&element.inner.Y, &zInvs[i])

		xbytes := X.Bytes()
		ybytes := Y.Bytes()
		copy(uncompressedPoints[i][:], xbytes[:])
		copy(uncompressedPoints[i][coordinateSize:], ybytes[:])
	}

	return uncompressedPoints
}

func (p *Element) setBytes(buf []byte, trusted bool) error {
	if len(buf) != CompressedSize {
		return errors.New("invalid compressed point size")
	}

	// set the buffer which is x * SignY as X
	var x fp.Element
	if err := x.SetBytesCanonical(buf); err != nil {
		return fmt.Errorf("invalid compressed point: %s", err)
	}

	point := bandersnatch.GetPointFromX(&x, true)
	if point == nil {
		return errors.New("point is not on the curve")
	}

	// If the source isn't trusted, we do the subgroup check.
	if !trusted {
		err := subgroupCheck(x)
		if err != nil {
			return err
		}
	}

	// We have a valid point, set it.
	*p = Element{inner: bandersnatch.PointProj{
		X: point.X,
		Y: point.Y,
		Z: fp.One(),
	}}

	return nil
}

// SetBytes deserializes a compressed group element from buf.
// This method does all the proper checks assuming the bytes come from an
// untrusted source.
func (p *Element) SetBytes(buf []byte) error {
	return p.setBytes(buf, false)
}

// SetBytesUnsafe deserializes a compressed group element from buf.
// **DO NOT** use this method if the bytes comes from an untrusted source.
func (p *Element) SetBytesUnsafe(buf []byte) error {
	return p.setBytes(buf, true)
}

// SetBytesUncompressed deserializes an uncompressed group element from buf.
// This method does all the proper checks assuming the bytes come from an
// untrusted source.
func (p *Element) SetBytesUncompressed(buf []byte, trusted bool) error {
	if len(buf) != UncompressedSize {
		return errors.New("invalid uncompressed point size")
	}

	var x fp.Element
	x.SetBytes(buf[:coordinateSize])

	// subgroup check
	if !trusted {
		err := subgroupCheck(x)
		if err != nil {
			return err
		}
	}

	var y fp.Element
	y.SetBytes(buf[coordinateSize:])

	*p = Element{inner: bandersnatch.PointProj{
		X: x,
		Y: y,
		Z: fp.One(),
	}}

	return nil
}

// computes X/Y
func (p Element) mapToBaseField() fp.Element {
	var res fp.Element
	res.Div(&p.inner.X, &p.inner.Y)
	return res
}

// MapToScalarField maps a group element to the scalar field.
func (p Element) MapToScalarField(res *fr.Element) {
	basefield := p.mapToBaseField()
	baseFieldBytes := fp.BytesLE(basefield)

	res.SetBytesLE(baseFieldBytes[:])
}

// BatchMapToScalarField maps a slice of group elements to the scalar field.
func BatchMapToScalarField(result []*fr.Element, elements []*Element) error {
	if len(result) != len(elements) {
		return errors.New("result and elements slices must be the same length")
	}

	// Collect all y co-ordinates
	ys := make([]fp.Element, len(elements))
	for i := 0; i < len(elements); i++ {
		ys[i] = elements[i].inner.Y
	}

	// Invert y co-ordinates
	yInvs := fp.BatchInvert(ys)

	// Multiply x by yInv
	for i := 0; i < len(elements); i++ {
		var mappedElement fp.Element

		mappedElement.Mul(&elements[i].inner.X, &yInvs[i])
		byts := fp.BytesLE(mappedElement)
		result[i].SetBytesLE(byts[:])
	}

	return nil
}

// Equal returns true if p and other represent the same point.
func (p *Element) Equal(other *Element) bool {
	x1 := p.inner.X
	y1 := p.inner.Y

	x2 := other.inner.X
	y2 := other.inner.Y

	if x1.IsZero() && y1.IsZero() {
		return false
	}
	if x2.IsZero() && y2.IsZero() {
		return false
	}

	// Recall that the equality check for Banderwagon has to test
	// the equivalence class {(x, y), (-x, -y)}, thus check: x1*y2 == x2*y2.
	// Note that both points being in projective form doesn't change the check,
	// since the z1 and z2 terms cancel out.
	var lhs fp.Element
	var rhs fp.Element
	lhs.Mul(&x1, &y2)
	rhs.Mul(&y1, &x2)

	return lhs.Equal(&rhs)
}

func subgroupCheck(x fp.Element) error {
	// Check that  (1 - ax^2) is a square, if not abort.
	var res, one, ax_sq fp.Element
	one.SetOne()
	ax_sq.Square(&x)
	ax_sq.Mul(&ax_sq, &bandersnatch.CurveParams.A)
	res.Sub(&one, &ax_sq)
	if res.Legendre() <= 0 {
		return errors.New("point is not in the correct subgroup")
	}
	return nil
}

// SetIdentity sets p to the identity element.
func (p *Element) SetIdentity() *Element {
	*p = Identity
	return p
}

// Double sets p to 2*p1.
func (p *Element) Double(p1 *Element) *Element {
	p.inner.Double(&p1.inner)
	return p
}

// Add sets p to p1+p2.
func (p *Element) Add(p1, p2 *Element) *Element {
	p.inner.Add(&p1.inner, &p2.inner)
	return p
}

// AddMixed sets p to p1+p2, where p2 is in affine form.
func (p *Element) AddMixed(p1 *Element, p2 bandersnatch.PointAffine) *Element {
	p.inner.MixedAdd(&p1.inner, &p2)
	return p
}

// Sub sets p to p1-p2.
func (p *Element) Sub(p1, p2 *Element) *Element {
	var neg_p2 Element
	neg_p2.Neg(p2)

	return p.Add(p1, &neg_p2)
}

// IsOnCurve returns true if p is on the curve.
func (p *Element) IsOnCurve() bool {
	// TODO: use projective curve equation to check
	var point_aff bandersnatch.PointAffine
	point_aff.FromProj(&p.inner)
	return point_aff.IsOnCurve()
}

// Normalize returns a point in affine form.
// If the point is at infinity, returns an error.
func (p *Element) Normalize() error {
	if p.inner.Z.IsZero() {
		return errors.New("can not normalize point at infinity")
	}

	var point_aff bandersnatch.PointAffine
	point_aff.FromProj(&p.inner)

	p.inner.X.Set(&point_aff.X)
	p.inner.Y.Set(&point_aff.Y)
	p.inner.Z.SetOne()

	return nil
}

// Set sets p to p1.
func (p *Element) Set(p1 *Element) *Element {
	p.inner.X.Set(&p1.inner.X)
	p.inner.Y.Set(&p1.inner.Y)
	p.inner.Z.Set(&p1.inner.Z)
	return p
}

// Neg sets p to -p1.
func (p *Element) Neg(p1 *Element) *Element {
	p.inner.Neg(&p1.inner)
	return p
}

// ScalarMul sets p to p1*s.
func (p *Element) ScalarMul(p1 *Element, scalarMont *fr.Element) *Element {
	var bigScalar big.Int
	scalarMont.ToBigIntRegular(&bigScalar)
	p.inner.ScalarMultiplication(&p1.inner, &bigScalar)
	return p
}
