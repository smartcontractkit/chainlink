package coordinator

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/libocr/commontypes"
	ocr2Types "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/ocr2vrf/dkg"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"
	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/proto"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	dkg_wrapper "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/dkg"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/vrf_beacon"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/vrf_coordinator"
	vrf_wrapper "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/vrf_coordinator"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	ocr2vrfconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2vrf/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils/mathutil"
)

var _ ocr2vrftypes.CoordinatorInterface = &coordinator{}

var (
	dkgABI            = evmtypes.MustGetABI(dkg_wrapper.DKGMetaData.ABI)
	vrfBeaconABI      = evmtypes.MustGetABI(vrf_beacon.VRFBeaconMetaData.ABI)
	vrfCoordinatorABI = evmtypes.MustGetABI(vrf_coordinator.VRFCoordinatorMetaData.ABI)
	counterBuckets    = []float64{
		0,
		1,
		2,
		4,
		8,
		16,
		32,
		64,
		128,
		256,
		512,
		1024,
		2048,
	}
	timingBuckets = []float64{
		float64(1 * time.Millisecond),
		float64(5 * time.Millisecond),
		float64(10 * time.Millisecond),
		float64(50 * time.Millisecond),
		float64(100 * time.Millisecond),
		float64(500 * time.Millisecond),
		float64(time.Second),
		float64(5 * time.Second),
		float64(10 * time.Second),
		float64(30 * time.Second),
	}
	promLabels         = []string{"evmChainID", "oracleID", "configDigest"}
	promBlocksToReport = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ocr2vrf_coordinator_blocks_to_report",
		Help:    "Number of unfulfilled and in-flight blocks that fit in current report in reportBlocks",
		Buckets: counterBuckets,
	}, promLabels)
	promCallbacksToReport = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ocr2vrf_coordinator_callbacks_to_report",
		Help:    "Number of unfulfilled and and in-flight callbacks fit in current report in reportBlocks",
		Buckets: counterBuckets,
	}, promLabels)
	promBlocksInReport = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ocr2vrf_coordinator_blocks_in_report",
		Help:    "Number of blocks found in reportWillBeTransmitted",
		Buckets: counterBuckets,
	}, promLabels)
	promCallbacksInReport = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ocr2vrf_coordinator_callbacks_in_report",
		Help:    "Number of callbacks found in reportWillBeTransmitted",
		Buckets: counterBuckets,
	}, promLabels)
	promMethodDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ocr2vrf_coordinator_method_time",
		Help:    "The amount of time elapsed for given method call",
		Buckets: timingBuckets,
	}, append(promLabels, "methodName"))
)

const (
	// VRF-only events.
	randomnessRequestedEvent            string = "RandomnessRequested"
	randomnessFulfillmentRequestedEvent string = "RandomnessFulfillmentRequested"
	randomWordsFulfilledEvent           string = "RandomWordsFulfilled"
	newTransmissionEvent                string = "NewTransmission"
	outputsServedEvent                  string = "OutputsServed"

	// Both VRF and DKG contracts emit this, it's an OCR event.
	configSetEvent = "ConfigSet"
)

// block is used to key into a set that tracks beacon blocks.
type block struct {
	blockNumber uint64
	confDelay   uint32
}

type blockInReport struct {
	block
	recentBlockHeight uint64
	recentBlockHash   common.Hash
}

type callback struct {
	blockNumber uint64
	requestID   uint64
}

type callbackInReport struct {
	callback
	recentBlockHeight uint64
	recentBlockHash   common.Hash
}

type coordinator struct {
	lggr logger.Logger

	lp logpoller.LogPoller
	topics
	finalityDepth uint32

	onchainRouter      VRFBeaconCoordinator
	coordinatorAddress common.Address
	beaconAddress      common.Address

	// We need to keep track of DKG ConfigSet events as well.
	dkgAddress common.Address

	evmClient evmclient.Client

	// set of blocks that have been scheduled for transmission.
	toBeTransmittedBlocks *ocrCache[blockInReport]
	// set of request id's that have been scheduled for transmission.
	toBeTransmittedCallbacks *ocrCache[callbackInReport]
	blockhashLookback        uint64
	coordinatorConfig        *ocr2vrftypes.CoordinatorConfig
	configDigest             ocr2Types.ConfigDigest
	oracleID                 commontypes.OracleID
}

// New creates a new CoordinatorInterface implementor.
func New(
	lggr logger.Logger,
	beaconAddress common.Address,
	coordinatorAddress common.Address,
	dkgAddress common.Address,
	client evmclient.Client,
	logPoller logpoller.LogPoller,
	finalityDepth uint32,
) (ocr2vrftypes.CoordinatorInterface, error) {
	onchainRouter, err := newRouter(lggr, beaconAddress, coordinatorAddress, client)
	if err != nil {
		return nil, errors.Wrap(err, "onchain router creation")
	}

	t := newTopics()

	// Add log filters for the log poller so that it can poll and find the logs that
	// we need.
	err = logPoller.RegisterFilter(logpoller.Filter{
		Name: filterName(beaconAddress, coordinatorAddress, dkgAddress),
		EventSigs: []common.Hash{
			t.randomnessRequestedTopic,
			t.randomnessFulfillmentRequestedTopic,
			t.randomWordsFulfilledTopic,
			t.configSetTopic,
			t.outputsServedTopic,
			t.newTransmissionTopic}, Addresses: []common.Address{beaconAddress, coordinatorAddress, dkgAddress}})
	if err != nil {
		return nil, err
	}

	cacheEvictionWindowSeconds := int64(60)
	cacheEvictionWindow := time.Duration(cacheEvictionWindowSeconds * int64(time.Second))
	lookbackBlocks := uint64(1_000)

	return &coordinator{
		onchainRouter:            onchainRouter,
		coordinatorAddress:       coordinatorAddress,
		beaconAddress:            beaconAddress,
		dkgAddress:               dkgAddress,
		lp:                       logPoller,
		topics:                   t,
		finalityDepth:            finalityDepth,
		evmClient:                client,
		lggr:                     lggr.Named("OCR2VRFCoordinator"),
		toBeTransmittedBlocks:    NewBlockCache[blockInReport](cacheEvictionWindow),
		toBeTransmittedCallbacks: NewBlockCache[callbackInReport](cacheEvictionWindow),
		blockhashLookback:        mathutil.Min(256, lookbackBlocks),
		// defaults
		coordinatorConfig: &ocr2vrftypes.CoordinatorConfig{
			CacheEvictionWindowSeconds: cacheEvictionWindowSeconds,
			BatchGasLimit:              5_000_000,
			CoordinatorOverhead:        50_000,
			BlockGasOverhead:           50_000,
			CallbackOverhead:           50_000,
			LookbackBlocks:             lookbackBlocks,
		},
	}, nil
}

