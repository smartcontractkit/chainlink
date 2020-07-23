package adapters_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"syscall"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestEthTxAdapter_Perform(t *testing.T) {
	t.Parallel()

	gasPrice := utils.NewBig(big.NewInt(187))
	gasLimit := uint64(911)

	tests := []struct {
		name         string
		input        string
		format       string
		receiptState strpkg.AttemptState
		output       string
		finalStatus  models.RunStatus
	}{
		{
			"safe",
			"0xf7fffff1",
			"",
			strpkg.Safe,
			"0x0000000000000000000000000000000000000000000000000000000000000000f7fffff1",
			models.RunStatusCompleted,
		},
		{
			"safe with bytes format",
			"cönfirmed",
			"bytes",
			strpkg.Safe,
			"0x" +
				"00000000" + // function selector
				"0000000000000000000000000000000000000000000000000000000000000020" + // offset
				"000000000000000000000000000000000000000000000000000000000000000a" + // length in bytes = 10, umlaut = 2 bytes
				"63c3b66e6669726d656400000000000000000000000000000000000000000000", // encoded string left padded
			models.RunStatusCompleted,
		},
		{
			"confirmed",
			"0x19999990",
			"",
			strpkg.Confirmed,
			"0x000000000000000000000000000000000000000000000000000000000000000019999990",
			models.RunStatusPendingOutgoingConfirmations,
		},
		{
			"confirmed with bytes format",
			"cönfirmed",
			"bytes",
			strpkg.Confirmed,
			"0x" +
				"00000000" + // function selector
				"0000000000000000000000000000000000000000000000000000000000000020" + // offset
				"000000000000000000000000000000000000000000000000000000000000000a" + // length in bytes = 10, umlaut = 2 bytes
				"63c3b66e6669726d656400000000000000000000000000000000000000000000", // encoded string left padded
			models.RunStatusPendingOutgoingConfirmations,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			txManager := new(mocks.TxManager)
			txManager.On("Connected").Once().Return(true)

			tx := &models.Tx{Attempts: []*models.TxAttempt{&models.TxAttempt{}}}
			txData := hexutil.MustDecode(test.output)
			txManager.On("CreateTxWithGas", mock.Anything, mock.Anything, txData, gasPrice.ToInt(), gasLimit).Once().Return(tx, nil)
			txManager.On("CheckAttempt", mock.Anything, mock.Anything).Once().Return(&types.Receipt{}, test.receiptState, nil)

			store.TxManager = txManager

			adapter := adapters.EthTx{DataFormat: test.format, GasPrice: gasPrice, GasLimit: gasLimit}
			input := cltest.NewRunInputWithResult(test.input)
			result := adapter.Perform(input, store)

			assert.NoError(t, result.Error())
			assert.Equal(t, test.finalStatus, result.Status())

			txManager.AssertExpectations(t)
		})
	}
}

func TestEthTxAdapter_Perform_BytesFormatWithDataPrefix(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	txManager := new(mocks.TxManager)
	tx := &models.Tx{Attempts: []*models.TxAttempt{&models.TxAttempt{}}}
	txManager.On("Connected").Maybe().Return(true)
	txManager.On("CreateTxWithGas", mock.Anything, mock.Anything,
		hexutil.MustDecode("0x"+
			"00000000"+ // function selector
			"88888888"+ // data prefix
			"0000000000000000000000000000000000000000000000000000000000000040"+ // offset
			"000000000000000000000000000000000000000000000000000000000000000a"+ // length in bytes
			"63c3b66e6669726d656400000000000000000000000000000000000000000000"), // encoded string left padded
		mock.Anything, mock.Anything).Return(tx, nil)
	txManager.On("CheckAttempt", mock.Anything, mock.Anything).Return(&types.Receipt{}, strpkg.Unconfirmed, nil)
	store.TxManager = txManager

	adapter := adapters.EthTx{DataFormat: "bytes", DataPrefix: hexutil.MustDecode("0x88888888")}
	input := cltest.NewRunInputWithResult("cönfirmed")
	result := adapter.Perform(input, store)

	assert.NoError(t, result.Error())
	assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, result.Status())

	txManager.AssertExpectations(t)
}

