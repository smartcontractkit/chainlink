package trigger_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	trigger "github.com/smartcontractkit/chainlink/v2/core/capabilities/triggers"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func makePendingEvent(triggerID, eventID string, payload []byte, offset time.Duration) capabilities.PendingEvent {
	now := time.Now().Add(offset)
	return capabilities.PendingEvent{
		TriggerId:  triggerID,
		EventId:    eventID,
		Payload:    append([]byte(nil), payload...),
		FirstAt:    now,
		LastSentAt: time.Time{},
		Attempts:   0,
	}
}

func TestTriggerEventStore_InsertListDelete_DeleteEventsForTrigger(t *testing.T) {
	ctx := context.Background()
	ds := pgtest.NewSqlxDB(t)
	store := trigger.NewTriggerEventStore(ds)

	// insert a few events across two triggers
	e1 := makePendingEvent("trigA", "e1", []byte("p1"), -3*time.Minute)
	e2 := makePendingEvent("trigA", "e2", []byte("p2"), -2*time.Minute)
	e3 := makePendingEvent("trigB", "e3", []byte("p3"), -1*time.Minute)

	require.NoError(t, store.Insert(ctx, e1))
	require.NoError(t, store.Insert(ctx, e2))
	require.NoError(t, store.Insert(ctx, e3))

	// List should return all three in first_at ascending order
	recs, err := store.List(ctx)
	require.NoError(t, err)
	require.Len(t, recs, 3)
	// assert order and payloads
	require.Equal(t, "e1", recs[0].EventId)
	require.Equal(t, []byte("p1"), recs[0].Payload)
	require.Equal(t, "e2", recs[1].EventId)
	require.Equal(t, []byte("p2"), recs[1].Payload)
	require.Equal(t, "e3", recs[2].EventId)
	require.Equal(t, []byte("p3"), recs[2].Payload)

	// Delete a single event and ensure it's gone
	require.NoError(t, store.DeleteEvent(ctx, "trigA", "e1"))
	recs, err = store.List(ctx)
	require.NoError(t, err)
	require.Len(t, recs, 2)
	for _, r := range recs {
		require.NotEqual(t, "e1", r.EventId)
	}

	// Delete all for triggerA
	require.NoError(t, store.DeleteEventsForTrigger(ctx, "trigA"))
	recs, err = store.List(ctx)
	require.NoError(t, err)
	require.Len(t, recs, 1)
	require.Equal(t, "e3", recs[0].EventId)
}

func TestTriggerEventStore_UpdateDelivery(t *testing.T) {
	ctx := context.Background()
	ds := pgtest.NewSqlxDB(t)
	store := trigger.NewTriggerEventStore(ds)

	triggerID := "trig-" + uuid.NewString()
	eventID := "evt-update"

	pe := makePendingEvent(triggerID, eventID, []byte("payload"), -1*time.Minute)
	require.NoError(t, store.Insert(ctx, pe))

	// Update delivery metadata
	attempts := 5
	lastSent := time.Now().UTC().Truncate(time.Second)
	require.NoError(t, store.UpdateDelivery(ctx, triggerID, eventID, lastSent, attempts))

	// Verify persisted fields
	recs, err := store.List(ctx)
	require.NoError(t, err)
	require.Len(t, recs, 1)
	require.Equal(t, attempts, recs[0].Attempts)
	require.True(t, recs[0].LastSentAt.Truncate(time.Second).Equal(lastSent))

	err = store.UpdateDelivery(ctx, "no-such-trigger", "no-such-event", time.Now(), 1)
	require.Error(t, err)
}
