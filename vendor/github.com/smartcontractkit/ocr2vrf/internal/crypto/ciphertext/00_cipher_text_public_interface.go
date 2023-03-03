package ciphertext

import (
	"io"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/sign/anon"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"
)

type CipherText struct {
	*cipherText
}

func Encrypt(
	domainSep []byte, group anon.Suite, f *share.PriPoly,
	receiver *player_idx.PlayerIdx, pk kyber.Point,
) (cipherText *CipherText, secretShare kyber.Scalar, err error) {
	ct, secretShare, err := newCipherText(domainSep, group, f, receiver, pk)
	if err != nil {
		return nil, nil, err
	}
	return &CipherText{ct}, secretShare, nil
}

func (c *CipherText) Verify(
	group anon.Suite, domainSep []byte, pk, sharePublicCommitment kyber.Point,
) error {
	return c.verify(group, domainSep, pk, sharePublicCommitment)
}

func (c *CipherText) Decrypt(
	sk kyber.Scalar,
	group anon.Suite, domainSep []byte, sharePublicCommitment kyber.Point,
	receiver player_idx.PlayerIdx,
) (kyber.Scalar, error) {
	return c.decrypt(sk, group, domainSep, sharePublicCommitment)
}

func (c *CipherText) Marshal() ([]byte, error) {
	return c.marshal()
}

func Unmarshal(suite anon.Suite, byteStream io.Reader) (*CipherText, error) {
	c, err := unmarshal(suite, byteStream)
	if err != nil {
		return nil, err
	}
	return &CipherText{c}, nil
}

func (c *CipherText) Equal(c2 *CipherText) bool {
	return c.equal(c2.cipherText)
}
