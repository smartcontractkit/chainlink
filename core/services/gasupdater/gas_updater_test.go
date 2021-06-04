package gasupdater_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	logger2 "github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/gasupdater"
	gumocks "github.com/smartcontractkit/chainlink/core/services/gasupdater/mocks"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	htmocks "github.com/smartcontractkit/chainlink/core/services/headtracker/mocks"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGasUpdater_Start(t *testing.T) {
	t.Parallel()

	config := new(gumocks.Config)

	var batchSize uint32 = 0
	var blockDelay uint16 = 0
	var historySize uint16 = 2
	var ethFinalityDepth uint = 42
	var percentile uint16 = 35
	minGasPrice := big.NewInt(1)

	config.On("GasUpdaterBatchSize").Return(batchSize)
	config.On("GasUpdaterBlockDelay").Return(blockDelay)
	config.On("GasUpdaterBlockHistorySize").Return(historySize)
	config.On("EthFinalityDepth").Return(ethFinalityDepth)
	config.On("GasUpdaterTransactionPercentile").Return(percentile)
	config.On("EthMinGasPriceWei").Return(minGasPrice)
	config.On("ChainID").Return(big.NewInt(0))

	t.Run("loads initial state", func(t *testing.T) {
		ethClient := new(mocks.Client)

		blocks := []headtracker.Block{*cltest.NewBlock(41, common.Hash{}), *cltest.NewBlock(42, common.Hash{})}

		guIface := createGasUpdaterWithBlocks(config, blocks)
		gu := gasupdater.GasUpdaterToStruct(guIface)

		err := guIface.Start()
		require.NoError(t, err)

		assert.Len(t, gu.RollingBlockHistory(), 2)
		assert.Equal(t, int(gu.RollingBlockHistory()[0].Number), 41)
		assert.Equal(t, int(gu.RollingBlockHistory()[1].Number), 42)

		ethClient.AssertExpectations(t)
	})

	t.Run("boots even if initial batch call returns nothing", func(t *testing.T) {

		blockClient := new(htmocks.BlockEthClient)

		h := &models.Head{Hash: cltest.NewHash(), Number: 42}
		blockMap := make(map[int64]headtracker.Block)
		blockClient.On("FetchLatestHead", mock.Anything).Return(h, nil)
		blockClient.On("FetchBlocksByNumbers", mock.Anything, mock.Anything).Return(blockMap, nil)

		gu := createGasUpdaterWithBlockClient(config, blockClient)
		err := gu.Start()
		require.NoError(t, err)

		blockClient.AssertExpectations(t)
	})

	t.Run("starts anyway if fetching latest head fails", func(t *testing.T) {
		blockClient := new(htmocks.BlockEthClient)
		blockClient.On("FetchLatestHead", mock.Anything).Return(nil, errors.New("something exploded"))

		gu := createGasUpdaterWithBlockClient(config, blockClient)
		err := gu.Start()
		require.NoError(t, err)

		blockClient.AssertExpectations(t)
	})
}

