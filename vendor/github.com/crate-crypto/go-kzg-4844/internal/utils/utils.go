package utils

import (
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

// The spec includes a method to compute the modular inverse.
// This method is named .Inverse on `fr.Element`
// When the element to invert is zero, this method will return zero
// however note that this is not utilized in the specs anywhere
// and so it is also fine to panic on zero.
//
// [bls_modular_inverse]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#bls_modular_inverse
// [div]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#div

// ComputePowers computes x^0 to x^n-1.
//
// More precisely, given x and n, returns a slice containing [x^0, ..., x^n-1]
// In particular, for n==0, an empty slice is returned
//
// [compute_powers]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#compute_powers
func ComputePowers(x fr.Element, n uint) []fr.Element {
	if n == 0 {
		return []fr.Element{}
	}

	powers := make([]fr.Element, n)
	powers[0].SetOne()
	for i := uint(1); i < n; i++ {
		powers[i].Mul(&powers[i-1], &x)
	}

	return powers
}

// IsPowerOfTwo returns true if `value` is a power of two.
//
// `0` will return false
//
// [is_power_of_two]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#is_power_of_two
func IsPowerOfTwo(value uint64) bool {
	return value > 0 && (value&(value-1) == 0)
}

func ReduceCanonicalBigEndian(serScalar []byte) (fr.Element, error) {
	var scalar fr.Element
	err := scalar.SetBytesCanonical(serScalar)

	return scalar, err
}
