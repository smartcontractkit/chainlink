package headtracker

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	promBlockDownloads = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "block_fetcher_block_downloads",
		Help: "Number of block downloads",
	}, []string{"method"})
)

//go:generate mockery --name BlockFetcherInterface --output ./mocks/ --case=underscore
type (
	BlockFetcherInterface interface {
		FetchLatestHead(ctx context.Context) (*models.Head, error)
		BlockRange(ctx context.Context, fromBlock int64, toBlock int64) ([]Block, error)
		BlocksWithoutCache(ctx context.Context, numbers []int64) (map[int64]Block, error)
		SyncLatestHead(ctx context.Context, head models.Head) error
		Chain(ctx context.Context, latestHead models.Head) (models.Head, error)
	}

	BlockFetcher struct {
		logger *logger.Logger
		orm    ORM
		config BlockFetcherConfig

		blockEthClient BlockEthClient
		recent         map[common.Hash]*Block
		latestBlockNum int64

		mut        sync.Mutex
		syncingMut sync.Mutex
	}
)

//go:generate mockery --name BlockFetcherConfig --output ./mocks/ --case=underscore
type BlockFetcherConfig interface {
	BlockFetcherBatchSize() uint32
	EthFinalityDepth() uint
	EthHeadTrackerHistoryDepth() uint
	BlockBackfillDepth() uint64
}

func NewBlockFetcher(orm ORM, config BlockFetcherConfig, logger *logger.Logger, blockEthClient BlockEthClient) *BlockFetcher {
	return &BlockFetcher{
		logger:         logger,
		orm:            orm,
		config:         config,
		recent:         make(map[common.Hash]*Block),
		blockEthClient: blockEthClient,
	}
}

func (bf *BlockFetcher) Backfill(ctx context.Context, latestHead models.Head) {
	from := latestHead.Number - int64(bf.config.EthHeadTrackerHistoryDepth()-1)
	if from < 0 {
		from = 0
	}
	bf.StartDownloadAsync(ctx, from, latestHead.Number)
}

// FetchLatestHead - Fetches the latest head from the blockchain, regardless of local cache
func (bf *BlockFetcher) FetchLatestHead(ctx context.Context) (*models.Head, error) {
	return bf.blockEthClient.FetchLatestHead(ctx)
}

// BlockRange - Returns a range of blocks, either from local memory or fetched
func (bf *BlockFetcher) BlockRange(ctx context.Context, fromBlock int64, toBlock int64) ([]Block, error) {
	bf.logger.Debugw("BlockFetcher#BlockRange requested range", "fromBlock", fromBlock, "toBlock", toBlock)

	blocks, err := bf.GetBlockRange(ctx, fromBlock, toBlock)
	if err != nil {
		return nil, errors.Wrapf(err, "BlockFetcher#GetBlockRange error for %v -> %v",
			fromBlock, toBlock)
	}
	blocksSlice := make([]Block, 0)
	for _, block := range blocks {
		blocksSlice = append(blocksSlice, *block)
	}

	return blocksSlice, nil
}

func (bf *BlockFetcher) BlocksWithoutCache(ctx context.Context, numbers []int64) (map[int64]Block, error) {
	bf.logger.Debugw("BlockFetcher#BlocksWithoutCache", "blockNumbers", numbers)
	blocks, err := bf.blockEthClient.FetchBlocksByNumbers(ctx, numbers)
	if blocks != nil && len(blocks) > 0 {
		promBlockDownloads.WithLabelValues("FetchBlocksByNumbers").Add(float64(len(blocks)))
	}
	return blocks, err
}

func (bf *BlockFetcher) Chain(ctx context.Context, latestHead models.Head) (models.Head, error) {
	bf.logger.Debugf("BlockFetcher#Chain for head: %v %v", latestHead.Number, latestHead.Hash)

	from := latestHead.Number - int64(bf.config.EthFinalityDepth()-1)
	if from < 0 {
		from = 0
	}

	// typically all the heads should be constructed into a chain from memory, unless there was a very recent reorg
	// and new blocks were not downloaded yet
	headWithChain, err := bf.syncLatestHead(ctx, latestHead)
	if err != nil {
		return models.Head{}, errors.Wrapf(err, "BlockFetcher#Chain error for syncLatestHead: %v", latestHead.Number)
	}
	return headWithChain, nil
}

