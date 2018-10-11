package store_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_Start(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	ethMock := app.MockEthClient()
	store := app.Store

	ethMock.Register("eth_getTransactionCount", `0x2D0`)
	assert.Nil(t, store.Start())
	ethMock.EventuallyAllCalled(t)
}

func TestStore_Start_CleansupPrematureShutdown(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", `0x2D0`)

	s := app.Store
	job, init := cltest.NewJobWithWebInitiator()
	require.NoError(t, s.SaveJob(&job))

	jr := job.NewRun(init)
	jr.Status = models.RunStatusInProgress
	require.NoError(t, s.Save(&jr))

	require.NoError(t, s.Start())
	var cleanedJr models.JobRun
	require.NoError(t, s.One("ID", jr.ID, &cleanedJr))

	assert.Equal(t, jr.ID, cleanedJr.ID)
	assert.Equal(t, string(models.RunStatusUnstarted), string(cleanedJr.Status))

	require.NoError(t, app.JobRunner.Start())
	cltest.WaitForJobRunToComplete(t, s, cleanedJr)
}

func TestStore_Close(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore()
	defer cleanup()

	s.RunChannel.Send("whatever", nil)
	s.RunChannel.Send("whatever", nil)

	rr, open := <-s.RunChannel.Receive()
	assert.True(t, open)

	rr, open = <-s.RunChannel.Receive()
	assert.True(t, open)

	assert.NoError(t, s.Close())

	rr, open = <-s.RunChannel.Receive()
	assert.Equal(t, store.RunRequest{}, rr)
	assert.False(t, open)
}

func TestQueuedRunChannel_Send(t *testing.T) {
	t.Parallel()

	rq := store.NewQueuedRunChannel()
	ibn1 := cltest.IndexableBlockNumber(17)

	assert.NoError(t, rq.Send("first", ibn1))
	rr1 := <-rq.Receive()
	assert.Equal(t, ibn1, rr1.BlockNumber)
}

func TestQueuedRunChannel_Send_afterClose(t *testing.T) {
	t.Parallel()

	rq := store.NewQueuedRunChannel()
	ibn1 := cltest.IndexableBlockNumber(17)

	rq.Close()

	assert.Error(t, rq.Send("first", ibn1))
}
