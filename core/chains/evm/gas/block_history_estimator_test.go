package gas_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func NewEvmHash() common.Hash {
	return utils.NewHash()
}

func newConfigWithEIP1559DynamicFeesEnabled(t *testing.T) *gas.MockConfig {
	cfg := gas.NewMockConfig()
	cfg.EvmEIP1559DynamicFeesF = true
	cfg.BlockHistoryEstimatorBlockHistorySizeF = 8
	return cfg
}

func newConfigWithEIP1559DynamicFeesDisabled(t *testing.T) *gas.MockConfig {
	cfg := gas.NewMockConfig()
	cfg.EvmEIP1559DynamicFeesF = false
	cfg.BlockHistoryEstimatorBlockHistorySizeF = 8
	return cfg
}

func newBlockHistoryEstimatorWithChainID(t *testing.T, c evmclient.Client, cfg gas.Config, cid big.Int) gas.EvmEstimator {
	return gas.NewBlockHistoryEstimator(logger.TestLogger(t), c, cfg, cid)
}

func newBlockHistoryEstimator(t *testing.T, c evmclient.Client, cfg gas.Config) *gas.BlockHistoryEstimator {
	iface := newBlockHistoryEstimatorWithChainID(t, c, cfg, cltest.FixtureChainID)
	return gas.BlockHistoryEstimatorFromInterface(iface)
}

func TestBlockHistoryEstimator_Start(t *testing.T) {
	t.Parallel()

	cfg := newConfigWithEIP1559DynamicFeesEnabled(t)

	var batchSize uint32
	var blockDelay uint16
	var historySize uint16 = 2
	var percentile uint16 = 35
	minGasPrice := assets.NewWeiI(1)
	maxGasPrice := assets.NewWeiI(100)

	cfg.BlockHistoryEstimatorBatchSizeF = batchSize
	cfg.BlockHistoryEstimatorBlockDelayF = blockDelay
	cfg.BlockHistoryEstimatorBlockHistorySizeF = historySize
	cfg.BlockHistoryEstimatorTransactionPercentileF = percentile
	cfg.EvmGasLimitMultiplierF = float32(1)
	cfg.EvmMinGasPriceWeiF = minGasPrice
	cfg.EvmMaxGasPriceWeiF = maxGasPrice

	t.Run("loads initial state", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		h := &evmtypes.Head{Hash: utils.NewHash(), Number: 42, BaseFeePerGas: assets.NewWeiI(420)}
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(h, nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x2a" && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&evmtypes.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == "0x29" && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&evmtypes.Block{})
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &evmtypes.Block{
				Number: 42,
				Hash:   utils.NewHash(),
			}
			elems[1].Result = &evmtypes.Block{
				Number: 41,
				Hash:   utils.NewHash(),
			}
		}).Once()

		err := bhe.Start(testutils.Context(t))
		require.NoError(t, err)

		assert.Len(t, gas.GetRollingBlockHistory(bhe), 2)
		assert.Equal(t, int(gas.GetRollingBlockHistory(bhe)[0].Number), 41)
		assert.Equal(t, int(gas.GetRollingBlockHistory(bhe)[1].Number), 42)

		assert.Equal(t, assets.NewWeiI(420), gas.GetLatestBaseFee(bhe))
	})

	t.Run("starts and loads partial history if fetch context times out", func(t *testing.T) {
		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)

		cfg.BlockHistoryEstimatorBatchSizeF = uint32(1)
		cfg.BlockHistoryEstimatorBlockDelayF = blockDelay
		cfg.BlockHistoryEstimatorBlockHistorySizeF = historySize
		cfg.BlockHistoryEstimatorTransactionPercentileF = percentile
		cfg.EvmGasLimitMultiplierF = float32(1)
		cfg.EvmMinGasPriceWeiF = minGasPrice
		cfg.EvmEIP1559DynamicFeesF = true
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		h := &evmtypes.Head{Hash: utils.NewHash(), Number: 42, BaseFeePerGas: assets.NewWeiI(420)}
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(h, nil)
		// First succeeds (42)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == gas.Int64ToHex(42) && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&evmtypes.Block{})
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &evmtypes.Block{
				Number: 42,
				Hash:   utils.NewHash(),
			}
		}).Once()
		// Second fails (41)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == gas.Int64ToHex(41) && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&evmtypes.Block{})
		})).Return(errors.Wrap(context.DeadlineExceeded, "some error message")).Once()

		err := bhe.Start(testutils.Context(t))
		require.NoError(t, err)

		require.Len(t, gas.GetRollingBlockHistory(bhe), 1)
		assert.Equal(t, int(gas.GetRollingBlockHistory(bhe)[0].Number), 42)

		assert.Equal(t, assets.NewWeiI(420), gas.GetLatestBaseFee(bhe))
	})

	t.Run("boots even if initial batch call returns nothing", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		h := &evmtypes.Head{Hash: utils.NewHash(), Number: 42}
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(h, nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == int(historySize)
		})).Return(nil)

		err := bhe.Start(testutils.Context(t))
		require.NoError(t, err)

		// non-eip1559 block
		assert.Nil(t, gas.GetLatestBaseFee(bhe))
	})

	t.Run("starts anyway if fetching latest head fails", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(nil, errors.New("something exploded"))

		err := bhe.Start(testutils.Context(t))
		require.NoError(t, err)

		assert.Nil(t, gas.GetLatestBaseFee(bhe))

		_, _, err = bhe.GetLegacyGas(testutils.Context(t), make([]byte, 0), 100, maxGasPrice)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has not finished the first gas estimation yet, likely because a failure on start")

		_, _, err = bhe.GetDynamicFee(testutils.Context(t), 100, maxGasPrice)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has not finished the first gas estimation yet, likely because a failure on start")
	})

	t.Run("starts anyway if fetching first fetch fails, but errors on estimation", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		h := &evmtypes.Head{Hash: utils.NewHash(), Number: 42, BaseFeePerGas: assets.NewWeiI(420)}
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(h, nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.Anything).Return(errors.New("something went wrong"))

		err := bhe.Start(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, assets.NewWeiI(420), gas.GetLatestBaseFee(bhe))

		_, _, err = bhe.GetLegacyGas(testutils.Context(t), make([]byte, 0), 100, maxGasPrice)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has not finished the first gas estimation yet, likely because a failure on start")

		_, _, err = bhe.GetDynamicFee(testutils.Context(t), 100, maxGasPrice)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has not finished the first gas estimation yet, likely because a failure on start")
	})

	t.Run("returns error if main context is cancelled", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		h := &evmtypes.Head{Hash: utils.NewHash(), Number: 42, BaseFeePerGas: assets.NewWeiI(420)}
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(h, nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.Anything).Return(errors.New("this error doesn't matter"))

		ctx, cancel := context.WithCancel(testutils.Context(t))
		cancel()
		err := bhe.Start(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})

	t.Run("starts anyway even if the fetch context is cancelled due to taking longer than the MaxStartTime", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		h := &evmtypes.Head{Hash: utils.NewHash(), Number: 42, BaseFeePerGas: assets.NewWeiI(420)}
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(h, nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.Anything).Return(errors.New("this error doesn't matter")).Run(func(_ mock.Arguments) {
			time.Sleep(gas.MaxStartTime + 1*time.Second)
		})

		err := bhe.Start(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, assets.NewWeiI(420), gas.GetLatestBaseFee(bhe))

		_, _, err = bhe.GetLegacyGas(testutils.Context(t), make([]byte, 0), 100, maxGasPrice)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has not finished the first gas estimation yet, likely because a failure on start")

		_, _, err = bhe.GetDynamicFee(testutils.Context(t), 100, maxGasPrice)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has not finished the first gas estimation yet, likely because a failure on start")
	})
}

func TestBlockHistoryEstimator_OnNewLongestChain(t *testing.T) {
	cfg := newConfigWithEIP1559DynamicFeesDisabled(t)
	bhe := newBlockHistoryEstimator(t, nil, cfg)

	assert.Nil(t, gas.GetLatestBaseFee(bhe))

	// non EIP-1559 block
	h := cltest.Head(1)
	bhe.OnNewLongestChain(testutils.Context(t), h)
	assert.Nil(t, gas.GetLatestBaseFee(bhe))

	// EIP-1559 block
	h = cltest.Head(2)
	h.BaseFeePerGas = assets.NewWeiI(500)
	bhe.OnNewLongestChain(testutils.Context(t), h)

	assert.Equal(t, assets.NewWeiI(500), gas.GetLatestBaseFee(bhe))
}

