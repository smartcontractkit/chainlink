package beholder_test

/* TODO: CRE-4422 Refactor: relocate this to durableemitter pkg
import (
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	beholdersvc "github.com/smartcontractkit/chainlink/v2/core/services/beholder"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

// truncateChipDurableEvents clears the table so ORM tests stay deterministic when using a
// shared CL_DATABASE_URL (e.g. interrupted runs or parallel packages leaving rows behind).
func truncateChipDurableEvents(t *testing.T, db *sqlx.DB) {
	t.Helper()
	ctx := t.Context()
	_, err := db.ExecContext(ctx, `TRUNCATE TABLE cre.chip_durable_events RESTART IDENTITY`)
	require.NoError(t, err)
}

func TestPgDurableEventStore_InsertDeleteRoundTrip(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	truncateChipDurableEvents(t, db)
	ctx := t.Context()
	store := beholdersvc.NewPgDurableEventStore(db)

	id, err := store.Insert(ctx, []byte("test-payload"))
	require.NoError(t, err)
	require.Positive(t, id)

	events, err := store.ListPending(ctx, time.Now().Add(time.Second), 10)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, id, events[0].ID)
	assert.Equal(t, []byte("test-payload"), events[0].Payload)

	require.NoError(t, store.Delete(ctx, id))

	events, err = store.ListPending(ctx, time.Now().Add(time.Second), 10)
	require.NoError(t, err)
	assert.Empty(t, events)
}

func TestPgDurableEventStore_ListPending_RespectsCreatedBefore(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	truncateChipDurableEvents(t, db)
	ctx := t.Context()
	store := beholdersvc.NewPgDurableEventStore(db)

	_, err := store.Insert(ctx, []byte("event-1"))
	require.NoError(t, err)

	// createdBefore in the past should return nothing (event was just created).
	events, err := store.ListPending(ctx, time.Now().Add(-time.Hour), 10)
	require.NoError(t, err)
	assert.Empty(t, events)

	// createdBefore in the future should return the event.
	events, err = store.ListPending(ctx, time.Now().Add(time.Hour), 10)
	require.NoError(t, err)
	assert.Len(t, events, 1)
}

func TestPgDurableEventStore_ListPending_RespectsLimit(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	truncateChipDurableEvents(t, db)
	ctx := t.Context()
	store := beholdersvc.NewPgDurableEventStore(db)

	for i := 0; i < 20; i++ {
		_, err := store.Insert(ctx, []byte(fmt.Sprintf("event-%d", i)))
		require.NoError(t, err)
	}

	events, err := store.ListPending(ctx, time.Now().Add(time.Second), 5)
	require.NoError(t, err)
	assert.Len(t, events, 5)
}

func TestPgDurableEventStore_DeleteExpired(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	truncateChipDurableEvents(t, db)
	ctx := t.Context()
	store := beholdersvc.NewPgDurableEventStore(db)

	_, err := store.Insert(ctx, []byte("will-expire"))
	require.NoError(t, err)

	// TTL of 1 hour — nothing should be deleted (event is <1s old).
	deleted, err := store.DeleteExpired(ctx, time.Hour)
	require.NoError(t, err)
	assert.Equal(t, int64(0), deleted)

	// TTL of 0 — everything should be deleted.
	deleted, err = store.DeleteExpired(ctx, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(1), deleted)
}

func TestPgDurableEventStore_ObserveDurableQueue(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	truncateChipDurableEvents(t, db)
	ctx := testutils.Context(t)
	store := beholdersvc.NewPgDurableEventStore(db)

	st, err := store.ObserveDurableQueue(ctx, time.Hour, time.Minute)
	require.NoError(t, err)
	assert.Equal(t, int64(0), st.Depth)

	_, err = store.Insert(ctx, []byte("payload-bytes"))
	require.NoError(t, err)
	st, err = store.ObserveDurableQueue(ctx, time.Hour, time.Minute)
	require.NoError(t, err)
	assert.Equal(t, int64(1), st.Depth)
	assert.Equal(t, int64(len("payload-bytes")), st.PayloadBytes)
	assert.Positive(t, st.OldestPendingAge)
}

func TestPgDurableEventStore_MarkDeliveredAndPurgeDelivered(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	truncateChipDurableEvents(t, db)
	ctx := testutils.Context(t)
	store := beholdersvc.NewPgDurableEventStore(db)

	id, err := store.Insert(ctx, []byte("payload"))
	require.NoError(t, err)

	pending, err := store.ListPending(ctx, time.Now().Add(time.Hour), 10)
	require.NoError(t, err)
	require.Len(t, pending, 1)

	require.NoError(t, store.MarkDelivered(ctx, id))
	require.NoError(t, store.MarkDelivered(ctx, id), "second mark is idempotent")

	pending, err = store.ListPending(ctx, time.Now().Add(time.Hour), 10)
	require.NoError(t, err)
	require.Empty(t, pending)

	var cnt int64
	require.NoError(t, db.GetContext(ctx, &cnt, `SELECT count(*) FROM cre.chip_durable_events`))
	require.Equal(t, int64(1), cnt, "row remains as tombstone until purge")

	n, err := store.PurgeDelivered(ctx, 10)
	require.NoError(t, err)
	require.Equal(t, int64(1), n)

	require.NoError(t, db.GetContext(ctx, &cnt, `SELECT count(*) FROM cre.chip_durable_events`))
	require.Equal(t, int64(0), cnt)
}
*/
