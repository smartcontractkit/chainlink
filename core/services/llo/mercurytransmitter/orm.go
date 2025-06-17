package mercurytransmitter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/lib/pq"

	"github.com/smartcontractkit/libocr/commontypes"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

// ORM is scoped to a single DON ID
type ORM interface {
	DonID() uint32
	Insert(ctx context.Context, transmissions []*Transmission) error
	Delete(ctx context.Context, hashes [][32]byte) error
	Get(ctx context.Context, serverURL string, limit int, maxAge time.Duration) ([]*Transmission, error)
	Prune(ctx context.Context, serverURL string, maxSize, batchSize int) (int64, error)
}

type orm struct {
	ds    sqlutil.DataSource
	donID uint32
}

func NewORM(ds sqlutil.DataSource, donID uint32) ORM {
	return &orm{ds, donID}
}

func (o *orm) DonID() uint32 {
	return o.donID
}

// Insert inserts the transmissions, ignoring duplicates
func (o *orm) Insert(ctx context.Context, transmissions []*Transmission) error {
	if len(transmissions) == 0 {
		return nil
	}

	type transmission struct {
		DonID            uint32                `db:"don_id"`
		ServerURL        string                `db:"server_url"`
		ConfigDigest     ocrtypes.ConfigDigest `db:"config_digest"`
		SeqNr            int64                 `db:"seq_nr"`
		Report           []byte                `db:"report"`
		LifecycleStage   string                `db:"lifecycle_stage"`
		ReportFormat     uint32                `db:"report_format"`
		Signatures       [][]byte              `db:"signatures"`
		Signers          []uint8               `db:"signers"`
		TransmissionHash []byte                `db:"transmission_hash"`
	}
	records := make([]transmission, len(transmissions))
	for i, t := range transmissions {
		signatures := make([][]byte, len(t.Sigs))
		signers := make([]uint8, len(t.Sigs))
		for j, sig := range t.Sigs {
			signatures[j] = sig.Signature
			signers[j] = uint8(sig.Signer)
		}
		h := t.Hash()
		if t.SeqNr > math.MaxInt64 {
			// this is to appease the linter but shouldn't ever happen
			return fmt.Errorf("seqNr is too large (got: %d, max: %d)", t.SeqNr, math.MaxInt64)
		}
		records[i] = transmission{
			DonID:            o.donID,
			ServerURL:        t.ServerURL,
			ConfigDigest:     t.ConfigDigest,
			SeqNr:            int64(t.SeqNr),
			Report:           t.Report.Report,
			LifecycleStage:   string(t.Report.Info.LifeCycleStage),
			ReportFormat:     uint32(t.Report.Info.ReportFormat),
			Signatures:       signatures,
			Signers:          signers,
			TransmissionHash: h[:],
		}
	}

	_, err := o.ds.NamedExecContext(ctx, `
	INSERT INTO llo_mercury_transmit_queue (don_id, server_url, config_digest, seq_nr, report, lifecycle_stage, report_format, signatures, signers, transmission_hash)
		VALUES (:don_id, :server_url, :config_digest, :seq_nr, :report, :lifecycle_stage, :report_format, :signatures, :signers, :transmission_hash)
		ON CONFLICT (transmission_hash) DO NOTHING
	`, records)

	if err != nil {
		return fmt.Errorf("llo orm: failed to insert transmissions: %w", err)
	}
	return nil
}

// Delete deletes the given transmissions
func (o *orm) Delete(ctx context.Context, hashes [][32]byte) error {
	if len(hashes) == 0 {
		return nil
	}

	var pqHashes pq.ByteaArray
	for _, hash := range hashes {
		pqHashes = append(pqHashes, hash[:])
	}

	_, err := o.ds.ExecContext(ctx, `
		DELETE FROM llo_mercury_transmit_queue
		WHERE transmission_hash = ANY($1)
	`, pqHashes)
	if err != nil {
		return fmt.Errorf("llo orm: failed to delete transmissions: %w", err)
	}
	return nil
}

