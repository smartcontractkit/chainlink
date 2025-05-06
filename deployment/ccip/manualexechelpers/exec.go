package manualexechelpers

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	ccipchangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm/manualexeclib"
)

const (
	// 14 days is the default lookback but it can be overridden.
	DefaultLookbackMessages = 14 * 24 * time.Hour

	// DefaultLookbackCommitReport is the default lookback for commit reports.
	DefaultLookbackCommitReport = 14 * 24 * time.Hour

	// DefaultStepDuration is the default duration for each filter logs query performed.
	// For lookbacks that are really far back, we need to break it up into smaller chunks.
	// This is the duration of each chunk.
	// For example, if the lookback is 14 days and the step duration is 24 hours,
	// we will make 14 queries, each for 24 hours worth of blocks.
	DefaultStepDuration = 24 * time.Hour
)

var (
	blockTimeSecondsPerChain = map[uint64]uint64{
		// simchains
		chainsel.GETH_TESTNET.Selector:  1,
		chainsel.GETH_DEVNET_2.Selector: 1,
		chainsel.GETH_DEVNET_3.Selector: 1,

		// arb
		chainsel.ETHEREUM_MAINNET_ARBITRUM_1.Selector:         1,
		chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector: 1,

		// op
		chainsel.ETHEREUM_MAINNET_OPTIMISM_1.Selector:         1,
		chainsel.ETHEREUM_TESTNET_SEPOLIA_OPTIMISM_1.Selector: 1,

		// base
		chainsel.ETHEREUM_MAINNET_BASE_1.Selector:         1,
		chainsel.ETHEREUM_TESTNET_SEPOLIA_BASE_1.Selector: 1,

		// matic
		chainsel.POLYGON_MAINNET.Selector:        2,
		chainsel.POLYGON_TESTNET_MUMBAI.Selector: 2,
		chainsel.POLYGON_TESTNET_AMOY.Selector:   2,

		// bsc
		chainsel.BINANCE_SMART_CHAIN_MAINNET.Selector: 2,
		chainsel.BINANCE_SMART_CHAIN_TESTNET.Selector: 2,

		// eth
		chainsel.ETHEREUM_MAINNET.Selector:         12,
		chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector: 12,
		chainsel.ETHEREUM_TESTNET_HOLESKY.Selector: 12,

		// avax
		chainsel.AVALANCHE_MAINNET.Selector:      3,
		chainsel.AVALANCHE_TESTNET_FUJI.Selector: 3,
	}
	// used if the chain isn't in the map above.
	defaultBlockTimeSeconds uint64 = 2
)

// getStartBlock gets the starting block of a filter logs query based on the current head block and a lookback duration.
// block time is used to calculate the number of blocks to go back.
func getStartBlock(srcChainSel uint64, currentHead uint64, lookbackDuration time.Duration) uint64 {
	blockTimeSeconds := blockTimeSecondsPerChain[srcChainSel]
	if blockTimeSeconds == 0 {
		blockTimeSeconds = defaultBlockTimeSeconds
	}

	toSub := uint64(lookbackDuration.Seconds()) / blockTimeSeconds
	if toSub > currentHead {
		return 1 // start from genesis - might happen for simchains.
	}

	start := currentHead - toSub
	return start
}

func durationToBlocks(srcChainSel uint64, lookbackDuration time.Duration) uint64 {
	blockTimeSeconds := blockTimeSecondsPerChain[srcChainSel]
	if blockTimeSeconds == 0 {
		blockTimeSeconds = defaultBlockTimeSeconds
	}

	return uint64(lookbackDuration.Seconds()) / blockTimeSeconds
}

