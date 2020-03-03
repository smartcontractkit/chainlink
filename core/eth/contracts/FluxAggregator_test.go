package contracts_test

import (
	"testing"

	"chainlink/core/eth"
	"chainlink/core/internal/cltest"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

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
