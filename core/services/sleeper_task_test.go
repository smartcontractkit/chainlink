package services

import (
	"sync"
	"testing"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

type testWorker struct {
	output chan struct{}
}

func (w *testWorker) Work() {
	w.output <- struct{}{}
}

type blockingWorker struct {
	output  chan struct{}
	mutex   sync.Mutex
	started bool
}

func (w *blockingWorker) Work() {
	w.mutex.Lock()
	w.started = true
	w.mutex.Unlock()
}

func TestSleeperTask(t *testing.T) {
	worker := testWorker{output: make(chan struct{})}
	sleeper := NewSleeperTask(&worker)

	sleeper.Start()
	sleeper.WakeUp()

	gomega.NewGomegaWithT(t).Eventually(worker.output).Should(gomega.Receive(&struct{}{}))

	sleeper.Stop()
}

func TestSleeperTask_WakeupBeforeStarted(t *testing.T) {
	worker := testWorker{output: make(chan struct{})}
	sleeper := NewSleeperTask(&worker)

	sleeper.WakeUp()
	sleeper.Start()

	gomega.NewGomegaWithT(t).Eventually(worker.output).Should(gomega.Receive(&struct{}{}))

	sleeper.Stop()
}

func TestSleeperTask_Restart(t *testing.T) {
	worker := testWorker{output: make(chan struct{})}
	sleeper := NewSleeperTask(&worker)

	sleeper.Start()
	sleeper.WakeUp()

	gomega.NewGomegaWithT(t).Eventually(worker.output).Should(gomega.Receive(&struct{}{}))

	sleeper.Stop()

	sleeper.Start()
	sleeper.WakeUp()

	gomega.NewGomegaWithT(t).Eventually(worker.output).Should(gomega.Receive(&struct{}{}))

	sleeper.Stop()
}

func TestSleeperTask_SenderNotBlockedWhileWorking(t *testing.T) {
	worker := testWorker{output: make(chan struct{})}
	sleeper := NewSleeperTask(&worker)

	sleeper.Start()

	sleeper.WakeUp()
	sleeper.WakeUp()

	gomega.NewGomegaWithT(t).Eventually(worker.output).Should(gomega.Receive(&struct{}{}))

	sleeper.Stop()
}

func TestSleeperTask_StopWaitsUntilWorkFinishes(t *testing.T) {
	worker := blockingWorker{output: make(chan struct{})}
	sleeper := NewSleeperTask(&worker)

	// Block worker from setting 'started=true'. It must be acquired
	// before sleeper.Start() to avoid a race condition
	worker.mutex.Lock()
	sleeper.Start()
	assert.Equal(t, false, worker.started)

	// Increments the wait group if the channel is not full
	sleeper.WakeUp()

	beforeStop := make(chan struct{}, 1)
	afterStop := make(chan struct{}, 1)
	go func() {
		beforeStop <- struct{}{}
		sleeper.Stop()
		afterStop <- struct{}{}
	}()

	<-beforeStop
	assert.Equal(t, false, worker.started)

	// Release the worker to do it's work which will result in the wait group counter being decremented
	worker.mutex.Unlock()
	// Ensure that Stop() has returned
	<-afterStop
	assert.Equal(t, true, worker.started)
}

func TestSleeperTask_StopWithoutStartNonBlocking(t *testing.T) {
	worker := testWorker{output: make(chan struct{})}
	sleeper := NewSleeperTask(&worker)

	sleeper.Start()
	sleeper.WakeUp()
	gomega.NewGomegaWithT(t).Eventually(worker.output).Should(gomega.Receive(&struct{}{}))

	sleeper.Stop()
	sleeper.Stop()
}

type slowWorker struct {
	mutex  sync.Mutex
	output chan struct{}
}

func (t *slowWorker) Work() {
	t.output <- struct{}{}
	t.mutex.Lock()
	t.mutex.Unlock()
}

func TestSleeperTask_WakeWhileWorkingRepeatsWork(t *testing.T) {
	worker := slowWorker{output: make(chan struct{})}
	sleeper := NewSleeperTask(&worker)

	sleeper.Start()

	// Lock the worker's mutex so it's blocked *after* sending to the output
	// channel, this guarantees that the worker blocks till we unlock the mutex
	worker.mutex.Lock()
	sleeper.WakeUp()
	// Make sure an item is received in the channel so we know the worker is blocking
	gomega.NewGomegaWithT(t).Eventually(worker.output).Should(gomega.Receive(&struct{}{}))

	// Wake up the sleeper
	sleeper.WakeUp()
	// Now release the worker
	worker.mutex.Unlock()
	gomega.NewGomegaWithT(t).Eventually(worker.output).Should(gomega.Receive(&struct{}{}))

	sleeper.Stop()
}