// getCommitRootAcceptedEvent retrieves the CommitRootAccepted event for the provided
// (srcChainSel, destChainSel, msgSeqNr) triplet.
// If the commit root is not found, it returns an error.
// The lookback duration is used to determine how far back to look for the event.
// Queries are performed in steps of stepDuration worth of blocks.
func getCommitRootAcceptedEvent(
	ctx context.Context,
	lggr logger.Logger,
	env deployment.Environment,
	state ccipchangeset.CCIPOnChainState,
	srcChainSel uint64,
	destChainSel uint64,
	msgSeqNr uint64,
	lookbackDuration,
	stepDuration time.Duration,
	cachedBlockNumber uint64,
) (offramp.InternalMerkleRoot, uint64, error) {
	hdr, err := env.Chains[destChainSel].Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return offramp.InternalMerkleRoot{}, 0, fmt.Errorf("failed to get header: %w", err)
	}

	start := getStartBlock(srcChainSel, hdr.Number.Uint64(), lookbackDuration)
	if cachedBlockNumber != 0 {
		lggr.Infow("using cached block number to start search for root", "cachedBlockNumber", cachedBlockNumber)
		start = cachedBlockNumber
	}
	step := durationToBlocks(destChainSel, stepDuration)

	lggr.Infow("Getting commit root accepted event", "startBlock", start, "step", step)
	for start <= hdr.Number.Uint64() {
		end := min(start+step, hdr.Number.Uint64())
		lggr.Infow("Querying with", "startBlock", start, "endBlock", end, "step", step)

		iter, err := state.Chains[destChainSel].OffRamp.FilterCommitReportAccepted(
			&bind.FilterOpts{
				Start: start,
				End:   &end,
			},
		)
		if err != nil {
			return offramp.InternalMerkleRoot{}, 0, fmt.Errorf("failed to filter commit report accepted: %w", err)
		}

		for iter.Next() {
			if len(iter.Event.BlessedMerkleRoots) == 0 && len(iter.Event.UnblessedMerkleRoots) == 0 {
				// price updates only, can skip this event.
				continue
			}

			if root, found := findCommitRoot(
				lggr,
				append(iter.Event.BlessedMerkleRoots, iter.Event.UnblessedMerkleRoots...),
				srcChainSel,
				msgSeqNr,
				iter.Event.Raw.TxHash,
			); found {
				return root, iter.Event.Raw.BlockNumber, nil
			}
		}

		start = end + 1
	}

	lggr.Infow("didn't find commit root, maybe increase lookback duration")

	return offramp.InternalMerkleRoot{}, 0, errors.New("commit root not found")
}

// Helper function to find a commit root in a list of merkle roots
func findCommitRoot(
	lggr logger.Logger,
	roots []offramp.InternalMerkleRoot,
	srcChainSel uint64,
	msgSeqNr uint64,
	txHash common.Hash,
) (offramp.InternalMerkleRoot, bool) {
	for _, root := range roots {
		if root.SourceChainSelector == srcChainSel {
			lggr.Infow("checking commit root",
				"minSeqNr", root.MinSeqNr,
				"maxSeqNr", root.MaxSeqNr,
				"txHash", txHash.Hex(),
			)
			if msgSeqNr >= root.MinSeqNr && msgSeqNr <= root.MaxSeqNr {
				lggr.Infow("found commit root",
					"root", root,
					"txHash", txHash.Hex(),
				)
				return root, true
			}
		}
	}
	return offramp.InternalMerkleRoot{}, false
}

