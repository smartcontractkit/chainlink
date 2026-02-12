package trigger

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

type TriggerEventStore interface {
	capabilities.EventStore
}

// triggerEventStore is a Postgres-backed implementation of capabilities.EventStore.
type triggerEventStore struct {
	ds sqlutil.DataSource
}

var _ TriggerEventStore = (*triggerEventStore)(nil)

func NewTriggerEventStore(ds sqlutil.DataSource) triggerEventStore {
	return triggerEventStore{ds: ds}
}

const triggerPendingEventsTable = "trigger_pending_events"

func (s triggerEventStore) Insert(ctx context.Context, rec capabilities.PendingEvent) error {
	const q = `
INSERT INTO ` + triggerPendingEventsTable + ` (
  trigger_id, event_id, any_type_url, payload, first_at, last_sent_at, attempts
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT (trigger_id, event_id) DO UPDATE SET
  any_type_url = EXCLUDED.any_type_url,
  payload      = EXCLUDED.payload
`
	var lastSent sql.NullTime
	if !rec.LastSentAt.IsZero() {
		lastSent = sql.NullTime{Time: rec.LastSentAt, Valid: true}
	}

	if _, err := s.ds.ExecContext(
		ctx, q,
		rec.TriggerId,
		rec.EventId,
		rec.AnyTypeURL,
		rec.Payload,
		rec.FirstAt,
		lastSent,
		rec.Attempts,
	); err != nil {
		return fmt.Errorf("failed to insert pending event trigger_id=%s event_id=%s: %w", rec.TriggerId, rec.EventId, err)
	}
	return nil
}

func (s triggerEventStore) List(ctx context.Context) ([]capabilities.PendingEvent, error) {
	const q = `
SELECT
  trigger_id,
  event_id,
  any_type_url,
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
		AnyTypeURL string       `db:"any_type_url"`
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
			AnyTypeURL: r.AnyTypeURL,
			Payload:    append([]byte(nil), r.Payload...),
			FirstAt:    r.FirstAt,
			LastSentAt: last,
			Attempts:   r.Attempts,
		})
	}
	return out, nil
}

func (s triggerEventStore) DeleteEvent(ctx context.Context, triggerId, eventId string) error {
	const q = `
DELETE FROM ` + triggerPendingEventsTable + `
WHERE trigger_id = $1 AND event_id = $2
`
	if _, err := s.ds.ExecContext(ctx, q, triggerId, eventId); err != nil {
		return fmt.Errorf("failed to delete pending event trigger_id=%s event_id=%s: %w", triggerId, eventId, err)
	}
	return nil
}

func (s triggerEventStore) DeleteEventsForTrigger(ctx context.Context, triggerID string) error {
	const q = `
DELETE FROM ` + triggerPendingEventsTable + `
WHERE trigger_id = $1
`
	if _, err := s.ds.ExecContext(ctx, q, triggerID); err != nil {
		return fmt.Errorf("failed to delete pending events for trigger_id=%s: %w", triggerID, err)
	}
	return nil
}
