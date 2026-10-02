package v2

import "sync"

// inFlight tracks the executions inside ExecuteTrigger, whichever goroutine they
// run on: draining reads its count, and closing waits for it to reach zero. Once
// closed it refuses new executions. The zero value is ready to use.
type inFlight struct {
	mu     sync.Mutex
	idle   sync.Cond
	n      int32
	closed bool
}

// enter records a new execution, or returns false once closeAndWait has been called.
func (f *inFlight) enter() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return false
	}
	f.n++
	return true
}

// exit is the counterpart to a successful enter.
func (f *inFlight) exit() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.n--
	if f.n == 0 {
		f.idle.L = &f.mu
		f.idle.Broadcast()
	}
}

func (f *inFlight) count() int32 {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.n
}

// closeAndWait refuses new executions, then blocks until the running ones exit.
func (f *inFlight) closeAndWait() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.closed = true
	f.idle.L = &f.mu
	for f.n > 0 {
		f.idle.Wait()
	}
}
