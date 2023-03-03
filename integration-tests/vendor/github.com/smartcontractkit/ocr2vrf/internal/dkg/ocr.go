package dkg

import (
	"context"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

func (d *dkg) Query(context.Context, types.ReportTimestamp,
) (types.Query, error) {
	return nil, nil
}

func (d *dkg) Observation(
	ctx context.Context, _ types.ReportTimestamp, _ types.Query,
) (o types.Observation, err error) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	var respondingShareRecord *shareRecord
	if !d.keyReportedOnchain(ctx) {
		respondingShareRecord = d.myShareRecord
	} else {

		keyData, err2 := d.contract.KeyData(ctx, d.keyID, d.cfgDgst)
		if err2 != nil {

			return nil, errors.Wrap(err2, "digest marked as complete, but key data is unavailable")
		}
		respondingShareRecordHash, err2 := d.shareSets.getRandom(keyData.Hashes, d.randomness)
		if err2 != nil {

			d.logger.Error(
				"could not choose random hash to send for data-availability phase",
				commontypes.LogFields{
					"required hashes":     keyData.Hashes,
					"existing share sets": d.shareSets,
				},
			)
			return types.Observation{}, nil
		}
		var ok bool
		respondingShareRecord, ok = d.shareSets[respondingShareRecordHash]
		if !ok {
			d.logger.Error(
				"could not choose random share set to send. No record for the hash.",
				commontypes.LogFields{
					"randomly chosen hash": respondingShareRecordHash,
				},
			)
			return types.Observation{}, nil
		}
	}
	o, err = respondingShareRecord.marshal()
	if err != nil {
		return nil, errors.Wrap(err, "could not construct observation")
	}
	if _, err := unmarshalSignedShareRecord(d, o); err != nil {
		panic(err)
	}
	return o, nil
}

func (d *dkg) Report(
	ctx context.Context, _ types.ReportTimestamp, _ types.Query,
	shares []types.AttributedObservation,
) (shouldReport bool, report types.Report, err error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	v, err := d.newValidShareRecords(ctx)
	if err != nil {
		return false, nil, errors.Wrap(
			err, "could not create record for valid shares",
		)
	}
	for _, aobs := range shares {
		senderField := commontypes.LogFields{"sender": aobs.Observer}
		d.logger.Debug("processing share set", senderField)

		v.processShareSet(aobs)
	}
	if d.keyReportedOnchain(ctx) {

		return false, nil, d.recoverDistributedKeyShare(ctx)
	}

	if !v.enoughShareSets() {
		d.logger.Warn(
			"need quorum of unique share sets to construct secure distributed key",
			commontypes.LogFields{
				"required": d.t + 1,
				"received": v.validShareCount,
				"players":  v.players,
			},
		)
		return false, nil, nil
	}
	report, err = v.report()
	if err != nil {
		return false, nil, errors.Wrap(err,
			"could not extract onchain report from share set we just constructed",
		)
	}
	return true, report, nil
}

func (d *dkg) ShouldAcceptFinalizedReport(
	_ context.Context, _ types.ReportTimestamp, _ types.Report) (bool, error) {

	return true, nil
}

func (d *dkg) ShouldTransmitAcceptedReport(
	c context.Context, t types.ReportTimestamp, r types.Report) (bool, error) {

	return !d.keyReportedOnchain(c), nil
}

func (d *dkg) Start() error {
	return nil
}

func (d *dkg) Close() error {
	return nil
}