func (c *coordinator) CurrentChainHeight(ctx context.Context) (uint64, error) {
	head, err := c.lp.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return 0, err
	}
	return uint64(head), nil
}

// ReportIsOnchain returns true iff a report for the given OCR epoch/round is
// present onchain.
func (c *coordinator) ReportIsOnchain(
	ctx context.Context,
	epoch uint32,
	round uint8,
	configDigest [32]byte,
) (presentOnchain bool, err error) {
	now := time.Now().UTC()
	defer c.logAndEmitFunctionDuration("ReportIsOnchain", now)

	// Check if a NewTransmission event was emitted on-chain with the
	// provided epoch and round.

	epochAndRound := toEpochAndRoundUint40(epoch, round)

	// this is technically NOT a hash in the regular meaning,
	// however it has the same size as a common.Hash. We need
	// to left-pad by bytes because it has to be 256 (or 32 bytes)
	// long in order to use as a topic filter.
	enrTopic := common.BytesToHash(common.LeftPadBytes(epochAndRound.Bytes(), 32))

	c.lggr.Info(fmt.Sprintf("epoch and round: %s %s", epochAndRound.String(), enrTopic.String()))
	logs, err := c.lp.IndexedLogs(
		c.topics.newTransmissionTopic,
		c.beaconAddress,
		2,
		[]common.Hash{
			enrTopic,
		},
		1,
		pg.WithParentCtx(ctx))
	if err != nil {
		return false, errors.Wrap(err, "log poller IndexedLogs")
	}

	// Filter for valid logs that match the current config digest.
	var logsWithCorrectConfigDigest []logpoller.Log
	for i := 0; i < len(logs); i++ {
		rawLog := toGethLog(logs[i])
		unpacked, err := c.onchainRouter.ParseLog(rawLog)
		if err != nil {
			c.lggr.Warnw("Incorrect log found in NewTransmissions", "log", logs[i], "err", err)
			continue
		}
		nt, ok := unpacked.(*vrf_beacon.VRFBeaconNewTransmission)
		if !ok {
			c.lggr.Warnw("Type error for log in NewTransmissisons", "log", logs[i], "err", err)
			continue
		}
		if nt.ConfigDigest == configDigest {
			logsWithCorrectConfigDigest = append(logsWithCorrectConfigDigest, logs[i])
		}
	}

	c.lggr.Info(fmt.Sprintf("NewTransmission logs: %+v", logsWithCorrectConfigDigest))

	return len(logsWithCorrectConfigDigest) >= 1, nil
}

