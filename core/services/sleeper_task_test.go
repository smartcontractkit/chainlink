package services_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"chainlink/core/services"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testWorker struct {
	output chan struct{}
}

func (t *testWorker) Work() {
	t.output <- struct{}{}
}

type longRunningWorker struct {
	output   chan struct{}
	finished bool
}

func (w *longRunningWorker) Work() {
	w.output <- struct{}{}
	time.Sleep(100 * time.Millisecond)
	w.finished = true
}

func TestSleeperTask(t *testing.T) {
	worker := testWorker{output: make(chan struct{})}
	sleeper := services.NewSleeperTask(&worker)

	sleeper.Start()
	sleeper.WakeUp()

	gomega.NewGomegaWithT(t).Eventually(worker.output).Should(gomega.Receive(&struct{}{}))

	sleeper.Stop()
}

func TestSleeperTask_WakeupBeforeStarted(t *testing.T) {
	worker := testWorker{output: make(chan struct{})}
	sleeper := services.NewSleeperTask(&worker)

	sleeper.WakeUp()
	sleeper.Start()

	gomega.NewGomegaWithT(t).Eventually(worker.output).Should(gomega.Receive(&struct{}{}))

	sleeper.Stop()
}

func TestSleeperTask_WakeupAfterStarted(t *testing.T) {
	worker := testWorker{output: make(chan struct{})}
	sleeper := services.NewSleeperTask(&worker)

	sleeper.Start()
	sleeper.WakeUp()

	gomega.NewGomegaWithT(t).Eventually(worker.output).Should(gomega.Receive(&struct{}{}))

	sleeper.Stop()
}

func TestSleeperTask_WakeupAfterStoppedPanics(t *testing.T) {
	worker := testWorker{output: make(chan struct{})}
	sleeper := services.NewSleeperTask(&worker)

	sleeper.Start()
	sleeper.Stop()

	assert.Panics(t, sleeper.WakeUp, "expected sleeper.WakeUp() to panic on stopped sleeper, but it didn't")
}

func TestSleeperTask_WakeupNotBlockedWhileWorking(t *testing.T) {
	worker := testWorker{output: make(chan struct{})}
	sleeper := services.NewSleeperTask(&worker)

	sleeper.Start()

	sleeper.WakeUp()
	sleeper.WakeUp()

	gomega.NewGomegaWithT(t).Eventually(worker.output).Should(gomega.Receive(&struct{}{}))

	sleeper.Stop()
}

func TestSleeperTask_StartAfterStoppedReturnsError(t *testing.T) {
	worker := testWorker{output: make(chan struct{})}
	sleeper := services.NewSleeperTask(&worker)

	sleeper.Start()
	sleeper.WakeUp()

	gomega.NewGomegaWithT(t).Eventually(worker.output).Should(gomega.Receive(&struct{}{}))

	sleeper.Stop()

	assert.Error(t, sleeper.Stop())
}

func TestSleeperTask_StopWaitsUntilWorkFinishes(t *testing.T) {
	worker := longRunningWorker{output: make(chan struct{})}
	sleeper := services.NewSleeperTask(&worker)

	sleeper.Start()
	sleeper.WakeUp()

	<-worker.output
	sleeper.Stop()

	assert.Equal(t, true, worker.finished)
}

func TestSleeperTask_StopWithoutStartNonBlocking(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	worker := testWorker{output: make(chan struct{})}
	sleeper := services.NewSleeperTask(&worker)

	sleeper.Start()
	sleeper.WakeUp()
	g.Eventually(worker.output).Should(gomega.Receive(&struct{}{}))

	err := sleeper.Stop()
	require.NoError(t, err)

	g.Eventually(sleeper.Stop).Should(gomega.Equal(errors.New("sleeper task is already stopped")))
}

func TestSleeperTask_CallingStartTwiceErrors(t *testing.T) {
	worker := testWorker{output: make(chan struct{})}
	sleeper := services.NewSleeperTask(&worker)
	defer sleeper.Stop()

	require.NoError(t, sleeper.Start())
	require.Error(t, sleeper.Start())
}

func TestSleeperTask_CallingStopTwiceErrors(t *testing.T) {
	worker := testWorker{output: make(chan struct{})}
	sleeper := services.NewSleeperTask(&worker)
	sleeper.Start()
	require.NoError(t, sleeper.Stop())
	require.Error(t, sleeper.Stop())
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
	sleeper := services.NewSleeperTask(&worker)

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