func TestBlockHistoryEstimator_FetchBlocks(t *testing.T) {
	t.Parallel()

	t.Run("with history size of 0, errors", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)
		var blockDelay uint16 = 3
		var historySize uint16
		cfg.BlockHistoryEstimatorBlockDelayF = blockDelay
		cfg.BlockHistoryEstimatorBlockHistorySizeF = historySize

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		head := cltest.Head(42)
		err := bhe.FetchBlocks(testutils.Context(t), head)
		require.Error(t, err)
		require.EqualError(t, err, "BlockHistoryEstimator: history size must be > 0, got: 0")
	})

	t.Run("with current block height less than block delay does nothing", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)
		var blockDelay uint16 = 3
		var historySize uint16 = 1
		cfg.BlockHistoryEstimatorBlockDelayF = blockDelay
		cfg.BlockHistoryEstimatorBlockHistorySizeF = historySize

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		for i := -1; i < 3; i++ {
			head := cltest.Head(i)
			err := bhe.FetchBlocks(testutils.Context(t), head)
			require.Error(t, err)
			require.EqualError(t, err, fmt.Sprintf("BlockHistoryEstimator: cannot fetch, current block height %v is lower than EVM.RPCBlockQueryDelay=3", i))
		}
	})

	t.Run("with error retrieving blocks returns error", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)
		var blockDelay uint16 = 3
		var historySize uint16 = 3
		var batchSize uint32
		cfg.BlockHistoryEstimatorBlockDelayF = blockDelay
		cfg.BlockHistoryEstimatorBlockHistorySizeF = historySize
		cfg.BlockHistoryEstimatorBatchSizeF = batchSize

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		ethClient.On("BatchCallContext", mock.Anything, mock.Anything).Return(errors.New("something exploded"))

		err := bhe.FetchBlocks(testutils.Context(t), cltest.Head(42))
		require.Error(t, err)
		assert.EqualError(t, err, "BlockHistoryEstimator#fetchBlocks error fetching blocks with BatchCallContext: something exploded")
	})

	t.Run("batch fetches heads and transactions and sets them on the block history estimator instance", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)
		var blockDelay uint16
		var historySize uint16 = 3
		var batchSize uint32 = 2
		cfg.BlockHistoryEstimatorBlockDelayF = blockDelay
		cfg.BlockHistoryEstimatorBlockHistorySizeF = historySize
		// Test batching
		cfg.BlockHistoryEstimatorBatchSizeF = batchSize

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		b41 := evmtypes.Block{
			Number:       41,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(1, 2),
		}
		b42 := evmtypes.Block{
			Number:       42,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(3),
		}
		b43 := evmtypes.Block{
			Number:       43,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(),
		}

		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == gas.Int64ToHex(43) && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&evmtypes.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == gas.Int64ToHex(42) && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&evmtypes.Block{})
		})).Once().Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &b43
			// This errored block (42) will be ignored
			elems[1].Error = errors.New("something went wrong")
		})
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == gas.Int64ToHex(41) && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&evmtypes.Block{})
		})).Once().Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &b41
		})

		err := bhe.FetchBlocks(testutils.Context(t), cltest.Head(43))
		require.NoError(t, err)

		require.Len(t, gas.GetRollingBlockHistory(bhe), 2)
		assert.Equal(t, 41, int(gas.GetRollingBlockHistory(bhe)[0].Number))
		// 42 is missing because the fetch errored
		assert.Equal(t, 43, int(gas.GetRollingBlockHistory(bhe)[1].Number))
		assert.Len(t, gas.GetRollingBlockHistory(bhe)[0].Transactions, 2)
		assert.Len(t, gas.GetRollingBlockHistory(bhe)[1].Transactions, 0)

		// On new fetch, rolls over the history and drops the old heads

		b44 := evmtypes.Block{
			Number:       44,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(4),
		}

		// We are gonna refetch blocks 42 and 44
		// 43 is skipped because it was already in the history
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == gas.Int64ToHex(44) && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&evmtypes.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == gas.Int64ToHex(42) && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&evmtypes.Block{})
		})).Once().Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &b44
			elems[1].Result = &b42
		})

		head := evmtypes.NewHead(big.NewInt(44), b44.Hash, b43.Hash, uint64(time.Now().Unix()), utils.NewBig(&cltest.FixtureChainID))
		err = bhe.FetchBlocks(testutils.Context(t), &head)
		require.NoError(t, err)

		require.Len(t, gas.GetRollingBlockHistory(bhe), 3)
		assert.Equal(t, 42, int(gas.GetRollingBlockHistory(bhe)[0].Number))
		assert.Equal(t, 43, int(gas.GetRollingBlockHistory(bhe)[1].Number))
		assert.Equal(t, 44, int(gas.GetRollingBlockHistory(bhe)[2].Number))
		assert.Len(t, gas.GetRollingBlockHistory(bhe)[0].Transactions, 1)
		assert.Len(t, gas.GetRollingBlockHistory(bhe)[1].Transactions, 0)
		assert.Len(t, gas.GetRollingBlockHistory(bhe)[2].Transactions, 1)
	})

	t.Run("does not refetch blocks below EVM.FinalityDepth", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)
		var blockDelay uint16
		var historySize uint16 = 3
		var batchSize uint32 = 2
		cfg.BlockHistoryEstimatorBlockDelayF = blockDelay
		cfg.BlockHistoryEstimatorBlockHistorySizeF = historySize
		cfg.BlockHistoryEstimatorBatchSizeF = batchSize

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		b0 := evmtypes.Block{
			Number:       0,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(9001),
		}
		b1 := evmtypes.Block{
			Number:       1,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(9002),
		}
		blocks := []evmtypes.Block{b0, b1}

		gas.SetRollingBlockHistory(bhe, blocks)

		b2 := evmtypes.Block{
			Number:       2,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(1, 2),
		}
		b3 := evmtypes.Block{
			Number:       3,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(1, 2),
		}

		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == gas.Int64ToHex(3) && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&evmtypes.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == gas.Int64ToHex(2) && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&evmtypes.Block{})
		})).Once().Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &b3
			elems[1].Result = &b2
		})

		head2 := evmtypes.NewHead(big.NewInt(2), b2.Hash, b1.Hash, uint64(time.Now().Unix()), utils.NewBig(&cltest.FixtureChainID))
		head3 := evmtypes.NewHead(big.NewInt(3), b3.Hash, b2.Hash, uint64(time.Now().Unix()), utils.NewBig(&cltest.FixtureChainID))
		head3.Parent = &head2
		err := bhe.FetchBlocks(testutils.Context(t), &head3)
		require.NoError(t, err)

		require.Len(t, gas.GetRollingBlockHistory(bhe), 3)
		assert.Equal(t, 1, int(gas.GetRollingBlockHistory(bhe)[0].Number))
		assert.Equal(t, 2, int(gas.GetRollingBlockHistory(bhe)[1].Number))
		assert.Equal(t, 3, int(gas.GetRollingBlockHistory(bhe)[2].Number))
	})

	t.Run("replaces blocks on re-org within EVM.FinalityDepth", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)
		var blockDelay uint16
		var historySize uint16 = 3
		var batchSize uint32 = 2
		cfg.BlockHistoryEstimatorBlockDelayF = blockDelay
		cfg.BlockHistoryEstimatorBlockHistorySizeF = historySize
		cfg.BlockHistoryEstimatorBatchSizeF = batchSize

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		b0 := evmtypes.Block{
			Number:       0,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(9001),
		}
		b1 := evmtypes.Block{
			Number:       1,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(9002),
		}
		b2 := evmtypes.Block{
			Number:       2,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(1, 2),
		}
		b3 := evmtypes.Block{
			Number:       3,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(1, 2),
		}
		blocks := []evmtypes.Block{b0, b1, b2, b3}

		gas.SetRollingBlockHistory(bhe, blocks)

		// RE-ORG, head2 and head3 have different hash than saved b2 and b3
		head2 := evmtypes.NewHead(big.NewInt(2), utils.NewHash(), b1.Hash, uint64(time.Now().Unix()), utils.NewBig(&cltest.FixtureChainID))
		head3 := evmtypes.NewHead(big.NewInt(3), utils.NewHash(), head2.Hash, uint64(time.Now().Unix()), utils.NewBig(&cltest.FixtureChainID))
		head3.Parent = &head2

		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == gas.Int64ToHex(3) && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&evmtypes.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == gas.Int64ToHex(2) && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&evmtypes.Block{})
		})).Once().Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			b2New := b2
			b2New.Hash = head2.Hash
			elems[1].Result = &b2New
			b3New := b3
			b3New.Hash = head3.Hash
			elems[0].Result = &b3New
		})

		err := bhe.FetchBlocks(testutils.Context(t), &head3)
		require.NoError(t, err)

		require.Len(t, gas.GetRollingBlockHistory(bhe), 3)
		assert.Equal(t, 1, int(gas.GetRollingBlockHistory(bhe)[0].Number))
		assert.Equal(t, 2, int(gas.GetRollingBlockHistory(bhe)[1].Number))
		assert.Equal(t, 3, int(gas.GetRollingBlockHistory(bhe)[2].Number))
		assert.Equal(t, b1.Hash.Hex(), gas.GetRollingBlockHistory(bhe)[0].Hash.Hex())
		assert.Equal(t, head2.Hash.Hex(), gas.GetRollingBlockHistory(bhe)[1].Hash.Hex())
		assert.Equal(t, head3.Hash.Hex(), gas.GetRollingBlockHistory(bhe)[2].Hash.Hex())
	})

	t.Run("uses locally cached blocks if they are in the chain", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)
		var blockDelay uint16
		var historySize uint16 = 3
		var batchSize uint32 = 2
		cfg.BlockHistoryEstimatorBlockDelayF = blockDelay
		cfg.BlockHistoryEstimatorBlockHistorySizeF = historySize
		cfg.BlockHistoryEstimatorBatchSizeF = batchSize

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		b0 := evmtypes.Block{
			Number:       0,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(9001),
		}
		b1 := evmtypes.Block{
			Number:       1,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(9002),
		}
		b2 := evmtypes.Block{
			Number:       2,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(1, 2),
		}
		b3 := evmtypes.Block{
			Number:       3,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(1, 2),
		}
		blocks := []evmtypes.Block{b0, b1, b2, b3}

		gas.SetRollingBlockHistory(bhe, blocks)

		// head2 and head3 have identical hash to saved blocks
		head2 := evmtypes.NewHead(big.NewInt(2), b2.Hash, b1.Hash, uint64(time.Now().Unix()), utils.NewBig(&cltest.FixtureChainID))
		head3 := evmtypes.NewHead(big.NewInt(3), b3.Hash, head2.Hash, uint64(time.Now().Unix()), utils.NewBig(&cltest.FixtureChainID))
		head3.Parent = &head2

		err := bhe.FetchBlocks(testutils.Context(t), &head3)
		require.NoError(t, err)

		require.Len(t, gas.GetRollingBlockHistory(bhe), 3)
		assert.Equal(t, 1, int(gas.GetRollingBlockHistory(bhe)[0].Number))
		assert.Equal(t, 2, int(gas.GetRollingBlockHistory(bhe)[1].Number))
		assert.Equal(t, 3, int(gas.GetRollingBlockHistory(bhe)[2].Number))
		assert.Equal(t, b1.Hash.Hex(), gas.GetRollingBlockHistory(bhe)[0].Hash.Hex())
		assert.Equal(t, head2.Hash.Hex(), gas.GetRollingBlockHistory(bhe)[1].Hash.Hex())
		assert.Equal(t, head3.Hash.Hex(), gas.GetRollingBlockHistory(bhe)[2].Hash.Hex())
	})

	t.Run("fetches max(BlockHistoryEstimatorCheckInclusionBlocks, BlockHistoryEstimatorBlockHistorySize)", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)
		var blockDelay uint16
		var historySize uint16 = 1
		var batchSize uint32 = 2
		var checkInclusionBlocks uint16 = 2
		cfg.BlockHistoryEstimatorBlockDelayF = blockDelay
		cfg.BlockHistoryEstimatorBlockHistorySizeF = historySize
		cfg.BlockHistoryEstimatorBatchSizeF = batchSize
		cfg.BlockHistoryEstimatorCheckInclusionBlocksF = checkInclusionBlocks

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		b42 := evmtypes.Block{
			Number:       42,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(3),
		}
		b43 := evmtypes.Block{
			Number:       43,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(),
		}

		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == gas.Int64ToHex(43) && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&evmtypes.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == gas.Int64ToHex(42) && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&evmtypes.Block{})
		})).Once().Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &b43
			elems[1].Result = &b42
		})

		err := bhe.FetchBlocks(testutils.Context(t), cltest.Head(43))
		require.NoError(t, err)

		require.Len(t, gas.GetRollingBlockHistory(bhe), 2)
		assert.Equal(t, 42, int(gas.GetRollingBlockHistory(bhe)[0].Number))
		assert.Equal(t, 43, int(gas.GetRollingBlockHistory(bhe)[1].Number))
		assert.Len(t, gas.GetRollingBlockHistory(bhe)[0].Transactions, 1)
		assert.Len(t, gas.GetRollingBlockHistory(bhe)[1].Transactions, 0)
	})
}

func TestBlockHistoryEstimator_FetchBlocksAndRecalculate_NoEIP1559(t *testing.T) {
	t.Parallel()

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	cfg := newConfigWithEIP1559DynamicFeesDisabled(t)

	cfg.BlockHistoryEstimatorBlockDelayF = uint16(0)
	cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(35)
	cfg.BlockHistoryEstimatorBlockHistorySizeF = uint16(3)
	cfg.EvmMaxGasPriceWeiF = assets.NewWeiI(1000)
	cfg.EvmMinGasPriceWeiF = assets.NewWeiI(0)
	cfg.BlockHistoryEstimatorBatchSizeF = uint32(0)

	bhe := newBlockHistoryEstimator(t, ethClient, cfg)

	b1 := evmtypes.Block{
		Number:       1,
		Hash:         utils.NewHash(),
		Transactions: cltest.LegacyTransactionsFromGasPrices(1),
	}
	b2 := evmtypes.Block{
		Number:       2,
		Hash:         utils.NewHash(),
		Transactions: cltest.LegacyTransactionsFromGasPrices(2),
	}
	b3 := evmtypes.Block{
		Number:       3,
		Hash:         utils.NewHash(),
		Transactions: cltest.LegacyTransactionsFromGasPrices(200, 300, 100, 100, 100, 100),
	}

	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 3 &&
			b[0].Args[0] == "0x3" &&
			b[1].Args[0] == "0x2" &&
			b[2].Args[0] == "0x1"
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		elems[0].Result = &b3
		elems[1].Result = &b2
		elems[2].Result = &b1
	})

	bhe.FetchBlocksAndRecalculate(testutils.Context(t), cltest.Head(3))

	price := gas.GetGasPrice(bhe)
	require.Equal(t, assets.NewWeiI(100), price)

	assert.Len(t, gas.GetRollingBlockHistory(bhe), 3)
}

