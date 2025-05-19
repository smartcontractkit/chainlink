package rpclib_test

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
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

			ec := clienttest.NewClient(t)
			bc := rpclib.NewDynamicLimitedBatchCaller(logger.TestLogger(t), ec, tc.maxBatchSize, tc.backOffMultiplier, 1)
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

func TestDefaultEvmBatchCaller_batchCallLimit(t *testing.T) {
	ctx := testutils.Context(t)

	testCases := []struct {
		numCalls              uint
		batchSize             uint
		parallelRpcCallsLimit uint
	}{
		{numCalls: 100, batchSize: 10, parallelRpcCallsLimit: 5},
		{numCalls: 10, batchSize: 100, parallelRpcCallsLimit: 10},
		{numCalls: 1, batchSize: 100, parallelRpcCallsLimit: 10},
		{numCalls: 1000, batchSize: 10, parallelRpcCallsLimit: 2},
		{numCalls: 1 + rand.Uint()%1000, batchSize: 1 + rand.Uint()%500, parallelRpcCallsLimit: 1 + rand.Uint()%500},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			ec := clienttest.NewClient(t)
			bc := rpclib.NewDynamicLimitedBatchCaller(logger.TestLogger(t), ec, tc.batchSize, 99999, tc.parallelRpcCallsLimit)

			// generate the abi and the rpc calls
			intTyp, err := abi.NewType("uint64", "uint64", nil)
			assert.NoError(t, err)
			calls := make([]rpclib.EvmCall, tc.numCalls)
			mockAbi := abihelpers.MustParseABI("[]")
			for i := range calls {
				name := fmt.Sprintf("method_%d", i)
				meth := abi.NewMethod(name, name, abi.Function, "nonpayable", true, false, abi.Arguments{abi.Argument{Name: "a", Type: intTyp}}, abi.Arguments{abi.Argument{Name: "b", Type: intTyp}})
				mockAbi.Methods[name] = meth
				calls[i] = rpclib.NewEvmCall(mockAbi, name, common.Address{}, uint64(i))
			}

			// mock the rpc call to batch call context
			// for simplicity we just set an error
			ec.On("BatchCallContext", mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					evmCalls := args.Get(1).([]rpc.BatchElem)
					for i := range evmCalls {
						arg := evmCalls[i].Args[0].(map[string]interface{})["data"].(hexutil.Bytes)
						arg = arg[len(arg)-10:]
						evmCalls[i].Error = fmt.Errorf("%s", arg)
					}
				}).Return(nil)

			// make the call and make sure the results are received in order
			results, _ := bc.BatchCall(ctx, 0, calls)
			assert.Len(t, results, len(calls))
			for i, res := range results {
				resNum, err := strconv.ParseInt(res.Err.Error()[2:], 16, 64)
				assert.NoError(t, err)
				assert.Equal(t, int64(i), resNum)
			}
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
			dataAndErr: rpclib.DataAndErr{Outputs: []any{"abc"}, Err: errors.New("some err")},
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
