package triggers

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// fakeQueue is a minimal Queue whose Observe pops a scripted sequence of
// results and counts calls. Once the script is exhausted it returns empty.
type fakeQueue struct {
	mu      sync.Mutex
	results []observeResult
	calls   int
}

type observeResult struct {
	events []CoordinatedEvent
	err    error
}

func (q *fakeQueue) Put(context.Context, CoordinatedEvent) error { return nil }

func (q *fakeQueue) Stats(context.Context) (Stats, error) { return Stats{}, nil }

func (q *fakeQueue) Observe(context.Context) ([]CoordinatedEvent, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.calls++
	if len(q.results) == 0 {
		return nil, nil
	}
	r := q.results[0]
	q.results = q.results[1:]
	return r.events, r.err
}

func (q *fakeQueue) observeCalls() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.calls
}

func (q *fakeQueue) remaining() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.results)
}

// batchOf returns a one-event batch tagged with id.
func batchOf(id string) []CoordinatedEvent { return []CoordinatedEvent{{WorkflowID: id}} }

// TestDrainer_Drain calls drain directly (no worker running), so the batches
// it hands off stay in the channel for inspection.
func TestDrainer_Drain(t *testing.T) {
	t.Parallel()

	// queued drains and returns the batches waiting for the hook worker, without
	// needing the worker to be running.
	queued := func(d *drainer) [][]CoordinatedEvent {
		var out [][]CoordinatedEvent
		for {
			select {
			case b := <-d.batches:
				out = append(out, b)
			default:
				return out
			}
		}
	}

	tests := []struct {
		name    string
		results []observeResult // one per drain call
		want    [][]CoordinatedEvent
	}{
		{
			name:    "queues an observed batch",
			results: []observeResult{{events: batchOf("a")}},
			want:    [][]CoordinatedEvent{batchOf("a")},
		},
		{
			name:    "empty observe queues nothing",
			results: []observeResult{{}},
		},
		{
			name:    "observe error queues nothing, even with events",
			results: []observeResult{{events: batchOf("ignored"), err: errors.New("boom")}},
		},
		{
			name:    "full buffer drops the overflow and keeps the oldest",
			results: []observeResult{{events: batchOf("a")}, {events: batchOf("b")}},
			want:    [][]CoordinatedEvent{batchOf("a")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := &fakeQueue{results: tt.results}
			var hookCalled atomic.Bool
			d := NewDrainer(q, time.Hour, drainerHooks{OnObservedEvents: func(context.Context, []CoordinatedEvent) {
				hookCalled.Store(true)
			}}, logger.Test(t))

			for range tt.results {
				d.drain(t.Context()) // must return even when the buffer is full
			}

			assert.Equal(t, len(tt.results), q.observeCalls(), "one Observe per drain")
			assert.Equal(t, tt.want, queued(d))
			assert.False(t, hookCalled.Load(), "drain never calls the hook itself; only the worker does")
		})
	}
}

// TestDrainer_Running starts the drainer and checks what its hook receives.
func TestDrainer_Running(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		results []observeResult
		want    []string // workflow IDs the hook should see
	}{
		{
			name:    "hook receives an observed batch",
			results: []observeResult{{events: batchOf("a")}},
			want:    []string{"a"},
		},
		{
			name:    "loop survives an observe error",
			results: []observeResult{{err: errors.New("boom")}, {events: batchOf("a")}},
			want:    []string{"a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := &fakeQueue{results: tt.results}
			received := make(chan string, len(tt.results))
			d := NewDrainer(q, 5*time.Millisecond, drainerHooks{OnObservedEvents: func(_ context.Context, events []CoordinatedEvent) {
				for _, e := range events {
					received <- e.WorkflowID
				}
			}}, logger.Test(t))

			require.NoError(t, d.Start(t.Context()))
			t.Cleanup(func() { assert.NoError(t, d.Close()) })

			var got []string
			for range tt.want {
				select {
				case id := <-received:
					got = append(got, id)
				case <-time.After(5 * time.Second):
					t.Fatalf("timed out; hook saw %v, want %v", got, tt.want)
				}
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestDrainer_BlockedHook uses a hook that blocks until its context is
// cancelled. It checks the drain loop is not stalled by the hook, and that
// Close cancels the hook and waits for it to return.
func TestDrainer_BlockedHook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		batches      int
		maxHookCalls int // the batch the worker holds plus whatever fits in the buffer
	}{
		{name: "single batch", batches: 1, maxHookCalls: 1},
		{name: "overflow batches are dropped", batches: 10, maxHookCalls: _buffer + 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			results := make([]observeResult, tt.batches)
			for i := range results {
				results[i] = observeResult{events: batchOf("x")}
			}
			q := &fakeQueue{results: results}

			var once sync.Once
			entered := make(chan struct{})
			var inFlight, calls atomic.Int32
			d := NewDrainer(q, 5*time.Millisecond, drainerHooks{OnObservedEvents: func(ctx context.Context, _ []CoordinatedEvent) {
				calls.Add(1)
				inFlight.Add(1)
				defer inFlight.Add(-1)
				once.Do(func() { close(entered) })
				<-ctx.Done()
			}}, logger.Test(t))

			require.NoError(t, d.Start(t.Context()))
			<-entered // the hook is now blocked

			// The loop keeps observing and drains the queue even though
			// the hook never returns.
			require.Eventually(t, func() bool { return q.remaining() == 0 },
				5*time.Second, 5*time.Millisecond)

			// unblock the hook by closing the drainer
			require.NoError(t, d.Close())
			assert.Zero(t, inFlight.Load(), "Close returned while the hook was still running")
			assert.LessOrEqual(t, int(calls.Load()), tt.maxHookCalls, "overflow batches must be dropped, not queued without bound")
		})
	}
}
