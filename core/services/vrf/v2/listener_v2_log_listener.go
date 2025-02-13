package v2

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/mathutil"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
)

func (lsn *listenerV2) runLogListener(
	pollPeriod time.Duration,
	minConfs uint32,
) {
	lsn.l.Infow("Listening for run requests via log poller",
		"minConfs", minConfs)
	ticker := time.NewTicker(pollPeriod)
	defer ticker.Stop()
	var (
		lastProcessedBlock int64
		startingUp         = true
	)
	filterName := lsn.getLogPollerFilterName()
	ctx, cancel := lsn.chStop.NewCtx()
	defer cancel()
	for {
		select {
		case <-lsn.chStop:
			return
		case <-ticker.C:
			start := time.Now()
			lsn.l.Debugw("log listener loop")

			// If filter has not already been successfully registered, register it.
			if !lsn.chain.LogPoller().HasFilter(filterName) {
				err := lsn.chain.LogPoller().RegisterFilter(ctx, logpoller.Filter{
					Name: filterName,
					EventSigs: evmtypes.HashArray{
						lsn.coordinator.RandomWordsFulfilledTopic(),
						lsn.coordinator.RandomWordsRequestedTopic(),
					},
					Addresses: evmtypes.AddressArray{
						lsn.coordinator.Address(),
					},
				})
				if err != nil {
					lsn.l.Errorw("error registering filter in log poller, retrying",
						"err", err,
						"elapsed", time.Since(start))
					continue
				}
			}

			// on startup we want to initialize the last processed block
			if startingUp {
				var err error
				lsn.l.Debugw("initializing last processed block on startup")
				lastProcessedBlock, err = lsn.initializeLastProcessedBlock(ctx)
				if err != nil {
					lsn.l.Errorw("error initializing last processed block, retrying",
						"err", err,
						"elapsed", time.Since(start))
					continue
				}
				startingUp = false
				lsn.l.Debugw("initialized last processed block", "lastProcessedBlock", lastProcessedBlock)
			}

			pending, err := lsn.pollLogs(ctx, minConfs, lastProcessedBlock)
			if err != nil {
				lsn.l.Errorw("error polling vrf logs, retrying",
					"err", err,
					"elapsed", time.Since(start))
				continue
			}

			// process pending requests and insert any fulfillments into the inflight cache
			lsn.processPendingVRFRequests(ctx, pending)

			lastProcessedBlock, err = lsn.updateLastProcessedBlock(ctx, lastProcessedBlock)
			if err != nil {
				lsn.l.Errorw("error updating last processed block, continuing anyway", "err", err)
			} else {
				lsn.l.Debugw("updated last processed block", "lastProcessedBlock", lastProcessedBlock)
			}
			lsn.l.Debugw("log listener loop done", "elapsed", time.Since(start))
		}
	}
}

func (lsn *listenerV2) getLogPollerFilterName() string {
	return logpoller.FilterName(
		"VRFListener",
		"version", lsn.coordinator.Version(),
		"keyhash", lsn.job.VRFSpec.PublicKey.MustHash(),
		"coordinatorAddress", lsn.coordinator.Address())
}

