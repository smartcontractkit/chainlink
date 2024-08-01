package ipa

import (
	"fmt"

	"github.com/crate-crypto/go-ipa/bandersnatch/fr"
	"github.com/crate-crypto/go-ipa/common"
)

// domainSize will always equal 256, which is the same
// as the degree of the polynomial (+1), we are committing to.
// This constant is defined here for semantic reasons.
const domainSize = common.VectorLength

// PrecomputedWeights contains precomputed coefficients for efficient
// usage of the Barycentric formula.
type PrecomputedWeights struct {
	// This stores A'(x_i) and 1/A'(x_i)
	barycentricWeights []fr.Element
	// This stores 1/k and -1/k for k \in [0, 255]
	invertedDomain []fr.Element
}

// NewPrecomputedWeights generates the precomputed weights for the barycentric formula.
func NewPrecomputedWeights() *PrecomputedWeights {
	// Imagine we have two arrays of the same length and we concatenate them together
	// This is how we will store the A'(x_i) and 1/A'(x_i)
	// This midpoint variable is used to compute the offset that we need
	// to place 1/A'(x_i)
	midpoint := uint64(domainSize)

	// Note there are DOMAIN_SIZE number of weights, but we are also storing their inverses
	// so we need double the amount of space
	barycentricWeights := make([]fr.Element, midpoint*2)
	for i := uint64(0); i < midpoint; i++ {
		weight := computeBarycentricWeightForElement(i)

		var invWeight fr.Element
		invWeight.Inverse(&weight)

		barycentricWeights[i] = weight
		barycentricWeights[i+midpoint] = invWeight
	}

	// Computing 1/k and -1/k for k \in [0, 255]
	// Note that since we cannot do 1/0, we have one less element
	midpoint = domainSize - 1
	invertedDomain := make([]fr.Element, midpoint*2)
	for i := uint64(1); i < domainSize; i++ {
		var k fr.Element
		k.SetUint64(i)
		k.Inverse(&k)

		var negative_k fr.Element
		zero := fr.Zero()
		negative_k.Sub(&zero, &k)

		invertedDomain[i-1] = k
		invertedDomain[(i-1)+midpoint] = negative_k
	}

	return &PrecomputedWeights{
		barycentricWeights: barycentricWeights,
		invertedDomain:     invertedDomain,
	}

}

// computes A'(x_j) where x_j must be an element in the domain
// This is computed as the product of x_j - x_i where x_i is an element in the domain
// and x_i is not equal to x_j
func computeBarycentricWeightForElement(element uint64) fr.Element {
	// let domain_element_fr = Fr::from(domain_element as u128);
	if element > domainSize {
		panic(fmt.Sprintf("the domain is [0,255], %d is not in the domain", element))
	}

	var domain_element_fr fr.Element
	domain_element_fr.SetUint64(element)

	total := fr.One()

	for i := uint64(0); i < domainSize; i++ {
		if i == element {
			continue
		}

		var i_fr fr.Element
		i_fr.SetUint64(i)

		var tmp fr.Element
		tmp.Sub(&domain_element_fr, &i_fr)

		total.Mul(&total, &tmp)
	}

	return total
}

// ComputeBarycentricCoefficients, computes the coefficients `bary_coeffs`
// for a point `z` such that when we have a polynomial `p` in lagrange
// basis, the inner product of `p` and `bary_coeffs` is equal to p(z)
// Note that `z` should not be in the domain.
// This can also be seen as the lagrange coefficients L_i(point)
func (preComp *PrecomputedWeights) ComputeBarycentricCoefficients(point fr.Element) []fr.Element {
	// Compute A'(x_i) * (point - x_i)
	lagrangeEvals := make([]fr.Element, domainSize)
	for i := uint64(0); i < domainSize; i++ {
		weight := preComp.barycentricWeights[i]

		var i_fr fr.Element
		i_fr.SetUint64(i)
		lagrangeEvals[i].Sub(&point, &i_fr)
		lagrangeEvals[i].Mul(&lagrangeEvals[i], &weight)
	}

	totalProd := fr.One()
	for i := uint64(0); i < domainSize; i++ {
		var i_fr fr.Element
		i_fr.SetUint64(i)

		var tmp fr.Element
		tmp.Sub(&point, &i_fr)
		totalProd.Mul(&totalProd, &tmp)
	}

	lagrangeEvals = fr.BatchInvert(lagrangeEvals)
	for i := uint64(0); i < domainSize; i++ {
		lagrangeEvals[i].Mul(&lagrangeEvals[i], &totalProd)
	}

	return lagrangeEvals
}

// DivideOnDomain computes f(x) - f(x_i) / x - x_i where x_i is an element in the domain
func (preComp *PrecomputedWeights) DivideOnDomain(index uint8, f []fr.Element) []fr.Element {
	quotient := make([]fr.Element, domainSize)

	y := f[index]

	for i := 0; i < domainSize; i++ {
		if i != int(index) {
			den := i - int(index)
			absDen, is_neg := absInt(den)

			denInv := preComp.getInvertedElement(absDen, is_neg)

			// compute q_i
			quotient[i].Sub(&f[i], &y)
			quotient[i].Mul(&quotient[i], &denInv)

			weightRatio := preComp.getRatioOfWeights(int(index), i)
			var tmp fr.Element
			tmp.Mul(&weightRatio, &quotient[i])
			quotient[index].Sub(&quotient[index], &tmp)
		}
	}

	return quotient
}

func (preComp *PrecomputedWeights) getInvertedElement(element int, is_neg bool) fr.Element {
	index := element - 1

	if is_neg {
		midpoint := len(preComp.invertedDomain) / 2
		index += midpoint
	}

	return preComp.invertedDomain[index]
}

func (preComp *PrecomputedWeights) getRatioOfWeights(numerator int, denominator int) fr.Element {

	a := preComp.barycentricWeights[numerator]
	midpoint := len(preComp.barycentricWeights) / 2
	b := preComp.barycentricWeights[denominator+midpoint]

	var result fr.Element
	result.Mul(&a, &b)
	return result
}

func (preComp *PrecomputedWeights) getInverseBarycentricWeight(i int) fr.Element {

	midpoint := len(preComp.barycentricWeights) / 2
	return preComp.barycentricWeights[i+midpoint]
}

// Returns the absolute value and true if
// the value was negative
func absInt(x int) (int, bool) {
	is_negative := x < 0

	if is_negative {
		return -x, is_negative
	}

	return x, is_negative
}