// getCCIPMessageSentEvents retrieves the CCIPMessageSent events for the provided
// (srcChainSel, destChainSel, merkleRoot) triplet.
// It returns a list of CCIPMessageSent events whose sequence numbers match the range in the provided root.
func getCCIPMessageSentEvents(
	ctx context.Context,
	lggr logger.Logger,
	env deployment.Environment,
	state ccipchangeset.CCIPOnChainState,
	srcChainSel uint64,
	destChainSel uint64,
	merkleRoot offramp.InternalMerkleRoot,
	lookbackDuration,
	stepDuration time.Duration,
	cachedBlockNumber uint64,
) ([]onramp.OnRampCCIPMessageSent, []uint64, error) {
	hdr, err := env.Chains[srcChainSel].Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get header: %w", err)
	}

	start := getStartBlock(srcChainSel, hdr.Number.Uint64(), lookbackDuration)
	if cachedBlockNumber != 0 {
		lggr.Infow("using cached block number to start search for messages", "cachedBlockNumber", cachedBlockNumber)
		start = cachedBlockNumber
	}
	step := durationToBlocks(srcChainSel, stepDuration)

	var seqNrs []uint64
	for i := merkleRoot.MinSeqNr; i <= merkleRoot.MaxSeqNr; i++ {
		seqNrs = append(seqNrs, i)
	}

	lggr.Infow("would query with",
		"seqNrs", seqNrs,
		"minSeqNr", merkleRoot.MinSeqNr,
		"maxSeqNr", merkleRoot.MaxSeqNr,
		"startBlock", start,
		"stepBlocks", step,
	)

	var (
		ret            []onramp.OnRampCCIPMessageSent
		merkleRootSize = merkleRoot.MaxSeqNr - merkleRoot.MinSeqNr + 1
	)

	// scan the chain until either:
	// 1. we find all the messages in the merkle root.
	// 2. we reach the current head block.
	// Usually (1) is completed before (2) but (2) is still needed to terminate the loop.
	var blockNumbers []uint64
	for uint64(len(ret)) < merkleRootSize && start <= hdr.Number.Uint64() {
		end := min(start+step, hdr.Number.Uint64())
		lggr.Infow("Querying for messages with", "startBlock", start, "endBlock", end, "stepBlocks", step)
		iter, err := state.Chains[srcChainSel].OnRamp.FilterCCIPMessageSent(
			&bind.FilterOpts{
				Start: start,
				End:   &end,
			},
			[]uint64{destChainSel},
			seqNrs,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to filter ccip message sent: %w", err)
		}

		for iter.Next() {
			if iter.Event.DestChainSelector == destChainSel {
				lggr.Infow("checking message",
					"seqNr", iter.Event.SequenceNumber,
					"destChain", iter.Event.DestChainSelector,
					"txHash", iter.Event.Raw.TxHash.String())
				if iter.Event.SequenceNumber >= merkleRoot.MinSeqNr &&
					iter.Event.SequenceNumber <= merkleRoot.MaxSeqNr {
					ret = append(ret, *iter.Event)
					blockNumbers = append(blockNumbers, iter.Event.Raw.BlockNumber)
				}
			}
		}

		start = end + 1
	}

	if len(ret) != len(seqNrs) {
		return nil, nil, fmt.Errorf("not all messages found, got: %d, expected: %d", len(ret), len(seqNrs))
	}
	if len(blockNumbers) != len(seqNrs) {
		return nil, nil, fmt.Errorf("not all block numbers found, got: %d, expected: %d", len(blockNumbers), len(seqNrs))
	}

	lggr.Infow("found all messages for root", "merkleRoot", merkleRoot, "messages", len(ret))

	return ret, blockNumbers, nil
}

