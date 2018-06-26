package store_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
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