func TestGasUpdater_FetchBlocks(t *testing.T) {
	t.Parallel()

	t.Run("with history size of 0, errors", func(t *testing.T) {
		config := new(gumocks.Config)
		guIface := createGasUpdater(config)
		gu := gasupdater.GasUpdaterToStruct(guIface)

		var blockDelay uint16 = 3
		var historySize uint16 = 0
		config.On("ChainID").Return(big.NewInt(0))
		config.On("GasUpdaterBlockDelay").Return(blockDelay)
		config.On("GasUpdaterBlockHistorySize").Return(historySize)

		head := cltest.Head(42)
		err := gu.FetchBlocks(context.Background(), *head)
		require.Error(t, err)
		require.EqualError(t, err, "GasUpdater: history size must be > 0, got: 0")
	})

	t.Run("with current block height less than block delay does nothing", func(t *testing.T) {
		config := new(gumocks.Config)
		var blockDelay uint16 = 3
		var historySize uint16 = 1
		config.On("GasUpdaterBlockDelay").Return(blockDelay)
		config.On("GasUpdaterBlockHistorySize").Return(historySize)

		ethClient := new(mocks.Client)

		guIface := createGasUpdater(config)

		gu := gasupdater.GasUpdaterToStruct(guIface)

		for i := -1; i < 3; i++ {
			head := cltest.Head(i)
			err := gu.FetchBlocks(context.Background(), *head)
			require.Error(t, err)
			require.EqualError(t, err, fmt.Sprintf("GasUpdater: cannot fetch, current block height %v is lower than GAS_UPDATER_BLOCK_DELAY=3", i))
		}

		ethClient.AssertExpectations(t)
	})

	t.Run("with error retrieving blocks returns error", func(t *testing.T) {
		config := new(gumocks.Config)
		var blockDelay uint16 = 3
		var historySize uint16 = 3
		var batchSize uint32 = 0
		config.On("GasUpdaterBlockDelay").Return(blockDelay)
		config.On("GasUpdaterBlockHistorySize").Return(historySize)
		config.On("GasUpdaterBatchSize").Return(batchSize)
		blockClient := new(htmocks.BlockEthClient)

		guIface := createGasUpdaterWithBlockClient(config, blockClient)
		gu := gasupdater.GasUpdaterToStruct(guIface)

		blockClient.On("FetchBlocksByNumbers", mock.Anything, mock.Anything).Return(nil, errors.New("something exploded"))

		err := gu.FetchBlocks(context.Background(), *cltest.Head(42))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "something exploded")

		blockClient.AssertExpectations(t)
	})

	t.Run("batch fetches heads and transactions and sets them on the gas updater instance", func(t *testing.T) {
		config := new(gumocks.Config)
		var blockDelay uint16 = 1
		var historySize uint16 = 3
		var batchSize uint32 = 2
		config.On("GasUpdaterBlockDelay").Return(blockDelay)
		config.On("GasUpdaterBlockHistorySize").Return(historySize)
		// Test batching
		config.On("GasUpdaterBatchSize").Return(batchSize)

		blockClient := new(htmocks.BlockEthClient)

		blockMap := make(map[int64]headtracker.Block)
		blockMap[41] = cltest.NewBlockWithTransactions(41, cltest.TransactionsFromGasPrices(1, 2))
		// block 42 is missing / errored
		blockMap[43] = cltest.NewBlockWithTransactions(43, cltest.TransactionsFromGasPrices())
		blockClient.On("FetchBlocksByNumbers", mock.Anything, mock.Anything).Return(blockMap, nil)

		guIface := createGasUpdaterWithBlockClient(config, blockClient)
		gu := gasupdater.GasUpdaterToStruct(guIface)

		err := gu.FetchBlocks(context.Background(), *cltest.Head(44))
		require.NoError(t, err)

		assert.Len(t, gu.RollingBlockHistory(), 2)
		assert.Equal(t, 41, int(gu.RollingBlockHistory()[0].Number))
		assert.Equal(t, 43, int(gu.RollingBlockHistory()[1].Number))
		assert.Len(t, gu.RollingBlockHistory()[0].Transactions, 2)
		assert.Len(t, gu.RollingBlockHistory()[1].Transactions, 0)

		// On new fetch, rolls over the history and drops the old heads
		blockMap2 := make(map[int64]headtracker.Block)
		blockMap[42] = cltest.NewBlockWithTransactions(42, cltest.TransactionsFromGasPrices(3))
		blockMap[44] = cltest.NewBlockWithTransactions(44, cltest.TransactionsFromGasPrices(4))
		blockClient.On("FetchBlocksByNumbers", mock.Anything, mock.Anything).Return(blockMap2, nil)

		err = gu.FetchBlocks(context.Background(), *cltest.Head(45))
		require.NoError(t, err)

		assert.Len(t, gu.RollingBlockHistory(), 3)
		assert.Equal(t, 42, int(gu.RollingBlockHistory()[0].Number))
		assert.Equal(t, 43, int(gu.RollingBlockHistory()[1].Number))
		assert.Equal(t, 44, int(gu.RollingBlockHistory()[2].Number))
		assert.Len(t, gu.RollingBlockHistory()[0].Transactions, 1)
		assert.Len(t, gu.RollingBlockHistory()[1].Transactions, 0)
		assert.Len(t, gu.RollingBlockHistory()[2].Transactions, 1)
	})
}