// initializeLastProcessedBlock returns the earliest block number that we need to
// process requests for. This is the block number of the earliest unfulfilled request
// or the latest finalized block, if there are no unfulfilled requests.
// TODO: add tests
func (lsn *listenerV2) initializeLastProcessedBlock(ctx context.Context) (lastProcessedBlock int64, err error) {
	lp := lsn.chain.LogPoller()
	start := time.Now()

	// will retry on error in the runLogListener loop
	latestBlock, err := lp.LatestBlock(ctx)
	if err != nil {
		return 0, fmt.Errorf("LogPoller.LatestBlock(): %w", err)
	}
	fromTimestamp := time.Now().UTC().Add(-lsn.job.VRFSpec.RequestTimeout)
	ll := lsn.l.With(
		"latestFinalizedBlock", latestBlock.FinalizedBlockNumber,
		"latestBlock", latestBlock.BlockNumber,
		"fromTimestamp", fromTimestamp)
	ll.Debugw("Initializing last processed block")
	defer func() {
		ll.Debugw("Done initializing last processed block", "elapsed", time.Since(start))
	}()

	numBlocksToReplay := numReplayBlocks(lsn.job.VRFSpec.RequestTimeout, lsn.chain.ID())
	replayStartBlock := mathutil.Max(latestBlock.FinalizedBlockNumber-numBlocksToReplay, 1)
	ll.Debugw("running replay on log poller",
		"numBlocksToReplay", numBlocksToReplay,
		"replayStartBlock", replayStartBlock,
		"requestTimeout", lsn.job.VRFSpec.RequestTimeout,
	)
	err = lp.Replay(ctx, replayStartBlock)
	if err != nil {
		return 0, fmt.Errorf("LogPoller.Replay: %w", err)
	}

	// get randomness requested logs with the appropriate keyhash
	// keyhash is specified in topic1
	requests, err := lp.IndexedLogsCreatedAfter(
		ctx,
		lsn.coordinator.RandomWordsRequestedTopic(), // event sig
		lsn.coordinator.Address(),                   // address
		1,                                           // topic index
		[]common.Hash{lsn.job.VRFSpec.PublicKey.MustHash()}, // topic values
		fromTimestamp,      // from time
		evmtypes.Finalized, // confs
	)
	if err != nil {
		return 0, fmt.Errorf("LogPoller.LogsCreatedAfter RandomWordsRequested logs: %w", err)
	}

	// fulfillments don't have keyhash indexed, we'll have to get all of them
	// TODO: can we instead write a single query that joins on request id's somehow?
	fulfillments, err := lp.LogsCreatedAfter(
		ctx,
		lsn.coordinator.RandomWordsFulfilledTopic(), // event sig
		lsn.coordinator.Address(),                   // address
		fromTimestamp,                               // from time
		evmtypes.Finalized,                          // confs
	)
	if err != nil {
		return 0, fmt.Errorf("LogPoller.LogsCreatedAfter RandomWordsFulfilled logs: %w", err)
	}

	unfulfilled, _, _ := lsn.getUnfulfilled(append(requests, fulfillments...), ll)
	// find request block of earliest unfulfilled request
	// even if this block is > latest finalized, we use latest finalized as earliest unprocessed
	// because re-orgs can occur on any unfinalized block.
	var earliestUnfulfilledBlock = latestBlock.FinalizedBlockNumber
	for _, req := range unfulfilled {
		if req.Raw().BlockNumber < uint64(earliestUnfulfilledBlock) {
			earliestUnfulfilledBlock = int64(req.Raw().BlockNumber)
		}
	}

	return earliestUnfulfilledBlock, nil
}

func (lsn *listenerV2) updateLastProcessedBlock(ctx context.Context, currLastProcessedBlock int64) (lastProcessedBlock int64, err error) {
	lp := lsn.chain.LogPoller()
	start := time.Now()

	latestBlock, err := lp.LatestBlock(ctx)
	if err != nil {
		lsn.l.Errorw("error getting latest block", "err", err)
		return 0, fmt.Errorf("LogPoller.LatestBlock(): %w", err)
	}
	ll := lsn.l.With(
		"currLastProcessedBlock", currLastProcessedBlock,
		"latestBlock", latestBlock.BlockNumber,
		"latestFinalizedBlock", latestBlock.FinalizedBlockNumber)
	ll.Debugw("updating last processed block")
	defer func() {
		ll.Debugw("done updating last processed block", "elapsed", time.Since(start))
	}()

	logs, err := lp.LogsWithSigs(
		ctx,
		currLastProcessedBlock,
		latestBlock.FinalizedBlockNumber,
		[]common.Hash{lsn.coordinator.RandomWordsFulfilledTopic(), lsn.coordinator.RandomWordsRequestedTopic()},
		lsn.coordinator.Address(),
	)
	if err != nil {
		return currLastProcessedBlock, fmt.Errorf("LogPoller.LogsWithSigs: %w", err)
	}

	unfulfilled, unfulfilledLP, _ := lsn.getUnfulfilled(logs, ll)
	// find request block of earliest unfulfilled request
	// even if this block is > latest finalized, we use latest finalized as earliest unprocessed
	// because re-orgs can occur on any unfinalized block.
	var earliestUnprocessedRequestBlock = latestBlock.FinalizedBlockNumber
	for i, req := range unfulfilled {
		// need to drop requests that have timed out otherwise the earliestUnprocessedRequestBlock
		// will be unnecessarily far back and our queries will be slower.
		if unfulfilledLP[i].CreatedAt.Before(time.Now().UTC().Add(-lsn.job.VRFSpec.RequestTimeout)) {
			// request timed out, don't process
			lsn.l.Debugw("request timed out, skipping",
				"reqID", req.RequestID(),
			)
			continue
		}
		if req.Raw().BlockNumber < uint64(earliestUnprocessedRequestBlock) {
			earliestUnprocessedRequestBlock = int64(req.Raw().BlockNumber)
		}
	}

	return earliestUnprocessedRequestBlock, nil
}