func TestEthTxAdapter_Perform_Preformatted(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	hexPayload := "b72f443a17edf4a55f766cf3c83469e6f96494b16823a41a4acb25800f30310300000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000001b76616c69646174696f6e2e747769747465722e757365726e616d650000000000000000000000000000000000000000000000000000000000000000000000001c76616c69646174696f6e2e747769747465722e7369676e617475726500000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000a64657261696e6265726b00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008430783965353831646633383765376138343433653636336435386663313736303034356433666362376165643234393633376163633737376262653837643561333934326431363130643265386538626437353066643533633230643466633661383536303737623235656439653538356439616161336439646535626365376238316200000000000000000000000000000000000000000000000000000000"
	fs := "0xdeadcafe"

	txManager := new(mocks.TxManager)
	tx := &models.Tx{Attempts: []*models.TxAttempt{&models.TxAttempt{}}}
	txManager.On("Connected").Maybe().Return(true)
	txManager.On("CreateTxWithGas", mock.Anything, mock.Anything, hexutil.MustDecode(fs+hexPayload), mock.Anything, mock.Anything).Return(tx, nil)
	txManager.On("CheckAttempt", mock.Anything, mock.Anything).Return(&types.Receipt{}, strpkg.Unconfirmed, nil)
	store.TxManager = txManager

	adapter := adapters.EthTx{
		FunctionSelector: models.HexToFunctionSelector(fs),
		DataFormat:       "preformatted",
	}
	input := cltest.NewRunInputWithResult("0x" + hexPayload)
	result := adapter.Perform(input, store)

	assert.NoError(t, result.Error())
	assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, result.Status())

	txManager.AssertExpectations(t)
}

func TestEthTxAdapter_Perform_FromPendingOutgoingConfirmations_StillPending(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	txManager := new(mocks.TxManager)
	txManager.On("Connected").Return(true)
	txManager.On("BumpGasUntilSafe", mock.Anything).Return(&types.Receipt{}, strpkg.Confirmed, nil)
	store.TxManager = txManager

	adapter := adapters.EthTx{}
	input := *models.NewRunInputWithResult(
		models.NewID(), *models.NewID(), cltest.NewHash(), models.RunStatusPendingOutgoingConfirmations,
	)
	output := adapter.Perform(input, store)

	require.NoError(t, output.Error())
	assert.True(t, output.Status().PendingOutgoingConfirmations())
	assert.Equal(t, input.Data(), output.Data())

	txManager.AssertExpectations(t)
}

func TestEthTxAdapter_Perform_FromPendingOutgoingConfirmations_Safe(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	txManager := new(mocks.TxManager)
	txManager.On("Connected").Return(true)
	receiptHash := cltest.NewHash()
	receipt := &types.Receipt{TxHash: receiptHash, BlockNumber: big.NewInt(129831), Logs: []*types.Log{}, PostState: []byte{}}
	txManager.On("BumpGasUntilSafe", mock.Anything).Return(receipt, strpkg.Safe, nil)
	store.TxManager = txManager

	adapter := adapters.EthTx{}
	input := *models.NewRunInputWithResult(
		models.NewID(), *models.NewID(), cltest.NewHash(), models.RunStatusPendingOutgoingConfirmations,
	)
	output := adapter.Perform(input, store)

	require.NoError(t, output.Error())
	assert.Equal(t, models.RunStatusCompleted, output.Status())
	assert.Equal(t, receiptHash.String(), output.Result().String())

	receiptsJSON := output.Get("ethereumReceipts").String()
	fmt.Println("JSON ~>", receiptsJSON)
	var receipts []types.Receipt
	require.NoError(t, json.Unmarshal([]byte(receiptsJSON), &receipts))
	require.Len(t, receipts, 1)
	assert.Equal(t, receipt, &receipts[0])

	latestOutgoingTxHash := output.Get("latestOutgoingTxHash").String()
	assert.Equal(t, receiptHash.String(), latestOutgoingTxHash)

	txManager.AssertExpectations(t)
}

