package evm

import (
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	commonmocks "github.com/smartcontractkit/chainlink-common/pkg/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
			maxBatchSize:                  defaultRpcBatchSizeLimit,
			backOffMultiplier:             defaultRpcBatchBackOffMultiplier,
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
			mockCodec := commonmocks.NewCodec(t)
			mockCodec.On("Encode", mock.Anything, mock.Anything, mock.Anything).Return([]byte{}, nil)

			bc := newDynamicLimitedBatchCaller(logger.TestLogger(t), mockCodec, ec, tc.maxBatchSize, tc.backOffMultiplier, 1)

			calls := make(BatchCall, tc.numCalls)
			for i := range calls {
				calls[i] = Call{}
			}

			ec.On("BatchCallContext", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				evmCalls := args.Get(1).([]rpc.BatchElem)
				batchSizes = append(batchSizes, len(evmCalls))
			}).Return(errors.New("some error"))

			ctx := testutils.Context(t)
			_, _ = bc.BatchCall(ctx, 123, calls)
			assert.Equal(t, tc.expectedBatchSizesOnEachRetry, batchSizes)
		})
	}

}

//func TestDefaultEvmBatchCaller_batchCallLimit(t *testing.T) {
//	ctx := testutils.Context(t)
//
//	testCases := []struct {
//		numCalls              uint
//		batchSize             uint
//		parallelRpcCallsLimit uint
//	}{
//		{numCalls: 100, batchSize: 10, parallelRpcCallsLimit: 5},
//		{numCalls: 10, batchSize: 100, parallelRpcCallsLimit: 10},
//		{numCalls: 1, batchSize: 100, parallelRpcCallsLimit: 10},
//		{numCalls: 1000, batchSize: 10, parallelRpcCallsLimit: 2},
//		{numCalls: rand.Uint() % 1000, batchSize: rand.Uint() % 500, parallelRpcCallsLimit: rand.Uint() % 500},
//	}
//
//	for _, tc := range testCases {
//		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
//			ec := mocks.NewClient(t)
//			bc := rpclib.NewDynamicLimitedBatchCaller(logger.TestLogger(t), ec, tc.batchSize, 99999, tc.parallelRpcCallsLimit)
//
//			// generate the abi and the rpc calls
//			intTyp, err := abi.NewType("uint64", "uint64", nil)
//			assert.NoError(t, err)
//			calls := make([]rpclib.EvmCall, tc.numCalls)
//			mockAbi := abihelpers.MustParseABI("[]")
//			for i := range calls {
//				name := fmt.Sprintf("method_%d", i)
//				meth := abi.NewMethod(name, name, abi.Function, "nonpayable", true, false, abi.Arguments{abi.Argument{Name: "a", Type: intTyp}}, abi.Arguments{abi.Argument{Name: "b", Type: intTyp}})
//				mockAbi.Methods[name] = meth
//				calls[i] = rpclib.NewEvmCall(mockAbi, name, common.Address{}, uint64(i))
//			}
//
//			// mock the rpc call to batch call context
//			// for simplicity we just set an error
//			ec.On("BatchCallContext", mock.Anything, mock.Anything).
//				Run(func(args mock.Arguments) {
//					evmCalls := args.Get(1).([]rpc.BatchElem)
//					for i := range evmCalls {
//						arg := evmCalls[i].Args[0].(map[string]interface{})["data"].(hexutil.Bytes)
//						arg = arg[len(arg)-10:]
//						evmCalls[i].Error = fmt.Errorf("%s", arg)
//					}
//				}).Return(nil)
//
//			// make the call and make sure the results are received in order
//			results, _ := bc.BatchCall(ctx, 0, calls)
//			assert.Len(t, results, len(calls))
//			for i, res := range results {
//				resNum, err := strconv.ParseInt(res.Err.Error()[2:], 16, 64)
//				assert.NoError(t, err)
//				assert.Equal(t, int64(i), resNum)
//			}
//		})
//	}
//}