// ReportBlocks returns the heights and hashes of the blocks which require VRF
// proofs in the current report, and the callback requests which should be
// served as part of processing that report. Everything returned by this
// should concern blocks older than the corresponding confirmationDelay.
// Blocks and callbacks it has returned previously may be returned again, as
// long as retransmissionDelay blocks have passed since they were last
// returned. The callbacks returned do not have to correspond to the blocks.
//
// The implementor is responsible for only returning well-funded callback
// requests, and blocks for which clients have actually requested random output
//
// This can be implemented on ethereum using the RandomnessRequested and
// RandomnessFulfillmentRequested events, to identify which blocks and
// callbacks need to be served, along with the NewTransmission and
// RandomWordsFulfilled events, to identify which have already been served.
func (c *coordinator) ReportBlocks(
	ctx context.Context,
	slotInterval uint16, // TODO: unused for now
	confirmationDelays map[uint32]struct{},
	retransmissionDelay time.Duration, // TODO: unused for now
	maxBlocks, // TODO: unused for now
	maxCallbacks int, // TODO: unused for now
) (
	blocks []ocr2vrftypes.Block,
	callbacks []ocr2vrftypes.AbstractCostedCallbackRequest,
	recentBlockHashesStartHeight uint64,
	recentBlockHashes []common.Hash,
	err error,
) {
	now := time.Now().UTC()
	defer c.logAndEmitFunctionDuration("ReportBlocks", now)

	// Instantiate the gas used by this batch.
	currentBatchGasLimit := c.coordinatorConfig.CoordinatorOverhead

	// TODO: use head broadcaster instead?
	currentHeight, err := c.CurrentChainHeight(ctx)
	if err != nil {
		err = errors.Wrap(err, "header by number")
		return
	}

	// Evict expired items from the cache.
	c.toBeTransmittedBlocks.EvictExpiredItems(now)
	c.toBeTransmittedCallbacks.EvictExpiredItems(now)

	c.lggr.Infow("current chain height", "currentHeight", currentHeight)

	logs, err := c.lp.LogsWithSigs(
		int64(currentHeight-c.coordinatorConfig.LookbackBlocks),
		int64(currentHeight),
		[]common.Hash{
			c.randomnessRequestedTopic,
			c.randomnessFulfillmentRequestedTopic,
			c.randomWordsFulfilledTopic,
			c.outputsServedTopic,
		},
		c.coordinatorAddress,
		pg.WithParentCtx(ctx))
	if err != nil {
		err = errors.Wrapf(err, "logs with topics. address: %s", c.coordinatorAddress)
		return
	}

	c.lggr.Tracew("logsWithSigs", "logs", logs)

	randomnessRequestedLogs,
		randomnessFulfillmentRequestedLogs,
		randomWordsFulfilledLogs,
		outputsServedLogs,
		err := c.unmarshalLogs(logs)
	if err != nil {
		err = errors.Wrap(err, "unmarshal logs")
		return
	}

	c.lggr.Tracew(
		"finished unmarshalLogs",
		"RandomnessRequested", randomnessRequestedLogs,
		"RandomnessFulfillmentRequested", randomnessFulfillmentRequestedLogs,
		"RandomWordsFulfilled", randomWordsFulfilledLogs,
		"OutputsServed", outputsServedLogs,
	)

	// Get start height for recent blockhashes.
	recentBlockHashesStartHeight = uint64(0)
	if currentHeight >= c.blockhashLookback {
		recentBlockHashesStartHeight = currentHeight - c.blockhashLookback + 1
	}

	// Get blockhashes that pertain to requested blocks.
	blockhashesMapping, err := c.getBlockhashesMappingFromRequests(
		ctx,
		randomnessRequestedLogs,
		randomnessFulfillmentRequestedLogs,
		currentHeight,
		recentBlockHashesStartHeight,
	)
	if err != nil {
		err = errors.Wrap(err, "get blockhashes in ReportBlocks")
		return
	}

	// TODO BELOW: Write tests for the new blockhash retrieval.
	// Obtain recent blockhashes, ordered by ascending block height.
	for i := recentBlockHashesStartHeight; i <= uint64(currentHeight); i++ {
		recentBlockHashes = append(recentBlockHashes, blockhashesMapping[i])
	}

	blocksRequested := make(map[block]struct{})
	unfulfilled, err := c.filterEligibleRandomnessRequests(randomnessRequestedLogs, confirmationDelays, currentHeight, blockhashesMapping)
	if err != nil {
		err = errors.Wrap(err, "filter requests in ReportBlocks")
		return
	}
	for _, uf := range unfulfilled {
		blocksRequested[uf] = struct{}{}
	}

	c.lggr.Tracew("filtered eligible randomness requests", "blocks", unfulfilled)

	callbacksRequested, unfulfilled, err := c.filterEligibleCallbacks(randomnessFulfillmentRequestedLogs, confirmationDelays, currentHeight, blockhashesMapping)
	if err != nil {
		err = errors.Wrap(err, "filter callbacks in ReportBlocks")
		return
	}
	for _, uf := range unfulfilled {
		blocksRequested[uf] = struct{}{}
	}

	c.lggr.Tracew("filtered eligible callbacks and blocks", "callbacks", callbacksRequested, "blocks", maps.Keys(blocksRequested))

	// Remove blocks that have already received responses so that we don't
	// respond to them again.
	fulfilledBlocks := c.getFulfilledBlocks(outputsServedLogs)
	for _, f := range fulfilledBlocks {
		delete(blocksRequested, f)
	}

	c.lggr.Tracew("got fulfilled blocks", "fulfilled", fulfilledBlocks)

	// Fill blocks slice with valid requested blocks.
	blocks = []ocr2vrftypes.Block{}
	for block := range blocksRequested {
		if c.coordinatorConfig.BatchGasLimit-currentBatchGasLimit >= c.coordinatorConfig.BlockGasOverhead {
			blocks = append(blocks, ocr2vrftypes.Block{
				Hash:              blockhashesMapping[block.blockNumber],
				Height:            block.blockNumber,
				ConfirmationDelay: block.confDelay,
			})
			currentBatchGasLimit += c.coordinatorConfig.BlockGasOverhead
		} else {
			break
		}
	}

	c.lggr.Tracew("got elligible blocks", "blocks", blocks)

	// Find unfulfilled callback requests by filtering out already fulfilled callbacks.
	fulfilledRequestIDs := c.getFulfilledRequestIDs(randomWordsFulfilledLogs)
	callbacks = c.filterUnfulfilledCallbacks(callbacksRequested, fulfilledRequestIDs, confirmationDelays, currentHeight, currentBatchGasLimit)
	c.emitReportBlocksMetrics(len(blocks), len(callbacks))

	// Pull request IDs from elligible callbacks for logging. There should only be
	// at most 100-200 elligible callbacks in a report.
	var reqIDs []uint64
	for _, c := range callbacks {
		reqIDs = append(reqIDs, c.RequestID)
	}
	c.lggr.Debugw("reporting blocks and callbacks", "blocks", blocks, "callbacks", reqIDs)
	c.lggr.Tracew("alreday fulfilled blocks and callbacks", "blocks", fulfilledBlocks, "callbacks", maps.Keys(fulfilledRequestIDs))

	return
}

// getBlockhashesMappingFromRequests returns the blockhashes for enqueued request blocks.
func (c *coordinator) getBlockhashesMappingFromRequests(
	ctx context.Context,
	randomnessRequestedLogs []*vrf_wrapper.VRFCoordinatorRandomnessRequested,
	randomnessFulfillmentRequestedLogs []*vrf_wrapper.VRFCoordinatorRandomnessFulfillmentRequested,
	currentHeight uint64,
	recentBlockHashesStartHeight uint64,
) (blockhashesMapping map[uint64]common.Hash, err error) {

	// Get all request + callback requests into a mapping.
	rawBlocksRequested := make(map[uint64]struct{})
	for _, l := range randomnessRequestedLogs {
		if isBlockEligible(l.NextBeaconOutputHeight, l.ConfDelay, currentHeight) {
			rawBlocksRequested[l.NextBeaconOutputHeight] = struct{}{}

			// Also get the blockhash for the most recent cached report on this block,
			// if one exists.
			cacheKey := getBlockCacheKey(l.NextBeaconOutputHeight, l.ConfDelay.Uint64())
			t := c.toBeTransmittedBlocks.GetItem(cacheKey)
			if t != nil {
				rawBlocksRequested[t.recentBlockHeight] = struct{}{}
			}
		}
	}
	for _, l := range randomnessFulfillmentRequestedLogs {
		if isBlockEligible(l.NextBeaconOutputHeight, l.ConfDelay, currentHeight) {
			rawBlocksRequested[l.NextBeaconOutputHeight] = struct{}{}

			// Also get the blockhash for the most recent cached report on this callback,
			// if one exists.
			cacheKey := getCallbackCacheKey(l.RequestID.Int64())
			t := c.toBeTransmittedCallbacks.GetItem(cacheKey)
			if t != nil {
				rawBlocksRequested[t.recentBlockHeight] = struct{}{}
			}
		}
	}

	// Fill a unique list of request blocks.
	requestedBlockNumbers := []uint64{}
	for k := range rawBlocksRequested {
		requestedBlockNumbers = append(requestedBlockNumbers, k)
	}

	// Get a mapping of block numbers to block hashes.
	blockhashesMapping, err = c.getBlockhashesMapping(ctx, append(requestedBlockNumbers, uint64(currentHeight), recentBlockHashesStartHeight))
	if err != nil {
		err = errors.Wrap(err, "get blockhashes for ReportBlocks")
	}
	return
}

