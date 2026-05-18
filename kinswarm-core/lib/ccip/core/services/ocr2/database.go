package ocr2

import (
	"context"
	"database/sql"
	"encoding/binary"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	ocrcommon "github.com/smartcontractkit/libocr/commontypes"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type db struct {
	ds           sqlutil.DataSource
	oracleSpecID int32
	pluginID     int32
	lggr         logger.SugaredLogger
}

var (
	_ ocrtypes.Database = &db{}
)

// NewDB returns a new DB scoped to this oracleSpecID
func NewDB(ds sqlutil.DataSource, oracleSpecID int32, pluginID int32, lggr logger.Logger) *db {
	return &db{
		ds:           ds,
		oracleSpecID: oracleSpecID,
		pluginID:     pluginID,
		lggr:         logger.Sugared(lggr),
	}
}

func (d *db) ReadState(ctx context.Context, cd ocrtypes.ConfigDigest) (ps *ocrtypes.PersistentState, err error) {
	stmt := `
	SELECT epoch, highest_sent_epoch, highest_received_epoch
	FROM ocr2_persistent_states
	WHERE ocr2_oracle_spec_id = $1 AND config_digest = $2
	LIMIT 1`

	ps = new(ocrtypes.PersistentState)

	var tmp []int64
	var highestSentEpochTmp int64

	err = d.ds.QueryRowxContext(ctx, stmt, d.oracleSpecID, cd).Scan(&ps.Epoch, &highestSentEpochTmp, pq.Array(&tmp))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "ReadState failed")
	}

	ps.HighestSentEpoch = uint32(highestSentEpochTmp)

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

	stmt := `
	INSERT INTO ocr2_persistent_states (
		ocr2_oracle_spec_id,
		config_digest,
		epoch,
		highest_sent_epoch,
		highest_received_epoch,
		created_at,
		updated_at
	)
	VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	ON CONFLICT (ocr2_oracle_spec_id, config_digest)
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
		)`

	ctx, cancel := context.WithTimeout(sqlutil.WithoutDefaultTimeout(ctx), time.Minute)
	defer cancel()
	_, err := d.ds.ExecContext(
		ctx, stmt, d.oracleSpecID, cd, state.Epoch, state.HighestSentEpoch, pq.Array(&highestReceivedEpoch),
	)

	return errors.Wrap(err, "WriteState failed")
}