// Get returns all transmissions in chronologically descending order
// NOTE: passing maxAge=0 disables any age filter
func (o *orm) Get(ctx context.Context, serverURL string, limit int, maxAge time.Duration) ([]*Transmission, error) {
	// The priority queue uses seqnr to sort transmissions so order by
	// the same fields here for optimal insertion into the pq.
	maxAgeClause := ""
	params := []interface{}{o.donID, serverURL, limit}
	if maxAge > 0 {
		maxAgeClause = "\nAND inserted_at >= NOW() - ($4 * INTERVAL '1 MICROSECOND')"
		params = append(params, maxAge.Microseconds())
	}
	q := fmt.Sprintf(`
		SELECT config_digest, seq_nr, report, lifecycle_stage, report_format, signatures, signers
		FROM llo_mercury_transmit_queue
		WHERE don_id = $1 AND server_url = $2%s
		ORDER BY seq_nr DESC, inserted_at DESC
		LIMIT $3
		`, maxAgeClause)
	rows, err := o.ds.QueryContext(ctx, q, params...)
	if err != nil {
		return nil, fmt.Errorf("llo orm: failed to get transmissions: %w", err)
	}
	defer rows.Close()

	var transmissions []*Transmission
	for rows.Next() {
		transmission := Transmission{
			ServerURL: serverURL,
		}
		var digest []byte
		var signatures pq.ByteaArray
		var signers pq.Int32Array

		err := rows.Scan(
			&digest,
			&transmission.SeqNr,
			&transmission.Report.Report,
			&transmission.Report.Info.LifeCycleStage,
			&transmission.Report.Info.ReportFormat,
			&signatures,
			&signers,
		)
		if err != nil {
			return nil, fmt.Errorf("llo orm: failed to scan transmission: %w", err)
		}
		transmission.ConfigDigest = ocrtypes.ConfigDigest(digest)
		if len(signatures) != len(signers) {
			return nil, errors.New("signatures and signers must have the same length")
		}
		for i, sig := range signatures {
			if signers[i] > math.MaxUint8 {
				// this is to appease the linter but shouldn't ever happen
				return nil, fmt.Errorf("signer is too large (got: %d, max: %d)", signers[i], math.MaxUint8)
			}
			transmission.Sigs = append(transmission.Sigs, ocrtypes.AttributedOnchainSignature{
				Signature: sig,
				Signer:    commontypes.OracleID(signers[i]), //nolint:gosec // G115 false positive
			})
		}

		transmissions = append(transmissions, &transmission)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("llo orm: failed to scan transmissions: %w", err)
	}

	return transmissions, nil
}

// Prune keeps at most maxSize rows for the given (donID, serverURL) pair by
// deleting the oldest transmissions.
func (o *orm) Prune(ctx context.Context, serverURL string, maxSize, batchSize int) (rowsDeleted int64, err error) {
	var oldest uint64
	err = o.ds.GetContext(ctx, &oldest, `SELECT seq_nr
		FROM llo_mercury_transmit_queue
		WHERE don_id = $1 AND server_url = $2
		ORDER BY seq_nr DESC
		OFFSET $3
		LIMIT 1`, o.donID, serverURL, maxSize)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("llo orm: failed to get oldest seq_nr: %w", err)
	}
	// Prune the requests with seq_nr older than this, in batches to avoid long
	// locking or queries
	for {
		var res sql.Result
		res, err = o.ds.ExecContext(ctx, `
DELETE FROM llo_mercury_transmit_queue AS q
USING (
    SELECT transmission_hash
    FROM llo_mercury_transmit_queue
    WHERE don_id = $1
      AND server_url = $2
      AND seq_nr < $3
    ORDER BY seq_nr ASC
    LIMIT $4
) AS to_delete
WHERE q.transmission_hash = to_delete.transmission_hash;
		`, o.donID, serverURL, oldest, batchSize)
		if err != nil {
			return rowsDeleted, fmt.Errorf("llo orm: batch delete failed to prune transmissions: %w", err)
		}
		var rowsAffected int64
		rowsAffected, err = res.RowsAffected()
		if err != nil {
			return rowsDeleted, fmt.Errorf("llo orm: batch delete failed to get rows affected: %w", err)
		}
		if rowsAffected == 0 {
			break
		}
		rowsDeleted += rowsAffected
	}

	// This query to trim off the final few rows to reach exactly maxSize with
	// should now be fast and efficient because of the batch deletes that
	// already completed above.
	res, err := o.ds.ExecContext(ctx, `
WITH to_delete AS (
    SELECT ctid
    FROM (
        SELECT ctid,
               ROW_NUMBER() OVER (PARTITION BY don_id, server_url ORDER BY seq_nr DESC, inserted_at DESC) AS row_num
        FROM llo_mercury_transmit_queue
		WHERE don_id = $1 AND server_url = $2
    ) sub
    WHERE row_num > $3
)
DELETE FROM llo_mercury_transmit_queue
WHERE ctid IN (SELECT ctid FROM to_delete);
`, o.donID, serverURL, maxSize)

	if err != nil {
		return rowsDeleted, fmt.Errorf("llo orm: final truncate failed to prune transmissions: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return rowsDeleted, fmt.Errorf("llo orm: final truncate failed to get rows affected: %w", err)
	}
	rowsDeleted += rowsAffected

	return rowsDeleted, nil
}