// pollLogs uses the log poller to poll for the latest VRF logs
func (lsn *listenerV2) pollLogs(ctx context.Context, minConfs uint32, lastProcessedBlock int64) (pending []pendingRequest, err error) {
	start := time.Now()
	lp := lsn.chain.LogPoller()

	// latest unfinalized block used on purpose to get bleeding edge logs
	// we don't really have the luxury to wait for finalization on most chains
	// if we want to fulfill on time.
	latestBlock, err := lp.LatestBlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("LogPoller.LatestBlock(): %w", err)
	}
	lsn.setLatestHead(latestBlock)
	ll := lsn.l.With(
		"lastProcessedBlock", lastProcessedBlock,
		"minConfs", minConfs,
		"latestBlock", latestBlock.BlockNumber,
		"latestFinalizedBlock", latestBlock.FinalizedBlockNumber)
	ll.Debugw("polling for logs")
	defer func() {
		ll.Debugw("done polling for logs", "elapsed", time.Since(start))
	}()

	// We don't specify confs because each request can have a different conf above
	// the minimum. So we do all conf handling in getConfirmedAt.
	logs, err := lp.LogsWithSigs(
		ctx,
		lastProcessedBlock,
		latestBlock.BlockNumber,
		[]common.Hash{lsn.coordinator.RandomWordsFulfilledTopic(), lsn.coordinator.RandomWordsRequestedTopic()},
		lsn.coordinator.Address(),
	)
	if err != nil {
		return nil, fmt.Errorf("LogPoller.LogsWithSigs: %w", err)
	}

	unfulfilled, unfulfilledLP, fulfilled := lsn.getUnfulfilled(logs, ll)
	if len(unfulfilled) > 0 {
		ll.Debugw("found unfulfilled logs", "unfulfilled", len(unfulfilled))
	} else {
		ll.Debugw("no unfulfilled logs found")
	}

	lsn.handleFulfilled(fulfilled)

	return lsn.handleRequested(unfulfilled, unfulfilledLP, minConfs), nil
}

func (lsn *listenerV2) getUnfulfilled(logs []logpoller.Log, ll logger.Logger) (unfulfilled []RandomWordsRequested, unfulfilledLP []logpoller.Log, fulfilled map[string]RandomWordsFulfilled) {
	var (
		requested       = make(map[string]RandomWordsRequested)
		requestedLP     = make(map[string]logpoller.Log)
		errs            error
		expectedKeyHash = lsn.job.VRFSpec.PublicKey.MustHash()
	)
	fulfilled = make(map[string]RandomWordsFulfilled)
	for _, l := range logs {
		if l.EventSig == lsn.coordinator.RandomWordsFulfilledTopic() {
			parsed, err2 := lsn.coordinator.ParseRandomWordsFulfilled(l.ToGethLog())
			if err2 != nil {
				// should never happen
				errs = multierr.Append(errs, err2)
				continue
			}
			fulfilled[parsed.RequestID().String()] = parsed
		} else if l.EventSig == lsn.coordinator.RandomWordsRequestedTopic() {
			parsed, err2 := lsn.coordinator.ParseRandomWordsRequested(l.ToGethLog())
			if err2 != nil {
				// should never happen
				errs = multierr.Append(errs, err2)
				continue
			}
			keyHash := parsed.KeyHash()
			if !bytes.Equal(keyHash[:], expectedKeyHash[:]) {
				// wrong keyhash, can ignore
				continue
			}
			requested[parsed.RequestID().String()] = parsed
			requestedLP[parsed.RequestID().String()] = l
		}
	}
	// should never happen, unsure if recoverable
	// may be worth a panic
	if errs != nil {
		ll.Errorw("encountered parse errors", "err", errs)
	}

	if len(fulfilled) > 0 || len(requested) > 0 {
		ll.Infow("found logs", "fulfilled", len(fulfilled), "requested", len(requested))
	} else {
		ll.Debugw("no logs found")
	}

	// find unfulfilled requests by comparing requested events with the fulfilled events
	for reqID, req := range requested {
		if _, isFulfilled := fulfilled[reqID]; !isFulfilled {
			unfulfilled = append(unfulfilled, req)
			unfulfilledLP = append(unfulfilledLP, requestedLP[reqID])
		}
	}

	return unfulfilled, unfulfilledLP, fulfilled
}

