package utils

import (
	"fmt"
	"time"
)

// SleeperTask represents a task that waits in the background to process some work.
type SleeperTask interface {
	Stop() error
	WakeUp()
	WakeUpIfStarted()
}

// Worker is a simple interface that represents some work to do repeatedly
type Worker interface {
	Work()
	Name() string
}

type sleeperTask struct {
	worker  Worker
	chQueue chan struct{}
	chStop  chan struct{}
	chDone  chan struct{}
	StartStopOnce
}

// NewSleeperTask takes a worker and returns a SleeperTask.
//
// SleeperTask is guaranteed to call Work on the worker at least once for every
// WakeUp call.
// If the Worker is busy when WakeUp is called, the Worker will be called again
// immediately after it is finished. For this reason you should take care to
// make sure that Worker is idempotent.
// WakeUp does not block.
func NewSleeperTask(worker Worker) SleeperTask {
	s := &sleeperTask{
		worker:  worker,
		chQueue: make(chan struct{}, 1),
		chStop:  make(chan struct{}),
		chDone:  make(chan struct{}),
	}

	_ = s.StartOnce("SleeperTask-"+worker.Name(), func() error {
		go s.workerLoop()
		return nil
	})

	return s
}

// Stop stops the SleeperTask
func (s *sleeperTask) Stop() error {
	return s.StopOnce("SleeperTask-"+s.worker.Name(), func() error {
		close(s.chStop)
		select {
		case <-s.chDone:
		case <-time.After(15 * time.Second):
			return fmt.Errorf("SleeperTask-%s took too long to stop", s.worker.Name())
		}
		return nil
	})
}

func (s *sleeperTask) WakeUpIfStarted() {
	s.IfStarted(func() {
		select {
		case s.chQueue <- struct{}{}:
		default:
		}
	})
}

// WakeUp wakes up the sleeper task, asking it to execute its Worker.
func (s *sleeperTask) WakeUp() {
	if s.StartStopOnce.State() == StartStopOnce_Stopped {
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
}

type sleeperTaskWorker struct {
	name string
	work func()
}

// SleeperFuncTask returns a Worker to execute the given work function.
func SleeperFuncTask(work func(), name string) Worker {
	return &sleeperTaskWorker{name: name, work: work}
}

func (w *sleeperTaskWorker) Name() string { return w.name }
func (w *sleeperTaskWorker) Work()        { w.work() }
