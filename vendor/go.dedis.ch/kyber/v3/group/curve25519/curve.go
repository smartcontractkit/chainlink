package curve25519

import (
	"crypto/cipher"
	"crypto/sha512"
	"errors"
	"fmt"
	"math/big"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/mod"
	"go.dedis.ch/kyber/v3/util/random"
)

var zero = big.NewInt(0)
var one = big.NewInt(1)

// Extension of Point interface for elliptic curve X,Y coordinate access
type point interface {
	kyber.Point

	initXY(x, y *big.Int, curve kyber.Group)

	getXY() (x, y *mod.Int)
}

// Generic "kyber.base class" for Edwards curves,
// embodying functionality independent of internal Point representation.
type curve struct {
	self      kyber.Group // "Self pointer" for derived class
	Param                 // Twisted Edwards curve parameters
	zero, one mod.Int     // Constant ModInts with correct modulus
	a, d      mod.Int     // Curve equation parameters as ModInts
	full      bool        // True if we're using the full group

	order  mod.Int // Order of appropriate subgroup as a ModInt
	cofact mod.Int // Group's cofactor as a ModInt

	null kyber.Point // Identity point for this group
}

func (c *curve) String() string {
	if c.full {
		return c.Param.String() + "-full"
	}
	return c.Param.String()
}

func (c *curve) IsPrimeOrder() bool {
	return !c.full
}

// Returns the size in bytes of an encoded Scalar for this curve.
func (c *curve) ScalarLen() int {
	return (c.order.V.BitLen() + 7) / 8
}

// Create a new Scalar for this curve.
func (c *curve) Scalar() kyber.Scalar {
	return mod.NewInt64(0, &c.order.V)
}

// Returns the size in bytes of an encoded Point on this curve.
// Uses compressed representation consisting of the y-coordinate
// and only the sign bit of the x-coordinate.
func (c *curve) PointLen() int {
	return (c.P.BitLen() + 7 + 1) / 8
}

// NewKey returns a formatted curve25519 key (avoiding subgroup attack by requiring
// it to be a multiple of 8). NewKey implements the kyber/util/key.Generator interface.
func (c *curve) NewKey(stream cipher.Stream) kyber.Scalar {
	var buffer [32]byte
	random.Bytes(buffer[:], stream)
	scalar := sha512.Sum512(buffer[:])
	scalar[0] &= 248
	scalar[31] &= 127
	scalar[31] |= 64

	secret := c.Scalar().SetBytes(scalar[:32])
	return secret
}

// Initialize a twisted Edwards curve with given parameters.
// Caller passes pointers to null and base point prototypes to be initialized.
func (c *curve) init(self kyber.Group, p *Param, fullGroup bool,
	null, base point) *curve {
	c.self = self
	c.Param = *p
	c.full = fullGroup
	c.null = null

	// Edwards curve parameters as ModInts for convenience
	c.a.Init(&p.A, &p.P)
	c.d.Init(&p.D, &p.P)

	// Cofactor
	c.cofact.Init64(int64(p.R), &c.P)

	// Determine the modulus for scalars on this curve.
	// Note that we do NOT initialize c.order with Init(),
	// as that would normalize to the modulus, resulting in zero.
	// Just to be sure it's never used, we leave c.order.M set to nil.
	// We want it to be in a ModInt so we can pass it to P.Mul(),
	// but the scalar's modulus isn't needed for point multiplication.
	if fullGroup {
		// Scalar modulus is prime-order times the ccofactor
		c.order.V.SetInt64(int64(p.R)).Mul(&c.order.V, &p.Q)
	} else {
		c.order.V.Set(&p.Q) // Prime-order subgroup
	}

	// Useful ModInt constants for this curve
	c.zero.Init64(0, &c.P)
	c.one.Init64(1, &c.P)

	// Identity element is (0,1)
	null.initXY(zero, one, self)

	// Base point B
	var bx, by *big.Int
	if !fullGroup {
		bx, by = &p.PBX, &p.PBY
	} else {
		bx, by = &p.FBX, &p.FBY
		base.initXY(&p.FBX, &p.FBY, self)
	}
	if by.Sign() == 0 {
		// No standard base point was defined, so pick one.
		// Find the lowest-numbered y-coordinate that works.
		//println("Picking base point:")
		var x, y mod.Int
		for y.Init64(2, &c.P); ; y.Add(&y, &c.one) {
			if !c.solveForX(&x, &y) {
				continue // try another y
			}
			if c.coordSign(&x) != 0 {
				x.Neg(&x) // try positive x first
			}
			base.initXY(&x.V, &y.V, self)
			if c.validPoint(base) {
				break // got one
			}
			x.Neg(&x) // try -bx
			if c.validPoint(base) {
				break // got one
			}
		}
		//println("BX: "+x.V.String())
		//println("BY: "+y.V.String())
		bx, by = &x.V, &y.V
	}
	base.initXY(bx, by, self)

	// Sanity checks
	if !c.validPoint(null) {
		panic("invalid identity point " + null.String())
	}
	if !c.validPoint(base) {
		panic("invalid base point " + base.String())
	}

	return c
}

