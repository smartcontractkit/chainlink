package beholder_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	beholdersvc "github.com/smartcontractkit/chainlink/v2/core/services/beholder"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func TestPgDurableEventStore_InsertDeleteRoundTrip(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	store := beholdersvc.NewPgDurableEventStore(db)

	id, err := store.Insert(ctx, []byte("test-payload"))
	require.NoError(t, err)
	require.Greater(t, id, int64(0))

	events, err := store.ListPending(ctx, time.Now().Add(time.Second), 10)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, id, events[0].ID)
	assert.Equal(t, []byte("test-payload"), events[0].Payload)

	require.NoError(t, store.Delete(ctx, id))

	events, err = store.ListPending(ctx, time.Now().Add(time.Second), 10)
	require.NoError(t, err)
	assert.Len(t, events, 0)
}

func TestPgDurableEventStore_ListPending_RespectsCreatedBefore(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	store := beholdersvc.NewPgDurableEventStore(db)

	_, err := store.Insert(ctx, []byte("event-1"))
	require.NoError(t, err)

	// createdBefore in the past should return nothing (event was just created).
	events, err := store.ListPending(ctx, time.Now().Add(-time.Hour), 10)
	require.NoError(t, err)
	assert.Len(t, events, 0)

	// createdBefore in the future should return the event.
	events, err = store.ListPending(ctx, time.Now().Add(time.Hour), 10)
	require.NoError(t, err)
	assert.Len(t, events, 1)
}

func TestPgDurableEventStore_ListPending_RespectsLimit(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
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
	ctx := testutils.Context(t)
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
	require.Len(t, pending, 0)

	var cnt int64
	require.NoError(t, db.GetContext(ctx, &cnt, `SELECT count(*) FROM cre.chip_durable_events`))
	require.Equal(t, int64(1), cnt, "row remains as tombstone until purge")

	n, err := store.PurgeDelivered(ctx, 10)
	require.NoError(t, err)
	require.Equal(t, int64(1), n)

	require.NoError(t, db.GetContext(ctx, &cnt, `SELECT count(*) FROM cre.chip_durable_events`))
	require.Equal(t, int64(0), cnt)
}

// ---------- Benchmarks ----------

func randomPayload(size int) []byte {
	buf := make([]byte, size)
	_, _ = rand.Read(buf)
	return buf
}

// Benchmark_Insert measures raw INSERT throughput for individual events.
func Benchmark_Insert(b *testing.B) {
	db := pgtest.NewSqlxDB(b)
	ctx := testutils.Context(b)
	store := beholdersvc.NewPgDurableEventStore(db)
	payload := randomPayload(256)

	b.ResetTimer()
	for b.Loop() {
		_, err := store.Insert(ctx, payload)
		require.NoError(b, err)
	}
}

// Benchmark_InsertDelete measures the insert + delete cycle (the hot path when
// events are delivered successfully on the first attempt).
func Benchmark_InsertDelete(b *testing.B) {
	db := pgtest.NewSqlxDB(b)
	ctx := testutils.Context(b)
	store := beholdersvc.NewPgDurableEventStore(db)
	payload := randomPayload(256)

	b.ResetTimer()
	for b.Loop() {
		id, err := store.Insert(ctx, payload)
		require.NoError(b, err)
		require.NoError(b, store.Delete(ctx, id))
	}
}

// Benchmark_InsertPayloadSizes measures INSERT throughput at different payload sizes
// to understand how payload size affects DB performance.
func Benchmark_InsertPayloadSizes(b *testing.B) {
	sizes := []int{64, 256, 1024, 4096}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%dB", size), func(b *testing.B) {
			db := pgtest.NewSqlxDB(b)
			ctx := testutils.Context(b)
			store := beholdersvc.NewPgDurableEventStore(db)
			payload := randomPayload(size)

			b.ResetTimer()
			for b.Loop() {
				_, err := store.Insert(ctx, payload)
				require.NoError(b, err)
			}
		})
	}
}

// Benchmark_ListPending measures query performance with varying store depths.
func Benchmark_ListPending(b *testing.B) {
	depths := []int{100, 1000}
	for _, depth := range depths {
		b.Run(fmt.Sprintf("depth_%d", depth), func(b *testing.B) {
			db := pgtest.NewSqlxDB(b)
			ctx := testutils.Context(b)
			store := beholdersvc.NewPgDurableEventStore(db)
			payload := randomPayload(256)

			for i := 0; i < depth; i++ {
				_, err := store.Insert(ctx, payload)
				require.NoError(b, err)
			}

			b.ResetTimer()
			for b.Loop() {
				_, err := store.ListPending(ctx, time.Now().Add(time.Second), 100)
				require.NoError(b, err)
			}
		})
	}
}

// ---------- Load tests ----------

