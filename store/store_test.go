package store_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
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

func TestStore_Close(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore()
	defer cleanup()
	want := models.RunResult{JobRunID: "whatever"}

	s.RunChannel.Send(want, nil)
	s.RunChannel.Send(want, nil)

	rr, open := <-s.RunChannel.Receive()
	assert.Equal(t, want, rr.Input)
	assert.True(t, open)

	rr, open = <-s.RunChannel.Receive()
	assert.Equal(t, want, rr.Input)
	assert.True(t, open)

	assert.NoError(t, s.Close())

	rr, open = <-s.RunChannel.Receive()
	assert.Equal(t, store.RunRequest{}, rr)
	assert.False(t, open)
}

func TestRunChannel_Send(t *testing.T) {
	t.Parallel()

	rq := store.NewRunChannel()
	input1 := models.RunResult{JobRunID: "first"}
	ibn1 := cltest.IndexableBlockNumber(17)

	assert.NoError(t, rq.Send(input1, ibn1))
	rr1 := <-rq.Receive()
	assert.Equal(t, input1, rr1.Input)
	assert.Equal(t, ibn1, rr1.BlockNumber)
}

func TestRunChannel_Send_afterClose(t *testing.T) {
	t.Parallel()

	rq := store.NewRunChannel()
	input1 := models.RunResult{JobRunID: "first"}
	ibn1 := cltest.IndexableBlockNumber(17)

	rq.Close()

	assert.Error(t, rq.Send(input1, ibn1))
}
