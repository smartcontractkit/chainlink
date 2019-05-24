package adapters_test

import (
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEthTxAdapter_Perform_Confirmed(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()
	store := app.Store
	config := store.Config

	address := cltest.NewAddress()
	fHash := models.HexToFunctionSelector("b3f98adc")
	dataPrefix := hexutil.Bytes(
		hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000045746736453745"))
	inputValue := "0x9786856756"

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", `0x0100`)
	assert.Nil(t, app.StartAndConnect())

	hash := cltest.NewHash()
	sentAt := uint64(23456)
	confirmed := sentAt + 1
	safe := confirmed + config.MinOutgoingConfirmations()
	ethMock.Register("eth_sendRawTransaction", hash,
		func(_ interface{}, data ...interface{}) error {
			rlp := data[0].([]interface{})[0].(string)
			tx, err := utils.DecodeEthereumTx(rlp)
			assert.NoError(t, err)
			assert.Equal(t, address.String(), tx.To().String())
			wantData := "0x" +
				"b3f98adc" +
				"0000000000000000000000000000000000000000000000000045746736453745" +
				"0000000000000000000000000000000000000000000000000000009786856756"
			assert.Equal(t, wantData, hexutil.Encode(tx.Data()))
			return nil
		})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
	receipt := models.TxReceipt{Hash: hash, BlockNumber: cltest.Int(confirmed)}
	ethMock.Register("eth_getTransactionReceipt", receipt)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(safe))

	adapter := adapters.EthTx{
		Address:          address,
		DataPrefix:       dataPrefix,
		FunctionSelector: fHash,
	}
	input := cltest.RunResultWithResult(inputValue)
	data := adapter.Perform(input, store)

	assert.False(t, data.HasError())

	from := cltest.GetAccountAddress(store)
	txs, err := store.TxFrom(from)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(txs))
	attempts, _ := store.TxAttemptsFor(txs[0].ID)
	assert.Equal(t, 1, len(attempts))

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_Perform_ConfirmedWithBytes(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()
	store := app.Store
	config := store.Config

	address := cltest.NewAddress()
	fHash := models.HexToFunctionSelector("b3f98adc")
	dataPrefix := hexutil.Bytes(
		hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000045746736453745"))
	inputValue := "cönfirmed" // contains diacritic acute to check bytes counted for length not chars

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", `0x0100`)
	require.NoError(t, app.StartAndConnect())

	hash := cltest.NewHash()
	sentAt := uint64(23456)
	confirmed := sentAt + 1
	safe := confirmed + config.MinOutgoingConfirmations()
	ethMock.Register("eth_sendRawTransaction", hash,
		func(_ interface{}, data ...interface{}) error {
			rlp := data[0].([]interface{})[0].(string)
			tx, err := utils.DecodeEthereumTx(rlp)
			assert.NoError(t, err)
			assert.Equal(t, address.String(), tx.To().String())
			wantData := "0x" +
				"b3f98adc" +
				"0000000000000000000000000000000000000000000000000045746736453745" +
				"0000000000000000000000000000000000000000000000000000000000000040" +
				"000000000000000000000000000000000000000000000000000000000000000a" +
				"63c3b66e6669726d656400000000000000000000000000000000000000000000"
			assert.Equal(t, wantData, hexutil.Encode(tx.Data()))
			return nil
		})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
	receipt := models.TxReceipt{Hash: hash, BlockNumber: cltest.Int(confirmed)}
	ethMock.Register("eth_getTransactionReceipt", receipt)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(safe))

	adapter := adapters.EthTx{
		Address:          address,
		DataPrefix:       dataPrefix,
		FunctionSelector: fHash,
		DataFormat:       adapters.DataFormatBytes,
	}
	input := cltest.RunResultWithResult(inputValue)
	data := adapter.Perform(input, store)

	assert.False(t, data.HasError())

	from := cltest.GetAccountAddress(store)
	txs, err := store.TxFrom(from)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(txs))
	attempts, _ := store.TxAttemptsFor(txs[0].ID)
	assert.Equal(t, 1, len(attempts))

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_Perform_ConfirmedWithBytesAndNoDataPrefix(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()
	store := app.Store
	config := store.Config

	address := cltest.NewAddress()
	fHash := models.HexToFunctionSelector("b3f98adc")
	// contains diacritic acute to check bytes counted for length not chars
	inputValue := "cönfirmed"

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", `0x0100`)
	require.NoError(t, app.StartAndConnect())

	hash := cltest.NewHash()
	sentAt := uint64(23456)
	confirmed := sentAt + 1
	safe := confirmed + config.MinOutgoingConfirmations()
	ethMock.Register("eth_sendRawTransaction", hash,
		func(_ interface{}, data ...interface{}) error {
			rlp := data[0].([]interface{})[0].(string)
			tx, err := utils.DecodeEthereumTx(rlp)
			assert.NoError(t, err)
			assert.Equal(t, address.String(), tx.To().String())
			wantData := "0x" +
				"b3f98adc" +
				"0000000000000000000000000000000000000000000000000000000000000020" +
				"000000000000000000000000000000000000000000000000000000000000000a" +
				"63c3b66e6669726d656400000000000000000000000000000000000000000000"
			assert.Equal(t, wantData, hexutil.Encode(tx.Data()))
			return nil
		})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
	receipt := models.TxReceipt{Hash: hash, BlockNumber: cltest.Int(confirmed)}
	ethMock.Register("eth_getTransactionReceipt", receipt)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(safe))

	adapter := adapters.EthTx{
		Address:          address,
		FunctionSelector: fHash,
		DataFormat:       adapters.DataFormatBytes,
	}
	input := cltest.RunResultWithResult(inputValue)
	data := adapter.Perform(input, store)

	assert.False(t, data.HasError())
	assert.Equal(t, string(models.RunStatusCompleted), string(data.Status))

	from := cltest.GetAccountAddress(store)
	var txs []models.Tx
	gomega.NewGomegaWithT(t).Eventually(func() []models.Tx {
		var err error
		txs, err = store.TxFrom(from)
		assert.NoError(t, err)
		return txs
	}).Should(gomega.HaveLen(1))
	attempts, _ := store.TxAttemptsFor(txs[0].ID)
	assert.Equal(t, 1, len(attempts))

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_Perform_FromPendingConfirmations_StillPending(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()
	store := app.Store
	config := store.Config

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", models.TxReceipt{})
	sentAt := uint64(23456)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold()-1))

	require.NoError(t, app.StartAndConnect())

	from := cltest.GetAccountAddress(store)
	tx := cltest.NewTx(from, sentAt)
	assert.Nil(t, store.SaveTx(tx))
	a, err := store.AddTxAttempt(tx, tx.EthTx(big.NewInt(1)), sentAt)
	assert.NoError(t, err)
	adapter := adapters.EthTx{}
	sentResult := cltest.RunResultWithResult(a.Hash.String())
	input := sentResult
	input.MarkPendingConfirmations()

	output := adapter.Perform(input, store)

	assert.False(t, output.HasError())
	assert.True(t, output.Status.PendingConfirmations())
	_, err = store.FindTx(tx.ID)
	assert.NoError(t, err)
	attempts, _ := store.TxAttemptsFor(tx.ID)
	assert.Equal(t, 1, len(attempts))

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_Perform_FromPendingConfirmations_BumpGas(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()
	store := app.Store
	config := store.Config

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", "0x0")
	require.NoError(t, app.StartAndConnect())

	sentAt := uint64(23456)
	ethMock.Context("ethtx perform", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionReceipt", models.TxReceipt{})
		ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold()))
		ethMock.Register("eth_sendRawTransaction", cltest.NewHash())
	})

	from := cltest.GetAccountAddress(store)
	tx := cltest.NewTx(from, sentAt)
	assert.Nil(t, store.SaveTx(tx))
	a, err := store.AddTxAttempt(tx, tx.EthTx(big.NewInt(1)), 1)
	assert.NoError(t, err)
	adapter := adapters.EthTx{}
	sentResult := cltest.RunResultWithResult(a.Hash.String())
	input := sentResult
	input.MarkPendingConfirmations()

	output := adapter.Perform(input, store)

	assert.False(t, output.HasError())
	assert.True(t, output.Status.PendingConfirmations())
	_, err = store.FindTx(tx.ID)
	assert.NoError(t, err)
	attempts, _ := store.TxAttemptsFor(tx.ID)
	assert.Equal(t, 2, len(attempts))

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_Perform_FromPendingConfirmations_ConfirmCompletes(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()

	store := app.Store
	config := store.Config
	sentAt := uint64(23456)

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", models.TxReceipt{})
	confirmedHash := cltest.NewHash()
	receipt := models.TxReceipt{Hash: confirmedHash, BlockNumber: cltest.Int(sentAt)}
	ethMock.Register("eth_getTransactionReceipt", receipt)
	confirmedAt := sentAt + config.MinOutgoingConfirmations() - 1 // confirmations are 0-based idx
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(confirmedAt))

	require.NoError(t, app.StartAndConnect())

	tx := cltest.NewTx(cltest.NewAddress(), sentAt)
	assert.Nil(t, store.SaveTx(tx))
	store.AddTxAttempt(tx, tx.EthTx(big.NewInt(1)), sentAt)
	store.AddTxAttempt(tx, tx.EthTx(big.NewInt(2)), sentAt+1)
	a3, _ := store.AddTxAttempt(tx, tx.EthTx(big.NewInt(3)), sentAt+2)
	adapter := adapters.EthTx{}
	sentResult := cltest.RunResultWithResult(a3.Hash.String())
	input := sentResult
	input.MarkPendingConfirmations()

	assert.False(t, tx.Confirmed)

	output := adapter.Perform(input, store)

	assert.True(t, output.Status.Completed())
	assert.False(t, output.HasError())
	value, err := output.ResultString()
	assert.Nil(t, err)
	assert.Equal(t, confirmedHash.String(), value)

	tx, err = store.FindTx(tx.ID)
	assert.NoError(t, err)
	assert.True(t, tx.Confirmed)
	attempts, _ := store.TxAttemptsFor(tx.ID)
	assert.False(t, attempts[0].Confirmed)
	assert.True(t, attempts[1].Confirmed)
	assert.False(t, attempts[2].Confirmed)

	receiptsJSON := output.Get("ethereumReceipts").String()
	var receipts []models.TxReceipt
	assert.NoError(t, json.Unmarshal([]byte(receiptsJSON), &receipts))
	assert.Equal(t, 1, len(receipts))
	assert.Equal(t, receipt, receipts[0])

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_Perform_AppendingTransactionReceipts(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()

	store := app.Store
	config := store.Config
	sentAt := uint64(23456)

	ethMock := app.MockEthClient()
	receipt := models.TxReceipt{Hash: cltest.NewHash(), BlockNumber: cltest.Int(sentAt)}
	ethMock.Register("eth_getTransactionReceipt", receipt)
	confirmedAt := sentAt + config.MinOutgoingConfirmations() - 1 // confirmations are 0-based idx
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(confirmedAt))

	require.NoError(t, app.StartAndConnect())

	tx := cltest.NewTx(cltest.NewAddress(), sentAt)
	assert.Nil(t, store.SaveTx(tx))
	a, err := store.AddTxAttempt(tx, tx.EthTx(big.NewInt(1)), sentAt)
	assert.NoError(t, err)
	adapter := adapters.EthTx{}
	sentResult := cltest.RunResultWithResult(a.Hash.String())

	input := sentResult
	input.MarkPendingConfirmations()
	previousReceipt := models.TxReceipt{Hash: cltest.NewHash(), BlockNumber: cltest.Int(sentAt - 10)}
	input.Add("ethereumReceipts", []models.TxReceipt{previousReceipt})

	output := adapter.Perform(input, store)
	assert.True(t, output.Status.Completed())
	receiptsJSON := output.Get("ethereumReceipts").String()
	var receipts []models.TxReceipt
	assert.NoError(t, json.Unmarshal([]byte(receiptsJSON), &receipts))
	assert.Equal(t, []models.TxReceipt{previousReceipt, receipt}, receipts)

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_Perform_WithError(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()

	store := app.Store
	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", `0x0100`)
	require.NoError(t, app.StartAndConnect())

	adapter := adapters.EthTx{
		Address:          cltest.NewAddress(),
		FunctionSelector: models.HexToFunctionSelector("0xb3f98adc"),
	}
	input := cltest.RunResultWithResult("0x9786856756")
	ethMock.RegisterError("eth_blockNumber", "Cannot connect to nodes")
	output := adapter.Perform(input, store)

	assert.True(t, output.HasError())
	assert.Equal(t, "Cannot connect to nodes", output.Error())
}

