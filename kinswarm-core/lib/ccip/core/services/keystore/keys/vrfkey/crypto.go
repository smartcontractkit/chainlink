package vrfkey

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"go.dedis.ch/kyber/v3"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/hex"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
	bm "github.com/smartcontractkit/chainlink/v2/core/utils/big_math"
)

// This file contains golang re-implementations of functions on the VRF solidity
// contract. They are used to verify correct operation of those functions, and
// also to efficiently compute zInv off-chain, which makes computing the linear
// combination of c*gamma+s*hash onchain much more efficient.

var (
	// FieldSize is number of elements in secp256k1's base field, i.e. GF(FieldSize)
	FieldSize = mustParseBig(
		"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F",
	)
	Secp256k1Curve       = &secp256k1.Secp256k1{}
	Generator            = Secp256k1Curve.Point().Base()
	eulersCriterionPower = bm.Div(bm.Sub(FieldSize, bm.One), bm.Two)
	sqrtPower            = bm.Div(bm.Add(FieldSize, bm.One), bm.Four)
	ErrCGammaEqualsSHash = fmt.Errorf("pick a different nonce; c*gamma = s*hash, with this one")
	// hashToCurveHashPrefix is domain-separation tag for initial HashToCurve hash.
	// Corresponds to HASH_TO_CURVE_HASH_PREFIX in VRF.sol.
	hashToCurveHashPrefix = common.BigToHash(bm.One).Bytes()
	// scalarFromCurveHashPrefix is a domain-separation tag for the hash taken in
	// ScalarFromCurve. Corresponds to SCALAR_FROM_CURVE_POINTS_HASH_PREFIX in
	// VRF.sol.
	scalarFromCurveHashPrefix = common.BigToHash(bm.Two).Bytes()
	// RandomOutputHashPrefix is a domain-separation tag for the hash used to
	// compute the final VRF random output
	RandomOutputHashPrefix = common.BigToHash(bm.Three).Bytes()
)

type fieldElt = *big.Int

// neg(f) is the negation of f in the base field
func neg(f fieldElt) fieldElt { return bm.Sub(FieldSize, f) }

// projectiveSub(x1, z1, x2, z2) is the projective coordinates of x1/z1 - x2/z2
func projectiveSub(x1, z1, x2, z2 fieldElt) (fieldElt, fieldElt) {
	num1 := bm.Mul(z2, x1)
	num2 := neg(bm.Mul(z1, x2))
	return bm.Mod(bm.Add(num1, num2), FieldSize), bm.Mod(bm.Mul(z1, z2), FieldSize)
}

// projectiveMul(x1, z1, x2, z2) is projective coordinates of (x1/z1)×(x2/z2)
func projectiveMul(x1, z1, x2, z2 fieldElt) (fieldElt, fieldElt) {
	return bm.Mul(x1, x2), bm.Mul(z1, z2)
}

// ProjectiveECAdd(px, py, qx, qy) duplicates the calculation in projective
// coordinates of VRF.sol#projectiveECAdd, so we can reliably get the
// denominator (i.e, z)
func ProjectiveECAdd(p, q kyber.Point) (x, y, z fieldElt) {
	px, py := secp256k1.Coordinates(p)
	qx, qy := secp256k1.Coordinates(q)
	pz, qz := bm.One, bm.One
	lx := bm.Sub(qy, py)
	lz := bm.Sub(qx, px)

	sx, dx := projectiveMul(lx, lz, lx, lz)
	sx, dx = projectiveSub(sx, dx, px, pz)
	sx, dx = projectiveSub(sx, dx, qx, qz)

	sy, dy := projectiveSub(px, pz, sx, dx)
	sy, dy = projectiveMul(sy, dy, lx, lz)
	sy, dy = projectiveSub(sy, dy, py, pz)

	var sz fieldElt
	if dx != dy {
		sx = bm.Mul(sx, dy)
		sy = bm.Mul(sy, dx)
		sz = bm.Mul(dx, dy)
	} else {
		sz = dx
	}
	return bm.Mod(sx, FieldSize), bm.Mod(sy, FieldSize), bm.Mod(sz, FieldSize)
}

// IsSquare returns true iff x = y^2 for some y in GF(p)
func IsSquare(x *big.Int) bool {
	return bm.Equal(bm.One, bm.Exp(x, eulersCriterionPower, FieldSize))
}

// SquareRoot returns a s.t. a^2=x, as long as x is a square
func SquareRoot(x *big.Int) *big.Int {
	return bm.Exp(x, sqrtPower, FieldSize)
}