func TestGasUpdater_FetchBlocksAndRecalculate(t *testing.T) {
	t.Parallel()

	config := new(gumocks.Config)

	config.On("GasUpdaterBlockDelay").Return(uint16(0))
	config.On("GasUpdaterTransactionPercentile").Return(uint16(35))
	config.On("GasUpdaterBlockHistorySize").Return(uint16(3))
	config.On("EthMaxGasPriceWei").Return(big.NewInt(1000))
	config.On("EthMinGasPriceWei").Return(big.NewInt(0))
	config.On("ChainID").Return(big.NewInt(0))

	blocks := []headtracker.Block{
		cltest.NewBlockWithTransactions(1, cltest.TransactionsFromGasPrices(1)),
		cltest.NewBlockWithTransactions(2, cltest.TransactionsFromGasPrices(2)),
		cltest.NewBlockWithTransactions(3, cltest.TransactionsFromGasPrices(200, 300, 100, 100, 100, 100)),
	}

	guIface := createGasUpdaterWithBlocks(config, blocks)
	gu := gasupdater.GasUpdaterToStruct(guIface)

	config.On("SetEthGasPriceDefault", big.NewInt(100)).Return(nil)

	gu.FetchBlocksAndRecalculate(context.Background(), *cltest.Head(3))

	assert.Len(t, gu.RollingBlockHistory(), 3)

	config.AssertExpectations(t)
}

