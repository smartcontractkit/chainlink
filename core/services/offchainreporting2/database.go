package offchainreporting2

import (
	"context"
	"database/sql"
	"encoding/binary"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	ocrcommon "github.com/smartcontractkit/libocr/commontypes"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type db struct {
	*sql.DB
	oracleSpecID int32
	lggr         logger.Logger
}

var (
	_ ocrtypes.Database = &db{}
)

// NewDB returns a new DB scoped to this oracleSpecID
func NewDB(sqldb *sql.DB, oracleSpecID int32, lggr logger.Logger) *db {
	return &db{sqldb, oracleSpecID, lggr}
}

func (d *db) ReadState(ctx context.Context, cd ocrtypes.ConfigDigest) (ps *ocrtypes.PersistentState, err error) {
	q := d.QueryRowContext(ctx, `
SELECT epoch, highest_sent_epoch, highest_received_epoch
FROM offchainreporting2_persistent_states
WHERE offchainreporting2_oracle_spec_id = $1 AND config_digest = $2
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
INSERT INTO offchainreporting2_persistent_states (
	offchainreporting2_oracle_spec_id,
	config_digest,
	epoch,
	highest_sent_epoch,
	highest_received_epoch,
	created_at,
	updated_at
)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
ON CONFLICT (offchainreporting2_oracle_spec_id, config_digest)
DO UPDATE SET (
		epoch,
		highest_sent_epoch,
		highest_received_epoch,
		updated_at
	) = (
	 EXCLUDED.epoch,
	 EXCLUDED.highest_sent_epoch,
	 EXCLUDED.highest_received_epoch,
	 NOW()
	)`, d.oracleSpecID, cd, state.Epoch, state.HighestSentEpoch, pq.Array(&highestReceivedEpoch))

	return errors.Wrap(err, "WriteState failed")
}

func (d *db) ReadConfig(ctx context.Context) (c *ocrtypes.ContractConfig, err error) {
	q := d.QueryRowContext(ctx, `
SELECT
	config_digest,
	config_count,
	signers,
	transmitters,
	f,
	onchain_config,
	offchain_config_version,
	offchain_config
FROM offchainreporting2_contract_configs
WHERE offchainreporting2_oracle_spec_id = $1
LIMIT 1`, d.oracleSpecID)

	c = new(ocrtypes.ContractConfig)

	digest := []byte{}
	signers := [][]byte{}
	transmitters := [][]byte{}

	err = q.Scan(
		&digest,
		&c.ConfigCount,
		(*pq.ByteaArray)(&signers),
		(*pq.ByteaArray)(&transmitters),
		&c.F,
		&c.OnchainConfig,
		&c.OffchainConfigVersion,
		&c.OffchainConfig,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "ReadConfig failed")
	}

	copy(c.ConfigDigest[:], digest)

	c.Signers = []ocrtypes.OnchainPublicKey{}
	for _, s := range signers {
		signer := ocrtypes.OnchainPublicKey{}
		copy(signer, s[:])
		c.Signers = append(c.Signers, signer)
	}

	c.Transmitters = []ocrtypes.Account{}
	for _, t := range transmitters {
		transmitter := ocrtypes.Account(t)
		c.Transmitters = append(c.Transmitters, transmitter)
	}

	return
}

func (d *db) WriteConfig(ctx context.Context, c ocrtypes.ContractConfig) error {
	var signers [][]byte
	for _, s := range c.Signers {
		signers = append(signers, []byte(s))
	}
	_, err := d.ExecContext(ctx, `
INSERT INTO offchainreporting2_contract_configs (
	offchainreporting2_oracle_spec_id,
	config_digest,
	config_count,
	signers,
	transmitters,
	f,
	onchain_config,
	offchain_config_version,
	offchain_config,
	created_at,
	updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
ON CONFLICT (offchainreporting2_oracle_spec_id) DO UPDATE SET
	config_digest = EXCLUDED.config_digest,
	config_count = EXCLUDED.config_count,
	signers = EXCLUDED.signers,
	transmitters = EXCLUDED.transmitters,
	f = EXCLUDED.f,
	onchain_config = EXCLUDED.onchain_config,
	offchain_config_version = EXCLUDED.offchain_config_version,
	offchain_config = EXCLUDED.offchain_config,
	updated_at = NOW()
`,
		d.oracleSpecID,
		c.ConfigDigest,
		c.ConfigCount,
		pq.ByteaArray(signers),
		c.Transmitters,
		c.F,
		c.OnchainConfig,
		c.OffchainConfigVersion,
		c.OffchainConfig,
	)

	return errors.Wrap(err, "WriteConfig failed")
}