func TestEthTxAdapter_Perform_AppendingTransactionReceipts(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	txManager := new(mocks.TxManager)
	txManager.On("Connected").Return(true)
	receiptHash := cltest.NewHash()
	receipt := &types.Receipt{TxHash: receiptHash, BlockNumber: big.NewInt(129831), Logs: []*types.Log{}}
	txManager.On("BumpGasUntilSafe", mock.Anything).Return(receipt, strpkg.Safe, nil)
	store.TxManager = txManager

	adapter := adapters.EthTx{}
	data := cltest.JSONFromString(t, `{
		"ethereumReceipts": [{
            "cumulativeGasUsed": "0x0",
            "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
            "logs": []
        }],
		"result":"0x3f839aaf5915da8714313a57b9c0a362d1a9a3fac1210190ace5cf3b008d780f"
	}`)
	input := *models.NewRunInput(
		models.NewID(), *models.NewID(), data, models.RunStatusPendingOutgoingConfirmations,
	)
	output := adapter.Perform(input, store)

	require.NoError(t, output.Error())
	assert.Equal(t, models.RunStatusCompleted, output.Status())
	assert.Equal(t, receiptHash.String(), output.Result().String())

	receiptsJSON := output.Get("ethereumReceipts").String()
	var receipts []types.Receipt
	require.NoError(t, json.Unmarshal([]byte(receiptsJSON), &receipts))
	require.Len(t, receipts, 2)

	latestOutgoingTxHash := output.Get("latestOutgoingTxHash").String()
	assert.Equal(t, receiptHash.String(), latestOutgoingTxHash)

	txManager.AssertExpectations(t)
}

func TestEthTxAdapter_Perform_WithError(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	txManager := new(mocks.TxManager)
	txManager.On("Connected").Return(true)
	txManager.On("CreateTxWithGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("Cannot connect to node"))
	store.TxManager = txManager

	adapter := adapters.EthTx{}
	input := cltest.NewRunInputWithResult("0x9786856756")
	output := adapter.Perform(input, store)
	assert.NoError(t, output.Error())

	txManager.AssertExpectations(t)
}

func TestEthTxAdapter_Perform_PendingOutgoingConfirmations_WithFatalErrorInTxManager(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	txManager := new(mocks.TxManager)
	txManager.On("Connected").Return(true)
	txManager.On("BumpGasUntilSafe", mock.Anything).Return(nil, strpkg.Unknown, errors.New("Fatal"))
	store.TxManager = txManager

	adapter := adapters.EthTx{}
	input := *models.NewRunInputWithResult(
		models.NewID(), *models.NewID(), cltest.NewHash().String(), models.RunStatusPendingOutgoingConfirmations,
	)
	output := adapter.Perform(input, store)

	assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, output.Status())
	assert.NoError(t, output.Error())

	txManager.AssertExpectations(t)
}

func TestEthTxAdapter_Perform_PendingOutgoingConfirmations_WithRecoverableErrorInTxManager(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	txManager := new(mocks.TxManager)
	txManager.On("Connected").Return(true)
	txManager.On("BumpGasUntilSafe", mock.Anything).Return(nil, strpkg.Confirmed, errors.New("Connection reset by peer"))
	store.TxManager = txManager

	adapter := adapters.EthTx{}
	input := *models.NewRunInputWithResult(
		models.NewID(), *models.NewID(), cltest.NewHash().String(), models.RunStatusPendingOutgoingConfirmations,
	)
	output := adapter.Perform(input, store)

	require.NoError(t, output.Error())
	assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, output.Status())
	assert.Equal(t, input.Data(), output.Data())

	txManager.AssertExpectations(t)
}

func TestEthTxAdapter_Perform_NotConnectedWhenPendingOutgoingConfirmations(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	txManager := new(mocks.TxManager)
	txManager.On("Connected").Return(false)
	store.TxManager = txManager

	adapter := adapters.EthTx{}
	input := *models.NewRunInputWithResult(models.NewID(), *models.NewID(), cltest.NewHash().String(), models.RunStatusPendingOutgoingConfirmations)
	output := adapter.Perform(input, store)

	require.NoError(t, output.Error())
	assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, output.Status())
	assert.Equal(t, input.Data(), output.Data())

	txManager.AssertExpectations(t)
}

