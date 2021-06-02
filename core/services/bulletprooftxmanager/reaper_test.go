package bulletprooftxmanager_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager/mocks"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
)

func TestReaper_ReapEthTxes(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

	_, from := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	var nonce int64 = 0
	oneDayAgo := time.Now().Add(-24 * time.Hour)

	t.Run("with nothing in the database, doesn't error", func(t *testing.T) {
		config := new(mocks.ReaperConfig)
		config.On("EthFinalityDepth").Return(uint(10))
		config.On("EthTxReaperThreshold").Return(1 * time.Hour)
		config.On("EthTxReaperInterval").Return(1 * time.Hour)

		r := bulletprooftxmanager.NewReaper(store.DB, config)

		err := r.ReapEthTxes(42)
		assert.NoError(t, err)
	})

	// Confirmed in block number 5
	cltest.MustInsertConfirmedEthTxWithReceipt(t, store, from, nonce, 5)
	nonce++

	t.Run("skips if threshold=0", func(t *testing.T) {
		config := new(mocks.ReaperConfig)
		config.On("EthFinalityDepth").Return(uint(10))
		config.On("EthTxReaperThreshold").Return(0 * time.Second)
		config.On("EthTxReaperInterval").Return(1 * time.Hour)

		r := bulletprooftxmanager.NewReaper(store.DB, config)

		err := r.ReapEthTxes(42)
		assert.NoError(t, err)

		cltest.AssertCount(t, store, models.EthTx{}, 1)
	})

	t.Run("deletes confirmed eth_txes that exceed the age threshold with at least ETH_FINALITY_DEPTH blocks above their receipt", func(t *testing.T) {
		config := new(mocks.ReaperConfig)
		config.On("EthFinalityDepth").Return(uint(10))
		config.On("EthTxReaperThreshold").Return(1 * time.Hour)
		config.On("EthTxReaperInterval").Return(1 * time.Hour)

		r := bulletprooftxmanager.NewReaper(store.DB, config)

		err := r.ReapEthTxes(42)
		assert.NoError(t, err)
		// Didn't delete because eth_tx was not old enough
		cltest.AssertCount(t, store, models.EthTx{}, 1)

		store.DB.Exec(`UPDATE eth_txes SET created_at=?`, oneDayAgo)

		err = r.ReapEthTxes(12)
		assert.NoError(t, err)
		// Didn't delete because eth_tx although old enough, was still within ETH_FINALITY_DEPTH of the current head
		cltest.AssertCount(t, store, models.EthTx{}, 1)

		err = r.ReapEthTxes(42)
		assert.NoError(t, err)
		// Now it deleted because the eth_tx was past ETH_FINALITY_DEPTH
		cltest.AssertCount(t, store, models.EthTx{}, 0)
	})

	cltest.MustInsertFatalErrorEthTx(t, store, from)

	t.Run("deletes errored eth_txes that exceed the age threshold", func(t *testing.T) {
		config := new(mocks.ReaperConfig)
		config.On("EthFinalityDepth").Return(uint(10))
		config.On("EthTxReaperThreshold").Return(1 * time.Hour)
		config.On("EthTxReaperInterval").Return(1 * time.Hour)

		r := bulletprooftxmanager.NewReaper(store.DB, config)

		err := r.ReapEthTxes(42)
		assert.NoError(t, err)
		// Didn't delete because eth_tx was not old enough
		cltest.AssertCount(t, store, models.EthTx{}, 1)

		store.DB.Exec(`UPDATE eth_txes SET created_at=?`, oneDayAgo)

		err = r.ReapEthTxes(42)
		assert.NoError(t, err)
		// Deleted because it is old enough now
		cltest.AssertCount(t, store, models.EthTx{}, 0)
	})
}
