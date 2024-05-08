package monitoring

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

const numPollerUpdates = 10
const numGoroutinesPerManaged = 10

func TestManager(t *testing.T) {
	t.Run("all goroutines are stopped before the new ones begin", func(t *testing.T) {
		// Poller fires 10 rounds of updates.
		// The manager identifies these updates, terminates the current running managed function and starts a new one.
		// The managed function in turn runs 10 noop goroutines and increments/decrements a goroutine counter.
		defer goleak.VerifyNone(t)

		var goRoutineCounter int64
		var subs utils.Subprocesses
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		poller := &fakePoller{
			numPollerUpdates,
			make(chan interface{}),
		}
		subs.Go(func() {
			poller.Run(ctx)
		})

		manager := NewManager(
			newNullLogger(),
			poller,
		)
		managed := func(ctx context.Context, _ RDDData) {
			var localSubs utils.Subprocesses
			defer localSubs.Wait()
			for i := 0; i < numGoroutinesPerManaged; i++ {
				localSubs.Go(func() {
					atomic.AddInt64(&goRoutineCounter, 1)
					<-ctx.Done()
					atomic.AddInt64(&goRoutineCounter, -1)
				})
			}
		}
		subs.Go(func() {
			manager.Run(ctx, managed)
		})

		subs.Wait()
		require.Equal(t, int64(0), goRoutineCounter, "all child goroutines are gone")
	})

	t.Run("should not restart the monitor if the feeds are the same", func(t *testing.T) {
		feeds := []FeedConfig{
			generateFeedConfig(),
			generateFeedConfig(),
		}
		nodes := []NodeConfig{generateNodeConfig()}
		rddPoller := &fakePoller{0, make(chan interface{})}
		manager := NewManager(
			newNullLogger(),
			rddPoller,
		)

		var countManagedFuncExecutions uint64
		var managedFunc = func(_ context.Context, _ RDDData) {
			atomic.AddUint64(&countManagedFuncExecutions, 1)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		var subs utils.Subprocesses
		subs.Go(func() {
			manager.Run(ctx, managedFunc)
		})

		// The rdd poller returns the same feed configs three times!
		for i := 0; i < 3; i++ {
			select {
			case rddPoller.ch <- RDDData{feeds, nodes}:
			case <-ctx.Done():
			}
		}

		cancel()
		subs.Wait()

		require.Equal(t, countManagedFuncExecutions, uint64(1))
	})

	t.Run("should expose the current feeds to http", func(t *testing.T) {
		feeds := []FeedConfig{generateFeedConfig()}
		nodes := []NodeConfig{generateNodeConfig()}
		manager := &managerImpl{
			newNullLogger(),
			&fakePoller{0, make(chan interface{})},
			RDDData{feeds, nodes},
			sync.Mutex{},
		}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/debug", nil)
		manager.HTTPHandler().ServeHTTP(rec, req)
		type rddData struct {
			Feeds []fakeFeedConfig
			Nodes []fakeNodeConfig
		}
		dec := json.NewDecoder(rec.Body)
		decodedData := rddData{}
		err := dec.Decode(&decodedData)
		require.NoError(t, err)
		require.Equal(t, len(decodedData.Feeds), len(feeds))
		require.Equal(t, len(decodedData.Nodes), len(nodes))
	})

	t.Run("manager can manage multiple functions", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		var subs utils.Subprocesses
		ctx, cancel := context.WithCancel(tests.Context(t))

		poller := &fakePoller{ch: make(chan interface{})}

		manager := NewManager(logger.Test(t), poller)

		// run two managed funcs
		var createWG sync.WaitGroup
		var closeWG sync.WaitGroup
		managed := func(ctx context.Context, _ RDDData) {
			createWG.Done()
			<-ctx.Done()
			closeWG.Done()
		}
		managedNonBlocking := func(_ context.Context, _ RDDData) {
			createWG.Done()
			closeWG.Done()
		}
		subs.Go(func() {
			manager.Run(ctx, managed, managed, managedNonBlocking)
		})

		// send RDD update to create multiple managed funcs
		createWG.Add(3) // expect to see 3 created
		closeWG.Add(3)  // expect to see 3 closed on restart
		poller.ch <- RDDData{Feeds: []FeedConfig{generateFeedConfig()}}
		createWG.Wait() // wait for created

		// ensure stop + restarting works
		createWG.Add(3)                                                 // expect 3 new created
		closeWG.Add(3)                                                  // expect 3 closed on shutdown
		poller.ch <- RDDData{Feeds: []FeedConfig{generateFeedConfig()}} // trigger restart
		createWG.Wait()

		cancel()       // shutdown
		closeWG.Wait() // wait for managed funcs
		subs.Wait()    // wait for manager
	})
}
