package headtracker

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name BlockFetcherInterface --output ./mocks/ --case=underscore
type (
	BlockFetcherInterface interface {
		FetchLatestHead(ctx context.Context) (*models.Head, error)
		BlockRange(ctx context.Context, fromBlock int64, toBlock int64) ([]Block, error)
		SyncLatestHead(ctx context.Context, head models.Head) error
		Chain(ctx context.Context, latestHead models.Head) (models.Head, error)
	}

	BlockFetcher struct {
		ethClient eth.Client
		logger    *logger.Logger
		config    BlockFetcherConfig

		blockEthClient BlockEthClient
		recent         map[common.Hash]*Block
		latestBlockNum int64

		notifications chan struct{}

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

//type BlockDownload struct {
//	StartedAt time.Time
//
//	Number int64
//
//	// may be nil
//	Hash *common.Hash
//
//	// may be nil
//	Head *models.Head
//}

func (bf *BlockFetcher) BlockCache() []*Block {
	return bf.RecentSorted()
}

func NewBlockFetcher(config BlockFetcherConfig, logger *logger.Logger, blockEthClient BlockEthClient) *BlockFetcher {

	return &BlockFetcher{
		logger:         logger,
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
	bf.logger.Debug("Returned from Chain")
	return headWithChain, nil
}

func (bf *BlockFetcher) SyncLatestHead(ctx context.Context, head models.Head) error {
	bf.mut.Lock()
	cacheSize := len(bf.recent)
	bf.mut.Unlock()

	if cacheSize == 0 {
		bf.logger.Debug("SyncLatestHead: cache is empty, will backfill using batch calls")
		bf.Backfill(ctx, head)
	}

	_, err := bf.syncLatestHead(ctx, head)
	bf.logger.Debug("Returned from SyncLatestHead")
	return err
}

func (bf *BlockFetcher) StartDownloadAsync(ctx context.Context, fromBlock int64, toBlock int64) {
	timeout := 2 * time.Minute

	go func() {
		utils.RetryWithBackoff(ctx, func() (retry bool) {
			ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			err := bf.downloadRange(ctxTimeout, fromBlock, toBlock)
			if err != nil {
				bf.logger.Errorw("BlockFetcher#StartDownload error while downloading blocks. Will retry",
					"err", err, "fromBlock", fromBlock, "toBlock", toBlock)
				return true
			}
			return false
		})
		bf.logger.Debug("Returned from StartDownloadAsync")
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
	bf.logger.Debugw("BlockFetcher#BlockRange requested range", "fromBlock", fromBlock, "toBlock", toBlock)

	if fromBlock < 0 || toBlock < fromBlock {
		return make([]*Block, 0), errors.Errorf("Invalid range: %d -> %d", fromBlock, toBlock)
	}

	blocks := make(map[int64]*Block)

	err := bf.downloadRange(ctx, fromBlock, toBlock)
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

	return blockRange, nil
	// check if currently being downloaded
	// trigger download of missing ones
	//for {
	//	select {
	//	case _ = <-bd.notifications:
	//		// case after?
	//	default:
	//
	//	}
	//}
}

func (bf *BlockFetcher) downloadRange(ctx context.Context, fromBlock int64, toBlock int64) error {

	bf.mut.Lock()

	existingBlocks := make(map[int64]*Block)

	for _, b := range bf.recent {
		block := b
		if block.Number >= fromBlock && block.Number <= toBlock {
			existingBlocks[block.Number] = block
		}
	}
	bf.mut.Unlock()

	var blockNumsToFetch []int64
	// schedule fetch of missing blocks
	for i := fromBlock; i <= toBlock; i++ {
		if _, exists := existingBlocks[i]; exists {
			continue
		}
		blockNumsToFetch = append(blockNumsToFetch, i)
	}

	if len(blockNumsToFetch) > 0 {
		blocksFetched, err := bf.blockEthClient.FetchBlocksByNumbers(ctx, blockNumsToFetch)
		if err != nil {
			bf.logger.Errorw("BlockFetcher#BlockRange error while fetching missing blocks", "err", err, "fromBlock", fromBlock, "toBlock", toBlock)
			return err
		}

		if len(blocksFetched) < len(blockNumsToFetch) {
			bf.logger.Warnw("BlockFetcher#BlockRange did not fetch all requested blocks",
				"fromBlock", fromBlock, "toBlock", toBlock, "requestedLen", len(blockNumsToFetch), "blocksFetched", len(blocksFetched))
		}

		bf.mut.Lock()

		for _, blockItem := range blocksFetched {
			block := blockItem

			existingBlocks[block.Number] = &block
			bf.recent[block.Hash] = &block

			if bf.latestBlockNum < block.Number {
				bf.latestBlockNum = block.Number
			}
		}

		bf.cleanRecent()
		bf.mut.Unlock()

		select {
		case bf.notifications <- struct{}{}:
		default:
		}
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

func fromEthBlock(ethBlock types.Block) Block {
	var block Block
	block.Number = ethBlock.Number().Int64()
	block.Hash = ethBlock.Hash()
	block.ParentHash = ethBlock.ParentHash()

	return block
}

func headFromBlock(ethBlock Block) models.Head {
	var head models.Head
	head.Number = ethBlock.Number
	head.Hash = ethBlock.Hash
	head.ParentHash = ethBlock.ParentHash
	return head
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

	bf.logger.Debugw("BlockFetcher: starting sync head",
		"blockNumber", head.Number,
		"fromBlockHeight", from,
		"toBlockHeight", head.Number)
	defer func() {
		if ctx.Err() != nil {
			return
		}
		bf.logger.Debugw("BlockFetcher: finished sync head",
			"fetched", fetched,
			"blockNumber", head.Number,
			"time", time.Since(mark),
			"fromBlockHeight", from,
			"toBlockHeight", head.Number)
	}()

	// first, check if the previous block already exists locally
	var existingPrevBlock = bf.findBlockByHash(head.ParentHash)
	if existingPrevBlock != nil {
		// if yes, just fetch the latest block
		block, err := bf.fetchAndSaveBlock(ctx, head.Hash)
		if err != nil {
			return models.Head{}, errors.Wrap(err, "Failed to fetch latest block")
		}

		return bf.sequentialConstructChain(ctx, block, from)
	} else {
		// we don't have the previous block or there was a re-org
		bf.logger.Debugf("Getting a range of blocks: %v to %v", from, head.Number)
		blocks, err := bf.GetBlockRange(ctx, from, head.Number)

		if len(blocks) == 0 {
			logger.Warnf("No blocks returned by range %v to %v", from, head.Number)
			return head, nil
		}
		sort.Slice(blocks, func(i, j int) bool {
			return blocks[i].Number < blocks[j].Number
		})
		if err != nil {
			return models.Head{}, errors.Wrap(err, "BlockByNumber failed")
		}
		return bf.sequentialConstructChain(ctx, *blocks[len(blocks)-1], from)
	}
}

func (bf *BlockFetcher) fetchAndSaveBlock(ctx context.Context, hash common.Hash) (Block, error) {
	bf.logger.Debugf("Fetching block by hash: %v", hash)
	blockPtr, err := bf.blockEthClient.FastBlockByHash(ctx, hash)
	if ctx.Err() != nil {
		return Block{}, nil
	}
	if err != nil {
		return Block{}, errors.Wrap(err, "FastBlockByHash failed")
	}

	bf.mut.Lock()
	bf.recent[blockPtr.Hash] = blockPtr
	bf.mut.Unlock()
	return *blockPtr, nil
}

func (bf *BlockFetcher) sequentialConstructChain(ctx context.Context, block Block, from int64) (models.Head, error) {
	fetched := 0

	var chainTip = headFromBlock(block)
	var currentHead = &chainTip

	bf.logger.Debugf("Latest block: %v, %v", block.Number, block.Hash)

	for i := block.Number - 1; i >= from; i-- {

		zeroHash := common.Hash{}
		if currentHead.ParentHash == zeroHash {
			bf.logger.Debugf("currentHead.ParentHash is zero - returning")
			break
		}
		var existingBlock = bf.findBlockByHash(currentHead.ParentHash)
		if existingBlock != nil {
			block = *existingBlock
			bf.logger.Debugf("Found block in cache: %v - %v", block.Number, block.Hash)
		} else {
			bf.logger.Debugf("Fetching BlockByNumber: %v, as existing block was not found by %v", i, currentHead.ParentHash)

			//TODO: perhaps implement FastBlockByNumber
			blockPtr, err := bf.blockEthClient.BlockByNumber(ctx, i)
			fetched++
			if ctx.Err() != nil {
				break
			} else if err != nil {
				return models.Head{}, errors.Wrap(err, "BlockByNumber failed")
			}

			bf.mut.Lock()
			bf.recent[block.Hash] = blockPtr
			bf.mut.Unlock()
		}

		head := headFromBlock(block)
		currentHead.Parent = &head
		currentHead = &head
	}
	bf.logger.Debugf("Returning chain of length %v", chainTip.ChainLength())
	return chainTip, nil
}
