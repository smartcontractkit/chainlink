package pvss

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/sign/anon"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"
	"github.com/smartcontractkit/ocr2vrf/internal/crypto/point_translation"
	"github.com/smartcontractkit/ocr2vrf/internal/util"
)

var _ = (&ShareSet{}).Marshal

func (s *ShareSet) marshal() (m []byte, err error) {
	if s == nil {
		return nil, errors.Errorf("attempt to marshal non-existent share set")
	}
	rv := make([][]byte, 4+len(s.shares))
	cursor := 0

	if s.dealer == nil {
		return nil, errors.Errorf("can't marshal share set with no dealer specified")
	}
	rv[cursor] = s.dealer.Marshal()
	cursor++

	rv[cursor], err = (&pubPoly{s.coeffCommitments}).marshal()
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal coefficient commitments")
	}
	cursor++

	rv[cursor], err = marshalKyberPointWithLen(s.pvssKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal translated PVSS public key")
	}
	cursor++

	if len(s.shares) > int(player_idx.MaxPlayer) {
		return nil, errors.Errorf("too many shares to marshal")
	}
	rv[cursor] = player_idx.RawMarshal(player_idx.Int(len(s.shares)))
	cursor++

	for _, sh := range s.shares {
		rv[cursor], err = sh.marshal()
		if err != nil {
			return nil, errors.Wrap(err, "could not marshal share in share-set")
		}
		cursor++
	}

	if cursor != len(rv) {
		return nil, errors.Errorf(
			"ShareSet marshal fields out of registration: cursor: %d, fields: %d",
			cursor, len(rv),
		)
	}
	return bytes.Join(rv, nil), nil
}

var _ = UnmarshalShareSet

func unmarshalShareSet(
	g anon.Suite, translationGroup kyber.Group, data []byte,
	translation point_translation.PubKeyTranslation, domainSep types.ConfigDigest,
	pks []kyber.Point,
) (ss *ShareSet, rem []byte, err error) {

	dealer, data, err := unmarshalDealer(data)
	if err != nil {
		return nil, nil, err
	}

	coeffCommitments, data, err := unmarshalPubPoly(g, data)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not unmarshal coefficient commitments for share set")
	}

	pvssKeyB, data, err := readLenPrefixedBytes(data, 1)
	if err != nil {
		return nil, nil, util.WrapError(err, "could not read PVSS key bytes")
	}
	pvssKey := translationGroup.Point()
	if err2 := pvssKey.UnmarshalBinary(pvssKeyB); err2 != nil {
		return nil, nil, errors.Wrap(err2, "could not read translated PVSS public key")
	}
	numShares, data, err := player_idx.RawUnmarshal(data)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get number of shares in marshalled share set")
	}

	shareSpots := make([]*share, numShares)
	ss = &ShareSet{
		dealer, coeffCommitments.PubPoly, pvssKey, shareSpots, translation, nil,
	}
	for i := range ss.shares {
		ss.shares[i], data, err = unmarshal(g, translationGroup, data, ss)
		if err != nil {
			return nil, nil, errors.Wrap(err, "could not unmarshal shares in share set")
		}
	}
	if err := ss.verify(g, domainSep, pks); err != nil {
		return nil, nil, errors.Wrap(err, "unmarshaled to invalid share set")
	}
	return ss, data, nil
}

var _ = UnmarshalDealer

func unmarshalDealer(data []byte) (*player_idx.PlayerIdx, []byte, error) {
	dealer, data, err := player_idx.Unmarshal(data)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not unmarshal share set dealer")
	}
	if err := dealer.NonZero(); err != nil {
		return nil, nil, errors.Wrap(err, "dealer index for marshalled share set is zero")
	}
	return dealer, data, nil
}