// Test the sign of an x or y coordinate.
// We use the least-significant bit of the coordinate as the sign bit.
func (c *curve) coordSign(i *mod.Int) uint {
	return i.V.Bit(0)
}

// Convert a point to string representation.
func (c *curve) pointString(x, y *mod.Int) string {
	return fmt.Sprintf("(%s,%s)", x.String(), y.String())
}

// Encode an Edwards curve point.
// We use little-endian encoding for consistency with Ed25519.
func (c *curve) encodePoint(x, y *mod.Int) []byte {

	// Encode the y-coordinate
	b, _ := y.MarshalBinary()

	// Encode the sign of the x-coordinate.
	if y.M.BitLen()&7 == 0 {
		// No unused bits at the top of y-coordinate encoding,
		// so we must prepend a whole byte.
		b = append(make([]byte, 1), b...)
	}
	if c.coordSign(x) != 0 {
		b[0] |= 0x80
	}

	// Convert to little-endian
	reverse(b, b)
	return b
}

// Decode an Edwards curve point into the given x,y coordinates.
// Returns an error if the input does not denote a valid curve point.
// Note that this does NOT check if the point is in the prime-order subgroup:
// an adversary could create an encoding denoting a point
// on the twist of the curve, or in a larger subgroup.
// However, the "safecurves" criteria (http://safecurves.cr.yp.to)
// ensure that none of these other subgroups are small
// other than the tiny ones represented by the cofactor;
// hence Diffie-Hellman exchange can be done without subgroup checking
// without exposing more than the least-significant bits of the scalar.
func (c *curve) decodePoint(bb []byte, x, y *mod.Int) error {

	// Convert from little-endian
	//fmt.Printf("decoding:\n%s\n", hex.Dump(bb))
	b := make([]byte, len(bb))
	reverse(b, bb)

	// Extract the sign of the x-coordinate
	xsign := uint(b[0] >> 7)
	b[0] &^= 0x80

	// Extract the y-coordinate
	y.V.SetBytes(b)
	y.M = &c.P

	// Compute the corresponding x-coordinate
	if !c.solveForX(x, y) {
		return errors.New("invalid elliptic curve point")
	}
	if c.coordSign(x) != xsign {
		x.Neg(x)
	}

	return nil
}

// Given a y-coordinate, solve for the x-coordinate on the curve,
// using the characteristic equation rewritten as:
//
//	x^2 = (1 - y^2)/(a - d*y^2)
//
// Returns true on success,
// false if there is no x-coordinate corresponding to the chosen y-coordinate.
//
func (c *curve) solveForX(x, y *mod.Int) bool {
	var yy, t1, t2 mod.Int

	yy.Mul(y, y)                     // yy = y^2
	t1.Sub(&c.one, &yy)              // t1 = 1 - y^-2
	t2.Mul(&c.d, &yy).Sub(&c.a, &t2) // t2 = a - d*y^2
	t2.Div(&t1, &t2)                 // t2 = x^2
	return x.Sqrt(&t2)               // may fail if not a square
}

