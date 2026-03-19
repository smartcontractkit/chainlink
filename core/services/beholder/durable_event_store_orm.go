package beholder

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

const chipDurableEventsTable = "cre.chip_durable_events"

// PgDurableEventStore is a Postgres-backed implementation of beholder.DurableEventStore.
type PgDurableEventStore struct {
	ds sqlutil.DataSource
}

var _ beholder.DurableEventStore = (*PgDurableEventStore)(nil)

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

func (s *PgDurableEventStore) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM ` + chipDurableEventsTable + ` WHERE id = $1`
	if _, err := s.ds.ExecContext(ctx, q, id); err != nil {
		return fmt.Errorf("failed to delete chip durable event id=%d: %w", id, err)
	}
	return nil
}

func (s *PgDurableEventStore) ListPending(ctx context.Context, createdBefore time.Time, limit int) ([]beholder.DurableEvent, error) {
	const q = `
SELECT id, payload, created_at
FROM ` + chipDurableEventsTable + `
WHERE created_at < $1
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
    WHERE created_at < now() - $1::interval
    RETURNING id
)
SELECT count(*) FROM deleted`

	var count int64
	if err := s.ds.GetContext(ctx, &count, q, ttl.String()); err != nil {
		return 0, fmt.Errorf("failed to delete expired chip durable events: %w", err)
	}
	return count, nil
}
