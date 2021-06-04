package headtracker

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --name BlockEthClient --output ./mocks/ --case=underscore
type BlockEthClient interface {
	BlockByNumber(ctx context.Context, number int64) (*Block, error)
	FastBlockByHash(ctx context.Context, hash common.Hash) (*Block, error)
	FetchBlocksByNumbers(ctx context.Context, numbers []int64) (map[int64]Block, error)
	FetchLatestHead(ctx context.Context) (*models.Head, error)
}

type BlockEthClientImpl struct {
	ethClient eth.Client
	logger    *logger.Logger
	batchSize uint32
}

func NewBlockEthClientImpl(ethClient eth.Client, logger *logger.Logger, batchSize uint32) *BlockEthClientImpl {
	return &BlockEthClientImpl{
		ethClient: ethClient,
		logger:    logger,
		batchSize: batchSize,
	}
}

func (bc *BlockEthClientImpl) BlockByNumber(ctx context.Context, number int64) (*Block, error) {
	block, err := bc.ethClient.BlockByNumber(ctx, big.NewInt(number))
	if err != nil {
		return nil, err
	}
	b := FromEthBlock(*block)
	return &b, nil
}

func (bc *BlockEthClientImpl) FastBlockByHash(ctx context.Context, hash common.Hash) (*Block, error) {
	block, err := bc.ethClient.FastBlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	b := FromEthBlock(*block)
	return &b, nil
}

func (bc *BlockEthClientImpl) FetchBlocksByNumbers(ctx context.Context, numbers []int64) (map[int64]Block, error) {
	var reqs []rpc.BatchElem
	for _, number := range numbers {
		req := rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{Int64ToHex(number), true},
			Result: &Block{},
		}
		reqs = append(reqs, req)
	}

	if err := bc.batchFetch(ctx, reqs); err != nil {
		return nil, err
	}

	blocks := make(map[int64]Block)
	for i, req := range reqs {
		result, err := req.Result, req.Error
		if err != nil {
			bc.logger.Warnw("BlockFetcher#fetchBlocksByNumbers error while fetching block", "err", err, "blockNum", numbers[i])
			continue
		}

		b, is := result.(*Block)
		if !is {
			return nil, errors.Errorf("expected result to be a %T, got %T", &Block{}, result)
		}
		if b == nil {
			//TODO: can this happen on "Fetching it too early results in an empty block." ?
			bc.logger.Warnw("BlockFetcher#fetchBlocksByNumbers got nil block", "blockNum", numbers[i], "index", i)
			continue
		}
		if b.Hash == (common.Hash{}) {
			bc.logger.Warnw("BlockFetcher#fetchBlocksByNumbers block was missing hash", "block", b, "blockNum", numbers[i], "erroredBlockNum", b.Number)
			continue
		}

		blocks[b.Number] = *b
	}
	return blocks, nil
}

func (bc *BlockEthClientImpl) batchFetch(ctx context.Context, reqs []rpc.BatchElem) error {
	batchSize := int(bc.batchSize)

	if batchSize == 0 {
		batchSize = len(reqs)
	}
	for i := 0; i < len(reqs); i += batchSize {
		j := i + batchSize
		if j > len(reqs) {
			j = len(reqs)
		}

		var numbers []int64
		for ii := i; ii < j; ii++ {
			numbers = append(numbers, HexToInt64(reqs[ii].Args[0]))
		}

		logger.Debugw(fmt.Sprintf("BlockFetcher: Batch fetching %v blocks with numbers: %v", len(numbers), numbers))

		err := bc.ethClient.BatchCallContext(ctx, reqs[i:j])
		if ctx.Err() != nil {
			break
		} else if err != nil {
			return errors.Wrap(err, "BlockFetcher: Error fetching blocks with BatchCallContext")
		}
	}
	return nil
}

func (bc *BlockEthClientImpl) FetchLatestHead(ctx context.Context) (*models.Head, error) {
	return bc.ethClient.HeaderByNumber(ctx, nil)
}

type FakeBlockEthClient struct {
	blocks []Block
}

func NewFakeBlockEthClient(blocks []Block) *FakeBlockEthClient {
	return &FakeBlockEthClient{
		blocks: blocks,
	}
}

func (bc *FakeBlockEthClient) SetBlocks(blocks []Block) {
	bc.blocks = blocks
}

func (bc *FakeBlockEthClient) BlockByNumber(ctx context.Context, number int64) (*Block, error) {
	for _, block := range bc.blocks {
		if block.Number == number {
			return &block, nil
		}
	}
	return nil, errors.Errorf("not found: %v", number)
}

func (bc *FakeBlockEthClient) FastBlockByHash(ctx context.Context, hash common.Hash) (*Block, error) {
	for _, block := range bc.blocks {
		if block.Hash == hash {
			return &block, nil
		}
	}
	return nil, errors.Errorf("not found: %v", hash)
}

func (bc *FakeBlockEthClient) FetchBlocksByNumbers(ctx context.Context, numbers []int64) (map[int64]Block, error) {
	mapp := make(map[int64]Block)
	for _, block := range bc.blocks {
		for _, number := range numbers {
			if block.Number == number {
				mapp[number] = block
			}
		}
	}
	return mapp, nil
}

func (bc *FakeBlockEthClient) FetchLatestHead(ctx context.Context) (*models.Head, error) {
	if len(bc.blocks) == 0 {
		return nil, errors.Errorf("no head found")
	}
	head := headFromBlock(bc.blocks[len(bc.blocks)-1])
	return &head, nil
}
