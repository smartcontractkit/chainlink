package blockheaderfeeder

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/blockhashstore"
)

type Client interface {
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
}

var (
	zeroHash [32]byte
)

// BatchBHS defines an interface for interacting with a BatchBlockhashStore contract.
type BatchBHS interface {
	// GetBlockhashes returns blockhashes for given blockNumbers
	GetBlockhashes(ctx context.Context, blockNumbers []*big.Int) ([][32]byte, error)

	// StoreVerifyHeader stores blockhashes on-chain by using block headers
	StoreVerifyHeader(ctx context.Context, blockNumbers []*big.Int, blockHeaders [][]byte) error
}

// NewBlockHeaderFeeder creates a new BlockHeaderFeeder instance.
func NewBlockHeaderFeeder(
	logger logger.Logger,
	coordinator blockhashstore.Coordinator,
	bhs blockhashstore.BHS,
	batchBHS BatchBHS,
	lookbackBlocks int,
	latestBlock func(ctx context.Context) (uint64, error),
	client Client,
) *BlockHeaderFeeder {
	return &BlockHeaderFeeder{
		lggr:           logger,
		coordinator:    coordinator,
		bhs:            bhs,
		lookbackBlocks: lookbackBlocks,
		latestBlock:    latestBlock,
		stored:         make(map[uint64]struct{}),
		lastRunBlock:   0,
		ec:             client,
	}
}

// BlockHeaderFeeder checks recent VRF coordinator events and stores any blockhashes for blocks within
// (latest - 256) and lookbackBlocks that have unfulfilled requests.
type BlockHeaderFeeder struct {
	lggr                      logger.Logger
	coordinator               blockhashstore.Coordinator
	bhs                       blockhashstore.BHS
	batchBHS                  BatchBHS
	lookbackBlocks            int
	latestBlock               func(ctx context.Context) (uint64, error)
	stored                    map[uint64]struct{}
	lastRunBlock              uint64
	getBlockhashesBatchSize   uint16
	storeBlockhashesBatchSize uint16
	ec                        Client
}

// Run the feeder.
func (f *BlockHeaderFeeder) Run(ctx context.Context) error {
	latestBlock, err := f.latestBlock(ctx)
	if err != nil {
		f.lggr.Errorw("Failed to fetch current block number", "error", err)
		return errors.Wrap(err, "fetching block number")
	}

	var (
		fromBlock = int(latestBlock) - f.lookbackBlocks
		// EVM BLOCKHASH opcode allows (current block - 256) to be fetched on chain
		// BlockHeaderFeeder is responsible for block hashes older than 256 blocks
		toBlock = int(latestBlock) - 256
	)
	if fromBlock < 0 {
		fromBlock = 0
	}
	if toBlock < 0 {
		// Nothing to process, no blocks are in range.
		return nil
	}

	lggr := f.lggr.With("latestBlock", latestBlock)

	blockToRequests, err := blockhashstore.GetUnfulfilledBlocksAndRequests(ctx, lggr, f.coordinator, uint64(fromBlock), uint64(toBlock))
	if err != nil {
		return err
	}

	minBlockNumber := f.findLowestBlockNumberWithoutBlockhash(ctx, lggr, blockToRequests)
	if minBlockNumber == nil {
		lggr.Debug("no blocks to store")
		return nil
	}

	earliestStoredBlockNumber, err := f.findEarliestBlockNumberWithBlockhash(ctx, lggr, minBlockNumber.Uint64()+1, uint64(toBlock))
	if err != nil {
		return errors.Wrap(err, "finding earliest blocknumber with blockhash")
	}

	if earliestStoredBlockNumber == nil {
		// store earliest blockhash and return
		// on next iteration, earliestStoredBlockNumber will be found and
		// will make progress in storing blockhashes using blockheader
		f.bhs.StoreEarliest(ctx)
		lggr.Info("Stored earliest block number")
		return nil
	}

	blocks, err := blockhashstore.DecreasingBlockRange(earliestStoredBlockNumber.Sub(earliestStoredBlockNumber, big.NewInt(1)), minBlockNumber)
	if err != nil {
		return err
	}

	for i := 0; i < len(blocks); i += int(f.storeBlockhashesBatchSize) {
		j := i + int(f.storeBlockhashesBatchSize)
		if j > len(blocks) {
			j = len(blocks)
		}
		blockRange := blocks[i:j]
		blockHeaders, err := f.getRlpHeadersBatch(ctx, blockRange)
		if err != nil {
			return errors.Wrap(err, "fetching block headers")
		}

		lggr.Debugw("storing block headers", "blockRange", blockRange)
		err = f.batchBHS.StoreVerifyHeader(ctx, blockRange, blockHeaders)
		if err != nil {
			return errors.Wrap(err, "storing block headers")
		}
	}

	if f.lastRunBlock != 0 {
		// Prune stored, anything older than fromBlock can be discarded
		for block := f.lastRunBlock - uint64(f.lookbackBlocks); block < uint64(fromBlock); block++ {
			if _, ok := f.stored[block]; ok {
				delete(f.stored, block)
				f.lggr.Debugw("Pruned block from stored cache",
					"block", block, "latestBlock", latestBlock)
			}
		}
	}
	f.lastRunBlock = latestBlock
	return nil
}