// TestLoad_SustainedInsertDelete simulates the durable emitter's steady-state:
// concurrent inserts with concurrent deletes, measuring achieved throughput
// and verifying the store drains cleanly.
func TestLoad_SustainedInsertDelete(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	store := beholdersvc.NewPgDurableEventStore(db)

	const (
		totalEvents = 2000
		concurrency = 10
	)

	payload := randomPayload(256)
	ids := make(chan int64, totalEvents)
	var insertCount, deleteCount atomic.Int64

	start := time.Now()

	// Producer goroutines: insert events.
	var insertWg sync.WaitGroup
	for w := 0; w < concurrency; w++ {
		insertWg.Add(1)
		go func() {
			defer insertWg.Done()
			for i := 0; i < totalEvents/concurrency; i++ {
				id, err := store.Insert(ctx, payload)
				if err != nil {
					t.Errorf("insert failed: %v", err)
					return
				}
				insertCount.Add(1)
				ids <- id
			}
		}()
	}

	// Consumer goroutines: delete events as they're inserted.
	var deleteWg sync.WaitGroup
	for w := 0; w < concurrency; w++ {
		deleteWg.Add(1)
		go func() {
			defer deleteWg.Done()
			for id := range ids {
				if err := store.Delete(ctx, id); err != nil {
					t.Errorf("delete failed: %v", err)
					return
				}
				deleteCount.Add(1)
			}
		}()
	}

	insertWg.Wait()
	close(ids)
	deleteWg.Wait()

	elapsed := time.Since(start)
	insertRate := float64(insertCount.Load()) / elapsed.Seconds()
	deleteRate := float64(deleteCount.Load()) / elapsed.Seconds()

	t.Logf("--- Load Test Results ---")
	t.Logf("Total events:     %d", totalEvents)
	t.Logf("Concurrency:      %d", concurrency)
	t.Logf("Elapsed:          %s", elapsed.Round(time.Millisecond))
	t.Logf("Insert rate:      %.0f events/sec", insertRate)
	t.Logf("Delete rate:      %.0f events/sec", deleteRate)
	t.Logf("Insert+Delete:    %.0f ops/sec (combined)", insertRate+deleteRate)

	assert.Equal(t, int64(totalEvents), insertCount.Load())
	assert.Equal(t, int64(totalEvents), deleteCount.Load())

	// Verify store is fully drained.
	remaining, err := store.ListPending(ctx, time.Now().Add(time.Hour), totalEvents)
	require.NoError(t, err)
	assert.Len(t, remaining, 0, "store should be empty after load test")
}

// TestLoad_BurstThenDrain simulates Chip going down: a burst of inserts with
// no deletes (events pile up), then a drain phase where everything is deleted
// via ListPending + batch Delete.
func TestLoad_BurstThenDrain(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	store := beholdersvc.NewPgDurableEventStore(db)

	const burstSize = 1000
	payload := randomPayload(512)

	// Phase 1: burst insert (simulates events arriving while Chip is down).
	burstStart := time.Now()
	for i := 0; i < burstSize; i++ {
		_, err := store.Insert(ctx, payload)
		require.NoError(t, err)
	}
	burstElapsed := time.Since(burstStart)
	t.Logf("Burst insert:  %d events in %s (%.0f events/sec)",
		burstSize, burstElapsed.Round(time.Millisecond),
		float64(burstSize)/burstElapsed.Seconds())

	// Phase 2: drain via ListPending + Delete (simulates retransmit loop).
	drainStart := time.Now()
	totalDrained := 0
	for {
		batch, err := store.ListPending(ctx, time.Now().Add(time.Second), 100)
		require.NoError(t, err)
		if len(batch) == 0 {
			break
		}
		for _, e := range batch {
			require.NoError(t, store.Delete(ctx, e.ID))
		}
		totalDrained += len(batch)
	}
	drainElapsed := time.Since(drainStart)
	t.Logf("Drain:         %d events in %s (%.0f events/sec)",
		totalDrained, drainElapsed.Round(time.Millisecond),
		float64(totalDrained)/drainElapsed.Seconds())

	assert.Equal(t, burstSize, totalDrained)
}

// TestLoad_ConcurrentInsertWithListPending simulates the real contention pattern:
// inserts happening concurrently with ListPending queries from the retransmit loop.
func TestLoad_ConcurrentInsertWithListPending(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	store := beholdersvc.NewPgDurableEventStore(db)

	const (
		duration    = 3 * time.Second
		concurrency = 5
	)

	payload := randomPayload(256)
	var insertCount, queryCount atomic.Int64

	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	var wg sync.WaitGroup

	// Inserters.
	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				if _, err := store.Insert(ctx, payload); err != nil {
					return // context cancelled
				}
				insertCount.Add(1)
			}
		}()
	}

	// ListPending poller (simulates retransmit loop).
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if _, err := store.ListPending(ctx, time.Now().Add(time.Second), 100); err != nil {
				return
			}
			queryCount.Add(1)
		}
	}()

	wg.Wait()

	t.Logf("--- Contention Test Results (%s) ---", duration)
	t.Logf("Inserts:          %d (%.0f/sec)", insertCount.Load(), float64(insertCount.Load())/duration.Seconds())
	t.Logf("ListPending calls: %d (%.0f/sec)", queryCount.Load(), float64(queryCount.Load())/duration.Seconds())
}
