package beholder

/* TODO: CRE-4422 Refactor: relocate this to durableemitter pkg
import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

const chipDurableEventsTable = "cre.chip_durable_events"

// PgDurableEventStore is a Postgres-backed implementation of beholder.DurableEventStore.
type PgDurableEventStore struct {
	ds sqlutil.DataSource
}

var (
	_ beholder.DurableEventStore    = (*PgDurableEventStore)(nil)
	_ beholder.DurableQueueObserver = (*PgDurableEventStore)(nil)
	_ beholder.BatchInserter        = (*PgDurableEventStore)(nil)
)

func NewPgDurableEventStore(ds sqlutil.DataSource) *PgDurableEventStore {
	return &PgDurableEventStore{ds: ds}
}

func (s *PgDurableEventStore) Insert(ctx context.Context, payload []byte) (int64, error) {
	const q = `INSERT INTO ` + chipDurableEventsTable + ` (payload) VALUES ($1) RETURNING id`
	var id int64
	if err := s.ds.GetContext(ctx, &id, q, payload); err != nil {
		return 0, fmt.Errorf("failed to insert chip durable event: %w", err)
	}
	return id, nil
}

func (s *PgDurableEventStore) InsertBatch(ctx context.Context, payloads [][]byte) ([]int64, error) {
	if len(payloads) == 0 {
		return nil, nil
	}
	placeholders := make([]string, len(payloads))
	args := make([]interface{}, len(payloads))
	for i, p := range payloads {
		placeholders[i] = fmt.Sprintf("($%d)", i+1)
		args[i] = p
	}
	q := fmt.Sprintf(
		"INSERT INTO %s (payload) VALUES %s RETURNING id",
		chipDurableEventsTable,
		strings.Join(placeholders, ","),
	)

	var ids []int64
	if err := s.ds.SelectContext(ctx, &ids, q, args...); err != nil {
		return nil, fmt.Errorf("failed to batch insert chip durable events: %w", err)
	}
	return ids, nil
}

func (s *PgDurableEventStore) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM ` + chipDurableEventsTable + ` WHERE id = $1`
	if _, err := s.ds.ExecContext(ctx, q, id); err != nil {
		return fmt.Errorf("failed to delete chip durable event id=%d: %w", id, err)
	}
	return nil
}

func (s *PgDurableEventStore) MarkDelivered(ctx context.Context, id int64) error {
	const q = `UPDATE ` + chipDurableEventsTable + ` SET delivered_at = now() WHERE id = $1 AND delivered_at IS NULL`
	if _, err := s.ds.ExecContext(ctx, q, id); err != nil {
		return fmt.Errorf("failed to mark chip durable event delivered id=%d: %w", id, err)
	}
	return nil
}

func (s *PgDurableEventStore) MarkDeliveredBatch(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	const q = `UPDATE ` + chipDurableEventsTable + ` SET delivered_at = now() WHERE id = ANY($1) AND delivered_at IS NULL`
	res, err := s.ds.ExecContext(ctx, q, pq.Array(ids))
	if err != nil {
		return 0, fmt.Errorf("failed to batch mark chip durable events delivered: %w", err)
	}
	n, _ := res.RowsAffected()
	return n, nil
}

func (s *PgDurableEventStore) PurgeDelivered(ctx context.Context, batchLimit int) (int64, error) {
	if batchLimit <= 0 {
		return 0, nil
	}
	const q = `
WITH picked AS (
    SELECT id FROM ` + chipDurableEventsTable + `
    WHERE delivered_at IS NOT NULL
    ORDER BY delivered_at ASC
    LIMIT $1
)
DELETE FROM ` + chipDurableEventsTable + ` AS t
USING picked WHERE t.id = picked.id`
	res, err := s.ds.ExecContext(ctx, q, batchLimit)
	if err != nil {
		return 0, fmt.Errorf("failed to purge delivered chip durable events: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("purge delivered rows affected: %w", err)
	}
	return n, nil
}

func (s *PgDurableEventStore) ListPending(ctx context.Context, createdBefore time.Time, limit int) ([]beholder.DurableEvent, error) {
	const q = `
SELECT id, payload, created_at
FROM ` + chipDurableEventsTable + `
WHERE delivered_at IS NULL
  AND created_at < $1
ORDER BY created_at ASC
LIMIT $2`

	type row struct {
		ID        int64     `db:"id"`
		Payload   []byte    `db:"payload"`
		CreatedAt time.Time `db:"created_at"`
	}

	var rows []row
	if err := s.ds.SelectContext(ctx, &rows, q, createdBefore, limit); err != nil {
		return nil, fmt.Errorf("failed to list pending chip durable events: %w", err)
	}

	out := make([]beholder.DurableEvent, 0, len(rows))
	for _, r := range rows {
		out = append(out, beholder.DurableEvent{
			ID:        r.ID,
			Payload:   r.Payload,
			CreatedAt: r.CreatedAt,
		})
	}
	return out, nil
}

func (s *PgDurableEventStore) DeleteExpired(ctx context.Context, ttl time.Duration) (int64, error) {
	const q = `
WITH deleted AS (
    DELETE FROM ` + chipDurableEventsTable + `
    WHERE created_at <= now() - $1::interval
    RETURNING id
)
SELECT count(*) FROM deleted`

	var count int64
	if err := s.ds.GetContext(ctx, &count, q, ttl.String()); err != nil {
		return 0, fmt.Errorf("failed to delete expired chip durable events: %w", err)
	}
	return count, nil
}

type chipDurableQueueAgg struct {
	Cnt        int64      `db:"cnt"`
	PayloadSum int64      `db:"payload_sum"`
	MinCreated *time.Time `db:"min_created"`
}

// ObserveDurableQueue implements beholder.DurableQueueObserver for queue depth / age gauges.
func (s *PgDurableEventStore) ObserveDurableQueue(ctx context.Context, eventTTL, nearExpiryLead time.Duration) (beholder.DurableQueueStats, error) {
	const qAgg = `
SELECT
	count(*)::bigint AS cnt,
	coalesce(sum(octet_length(payload)), 0)::bigint AS payload_sum,
	min(created_at) AS min_created
FROM ` + chipDurableEventsTable + `
WHERE delivered_at IS NULL`

	var row chipDurableQueueAgg
	if err := s.ds.GetContext(ctx, &row, qAgg); err != nil {
		return beholder.DurableQueueStats{}, fmt.Errorf("durable queue aggregate: %w", err)
	}
	var st beholder.DurableQueueStats
	st.Depth = row.Cnt
	st.PayloadBytes = row.PayloadSum
	if row.MinCreated != nil {
		st.OldestPendingAge = time.Since(*row.MinCreated)
	}
	if eventTTL > 0 && nearExpiryLead > 0 && nearExpiryLead < eventTTL {
		ttlSec := int64(eventTTL.Round(time.Second) / time.Second)
		leadSec := int64(nearExpiryLead.Round(time.Second) / time.Second)
		const qNear = `
SELECT count(*)::bigint
FROM ` + chipDurableEventsTable + `
WHERE delivered_at IS NULL
  AND created_at >= now() - ($1::bigint * interval '1 second')
  AND created_at < now() - (($1::bigint - $2::bigint) * interval '1 second')`
		if err := s.ds.GetContext(ctx, &st.NearTTLCount, qNear, ttlSec, leadSec); err != nil {
			return beholder.DurableQueueStats{}, fmt.Errorf("durable queue near-ttl: %w", err)
		}
	}
	return st, nil
}
*/
