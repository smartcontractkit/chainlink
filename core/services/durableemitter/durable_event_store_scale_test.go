package durableemitter_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/durableemitter"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

// Scale check for the batched expiry path: a backlog on the order of a real
// outage (250k rows, ~1KB payloads) drained in 5000-row slices. Opt-in because
// it takes tens of seconds and is a measurement, not a correctness gate.
// Run: PURGE_SCALE=1 CL_DATABASE_URL=... go test ./core/services/durableemitter/ -run Scale -v
func TestPgDurableEventStore_DeleteExpiredBatch_Scale(t *testing.T) {
	if os.Getenv("PURGE_SCALE") == "" {
		t.Skip("set PURGE_SCALE=1 to run")
	}
	db := pgtest.NewSqlxDB(t)
	truncateChipDurableEvents(t, db)
	ctx := t.Context()
	store := durableemitter.NewPgDurableEventStore(db)

	const rows = 250_000
	t0 := time.Now()
	_, err := db.ExecContext(ctx, `
INSERT INTO cre.chip_durable_events (payload, created_at)
SELECT decode(repeat('ab', 512), 'hex'), now() - interval '2 hours'
FROM generate_series(1, $1)`, rows)
	require.NoError(t, err)
	t.Logf("seeded %d rows in %s", rows, time.Since(t0).Round(time.Millisecond))

	const batch = 5000
	var total int
	var slowest time.Duration
	start := time.Now()
	for {
		s := time.Now()
		payloads, err := store.DeleteExpiredBatch(ctx, time.Hour, batch)
		require.NoError(t, err)
		d := time.Since(s)
		if d > slowest {
			slowest = d
		}
		total += len(payloads)
		if len(payloads) < batch {
			break
		}
	}
	elapsed := time.Since(start)
	require.Equal(t, rows, total)
	t.Logf("drained %d rows in %s across %d slices; slowest slice %s; mean %s/slice",
		total, elapsed.Round(time.Millisecond), (rows+batch-1)/batch, slowest.Round(time.Millisecond),
		(elapsed / time.Duration((rows+batch-1)/batch)).Round(time.Millisecond))
}
