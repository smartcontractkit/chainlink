// The block header feeder package enables automated lookback and blockhash filling beyond the
// EVM 256 block lookback window to catch missed block hashes.
package blockheaderfeeder

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
)

var (
	zeroHash [32]byte
)

type BlockHeaderProvider interface {
	RlpHeadersBatch(ctx context.Context, blockRange []*big.Int) ([][]byte, error)
}

// BatchBHS defines an interface for interacting with a BatchBlockhashStore contract.
type BatchBHS interface {
	// GetBlockhashes returns blockhashes for given blockNumbers
	GetBlockhashes(ctx context.Context, blockNumbers []*big.Int) ([][32]byte, error)

	// StoreVerifyHeader stores blockhashes on-chain by using block headers
	StoreVerifyHeader(ctx context.Context, blockNumbers []*big.Int, blockHeaders [][]byte, fromAddress common.Address) error
}

// NewBlockHeaderFeeder creates a new BlockHeaderFeeder instance.
func NewBlockHeaderFeeder(
	logger logger.Logger,
	coordinator blockhashstore.Coordinator,
	bhs blockhashstore.BHS,
	batchBHS BatchBHS,
	blockHeaderProvider BlockHeaderProvider,
	waitBlocks int,
	lookbackBlocks int,
	latestBlock func(ctx context.Context) (uint64, error),
	gethks keystore.Eth,
	getBlockhashesBatchSize uint16,
	storeBlockhashesBatchSize uint16,
	fromAddresses []ethkey.EIP55Address,
	chainID *big.Int,
) *BlockHeaderFeeder {
	return &BlockHeaderFeeder{
		lggr:                      logger,
		coordinator:               coordinator,
		bhs:                       bhs,
		batchBHS:                  batchBHS,
		waitBlocks:                waitBlocks,
		lookbackBlocks:            lookbackBlocks,
		latestBlock:               latestBlock,
		stored:                    make(map[uint64]struct{}),
		lastRunBlock:              0,
		getBlockhashesBatchSize:   getBlockhashesBatchSize,
		storeBlockhashesBatchSize: storeBlockhashesBatchSize,
		blockHeaderProvider:       blockHeaderProvider,
		gethks:                    gethks,
		fromAddresses:             fromAddresses,
		chainID:                   chainID,
	}
}

// BlockHeaderFeeder checks recent VRF coordinator events and stores any blockhashes for blocks within
// waitBlocks and lookbackBlocks that have unfulfilled requests.
type BlockHeaderFeeder struct {
	lggr                      logger.Logger
	coordinator               blockhashstore.Coordinator
	bhs                       blockhashstore.BHS
	batchBHS                  BatchBHS
	waitBlocks                int
	lookbackBlocks            int
	latestBlock               func(ctx context.Context) (uint64, error)
	stored                    map[uint64]struct{}
	blockHeaderProvider       BlockHeaderProvider
	lastRunBlock              uint64
	getBlockhashesBatchSize   uint16
	storeBlockhashesBatchSize uint16
	gethks                    keystore.Eth
	fromAddresses             []ethkey.EIP55Address
	chainID                   *big.Int
}

