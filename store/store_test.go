package store_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
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

func TestNewStore_Stop(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore()
	defer cleanup()
	want := store.RunRequest{Input: models.RunResult{JobRunID: "whatever"}}

	s.RunQueue <- want
	s.RunQueue <- want

	rr, open := <-s.RunQueue
	assert.Equal(t, want, rr)
	assert.True(t, open)

	s.Stop()

	rr, open = <-s.RunQueue
	assert.Equal(t, want, rr)
	assert.True(t, open)

	rr, open = <-s.RunQueue
	assert.Equal(t, store.RunRequest{}, rr)
	assert.False(t, open)
}