func (c *coordinator) getFulfilledBlocks(outputsServedLogs []*vrf_coordinator.VRFCoordinatorOutputsServed) (fulfilled []block) {
	for _, r := range outputsServedLogs {
		for _, o := range r.OutputsServed {
			fulfilled = append(fulfilled, block{
				blockNumber: o.Height,
				confDelay:   uint32(o.ConfirmationDelay.Uint64()),
			})
		}
	}
	return
}

// getBlockhashesMapping returns the blockhashes corresponding to a slice of block numbers.
func (c *coordinator) getBlockhashesMapping(
	ctx context.Context,
	blockNumbers []uint64,
) (blockhashesMapping map[uint64]common.Hash, err error) {
	// GetBlocks doesn't necessarily need a sorted blockNumbers array,
	// but sorting it is helpful for testing.
	sort.Slice(blockNumbers, func(a, b int) bool {
		return blockNumbers[a] < blockNumbers[b]
	})

	heads, err := c.lp.GetBlocksRange(ctx, blockNumbers, pg.WithParentCtx(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "logpoller.GetBlocks")
	}

	blockhashesMapping = make(map[uint64]common.Hash)
	for _, head := range heads {
		blockhashesMapping[uint64(head.BlockNumber)] = head.BlockHash
	}

	// Ensure that every requested block received a blockhash.
	for _, b := range blockNumbers {
		if _, ok := blockhashesMapping[b]; !ok {
			err = fmt.Errorf("could not find all heads in db: want %d got %d", len(blockNumbers), len(heads))
			return
		}
	}
	return
}

// getFulfilledRequestIDs returns the request IDs referenced by the given RandomWordsFulfilled logs slice.
func (c *coordinator) getFulfilledRequestIDs(randomWordsFulfilledLogs []*vrf_wrapper.VRFCoordinatorRandomWordsFulfilled) map[uint64]struct{} {
	fulfilledRequestIDs := make(map[uint64]struct{})
	for _, r := range randomWordsFulfilledLogs {
		for _, requestID := range r.RequestIDs {
			fulfilledRequestIDs[requestID.Uint64()] = struct{}{}
		}
	}
	return fulfilledRequestIDs
}

// filterUnfulfilledCallbacks returns unfulfilled callback requests given the
// callback request logs and the already fulfilled callback request IDs.
func (c *coordinator) filterUnfulfilledCallbacks(
	callbacksRequested []*vrf_wrapper.VRFCoordinatorRandomnessFulfillmentRequested,
	fulfilledRequestIDs map[uint64]struct{},
	confirmationDelays map[uint32]struct{},
	currentHeight uint64,
	currentBatchGasLimit int64,
) (callbacks []ocr2vrftypes.AbstractCostedCallbackRequest) {

	/**
	 * Callback batch ordering:
	 * - Callbacks are first ordered by beacon output + confirmation delay (ascending), in other words
	 *   the fulfillments at the oldest block are first in line.
	 * - Within the same block, fulfillments are ordered by gasAllowance (ascending), i.e the callbacks with
	 *   the lowest gasAllowance are first in line.
	 * - This ordering ensures that the oldest callbacks can be picked up first, and that as many callbacks as
	 *   possible can be fit into a batch.
	 *
	 * Example:
	 * Unsorted: (outputHeight: 1, gasAllowance: 200k), (outputHeight: 3, gasAllowance: 100k), (outputHeight: 1, gasAllowance: 100k)
	 * Sorted: (outputHeight: 1, gasAllowance: 100k), (outputHeight: 1, gasAllowance: 200k), (outputHeight: 3, gasAllowance: 100k)
	 *
	 */
	sort.Slice(callbacksRequested, func(a, b int) bool {
		aHeight := callbacksRequested[a].NextBeaconOutputHeight + callbacksRequested[a].ConfDelay.Uint64()
		bHeight := callbacksRequested[b].NextBeaconOutputHeight + callbacksRequested[b].ConfDelay.Uint64()
		if aHeight == bHeight {
			return callbacksRequested[a].GasAllowance < callbacksRequested[b].GasAllowance
		}
		return aHeight < bHeight
	})

	for _, r := range callbacksRequested {
		// Check if there is room left in the batch. If there is no room left, the coordinator
		// will keep iterating, until it either finds a callback in a subsequent output height that
		// can fit into the current batch or reaches the end of the sorted callbacks slice.
		if c.coordinatorConfig.BatchGasLimit-currentBatchGasLimit < (int64(r.GasAllowance) + c.coordinatorConfig.CallbackOverhead) {
			continue
		}

		requestID := r.RequestID
		if _, ok := fulfilledRequestIDs[requestID.Uint64()]; !ok {
			// The on-chain machinery will revert requests that specify an unsupported
			// confirmation delay, so this is more of a sanity check than anything else.
			if _, ok := confirmationDelays[uint32(r.ConfDelay.Uint64())]; !ok {
				// if we can't find the conf delay in the map then just ignore this request
				c.lggr.Errorw("ignoring bad request, unsupported conf delay",
					"confDelay", r.ConfDelay.String(),
					"supportedConfDelays", confirmationDelays)
				continue
			}

			// NOTE: we already check if the callback has been fulfilled in filterEligibleCallbacks,
			// so we don't need to do that again here.
			if isBlockEligible(r.NextBeaconOutputHeight, r.ConfDelay, currentHeight) {
				callbacks = append(callbacks, ocr2vrftypes.AbstractCostedCallbackRequest{
					BeaconHeight:      r.NextBeaconOutputHeight,
					ConfirmationDelay: uint32(r.ConfDelay.Uint64()),
					SubscriptionID:    r.SubID,
					Price:             big.NewInt(0), // TODO: no price tracking
					RequestID:         requestID.Uint64(),
					NumWords:          r.NumWords,
					Requester:         r.Requester,
					Arguments:         r.Arguments,
					GasAllowance:      big.NewInt(int64(r.GasAllowance)),
					GasPrice:          r.GasPrice,
					WeiPerUnitLink:    r.WeiPerUnitLink,
				})
				currentBatchGasLimit += int64(r.GasAllowance)
				c.lggr.Debugw("Request is unfulfilled", "requestID", requestID)
			}
		}
	}
	return callbacks
}

