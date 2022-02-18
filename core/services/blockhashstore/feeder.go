package blockhashstore

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/logger"
)

// Event contains metadata about a VRF randomness request or fulfillment.
type Event struct {
	// ID of the relevant VRF request. For a VRF V1 request, this will an encoded 32 byte array.
	// For VRF V2, it will be an integer in string form.
	ID string

	// Block that the request or fulfillment was included in.
	Block uint64
}

// Coordinator defines an interface for fetching request and fulfillment metadata from a VRF
// coordinator.
type Coordinator interface {
	// Requests fetches VRF requests that occurred within the specified blocks.
	Requests(ctx context.Context, fromBlock uint64, toBlock uint64) ([]Event, error)

	// Fulfillments fetches VRF fulfillments that occurred since the specified block.
	Fulfillments(ctx context.Context, fromBlock uint64) ([]Event, error)
}

// BHS defines an interface for interacting with a BlockhashStore contract.
type BHS interface {
	// Store the hash associated with blockNum.
	Store(ctx context.Context, blockNum uint64) error

	// IsStored checks whether the hash associated with blockNum is already stored.
	IsStored(ctx context.Context, blockNum uint64) (bool, error)
}

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

	var (
		fromBlock        = int(latestBlock) - f.lookbackBlocks
		toBlock          = int(latestBlock) - f.waitBlocks
		blockToRequests  = make(map[uint64]map[string]struct{})
		requestIDToBlock = make(map[string]uint64)
	)
	if fromBlock < 0 {
		fromBlock = 0
	}
	if toBlock < 0 {
		// Nothing to process, no blocks are in range.
		return nil
	}
	reqs, err := f.coordinator.Requests(ctx, uint64(fromBlock), uint64(toBlock))
	if err != nil {
		f.lggr.Errorw("Failed to fetch VRF requests",
			"error", err,
			"latestBlock", latestBlock,
			"fromBlock", fromBlock,
			"toBlock", toBlock)
		return errors.Wrap(err, "fetching VRF requests")
	}
	for _, req := range reqs {
		if _, ok := blockToRequests[req.Block]; !ok {
			blockToRequests[req.Block] = make(map[string]struct{})
		}
		blockToRequests[req.Block][req.ID] = struct{}{}
		requestIDToBlock[req.ID] = req.Block
	}

	fuls, err := f.coordinator.Fulfillments(ctx, uint64(fromBlock))
	if err != nil {
		f.lggr.Errorw("Failed to fetch VRF fulfillments",
			"error", err,
			"latestBlock", latestBlock,
			"fromBlock", fromBlock,
			"toBlock", toBlock)
		return errors.Wrap(err, "fetching VRF fulfillments")
	}
	for _, ful := range fuls {
		requestBlock, ok := requestIDToBlock[ful.ID]
		if !ok {
			continue
		}
		delete(blockToRequests[requestBlock], ful.ID)
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
				"unfulfilledReqIDs", limitReqIDs(unfulfilledReqs))
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
			"unfulfilledReqIDs", limitReqIDs(unfulfilledReqs))
		f.stored[block] = struct{}{}
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
	return errs
}

// limitReqIDs converts a set of request IDs to a slice limited to 50 IDs max.
func limitReqIDs(reqs map[string]struct{}) []string {
	var reqIDs []string
	for id := range reqs {
		reqIDs = append(reqIDs, id)
		if len(reqIDs) >= 50 {
			break
		}
	}
	return reqIDs
}
