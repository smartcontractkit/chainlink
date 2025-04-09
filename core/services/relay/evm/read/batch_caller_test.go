package read_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/read"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
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
			maxBatchSize:                  read.DefaultRpcBatchSizeLimit,
			backOffMultiplier:             read.DefaultRpcBatchBackOffMultiplier,
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

	mockCodec := mocks.NewCodec(t)
	mockCodec.On("Encode", mock.Anything, mock.Anything, mock.Anything).Return([]byte{}, nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			batchSizes := make([]int, 0)
			ec := clienttest.NewClient(t)
			ec.On("BatchCallContext", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				evmCalls := args.Get(1).([]rpc.BatchElem)
				batchSizes = append(batchSizes, len(evmCalls))
			}).Return(errors.New("some error"))

			calls := make(read.BatchCall, tc.numCalls)
			for i := range calls {
				calls[i] = read.Call{}
			}

			bc := read.NewDynamicLimitedBatchCaller(logger.Test(t), mockCodec, ec, tc.maxBatchSize, tc.backOffMultiplier, 1)
			_, _ = bc.BatchCall(testutils.Context(t), 123, calls)
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
	}

	type MethodParam struct {
		A uint64
	}
	type MethodReturn struct {
		B uint64
	}
	paramABI := `[{"type":"uint64","name":"A"}]`
	returnABI := `[{"type":"uint64","name":"B"}]`
	codecConfig := evmtypes.CodecConfig{Configs: map[string]evmtypes.ChainCodecConfig{}}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			ec := clienttest.NewClient(t)
			calls := make(read.BatchCall, tc.numCalls)
			for j := range calls {
				contractName := fmt.Sprintf("testCase_%d", i)
				methodName := fmt.Sprintf("method_%d", j)
				codecConfig.Configs[fmt.Sprintf("params.%s.%s", contractName, methodName)] = evmtypes.ChainCodecConfig{TypeABI: paramABI}
				codecConfig.Configs[fmt.Sprintf("return.%s.%s", contractName, methodName)] = evmtypes.ChainCodecConfig{TypeABI: returnABI}

				params := MethodParam{A: uint64(j)}
				var returnVal MethodReturn
				calls[j] = read.Call{
					ContractName: contractName,
					ReadName:     methodName,
					Params:       &params,
					ReturnVal:    &returnVal,
				}
			}

			ec.On("BatchCallContext", mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					evmCalls := args.Get(1).([]rpc.BatchElem)
					for i := range evmCalls {
						arg := evmCalls[i].Args[0].(map[string]interface{})["data"].(hexutil.Bytes)
						str, isOk := evmCalls[i].Result.(*string)
						require.True(t, isOk)
						*str = fmt.Sprintf("0x%064x", new(big.Int).SetBytes(arg[24:]).Uint64())
					}
				}).Return(nil)

			testCodec, err := codec.NewCodec(codecConfig)
			require.NoError(t, err)
			bc := read.NewDynamicLimitedBatchCaller(logger.Test(t), testCodec, ec, tc.batchSize, 99999, tc.parallelRpcCallsLimit)

			// make the call and make sure the results are there
			results, err := bc.BatchCall(ctx, 0, calls)
			require.NoError(t, err)
			for _, call := range calls {
				contractResults, ok := results[call.ContractName]
				if !ok {
					t.Errorf("missing contract name %s", call.ContractName)
				}
				hasResult := false
				for j, result := range contractResults {
					if hasResult = result.MethodName == call.ReadName; hasResult {
						require.NoError(t, result.Err)
						resNum, isOk := result.ReturnValue.(*MethodReturn)
						require.True(t, isOk)
						require.Equal(t, uint64(j), resNum.B)
						break
					}
				}
				if !hasResult {
					t.Errorf("missing method name %s", call.ReadName)
				}
			}
		})
	}
}
