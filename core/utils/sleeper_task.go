package utils

// SleeperTask represents a task that waits in the background to process some work.
type SleeperTask interface {
	Stop() error
	WakeUp(WorkFunc)
}

// Worker is a simple interface that represents some work to do repeatedly
type WorkFunc func()

type sleeperTask struct {
	chQueue     chan WorkFunc
	chQueueDone chan struct{}
}

// NewSleeperTask returns a SleeperTask.
//
// SleeperTask is guaranteed to call the given WorkFunc at least once for every
// WakeUp call.
// If the WorkFunc is busy when WakeUp is called, the Worker will be called again
// immediately after it is finished. For this reason you should take care to
// make sure that WorkFunc is idempotent.
// WakeUp does not block.
//
func NewSleeperTask() SleeperTask {
	s := &sleeperTask{
		chQueue:     make(chan WorkFunc, 1),
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
func (s *sleeperTask) WakeUp(fn WorkFunc) {
	select {
	case s.chQueue <- fn:
	default:
	}
}

func (s *sleeperTask) workerLoop() {
	defer close(s.chQueueDone)

	for fn := range s.chQueue {
		fn()
	}

	if len(s.chQueue) > 0 {
		fn := <-s.chQueue
		fn()
	}
}

type SleeperTaskReceiver interface {
}

type InterfaceSleeperTask struct {
	*sleeperTask
	receiver SleeperTaskReceiver
}

func NewInterfaceSleeperTask(r SleeperTaskReceiver) SleeperTask {
	return &InterfaceSleeperTask{&sleeperTask{}, r}
}
