package pvss

import (
	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"
	"github.com/smartcontractkit/ocr2vrf/internal/crypto/point_translation"

	"go.dedis.ch/kyber/v3"
	kshare "go.dedis.ch/kyber/v3/share"
)

type ShareSet struct {
	dealer *player_idx.PlayerIdx

	coeffCommitments *kshare.PubPoly

	pvssKey kyber.Point

	shares []*share

	translation point_translation.PubKeyTranslation

	xXXToxicWaste *kshare.PriPoly
}
