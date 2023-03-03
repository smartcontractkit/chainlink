package pvss

import (
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"
	"github.com/smartcontractkit/ocr2vrf/internal/crypto/point_translation"

	"go.dedis.ch/kyber/v3"
	kshare "go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/sign/anon"
)

func NewShareSet(domainSep types.ConfigDigest, threshold player_idx.Int,
	dealer *player_idx.PlayerIdx, group anon.Suite,
	translation point_translation.PubKeyTranslation, pks []kyber.Point,
) (*ShareSet, error) {
	return newShareSet(domainSep, threshold, dealer, group, translation, pks)
}

func (s *ShareSet) Marshal() (m []byte, err error) {
	return s.marshal()
}

func UnmarshalShareSet(
	g anon.Suite, translationGroup kyber.Group, data []byte,
	translation point_translation.PubKeyTranslation, domainSep types.ConfigDigest,
	pks []kyber.Point,
) (ss *ShareSet, rem []byte, err error) {
	return unmarshalShareSet(g, translationGroup, data, translation, domainSep, pks)
}

func (s *ShareSet) Verify(group anon.Suite, domainSep types.ConfigDigest, pks []kyber.Point) error {
	return s.verify(group, domainSep, pks)
}

func (s *ShareSet) Decrypt(
	playerIdx player_idx.PlayerIdx,
	sk kyber.Scalar,
	keyGroup anon.Suite,
	domainSep types.ConfigDigest,
) (kshare.PriShare, error) {
	sharePublicCommitment := playerIdx.EvalPoint(s.coeffCommitments)
	return s.decrypt(
		playerIdx, sk, keyGroup, domainSep, sharePublicCommitment,
	)
}

func (s *ShareSet) Dealer() (*player_idx.PlayerIdx, error) {
	return s.dealer.Check()
}

func UnmarshalDealer(data []byte) (*player_idx.PlayerIdx, error) {
	rv, _, err := unmarshalDealer(data)
	return rv, err
}

func (s *ShareSet) PublicKey() kyber.Point {
	return s.pvssKey
}

func (s *ShareSet) PublicShares() []kyber.Point {
	return s.publicShares()
}
