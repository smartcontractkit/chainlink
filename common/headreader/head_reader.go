package headreader

import (
	"context"
	"database/sql"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"
	htrktypes "github.com/smartcontractkit/chainlink/v2/common/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)


type LogsPoller[BLOCK_HASH types.Hashable] interface {
	Backfill(ctx context.Context, currentBlockNumber int64, lastSafeBackfillBlock int64) error
	FilterLogs(ctx context.Context, hash BLOCK_HASH) ([]Log, error)
}

type HeadReaderConfig interface {
	HistoryDepth() uint32
	MaxBufferSize() uint32
	SamplingInterval() time.Duration
}

type HeadReaderClient[
	HTH htrktypes.Head[BLOCK_HASH, ID],
	ID types.ID,
	BLOCK_HASH types.Hashable,
] interface {
	HeadByNumber(ctx context.Context, n *big.Int) (HTH, error)
	HeadByHash(ctx context.Context, n BLOCK_HASH) (HTH, error)
	LatestBlockByType(ctx context.Context, finalityType string) (HTH, error)
	ConfiguredChainID() *big.Int
}

type HeadReader[
	HTH htrktypes.Head[BLOCK_HASH, ID],
	ID types.ID,
	BLOCK_HASH types.Hashable,
] interface {
	Start(ctx context.Context) error
	RetrieveAndSaveHeads(ctx context.Context, currentBlockNumber int64) HTH
	LatestBlockProcessed() (HTH, error)
}

// HeadsBufferSize - The buffer is used when heads sampling is disabled, to ensure the callback is run for every head
const HeadsBufferSize = 10

type headReader[
	HTH htrktypes.Head[BLOCK_HASH, ID],
	ID types.ID,
	BLOCK_HASH types.Hashable,
] struct {
	utils.StartStopOnce
	hs              HeadSaver[HTH, ID, BLOCK_HASH]
	lggr            logger.Logger
	hrConfig        HeadReaderConfig
	client          HeadReaderClient[HTH, ID, BLOCK_HASH]
	headBroadcaster types.HeadBroadcaster[HTH, BLOCK_HASH]
	headListener    types.HeadListener[HTH, BLOCK_HASH]
	broadcastMB     *utils.Mailbox[HTH]
	chStop          utils.StopChan
	ctx             context.Context
	logsPoller      LogsPoller[BLOCK_HASH]
	getNilHead      func() HTH
	wgDone          sync.WaitGroup
}

func NewHeadReader[
	HTH htrktypes.Head[BLOCK_HASH, ID],
	S types.Subscription,
	ID types.ID,
	BLOCK_HASH types.Hashable,
](hs HeadSaver[HTH, ID, BLOCK_HASH],
	lggr logger.Logger,
	hrConfig HeadReaderConfig,
	client HeadReaderClient[HTH, ID, BLOCK_HASH],
	headBroadcaster types.HeadBroadcaster[HTH, BLOCK_HASH],
	headListener types.HeadListener[HTH, BLOCK_HASH],
	logsPoller LogsPoller[BLOCK_HASH],
	getNilHead func() HTH,
) *headReader[HTH, ID, BLOCK_HASH] {
	chStop := make(chan struct{})
	return &headReader[HTH, ID, BLOCK_HASH]{
		hs:              hs,
		lggr:            lggr.Named("HeadReader"),
		hrConfig:        hrConfig,
		client:          client,
		headBroadcaster: headBroadcaster,
		headListener:    headListener,
		broadcastMB:     utils.NewMailbox[HTH](HeadsBufferSize),
		chStop:          chStop,
		logsPoller:      logsPoller,
		getNilHead:      getNilHead,
	}
}

func (hr *headReader[HTH, ID, BLOCK_HASH]) Start(ctx context.Context) error {
	return hr.StartOnce("HeadReader", func() error {
		initialHead, err := hr.getInitialHead(ctx)
		if err != nil {
			if errors.Is(err, ctx.Err()) {
				return nil
			}
			hr.lggr.Errorw("Error getting initial head", "err", err)
		} else if initialHead.IsValid() {
			if err := hr.handleNewHead(ctx, initialHead); err != nil {
				return errors.Wrap(err, "error handling initial head")
			}
		} else {
			hr.lggr.Debug("Got nil initial head")
		}

		hr.wgDone.Add(2)
		go hr.headListener.ListenForNewHeads(hr.handleNewHead, hr.wgDone.Done)
		go hr.broadcastLoop()

		return nil
	})
}

func (hr *headReader[HTH, ID, BLOCK_HASH]) getInitialHead(ctx context.Context) (HTH, error) {
	head, err := hr.client.HeadByNumber(ctx, nil)
	if err != nil {
		return hr.getNilHead(), errors.Wrap(err, "failed to fetch initial head")
	}
	loggerFields := []interface{}{"head", head}
	if head.IsValid() {
		loggerFields = append(loggerFields, "blockNumber", head.BlockNumber(), "blockHash", head.BlockHash())
	}
	hr.lggr.Debugw("Got initial head", loggerFields...)
	return head, nil
}

func (hr *headReader[HTH, ID, BLOCK_HASH]) handleNewHead(ctx context.Context, latestHead HTH) error {
	var start int64
	var latestFinalizedBlock HTH
	lastProcessed, err := hr.hs.SelectLatestBlock()
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			// Assume transient db reading issue, retry forever.
			hr.lggr.Errorw("unable to get starting block", "err", err)
			return err
		}
		// Otherwise this is the first poll _ever_ on a new chain.
		// Only safe thing to do is to start at the first finalized block.
		latest := latestHead
		if !latest.IsValid() {
			hr.lggr.Warnw("Unable to get latest for first poll", "err", err)
			return err
		}
		latestNum := latest.BlockNumber()
		latestFinalizedBlock, err = hr.client.LatestBlockByType(hr.ctx, "finalized")
		if err != nil {
			hr.lggr.Warnw("Unable to get latestFinalizedBlock", "err", err)
		}
		// Do not support polling chains which don't even have finality depth worth of blocks.
		// Could conceivably support this but not worth the effort.
		// Need finality depth + 1, no block 0.
		if latestNum <= latestFinalizedBlock.BlockNumber() {
			hr.lggr.Warnw("Insufficient number of blocks on chain, waiting for finality depth", "err", err, "latest", latestNum, "finaled", latestFinalizedBlock.BlockNumber())
			return err
		}
		// Starting at the first finalized block. We do backfill the first finalized block.
		start = latestNum - latestFinalizedBlock.BlockNumber()
	} else {
		start = lastProcessed.BlockNumber() + 1
	}
	hr.RetrieveAndSaveHeads(hr.ctx, start, latestHead, latestFinalizedBlock)
	return nil
}