func TestGasUpdater_Recalculate(t *testing.T) {
	t.Parallel()

	maxGasPrice := big.NewInt(100)
	minGasPrice := big.NewInt(10)

	t.Run("does not crash or set gas price to zero if there are no transactions", func(t *testing.T) {
		config := new(gumocks.Config)

		config.On("GasUpdaterTransactionPercentile").Return(uint16(35))
		config.On("EthMinGasPriceWei").Return(big.NewInt(1))
		config.On("ChainID").Return(big.NewInt(0))

		guIface := createGasUpdater(config)
		gu := gasupdater.GasUpdaterToStruct(guIface)

		blocks := []headtracker.Block{}
		gasupdater.SetRollingBlockHistory(gu, blocks)
		gu.Recalculate(*cltest.Head(1))

		blocks = []headtracker.Block{headtracker.Block{}}
		gasupdater.SetRollingBlockHistory(gu, blocks)
		gu.Recalculate(*cltest.Head(1))

		blocks = []headtracker.Block{headtracker.Block{Transactions: []headtracker.Transaction{}}}
		gasupdater.SetRollingBlockHistory(gu, blocks)
		gu.Recalculate(*cltest.Head(1))

		config.AssertExpectations(t)
	})

	t.Run("sets gas price to ETH_MAX_GAS_PRICE_WEI if the calculation would otherwise exceed it", func(t *testing.T) {
		config := new(gumocks.Config)

		config.On("EthMaxGasPriceWei").Return(maxGasPrice)
		config.On("EthMinGasPriceWei").Return(minGasPrice)
		config.On("GasUpdaterTransactionPercentile").Return(uint16(35))
		config.On("ChainID").Return(big.NewInt(0))

		guIface := createGasUpdater(config)
		gu := gasupdater.GasUpdaterToStruct(guIface)

		blocks := []headtracker.Block{
			headtracker.Block{
				Number:       0,
				Hash:         cltest.NewHash(),
				Transactions: cltest.TransactionsFromGasPrices(9001),
			},
			headtracker.Block{
				Number:       1,
				Hash:         cltest.NewHash(),
				Transactions: cltest.TransactionsFromGasPrices(9002),
			},
		}

		config.On("SetEthGasPriceDefault", maxGasPrice).Return(nil)
		gasupdater.SetRollingBlockHistory(gu, blocks)

		gu.Recalculate(*cltest.Head(1))

		config.AssertExpectations(t)
	})

	t.Run("sets gas price to ETH_MIN_GAS_PRICE_WEI if the calculation would otherwise fall below it", func(t *testing.T) {
		config := new(gumocks.Config)

		config.On("EthMaxGasPriceWei").Return(maxGasPrice)
		config.On("EthMinGasPriceWei").Return(minGasPrice)
		config.On("GasUpdaterTransactionPercentile").Return(uint16(35))
		config.On("ChainID").Return(big.NewInt(0))

		guIface := createGasUpdater(config)
		gu := gasupdater.GasUpdaterToStruct(guIface)

		blocks := []headtracker.Block{
			headtracker.Block{
				Number:       0,
				Hash:         cltest.NewHash(),
				Transactions: cltest.TransactionsFromGasPrices(5),
			},
			headtracker.Block{
				Number:       1,
				Hash:         cltest.NewHash(),
				Transactions: cltest.TransactionsFromGasPrices(7),
			},
		}

		config.On("SetEthGasPriceDefault", minGasPrice).Return(nil)
		gasupdater.SetRollingBlockHistory(gu, blocks)

		gu.Recalculate(*cltest.Head(1))

		config.AssertExpectations(t)
	})

	t.Run("ignores any transaction with a zero gas limit", func(t *testing.T) {
		config := new(gumocks.Config)

		config.On("EthMaxGasPriceWei").Return(maxGasPrice)
		config.On("EthMinGasPriceWei").Return(minGasPrice)
		config.On("GasUpdaterTransactionPercentile").Return(uint16(100))
		config.On("ChainID").Return(big.NewInt(0))

		guIface := createGasUpdater(config)

		gu := gasupdater.GasUpdaterToStruct(guIface)

		b1Hash := cltest.NewHash()
		b2Hash := cltest.NewHash()

		blocks := []headtracker.Block{
			headtracker.Block{
				Number:       0,
				Hash:         b1Hash,
				ParentHash:   common.Hash{},
				Transactions: cltest.TransactionsFromGasPrices(50),
			},
			headtracker.Block{
				Number:       1,
				Hash:         b2Hash,
				ParentHash:   b1Hash,
				Transactions: []headtracker.Transaction{headtracker.Transaction{GasPrice: big.NewInt(70), GasLimit: 42}},
			},
			headtracker.Block{
				Number:       2,
				Hash:         cltest.NewHash(),
				ParentHash:   b2Hash,
				Transactions: []headtracker.Transaction{headtracker.Transaction{GasPrice: big.NewInt(90), GasLimit: 0}},
			},
		}

		config.On("SetEthGasPriceDefault", big.NewInt(70)).Return(nil)
		gasupdater.SetRollingBlockHistory(gu, blocks)

		gu.Recalculate(*cltest.Head(2))

		config.AssertExpectations(t)
	})

	t.Run("takes into account zero priced transctions if chain is not xDai", func(t *testing.T) {
		// Because everyone loves free gas!
		config := new(gumocks.Config)

		config.On("EthMaxGasPriceWei").Return(maxGasPrice)
		config.On("EthMinGasPriceWei").Return(big.NewInt(0))
		config.On("GasUpdaterTransactionPercentile").Return(uint16(50))
		config.On("ChainID").Return(big.NewInt(0))

		guIface := createGasUpdater(config)
		gu := gasupdater.GasUpdaterToStruct(guIface)

		b1Hash := cltest.NewHash()

		blocks := []headtracker.Block{
			headtracker.Block{
				Number:       0,
				Hash:         b1Hash,
				ParentHash:   common.Hash{},
				Transactions: cltest.TransactionsFromGasPrices(0, 0, 0, 0, 100),
			},
		}

		config.On("SetEthGasPriceDefault", big.NewInt(0)).Return(nil)
		gasupdater.SetRollingBlockHistory(gu, blocks)

		gu.Recalculate(*cltest.Head(0))

		config.AssertExpectations(t)
	})

	t.Run("ignores zero priced transactions on xDai", func(t *testing.T) {
		config := new(gumocks.Config)

		config.On("EthMaxGasPriceWei").Return(maxGasPrice)
		config.On("EthMinGasPriceWei").Return(big.NewInt(100))
		config.On("GasUpdaterTransactionPercentile").Return(uint16(50))
		config.On("ChainID").Return(big.NewInt(100))

		guIface := createGasUpdater(config)
		gu := gasupdater.GasUpdaterToStruct(guIface)

		b1Hash := cltest.NewHash()

		blocks := []headtracker.Block{
			headtracker.Block{
				Number:       0,
				Hash:         b1Hash,
				ParentHash:   common.Hash{},
				Transactions: cltest.TransactionsFromGasPrices(0, 0, 0, 0, 100),
			},
		}

		config.On("SetEthGasPriceDefault", big.NewInt(100)).Return(nil)
		gasupdater.SetRollingBlockHistory(gu, blocks)

		gu.Recalculate(*cltest.Head(0))

		config.AssertExpectations(t)
	})

	t.Run("handles unreasonably large gas prices (larger than a 64 bit int can hold)", func(t *testing.T) {
		// Seems unlikely we will ever experience gas prices > 9 Petawei on mainnet (praying to the eth Gods 🙏)
		// But other chains could easily use a different base of account
		config := new(gumocks.Config)

		reasonablyHugeGasPrice := big.NewInt(0).Mul(big.NewInt(math.MaxInt64), big.NewInt(1000))

		config.On("EthMaxGasPriceWei").Return(reasonablyHugeGasPrice)
		config.On("EthMinGasPriceWei").Return(big.NewInt(10))
		config.On("GasUpdaterTransactionPercentile").Return(uint16(50))
		config.On("ChainID").Return(big.NewInt(0))

		guIface := createGasUpdater(config)
		gu := gasupdater.GasUpdaterToStruct(guIface)

		unreasonablyHugeGasPrice := big.NewInt(0).Mul(big.NewInt(math.MaxInt64), big.NewInt(1000000))

		b1Hash := cltest.NewHash()

		blocks := []headtracker.Block{
			headtracker.Block{
				Number:     0,
				Hash:       b1Hash,
				ParentHash: common.Hash{},
				Transactions: []headtracker.Transaction{
					headtracker.Transaction{GasPrice: big.NewInt(50), GasLimit: 42},
					headtracker.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					headtracker.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					headtracker.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					headtracker.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					headtracker.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					headtracker.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					headtracker.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
					headtracker.Transaction{GasPrice: unreasonablyHugeGasPrice, GasLimit: 42},
				},
			},
		}

		config.On("SetEthGasPriceDefault", reasonablyHugeGasPrice).Return(nil)
		gasupdater.SetRollingBlockHistory(gu, blocks)

		gu.Recalculate(*cltest.Head(0))

		config.AssertExpectations(t)
	})
}

