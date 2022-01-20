package monitoring

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const numPollerUpdates = 10
const numGoroutinesPerManaged = 10

func TestManager(t *testing.T) {
	t.Run("all goroutines are stopped before the new ones begin", func(t *testing.T) {
		// Poller fires 10 rounds of updates.
		// The manager identifies these updates, terminates the current running managed function and starts a new one.
		// The managed function in turn runs 10 noop goroutines and increments/decrements a goroutine counter.

		var goRoutineCounter int64 = 0
		wg := &sync.WaitGroup{}
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		poller := &fakePoller{
			numPollerUpdates,
			make(chan interface{}),
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			poller.Start(ctx)
		}()

		manager := NewManager(
			newNullLogger(),
			poller,
		)
		managed := func(ctx context.Context, localWg *sync.WaitGroup, _ []FeedConfig) {
			localWg.Add(numGoroutinesPerManaged)
			for i := 0; i < numGoroutinesPerManaged; i++ {
				go func(i int, ctx context.Context) {
					defer localWg.Done()
					atomic.AddInt64(&goRoutineCounter, 1)
					<-ctx.Done()
					atomic.AddInt64(&goRoutineCounter, -1)
				}(i, ctx)
			}
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			manager.Start(ctx, wg, managed)
		}()

		wg.Wait()
		require.Equal(t, int64(0), goRoutineCounter, "all child goroutines are gone")
	})
}
