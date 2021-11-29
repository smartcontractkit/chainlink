package utils_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
)

type countingWorker struct {
	numJobsPerformed atomic.Int32
	delay            time.Duration
}

func (t *countingWorker) Name() string {
	return "CountingWorker"
}

func (t *countingWorker) Work() {
	if t.delay != 0 {
		time.Sleep(t.delay)
	}
	// Without an atomic, the race detector fails
	t.numJobsPerformed.Inc()
}

func (t *countingWorker) getNumJobsPerformed() int {
	return int(t.numJobsPerformed.Load())
}

func TestSleeperTask_WakeupAfterStopPanics(t *testing.T) {
	t.Parallel()

	worker := &countingWorker{}
	sleeper := utils.NewSleeperTask(worker)

	require.NoError(t, sleeper.Stop())

	require.Panics(t, func() {
		sleeper.WakeUp()
	})
	gomega.NewWithT(t).Eventually(worker.getNumJobsPerformed).Should(gomega.Equal(0))
}

func TestSleeperTask_CallingStopTwiceFails(t *testing.T) {
	t.Parallel()

	worker := &countingWorker{}
	sleeper := utils.NewSleeperTask(worker)
	require.NoError(t, sleeper.Stop())
	require.Error(t, sleeper.Stop())
}

func TestSleeperTask_WakeupPerformsWork(t *testing.T) {
	t.Parallel()

	worker := &countingWorker{}
	sleeper := utils.NewSleeperTask(worker)

	sleeper.WakeUp()
	gomega.NewWithT(t).Eventually(worker.getNumJobsPerformed).Should(gomega.Equal(1))
	require.NoError(t, sleeper.Stop())
}

type controllableWorker struct {
	countingWorker
	awaitWorkStarted chan struct{}
	allowResumeWork  chan struct{}
	ignoreSignals    bool
}

func (w *controllableWorker) Work() {
	if !w.ignoreSignals {
		w.awaitWorkStarted <- struct{}{}
		<-w.allowResumeWork
	}
	w.countingWorker.Work()
}

func TestSleeperTask_WakeupEnqueuesMaxTwice(t *testing.T) {
	t.Parallel()

	worker := &controllableWorker{awaitWorkStarted: make(chan struct{}), allowResumeWork: make(chan struct{})}
	sleeper := utils.NewSleeperTask(worker)

	sleeper.WakeUp()
	<-worker.awaitWorkStarted
	sleeper.WakeUp()
	sleeper.WakeUp()
	sleeper.WakeUp()
	sleeper.WakeUp()
	sleeper.WakeUp()
	worker.ignoreSignals = true
	worker.allowResumeWork <- struct{}{}

	gomega.NewWithT(t).Eventually(worker.getNumJobsPerformed).Should(gomega.Equal(2))
	gomega.NewWithT(t).Consistently(worker.getNumJobsPerformed).Should(gomega.BeNumerically("<", 3))
	require.NoError(t, sleeper.Stop())
}

func TestSleeperTask_StopWaitsUntilWorkFinishes(t *testing.T) {
	t.Parallel()

	worker := &controllableWorker{awaitWorkStarted: make(chan struct{}), allowResumeWork: make(chan struct{})}
	sleeper := utils.NewSleeperTask(worker)

	sleeper.WakeUp()
	<-worker.awaitWorkStarted
	require.Equal(t, 0, worker.getNumJobsPerformed())
	worker.allowResumeWork <- struct{}{}

	require.NoError(t, sleeper.Stop())
	require.Equal(t, worker.getNumJobsPerformed(), 1)
}
