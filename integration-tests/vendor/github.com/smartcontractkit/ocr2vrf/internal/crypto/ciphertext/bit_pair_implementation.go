package ciphertext

import (
	"bytes"
	"reflect"

	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/sign/anon"
)

type elGamalBitPair struct {
	blindingCommitment, cipherTextTerm kyber.Point

	proof bitPairProof

	suite anon.Suite
}

func newElGamalBitPair(suite anon.Suite, domainSep []byte, b int, pk kyber.Point,
) (bp *elGamalBitPair, blindingSecret kyber.Scalar, err error) {
	if reflect.TypeOf(pk) != reflect.TypeOf(suite.Point()) {
		return nil, nil, errors.Errorf(
			"public key (of type %T) does not match group points (of type %T)",
			pk, suite.Point(),
		)
	}
	if (b > 3) || (b < 0) {
		return nil, nil, errors.Errorf("can only encode 0b00, 0b01, 0b10 or 0b11, got 0b%#b", b)
	}
	x := suite.Scalar().Pick(suite.RandomStream())
	blindingCommitment := suite.Point().Mul(x, nil)
	plaintextTerm := suite.Point().Mul(suite.Scalar().SetInt64(int64(b)), nil)
	blindingTerm := suite.Point().Mul(x, pk)
	cipherTextTerm := suite.Point().Add(blindingTerm, plaintextTerm)
	proof, err := proveBitPair(suite, domainSep, pk, blindingCommitment, cipherTextTerm, b, x)
	if err != nil {

		return nil, nil, errors.Wrapf(err, "while constructing encrypted bit pair")
	}
	rv := &elGamalBitPair{blindingCommitment, cipherTextTerm, proof, suite}
	if err := rv.verify(domainSep, pk); err != nil {

		pkm, err2 := pk.MarshalBinary()
		if err2 != nil {
			return nil, nil,
				errors.Wrapf(
					err2, "error while marshalling public key in newElGamalBitPair: %s", err,
				)
		}
		return nil, nil,
			errors.Wrapf(
				err,
				"made unverifiable bit pair! domainSep: 0x%x pk: 0x%x %s x: %s",
				domainSep, pkm, err, x)
	}
	return rv, x, nil
}

var _ = (&CipherText{}).Verify

func (e *elGamalBitPair) verify(domainSep []byte, encryptionKey kyber.Point) error {
	pgroup, err := makeProductGroup(e.suite, encryptionKey)
	if err != nil {
		return errors.Wrapf(err, "while verifying bit pair encryption")
	}
	return e.proof.verify(e.suite, domainSep, pgroup, e.blindingCommitment, e.cipherTextTerm,
		encryptionKey)
}

func (e *elGamalBitPair) decrypt(sk kyber.Scalar) (int, error) {
	if reflect.TypeOf(sk) != reflect.TypeOf(e.suite.Scalar()) {
		return 0, errors.Errorf("need scalar of type %T, got type %T", e.suite.Scalar(), sk)
	}
	plainText := e.suite.Point()

	plainText.Sub(e.cipherTextTerm, plainText.Mul(sk, e.blindingCommitment))
	for i, pt := range memPlainTexts(e.suite) {
		if pt.Equal(plainText) {
			return i, nil
		}
	}
	return 0, errors.Errorf("plaintext unknown")
}

func (e *elGamalBitPair) marshal() (m []byte, err error) {
	rv := make([][]byte, 3)
	cursor := 0

	rv[cursor], err = e.blindingCommitment.MarshalBinary()
	cursor++
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal blinding commitment")
	}

	rv[cursor], err = e.cipherTextTerm.MarshalBinary()
	cursor++
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal ciphertext point")
	}

	rv[cursor] = e.proof[:]
	if cursor != len(rv)-1 {
		panic("return values out of registration")
	}

	return bytes.Join(rv, nil), nil
}

func unmarshalElGamalBitPair(suite anon.Suite, d []byte,
) (e *elGamalBitPair, err error) {
	if len(d) < elGamalBitPairMarshalLength(suite) {
		return nil, errors.Errorf("marshal data too short to contain elGamalBitPair")
	}
	e = &elGamalBitPair{suite: suite}
	e.blindingCommitment = suite.Point()
	pointLen := e.blindingCommitment.MarshalSize()

	if err := e.blindingCommitment.UnmarshalBinary(d[:pointLen]); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal blinding commitment")
	}
	remainder := d[pointLen:]

	e.cipherTextTerm = suite.Point()
	if err := e.cipherTextTerm.UnmarshalBinary(remainder[:pointLen]); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal ciphertext point")
	}
	remainder = remainder[pointLen:]

	proofLen := bitPairProofLen(suite)
	e.proof = append([]byte{}, remainder[:proofLen]...)

	return e, nil
}

func elGamalBitPairMarshalLength(suite anon.Suite) int {
	pointLen := suite.PointLen()
	return pointLen +
		pointLen +
		bitPairProofLen(suite)
}

func (e *elGamalBitPair) equal(e2 *elGamalBitPair) bool {
	return e.blindingCommitment.Equal(e2.blindingCommitment) &&
		e.cipherTextTerm.Equal(e2.cipherTextTerm) &&
		bytes.Equal(e.proof[:], e2.proof[:]) &&
		e.suite.String() == e2.suite.String()
}