func TestBlockHistoryEstimator_Recalculate_NoEIP1559(t *testing.T) {
	t.Parallel()

	maxGasPrice := assets.NewWeiI(100)
	minGasPrice := assets.NewWeiI(10)

	t.Run("does not crash or set gas price to zero if there are no transactions", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

		cfg := newConfigWithEIP1559DynamicFeesDisabled(t)

		cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(35)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		blocks := []evmtypes.Block{}
		gas.SetRollingBlockHistory(bhe, blocks)
		bhe.Recalculate(cltest.Head(1))

		blocks = []evmtypes.Block{evmtypes.Block{}}
		gas.SetRollingBlockHistory(bhe, blocks)
		bhe.Recalculate(cltest.Head(1))

		blocks = []evmtypes.Block{evmtypes.Block{Transactions: []evmtypes.Transaction{}}}
		gas.SetRollingBlockHistory(bhe, blocks)
		bhe.Recalculate(cltest.Head(1))
	})

	t.Run("sets gas price to EVM.GasEstimator.PriceMax if the calculation would otherwise exceed it", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesDisabled(t)

		cfg.EvmMaxGasPriceWeiF = maxGasPrice
		cfg.EvmMinGasPriceWeiF = minGasPrice
		cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(35)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		blocks := []evmtypes.Block{
			evmtypes.Block{
				Number:       0,
				Hash:         utils.NewHash(),
				Transactions: cltest.LegacyTransactionsFromGasPrices(9001),
			},
			evmtypes.Block{
				Number:       1,
				Hash:         utils.NewHash(),
				Transactions: cltest.LegacyTransactionsFromGasPrices(9002),
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(1))

		price := gas.GetGasPrice(bhe)
		require.Equal(t, maxGasPrice, price)
	})

	t.Run("sets gas price to EVM.Transactions.PriceMin if the calculation would otherwise fall below it", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesDisabled(t)

		cfg.EvmMaxGasPriceWeiF = maxGasPrice
		cfg.EvmMinGasPriceWeiF = minGasPrice
		cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(35)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		blocks := []evmtypes.Block{
			evmtypes.Block{
				Number:       0,
				Hash:         utils.NewHash(),
				Transactions: cltest.LegacyTransactionsFromGasPrices(5),
			},
			evmtypes.Block{
				Number:       1,
				Hash:         utils.NewHash(),
				Transactions: cltest.LegacyTransactionsFromGasPrices(7),
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(1))

		price := gas.GetGasPrice(bhe)
		require.Equal(t, minGasPrice, price)
	})

	t.Run("ignores any transaction with a zero gas limit", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesDisabled(t)

		cfg.EvmMaxGasPriceWeiF = maxGasPrice
		cfg.EvmMinGasPriceWeiF = minGasPrice
		cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(100)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		b1Hash := utils.NewHash()
		b2Hash := utils.NewHash()

		blocks := []evmtypes.Block{
			{
				Number:       0,
				Hash:         b1Hash,
				ParentHash:   common.Hash{},
				Transactions: cltest.LegacyTransactionsFromGasPrices(50),
			},
			{
				Number:       1,
				Hash:         b2Hash,
				ParentHash:   b1Hash,
				Transactions: []evmtypes.Transaction{evmtypes.Transaction{GasPrice: assets.NewWeiI(70), GasLimit: 42}},
			},
			{
				Number:       2,
				Hash:         utils.NewHash(),
				ParentHash:   b2Hash,
				Transactions: []evmtypes.Transaction{evmtypes.Transaction{GasPrice: assets.NewWeiI(90), GasLimit: 0}},
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(2))

		price := gas.GetGasPrice(bhe)
		require.Equal(t, assets.NewWeiI(70), price)
	})

	t.Run("takes into account zero priced transactions if chain is not xDai", func(t *testing.T) {
		// Because everyone loves free gas!
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesDisabled(t)

		cfg.EvmMaxGasPriceWeiF = maxGasPrice
		cfg.EvmMinGasPriceWeiF = assets.NewWeiI(0)
		cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(50)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		b1Hash := utils.NewHash()

		blocks := []evmtypes.Block{
			evmtypes.Block{
				Number:       0,
				Hash:         b1Hash,
				ParentHash:   common.Hash{},
				Transactions: cltest.LegacyTransactionsFromGasPrices(0, 0, 0, 0, 100),
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(0))

		price := gas.GetGasPrice(bhe)
		require.Equal(t, assets.NewWeiI(0), price)
	})

	t.Run("ignores zero priced transactions on xDai", func(t *testing.T) {
		chainID := big.NewInt(100)

		ethClient := evmtest.NewEthClientMock(t)
		cfg := newConfigWithEIP1559DynamicFeesDisabled(t)

		cfg.EvmMaxGasPriceWeiF = maxGasPrice
		cfg.EvmMinGasPriceWeiF = assets.NewWeiI(100)
		cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(50)

		ibhe := newBlockHistoryEstimatorWithChainID(t, ethClient, cfg, *chainID)
		bhe := gas.BlockHistoryEstimatorFromInterface(ibhe)

		b1Hash := utils.NewHash()

		blocks := []evmtypes.Block{
			evmtypes.Block{
				Number:       0,
				Hash:         b1Hash,
				ParentHash:   common.Hash{},
				Transactions: cltest.LegacyTransactionsFromGasPrices(0, 0, 0, 0, 100),
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(0))

		price := gas.GetGasPrice(bhe)
		require.Equal(t, assets.NewWeiI(100), price)
	})

	t.Run("handles unreasonably large gas prices (larger than a 64 bit int can hold)", func(t *testing.T) {
		// Seems unlikely we will ever experience gas prices > 9 Petawei on mainnet (praying to the eth Gods 🙏)
		// But other chains could easily use a different base of account
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesDisabled(t)

		reasonablyHugeGasPrice := assets.NewWeiI(1000).Mul(big.NewInt(math.MaxInt64))

		cfg.EvmMaxGasPriceWeiF = reasonablyHugeGasPrice
		cfg.EvmMinGasPriceWeiF = assets.NewWeiI(10)
		cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(50)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		unreasonablyHugeGasPrice := assets.NewWeiI(1000000).Mul(big.NewInt(math.MaxInt64))

		b1Hash := utils.NewHash()

		blocks := []evmtypes.Block{
			evmtypes.Block{
				Number:     0,
				Hash:       b1Hash,
				ParentHash: common.Hash{},
				Transactions: []evmtypes.Transaction{
					evmtypes.Transaction{GasPrice: assets.NewWeiI(50), GasLimit: 42},
					evmtypes.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					evmtypes.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					evmtypes.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					evmtypes.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					evmtypes.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					evmtypes.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					evmtypes.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					evmtypes.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
				},
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(0))

		price := gas.GetGasPrice(bhe)
		require.Equal(t, reasonablyHugeGasPrice, price)
	})

	t.Run("doesn't panic if gas price is nil (although I'm still unsure how this can happen)", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesDisabled(t)

		cfg.EvmMaxGasPriceWeiF = maxGasPrice
		cfg.EvmMinGasPriceWeiF = assets.NewWeiI(100)
		cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(50)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		b1Hash := utils.NewHash()

		blocks := []evmtypes.Block{
			evmtypes.Block{
				Number:     0,
				Hash:       b1Hash,
				ParentHash: common.Hash{},
				Transactions: []evmtypes.Transaction{
					{GasPrice: nil, GasLimit: 42, Hash: utils.NewHash()},
					{GasPrice: assets.NewWeiI(100), GasLimit: 42, Hash: utils.NewHash()},
				},
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(0))

		price := gas.GetGasPrice(bhe)
		require.Equal(t, assets.NewWeiI(100), price)
	})
}

func newBlockWithBaseFee() evmtypes.Block {
	return evmtypes.Block{BaseFeePerGas: assets.GWei(5)}
}

func TestBlockHistoryEstimator_Recalculate_EIP1559(t *testing.T) {
	t.Parallel()

	maxGasPrice := assets.NewWeiI(100)

	t.Run("does not crash or set gas price to zero if there are no transactions", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)

		cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(35)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		blocks := []evmtypes.Block{}
		gas.SetRollingBlockHistory(bhe, blocks)
		bhe.Recalculate(cltest.Head(1))

		blocks = []evmtypes.Block{evmtypes.Block{}} // No base fee (doesn't crash)
		gas.SetRollingBlockHistory(bhe, blocks)
		bhe.Recalculate(cltest.Head(1))

		blocks = []evmtypes.Block{newBlockWithBaseFee()}
		gas.SetRollingBlockHistory(bhe, blocks)
		bhe.Recalculate(cltest.Head(1))

		empty := newBlockWithBaseFee()
		empty.Transactions = []evmtypes.Transaction{}
		blocks = []evmtypes.Block{empty}
		gas.SetRollingBlockHistory(bhe, blocks)
		bhe.Recalculate(cltest.Head(1))

		withOnlyLegacyTransactions := newBlockWithBaseFee()
		withOnlyLegacyTransactions.Transactions = cltest.LegacyTransactionsFromGasPrices(9001)
		blocks = []evmtypes.Block{withOnlyLegacyTransactions}
		gas.SetRollingBlockHistory(bhe, blocks)
		bhe.Recalculate(cltest.Head(1))
	})

	t.Run("does not set tip higher than EVM.GasEstimator.PriceMax", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)

		cfg.EvmMaxGasPriceWeiF = maxGasPrice
		cfg.EvmMinGasPriceWeiF = assets.NewWeiI(0)
		cfg.EvmGasTipCapMinimumF = assets.NewWeiI(0)
		cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(35)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		blocks := []evmtypes.Block{
			evmtypes.Block{
				BaseFeePerGas: assets.NewWeiI(1),
				Number:        0,
				Hash:          utils.NewHash(),
				Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(9001),
			},
			evmtypes.Block{
				BaseFeePerGas: assets.NewWeiI(1),
				Number:        1,
				Hash:          utils.NewHash(),
				Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(9002),
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(1))

		tipCap := gas.GetTipCap(bhe)
		require.Equal(t, tipCap.Int64(), maxGasPrice.Int64())
	})

	t.Run("sets tip cap to EVM.Transactions.PriceMin if the calculation would otherwise fall below it", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)

		cfg.EvmMaxGasPriceWeiF = maxGasPrice
		cfg.EvmMinGasPriceWeiF = assets.NewWeiI(0)
		cfg.EvmGasTipCapMinimumF = assets.NewWeiI(10)
		cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(35)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		blocks := []evmtypes.Block{
			evmtypes.Block{
				BaseFeePerGas: assets.NewWeiI(1),
				Number:        0,
				Hash:          utils.NewHash(),
				Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(5),
			},
			evmtypes.Block{
				BaseFeePerGas: assets.NewWeiI(1),
				Number:        1,
				Hash:          utils.NewHash(),
				Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(7),
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(1))

		price := gas.GetTipCap(bhe)
		require.Equal(t, assets.NewWeiI(10), price)
	})

	t.Run("ignores any transaction with a zero gas limit", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)

		cfg.EvmMaxGasPriceWeiF = maxGasPrice
		cfg.EvmMinGasPriceWeiF = assets.NewWeiI(0)
		cfg.EvmGasTipCapMinimumF = assets.NewWeiI(10)
		cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(95)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		b1Hash := utils.NewHash()
		b2Hash := utils.NewHash()

		blocks := []evmtypes.Block{
			{
				Number:       0,
				Hash:         b1Hash,
				ParentHash:   common.Hash{},
				Transactions: cltest.LegacyTransactionsFromGasPrices(50),
			},
			{
				BaseFeePerGas: assets.NewWeiI(10),
				Number:        1,
				Hash:          b2Hash,
				ParentHash:    b1Hash,
				Transactions:  []evmtypes.Transaction{evmtypes.Transaction{Type: 0x2, MaxFeePerGas: assets.NewWeiI(1000), MaxPriorityFeePerGas: assets.NewWeiI(60), GasLimit: 42}},
			},
			{
				Number:       2,
				Hash:         utils.NewHash(),
				ParentHash:   b2Hash,
				Transactions: []evmtypes.Transaction{evmtypes.Transaction{Type: 0x2, MaxFeePerGas: assets.NewWeiI(1000), MaxPriorityFeePerGas: assets.NewWeiI(80), GasLimit: 0}},
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(2))

		price := gas.GetTipCap(bhe)
		require.Equal(t, assets.NewWeiI(60), price)
	})

	t.Run("respects minimum gas tip cap", func(t *testing.T) {
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)

		cfg.EvmMaxGasPriceWeiF = maxGasPrice
		cfg.EvmMinGasPriceWeiF = assets.NewWeiI(0)
		cfg.EvmGasTipCapMinimumF = assets.NewWeiI(1)
		cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(35)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		b1Hash := utils.NewHash()

		blocks := []evmtypes.Block{
			evmtypes.Block{
				BaseFeePerGas: assets.NewWeiI(10),
				Number:        0,
				Hash:          b1Hash,
				ParentHash:    common.Hash{},
				Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(0, 0, 0, 0, 100),
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(0))

		price := gas.GetTipCap(bhe)
		assert.Equal(t, assets.NewWeiI(1), price)
	})

	t.Run("allows to set zero tip cap if minimum allows it", func(t *testing.T) {
		// Because everyone loves *cheap* gas!
		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)

		cfg.EvmMaxGasPriceWeiF = maxGasPrice
		cfg.EvmMinGasPriceWeiF = assets.NewWeiI(0)
		cfg.EvmGasTipCapMinimumF = assets.NewWeiI(0)
		cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(35)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		b1Hash := utils.NewHash()

		blocks := []evmtypes.Block{
			evmtypes.Block{
				BaseFeePerGas: assets.NewWeiI(10),
				Number:        0,
				Hash:          b1Hash,
				ParentHash:    common.Hash{},
				Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(0, 0, 0, 0, 100),
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(0))

		price := gas.GetTipCap(bhe)
		require.Equal(t, assets.NewWeiI(0), price)
	})
}

func TestBlockHistoryEstimator_EffectiveTipCap(t *testing.T) {
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	cfg := newConfigWithEIP1559DynamicFeesEnabled(t)

	bhe := newBlockHistoryEstimator(t, ethClient, cfg)

	block := evmtypes.Block{
		Number:     0,
		Hash:       utils.NewHash(),
		ParentHash: common.Hash{},
	}

	eipblock := block
	eipblock.BaseFeePerGas = assets.NewWeiI(100)

	t.Run("returns nil if block is missing base fee", func(t *testing.T) {
		tx := evmtypes.Transaction{Type: 0x0, GasPrice: assets.NewWeiI(42), GasLimit: 42, Hash: utils.NewHash()}
		res := bhe.EffectiveTipCap(block, tx)
		assert.Nil(t, res)
	})
	t.Run("legacy transaction type infers tip cap from tx.gas_price - block.base_fee_per_gas", func(t *testing.T) {
		tx := evmtypes.Transaction{Type: 0x0, GasPrice: assets.NewWeiI(142), GasLimit: 42, Hash: utils.NewHash()}
		res := bhe.EffectiveTipCap(eipblock, tx)
		assert.Equal(t, "42 wei", res.String())
	})
	t.Run("tx type 2 should calculate gas price", func(t *testing.T) {
		// 0x2 transaction (should use MaxPriorityFeePerGas)
		tx := evmtypes.Transaction{Type: 0x2, MaxPriorityFeePerGas: assets.NewWeiI(200), MaxFeePerGas: assets.NewWeiI(250), GasLimit: 42, Hash: utils.NewHash()}
		res := bhe.EffectiveTipCap(eipblock, tx)
		assert.Equal(t, "200 wei", res.String())
		// 0x2 transaction (should use MaxPriorityFeePerGas, ignoring gas price)
		tx = evmtypes.Transaction{Type: 0x2, GasPrice: assets.NewWeiI(400), MaxPriorityFeePerGas: assets.NewWeiI(200), MaxFeePerGas: assets.NewWeiI(350), GasLimit: 42, Hash: utils.NewHash()}
		res = bhe.EffectiveTipCap(eipblock, tx)
		assert.Equal(t, "200 wei", res.String())
	})
	t.Run("missing field returns nil", func(t *testing.T) {
		tx := evmtypes.Transaction{Type: 0x2, GasPrice: assets.NewWeiI(132), MaxFeePerGas: assets.NewWeiI(200), GasLimit: 42, Hash: utils.NewHash()}
		res := bhe.EffectiveTipCap(eipblock, tx)
		assert.Nil(t, res)
	})
	t.Run("unknown type returns nil", func(t *testing.T) {
		tx := evmtypes.Transaction{Type: 0x3, GasPrice: assets.NewWeiI(55555), MaxPriorityFeePerGas: assets.NewWeiI(200), MaxFeePerGas: assets.NewWeiI(250), GasLimit: 42, Hash: utils.NewHash()}
		res := bhe.EffectiveTipCap(eipblock, tx)
		assert.Nil(t, res)
	})
}

