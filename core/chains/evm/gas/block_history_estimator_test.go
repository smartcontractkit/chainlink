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

	"github.com/smartcontractkit/chainlink/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	gumocks "github.com/smartcontractkit/chainlink/core/chains/evm/gas/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	cfg "github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func newConfigWithEIP1559DynamicFeesEnabled(t *testing.T) *gumocks.Config {
	config := new(gumocks.Config)
	config.Test(t)
	config.On("EvmEIP1559DynamicFees").Maybe().Return(true)
	config.On("ChainType").Maybe().Return(cfg.ChainType(""))
	return config
}

func newConfigWithEIP1559DynamicFeesDisabled(t *testing.T) *gumocks.Config {
	config := new(gumocks.Config)
	config.Test(t)
	config.On("EvmEIP1559DynamicFees").Maybe().Return(false)
	config.On("ChainType").Maybe().Return(cfg.ChainType(""))
	return config
}

func newBlockHistoryEstimatorWithChainID(t *testing.T, c evmclient.Client, cfg gas.Config, cid big.Int) gas.Estimator {
	return gas.NewBlockHistoryEstimator(logger.TestLogger(t), c, cfg, cid)
}

func newBlockHistoryEstimator(t *testing.T, c evmclient.Client, cfg gas.Config) *gas.BlockHistoryEstimator {
	iface := newBlockHistoryEstimatorWithChainID(t, c, cfg, cltest.FixtureChainID)
	return gas.BlockHistoryEstimatorFromInterface(iface)
}

