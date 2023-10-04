package utils

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

type errNotStarted struct {
	state startStopOnceState
}

func (e *errNotStarted) Error() string {
	return fmt.Sprintf("service is %q, not started", e.state)
}

// startStopOnceState holds the state for StartStopOnce
type startStopOnceState int32

const (
	startStopOnceUnstarted startStopOnceState = iota
	startStopOnceStarted
	startStopOnceStarting
	startStopOnceStopping
	startStopOnceStopped
)

func (s startStopOnceState) String() string {
	switch s {
	case startStopOnceUnstarted:
		return "Unstarted"
	case startStopOnceStarted:
		return "Started"
	case startStopOnceStarting:
		return "Starting"
	case startStopOnceStopping:
		return "Stopping"
	case startStopOnceStopped:
		return "Stopped"
	default:
		return fmt.Sprintf("unrecognized state: %d", s)
	}
}

// StartStopOnce can be embedded in a struct to help implement types.Service.
// Deprecated: use services.StateMachine
type StartStopOnce struct {
	state        atomic.Int32
	sync.RWMutex // lock is held during startup/shutdown, RLock is held while executing functions dependent on a particular state
}

// StartOnce sets the state to Started
func (s *StartStopOnce) StartOnce(name string, fn func() error) error {
	// SAFETY: We do this compare-and-swap outside of the lock so that
	// concurrent StartOnce() calls return immediately.
	success := s.state.CompareAndSwap(int32(startStopOnceUnstarted), int32(startStopOnceStarting))

	if !success {
		return errors.Errorf("%v has already started once", name)
	}

	s.Lock()
	defer s.Unlock()

	err := fn()

	success = s.state.CompareAndSwap(int32(startStopOnceStarting), int32(startStopOnceStarted))

	if !success {
		// SAFETY: If this is reached, something must be very wrong: once.state
		// was tampered with outside of the lock.
		panic(fmt.Sprintf("%v entered unreachable state, unable to set state to started", name))
	}

	return err
}

// StopOnce sets the state to Stopped
func (s *StartStopOnce) StopOnce(name string, fn func() error) error {
	// SAFETY: We hold the lock here so that Stop blocks until StartOnce
	// executes. This ensures that a very fast call to Stop will wait for the
	// code to finish starting up before teardown.
	s.Lock()
	defer s.Unlock()

	success := s.state.CompareAndSwap(int32(startStopOnceStarted), int32(startStopOnceStopping))

	if !success {
		return errors.Errorf("%v is unstarted or has already stopped once", name)
	}

	err := fn()

	success = s.state.CompareAndSwap(int32(startStopOnceStopping), int32(startStopOnceStopped))

	if !success {
		// SAFETY: If this is reached, something must be very wrong: once.state
		// was tampered with outside of the lock.
		panic(fmt.Sprintf("%v entered unreachable state, unable to set state to stopped", name))
	}

	return err
}

// State retrieves the current state
func (s *StartStopOnce) State() startStopOnceState {
	state := s.state.Load()
	return startStopOnceState(state)
}

// IfStarted runs the func and returns true only if started, otherwise returns false
func (s *StartStopOnce) IfStarted(f func()) (ok bool) {
	s.RLock()
	defer s.RUnlock()

	state := s.state.Load()

	if startStopOnceState(state) == startStopOnceStarted {
		f()
		return true
	}
	return false
}

// IfNotStopped runs the func and returns true if in any state other than Stopped
func (s *StartStopOnce) IfNotStopped(f func()) (ok bool) {
	s.RLock()
	defer s.RUnlock()

	state := s.state.Load()

	if startStopOnceState(state) == startStopOnceStopped {
		return false
	}
	f()
	return true
}

// Ready returns ErrNotStarted if the state is not started.
func (s *StartStopOnce) Ready() error {
	state := s.State()
	if state == startStopOnceStarted {
		return nil
	}
	return &errNotStarted{state: state}
}

// Healthy returns ErrNotStarted if the state is not started.
// Override this per-service with more specific implementations.
func (s *StartStopOnce) Healthy() error {
	state := s.State()
	if state == startStopOnceStarted {
		return nil
	}
	return &errNotStarted{state: state}
}
