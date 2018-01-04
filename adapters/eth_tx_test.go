package adapters_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink-go/adapters"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	strpkg "github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/store/models"
	"github.com/smartcontractkit/chainlink-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestEthTxAdapterConfirmed(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplicationWithKeyStore()
	defer app.Stop()
	store := app.Store
	config := store.Config

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", `0x0100`)
	hash := cltest.NewTxHash()
	sentAt := uint64(23456)
	confirmed := sentAt + 1
	safe := confirmed + config.EthMinConfirmations
	ethMock.Register("eth_sendRawTransaction", hash)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{Hash: hash, BlockNumber: confirmed})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(safe))

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
	attempts, _ := store.AttemptsFor(txs[0].ID)
	assert.Equal(t, 1, len(attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthTxAdapterFromPending(t *testing.T) {
	t.Parallel()
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
	assert.Nil(t, store.Save(txr))
	a, err := store.AddAttempt(txr, txr.Signable(big.NewInt(1)), sentAt)
	assert.Nil(t, err)
	adapter := adapters.EthTx{}
	input := models.RunResultPending(models.RunResultWithValue(a.Hash))

	output := adapter.Perform(input, store)

	assert.False(t, output.HasError())
	assert.True(t, output.Pending)
	assert.Nil(t, store.One("ID", txr.ID, txr))
	attempts, _ := store.AttemptsFor(txr.ID)
	assert.Equal(t, 1, len(attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthTxAdapterFromPendingBumpGas(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplicationWithKeyStore()
	defer app.Stop()
	store := app.Store
	config := store.Config

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	sentAt := uint64(23456)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold))
	ethMock.Register("eth_sendRawTransaction", cltest.NewTxHash())

	from := store.KeyStore.GetAccount().Address.String()
	txr := cltest.NewEthTx(from, sentAt)
	assert.Nil(t, store.Save(txr))
	a, err := store.AddAttempt(txr, txr.Signable(big.NewInt(1)), 1)
	assert.Nil(t, err)
	adapter := adapters.EthTx{}
	input := models.RunResultPending(models.RunResultWithValue(a.Hash))

	output := adapter.Perform(input, store)

	assert.True(t, output.Pending)
	assert.Nil(t, store.One("ID", txr.ID, txr))
	attempts, _ := store.AttemptsFor(txr.ID)
	assert.Equal(t, 2, len(attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthTxAdapterFromPendingConfirm(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplicationWithKeyStore()
	defer app.Stop()
	store := app.Store
	config := store.Config

	sentAt := uint64(23456)

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{
		Hash:        cltest.NewTxHash(),
		BlockNumber: sentAt,
	})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthMinConfirmations))

	txr := cltest.NewEthTx(cltest.NewEthAddress(), sentAt)
	assert.Nil(t, store.Save(txr))
	store.AddAttempt(txr, txr.Signable(big.NewInt(1)), sentAt)
	store.AddAttempt(txr, txr.Signable(big.NewInt(2)), sentAt+1)
	a3, _ := store.AddAttempt(txr, txr.Signable(big.NewInt(3)), sentAt+2)
	adapter := adapters.EthTx{}
	input := models.RunResultPending(models.RunResultWithValue(a3.Hash))

	assert.False(t, txr.Confirmed)

	output := adapter.Perform(input, store)

	assert.False(t, output.Pending)
	assert.False(t, output.HasError())

	assert.Nil(t, store.One("ID", txr.ID, txr))
	assert.True(t, txr.Confirmed)
	attempts, _ := store.AttemptsFor(txr.ID)
	assert.False(t, attempts[0].Confirmed)
	assert.True(t, attempts[1].Confirmed)
	assert.False(t, attempts[2].Confirmed)

	assert.True(t, ethMock.AllCalled())
}

func TestEthTxAdapterWithError(t *testing.T) {
	t.Parallel()
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
