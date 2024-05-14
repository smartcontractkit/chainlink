package mercury

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	pkgerrors "github.com/pkg/errors"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

type ORM interface {
	InsertTransmitRequest(ctx context.Context, serverURLs []string, req *pb.TransmitRequest, jobID int32, reportCtx ocrtypes.ReportContext) error
	DeleteTransmitRequests(ctx context.Context, serverURL string, reqs []*pb.TransmitRequest) error
	GetTransmitRequests(ctx context.Context, serverURL string, jobID int32) ([]*Transmission, error)
	PruneTransmitRequests(ctx context.Context, serverURL string, jobID int32, maxSize int) error
	LatestReport(ctx context.Context, feedID [32]byte) (report []byte, err error)
}

func FeedIDFromReport(report ocrtypes.Report) (feedID utils.FeedID, err error) {
	if n := copy(feedID[:], report); n != 32 {
		return feedID, pkgerrors.Errorf("invalid length for report: %d", len(report))
	}
	return feedID, nil
}

type orm struct {
	ds sqlutil.DataSource
}

func NewORM(ds sqlutil.DataSource) ORM {
	return &orm{ds: ds}
}

// InsertTransmitRequest inserts one transmit request if the payload does not exist already.
func (o *orm) InsertTransmitRequest(ctx context.Context, serverURLs []string, req *pb.TransmitRequest, jobID int32, reportCtx ocrtypes.ReportContext) error {
	feedID, err := FeedIDFromReport(req.Payload)
	if err != nil {
		return err
	}
	if len(serverURLs) == 0 {
		return errors.New("no server URLs provided")
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var err1, err2 error

	go func() {
		defer wg.Done()

		values := make([]string, len(serverURLs))
		args := []interface{}{
			req.Payload,
			hashPayload(req.Payload),
			reportCtx.ConfigDigest[:],
			reportCtx.Epoch,
			reportCtx.Round,
			reportCtx.ExtraHash[:],
			jobID,
			feedID[:],
		}
		for i, serverURL := range serverURLs {
			// server url is the only thing that changes, might as well re-use
			// the same parameters for each insert
			values[i] = fmt.Sprintf("($1, $2, $3, $4, $5, $6, $7, $8, $%d)", i+9)
			args = append(args, serverURL)
		}

		_, err1 = o.ds.ExecContext(ctx, fmt.Sprintf(`
		INSERT INTO mercury_transmit_requests (payload, payload_hash, config_digest, epoch, round, extra_hash, job_id, feed_id, server_url)
		VALUES %s
		ON CONFLICT (server_url, payload_hash) DO NOTHING
	`, strings.Join(values, ",")), args...)
	}()

	go func() {
		defer wg.Done()
		_, err2 = o.ds.ExecContext(ctx, `
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
func (o *orm) DeleteTransmitRequests(ctx context.Context, serverURL string, reqs []*pb.TransmitRequest) error {
	if len(reqs) == 0 {
		return nil
	}

	var hashes pq.ByteaArray
	for _, req := range reqs {
		hashes = append(hashes, hashPayload(req.Payload))
	}

	_, err := o.ds.ExecContext(ctx, `
		DELETE FROM mercury_transmit_requests
		WHERE server_url = $1 AND payload_hash = ANY($2)
	`, serverURL, hashes)
	return err
}

// GetTransmitRequests returns all transmit requests in chronologically descending order.
func (o *orm) GetTransmitRequests(ctx context.Context, serverURL string, jobID int32) ([]*Transmission, error) {
	// The priority queue uses epoch and round to sort transmissions so order by
	// the same fields here for optimal insertion into the pq.
	rows, err := o.ds.QueryContext(ctx, `
		SELECT payload, config_digest, epoch, round, extra_hash
		FROM mercury_transmit_requests
		WHERE job_id = $1 AND server_url = $2
		ORDER BY epoch DESC, round DESC
	`, jobID, serverURL)
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
func (o *orm) PruneTransmitRequests(ctx context.Context, serverURL string, jobID int32, maxSize int) error {
	// Prune the oldest requests by epoch and round.
	_, err := o.ds.ExecContext(ctx, `
		DELETE FROM mercury_transmit_requests
		WHERE job_id = $1 AND server_url = $2 AND
		payload_hash NOT IN (
		    SELECT payload_hash
			FROM mercury_transmit_requests
			WHERE job_id = $1 AND server_url = $2
			ORDER BY epoch DESC, round DESC
			LIMIT $3
		)
	`, jobID, serverURL, maxSize)
	return err
}

func (o *orm) LatestReport(ctx context.Context, feedID [32]byte) (report []byte, err error) {
	err = o.ds.GetContext(ctx, &report, `SELECT report FROM feed_latest_reports WHERE feed_id = $1`, feedID[:])
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return report, err
}

func hashPayload(payload []byte) []byte {
	checksum := sha256.Sum256(payload)
	return checksum[:]
}