func (bf *BlockFetcher) SyncLatestHead(ctx context.Context, head models.Head) error {
	_, err := bf.syncLatestHead(ctx, head)
	return err
}

func (bf *BlockFetcher) StartDownloadAsync(ctx context.Context, fromBlock int64, toBlock int64) {
	timeout := 2 * time.Minute

	go func() {
		utils.RetryWithBackoff(ctx, func() (retry bool) {
			ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			err := bf.downloadRange(ctxTimeout, fromBlock, toBlock, true)
			if err != nil {
				bf.logger.Errorw("BlockFetcher: error while downloading blocks. Will retry",
					"err", err, "fromBlock", fromBlock, "toBlock", toBlock)
				return true
			}
			return false
		})
		bf.logger.Debug("BlockFetcher: Finished async download of blocks")
	}()
}

func (bf *BlockFetcher) RecentSorted() []*Block {

	toReturn := make([]*Block, 0)

	for _, b := range bf.recent {
		block := b
		toReturn = append(toReturn, block)
	}

	sort.Slice(toReturn, func(i, j int) bool {
		return toReturn[i].Number < toReturn[j].Number
	})
	return toReturn
}

func (bf *BlockFetcher) GetBlockRange(ctx context.Context, fromBlock int64, toBlock int64) ([]*Block, error) {
	return bf.getBlockRange(ctx, fromBlock, toBlock, true)
}
func (bf *BlockFetcher) getBlockRange(ctx context.Context, fromBlock int64, toBlock int64, useCache bool) ([]*Block, error) {

	if fromBlock < 0 || toBlock < fromBlock {
		return make([]*Block, 0), errors.Errorf("Invalid range: %d -> %d", fromBlock, toBlock)
	}

	blocks := make(map[int64]*Block)

	err := bf.downloadRange(ctx, fromBlock, toBlock, useCache)
	if err != nil {
		return make([]*Block, 0), errors.Wrapf(err, "BlockFetcher#GetBlockRange error while downloading blocks %v -> %v",
			fromBlock, toBlock)
	}
	bf.mut.Lock()

	for _, b := range bf.recent {
		block := b
		if block.Number >= fromBlock && block.Number <= toBlock {
			blocks[block.Number] = block
		}
	}
	bf.mut.Unlock()

	blockRange := make([]*Block, 0)
	for _, b := range blocks {
		block := b
		blockRange = append(blockRange, block)
	}

	sort.Slice(blockRange, func(i, j int) bool {
		return blockRange[i].Number < blockRange[j].Number
	})

	return blockRange, nil
}

func (bf *BlockFetcher) downloadRange(ctx context.Context, fromBlock int64, toBlock int64, useCache bool) error {

	existingBlocks := make(map[int64]*Block)

	if useCache {
		bf.mut.Lock()
		for _, b := range bf.recent {
			block := b
			if block.Number >= fromBlock && block.Number <= toBlock {
				existingBlocks[block.Number] = block
			}
		}
		bf.mut.Unlock()
	}

	var blockNumsToFetch []int64
	for i := fromBlock; i <= toBlock; i++ {
		if useCache {
			if _, exists := existingBlocks[i]; exists {
				bf.logger.Debugf("BlockFetcher: %v already exists", i)
				continue
			}
		}
		blockNumsToFetch = append(blockNumsToFetch, i)
	}

	if len(blockNumsToFetch) > 0 {
		blocksFetched, err := bf.blockEthClient.FetchBlocksByNumbers(ctx, blockNumsToFetch)
		if err != nil {
			bf.logger.Errorw("BlockFetcher: error while fetching missing blocks", "err", err, "fromBlock", fromBlock, "toBlock", toBlock)
			return err
		}
		promBlockDownloads.WithLabelValues("FetchBlocksByNumbers").Add(float64(len(blocksFetched)))

		if len(blocksFetched) < len(blockNumsToFetch) {
			bf.logger.Warnw("BlockFetcher: did not fetch all requested blocks",
				"fromBlock", fromBlock, "toBlock", toBlock, "requestedLen", len(blockNumsToFetch), "blocksFetched", len(blocksFetched))
		}

		bf.mut.Lock()

		for _, blockItem := range blocksFetched {
			block := blockItem

			bf.recent[block.Hash] = &block

			if bf.latestBlockNum < block.Number {
				bf.latestBlockNum = block.Number
			}
		}

		bf.cleanRecent()
		bf.mut.Unlock()
	}
	return nil
}