func TestBlockHistoryEstimator_EffectiveGasPrice(t *testing.T) {
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	cfg := newConfigWithEIP1559DynamicFeesDisabled(t)

	bhe := newBlockHistoryEstimator(t, ethClient, cfg)

	block := evmtypes.Block{
		Number:     0,
		Hash:       utils.NewHash(),
		ParentHash: common.Hash{},
	}

	eipblock := block
	eipblock.BaseFeePerGas = assets.NewWeiI(100)

	t.Run("legacy transaction type should use GasPrice", func(t *testing.T) {
		tx := evmtypes.Transaction{Type: 0x0, GasPrice: assets.NewWeiI(42), GasLimit: 42, Hash: utils.NewHash()}
		res := bhe.EffectiveGasPrice(eipblock, tx)
		assert.Equal(t, "42 wei", res.String())
		tx = evmtypes.Transaction{Type: 0x0, GasLimit: 42, Hash: utils.NewHash()}
		res = bhe.EffectiveGasPrice(eipblock, tx)
		assert.Nil(t, res)
		tx = evmtypes.Transaction{Type: 0x1, GasPrice: assets.NewWeiI(42), GasLimit: 42, Hash: utils.NewHash()}
		res = bhe.EffectiveGasPrice(eipblock, tx)
		assert.Equal(t, "42 wei", res.String())
	})
	t.Run("tx type 2 should calculate gas price", func(t *testing.T) {
		// 0x2 transaction (should calculate to 250)
		tx := evmtypes.Transaction{Type: 0x2, MaxPriorityFeePerGas: assets.NewWeiI(200), MaxFeePerGas: assets.NewWeiI(250), GasLimit: 42, Hash: utils.NewHash()}
		res := bhe.EffectiveGasPrice(eipblock, tx)
		assert.Equal(t, "250 wei", res.String())
		// 0x2 transaction (should calculate to 300)
		tx = evmtypes.Transaction{Type: 0x2, MaxPriorityFeePerGas: assets.NewWeiI(200), MaxFeePerGas: assets.NewWeiI(350), GasLimit: 42, Hash: utils.NewHash()}
		res = bhe.EffectiveGasPrice(eipblock, tx)
		assert.Equal(t, "300 wei", res.String())
		// 0x2 transaction (should calculate to 300, ignoring gas price)
		tx = evmtypes.Transaction{Type: 0x2, MaxPriorityFeePerGas: assets.NewWeiI(200), MaxFeePerGas: assets.NewWeiI(350), GasLimit: 42, Hash: utils.NewHash()}
		res = bhe.EffectiveGasPrice(eipblock, tx)
		assert.Equal(t, "300 wei", res.String())
		// 0x2 transaction (should fall back to gas price since MaxFeePerGas is missing)
		tx = evmtypes.Transaction{Type: 0x2, GasPrice: assets.NewWeiI(32), MaxPriorityFeePerGas: assets.NewWeiI(200), GasLimit: 42, Hash: utils.NewHash()}
		res = bhe.EffectiveGasPrice(eipblock, tx)
		assert.Equal(t, "32 wei", res.String())
	})
	t.Run("tx type 2 has block missing base fee (should never happen but must handle gracefully)", func(t *testing.T) {
		// 0x2 transaction (should calculate to 250)
		tx := evmtypes.Transaction{Type: 0x2, GasPrice: assets.NewWeiI(55555), MaxPriorityFeePerGas: assets.NewWeiI(200), MaxFeePerGas: assets.NewWeiI(250), GasLimit: 42, Hash: utils.NewHash()}
		res := bhe.EffectiveGasPrice(block, tx)
		assert.Equal(t, "55.555 kwei", res.String())
	})
	t.Run("unknown type returns nil", func(t *testing.T) {
		tx := evmtypes.Transaction{Type: 0x3, GasPrice: assets.NewWeiI(55555), MaxPriorityFeePerGas: assets.NewWeiI(200), MaxFeePerGas: assets.NewWeiI(250), GasLimit: 42, Hash: utils.NewHash()}
		res := bhe.EffectiveGasPrice(block, tx)
		assert.Nil(t, res)
	})
}

func TestBlockHistoryEstimator_Block_Unmarshal(t *testing.T) {
	blockJSON := `
{
    "author": "0x1438087186fdbfd4c256fa2df446921e30e54df8",
    "difficulty": "0xfffffffffffffffffffffffffffffffd",
    "extraData": "0xdb830302058c4f70656e457468657265756d86312e35312e30826c69",
    "gasLimit": "0xbebc20",
    "gasUsed": "0xbb58ce",
    "hash": "0x317cfd032b5d6657995f17fe768f7cc4ea0ada27ad421c4caa685a9071ea955c",
    "logsBloom": "0x0004000021000004000020200088810004110800400030002140000020801020120020000000000108002087c030000a80402800001600080400000c00010002100001881002008000004809126000002802a0a801004001000012100000000010000000120000068000000010200800400000004400010400010098540440400044200020008480000000800040000000000c818000510002200c000020000400800221d20100000081800101840000080100041000002080080000408243424280020200680000000201224500000c120008000800220000800009080028088020400000000040002000400000046000000000400000000000000802008000",
    "miner": "0x1438087186fdbfd4c256fa2df446921e30e54df8",
    "number": "0xf47e79",
    "parentHash": "0xb47ab3b1dc5c2c090dcecdc744a65a279ea6bb8dec11fb3c247df4cc2f584848",
    "receiptsRoot": "0x6c0a0e448f63da4b6552333aaead47a9702cd5d08c9c42edbdc30622706c840b",
    "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "signature": "0x30c7bfa28eceacb9f6b7c4acbb5b82e21792825ab20db8ecd3570b7e106f362b715b51e98f85aa9bb02e411fa1916c3cbb6a0ca34cc66d32e1142ec5282d829500",
    "size": "0x10fd",
    "stateRoot": "0x32cfd26ec2360c44797fc631c2e2d0395befb8369601bd16d482e3e7be4ebf2c",
    "step": 324172559,
    "totalDifficulty": "0xf47e78ffffffffffffffffffffffffebbb0678",
    "timestamp": "0x609c674b",
    "transactions": [
			{
        "hash": "0x3f8e13d8c15d929bd3f7d99be94484eb82f328bbb76052c9464614c12f10b990",
        "nonce": "0x2bb04",
        "blockHash": "0x317cfd032b5d6657995f17fe768f7cc4ea0ada27ad421c4caa685a9071ea955c",
        "blockNumber": "0xf47e79",
        "transactionIndex": "0x0",
        "from": "0x1438087186fdbfd4c256fa2df446921e30e54df8",
        "to": "0x5870b0527dedb1cfbd9534343feda1a41ce47766",
        "value": "0x0",
        "gasPrice": "0x1",
        "gas": "0x1",
        "data": "0x0b61ba8554b40c84fe2c9b5aad2fb692bdc00a9ba7f87d0abd35c68715bb347440c841d9000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000910411107ae9ec4e54f9b9e76d2a269a75dfab916c1edb866159e152e370f1ca8f72e95bf922fa069af9d532bef4fee8c89a401a501c622d763e4944ecacad16b4ace8dd0d532124b7c376cb5b04e63c4bf43b704eeb7ca822ec4258d8b0c2b2f5ef3680b858d15bcdf2f3632ad9e92963f37234c51f809981f3d4e34519d1f853408bbbe015e9572f9fcd55e9c0c38333ff000000000000000000000000000000",
        "input": "0x0b61ba8554b40c84fe2c9b5aad2fb692bdc00a9ba7f87d0abd35c68715bb347440c841d9000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000910411107ae9ec4e54f9b9e76d2a269a75dfab916c1edb866159e152e370f1ca8f72e95bf922fa069af9d532bef4fee8c89a401a501c622d763e4944ecacad16b4ace8dd0d532124b7c376cb5b04e63c4bf43b704eeb7ca822ec4258d8b0c2b2f5ef3680b858d15bcdf2f3632ad9e92963f37234c51f809981f3d4e34519d1f853408bbbe015e9572f9fcd55e9c0c38333ff000000000000000000000000000000",
        "v": "0xeb",
        "s": "0x7bbc91758d2485a0d97e92bc4f0c226bf961c8aeb7db59d152206995937cd907",
        "r": "0xe34e3a2a8f3159238dc843250d4ae0507d12ef49dec7bcf3057e6bd7b8560ae"
      },
      {
        "hash": "0x3f8e13d8c15d929bd3f7d99be94484eb82f328bbb76052c9464614c12f10b990",
        "nonce": "0x2bb04",
        "blockHash": "0x317cfd032b5d6657995f17fe768f7cc4ea0ada27ad421c4caa685a9071ea955c",
        "blockNumber": "0xf47e79",
        "transactionIndex": "0x0",
        "from": "0x1438087186fdbfd4c256fa2df446921e30e54df8",
        "to": "0x5870b0527dedb1cfbd9534343feda1a41ce47766",
        "value": "0x0",
        "gasPrice": "0x0",
        "gas": "0x0",
        "data": "0x0b61ba8554b40c84fe2c9b5aad2fb692bdc00a9ba7f87d0abd35c68715bb347440c841d9000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000910411107ae9ec4e54f9b9e76d2a269a75dfab916c1edb866159e152e370f1ca8f72e95bf922fa069af9d532bef4fee8c89a401a501c622d763e4944ecacad16b4ace8dd0d532124b7c376cb5b04e63c4bf43b704eeb7ca822ec4258d8b0c2b2f5ef3680b858d15bcdf2f3632ad9e92963f37234c51f809981f3d4e34519d1f853408bbbe015e9572f9fcd55e9c0c38333ff000000000000000000000000000000",
        "input": "0x0b61ba8554b40c84fe2c9b5aad2fb692bdc00a9ba7f87d0abd35c68715bb347440c841d9000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000910411107ae9ec4e54f9b9e76d2a269a75dfab916c1edb866159e152e370f1ca8f72e95bf922fa069af9d532bef4fee8c89a401a501c622d763e4944ecacad16b4ace8dd0d532124b7c376cb5b04e63c4bf43b704eeb7ca822ec4258d8b0c2b2f5ef3680b858d15bcdf2f3632ad9e92963f37234c51f809981f3d4e34519d1f853408bbbe015e9572f9fcd55e9c0c38333ff000000000000000000000000000000",
        "type": "0x00",
        "v": "0xeb",
        "s": "0x7bbc91758d2485a0d97e92bc4f0c226bf961c8aeb7db59d152206995937cd907",
        "r": "0xe34e3a2a8f3159238dc843250d4ae0507d12ef49dec7bcf3057e6bd7b8560ae"
      },
      {
        "hash": "0x238423bddc38e241f35ea3ed52cb096352c71d423b9ea3441937754f4edcb312",
        "nonce": "0xb847",
        "blockHash": "0x317cfd032b5d6657995f17fe768f7cc4ea0ada27ad421c4caa685a9071ea955c",
        "blockNumber": "0xf47e79",
        "transactionIndex": "0x1",
        "from": "0x25461d55ca1ddf4317160fd917192fe1d981b908",
        "to": "0x5d9593586b4b5edbd23e7eba8d88fd8f09d83ebd",
        "value": "0x0",
        "gasPrice": "0x42725ae1000",
        "gas": "0x1e8480",
        "data": "0x893d242d000000000000000000000000eac6cee594edd353351babc145c624849bb70b1100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001e57396fe60670c00000000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000000000000000000000000000000000000000000000",
        "input": "0x893d242d000000000000000000000000eac6cee594edd353351babc145c624849bb70b1100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001e57396fe60670c00000000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000000000000000000000000000000000000000000000",
        "type": "0x00",
        "v": "0xeb",
        "s": "0x7f795b5cb15410b41c1518edc1aed2f1e984b8c93e357bdee79b23bba8dc841d",
        "r": "0x958db39caa6dd066d3b010a4d9e6427399601738e0071470d822594e4565aa99"
      }
	]
}
`

	var block evmtypes.Block
	err := json.Unmarshal([]byte(blockJSON), &block)
	assert.NoError(t, err)

	assert.Equal(t, int64(16023161), block.Number)
	assert.Equal(t, common.HexToHash("0x317cfd032b5d6657995f17fe768f7cc4ea0ada27ad421c4caa685a9071ea955c"), block.Hash)
	assert.Equal(t, common.HexToHash("0xb47ab3b1dc5c2c090dcecdc744a65a279ea6bb8dec11fb3c247df4cc2f584848"), block.ParentHash)

	require.Len(t, block.Transactions, 3)

	assert.Equal(t, int64(1), block.Transactions[0].GasPrice.Int64())
	assert.Equal(t, uint32(1), block.Transactions[0].GasLimit)

	assert.Equal(t, int64(0), block.Transactions[1].GasPrice.Int64())
	assert.Equal(t, uint32(0), block.Transactions[1].GasLimit)

	assert.Equal(t, assets.NewWeiI(4566182400000), block.Transactions[2].GasPrice)
	assert.Equal(t, uint32(2000000), block.Transactions[2].GasLimit)
}