func TestEthTxAdapter_Perform_WithErrorInvalidInput(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()

	store := app.Store
	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", `0x0100`)
	require.NoError(t, app.StartAndConnect())

	adapter := adapters.EthTx{
		Address:          cltest.NewAddress(),
		FunctionSelector: models.HexToFunctionSelector("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF1"),
	}
	input := cltest.RunResultWithResult("0x9786856756")
	ethMock.RegisterError("eth_blockNumber", "Cannot connect to nodes")
	output := adapter.Perform(input, store)

	assert.True(t, output.HasError())
	assert.Equal(t, "Cannot connect to nodes", output.Error())
}

func TestEthTxAdapter_Perform_PendingConfirmations_WithErrorInTxManager(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()

	store := app.Store
	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", `0x0100`)
	assert.Nil(t, app.Start())

	adapter := adapters.EthTx{
		Address:          cltest.NewAddress(),
		FunctionSelector: models.HexToFunctionSelector("0xb3f98adc"),
	}
	input := cltest.RunResultWithResult("")
	input.Status = models.RunStatusPendingConfirmations
	ethMock.RegisterError("eth_blockNumber", "Cannot connect to nodes")
	output := adapter.Perform(input, store)

	assert.False(t, output.HasError())
}

func TestEthTxAdapter_DeserializationBytesFormat(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()
	ctrl := gomock.NewController(t)
	txmMock := mocks.NewMockTxManager(ctrl)
	store.TxManager = txmMock
	txmMock.EXPECT().Register(gomock.Any())
	txmMock.EXPECT().Connected().Return(true).AnyTimes()
	txmMock.EXPECT().CreateTxWithGas(gomock.Any(), hexutil.MustDecode(
		"0x00000000"+
			"0000000000000000000000000000000000000000000000000000000000000020"+
			"000000000000000000000000000000000000000000000000000000000000000b"+
			"68656c6c6f20776f726c64000000000000000000000000000000000000000000"),
		gomock.Any(), gomock.Any()).Return(&models.Tx{}, nil)
	txmMock.EXPECT().BumpGasUntilSafe(gomock.Any())

	task := models.TaskSpec{}
	err := json.Unmarshal([]byte(`{"type": "EthTx", "params": {"format": "bytes"}}`), &task)
	assert.NoError(t, err)
	assert.Equal(t, task.Type, adapters.TaskTypeEthTx)

	adapter, err := adapters.For(task, store)
	assert.NoError(t, err)
	ethtx, ok := adapter.BaseAdapter.(*adapters.EthTx)
	assert.True(t, ok)
	assert.Equal(t, ethtx.DataFormat, adapters.DataFormatBytes)

	input := models.RunResult{
		Data:   cltest.JSONFromString(t, `{"result": "hello world"}`),
		Status: models.RunStatusInProgress,
	}
	result := adapter.Perform(input, store)
	assert.False(t, result.HasError())
	assert.Equal(t, result.Error(), "")
}

