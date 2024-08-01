package ipa

import (
	"fmt"

	"github.com/crate-crypto/go-ipa/bandersnatch/fr"
	"github.com/crate-crypto/go-ipa/banderwagon"
	"github.com/crate-crypto/go-ipa/common"
)

// CheckIPAProof verifies an IPA proof for a committed polynomial in evaluation form.
// It verifies that `proof` is a valid proof for the polynomial at the evaluation
// point `evalPoint` with result `result`
func CheckIPAProof(transcript *common.Transcript, ic *IPAConfig, commitment banderwagon.Element, proof IPAProof, evalPoint fr.Element, result fr.Element) (bool, error) {
	transcript.DomainSep(labelDomainSep)

	if len(proof.L) != len(proof.R) {
		return false, fmt.Errorf("vectors L and R should be the same size")
	}
	if len(proof.L) != int(ic.numRounds) {
		return false, fmt.Errorf("the number of points for L and R should be equal to the number of rounds")
	}

	b := computeBVector(ic, evalPoint)

	transcript.AppendPoint(&commitment, labelC)
	transcript.AppendScalar(&evalPoint, labelInputPoint)
	transcript.AppendScalar(&result, labelOutputPoint)

	w := transcript.ChallengeScalar(labelW)

	// Rescaling of q.
	var q banderwagon.Element
	q.ScalarMul(&ic.Q, &w)

	var qy banderwagon.Element
	qy.ScalarMul(&q, &result)
	commitment.Add(&commitment, &qy)

	challenges := generateChallenges(transcript, &proof)
	challengesInv := fr.BatchInvert(challenges)

	// Compute expected commitment
	var err error
	for i := 0; i < len(challenges); i++ {
		x := challenges[i]
		L := proof.L[i]
		R := proof.R[i]

		commitment, err = commit([]banderwagon.Element{commitment, L, R}, []fr.Element{fr.One(), x, challengesInv[i]})
		if err != nil {
			return false, fmt.Errorf("could not compute commitment+x*L+x^-1*R: %w", err)
		}
	}

	g := ic.SRS

	// We compute the folding-scalars for g and b.
	foldingScalars := make([]fr.Element, len(g))
	for i := 0; i < len(g); i++ {
		scalar := fr.One()

		for challengeIdx := 0; challengeIdx < len(challenges); challengeIdx++ {
			if i&(1<<(7-challengeIdx)) > 0 {
				scalar.Mul(&scalar, &challengesInv[challengeIdx])
			}
		}
		foldingScalars[i] = scalar
	}
	g0, err := MultiScalar(g, foldingScalars)
	if err != nil {
		return false, fmt.Errorf("could not compute g0: %w", err)
	}
	b0, err := InnerProd(b, foldingScalars)
	if err != nil {
		return false, fmt.Errorf("could not compute b0: %w", err)
	}

	var got banderwagon.Element
	//  g0 * a + (a * b) * Q;
	var part_1 banderwagon.Element
	part_1.ScalarMul(&g0, &proof.A_scalar)

	var part_2 banderwagon.Element
	var part_2a fr.Element

	part_2a.Mul(&b0, &proof.A_scalar)
	part_2.ScalarMul(&q, &part_2a)

	got.Add(&part_1, &part_2)

	return got.Equal(&commitment), nil
}

func generateChallenges(transcript *common.Transcript, proof *IPAProof) []fr.Element {

	challenges := make([]fr.Element, len(proof.L))
	for i := 0; i < len(proof.L); i++ {
		transcript.AppendPoint(&proof.L[i], labelL)
		transcript.AppendPoint(&proof.R[i], labelR)
		challenges[i] = transcript.ChallengeScalar(labelX)
	}
	return challenges
}
