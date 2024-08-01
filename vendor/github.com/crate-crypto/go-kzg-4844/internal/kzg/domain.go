package kzg

import (
	"fmt"
	"math/big"
	"math/bits"

	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/crate-crypto/go-kzg-4844/internal/utils"
)

// Domain is a struct defining the set of points that polynomials are evaluated over.
// To enable efficient FFT-based algorithms, these points are chosen as 2^i'th roots of unity and we precompute and store
// certain values related to that inside the struct.
type Domain struct {
	// Size of the domain as a uint64. This must be a power of 2.
	// Since the base field has 2^i'th roots of unity for i<=32, Cardinality is <= 2^32)
	Cardinality uint64
	// Inverse of the size of the domain as
	// a field element. This is useful for
	// inverse FFTs.
	CardinalityInv fr.Element
	// Generator for the multiplicative subgroup
	// Not a primitive element (i.e. generator) for the *whole* field.
	//
	// This generator will have order equal to the
	// cardinality of the domain.
	Generator fr.Element
	// Inverse of the Generator. This is precomputed
	// and useful for inverse FFTs.
	GeneratorInv fr.Element

	// Roots of unity for the multiplicative subgroup
	// Note that these may or may not be in bit-reversed order.
	Roots []fr.Element

	// Precomputed inverses of the domain which
	// we will use to speed up the computation of
	// f(x)/g(x) where g(x) is a linear polynomial
	// which vanishes on a point on the domain
	PreComputedInverses []fr.Element
}

// NewDomain returns a new domain with the desired number of points x.
//
// We only support powers of 2 for x.
//
// Modified from [gnark-crypto].
//
// [gnark-crypto]: https://github.com/ConsenSys/gnark-crypto/blob/8f7ca09273c24ed9465043566906cbecf5dcee91/ecc/bls12-381/fr/fft/domain.go#L66
func NewDomain(x uint64) *Domain {
	if bits.OnesCount64(x) != 1 {
		panic(fmt.Sprintf("x (%d) is not a power of 2. This library only supports domain sizes that are powers of two", x))
	}
	domain := &Domain{}
	domain.Cardinality = x

	// Generator of the largest 2-adic subgroup.
	// This particular element has order 2^maxOrderRoot == 2^32.
	var rootOfUnity fr.Element
	_, err := rootOfUnity.SetString("10238227357739495823651030575849232062558860180284477541189508159991286009131")
	if err != nil {
		panic("failed to initialize root of unity")
	}
	const maxOrderRoot uint64 = 32

	// Find generator subgroup of order x.
	// This can be constructed by powering a generator of the largest 2-adic subgroup of order 2^32 by an exponent
	// of (2^32)/x, provided x is <= 2^32.
	logx := uint64(bits.TrailingZeros64(x))
	if logx > maxOrderRoot {
		panic(fmt.Sprintf("x (%d) is too big: the required root of unity does not exist", x))
	}
	expo := uint64(1 << (maxOrderRoot - logx))
	domain.Generator.Exp(rootOfUnity, big.NewInt(int64(expo))) // Domain.Generator has order x now.

	// Store Inverse of the generator and inverse of the domain size (as field elements).
	domain.GeneratorInv.Inverse(&domain.Generator)
	domain.CardinalityInv.SetUint64(x)
	domain.CardinalityInv.Inverse(&domain.CardinalityInv)

	// Compute all relevant roots of unity, i.e. the multiplicative subgroup of size x.
	domain.Roots = make([]fr.Element, x)
	current := fr.One()
	for i := uint64(0); i < x; i++ {
		domain.Roots[i] = current
		current.Mul(&current, &domain.Generator)
	}

	// Compute precomputed inverses: 1 / w^i
	// Note: domain.PreComputedInverses[i] == domain.Roots[x-i mod x], so
	// these are redundant, but simplify writing down some algorithms
	// and not deal with the case where the roots are bit-reversed.
	// We use BatchInvert instead of the above for clarity.
	domain.PreComputedInverses = fr.BatchInvert(domain.Roots)

	return domain
}

/*
Taken from a chat with Dr Dankrad Feist:
- Samples are going to be contiguous when we switch on full sharding.
- Technically there is nothing that requires samples to be contiguous
pieces of data, but it seems a lot nicer.
- also the relationship between original and interpolated data would
look really strange, with them being interleaved.
- Everything is just nice in brp and looks really strange in direct order
once you introduce sharding. So best to use it from the start and not have
to think about all these when you add DAS.
*/

