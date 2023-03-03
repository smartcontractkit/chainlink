package ciphertext

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/product_group"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/sign/anon"
)

type bitPairProof []byte

func proveBitPair(
	suite anon.Suite, domainSep []byte,
	pk, blindingCommitment, cipherTextTerm kyber.Point,
	bitPair int, secret kyber.Scalar,
) (rv bitPairProof, err error) {
	if (bitPair < 0) || (bitPair > 3) {
		return rv, errors.Errorf("bitPair must be in {0,1,2,3}, got %d", bitPair)
	}

	pgroup, err := makeProductGroup(suite, pk)
	if err != nil {
		return rv, errors.Wrap(err, "while proving common discrete log")
	}

	possiblePKs, err := blindingPKs(suite, pgroup, blindingCommitment, cipherTextTerm)
	if err != nil {
		return rv, errors.Wrapf(err, "while constructing bit-pair proof")
	}
	pPKs := pksToPoints(possiblePKs)

	sig := anon.Sign(pgroup, domainSep, pPKs, nil, bitPair, secret)
	if len(sig) != bitPairProofLen(suite) {
		return rv, errors.Errorf("signature wrong length")
	}
	return sig, nil
}

func bitPairProofLen(suite anon.Suite) int {

	return (4 + 1) * suite.ScalarLen()
}

func makeProductGroup(suite anon.Suite, pk kyber.Point) (*product_group.ProductGroup, error) {
	return product_group.NewProductGroup(
		[]kyber.Group{suite, suite}, []kyber.Point{suite.Point().Base(), pk}, suite.RandomStream(),
	)
}

func (p bitPairProof) verify(suite anon.Suite,
	domainSep []byte, pgroup *product_group.ProductGroup,
	blindingCommitment, cipherTextTerm, pk kyber.Point,
) error {
	possibleBlindingTerms, err := blindingPKs(suite, pgroup, blindingCommitment, cipherTextTerm)
	if err != nil {
		return errors.Wrap(err, "while constructing possible keys for bit-pair proof")
	}
	possiblePKs := pksToPoints(possibleBlindingTerms)

	_, err = anon.Verify(pgroup, domainSep, possiblePKs, nil, p)
	if err != nil {
		return errors.Wrap(err, "while verifying bit-pair proof")
	}
	return nil
}

func blindingPKs(
	suite anon.Suite, pgroup *product_group.ProductGroup,
	blindingCommitment, cipherTextTerm kyber.Point,
) ([]*product_group.PointTuple, error) {
	var possibleBlindingTerms []*product_group.PointTuple
	for _, possiblePlaintextTerm := range memPlainTexts(suite) {
		possibleBlindingTerm := cipherTextTerm.Clone().Sub(cipherTextTerm, possiblePlaintextTerm)
		possibleKnownDLPoint, err := pgroup.NewPoint(
			[]kyber.Point{blindingCommitment, possibleBlindingTerm},
		)
		if err != nil {
			return nil, errors.Wrapf(err, "while constructing possible binding term")
		}
		possibleBlindingTerms = append(possibleBlindingTerms, possibleKnownDLPoint)
	}
	return possibleBlindingTerms, nil
}

func pksToPoints(pks []*product_group.PointTuple) (rv []kyber.Point) {
	for _, pk := range pks {
		rv = append(rv, kyber.Point(pk))
	}
	return rv
}