func TestEthTxAdapter_Perform_NotConnected(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	txManager := new(mocks.TxManager)
	txManager.On("Connected").Return(false)
	store.TxManager = txManager

	adapter := adapters.EthTx{}
	data := adapter.Perform(models.RunInput{}, store)

	require.NoError(t, data.Error())
	assert.Equal(t, models.RunStatusPendingConnection, data.Status())

	txManager.AssertExpectations(t)
}

func TestEthTxAdapter_Perform_CreateTxWithGasErrorTreatsAsNotConnected(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	txManager := new(mocks.TxManager)
	txManager.On("Connected").Return(true)
	txManager.On("CreateTxWithGas",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil, syscall.ETIMEDOUT)
	store.TxManager = txManager

	adapter := adapters.EthTx{}
	data := adapter.Perform(models.RunInput{}, store)

	require.NoError(t, data.Error())
	assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, data.Status())

	txManager.AssertExpectations(t)
}

func TestEthTxAdapter_Perform_CheckAttemptErrorTreatsAsNotConnected(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	txManager := new(mocks.TxManager)
	txManager.On("Connected").Return(true)
	txManager.On("CreateTxWithGas",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(&models.Tx{
		Attempts: []*models.TxAttempt{&models.TxAttempt{}},
	}, nil)
	txManager.On("CheckAttempt", mock.Anything, mock.Anything).Return(nil, strpkg.Unknown, syscall.EWOULDBLOCK)
	store.TxManager = txManager

	adapter := adapters.EthTx{}
	data := adapter.Perform(models.RunInput{}, store)

	require.NoError(t, data.Error())
	assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, data.Status())

	txManager.AssertExpectations(t)
}

func TestEthTxAdapter_Perform_CreateTxWithEmptyResponseErrorTreatsAsPendingOutgoingConfirmations(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	from := cltest.NewAddress()
	tx := cltest.CreateTx(t, store, from, 1)

	badResponseErr := errors.New("Bad response on request: [ TransactionIndex ]. Error cause was EmptyResponse, (majority count: 94 / total: 94)")
	txManager := new(mocks.TxManager)
	txManager.On("Connected").Return(true)
	txManager.On("CreateTxWithGas",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(tx, nil)
	txManager.On("CheckAttempt", mock.Anything, mock.Anything).Return(nil, strpkg.Unknown, badResponseErr)
	store.TxManager = txManager

	adapter := adapters.EthTx{}
	output := adapter.Perform(models.RunInput{}, store)

	require.NoError(t, output.Error())
	assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, output.Status())

	// Have a head come through with the same empty response
	txManager.On("Connected").Return(true)
	txManager.On("BumpGasUntilSafe", mock.Anything).Return(nil, strpkg.Unknown, badResponseErr)

	input := *models.NewRunInput(models.NewID(), *models.NewID(), output.Data(), output.Status())
	output = adapter.Perform(input, store)
	require.NoError(t, output.Error())
	assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, output.Status())

	txManager.AssertExpectations(t)
}

func TestEthTxAdapter_Perform_NoDoubleSpendOnSendTransactionFail(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	var sentData []byte

	txManager := new(mocks.TxManager)
	txManager.On("Connected").Return(true)
	txManager.On("CreateTxWithGas",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(data []byte) bool {
			sentData = data
			return len(data) > 0
		}),
		mock.Anything,
		mock.Anything).Once().Return(nil, errors.New("no bueno"))
	store.TxManager = txManager

	adapter := adapters.EthTx{}
	input := cltest.NewRunInputWithResult("0x9786856756")
	result := adapter.Perform(input, store)
	require.NoError(t, result.Error())

	txAttempt := &models.TxAttempt{}
	tx := &models.Tx{Attempts: []*models.TxAttempt{txAttempt}}
	txManager.On("CreateTxWithGas",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(data []byte) bool {
			return bytes.Equal(sentData, data)
		}),
		mock.Anything,
		mock.Anything).Once().Return(tx, nil)
	txManager.On("CheckAttempt", txAttempt, uint64(0)).Return(&types.Receipt{}, strpkg.Confirmed, nil)

	result = adapter.Perform(input, store)
	require.NoError(t, result.Error())

	txManager.AssertExpectations(t)
}