// bitReverse applies the bit-reversal permutation to `list`.
// `len(list)` must be a power of 2
//
// This means that for post-state list output and pre-state list input,
// we have output[i] == input[bitreverse(i)], where bitreverse reverses the bit-pattern
// of i, interpreted as a log2(len(list))-bit integer.
//
// This is in no way needed for basic KZG and is included in this library as
// a stepping-stone to full Dank-sharding.
//
// Modified from [gnark-crypto].
//
// [gnark-crypto]: https://github.com/ConsenSys/gnark-crypto/blob/8f7ca09273c24ed9465043566906cbecf5dcee91/ecc/bls12-381/fr/fft/fft.go#L245
//
// [reverse_bits]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#reverse_bits
func bitReverse[K interface{}](list []K) {
	n := uint64(len(list))
	if !utils.IsPowerOfTwo(n) {
		panic("size of list given to bitReverse must be a power of two")
	}

	// The standard library's bits.Reverse64 inverts its input as a 64-bit unsigned integer.
	// However, we need to invert it as a log2(len(list))-bit integer, so we need to correct this by
	// shifting appropriately.
	shiftCorrection := uint64(64 - bits.TrailingZeros64(n))

	for i := uint64(0); i < n; i++ {
		// Find index irev, such that i and irev get swapped
		irev := bits.Reverse64(i) >> shiftCorrection
		if irev > i {
			list[i], list[irev] = list[irev], list[i]
		}
	}
}

// ReverseRoots applies the bit-reversal permutation to the list of precomputed roots of unity and their inverses in the domain.
//
// [bit_reversal_permutation]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#bit_reversal_permutation
func (domain *Domain) ReverseRoots() {
	bitReverse(domain.Roots)
	bitReverse(domain.PreComputedInverses)
}

// findRootIndex returns the index of the element in the domain or -1 if not found.
//
//   - If point is in the domain (meaning that point is a domain.Cardinality'th root of unity), returns the index of the point in the domain.
//   - If point is not in the domain, returns -1.
func (domain *Domain) findRootIndex(point fr.Element) int64 {
	for i := int64(0); i < int64(domain.Cardinality); i++ {
		if point.Equal(&domain.Roots[i]) {
			return i
		}
	}

	return -1
}

// EvaluateLagrangePolynomial evaluates a Lagrange polynomial at the given point of evaluation.
//
// The input polynomial is given in evaluation form, meaning a list of evaluations at the points in the domain.
// If len(poly) != domain.Cardinality, returns an error.
//
// [evaluate_polynomial_in_evaluation_form]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#evaluate_polynomial_in_evaluation_form
func (domain *Domain) EvaluateLagrangePolynomial(poly Polynomial, evalPoint fr.Element) (*fr.Element, error) {
	outputPoint, _, err := domain.evaluateLagrangePolynomial(poly, evalPoint)

	return outputPoint, err
}

// evaluateLagrangePolynomial is the implementation for [EvaluateLagrangePolynomial].
//
// It evaluates a Lagrange polynomial at the given point of evaluation and reports whether the given point was among the points of the domain:
//   - The input polynomial is given in evaluation form, that is, a list of evaluations at the points in the domain.
//   - The evaluationResult is the result of evaluation at evalPoint.
//   - indexInDomain is the index inside domain.Roots, if evalPoint is among them, -1 otherwise
//
// This semantics was copied from the go library, see: https://cs.opensource.google/go/x/exp/+/522b1b58:slices/slices.go;l=117
func (domain *Domain) evaluateLagrangePolynomial(poly Polynomial, evalPoint fr.Element) (*fr.Element, int64, error) {
	var indexInDomain int64 = -1

	if domain.Cardinality != uint64(len(poly)) {
		return nil, indexInDomain, ErrPolynomialMismatchedSizeDomain
	}

	// If the evaluation point is in the domain
	// then evaluation of the polynomial in lagrange form
	// is the same as indexing it with the position
	// that the evaluation point is in, in the domain
	indexInDomain = domain.findRootIndex(evalPoint)
	if indexInDomain != -1 {
		return &poly[indexInDomain], indexInDomain, nil
	}

	denom := make([]fr.Element, domain.Cardinality)
	for i := range denom {
		denom[i].Sub(&evalPoint, &domain.Roots[i])
	}
	invDenom := fr.BatchInvert(denom)

	var result fr.Element
	for i := 0; i < int(domain.Cardinality); i++ {
		var num fr.Element
		num.Mul(&poly[i], &domain.Roots[i])

		var div fr.Element
		div.Mul(&num, &invDenom[i])

		result.Add(&result, &div)
	}

	// result * (x^width - 1) * 1/width
	var tmp fr.Element
	tmp.Exp(evalPoint, big.NewInt(0).SetUint64(domain.Cardinality))
	one := fr.One()
	tmp.Sub(&tmp, &one)
	tmp.Mul(&tmp, &domain.CardinalityInv)
	result.Mul(&tmp, &result)

	return &result, indexInDomain, nil
}
