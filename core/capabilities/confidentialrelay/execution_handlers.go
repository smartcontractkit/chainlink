package confidentialrelay

import (
	"context"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/host"
)

// ExecutionHandlers tracks the confidential executions in flight on this node so
// the relay handler can service the enclave's secrets-get / capability-execute
// callbacks. An execution is registered for the lifetime of its
// ConfidentialModule.Execute call (AddExecution on entry, RemoveExecution on
// return).
//
// The enclave runs a single, DON-shared execution and only needs a relay quorum
// to proceed, so its callback (fanned out by the gateway to every DON member)
// can reach a node before that node has started its own copy of the execution —
// i.e. before AddExecution has run. GetExecutionWithWait bridges that start-edge
// race by briefly waiting for the handler to appear instead of failing the
// callback outright, which would drop that node's signature and erode the relay
// quorum. The zero value is ready to use.
type ExecutionHandlers struct {
	mu       sync.Mutex
	handlers map[string]host.ExecutionHelperWithRawSecrets
	// waiters holds the wake channels of callers parked in GetExecutionWithWait
	// for a not-yet-registered key. AddExecution closes and clears them.
	waiters map[string][]chan struct{}
}

func (e *ExecutionHandlers) AddExecution(workflowID, execID string, helper host.ExecutionHelperWithRawSecrets) {
	key := wfexecID(workflowID, execID)
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.handlers == nil {
		e.handlers = make(map[string]host.ExecutionHelperWithRawSecrets)
	}
	e.handlers[key] = helper
	// Wake any callbacks that arrived before this execution registered.
	for _, ch := range e.waiters[key] {
		close(ch)
	}
	delete(e.waiters, key)
}

func (e *ExecutionHandlers) RemoveExecution(workflowID, execID string) {
	key := wfexecID(workflowID, execID)
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.handlers, key)
}

// GetExecution returns the registered handler for (workflowID, execID) without
// waiting. Prefer GetExecutionWithWait on the relay callback path.
func (e *ExecutionHandlers) GetExecution(workflowID, execID string) (host.ExecutionHelperWithRawSecrets, bool) {
	key := wfexecID(workflowID, execID)
	e.mu.Lock()
	defer e.mu.Unlock()
	helper, ok := e.handlers[key]
	return helper, ok
}

// GetExecutionWithWait returns the handler for (workflowID, execID), waiting for
// it to be registered if it is not present yet, until ctx is done. The caller is
// expected to bound ctx (a few seconds, well under the gateway relay request
// timeout) so a callback for an execution this node never runs still fails in
// bounded time. It returns (nil, false) on timeout or cancellation.
func (e *ExecutionHandlers) GetExecutionWithWait(ctx context.Context, workflowID, execID string) (host.ExecutionHelperWithRawSecrets, bool) {
	key := wfexecID(workflowID, execID)

	// Check and, on a miss, register the waiter under a single lock hold. Because
	// AddExecution takes the same lock, it cannot slip in between the check and
	// the registration: it either already ran (handled by the hit below) or it
	// will observe our channel and close it. This closes the lost-wakeup window.
	e.mu.Lock()
	if helper, ok := e.handlers[key]; ok {
		e.mu.Unlock()
		return helper, true
	}
	ch := make(chan struct{})
	if e.waiters == nil {
		e.waiters = make(map[string][]chan struct{})
	}
	e.waiters[key] = append(e.waiters[key], ch)
	e.mu.Unlock()

	select {
	case <-ch:
		// Registered while we waited. Re-read rather than assume it is still
		// present: a very short execution may have already been removed.
		return e.GetExecution(workflowID, execID)
	case <-ctx.Done():
		e.removeWaiter(key, ch)
		// AddExecution may have raced the deadline; prefer a hit over a spurious
		// miss.
		return e.GetExecution(workflowID, execID)
	}
}

// removeWaiter drops ch from the waiter list for key. It is a no-op if
// AddExecution already consumed (closed and cleared) it, so ch is never
// double-closed.
func (e *ExecutionHandlers) removeWaiter(key string, ch chan struct{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	ws := e.waiters[key]
	for i, c := range ws {
		if c == ch {
			e.waiters[key] = append(ws[:i], ws[i+1:]...)
			break
		}
	}
	if len(e.waiters[key]) == 0 {
		delete(e.waiters, key)
	}
}

func wfexecID(workflowID, execID string) string {
	return workflowID + "." + execID
}