// manuallyExecuteSingle manually executes a single message on the destination chain.
func manuallyExecuteSingle(
	ctx context.Context,
	lggr logger.Logger,
	state ccipchangeset.CCIPOnChainState,
	env deployment.Environment,
	srcChainSel uint64,
	destChainSel uint64,
	msgSeqNr uint64,
	lookbackDuration,
	lookbackDurationCommitReport,
	stepDuration time.Duration,
	reExecuteIfFailed bool,
	extraDataCodec ccipcommon.ExtraDataCodec,
	messageSentCache *MessageSentCache,
	commitRootCache *RootCache,
) error {
	onRampAddress := state.Chains[srcChainSel].OnRamp.Address()

	execState, err := state.Chains[destChainSel].OffRamp.GetExecutionState(&bind.CallOpts{
		Context: ctx,
	}, srcChainSel, msgSeqNr)
	if err != nil {
		return fmt.Errorf("failed to get execution state: %w", err)
	}

	if execState == testhelpers.EXECUTION_STATE_SUCCESS ||
		(execState == testhelpers.EXECUTION_STATE_FAILURE && !reExecuteIfFailed) {
		lggr.Infow("message already executed", "execState", execState, "msgSeqNr", msgSeqNr)
		return nil
	}

	lggr.Infow("contract addresses",
		"offRampAddress", state.Chains[destChainSel].OffRamp.Address(),
		"onRampAddress", onRampAddress,
		"execState", execState,
	)

	merkleRoot, inCache := commitRootCache.Get(msgSeqNr)
	if !inCache {
		latestBlockNumber := commitRootCache.GetLatestBlockNumber()
		lggr.Infow("merkle root not found in cache, fetching from the chain",
			"msgSeqNr", msgSeqNr,
			"latestBlockNumber", latestBlockNumber,
		)
		var err error
		var blockNumber uint64
		merkleRoot, blockNumber, err = getCommitRootAcceptedEvent(
			ctx,
			lggr,
			env,
			state,
			srcChainSel,
			destChainSel,
			msgSeqNr,
			lookbackDurationCommitReport,
			stepDuration,
			latestBlockNumber,
		)
		if err != nil {
			return fmt.Errorf("failed to get merkle root: %w", err)
		}

		// add to the cache for faster fetching.
		commitRootCache.Add(RootCacheEntry{
			Root:        merkleRoot,
			BlockNumber: blockNumber,
		})
		commitRootCache.Build()
	} else {
		lggr.Infow("found merkle root in cache", "msgSeqNr", msgSeqNr, "merkleRoot", merkleRoot)
	}

	lggr.Infow("merkle root",
		"merkleRoot", hexutil.Encode(merkleRoot.MerkleRoot[:]),
		"minSeqNr", merkleRoot.MinSeqNr,
		"maxSeqNr", merkleRoot.MaxSeqNr,
		"sourceChainSel", merkleRoot.SourceChainSelector,
	)

	merkleRootSize := merkleRoot.MaxSeqNr - merkleRoot.MinSeqNr + 1
	var ccipMessageSentEvents []onramp.OnRampCCIPMessageSent
	if _, ok := messageSentCache.Get(msgSeqNr); ok {
		lggr.Infow("found message in cache, fetching the rest", "msgSeqNr", msgSeqNr)
		for start := merkleRoot.MinSeqNr; start <= merkleRoot.MaxSeqNr; start++ {
			message, ok := messageSentCache.Get(start)
			if !ok {
				lggr.Infow("message not found in cache, fetching from the chain", "msgSeqNr", start)
				break
			}
			ccipMessageSentEvents = append(ccipMessageSentEvents, message)
		}

		if uint64(len(ccipMessageSentEvents)) != merkleRootSize {
			latestBlockNumber := messageSentCache.GetLatestBlockNumber()
			lggr.Infow("not all messages found in cache, fetching from the chain", "msgSeqNr", msgSeqNr, "latestBlockNumber", latestBlockNumber)
			var err error
			var blockNumbers []uint64
			ccipMessageSentEvents, blockNumbers, err = getCCIPMessageSentEvents(
				ctx,
				lggr,
				env,
				state,
				srcChainSel,
				destChainSel,
				merkleRoot,
				lookbackDuration,
				stepDuration,
				latestBlockNumber,
			)
			if err != nil {
				return fmt.Errorf("failed to get ccip message sent event: %w", err)
			}

			if len(ccipMessageSentEvents) != len(blockNumbers) {
				return fmt.Errorf("unexpected mismatch in message count (%d) and block number count (%d), msgSeqNr: %d",
					len(ccipMessageSentEvents), len(blockNumbers), msgSeqNr)
			}

			// add to the cache
			for i, event := range ccipMessageSentEvents {
				messageSentCache.Add(event, blockNumbers[i])
			}
		}
	} else {
		latestBlockNumber := messageSentCache.GetLatestBlockNumber()
		lggr.Infow("not found in cache, fetching from the chain", "msgSeqNr", msgSeqNr, "latestBlockNumber", latestBlockNumber)
		var err error
		var blockNumbers []uint64
		ccipMessageSentEvents, blockNumbers, err = getCCIPMessageSentEvents(
			ctx,
			lggr,
			env,
			state,
			srcChainSel,
			destChainSel,
			merkleRoot,
			lookbackDuration,
			stepDuration,
			latestBlockNumber,
		)
		if err != nil {
			return fmt.Errorf("failed to get ccip message sent event: %w", err)
		}

		if len(ccipMessageSentEvents) != len(blockNumbers) {
			return fmt.Errorf("unexpected mismatch in message count (%d) and block number count (%d), msgSeqNr: %d",
				len(ccipMessageSentEvents), len(blockNumbers), msgSeqNr)
		}

		// add to the cache
		for i, event := range ccipMessageSentEvents {
			messageSentCache.Add(event, blockNumbers[i])
		}
	}

	messageHashes, err := manualexeclib.GetMessageHashes(
		ctx,
		lggr,
		onRampAddress,
		ccipMessageSentEvents,
		extraDataCodec,
	)
	if err != nil {
		return fmt.Errorf("failed to get message hashes: %w", err)
	}

	hashes, flags, err := manualexeclib.GetMerkleProof(
		lggr,
		merkleRoot,
		messageHashes,
		msgSeqNr,
	)
	if err != nil {
		return fmt.Errorf("failed to get merkle proof: %w", err)
	}

	lggr.Infow("got hashes and flags", "hashes", hashes, "flags", flags)

	// since we're only executing one message, we need to only include that message
	// in the report.
	var filteredMsgSentEvents []onramp.OnRampCCIPMessageSent
	for i, event := range ccipMessageSentEvents {
		if event.Message.Header.SequenceNumber == msgSeqNr && event.Message.Header.SourceChainSelector == srcChainSel {
			filteredMsgSentEvents = append(filteredMsgSentEvents, ccipMessageSentEvents[i])
		}
	}

	// sanity check, should not be possible at this point.
	if len(filteredMsgSentEvents) == 0 {
		return fmt.Errorf("no message found for seqNr %d", msgSeqNr)
	}

	// TODO: CCIP-5702
	execReport, err := manualexeclib.CreateExecutionReport(
		srcChainSel,
		onRampAddress,
		filteredMsgSentEvents,
		hashes,
		flags,
		extraDataCodec,
	)
	if err != nil {
		return fmt.Errorf("failed to create execution report: %w", err)
	}

	txOpts := &bind.TransactOpts{
		From:   env.Chains[destChainSel].DeployerKey.From,
		Nonce:  nil,
		Signer: env.Chains[destChainSel].DeployerKey.Signer,
		Value:  big.NewInt(0),
		// We manually set the gas limit here because estimateGas doesn't take into account
		// internal reverts (such as those that could happen on ERC165 interface checks).
		// This is just a big value for now, we can investigate something more efficient later.
		GasLimit: 1e6,
	}
	tx, err := state.Chains[destChainSel].OffRamp.ManuallyExecute(
		txOpts,
		[]offramp.InternalExecutionReport{execReport},
		[][]offramp.OffRampGasLimitOverride{
			{
				{
					ReceiverExecutionGasLimit: big.NewInt(200_000),
					TokenGasOverrides:         nil,
				},
			},
		},
	)
	_, err = deployment.ConfirmIfNoErrorWithABI(env.Chains[destChainSel], tx, offramp.OffRampABI, err)
	if err != nil {
		return fmt.Errorf("failed to execute message: %w", err)
	}

	lggr.Infow("successfully manually executed msg", "msgSeqNr", msgSeqNr)

	return nil
}

