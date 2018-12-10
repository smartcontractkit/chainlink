package services

import (
	"sync"

	"github.com/smartcontractkit/chainlink/store"
)

// SleeperTask represents a task that waits in the background to process some work.
type SleeperTask interface {
	Start() error
	Stop() error
	WakeUp()
}

// Worker is a simple interface that represents some work to do repeatedly
type Worker interface {
	Work()
}

type sleeperTask struct {
	worker  Worker
	store   *store.Store
	cond    *sync.Cond
	started bool
}

// NewSleeperTask takes a worker and retruns a SleeperTask
func NewSleeperTask(worker Worker) SleeperTask {
	var m sync.Mutex
	return &sleeperTask{
		worker: worker,
		cond:   sync.NewCond(&m),
	}
}

// Start begins the SleeperTask
func (s *sleeperTask) Start() error {
	var wg sync.WaitGroup
	s.cond.L.Lock()
	s.started = true
	s.cond.L.Unlock()
	wg.Add(1)
	go s.workerLoop(&wg)
	wg.Wait()
	return nil
}

// Stop stops the SleeperTask
func (s *sleeperTask) Stop() error {
	s.cond.L.Lock()
	s.started = false
	s.cond.Signal()
	s.cond.L.Unlock()
	return nil
}

// WakeUp wakes up the sleeper task, asking it to execute its Worker.
func (s *sleeperTask) WakeUp() {
	go s.cond.Signal()
}

// workerLoop is the goroutine behind the sleeper task that waits for a signal
// before kicking off the worker
func (s *sleeperTask) workerLoop(wg *sync.WaitGroup) {
	for {
		s.cond.L.Lock()
		if wg != nil {
			wg.Done()
			wg = nil
		}
		s.cond.Wait()
		if s.started == false {
			s.cond.L.Unlock()
			return
		}
		s.worker.Work()
		s.cond.L.Unlock()
	}
}
