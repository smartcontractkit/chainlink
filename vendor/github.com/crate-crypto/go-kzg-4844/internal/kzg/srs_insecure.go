// Methods in this file should not be used in production.
// They are used in order to create trusted setup instances
// for testing and or development.

package kzg

import (
	"math/big"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

// newLagrangeSRSInsecure creates a new SRS object with the secret `bAlpha`.
// The resulting SRS is in Lagrange basis.
//
// This method should not be used in production because as the secret is supplied as input.
func newLagrangeSRSInsecure(domain Domain, bAlpha *big.Int) (*SRS, error) {
	return newSRSInsecure(domain, bAlpha, true)
}

// newMonomialSRSInsecure creates a new SRS object with the secret `bAlpha`.
// The resulting SRS is in Monomial basis.
//
// This method should not be used in production because as the secret is supplied as input.
func newMonomialSRSInsecure(domain Domain, bAlpha *big.Int) (*SRS, error) {
	return newSRSInsecure(domain, bAlpha, false)
}

// newSRSInsecure creates a new SRS object with the secret `bAlpha`.
// convertToLagrange controls whether the result is in monomial or Lagrange basis.
//
// This method should not be used in production because as the secret is supplied as input.
func newSRSInsecure(domain Domain, bAlpha *big.Int, convertToLagrange bool) (*SRS, error) {
	srs, err := newMonomialSRSInsecureUint64(domain.Cardinality, bAlpha)
	if err != nil {
		return nil, err
	}

	if convertToLagrange {
		// Convert SRS from monomial form to lagrange form
		lagrangeG1 := domain.IfftG1(srs.CommitKey.G1)
		srs.CommitKey.G1 = lagrangeG1
	}

	return srs, nil
}

// newMonomialSRSInsecureUint64 creates a new SRS object with the secret `bAlpha` in monomial basis.
//
// Note that the function name ends with Uint64, because we provide the size argument as a
// uint64 rather than a Domain. A newMonomialSRSInsecure functions taking a Domain as input
// to match the other functions is defined in the testing code.
//
// This method should not be used in production because as the secret is supplied as input.
//
// Copied from [gnark-crypto].
//
// [gnark-crypto]: https://github.com/ConsenSys/gnark-crypto/blob/8f7ca09273c24ed9465043566906cbecf5dcee91/ecc/bls12-381/fr/kzg/kzg.go#L65
func newMonomialSRSInsecureUint64(size uint64, bAlpha *big.Int) (*SRS, error) {
	if size < 2 {
		return nil, ErrMinSRSSize
	}

	var commitKey CommitKey
	var openKey OpeningKey
	commitKey.G1 = make([]bls12381.G1Affine, size)

	var alpha fr.Element
	alpha.SetBigInt(bAlpha)

	_, _, gen1Aff, gen2Aff := bls12381.Generators()
	commitKey.G1[0] = gen1Aff
	openKey.GenG1 = gen1Aff
	openKey.GenG2 = gen2Aff
	openKey.AlphaG2.ScalarMultiplication(&gen2Aff, bAlpha)

	alphas := make([]fr.Element, size-1)
	alphas[0] = alpha
	for i := 1; i < len(alphas); i++ {
		alphas[i].Mul(&alphas[i-1], &alpha)
	}
	g1s := bls12381.BatchScalarMultiplicationG1(&gen1Aff, alphas)
	copy(commitKey.G1[1:], g1s)

	return &SRS{
		CommitKey:  commitKey,
		OpeningKey: openKey,
	}, nil
}