// filterEligibleCallbacks extracts valid callback requests from the given logs,
// based on their readiness to be fulfilled. It also returns any unfulfilled blocks
// associated with those callbacks.
func (c *coordinator) filterEligibleCallbacks(
	randomnessFulfillmentRequestedLogs []*vrf_wrapper.VRFCoordinatorRandomnessFulfillmentRequested,
	confirmationDelays map[uint32]struct{},
	currentHeight uint64,
	blockhashesMapping map[uint64]common.Hash,
) (callbacks []*vrf_wrapper.VRFCoordinatorRandomnessFulfillmentRequested, unfulfilled []block, err error) {

	for _, r := range randomnessFulfillmentRequestedLogs {
		// The on-chain machinery will revert requests that specify an unsupported
		// confirmation delay, so this is more of a sanity check than anything else.
		if _, ok := confirmationDelays[uint32(r.ConfDelay.Uint64())]; !ok {
			// if we can't find the conf delay in the map then just ignore this request
			c.lggr.Errorw("ignoring bad request, unsupported conf delay",
				"confDelay", r.ConfDelay.String(),
				"supportedConfDelays", confirmationDelays)
			continue
		}

		// Check that the callback is elligible.
		if isBlockEligible(r.NextBeaconOutputHeight, r.ConfDelay, currentHeight) {
			cacheKey := getCallbackCacheKey(r.RequestID.Int64())
			t := c.toBeTransmittedCallbacks.GetItem(cacheKey)
			// If the callback is found in the cache and the recentBlockHash from the report containing the callback
			// is correct, then the callback is in-flight and should not be included in the current observation. If that
			// report gets re-orged, then the recentBlockHash of the report will become invalid, in which case
			// the cached callback is ignored, and the callback is added to the current observation.
			inflightTransmission := (t != nil) && (t.recentBlockHash == blockhashesMapping[t.recentBlockHeight])
			if inflightTransmission {
				c.lggr.Debugw("Request is in-flight", "requestID", r.RequestID)
				continue
			}

			callbacks = append(callbacks, r)

			// We could have a callback request that was made in a different block than what we
			// have possibly already received from regular requests.
			unfulfilled = append(unfulfilled, block{
				blockNumber: r.NextBeaconOutputHeight,
				confDelay:   uint32(r.ConfDelay.Uint64()),
			})
			c.lggr.Debugw("Request is eligible", "requestID", r.RequestID)
		}
	}
	return
}

// filterEligibleRandomnessRequests extracts valid randomness requests from the given logs,
// based on their readiness to be fulfilled.
func (c *coordinator) filterEligibleRandomnessRequests(
	randomnessRequestedLogs []*vrf_wrapper.VRFCoordinatorRandomnessRequested,
	confirmationDelays map[uint32]struct{},
	currentHeight uint64,
	blockhashesMapping map[uint64]common.Hash,
) (unfulfilled []block, err error) {

	for _, r := range randomnessRequestedLogs {
		// The on-chain machinery will revert requests that specify an unsupported
		// confirmation delay, so this is more of a sanity check than anything else.
		if _, ok := confirmationDelays[uint32(r.ConfDelay.Uint64())]; !ok {
			// if we can't find the conf delay in the map then just ignore this request
			c.lggr.Errorw("ignoring bad request, unsupported conf delay",
				"confDelay", r.ConfDelay.String(),
				"supportedConfDelays", confirmationDelays)
			continue
		}

		// Check that the block is elligible.
		if isBlockEligible(r.NextBeaconOutputHeight, r.ConfDelay, currentHeight) {
			cacheKey := getBlockCacheKey(r.NextBeaconOutputHeight, r.ConfDelay.Uint64())
			t := c.toBeTransmittedBlocks.GetItem(cacheKey)
			// If the block is found in the cache and the recentBlockHash from the report containing the block
			// is correct, then the block is in-flight and should not be included in the current observation. If that
			// report gets re-orged, then the recentBlockHash of the report will become invalid, in which case
			// the cached block is ignored and the block is added to the current observation.
			validTransmission := (t != nil) && (t.recentBlockHash == blockhashesMapping[t.recentBlockHeight])
			if validTransmission {
				c.lggr.Debugw("Block is in-flight", "blockNumber", r.NextBeaconOutputHeight, "confDelay", r.ConfDelay)
				continue
			}

			unfulfilled = append(unfulfilled, block{
				blockNumber: r.NextBeaconOutputHeight,
				confDelay:   uint32(r.ConfDelay.Uint64()),
			})
			c.lggr.Debugw("Block is eligible", "blockNumber", r.NextBeaconOutputHeight, "confDelay", r.ConfDelay)
		}
	}
	return
}