// YSquared returns x^3+7 mod fieldSize, the right-hand side of the secp256k1
// curve equation.
func YSquared(x *big.Int) *big.Int {
	return bm.Mod(bm.Add(bm.Exp(x, bm.Three, FieldSize), bm.Seven), FieldSize)
}

// IsCurveXOrdinate returns true iff there is y s.t. y^2=x^3+7
func IsCurveXOrdinate(x *big.Int) bool {
	return IsSquare(YSquared(x))
}

// FieldHash hashes xs uniformly into {0, ..., fieldSize-1}. msg is assumed to
// already be a 256-bit hash
func FieldHash(msg []byte) *big.Int {
	rv := utils.MustHash(string(msg)).Big()
	// Hash recursively until rv < q. P(success per iteration) >= 0.5, so
	// number of extra hashes is geometrically distributed, with mean < 1.
	for rv.Cmp(FieldSize) >= 0 {
		rv = utils.MustHash(string(common.BigToHash(rv).Bytes())).Big()
	}
	return rv
}

// linearCombination returns c*p1+s*p2
func linearCombination(c *big.Int, p1 kyber.Point,
	s *big.Int, p2 kyber.Point) kyber.Point {
	return Secp256k1Curve.Point().Add(
		Secp256k1Curve.Point().Mul(secp256k1.IntToScalar(c), p1),
		Secp256k1Curve.Point().Mul(secp256k1.IntToScalar(s), p2))
}

// checkCGammaNotEqualToSHash checks c*gamma ≠ s*hash, as required by solidity
// verifier
func checkCGammaNotEqualToSHash(c *big.Int, gamma kyber.Point, s *big.Int,
	hash kyber.Point) error {
	cGamma := Secp256k1Curve.Point().Mul(secp256k1.IntToScalar(c), gamma)
	sHash := Secp256k1Curve.Point().Mul(secp256k1.IntToScalar(s), hash)
	if cGamma.Equal(sHash) {
		return ErrCGammaEqualsSHash
	}
	return nil
}

// HashToCurve is a cryptographic hash function which outputs a secp256k1 point,
// or an error. It passes each candidate x ordinate to ordinates function.
func HashToCurve(p kyber.Point, input *big.Int, ordinates func(x *big.Int),
) (kyber.Point, error) {
	if !(secp256k1.ValidPublicKey(p) && input.BitLen() <= 256 && input.Cmp(bm.Zero) >= 0) {
		return nil, fmt.Errorf("bad input to vrf.HashToCurve")
	}
	x := FieldHash(append(hashToCurveHashPrefix, append(secp256k1.LongMarshal(p),
		utils.Uint256ToBytes32(input)...)...))
	ordinates(x)
	for !IsCurveXOrdinate(x) { // Hash recursively until x^3+7 is a square
		x.Set(FieldHash(common.BigToHash(x).Bytes()))
		ordinates(x)
	}
	y := SquareRoot(YSquared(x))
	rv := secp256k1.SetCoordinates(x, y)
	if bm.Equal(bm.I().Mod(y, bm.Two), bm.One) { // Negate response if y odd
		rv = rv.Neg(rv)
	}
	return rv, nil
}

// ScalarFromCurve returns a hash for the curve points. Corresponds to the
// hash computed in VRF.sol#ScalarFromCurvePoints
func ScalarFromCurvePoints(
	hash, pk, gamma kyber.Point, uWitness [20]byte, v kyber.Point) *big.Int {
	if !(secp256k1.ValidPublicKey(hash) && secp256k1.ValidPublicKey(pk) &&
		secp256k1.ValidPublicKey(gamma) && secp256k1.ValidPublicKey(v)) {
		panic("bad arguments to vrf.ScalarFromCurvePoints")
	}
	// msg will contain abi.encodePacked(hash, pk, gamma, v, uWitness)
	msg := scalarFromCurveHashPrefix
	for _, p := range []kyber.Point{hash, pk, gamma, v} {
		msg = append(msg, secp256k1.LongMarshal(p)...)
	}
	msg = append(msg, uWitness[:]...)
	return bm.I().SetBytes(utils.MustHash(string(msg)).Bytes())
}

func mustParseBig(hx string) *big.Int {
	n, err := hex.ParseBig(hx)
	if err != nil {
		panic(fmt.Errorf(`failed to convert "%s" as hex to big.Int`, hx))
	}
	return n
}
