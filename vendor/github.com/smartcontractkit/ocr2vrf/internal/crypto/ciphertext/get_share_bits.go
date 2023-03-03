package ciphertext

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"
)

func getShareBits(receiver *player_idx.PlayerIdx, secretPoly *share.PriPoly) (kyber.Scalar, []byte, error) {
	privateShare := receiver.Eval(secretPoly)
	rawShare, err := privateShare.MarshalBinary()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "while computing bit reperesentation of secret share")
	}
	if err := verifyMarshalOutputBigEndian(privateShare, rawShare); err != nil {
		return nil, nil, err
	}
	return privateShare, rawShare, nil
}

func verifyMarshalOutputBigEndian(share kyber.Scalar, rawShare []byte) error {
	tot := share.Clone().Zero()
	for bitIdx := 0; bitIdx < len(rawShare)*8; bitIdx++ {
		tot = tot.Clone().Add(tot, tot)
		if rawShare[bitIdx/8]&(1<<(7-(bitIdx%8))) > 0 {
			tot = tot.Clone().Add(tot, tot.Clone().One())
		}
	}
	if !tot.Equal(share) {
		return errors.Errorf("scalars do not marshal to big-endian representation")
	}
	return nil
}
