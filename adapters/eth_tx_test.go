package adapters_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	strpkg "github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

func TestEthTxAdapter_Perform_Confirmed(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store
	config := store.Config

	address := cltest.NewAddress()
	fHash := models.HexToFunctionSelector("b3f98adc")
	dataPrefix := hexutil.Bytes(hexutil.MustDecode("0x45746736453745"))
	inputValue := "0x9786856756"

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", `0x0100`)
	hash := cltest.NewHash()
	sentAt := uint64(23456)
	confirmed := sentAt + 1
	safe := confirmed + config.EthMinConfirmations
	ethMock.Register("eth_sendRawTransaction", hash,
		func(_ interface{}, data ...interface{}) error {
			rlp := data[0].([]interface{})[0].(string)
			tx, err := utils.DecodeEthereumTx(rlp)
			assert.Nil(t, err)
			assert.Equal(t, address.String(), tx.To().String())
			wantData := utils.HexConcat(fHash.String(), dataPrefix.String(), inputValue)
			assert.Equal(t, wantData, hexutil.Encode(tx.Data()))
			return nil
		})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
	receipt := strpkg.TxReceipt{Hash: hash, BlockNumber: cltest.BigHexInt(confirmed)}
	ethMock.Register("eth_getTransactionReceipt", receipt)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(safe))

	adapter := adapters.EthTx{
		Address:          address,
		DataPrefix:       dataPrefix,
		FunctionSelector: fHash,
	}
	input := cltest.RunResultWithValue(inputValue)
	data := adapter.Perform(input, store)

	assert.False(t, data.HasError())

	from := store.KeyStore.GetAccount().Address
	txs := []models.Tx{}
	assert.Nil(t, store.Where("From", from, &txs))
	assert.Equal(t, 1, len(txs))
	attempts, _ := store.AttemptsFor(txs[0].ID)
	assert.Equal(t, 1, len(attempts))

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_Perform_FromPending(t *testing.T) {
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
	sentResult := cltest.RunResultWithValue(a.Hash.String())
	input := sentResult.MarkPending()

	output := adapter.Perform(input, store)

	assert.False(t, output.HasError())
	assert.True(t, output.Pending)
	assert.Nil(t, store.One("ID", tx.ID, tx))
	attempts, _ := store.AttemptsFor(tx.ID)
	assert.Equal(t, 1, len(attempts))

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_Perform_FromPendingBumpGas(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store
	config := store.Config

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	sentAt := uint64(23456)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold))
	ethMock.Register("eth_sendRawTransaction", cltest.NewHash())

	from := store.KeyStore.GetAccount().Address
	tx := cltest.NewTx(from, sentAt)
	assert.Nil(t, store.Save(tx))
	a, err := store.AddAttempt(tx, tx.EthTx(big.NewInt(1)), 1)
	assert.Nil(t, err)
	adapter := adapters.EthTx{}
	sentResult := cltest.RunResultWithValue(a.Hash.String())
	input := sentResult.MarkPending()

	output := adapter.Perform(input, store)

	assert.False(t, output.HasError())
	assert.True(t, output.Pending)
	assert.Nil(t, store.One("ID", tx.ID, tx))
	attempts, _ := store.AttemptsFor(tx.ID)
	assert.Equal(t, 2, len(attempts))

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_Perform_FromPendingConfirm(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store
	config := store.Config

	sentAt := uint64(23456)

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{
		Hash:        cltest.NewHash(),
		BlockNumber: cltest.BigHexInt(sentAt),
	})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthMinConfirmations))

	tx := cltest.NewTx(cltest.NewAddress(), sentAt)
	assert.Nil(t, store.Save(tx))
	store.AddAttempt(tx, tx.EthTx(big.NewInt(1)), sentAt)
	store.AddAttempt(tx, tx.EthTx(big.NewInt(2)), sentAt+1)
	a3, _ := store.AddAttempt(tx, tx.EthTx(big.NewInt(3)), sentAt+2)
	adapter := adapters.EthTx{}
	sentResult := cltest.RunResultWithValue(a3.Hash.String())
	input := sentResult.MarkPending()

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

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_Perform_WithError(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	store := app.Store
	ethMock := app.MockEthClient()
	ethMock.RegisterError("eth_getTransactionCount", "Cannot connect to nodes")

	adapter := adapters.EthTx{
		Address:          cltest.NewAddress(),
		FunctionSelector: models.HexToFunctionSelector("0xb3f98adc"),
	}
	input := cltest.RunResultWithValue("")
	output := adapter.Perform(input, store)

	assert.True(t, output.HasError())
	assert.Equal(t, output.Error(), "Cannot connect to nodes")
}