func (lsn *listenerV2) getConfirmedAt(req RandomWordsRequested, nodeMinConfs uint32) uint64 {
	// Take the max(nodeMinConfs, requestedConfs + requestedConfsDelay).
	// Add the requested confs delay if provided in the jobspec so that we avoid an edge case
	// where the primary and backup VRF v2 nodes submit a proof at the same time.
	minConfs := nodeMinConfs
	if uint32(req.MinimumRequestConfirmations())+uint32(lsn.job.VRFSpec.RequestedConfsDelay) > nodeMinConfs {
		minConfs = uint32(req.MinimumRequestConfirmations()) + uint32(lsn.job.VRFSpec.RequestedConfsDelay)
	}
	newConfs := uint64(minConfs) * (1 << lsn.respCount[req.RequestID().String()])
	// We cap this at 200 because solidity only supports the most recent 256 blocks
	// in the contract so if it was older than that, fulfillments would start failing
	// without the blockhash store feeder. We use 200 to give the node plenty of time
	// to fulfill even on fast chains.
	if newConfs > 200 {
		newConfs = 200
	}
	if lsn.respCount[req.RequestID().String()] > 0 {
		lsn.l.Warnw("Duplicate request found after fulfillment, doubling incoming confirmations",
			"txHash", req.Raw().TxHash,
			"blockNumber", req.Raw().BlockNumber,
			"blockHash", req.Raw().BlockHash,
			"reqID", req.RequestID().String(),
			"newConfs", newConfs)
		vrfcommon.IncDupeReqs(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, lsn.coordinator.Version())
	}
	return req.Raw().BlockNumber + newConfs
}

func (lsn *listenerV2) handleFulfilled(fulfilled map[string]RandomWordsFulfilled) {
	for _, v := range fulfilled {
		// don't process same log over again
		// log key includes block number and blockhash, so on re-orgs it would return true
		// and we would re-process the re-orged request.
		if !lsn.fulfillmentLogDeduper.ShouldDeliver(v.Raw()) {
			continue
		}
		lsn.l.Debugw("Received fulfilled log", "reqID", v.RequestID(), "success", v.Success())
		lsn.respCount[v.RequestID().String()]++
		lsn.blockNumberToReqID.Insert(fulfilledReqV2{
			blockNumber: v.Raw().BlockNumber,
			reqID:       v.RequestID().String(),
		})
	}
}

func (lsn *listenerV2) handleRequested(requested []RandomWordsRequested, requestedLP []logpoller.Log, minConfs uint32) (pendingRequests []pendingRequest) {
	for i, req := range requested {
		// don't process same log over again
		// log key includes block number and blockhash, so on re-orgs it would return true
		// and we would re-process the re-orged request.
		if lsn.inflightCache.Contains(req.Raw()) {
			continue
		}

		confirmedAt := lsn.getConfirmedAt(req, minConfs)
		lsn.l.Debugw("VRFListenerV2: Received log request",
			"reqID", req.RequestID(),
			"reqBlockNumber", req.Raw().BlockNumber,
			"reqBlockHash", req.Raw().BlockHash,
			"reqTxHash", req.Raw().TxHash,
			"confirmedAt", confirmedAt,
			"subID", req.SubID(),
			"sender", req.Sender())
		pendingRequests = append(pendingRequests, pendingRequest{
			confirmedAtBlock: confirmedAt,
			req:              req,
			utcTimestamp:     requestedLP[i].CreatedAt.UTC(),
		})
		lsn.reqAdded()
	}

	return pendingRequests
}

// numReplayBlocks returns the number of blocks to replay on startup
// given the request timeout and the chain ID.
// if the chain ID is not recognized it assumes a block time of 1 second
// and returns the number of blocks in a day.
func numReplayBlocks(requestTimeout time.Duration, chainID *big.Int) int64 {
	var timeoutSeconds = int64(requestTimeout.Seconds())
	switch chainID.String() {
	case
		"1",        // eth mainnet
		"3",        // eth robsten
		"4",        // eth rinkeby
		"5",        // eth goerli
		"11155111": // eth sepolia
		// block time is 12s
		return timeoutSeconds / 12
	case
		"137",   // polygon mainnet
		"80001", // polygon mumbai
		"80002": // polygon amoy
		// block time is 2s
		return timeoutSeconds / 2
	case
		"56", // bsc mainnet
		"97": // bsc testnet
		// block time is 2s
		return timeoutSeconds / 2
	case
		"43114", // avalanche mainnet
		"43113": // avalanche fuji
		// block time is 1s
		return timeoutSeconds
	case
		"250",  // fantom mainnet
		"4002": // fantom testnet
		// block time is 1s
		return timeoutSeconds
	case
		"42161",  // arbitrum mainnet
		"421613", // arbitrum goerli
		"421614": // arbitrum sepolia
		// block time is 0.25s in the worst case
		return timeoutSeconds * 4
	case
		"10",       // optimism mainnet
		"69",       // optimism kovan
		"420",      // optimism goerli
		"11155420": // optimism sepolia
		// block time is 2s
		return timeoutSeconds / 2
	case
		"8453",  // base mainnet
		"84531", // base goerli
		"84532": // base sepolia
		// block time is 2s
		return timeoutSeconds / 2
	default:
		// assume block time of 1s
		return timeoutSeconds
	}
}
