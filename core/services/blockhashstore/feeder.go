package blockhashstore

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const trustedTimeout = 1 * time.Second

// NewFeeder creates a new Feeder instance.
func NewFeeder(
	logger logger.Logger,
	coordinator Coordinator,
	bhs BHS,
	lp logpoller.LogPoller,
	trustedBHSBatchSize int32,
	waitBlocks int,
	lookbackBlocks int,
	heartbeatPeriod time.Duration,
	latestBlock func(ctx context.Context) (uint64, error),
) *Feeder {
	return &Feeder{
		lggr:                logger,
		coordinator:         coordinator,
		bhs:                 bhs,
		lp:                  lp,
		trustedBHSBatchSize: trustedBHSBatchSize,
		waitBlocks:          waitBlocks,
		lookbackBlocks:      lookbackBlocks,
		latestBlock:         latestBlock,
		stored:              make(map[uint64]struct{}),
		storedTrusted:       make(map[uint64]common.Hash),
		lastRunBlock:        0,
		wgStored:            sync.WaitGroup{},
		heartbeatPeriod:     heartbeatPeriod,
	}
}

// Feeder checks recent VRF coordinator events and stores any blockhashes for blocks within
// waitBlocks and lookbackBlocks that have unfulfilled requests.
type Feeder struct {
	lggr                logger.Logger
	coordinator         Coordinator
	bhs                 BHS
	lp                  logpoller.LogPoller
	trustedBHSBatchSize int32
	waitBlocks          int
	lookbackBlocks      int
	latestBlock         func(ctx context.Context) (uint64, error)

	// heartbeatPeriodTime is a heartbeat period in seconds by which
	// the feeder will always store a blockhash, even if there are no
	// unfulfilled requests. This is to ensure that there are blockhashes
	// in the store to start from if we ever need to run backwards mode.
	heartbeatPeriod time.Duration

	stored        map[uint64]struct{}    // used for trustless feeder
	storedTrusted map[uint64]common.Hash // used for trusted feeder
	lastRunBlock  uint64
	wgStored      sync.WaitGroup
	batchLock     sync.Mutex
	errsLock      sync.Mutex
}

//go:generate mockery --quiet --name Timer --output ./mocks/ --case=underscore
type Timer interface {
	After(d time.Duration) <-chan time.Time
}

type realTimer struct{}