func createGasUpdater(config *gumocks.Config) gasupdater.GasUpdater {
	blockClient := headtracker.NewFakeBlockEthClient([]headtracker.Block{})
	bfConfig := new(htmocks.BlockFetcherConfig)
	bfConfig.On("BlockBackfillDepth").Return(uint64(50))
	bfConfig.On("BlockFetcherBatchSize").Return(uint32(4))
	bfConfig.On("EthFinalityDepth").Return(uint(50))
	bfConfig.On("EthHeadTrackerHistoryDepth").Return(uint(10))
	blockFetcher := headtracker.NewBlockFetcher(bfConfig, logger2.Default, blockClient)
	return gasupdater.NewGasUpdater(blockFetcher, config)
}

func createGasUpdaterWithBlocks(config *gumocks.Config, blocks []headtracker.Block) gasupdater.GasUpdater {
	blockClient := headtracker.NewFakeBlockEthClient(blocks)
	bfConfig := new(htmocks.BlockFetcherConfig)
	bfConfig.On("BlockBackfillDepth").Return(uint64(50))
	bfConfig.On("BlockFetcherBatchSize").Return(uint32(4))
	bfConfig.On("EthFinalityDepth").Return(uint(50))
	bfConfig.On("EthHeadTrackerHistoryDepth").Return(uint(10))

	blockFetcher := headtracker.NewBlockFetcher(bfConfig, logger2.Default, blockClient)
	return gasupdater.NewGasUpdater(blockFetcher, config)
}

func createGasUpdaterWithBlockClient(config *gumocks.Config, blockClient headtracker.BlockEthClient) gasupdater.GasUpdater {
	bfConfig := new(htmocks.BlockFetcherConfig)
	bfConfig.On("BlockBackfillDepth").Return(uint64(50))
	bfConfig.On("BlockFetcherBatchSize").Return(uint32(4))
	bfConfig.On("EthFinalityDepth").Return(uint(50))
	bfConfig.On("EthHeadTrackerHistoryDepth").Return(uint(10))

	blockFetcher := headtracker.NewBlockFetcher(bfConfig, logger2.Default, blockClient)
	return gasupdater.NewGasUpdater(blockFetcher, config)
}

func TestGasUpdater_Block(t *testing.T) {
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

	var block headtracker.Block
	err := json.Unmarshal([]byte(blockJSON), &block)
	assert.NoError(t, err)

	assert.Equal(t, int64(16023161), block.Number)
	assert.Equal(t, common.HexToHash("0x317cfd032b5d6657995f17fe768f7cc4ea0ada27ad421c4caa685a9071ea955c"), block.Hash)
	assert.Equal(t, common.HexToHash("0xb47ab3b1dc5c2c090dcecdc744a65a279ea6bb8dec11fb3c247df4cc2f584848"), block.ParentHash)

	require.Len(t, block.Transactions, 2)

	assert.Equal(t, int64(0), block.Transactions[0].GasPrice.Int64())
	assert.Equal(t, uint64(0), block.Transactions[0].GasLimit)
	assert.Equal(t, big.NewInt(4566182400000), block.Transactions[1].GasPrice)
	assert.Equal(t, uint64(2000000), block.Transactions[1].GasLimit)
}
