package contracts_test

import (
	"encoding"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/eth/contracts"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFluxAggregatorClient_RoundState(t *testing.T) {
	aggregatorAddress := cltest.NewAddress()

	nodeAddr := cltest.NewAddress()
	selector := make([]byte, 16)
	rsHash := utils.MustHash("oracleRoundState(address,uint32)")
	copy(selector, rsHash.Bytes()[:4])
	data := append(selector, nodeAddr[:]...)
	data = append(data, utils.EVMWordUint64(0)...)
	expectedCallArgs := eth.CallArgs{
		To:   aggregatorAddress,
		Data: data,
	}

	rawReturnData := `0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000f0000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000000f000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000110000000000000000000000000000000000000000000000000000000000000100`

	tests := []struct {
		name                   string
		response               []byte
		expectedRoundID        uint32
		expectedEligible       bool
		expectedAnswer         *big.Int
		expectedTimesOutAt     uint64
		expectedAvailableFunds uint64
		expectedPaymentAmount  uint64
	}{
		{"zero, false", cltest.MakeRoundStateReturnData(0, false, 0, 0, 0, 0, 0, 17), 0, false, big.NewInt(0), 0, 0, 0},
		{"non-zero, false", cltest.MakeRoundStateReturnData(1, false, 23, 1230, 4, 36, 72, 17), 1, false, big.NewInt(23), 1234, 36, 72},
		{"zero, true", cltest.MakeRoundStateReturnData(0, true, 0, 0, 0, 0, 0, 17), 0, true, big.NewInt(0), 0, 0, 0},
		{"non-zero true", cltest.MakeRoundStateReturnData(12, true, 91, 9870, 6, 45, 999, 17), 12, true, big.NewInt(91), 9876, 45, 999},
		{"real call data", rawReturnData, 3, true, big.NewInt(15), (22 + 15), 10, 256},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ethClient := new(mocks.Client)

			ethClient.On("Call", mock.Anything, "eth_call", expectedCallArgs, "latest").Return(nil).
				Run(func(args mock.Arguments) {
					res := args.Get(0)
					err := res.(encoding.TextUnmarshaler).UnmarshalText(test.response)
					require.NoError(t, err)
				})

			fa, err := contracts.NewFluxAggregator(aggregatorAddress, ethClient, nil)
			require.NoError(t, err)

			roundState, err := fa.RoundState(nodeAddr, 0)
			require.NoError(t, err)
			assert.Equal(t, test.expectedRoundID, roundState.ReportableRoundID)
			assert.Equal(t, test.expectedEligible, roundState.EligibleToSubmit)
			assert.True(t, test.expectedAnswer.Cmp(roundState.LatestAnswer) == 0)
			assert.Equal(t, test.expectedTimesOutAt, roundState.TimesOutAt())
			assert.Equal(t, test.expectedAvailableFunds, roundState.AvailableFunds.Uint64())
			assert.Equal(t, test.expectedPaymentAmount, roundState.PaymentAmount.Uint64())
			ethClient.AssertExpectations(t)
		})
	}
}

func TestFluxAggregatorClient_DecodesLogs(t *testing.T) {
	fa, err := contracts.NewFluxAggregator(common.Address{}, nil, nil)
	require.NoError(t, err)

	newRoundLogRaw := cltest.LogFromFixture(t, "../../testdata/new_round_log.json")
	var newRoundLog contracts.LogNewRound
	err = fa.UnpackLog(&newRoundLog, "NewRound", newRoundLogRaw)
	require.NoError(t, err)
	require.Equal(t, int64(1), newRoundLog.RoundId.Int64())
	require.Equal(t, common.HexToAddress("f17f52151ebef6c7334fad080c5704d77216b732"), newRoundLog.StartedBy)
	require.Equal(t, int64(15), newRoundLog.StartedAt.Int64())

	type BadLogNewRound struct {
		RoundID   *big.Int
		StartedBy common.Address
		StartedAt *big.Int
	}
	var badNewRoundLog BadLogNewRound
	err = fa.UnpackLog(&badNewRoundLog, "NewRound", newRoundLogRaw)
	require.Error(t, err)

	answerUpdatedLogRaw := cltest.LogFromFixture(t, "../../testdata/answer_updated_log.json")
	var answerUpdatedLog contracts.LogAnswerUpdated
	err = fa.UnpackLog(&answerUpdatedLog, "AnswerUpdated", answerUpdatedLogRaw)
	require.NoError(t, err)
	require.Equal(t, int64(1), answerUpdatedLog.Current.Int64())
	require.Equal(t, int64(2), answerUpdatedLog.RoundId.Int64())
	require.Equal(t, int64(3), answerUpdatedLog.UpdatedAt.Int64())

	type BadLogAnswerUpdated struct {
		Current   *big.Int
		RoundID   *big.Int
		Timestamp *big.Int
	}
	var badAnswerUpdatedLog BadLogAnswerUpdated
	err = fa.UnpackLog(&badAnswerUpdatedLog, "AnswerUpdated", answerUpdatedLogRaw)
	require.Error(t, err)
}
