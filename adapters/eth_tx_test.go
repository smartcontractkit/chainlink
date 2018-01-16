package adapters_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	strpkg "github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

func TestEthTxAdapterConfirmed(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
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
	receipt := strpkg.TxReceipt{Hash: hash, BlockNumber: confirmed}
	ethMock.Register("eth_getTransactionReceipt", receipt)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(safe))

	adapter := adapters.EthTx{
		Address:    cltest.NewEthAddress(),
		FunctionID: models.HexToFunctionID("b3f98adc"),
	}
	input := models.RunResultWithValue("")
	output := adapter.Perform(input, store)

	assert.False(t, output.HasError())

	from := store.KeyStore.GetAccount().Address
	txs := []models.Tx{}
	assert.Nil(t, store.Where("From", from, &txs))
	assert.Equal(t, 1, len(txs))
	attempts, _ := store.AttemptsFor(txs[0].ID)
	assert.Equal(t, 1, len(attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthTxAdapterFromPending(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store
	config := store.Config

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	sentAt := uint64(23456)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold-1))

	from := store.KeyStore.GetAccount().Address
	tx := cltest.NewTx(from, sentAt)
	assert.Nil(t, store.Save(tx))
	a, err := store.AddAttempt(tx, tx.EthTx(big.NewInt(1)), sentAt)
	assert.Nil(t, err)
	adapter := adapters.EthTx{}
	sentResult := models.RunResultWithValue(a.Hash.String())
	input := models.RunResultPending(sentResult)

	output := adapter.Perform(input, store)

	assert.False(t, output.HasError())
	assert.True(t, output.Pending)
	assert.Nil(t, store.One("ID", tx.ID, tx))
	attempts, _ := store.AttemptsFor(tx.ID)
	assert.Equal(t, 1, len(attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthTxAdapterFromPendingBumpGas(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store
	config := store.Config

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	sentAt := uint64(23456)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold))
	ethMock.Register("eth_sendRawTransaction", cltest.NewTxHash())

	from := store.KeyStore.GetAccount().Address
	tx := cltest.NewTx(from, sentAt)
	assert.Nil(t, store.Save(tx))
	a, err := store.AddAttempt(tx, tx.EthTx(big.NewInt(1)), 1)
	assert.Nil(t, err)
	adapter := adapters.EthTx{}
	sentResult := models.RunResultWithValue(a.Hash.String())
	input := models.RunResultPending(sentResult)

	output := adapter.Perform(input, store)

	assert.False(t, output.HasError())
	assert.True(t, output.Pending)
	assert.Nil(t, store.One("ID", tx.ID, tx))
	attempts, _ := store.AttemptsFor(tx.ID)
	assert.Equal(t, 2, len(attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthTxAdapterFromPendingConfirm(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
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

	tx := cltest.NewTx(cltest.NewEthAddress(), sentAt)
	assert.Nil(t, store.Save(tx))
	store.AddAttempt(tx, tx.EthTx(big.NewInt(1)), sentAt)
	store.AddAttempt(tx, tx.EthTx(big.NewInt(2)), sentAt+1)
	a3, _ := store.AddAttempt(tx, tx.EthTx(big.NewInt(3)), sentAt+2)
	adapter := adapters.EthTx{}
	sentResult := models.RunResultWithValue(a3.Hash.String())
	input := models.RunResultPending(sentResult)

	assert.False(t, tx.Confirmed)

	output := adapter.Perform(input, store)

	assert.False(t, output.Pending)
	assert.False(t, output.HasError())

	assert.Nil(t, store.One("ID", tx.ID, tx))
	assert.True(t, tx.Confirmed)
	attempts, _ := store.AttemptsFor(tx.ID)
	assert.False(t, attempts[0].Confirmed)
	assert.True(t, attempts[1].Confirmed)
	assert.False(t, attempts[2].Confirmed)

	assert.True(t, ethMock.AllCalled())
}

func TestEthTxAdapterWithError(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	store := app.Store
	eth := app.MockEthClient()
	eth.RegisterError("eth_getTransactionCount", "Cannot connect to nodes")

	adapter := adapters.EthTx{
		Address:    cltest.NewEthAddress(),
		FunctionID: models.HexToFunctionID("0xb3f98adc"),
	}
	input := models.RunResultWithValue("")
	output := adapter.Perform(input, store)

	assert.True(t, output.HasError())
	assert.Equal(t, output.Error(), "Cannot connect to nodes")
}
