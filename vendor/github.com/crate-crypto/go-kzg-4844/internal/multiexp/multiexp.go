package multiexp

import (
	"github.com/consensys/gnark-crypto/ecc"
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

// MultiExp computes a multi exponentiation -- That is, an inner product between points and scalars.
//
// More precisely, the result is set to scalars[0]*points[0] + ... + scalars[n-1]*points[n-1], where n is the length of both slices
// If the slices differ in length, this function returns an error.
//
// numGoRoutines is used to configure the amount of concurrency needed. Setting this
// value to a negative number or 0 will make it default to the number of CPUs.
//
// Returns an error if the numGoRoutines exceeds 1024.
//
// [g1_lincomb]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#g1_lincomb
func MultiExp(scalars []fr.Element, points []bls12381.G1Affine, numGoRoutines int) (*bls12381.G1Affine, error) {
	err := isValidNumGoRoutines(numGoRoutines)
	if err != nil {
		return nil, err
	}
	return new(bls12381.G1Affine).MultiExp(points, scalars, ecc.MultiExpConfig{NbTasks: numGoRoutines})
}

// isValidNumGoRoutines will return an error if the number
// of go routines to be used is not Valid.
//
// Valid meaning that is less than 1024.
//
// 1024 is chosen here as the underlying gnark-crypto library will
// return an error for more than 1024.
// Instead of waiting until the user tries to call an algorithm
// which requires numGoRoutines, we return the error here instead.
func isValidNumGoRoutines(value int) error {
	if value >= 1024 {
		return ErrTooManyGoRoutines
	}
	return nil
}
