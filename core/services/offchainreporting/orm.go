package offchainreporting

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"

	ocrtypes "github.com/smartcontractkit/offchain-reporting-design/prototype/offchainreporting/types"
)

// database is an abstraction that conforms to the Database interface in the
// offchain reporting prototype, which is unaware of job IDs.
type orm struct {
	db    *gorm.DB
	jobID int32
}

var _ ocrtypes.Database = orm{}

type persistentState struct {
	JobSpecID                int32 `gorm:"primary_key"`
	ConfigDigest             ocrtypes.ConfigDigest
	ocrtypes.PersistentState `gorm:"embedded"`
}

func (o orm) ReadState(configDigest ocrtypes.ConfigDigest) (*ocrtypes.PersistentState, error) {
	var state persistentState
	err := o.db.
		Where("job_spec_id = ? AND config_digest = ?", o.jobID, configDigest).
		Find(&state).
		Error
	if err != nil {
		return nil, err
	}
	return &state.PersistentState, nil
}

func (o orm) WriteState(configDigest ocrtypes.ConfigDigest, state ocrtypes.PersistentState) error {
	return o.db.Exec(`
        INSERT INTO offchain_reporting_persistent_states (
            job_spec_id, group_id, epoch, highest_sent_epoch, highest_received_epoch
        ) VALUES (
            ?, ?, ?, ?, ?
        ) ON CONFLICT (job_spec_id, config_digest) DO
        UPDATE SET epoch = EXCLUDED.epoch, highest_sent_epoch = EXCLUDED.highest_sent_epoch,
            highest_received_epoch = EXCLUDED.highest_received_epoch
        WHERE job_spec_id = ? AND config_digest = ?
    `, o.jobID, configDigest, state.Epoch, state.HighestSentEpoch, state.HighestReceivedEpoch,
		state.Epoch, state.HighestSentEpoch, state.HighestReceivedEpoch,
		o.jobID, configDigest,
	).Error
}

type contractConfig struct {
	JobSpecID    int32                 `gorm:"primary_key"`
	ConfigDigest ocrtypes.ConfigDigest `gorm:"type:bytea"`
	Signers      []common.Address      `gorm:"type:bytea[]"`
	Transmitters []common.Address      `gorm:"type:bytea[]"`
	Threshold    uint8
	Encoded      []byte
}

func (o orm) ReadConfig() (*ocrtypes.ContractConfig, error) {
	var config contractConfig
	err := o.db.
		Where("job_spec_id = ?", o.jobID).
		Find(&config).
		Error
	if err != nil {
		return nil, err
	}
	return &ocrtypes.ContractConfig{
		ConfigDigest: config.ConfigDigest,
		Signers:      config.Signers,
		Transmitters: config.Transmitters,
		Threshold:    config.Threshold,
		Encoded:      config.Encoded,
	}, nil
}

func (o orm) WriteConfig(config ocrtypes.ContractConfig) error {
	err := o.db.Exec(`
        INSERT INTO offchain_reporting_configs (
            job_spec_id, config_digest, signers, transmitters, threshold, encoded
        ) VALUES (
            ?, ?, ?, ?, ?, ?
        )
        ON CONFLICT (job_spec_id, config_digest) DO
        UPDATE offchain_reporting_configs
        SET signers = EXCLUDED.signers, transmitters = EXCLUDED.transmitters, threshold = EXCLUDED.threshold,
            encoded = EXCLUDED.encoded
        WHERE job_spec_id = ? AND config_digest = ?
    `, o.jobID, config.ConfigDigest, config.Signers, config.Transmitters, config.Threshold, config.Encoded,
		o.jobID, config.ConfigDigest,
	).Error
	if err != nil {
		return err
	}
	return nil
}

func (o orm) StorePendingTransmission(
	key ocrtypes.PendingTransmissionKey,
	t ocrtypes.PendingTransmission,
) error {
	return o.db.Exec(`
        INSERT INTO offchain_reporting_pending_transmissions (
            job_spec_id, config_digest, epoch, round, time, median, serialized_report, rs, ss, vs
        ) VALUES (
            ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
        )
        ON CONFLICT (job_spec_id, config_digest, epoch, round) DO
        UPDATE offchain_reporting_pending_transmissions
        SET job_spec_id = EXCLUDED.job_spec_id, config_digest = EXCLUDED.config_digest, epoch = EXCLUDED.epoch
            round = EXCLUDED.round, time = EXCLUDED.time, median = EXCLUDED.median,
            serialized_report = EXCLUDED.serialized_report, rs = EXCLUDED.rs, ss = EXCLUDED.ss, vs = EXCLUDED.vs
        WHERE job_spec_id = ? AND config_digest = ? AND epoch = ? AND round = ?
    `, o.jobID, key.ConfigDigest, key.Epoch, key.Round, t.Time, t.Median, t.SerializedReport, t.Rs,
		t.Ss, t.Vs, o.jobID, key.ConfigDigest, key.Epoch, key.Round).
		Error
}

func (o orm) PendingTransmissionsWithConfigDigest(
	configDigest ocrtypes.ConfigDigest,
) (map[ocrtypes.PendingTransmissionKey]ocrtypes.PendingTransmission, error) {
	type pendingTransmission struct {
		ConfigDigest ocrtypes.ConfigDigest
		Epoch        int
		Round        int

		Time             time.Time
		Median           ocrtypes.Observation
		SerializedReport []byte
		Rs               [][32]byte `gorm:"type:bytea[]"`
		Ss               [][32]byte `gorm:"type:bytea[]"`
		Vs               [32]byte
	}

	var pendingTransmissions []pendingTransmission
	err := o.db.
		Where("job_spec_id = ? AND config_digest = ?").
		Find(&pendingTransmissions).
		Error
	if err != nil {
		return nil, err
	}

	m := make(map[ocrtypes.PendingTransmissionKey]ocrtypes.PendingTransmission)
	for _, t := range pendingTransmissions {
		key := ocrtypes.PendingTransmissionKey{
			ConfigDigest: configDigest,
			Epoch:        t.Epoch,
			Round:        t.Round,
		}
		m[key] = ocrtypes.PendingTransmission{
			Time:             t.Time,
			Median:           t.Median,
			SerializedReport: t.SerializedReport,
			Rs:               t.Rs,
			Ss:               t.Ss,
			Vs:               t.Vs,
		}
	}
	return m, nil
}

func (o orm) DeletePendingTransmission(key ocrtypes.PendingTransmissionKey) error {
	return o.db.Exec(`
        DELETE FROM offchain_reporting_pending_transmissions
        WHERE job_id = ? AND config_digest = ? AND epoch = ? AND round = ?
    `, o.jobID, key.ConfigDigest, key.Epoch, key.Round).Error
}

func (o orm) DeletePendingTransmissionsOlderThan(t time.Time) error {
	return o.db.Exec(`
        DELETE FROM offchain_reporting_pending_transmissions
        WHERE job_spec_id = ? AND time < ?
    `, o.jobID, t).Error
}
