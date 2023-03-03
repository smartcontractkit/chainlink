package ciphertext

import (
	"github.com/pkg/errors"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/sign/anon"
)

func encrypt(
	domainSep []byte, group anon.Suite, plaintext []byte, pk kyber.Point,
) (cipherText []*elGamalBitPair, totalBlindingSecret kyber.Scalar, err error) {
	if len(plaintext) > plaintextMaxSizeBytes {
		return nil, nil, errors.Errorf("plaintext longer than 256 bits")
	}
	totalBlindingSecret = group.Scalar().Zero()
	one := group.Scalar().One()
	two := group.Scalar().Add(one, one)
	four := two.Add(two, two)
	fourPower := one

	for byteIdx := len(plaintext) - 1; byteIdx >= 0; byteIdx-- {
		cbyte := plaintext[byteIdx]
		for pairIdx := 0; pairIdx < 4; pairIdx++ {

			bitPair := int(cbyte & 0b11)
			cbyte >>= 2

			pairCipherText, blindingSecret, err := newElGamalBitPair(
				group,
				encryptDomainSep(domainSep, uint8(len(cipherText))),
				bitPair, pk,
			)
			if err != nil {

				return nil, nil, errors.Wrapf(err, "while encrypting secret share")
			}
			totalBlindingSecret.Add(
				totalBlindingSecret.Clone(), blindingSecret.Mul(fourPower, blindingSecret.Clone()),
			)
			cipherText = append(cipherText, pairCipherText)
			fourPower = fourPower.Mul(fourPower, four)
		}
	}
	return cipherText, totalBlindingSecret, nil
}

func encryptDomainSep(domainSep []byte, pairIdx uint8) []byte {
	return append(domainSep, pairIdx)
}

func combinedCipherTexts(cipherText []*elGamalBitPair, s anon.Suite) kyber.Point {
	rv := cipherText[0].cipherTextTerm.Clone().Null()
	fourPower := s.Scalar().One()
	two := s.Scalar().Add(fourPower, fourPower)
	four := two.Add(two, two)
	for _, bitpair := range cipherText {
		rv = rv.Add(rv, rv.Clone().Mul(fourPower, bitpair.cipherTextTerm))
		fourPower = fourPower.Mul(fourPower, four)
	}
	return rv
}

const plaintextMaxSizeBytes = 32