func TestBlockHistoryEstimator_EIP1559Block_Unmarshal(t *testing.T) {
	blockJSON := `
{
    "baseFeePerGas": "0xa1894585c",
    "difficulty": "0x1cc4a2d7045f39",
    "extraData": "0x73656f32",
    "gasLimit": "0x1c9c380",
    "gasUsed": "0x1c9c203",
    "hash": "0x11ac873a6cd8b8b7b57ec1efe3984b706362aa5e8f5749a5ec9b1f64bb4615f0",
    "logsBloom": "0x2b181cd7982005346543c60498149414cc92419055218c5111988a6c81c7560105c91c82ec3348283288c2187b0111407e28c08c4b45b4ea2e980893c050002588606218aa083c0c0824e46923b850d07048da924052828c26082c910663fac682070310ba3189bed51194261220990c2920cc434d042c06a1941158dfc91eeb572107e1c5595a0032051109c500ba42a093398850ad020b1118d41716d371286ba348e041685144210401078b8901281001e840290d0e9391c00138cf00120d92499ca250d3026003e13c1e10bac2a3a57499007a2213002714a2a2f24f24480d0539c30142f2ed09105d5b10038330ac1622cc188a00f0c3108801455882cc",
    "miner": "0x3ecef08d0e2dad803847e052249bb4f8bff2d5bb",
    "mixHash": "0x57f4a273c69c4028916abfaa57252035fb7e71ce8444034764b8988d9a89c7b6",
    "nonce": "0x015e0d851f990730",
    "number": "0xc65d68",
    "parentHash": "0x1ae6168805dfd2e48311181774019c17fb09b24ab75dcad6566d18d38d5c4071",
    "receiptsRoot": "0x3ced645d38426647aad078b8e4bc62ff03571a74b099c983133eb34808240309",
    "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "size": "0x2655",
    "stateRoot": "0x073e7b70e9b1357329cbf0b19a10a981057a29accbafcc34d52b592dc0be9848",
    "timestamp": "0x6112f709",
    "totalDifficulty": "0x6171fd1e7626bc65d9b",
    "transactions": [
      {
        "blockHash": "0x11ac873a6cd8b8b7b57ec1efe3984b706362aa5e8f5749a5ec9b1f64bb4615f0",
        "blockNumber": "0xc65d68",
        "from": "0x305bf59bbd7a89ca9ce4d460b0efb54266d9e6c3",
        "gas": "0xdbba0",
        "gasPrice": "0x9f05f8ee00",
        "hash": "0x8e58af889f4e831ef9a67df84058bcfb7090cbcb5c6f1046c211dafee6050944",
        "input": "0xc18a84bc0000000000000000000000007ae132b71ddc6f4866fbf103be655830d9ca666c00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000124e94584ee00000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000bb0e17ef65f82ab018d8edd776e8dd940327b28b00000000000000000000000000000000000000000000002403ecad7d36e5bda0000000000000000000000000000000000000000000000000af7c8acfe5037ea80000000000000000000000000000000000000000000000000000000000c65d680000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002bbb0e17ef65f82ab018d8edd776e8dd940327b28b000bb8c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
        "nonce": "0x6654",
        "to": "0x4d246be90c2f36730bb853ad41d0a189061192d3",
        "transactionIndex": "0x0",
        "value": "0x0",
        "type": "0x0",
        "v": "0x25",
        "r": "0x9f8af9e6424f264daaba992c09c2b38d05444cbb5e6bd5e26c965393e287c9fa",
        "s": "0x76802388299eb0baa80a678831ef0722c5b1e1212f5eca26a5e911cb81388b2b"
      },
      {
        "blockHash": "0x11ac873a6cd8b8b7b57ec1efe3984b706362aa5e8f5749a5ec9b1f64bb4615f0",
        "blockNumber": "0xc65d68",
        "from": "0xef3f063136fe5002065bf7c4a2d85ff34cfb0ac0",
        "gas": "0xdfeae",
        "gasPrice": "0x2ba7def3000",
        "hash": "0x0190f436ce165abb741b8513f64d194682677e1db72422f0f533fe6c0248e59a",
        "input": "0x926427440000000000000000000000000000000000000000000000000000000000000005",
        "nonce": "0x267",
        "to": "0xad9fd7cb4fc7a0fbce08d64068f60cbde22ed34c",
        "transactionIndex": "0x1",
        "value": "0x62967a5c8460000",
        "type": "0x0",
        "v": "0x26",
        "r": "0xd06f53ad57d61543526b529c2532903ac0d45b1d727567d04dc9b2f4e6340521",
        "s": "0x6332bcec6a66abf4bed4df24e25e1e4dfc61c5d5bc32a441033c285c14c402d"
      },
	  {
        "blockHash": "0x11ac873a6cd8b8b7b57ec1efe3984b706362aa5e8f5749a5ec9b1f64bb4615f0",
        "blockNumber": "0xc65d68",
        "from": "0xff54553ff5edf0e93d58555303291805770e5793",
        "gas": "0x5208",
        "gasPrice": "0x746a528800",
        "maxFeePerGas": "0x746a528800",
        "maxPriorityFeePerGas": "0x746a528800",
        "hash": "0x136aa666e6b8109b2b4aca8008ecad8df2047f4e2aced4808248fa8927a13395",
        "input": "0x",
        "nonce": "0x1",
        "to": "0xb5d85cbf7cb3ee0d56b3bb207d5fc4b82f43f511",
        "transactionIndex": "0x3b",
        "value": "0x1302a5a6ad330400",
        "type": "0x2",
        "accessList": [],
        "chainId": "0x1",
        "v": "0x1",
        "r": "0x2806aa357d15790319e1def013902135dc8fa191182e2f87edae352e50ef281",
        "s": "0x61d160d7de9af375c7fc40aed956e711af3af20146afe27d5122adf28cd25c9"
      },
      {
        "blockHash": "0x11ac873a6cd8b8b7b57ec1efe3984b706362aa5e8f5749a5ec9b1f64bb4615f0",
        "blockNumber": "0xc65d68",
        "from": "0xb090838386b9207994a42f740217066af2de53ad",
        "gas": "0x5208",
        "maxFeePerGas": "0x746a528800",
        "maxPriorityFeePerGas": "0x746a528800",
        "hash": "0x13d4ecea98e37359e63e39e350ed0b1456e1acbf985eb8d4a0ef0e89a705c10d",
        "input": "0x",
        "nonce": "0x1",
        "to": "0xb5d85cbf7cb3ee0d56b3bb207d5fc4b82f43f511",
        "transactionIndex": "0x3c",
        "value": "0xe95497bc358fe60",
        "type": "0x2",
        "accessList": [],
        "chainId": "0x1",
        "v": "0x1",
        "r": "0xa0d09f41bb4279d73e4255a1c1ce6cb10cb1fba04b4eca4af582ab2928201b27",
        "s": "0x682f2a7a734b7c5887c5e228d35af4d3d3ad240c2c14f97aa9145a6c9edcd0a1"
      }
	]
}
`

	var block evmtypes.Block
	err := json.Unmarshal([]byte(blockJSON), &block)
	assert.NoError(t, err)

	assert.Equal(t, int64(13000040), block.Number)
	assert.Equal(t, "43.362048092 gwei", block.BaseFeePerGas.String())
	assert.Equal(t, common.HexToHash("0x11ac873a6cd8b8b7b57ec1efe3984b706362aa5e8f5749a5ec9b1f64bb4615f0"), block.Hash)
	assert.Equal(t, common.HexToHash("0x1ae6168805dfd2e48311181774019c17fb09b24ab75dcad6566d18d38d5c4071"), block.ParentHash)

	require.Len(t, block.Transactions, 4)

	assert.Equal(t, int64(683000000000), block.Transactions[0].GasPrice.Int64())
	assert.Equal(t, 900000, int(block.Transactions[0].GasLimit))
	assert.Nil(t, block.Transactions[0].MaxFeePerGas)
	assert.Nil(t, block.Transactions[0].MaxPriorityFeePerGas)
	assert.Equal(t, evmtypes.TxType(0x0), block.Transactions[0].Type)
	assert.Equal(t, "0x8e58af889f4e831ef9a67df84058bcfb7090cbcb5c6f1046c211dafee6050944", block.Transactions[0].Hash.String())

	assert.Equal(t, assets.NewWeiI(3000000000000), block.Transactions[1].GasPrice)
	assert.Equal(t, "0x0190f436ce165abb741b8513f64d194682677e1db72422f0f533fe6c0248e59a", block.Transactions[1].Hash.String())

	assert.Equal(t, int64(500000000000), block.Transactions[2].GasPrice.Int64())
	assert.Equal(t, 21000, int(block.Transactions[2].GasLimit))
	assert.Equal(t, int64(500000000000), block.Transactions[2].MaxFeePerGas.Int64())
	assert.Equal(t, int64(500000000000), block.Transactions[2].MaxPriorityFeePerGas.Int64())
	assert.Equal(t, evmtypes.TxType(0x2), block.Transactions[2].Type)
	assert.Equal(t, "0x136aa666e6b8109b2b4aca8008ecad8df2047f4e2aced4808248fa8927a13395", block.Transactions[2].Hash.String())

	assert.Nil(t, block.Transactions[3].GasPrice)
	assert.Equal(t, 21000, int(block.Transactions[3].GasLimit))
	assert.Equal(t, "0x13d4ecea98e37359e63e39e350ed0b1456e1acbf985eb8d4a0ef0e89a705c10d", block.Transactions[3].Hash.String())
}

func TestBlockHistoryEstimator_GetLegacyGas(t *testing.T) {
	t.Parallel()

	cfg := newConfigWithEIP1559DynamicFeesDisabled(t)

	maxGasPrice := assets.NewWeiI(1000000)
	cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(35)
	cfg.EvmEIP1559DynamicFeesF = false
	cfg.EvmGasLimitMultiplierF = float32(1)
	cfg.EvmMaxGasPriceWeiF = maxGasPrice
	cfg.EvmMinGasPriceWeiF = assets.NewWeiI(0)
	cfg.BlockHistoryEstimatorCheckInclusionBlocksF = uint16(0)
	cfg.BlockHistoryEstimatorBlockHistorySizeF = uint16(8)

	bhe := newBlockHistoryEstimator(t, nil, cfg)

	blocks := []evmtypes.Block{
		{
			Number:       0,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(1000),
		},
		{
			Number:       1,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(1200),
		},
	}

	gas.SetRollingBlockHistory(bhe, blocks)
	bhe.Recalculate(cltest.Head(1))
	gas.SimulateStart(t, bhe)

	t.Run("if gas price is lower than global max and user specified max gas price", func(t *testing.T) {
		fee, limit, err := bhe.GetLegacyGas(testutils.Context(t), make([]byte, 0), 10000, maxGasPrice)
		require.NoError(t, err)

		assert.Equal(t, assets.NewWeiI(1000), fee)
		assert.Equal(t, 10000, int(limit))
	})

	t.Run("if gas price is higher than user-specified max", func(t *testing.T) {
		fee, limit, err := bhe.GetLegacyGas(testutils.Context(t), make([]byte, 0), 10000, assets.NewWeiI(800))
		require.NoError(t, err)

		assert.Equal(t, assets.NewWeiI(800), fee)
		assert.Equal(t, 10000, int(limit))
	})

	cfg = newConfigWithEIP1559DynamicFeesDisabled(t)
	cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(35)
	cfg.EvmEIP1559DynamicFeesF = false
	cfg.EvmGasLimitMultiplierF = float32(1)
	cfg.EvmMaxGasPriceWeiF = assets.NewWeiI(700)
	cfg.EvmMinGasPriceWeiF = assets.NewWeiI(0)
	bhe = newBlockHistoryEstimator(t, nil, cfg)
	gas.SetRollingBlockHistory(bhe, blocks)
	bhe.Recalculate(cltest.Head(1))
	gas.SimulateStart(t, bhe)

	t.Run("if gas price is higher than global max", func(t *testing.T) {
		fee, limit, err := bhe.GetLegacyGas(testutils.Context(t), make([]byte, 0), 10000, maxGasPrice)
		require.NoError(t, err)

		assert.Equal(t, assets.NewWeiI(700), fee)
		assert.Equal(t, 10000, int(limit))
	})
}