// Run the feeder.
func (f *BlockHeaderFeeder) Run(ctx context.Context) error {
	latestBlockNumber, err := f.latestBlock(ctx)
	if err != nil {
		f.lggr.Errorw("Failed to fetch current block number", "error", err)
		return errors.Wrap(err, "fetching block number")
	}

	fromBlock, toBlock := blockhashstore.GetSearchWindow(int(latestBlockNumber), f.waitBlocks, f.lookbackBlocks)
	if toBlock == 0 {
		// Nothing to process, no blocks are in range.
		return nil
	}

	lggr := f.lggr.With("latestBlock", latestBlockNumber, "fromBlock", fromBlock, "toBlock", toBlock)
	lggr.Debug("searching for unfulfilled blocks")

	blockToRequests, err := blockhashstore.GetUnfulfilledBlocksAndRequests(ctx, lggr, f.coordinator, fromBlock, toBlock)
	if err != nil {
		return err
	}

	minBlockNumber := f.findLowestBlockNumberWithoutBlockhash(ctx, lggr, blockToRequests)
	if minBlockNumber == nil {
		lggr.Debug("no blocks to store")
		return nil
	}

	lggr.Debugw("found lowest block number without blockhash", "minBlockNumber", minBlockNumber)

	earliestStoredBlockNumber, err := f.findEarliestBlockNumberWithBlockhash(ctx, lggr, minBlockNumber.Uint64()+1, uint64(toBlock))
	if err != nil {
		return errors.Wrap(err, "finding earliest blocknumber with blockhash")
	}

	lggr.Debugw("found earliest block number with blockhash", "earliestStoredBlockNumber", earliestStoredBlockNumber)

	if earliestStoredBlockNumber == nil {
		// store earliest blockhash and return
		// on next iteration, earliestStoredBlockNumber will be found and
		// will make progress in storing blockhashes using blockheader.
		// In this scenario, f.stored is not updated until the next iteration
		// because we do not know which block number will be stored in the current iteration
		err = f.bhs.StoreEarliest(ctx)
		if err != nil {
			return errors.Wrap(err, "storing earliest")
		}
		lggr.Info("Stored earliest block number")
		return nil
	}

	// get the block range from (earliestStoredBlockNumber - 1) (inclusive) to minBlockNumber (inclusive) in descending order
	blocks, err := blockhashstore.DecreasingBlockRange(earliestStoredBlockNumber.Sub(earliestStoredBlockNumber, big.NewInt(1)), minBlockNumber)
	if err != nil {
		return err
	}

	// use 1 sending key for all batches because ordering matters for StoreVerifyHeader
	fromAddress, err := f.gethks.GetRoundRobinAddress(f.chainID, blockhashstore.SendingKeys(f.fromAddresses)...)
	if err != nil {
		return errors.Wrap(err, "getting round robin address")
	}

	for i := 0; i < len(blocks); i += int(f.storeBlockhashesBatchSize) {
		j := i + int(f.storeBlockhashesBatchSize)
		if j > len(blocks) {
			j = len(blocks)
		}
		blockRange := blocks[i:j]
		blockHeaders, err := f.blockHeaderProvider.RlpHeadersBatch(ctx, blockRange)
		if err != nil {
			return errors.Wrap(err, "fetching block headers")
		}

		lggr.Debugw("storing block headers", "blockRange", blockRange)
		err = f.batchBHS.StoreVerifyHeader(ctx, blockRange, blockHeaders, fromAddress)
		if err != nil {
			return errors.Wrap(err, "store block headers")
		}
		for _, blockNumber := range blockRange {
			f.stored[blockNumber.Uint64()] = struct{}{}
		}
	}

	if f.lastRunBlock != 0 {
		// Prune stored, anything older than fromBlock can be discarded
		for block := f.lastRunBlock - uint64(f.lookbackBlocks); block < fromBlock; block++ {
			if _, ok := f.stored[block]; ok {
				delete(f.stored, block)
				lggr.Debugw("Pruned block from stored cache",
					"block", block)
			}
		}
	}
	// lastRunBlock is only used for pruning
	// only time we update lastRunBlock is when the run reaches completion, indicating
	// that new block has been stored
	f.lastRunBlock = latestBlockNumber
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
			continue
		} else if stored {
			lggr.Infow("Blockhash already stored",
				"block", block, "unfulfilledReqIDs", blockhashstore.LimitReqIDs(unfulfilledReqs, 50))
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

// findEarliestBlockNumberWithBlockhash searches [startBlock, toBlock) where startBlock is inclusive and toBlock is exclusive
// and returns the first block that has blockhash already stored. Returns nil if no blockhashes are found
func (f *BlockHeaderFeeder) findEarliestBlockNumberWithBlockhash(ctx context.Context, lggr logger.Logger, startBlock, toBlock uint64) (*big.Int, error) {
	for i := startBlock; i < toBlock; i += uint64(f.getBlockhashesBatchSize) {
		j := i + uint64(f.getBlockhashesBatchSize)
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
				lggr.Infow("found earliest block number with blockhash", "earliestBlockNumber", earliestBlockNumber, "blockhash", hex.EncodeToString(bh[:]))
				f.stored[blockNumber] = struct{}{}
				return big.NewInt(0).SetUint64(earliestBlockNumber), nil
			}
		}
	}
	return nil, nil
}