func TestEthTxAdapter_Perform_CustomGas(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	gasPrice := big.NewInt(187)
	gasLimit := uint64(911)

	ctrl := gomock.NewController(t)
	txmMock := mocks.NewMockTxManager(ctrl)
	store.TxManager = txmMock
	txmMock.EXPECT().Register(gomock.Any())
	txmMock.EXPECT().Connected()
	txmMock.EXPECT().CreateTxWithGas(
		gomock.Any(),
		gomock.Any(),
		gasPrice,
		gasLimit,
	).Return(&models.Tx{}, nil)
	txmMock.EXPECT().BumpGasUntilSafe(gomock.Any())

	adapter := adapters.EthTx{
		Address:          cltest.NewAddress(),
		FunctionSelector: models.HexToFunctionSelector("0xb3f98adc"),
		GasPrice:         models.NewBig(gasPrice),
		GasLimit:         gasLimit,
	}

	input := models.RunResult{
		Data:   cltest.JSONFromString(t, `{"result": "hello world"}`),
		Status: models.RunStatusInProgress,
	}

	result := adapter.Perform(input, store)
	assert.False(t, result.HasError())
}

func TestEthTxAdapter_Perform_NotConnected(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()
	store := app.Store

	adapter := adapters.EthTx{}
	input := models.RunResult{}
	data := adapter.Perform(input, store)

	assert.False(t, data.HasError())
	assert.Equal(t, models.RunStatusPendingConnection, data.Status)
}