func (bf *BlockFetcher) cleanRecent() {
	var blockNumsToDelete []common.Hash
	for _, b := range bf.recent {
		block := b
		if block.Number < bf.latestBlockNum-int64(bf.config.EthHeadTrackerHistoryDepth()) {
			blockNumsToDelete = append(blockNumsToDelete, block.Hash)
		}
	}
	for _, toDelete := range blockNumsToDelete {
		delete(bf.recent, toDelete)
	}
}

func (bf *BlockFetcher) findBlockByHash(hash common.Hash) *Block {
	bf.mut.Lock()
	defer bf.mut.Unlock()
	block, ok := bf.recent[hash]
	if ok {
		return block
	}
	return nil
}

func (bf *BlockFetcher) syncLatestHead(ctx context.Context, head models.Head) (models.Head, error) {
	bf.syncingMut.Lock()
	defer bf.syncingMut.Unlock()

	from := head.Number - int64(bf.config.EthHeadTrackerHistoryDepth()-1)
	if from < 0 {
		from = 0
	}

	mark := time.Now()
	fetched := 0

	bf.logger.Debugw("BlockFetcher: Starting sync head",
		"blockNumber", head.Number,
		"blockHash", head.Hash,
		"fromBlockHeight", from,
		"toBlockHeight", head.Number)
	defer func() {
		if ctx.Err() != nil {
			return
		}
		bf.logger.Debugw("BlockFetcher: Finished sync head",
			"fetched", fetched,
			"blockNumber", head.Number,
			"time", time.Since(mark),
			"fromBlockHeight", from,
			"toBlockHeight", head.Number)
	}()

	// first, check if the previous block already exists locally
	logger.Debugf("Block by parent hash %v", head.ParentHash)
	var existingPrevBlock = bf.findBlockByHash(head.ParentHash)
	if existingPrevBlock != nil {
		logger.Debugf("Previous block already exists locally: %v, %v", existingPrevBlock.Number, existingPrevBlock.Hash)
		// if yes, just fetch the latest block
		block, err := bf.fetchAndSaveBlock(ctx, head.Hash)
		if err != nil {
			return models.Head{}, errors.Wrap(err, "Failed to fetch latest block")
		}

		return bf.sequentialConstructChain(ctx, block, from)
	}

	// We don't have the previous block at all, or there was a re-org
	// For efficiency, fetch ranges of previous blocks until an ancestor already known
	block, err := bf.fetchInBatchesUntilKnownAncestor(ctx, head, from)
	if err != nil {
		return models.Head{}, errors.Wrap(err, "Failed fetching blocks until known ancestor")
	}

	if block == nil {
		logger.Warnw("BlockFetcher: No latest block returned from fetchInBatchesUntilKnownAncestor", "latestHead", head.Number, "from", from)
		return head, nil
	}
	return bf.sequentialConstructChain(ctx, *block, from)
}

