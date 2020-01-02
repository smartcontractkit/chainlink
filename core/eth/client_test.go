package eth_test

import (
	"encoding/json"
	"testing"

	"math/big"

	"chainlink/core/eth"
	"chainlink/core/internal/cltest"
	"chainlink/core/internal/mocks"
	strpkg "chainlink/core/store"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCallerSubscriberClient_GetTxReceipt(t *testing.T) {
	response := cltest.MustReadFile(t, "testdata/getTransactionReceipt.json")
	mockServer, wsCleanup := cltest.NewWSServer(string(response))
	defer wsCleanup()
	config := cltest.NewConfigWithWSServer(t, mockServer)
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	ec := store.TxManager.(*strpkg.EthTxManager).Client

	hash := common.HexToHash("0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238")
	receipt, err := ec.GetTxReceipt(hash)
	assert.NoError(t, err)
	assert.Equal(t, hash, receipt.Hash)
	assert.Equal(t, cltest.Int(uint64(11)), receipt.BlockNumber)
}

func TestTxReceipt_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantLogLen int
	}{
		{"basic", "testdata/getTransactionReceipt.json", 0},
		{"runlog request", "testdata/runlogReceipt.json", 4},
		{"runlog response", "testdata/responseReceipt.json", 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonStr := cltest.JSONFromFixture(t, test.path).Get("result").String()
			var receipt eth.TxReceipt
			err := json.Unmarshal([]byte(jsonStr), &receipt)
			require.NoError(t, err)

			assert.Equal(t, test.wantLogLen, len(receipt.Logs))
		})
	}
}

func TestTxReceipt_FulfilledRunlog(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"basic", "testdata/getTransactionReceipt.json", false},
		{"runlog request", "testdata/runlogReceipt.json", false},
		{"runlog response", "testdata/responseReceipt.json", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			receipt := cltest.TxReceiptFromFixture(t, test.path)
			assert.Equal(t, test.want, receipt.FulfilledRunLog())
		})
	}
}

func TestCallerSubscriberClient_GetNonce(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	require.NoError(t, app.Start())

	ethMock := app.MockCallerSubscriberClient()
	ethClientObject := app.Store.TxManager.(*strpkg.EthTxManager).Client
	ethMock.Register("eth_getTransactionCount", "0x0100")
	result, err := ethClientObject.GetNonce(cltest.NewAddress())
	assert.NoError(t, err)
	var expected uint64 = 256
	assert.Equal(t, result, expected)
}

func TestCallerSubscriberClient_SendRawTx(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	require.NoError(t, app.Start())

	ethMock := app.MockCallerSubscriberClient()
	ethClientObject := app.Store.TxManager.(*strpkg.EthTxManager).Client
	ethMock.Register("eth_sendRawTransaction", common.Hash{1})
	result, err := ethClientObject.SendRawTx("test")
	assert.NoError(t, err)
	assert.Equal(t, result, common.Hash{1})
}

func TestCallerSubscriberClient_GetEthBalance(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	require.NoError(t, app.Start())

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"basic", "0x0100", "0.000000000000000256"},
		{"larger than signed 64 bit integer", "0x4b3b4ca85a86c47a098a224000000000", "100000000000000000000.000000000000000000"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ethMock := app.MockCallerSubscriberClient()
			ethClientObject := app.Store.TxManager.(*strpkg.EthTxManager).Client

			ethMock.Register("eth_getBalance", test.input)
			result, err := ethClientObject.GetEthBalance(cltest.NewAddress())
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result.String())
		})
	}
}

func TestCallerSubscriberClient_GetERC20Balance(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	require.NoError(t, app.Start())

	ethMock := app.MockCallerSubscriberClient()
	ethClientObject := app.Store.TxManager.(*strpkg.EthTxManager).Client

	ethMock.Register("eth_call", "0x0100") // 256
	result, err := ethClientObject.GetERC20Balance(cltest.NewAddress(), cltest.NewAddress())
	assert.NoError(t, err)
	expected := big.NewInt(256)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	ethMock.Register("eth_call", "0x4b3b4ca85a86c47a098a224000000000") // 1e38
	result, err = ethClientObject.GetERC20Balance(cltest.NewAddress(), cltest.NewAddress())
	expected = big.NewInt(0)
	expected.SetString("100000000000000000000000000000000000000", 10)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCallerSubscriberClient_GetAggregatorPrice(t *testing.T) {
	caller := new(mocks.CallerSubscriber)
	ethClient := &eth.CallerSubscriberClient{CallerSubscriber: caller}
	address := cltest.NewAddress()

	// aggregatorLatestAnswerID is the first 4 bytes of the keccak256 of
	// Chainlink's aggregator latestAnswer function.
	const aggregatorLatestAnswerID = "50d25bcd"
	aggregatorLatestAnswerSelector := eth.HexToFunctionSelector(aggregatorLatestAnswerID)

	expectedCallArgs := eth.CallArgs{
		To:   address,
		Data: aggregatorLatestAnswerSelector.Bytes(),
	}

	tests := []struct {
		name, response string
		precision      int32
		expectation    decimal.Decimal
	}{
		{"hex", "0x0100", 2, decimal.NewFromFloat(2.56)},
		{"decimal", "10000000000000", 11, decimal.NewFromInt(100)},
		{"large decimal", "52050000000000000000", 11, decimal.RequireFromString("520500000")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			caller.On("Call", mock.Anything, "eth_call", expectedCallArgs, "latest").Return(nil).
				Run(func(args mock.Arguments) {
					res := args.Get(0).(*string)
					*res = test.response
				})
			result, err := ethClient.GetAggregatorPrice(address, test.precision)
			require.NoError(t, err)
			assert.True(t, test.expectation.Equal(result))
			caller.AssertExpectations(t)
		})
	}
}
