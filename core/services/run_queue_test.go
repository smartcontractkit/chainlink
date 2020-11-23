package services_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

func TestRunQueue(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)

	runExecutor := new(mocks.RunExecutor)
	runQueue := services.NewRunQueue(runExecutor)

	executeJobChannel := make(chan struct{})

	runQueue.Start()
	defer runQueue.Stop()

	runExecutor.On("Execute", mock.Anything).
		Return(nil, nil).
		Run(func(mock.Arguments) {
			executeJobChannel <- struct{}{}
		})

	runQueue.Run(models.NewID())

	g.Eventually(func() int {
		return runQueue.WorkerCount()
	}).Should(gomega.Equal(1))

	cltest.CallbackOrTimeout(t, "Execute", func() {
		<-executeJobChannel
	})

	runExecutor.AssertExpectations(t)

	g.Eventually(func() int {
		return runQueue.WorkerCount()
	}).Should(gomega.Equal(0))
}

func TestRunQueue_OneWorkerPerRun(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)

	runExecutor := new(mocks.RunExecutor)
	runQueue := services.NewRunQueue(runExecutor)

	executeJobChannel := make(chan struct{})

	runQueue.Start()
	defer runQueue.Stop()

	runExecutor.On("Execute", mock.Anything).
		Return(nil, nil).
		Run(func(mock.Arguments) {
			executeJobChannel <- struct{}{}
		})

	runQueue.Run(models.NewID())
	runQueue.Run(models.NewID())

	g.Eventually(func() int {
		return runQueue.WorkerCount()
	}).Should(gomega.Equal(2))

	cltest.CallbackOrTimeout(t, "Execute", func() {
		<-executeJobChannel
		<-executeJobChannel
	})

	runExecutor.AssertExpectations(t)

	g.Eventually(func() int {
		return runQueue.WorkerCount()
	}).Should(gomega.Equal(0))
}

func TestRunQueue_OneWorkerForSameRunTriggeredMultipleTimes(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)

	runExecutor := new(mocks.RunExecutor)
	runQueue := services.NewRunQueue(runExecutor)

	executeJobChannel := make(chan struct{})

	runQueue.Start()
	defer runQueue.Stop()

	runExecutor.On("Execute", mock.Anything).
		Return(nil, nil).
		Run(func(mock.Arguments) {
			executeJobChannel <- struct{}{}
		})

	id := models.NewID()
	runQueue.Run(id)
	runQueue.Run(id)

	g.Eventually(func() int {
		return runQueue.WorkerCount()
	}).Should(gomega.Equal(1))

	g.Consistently(func() int {
		return runQueue.WorkerCount()
	}).Should(gomega.BeNumerically("<", 2))

	cltest.CallbackOrTimeout(t, "Execute", func() {
		<-executeJobChannel
		<-executeJobChannel
	})

	runExecutor.AssertExpectations(t)

	g.Eventually(func() int {
		return runQueue.WorkerCount()
	}).Should(gomega.Equal(0))
}
