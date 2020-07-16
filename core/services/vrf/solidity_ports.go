package vrf

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"

	"go.dedis.ch/kyber/v3"
)

// This file contains golang re-implementations of functions on the VRF solidity
// contract. They are used to verify correct operation of those functions, and
// also to efficiently compute zInv off-chain, which makes computing the linear
// combination of c*gamma+s*hash onchain much more efficient.

type fieldElt = *big.Int

// neg(f) is the negation of f in the base field
func neg(f fieldElt) fieldElt { return sub(FieldSize, f) }

// projectiveSub(x1, z1, x2, z2) is the projective coordinates of x1/z1 - x2/z2
func projectiveSub(x1, z1, x2, z2 fieldElt) (fieldElt, fieldElt) {
	num1 := mul(z2, x1)
	num2 := neg(mul(z1, x2))
	return mod(add(num1, num2), FieldSize), mod(mul(z1, z2), FieldSize)
}

// projectiveMul(x1, z1, x2, z2) is projective coordinates of (x1/z1)×(x2/z2)
func projectiveMul(x1, z1, x2, z2 fieldElt) (fieldElt, fieldElt) {
	return mul(x1, x2), mul(z1, z2)
}

// ProjectiveECAdd(px, py, qx, qy) duplicates the calculation in projective
// coordinates of VRF.sol#projectiveECAdd, so we can reliably get the
// denominator (i.e, z)
func ProjectiveECAdd(p, q kyber.Point) (x, y, z fieldElt) {
	px, py := secp256k1.Coordinates(p)
	qx, qy := secp256k1.Coordinates(q)
	pz, qz := one, one
	lx := sub(qy, py)
	lz := sub(qx, px)

	sx, dx := projectiveMul(lx, lz, lx, lz)
	sx, dx = projectiveSub(sx, dx, px, pz)
	sx, dx = projectiveSub(sx, dx, qx, qz)

	sy, dy := projectiveSub(px, pz, sx, dx)
	sy, dy = projectiveMul(sy, dy, lx, lz)
	sy, dy = projectiveSub(sy, dy, py, pz)

	var sz fieldElt
	if dx != dy {
		sx = mul(sx, dy)
		sy = mul(sy, dx)
		sz = mul(dx, dy)
	} else {
		sz = dx
	}
	return mod(sx, FieldSize), mod(sy, FieldSize), mod(sz, FieldSize)
}
