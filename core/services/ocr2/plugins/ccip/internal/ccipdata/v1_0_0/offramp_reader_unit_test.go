package v1_0_0

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	mock_contracts "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/mocks/v1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib/rpclibmocks"
)

func TestOffRampGetDestinationTokensFromSourceTokens(t *testing.T) {
	ctx := testutils.Context(t)
	const numSrcTokens = 20

	testCases := []struct {
		name           string
		outputChangeFn func(outputs []rpclib.DataAndErr) []rpclib.DataAndErr
		expErr         bool
	}{
		{
			name:           "happy path",
			outputChangeFn: func(outputs []rpclib.DataAndErr) []rpclib.DataAndErr { return outputs },
			expErr:         false,
		},
		{
			name: "rpc error",
			outputChangeFn: func(outputs []rpclib.DataAndErr) []rpclib.DataAndErr {
				outputs[2].Err = fmt.Errorf("some error")
				return outputs
			},
			expErr: true,
		},
		{
			name: "unexpected outputs length should be fine if the type is correct",
			outputChangeFn: func(outputs []rpclib.DataAndErr) []rpclib.DataAndErr {
				outputs[0].Outputs = append(outputs[0].Outputs, "unexpected", 123)
				return outputs
			},
			expErr: false,
		},
		{
			name: "different compatible type",
			outputChangeFn: func(outputs []rpclib.DataAndErr) []rpclib.DataAndErr {
				outputs[0].Outputs = []any{outputs[0].Outputs[0].(common.Address)}
				return outputs
			},
			expErr: false,
		},
		{
			name: "different incompatible type",
			outputChangeFn: func(outputs []rpclib.DataAndErr) []rpclib.DataAndErr {
				outputs[0].Outputs = []any{outputs[0].Outputs[0].(common.Address).Bytes()}
				return outputs
			},
			expErr: true,
		},
	}

	lp := mocks.NewLogPoller(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			batchCaller := rpclibmocks.NewEvmBatchCaller(t)
			o := &OffRamp{evmBatchCaller: batchCaller, lp: lp}
			srcTks, dstTks, outputs := generateTokensAndOutputs(numSrcTokens)
			outputs = tc.outputChangeFn(outputs)
			batchCaller.On("BatchCall", mock.Anything, mock.Anything, mock.Anything).Return(outputs, nil)
			genericAddrs := ccipcalc.EvmAddrsToGeneric(srcTks...)
			actualDstTokens, err := o.getDestinationTokensFromSourceTokens(ctx, genericAddrs)

			if tc.expErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, ccipcalc.EvmAddrsToGeneric(dstTks...), actualDstTokens)
		})
	}
}

func TestCachedOffRampTokens(t *testing.T) {
	// Test data.
	srcTks, dstTks, outputs := generateTokensAndOutputs(3)

	// Mock contract wrapper.
	mockOffRamp := mock_contracts.NewEVM2EVMOffRampInterface(t)
	mockOffRamp.On("GetDestinationTokens", mock.Anything).Return(dstTks, nil)
	mockOffRamp.On("GetSupportedTokens", mock.Anything).Return(srcTks, nil)
	mockOffRamp.On("Address").Return(utils.RandomAddress())

	lp := mocks.NewLogPoller(t)
	lp.On("LatestBlock", mock.Anything).Return(logpoller.LogPollerBlock{BlockNumber: rand.Int63()}, nil)

	ec := evmclimocks.NewClient(t)

	batchCaller := rpclibmocks.NewEvmBatchCaller(t)
	batchCaller.On("BatchCall", mock.Anything, mock.Anything, mock.Anything).Return(outputs, nil)

	offRamp := OffRamp{
		offRampV100:    mockOffRamp,
		lp:             lp,
		Logger:         logger.TestLogger(t),
		Client:         ec,
		evmBatchCaller: batchCaller,
		cachedOffRampTokens: cache.NewLogpollerEventsBased[cciptypes.OffRampTokens](
			lp,
			offRamp_poolAddedPoolRemovedEvents,
			mockOffRamp.Address(),
		),
	}

	ctx := testutils.Context(t)
	tokens, err := offRamp.GetTokens(ctx)
	require.NoError(t, err)

	// Verify data is properly loaded in the cache.
	expectedPools := make(map[cciptypes.Address]cciptypes.Address)
	for i := range dstTks {
		expectedPools[cciptypes.Address(dstTks[i].String())] = cciptypes.Address(dstTks[i].String())
	}
	require.Equal(t, cciptypes.OffRampTokens{
		DestinationTokens: ccipcalc.EvmAddrsToGeneric(dstTks...),
		SourceTokens:      ccipcalc.EvmAddrsToGeneric(srcTks...),
		DestinationPool:   expectedPools,
	}, tokens)
}

func generateTokensAndOutputs(nbTokens uint) ([]common.Address, []common.Address, []rpclib.DataAndErr) {
	srcTks := make([]common.Address, nbTokens)
	dstTks := make([]common.Address, nbTokens)
	outputs := make([]rpclib.DataAndErr, nbTokens)
	for i := range srcTks {
		srcTks[i] = utils.RandomAddress()
		dstTks[i] = utils.RandomAddress()
		outputs[i] = rpclib.DataAndErr{
			Outputs: []any{dstTks[i]}, Err: nil,
		}
	}
	return srcTks, dstTks, outputs
}
