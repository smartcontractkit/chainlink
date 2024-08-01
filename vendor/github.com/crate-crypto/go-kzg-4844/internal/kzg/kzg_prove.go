package kzg

import (
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

// Open verifies that a polynomial f(x) when evaluated at a point `z` is equal to `f(z)`
//
// numGoRoutines is used to configure the amount of concurrency needed. Setting this
// value to a negative number or 0 will make it default to the number of CPUs.
//
// [compute_kzg_proof_impl]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#compute_kzg_proof_impl
func Open(domain *Domain, p Polynomial, evaluationPoint fr.Element, ck *CommitKey, numGoRoutines int) (OpeningProof, error) {
	if len(p) == 0 || len(p) > len(ck.G1) {
		return OpeningProof{}, ErrInvalidPolynomialSize
	}

	outputPoint, indexInDomain, err := domain.evaluateLagrangePolynomial(p, evaluationPoint)
	if err != nil {
		return OpeningProof{}, err
	}

	// Compute the quotient polynomial
	quotientPoly, err := domain.computeQuotientPoly(p, indexInDomain, *outputPoint, evaluationPoint)
	if err != nil {
		return OpeningProof{}, err
	}

	// Commit to Quotient polynomial
	quotientCommit, err := Commit(quotientPoly, ck, numGoRoutines)
	if err != nil {
		return OpeningProof{}, err
	}

	res := OpeningProof{
		InputPoint:   evaluationPoint,
		ClaimedValue: *outputPoint,
	}

	res.QuotientCommitment.Set(quotientCommit)

	return res, nil
}

// computeQuotientPoly computes q(X) = (f(X) - f(z)) / (X - z) in Lagrange form.
//
// We refer to the result q(X) as the quotient polynomial.
//
// The division needs to be handled differently if `z` is an element in the domain
// because the naive formula would compute 0/0. Hence, you will observe that this function
// will follow a different code-path depending on this condition.
//
// In our situation, both f(z) and whether z is inside the domain are always known to the caller,
// so we just take is as input rather than (re-)computing it ourself. The method does not check that those
// values provided are correct.
//
// indexInDomain needs to be set to -1 to indicate that z is not in the domain and to the index in the domain if it is.
//
// The matching code for this method is in `compute_kzg_proof_impl` where the quotient polynomial
// is computed.
func (domain *Domain) computeQuotientPoly(f Polynomial, indexInDomain int64, fz, z fr.Element) (Polynomial, error) {
	if domain.Cardinality != uint64(len(f)) {
		return nil, ErrPolynomialMismatchedSizeDomain
	}

	if indexInDomain != -1 {
		// Note: the uint64 conversion is both semantically correct and safer
		// than accepting an `int``, since we know it shouldn't be negative
		// and it should cause a panic, if not checked; uint64(-1) = 2^64 -1
		return domain.computeQuotientPolyOnDomain(f, uint64(indexInDomain))
	}

	return domain.computeQuotientPolyOutsideDomain(f, fz, z)
}

// computeQuotientPolyOutsideDomain computes q(X) = (f(X) - f(z)) / (X - z) in lagrange form where `z` is not in the domain.
//
// This is the implementation of computeQuotientPoly for the case where z is not in the domain.
// Since both input and output polynomials are given in evaluation form, this method just performs the desired operation pointwise.
func (domain *Domain) computeQuotientPolyOutsideDomain(f Polynomial, fz, z fr.Element) (Polynomial, error) {
	// Compute the lagrange form the of the numerator f(X) - f(z)
	// Since f(X) is already in lagrange form, we can compute f(X) - f(z)
	// by shifting all elements in f(X) by f(z)
	numerator := make(Polynomial, len(f))
	for i := 0; i < len(f); i++ {
		numerator[i].Sub(&f[i], &fz)
	}

	// Compute the lagrange form of the denominator X - z.
	// This means that we need to compute w - z for all points w in the domain.
	denominator := make(Polynomial, len(f))
	for i := 0; i < len(f); i++ {
		denominator[i].Sub(&domain.Roots[i], &z)
	}

	// To invert the denominator polynomial at each point of the domain, we perform a batch-inversion.
	// Since `z` is not in the domain, we are sure that there are no zeroes in this inversion.
	//
	// Note: if there was a zero, the gnark-crypto library would skip
	// it and not panic.
	denominator = fr.BatchInvert(denominator)

	// Compute the quotient q(X)
	for i := 0; i < len(f); i++ {
		denominator[i].Mul(&denominator[i], &numerator[i])
	}

	return denominator, nil
}

// computeQuotientPolyOnDomain computes (f(X) - f(z)) / (X - z) in Lagrange form where `z` is in the domain.
//
// This is the implementation of computeQuotientPoly for the case where the evaluation point is in the domain.
//
// [compute_quotient_eval_within_domain]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#compute_quotient_eval_within_domain
func (domain *Domain) computeQuotientPolyOnDomain(f Polynomial, index uint64) (Polynomial, error) {
	fz := f[index]
	z := domain.Roots[index]
	invZ := domain.PreComputedInverses[index]

	// Compute the evaluation of X - z at every point in the domain.
	rootsMinusZ := make([]fr.Element, domain.Cardinality)
	for i := 0; i < int(domain.Cardinality); i++ {
		rootsMinusZ[i].Sub(&domain.Roots[i], &z)
	}

	// Since we know that `z` is in the domain, rootsMinusZ[index] will be zero.
	// We set this value to `1` instead to compute the batch inversion without having to special-case here.
	// This way, the value of rootsMinusZ[index] will stay untouched.
	// Note: The underlying gnark-crypto library will not panic if
	// one of the elements is zero, but this is not common across libraries so we just set it to one.
	rootsMinusZ[index].SetOne()

	// Evaluation of 1/(X-z) at every point of the domain, except for index.
	invRootsMinusZ := fr.BatchInvert(rootsMinusZ)

	quotientPoly := make(Polynomial, domain.Cardinality)
	for j := 0; j < int(domain.Cardinality); j++ {
		// Check if we are on the current root of unity
		// Note: For notations below, we use `m` to denote `index`
		if uint64(j) == index {
			continue
		}

		// Compute q_j = f_j / w^j - w^m for j != m.
		// This is exactly the same as in the computeQuotientPolyOutsideDomain - case.
		//
		// Note: f_j is the numerator of the quotient polynomial ie f_j = f[j] - f(z)
		//
		//
		var q_j fr.Element
		q_j.Sub(&f[j], &fz)
		q_j.Mul(&q_j, &invRootsMinusZ[j])
		quotientPoly[j] = q_j

		// Compute the contribution to q_m coming from the j'th term of the input.
		// This term is given by
		// q_m_j = (f_j / w^m - w^j) * (w^j/w^m) , where w^m = z
		//		 = - q_j * w^{j-m}
		//
		// We _could_ find 1 / w^{j-m} via a lookup table
		// but we want to avoid lookup tables because
		// the roots are bit-reversed which can make the
		// code less readable.
		var q_m_j fr.Element
		q_m_j.Neg(&q_j)
		q_m_j.Mul(&q_m_j, &domain.Roots[j])
		q_m_j.Mul(&q_m_j, &invZ)

		quotientPoly[index].Add(&quotientPoly[index], &q_m_j)
	}

	return quotientPoly, nil
}