// Test if a supposed point is on the curve,
// by checking the characteristic equation for Edwards curves:
//
//	a*x^2 + y^2 = 1 + d*x^2*y^2
//
func (c *curve) onCurve(x, y *mod.Int) bool {
	var xx, yy, l, r mod.Int

	xx.Mul(x, x) // xx = x^2
	yy.Mul(y, y) // yy = y^2

	l.Mul(&c.a, &xx).Add(&l, &yy) // l = a*x^2 + y^2
	r.Mul(&c.d, &xx).Mul(&r, &yy).Add(&c.one, &r)
	// r = 1 + d*x^2*y^2
	return l.Equal(&r)
}

// Sanity-check a point to ensure that it is on the curve
// and within the appropriate subgroup.
func (c *curve) validPoint(P point) bool {

	// Check on-curve
	x, y := P.getXY()
	if !c.onCurve(x, y) {
		return false
	}

	// Check in-subgroup by multiplying by subgroup order
	Q := c.self.Point()
	Q.Mul(&c.order, P)
	if !Q.Equal(c.null) {
		return false
	}

	return true
}

// Return number of bytes that can be embedded into points on this curve.
func (c *curve) embedLen() int {
	// Reserve at least 8 most-significant bits for randomness,
	// and the least-significant 8 bits for embedded data length.
	// (Hopefully it's unlikely we'll need >=2048-bit curves soon.)
	return (c.P.BitLen() - 8 - 8) / 8
}

// Pick a [pseudo-]random curve point with optional embedded data,
// filling in the point's x,y coordinates
func (c *curve) embed(P point, data []byte, rand cipher.Stream) {

	// How much data to embed?
	dl := c.embedLen()
	if dl > len(data) {
		dl = len(data)
	}

	// Retry until we find a valid point
	var x, y mod.Int
	var Q kyber.Point
	for {
		// Get random bits the size of a compressed Point encoding,
		// in which the topmost bit is reserved for the x-coord sign.
		l := c.PointLen()
		b := make([]byte, l)
		rand.XORKeyStream(b, b) // Interpret as little-endian
		if data != nil {
			b[0] = byte(dl)       // Encode length in low 8 bits
			copy(b[1:1+dl], data) // Copy in data to embed
		}
		reverse(b, b) // Convert to big-endian form

		xsign := b[0] >> 7                    // save x-coordinate sign bit
		b[0] &^= 0xff << uint(c.P.BitLen()&7) // clear high bits

		y.M = &c.P // set y-coordinate
		y.SetBytes(b)

		if !c.solveForX(&x, &y) { // Corresponding x-coordinate?
			continue // none, retry
		}

		// Pick a random sign for the x-coordinate
		if c.coordSign(&x) != uint(xsign) {
			x.Neg(&x)
		}

		// Initialize the point
		P.initXY(&x.V, &y.V, c.self)
		if c.full {
			// If we're using the full group,
			// we just need any point on the curve, so we're done.
			return
		}

		// We're using the prime-order subgroup,
		// so we need to make sure the point is in that subgroup.
		// If we're not trying to embed data,
		// we can convert our point into one in the subgroup
		// simply by multiplying it by the cofactor.
		if data == nil {
			P.Mul(&c.cofact, P) // multiply by cofactor
			if P.Equal(c.null) {
				continue // unlucky; try again
			}
			return
		}

		// Since we need the point's y-coordinate to make sense,
		// we must simply check if the point is in the subgroup
		// and retry point generation until it is.
		if Q == nil {
			Q = c.self.Point()
		}
		Q.Mul(&c.order, P)
		if Q.Equal(c.null) {
			return
		}

		// Keep trying...
	}
}

// Extract embedded data from a point group element,
// or an error if embedded data is invalid or not present.
func (c *curve) data(x, y *mod.Int) ([]byte, error) {
	b := c.encodePoint(x, y)
	dl := int(b[0])
	if dl > c.embedLen() {
		return nil, errors.New("invalid embedded data length")
	}
	return b[1 : 1+dl], nil
}

// reverse copies src into dst in byte-reversed order and returns dst,
// such that src[0] goes into dst[len-1] and vice versa.
// dst and src may be the same slice but otherwise must not overlap.
func reverse(dst, src []byte) []byte {
	l := len(dst)
	for i, j := 0, l-1; i < (l+1)/2; {
		dst[i], dst[j] = src[j], src[i]
		i++
		j--
	}
	return dst
}