func (c *coordinator) unmarshalLogs(
	logs []logpoller.Log,
) (
	randomnessRequestedLogs []*vrf_wrapper.VRFCoordinatorRandomnessRequested,
	randomnessFulfillmentRequestedLogs []*vrf_wrapper.VRFCoordinatorRandomnessFulfillmentRequested,
	randomWordsFulfilledLogs []*vrf_wrapper.VRFCoordinatorRandomWordsFulfilled,
	outputsServedLogs []*vrf_wrapper.VRFCoordinatorOutputsServed,
	err error,
) {
	for _, lg := range logs {
		rawLog := toGethLog(lg)
		switch lg.EventSig {
		case c.randomnessRequestedTopic:
			unpacked, err2 := c.onchainRouter.ParseLog(rawLog)
			if err2 != nil {
				// should never happen
				err = errors.Wrap(err2, "unmarshal RandomnessRequested log")
				return
			}
			rr, ok := unpacked.(*vrf_wrapper.VRFCoordinatorRandomnessRequested)
			if !ok {
				// should never happen
				err = errors.New("cast to *VRFCoordinatorRandomnessRequested")
				return
			}
			randomnessRequestedLogs = append(randomnessRequestedLogs, rr)
		case c.randomnessFulfillmentRequestedTopic:
			unpacked, err2 := c.onchainRouter.ParseLog(rawLog)
			if err2 != nil {
				// should never happen
				err = errors.Wrap(err2, "unmarshal RandomnessFulfillmentRequested log")
				return
			}
			rfr, ok := unpacked.(*vrf_wrapper.VRFCoordinatorRandomnessFulfillmentRequested)
			if !ok {
				// should never happen
				err = errors.New("cast to *VRFCoordinatorRandomnessFulfillmentRequested")
				return
			}
			randomnessFulfillmentRequestedLogs = append(randomnessFulfillmentRequestedLogs, rfr)
		case c.randomWordsFulfilledTopic:
			unpacked, err2 := c.onchainRouter.ParseLog(rawLog)
			if err2 != nil {
				// should never happen
				err = errors.Wrap(err2, "unmarshal RandomWordsFulfilled log")
				return
			}
			rwf, ok := unpacked.(*vrf_wrapper.VRFCoordinatorRandomWordsFulfilled)
			if !ok {
				// should never happen
				err = errors.New("cast to *VRFCoordinatorRandomWordsFulfilled")
				return
			}
			randomWordsFulfilledLogs = append(randomWordsFulfilledLogs, rwf)
		case c.outputsServedTopic:
			unpacked, err2 := c.onchainRouter.ParseLog(rawLog)
			if err2 != nil {
				// should never happen
				err = errors.Wrap(err2, "unmarshal OutputsServed log")
				return
			}
			nt, ok := unpacked.(*vrf_coordinator.VRFCoordinatorOutputsServed)
			if !ok {
				// should never happen
				err = errors.New("cast to *vrf_coordinator.VRFCoordinatorOutputsServed")
			}
			outputsServedLogs = append(outputsServedLogs, nt)
		default:
			c.lggr.Error(fmt.Sprintf("Unexpected event sig: %s", lg.EventSig))
			c.lggr.Error(fmt.Sprintf("expected one of: %s (RandomnessRequested) %s (RandomnessFulfillmentRequested) %s (RandomWordsFulfilled) %s (OutputsServed), got %s",
				hexutil.Encode(c.randomnessRequestedTopic[:]),
				hexutil.Encode(c.randomnessFulfillmentRequestedTopic[:]),
				hexutil.Encode(c.randomWordsFulfilledTopic[:]),
				hexutil.Encode(c.outputsServedTopic[:]),
				lg.EventSig))
		}
	}
	return
}

// ReportWillBeTransmitted registers to the CoordinatorInterface that the
// local node has accepted the AbstractReport for transmission, so that its
// blocks and callbacks can be tracked for possible later retransmission
func (c *coordinator) ReportWillBeTransmitted(ctx context.Context, report ocr2vrftypes.AbstractReport) error {
	now := time.Now().UTC()
	defer c.logAndEmitFunctionDuration("ReportWillBeTransmitted", now)

	// Evict expired items from the cache.
	c.toBeTransmittedBlocks.EvictExpiredItems(now)
	c.toBeTransmittedCallbacks.EvictExpiredItems(now)

	// Check for a re-org, and return an error if one is present.
	blockhashesMapping, err := c.getBlockhashesMapping(ctx, []uint64{report.RecentBlockHeight})
	if err != nil {
		return errors.Wrap(err, "getting blockhash mapping in ReportWillBeTransmitted")
	}
	if blockhashesMapping[report.RecentBlockHeight] != report.RecentBlockHash {
		return errors.Errorf("blockhash of report does not match most recent blockhash in ReportWillBeTransmitted")
	}

	blocksRequested := []blockInReport{}
	callbacksRequested := []callbackInReport{}

	// Get all requested blocks and callbacks.
	for _, output := range report.Outputs {
		// If the VRF proof size is 0, the block is not included in this output. We still
		// check for callbacks in the ouptut.
		if len(output.VRFProof) > 0 {
			bR := blockInReport{
				block: block{
					blockNumber: output.BlockHeight,
					confDelay:   output.ConfirmationDelay,
				},
				recentBlockHeight: report.RecentBlockHeight,
				recentBlockHash:   report.RecentBlockHash,
			}
			// Store block in blocksRequested.br
			blocksRequested = append(blocksRequested, bR)
		}

		// Iterate through callbacks for output.
		for _, cb := range output.Callbacks {
			cbR := callbackInReport{
				callback: callback{
					blockNumber: cb.BeaconHeight,
					requestID:   cb.RequestID,
				},
				recentBlockHeight: report.RecentBlockHeight,
				recentBlockHash:   report.RecentBlockHash,
			}

			// Add callback to callbacksRequested.
			callbacksRequested = append(callbacksRequested, cbR)
		}
	}

	// Apply blockhashes to blocks and mark them as transmitted.
	for _, b := range blocksRequested {
		cacheKey := getBlockCacheKey(b.blockNumber, uint64(b.confDelay))
		c.toBeTransmittedBlocks.CacheItem(b, cacheKey, now)
		c.lggr.Debugw("Block is being transmitted", "blockNumber", b.blockNumber, "confDelay", b.confDelay)
	}

	// Add the corresponding blockhashes to callbacks and mark them as transmitted.
	for _, cb := range callbacksRequested {
		cacheKey := getCallbackCacheKey(int64(cb.requestID))
		c.toBeTransmittedCallbacks.CacheItem(cb, cacheKey, now)
		c.lggr.Debugw("Request is being transmitted", "requestID", cb.requestID)
	}

	c.emitReportWillBeTransmittedMetrics(len(blocksRequested), len(callbacksRequested))

	return nil
}

