package store_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store"
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

	s.RunChannel.Send("whatever")
	s.RunChannel.Send("whatever")

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

	assert.NoError(t, rq.Send("first"))
	rr1 := <-rq.Receive()
	assert.NotNil(t, rr1)
}

func TestQueuedRunChannel_Send_afterClose(t *testing.T) {
	t.Parallel()

	rq := store.NewQueuedRunChannel()
	rq.Close()

	assert.Error(t, rq.Send("first"))
}
