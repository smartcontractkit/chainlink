package store_test

import (
	"encoding/hex"
	"testing"

	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	strpkg "github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/store/models"
	"github.com/smartcontractkit/chainlink-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestEthCreateTx(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store
	manager := store.Eth

	to := cltest.NewEthAddress()
	data, err := hex.DecodeString("0000abcdef")
	assert.Nil(t, err)
	hash := cltest.NewTxHash()
	sentAt := uint64(23456)
	nonce := uint64(256)
	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce))
	ethMock.Register("eth_sendRawTransaction", hash)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))

	a, err := manager.CreateTx(to, data)
	assert.Nil(t, err)
	tx := models.Tx{}
	assert.Nil(t, store.One("ID", a.TxID, &tx))
	assert.Nil(t, err)
	assert.Equal(t, nonce, tx.Nonce)
	assert.Equal(t, data, tx.Data)
	assert.Equal(t, to, tx.To)

	assert.Nil(t, store.One("From", tx.From, &tx))
	assert.Equal(t, nonce, tx.Nonce)
	attempts, err := store.AttemptsFor(tx.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthEnsureTxConfirmedBeforeThreshold(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store
	config := store.Config
	eth := store.Eth

	sentAt := uint64(23456)
	from := store.KeyStore.GetAccount().Address

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold-1))

	tx := cltest.CreateTxAndAttempt(store, from, sentAt)
	attempts, err := store.AttemptsFor(tx.ID)
	assert.Nil(t, err)
	a := attempts[0]

	confirmed, err := eth.EnsureTxConfirmed(a.Hash)
	assert.Nil(t, err)
	assert.False(t, confirmed)
	assert.Nil(t, store.One("ID", tx.ID, tx))
	attempts, err = store.AttemptsFor(tx.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthEnsureTxConfirmedAtThreshold(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store
	config := store.Config
	eth := store.Eth

	sentAt := uint64(23456)
	from := store.KeyStore.GetAccount().Address

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold))
	ethMock.Register("eth_sendRawTransaction", cltest.NewTxHash())

	tx := cltest.CreateTxAndAttempt(store, from, sentAt)
	attempts, err := store.AttemptsFor(tx.ID)
	assert.Nil(t, err)
	a := attempts[0]

	confirmed, err := eth.EnsureTxConfirmed(a.Hash)
	assert.Nil(t, err)
	assert.False(t, confirmed)
	assert.Nil(t, store.One("ID", tx.ID, tx))
	attempts, err = store.AttemptsFor(tx.ID)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthEnsureTxConfirmedWhenSafe(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store
	config := store.Config
	eth := store.Eth

	sentAt := uint64(23456)
	from := store.KeyStore.GetAccount().Address

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{
		Hash:        cltest.NewTxHash(),
		BlockNumber: sentAt,
	})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthMinConfirmations))

	tx := cltest.CreateTxAndAttempt(store, from, sentAt)
	a := tx.TxAttempt

	confirmed, err := eth.EnsureTxConfirmed(a.Hash)
	assert.Nil(t, err)
	assert.True(t, confirmed)
	assert.Nil(t, store.One("ID", tx.ID, tx))
	attempts, err := store.AttemptsFor(tx.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthEnsureTxConfirmedWhenWithConfsButNotSafe(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store
	config := store.Config
	eth := store.Eth

	sentAt := uint64(23456)
	from := store.KeyStore.GetAccount().Address

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{
		Hash:        cltest.NewTxHash(),
		BlockNumber: sentAt,
	})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthMinConfirmations-1))

	tx := cltest.CreateTxAndAttempt(store, from, sentAt)
	a := tx.TxAttempt

	confirmed, err := eth.EnsureTxConfirmed(a.Hash)
	assert.Nil(t, err)
	assert.False(t, confirmed)
	assert.Nil(t, store.One("ID", tx.ID, tx))
	attempts, err := store.AttemptsFor(tx.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(attempts))

	assert.True(t, ethMock.AllCalled())
}