func TestBlockHistoryEstimator_Start(t *testing.T) {
	t.Parallel()

	config := newConfigWithEIP1559DynamicFeesEnabled(t)

	var batchSize uint32 = 0
	var blockDelay uint16 = 0
	var historySize uint16 = 2
	var percentile uint16 = 35
	minGasPrice := big.NewInt(1)

	config.On("BlockHistoryEstimatorBatchSize").Return(batchSize)
	config.On("BlockHistoryEstimatorBlockDelay").Return(blockDelay)
	config.On("BlockHistoryEstimatorBlockHistorySize").Return(historySize)
	config.On("BlockHistoryEstimatorTransactionPercentile").Maybe().Return(percentile)
	config.On("EvmGasLimitMultiplier").Maybe().Return(float32(1))
	config.On("EvmMinGasPriceWei").Maybe().Return(minGasPrice)
	config.On("EvmEIP1559DynamicFees").Maybe().Return(true)

	t.Run("loads initial state", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		h := &evmtypes.Head{Hash: utils.NewHash(), Number: 42, BaseFeePerGas: utils.NewBigI(420)}
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(h, nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == "0x2a" && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&gas.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == "0x29" && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&gas.Block{})
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &gas.Block{
				Number: 42,
				Hash:   utils.NewHash(),
			}
			elems[1].Result = &gas.Block{
				Number: 41,
				Hash:   utils.NewHash(),
			}
		}).Once()

		err := bhe.Start(testutils.Context(t))
		require.NoError(t, err)

		assert.Len(t, bhe.RollingBlockHistory(), 2)
		assert.Equal(t, int(bhe.RollingBlockHistory()[0].Number), 41)
		assert.Equal(t, int(bhe.RollingBlockHistory()[1].Number), 42)

		assert.Equal(t, big.NewInt(420), gas.GetLatestBaseFee(bhe))

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("starts and loads partial history if fetch context times out", func(t *testing.T) {
		cfg := newConfigWithEIP1559DynamicFeesEnabled(t)

		cfg.On("BlockHistoryEstimatorBatchSize").Return(uint32(1))
		cfg.On("BlockHistoryEstimatorBlockDelay").Return(blockDelay)
		cfg.On("BlockHistoryEstimatorBlockHistorySize").Return(historySize)
		cfg.On("BlockHistoryEstimatorTransactionPercentile").Maybe().Return(percentile)
		cfg.On("EvmGasLimitMultiplier").Maybe().Return(float32(1))
		cfg.On("EvmMinGasPriceWei").Maybe().Return(minGasPrice)
		cfg.On("EvmEIP1559DynamicFees").Maybe().Return(true)
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		bhe := newBlockHistoryEstimator(t, ethClient, cfg)

		h := &evmtypes.Head{Hash: utils.NewHash(), Number: 42, BaseFeePerGas: utils.NewBigI(420)}
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(h, nil)
		// First succeeds (42)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == gas.Int64ToHex(42) && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&gas.Block{})
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &gas.Block{
				Number: 42,
				Hash:   utils.NewHash(),
			}
		}).Once()
		// Second fails (41)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == gas.Int64ToHex(41) && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&gas.Block{})
		})).Return(errors.Wrap(context.DeadlineExceeded, "some error message")).Once()

		err := bhe.Start(testutils.Context(t))
		require.NoError(t, err)

		require.Len(t, bhe.RollingBlockHistory(), 1)
		assert.Equal(t, int(bhe.RollingBlockHistory()[0].Number), 42)

		assert.Equal(t, big.NewInt(420), gas.GetLatestBaseFee(bhe))

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("boots even if initial batch call returns nothing", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		h := &evmtypes.Head{Hash: utils.NewHash(), Number: 42}
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(h, nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == int(historySize)
		})).Return(nil)

		err := bhe.Start(testutils.Context(t))
		require.NoError(t, err)

		// non-eip1559 block
		assert.Nil(t, gas.GetLatestBaseFee(bhe))

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("starts anyway if fetching latest head fails", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(nil, errors.New("something exploded"))

		err := bhe.Start(testutils.Context(t))
		require.NoError(t, err)

		assert.Nil(t, gas.GetLatestBaseFee(bhe))

		_, _, err = bhe.GetLegacyGas(make([]byte, 0), 100)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has not finished the first gas estimation yet")

		_, _, err = bhe.GetDynamicFee(100)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has not finished the first gas estimation yet")

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("starts anyway if fetching first fetch fails, but errors on estimation", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		h := &evmtypes.Head{Hash: utils.NewHash(), Number: 42, BaseFeePerGas: utils.NewBigI(420)}
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(h, nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.Anything).Return(errors.New("something went wrong"))

		err := bhe.Start(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, big.NewInt(420), gas.GetLatestBaseFee(bhe))

		_, _, err = bhe.GetLegacyGas(make([]byte, 0), 100)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has not finished the first gas estimation yet")

		_, _, err = bhe.GetDynamicFee(100)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has not finished the first gas estimation yet")

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("returns error if main context is cancelled", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		h := &evmtypes.Head{Hash: utils.NewHash(), Number: 42, BaseFeePerGas: utils.NewBigI(420)}
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(h, nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.Anything).Return(errors.New("this error doesn't matter"))

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := bhe.Start(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "context canceled")

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("starts anyway even if the fetch context is cancelled due to taking longer than the MaxStartTime", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		h := &evmtypes.Head{Hash: utils.NewHash(), Number: 42, BaseFeePerGas: utils.NewBigI(420)}
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(h, nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.Anything).Return(errors.New("this error doesn't matter")).Run(func(_ mock.Arguments) {
			time.Sleep(gas.MaxStartTime + 1*time.Second)
		})

		err := bhe.Start(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, big.NewInt(420), gas.GetLatestBaseFee(bhe))

		_, _, err = bhe.GetLegacyGas(make([]byte, 0), 100)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has not finished the first gas estimation yet")

		_, _, err = bhe.GetDynamicFee(100)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has not finished the first gas estimation yet")

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})
}

func TestBlockHistoryEstimator_OnNewLongestChain(t *testing.T) {
	bhe := newBlockHistoryEstimator(t, nil, nil)

	assert.Nil(t, gas.GetLatestBaseFee(bhe))

	// non EIP-1559 block
	h := cltest.Head(1)
	bhe.OnNewLongestChain(context.Background(), h)
	assert.Nil(t, gas.GetLatestBaseFee(bhe))

	// EIP-1559 block
	h = cltest.Head(2)
	h.BaseFeePerGas = utils.NewBigI(500)
	bhe.OnNewLongestChain(context.Background(), h)

	assert.Equal(t, big.NewInt(500), gas.GetLatestBaseFee(bhe))
}

func TestBlockHistoryEstimator_FetchBlocks(t *testing.T) {
	t.Parallel()

	t.Run("with history size of 0, errors", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		config := newConfigWithEIP1559DynamicFeesEnabled(t)
		bhe := newBlockHistoryEstimator(t, ethClient, config)

		var blockDelay uint16 = 3
		var historySize uint16 = 0
		config.On("BlockHistoryEstimatorBlockDelay").Return(blockDelay)
		config.On("BlockHistoryEstimatorBlockHistorySize").Return(historySize)

		head := cltest.Head(42)
		err := bhe.FetchBlocks(context.Background(), head)
		require.Error(t, err)
		require.EqualError(t, err, "BlockHistoryEstimator: history size must be > 0, got: 0")
	})

	t.Run("with current block height less than block delay does nothing", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		config := newConfigWithEIP1559DynamicFeesEnabled(t)
		bhe := newBlockHistoryEstimator(t, ethClient, config)

		var blockDelay uint16 = 3
		var historySize uint16 = 1
		config.On("BlockHistoryEstimatorBlockDelay").Return(blockDelay)
		config.On("BlockHistoryEstimatorBlockHistorySize").Return(historySize)

		for i := -1; i < 3; i++ {
			head := cltest.Head(i)
			err := bhe.FetchBlocks(context.Background(), head)
			require.Error(t, err)
			require.EqualError(t, err, fmt.Sprintf("BlockHistoryEstimator: cannot fetch, current block height %v is lower than GAS_UPDATER_BLOCK_DELAY=3", i))
		}

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("with error retrieving blocks returns error", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		config := newConfigWithEIP1559DynamicFeesEnabled(t)
		bhe := newBlockHistoryEstimator(t, ethClient, config)

		var blockDelay uint16 = 3
		var historySize uint16 = 3
		var batchSize uint32 = 0
		config.On("BlockHistoryEstimatorBlockDelay").Return(blockDelay)
		config.On("BlockHistoryEstimatorBlockHistorySize").Return(historySize)
		config.On("BlockHistoryEstimatorBatchSize").Return(batchSize)

		ethClient.On("BatchCallContext", mock.Anything, mock.Anything).Return(errors.New("something exploded"))

		err := bhe.FetchBlocks(context.Background(), cltest.Head(42))
		require.Error(t, err)
		assert.EqualError(t, err, "BlockHistoryEstimator#fetchBlocks error fetching blocks with BatchCallContext: something exploded")

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("batch fetches heads and transactions and sets them on the block history estimator instance", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		config := newConfigWithEIP1559DynamicFeesEnabled(t)
		bhe := newBlockHistoryEstimator(t, ethClient, config)

		var blockDelay uint16 = 0
		var historySize uint16 = 3
		var batchSize uint32 = 2
		config.On("BlockHistoryEstimatorBlockDelay").Return(blockDelay)
		config.On("BlockHistoryEstimatorBlockHistorySize").Return(historySize)
		// Test batching
		config.On("BlockHistoryEstimatorBatchSize").Return(batchSize)

		b41 := gas.Block{
			Number:       41,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(1, 2),
		}
		b42 := gas.Block{
			Number:       42,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(3),
		}
		b43 := gas.Block{
			Number:       43,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(),
		}

		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == gas.Int64ToHex(43) && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&gas.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == gas.Int64ToHex(42) && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&gas.Block{})
		})).Once().Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &b43
			// This errored block (42) will be ignored
			elems[1].Error = errors.New("something went wrong")
		})
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == gas.Int64ToHex(41) && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&gas.Block{})
		})).Once().Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &b41
		})

		err := bhe.FetchBlocks(context.Background(), cltest.Head(43))
		require.NoError(t, err)

		require.Len(t, bhe.RollingBlockHistory(), 2)
		assert.Equal(t, 41, int(bhe.RollingBlockHistory()[0].Number))
		// 42 is missing because the fetch errored
		assert.Equal(t, 43, int(bhe.RollingBlockHistory()[1].Number))
		assert.Len(t, bhe.RollingBlockHistory()[0].Transactions, 2)
		assert.Len(t, bhe.RollingBlockHistory()[1].Transactions, 0)

		ethClient.AssertExpectations(t)

		// On new fetch, rolls over the history and drops the old heads

		b44 := gas.Block{
			Number:       44,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(4),
		}

		// We are gonna refetch blocks 42 and 44
		// 43 is skipped because it was already in the history
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == gas.Int64ToHex(44) && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&gas.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == gas.Int64ToHex(42) && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&gas.Block{})
		})).Once().Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &b44
			elems[1].Result = &b42
		})

		head := evmtypes.NewHead(big.NewInt(44), b44.Hash, b43.Hash, uint64(time.Now().Unix()), utils.NewBig(&cltest.FixtureChainID))
		err = bhe.FetchBlocks(context.Background(), &head)
		require.NoError(t, err)

		ethClient.AssertExpectations(t)

		require.Len(t, bhe.RollingBlockHistory(), 3)
		assert.Equal(t, 42, int(bhe.RollingBlockHistory()[0].Number))
		assert.Equal(t, 43, int(bhe.RollingBlockHistory()[1].Number))
		assert.Equal(t, 44, int(bhe.RollingBlockHistory()[2].Number))
		assert.Len(t, bhe.RollingBlockHistory()[0].Transactions, 1)
		assert.Len(t, bhe.RollingBlockHistory()[1].Transactions, 0)
		assert.Len(t, bhe.RollingBlockHistory()[2].Transactions, 1)

		config.AssertExpectations(t)
	})

	t.Run("does not refetch blocks below EVM_FINALITY_DEPTH", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		config := newConfigWithEIP1559DynamicFeesEnabled(t)
		bhe := newBlockHistoryEstimator(t, ethClient, config)

		var blockDelay uint16 = 0
		var historySize uint16 = 3
		var batchSize uint32 = 2
		config.On("BlockHistoryEstimatorBlockDelay").Return(blockDelay)
		config.On("BlockHistoryEstimatorBlockHistorySize").Return(historySize)
		config.On("BlockHistoryEstimatorBatchSize").Return(batchSize)

		b0 := gas.Block{
			Number:       0,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(9001),
		}
		b1 := gas.Block{
			Number:       1,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(9002),
		}
		blocks := []gas.Block{b0, b1}

		gas.SetRollingBlockHistory(bhe, blocks)

		b2 := gas.Block{
			Number:       2,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(1, 2),
		}
		b3 := gas.Block{
			Number:       3,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(1, 2),
		}

		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == gas.Int64ToHex(3) && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&gas.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == gas.Int64ToHex(2) && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&gas.Block{})
		})).Once().Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &b3
			elems[1].Result = &b2
		})

		head2 := evmtypes.NewHead(big.NewInt(2), b2.Hash, b1.Hash, uint64(time.Now().Unix()), utils.NewBig(&cltest.FixtureChainID))
		head3 := evmtypes.NewHead(big.NewInt(3), b3.Hash, b2.Hash, uint64(time.Now().Unix()), utils.NewBig(&cltest.FixtureChainID))
		head3.Parent = &head2
		err := bhe.FetchBlocks(context.Background(), &head3)
		require.NoError(t, err)

		require.Len(t, bhe.RollingBlockHistory(), 3)
		assert.Equal(t, 1, int(bhe.RollingBlockHistory()[0].Number))
		assert.Equal(t, 2, int(bhe.RollingBlockHistory()[1].Number))
		assert.Equal(t, 3, int(bhe.RollingBlockHistory()[2].Number))

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("replaces blocks on re-org within EVM_FINALITY_DEPTH", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		config := newConfigWithEIP1559DynamicFeesEnabled(t)
		bhe := newBlockHistoryEstimator(t, ethClient, config)

		var blockDelay uint16 = 0
		var historySize uint16 = 3
		var batchSize uint32 = 2
		config.On("BlockHistoryEstimatorBlockDelay").Return(blockDelay)
		config.On("BlockHistoryEstimatorBlockHistorySize").Return(historySize)
		config.On("BlockHistoryEstimatorBatchSize").Return(batchSize)

		b0 := gas.Block{
			Number:       0,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(9001),
		}
		b1 := gas.Block{
			Number:       1,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(9002),
		}
		b2 := gas.Block{
			Number:       2,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(1, 2),
		}
		b3 := gas.Block{
			Number:       3,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(1, 2),
		}
		blocks := []gas.Block{b0, b1, b2, b3}

		gas.SetRollingBlockHistory(bhe, blocks)

		// RE-ORG, head2 and head3 have different hash than saved b2 and b3
		head2 := evmtypes.NewHead(big.NewInt(2), utils.NewHash(), b1.Hash, uint64(time.Now().Unix()), utils.NewBig(&cltest.FixtureChainID))
		head3 := evmtypes.NewHead(big.NewInt(3), utils.NewHash(), head2.Hash, uint64(time.Now().Unix()), utils.NewBig(&cltest.FixtureChainID))
		head3.Parent = &head2

		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				b[0].Method == "eth_getBlockByNumber" && b[0].Args[0] == gas.Int64ToHex(3) && b[0].Args[1] == true && reflect.TypeOf(b[0].Result) == reflect.TypeOf(&gas.Block{}) &&
				b[1].Method == "eth_getBlockByNumber" && b[1].Args[0] == gas.Int64ToHex(2) && b[1].Args[1] == true && reflect.TypeOf(b[1].Result) == reflect.TypeOf(&gas.Block{})
		})).Once().Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			b2New := b2
			b2New.Hash = head2.Hash
			elems[1].Result = &b2New
			b3New := b3
			b3New.Hash = head3.Hash
			elems[0].Result = &b3New
		})

		err := bhe.FetchBlocks(context.Background(), &head3)
		require.NoError(t, err)

		ethClient.AssertExpectations(t)

		require.Len(t, bhe.RollingBlockHistory(), 3)
		assert.Equal(t, 1, int(bhe.RollingBlockHistory()[0].Number))
		assert.Equal(t, 2, int(bhe.RollingBlockHistory()[1].Number))
		assert.Equal(t, 3, int(bhe.RollingBlockHistory()[2].Number))
		assert.Equal(t, b1.Hash.Hex(), bhe.RollingBlockHistory()[0].Hash.Hex())
		assert.Equal(t, head2.Hash.Hex(), bhe.RollingBlockHistory()[1].Hash.Hex())
		assert.Equal(t, head3.Hash.Hex(), bhe.RollingBlockHistory()[2].Hash.Hex())

		config.AssertExpectations(t)
	})

	t.Run("uses locally cached blocks if they are in the chain", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		config := newConfigWithEIP1559DynamicFeesEnabled(t)
		bhe := newBlockHistoryEstimator(t, ethClient, config)

		var blockDelay uint16 = 0
		var historySize uint16 = 3
		var batchSize uint32 = 2
		config.On("BlockHistoryEstimatorBlockDelay").Return(blockDelay)
		config.On("BlockHistoryEstimatorBlockHistorySize").Return(historySize)
		config.On("BlockHistoryEstimatorBatchSize").Return(batchSize)

		b0 := gas.Block{
			Number:       0,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(9001),
		}
		b1 := gas.Block{
			Number:       1,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(9002),
		}
		b2 := gas.Block{
			Number:       2,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(1, 2),
		}
		b3 := gas.Block{
			Number:       3,
			Hash:         utils.NewHash(),
			Transactions: cltest.LegacyTransactionsFromGasPrices(1, 2),
		}
		blocks := []gas.Block{b0, b1, b2, b3}

		gas.SetRollingBlockHistory(bhe, blocks)

		// head2 and head3 have identical hash to saved blocks
		head2 := evmtypes.NewHead(big.NewInt(2), b2.Hash, b1.Hash, uint64(time.Now().Unix()), utils.NewBig(&cltest.FixtureChainID))
		head3 := evmtypes.NewHead(big.NewInt(3), b3.Hash, head2.Hash, uint64(time.Now().Unix()), utils.NewBig(&cltest.FixtureChainID))
		head3.Parent = &head2

		err := bhe.FetchBlocks(context.Background(), &head3)
		require.NoError(t, err)

		ethClient.AssertExpectations(t)

		require.Len(t, bhe.RollingBlockHistory(), 3)
		assert.Equal(t, 1, int(bhe.RollingBlockHistory()[0].Number))
		assert.Equal(t, 2, int(bhe.RollingBlockHistory()[1].Number))
		assert.Equal(t, 3, int(bhe.RollingBlockHistory()[2].Number))
		assert.Equal(t, b1.Hash.Hex(), bhe.RollingBlockHistory()[0].Hash.Hex())
		assert.Equal(t, head2.Hash.Hex(), bhe.RollingBlockHistory()[1].Hash.Hex())
		assert.Equal(t, head3.Hash.Hex(), bhe.RollingBlockHistory()[2].Hash.Hex())

		config.AssertExpectations(t)
	})
}