func TestEthTxAdapter_Perform_BPTXM(t *testing.T) {
	t.Parallel()

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	config.Config.Set("ENABLE_BULLETPROOF_TX_MANAGER", true)
	store, cleanup := cltest.NewStoreWithConfig(config)
	notifier := new(mocks.NotifyNewEthTx)
	notifier.On("Trigger").Return()
	store.NotifyNewEthTx = notifier
	defer cleanup()

	toAddress := cltest.NewAddress()
	gasLimit := uint64(42)
	functionSelector := models.HexToFunctionSelector("0x70a08231") // balanceOf(address)
	dataPrefix := hexutil.MustDecode("0x88888888")

	t.Run("with valid data and empty DataFormat writes to database and returns run output pending outgoing confirmations", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:        toAddress,
			GasLimit:         gasLimit,
			FunctionSelector: functionSelector,
			DataPrefix:       dataPrefix,
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		input := models.NewRunInputWithResult(jobRunID, taskRunID, "0x9786856756", models.RunStatusUnstarted)
		runOutput := adapter.Perform(*input, store)
		require.NoError(t, runOutput.Error())
		assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, runOutput.Status())

		etrt, err := store.FindEthTaskRunTxByTaskRunID(input.TaskRunID().UUID())
		require.NoError(t, err)

		assert.Equal(t, taskRunID.UUID(), etrt.TaskRunID)
		require.NotNil(t, etrt.EthTx)
		assert.Nil(t, etrt.EthTx.Nonce)
		assert.Equal(t, toAddress, etrt.EthTx.ToAddress)
		assert.Equal(t, "70a08231888888880000000000000000000000000000000000000000000000000000009786856756", hex.EncodeToString(etrt.EthTx.EncodedPayload))
		assert.Equal(t, gasLimit, etrt.EthTx.GasLimit)
		assert.Equal(t, models.EthTxUnstarted, etrt.EthTx.State)

		notifier.AssertExpectations(t)
	})

	t.Run("if FromAddresses is provided but no key matches, returns job error", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:        toAddress,
			FromAddresses:    []gethCommon.Address{cltest.NewAddress()},
			GasLimit:         gasLimit,
			FunctionSelector: functionSelector,
			DataPrefix:       dataPrefix,
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		input := models.NewRunInputWithResult(jobRunID, taskRunID, "0x9786856756", models.RunStatusUnstarted)
		runOutput := adapter.Perform(*input, store)
		require.EqualError(t, runOutput.Error(), "insertEthTx failed to pickFromAddress: no keys available")
		assert.Equal(t, models.RunStatusErrored, runOutput.Status())
	})

	t.Run("with bytes DataFormat writes correct encoded data to database", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:        toAddress,
			GasLimit:         gasLimit,
			FunctionSelector: functionSelector,
			DataPrefix:       dataPrefix,
			DataFormat:       "bytes",
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		input := models.NewRunInputWithResult(jobRunID, taskRunID, "cönfirmed", models.RunStatusUnstarted)
		runOutput := adapter.Perform(*input, store)
		require.NoError(t, runOutput.Error())
		assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, runOutput.Status())

		expectedData := hexutil.MustDecode(
			functionSelector.String() +
				"88888888" + // dataPrefix
				"0000000000000000000000000000000000000000000000000000000000000040" + // offset
				"000000000000000000000000000000000000000000000000000000000000000a" + // length in bytes
				"63c3b66e6669726d656400000000000000000000000000000000000000000000") // encoded string left padded

		etrt, err := store.FindEthTaskRunTxByTaskRunID(input.TaskRunID().UUID())
		require.NoError(t, err)

		assert.Equal(t, expectedData, etrt.EthTx.EncodedPayload)
	})

	t.Run("with invalid data returns run output error and does not write to DB", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:        toAddress,
			GasLimit:         gasLimit,
			FunctionSelector: functionSelector,
			DataPrefix:       dataPrefix,
			DataFormat:       "some old bollocks",
		}
		jobRunID := models.NewID()
		taskRunID := models.NewID()
		input := models.NewRunInputWithResult(jobRunID, *taskRunID, "0x9786856756", models.RunStatusUnstarted)
		runOutput := adapter.Perform(*input, store)
		assert.Contains(t, runOutput.Error().Error(), "while constructing EthTx data: unsupported format: some old bollocks")
		assert.Equal(t, models.RunStatusErrored, runOutput.Status())

		trtx, err := store.FindEthTaskRunTxByTaskRunID(input.TaskRunID().UUID())
		require.NoError(t, err)
		require.Nil(t, trtx)
	})

	t.Run("with unconfirmed transaction returns output pending confirmations", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:        toAddress,
			GasLimit:         gasLimit,
			FunctionSelector: functionSelector,
			DataPrefix:       dataPrefix,
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 0)
		store.DB.Exec(`INSERT INTO eth_task_run_txes (task_run_id, eth_tx_id) VALUES ($1, $2)`, taskRunID.UUID(), etx.ID)
		input := models.NewRunInputWithResult(jobRunID, taskRunID, "0x9786856756", models.RunStatusUnstarted)

		// Do the thing
		runOutput := adapter.Perform(*input, store)

		require.NoError(t, runOutput.Error())
		assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, runOutput.Status())
	})

	t.Run("with confirmed transaction returns pending outgoing confirmations if receipt is missing (invariant violation, should never happen)", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:        toAddress,
			GasLimit:         gasLimit,
			FunctionSelector: functionSelector,
			DataPrefix:       dataPrefix,
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 1, 1)
		store.DB.Exec(`INSERT INTO eth_task_run_txes (task_run_id, eth_tx_id) VALUES ($1, $2)`, taskRunID.UUID(), etx.ID)
		input := models.NewRunInputWithResult(jobRunID, taskRunID, "0x9786856756", models.RunStatusUnstarted)

		// Do the thing
		runOutput := adapter.Perform(*input, store)

		require.NoError(t, runOutput.Error())
		assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, runOutput.Status())
	})

	t.Run("with confirmed transaction with exactly one attempt with exactly one receipt that is younger than minRequiredOutgoingConfirmations, returns output pending_outgoing_confirmations", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:                        toAddress,
			GasLimit:                         gasLimit,
			FunctionSelector:                 functionSelector,
			DataPrefix:                       dataPrefix,
			MinRequiredOutgoingConfirmations: 12,
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 2, 1)

		confirmedAttemptHash := etx.EthTxAttempts[0].Hash

		cltest.MustInsertEthReceipt(t, store, 1, cltest.NewHash(), confirmedAttemptHash)
		require.NoError(t, store.IdempotentInsertHead(models.Head{
			Hash:   cltest.NewHash(),
			Number: 12,
		}))
		store.DB.Exec(`INSERT INTO eth_task_run_txes (task_run_id, eth_tx_id) VALUES ($1, $2)`, taskRunID.UUID(), etx.ID)
		input := models.NewRunInputWithResult(jobRunID, taskRunID, "0x9786856756", models.RunStatusUnstarted)

		// Do the thing
		runOutput := adapter.Perform(*input, store)

		require.NoError(t, runOutput.Error())
		assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, runOutput.Status())
	})

	t.Run("with confirmed transaction with exactly one attempt with exactly one receipt that is equal to minRequiredOutgoingConfirmations, returns output complete with transaction hash pulled from receipt", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:                        toAddress,
			GasLimit:                         gasLimit,
			FunctionSelector:                 functionSelector,
			DataPrefix:                       dataPrefix,
			MinRequiredOutgoingConfirmations: 12,
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 3, 1)

		confirmedAttemptHash := etx.EthTxAttempts[0].Hash

		cltest.MustInsertEthReceipt(t, store, 1, cltest.NewHash(), confirmedAttemptHash)
		require.NoError(t, store.IdempotentInsertHead(models.Head{
			Hash:   cltest.NewHash(),
			Number: 13,
		}))
		store.DB.Exec(`INSERT INTO eth_task_run_txes (task_run_id, eth_tx_id) VALUES ($1, $2)`, taskRunID.UUID(), etx.ID)
		data := cltest.JSONFromString(t, `{"foo": "bar", "result": "some old bollocks"}`)
		input := models.NewRunInput(jobRunID, taskRunID, data, models.RunStatusUnstarted)

		// Do the thing
		runOutput := adapter.Perform(*input, store)

		require.NoError(t, runOutput.Error())
		assert.Equal(t, models.RunStatusCompleted, runOutput.Status())
		assert.Equal(t, confirmedAttemptHash.Hex(), runOutput.Result().String())
		// Does not clobber previously assigned data
		assert.Equal(t, "bar", runOutput.Get("foo").String())
		// Assigns latestOutgoingTxHash for legacy compatibility
		assert.Equal(t, confirmedAttemptHash.Hex(), runOutput.Get("latestOutgoingTxHash").String())
	})

	t.Run("with confirmed transaction with exactly one attempt with exactly one receipt that is older than minRequiredOutgoingConfirmations, returns output complete with transaction hash pulled from receipt", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:                        toAddress,
			GasLimit:                         gasLimit,
			FunctionSelector:                 functionSelector,
			DataPrefix:                       dataPrefix,
			MinRequiredOutgoingConfirmations: 12,
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 4, 1)

		confirmedAttemptHash := etx.EthTxAttempts[0].Hash

		cltest.MustInsertEthReceipt(t, store, 1, cltest.NewHash(), confirmedAttemptHash)
		require.NoError(t, store.IdempotentInsertHead(models.Head{
			Hash:   cltest.NewHash(),
			Number: 14,
		}))
		store.DB.Exec(`INSERT INTO eth_task_run_txes (task_run_id, eth_tx_id) VALUES ($1, $2)`, taskRunID.UUID(), etx.ID)
		input := models.NewRunInputWithResult(jobRunID, taskRunID, "0x9786856756", models.RunStatusUnstarted)

		// Do the thing
		runOutput := adapter.Perform(*input, store)

		require.NoError(t, runOutput.Error())
		assert.Equal(t, models.RunStatusCompleted, runOutput.Status())
		assert.Equal(t, confirmedAttemptHash.Hex(), runOutput.Result().String())
	})

	t.Run("with confirmed transaction with two attempts, one of which has exactly one receipt that is older than the default MinRequiredOutgoingConfirmations, returns output complete with transaction hash pulled from receipt", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:        toAddress,
			GasLimit:         gasLimit,
			FunctionSelector: functionSelector,
			DataPrefix:       dataPrefix,
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 5, 1)
		attempt2 := cltest.MustInsertBroadcastEthTxAttempt(t, etx.ID, store, 2)

		confirmedAttemptHash := attempt2.Hash

		cltest.MustInsertEthReceipt(t, store, 1, cltest.NewHash(), confirmedAttemptHash)
		require.NoError(t, store.IdempotentInsertHead(models.Head{
			Hash:   cltest.NewHash(),
			Number: int64(store.Config.MinRequiredOutgoingConfirmations()) + 2,
		}))
		store.DB.Exec(`INSERT INTO eth_task_run_txes (task_run_id, eth_tx_id) VALUES ($1, $2)`, taskRunID.UUID(), etx.ID)
		input := models.NewRunInputWithResult(jobRunID, taskRunID, "0x9786856756", models.RunStatusUnstarted)

		// Do the thing
		runOutput := adapter.Perform(*input, store)

		require.NoError(t, runOutput.Error())
		assert.Equal(t, models.RunStatusCompleted, runOutput.Status())
		assert.Equal(t, confirmedAttemptHash.Hex(), runOutput.Result().String())
	})

	t.Run("with transaction that ended up in fatal_error state returns job run error", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:        toAddress,
			GasLimit:         gasLimit,
			FunctionSelector: functionSelector,
			DataPrefix:       dataPrefix,
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		etx := cltest.MustInsertFatalErrorEthTx(t, store)
		require.NoError(t, store.DB.Exec(`INSERT INTO eth_task_run_txes (task_run_id, eth_tx_id) VALUES ($1, $2)`, taskRunID.UUID(), etx.ID).Error)

		input := models.NewRunInputWithResult(jobRunID, taskRunID, "0x9786856756", models.RunStatusUnstarted)

		// Do the thing
		runOutput := adapter.Perform(*input, store)

		require.EqualError(t, runOutput.Error(), "something exploded")
		assert.Equal(t, models.RunStatusErrored, runOutput.Status())
		assert.Equal(t, "", runOutput.Result().String())
	})

}
