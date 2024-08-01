package common

import (
	"fmt"
	"io"

	"github.com/crate-crypto/go-ipa/bandersnatch/fr"
	"github.com/crate-crypto/go-ipa/banderwagon"
)

// VectorLength is the number of elements in the vector. This value is fixed.
// Note that this means that the degree of the polynomial is one less than this value.
const VectorLength = 256

// Returns powers of x from 0 to degree-1
// <1, x, x^2, x^3, x^4,...,x^(degree-1)>
// TODO This method is used in two places; one is to evaluate a polynomial (test), and the other is to
// TODO compute powers of challenges.
// TODO the first one we can use the bls package for
// TODO The second we _could_ just multiply on each iteration, (depends on how readable it is)
func PowersOf(x fr.Element, degree int) []fr.Element {
	result := make([]fr.Element, degree)
	result[0] = fr.One()

	for i := 1; i < degree; i++ {
		result[i].Mul(&result[i-1], &x)
	}

	return result
}

func ReadPoint(r io.Reader) (*banderwagon.Element, error) {
	var x = make([]byte, 32)
	if _, err := io.ReadAtLeast(r, x, 32); err != nil {
		return nil, fmt.Errorf("reading x coordinate: %w", err)
	}
	var p = &banderwagon.Element{}
	if err := p.SetBytes(x); err != nil {
		return nil, fmt.Errorf("deserializing point: %w", err)
	}
	return p, nil
}

func ReadScalar(r io.Reader) (*fr.Element, error) {
	var x = make([]byte, 32)
	if _, err := io.ReadAtLeast(r, x, 32); err != nil {
		return nil, fmt.Errorf("reading scalar: %w", err)
	}
	var scalar = &fr.Element{}
	if _, err := scalar.SetBytesLECanonical(x); err != nil {
		return nil, fmt.Errorf("deserializing scalar: %s", err)
	}

	return scalar, nil
}
