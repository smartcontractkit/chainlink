package confidentialrelay

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/host"
)

func TestExecutionHandlers_AddGetRemove(t *testing.T) {
	t.Parallel()

	var eh ExecutionHandlers
	helper := &mockExecutionHelper{}

	eh.AddExecution("wf", "exec", helper)
	got, ok := eh.GetExecution("wf", "exec")
	require.True(t, ok)
	assert.Equal(t, helper, got)

	eh.RemoveExecution("wf", "exec")
	_, ok = eh.GetExecution("wf", "exec")
	assert.False(t, ok)
}

func TestExecutionHandlers_GetExecutionWithWait_AlreadyPresent(t *testing.T) {
	t.Parallel()

	var eh ExecutionHandlers
	helper := &mockExecutionHelper{}
	eh.AddExecution("wf", "exec", helper)

	// A generous deadline that must not be consumed: the handler is already there.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	start := time.Now()
	got, ok := eh.GetExecutionWithWait(ctx, "wf", "exec")
	require.True(t, ok)
	assert.Equal(t, helper, got)
	assert.Less(t, time.Since(start), time.Second, "present handler should return without waiting")
}

// The enclave's relay callback can arrive before the node has started its
// execution; GetExecutionWithWait must return the handler once AddExecution
// registers it, without a lost wakeup.
func TestExecutionHandlers_GetExecutionWithWait_AppearsWhileWaiting(t *testing.T) {
	t.Parallel()

	var eh ExecutionHandlers
	helper := &mockExecutionHelper{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type result struct {
		helper host.ExecutionHelperWithRawSecrets
		ok     bool
	}
	resCh := make(chan result, 1)
	go func() {
		h, ok := eh.GetExecutionWithWait(ctx, "wf", "exec")
		resCh <- result{h, ok}
	}()

	time.Sleep(50 * time.Millisecond) // let the waiter park (correctness holds even if it doesn't)
	eh.AddExecution("wf", "exec", helper)

	select {
	case r := <-resCh:
		require.True(t, r.ok)
		assert.Equal(t, helper, r.helper)
	case <-time.After(2 * time.Second):
		t.Fatal("GetExecutionWithWait did not return after AddExecution")
	}
}

// A callback for an execution this node never runs must fail in bounded time and
// must not leak the waiter.
func TestExecutionHandlers_GetExecutionWithWait_Timeout(t *testing.T) {
	t.Parallel()

	var eh ExecutionHandlers
	// Capture start before the deadline is set so the elapsed lower-bound check
	// cannot flake on the sliver of time between WithTimeout and the call.
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, ok := eh.GetExecutionWithWait(ctx, "wf", "missing")
	elapsed := time.Since(start)
	assert.False(t, ok)
	assert.GreaterOrEqual(t, elapsed, 100*time.Millisecond)
	assert.Less(t, elapsed, 5*time.Second)

	eh.mu.Lock()
	_, present := eh.waiters[wfexecID("wf", "missing")]
	eh.mu.Unlock()
	assert.False(t, present, "timed-out waiter should be cleaned up")
}

// AddExecution must wake every parked waiter for the key.
func TestExecutionHandlers_GetExecutionWithWait_ManyWaiters(t *testing.T) {
	t.Parallel()

	var eh ExecutionHandlers
	helper := &mockExecutionHelper{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const n = 10
	var wg sync.WaitGroup
	oks := make([]bool, n)
	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, ok := eh.GetExecutionWithWait(ctx, "wf", "exec")
			oks[i] = ok
		}(i)
	}

	time.Sleep(50 * time.Millisecond)
	eh.AddExecution("wf", "exec", helper)
	wg.Wait()

	for i := range n {
		assert.True(t, oks[i], "waiter %d should have been woken", i)
	}
}
