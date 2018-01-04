package store_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	strpkg "github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestEthCreateTx(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplicationWithKeyStore()
	store := app.Store
	defer app.Stop()
	manager := store.Eth

	to := "0xb70a511baC46ec6442aC6D598eaC327334e634dB"
	data := "0000abcdef"
	txid := "0x86300ee06a57eb27fbd8a6d5380783d4f8cb7210747689fe452e40f049d3de08"
	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", "0x0100") // 256
	ethMock.Register("eth_sendRawTransaction", txid)

	tx, err := manager.CreateTx(to, data)
	assert.Nil(t, err)
	assert.Equal(t, uint64(256), tx.Nonce)
	assert.Equal(t, data, tx.Data)
	assert.Equal(t, to, tx.To)

	assert.Nil(t, store.One("From", tx.From, tx))
	assert.Equal(t, uint64(256), tx.Nonce)
	assert.Equal(t, 1, len(tx.Attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthEnsureTxConfirmedBeforeThreshold(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplicationWithKeyStore()
	store := app.Store
	defer app.Stop()
	config := store.Config
	eth := store.Eth

	sentAt := uint64(23456)
	from := store.KeyStore.GetAccount().Address.String()

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold-1))

	txr := cltest.NewEthTx(from, sentAt)
	assert.Nil(t, store.SaveTx(txr))

	confirmed, err := eth.EnsureTxConfirmed(txr.TxID())
	assert.Nil(t, err)
	assert.False(t, confirmed)
	assert.Nil(t, store.One("ID", txr.ID, txr))
	assert.Equal(t, 1, len(txr.Attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthEnsureTxConfirmedAtThreshold(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplicationWithKeyStore()
	store := app.Store
	defer app.Stop()
	config := store.Config
	eth := store.Eth

	sentAt := uint64(23456)
	from := store.KeyStore.GetAccount().Address.String()

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold))
	ethMock.Register("eth_sendRawTransaction", cltest.NewTxID())

	txr := cltest.NewEthTx(from, sentAt)
	assert.Nil(t, store.SaveTx(txr))

	confirmed, err := eth.EnsureTxConfirmed(txr.TxID())
	assert.Nil(t, err)
	assert.False(t, confirmed)
	assert.Nil(t, store.One("ID", txr.ID, txr))
	assert.Equal(t, 2, len(txr.Attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthEnsureTxConfirmedWhenSafe(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplicationWithKeyStore()
	store := app.Store
	defer app.Stop()
	config := store.Config
	eth := store.Eth

	sentAt := uint64(23456)
	from := store.KeyStore.GetAccount().Address.String()

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{
		TxID:        cltest.NewTxID(),
		BlockNumber: sentAt,
	})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthMinConfirmations))

	txr := cltest.NewEthTx(from, sentAt)
	assert.Nil(t, store.SaveTx(txr))

	confirmed, err := eth.EnsureTxConfirmed(txr.TxID())
	assert.Nil(t, err)
	assert.True(t, confirmed)
	assert.Nil(t, store.One("ID", txr.ID, txr))
	assert.Equal(t, 1, len(txr.Attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthEnsureTxConfirmedWhenWithConfsButNotSafe(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplicationWithKeyStore()
	store := app.Store
	defer app.Stop()
	config := store.Config
	eth := store.Eth

	sentAt := uint64(23456)
	from := store.KeyStore.GetAccount().Address.String()

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{
		TxID:        cltest.NewTxID(),
		BlockNumber: sentAt,
	})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthMinConfirmations-1))

	txr := cltest.NewEthTx(from, sentAt)
	assert.Nil(t, store.SaveTx(txr))

	confirmed, err := eth.EnsureTxConfirmed(txr.TxID())
	assert.Nil(t, err)
	assert.False(t, confirmed)
	assert.Nil(t, store.One("ID", txr.ID, txr))
	assert.Equal(t, 1, len(txr.Attempts))

	assert.True(t, ethMock.AllCalled())
}