func (hr *headReader[HTH, ID, BLOCK_HASH]) RetrieveAndSaveHeads(ctx context.Context, currentBlockNumber int64, latestHead HTH, latestFinalizedBlock HTH) {
	hr.lggr.Debugw("Retrieving heads", "currentBlockNumber", currentBlockNumber)
	var err error
	latestBlock := latestHead
	if !latestBlock.IsValid() {
		hr.lggr.Warnw("Unable to get latestBlockNumber block", "err", err, "currentBlockNumber", currentBlockNumber)
		return
	}
	if !latestFinalizedBlock.IsValid() {
		latestFinalizedBlock, err = hr.client.LatestBlockByType(ctx, "finalized")
		if err != nil {
			hr.lggr.Warnw("Unable to get latestFinalizedBlock", "err", err)
		}
	}

	latestBlockNumber := latestBlock.BlockNumber()
	if currentBlockNumber > latestBlockNumber {
		// Note there can also be a reorg "shortening" i.e. chain height decreases but TDD increases. In that case
		// we also just wait until the new tip is longer and then detect the reorg.
		hr.lggr.Debugw("No new blocks since last poll", "currentBlockNumber", currentBlockNumber, "latestBlockNumber", latestBlockNumber)
		return
	}
	var currentBlock HTH
	if currentBlockNumber == latestBlockNumber {
		// Can re-use our currentBlock and avoid an extra RPC call.
		currentBlock = latestBlock
	}
	// Possibly handle a reorg. For example if we crash, we'll be in the middle of processing unfinalized blocks.
	// Returns (currentBlock || LCA+1 if reorg detected, error)
	currentBlock, err = hr.getCurrentBlockMaybeHandleReorg(ctx, currentBlockNumber, currentBlock, latestFinalizedBlock)
	if err != nil {
		// If there's an error handling the reorg, we can't be sure what state the db was left in.
		// Resume from the latest block saved and retry.
		hr.lggr.Errorw("Unable to get current block, retrying", "err", err)
		return
	}
	currentBlockNumber = currentBlock.BlockNumber()

	// backfill finalized blocks if we can for performance. If we crash during backfill, we
	// may reprocess logs.  Log insertion is idempotent so this is ok.
	// E.g. 1<-2<-3(currentBlockNumber)<-4<-5<-6<-7(latestBlockNumber), finality is 2. So 3,4 can be batched.
	// Although 5 is finalized, we still need to save it to the db for reorg detection if 6 is a reorg.
	// start = currentBlockNumber = 3, end = latestBlockNumber - finality - 1 = 7-2-1 = 4 (inclusive range).
	lastSafeBackfillBlock := latestBlockNumber - latestFinalizedBlock.BlockNumber() - 1
	if lastSafeBackfillBlock >= currentBlockNumber {
		hr.lggr.Infow("Backfilling logs", "start", currentBlockNumber, "end", lastSafeBackfillBlock)
		if err = hr.logsPoller.Backfill(ctx, currentBlockNumber, lastSafeBackfillBlock); err != nil {
			// If there's an error backfilling, we can just return and retry from the last block saved
			// since we don't save any blocks on backfilling. We may re-insert the same logs but thats ok.
			hr.lggr.Warnw("Unable to backfill finalized logs, retrying later", "err", err)
			return
		}
		currentBlockNumber = lastSafeBackfillBlock + 1
	}

	for {
		h := currentBlock
		var logs []Log
		logs, err = hr.logsPoller.FilterLogs(ctx, h.BlockHash())
		if err != nil {
			hr.lggr.Warnw("Unable to query for logs, retrying", "err", err, "block", currentBlockNumber)
			return
		}
		hr.lggr.Debugw("Unfinalized log query", "logs", len(logs),
			"currentBlockNumber", currentBlockNumber, "blockHash", currentBlock.BlockHash(), "timestamp", currentBlock.BlockTimestamp().Unix())
		err = hr.hs.Store(h, logs)
		if err != nil {
			hr.lggr.Warnw("Unable to save logs resuming from last saved block + 1", "err", err, "block", currentBlockNumber)
			return
		}
		// Update current block.
		// Same reorg detection on unfinalized blocks.
		currentBlockNumber++
		if currentBlockNumber > latestBlockNumber {
			break
		}
		var emptyHead HTH
		currentBlock, err = hr.getCurrentBlockMaybeHandleReorg(ctx, currentBlockNumber, emptyHead, latestFinalizedBlock)
		if err != nil {
			// If there's an error handling the reorg, we can't be sure what state the db was left in.
			// Resume from the latest block saved.
			hr.lggr.Errorw("Unable to get current block", "err", err)
			return
		}
		currentBlockNumber = currentBlock.BlockNumber()
	}
	// headWithChain
	var emptyHead HTH
	hr.broadcastMB.Deliver(emptyHead)
}

