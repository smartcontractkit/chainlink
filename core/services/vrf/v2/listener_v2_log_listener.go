package v2

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func (lsn *listenerV2) runLogListener(
	pollPeriod time.Duration,
	minConfs uint32,
) {
	lsn.l.Infow("Listening for run requests via log poller",
		"minConfs", minConfs)
	ticker := time.NewTicker(utils.WithJitter(pollPeriod / 2))
	defer ticker.Stop()
	var (
		lastProcessedBlock int64
		startingUp         = true
	)
	ctx, cancel := lsn.chStop.NewCtx()
	defer cancel()
	for {
		select {
		case <-lsn.chStop:
			return
		case <-ticker.C:
			// Filter registration is idempotent, so we can just call it every time
			// and retry on errors using the ticker.
			err := lsn.chain.LogPoller().RegisterFilter(logpoller.Filter{
				Name: fmt.Sprintf("vrf_%s_keyhash_%s_job_%d", lsn.coordinator.Version(), lsn.job.VRFSpec.PublicKey.MustHash().String(), lsn.job.ID),
				EventSigs: evmtypes.HashArray{
					lsn.coordinator.RandomWordsFulfilledTopic(),
					lsn.coordinator.RandomWordsRequestedTopic(),
				},
				Addresses: evmtypes.AddressArray{
					lsn.coordinator.Address(),
				},
			})
			if err != nil {
				lsn.l.Errorw("error registering filter in log poller, retrying", "err", err)
				continue
			}

			// on startup we want to initialize the last processed block
			if startingUp {
				lsn.l.Infow("initializing last processed block on startup")
				lastProcessedBlock, err = lsn.initializeLastProcessedBlock()
				if err != nil {
					lsn.l.Errorw("error initializing last processed block, retrying", "err", err)
					continue
				}
				startingUp = false
			}

			pending, err := lsn.pollLogs(ctx, minConfs, lastProcessedBlock)
			if err != nil {
				lsn.l.Errorw("error polling vrf logs, retrying", "err", err)
				continue
			}

			// process pending requests and insert any fulfillments into the inflight cache
			lsn.processPendingVRFRequests(ctx, pending)

			lastProcessedBlock, err = lsn.updateLastProcessedBlock(ctx, lastProcessedBlock)
			if err != nil {
				lsn.l.Errorw("error updating last processed block, continuing anyway", "err", err)
			} else {
				lsn.l.Infow("updated last processed block", "lastProcessedBlock", lastProcessedBlock)
			}
		}
	}
}

// initializeLastProcessedBlock returns the earliest block number that we need to
// process requests for. This is the block number of the earliest unfulfilled request
// or the latest finalized block, if there are no unfulfilled requests.
func (lsn *listenerV2) initializeLastProcessedBlock() (lastProcessedBlock int64, err error) {
	lp := lsn.chain.LogPoller()

	// will retry on error in the runLogListener loop
	latestBlock, err := lp.LatestBlock()
	if err != nil {
		return 0, errors.Wrap(err, "LogPoller.LatestBlock()")
	}

	ll := lsn.l.With("latestFinalizedBlock", latestBlock.FinalizedBlockNumber, "latestBlock", latestBlock.BlockNumber)
	ll.Infow("Initializing last processed block")

	fromTimestamp := time.Now().UTC().Add(-lsn.job.VRFSpec.RequestTimeout)
	// get randomness requested logs with the appropriate keyhash
	// keyhash is specified in topic1
	requests, err := lp.IndexedLogsCreatedAfter(
		lsn.coordinator.RandomWordsRequestedTopic(), // event sig
		lsn.coordinator.Address(),                   // address
		1,                                           // topic index
		[]common.Hash{lsn.job.VRFSpec.PublicKey.MustHash()}, // topic values
		fromTimestamp,       // from time
		logpoller.Finalized, // confs
	)
	if err != nil {
		return 0, errors.Wrap(err, "LogPoller.LogsCreatedAfter RandomWordsRequested logs")
	}

	// fulfillments don't have keyhash indexed, we'll have to get all of them
	fulfillments, err := lp.LogsCreatedAfter(
		lsn.coordinator.RandomWordsFulfilledTopic(), // event sig
		lsn.coordinator.Address(),                   // address
		fromTimestamp,                               // from time
		logpoller.Finalized,                         // confs
	)
	if err != nil {
		return 0, errors.Wrap(err, "LogPoller.LogsCreatedAfter RandomWordsFulfilled logs")
	}

	unfulfilled, _, _ := lsn.getUnfulfilled(append(requests, fulfillments...), ll)
	var earliestUnfulfilledBlock int64 = math.MaxInt64
	for _, req := range unfulfilled {
		if req.Raw().BlockNumber < uint64(earliestUnfulfilledBlock) {
			earliestUnfulfilledBlock = int64(req.Raw().BlockNumber)
		}
	}
	if earliestUnfulfilledBlock == math.MaxInt64 {
		// no unfulfilled requests
		return latestBlock.FinalizedBlockNumber, nil
	}

	// earliestUnfulfilledBlock is <= latestBlock.FinalizedBlockNumber because in our queries we specify
	// logpoller.Finalized as the number of confirmations, so it's impossible for this not to be the case.
	return earliestUnfulfilledBlock, nil
}

