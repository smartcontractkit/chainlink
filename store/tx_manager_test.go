package store_test

import (
	"encoding/hex"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	strpkg "github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

func TestTxManager_CreateTx(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store
	manager := store.TxManager

	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	assert.Nil(t, err)
	hash := cltest.NewHash()
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

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_MeetsMinConfirmations_BeforeThreshold(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store
	config := store.Config
	txm := store.TxManager

	sentAt := uint64(23456)
	from := cltest.GetAccountAddress(store)

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold-1))

	tx := cltest.CreateTxAndAttempt(store, from, sentAt)
	attempts, err := store.AttemptsFor(tx.ID)
	assert.Nil(t, err)
	a := attempts[0]

	confirmed, err := txm.MeetsMinConfirmations(a.Hash)
	assert.Nil(t, err)
	assert.False(t, confirmed)
	assert.Nil(t, store.One("ID", tx.ID, tx))
	attempts, err = store.AttemptsFor(tx.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(attempts))

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_MeetsMinConfirmations_AtThreshold(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store
	config := store.Config
	txm := store.TxManager

	sentAt := uint64(23456)
	from := cltest.GetAccountAddress(store)

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold))
	ethMock.Register("eth_sendRawTransaction", cltest.NewHash())

	tx := cltest.CreateTxAndAttempt(store, from, sentAt)
	attempts, err := store.AttemptsFor(tx.ID)
	assert.Nil(t, err)
	a := attempts[0]

	confirmed, err := txm.MeetsMinConfirmations(a.Hash)
	assert.Nil(t, err)
	assert.False(t, confirmed)
	assert.Nil(t, store.One("ID", tx.ID, tx))
	attempts, err = store.AttemptsFor(tx.ID)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(attempts))

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_MeetsMinConfirmations_confirmed(t *testing.T) {
	t.Parallel()

	config, configCleanup := cltest.NewConfig()
	defer configCleanup()

	sentAt := uint64(1)
	receiptAt := uint64(2)
	config.TxMinConfirmations = 2

	app, cleanup := cltest.NewApplicationWithConfigAndKeyStore(config)
	defer cleanup()
	store := app.Store
	txm := store.TxManager

	from := cltest.GetAccountAddress(store)

	tests := []struct {
		name          string
		currentHeight uint64
		want          bool
	}{
		{"less than min confs", 2, false},
		{"equal min confs", 3, true},
		{"1 greater than min confs", 4, true},
		{"2 greater than min confs", 5, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ethMock := app.MockEthClient()
			confirmationReceipt := strpkg.TxReceipt{
				Hash:        cltest.NewHash(),
				BlockNumber: cltest.BigHexInt(receiptAt),
			}
			ethMock.Register("eth_getTransactionReceipt", confirmationReceipt)
			ethMock.Register("eth_blockNumber", utils.Uint64ToHex(test.currentHeight))

			tx := cltest.CreateTxAndAttempt(store, from, sentAt)
			a := tx.TxAttempt

			actual, err := txm.MeetsMinConfirmations(a.Hash)
			assert.Nil(t, err)
			assert.Equal(t, test.want, actual)

			attempts, err := store.AttemptsFor(tx.ID)
			assert.Nil(t, err)
			assert.Equal(t, 1, len(attempts))

			ethMock.EventuallyAllCalled(t)
		})
	}
}
