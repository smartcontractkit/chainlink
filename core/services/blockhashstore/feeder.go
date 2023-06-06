package blockhashstore

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// NewFeeder creates a new Feeder instance.
func NewFeeder(
	logger logger.Logger,
	coordinator Coordinator,
	bhs BHS,
	waitBlocks int,
	lookbackBlocks int,
	latestBlock func(ctx context.Context) (uint64, error),
) *Feeder {
	return &Feeder{
		lggr:           logger,
		coordinator:    coordinator,
		bhs:            bhs,
		waitBlocks:     waitBlocks,
		lookbackBlocks: lookbackBlocks,
		latestBlock:    latestBlock,
		stored:         make(map[uint64]struct{}),
		lastRunBlock:   0,
	}
}

// Feeder checks recent VRF coordinator events and stores any blockhashes for blocks within
// waitBlocks and lookbackBlocks that have unfulfilled requests.
type Feeder struct {
	lggr           logger.Logger
	coordinator    Coordinator
	bhs            BHS
	waitBlocks     int
	lookbackBlocks int
	latestBlock    func(ctx context.Context) (uint64, error)

	stored       map[uint64]struct{}
	lastRunBlock uint64
}

// Run the feeder.
func (f *Feeder) Run(ctx context.Context) error {
	latestBlock, err := f.latestBlock(ctx)
	if err != nil {
		f.lggr.Errorw("Failed to fetch current block number", "error", err)
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
				"error", err,
				"block", block)
			errs = multierr.Append(errs, errors.Wrap(err, "checking if stored"))
		} else if stored {
			f.lggr.Infow("Blockhash already stored",
				"block", block, "latestBlock", latestBlock,
				"unfulfilledReqIDs", LimitReqIDs(unfulfilledReqs, 50))
			f.stored[block] = struct{}{}
			continue
		}

		// Block needs to be stored
		err = f.bhs.Store(ctx, block)
		if err != nil {
			f.lggr.Errorw("Failed to store block", "error", err, "block", block)
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