func (d *db) StorePendingTransmission(ctx context.Context, t ocrtypes.ReportTimestamp, tx ocrtypes.PendingTransmission) error {
	var signatures [][]byte
	for _, s := range tx.AttributedSignatures {
		signatures = append(signatures, s.Signature)
		buffer := make([]byte, binary.MaxVarintLen64)
		binary.PutVarint(buffer, int64(s.Signer))
		signatures = append(signatures, buffer)
	}

	digest := make([]byte, 32)
	copy(digest, t.ConfigDigest[:])

	extraHash := make([]byte, 32)
	copy(extraHash[:], tx.ExtraHash[:])

	_, err := d.ExecContext(ctx, `
INSERT INTO offchainreporting2_pending_transmissions (
	offchainreporting2_oracle_spec_id,
	config_digest,
	epoch,
	round,

	time,
	extra_hash,
	report,
	attributed_signatures,

	created_at,
	updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
ON CONFLICT (offchainreporting2_oracle_spec_id, config_digest, epoch, round) DO UPDATE SET
	offchainreporting2_oracle_spec_id = EXCLUDED.offchainreporting2_oracle_spec_id,
	config_digest = EXCLUDED.config_digest,
	epoch = EXCLUDED.epoch,
	round = EXCLUDED.round,

	time = EXCLUDED.time,
	extra_hash = EXCLUDED.extra_hash,
	report = EXCLUDED.report,
	attributed_signatures = EXCLUDED.attributed_signatures,

	updated_at = NOW()
`,
		d.oracleSpecID,
		digest,
		t.Epoch,
		t.Round,
		tx.Time,
		extraHash,
		tx.Report,
		pq.ByteaArray(signatures),
	)

	return errors.Wrap(err, "StorePendingTransmission failed")
}

func (d *db) PendingTransmissionsWithConfigDigest(ctx context.Context, cd ocrtypes.ConfigDigest) (map[ocrtypes.ReportTimestamp]ocrtypes.PendingTransmission, error) {
	rows, err := d.QueryContext(ctx, `
SELECT
	config_digest,
	epoch,
	round,
	time,
	extra_hash,
	report,
	attributed_signatures
FROM offchainreporting2_pending_transmissions
WHERE offchainreporting2_oracle_spec_id = $1 AND config_digest = $2
`, d.oracleSpecID, cd)
	if err != nil {
		return nil, errors.Wrap(err, "PendingTransmissionsWithConfigDigest failed to query rows")
	}
	defer d.lggr.ErrorIfClosing(rows, "offchainreporting2_pending_transmissions rows")

	m := make(map[ocrtypes.ReportTimestamp]ocrtypes.PendingTransmission)

	for rows.Next() {
		k := ocrtypes.ReportTimestamp{}
		p := ocrtypes.PendingTransmission{}

		signatures := [][]byte{}
		digest := []byte{}
		extraHash := []byte{}
		report := []byte{}

		if err := rows.Scan(&digest, &k.Epoch, &k.Round, &p.Time, &extraHash, &report, (*pq.ByteaArray)(&signatures)); err != nil {
			return nil, errors.Wrap(err, "PendingTransmissionsWithConfigDigest failed to scan row")
		}

		copy(k.ConfigDigest[:], digest)
		copy(p.ExtraHash[:], extraHash)
		p.Report = make([]byte, len(report))
		copy(p.Report[:], report)

		for index := 0; index < len(signatures); index += 2 {
			signature := signatures[index]
			signer, _ := binary.Varint(signatures[index+1])
			sig := ocrtypes.AttributedOnchainSignature{
				Signature: signature,
				Signer:    ocrcommon.OracleID(signer),
			}
			p.AttributedSignatures = append(p.AttributedSignatures, sig)
		}
		m[k] = p
	}

	if err := rows.Err(); err != nil {
		return m, err
	}

	return m, nil
}

func (d *db) DeletePendingTransmission(ctx context.Context, t ocrtypes.ReportTimestamp) (err error) {
	_, err = d.ExecContext(ctx, `
DELETE FROM offchainreporting2_pending_transmissions
WHERE offchainreporting2_oracle_spec_id = $1 AND  config_digest = $2 AND epoch = $3 AND round = $4
`, d.oracleSpecID, t.ConfigDigest, t.Epoch, t.Round)

	err = errors.Wrap(err, "DeletePendingTransmission failed")

	return
}

func (d *db) DeletePendingTransmissionsOlderThan(ctx context.Context, t time.Time) (err error) {
	_, err = d.ExecContext(ctx, `
DELETE FROM offchainreporting2_pending_transmissions
WHERE offchainreporting2_oracle_spec_id = $1 AND time < $2
`, d.oracleSpecID, t)

	err = errors.Wrap(err, "DeletePendingTransmissionsOlderThan failed")

	return
}
