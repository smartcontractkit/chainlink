package dkg

import (
	"context"

	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"
	"github.com/smartcontractkit/ocr2vrf/internal/dkg/contract"
	"github.com/smartcontractkit/ocr2vrf/internal/pvss"
	"github.com/smartcontractkit/ocr2vrf/types/hash"
)

var _ = (&dkg{}).Report

func (d *dkg) recoverShareRecord(report []byte) (r *shareRecord, h hash.Hash, err error) {
	h = hash.GetHash(report)
	r, present := d.shareSets[h]
	if !present {
		r, err = unmarshalSignedShareRecord(d, report)
		if err != nil {
			return nil, hash.Zero, err
		}
		d.shareSets[h] = r
	}
	return r, h, nil
}

type validShareRecords struct {
	includedDealers map[player_idx.PlayerIdx]bool

	includedHashes hash.Hashes

	validShareCount int
	players         []*player_idx.PlayerIdx
	d               *dkg

	aggregatePublicKey kyber.Point

	context context.Context
}

func (d *dkg) newValidShareRecords(
	ctx context.Context,
) (*validShareRecords, error) {
	players, err := player_idx.PlayerIdxs(player_idx.Int(len(d.epks)))
	if err != nil {
		return nil, errors.Wrap(err, "could not construct player list")
	}
	return &validShareRecords{
		map[player_idx.PlayerIdx]bool{},
		hash.MakeHashes(),
		0,
		players,
		d,
		d.translationGroup.Point(),
		ctx,
	}, nil
}

func (v *validShareRecords) validateShareRecord(
	r *shareRecord,
	sender *player_idx.PlayerIdx,
) (reportedDealer *player_idx.PlayerIdx, err error) {
	reportedDealer, err = r.shareSet.Dealer()
	if err != nil {
		return nil, errors.Wrap(
			err, "bad dealer on prospective share record for report",
		)
	}
	if !v.d.keyReportedOnchain(v.context) {
		if v.includedDealers[*reportedDealer] {
			return nil, errors.Errorf(
				"excluding share set from dealer %d, which already has a share set "+
					"incorporated in the current report",
				reportedDealer,
			)
		}
		if !sender.Equal(reportedDealer) && !v.d.keyReportedOnchain(v.context) {
			return nil, errors.Errorf(
				"share-set Observer (%v) and report's self-attributed player (%v) "+
					"indices don't match, during DKG-construction phase",
				sender, reportedDealer,
			)
		}
	}
	return reportedDealer, nil
}

func (v *validShareRecords) storeValidShareSet(
	marshaledShareSet []byte,
	reportedDealer player_idx.PlayerIdx,
	shareSet *pvss.ShareSet,
	h *hash.Hash,
) {
	v.includedDealers[reportedDealer] = true

	agg := v.aggregatePublicKey.Clone()
	_ = v.aggregatePublicKey.Add(agg, shareSet.PublicKey())
	v.includedHashes.Add(*h)
	go v.persistShares(reportedDealer, marshaledShareSet, *h)
	v.validShareCount++
}

func (v *validShareRecords) processShareSet(aobs types.AttributedObservation) {
	if int(aobs.Observer) >= len(v.players) {
		v.d.logger.Debug("observer index out of range", commontypes.LogFields{
			"observer index": aobs.Observer, "max index": len(v.players) - 1,
		})
		return
	}
	sender := v.players[aobs.Observer]

	r, h, err := v.d.recoverShareRecord(aobs.Observation)
	if err != nil {
		v.d.logger.Warn("excluding invalid share set from report",
			commontypes.LogFields{"err": err, "sender": sender})
		return
	}

	reportedDealer, err := v.validateShareRecord(r, sender)
	if err != nil {
		v.d.logger.Warn("invalid share set", commontypes.LogFields{"err": err})
		return
	}
	v.storeValidShareSet(aobs.Observation, *reportedDealer, r.shareSet, &h)
}

func (v *validShareRecords) enoughShareSets() bool {
	return v.validShareCount > int(v.d.t)
}

func (v *validShareRecords) report() (rv []byte, err error) {

	kd := &contract.KeyData{v.aggregatePublicKey, v.includedHashes}
	kb, err := kd.MarshalBinary(v.d.keyID)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal key for onchain report")
	}
	return kb, nil
}

func unmarshalSignedShareRecord(d *dkg, report []byte) (*shareRecord, error) {
	r, rem, err := unmarshalShareRecord(
		SigningGroup, d.encryptionGroup, d.translationGroup,
		report, d.translator, d.cfgDgst, d.epks, d.spks,
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal unknown share record")
	}
	if len(rem) > 0 {
		return nil, errors.Errorf(
			"overage of %d bytes in marshalled share record: 0x%x", len(rem), rem,
		)
	}
	return r, nil
}
