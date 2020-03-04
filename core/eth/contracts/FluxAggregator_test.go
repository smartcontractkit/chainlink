package contracts_test

import (
	"encoding"
	"fmt"
	"math/big"
	"testing"

	"chainlink/core/eth"
	"chainlink/core/eth/contracts"
	"chainlink/core/internal/cltest"
	"chainlink/core/internal/mocks"
	"chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFluxAggregatorClient_LatestAnswer(t *testing.T) {
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
		{"hex - Zero", cltest.MustEVMUintHexFromBase10String(t, "0"), 2, decimal.NewFromFloat(0)},
		{"hex", cltest.MustEVMUintHexFromBase10String(t, "256"), 2, decimal.NewFromFloat(2.56)},
		{"decimal", cltest.MustEVMUintHexFromBase10String(t, "10000000000000"), 11, decimal.NewFromInt(100)},
		{"large decimal", cltest.MustEVMUintHexFromBase10String(t, "52050000000000000000"), 11, decimal.RequireFromString("520500000")},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ethClient := new(mocks.Client)

			ethClient.On("Call", mock.Anything, "eth_call", expectedCallArgs, "latest").Return(nil).
				Run(func(args mock.Arguments) {
					res := args.Get(0)
					err := res.(encoding.TextUnmarshaler).UnmarshalText([]byte(test.response))
					require.NoError(t, err)
				})

			fa, err := contracts.NewFluxAggregator(ethClient, address)
			require.NoError(t, err)

			result, err := fa.LatestAnswer(test.precision)
			require.NoError(t, err)
			assert.True(t, test.expectation.Equal(result))
			ethClient.AssertExpectations(t)
		})
	}
}

func TestFluxAggregatorClient_LatestRound(t *testing.T) {
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
		{"zero", cltest.MustEVMUintHexFromBase10String(t, "0"), big.NewInt(0)},
		{"small", cltest.MustEVMUintHexFromBase10String(t, "12"), big.NewInt(12)},
		{"large", cltest.MustEVMUintHexFromBase10String(t, "52050000000000000000"), large},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ethClient := new(mocks.Client)

			ethClient.On("Call", mock.Anything, "eth_call", expectedCallArgs, "latest").Return(nil).
				Run(func(args mock.Arguments) {
					res := args.Get(0)
					err := res.(encoding.TextUnmarshaler).UnmarshalText([]byte(test.response))
					require.NoError(t, err)
				})

			fa, err := contracts.NewFluxAggregator(ethClient, address)
			require.NoError(t, err)

			result, err := fa.LatestRound()
			require.NoError(t, err)
			require.True(t, test.expectation.Cmp(result) == 0)
			ethClient.AssertExpectations(t)
		})
	}
}

func TestFluxAggregatorClient_ReportingRound(t *testing.T) {
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
		{"zero", cltest.MustEVMUintHexFromBase10String(t, "0"), big.NewInt(0)},
		{"small", cltest.MustEVMUintHexFromBase10String(t, "12"), big.NewInt(12)},
		{"large", cltest.MustEVMUintHexFromBase10String(t, "52050000000000000000"), large},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ethClient := new(mocks.Client)

			ethClient.On("Call", mock.Anything, "eth_call", expectedCallArgs, "latest").Return(nil).
				Run(func(args mock.Arguments) {
					res := args.Get(0)
					err := res.(encoding.TextUnmarshaler).UnmarshalText([]byte(test.response))
					require.NoError(t, err)
				})

			fa, err := contracts.NewFluxAggregator(ethClient, address)
			require.NoError(t, err)

			result, err := fa.ReportingRound()
			require.NoError(t, err)
			require.True(t, test.expectation.Cmp(result) == 0)
			ethClient.AssertExpectations(t)
		})
	}
}

func TestFluxAggregatorClient_TimedOutStatus(t *testing.T) {
	const aggregatorTimedOutStatusID = "25b6ae00"
	address := cltest.NewAddress()
	aggregatorTimedOutStatusSelector := eth.HexToFunctionSelector(aggregatorTimedOutStatusID)
	roundBytes := common.Hex2BytesFixed(hexutil.EncodeUint64(0), 32)
	callData := utils.ConcatBytes(aggregatorTimedOutStatusSelector.Bytes(), roundBytes)

	expectedCallArgs := eth.CallArgs{
		To:   address,
		Data: callData,
	}

	var evmFalse = "0x0000000000000000000000000000000000000000000000000000000000000000"
	var evmTrue = "0x0000000000000000000000000000000000000000000000000000000000000001"

	tests := []struct {
		name        string
		response    string
		expectation bool
	}{
		{"true", evmTrue, true},
		{"false", evmFalse, false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ethClient := new(mocks.Client)

			ethClient.On("Call", mock.Anything, "eth_call", expectedCallArgs, "latest").Return(nil).
				Run(func(args mock.Arguments) {
					res := args.Get(0)
					err := res.(encoding.TextUnmarshaler).UnmarshalText([]byte(test.response))
					require.NoError(t, err)
				})

			fa, err := contracts.NewFluxAggregator(ethClient, address)
			require.NoError(t, err)

			result, err := fa.TimedOutStatus(big.NewInt(0))
			require.NoError(t, err)
			assert.Equal(t, test.expectation, result)
			ethClient.AssertExpectations(t)
		})
	}
}

func TestFluxAggregatorClient_LatestSubmission(t *testing.T) {
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
		response       string
		expectedAnswer *big.Int
		expectedRound  *big.Int
	}{
		{"zero", cltest.MustEVMUintHexFromBase10String(t, "0", "0"), big.NewInt(0), big.NewInt(0)},
		{"small", cltest.MustEVMUintHexFromBase10String(t, "8", "12"), big.NewInt(8), big.NewInt(12)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ethClient := new(mocks.Client)

			ethClient.On("Call", mock.Anything, "eth_call", expectedCallArgs, "latest").Return(nil).
				Run(func(args mock.Arguments) {
					res := args.Get(0)
					fmt.Println("ZZZZ ~>", test.response)
					err := res.(encoding.TextUnmarshaler).UnmarshalText([]byte(test.response))
					require.NoError(t, err)
				})

			fa, err := contracts.NewFluxAggregator(ethClient, aggregatorAddress)
			require.NoError(t, err)

			answer, round, err := fa.LatestSubmission(oracleAddress)
			require.NoError(t, err)
			assert.Equal(t, test.expectedAnswer.String(), answer.String())
			assert.Equal(t, test.expectedRound.String(), round.String())
			ethClient.AssertExpectations(t)
		})
	}
}
