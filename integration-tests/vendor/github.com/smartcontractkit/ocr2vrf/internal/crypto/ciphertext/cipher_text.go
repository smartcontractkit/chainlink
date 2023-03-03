package ciphertext

import (
	"bytes"

	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/sign/anon"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"
)

type cipherText struct {
	cipherText            []*elGamalBitPair
	receiver              *player_idx.PlayerIdx
	encryptionKey         kyber.Point
	encodesShareProof     dLKnowledgeProof
	suite                 anon.Suite
	sharePublicCommitment kyber.Point
}

var _ = Encrypt

func newCipherText(
	domainSep []byte, group anon.Suite, f *share.PriPoly,
	receiver *player_idx.PlayerIdx, pk kyber.Point,
) (rv *cipherText, secretShare kyber.Scalar, err error) {
	rv = &cipherText{suite: group, receiver: receiver, encryptionKey: pk}
	secretShare, rawShare, err := getShareBits(receiver, f)
	if err != nil {
		return nil, nil, err
	}
	var totalBlindingSecret kyber.Scalar

	rv.cipherText, totalBlindingSecret, err = encrypt(domainSep, group, rawShare, pk)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "could not encrypt secret share")
	}

	if err := rv.proveFinalDLKnowledge(domainSep, group, pk, totalBlindingSecret); err != nil {

		return nil, nil, errors.Wrap(err, "could not prove ciphertext encodes secret share")
	}
	blindingPK := group.Point().Mul(secretShare, nil)
	if err := rv.verify(group, domainSep, pk, blindingPK); err != nil {

		panic(err)
	}
	return rv, secretShare, nil

}

func (c *cipherText) cipherTextDomainSep(domainSep []byte) (ds []byte, err error) {
	encm := make([][]byte, len(c.cipherText)+1)
	encm[0] = domainSep
	offset := 1
	for ctidx, ct := range c.cipherText {
		encm[ctidx+offset], err = ct.marshal()
		if err != nil {
			return nil, errors.Wrapf(err, "could not marshal ciphertext for domain separator")
		}
	}
	return bytes.Join(encm, nil), nil
}

var _ = ((*CipherText)(nil)).Verify

func (c *cipherText) verify(
	group anon.Suite, domainSep []byte, encryptionPK, sharePublicCommitment kyber.Point,
) error {
	if len(c.cipherText) > plaintextMaxSizeBytes*4 {
		return errors.Errorf("ciphertext too large (max %d pairs)",
			plaintextMaxSizeBytes*4,
		)
	}

	combinedBlindingFactors := group.Point().Sub(
		combinedCipherTexts(c.cipherText, group),
		sharePublicCommitment,
	)
	edomain, err := c.cipherTextDomainSep(domainSep)
	if err != nil {
		return err
	}

	err = c.encodesShareProof.verify(
		group, edomain, encryptionPK, combinedBlindingFactors,
	)
	if err != nil {
		return errors.Wrapf(err, "could not verify overall share-encoding proof")
	}
	for pairIdx, bitPair := range c.cipherText {
		err := bitPair.verify(
			encryptDomainSep(domainSep, uint8(pairIdx)),
			c.encryptionKey,
		)
		if err != nil {

			return errors.Wrapf(err, "part of ciphertext failed to verify")
		}
	}
	return nil
}

func (c *cipherText) proveFinalDLKnowledge(
	domainSep []byte, group anon.Suite, pk kyber.Point,
	totalBlindingSecret kyber.Scalar,
) error {

	dlDomainSep, err := c.cipherTextDomainSep(domainSep)
	if err != nil {
		return err
	}
	c.encodesShareProof, err = newDLKnowledgeProof(dlDomainSep, group, pk, totalBlindingSecret)
	if err != nil {
		return errors.Wrapf(err, "could not construct proof that ciphertext contains secret share")
	}
	return nil
}

func (c *cipherText) equal(c2 *cipherText) bool {
	if c.receiver.Equal(c2.receiver) &&
		c.encryptionKey.Equal(c2.encryptionKey) &&
		c.encodesShareProof.equal(c2.encodesShareProof) &&
		c.suite.String() == c2.suite.String() &&
		len(c.cipherText) == len(c2.cipherText) {

		for i, ct := range c.cipherText {
			if !ct.equal(c2.cipherText[i]) {
				return false
			}
		}
		return true
	}
	return false
}
