package trigger

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

// triggerEventStore is a Postgres-backed implementation of capabilities.EventStore.
type triggerEventStore struct {
	ds sqlutil.DataSource
}

var _ capabilities.EventStore = (*triggerEventStore)(nil)

func NewTriggerEventStore(ds sqlutil.DataSource) *triggerEventStore {
	return &triggerEventStore{ds: ds}
}

const triggerPendingEventsTable = "cre.trigger_pending_events"

func (s *triggerEventStore) Insert(ctx context.Context, rec capabilities.PendingEvent) error {
	const q = `
INSERT INTO ` + triggerPendingEventsTable + ` (
  trigger_id, event_id, payload, first_at, last_sent_at, attempts
) VALUES ($1, $2, $3, $4, $5, $6)`
	var lastSent sql.NullTime
	if !rec.LastSentAt.IsZero() {
		lastSent = sql.NullTime{Time: rec.LastSentAt, Valid: true}
	}

	if _, err := s.ds.ExecContext(
		ctx, q,
		rec.TriggerId,
		rec.EventId,
		rec.Payload,
		rec.FirstAt,
		lastSent,
		rec.Attempts,
	); err != nil {
		return fmt.Errorf("failed to insert pending event trigger_id=%s event_id=%s: %w", rec.TriggerId, rec.EventId, err)
	}
	return nil
}

func (s *triggerEventStore) UpdateDelivery(ctx context.Context, triggerID string, eventID string, lastSentAt time.Time, attempts int) error {
	const q = `
UPDATE ` + triggerPendingEventsTable + `
SET last_sent_at = $3, attempts = $4
WHERE trigger_id = $1 AND event_id = $2
`
	var lastSent interface{}
	if !lastSentAt.IsZero() {
		lastSent = lastSentAt
	}
	res, err := s.ds.ExecContext(ctx, q, triggerID, eventID, lastSent, attempts)
	if err != nil {
		return fmt.Errorf("failed to update delivery metadata for trigger_id=%s event_id=%s: %w", triggerID, eventID, err)
	}

	// verify an event was actually updated
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected while updating delivery for trigger_id=%s event_id=%s: %w", triggerID, eventID, err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *triggerEventStore) List(ctx context.Context) ([]capabilities.PendingEvent, error) {
	const q = `
SELECT
  trigger_id,
  event_id,
  payload,
  first_at,
  last_sent_at,
  attempts
FROM ` + triggerPendingEventsTable + `
ORDER BY first_at ASC
`

	type row struct {
		TriggerID  string       `db:"trigger_id"`
		EventID    string       `db:"event_id"`
		Payload    []byte       `db:"payload"`
		FirstAt    time.Time    `db:"first_at"`
		LastSentAt sql.NullTime `db:"last_sent_at"`
		Attempts   int          `db:"attempts"`
	}

	var rows []row
	if err := s.ds.SelectContext(ctx, &rows, q); err != nil {
		return nil, fmt.Errorf("failed to list pending events: %w", err)
	}

	out := make([]capabilities.PendingEvent, 0, len(rows))
	for _, r := range rows {
		var last time.Time
		if r.LastSentAt.Valid {
			last = r.LastSentAt.Time
		}
		out = append(out, capabilities.PendingEvent{
			TriggerId:  r.TriggerID,
			EventId:    r.EventID,
			Payload:    append([]byte(nil), r.Payload...),
			FirstAt:    r.FirstAt,
			LastSentAt: last,
			Attempts:   r.Attempts,
		})
	}
	return out, nil
}

func (s *triggerEventStore) DeleteEvent(ctx context.Context, triggerID, eventID string) error {
	const q = `
DELETE FROM ` + triggerPendingEventsTable + `
WHERE trigger_id = $1 AND event_id = $2
`
	if _, err := s.ds.ExecContext(ctx, q, triggerID, eventID); err != nil {
		return fmt.Errorf("failed to delete pending event trigger_id=%s event_id=%s: %w", triggerID, eventID, err)
	}
	return nil
}

func (s *triggerEventStore) DeleteEventsForTrigger(ctx context.Context, triggerID string) error {
	const q = `
DELETE FROM ` + triggerPendingEventsTable + `
WHERE trigger_id = $1
`
	if _, err := s.ds.ExecContext(ctx, q, triggerID); err != nil {
		return fmt.Errorf("failed to delete pending events for trigger_id=%s: %w", triggerID, err)
	}
	return nil
}