func (r *realTimer) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (f *Feeder) StartHeartbeats(ctx context.Context, timer Timer) {
	if f.heartbeatPeriod == 0 {
		f.lggr.Infow("Not starting heartbeat blockhash using storeEarliest")
		return
	}
	f.lggr.Infow(fmt.Sprintf("Starting heartbeat blockhash using storeEarliest every %s", f.heartbeatPeriod.String()))
	for {
		after := timer.After(f.heartbeatPeriod)
		select {
		case <-after:
			f.lggr.Infow("storing heartbeat blockhash using storeEarliest",
				"heartbeatPeriodSeconds", f.heartbeatPeriod.Seconds())
			if err := f.bhs.StoreEarliest(ctx); err != nil {
				f.lggr.Infow("failed to store heartbeat blockhash using storeEarliest",
					"heartbeatPeriodSeconds", f.heartbeatPeriod.Seconds(),
					"err", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

// Run the feeder.
func (f *Feeder) Run(ctx context.Context) error {
	latestBlock, err := f.latestBlock(ctx)
	if err != nil {
		f.lggr.Errorw("Failed to fetch current block number", "err", err)
		return errors.Wrap(err, "fetching block number")
	}

	fromBlock, toBlock := GetSearchWindow(int(latestBlock), f.waitBlocks, f.lookbackBlocks)
	if toBlock == 0 {
		// Nothing to process, no blocks are in range.
		return nil
	}

	lggr := f.lggr.With("latestBlock", latestBlock, "fromBlock", fromBlock, "toBlock", toBlock)
	blockToRequests, err := GetUnfulfilledBlocksAndRequests(ctx, lggr, f.coordinator, fromBlock, toBlock)
	if err != nil {
		return err
	}

	// For a trusted BHS, run our trusted logic.
	if f.bhs.IsTrusted() {
		return f.runTrusted(ctx, latestBlock, fromBlock, blockToRequests)
	}

	var errs error
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
			f.lggr.Errorw("Failed to check if block is already stored, attempting to store anyway",
				"err", err,
				"block", block)
			errs = multierr.Append(errs, errors.Wrap(err, "checking if stored"))
		} else if stored {
			// IsStored() can be based on unfinalized blocks. Therefore, f.stored mapping is not updated
			f.lggr.Infow("Blockhash already stored",
				"block", block, "latestBlock", latestBlock,
				"unfulfilledReqIDs", LimitReqIDs(unfulfilledReqs, 50))
			continue
		}

		// Block needs to be stored
		err = f.bhs.Store(ctx, block)
		if err != nil {
			f.lggr.Errorw("Failed to store block", "err", err, "block", block)
			errs = multierr.Append(errs, errors.Wrap(err, "storing block"))
			continue
		}

		f.lggr.Infow("Stored blockhash",
			"block", block, "latestBlock", latestBlock,
			"unfulfilledReqIDs", LimitReqIDs(unfulfilledReqs, 50))
		f.stored[block] = struct{}{}
	}

	if f.lastRunBlock != 0 {
		// Prune stored, anything older than fromBlock can be discarded
		for block := f.lastRunBlock - uint64(f.lookbackBlocks); block < fromBlock; block++ {
			if _, ok := f.stored[block]; ok {
				delete(f.stored, block)
				f.lggr.Debugw("Pruned block from stored cache",
					"block", block, "latestBlock", latestBlock)
			}
		}
	}
	f.lastRunBlock = latestBlock
	return errs
}

func (f *Feeder) runTrusted(
	ctx context.Context,
	latestBlock uint64,
	fromBlock uint64,
	blockToRequests map[uint64]map[string]struct{},
) error {
	var errs error

	// Iterate through each request block via waitGroup.
	// For blocks with pending requests, add them to the batch to be stored.
	// Note: Golang maps sort items in a range randomly, so although the batch size is used
	// to limit blocks-per-batch, every block has an equal chance of getting picked up
	// on each run.
	var batch = make(map[uint64]struct{})
	for blockKey, unfulfilledReqs := range blockToRequests {
		f.wgStored.Add(1)
		var unfulfilled = unfulfilledReqs
		var block = blockKey
		go func() {
			defer f.wgStored.Done()
			if len(unfulfilled) == 0 {
				return
			}

			// Do not store a block if it has been marked as stored; otherwise, store it even
			// if the RPC call errors, as to be conservative.
			timeoutCtx, cancel := context.WithTimeout(ctx, trustedTimeout)
			defer cancel()
			stored, err := f.bhs.IsStored(timeoutCtx, block)
			if err != nil {
				f.lggr.Errorw("Failed to check if block is already stored, attempting to store anyway",
					"err", err,
					"block", block)
				f.errsLock.Lock()
				errs = multierr.Append(errs, errors.Wrap(err, "checking if stored"))
				f.errsLock.Unlock()
			} else if stored {
				f.lggr.Infow("Blockhash already stored",
					"block", block, "latestBlock", latestBlock,
					"unfulfilledReqIDs", LimitReqIDs(unfulfilled, 50))
				return
			}

			// If there's room, store the block in the batch. Threadsafe.
			f.batchLock.Lock()
			if len(batch) < int(f.trustedBHSBatchSize) {
				batch[block] = struct{}{}
			}
			f.batchLock.Unlock()
		}()
	}

	// Ensure all blocks are checked before storing the batch.
	f.wgStored.Wait()

	// For a non-empty batch, store all blocks.
	if len(batch) != 0 {
		var blocksToStore []uint64
		var blockhashesToStore []common.Hash
		var latestBlockhash common.Hash

		// Get all logpoller blocks for the range including the batch and the latest block,
		// as to include the recent blockhash.
		lpBlocks, err := f.lp.GetBlocksRange(ctx, append(maps.Keys(batch), latestBlock))
		if err != nil {
			f.lggr.Errorw("Failed to get blocks range",
				"err", err,
				"blocks", batch)
			errs = multierr.Append(errs, errors.Wrap(err, "log poller get blocks range"))
			return errs
		}

		// If the log poller block's blocknumber is included in the desired batch,
		// append its blockhash to our blockhashes we want to store.
		// If it is the log poller block pertaining to our recent block number, assig it.
		for _, b := range lpBlocks {
			if b.BlockNumber == int64(latestBlock) {
				latestBlockhash = b.BlockHash
			}
			if f.storedTrusted[uint64(b.BlockNumber)] == b.BlockHash {
				// blockhash is already stored. skip to save gas
				continue
			}
			if _, ok := batch[uint64(b.BlockNumber)]; ok {
				blocksToStore = append(blocksToStore, uint64(b.BlockNumber))
				blockhashesToStore = append(blockhashesToStore, b.BlockHash)
			}
		}

		if len(blocksToStore) == 0 {
			f.lggr.Debugw("no blocks to store", "latestBlock", latestBlock)
			return errs
		}
		// Store the batch of blocks and their blockhashes.
		err = f.bhs.StoreTrusted(ctx, blocksToStore, blockhashesToStore, latestBlock, latestBlockhash)
		if err != nil {
			f.lggr.Errorw("Failed to store trusted",
				"err", err,
				"blocks", blocksToStore,
				"blockhashesToStore", blockhashesToStore,
				"latestBlock", latestBlock,
				"latestBlockhash", latestBlockhash,
			)
			errs = multierr.Append(errs, errors.Wrap(err, "checking if stored"))
			return errs
		}
		for i, block := range blocksToStore {
			f.storedTrusted[block] = blockhashesToStore[i]
		}
	}

	// Prune storedTrusted, anything older than fromBlock can be discarded.
	for b := range f.storedTrusted {
		if b < fromBlock {
			delete(f.storedTrusted, b)
		}
	}

	return errs
}
