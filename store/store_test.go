package store_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/stretchr/testify/assert"
)

func TestNewStore_Start(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	ethMock := app.MockEthClient()
	store := app.Store

	ethMock.Register("eth_getTransactionCount", `0x2D0`)
	assert.Nil(t, store.Start())
	ethMock.EventuallyAllCalled(t)
}

func TestRunManager_WorkerChannelFor(t *testing.T) {
	t.Parallel()
	rm := store.NewRunManager()

	chan1 := rm.WorkerChannelFor("foo")
	chan2 := rm.WorkerChannelFor("bar")
	chan3 := rm.WorkerChannelFor("foo")

	assert.NotEqual(t, chan1, chan2)
	assert.Equal(t, chan1, chan3)
	assert.NotEqual(t, chan2, chan3)
}