// DKGVRFCommittees returns the addresses of the signers and transmitters
// for the DKG and VRF OCR committees. On ethereum, these can be retrieved
// from the most recent ConfigSet events for each contract.
func (c *coordinator) DKGVRFCommittees(ctx context.Context) (dkgCommittee, vrfCommittee ocr2vrftypes.OCRCommittee, err error) {
	startTime := time.Now().UTC()
	defer c.logAndEmitFunctionDuration("DKGVRFCommittees", startTime)

	latestVRF, err := c.lp.LatestLogByEventSigWithConfs(
		c.configSetTopic,
		c.beaconAddress,
		int(c.finalityDepth),
	)
	if err != nil {
		err = errors.Wrap(err, "latest vrf ConfigSet by sig with confs")
		return
	}

	latestDKG, err := c.lp.LatestLogByEventSigWithConfs(
		c.configSetTopic,
		c.dkgAddress,
		int(c.finalityDepth),
	)
	if err != nil {
		err = errors.Wrap(err, "latest dkg ConfigSet by sig with confs")
		return
	}

	var vrfConfigSetLog vrf_beacon.VRFBeaconConfigSet
	err = vrfBeaconABI.UnpackIntoInterface(&vrfConfigSetLog, configSetEvent, latestVRF.Data)
	if err != nil {
		err = errors.Wrap(err, "unpack vrf ConfigSet into interface")
		return
	}

	var dkgConfigSetLog dkg_wrapper.DKGConfigSet
	err = dkgABI.UnpackIntoInterface(&dkgConfigSetLog, configSetEvent, latestDKG.Data)
	if err != nil {
		err = errors.Wrap(err, "unpack dkg ConfigSet into interface")
		return
	}

	// len(signers) == len(transmitters), this is guaranteed by libocr.
	for i := range vrfConfigSetLog.Signers {
		vrfCommittee.Signers = append(vrfCommittee.Signers, vrfConfigSetLog.Signers[i])
		vrfCommittee.Transmitters = append(vrfCommittee.Transmitters, vrfConfigSetLog.Transmitters[i])
	}

	for i := range dkgConfigSetLog.Signers {
		dkgCommittee.Signers = append(dkgCommittee.Signers, dkgConfigSetLog.Signers[i])
		dkgCommittee.Transmitters = append(dkgCommittee.Transmitters, dkgConfigSetLog.Transmitters[i])
	}

	return
}

// ProvingKeyHash returns the VRF current proving block, in view of the local
// node. On ethereum this can be retrieved from the VRF contract's attribute
// s_provingKeyHash
func (c *coordinator) ProvingKeyHash(ctx context.Context) (common.Hash, error) {
	h, err := c.onchainRouter.SProvingKeyHash(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "get proving block hash")
	}

	return h, nil
}

// BeaconPeriod returns the period used in the coordinator's contract
func (c *coordinator) BeaconPeriod(ctx context.Context) (uint16, error) {
	beaconPeriodBlocks, err := c.onchainRouter.IBeaconPeriodBlocks(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return 0, errors.Wrap(err, "get beacon period blocks")
	}

	return uint16(beaconPeriodBlocks.Int64()), nil
}

// ConfirmationDelays returns the list of confirmation delays defined in the coordinator's contract
func (c *coordinator) ConfirmationDelays(ctx context.Context) ([]uint32, error) {
	confDelays, err := c.onchainRouter.GetConfirmationDelays(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not get confirmation delays")
	}
	var result []uint32
	for _, c := range confDelays {
		result = append(result, uint32(c.Uint64()))
	}
	return result, nil
}

// KeyID returns the key ID from coordinator's contract
func (c *coordinator) KeyID(ctx context.Context) (dkg.KeyID, error) {
	keyID, err := c.onchainRouter.SKeyID(&bind.CallOpts{Context: ctx})
	if err != nil {
		return dkg.KeyID{}, errors.Wrap(err, "could not get key ID")
	}
	return keyID, nil
}

// isBlockEligible returns true if and only if the nextBeaconOutputHeight plus
// the confDelay is less than the current blockchain height, meaning that the beacon
// output height has enough confirmations.
//
// NextBeaconOutputHeight is always greater than the request block, therefore
// a number of confirmations on the beacon block is always enough confirmations
// for the request block.
func isBlockEligible(nextBeaconOutputHeight uint64, confDelay *big.Int, currentHeight uint64) bool {
	cond := confDelay.Uint64() < currentHeight // Edge case: for simulated chains with low block numbers
	cond = cond && (nextBeaconOutputHeight+confDelay.Uint64()) < currentHeight
	return cond
}

// toEpochAndRoundUint40 returns a single unsigned 40 bit big.Int object
// that has the epoch in the first 32 bytes and the round in the last 8 bytes,
// in a big-endian fashion.
func toEpochAndRoundUint40(epoch uint32, round uint8) *big.Int {
	return big.NewInt((int64(epoch) << 8) + int64(round))
}

