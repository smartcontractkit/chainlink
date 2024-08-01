package kzg

import (
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/crate-crypto/go-kzg-4844/internal/multiexp"
)

// OpeningKey is the key used to verify opening proofs
type OpeningKey struct {
	// This is the degree-0 G_1 element in the trusted setup.
	// In the specs, this is denoted as `KZG_SETUP_G1[0]`
	GenG1 bls12381.G1Affine
	// This is the degree-0 G_2 element in the trusted setup.
	// In the specs, this is denoted as `KZG_SETUP_G2[0]`
	GenG2 bls12381.G2Affine
	// This is the degree-1 G_2 element in the trusted setup.
	// In the specs, this is denoted as `KZG_SETUP_G2[1]`
	AlphaG2 bls12381.G2Affine
}

// CommitKey holds the data needed to commit to polynomials and by proxy make opening proofs
type CommitKey struct {
	// These are the G1 elements from the trusted setup.
	// In the specs this is denoted as `KZG_SETUP_G1` before
	// we processed it with `ifftG1`. Once we compute `ifftG1`
	// then this list is denoted as `KZG_SETUP_LAGRANGE` in the specs.
	G1 []bls12381.G1Affine
}

// ReversePoints applies the bit reversal permutation
// to the G1 points stored inside the CommitKey c.
func (c *CommitKey) ReversePoints() {
	bitReverse(c.G1)
}

// SRS holds the structured reference string (SRS) for making
// and verifying KZG proofs
//
// This codebase is only concerned with polynomials in Lagrange
// form, so we only expose methods to create the SRS in Lagrange form
//
// The monomial SRS methods are solely used for testing.
type SRS struct {
	CommitKey  CommitKey
	OpeningKey OpeningKey
}

// Commit commits to a polynomial using a multi exponentiation with the
// Commitment key.
//
// numGoRoutines is used to configure the amount of concurrency needed. Setting this
// value to a negative number or 0 will make it default to the number of CPUs.
func Commit(p Polynomial, ck *CommitKey, numGoRoutines int) (*Commitment, error) {
	if len(p) == 0 || len(p) > len(ck.G1) {
		return nil, ErrInvalidPolynomialSize
	}

	return multiexp.MultiExp(p, ck.G1[:len(p)], numGoRoutines)
}
