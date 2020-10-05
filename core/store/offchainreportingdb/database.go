package offchainreportingdb

import (
	"context"
	"database/sql"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
	ocrtypes "github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

type (
	db struct {
		*sql.DB
		oracleSpecID int
	}
)

// NewDB returns a new DB scoped to this oracleSpecID
func NewDB(sqldb *sql.DB, oracleSpecID int) ocrtypes.Database {
	return &db{sqldb, oracleSpecID}
}

func (d *db) ReadState(ctx context.Context, cd ocrtypes.ConfigDigest) (ps *ocrtypes.PersistentState, err error) {
	q := d.QueryRowContext(ctx, `
SELECT epoch, highest_sent_epoch, highest_received_epoch
FROM offchainreporting_persistent_states
WHERE offchainreporting_oracle_spec_id = $1 AND config_digest = $2
LIMIT 1`, d.oracleSpecID, cd)

	ps = new(ocrtypes.PersistentState)

	var tmp []int64
	err = q.Scan(&ps.Epoch, &ps.HighestSentEpoch, pq.Array(&tmp))

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "ReadState failed")
	}

	for _, v := range tmp {
		ps.HighestReceivedEpoch = append(ps.HighestReceivedEpoch, uint32(v))
	}

	return ps, nil
}

func (d *db) WriteState(ctx context.Context, cd ocrtypes.ConfigDigest, state ocrtypes.PersistentState) error {
	var highestReceivedEpoch []int64
	for _, v := range state.HighestReceivedEpoch {
		highestReceivedEpoch = append(highestReceivedEpoch, int64(v))
	}
	_, err := d.ExecContext(ctx, `
INSERT INTO offchainreporting_persistent_states (offchainreporting_oracle_spec_id, config_digest, epoch, highest_sent_epoch, highest_received_epoch, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
ON CONFLICT (offchainreporting_oracle_spec_id, config_digest) DO UPDATE SET
	(epoch, highest_sent_epoch, highest_received_epoch, updated_at)
	=
	(
	 EXCLUDED.epoch,
	 EXCLUDED.highest_sent_epoch,
	 EXCLUDED.highest_received_epoch,
	 NOW()
	)
`, d.oracleSpecID, cd, state.Epoch, state.HighestSentEpoch, pq.Array(&highestReceivedEpoch))

	return errors.Wrap(err, "WriteState failed")
}

func (d *db) ReadConfig(ctx context.Context) (c *ocrtypes.ContractConfig, err error) {
	q := d.QueryRowContext(ctx, `
SELECT config_digest, signers, transmitters, threshold, encoded_config_version, encoded
FROM offchainreporting_contract_configs
WHERE offchainreporting_oracle_spec_id = $1
LIMIT 1`, d.oracleSpecID)

	c = new(ocrtypes.ContractConfig)

	var signers [][]byte
	var transmitters [][]byte

	err = q.Scan(&c.ConfigDigest, (*pq.ByteaArray)(&signers), (*pq.ByteaArray)(&transmitters), &c.Threshold, &c.EncodedConfigVersion, &c.Encoded)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "ReadConfig failed")
	}

	for _, s := range signers {
		c.Signers = append(c.Signers, common.BytesToAddress(s))
	}
	for _, t := range transmitters {
		c.Transmitters = append(c.Transmitters, common.BytesToAddress(t))
	}

	return
}

func (d *db) WriteConfig(ctx context.Context, c ocrtypes.ContractConfig) error {
	var signers [][]byte
	var transmitters [][]byte
	for _, s := range c.Signers {
		signers = append(signers, s.Bytes())
	}
	for _, t := range c.Transmitters {
		transmitters = append(transmitters, t.Bytes())
	}
	_, err := d.ExecContext(ctx, `
INSERT INTO offchainreporting_contract_configs (offchainreporting_oracle_spec_id, config_digest, signers, transmitters, threshold, encoded_config_version, encoded, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
ON CONFLICT (offchainreporting_oracle_spec_id) DO UPDATE SET
	config_digest = EXCLUDED.config_digest,
	signers = EXCLUDED.signers,
	transmitters = EXCLUDED.transmitters,
	threshold = EXCLUDED.threshold,
	encoded_config_version = EXCLUDED.encoded_config_version,
	encoded = EXCLUDED.encoded,
	updated_at = NOW()
`, d.oracleSpecID, c.ConfigDigest, pq.ByteaArray(signers), pq.ByteaArray(transmitters), c.Threshold, int(c.EncodedConfigVersion), c.Encoded)

	return errors.Wrap(err, "WriteConfig failed")
}

