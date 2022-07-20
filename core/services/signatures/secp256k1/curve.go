// Package secp256k1 is an implementation of the kyber.{Group,Point,Scalar}
////////////////////////////////////////////////////////////////////////////////
//       XXX: Do not use in production until this code has been audited.
////////////////////////////////////////////////////////////////////////////////
// interfaces, based on btcd/btcec and kyber/group/mod
//
// XXX: NOT CONSTANT TIME!
package secp256k1

import (
	"math/big"

	secp256k1BTCD "github.com/btcsuite/btcd/btcec"

	"go.dedis.ch/kyber/v3"
)

// Secp256k1 represents the secp256k1 group.
// There are no parameters and no initialization is required
// because it supports only this one specific curve.
type Secp256k1 struct{}

// s256 is the btcec representation of secp256k1.
var s256 *secp256k1BTCD.KoblitzCurve = secp256k1BTCD.S256()

// String returns the name of the curve
func (*Secp256k1) String() string { return "Secp256k1" }

var egScalar kyber.Scalar = newScalar(big.NewInt(0))
var egPoint kyber.Point = &secp256k1Point{newFieldZero(), newFieldZero()}

// ScalarLen returns the length of a marshalled Scalar
func (*Secp256k1) ScalarLen() int { return egScalar.MarshalSize() }

// Scalar creates a new Scalar for the prime-order group on the secp256k1 curve
func (*Secp256k1) Scalar() kyber.Scalar { return newScalar(big.NewInt(0)) }

// PointLen returns the length of a marshalled Point
func (*Secp256k1) PointLen() int { return egPoint.MarshalSize() }

// Point returns a new secp256k1 point
func (*Secp256k1) Point() kyber.Point {
	return &secp256k1Point{newFieldZero(), newFieldZero()}
}
