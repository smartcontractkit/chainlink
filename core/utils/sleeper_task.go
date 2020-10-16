package utils

import (
	"sync"
)

// SleeperTask represents a task that waits in the background to process some work.
type SleeperTask interface {
	Stop() error
	WakeUp()
}

// Worker is a simple interface that represents some work to do repeatedly
type Worker interface {
	Work()
}

type sleeperTask struct {
	worker  Worker
	chQueue chan struct{}
	chStop  chan struct{}
	chDone  chan struct{}
	stopped bool
	stopMu  *sync.RWMutex
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
	s := &sleeperTask{
		worker:  worker,
		chQueue: make(chan struct{}, 1),
		chStop:  make(chan struct{}),
		chDone:  make(chan struct{}),
		stopMu:  new(sync.RWMutex),
	}

	go s.workerLoop()

	return s
}

// Stop stops the SleeperTask
// It never returns an error, this is simply to comply with the interface
func (s *sleeperTask) Stop() error {
	s.stopMu.Lock()
	defer s.stopMu.Unlock()
	if s.stopped {
		panic("already stopped")
	}
	close(s.chStop)
	<-s.chDone
	s.stopped = true
	return nil
}

// WakeUp wakes up the sleeper task, asking it to execute its Worker.
func (s *sleeperTask) WakeUp() {
	s.stopMu.RLock()
	defer s.stopMu.RUnlock()
	if s.stopped {
		panic("cannot wake up stopped sleeper task")
	}
	select {
	case s.chQueue <- struct{}{}:
	default:
	}
}

func (s *sleeperTask) workerLoop() {
	defer close(s.chDone)

	for {
		select {
		case <-s.chQueue:
			s.worker.Work()
		case <-s.chStop:
			return
		}
	}

	// FIXME: Unreachable code??

	if len(s.chQueue) > 0 {
		s.worker.Work()
	}
}

type SleeperTaskFuncWorker func()

func (fn SleeperTaskFuncWorker) Work() { fn() }
