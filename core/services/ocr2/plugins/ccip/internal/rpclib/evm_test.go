package rpclib_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib"
)

func TestDefaultEvmBatchCaller_BatchCallDynamicLimit(t *testing.T) {
	testCases := []struct {
		name                          string
		maxBatchSize                  uint
		backOffMultiplier             uint
		numCalls                      int
		expectedBatchSizesOnEachRetry []int
	}{
		{
			name:                          "defaults",
			maxBatchSize:                  rpclib.DefaultRpcBatchSizeLimit,
			backOffMultiplier:             rpclib.DefaultRpcBatchBackOffMultiplier,
			numCalls:                      200,
			expectedBatchSizesOnEachRetry: []int{100, 20, 4, 1},
		},
		{
			name:                          "base simple scenario",
			maxBatchSize:                  20,
			backOffMultiplier:             2,
			numCalls:                      100,
			expectedBatchSizesOnEachRetry: []int{20, 10, 5, 2, 1},
		},
		{
			name:                          "remainder",
			maxBatchSize:                  99,
			backOffMultiplier:             5,
			numCalls:                      100,
			expectedBatchSizesOnEachRetry: []int{99, 19, 3, 1},
		},
		{
			name:                          "large back off multiplier",
			maxBatchSize:                  20,
			backOffMultiplier:             18,
			numCalls:                      100,
			expectedBatchSizesOnEachRetry: []int{20, 1},
		},
		{
			name:                          "back off equal to batch size",
			maxBatchSize:                  20,
			backOffMultiplier:             20,
			numCalls:                      100,
			expectedBatchSizesOnEachRetry: []int{20, 1},
		},
		{
			name:                          "back off larger than batch size",
			maxBatchSize:                  20,
			backOffMultiplier:             220,
			numCalls:                      100,
			expectedBatchSizesOnEachRetry: []int{20, 1},
		},
		{
			name:                          "back off 1",
			maxBatchSize:                  20,
			backOffMultiplier:             1,
			numCalls:                      100,
			expectedBatchSizesOnEachRetry: []int{20, 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			batchSizes := make([]int, 0)

			ec := mocks.NewClient(t)
			bc := rpclib.NewDynamicLimitedBatchCaller(logger.TestLogger(t), ec, tc.maxBatchSize, tc.backOffMultiplier)
			ctx := testutils.Context(t)
			calls := make([]rpclib.EvmCall, tc.numCalls)
			emptyAbi := abihelpers.MustParseABI("[]")
			for i := range calls {
				calls[i] = rpclib.NewEvmCall(emptyAbi, "", common.Address{})
			}
			ec.On("BatchCallContext", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				evmCalls := args.Get(1).([]rpc.BatchElem)
				batchSizes = append(batchSizes, len(evmCalls))
			}).Return(errors.New("some error"))
			_, _ = bc.BatchCall(ctx, 123, calls)

			assert.Equal(t, tc.expectedBatchSizesOnEachRetry, batchSizes)
		})
	}

}

func TestParseOutput(t *testing.T) {
	type testCase[T any] struct {
		name       string
		dataAndErr rpclib.DataAndErr
		outputIdx  int
		expRes     T
		expErr     bool
	}

	testCases := []testCase[string]{
		{
			name:       "success",
			dataAndErr: rpclib.DataAndErr{Outputs: []any{"abc"}, Err: nil},
			outputIdx:  0,
			expRes:     "abc",
			expErr:     false,
		},
		{
			name:       "index error on empty list",
			dataAndErr: rpclib.DataAndErr{Outputs: []any{}, Err: nil},
			outputIdx:  0,
			expErr:     true,
		},
		{
			name:       "index error on non-empty list",
			dataAndErr: rpclib.DataAndErr{Outputs: []any{"a", "b"}, Err: nil},
			outputIdx:  2,
			expErr:     true,
		},
		{
			name:       "negative index",
			dataAndErr: rpclib.DataAndErr{Outputs: []any{"a", "b"}, Err: nil},
			outputIdx:  -1,
			expErr:     true,
		},
		{
			name:       "wrong type",
			dataAndErr: rpclib.DataAndErr{Outputs: []any{1234}, Err: nil},
			outputIdx:  0,
			expErr:     true,
		},
		{
			name:       "has err",
			dataAndErr: rpclib.DataAndErr{Outputs: []any{"abc"}, Err: fmt.Errorf("some err")},
			outputIdx:  0,
			expErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := rpclib.ParseOutput[string](tc.dataAndErr, tc.outputIdx)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expRes, res)
		})
	}
}
