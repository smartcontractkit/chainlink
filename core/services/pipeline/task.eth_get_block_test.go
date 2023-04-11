package pipeline_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	htmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/mocks"
	evmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func Test_ETHGetBlockTask(t *testing.T) {
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {})

	lggr := logger.TestLogger(t)
	var vars pipeline.Vars
	var inputs []pipeline.Result

	h := evmtypes.Head{
		Number:           testutils.NewRandomPositiveInt64(),
		Hash:             utils.NewHash(),
		ParentHash:       utils.NewHash(),
		Timestamp:        time.Now(),
		BaseFeePerGas:    assets.NewWeiI(testutils.NewRandomPositiveInt64()),
		ReceiptsRoot:     utils.NewHash(),
		TransactionsRoot: utils.NewHash(),
		StateRoot:        utils.NewHash(),
	}

	t.Run("returns head from headtracker if present", func(t *testing.T) {
		headTracker := htmocks.NewHeadTracker(t)
		chain := evmmocks.NewChain(t)
		chain.On("HeadTracker").Return(headTracker)

		cc := evmtest.NewMockChainSetWithChain(t, chain)

		task := pipeline.ETHGetBlockTask{}
		task.HelperSetDependencies(cc, cfg)

		headTracker.On("LatestChain").Return(&h, nil)
		res, ri := task.Run(testutils.Context(t), lggr, vars, inputs)

		assert.Nil(t, res.Error)
		hVal, is := res.Value.(map[string]interface{})
		require.True(t, is, "expected %T to be map[string]interface{}", res.Value)
		assert.Equal(t, h.Number, hVal["number"])
		assert.Equal(t, h.Hash, hVal["hash"])
		assert.Equal(t, h.ParentHash, hVal["parentHash"])
		assert.Equal(t, h.Timestamp, hVal["timestamp"])
		assert.Equal(t, h.BaseFeePerGas, hVal["baseFeePerGas"])
		assert.Equal(t, h.ReceiptsRoot, hVal["receiptsRoot"])
		assert.Equal(t, h.TransactionsRoot, hVal["transactionsRoot"])
		assert.Equal(t, h.StateRoot, hVal["stateRoot"])
		assert.Equal(t, pipeline.RunInfo(pipeline.RunInfo{IsRetryable: false, IsPending: false}), ri)

		chain.AssertExpectations(t)
		headTracker.AssertExpectations(t)
	})

	t.Run("if headtracker returns nil head and eth call succeeds", func(t *testing.T) {
		ethClient := evmclimocks.NewClient(t)
		headTracker := htmocks.NewHeadTracker(t)
		chain := evmmocks.NewChain(t)
		chain.On("Client").Return(ethClient)
		chain.On("HeadTracker").Return(headTracker)

		cc := evmtest.NewMockChainSetWithChain(t, chain)

		task := pipeline.ETHGetBlockTask{}
		task.HelperSetDependencies(cc, cfg)

		// This can happen in some cases e.g. RPC node is offline
		headTracker.On("LatestChain").Return(nil)
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(&h, nil)

		res, ri := task.Run(testutils.Context(t), lggr, vars, inputs)

		assert.Nil(t, res.Error)
		hVal, is := res.Value.(map[string]interface{})
		require.True(t, is, "expected %T to be map[string]interface{}", res.Value)
		assert.Equal(t, h.Number, hVal["number"])
		assert.Equal(t, h.Hash, hVal["hash"])
		assert.Equal(t, h.ParentHash, hVal["parentHash"])
		assert.Equal(t, h.Timestamp, hVal["timestamp"])
		assert.Equal(t, h.BaseFeePerGas, hVal["baseFeePerGas"])
		assert.Equal(t, h.ReceiptsRoot, hVal["receiptsRoot"])
		assert.Equal(t, h.TransactionsRoot, hVal["transactionsRoot"])
		assert.Equal(t, h.StateRoot, hVal["stateRoot"])
		assert.Equal(t, pipeline.RunInfo(pipeline.RunInfo{IsRetryable: false, IsPending: false}), ri)

		chain.AssertExpectations(t)
		ethClient.AssertExpectations(t)
		headTracker.AssertExpectations(t)
	})

	t.Run("if headtracker returns nil head and eth call fails", func(t *testing.T) {
		ethClient := evmclimocks.NewClient(t)
		headTracker := htmocks.NewHeadTracker(t)
		chain := evmmocks.NewChain(t)
		chain.On("Client").Return(ethClient)
		chain.On("HeadTracker").Return(headTracker)

		cc := evmtest.NewMockChainSetWithChain(t, chain)

		task := pipeline.ETHGetBlockTask{}
		task.HelperSetDependencies(cc, cfg)

		// This can happen in some cases e.g. RPC node is offline
		headTracker.On("LatestChain").Return(nil)
		err := errors.New("foo")
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(nil, err)
		res, ri := task.Run(testutils.Context(t), lggr, vars, inputs)

		assert.Equal(t, pipeline.Result(pipeline.Result{Value: interface{}(nil), Error: err}), res)
		assert.Equal(t, pipeline.RunInfo(pipeline.RunInfo{IsRetryable: false, IsPending: false}), ri)

		chain.AssertExpectations(t)
		ethClient.AssertExpectations(t)
		headTracker.AssertExpectations(t)
	})
}
