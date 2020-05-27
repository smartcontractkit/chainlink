package adapters_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/big"
	"syscall"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/eth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common/hexutil"
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
			txManager.On("CheckAttempt", mock.Anything, mock.Anything).Once().Return(&eth.TxReceipt{}, test.receiptState, nil)

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
	txManager.On("CheckAttempt", mock.Anything, mock.Anything).Return(&eth.TxReceipt{}, strpkg.Unconfirmed, nil)
	store.TxManager = txManager

	adapter := adapters.EthTx{DataFormat: "bytes", DataPrefix: hexutil.MustDecode("0x88888888")}
	input := cltest.NewRunInputWithResult("cönfirmed")
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
	txManager.On("BumpGasUntilSafe", mock.Anything).Return(&eth.TxReceipt{}, strpkg.Confirmed, nil)
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
	receipt := &eth.TxReceipt{Hash: receiptHash, BlockNumber: cltest.Int(129831)}
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
	var receipts []eth.TxReceipt
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
	receipt := &eth.TxReceipt{Hash: receiptHash, BlockNumber: cltest.Int(129831)}
	txManager.On("BumpGasUntilSafe", mock.Anything).Return(receipt, strpkg.Safe, nil)
	store.TxManager = txManager

	adapter := adapters.EthTx{}
	data := cltest.JSONFromString(t, `{
		"ethereumReceipts": [{}],
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
	var receipts []eth.TxReceipt
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
	txManager.On("CheckAttempt", txAttempt, uint64(0)).Return(&eth.TxReceipt{}, strpkg.Confirmed, nil)

	result = adapter.Perform(input, store)
	require.NoError(t, result.Error())

	txManager.AssertExpectations(t)
}