func TestBlockHistoryEstimator_UseDefaultPriceAsFallback(t *testing.T) {
	t.Parallel()

	var batchSize uint32
	var blockDelay uint16
	var historySize uint16 = 3
	var specialTxTypeCode evmtypes.TxType = 0x7e

	t.Run("fallbacks to EvmGasPriceDefault if there aren't any valid transactions to estimate from.", func(t *testing.T) {
		cfg := newConfigWithEIP1559DynamicFeesDisabled(t)

		cfg.BlockHistoryEstimatorBatchSizeF = batchSize
		cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(35)
		cfg.EvmEIP1559DynamicFeesF = false
		cfg.EvmGasLimitMultiplierF = float32(1)
		cfg.EvmMaxGasPriceWeiF = assets.NewWeiI(1000000)
		cfg.EvmGasPriceDefaultF = assets.NewWeiI(100)
		cfg.BlockHistoryEstimatorBlockDelayF = blockDelay
		cfg.BlockHistoryEstimatorBlockHistorySizeF = historySize

		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		h := &evmtypes.Head{Hash: utils.NewHash(), Number: 42, BaseFeePerGas: nil}
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(h, nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 3 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == gas.Int64ToHex(42) && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&evmtypes.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == gas.Int64ToHex(41) && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&evmtypes.Block{}) &&
				b[2].Method == "eth_getBlockByNumber" && b[2].Args[0] == gas.Int64ToHex(40) && b[1].Args[1] == true && reflect.TypeOf(b[2].Result) == reflect.TypeOf(&evmtypes.Block{})
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &evmtypes.Block{
				Number: 42,
				Hash:   utils.NewHash(),
			}
			elems[1].Result = &evmtypes.Block{
				Number: 41,
				Hash:   utils.NewHash(),
			}
			elems[2].Result = &evmtypes.Block{
				Number:       40,
				Hash:         utils.NewHash(),
				Transactions: cltest.LegacyTransactionsFromGasPricesTxType(specialTxTypeCode, 1),
			}
		}).Once()

		err := bhe.Start(testutils.Context(t))
		require.NoError(t, err)

		fee, limit, err := bhe.GetLegacyGas(testutils.Context(t), make([]byte, 0), 10000, assets.NewWeiI(800))
		require.NoError(t, err)
		require.Equal(t, cfg.EvmGasPriceDefault(), fee)
		assert.Equal(t, 10000, int(limit))
	})

	t.Run("fallbacks to EvmGasTipCapDefault if there aren't any valid transactions to estimate from.", func(t *testing.T) {
		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)
		cfg.BlockHistoryEstimatorBatchSizeF = batchSize
		cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(35)
		cfg.EvmEIP1559DynamicFeesF = true
		cfg.EvmGasLimitMultiplierF = float32(1)
		cfg.EvmMaxGasPriceWeiF = assets.NewWeiI(1000000)
		cfg.EvmGasPriceDefaultF = assets.NewWeiI(100)
		cfg.BlockHistoryEstimatorBlockDelayF = blockDelay
		cfg.BlockHistoryEstimatorBlockHistorySizeF = historySize
		cfg.BlockHistoryEstimatorEIP1559FeeCapBufferBlocksF = uint16(4)
		cfg.EvmGasTipCapDefaultF = assets.NewWeiI(50)
		cfg.EvmGasBumpThresholdF = uint64(1)

		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		h := &evmtypes.Head{Hash: utils.NewHash(), Number: 42, BaseFeePerGas: assets.NewWeiI(40)}
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(h, nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 3 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == gas.Int64ToHex(42) && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&evmtypes.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == gas.Int64ToHex(41) && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&evmtypes.Block{}) &&
				b[2].Method == "eth_getBlockByNumber" && b[2].Args[0] == gas.Int64ToHex(40) && b[1].Args[1] == true && reflect.TypeOf(b[2].Result) == reflect.TypeOf(&evmtypes.Block{})
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &evmtypes.Block{
				Number: 42,
				Hash:   utils.NewHash(),
			}
			elems[1].Result = &evmtypes.Block{
				Number: 41,
				Hash:   utils.NewHash(),
			}
			elems[2].Result = &evmtypes.Block{
				Number:       40,
				Hash:         utils.NewHash(),
				Transactions: cltest.DynamicFeeTransactionsFromTipCapsTxType(specialTxTypeCode, 1),
			}
		}).Once()

		err := bhe.Start(testutils.Context(t))
		require.NoError(t, err)
		fee, limit, err := bhe.GetDynamicFee(testutils.Context(t), 100000, assets.NewWeiI(200))
		require.NoError(t, err)

		assert.Equal(t, gas.DynamicFee{FeeCap: assets.NewWeiI(114), TipCap: cfg.EvmGasTipCapDefault()}, fee)
		assert.Equal(t, 100000, int(limit))
	})
}

func TestBlockHistoryEstimator_GetDynamicFee(t *testing.T) {
	t.Parallel()

	cfg := newConfigWithEIP1559DynamicFeesEnabled(t)
	maxGasPrice := assets.NewWeiI(1000000)
	cfg.BlockHistoryEstimatorEIP1559FeeCapBufferBlocksF = uint16(4)
	cfg.BlockHistoryEstimatorTransactionPercentileF = uint16(35)
	cfg.EvmEIP1559DynamicFeesF = true
	cfg.EvmGasLimitMultiplierF = float32(1)
	cfg.EvmMaxGasPriceWeiF = maxGasPrice
	cfg.EvmGasTipCapMinimumF = assets.NewWeiI(0)
	cfg.EvmMinGasPriceWeiF = assets.NewWeiI(0)

	bhe := newBlockHistoryEstimator(t, nil, cfg)

	blocks := []evmtypes.Block{
		evmtypes.Block{
			BaseFeePerGas: assets.NewWeiI(88889),
			Number:        0,
			Hash:          utils.NewHash(),
			Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(5000, 6000, 6000),
		},
		evmtypes.Block{
			BaseFeePerGas: assets.NewWeiI(100000),
			Number:        1,
			Hash:          utils.NewHash(),
			Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(10000),
		},
	}
	gas.SetRollingBlockHistory(bhe, blocks)

	bhe.Recalculate(cltest.Head(1))
	gas.SimulateStart(t, bhe)

	t.Run("if estimator is missing base fee and gas bumping is enabled", func(t *testing.T) {
		cfg.EvmGasBumpThresholdF = uint64(1)

		_, _, err := bhe.GetDynamicFee(testutils.Context(t), 100000, maxGasPrice)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "BlockHistoryEstimator: no value for latest block base fee; cannot estimate EIP-1559 base fee. Are you trying to run with EIP1559 enabled on a non-EIP1559 chain?")
	})

	t.Run("if estimator is missing base fee and gas bumping is disabled", func(t *testing.T) {
		cfg.EvmGasBumpThresholdF = uint64(0)

		fee, limit, err := bhe.GetDynamicFee(testutils.Context(t), 100000, maxGasPrice)
		require.NoError(t, err)
		assert.Equal(t, gas.DynamicFee{FeeCap: maxGasPrice, TipCap: assets.NewWeiI(6000)}, fee)
		assert.Equal(t, 100000, int(limit))
	})

	h := cltest.Head(1)
	h.BaseFeePerGas = assets.NewWeiI(112500)
	bhe.OnNewLongestChain(testutils.Context(t), h)

	t.Run("if gas bumping is enabled", func(t *testing.T) {
		cfg.EvmGasBumpThresholdF = uint64(1)

		fee, limit, err := bhe.GetDynamicFee(testutils.Context(t), 100000, maxGasPrice)
		require.NoError(t, err)

		assert.Equal(t, gas.DynamicFee{FeeCap: assets.NewWeiI(186203), TipCap: assets.NewWeiI(6000)}, fee)
		assert.Equal(t, 100000, int(limit))
	})

	t.Run("if gas bumping is disabled", func(t *testing.T) {
		cfg.EvmGasBumpThresholdF = uint64(0)

		fee, limit, err := bhe.GetDynamicFee(testutils.Context(t), 100000, maxGasPrice)
		require.NoError(t, err)

		assert.Equal(t, gas.DynamicFee{FeeCap: maxGasPrice, TipCap: assets.NewWeiI(6000)}, fee)
		assert.Equal(t, 100000, int(limit))
	})

	t.Run("if gas bumping is enabled and local max gas price set", func(t *testing.T) {
		cfg.EvmGasBumpThresholdF = uint64(1)

		fee, limit, err := bhe.GetDynamicFee(testutils.Context(t), 100000, assets.NewWeiI(180000))
		require.NoError(t, err)

		assert.Equal(t, gas.DynamicFee{FeeCap: assets.NewWeiI(180000), TipCap: assets.NewWeiI(6000)}, fee)
		assert.Equal(t, 100000, int(limit))
	})

	t.Run("if bump threshold is 0 and local max gas price set", func(t *testing.T) {
		cfg.EvmGasBumpThresholdF = uint64(0)

		fee, limit, err := bhe.GetDynamicFee(testutils.Context(t), 100000, assets.NewWeiI(100))
		require.NoError(t, err)

		assert.Equal(t, gas.DynamicFee{FeeCap: assets.NewWeiI(100), TipCap: assets.NewWeiI(6000)}, fee)
		assert.Equal(t, 100000, int(limit))
	})

	h = cltest.Head(1)
	h.BaseFeePerGas = assets.NewWeiI(900000)
	bhe.OnNewLongestChain(testutils.Context(t), h)

	t.Run("if gas bumping is enabled and global max gas price lower than local max gas price", func(t *testing.T) {
		cfg.EvmGasBumpThresholdF = uint64(1)

		fee, limit, err := bhe.GetDynamicFee(testutils.Context(t), 100000, assets.NewWeiI(1200000))
		require.NoError(t, err)

		assert.Equal(t, gas.DynamicFee{FeeCap: assets.NewWeiI(1000000), TipCap: assets.NewWeiI(6000)}, fee)
		assert.Equal(t, 100000, int(limit))
	})
}

var _ txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash] = &MockAttempt{}

type MockAttempt struct {
	BroadcastBeforeBlockNum *int64
	Hash                    common.Hash
	TxType                  int
	GasPrice                *assets.Wei
	GasFeeCap               *assets.Wei
	GasTipCap               *assets.Wei
}

func (m *MockAttempt) Fee() (f gas.EvmFee) {
	f.Legacy = m.getGasPrice()

	d := m.dynamicFee()
	if d.FeeCap != nil && d.TipCap != nil {
		f.Dynamic = &d
	}
	return f
}

func (m *MockAttempt) getGasPrice() *assets.Wei {
	return m.GasPrice
}

func (m *MockAttempt) dynamicFee() gas.DynamicFee {
	return gas.DynamicFee{
		FeeCap: m.GasFeeCap,
		TipCap: m.GasTipCap,
	}

}

func (m *MockAttempt) GetChainSpecificGasLimit() uint32 {
	panic("not implemented") // TODO: Implement
}

func (m *MockAttempt) GetBroadcastBeforeBlockNum() *int64 {
	return m.BroadcastBeforeBlockNum
}

func (m *MockAttempt) GetHash() common.Hash {
	return m.Hash
}

func (m *MockAttempt) GetTxType() int {
	return m.TxType
}