// getCurrentBlockMaybeHandleReorg accepts a block number
// and will return that block if its parent points to our last saved block.
// One can optionally pass the block header if it has already been queried to avoid an extra RPC call.
// If its parent does not point to our last saved block we know a reorg has occurred,
// so we:
// 1. Find the LCA by following parent hashes.
// 2. Delete all logs and blocks after the LCA
// 3. Return the LCA+1, i.e. our new current (unprocessed) block.
func (hr *headReader[HTH, ID, BLOCK_HASH]) getCurrentBlockMaybeHandleReorg(ctx context.Context, currentBlockNumber int64, currentBlock HTH, latestFinalizedBlock HTH) (HTH, error) {
	var err1 error
	var emptyHead HTH
	if !currentBlock.IsValid() {
		// If we don't have the current block already, lets get it.
		currentBlock, err1 = hr.client.HeadByNumber(ctx, big.NewInt(currentBlockNumber))
		if err1 != nil {
			hr.lggr.Warnw("Unable to get currentBlock", "err", err1, "currentBlockNumber", currentBlockNumber)
			return emptyHead, err1
		}
		// Additional sanity checks, don't necessarily trust the RPC.
		if !currentBlock.IsValid() {
			hr.lggr.Errorf("Unexpected nil block from RPC", "currentBlockNumber", currentBlockNumber)
			return emptyHead, errors.Errorf("Got nil block for %d", currentBlockNumber)
		}
		if currentBlock.BlockNumber() != currentBlockNumber {
			hr.lggr.Warnw("Unable to get currentBlock, rpc returned incorrect block", "currentBlockNumber", currentBlockNumber, "got", currentBlock.BlockNumber())
			return emptyHead, errors.Errorf("Block mismatch have %d want %d", currentBlock.BlockNumber(), currentBlockNumber)
		}
	}
	// Does this currentBlock point to the same parent that we have saved?
	// If not, there was a reorg, so we need to rewind.
	expectedParent, err1 := hr.hs.SelectBlockByNumber(currentBlockNumber - 1)
	if err1 != nil && !errors.Is(err1, sql.ErrNoRows) {
		// If err is not a 'no rows' error, assume transient db issue and retry
		hr.lggr.Warnw("Unable to read latestBlockNumber currentBlock saved", "err", err1, "currentBlockNumber", currentBlockNumber)
		return emptyHead, errors.New("Unable to read latestBlockNumber currentBlock saved")
	}
	// We will not have the previous currentBlock on initial poll.
	havePreviousBlock := err1 == nil
	if !havePreviousBlock {
		hr.lggr.Infow("Do not have previous block, first poll ever on new chain or after backfill", "currentBlockNumber", currentBlockNumber)
		return currentBlock, nil
	}
	// Check for reorg.
	if currentBlock.GetParentHash() != expectedParent.BlockHash() {
		// There can be another reorg while we're finding the LCA.
		// That is ok, since we'll detect it on the next iteration.
		// Since we go currentBlock by currentBlock for unfinalized logs, the mismatch starts at currentBlockNumber - 1.
		blockAfterLCA, err2 := hr.findBlockAfterLCA(ctx, currentBlock, latestFinalizedBlock)
		if err2 != nil {
			hr.lggr.Warnw("Unable to find LCA after reorg, retrying", "err", err2)
			return emptyHead, errors.New("Unable to find LCA after reorg, retrying")
		}

		hr.lggr.Infow("Reorg detected", "blockAfterLCA", blockAfterLCA.BlockNumber(), "currentBlockNumber", currentBlockNumber)
		// We truncate all the blocks and logs after the LCA.
		// We could preserve the logs for forensics, since its possible
		// that applications see them and take action upon it, however that
		// results in significantly slower reads since we must then compute
		// the canonical set per read. Typically, if an application took action on a log
		// it would be saved elsewhere e.g. eth_txes, so it seems better to just support the fast reads.
		// Its also nicely analogous to reading from the chain itself.
		err2 = hr.hs.Delete(blockAfterLCA.BlockNumber())
		if err2 != nil {
			// If we error on db commit, we can't know if the tx went through or not.
			// We return an error here which will cause us to restart polling from lastBlockSaved + 1
			return emptyHead, err2
		}
		return blockAfterLCA, nil
	}
	// No reorg, return current block.
	return currentBlock, nil
}