func (d *db) ReadConfig(ctx context.Context) (c *ocrtypes.ContractConfig, err error) {
	stmt := `
	SELECT
		config_digest,
		config_count,
		signers,
		transmitters,
		f,
		onchain_config,
		offchain_config_version,
		offchain_config
	FROM ocr2_contract_configs
	WHERE ocr2_oracle_spec_id = $1 AND plugin_id = $2
	LIMIT 1`

	c = new(ocrtypes.ContractConfig)

	digest := []byte{}
	signers := [][]byte{}
	transmitters := [][]byte{}

	err = d.ds.QueryRowxContext(ctx, stmt, d.oracleSpecID, d.pluginID).Scan(
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
	}
	if err != nil {
		return nil, errors.Wrap(err, "ReadConfig failed")
	}

	copy(c.ConfigDigest[:], digest)

	c.Signers = []ocrtypes.OnchainPublicKey{}
	for _, s := range signers {
		c.Signers = append(c.Signers, s)
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
	stmt := `
	INSERT INTO ocr2_contract_configs (
		ocr2_oracle_spec_id,
		plugin_id,
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
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
	ON CONFLICT (ocr2_oracle_spec_id, plugin_id) DO UPDATE SET
		config_digest = EXCLUDED.config_digest,
		config_count = EXCLUDED.config_count,
		signers = EXCLUDED.signers,
		transmitters = EXCLUDED.transmitters,
		f = EXCLUDED.f,
		onchain_config = EXCLUDED.onchain_config,
		offchain_config_version = EXCLUDED.offchain_config_version,
		offchain_config = EXCLUDED.offchain_config,
		updated_at = NOW()
	`
	_, err := d.ds.ExecContext(ctx, stmt,
		d.oracleSpecID,
		d.pluginID,
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

	stmt := `
	INSERT INTO ocr2_pending_transmissions (
		ocr2_oracle_spec_id,
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
	ON CONFLICT (ocr2_oracle_spec_id, config_digest, epoch, round) DO UPDATE SET
		ocr2_oracle_spec_id = EXCLUDED.ocr2_oracle_spec_id,
		config_digest = EXCLUDED.config_digest,
		epoch = EXCLUDED.epoch,
		round = EXCLUDED.round,
	
		time = EXCLUDED.time,
		extra_hash = EXCLUDED.extra_hash,
		report = EXCLUDED.report,
		attributed_signatures = EXCLUDED.attributed_signatures,
	
		updated_at = NOW()
	`

	_, err := d.ds.ExecContext(ctx, stmt,
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
	stmt := `
	SELECT
		config_digest,
		epoch,
		round,
		time,
		extra_hash,
		report,
		attributed_signatures
	FROM ocr2_pending_transmissions
	WHERE ocr2_oracle_spec_id = $1 AND config_digest = $2
	`
	rows, err := d.ds.QueryxContext(ctx, stmt, d.oracleSpecID, cd) //nolint sqlclosecheck false positive
	if err != nil {
		return nil, errors.Wrap(err, "PendingTransmissionsWithConfigDigest failed to query rows")
	}
	defer d.lggr.ErrorIfFn(rows.Close, "Error closing ocr2_pending_transmissions rows")

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
	ctx, cancel := context.WithTimeout(sqlutil.WithoutDefaultTimeout(ctx), time.Minute)
	defer cancel()
	_, err = d.ds.ExecContext(ctx, `
DELETE FROM ocr2_pending_transmissions
WHERE ocr2_oracle_spec_id = $1 AND  config_digest = $2 AND epoch = $3 AND round = $4
`, d.oracleSpecID, t.ConfigDigest, t.Epoch, t.Round)

	err = errors.Wrap(err, "DeletePendingTransmission failed")

	return
}

func (d *db) DeletePendingTransmissionsOlderThan(ctx context.Context, t time.Time) (err error) {
	ctx, cancel := context.WithTimeout(sqlutil.WithoutDefaultTimeout(ctx), time.Minute)
	defer cancel()
	_, err = d.ds.ExecContext(ctx, `
DELETE FROM ocr2_pending_transmissions
WHERE ocr2_oracle_spec_id = $1 AND time < $2
`, d.oracleSpecID, t)

	err = errors.Wrap(err, "DeletePendingTransmissionsOlderThan failed")

	return
}

func (d *db) ReadProtocolState(ctx context.Context, configDigest ocrtypes.ConfigDigest, key string) (value []byte, err error) {
	err = d.ds.GetContext(ctx, &value, `
SELECT value FROM ocr_protocol_states
WHERE config_digest = $1 AND key = $2;
`, configDigest, key)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	err = errors.Wrapf(err, "ReadProtocolState failed for job %d", d.oracleSpecID)

	return
}

func (d *db) WriteProtocolState(ctx context.Context, configDigest ocrtypes.ConfigDigest, key string, value []byte) (err error) {
	if value == nil {
		_, err = d.ds.ExecContext(ctx, `DELETE FROM ocr_protocol_states WHERE config_digest = $1 AND key = $2;`, configDigest, key)
	} else {
		_, err = d.ds.ExecContext(ctx, `
INSERT INTO ocr_protocol_states (config_digest, key, value) VALUES ($1, $2, $3)
ON CONFLICT (config_digest, key) DO UPDATE SET value = $3;`, configDigest, key, value)
	}

	err = errors.Wrapf(err, "WriteProtocolState failed for job %d", d.oracleSpecID)

	return
}
