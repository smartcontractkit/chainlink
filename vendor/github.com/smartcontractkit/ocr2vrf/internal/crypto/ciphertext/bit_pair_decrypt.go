package ciphertext

import (
	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/sign/anon"
)

var _ = ((*CipherText)(nil)).Decrypt

func (c *cipherText) decrypt(
	sk kyber.Scalar,
	group anon.Suite, domainSep []byte, sharePublicCommitment kyber.Point,
) (plaintextShare kyber.Scalar, err error) {
	if len(c.cipherText) > plaintextMaxSizeBytes*4 {
		return nil, errors.Errorf("ciphertext too large (max %d pairs)",
			plaintextMaxSizeBytes*4,
		)
	}
	encryptionPK := group.Point().Mul(sk, nil)
	err = c.verify(group, domainSep, encryptionPK, sharePublicCommitment)
	if err != nil {
		return nil, errors.Wrap(err, "refusing to decrypt unverifiable share")
	}

	plaintextShare = group.Scalar()

	zero := group.Scalar().Zero()
	one := group.Scalar().One()
	two := group.Scalar().Add(one, one)
	three := group.Scalar().Add(two, one)
	four := group.Scalar().Add(three, one)
	bitPairs := map[int]kyber.Scalar{0: zero, 1: one, 2: two, 3: three}

	fourPower := group.Scalar().One()

	for _, bitPair := range c.cipherText {
		numericPair, err := bitPair.decrypt(sk)
		if err != nil {

			return nil, errors.Wrapf(err, "could not decrypt %+v as part of share", bitPair)
		}
		if numericPair < 0 || numericPair > 3 {

			return nil, errors.Errorf(
				"bit pair decryption (%d) out of range. Must be less than 4",
				numericPair,
			)
		}

		shiftedBitPair := group.Scalar().Mul(bitPairs[numericPair], fourPower)
		plaintextShare = group.Scalar().Add(plaintextShare, shiftedBitPair)
		fourPower = group.Scalar().Mul(fourPower, four)
	}
	return plaintextShare, nil
}