// Find the first place where our chain and their chain have the same block,
// that block number is the LCA. Return the block after that, where we want to resume polling.
func (hr *headReader[HTH, ID, BLOCK_HASH]) findBlockAfterLCA(ctx context.Context, current HTH, latestFinalizedBlock HTH) (HTH, error) {
	var emptyHead HTH
	// Current is where the mismatch starts.
	// Check its parent to see if its the same as ours saved.
	parent, err := hr.client.HeadByHash(ctx, current.GetParentHash())
	if err != nil {
		return emptyHead, err
	}
	blockAfterLCA := current
	reorgStart := parent.BlockNumber()
	// We expect reorgs up to the block after (current - finalityDepth),
	// since the block at (current - finalityDepth) is finalized.
	// We loop via parent instead of current so current always holds the LCA+1.
	// If the parent block number becomes < the first finalized block our reorg is too deep.
	for parent.BlockNumber() >= (reorgStart - latestFinalizedBlock.BlockNumber()) {
		ourParentBlockHash, err := hr.hs.SelectBlockByNumber(parent.BlockNumber())
		if err != nil {
			return emptyHead, err
		}
		if parent.BlockHash() == ourParentBlockHash.BlockHash() {
			// If we do have the blockhash, return blockAfterLCA
			return blockAfterLCA, nil
		}
		// Otherwise get a new parent and update blockAfterLCA.
		blockAfterLCA = parent
		parent, err = hr.client.HeadByHash(ctx, parent.GetParentHash())
		if err != nil {
			return emptyHead, err
		}
	}
	hr.lggr.Criticalw("Reorg greater than finality depth detected", "max reorg depth", latestFinalizedBlock.BlockNumber())
	rerr := errors.New("Reorg greater than finality depth")
	hr.SvcErrBuffer.Append(rerr)
	return emptyHead, rerr
}

func (hr *headReader[HTH, ID, BLOCK_HASH]) broadcastLoop() {
	defer hr.wgDone.Done()

	samplingInterval := hr.hrConfig.SamplingInterval()
	if samplingInterval > 0 {
		hr.lggr.Debugf("Head sampling is enabled - sampling interval is set to: %v", samplingInterval)
		debounceHead := time.NewTicker(samplingInterval)
		defer debounceHead.Stop()
		for {
			select {
			case <-hr.chStop:
				return
			case <-debounceHead.C:
				item := hr.broadcastMB.RetrieveLatestAndClear()
				if !item.IsValid() {
					continue
				}
				hr.headBroadcaster.BroadcastNewLongestChain(item)
			}
		}
	} else {
		hr.lggr.Info("Head sampling is disabled - callback will be called on every head")
		for {
			select {
			case <-hr.chStop:
				return
			case <-hr.broadcastMB.Notify():
				for {
					item, exists := hr.broadcastMB.Retrieve()
					if !exists {
						break
					}
					hr.headBroadcaster.BroadcastNewLongestChain(item)
				}
			}
		}
	}
}
