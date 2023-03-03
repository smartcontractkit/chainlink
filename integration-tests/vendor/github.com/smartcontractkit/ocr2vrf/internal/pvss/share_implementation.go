package pvss

import (
	"bytes"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/ciphertext"
	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"

	"go.dedis.ch/kyber/v3"
	kshare "go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/sign/anon"
)

type share struct {
	cipherText *ciphertext.CipherText

	encryptionKey kyber.Point

	subKeyTranslation kyber.Point

	shareSet *ShareSet
}

func newShare(domainSep []byte, group anon.Suite,
	secretPoly *kshare.PriPoly, shareSet *ShareSet, receiver *player_idx.PlayerIdx,
	pk kyber.Point,
) (rv *share, err error) {
	rv = &share{shareSet: shareSet, encryptionKey: pk}
	domainSep = rv.domainSep(domainSep, receiver)
	var secretShare kyber.Scalar
	rv.cipherText, secretShare, err = ciphertext.Encrypt(
		domainSep, group, secretPoly, receiver, pk,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encrypt PVSS share")
	}
	rv.subKeyTranslation, err = shareSet.translation.TranslateKey(secretShare)
	if err != nil {
		return nil, errors.Wrapf(err, "could not translate pub key for secret share")
	}
	if err2 := rv.verify(group, domainSep, receiver); err2 != nil {
		panic("could not verify share just constructed")
	}
	return rv, err
}

func (s *share) decrypt(
	sk kyber.Scalar,
	keyGroup anon.Suite,
	domainSep []byte,
	sharePublicCommitment kyber.Point,
	receiver player_idx.PlayerIdx,
) (kshare.PriShare, error) {
	plaintextShare, err := s.cipherText.Decrypt(
		sk, keyGroup, s.domainSep(domainSep, &receiver), sharePublicCommitment,
		receiver,
	)
	if err != nil {
		return kshare.PriShare{}, errors.Wrap(err, "could not decrypt share")
	}
	return receiver.PriShare(plaintextShare), nil
}

func (s *share) verify(
	group anon.Suite, domainSep []byte, receiver *player_idx.PlayerIdx,
) error {

	sharePublicCommitment := receiver.EvalPoint(s.shareSet.coeffCommitments)
	err := s.shareSet.translation.VerifyTranslation(sharePublicCommitment, s.subKeyTranslation)
	if err != nil {
		return errors.Wrapf(err, "bad translation of share public key")
	}
	return s.cipherText.Verify(group, domainSep, s.encryptionKey, sharePublicCommitment)
}

func (s *share) domainSep(domainSep []byte, receiver *player_idx.PlayerIdx) []byte {
	return bytes.Join([][]byte{domainSep, receiver.Marshal()}, nil)
}

func (s *share) equal(s2 *share) bool {
	return s.cipherText.Equal(s2.cipherText) &&
		s.encryptionKey.Equal(s2.encryptionKey) &&
		((s.subKeyTranslation == nil && s2.subKeyTranslation == nil) ||
			s.subKeyTranslation.Equal(s2.subKeyTranslation))
}
