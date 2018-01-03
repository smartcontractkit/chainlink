package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-go/adapters"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	strpkg "github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/store/models"
	"github.com/smartcontractkit/chainlink-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestEthTxAdapterConfirmed(t *testing.T) {
	app := cltest.NewApplicationWithKeyStore()
	defer app.Stop()
	store := app.Store
	config := store.Config

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", `0x0100`)
	txid := cltest.NewTxID()
	confed := uint64(23456)
	ethMock.Register("eth_sendRawTransaction", txid)
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{TXID: txid, BlockNumber: confed})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(confed+config.EthMinConfirmations))

	adapter := adapters.EthTx{
		Address:    cltest.NewEthAddress(),
		FunctionID: "12345678",
	}
	input := models.RunResultWithValue("")
	output := adapter.Perform(input, store)

	assert.False(t, output.HasError())

	from := store.KeyStore.GetAccount().Address.String()
	txs := []models.EthTx{}
	assert.Nil(t, store.Where("From", from, &txs))
	assert.Equal(t, 1, len(txs))
	assert.Equal(t, 1, len(txs[0].Attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthTxAdapterFromPending(t *testing.T) {
	app := cltest.NewApplicationWithKeyStore()
	defer app.Stop()
	store := app.Store
	config := store.Config

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	sentAt := uint64(23456)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold-1))

	from := store.KeyStore.GetAccount().Address.String()
	txr := cltest.NewEthTx(from, sentAt)
	assert.Nil(t, store.SaveTx(txr))
	adapter := adapters.EthTx{Address: cltest.NewEthAddress(), FunctionID: "12345678"}
	input := models.RunResultPending(models.RunResultWithValue(txr.TxID()))

	output := adapter.Perform(input, store)

	assert.True(t, output.Pending)
	assert.Nil(t, store.One("ID", txr.ID, txr))
	assert.Equal(t, 1, len(txr.Attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthTxAdapterFromPendingBumpGas(t *testing.T) {
	app := cltest.NewApplicationWithKeyStore()
	defer app.Stop()
	store := app.Store
	config := store.Config

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	sentAt := uint64(23456)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold))
	ethMock.Register("eth_sendRawTransaction", cltest.NewTxID())

	from := store.KeyStore.GetAccount().Address.String()
	txr := cltest.NewEthTx(from, sentAt)
	assert.Nil(t, store.SaveTx(txr))
	adapter := adapters.EthTx{Address: cltest.NewEthAddress(), FunctionID: "12345678"}
	input := models.RunResultPending(models.RunResultWithValue(txr.TxID()))

	output := adapter.Perform(input, store)

	assert.True(t, output.Pending)
	assert.Nil(t, store.One("ID", txr.ID, txr))
	assert.Equal(t, 2, len(txr.Attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthTxAdapterWithError(t *testing.T) {
	app := cltest.NewApplicationWithKeyStore()
	defer app.Stop()
	store := app.Store
	eth := app.MockEthClient()
	eth.RegisterError("eth_getTransactionCount", "Cannot connect to nodes")

	adapter := adapters.EthTx{
		Address:    "recipient",
		FunctionID: "fid",
	}
	input := models.RunResultWithValue("Hello World!")
	output := adapter.Perform(input, store)

	assert.True(t, output.HasError())
	assert.Equal(t, output.Error(), "Cannot connect to nodes")
}