func TestBlockHistoryEstimator_CheckConnectivity(t *testing.T) {
	cfg := newConfigWithEIP1559DynamicFeesDisabled(t)
	cfg.BlockHistoryEstimatorCheckInclusionBlocksF = uint16(4)
	lggr, obs := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	bhe := gas.BlockHistoryEstimatorFromInterface(
		gas.NewBlockHistoryEstimator(lggr, nil, cfg, *testutils.NewRandomEVMChainID()),
	)

	attempts := []txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash]{
		&MockAttempt{TxType: 0x0, Hash: NewEvmHash()},
	}

	t.Run("skips connectivity check if latest block is not present", func(t *testing.T) {
		err := bhe.CheckConnectivity(attempts)
		require.NoError(t, err)

		testutils.WaitForLogMessage(t, obs, "Latest block is unknown; skipping inclusion check")
	})

	h := cltest.Head(1)
	h.BaseFeePerGas = assets.NewWeiI(112500)
	bhe.OnNewLongestChain(testutils.Context(t), h)

	b0 := evmtypes.Block{
		Number:       0,
		Hash:         utils.NewHash(),
		Transactions: cltest.LegacyTransactionsFromGasPrices(),
	}
	gas.SetRollingBlockHistory(bhe, []evmtypes.Block{b0})

	t.Run("skips connectivity check if block history has insufficient size", func(t *testing.T) {
		err := bhe.CheckConnectivity(attempts)
		require.NoError(t, err)

		testutils.WaitForLogMessage(t, obs, "Block history in memory with length 1 is insufficient to determine whether transaction should have been included within the past 4 blocks")
	})

	t.Run("skips connectivity check if attempts is nil or empty", func(t *testing.T) {
		err := bhe.CheckConnectivity(nil)
		require.NoError(t, err)
	})

	b1 := evmtypes.Block{
		Number:       1,
		Hash:         utils.NewHash(),
		ParentHash:   b0.Hash,
		Transactions: cltest.LegacyTransactionsFromGasPrices(),
	}
	b2 := evmtypes.Block{
		Number:       2,
		Hash:         utils.NewHash(),
		ParentHash:   b1.Hash,
		Transactions: cltest.LegacyTransactionsFromGasPrices(),
	}
	b3 := evmtypes.Block{
		Number:       3,
		Hash:         utils.NewHash(),
		ParentHash:   b2.Hash,
		Transactions: cltest.LegacyTransactionsFromGasPrices(),
	}
	gas.SetRollingBlockHistory(bhe, []evmtypes.Block{b0, b1, b2, b3})
	h = cltest.Head(5)
	h.BaseFeePerGas = assets.NewWeiI(112500)
	bhe.OnNewLongestChain(testutils.Context(t), h)

	t.Run("returns error if one of the supplied attempts is missing BroadcastBeforeBlockNum", func(t *testing.T) {
		err := bhe.CheckConnectivity(attempts)
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("BroadcastBeforeBlockNum was unexpectedly nil for attempt %s", attempts[0].GetHash()))
	})

	num := int64(0)
	hash := utils.NewHash()
	attempts = []txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash]{
		&MockAttempt{TxType: 0x3, BroadcastBeforeBlockNum: &num, Hash: hash},
	}

	t.Run("returns error if one of the supplied attempts has an unknown transaction type", func(t *testing.T) {
		err := bhe.CheckConnectivity(attempts)
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("attempt %s has unknown transaction type 0x3", hash))
	})

	attempts = []txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash]{
		&MockAttempt{TxType: 0x0, BroadcastBeforeBlockNum: &num, Hash: hash},
	}

	t.Run("skips connectivity check if no transactions are suitable", func(t *testing.T) {
		err := bhe.CheckConnectivity(attempts)
		require.NoError(t, err)

		testutils.WaitForLogMessage(t, obs, fmt.Sprintf("no suitable transactions found to verify if transaction %s has been included within expected inclusion blocks of 4", hash))
	})

	t.Run("in legacy mode", func(t *testing.T) {
		b0 = evmtypes.Block{
			Number:       0,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(1000),
		}
		b1 = evmtypes.Block{
			Number:       1,
			Hash:         utils.NewHash(),
			ParentHash:   b0.Hash,
			Transactions: cltest.LegacyTransactionsFromGasPrices(2, 3, 4, 5, 6),
		}
		b2 = evmtypes.Block{
			Number:       2,
			Hash:         utils.NewHash(),
			ParentHash:   b1.Hash,
			Transactions: cltest.LegacyTransactionsFromGasPrices(4, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 9),
		}
		b3 = evmtypes.Block{
			Number:       3,
			Hash:         utils.NewHash(),
			ParentHash:   b2.Hash,
			Transactions: cltest.LegacyTransactionsFromGasPrices(3, 4, 5, 6, 7),
		}
		gas.SetRollingBlockHistory(bhe, []evmtypes.Block{b0, b1, b2, b3})

		attempts = []txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash]{
			&MockAttempt{TxType: 0x0, Hash: NewEvmHash(), GasPrice: assets.NewWeiI(1000), BroadcastBeforeBlockNum: testutils.Ptr(int64(4))}, // This is very expensive but will be ignored due to BroadcastBeforeBlockNum being too recent
			&MockAttempt{TxType: 0x0, Hash: NewEvmHash(), GasPrice: assets.NewWeiI(3), BroadcastBeforeBlockNum: testutils.Ptr(int64(0))},
			&MockAttempt{TxType: 0x0, Hash: NewEvmHash(), GasPrice: assets.NewWeiI(5), BroadcastBeforeBlockNum: testutils.Ptr(int64(1))},
			&MockAttempt{TxType: 0x0, Hash: NewEvmHash(), GasPrice: assets.NewWeiI(7), BroadcastBeforeBlockNum: testutils.Ptr(int64(1))},
		}

		t.Run("passes check if all blocks have percentile price higher or exactly at the highest transaction gas price", func(t *testing.T) {
			cfg.BlockHistoryEstimatorCheckInclusionBlocksF = 3
			cfg.BlockHistoryEstimatorCheckInclusionPercentileF = 80 // percentile price is 7 wei

			err := bhe.CheckConnectivity(attempts)
			require.NoError(t, err)
		})
		t.Run("fails check if one or more blocks has percentile price higher than highest transaction gas price", func(t *testing.T) {
			cfg.BlockHistoryEstimatorCheckInclusionBlocksF = 3
			cfg.BlockHistoryEstimatorCheckInclusionPercentileF = 40 // percentile price is 5 wei

			err := bhe.CheckConnectivity(attempts)
			require.Error(t, err)
			assert.Contains(t, err.Error(), fmt.Sprintf("transaction %s has gas price of 7 wei, which is above percentile=40%% (percentile price: 5 wei) for blocks 2 thru 0 (checking 3 blocks)", attempts[3].GetHash()))
			require.ErrorIs(t, err, gas.ErrConnectivity)
		})

		t.Run("fails check if one or more blocks has percentile price higher than any transaction gas price", func(t *testing.T) {
			cfg.BlockHistoryEstimatorCheckInclusionBlocksF = 3
			cfg.BlockHistoryEstimatorCheckInclusionPercentileF = 30 // percentile price is 4 wei

			err := bhe.CheckConnectivity(attempts)
			require.Error(t, err)
			assert.Contains(t, err.Error(), fmt.Sprintf("transaction %s has gas price of 5 wei, which is above percentile=30%% (percentile price: 4 wei) for blocks 2 thru 0 (checking 3 blocks)", attempts[2].GetHash()))

			cfg.BlockHistoryEstimatorCheckInclusionBlocksF = 3
			cfg.BlockHistoryEstimatorCheckInclusionPercentileF = 5 // percentile price is 2 wei

			err = bhe.CheckConnectivity(attempts)
			require.Error(t, err)
			assert.Contains(t, err.Error(), fmt.Sprintf("transaction %s has gas price of 3 wei, which is above percentile=5%% (percentile price: 2 wei) for blocks 2 thru 0 (checking 3 blocks)", attempts[1].GetHash()))
			require.ErrorIs(t, err, gas.ErrConnectivity)
		})
	})

	t.Run("handles mixed legacy and EIP-1559 transactions", func(t *testing.T) {
		b0 = evmtypes.Block{
			BaseFeePerGas: assets.NewWeiI(1),
			Number:        3,
			Hash:          utils.NewHash(),
			Transactions:  append(cltest.LegacyTransactionsFromGasPrices(1, 2, 3, 4, 5), cltest.DynamicFeeTransactionsFromTipCaps(6, 7, 8, 9, 10)...),
		}
		gas.SetRollingBlockHistory(bhe, []evmtypes.Block{b0})

		attempts = []txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash]{
			&MockAttempt{TxType: 0x2, Hash: NewEvmHash(), GasFeeCap: assets.NewWeiI(1), GasTipCap: assets.NewWeiI(3), BroadcastBeforeBlockNum: testutils.Ptr(int64(0))},
			&MockAttempt{TxType: 0x0, Hash: NewEvmHash(), GasPrice: assets.NewWeiI(10), BroadcastBeforeBlockNum: testutils.Ptr(int64(0))},
		}

		t.Run("passes check if both transactions are ok", func(t *testing.T) {
			cfg.BlockHistoryEstimatorCheckInclusionBlocksF = 1
			cfg.BlockHistoryEstimatorCheckInclusionPercentileF = 90 // percentile price is 5 wei

			err := bhe.CheckConnectivity(attempts)
			require.NoError(t, err)
		})
		t.Run("fails check if legacy transaction fails", func(t *testing.T) {
			cfg.BlockHistoryEstimatorCheckInclusionBlocksF = 1
			cfg.BlockHistoryEstimatorCheckInclusionPercentileF = 60

			err := bhe.CheckConnectivity(attempts)
			require.Error(t, err)
			assert.Contains(t, err.Error(), fmt.Sprintf("transaction %s has gas price of 10 wei, which is above percentile=60%% (percentile price: 7 wei) for blocks 3 thru 3 (checking 1 blocks)", attempts[1].GetHash()))
			require.ErrorIs(t, err, gas.ErrConnectivity)
		})

		attempts = []txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash]{
			&MockAttempt{TxType: 0x2, Hash: NewEvmHash(), GasFeeCap: assets.NewWeiI(11), GasTipCap: assets.NewWeiI(10), BroadcastBeforeBlockNum: testutils.Ptr(int64(0))},
			&MockAttempt{TxType: 0x0, Hash: NewEvmHash(), GasPrice: assets.NewWeiI(3), BroadcastBeforeBlockNum: testutils.Ptr(int64(0))},
		}

		t.Run("fails check if dynamic fee transaction fails", func(t *testing.T) {
			gas.SetRollingBlockHistory(bhe, []evmtypes.Block{b0})
			cfg.BlockHistoryEstimatorCheckInclusionBlocksF = 1
			cfg.BlockHistoryEstimatorCheckInclusionPercentileF = 60

			err := bhe.CheckConnectivity(attempts)
			require.Error(t, err)
			assert.Contains(t, err.Error(), fmt.Sprintf("transaction %s has tip cap of 10 wei, which is above percentile=60%% (percentile tip cap: 6 wei) for blocks 3 thru 3 (checking 1 blocks)", attempts[0].GetHash()))
			require.ErrorIs(t, err, gas.ErrConnectivity)
		})

	})

	t.Run("in EIP-1559 mode", func(t *testing.T) {
		b0 = evmtypes.Block{
			BaseFeePerGas: assets.NewWeiI(5),
			Number:        0,
			Hash:          utils.NewHash(),
			Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(1000),
		}
		b1 = evmtypes.Block{
			BaseFeePerGas: assets.NewWeiI(8),
			Number:        1,
			Hash:          utils.NewHash(),
			ParentHash:    b0.Hash,
			Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(2, 3, 4, 5, 6),
		}
		b2 = evmtypes.Block{
			BaseFeePerGas: assets.NewWeiI(13),
			Number:        2,
			Hash:          utils.NewHash(),
			ParentHash:    b1.Hash,
			Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(4, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 9),
		}
		b3 = evmtypes.Block{
			BaseFeePerGas: assets.NewWeiI(21),
			Number:        3,
			Hash:          utils.NewHash(),
			ParentHash:    b2.Hash,
			Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(3, 4, 5, 6, 7),
		}
		blocks := []evmtypes.Block{b0, b1, b2, b3}
		gas.SetRollingBlockHistory(bhe, blocks)

		attempts = []txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash]{
			&MockAttempt{TxType: 0x2, Hash: NewEvmHash(), GasFeeCap: assets.NewWeiI(30), GasTipCap: assets.NewWeiI(1000), BroadcastBeforeBlockNum: testutils.Ptr(int64(4))}, // This is very expensive but will be ignored due to BroadcastBeforeBlockNum being too recent
			&MockAttempt{TxType: 0x2, Hash: NewEvmHash(), GasFeeCap: assets.NewWeiI(30), GasTipCap: assets.NewWeiI(3), BroadcastBeforeBlockNum: testutils.Ptr(int64(0))},
			&MockAttempt{TxType: 0x2, Hash: NewEvmHash(), GasFeeCap: assets.NewWeiI(30), GasTipCap: assets.NewWeiI(5), BroadcastBeforeBlockNum: testutils.Ptr(int64(1))},
			&MockAttempt{TxType: 0x2, Hash: NewEvmHash(), GasFeeCap: assets.NewWeiI(30), GasTipCap: assets.NewWeiI(7), BroadcastBeforeBlockNum: testutils.Ptr(int64(1))},
		}

		t.Run("passes check if all blocks have 90th percentile price higher than highest transaction tip cap", func(t *testing.T) {
			cfg.BlockHistoryEstimatorCheckInclusionBlocksF = 3
			cfg.BlockHistoryEstimatorCheckInclusionPercentileF = 80

			err := bhe.CheckConnectivity(attempts)
			require.NoError(t, err)
		})

		t.Run("fails check if one or more blocks has percentile tip cap higher than any transaction tip cap, and base fee higher than the block base fee", func(t *testing.T) {
			cfg.BlockHistoryEstimatorCheckInclusionBlocksF = 3
			cfg.BlockHistoryEstimatorCheckInclusionPercentileF = 20

			err := bhe.CheckConnectivity(attempts)
			require.Error(t, err)
			assert.Contains(t, err.Error(), fmt.Sprintf("transaction %s has tip cap of 5 wei, which is above percentile=20%% (percentile tip cap: 4 wei) for blocks 2 thru 0 (checking 3 blocks)", attempts[2].GetHash()))
			require.ErrorIs(t, err, gas.ErrConnectivity)

			cfg.BlockHistoryEstimatorCheckInclusionBlocksF = 3
			cfg.BlockHistoryEstimatorCheckInclusionPercentileF = 5

			err = bhe.CheckConnectivity(attempts)
			require.Error(t, err)
			assert.Contains(t, err.Error(), fmt.Sprintf("transaction %s has tip cap of 3 wei, which is above percentile=5%% (percentile tip cap: 2 wei) for blocks 2 thru 0 (checking 3 blocks)", attempts[1].GetHash()))
			require.ErrorIs(t, err, gas.ErrConnectivity)
		})

		t.Run("passes check if, for at least one block, feecap < tipcap+basefee, even if percentile is not reached", func(t *testing.T) {
			cfg.BlockHistoryEstimatorCheckInclusionBlocksF = 3
			cfg.BlockHistoryEstimatorCheckInclusionPercentileF = 5

			attempts = []txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash]{
				&MockAttempt{TxType: 0x2, Hash: NewEvmHash(), GasFeeCap: assets.NewWeiI(4), GasTipCap: assets.NewWeiI(7), BroadcastBeforeBlockNum: testutils.Ptr(int64(1))},
			}

			err := bhe.CheckConnectivity(attempts)
			require.NoError(t, err)
		})
	})
}

