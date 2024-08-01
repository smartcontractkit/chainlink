package kzg

import (
	"math/big"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

// In this file we implement a simple version of the fft algorithm
// without any optimizations. This is sufficient as the fft algorithm is
// not on the hot path; we only need it to compute the lagrange version
// of the SRS, this can be done once at startup. Even if not cached,
// this process takes two to three seconds.
//
// See: https://faculty.sites.iastate.edu/jia/files/inline-files/polymultiply.pdf
// for a reference.

// Computes an FFT (Fast Fourier Transform) of the G1 elements.
//
// The elements are returned in order as opposed to being returned in
// bit-reversed order.
func (domain *Domain) FftG1(values []bls12381.G1Affine) []bls12381.G1Affine {
	return fftG1(values, domain.Generator)
}

// Computes an IFFT(Inverse Fast Fourier Transform) of the G1 elements.
//
// The elements are returned in order as opposed to being returned in
// bit-reversed order.
func (domain *Domain) IfftG1(values []bls12381.G1Affine) []bls12381.G1Affine {
	var invDomainBI big.Int
	domain.CardinalityInv.BigInt(&invDomainBI)

	inverseFFT := fftG1(values, domain.GeneratorInv)

	// scale by the inverse of the domain size
	for i := 0; i < len(inverseFFT); i++ {
		inverseFFT[i].ScalarMultiplication(&inverseFFT[i], &invDomainBI)
	}

	return inverseFFT
}

// fftG1 computes an FFT (Fast Fourier Transform) of the G1 elements.
//
// This is the actual implementation of [FftG1] with the same convention.
// That is, the returned slice is in "normal", rather than bit-reversed order.
// We assert that values is a slice of length n==2^i and nthRootOfUnity is a primitive n'th root of unity.
func fftG1(values []bls12381.G1Affine, nthRootOfUnity fr.Element) []bls12381.G1Affine {
	n := len(values)
	if n == 1 {
		return values
	}

	var generatorSquared fr.Element
	generatorSquared.Square(&nthRootOfUnity) // generator with order n/2

	// split the input slice into a (copy of) the values at even resp. odd indices.
	even, odd := takeEvenOdd(values)

	// perform FFT recursively on those parts.
	fftEven := fftG1(even, generatorSquared)
	fftOdd := fftG1(odd, generatorSquared)

	// combine them to get the result
	// - evaluations[k] = fftEven[k] + w^k * fftOdd[k]
	// - evaluations[k] = fftEven[k] - w^k * fftOdd[k]
	// where w is a n'th primitive root of unity.
	inputPoint := fr.One()
	evaluations := make([]bls12381.G1Affine, n)
	for k := 0; k < n/2; k++ {
		var tmp bls12381.G1Affine

		var inputPointBI big.Int
		inputPoint.BigInt(&inputPointBI)

		if inputPoint.IsOne() {
			tmp.Set(&fftOdd[k])
		} else {
			tmp.ScalarMultiplication(&fftOdd[k], &inputPointBI)
		}

		evaluations[k].Add(&fftEven[k], &tmp)
		evaluations[k+n/2].Sub(&fftEven[k], &tmp)

		// we could take this from precomputed values in Domain (as domain.roots[n*k]), but then we would need to pass the domain.
		// At any rate, we don't really need to optimize here.
		inputPoint.Mul(&inputPoint, &nthRootOfUnity)
	}

	return evaluations
}

// takeEvenOdd Takes a slice and return two slices
// The first slice contains (a copy of) all of the elements
// at even indices, the second slice contains
// (a copy of) all of the elements at odd indices
//
// We assume that the length of the given values slice is even
// so the returned arrays will be the same length.
// This is the case for a radix-2 FFT
func takeEvenOdd[T interface{}](values []T) ([]T, []T) {
	n := len(values)
	even := make([]T, 0, n/2)
	odd := make([]T, 0, n/2)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			even = append(even, values[i])
		} else {
			odd = append(odd, values[i])
		}
	}

	return even, odd
}
