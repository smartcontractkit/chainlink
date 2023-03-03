package dkg

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"
	"github.com/smartcontractkit/ocr2vrf/internal/pvss"
	"github.com/smartcontractkit/ocr2vrf/internal/util"
	"github.com/smartcontractkit/ocr2vrf/types"
	"github.com/smartcontractkit/ocr2vrf/types/hash"

	"go.dedis.ch/kyber/v3/sign/anon"
)

func (sr shareRecord) PersistentRecord() (types.PersistentShareSetRecord, error) {
	marshaled, err := sr.marshal()
	if err != nil {
		errMsg := "could not marshal share records for local persistence"
		return types.PersistentShareSetRecord{},
			util.WrapError(err, errMsg)
	}
	dealer, err := sr.shareSet.Dealer()
	if err != nil {
		errMsg := "could not marshal share record for local persistence"
		return types.PersistentShareSetRecord{},
			util.WrapError(err, errMsg)
	}
	rv := types.PersistentShareSetRecord{
		*dealer,
		marshaled,
		hash.Hash{},
	}
	return rv, nil
}

func (d *dkg) initializeShareSets(signingGroup anon.Suite) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	persistedShares, err := d.db.ReadShareRecords(d.cfgDgst, d.keyID)
	if err != nil {
		return util.WrapError(err, "could not recover persisted shares")
	}
	myShareRecovered := false
	for _, storedShare := range persistedShares {
		m := storedShare.MarshaledShareRecord

		mHash := hash.GetHash(m)
		if storedShare.Hash != mHash {
			dealer := storedShare.Dealer
			return fmt.Errorf("hash mismatch on record from %d", dealer)
		}
		share, err := unmarshalSignedShareRecord(d, m)
		if err != nil {
			errMsg := "could not unmarshal persisted share record"
			return util.WrapError(err, errMsg)
		}
		if err := d.shareSets.set(share, mHash); err != nil {
			errMsg := "could not record persisted share from %s"
			return util.WrapErrorf(err, errMsg, storedShare.Dealer)
		}
		if storedShare.Dealer.Equal(d.selfIdx) {
			myShareRecovered = true
			d.myShareRecord = share
		}
	}
	if !myShareRecovered {

		shareSet, err := pvss.NewShareSet(
			d.cfgDgst, d.t, d.selfIdx, d.encryptionGroup, d.translator, d.epks,
		)
		if err != nil {
			return util.WrapError(err, "could not create own share set")
		}
		msr, err := newShareRecord(signingGroup, shareSet, d.ssk, d.cfgDgst)
		if err != nil {
			return util.WrapError(err, "could not create own share record")
		}
		d.myShareRecord = msr
		psr, err := d.myShareRecord.PersistentRecord()
		if err != nil {
			errMsg := "could not construct own persistent share record"
			return util.WrapError(err, errMsg)
		}
		lpsr := []types.PersistentShareSetRecord{psr}
		err = d.db.WriteShareRecords(context.Background(), d.cfgDgst, d.keyID, lpsr)
		if err != nil {
			return util.WrapError(err, "could not write own share record")
		}
		err = d.shareSets.set(d.myShareRecord, hash.Zero)
		if err != nil {
			return util.WrapError(err, "could not set own share record")
		}
	}
	return nil
}

func (v *validShareRecords) persistShares(
	reportedDealer player_idx.PlayerIdx,
	marshaledShareSet []byte,
	h hash.Hash,
) {
	psr := types.PersistentShareSetRecord{reportedDealer, marshaledShareSet, h}
	lpsr := []types.PersistentShareSetRecord{psr}
	v.d.lock.Lock()
	defer v.d.lock.Unlock()
	err := v.d.db.WriteShareRecords(v.context, v.d.cfgDgst, v.d.keyID, lpsr)
	if err != nil {
		v.d.logger.Warn("failed to persist share", commontypes.LogFields{
			"player":     reportedDealer,
			"err":        err,
			"share hash": h,
		})
	}
}
