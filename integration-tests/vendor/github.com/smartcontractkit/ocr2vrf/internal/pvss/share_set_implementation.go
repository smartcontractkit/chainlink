package pvss

import (
	"bytes"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"
	"github.com/smartcontractkit/ocr2vrf/internal/crypto/point_translation"
	"github.com/smartcontractkit/ocr2vrf/internal/util"

	"go.dedis.ch/kyber/v3"
	kshare "go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/sign/anon"
)

var MinPlayers = 5

var _ = NewShareSet

func newShareSet(domainSep types.ConfigDigest, threshold player_idx.Int, dealerIdx *player_idx.PlayerIdx,
	group anon.Suite, translation point_translation.PubKeyTranslation,
	pks []kyber.Point, toxicWasteNeverUseThisParam ...interface{},
) (*ShareSet, error) {
	if len(pks) < MinPlayers {
		return nil, errors.Errorf("%d is not enough players", len(pks))
	}
	if int(threshold) >= len(pks) {
		return nil, errors.Errorf(
			"threshold %d cannot exceed number of players %d", threshold, len(pks),
		)
	}
	players, err := player_idx.PlayerIdxs(player_idx.Int(len(pks)))
	if err != nil {
		return nil, err
	}
	secretPoly := kshare.NewPriPoly(group, int(threshold)+1, nil, group.RandomStream())
	coeffCommits := secretPoly.Commit(nil)
	pvssKey, err := translation.TranslateKey(secretPoly.Secret())
	if err != nil {
		return nil, errors.Wrapf(err, "could not translate dealer's additive share")
	}
	rv := &ShareSet{
		dealerIdx, coeffCommits, pvssKey, make([]*share, len(pks)), translation, nil,
	}

	if len(toxicWasteNeverUseThisParam) == 1 && toxicWasteNeverUseThisParam[0] ==
		"⚠⚠⚠☣☢️☠ Yes, please inject me with that yummy plutonium-salt solution. "+
			"I long for the sweet release of a lingering death by radiation "+
			"sickness ☠☢️☣️⚠⚠⚠" {
		rv.xXXToxicWaste = secretPoly
	}
	edomain, err := rv.domainSep(domainSep)
	if err != nil {
		return nil, err
	}
	for shareIdx, receiver := range players {
		pk := receiver.Index(pks).(kyber.Point)
		rv.shares[shareIdx], err = newShare(
			edomain, group, secretPoly, rv, receiver, pk,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "while constructing share set")
		}
	}
	return rv, nil
}

var _ = (&ShareSet{}).Verify

func (s *ShareSet) verify(group anon.Suite, domainSep types.ConfigDigest, pks []kyber.Point) error {
	if _, commits := s.coeffCommitments.Info(); len(commits) < 1 {
		return errors.Errorf("need at least one coefficient commitment in a valid share")
	}
	if err := s.translation.VerifyTranslation(s.coeffCommitments.Commit(), s.pvssKey); err != nil {
		return errors.Wrapf(err, "bad translation of additive key share in share set")
	}
	numPlayers := len(pks)
	if numPlayers > int(player_idx.MaxPlayer) {
		return errors.Errorf("Can't handle %d players; %d is max", len(pks), player_idx.MaxPlayer)
	}
	edomain, err := s.domainSep(domainSep)
	if err != nil {
		return errors.Wrapf(err, "failed to construct domain separator for PVSS proofs")
	}
	players, err := player_idx.PlayerIdxs(player_idx.Int(numPlayers))
	if err != nil {
		return util.WrapError(err, "could not get list of player indices in ShareSet.verify")
	}
	for shareIdx, share := range s.shares {
		p := players[shareIdx]
		if err := share.verify(group, append(edomain, p.Marshal()...), p); err != nil {
			return errors.Wrapf(err, "could not verify every share in share set")
		}
	}
	return nil
}

func (s *ShareSet) domainSep(domainSep types.ConfigDigest) ([]byte, error) {
	_, commits := s.coeffCommitments.Info()
	components := make([][]byte, len(commits)+2)
	components[0] = domainSep[:]
	offset := 1
	var err error
	for ci, c := range commits {
		components[ci+offset], err = c.MarshalBinary()
		if err != nil {
			return nil, errors.Errorf(
				"while including secret polynomial coefficient commits in domain separator",
			)
		}
	}

	components[len(components)-1] = s.dealer.Marshal()
	return bytes.Join(components, nil), nil
}

func (s *ShareSet) Equal(s2 *ShareSet) bool {
	if s == nil || s2 == nil {
		return false
	}
	if !s.dealer.Equal(s2.dealer) {
		return false
	}
	_, scommits := s.coeffCommitments.Info()
	_, s2commits := s2.coeffCommitments.Info()
	if len(scommits) == len(s2commits) &&
		s.coeffCommitments.Equal(s2.coeffCommitments) &&
		s.pvssKey.Equal(s2.pvssKey) &&
		(len(s.shares) == len(s2.shares)) {
		for i, sh := range s.shares {
			if !sh.equal(s2.shares[i]) {
				return false
			}
		}
		return true
	}
	return false
}

var _ = (*ShareSet)(nil).Decrypt

func (s *ShareSet) decrypt(
	playerIdx player_idx.PlayerIdx, sk kyber.Scalar,
	keyGroup anon.Suite, domainSep types.ConfigDigest, sharePublicCommitment kyber.Point,
) (kshare.PriShare, error) {
	playerShare := playerIdx.Index(s.shares).(*share)
	edomain, err := s.domainSep(domainSep)
	if err != nil {
		return kshare.PriShare{}, err
	}
	return playerShare.decrypt(
		sk, keyGroup, edomain, sharePublicCommitment, playerIdx,
	)
}

var _ = (*ShareSet)(nil).PublicShares

func (s *ShareSet) publicShares() (rv []kyber.Point) {
	for _, share := range s.shares {
		rv = append(rv, share.subKeyTranslation)
	}
	return rv
}
