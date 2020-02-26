package eth_test

import (
	"encoding/json"
	"testing"

	"math/big"

	"chainlink/core/eth"
	"chainlink/core/internal/cltest"
	"chainlink/core/internal/mocks"
	strpkg "chainlink/core/store"
	"chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
		test := test
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
		test := test
		t.Run(test.name, func(t *testing.T) {
			receipt := cltest.TxReceiptFromFixture(t, test.path)
			assert.Equal(t, test.want, receipt.FulfilledRunLog())
		})
	}
}

func TestCallerSubscriberClient_GetNonce(t *testing.T) {
	t.Parallel()

	ethClientMock := new(mocks.CallerSubscriber)
	ethClient := &eth.CallerSubscriberClient{CallerSubscriber: ethClientMock}
	address := cltest.NewAddress()
	response := "0x0100"

	ethClientMock.On("Call", mock.Anything, "eth_getTransactionCount", address.String(), "pending").
		Return(nil).
		Run(func(args mock.Arguments) {
			res := args.Get(0).(*string)
			*res = response
		})

	result, err := ethClient.GetNonce(address)
	require.NoError(t, err)

	var expected uint64 = 256
	require.Equal(t, result, expected)
}

func TestCallerSubscriberClient_SendRawTx(t *testing.T) {
	t.Parallel()

	ethClientMock := new(mocks.CallerSubscriber)
	ethClient := &eth.CallerSubscriberClient{CallerSubscriber: ethClientMock}
	txData := "0xdeadbeef"
	returnedHash := cltest.NewHash()

	ethClientMock.On("Call", mock.Anything, "eth_sendRawTransaction", txData).
		Return(nil).
		Run(func(args mock.Arguments) {
			res := args.Get(0).(*common.Hash)
			*res = returnedHash
		})

	result, err := ethClient.SendRawTx(txData)
	assert.NoError(t, err)
	assert.Equal(t, result, returnedHash)
}

func TestCallerSubscriberClient_GetEthBalance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"basic", "0x0100", "0.000000000000000256"},
		{"larger than signed 64 bit integer", "0x4b3b4ca85a86c47a098a224000000000", "100000000000000000000.000000000000000000"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ethClientMock := new(mocks.CallerSubscriber)
			ethClient := &eth.CallerSubscriberClient{CallerSubscriber: ethClientMock}

			ethClientMock.On("Call", mock.Anything, "eth_getBalance", mock.Anything, "latest").
				Return(nil).
				Run(func(args mock.Arguments) {
					res := args.Get(0).(*string)
					*res = test.input
				})

			result, err := ethClient.GetEthBalance(cltest.NewAddress())
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result.String())
		})
	}
}