func (d *db) StorePendingTransmission(ctx context.Context, k ocrtypes.PendingTransmissionKey, p ocrtypes.PendingTransmission) error {
	median := utils.NewBig(p.Median)
	var rs [][]byte
	var ss [][]byte
	for _, v := range p.Rs {
		rs = append(rs, v[:])
	}
	for _, v := range p.Ss {
		ss = append(ss, v[:])
	}

	_, err := d.ExecContext(ctx, `
INSERT INTO offchainreporting_pending_transmissions (
	offchainreporting_oracle_spec_id,
	config_digest,
	epoch,
	round,
	time,
	median,
	serialized_report,
	rs,
	ss,
	vs,
	created_at,
	updated_at
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,NOW(),NOW())
ON CONFLICT (offchainreporting_oracle_spec_id, config_digest, epoch, round) DO UPDATE SET
	time = EXCLUDED.time,
	median = EXCLUDED.median,
	serialized_report = EXCLUDED.serialized_report,
	rs = EXCLUDED.rs,
	ss = EXCLUDED.ss,
	vs = EXCLUDED.vs,
	updated_at = NOW()
`, d.oracleSpecID, k.ConfigDigest, k.Epoch, k.Round, p.Time, median, p.SerializedReport, pq.ByteaArray(rs), pq.ByteaArray(ss), p.Vs[:])

	return errors.Wrap(err, "StorePendingTransmission failed")
}

func (d *db) PendingTransmissionsWithConfigDigest(ctx context.Context, cd ocrtypes.ConfigDigest) (map[ocrtypes.PendingTransmissionKey]ocrtypes.PendingTransmission, error) {
	rows, err := d.QueryContext(ctx, `
SELECT 
	config_digest,
	epoch,
	round,
	time,
	median,
	serialized_report,
	rs,
	ss,
	vs
FROM offchainreporting_pending_transmissions
WHERE offchainreporting_oracle_spec_id = $1 AND config_digest = $2
`, d.oracleSpecID, cd)
	if err != nil {
		return nil, errors.Wrap(err, "PendingTransmissionsWithConfigDigest failed to query rows")
	}
	defer logger.ErrorIfCalling(rows.Close)

	m := make(map[ocrtypes.PendingTransmissionKey]ocrtypes.PendingTransmission)

	for rows.Next() {
		k := ocrtypes.PendingTransmissionKey{}
		p := ocrtypes.PendingTransmission{}

		var median utils.Big
		var rs [][]byte
		var ss [][]byte
		var vs []byte
		if err := rows.Scan(&k.ConfigDigest, &k.Epoch, &k.Round, &p.Time, &median, &p.SerializedReport, (*pq.ByteaArray)(&rs), (*pq.ByteaArray)(&ss), &vs); err != nil {
			return nil, errors.Wrap(err, "PendingTransmissionsWithConfigDigest failed to scan row")
		}
		p.Median = median.ToInt()
		for i, v := range rs {
			var r [32]byte
			if n := copy(r[:], v); n != 32 {
				return nil, errors.Errorf("expected 32 bytes for rs value at index %v, got %v bytes", i, n)
			}
			p.Rs = append(p.Rs, r)
		}
		for i, v := range ss {
			var s [32]byte
			if n := copy(s[:], v); n != 32 {
				return nil, errors.Errorf("expected 32 bytes for ss value at index %v, got %v bytes", i, n)
			}
			p.Ss = append(p.Ss, s)
		}
		if n := copy(p.Vs[:], vs); n != 32 {
			return nil, errors.Errorf("expected 32 bytes for vs, got %v bytes", n)
		}
		m[k] = p
	}

	return m, nil
}

func (d *db) DeletePendingTransmission(ctx context.Context, k ocrtypes.PendingTransmissionKey) (err error) {
	_, err = d.ExecContext(ctx, `
DELETE FROM offchainreporting_pending_transmissions
WHERE offchainreporting_oracle_spec_id = $1 AND  config_digest = $2 AND epoch = $3 AND round = $4
`, d.oracleSpecID, k.ConfigDigest, k.Epoch, k.Round)

	err = errors.Wrap(err, "DeletePendingTransmission failed")

	return
}

func (d *db) DeletePendingTransmissionsOlderThan(ctx context.Context, t time.Time) (err error) {
	_, err = d.ExecContext(ctx, `
DELETE FROM offchainreporting_pending_transmissions
WHERE offchainreporting_oracle_spec_id = $1 AND time < $2
`, d.oracleSpecID, t)

	err = errors.Wrap(err, "DeletePendingTransmissionsOlderThan failed")

	return
}
