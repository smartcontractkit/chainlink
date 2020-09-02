package services

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
	worker      Worker
	chQueue     chan struct{}
	chQueueDone chan struct{}
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
		worker:      worker,
		chQueue:     make(chan struct{}, 1),
		chQueueDone: make(chan struct{}),
	}

	go s.workerLoop()

	return s
}

// Stop stops the SleeperTask.  It never returns an error.  Its error return
// exists so as to satisfy other interfaces.
func (s *sleeperTask) Stop() error {
	close(s.chQueue)
	<-s.chQueueDone
	return nil
}

// WakeUp wakes up the sleeper task, asking it to execute its Worker.
func (s *sleeperTask) WakeUp() {
	select {
	case s.chQueue <- struct{}{}:
	default:
	}
}

func (s *sleeperTask) workerLoop() {
	defer close(s.chQueueDone)

	for range s.chQueue {
		s.worker.Work()
	}

	if len(s.chQueue) > 0 {
		s.worker.Work()
	}
}

type SleeperTaskFuncWorker func()

func (fn SleeperTaskFuncWorker) Work() { fn() }