// Fetch ranges of previous blocks, until reaching an ancestor already existing in cache
func (bf *BlockFetcher) fetchInBatchesUntilKnownAncestor(ctx context.Context, head models.Head, untilFrom int64) (*Block, error) {
	var latestBlock *Block
	latestToFetch := head.Number
	for {
		currentFrom := latestToFetch - int64(bf.config.BlockFetcherBatchSize()) + 1
		currentTo := latestToFetch

		logger.Debugw("BlockFetcher: Getting Range",
			"currentFrom", currentFrom, "currentTo", currentTo, "untilFrom", untilFrom)

		if currentFrom < untilFrom {
			currentFrom = untilFrom
		}
		if currentFrom < 0 {
			currentFrom = 0
		}
		if currentFrom > currentTo {
			break
		}

		blocks, err := bf.getBlockRange(ctx, currentFrom, currentTo, false)
		if err != nil {
			return nil, errors.Wrap(err, "BlockByNumber failed")
		}

		if len(blocks) == 0 {
			logger.Warnf("BlockFetcher: No blocks returned by range %v to %v. ", currentFrom, currentTo)
			break
		}

		if latestBlock == nil {
			latestBlock = blocks[len(blocks)-1]
		}

		earliestFetched := blocks[0]
		existingParent := bf.findBlockByHash(earliestFetched.ParentHash)
		if existingParent != nil {
			logger.Debugw("BlockFetcher: Found already known ancestor block",
				"existingParentNumber", existingParent.Number, "existingParentHash", existingParent.Hash)
			break
		}

		if currentFrom == 0 || currentFrom <= untilFrom {
			// done
			break
		}
		latestToFetch = currentFrom - 1
	}
	return latestBlock, nil
}

func (bf *BlockFetcher) fetchAndSaveBlock(ctx context.Context, hash common.Hash) (Block, error) {
	bf.logger.Debugf("BlockFetcher: Fetching block by hash: %v", hash)
	blockPtr, err := bf.blockEthClient.FastBlockByHash(ctx, hash)
	promBlockDownloads.WithLabelValues("FastBlockByHash").Inc()

	if ctx.Err() != nil {
		return Block{}, nil
	}
	if err != nil {
		return Block{}, errors.Wrap(err, "FastBlockByHash failed")
	}

	err = bf.storeBlock(ctx, blockPtr)
	if err != nil {
		return Block{}, errors.Wrap(err, "storeBlock failed")
	}
	return *blockPtr, nil
}

func (bf *BlockFetcher) sequentialConstructChain(ctx context.Context, block Block, from int64) (models.Head, error) {

	var chainTip = HeadFromBlock(block)
	var currentHead = &chainTip

	bf.logger.Debugf("BlockFetcher: Constructing the chain until the latest block: %v, %v", block.Number, block.Hash)

	currentBlock := &block
	for i := block.Number - 1; i >= from; i-- {

		zeroHash := common.Hash{}
		if currentHead.ParentHash == zeroHash {
			bf.logger.Debugf("BlockFetcher: currentHead.ParentHash is zero - returning")
			break
		}
		var existingBlock = bf.findBlockByHash(currentHead.ParentHash)
		if existingBlock != nil {
			currentBlock = existingBlock
		} else {
			bf.logger.Debugf("BlockFetcher: Fetching block by number: %v, as existing block was not found by %v", i, currentHead.ParentHash)

			blockPtr, err := bf.blockEthClient.BlockByNumber(ctx, i)
			promBlockDownloads.WithLabelValues("BlockByNumber").Inc()
			if ctx.Err() != nil {
				break
			} else if err != nil {
				return models.Head{}, errors.Wrap(err, "BlockByNumber failed")
			}

			currentBlock = blockPtr

			err = bf.storeBlock(ctx, blockPtr)
			if err != nil {
				return models.Head{}, errors.Wrap(err, "storeBlock failed")
			}
		}

		head := HeadFromBlock(*currentBlock)
		currentHead.Parent = &head
		currentHead = &head
	}
	bf.logger.Debugf("BlockFetcher: Returning chain of length %v", chainTip.ChainLength())
	return chainTip, nil
}

func (bf *BlockFetcher) storeBlock(ctx context.Context, block *Block) error {
	bf.mut.Lock()
	bf.recent[block.Hash] = block
	bf.mut.Unlock()

	return bf.orm.SaveBlock(ctx, block)
}
