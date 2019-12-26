package services_test

import (
	"math/big"
	"sync"
	"testing"

	ethpkg "chainlink/core/eth"
	"chainlink/core/internal/cltest"
	"chainlink/core/internal/mocks"
	"chainlink/core/services"
	"chainlink/core/store/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestJobSubscriber_OnNewHead(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runManager := new(mocks.RunManager)
	jobSubscriber := services.NewJobSubscriber(store, runManager)

	wg := sync.WaitGroup{}
	wg.Add(1)
	blockResumeAllConfirmingChannel := make(chan struct{})

	runManager.On("ResumeAllConfirming", big.NewInt(1337)).
		Return(nil).
		Once().
		Run(func(mock.Arguments) {
			wg.Done()
			blockResumeAllConfirmingChannel <- struct{}{}
		})
	runManager.On("ResumeAllConfirming", big.NewInt(1338)).
		Return(nil).
		Once().
		Run(func(mock.Arguments) {
			blockResumeAllConfirmingChannel <- struct{}{}
		})
	jobSubscriber.OnNewHead(cltest.Head(1337))
	jobSubscriber.OnNewHead(cltest.Head(1338))

	wg.Wait()

	// JobSubscriber on new head channel is now blocked, this head should get dropped
	jobSubscriber.OnNewHead(cltest.Head(1339))

	// Unblock the channel
	cltest.CallbackOrTimeout(t, "ResumeAllConfirming", func() {
		<-blockResumeAllConfirmingChannel
		<-blockResumeAllConfirmingChannel
	})

	// Make sure after dropping a head (because of congestion) that it resumes again
	runManager.On("ResumeAllConfirming", big.NewInt(1340)).
		Return(nil).
		Once().
		Run(func(mock.Arguments) {
			blockResumeAllConfirmingChannel <- struct{}{}
		})
	jobSubscriber.OnNewHead(cltest.Head(1340))

	cltest.CallbackOrTimeout(t, "ResumeAllConfirming #2", func() {
		<-blockResumeAllConfirmingChannel
	})

	runManager.AssertExpectations(t)
}

func TestJobSubscriber_AddJob_RemoveJob(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	cltest.MockEthOnStore(t, store)

	runManager := new(mocks.RunManager)
	jobSubscriber := services.NewJobSubscriber(store, runManager)

	jobSpec := cltest.NewJobWithLogInitiator()
	err := jobSubscriber.AddJob(jobSpec, cltest.Head(321))
	require.NoError(t, err)

	assert.Len(t, jobSubscriber.Jobs(), 1)

	err = jobSubscriber.RemoveJob(jobSpec.ID)
	require.NoError(t, err)

	assert.Len(t, jobSubscriber.Jobs(), 0)

	runManager.AssertExpectations(t)

}

func TestJobSubscriber_AddJob_NotLogInitiatedError(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runManager := new(mocks.RunManager)
	jobSubscriber := services.NewJobSubscriber(store, runManager)

	job := models.JobSpec{}
	err := jobSubscriber.AddJob(job, cltest.Head(1))
	require.NoError(t, err)
}

func TestJobSubscriber_RemoveJob_NotFoundError(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runManager := new(mocks.RunManager)
	jobSubscriber := services.NewJobSubscriber(store, runManager)

	err := jobSubscriber.RemoveJob(models.NewID())
	require.Error(t, err)
}

func TestJobSubscriber_Connect_Disconnect(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runManager := new(mocks.RunManager)
	jobSubscriber := services.NewJobSubscriber(store, runManager)

	eth := cltest.MockEthOnStore(t, store)
	eth.Register("eth_getLogs", []ethpkg.Log{})
	eth.Register("eth_getLogs", []ethpkg.Log{})

	jobSpec1 := cltest.NewJobWithLogInitiator()
	jobSpec2 := cltest.NewJobWithLogInitiator()
	assert.Nil(t, store.CreateJob(&jobSpec1))
	assert.Nil(t, store.CreateJob(&jobSpec2))
	eth.RegisterSubscription("logs")
	eth.RegisterSubscription("logs")

	assert.Nil(t, jobSubscriber.Connect(cltest.Head(491)))
	eth.EventuallyAllCalled(t)

	assert.Len(t, jobSubscriber.Jobs(), 2)

	jobSubscriber.Disconnect()

	assert.Len(t, jobSubscriber.Jobs(), 0)
}