func toGethLog(lg logpoller.Log) types.Log {
	var topics []common.Hash
	for _, b := range lg.Topics {
		topics = append(topics, common.BytesToHash(b))
	}
	return types.Log{
		Data:        lg.Data,
		Address:     lg.Address,
		BlockHash:   lg.BlockHash,
		BlockNumber: uint64(lg.BlockNumber),
		Topics:      topics,
		TxHash:      lg.TxHash,
		Index:       uint(lg.LogIndex),
	}
}

// getBlockCacheKey returns a cache key for a requested block
// The blockhash of the block does not need to be included in the key. Instead,
// the block cached at a given key contains a blockhash that is checked for validity
// against the log poller's current state.
func getBlockCacheKey(blockNumber uint64, confDelay uint64) common.Hash {
	var blockNumberBytes [8]byte
	var confDelayBytes [8]byte

	binary.BigEndian.PutUint64(blockNumberBytes[:], blockNumber)
	binary.BigEndian.PutUint64(confDelayBytes[:], confDelay)

	return common.BytesToHash(bytes.Join([][]byte{blockNumberBytes[:], confDelayBytes[:]}, nil))
}

// getBlockCacheKey returns a cache key for a requested callback
// The blockhash of the callback does not need to be included in the key. Instead,
// the callback cached at a given key contains a blockhash that is checked for validity
// against the log poller's current state.
func getCallbackCacheKey(requestID int64) common.Hash {
	return common.BigToHash(big.NewInt(requestID))
}

// logAndEmitFunctionDuration logs the time in milliseconds and emits metrics in nanosecond for function duration
func (c *coordinator) logAndEmitFunctionDuration(funcName string, startTime time.Time) {
	elapsed := time.Now().UTC().Sub(startTime)
	c.lggr.Debugf("%s took %d milliseconds to complete", funcName, elapsed.Milliseconds())
	promMethodDuration.WithLabelValues(
		append(c.labelValues(), funcName)...,
	).Observe(float64(elapsed.Nanoseconds()))
}

func (c *coordinator) UpdateConfiguration(
	b []byte,
	configDigest ocr2Types.ConfigDigest,
	oracleID commontypes.OracleID,
) error {
	// Update config digest & oracle ID for epoch.
	c.configDigest = configDigest
	c.oracleID = oracleID

	// Unmarshal off-chain config.
	err := proto.Unmarshal(b, c.coordinatorConfig)
	if err != nil {
		return errors.Wrap(err, "error setting offchain config on coordinator")
	}

	// Update local caches with new eviction window.
	cacheEvictionWindowSeconds := c.coordinatorConfig.CacheEvictionWindowSeconds
	cacheEvictionWindow := time.Duration(cacheEvictionWindowSeconds * int64(time.Second))
	c.toBeTransmittedBlocks.SetEvictonWindow(cacheEvictionWindow)
	c.toBeTransmittedCallbacks.SetEvictonWindow(cacheEvictionWindow)

	c.blockhashLookback = mathutil.Min(256, c.coordinatorConfig.LookbackBlocks)
	c.lggr.Infow("set offchain config",
		offchainConfigFields(c.coordinatorConfig)...,
	)

	return nil
}

func offchainConfigFields(coordinatorConfig *ocr2vrftypes.CoordinatorConfig) []any {
	return []any{
		"cacheEvictionWindowSeconds", coordinatorConfig.CacheEvictionWindowSeconds,
		"batchGasLimit", coordinatorConfig.BatchGasLimit,
		"coordinatorOverhead", coordinatorConfig.CoordinatorOverhead,
		"lookbackBlocks", coordinatorConfig.LookbackBlocks,
		"blockGasOverhead", coordinatorConfig.BlockGasOverhead,
		"callbackOverhead", coordinatorConfig.CallbackOverhead,
	}
}

func (c *coordinator) labelValues() []string {
	chainId := c.evmClient.ConfiguredChainID()
	return []string{chainId.String(), fmt.Sprintf("%d", c.oracleID), common.Bytes2Hex(c.configDigest[:])}
}

func (c *coordinator) emitReportBlocksMetrics(
	numBlocks int,
	numCallbacks int) {
	promBlocksToReport.WithLabelValues(c.labelValues()...).Observe(float64(numBlocks))
	promCallbacksToReport.WithLabelValues(c.labelValues()...).Observe(float64(numCallbacks))
}

func (c *coordinator) emitReportWillBeTransmittedMetrics(
	numBlocks int,
	numCallbacks int) {
	promBlocksInReport.WithLabelValues(c.labelValues()...).Observe(float64(numBlocks))
	promCallbacksInReport.WithLabelValues(c.labelValues()...).Observe(float64(numCallbacks))
}

func filterName(beaconAddress, coordinatorAddress, dkgAddress common.Address) string {
	return logpoller.FilterName("VRF Coordinator", beaconAddress, coordinatorAddress, dkgAddress)
}

func FilterNamesFromSpec(spec *job.OCR2OracleSpec) (names []string, err error) {
	var cfg ocr2vrfconfig.PluginConfig
	var beaconAddress, coordinatorAddress, dkgAddress ethkey.EIP55Address

	if err = json.Unmarshal(spec.PluginConfig.Bytes(), &cfg); err != nil {
		err = errors.Wrap(err, "failed to unmarshal ocr2vrf plugin config")
		return nil, err
	}

	if beaconAddress, err = ethkey.NewEIP55Address(spec.ContractID); err == nil {
		if coordinatorAddress, err = ethkey.NewEIP55Address(cfg.VRFCoordinatorAddress); err == nil {
			if dkgAddress, err = ethkey.NewEIP55Address(cfg.DKGContractAddress); err == nil {
				return []string{filterName(beaconAddress.Address(), coordinatorAddress.Address(), dkgAddress.Address())}, nil
			}
		}
	}

	return nil, err
}
