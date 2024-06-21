package evm

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	commonmocks "github.com/smartcontractkit/chainlink-common/pkg/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

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

func TestDefaultEvmBatchCaller_batchCallLimit(t *testing.T) {
	ctx := testutils.Context(t)

	testCases := []struct {
		numCalls              uint
		batchSize             uint
		parallelRpcCallsLimit uint
	}{
		{numCalls: 100, batchSize: 10, parallelRpcCallsLimit: 5},
		//{numCalls: 10, batchSize: 100, parallelRpcCallsLimit: 10},
		//{numCalls: 1, batchSize: 100, parallelRpcCallsLimit: 10},
		//{numCalls: 1000, batchSize: 10, parallelRpcCallsLimit: 2},
		//{numCalls: rand.Uint() % 1000, batchSize: rand.Uint() % 500, parallelRpcCallsLimit: rand.Uint() % 500},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			ec := mocks.NewClient(t)
			tcb := NewTestCodecBuilder()

			// generate the abi and the rpc calls
			calls := make(BatchCall, tc.numCalls)
			intTyp, err := abi.NewType("uint64", "uint64", nil)
			assert.NoError(t, err)

			for j := range calls {
				contractName := fmt.Sprintf("testCase_%d", i)
				methodName := fmt.Sprintf("method_%d", j)
				meth := abi.NewMethod(methodName, methodName, abi.Function, "nonpayable", true, false, abi.Arguments{abi.Argument{Name: "", Type: intTyp}}, abi.Arguments{abi.Argument{Name: "b", Type: intTyp}})
				require.NoError(t, tcb.addEncoderDef(contractName, methodName, meth.Inputs, meth.ID, nil))
				require.NoError(t, tcb.addDecoderDef(contractName, methodName, meth.Outputs, nil))
				fmt.Println("abi method ", meth.Inputs[0].Type)
				var returnVal uint64
				params := uint64(j)
				calls[j] = Call{
					contractName: contractName,
					methodName:   methodName,
					params:       &params,
					returnVal:    returnVal,
				}
			}

			testCodec, err := tcb.toCodec()
			require.NoError(t, err)
			bc := newDynamicLimitedBatchCaller(logger.TestLogger(t), testCodec, ec, tc.batchSize, 99999, tc.parallelRpcCallsLimit)

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

			// make the call and make sure the results are there
			results, err := bc.BatchCall(ctx, 0, calls)
			require.NoError(t, err)

			fmt.Println("results ", results)
			for _, call := range calls {
				contractResults, ok := results[call.contractName]
				if !ok {
					t.Errorf("missing contract name %s", call.contractName)
				}
				hasResult := false
				for j, result := range contractResults {
					if hasResult = result.methodName == call.methodName; hasResult {
						resNum, err := strconv.ParseInt(result.err.Error()[2:], 16, 64)
						assert.NoError(t, err)
						assert.Equal(t, int64(j), resNum)
						break
					}
				}
				if !hasResult {
					t.Errorf("missing method name %s", call.methodName)
				}
			}
		})
	}
}

type TestCodecBuilder struct {
	*parsedTypes
}

func NewTestCodecBuilder() *TestCodecBuilder {
	return &TestCodecBuilder{
		parsedTypes: &parsedTypes{
			encoderDefs: make(map[string]evmtypes.CodecEntry),
			decoderDefs: make(map[string]evmtypes.CodecEntry),
		},
	}
}

func (tcb *TestCodecBuilder) addEncoderDef(contractName, itemType string, args abi.Arguments, prefix []byte, inputModifications codec.ModifiersConfig) error {
	// ABI.Pack prepends the method.ID to the encodings, we'll need the encoder to do the same.
	fmt.Println(" inputModifications ", inputModifications)
	inputMod, err := inputModifications.ToModifier(evmDecoderHooks...)
	if err != nil {
		return err
	}
	input := evmtypes.NewCodecEntry(args, prefix, inputMod)

	if err = input.Init(); err != nil {
		return err
	}

	tcb.parsedTypes.encoderDefs[wrapItemType(contractName, itemType, true)] = input
	return nil
}

func (tcb *TestCodecBuilder) addDecoderDef(contractName, itemType string, outputs abi.Arguments, outputModifications codec.ModifiersConfig) error {
	mod, err := outputModifications.ToModifier(evmDecoderHooks...)
	if err != nil {
		return err
	}
	output := evmtypes.NewCodecEntry(outputs, nil, mod)
	tcb.parsedTypes.decoderDefs[wrapItemType(contractName, itemType, false)] = output
	return output.Init()
}

func (tcb *TestCodecBuilder) testCodec() (types.Codec, error) {
	return tcb.parsedTypes.toCodec()
}
