package adapters_test

import (
	"encoding/json"
	"errors"
	"math/big"
	"syscall"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEthTxAdapter_Perform_Confirmed(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	store := app.Store

	address := cltest.NewAddress()
	fHash := models.HexToFunctionSelector("b3f98adc")
	dataPrefix := hexutil.Bytes(
		hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000045746736453745"))
	inputValue := "0x9786856756"

	ethMock, err := app.MockStartAndConnect()
	require.NoError(t, err)

	hash := cltest.NewHash()
	sentAt := uint64(23456)
	confirmed := sentAt + 1
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

	adapter := adapters.EthTx{
		Address:          address,
		DataPrefix:       dataPrefix,
		FunctionSelector: fHash,
	}
	input := cltest.RunResultWithResult(inputValue)
	data := adapter.Perform(input, store)

	assert.False(t, data.HasError())

	from := cltest.GetAccountAddress(t, store)
	txs, err := store.TxFrom(from)
	assert.NoError(t, err)
	require.Len(t, txs, 1)
	assert.Len(t, txs[0].Attempts, 1)

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_Perform_ConfirmedWithBytes(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	store := app.Store

	address := cltest.NewAddress()
	fHash := models.HexToFunctionSelector("b3f98adc")
	dataPrefix := hexutil.Bytes(
		hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000045746736453745"))
	inputValue := "cönfirmed" // contains diacritic acute to check bytes counted for length not chars

	ethMock, err := app.MockStartAndConnect()
	require.NoError(t, err)

	hash := cltest.NewHash()
	sentAt := uint64(23456)
	confirmed := sentAt + 1
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

	adapter := adapters.EthTx{
		Address:          address,
		DataPrefix:       dataPrefix,
		FunctionSelector: fHash,
		DataFormat:       adapters.DataFormatBytes,
	}
	input := cltest.RunResultWithResult(inputValue)
	data := adapter.Perform(input, store)

	assert.False(t, data.HasError())

	from := cltest.GetAccountAddress(t, store)
	txs, err := store.TxFrom(from)
	assert.NoError(t, err)
	require.Len(t, txs, 1)
	assert.Len(t, txs[0].Attempts, 1)

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_Perform_SafeWithBytesAndNoDataPrefix(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	store := app.Store

	address := cltest.NewAddress()
	fHash := models.HexToFunctionSelector("b3f98adc")
	// contains diacritic acute to check bytes counted for length not chars
	inputValue := "cönfirmed"

	ethMock, err := app.MockStartAndConnect()
	require.NoError(t, err)

	hash := cltest.NewHash()
	currentHeight := uint64(23456)
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
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(currentHeight))
	safe := currentHeight - store.Config.MinOutgoingConfirmations()
	receipt := models.TxReceipt{Hash: hash, BlockNumber: cltest.Int(safe)}
	ethMock.Register("eth_getTransactionReceipt", receipt)

	adapter := adapters.EthTx{
		Address:          address,
		FunctionSelector: fHash,
		DataFormat:       adapters.DataFormatBytes,
	}
	input := cltest.RunResultWithResult(inputValue)
	data := adapter.Perform(input, store)

	assert.False(t, data.HasError())
	assert.Equal(t, string(models.RunStatusCompleted), string(data.Status))

	from := cltest.GetAccountAddress(t, store)
	var txs []models.Tx
	gomega.NewGomegaWithT(t).Eventually(func() []models.Tx {
		var err error
		txs, err = store.TxFrom(from)
		assert.NoError(t, err)
		return txs
	}).Should(gomega.HaveLen(1))
	require.Len(t, txs, 1)
	assert.Len(t, txs[0].Attempts, 1)

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_Perform_FromPendingConfirmations_StillPending(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	store := app.Store
	config := store.Config
	ethMock := app.MockEthClient()

	ethMock.Register("eth_getTransactionReceipt", models.TxReceipt{})
	sentAt := uint64(23456)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold()-1))
	ethMock.Register("eth_chainId", config.ChainID())

	require.NoError(t, app.StartAndConnect())

	from := cltest.GetAccountAddress(t, store)
	tx := cltest.CreateTx(t, store, from, sentAt)
	a := tx.Attempts[0]
	adapter := adapters.EthTx{}
	sentResult := cltest.RunResultWithResult(a.Hash.String())
	sentResult.MarkPendingConfirmations()

	output := adapter.Perform(sentResult, store)

	assert.False(t, output.HasError())
	assert.True(t, output.Status.PendingConfirmations())
	tx, err := store.FindTx(tx.ID)
	require.NoError(t, err)
	assert.Len(t, tx.Attempts, 1)

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_Perform_FromPendingConfirmations_BumpGas(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	store := app.Store
	config := store.Config

	ethMock, err := app.MockStartAndConnect()
	require.NoError(t, err)

	sentAt := uint64(23456)
	ethMock.Context("ethtx perform", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionReceipt", models.TxReceipt{})
		ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold()))
		ethMock.Register("eth_sendRawTransaction", cltest.NewHash())
	})

	from := cltest.GetAccountAddress(t, store)
	tx := cltest.CreateTx(t, store, from, sentAt)
	a := tx.Attempts[0]

	adapter := adapters.EthTx{}
	sentResult := cltest.RunResultWithResult(a.Hash.String())
	sentResult.MarkPendingConfirmations()

	output := adapter.Perform(sentResult, store)
	assert.False(t, output.HasError())
	assert.True(t, output.Status.PendingConfirmations())
	tx, err = store.FindTx(tx.ID)
	require.NoError(t, err)
	assert.Len(t, tx.Attempts, 2)

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_Perform_FromPendingConfirmations_ConfirmCompletes(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	store := app.Store
	config := store.Config
	sentAt := uint64(23456)

	ethMock := app.MockEthClient(cltest.Strict)
	ethMock.Register("eth_getTransactionCount", `0x100`)
	ethMock.Register("eth_call", "0x1")
	ethMock.Register("eth_getBalance", "0x100")

	ethMock.Register("eth_chainId", store.Config.ChainID())
	ethMock.Register("eth_getTransactionReceipt", models.TxReceipt{})
	confirmedHash := cltest.NewHash()
	receipt := models.TxReceipt{Hash: confirmedHash, BlockNumber: cltest.Int(sentAt)}
	ethMock.Register("eth_getTransactionReceipt", receipt)
	confirmedAt := sentAt + config.MinOutgoingConfirmations() - 1 // confirmations are 0-based idx
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(confirmedAt))

	require.NoError(t, app.StartAndConnect())

	tx := cltest.NewTx(cltest.NewAddress(), sentAt)
	tx.GasPrice = models.NewBig(config.EthGasPriceDefault())
	require.NoError(t, store.DB.Save(tx).Error)
	store.AddTxAttempt(tx, tx.EthTx(big.NewInt(2)), sentAt+1)
	a3, _ := store.AddTxAttempt(tx, tx.EthTx(big.NewInt(3)), sentAt+2)
	adapter := adapters.EthTx{}
	sentResult := cltest.RunResultWithResult(a3.Hash.String())
	sentResult.MarkPendingConfirmations()

	assert.False(t, tx.Confirmed)

	output := adapter.Perform(sentResult, store)

	assert.True(t, output.Status.Completed())
	assert.False(t, output.HasError())
	value, err := output.ResultString()
	assert.Nil(t, err)
	assert.Equal(t, confirmedHash.String(), value)

	tx, err = store.FindTx(tx.ID)
	require.NoError(t, err)
	assert.True(t, tx.Confirmed)
	require.Len(t, tx.Attempts, 2)
	assert.False(t, tx.Attempts[0].Confirmed)
	assert.True(t, tx.Attempts[1].Confirmed)

	receiptsJSON := output.Get("ethereumReceipts").String()
	var receipts []models.TxReceipt
	assert.NoError(t, json.Unmarshal([]byte(receiptsJSON), &receipts))
	assert.Equal(t, 1, len(receipts))
	assert.Equal(t, receipt, receipts[0])

	confirmedTxHex := output.Get("latestOutgoingTxHash").String()
	assert.Equal(t, confirmedHash, common.HexToHash(confirmedTxHex))

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_Perform_AppendingTransactionReceipts(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	store := app.Store
	config := store.Config
	sentAt := uint64(23456)

	ethMock := app.MockEthClient()
	receipt := models.TxReceipt{Hash: cltest.NewHash(), BlockNumber: cltest.Int(sentAt)}
	ethMock.Register("eth_getTransactionReceipt", receipt)
	confirmedAt := sentAt + config.MinOutgoingConfirmations() - 1 // confirmations are 0-based idx
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(confirmedAt))
	ethMock.Register("eth_chainId", config.ChainID())

	require.NoError(t, app.StartAndConnect())

	tx := cltest.NewTx(cltest.NewAddress(), sentAt)
	tx.GasPrice = models.NewBig(config.EthGasPriceDefault())
	require.NoError(t, store.DB.Save(tx).Error)
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

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	store := app.Store
	ethMock, err := app.MockStartAndConnect()
	require.NoError(t, err)

	adapter := adapters.EthTx{
		Address:          cltest.NewAddress(),
		FunctionSelector: models.HexToFunctionSelector("0xb3f98adc"),
	}
	input := cltest.RunResultWithResult("0x9786856756")
	ethMock.RegisterError("eth_blockNumber", "Cannot connect to nodes")
	output := adapter.Perform(input, store)

	assert.True(t, output.HasError())
	assert.Contains(t, output.Error(), "Cannot connect to nodes")
}

func TestEthTxAdapter_Perform_WithErrorInvalidInput(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	store := app.Store
	ethMock := app.MockEthClient()
	ethMock.Register("eth_chainId", store.Config.ChainID())
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
	assert.Contains(t, output.Error(), "Cannot connect to nodes")
}

func TestEthTxAdapter_Perform_PendingConfirmations_WithFatalErrorInTxManager(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	store := app.Store
	ethMock := app.MockEthClient(cltest.Strict)
	ethMock.Register("eth_getTransactionCount", `0x17`)
	ethMock.Register("eth_chainId", store.Config.ChainID())
	assert.Nil(t, app.Start())

	require.NoError(t, app.WaitForConnection())

	adapter := adapters.EthTx{
		Address:          cltest.NewAddress(),
		FunctionSelector: models.HexToFunctionSelector("0xb3f98adc"),
	}
	input := cltest.RunResultWithResult(cltest.NewHash().String())
	input.Status = models.RunStatusPendingConfirmations
	ethMock.RegisterError("eth_blockNumber", "Invalid node id")
	output := adapter.Perform(input, store)

	ethMock.AssertAllCalled()

	assert.Equal(t, models.RunStatusErrored, output.Status)
	assert.NotNil(t, output.Error())
}

func TestEthTxAdapter_Perform_PendingConfirmations_WithRecoverableErrorInTxManager(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	store := app.Store
	ethMock := app.MockEthClient(cltest.Strict)
	ethMock.Register("eth_getTransactionCount", `0x12`)
	ethMock.Register("eth_chainId", store.Config.ChainID())
	assert.Nil(t, app.Start())

	from := cltest.GetAccountAddress(t, store)
	tx := cltest.CreateTx(t, store, from, uint64(14372))
	input := cltest.RunResultWithResult(tx.Attempts[0].Hash.String())
	input.Status = models.RunStatusPendingConfirmations

	ethMock.Register("eth_blockNumber", "0x100")
	ethMock.RegisterError("eth_getTransactionReceipt", "Connection reset by peer")

	require.NoError(t, app.WaitForConnection())

	adapter := adapters.EthTx{
		Address:          cltest.NewAddress(),
		FunctionSelector: models.HexToFunctionSelector("0xb3f98adc"),
	}
	output := adapter.Perform(input, store)

	ethMock.AssertAllCalled()

	assert.Equal(t, models.RunStatusPendingConfirmations, output.Status)
	assert.NoError(t, output.GetError())
}

func TestEthTxAdapter_DeserializationBytesFormat(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ctrl := gomock.NewController(t)
	txmMock := mocks.NewMockTxManager(ctrl)
	store.TxManager = txmMock

	txAttempt := &models.TxAttempt{}
	tx := &models.Tx{Attempts: []*models.TxAttempt{txAttempt}}
	txmMock.EXPECT().Register(gomock.Any())
	txmMock.EXPECT().Connected().Return(true).AnyTimes()
	txmMock.EXPECT().CreateTxWithGas(gomock.Any(), gomock.Any(), hexutil.MustDecode(
		"0x00000000"+
			"0000000000000000000000000000000000000000000000000000000000000020"+
			"000000000000000000000000000000000000000000000000000000000000000b"+
			"68656c6c6f20776f726c64000000000000000000000000000000000000000000"),
		gomock.Any(), gomock.Any()).Return(tx, nil)
	txmMock.EXPECT().CheckAttempt(txAttempt, uint64(0)).Return(&models.TxReceipt{}, strpkg.Unconfirmed, nil)

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

	store, cleanup := cltest.NewStore(t)
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

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	store := app.Store

	adapter := adapters.EthTx{}
	input := models.RunResult{}
	data := adapter.Perform(input, store)

	assert.False(t, data.HasError())
	assert.Equal(t, models.RunStatusPendingConnection, data.Status)
}

func TestEthTxAdapter_Perform_CreateTxWithGasErrorTreatsAsNotConnected(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	store := app.Store

	ctrl := gomock.NewController(t)
	txmMock := mocks.NewMockTxManager(ctrl)
	store.TxManager = txmMock
	txmMock.EXPECT().Register(gomock.Any())
	txmMock.EXPECT().Connected().Return(true)
	txmMock.EXPECT().CreateTxWithGas(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(nil, syscall.ETIMEDOUT)

	adapter := adapters.EthTx{}
	input := models.RunResult{}
	data := adapter.Perform(input, store)

	assert.False(t, data.HasError())
	assert.Equal(t, models.RunStatusPendingConnection, data.Status)
}

func TestEthTxAdapter_Perform_CheckAttemptErrorTreatsAsNotConnected(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	store := app.Store

	ctrl := gomock.NewController(t)
	txmMock := mocks.NewMockTxManager(ctrl)
	store.TxManager = txmMock
	txmMock.EXPECT().Register(gomock.Any())
	txmMock.EXPECT().Connected().Return(true)
	txmMock.EXPECT().CreateTxWithGas(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&models.Tx{
		Attempts: []*models.TxAttempt{&models.TxAttempt{}},
	}, nil)
	txmMock.EXPECT().CheckAttempt(gomock.Any(), gomock.Any()).Return(nil, strpkg.Unknown, syscall.EWOULDBLOCK)

	adapter := adapters.EthTx{}
	input := models.RunResult{}
	data := adapter.Perform(input, store)

	assert.False(t, data.HasError())
	assert.Equal(t, models.RunStatusPendingConnection, data.Status)
}

func TestEthTxAdapter_Perform_CreateTxWithEmptyResponseErrorTreatsAsPendingConfirmations(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	store := app.Store

	from := cltest.GetAccountAddress(t, store)
	tx := cltest.CreateTx(t, store, from, 1)

	ctrl := gomock.NewController(t)
	txmMock := mocks.NewMockTxManager(ctrl)
	store.TxManager = txmMock
	txmMock.EXPECT().Register(gomock.Any())
	txmMock.EXPECT().Connected().Return(true)
	txmMock.EXPECT().CreateTxWithGas(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(tx, nil)
	txmMock.EXPECT().CheckAttempt(
		gomock.Any(),
		gomock.Any(),
	).Return(nil, strpkg.Unknown, errors.New("Bad response on request: [ TransactionIndex ]. Error cause was EmptyResponse, (majority count: 94 / total: 94)"))

	adapter := adapters.EthTx{}
	input := models.RunResult{}
	input = adapter.Perform(input, store)

	assert.False(t, input.HasError())
	assert.Equal(t, models.RunStatusPendingConfirmations, input.Status)

	// Have a head come through with the same empty response
	txmMock.EXPECT().Connected().Return(true)
	txmMock.EXPECT().BumpGasUntilSafe(
		gomock.Any(),
	).Return(nil, strpkg.Unknown, errors.New("Bad response on request: [ TransactionIndex ]. Error cause was EmptyResponse, (majority count: 94 / total: 94)"))

	input = adapter.Perform(input, store)
	assert.False(t, input.HasError())
	assert.Equal(t, models.RunStatusPendingConfirmations, input.Status)
}

func TestEthTxAdapter_Perform_NoDoubleSpendOnSendTransactionFail(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	store := app.Store
	ethMock := app.MockEthClient(cltest.Strict)
	ethMock.Register("eth_getTransactionCount", `0x1`)
	ethMock.Register("eth_chainId", store.Config.ChainID())

	address := cltest.NewAddress()
	fHash := models.HexToFunctionSelector("b3f98adc")
	dataPrefix := hexutil.Bytes(
		hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000045746736453745"))
	inputValue := "0x9786856756"

	assert.Nil(t, app.StartAndConnect())

	// Run the adapter, but make sure the transaction sending fails

	hash := cltest.NewHash()
	sentAt := uint64(9183)

	var firstTxData []interface{}
	ethMock.Register("eth_sendRawTransaction", hash,
		func(_ interface{}, data ...interface{}) error {
			firstTxData = data
			return errors.New("no bueno")
		})
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))

	adapter := adapters.EthTx{
		Address:          address,
		DataPrefix:       dataPrefix,
		FunctionSelector: fHash,
	}
	input := cltest.RunResultWithResult(inputValue)
	input.CachedJobRunID = models.NewID()
	data := adapter.Perform(input, store)
	assert.Error(t, data.GetError())

	// Run the adapter again

	confirmed := sentAt + 1
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(confirmed))

	var secondTxData []interface{}
	ethMock.Register("eth_sendRawTransaction", hash,
		func(_ interface{}, data ...interface{}) error {
			secondTxData = data
			return nil
		})
	receipt := models.TxReceipt{Hash: hash, BlockNumber: cltest.Int(confirmed)}
	ethMock.Register("eth_getTransactionReceipt", receipt)

	data = adapter.Perform(input, store)
	assert.NoError(t, data.GetError())

	// The first and second transaction should have the same data
	assert.Equal(t, firstTxData, secondTxData)

	addresses := cltest.GetAccountAddresses(store)
	require.Len(t, addresses, 1)

	// There should only be one transaction with one attempt
	transactions, err := store.TxFrom(addresses[0])
	require.NoError(t, err)
	require.Len(t, transactions, 1)
	assert.Len(t, transactions[0].Attempts, 1)

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_Perform_NoDoubleSpendOnSendTransactionFailAndNonceChange(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	store := app.Store

	ethMock := app.MockEthClient(cltest.Strict)
	ethMock.Register("eth_getTransactionCount", `0x1`)
	ethMock.Register("eth_getTransactionCount", `0x2`)
	ethMock.Register("eth_chainId", store.Config.ChainID())
	app.AddUnlockedKey()

	addresses := cltest.GetAccountAddresses(store)
	require.Len(t, addresses, 2)

	address := cltest.NewAddress()
	fHash := models.HexToFunctionSelector("b3f98adc")
	dataPrefix := hexutil.Bytes(
		hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000045746736453745"))
	inputValue := "0x9786856756"

	assert.Nil(t, app.StartAndConnect())

	// Run the adapter, but make sure the transaction sending fails

	hash := cltest.NewHash()
	sentAt := uint64(9183)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
	var firstTxData []interface{}
	ethMock.Register("eth_sendRawTransaction", hash,
		func(_ interface{}, data ...interface{}) error {
			firstTxData = data
			return errors.New("no bueno")
		})
	receipt := models.TxReceipt{Hash: hash, BlockNumber: cltest.Int(sentAt)}
	ethMock.Register("eth_getTransactionReceipt", receipt)

	adapter := adapters.EthTx{
		Address:          address,
		DataPrefix:       dataPrefix,
		FunctionSelector: fHash,
	}
	input := cltest.RunResultWithResult(inputValue)
	input.CachedJobRunID = models.NewID()
	data := adapter.Perform(input, store)
	assert.Error(t, data.GetError())

	// Run the adapter again

	confirmed := sentAt + 1
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(confirmed))

	var secondTxData []interface{}
	ethMock.Register("eth_sendRawTransaction", hash,
		func(_ interface{}, data ...interface{}) error {
			secondTxData = data
			return nil
		})

	data = adapter.Perform(input, store)
	assert.NoError(t, data.GetError())

	// Since the nonce (and from address) changed, the data should also change
	assert.NotEqual(t, firstTxData, secondTxData)

	// The original account should have no txes, because it was reassigned
	txs, err := store.TxFrom(addresses[0])
	require.NoError(t, err)
	assert.Len(t, txs, 0)

	// The second account should have only one tx
	txs, err = store.TxFrom(addresses[1])
	require.NoError(t, err)
	require.Len(t, txs, 1)
	assert.Len(t, txs[0].Attempts, 1)

	ethMock.EventuallyAllCalled(t)
}

func TestEthTxAdapter_IsClientRetriable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		error     error
		retriable bool
	}{
		{"nil error", nil, false},
		{"http invalid method", errors.New("net/http: invalid method SGET"), false},
		{"syscall.ECONNRESET", syscall.ECONNRESET, false},
		{"syscall.ECONNABORTED", syscall.ECONNABORTED, false},
		{"syscall.EWOULDBLOCK", syscall.EWOULDBLOCK, true},
		{"syscall.ETIMEDOUT", syscall.ETIMEDOUT, true},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.retriable, adapters.IsClientRetriable(test.error))
		})
	}
}
