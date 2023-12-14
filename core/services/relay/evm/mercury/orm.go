package mercury

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	pkgerrors "github.com/pkg/errors"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

type ORM interface {
	InsertTransmitRequest(req *pb.TransmitRequest, jobID int32, reportCtx ocrtypes.ReportContext, qopts ...pg.QOpt) error
	DeleteTransmitRequests(reqs []*pb.TransmitRequest, qopts ...pg.QOpt) error
	GetTransmitRequests(jobID int32, qopts ...pg.QOpt) ([]*Transmission, error)
	PruneTransmitRequests(jobID int32, maxSize int, qopts ...pg.QOpt) error
	LatestReport(ctx context.Context, feedID [32]byte, qopts ...pg.QOpt) (report []byte, err error)
}

func FeedIDFromReport(report ocrtypes.Report) (feedID utils.FeedID, err error) {
	if n := copy(feedID[:], report); n != 32 {
		return feedID, pkgerrors.Errorf("invalid length for report: %d", len(report))
	}
	return feedID, nil
}

type orm struct {
	q pg.Q
}

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) ORM {
	namedLogger := lggr.Named("MercuryORM")
	q := pg.NewQ(db, namedLogger, cfg)
	return &orm{
		q: q,
	}
}

// InsertTransmitRequest inserts one transmit request if the payload does not exist already.
func (o *orm) InsertTransmitRequest(req *pb.TransmitRequest, jobID int32, reportCtx ocrtypes.ReportContext, qopts ...pg.QOpt) error {
	feedID, err := FeedIDFromReport(req.Payload)
	if err != nil {
		return err
	}

	q := o.q.WithOpts(qopts...)
	var wg sync.WaitGroup
	wg.Add(2)
	var err1, err2 error

	go func() {
		defer wg.Done()
		err1 = q.ExecQ(`
		INSERT INTO mercury_transmit_requests (payload, payload_hash, config_digest, epoch, round, extra_hash, job_id, feed_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (payload_hash) DO NOTHING
	`, req.Payload, hashPayload(req.Payload), reportCtx.ConfigDigest[:], reportCtx.Epoch, reportCtx.Round, reportCtx.ExtraHash[:], jobID, feedID[:])
	}()

	go func() {
		defer wg.Done()
		err2 = q.ExecQ(`
		INSERT INTO feed_latest_reports (feed_id, report, epoch, round, updated_at, job_id)
		VALUES ($1, $2, $3, $4, NOW(), $5)
		ON CONFLICT (feed_id) DO UPDATE
		SET feed_id=$1, report=$2, epoch=$3, round=$4, updated_at=NOW()
		WHERE excluded.epoch > feed_latest_reports.epoch OR (excluded.epoch = feed_latest_reports.epoch AND excluded.round > feed_latest_reports.round)
		`, feedID[:], req.Payload, reportCtx.Epoch, reportCtx.Round, jobID)
	}()
	wg.Wait()
	return errors.Join(err1, err2)
}

// DeleteTransmitRequest deletes the given transmit requests if they exist.
func (o *orm) DeleteTransmitRequests(reqs []*pb.TransmitRequest, qopts ...pg.QOpt) error {
	if len(reqs) == 0 {
		return nil
	}

	var hashes pq.ByteaArray
	for _, req := range reqs {
		hashes = append(hashes, hashPayload(req.Payload))
	}

	q := o.q.WithOpts(qopts...)
	err := q.ExecQ(`
		DELETE FROM mercury_transmit_requests
		WHERE payload_hash = ANY($1)
	`, hashes)
	return err
}

// GetTransmitRequests returns all transmit requests in chronologically descending order.
func (o *orm) GetTransmitRequests(jobID int32, qopts ...pg.QOpt) ([]*Transmission, error) {
	q := o.q.WithOpts(qopts...)
	// The priority queue uses epoch and round to sort transmissions so order by
	// the same fields here for optimal insertion into the pq.
	rows, err := q.QueryContext(q.ParentCtx, `
		SELECT payload, config_digest, epoch, round, extra_hash
		FROM mercury_transmit_requests
		WHERE job_id = $1
		ORDER BY epoch DESC, round DESC
	`, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transmissions []*Transmission
	for rows.Next() {
		transmission := &Transmission{Req: &pb.TransmitRequest{}}
		var digest, extraHash common.Hash

		err := rows.Scan(
			&transmission.Req.Payload,
			&digest,
			&transmission.ReportCtx.Epoch,
			&transmission.ReportCtx.Round,
			&extraHash,
		)
		if err != nil {
			return nil, err
		}
		transmission.ReportCtx.ConfigDigest = ocrtypes.ConfigDigest(digest)
		transmission.ReportCtx.ExtraHash = extraHash

		transmissions = append(transmissions, transmission)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transmissions, nil
}

// PruneTransmitRequests keeps at most maxSize rows for the given job ID,
// deleting the oldest transactions.
func (o *orm) PruneTransmitRequests(jobID int32, maxSize int, qopts ...pg.QOpt) error {
	q := o.q.WithOpts(qopts...)
	// Prune the oldest requests by epoch and round.
	return q.ExecQ(`
		DELETE FROM mercury_transmit_requests
		WHERE job_id = $1 AND
		payload_hash NOT IN (
		    SELECT payload_hash
			FROM mercury_transmit_requests
			WHERE job_id = $1
			ORDER BY epoch DESC, round DESC
			LIMIT $2
		)
	`, jobID, maxSize)
}

func (o *orm) LatestReport(ctx context.Context, feedID [32]byte, qopts ...pg.QOpt) (report []byte, err error) {
	q := o.q.WithOpts(qopts...)
	err = q.GetContext(ctx, &report, `SELECT report FROM feed_latest_reports WHERE feed_id = $1`, feedID[:])
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return report, err
}

func hashPayload(payload []byte) []byte {
	checksum := sha256.Sum256(payload)
	return checksum[:]
}