func (lsn *listenerV2) updateLastProcessedBlock(ctx context.Context, currLastProcessedBlock int64) (lastProcessedBlock int64, err error) {
	lp := lsn.chain.LogPoller()
	ll := lsn.l.With("currLastProcessedBlock", currLastProcessedBlock)

	latestBlock, err := lp.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		lsn.l.Errorw("error getting latest block", "err", err)
		return 0, errors.Wrap(err, "LogPoller.LatestBlock()")
	}

	logs, err := lp.LogsWithSigs(
		currLastProcessedBlock,
		latestBlock.FinalizedBlockNumber,
		[]common.Hash{lsn.coordinator.RandomWordsFulfilledTopic(), lsn.coordinator.RandomWordsRequestedTopic()},
		lsn.coordinator.Address(),
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return currLastProcessedBlock, errors.Wrap(err, "LogPoller.LogsWithSigs")
	}

	unfulfilled, _, _ := lsn.getUnfulfilled(logs, ll)
	// find request block of earliest unfulfilled request
	var earliestUnprocessedRequestBlock int64 = math.MaxInt64
	for _, req := range unfulfilled {
		if req.Raw().BlockNumber < uint64(earliestUnprocessedRequestBlock) {
			earliestUnprocessedRequestBlock = int64(req.Raw().BlockNumber)
		}
	}

	// cases:
	// 1. pending requests exist, earliestUnprocessedRequestBlock < latestBlock.FinalizedBlockNumber
	// 2. pending requests exist, earliestUnprocessedRequestBlock > latestBlock.FinalizedBlockNumber
	// 3. no pending requests, earliestUnprocessedRequestBlock == math.MaxInt64
	// case 2 or 3
	if earliestUnprocessedRequestBlock == math.MaxInt64 || earliestUnprocessedRequestBlock > latestBlock.FinalizedBlockNumber {
		return latestBlock.FinalizedBlockNumber, nil
	}
	// case 1
	return earliestUnprocessedRequestBlock, nil
}

// pollLogs uses the log poller to poll for the latest VRF logs
func (lsn *listenerV2) pollLogs(ctx context.Context, minConfs uint32, lastProcessedBlock int64) (pending []pendingRequest, err error) {
	lp := lsn.chain.LogPoller()
	ll := lsn.l.With("lastProcessedBlock", lastProcessedBlock, "minConfs", minConfs)

	// latest unfinalized block used on purpose to get bleeding edge logs
	// we don't really have the luxury to wait for finalization on most chains
	// if we want to fulfill on time.
	latestBlock, err := lp.LatestBlock()
	if err != nil {
		return nil, errors.Wrap(err, "LogPoller.LatestBlock()")
	}
	lsn.setLatestHead(latestBlock)

	// We don't specify confs because each request can have a different conf above
	// the minimum. So we do all conf handling in getConfirmedAt.
	logs, err := lp.LogsWithSigs(
		lastProcessedBlock,
		latestBlock.BlockNumber,
		[]common.Hash{lsn.coordinator.RandomWordsFulfilledTopic(), lsn.coordinator.RandomWordsRequestedTopic()},
		lsn.coordinator.Address(),
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, errors.Wrap(err, "LogPoller.LogsWithSigs")
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
		lsn.l.Infow("VRFListenerV2: Received log request", "reqID", req.RequestID(), "confirmedAt", confirmedAt, "subID", req.SubID(), "sender", req.Sender())
		pendingRequests = append(pendingRequests, pendingRequest{
			confirmedAtBlock: confirmedAt,
			req:              req,
			utcTimestamp:     requestedLP[i].CreatedAt.UTC(),
		})
		lsn.reqAdded()
	}

	return pendingRequests
}
