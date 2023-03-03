package dkg

import (
	"bytes"
	"encoding/binary"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/ciphertext/schnorr"
	"github.com/smartcontractkit/ocr2vrf/internal/crypto/point_translation"
	"github.com/smartcontractkit/ocr2vrf/internal/pvss"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/sign/anon"
)

type shareRecord struct {
	shareSet *pvss.ShareSet

	marshaledShareRecord []byte

	sig signature
}

type signature struct{ sig []byte }

func newShareRecord(
	suite schnorr.Suite, shareSet *pvss.ShareSet, sk kyber.Scalar,
	domainSep types.ConfigDigest,
) (*shareRecord, error) {
	rv := &shareRecord{shareSet: shareSet}

	if err := rv.sign(suite, domainSep, sk); err != nil {
		return nil, errors.Wrapf(err, "could not sign new share record")
	}
	return rv, nil
}

const shareLenBound = 10_000_000

func (r *shareRecord) marshal() ([]byte, error) {
	msr := r.marshaledShareRecord
	if msr == nil {
		ss, err := r.shareSet.Marshal()
		if err != nil {
			return nil, errors.Wrap(err, "could not marshal share record")
		}
		if len(ss) > shareLenBound {
			return nil, errors.Wrap(err, "could not marshal share record: marshalled share set too long")
		}
		msr = msrComponents(ss, r.sig.sig)
		r.marshaledShareRecord = msr
	}
	return msr, nil
}

func msrComponents(ssBytes, sig []byte) []byte {
	ssLenData := make([]byte, 4)
	binary.BigEndian.PutUint32(ssLenData, uint32(len(ssBytes)))
	return bytes.Join([][]byte{ssLenData, ssBytes, sig}, nil)
}

func unmarshalShareRecord(
	sigSuite schnorr.Suite, g anon.Suite, translationGroup kyber.Group,
	data []byte, translation point_translation.PubKeyTranslation,
	cfgDgst types.ConfigDigest, pks []kyber.Point, spks []kyber.Point,
) (*shareRecord, []byte, error) {
	if len(data) < 4 {
		return nil, nil, errors.Errorf("marshalled share record too short, %d bytes", len(data))
	}
	if len(data) > shareLenBound {
		return nil, nil, errors.Errorf("marshalled share record too long, %d bytes", len(data))
	}
	ssLenData, data := data[:4], data[4:]
	ssLen := binary.BigEndian.Uint32(ssLenData)
	if int(ssLen) > len(data)+1 {
		return nil, nil, errors.Errorf(
			"marshalled share record too short, %d bytes, need %d bytes",
			len(data), ssLen,
		)
	}
	if ssLen > shareLenBound {
		return nil, nil, errors.Errorf("marshalled share record too long, %d bytes", len(data))
	}
	ssBytes, data := data[:ssLen], data[ssLen:]

	dealer, err := pvss.UnmarshalDealer(ssBytes)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get signer identity for share record")
	}
	if !dealer.AtMost(uint8(len(spks))) {
		return nil, nil, errors.Errorf("dealer out of range")
	}

	sig, data := data[:64], data[64:]
	dealerPK := dealer.Index(spks).(kyber.Point)

	msg := append(cfgDgst[:], ssBytes...)
	err = schnorr.Verify(sigSuite, dealerPK, msg, sig)
	if err != nil {
		return nil, nil, errors.Wrap(err, "invalid signature on marshalled share set")
	}

	shareSet, _, err := pvss.UnmarshalShareSet(
		g, translationGroup, ssBytes, translation, cfgDgst, pks,
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not unmarshal share record")
	}

	msr := msrComponents(ssBytes, sig)
	return &shareRecord{shareSet, msr, signature{sig}}, data, nil
}

func (r *shareRecord) sign(suite schnorr.Suite,
	domainSep types.ConfigDigest,
	sk kyber.Scalar) error {
	ss, err := r.shareSet.Marshal()
	if err != nil {
		return errors.Wrap(err, "could not marshal share set for signing")
	}
	msg := append(domainSep[:], ss...)
	r.sig.sig, err = schnorr.Sign(suite, sk, msg)
	if err != nil {
		return errors.Wrap(err, "could sign share set")
	}

	pk := suite.Point().Mul(sk, nil)

	if err := schnorr.Verify(suite, pk, msg, r.sig.sig); err != nil {
		panic("failed to verify own signature: " + err.Error())
	}
	return nil
}