func TestBlockHistoryEstimator_Bumps(t *testing.T) {
	t.Parallel()
	maxGasPrice := assets.NewWeiI(1000000)

	t.Run("BumpLegacyGas checks connectivity", func(t *testing.T) {
		cfg := newConfigWithEIP1559DynamicFeesDisabled(t)
		cfg.BlockHistoryEstimatorCheckInclusionBlocksF = 1
		cfg.BlockHistoryEstimatorCheckInclusionPercentileF = 10
		cfg.EvmGasBumpPercentF = 10
		cfg.EvmGasBumpWeiF = assets.NewWeiI(150)
		cfg.EvmMaxGasPriceWeiF = maxGasPrice
		cfg.EvmGasLimitMultiplierF = float32(1.1)
		bhe := newBlockHistoryEstimator(t, nil, cfg)

		b1 := evmtypes.Block{
			Number:       1,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(1),
		}
		gas.SetRollingBlockHistory(bhe, []evmtypes.Block{b1})
		head := cltest.Head(1)
		bhe.OnNewLongestChain(testutils.Context(t), head)

		attempts := []txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash]{
			&MockAttempt{TxType: 0x0, Hash: NewEvmHash(), GasPrice: assets.NewWeiI(1000), BroadcastBeforeBlockNum: testutils.Ptr(int64(0))},
		}

		_, _, err := bhe.BumpLegacyGas(testutils.Context(t), assets.NewWeiI(42), 100000, maxGasPrice, gas.MakeEvmPriorAttempts(attempts))
		require.Error(t, err)
		assert.True(t, errors.Is(err, gas.ErrConnectivity))
		assert.Contains(t, err.Error(), fmt.Sprintf("transaction %s has gas price of 1 kwei, which is above percentile=10%% (percentile price: 1 wei) for blocks 1 thru 1 (checking 1 blocks)", attempts[0].GetHash()))
	})

	t.Run("BumpLegacyGas calls BumpLegacyGasPriceOnly with proper current gas price", func(t *testing.T) {
		cfg := newConfigWithEIP1559DynamicFeesDisabled(t)
		cfg.EvmGasBumpPercentF = 10
		cfg.EvmGasBumpWeiF = assets.NewWeiI(150)
		cfg.EvmMaxGasPriceWeiF = maxGasPrice
		cfg.EvmGasLimitMultiplierF = float32(1.1)
		bhe := newBlockHistoryEstimator(t, nil, cfg)

		t.Run("ignores nil current gas price", func(t *testing.T) {
			gasPrice, gasLimit, err := bhe.BumpLegacyGas(testutils.Context(t), assets.NewWeiI(42), 100000, maxGasPrice, nil)
			require.NoError(t, err)

			expectedGasPrice, expectedGasLimit, err := gas.BumpLegacyGasPriceOnly(cfg, logger.TestLogger(t), nil, assets.NewWeiI(42), 100000, maxGasPrice)
			require.NoError(t, err)

			assert.Equal(t, expectedGasLimit, gasLimit)
			assert.Equal(t, expectedGasPrice, gasPrice)
		})

		t.Run("ignores current gas price > max gas price", func(t *testing.T) {
			gasPrice, gasLimit, err := bhe.BumpLegacyGas(testutils.Context(t), assets.NewWeiI(42), 100000, maxGasPrice, nil)
			require.NoError(t, err)

			massive := assets.NewWeiI(100000000000000)
			gas.SetGasPrice(bhe, massive)

			expectedGasPrice, expectedGasLimit, err := gas.BumpLegacyGasPriceOnly(cfg, logger.TestLogger(t), massive, assets.NewWeiI(42), 100000, maxGasPrice)
			require.NoError(t, err)

			assert.Equal(t, expectedGasLimit, gasLimit)
			assert.Equal(t, expectedGasPrice, gasPrice)
		})

		t.Run("ignores current gas price < bumped gas price", func(t *testing.T) {
			gas.SetGasPrice(bhe, assets.NewWeiI(191))

			gasPrice, gasLimit, err := bhe.BumpLegacyGas(testutils.Context(t), assets.NewWeiI(42), 100000, maxGasPrice, nil)
			require.NoError(t, err)

			assert.Equal(t, 110000, int(gasLimit))
			assert.Equal(t, assets.NewWeiI(192), gasPrice)
		})

		t.Run("uses current gas price > bumped gas price", func(t *testing.T) {
			gas.SetGasPrice(bhe, assets.NewWeiI(193))

			gasPrice, gasLimit, err := bhe.BumpLegacyGas(testutils.Context(t), assets.NewWeiI(42), 100000, maxGasPrice, nil)
			require.NoError(t, err)

			assert.Equal(t, 110000, int(gasLimit))
			assert.Equal(t, assets.NewWeiI(193), gasPrice)
		})

		t.Run("bumped gas price > max gas price", func(t *testing.T) {
			gas.SetGasPrice(bhe, assets.NewWeiI(191))

			gasPrice, gasLimit, err := bhe.BumpLegacyGas(testutils.Context(t), assets.NewWeiI(42), 100000, assets.NewWeiI(100), nil)
			require.Error(t, err)

			assert.Nil(t, gasPrice)
			assert.Equal(t, 0, int(gasLimit))
			assert.Contains(t, err.Error(), "bumped gas price of 192 wei would exceed configured max gas price of 100 wei (original price was 42 wei).")
		})

		t.Run("current gas price > max gas price", func(t *testing.T) {
			gas.SetGasPrice(bhe, assets.NewWeiI(193))

			gasPrice, gasLimit, err := bhe.BumpLegacyGas(testutils.Context(t), assets.NewWeiI(42), 100000, assets.NewWeiI(100), nil)
			require.Error(t, err)

			assert.Nil(t, gasPrice)
			assert.Equal(t, 0, int(gasLimit))
			assert.Contains(t, err.Error(), "bumped gas price of 192 wei would exceed configured max gas price of 100 wei (original price was 42 wei).")
		})
	})

	t.Run("BumpDynamicFee checks connectivity", func(t *testing.T) {
		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)
		cfg.BlockHistoryEstimatorCheckInclusionBlocksF = 1
		cfg.BlockHistoryEstimatorCheckInclusionPercentileF = 10
		cfg.EvmGasBumpPercentF = 10
		cfg.EvmGasBumpWeiF = assets.NewWeiI(150)
		cfg.EvmMaxGasPriceWeiF = maxGasPrice
		cfg.EvmGasLimitMultiplierF = float32(1.1)
		bhe := newBlockHistoryEstimator(t, nil, cfg)

		b1 := evmtypes.Block{
			BaseFeePerGas: assets.NewWeiI(1),
			Number:        1,
			Hash:          utils.NewHash(),
			Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(1),
		}
		gas.SetRollingBlockHistory(bhe, []evmtypes.Block{b1})
		head := cltest.Head(1)
		bhe.OnNewLongestChain(testutils.Context(t), head)

		originalFee := gas.DynamicFee{FeeCap: assets.NewWeiI(100), TipCap: assets.NewWeiI(25)}
		attempts := []txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash]{
			&MockAttempt{TxType: 0x2, Hash: NewEvmHash(), GasTipCap: originalFee.TipCap, GasFeeCap: originalFee.FeeCap, BroadcastBeforeBlockNum: testutils.Ptr(int64(0))},
		}

		_, _, err := bhe.BumpDynamicFee(testutils.Context(t), originalFee, 100000, maxGasPrice, gas.MakeEvmPriorAttempts(attempts))
		require.Error(t, err)
		assert.True(t, errors.Is(err, gas.ErrConnectivity))
		assert.Contains(t, err.Error(), fmt.Sprintf("transaction %s has tip cap of 25 wei, which is above percentile=10%% (percentile tip cap: 1 wei) for blocks 1 thru 1 (checking 1 blocks)", attempts[0].GetHash()))
	})

	t.Run("BumpDynamicFee bumps the fee", func(t *testing.T) {
		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)
		cfg.EvmGasBumpPercentF = 10
		cfg.EvmGasBumpWeiF = assets.NewWeiI(150)
		cfg.EvmMaxGasPriceWeiF = maxGasPrice
		cfg.EvmGasLimitMultiplierF = float32(1.1)
		cfg.EvmGasTipCapDefaultF = assets.NewWeiI(52)

		bhe := newBlockHistoryEstimator(t, nil, cfg)

		t.Run("when current tip cap is nil", func(t *testing.T) {
			originalFee := gas.DynamicFee{FeeCap: assets.NewWeiI(100), TipCap: assets.NewWeiI(25)}
			fee, gasLimit, err := bhe.BumpDynamicFee(testutils.Context(t), originalFee, 100000, maxGasPrice, nil)
			require.NoError(t, err)

			assert.Equal(t, 110000, int(gasLimit))
			assert.Equal(t, gas.DynamicFee{FeeCap: assets.NewWeiI(250), TipCap: assets.NewWeiI(202)}, fee)
		})
		t.Run("ignores current tip cap that is smaller than original fee with bump applied", func(t *testing.T) {
			gas.SetTipCap(bhe, assets.NewWeiI(201))

			originalFee := gas.DynamicFee{FeeCap: assets.NewWeiI(100), TipCap: assets.NewWeiI(25)}
			fee, gasLimit, err := bhe.BumpDynamicFee(testutils.Context(t), originalFee, 100000, maxGasPrice, nil)
			require.NoError(t, err)

			assert.Equal(t, 110000, int(gasLimit))
			assert.Equal(t, gas.DynamicFee{FeeCap: assets.NewWeiI(250), TipCap: assets.NewWeiI(202)}, fee)
		})
		t.Run("uses current tip cap that is larger than original fee with bump applied", func(t *testing.T) {
			gas.SetTipCap(bhe, assets.NewWeiI(203))

			originalFee := gas.DynamicFee{FeeCap: assets.NewWeiI(100), TipCap: assets.NewWeiI(25)}
			fee, gasLimit, err := bhe.BumpDynamicFee(testutils.Context(t), originalFee, 100000, maxGasPrice, nil)
			require.NoError(t, err)

			assert.Equal(t, 110000, int(gasLimit))
			assert.Equal(t, gas.DynamicFee{FeeCap: assets.NewWeiI(250), TipCap: assets.NewWeiI(203)}, fee)
		})
		t.Run("ignores absurdly large current tip cap", func(t *testing.T) {
			gas.SetTipCap(bhe, assets.NewWeiI(1000000000000000))

			originalFee := gas.DynamicFee{FeeCap: assets.NewWeiI(100), TipCap: assets.NewWeiI(25)}
			fee, gasLimit, err := bhe.BumpDynamicFee(testutils.Context(t), originalFee, 100000, maxGasPrice, nil)
			require.NoError(t, err)

			assert.Equal(t, 110000, int(gasLimit))
			assert.Equal(t, gas.DynamicFee{FeeCap: assets.NewWeiI(250), TipCap: assets.NewWeiI(202)}, fee)
		})

		t.Run("bumped tip cap price > max gas price", func(t *testing.T) {
			gas.SetTipCap(bhe, assets.NewWeiI(203))

			originalFee := gas.DynamicFee{FeeCap: assets.NewWeiI(100), TipCap: assets.NewWeiI(990000)}
			fee, gasLimit, err := bhe.BumpDynamicFee(testutils.Context(t), originalFee, 100000, maxGasPrice, nil)
			require.Error(t, err)

			assert.Equal(t, 0, int(gasLimit))
			assert.Equal(t, gas.DynamicFee{}, fee)
			assert.Contains(t, err.Error(), "bumped tip cap of 1.089 mwei would exceed configured max gas price of 1 mwei (original fee: tip cap 990 kwei, fee cap 100 wei)")
		})

		t.Run("bumped fee cap price > max gas price", func(t *testing.T) {
			gas.SetTipCap(bhe, assets.NewWeiI(203))

			originalFee := gas.DynamicFee{FeeCap: assets.NewWeiI(990000), TipCap: assets.NewWeiI(25)}
			fee, gasLimit, err := bhe.BumpDynamicFee(testutils.Context(t), originalFee, 100000, maxGasPrice, nil)
			require.Error(t, err)

			assert.Equal(t, 0, int(gasLimit))
			assert.Equal(t, gas.DynamicFee{}, fee)
			assert.Contains(t, err.Error(), "bumped fee cap of 1.089 mwei would exceed configured max gas price of 1 mwei (original fee: tip cap 25 wei, fee cap 990 kwei)")
		})
	})
}