// ManuallyExecuteAll will manually execute the provided messages if they were not already executed.
// At the moment offchain token data (i.e USDC/Lombard attestations) is not supported,
// and only EVM is supported as a source or dest chain.
func ManuallyExecuteAll(
	ctx context.Context,
	lggr logger.Logger,
	state ccipchangeset.CCIPOnChainState,
	env deployment.Environment,
	srcChainSel uint64,
	destChainSel uint64,
	msgSeqNrs []int64,
	lookbackDurationMsgs,
	lookbackDurationCommitReport,
	stepDuration time.Duration,
	reExecuteIfFailed bool,
) error {
	extraDataCodec := ccipcommon.NewExtraDataCodec(ccipcommon.NewExtraDataCodecParams(ccipevm.ExtraDataDecoder{}, ccipsolana.ExtraDataDecoder{}))

	// for large backfills, these caches can speed things up because we don't need to query
	// the chain multiple times for the same root/messages.
	messageSentCache := NewMessageSentCache()
	commitRootCache := NewRootCache()
	for _, seqNr := range msgSeqNrs {
		err := manuallyExecuteSingle(
			ctx,
			lggr,
			state,
			env,
			srcChainSel,
			destChainSel,
			uint64(seqNr), //nolint:gosec // seqNr is never <= 0.
			lookbackDurationMsgs,
			lookbackDurationCommitReport,
			stepDuration,
			reExecuteIfFailed,
			extraDataCodec,
			messageSentCache,
			commitRootCache,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckAlreadyExecuted will check the execution state of the provided messages and log if they were already executed.
func CheckAlreadyExecuted(
	ctx context.Context,
	lggr logger.Logger,
	state ccipchangeset.CCIPOnChainState,
	srcChainSel uint64,
	destChainSel uint64,
	msgSeqNrs []int64,
) error {
	for _, seqNr := range msgSeqNrs {
		execState, err := state.Chains[destChainSel].OffRamp.GetExecutionState(
			&bind.CallOpts{Context: ctx},
			srcChainSel,
			uint64(seqNr), //nolint:gosec // seqNr is never <= 0.
		)
		if err != nil {
			return fmt.Errorf("failed to get execution state: %w", err)
		}

		if execState == testhelpers.EXECUTION_STATE_SUCCESS || execState == testhelpers.EXECUTION_STATE_FAILURE {
			lggr.Infow("message already executed", "execState", execState, "msgSeqNr", seqNr)
		} else {
			lggr.Infow("message not executed", "execState", execState, "msgSeqNr", seqNr)
		}
	}

	return nil
}

type RootCacheEntry struct {
	Root        offramp.InternalMerkleRoot
	BlockNumber uint64
}

// RootCache caches merkle roots for fast lookups via a single sequence number.
// For example, if we have
// cache = {(1, 10), (20, 30), (40, 50)}
// Get(5) = (1, 10)
// Get(25) = (20, 30)
// Get(35) = nil
type RootCache struct {
	roots []RootCacheEntry
}

// NewRootCache creates a new RootCache
// with an empty list of ranges.
func NewRootCache() *RootCache {
	return &RootCache{
		roots: []RootCacheEntry{},
	}
}

// Add a root to the RootCache
func (rm *RootCache) Add(entry RootCacheEntry) {
	rm.roots = append(rm.roots, entry)
}

// Build prepares the RootCache for fast lookups (sorts by MinSeqNr)
func (rm *RootCache) Build() {
	sort.Slice(rm.roots, func(i, j int) bool {
		return rm.roots[i].Root.MinSeqNr < rm.roots[j].Root.MinSeqNr
	})
}

// Get finds the commit root associated that the given sequence number belongs to
// returns the root and a boolean indicating if it was found.
func (rm *RootCache) Get(seqNr uint64) (offramp.InternalMerkleRoot, bool) {
	// Use binary search to find the range
	idx := sort.Search(len(rm.roots), func(i int) bool {
		return rm.roots[i].Root.MinSeqNr > seqNr
	}) - 1

	if idx >= 0 && idx < len(rm.roots) && seqNr >= rm.roots[idx].Root.MinSeqNr && seqNr <= rm.roots[idx].Root.MaxSeqNr {
		return rm.roots[idx].Root, true
	}
	return offramp.InternalMerkleRoot{}, false
}

func (rm *RootCache) GetLatestBlockNumber() uint64 {
	if len(rm.roots) == 0 {
		return 0
	}

	return rm.roots[len(rm.roots)-1].BlockNumber
}

type MessageSentCache struct {
	messageSentCache  map[uint64]onramp.OnRampCCIPMessageSent
	latestBlockNumber uint64
}

func NewMessageSentCache() *MessageSentCache {
	return &MessageSentCache{messageSentCache: make(map[uint64]onramp.OnRampCCIPMessageSent)}
}

func (msc *MessageSentCache) Add(entry onramp.OnRampCCIPMessageSent, blockNumber uint64) {
	msc.messageSentCache[entry.Message.Header.SequenceNumber] = entry
	if blockNumber > msc.latestBlockNumber {
		msc.latestBlockNumber = blockNumber
	}
}

func (msc *MessageSentCache) GetLatestBlockNumber() uint64 {
	return msc.latestBlockNumber
}

func (msc *MessageSentCache) Get(seqNr uint64) (onramp.OnRampCCIPMessageSent, bool) {
	entry, ok := msc.messageSentCache[seqNr]
	return entry, ok
}
