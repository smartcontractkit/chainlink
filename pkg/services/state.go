package services

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	pkgerrors "github.com/pkg/errors"
)

// defaultErrorBufferCap is the default cap on the errors an error buffer can store at any time
const defaultErrorBufferCap = 50

type errNotStarted struct {
	state State
}

func (e *errNotStarted) Error() string {
	return fmt.Sprintf("service is %q, not started", e.state)
}

var (
	ErrAlreadyStopped      = errors.New("already stopped")
	ErrCannotStopUnstarted = errors.New("cannot stop unstarted service")
)

// StateMachine contains a State integer
type StateMachine struct {
	state        atomic.Int32
	sync.RWMutex // lock is held during startup/shutdown, RLock is held while executing functions dependent on a particular state

	// SvcErrBuffer is an ErrorBuffer that let service owners track critical errors happening in the service.
	//
	// SvcErrBuffer.SetCap(int) Overrides buffer limit from defaultErrorBufferCap
	// SvcErrBuffer.Append(error) Appends an error to the buffer
	// SvcErrBuffer.Flush() error returns all tracked errors as a single joined error
	SvcErrBuffer ErrorBuffer
}

// State holds the state for StateMachine
type State int32

// nolint
const (
	stateUnstarted State = iota
	stateStarted
	stateStarting
	stateStartFailed
	stateStopping
	stateStopped
	stateStopFailed
)

func (s State) String() string {
	switch s {
	case stateUnstarted:
		return "Unstarted"
	case stateStarted:
		return "Started"
	case stateStarting:
		return "Starting"
	case stateStartFailed:
		return "StartFailed"
	case stateStopping:
		return "Stopping"
	case stateStopped:
		return "Stopped"
	case stateStopFailed:
		return "StopFailed"
	default:
		return fmt.Sprintf("unrecognized state: %d", s)
	}
}

// StartOnce sets the state to Started
func (once *StateMachine) StartOnce(name string, fn func() error) error {
	// SAFETY: We do this compare-and-swap outside of the lock so that
	// concurrent StartOnce() calls return immediately.
	success := once.state.CompareAndSwap(int32(stateUnstarted), int32(stateStarting))

	if !success {
		return pkgerrors.Errorf("%v has already been started once; state=%v", name, State(once.state.Load()))
	}

	once.Lock()
	defer once.Unlock()

	// Setting cap before calling startup fn in case of crits in startup
	once.SvcErrBuffer.SetCap(defaultErrorBufferCap)
	err := fn()

	if err == nil {
		success = once.state.CompareAndSwap(int32(stateStarting), int32(stateStarted))
	} else {
		success = once.state.CompareAndSwap(int32(stateStarting), int32(stateStartFailed))
	}

	if !success {
		// SAFETY: If this is reached, something must be very wrong: once.state
		// was tampered with outside of the lock.
		panic(fmt.Sprintf("%v entered unreachable state, unable to set state to started", name))
	}

	return err
}

// StopOnce sets the state to Stopped
func (once *StateMachine) StopOnce(name string, fn func() error) error {
	// SAFETY: We hold the lock here so that Stop blocks until StartOnce
	// executes. This ensures that a very fast call to Stop will wait for the
	// code to finish starting up before teardown.
	once.Lock()
	defer once.Unlock()

	success := once.state.CompareAndSwap(int32(stateStarted), int32(stateStopping))

	if !success {
		state := once.state.Load()
		switch state {
		case int32(stateStopped):
			return pkgerrors.Wrapf(ErrAlreadyStopped, "%s has already been stopped", name)
		case int32(stateUnstarted):
			return pkgerrors.Wrapf(ErrCannotStopUnstarted, "%s has not been started", name)
		default:
			return pkgerrors.Errorf("%v cannot be stopped from this state; state=%v", name, State(state))
		}
	}

	err := fn()

	if err == nil {
		success = once.state.CompareAndSwap(int32(stateStopping), int32(stateStopped))
	} else {
		success = once.state.CompareAndSwap(int32(stateStopping), int32(stateStopFailed))
	}

	if !success {
		// SAFETY: If this is reached, something must be very wrong: once.state
		// was tampered with outside of the lock.
		panic(fmt.Sprintf("%v entered unreachable state, unable to set state to stopped", name))
	}

	return err
}

// State retrieves the current state
func (once *StateMachine) State() State {
	state := once.state.Load()
	return State(state)
}

// IfStarted runs the func and returns true only if started, otherwise returns false
func (once *StateMachine) IfStarted(f func()) (ok bool) {
	once.RLock()
	defer once.RUnlock()

	state := once.state.Load()

	if State(state) == stateStarted {
		f()
		return true
	}
	return false
}

// IfNotStopped runs the func and returns true if in any state other than Stopped
func (once *StateMachine) IfNotStopped(f func()) (ok bool) {
	once.RLock()
	defer once.RUnlock()

	state := once.state.Load()

	if State(state) == stateStopped {
		return false
	}
	f()
	return true
}

// Ready returns ErrNotStarted if the state is not started.
func (once *StateMachine) Ready() error {
	state := once.State()
	if state == stateStarted {
		return nil
	}
	return &errNotStarted{state: state}
}

// Healthy returns ErrNotStarted if the state is not started.
// Override this per-service with more specific implementations.
func (once *StateMachine) Healthy() error {
	state := once.State()
	if state == stateStarted {
		return once.SvcErrBuffer.Flush()
	}
	return &errNotStarted{state: state}
}

// ErrorBuffer uses joinedErrors interface to join multiple errors into a single error.
// This is useful to track the most recent N errors in a service and flush them as a single error.
type ErrorBuffer struct {
	// buffer is a slice of errors
	buffer []error

	// cap is the maximum number of errors that the buffer can hold.
	// Exceeding the cap results in discarding the oldest error
	cap int

	mu sync.RWMutex
}

func (eb *ErrorBuffer) Flush() (err error) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	err = errors.Join(eb.buffer...)
	eb.buffer = nil
	return
}

func (eb *ErrorBuffer) Append(incoming error) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if len(eb.buffer) == eb.cap && eb.cap != 0 {
		eb.buffer = append(eb.buffer[1:], incoming)
		return
	}
	eb.buffer = append(eb.buffer, incoming)
}

func (eb *ErrorBuffer) SetCap(cap int) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	if len(eb.buffer) > cap {
		eb.buffer = eb.buffer[len(eb.buffer)-cap:]
	}
	eb.cap = cap
}