func TestCallerSubscriberClient_GetERC20Balance(t *testing.T) {
	t.Parallel()

	expectedBig, _ := big.NewInt(0).SetString("100000000000000000000000000000000000000", 10)

	tests := []struct {
		name     string
		input    string
		expected *big.Int
	}{
		{"small", "0x0100", big.NewInt(256)},
		{"big", "0x4b3b4ca85a86c47a098a224000000000", expectedBig},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {

			ethClientMock := new(mocks.CallerSubscriber)
			ethClient := &eth.CallerSubscriberClient{CallerSubscriber: ethClientMock}

			contractAddress := cltest.NewAddress()
			userAddress := cltest.NewAddress()

			functionSelector := eth.HexToFunctionSelector("0x70a08231") // balanceOf(address)
			data := utils.ConcatBytes(functionSelector.Bytes(), common.LeftPadBytes(userAddress.Bytes(), utils.EVMWordByteLen))
			callArgs := eth.CallArgs{
				To:   contractAddress,
				Data: data,
			}

			ethClientMock.On("Call", mock.Anything, "eth_call", callArgs, "latest").
				Return(nil).
				Run(func(args mock.Arguments) {
					res := args.Get(0).(*string)
					*res = test.input
				})

			result, err := ethClient.GetERC20Balance(userAddress, contractAddress)
			assert.NoError(t, err)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestCallerSubscriberClient_GetAggregatorPrice(t *testing.T) {
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
		{"hex - Zero", "0x", 2, decimal.NewFromFloat(0)},
		{"hex", "0x0100", 2, decimal.NewFromFloat(2.56)},
		{"decimal", "10000000000000", 11, decimal.NewFromInt(100)},
		{"large decimal", "52050000000000000000", 11, decimal.RequireFromString("520500000")},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			caller := new(mocks.CallerSubscriber)
			ethClient := &eth.CallerSubscriberClient{CallerSubscriber: caller}

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

func TestCallerSubscriberClient_GetAggregatorLatestRound(t *testing.T) {
	address := cltest.NewAddress()

	const aggregatorLatestRoundID = "668a0f02"
	aggregatorLatestRoundSelector := eth.HexToFunctionSelector(aggregatorLatestRoundID)

	expectedCallArgs := eth.CallArgs{
		To:   address,
		Data: aggregatorLatestRoundSelector.Bytes(),
	}
	large, ok := new(big.Int).SetString("52050000000000000000", 10)
	require.True(t, ok)

	tests := []struct {
		name, response string
		expectation    *big.Int
	}{
		{"zero", "0", big.NewInt(0)},
		{"small", "12", big.NewInt(12)},
		{"large", "52050000000000000000", large},
		{"hex zero default", "0x", big.NewInt(0)},
		{"hex zero", "0x0", big.NewInt(0)},
		{"hex", "0x0100", big.NewInt(256)},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			caller := new(mocks.CallerSubscriber)
			ethClient := &eth.CallerSubscriberClient{CallerSubscriber: caller}

			caller.On("Call", mock.Anything, "eth_call", expectedCallArgs, "latest").Return(nil).
				Run(func(args mock.Arguments) {
					res := args.Get(0).(*string)
					*res = test.response
				})
			result, err := ethClient.GetAggregatorLatestRound(address)
			require.NoError(t, err)
			assert.Equal(t, test.expectation, result)
			caller.AssertExpectations(t)
		})
	}
}

func TestCallerSubscriberClient_GetAggregatorReportingRound(t *testing.T) {
	address := cltest.NewAddress()

	const aggregatorReportingRoundID = "6fb4bb4e"
	aggregatorReportingRoundSelector := eth.HexToFunctionSelector(aggregatorReportingRoundID)

	expectedCallArgs := eth.CallArgs{
		To:   address,
		Data: aggregatorReportingRoundSelector.Bytes(),
	}
	large, ok := new(big.Int).SetString("52050000000000000000", 10)
	require.True(t, ok)

	tests := []struct {
		name, response string
		expectation    *big.Int
	}{
		{"zero", "0", big.NewInt(0)},
		{"small", "12", big.NewInt(12)},
		{"large", "52050000000000000000", large},
		{"hex zero default", "0x", big.NewInt(0)},
		{"hex zero", "0x0", big.NewInt(0)},
		{"hex", "0x0100", big.NewInt(256)},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			caller := new(mocks.CallerSubscriber)
			ethClient := &eth.CallerSubscriberClient{CallerSubscriber: caller}

			caller.On("Call", mock.Anything, "eth_call", expectedCallArgs, "latest").Return(nil).
				Run(func(args mock.Arguments) {
					res := args.Get(0).(*string)
					*res = test.response
				})
			result, err := ethClient.GetAggregatorReportingRound(address)
			require.NoError(t, err)
			assert.Equal(t, test.expectation, result)
			caller.AssertExpectations(t)
		})
	}
}

