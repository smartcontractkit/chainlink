package ciphertext

import (
	"bytes"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/change_base_group"
	"github.com/smartcontractkit/ocr2vrf/internal/crypto/ciphertext/schnorr"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/sign/anon"
)

type dLKnowledgeProof []byte

func newDLKnowledgeProof(
	domainSep []byte, group anon.Suite, h kyber.Point, b kyber.Scalar,
) (dLKnowledgeProof, error) {

	g, err := change_base_group.NewChangeBaseGroup(group, h)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create cyclic group from given generator")
	}

	rv, err := schnorr.Sign(g, b, domainSep)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate knowledge proof")
	}

	pk := g.Point().Mul(b, nil)
	msg := domainSep
	sig := rv
	if err2 := schnorr.Verify(g, pk, msg, sig); err2 != nil {

		panic(errors.Wrapf(err2, "created unverifiable signature"))
	}
	return rv, err
}

func (p dLKnowledgeProof) verify(
	group anon.Suite, domainSep []byte, blindingPK, signingPK kyber.Point,
) error {
	g, err := change_base_group.NewChangeBaseGroup(group, blindingPK)
	if err != nil {
		return errors.Wrapf(err, "could not create cyclic group from given generator")
	}
	pk, err := g.Lift(signingPK)
	if err != nil {
		return errors.Wrapf(err, "could not create point pair for DL proof")
	}

	sig := []byte(p)
	msg := domainSep
	return errors.Wrap(schnorr.Verify(g, pk, msg, sig),
		"could not verify share proof")
}

func (p dLKnowledgeProof) equal(p2 dLKnowledgeProof) bool {
	return bytes.Equal(p, p2)
}
