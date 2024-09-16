package mercurytransmitter

import (
	"context"
	"errors"
	"fmt"
	"math"

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
	Get(ctx context.Context, serverURL string) ([]*Transmission, error)
	Prune(ctx context.Context, serverURL string, maxSize int) error
	Cleanup(ctx context.Context) error
}

type orm struct {
	ds    sqlutil.DataSource
	donID uint32
}

func NewORM(ds sqlutil.DataSource, donID uint32) ORM {
	return &orm{ds: ds, donID: donID}
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
			SeqNr:            int64(t.SeqNr), //nolint
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
func (o *orm) Get(ctx context.Context, serverURL string) ([]*Transmission, error) {
	// The priority queue uses seqnr to sort transmissions so order by
	// the same fields here for optimal insertion into the pq.
	rows, err := o.ds.QueryContext(ctx, `
		SELECT config_digest, seq_nr, report, lifecycle_stage, report_format, signatures, signers
		FROM llo_mercury_transmit_queue
		WHERE don_id = $1 AND server_url = $2
		ORDER BY seq_nr DESC, transmission_hash DESC
	`, o.donID, serverURL)
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
				Signer:    commontypes.OracleID(signers[i]), //nolint
			})
		}

		transmissions = append(transmissions, &transmission)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("llo orm: failed to scan transmissions: %w", err)
	}

	return transmissions, nil
}

// Prune keeps at most maxSize rows for the given job ID,
// deleting the oldest transactions.
func (o *orm) Prune(ctx context.Context, serverURL string, maxSize int) error {
	// Prune the oldest requests by epoch and round.
	_, err := o.ds.ExecContext(ctx, `
		DELETE FROM llo_mercury_transmit_queue
		WHERE don_id = $1 AND server_url = $2 AND
		transmission_hash NOT IN (
		    SELECT transmission_hash
			FROM llo_mercury_transmit_queue
			WHERE don_id = $1 AND server_url = $2
			ORDER BY seq_nr DESC, transmission_hash DESC
			LIMIT $3
		)
	`, o.donID, serverURL, maxSize)
	if err != nil {
		return fmt.Errorf("llo orm: failed to prune transmissions: %w", err)
	}
	return nil
}

func (o *orm) Cleanup(ctx context.Context) error {
	_, err := o.ds.ExecContext(ctx, `DELETE FROM llo_mercury_transmit_queue WHERE don_id = $1`, o.donID)
	if err != nil {
		return fmt.Errorf("llo orm: failed to cleanup transmissions: %w", err)
	}
	return nil
}
