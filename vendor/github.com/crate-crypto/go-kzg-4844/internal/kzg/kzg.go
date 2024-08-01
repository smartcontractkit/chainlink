package kzg

import (
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

// A polynomial in lagrange form
type Polynomial = []fr.Element

// A commitment to a polynomial
// Excluding tests, this will be produced
// by committing to a polynomial in lagrange form
type Commitment = bls12381.G1Affine