func TestCallerSubscriberClient_GetAggregatorTimeout(t *testing.T) {
	address := cltest.NewAddress()

	const aggregatorTimeoutID = "70dea79a"
	aggregatorTimeoutSelector := eth.HexToFunctionSelector(aggregatorTimeoutID)

	expectedCallArgs := eth.CallArgs{
		To:   address,
		Data: aggregatorTimeoutSelector.Bytes(),
	}
	large, ok := new(big.Int).SetString("52050000000000000000", 10)
	require.True(t, ok)

	tests := []struct {
		name, response string
		expectation    *big.Int
	}{
		{"zero", "0", big.NewInt(0)},
		{"small", "12", big.NewInt(12)},
		{"large", "52050000000000000000", large},
		{"hex zero default", "0x", big.NewInt(0)},
		{"hex zero", "0x0", big.NewInt(0)},
		{"hex", "0x0100", big.NewInt(256)},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			caller := new(mocks.CallerSubscriber)
			ethClient := &eth.CallerSubscriberClient{CallerSubscriber: caller}

			caller.On("Call", mock.Anything, "eth_call", expectedCallArgs, "latest").Return(nil).
				Run(func(args mock.Arguments) {
					res := args.Get(0).(*string)
					*res = test.response
				})
			result, err := ethClient.GetAggregatorTimeout(address)
			require.NoError(t, err)
			assert.Equal(t, test.expectation, result)
			caller.AssertExpectations(t)
		})
	}
}

func TestCallerSubscriberClient_GetAggregatorTimedOutStatus(t *testing.T) {
	const aggregatorTimedOutStatusID = "25b6ae00"
	address := cltest.NewAddress()
	aggregatorTimedOutStatusSelector := eth.HexToFunctionSelector(aggregatorTimedOutStatusID)
	roundBytes := common.Hex2BytesFixed(hexutil.EncodeUint64(0), 32)
	callData := utils.ConcatBytes(aggregatorTimedOutStatusSelector.Bytes(), roundBytes)

	expectedCallArgs := eth.CallArgs{
		To:   address,
		Data: callData,
	}

	tests := []struct {
		name        string
		response    bool
		expectation bool
	}{
		{"true", true, true},
		{"false", false, false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			caller := new(mocks.CallerSubscriber)
			ethClient := &eth.CallerSubscriberClient{CallerSubscriber: caller}

			caller.On("Call", mock.Anything, "eth_call", expectedCallArgs, "latest").Return(nil).
				Run(func(args mock.Arguments) {
					res := args.Get(0).(*bool)
					*res = test.response
				})
			result, err := ethClient.GetAggregatorTimedOutStatus(address, big.NewInt(0))
			require.NoError(t, err)
			assert.Equal(t, test.expectation, result)
			caller.AssertExpectations(t)
		})
	}
}

func TestCallerSubscriberClient_GetAggregatorLatestSubmission(t *testing.T) {
	caller := new(mocks.CallerSubscriber)
	ethClient := &eth.CallerSubscriberClient{CallerSubscriber: caller}
	aggregatorAddress := cltest.NewAddress()
	oracleAddress := cltest.NewAddress()

	const aggregatorLatestSubmission = "bb07bacd"
	aggregatorLatestSubmissionSelector := eth.HexToFunctionSelector(aggregatorLatestSubmission)

	callData := utils.ConcatBytes(aggregatorLatestSubmissionSelector.Bytes(), oracleAddress.Hash().Bytes())

	expectedCallArgs := eth.CallArgs{
		To:   aggregatorAddress,
		Data: callData,
	}

	tests := []struct {
		name           string
		answer         int64
		round          int64
		expectedAnswer *big.Int
		expectedRound  *big.Int
	}{
		{"zero", 0, 0, big.NewInt(0), big.NewInt(0)},
		{"small", 8, 12, big.NewInt(8), big.NewInt(12)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			caller.On("Call", mock.Anything, "eth_call", expectedCallArgs, "latest").Return(nil).
				Run(func(args mock.Arguments) {
					res := args.Get(0).(*string)
					answerBytes, err := utils.EVMWordSignedBigInt(big.NewInt(test.answer))
					require.NoError(t, err)
					roundBytes, err := utils.EVMWordBigInt(big.NewInt(test.round))
					require.NoError(t, err)
					*res = hexutil.Encode(append(answerBytes, roundBytes...))
				})
			answer, round, err := ethClient.GetAggregatorLatestSubmission(aggregatorAddress, oracleAddress)
			require.NoError(t, err)
			assert.Equal(t, test.expectedAnswer.String(), answer.String())
			assert.Equal(t, test.expectedRound.String(), round.String())
			caller.AssertExpectations(t)
		})
	}
}