func (f *BlockHeaderFeeder) findLowestBlockNumberWithoutBlockhash(ctx context.Context, lggr logger.Logger, blockToRequests map[uint64]map[string]struct{}) *big.Int {
	var min *big.Int
	for block, unfulfilledReqs := range blockToRequests {
		if len(unfulfilledReqs) == 0 {
			continue
		}
		if _, ok := f.stored[block]; ok {
			// Already stored
			continue
		}
		stored, err := f.bhs.IsStored(ctx, block)
		if err != nil {
			lggr.Warnw("Failed to check if block is already stored",
				"error", err,
				"block", block)
		} else if stored {
			lggr.Infow("Blockhash already stored",
				"block", block, "unfulfilledReqIDs", blockhashstore.LimitReqIDs(unfulfilledReqs))
			f.stored[block] = struct{}{}
			continue
		}
		blockNumber := big.NewInt(0).SetUint64(block)
		if min == nil || min.Cmp(blockNumber) >= 0 {
			min = blockNumber
		}
	}
	return min
}

func (f *BlockHeaderFeeder) findEarliestBlockNumberWithBlockhash(ctx context.Context, lggr logger.Logger, startBlock, toBlock uint64) (*big.Int, error) {
	from := startBlock
	for i := startBlock; i < toBlock; i += uint64(f.getBlockhashesBatchSize) {
		j := from + uint64(f.getBlockhashesBatchSize)
		if j > toBlock {
			j = toBlock
		}

		lggr.Debug(fmt.Sprintf("Looking for earliest block number with blockhash %v thru %v", i, j))

		blockNumber := i
		var blocks []*big.Int
		for blockNumber < j {
			blocks = append(blocks, big.NewInt(0).SetUint64(blockNumber))
			blockNumber++
		}

		blockhashes, err := f.batchBHS.GetBlockhashes(ctx, blocks)
		if err != nil {
			return nil, errors.Wrap(err, "fetching blockhashes")
		}

		for idx, bh := range blockhashes {
			if !bytes.Equal(bh[:], zeroHash[:]) {
				earliestBlockNumber := i + uint64(idx)
				lggr.Infow("found earliest block number with blockhash", "earliestBlockNumber", earliestBlockNumber, "blockhash", bh)
				return big.NewInt(0).SetUint64(earliestBlockNumber), nil
			}
		}
	}
	return nil, nil
}

func (f *BlockHeaderFeeder) getRlpHeadersBatch(ctx context.Context, blockRange []*big.Int) ([][]byte, error) {
	var reqs []rpc.BatchElem
	for _, num := range blockRange {
		req := rpc.BatchElem{
			Method: "eth_getHeaderByNumber",
			// Get child block since it's the one that has the parent hash in its header.
			Args:   []interface{}{num.Add(num, big.NewInt(1))},
			Result: &types.Header{},
		}
		reqs = append(reqs, req)
	}
	err := f.ec.BatchCallContext(ctx, reqs)
	if err != nil {
		return nil, err
	}

	var headers [][]byte
	for _, req := range reqs {
		header, err := rlp.EncodeToBytes(req.Result)
		if err != nil {
			return nil, err
		}
		headers = append(headers, header)
	}

	return headers, nil
}
