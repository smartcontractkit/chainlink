package dkg

import (
	"context"
	"io"
	"reflect"
	"sync"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"
	"github.com/smartcontractkit/ocr2vrf/internal/crypto/point_translation"
	"github.com/smartcontractkit/ocr2vrf/internal/dkg/contract"
	"github.com/smartcontractkit/ocr2vrf/internal/pvss"
	dkg_types "github.com/smartcontractkit/ocr2vrf/types"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	kshare "go.dedis.ch/kyber/v3/share"
	"go.dedis.ch/kyber/v3/sign/anon"
)

type dkg struct {
	t player_idx.Int

	lock sync.RWMutex

	selfIdx *player_idx.PlayerIdx

	cfgDgst types.ConfigDigest

	keyID contract.KeyID

	keyConsumer KeyConsumer

	shareSets shareRecords

	myShareRecord *shareRecord

	esk  kyber.Scalar
	epks []kyber.Point
	ssk  kyber.Scalar
	spks []kyber.Point

	encryptionGroup anon.Suite

	translationGroup kyber.Group

	translator point_translation.PubKeyTranslation

	lastKeyData *KeyData

	contract onchainContract

	completed   bool
	dkgComplete func()

	db dkg_types.DKGSharePersistence

	logger commontypes.Logger

	randomness io.Reader
}

var _ types.ReportingPlugin = (*dkg)(nil)

func (a *NewDKGArgs) SanityCheckArgs() error {
	if a.encryptionGroup == nil {
		return errors.Errorf("encryption group not set")
	}
	if a.esk == nil {
		return errors.Errorf("encryption secret key not set")
	}
	if a.ssk == nil {
		return errors.Errorf("signing secret key not set")
	}
	n := len(a.epks)
	if n > int(player_idx.MaxPlayer) {
		return errors.Errorf("too many players")
	}
	if int(a.t) > n {
		return errors.Errorf("threshold can't exceed number of players")
	}
	if _, err := a.selfIdx.Check(); err != nil {
		return errors.Wrap(err, "self index invalid")
	}
	if !a.selfIdx.AtMost(player_idx.Int(n)) {
		return errors.Errorf("player index can't exceed number of players")
	}
	if reflect.TypeOf(a.esk) != reflect.TypeOf(a.encryptionGroup.Scalar()) {
		return errors.Errorf("encryption secret key must come from encryption group")
	}
	exampleEncryptionPK := a.encryptionGroup.Point()
	for _, epk := range a.epks {
		if reflect.TypeOf(epk) != reflect.TypeOf(exampleEncryptionPK) {
			return errors.Errorf(
				"encryption public key must be of type returned by encryptionGroup %v, %T",
				a.encryptionGroup, exampleEncryptionPK,
			)
		}
	}
	if !a.selfIdx.Index(a.epks).(kyber.Point).Equal(a.encryptionGroup.Point().Mul(a.esk, nil)) {
		return errors.Errorf("secret encryption key does not match public encryption key")
	}
	if len(a.spks) != n {
		return errors.Errorf("must be exactly one signing key for each player")
	}
	if reflect.TypeOf(a.ssk) != reflect.TypeOf(a.signingGroup().Scalar()) {
		return errors.Errorf("signing secret key must be of type %T", a.signingGroup().Scalar())
	}
	exampleSigningPK := a.signingGroup().Point()
	for _, spk := range a.spks {
		if reflect.TypeOf(spk) != reflect.TypeOf(exampleSigningPK) {
			return errors.Errorf("signing public keys must be of type %T", exampleSigningPK)
		}
	}
	if !a.selfIdx.Index(a.spks).(kyber.Point).Equal(a.signingGroup().Point().Mul(a.ssk, nil)) {
		return errors.Errorf("secret signing key does not match public signing key")
	}
	if n < pvss.MinPlayers {
		return errors.Errorf("not enough players (need at least %d)", pvss.MinPlayers)
	}
	return nil
}

func (a *NewDKGArgs) signingGroup() anon.Suite {
	if a.xxxTestingOnlySigningGroup != nil {
		return a.xxxTestingOnlySigningGroup
	}
	return SigningGroup
}

func (d *dkg) keyReportedOnchain(ctx context.Context) bool {
	kd, err := d.contract.KeyData(ctx, d.keyID, d.cfgDgst)
	return err == nil && kd.PublicKey != nil && len(kd.Hashes) > 0
}

func (d *dkg) recoverDistributedKeyShare(ctx context.Context) (err error) {
	kd, err := d.contract.KeyData(ctx, d.keyID, d.cfgDgst)
	if err != nil {
		return errors.Wrap(err, "could not get key data while recovering key shares")
	}
	if d.shareSets.allKeysPresent(kd.Hashes) {

		finalShare, err := d.shareSets.recoverDistributedKeyShare(
			d.esk, *d.selfIdx, &kd, d.encryptionGroup, d.cfgDgst,
		)
		if err != nil {
			return errors.Wrap(err, "could not recover distribute key from shares")
		}

		shares, err := d.shareSets.recoverPublicShares(&kd)
		if err != nil {
			return errors.Wrap(err, "could not get public shares to report to consumer")
		}
		players, err := player_idx.PlayerIdxs(player_idx.Int(len(shares)))
		if err != nil {
			return errors.Wrap(err, "could not construct players for pubshares")
		}
		pubShares := make([]kshare.PubShare, len(players))
		for i, playerIdx := range players {

			pubShares[i] = playerIdx.PubShare(shares[i])
		}

		keyData := &KeyData{
			kd.PublicKey,
			pubShares,
			&SecretShare{*d.selfIdx, finalShare.V},
			d.t,
			true,
		}

		d.keyConsumer.NewKey(d.keyID, keyData)
		d.completed = true
		d.dkgComplete()
		return nil
	}
	return errors.Errorf(
		"do not yet have all shares required for reconstruction of given key",
	)
}

var SigningGroup anon.Suite = edwards25519.NewBlakeSHA256Ed25519()

type onchainContract interface {
	KeyData(context.Context, contract.KeyID, types.ConfigDigest) (contract.KeyData, error)
}
