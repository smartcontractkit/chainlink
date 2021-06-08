package services_test

import (
	"context"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestJobSubscriber_OnNewLongestChain(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runManager := new(mocks.RunManager)
	jobSubscriber := services.NewJobSubscriber(store, runManager)
	defer jobSubscriber.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)
	resumeJobChannel := make(chan struct{})

	runManager.On("ResumeAllPendingNextBlock", big.NewInt(1337)).
		Return(nil).
		Once().
		Run(func(mock.Arguments) {
			wg.Done()
			resumeJobChannel <- struct{}{}
		})
	runManager.On("ResumeAllPendingNextBlock", big.NewInt(1339)).
		Return(nil).
		Once().
		Run(func(mock.Arguments) {
			resumeJobChannel <- struct{}{}
		})
	jobSubscriber.OnNewLongestChain(context.TODO(), *cltest.Head(1337))

	// Make sure ResumeAllPendingNextBlock is reached before sending the next head
	wg.Wait()

	// This head should get dropped
	jobSubscriber.OnNewLongestChain(context.TODO(), *cltest.Head(1338))

	// This head should get processed
	jobSubscriber.OnNewLongestChain(context.TODO(), *cltest.Head(1339))

	// Unblock the channel
	cltest.CallbackOrTimeout(t, "ResumeAllPendingNextBlock", func() {
		<-resumeJobChannel
		<-resumeJobChannel
	})

	// Make sure after dropping a head (because of congestion) that it resumes again
	runManager.On("ResumeAllPendingNextBlock", big.NewInt(1340)).
		Return(nil).
		Once().
		Run(func(mock.Arguments) {
			resumeJobChannel <- struct{}{}
		})
	jobSubscriber.OnNewLongestChain(context.TODO(), *cltest.Head(1340))

	cltest.CallbackOrTimeout(t, "ResumeAllPendingNextBlock #2", func() {
		<-resumeJobChannel
	})

	runManager.AssertExpectations(t)
}

func TestJobSubscriber_AddJob_RemoveJob(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ethClient, _, assertMocksCalled := cltest.NewEthMocks(t)
	defer assertMocksCalled()
	store.EthClient = ethClient
	b := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(2),
	})
	ethClient.On("BlockByNumber", mock.Anything, mock.Anything).Maybe().Return(b, nil)
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(cltest.EmptyMockSubscription(), nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Maybe().Return([]types.Log{}, nil)

	runManager := new(mocks.RunManager)
	jobSubscriber := services.NewJobSubscriber(store, runManager)
	defer jobSubscriber.Close()

	jobSpec := cltest.NewJobWithLogInitiator()
	err := jobSubscriber.AddJob(jobSpec, cltest.Head(321))
	require.NoError(t, err)

	// Re-adding the same jobID should be idempotent
	// and NOT create a new subscription and overwrite
	jobSpec2 := cltest.NewJobWithLogInitiator()
	jobSpec2.ID = jobSpec.ID
	jobSpec2.Name = "should not overwrite"
	err = jobSubscriber.AddJob(jobSpec, cltest.Head(321))
	require.NoError(t, err)
	jbs := jobSubscriber.Jobs()
	require.Equal(t, 1, len(jbs))
	require.Equal(t, "", jbs[0].Name)

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
	defer jobSubscriber.Close()

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
	defer jobSubscriber.Close()

	err := jobSubscriber.RemoveJob(models.NewJobID())
	require.Error(t, err)
}

func TestJobSubscriber_Connect_Disconnect(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runManager := new(mocks.RunManager)
	jobSubscriber := services.NewJobSubscriber(store, runManager)

	ethClient := new(mocks.Client)
	defer ethClient.AssertExpectations(t)
	store.EthClient = ethClient
	b := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(500),
	})
	ethClient.On("BlockByNumber", mock.Anything, mock.Anything).Maybe().Return(b, nil)
	ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(cltest.EmptyMockSubscription(), nil)
	ethClient.On("FilterLogs", mock.Anything, mock.Anything).Maybe().Return([]types.Log{}, nil)

	jobSpec1 := cltest.NewJobWithLogInitiator()
	jobSpec2 := cltest.NewJobWithLogInitiator()
	require.Nil(t, store.CreateJob(&jobSpec1))
	require.Nil(t, store.CreateJob(&jobSpec2))

	require.Nil(t, jobSubscriber.Connect(cltest.Head(491)))

	assert.Len(t, jobSubscriber.Jobs(), 2)

	jobSubscriber.Close()

	assert.Len(t, jobSubscriber.Jobs(), 0)
}
