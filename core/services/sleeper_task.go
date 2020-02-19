package services

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	stopWaitTime = 5 * time.Second
)

// SleeperTask represents a task that waits in the background to process some work.
type SleeperTask interface {
	Start() error
	Stop() error
	WakeUp() error
}

// Worker is a simple interface that represents some work to do repeatedly
type Worker interface {
	Work()
}

type sleeperTask struct {
	worker  Worker
	waker   chan struct{}
	closer  chan struct{}
	closed  chan struct{}
	started bool
	stopped bool
	mutex   sync.Mutex
}

// NewSleeperTask takes a worker and returns a SleeperTask.
//
// SleeperTask is guaranteed to call Work on the worker at least once for every
// WakeUp call.
// If the Worker is busy when WakeUp is called, the Worker will be called again
// immediately after it is finished. For this reason you should take care to
// make sure that Worker is idempotent.
// WakeUp does not block.
//
func NewSleeperTask(worker Worker) SleeperTask {
	return &sleeperTask{
		worker:  worker,
		waker:   make(chan struct{}, 1),
		closer:  make(chan struct{}, 1),
		closed:  make(chan struct{}, 1),
		started: false,
		stopped: false,
		mutex:   sync.Mutex{},
	}
}

// Start begins the SleeperTask
// Calling start on a started task will return an error
func (s *sleeperTask) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.stopped {
		return errors.New("cannot start a sleeper task that has been stopped")
	} else if s.started {
		return errors.New("sleeper task is already started")
	}
	s.started = true
	go s.workerLoop()
	return nil
}

// Stop stops the SleeperTask
// It can only be called once
// It will block until all tasks have completed, or the timeout is reached
// Stopped tasks cannot be restarted
func (s *sleeperTask) Stop() error {
	s.mutex.Lock()
	if s.stopped {
		s.mutex.Unlock()
		return errors.New("sleeper task is already stopped")
	}
	s.stopped = true
	s.mutex.Unlock()

	s.closer <- struct{}{}
	// NOTE: Closing the channels will cause the rogue task to panic if it does eventually complete
	defer close(s.waker)
	defer close(s.closer)
	defer close(s.closed)
	select {
	case <-s.closed:
		return nil
	case <-time.After(stopWaitTime):
		return fmt.Errorf("task did not complete within %s", stopWaitTime.String())
	}
}

// WakeUp wakes up the sleeper task, asking it to execute its Worker.
// Idempotent, can be called multiple times but will only wake the worker once (until the worker finishes again)
func (s *sleeperTask) WakeUp() error {
	s.mutex.Lock()
	if s.stopped {
		s.mutex.Unlock()
		return errors.New("cannot wake up stopped sleeper task")
	}
	s.mutex.Unlock()

	select {
	case s.waker <- struct{}{}:
	default:
	}
	return nil
}

// workerLoop is the goroutine behind the sleeper task that waits for a signal
// before kicking off the worker
func (s *sleeperTask) workerLoop() {
	for {
		select {
		case <-s.waker:
			s.worker.Work()
		case <-s.closer:
			s.closed <- struct{}{}
			return
		}
	}
}