func TestEthTxAdapter_Perform_NoDoubleSpendOnSendTransactionFail(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()
	store := app.Store
	config := store.Config

	ethMock := app.MockEthClient(cltest.Strict)
	ethMock.Register("eth_getTransactionCount", `0x1`)
	ethMock.Register("eth_getTransactionCount", `0x2`)
	app.AddUnlockedKey()

	address := cltest.NewAddress()
	fHash := models.HexToFunctionSelector("b3f98adc")
	dataPrefix := hexutil.Bytes(
		hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000045746736453745"))
	inputValue := "0x9786856756"

	assert.Nil(t, app.StartAndConnect())

	// Run the adapter, but make sure the transaction sending fails

	hash := cltest.NewHash()
	sentAt := uint64(9183)
	ethMock.Register("eth_getBalance", "0x100")
	ethMock.Register("eth_call", "0x100")
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
	var firstTxData []interface{}
	ethMock.Register("eth_sendRawTransaction", hash,
		func(_ interface{}, data ...interface{}) error {
			firstTxData = data
			return errors.New("no bueno")
		})

	adapter := adapters.EthTx{
		Address:          address,
		DataPrefix:       dataPrefix,
		FunctionSelector: fHash,
	}
	input := cltest.RunResultWithResult(inputValue)
	data := adapter.Perform(input, store)
	assert.Error(t, data.GetError())

	// Run the adapter again

	confirmed := sentAt + 1
	safe := confirmed + config.MinOutgoingConfirmations()
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(confirmed))

	var secondTxData []interface{}
	ethMock.Register("eth_sendRawTransaction", hash,
		func(_ interface{}, data ...interface{}) error {
			secondTxData = data
			return nil
		})
	receipt := models.TxReceipt{Hash: hash, BlockNumber: cltest.Int(confirmed)}
	ethMock.Register("eth_getTransactionReceipt", receipt)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(safe))

	data = adapter.Perform(input, store)
	assert.NoError(t, data.GetError())

	// The first and second transaction should have the same data
	assert.Equal(t, firstTxData, secondTxData)

	from := cltest.GetAccountAddress(store)
	txs, err := store.TxFrom(from)
	require.NoError(t, err)
	require.Len(t, txs, 1)
	attempts, _ := store.TxAttemptsFor(txs[0].ID)
	assert.Len(t, attempts, 1)

	ethMock.EventuallyAllCalled(t)
}