func TestBlockHistoryEstimator_FetchBlocksAndRecalculate_NoEIP1559(t *testing.T) {
	t.Parallel()

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	config := newConfigWithEIP1559DynamicFeesDisabled(t)

	config.On("BlockHistoryEstimatorBlockDelay").Return(uint16(0))
	config.On("BlockHistoryEstimatorTransactionPercentile").Return(uint16(35))
	config.On("BlockHistoryEstimatorBlockHistorySize").Return(uint16(3))
	config.On("EvmMaxGasPriceWei").Return(big.NewInt(1000))
	config.On("EvmMinGasPriceWei").Return(big.NewInt(0))
	config.On("BlockHistoryEstimatorBatchSize").Return(uint32(0))

	bhe := newBlockHistoryEstimator(t, ethClient, config)

	b1 := gas.Block{
		Number:       1,
		Hash:         utils.NewHash(),
		Transactions: cltest.LegacyTransactionsFromGasPrices(1),
	}
	b2 := gas.Block{
		Number:       2,
		Hash:         utils.NewHash(),
		Transactions: cltest.LegacyTransactionsFromGasPrices(2),
	}
	b3 := gas.Block{
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

	bhe.FetchBlocksAndRecalculate(context.Background(), cltest.Head(3))

	price := gas.GetGasPrice(bhe)
	require.Equal(t, big.NewInt(100), price)

	assert.Len(t, bhe.RollingBlockHistory(), 3)

	config.AssertExpectations(t)
	ethClient.AssertExpectations(t)
}

func TestBlockHistoryEstimator_Recalculate_NoEIP1559(t *testing.T) {
	t.Parallel()

	maxGasPrice := big.NewInt(100)
	minGasPrice := big.NewInt(10)

	t.Run("does not crash or set gas price to zero if there are no transactions", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		config := newConfigWithEIP1559DynamicFeesDisabled(t)

		config.On("BlockHistoryEstimatorTransactionPercentile").Return(uint16(35))

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		blocks := []gas.Block{}
		gas.SetRollingBlockHistory(bhe, blocks)
		bhe.Recalculate(cltest.Head(1))

		blocks = []gas.Block{gas.Block{}}
		gas.SetRollingBlockHistory(bhe, blocks)
		bhe.Recalculate(cltest.Head(1))

		blocks = []gas.Block{gas.Block{Transactions: []gas.Transaction{}}}
		gas.SetRollingBlockHistory(bhe, blocks)
		bhe.Recalculate(cltest.Head(1))

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("sets gas price to ETH_MAX_GAS_PRICE_WEI if the calculation would otherwise exceed it", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		config := newConfigWithEIP1559DynamicFeesDisabled(t)

		config.On("EvmMaxGasPriceWei").Return(maxGasPrice)
		config.On("EvmMinGasPriceWei").Return(minGasPrice)
		config.On("BlockHistoryEstimatorTransactionPercentile").Return(uint16(35))

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		blocks := []gas.Block{
			gas.Block{
				Number:       0,
				Hash:         utils.NewHash(),
				Transactions: cltest.LegacyTransactionsFromGasPrices(9001),
			},
			gas.Block{
				Number:       1,
				Hash:         utils.NewHash(),
				Transactions: cltest.LegacyTransactionsFromGasPrices(9002),
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(1))

		price := gas.GetGasPrice(bhe)
		require.Equal(t, maxGasPrice, price)

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("sets gas price to ETH_MIN_GAS_PRICE_WEI if the calculation would otherwise fall below it", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		config := newConfigWithEIP1559DynamicFeesDisabled(t)

		config.On("EvmMaxGasPriceWei").Return(maxGasPrice)
		config.On("EvmMinGasPriceWei").Return(minGasPrice)
		config.On("BlockHistoryEstimatorTransactionPercentile").Return(uint16(35))

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		blocks := []gas.Block{
			gas.Block{
				Number:       0,
				Hash:         utils.NewHash(),
				Transactions: cltest.LegacyTransactionsFromGasPrices(5),
			},
			gas.Block{
				Number:       1,
				Hash:         utils.NewHash(),
				Transactions: cltest.LegacyTransactionsFromGasPrices(7),
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(1))

		price := gas.GetGasPrice(bhe)
		require.Equal(t, minGasPrice, price)

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("ignores any transaction with a zero gas limit", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		config := newConfigWithEIP1559DynamicFeesDisabled(t)

		config.On("EvmMaxGasPriceWei").Return(maxGasPrice)
		config.On("EvmMinGasPriceWei").Return(minGasPrice)
		config.On("BlockHistoryEstimatorTransactionPercentile").Return(uint16(100))

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		b1Hash := utils.NewHash()
		b2Hash := utils.NewHash()

		blocks := []gas.Block{
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
				Transactions: []gas.Transaction{gas.Transaction{GasPrice: big.NewInt(70), GasLimit: 42}},
			},
			{
				Number:       2,
				Hash:         utils.NewHash(),
				ParentHash:   b2Hash,
				Transactions: []gas.Transaction{gas.Transaction{GasPrice: big.NewInt(90), GasLimit: 0}},
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(2))

		price := gas.GetGasPrice(bhe)
		require.Equal(t, big.NewInt(70), price)

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("takes into account zero priced transctions if chain is not xDai", func(t *testing.T) {
		// Because everyone loves free gas!
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		config := newConfigWithEIP1559DynamicFeesDisabled(t)

		config.On("EvmMaxGasPriceWei").Return(maxGasPrice)
		config.On("EvmMinGasPriceWei").Return(big.NewInt(0))
		config.On("BlockHistoryEstimatorTransactionPercentile").Return(uint16(50))

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		b1Hash := utils.NewHash()

		blocks := []gas.Block{
			gas.Block{
				Number:       0,
				Hash:         b1Hash,
				ParentHash:   common.Hash{},
				Transactions: cltest.LegacyTransactionsFromGasPrices(0, 0, 0, 0, 100),
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(0))

		price := gas.GetGasPrice(bhe)
		require.Equal(t, big.NewInt(0), price)

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("ignores zero priced transactions on xDai", func(t *testing.T) {
		chainID := big.NewInt(100)

		ethClient := cltest.NewEthClientMock(t)
		config := newConfigWithEIP1559DynamicFeesDisabled(t)

		config.On("EvmMaxGasPriceWei").Return(maxGasPrice)
		config.On("EvmMinGasPriceWei").Return(big.NewInt(100))
		config.On("BlockHistoryEstimatorTransactionPercentile").Return(uint16(50))

		ibhe := newBlockHistoryEstimatorWithChainID(t, ethClient, config, *chainID)
		bhe := gas.BlockHistoryEstimatorFromInterface(ibhe)

		b1Hash := utils.NewHash()

		blocks := []gas.Block{
			gas.Block{
				Number:       0,
				Hash:         b1Hash,
				ParentHash:   common.Hash{},
				Transactions: cltest.LegacyTransactionsFromGasPrices(0, 0, 0, 0, 100),
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(0))

		price := gas.GetGasPrice(bhe)
		require.Equal(t, big.NewInt(100), price)

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("handles unreasonably large gas prices (larger than a 64 bit int can hold)", func(t *testing.T) {
		// Seems unlikely we will ever experience gas prices > 9 Petawei on mainnet (praying to the eth Gods 🙏)
		// But other chains could easily use a different base of account
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		config := newConfigWithEIP1559DynamicFeesDisabled(t)

		reasonablyHugeGasPrice := big.NewInt(0).Mul(big.NewInt(math.MaxInt64), big.NewInt(1000))

		config.On("EvmMaxGasPriceWei").Return(reasonablyHugeGasPrice)
		config.On("EvmMinGasPriceWei").Return(big.NewInt(10))
		config.On("BlockHistoryEstimatorTransactionPercentile").Return(uint16(50))

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		unreasonablyHugeGasPrice := big.NewInt(0).Mul(big.NewInt(math.MaxInt64), big.NewInt(1000000))

		b1Hash := utils.NewHash()

		blocks := []gas.Block{
			gas.Block{
				Number:     0,
				Hash:       b1Hash,
				ParentHash: common.Hash{},
				Transactions: []gas.Transaction{
					gas.Transaction{GasPrice: big.NewInt(50), GasLimit: 42},
					gas.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					gas.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					gas.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					gas.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					gas.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					gas.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					gas.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					gas.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
				},
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(0))

		price := gas.GetGasPrice(bhe)
		require.Equal(t, reasonablyHugeGasPrice, price)

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("doesn't panic if gas price is nil (although I'm still unsure how this can happen)", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		config := newConfigWithEIP1559DynamicFeesDisabled(t)

		config.On("EvmMaxGasPriceWei").Return(maxGasPrice)
		config.On("EvmMinGasPriceWei").Return(big.NewInt(100))
		config.On("BlockHistoryEstimatorTransactionPercentile").Return(uint16(50))

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		b1Hash := utils.NewHash()

		blocks := []gas.Block{
			gas.Block{
				Number:     0,
				Hash:       b1Hash,
				ParentHash: common.Hash{},
				Transactions: []gas.Transaction{
					{GasPrice: nil, GasLimit: 42, Hash: utils.NewHash()},
					{GasPrice: big.NewInt(100), GasLimit: 42, Hash: utils.NewHash()},
				},
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(0))

		price := gas.GetGasPrice(bhe)
		require.Equal(t, big.NewInt(100), price)

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})
}

func newBlockWithBaseFee() gas.Block {
	return gas.Block{BaseFeePerGas: assets.GWei(5)}
}

func TestBlockHistoryEstimator_Recalculate_EIP1559(t *testing.T) {
	t.Parallel()

	maxGasPrice := big.NewInt(100)

	t.Run("does not crash or set gas price to zero if there are no transactions", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)

		config := newConfigWithEIP1559DynamicFeesEnabled(t)

		config.On("BlockHistoryEstimatorTransactionPercentile").Return(uint16(35))

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		blocks := []gas.Block{}
		gas.SetRollingBlockHistory(bhe, blocks)
		bhe.Recalculate(cltest.Head(1))

		blocks = []gas.Block{gas.Block{}} // No base fee (doesn't crash)
		gas.SetRollingBlockHistory(bhe, blocks)
		bhe.Recalculate(cltest.Head(1))

		blocks = []gas.Block{newBlockWithBaseFee()}
		gas.SetRollingBlockHistory(bhe, blocks)
		bhe.Recalculate(cltest.Head(1))

		empty := newBlockWithBaseFee()
		empty.Transactions = []gas.Transaction{}
		blocks = []gas.Block{empty}
		gas.SetRollingBlockHistory(bhe, blocks)
		bhe.Recalculate(cltest.Head(1))

		withOnlyLegacyTransactions := newBlockWithBaseFee()
		withOnlyLegacyTransactions.Transactions = cltest.LegacyTransactionsFromGasPrices(9001)
		blocks = []gas.Block{withOnlyLegacyTransactions}
		gas.SetRollingBlockHistory(bhe, blocks)
		bhe.Recalculate(cltest.Head(1))

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("can set tip higher than ETH_MAX_GAS_PRICE_WEI (we rely on fee cap to limit it)", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		config := newConfigWithEIP1559DynamicFeesEnabled(t)

		config.On("EvmMaxGasPriceWei").Return(maxGasPrice)
		config.On("EvmMinGasPriceWei").Return(big.NewInt(0))
		config.On("EvmGasTipCapMinimum").Return(big.NewInt(0))
		config.On("BlockHistoryEstimatorTransactionPercentile").Return(uint16(35))

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		blocks := []gas.Block{
			gas.Block{
				BaseFeePerGas: big.NewInt(1),
				Number:        0,
				Hash:          utils.NewHash(),
				Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(9001),
			},
			gas.Block{
				BaseFeePerGas: big.NewInt(1),
				Number:        1,
				Hash:          utils.NewHash(),
				Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(9002),
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(1))

		tipCap := gas.GetTipCap(bhe)
		require.Greater(t, tipCap.Int64(), maxGasPrice.Int64())

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("sets tip cap to ETH_MIN_GAS_PRICE_WEI if the calculation would otherwise fall below it", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		config := newConfigWithEIP1559DynamicFeesEnabled(t)

		config.On("EvmMaxGasPriceWei").Return(maxGasPrice)
		config.On("EvmMinGasPriceWei").Return(big.NewInt(0))
		config.On("EvmGasTipCapMinimum").Return(big.NewInt(10))
		config.On("BlockHistoryEstimatorTransactionPercentile").Return(uint16(35))

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		blocks := []gas.Block{
			gas.Block{
				BaseFeePerGas: big.NewInt(1),
				Number:        0,
				Hash:          utils.NewHash(),
				Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(5),
			},
			gas.Block{
				BaseFeePerGas: big.NewInt(1),
				Number:        1,
				Hash:          utils.NewHash(),
				Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(7),
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(1))

		price := gas.GetTipCap(bhe)
		require.Equal(t, big.NewInt(10), price)

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("ignores any transaction with a zero gas limit", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		config := newConfigWithEIP1559DynamicFeesEnabled(t)

		config.On("EvmMaxGasPriceWei").Return(maxGasPrice)
		config.On("EvmMinGasPriceWei").Return(big.NewInt(0))
		config.On("EvmGasTipCapMinimum").Return(big.NewInt(10))
		config.On("BlockHistoryEstimatorTransactionPercentile").Return(uint16(95))

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		b1Hash := utils.NewHash()
		b2Hash := utils.NewHash()

		blocks := []gas.Block{
			{
				Number:       0,
				Hash:         b1Hash,
				ParentHash:   common.Hash{},
				Transactions: cltest.LegacyTransactionsFromGasPrices(50),
			},
			{
				BaseFeePerGas: big.NewInt(10),
				Number:        1,
				Hash:          b2Hash,
				ParentHash:    b1Hash,
				Transactions:  []gas.Transaction{gas.Transaction{Type: 0x2, MaxFeePerGas: big.NewInt(1000), MaxPriorityFeePerGas: big.NewInt(60), GasLimit: 42}},
			},
			{
				Number:       2,
				Hash:         utils.NewHash(),
				ParentHash:   b2Hash,
				Transactions: []gas.Transaction{gas.Transaction{Type: 0x2, MaxFeePerGas: big.NewInt(1000), MaxPriorityFeePerGas: big.NewInt(80), GasLimit: 0}},
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(2))

		price := gas.GetTipCap(bhe)
		require.Equal(t, big.NewInt(60), price)

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("respects minimum gas tip cap", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		config := newConfigWithEIP1559DynamicFeesEnabled(t)

		config.On("EvmMaxGasPriceWei").Return(maxGasPrice)
		config.On("EvmMinGasPriceWei").Return(big.NewInt(0))
		config.On("EvmGasTipCapMinimum").Return(big.NewInt(1))
		config.On("BlockHistoryEstimatorTransactionPercentile").Return(uint16(35))

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		b1Hash := utils.NewHash()

		blocks := []gas.Block{
			gas.Block{
				BaseFeePerGas: big.NewInt(10),
				Number:        0,
				Hash:          b1Hash,
				ParentHash:    common.Hash{},
				Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(0, 0, 0, 0, 100),
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(0))

		price := gas.GetTipCap(bhe)
		assert.Equal(t, big.NewInt(1), price)

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})

	t.Run("allows to set zero tip cap if minimum allows it", func(t *testing.T) {
		// Because everyone loves *cheap* gas!
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		config := newConfigWithEIP1559DynamicFeesEnabled(t)

		config.On("EvmMaxGasPriceWei").Return(maxGasPrice)
		config.On("EvmMinGasPriceWei").Return(big.NewInt(0))
		config.On("EvmGasTipCapMinimum").Return(big.NewInt(0))
		config.On("BlockHistoryEstimatorTransactionPercentile").Return(uint16(35))

		bhe := newBlockHistoryEstimator(t, ethClient, config)

		b1Hash := utils.NewHash()

		blocks := []gas.Block{
			gas.Block{
				BaseFeePerGas: big.NewInt(10),
				Number:        0,
				Hash:          b1Hash,
				ParentHash:    common.Hash{},
				Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(0, 0, 0, 0, 100),
			},
		}

		gas.SetRollingBlockHistory(bhe, blocks)

		bhe.Recalculate(cltest.Head(0))

		price := gas.GetTipCap(bhe)
		require.Equal(t, big.NewInt(0), price)

		ethClient.AssertExpectations(t)
		config.AssertExpectations(t)
	})
}

func TestBlockHistoryEstimator_EffectiveTipCap(t *testing.T) {
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	config := newConfigWithEIP1559DynamicFeesEnabled(t)

	bhe := newBlockHistoryEstimator(t, ethClient, config)

	block := gas.Block{
		Number:     0,
		Hash:       utils.NewHash(),
		ParentHash: common.Hash{},
	}

	eipblock := block
	eipblock.BaseFeePerGas = big.NewInt(100)

	t.Run("returns nil if block is missing base fee", func(t *testing.T) {
		tx := gas.Transaction{Type: 0x0, GasPrice: big.NewInt(42), GasLimit: 42, Hash: utils.NewHash()}
		res := bhe.EffectiveTipCap(block, tx)
		assert.Nil(t, res)
	})
	t.Run("legacy transaction type infers tip cap from tx.gas_price - block.base_fee_per_gas", func(t *testing.T) {
		tx := gas.Transaction{Type: 0x0, GasPrice: big.NewInt(142), GasLimit: 42, Hash: utils.NewHash()}
		res := bhe.EffectiveTipCap(eipblock, tx)
		assert.Equal(t, "42", res.String())
	})
	t.Run("tx type 2 should calculate gas price", func(t *testing.T) {
		// 0x2 transaction (should use MaxPriorityFeePerGas)
		tx := gas.Transaction{Type: 0x2, MaxPriorityFeePerGas: big.NewInt(200), MaxFeePerGas: big.NewInt(250), GasLimit: 42, Hash: utils.NewHash()}
		res := bhe.EffectiveTipCap(eipblock, tx)
		assert.Equal(t, "200", res.String())
		// 0x2 transaction (should use MaxPriorityFeePerGas, ignoring gas price)
		tx = gas.Transaction{Type: 0x2, GasPrice: big.NewInt(400), MaxPriorityFeePerGas: big.NewInt(200), MaxFeePerGas: big.NewInt(350), GasLimit: 42, Hash: utils.NewHash()}
		res = bhe.EffectiveTipCap(eipblock, tx)
		assert.Equal(t, "200", res.String())
	})
	t.Run("missing field returns nil", func(t *testing.T) {
		tx := gas.Transaction{Type: 0x2, GasPrice: big.NewInt(132), MaxFeePerGas: big.NewInt(200), GasLimit: 42, Hash: utils.NewHash()}
		res := bhe.EffectiveTipCap(eipblock, tx)
		assert.Nil(t, res)
	})
	t.Run("unknown type returns nil", func(t *testing.T) {
		tx := gas.Transaction{Type: 0x3, GasPrice: big.NewInt(55555), MaxPriorityFeePerGas: big.NewInt(200), MaxFeePerGas: big.NewInt(250), GasLimit: 42, Hash: utils.NewHash()}
		res := bhe.EffectiveTipCap(eipblock, tx)
		assert.Nil(t, res)
	})

	ethClient.AssertExpectations(t)
	config.AssertExpectations(t)
}

func TestBlockHistoryEstimator_EffectiveGasPrice(t *testing.T) {
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	config := newConfigWithEIP1559DynamicFeesDisabled(t)

	bhe := newBlockHistoryEstimator(t, ethClient, config)

	block := gas.Block{
		Number:     0,
		Hash:       utils.NewHash(),
		ParentHash: common.Hash{},
	}

	eipblock := block
	eipblock.BaseFeePerGas = big.NewInt(100)

	t.Run("legacy transaction type should use GasPrice", func(t *testing.T) {
		tx := gas.Transaction{Type: 0x0, GasPrice: big.NewInt(42), GasLimit: 42, Hash: utils.NewHash()}
		res := bhe.EffectiveGasPrice(eipblock, tx)
		assert.Equal(t, "42", res.String())
		tx = gas.Transaction{Type: 0x0, GasLimit: 42, Hash: utils.NewHash()}
		res = bhe.EffectiveGasPrice(eipblock, tx)
		assert.Nil(t, res)
		tx = gas.Transaction{Type: 0x1, GasPrice: big.NewInt(42), GasLimit: 42, Hash: utils.NewHash()}
		res = bhe.EffectiveGasPrice(eipblock, tx)
		assert.Equal(t, "42", res.String())
	})
	t.Run("tx type 2 should calculate gas price", func(t *testing.T) {
		// 0x2 transaction (should calculate to 250)
		tx := gas.Transaction{Type: 0x2, MaxPriorityFeePerGas: big.NewInt(200), MaxFeePerGas: big.NewInt(250), GasLimit: 42, Hash: utils.NewHash()}
		res := bhe.EffectiveGasPrice(eipblock, tx)
		assert.Equal(t, "250", res.String())
		// 0x2 transaction (should calculate to 300)
		tx = gas.Transaction{Type: 0x2, MaxPriorityFeePerGas: big.NewInt(200), MaxFeePerGas: big.NewInt(350), GasLimit: 42, Hash: utils.NewHash()}
		res = bhe.EffectiveGasPrice(eipblock, tx)
		assert.Equal(t, "300", res.String())
		// 0x2 transaction (should calculate to 300, ignoring gas price)
		tx = gas.Transaction{Type: 0x2, MaxPriorityFeePerGas: big.NewInt(200), MaxFeePerGas: big.NewInt(350), GasLimit: 42, Hash: utils.NewHash()}
		res = bhe.EffectiveGasPrice(eipblock, tx)
		assert.Equal(t, "300", res.String())
		// 0x2 transaction (should fall back to gas price since MaxFeePerGas is missing)
		tx = gas.Transaction{Type: 0x2, GasPrice: big.NewInt(32), MaxPriorityFeePerGas: big.NewInt(200), GasLimit: 42, Hash: utils.NewHash()}
		res = bhe.EffectiveGasPrice(eipblock, tx)
		assert.Equal(t, "32", res.String())
	})
	t.Run("tx type 2 has block missing base fee (should never happen but must handle gracefully)", func(t *testing.T) {
		// 0x2 transaction (should calculate to 250)
		tx := gas.Transaction{Type: 0x2, GasPrice: big.NewInt(55555), MaxPriorityFeePerGas: big.NewInt(200), MaxFeePerGas: big.NewInt(250), GasLimit: 42, Hash: utils.NewHash()}
		res := bhe.EffectiveGasPrice(block, tx)
		assert.Equal(t, "55555", res.String())
	})
	t.Run("unknown type returns nil", func(t *testing.T) {
		tx := gas.Transaction{Type: 0x3, GasPrice: big.NewInt(55555), MaxPriorityFeePerGas: big.NewInt(200), MaxFeePerGas: big.NewInt(250), GasLimit: 42, Hash: utils.NewHash()}
		res := bhe.EffectiveGasPrice(block, tx)
		assert.Nil(t, res)
	})

	ethClient.AssertExpectations(t)
	config.AssertExpectations(t)
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

	var block gas.Block
	err := json.Unmarshal([]byte(blockJSON), &block)
	assert.NoError(t, err)

	assert.Equal(t, int64(16023161), block.Number)
	assert.Equal(t, common.HexToHash("0x317cfd032b5d6657995f17fe768f7cc4ea0ada27ad421c4caa685a9071ea955c"), block.Hash)
	assert.Equal(t, common.HexToHash("0xb47ab3b1dc5c2c090dcecdc744a65a279ea6bb8dec11fb3c247df4cc2f584848"), block.ParentHash)

	require.Len(t, block.Transactions, 3)

	assert.Equal(t, int64(1), block.Transactions[0].GasPrice.Int64())
	assert.Equal(t, uint64(1), block.Transactions[0].GasLimit)

	assert.Equal(t, int64(0), block.Transactions[1].GasPrice.Int64())
	assert.Equal(t, uint64(0), block.Transactions[1].GasLimit)

	assert.Equal(t, big.NewInt(4566182400000), block.Transactions[2].GasPrice)
	assert.Equal(t, uint64(2000000), block.Transactions[2].GasLimit)
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

	var block gas.Block
	err := json.Unmarshal([]byte(blockJSON), &block)
	assert.NoError(t, err)

	assert.Equal(t, int64(13000040), block.Number)
	assert.Equal(t, "43362048092", block.BaseFeePerGas.String())
	assert.Equal(t, common.HexToHash("0x11ac873a6cd8b8b7b57ec1efe3984b706362aa5e8f5749a5ec9b1f64bb4615f0"), block.Hash)
	assert.Equal(t, common.HexToHash("0x1ae6168805dfd2e48311181774019c17fb09b24ab75dcad6566d18d38d5c4071"), block.ParentHash)

	require.Len(t, block.Transactions, 4)

	assert.Equal(t, int64(683000000000), block.Transactions[0].GasPrice.Int64())
	assert.Equal(t, 900000, int(block.Transactions[0].GasLimit))
	assert.Nil(t, block.Transactions[0].MaxFeePerGas)
	assert.Nil(t, block.Transactions[0].MaxPriorityFeePerGas)
	assert.Equal(t, gas.TxType(0x0), block.Transactions[0].Type)
	assert.Equal(t, "0x8e58af889f4e831ef9a67df84058bcfb7090cbcb5c6f1046c211dafee6050944", block.Transactions[0].Hash.String())

	assert.Equal(t, big.NewInt(3000000000000), block.Transactions[1].GasPrice)
	assert.Equal(t, "0x0190f436ce165abb741b8513f64d194682677e1db72422f0f533fe6c0248e59a", block.Transactions[1].Hash.String())

	assert.Equal(t, int64(500000000000), block.Transactions[2].GasPrice.Int64())
	assert.Equal(t, 21000, int(block.Transactions[2].GasLimit))
	assert.Equal(t, int64(500000000000), block.Transactions[2].MaxFeePerGas.Int64())
	assert.Equal(t, int64(500000000000), block.Transactions[2].MaxPriorityFeePerGas.Int64())
	assert.Equal(t, gas.TxType(0x2), block.Transactions[2].Type)
	assert.Equal(t, "0x136aa666e6b8109b2b4aca8008ecad8df2047f4e2aced4808248fa8927a13395", block.Transactions[2].Hash.String())

	assert.Nil(t, block.Transactions[3].GasPrice)
	assert.Equal(t, 21000, int(block.Transactions[3].GasLimit))
	assert.Equal(t, "0x13d4ecea98e37359e63e39e350ed0b1456e1acbf985eb8d4a0ef0e89a705c10d", block.Transactions[3].Hash.String())
}

func TestBlockHistoryEstimator_GetDynamicFee(t *testing.T) {
	t.Parallel()

	cfg := newConfigWithEIP1559DynamicFeesEnabled(t)
	maxGasPrice := big.NewInt(1000000)
	cfg.On("BlockHistoryEstimatorEIP1559FeeCapBufferBlocks").Return(uint16(4))
	cfg.On("BlockHistoryEstimatorTransactionPercentile").Return(uint16(35))
	cfg.On("EvmEIP1559DynamicFees").Return(true)
	cfg.On("EvmGasLimitMultiplier").Return(float32(1))
	cfg.On("EvmGasTipCapMinimum").Return(big.NewInt(0))
	cfg.On("EvmMaxGasPriceWei").Return(maxGasPrice)
	cfg.On("EvmMinGasPriceWei").Return(big.NewInt(0))

	bhe := newBlockHistoryEstimator(t, nil, cfg)

	blocks := []gas.Block{
		gas.Block{
			BaseFeePerGas: big.NewInt(88889),
			Number:        0,
			Hash:          utils.NewHash(),
			Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(5000, 6000, 6000),
		},
		gas.Block{
			BaseFeePerGas: big.NewInt(100000),
			Number:        1,
			Hash:          utils.NewHash(),
			Transactions:  cltest.DynamicFeeTransactionsFromTipCaps(10000),
		},
	}
	gas.SetRollingBlockHistory(bhe, blocks)

	bhe.Recalculate(cltest.Head(1))
	gas.SimulateStart(bhe)

	t.Run("if estimator is missing base fee and gas bumping is enabled", func(t *testing.T) {
		cfg.On("EvmGasBumpThreshold").Return(uint64(1)).Once()

		_, _, err := bhe.GetDynamicFee(100000)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "f")
	})

	t.Run("if estimator is missing base fee and gas bumping is disabled", func(t *testing.T) {
		cfg.On("EvmGasBumpThreshold").Return(uint64(0)).Once()

		fee, limit, err := bhe.GetDynamicFee(100000)
		require.NoError(t, err)
		assert.Equal(t, gas.DynamicFee{FeeCap: maxGasPrice, TipCap: big.NewInt(6000)}, fee)
		assert.Equal(t, 100000, int(limit))
	})

	h := cltest.Head(1)
	h.BaseFeePerGas = utils.NewBigI(112500)
	bhe.OnNewLongestChain(context.Background(), h)

	t.Run("if gas bumping is enabled", func(t *testing.T) {
		cfg.On("EvmGasBumpThreshold").Return(uint64(1)).Once()

		fee, limit, err := bhe.GetDynamicFee(100000)
		require.NoError(t, err)

		assert.Equal(t, gas.DynamicFee{FeeCap: big.NewInt(186203), TipCap: big.NewInt(6000)}, fee)
		assert.Equal(t, 100000, int(limit))
	})

	t.Run("if gas bumping is disabled", func(t *testing.T) {
		cfg.On("EvmGasBumpThreshold").Return(uint64(0)).Once()

		fee, limit, err := bhe.GetDynamicFee(100000)
		require.NoError(t, err)

		assert.Equal(t, gas.DynamicFee{FeeCap: maxGasPrice, TipCap: big.NewInt(6000)}, fee)
		assert.Equal(t, 100000, int(limit))
	})

	cfg.AssertExpectations(t)
}

func TestBlockHistoryEstimator_Bumps(t *testing.T) {
	t.Parallel()

	t.Run("BumpLegacyGas calls BumpLegacyGasPriceOnly with proper current gas price", func(t *testing.T) {
		config := newConfigWithEIP1559DynamicFeesDisabled(t)
		bhe := newBlockHistoryEstimator(t, nil, config)

		config.On("EvmGasBumpPercent").Return(uint16(10))
		config.On("EvmGasBumpWei").Return(big.NewInt(150))
		config.On("EvmMaxGasPriceWei").Return(big.NewInt(1000000))
		config.On("EvmGasLimitMultiplier").Return(float32(1.1))

		t.Run("ignores nil current gas price", func(t *testing.T) {
			gasPrice, gasLimit, err := bhe.BumpLegacyGas(big.NewInt(42), 100000)
			require.NoError(t, err)

			expectedGasPrice, expectedGasLimit, err := gas.BumpLegacyGasPriceOnly(config, logger.TestLogger(t), nil, big.NewInt(42), 100000)
			require.NoError(t, err)

			assert.Equal(t, expectedGasLimit, gasLimit)
			assert.Equal(t, expectedGasPrice, gasPrice)
		})

		t.Run("ignores current gas price > max gas price", func(t *testing.T) {
			gasPrice, gasLimit, err := bhe.BumpLegacyGas(big.NewInt(42), 100000)
			require.NoError(t, err)

			massive := big.NewInt(100000000000000)
			gas.SetGasPrice(bhe, massive)

			expectedGasPrice, expectedGasLimit, err := gas.BumpLegacyGasPriceOnly(config, logger.TestLogger(t), massive, big.NewInt(42), 100000)
			require.NoError(t, err)

			assert.Equal(t, expectedGasLimit, gasLimit)
			assert.Equal(t, expectedGasPrice, gasPrice)
		})

		t.Run("ignores current gas price < bumped gas price", func(t *testing.T) {
			gas.SetGasPrice(bhe, big.NewInt(191))

			gasPrice, gasLimit, err := bhe.BumpLegacyGas(big.NewInt(42), 100000)
			require.NoError(t, err)

			assert.Equal(t, 110000, int(gasLimit))
			assert.Equal(t, big.NewInt(192), gasPrice)
		})

		t.Run("uses current gas price > bumped gas price", func(t *testing.T) {
			gas.SetGasPrice(bhe, big.NewInt(193))

			gasPrice, gasLimit, err := bhe.BumpLegacyGas(big.NewInt(42), 100000)
			require.NoError(t, err)

			assert.Equal(t, 110000, int(gasLimit))
			assert.Equal(t, big.NewInt(193), gasPrice)
		})

		config.AssertExpectations(t)
	})

	t.Run("BumpDynamicFee bumps the fee", func(t *testing.T) {
		config := newConfigWithEIP1559DynamicFeesEnabled(t)
		bhe := newBlockHistoryEstimator(t, nil, config)

		config.On("EvmGasBumpPercent").Return(uint16(10))
		config.On("EvmGasBumpWei").Return(big.NewInt(150))
		config.On("EvmMaxGasPriceWei").Return(big.NewInt(1000000))
		config.On("EvmGasLimitMultiplier").Return(float32(1.1))
		config.On("EvmGasTipCapDefault").Return(big.NewInt(52))

		t.Run("when current tip cap is nil", func(t *testing.T) {
			originalFee := gas.DynamicFee{FeeCap: big.NewInt(100), TipCap: big.NewInt(25)}
			fee, gasLimit, err := bhe.BumpDynamicFee(originalFee, 100000)
			require.NoError(t, err)

			assert.Equal(t, 110000, int(gasLimit))
			assert.Equal(t, gas.DynamicFee{FeeCap: big.NewInt(250), TipCap: big.NewInt(202)}, fee)
		})
		t.Run("ignores current tip cap that is smaller than original fee with bump applied", func(t *testing.T) {
			gas.SetTipCap(bhe, big.NewInt(201))

			originalFee := gas.DynamicFee{FeeCap: big.NewInt(100), TipCap: big.NewInt(25)}
			fee, gasLimit, err := bhe.BumpDynamicFee(originalFee, 100000)
			require.NoError(t, err)

			assert.Equal(t, 110000, int(gasLimit))
			assert.Equal(t, gas.DynamicFee{FeeCap: big.NewInt(250), TipCap: big.NewInt(202)}, fee)
		})
		t.Run("uses current tip cap that is larger than original fee with bump applied", func(t *testing.T) {
			gas.SetTipCap(bhe, big.NewInt(203))

			originalFee := gas.DynamicFee{FeeCap: big.NewInt(100), TipCap: big.NewInt(25)}
			fee, gasLimit, err := bhe.BumpDynamicFee(originalFee, 100000)
			require.NoError(t, err)

			assert.Equal(t, 110000, int(gasLimit))
			assert.Equal(t, gas.DynamicFee{FeeCap: big.NewInt(250), TipCap: big.NewInt(203)}, fee)
		})
		t.Run("ignores absurdly large current tip cap", func(t *testing.T) {
			gas.SetTipCap(bhe, big.NewInt(1000000000000000))

			originalFee := gas.DynamicFee{FeeCap: big.NewInt(100), TipCap: big.NewInt(25)}
			fee, gasLimit, err := bhe.BumpDynamicFee(originalFee, 100000)
			require.NoError(t, err)

			assert.Equal(t, 110000, int(gasLimit))
			assert.Equal(t, gas.DynamicFee{FeeCap: big.NewInt(250), TipCap: big.NewInt(202)}, fee)
		})

		config.AssertExpectations(t)
	})
}
