package v2

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"
	"math"
	"math/big"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/hex"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
	"github.com/smartcontractkit/chainlink-integrations/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2plus_interface"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	bigmath "github.com/smartcontractkit/chainlink/v2/core/utils/big_math"
)

// Returns all the confirmed logs from the provided pending queue by subscription
func (lsn *listenerV2) getConfirmedLogsBySub(latestHead uint64, pendingRequests []pendingRequest) map[string][]pendingRequest {
	vrfcommon.UpdateQueueSize(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, lsn.coordinator.Version(), uniqueReqs(pendingRequests))
	var toProcess = make(map[string][]pendingRequest)
	for _, request := range pendingRequests {
		if lsn.ready(request, latestHead) {
			toProcess[request.req.SubID().String()] = append(toProcess[request.req.SubID().String()], request)
		}
	}
	return toProcess
}

func (lsn *listenerV2) ready(req pendingRequest, latestHead uint64) bool {
	// Request is not eligible for fulfillment yet
	if req.confirmedAtBlock > latestHead {
		return false
	}

	if lsn.job.VRFSpec.BackoffInitialDelay == 0 || req.attempts == 0 {
		// Backoff is disabled, or this is the first try
		return true
	}

	return time.Now().UTC().After(
		nextTry(
			req.attempts,
			lsn.job.VRFSpec.BackoffInitialDelay,
			lsn.job.VRFSpec.BackoffMaxDelay,
			req.lastTry))
}

func nextTry(retries int, initial, max time.Duration, last time.Time) time.Time {
	expBackoffFactor := math.Pow(backoffFactor, float64(retries-1))

	var delay time.Duration
	if expBackoffFactor > float64(max/initial) {
		delay = max
	} else {
		delay = time.Duration(float64(initial) * expBackoffFactor)
	}
	return last.Add(delay)
}

// Remove all entries 10000 blocks or older
// to avoid a memory leak.
func (lsn *listenerV2) pruneConfirmedRequestCounts() {
	min := lsn.blockNumberToReqID.FindMin()
	for min != nil {
		m := min.(fulfilledReqV2)
		if m.blockNumber > (lsn.getLatestHead() - 10000) {
			break
		}
		delete(lsn.respCount, m.reqID)
		lsn.blockNumberToReqID.DeleteMin()
		min = lsn.blockNumberToReqID.FindMin()
	}
}

// Determine a set of logs that are confirmed
// and the subscription has sufficient balance to fulfill,
// given a eth call with the max gas price.
// Note we have to consider the pending reqs already in the txm as already "spent" link or native,
// using a max link or max native consumed in their metadata.
// A user will need a minBalance capable of fulfilling a single req at the max gas price or nothing will happen.
// This is acceptable as users can choose different keyhashes which have different max gas prices.
// Other variables which can change the bill amount between our eth call simulation and tx execution:
// - Link/eth price fluctuation
// - Falling back to BHS
// However the likelihood is vanishingly small as
// 1) the window between simulation and tx execution is tiny.
// 2) the max gas price provides a very large buffer most of the time.
// Its easier to optimistically assume it will go though and in the rare case of a reversion
// we simply retry TODO: follow up where if we see a fulfillment revert, return log to the queue.
func (lsn *listenerV2) processPendingVRFRequests(ctx context.Context, pendingRequests []pendingRequest) {
	confirmed := lsn.getConfirmedLogsBySub(lsn.getLatestHead(), pendingRequests)
	var processedMu sync.Mutex
	processed := make(map[string]struct{})
	start := time.Now()

	defer func() {
		for _, subReqs := range confirmed {
			for _, req := range subReqs {
				if _, ok := processed[req.req.RequestID().String()]; ok {
					// add to the inflight cache so that we don't re-process this request
					lsn.inflightCache.Add(req.req.Raw())
				}
			}
		}
		lsn.l.Infow("Finished processing pending requests",
			"totalProcessed", len(processed),
			"totalFailed", len(pendingRequests)-len(processed),
			"total", len(pendingRequests),
			"time", time.Since(start).String(),
			"inflightCacheSize", lsn.inflightCache.Size())
	}()

	if len(confirmed) == 0 {
		lsn.l.Infow("No pending requests ready for processing")
		return
	}
	for subID, reqs := range confirmed {
		l := lsn.l.With("subID", subID, "startTime", time.Now(), "numReqsForSub", len(reqs))
		// Get the balance of the subscription and also it's active status.
		// The reason we need both is that we cannot determine if a subscription
		// is active solely by it's balance, since an active subscription could legitimately
		// have a zero balance.
		var (
			startLinkBalance *big.Int
			startEthBalance  *big.Int
			subIsActive      bool
		)
		sID, ok := new(big.Int).SetString(subID, 10)
		if !ok {
			l.Criticalf("Unable to convert %s to Int", subID)
			return
		}
		sub, err := lsn.coordinator.GetSubscription(&bind.CallOpts{
			Context: ctx}, sID)

		if err != nil {
			if !strings.Contains(err.Error(), "execution reverted") {
				// Most likely this is an RPC error, so we re-try later.
				l.Errorw("Unable to read subscription balance", "err", err)
				return
			}
			// "execution reverted" indicates that the subscription no longer exists.
			// We can no longer just mark these as processed and continue,
			// since it could be that the subscription was canceled while there
			// were still unfulfilled requests.
			// The simplest approach to handle this is to enter the processRequestsPerSub
			// loop rather than create a bunch of largely duplicated code
			// to handle this specific situation, since we need to run the pipeline to get
			// the VRF proof, abi-encode it, etc.
			l.Warnw("Subscription not found - setting start balance to zero", "subID", subID, "err", err)
			startLinkBalance = big.NewInt(0)
		} else {
			// Happy path - sub is active.
			startLinkBalance = sub.Balance()
			if sub.Version() == vrfcommon.V2Plus {
				startEthBalance = sub.NativeBalance()
			}
			subIsActive = true
		}

		// Sort requests in ascending order by CallbackGasLimit
		// so that we process the "cheapest" requests for each subscription
		// first. This allows us to break out of the processing loop as early as possible
		// in the event that a subscription is too underfunded to have it's
		// requests processed.
		slices.SortFunc(reqs, func(a, b pendingRequest) int {
			return cmp.Compare(a.req.CallbackGasLimit(), b.req.CallbackGasLimit())
		})

		p := lsn.processRequestsPerSub(ctx, sID, startLinkBalance, startEthBalance, reqs, subIsActive)
		processedMu.Lock()
		for reqID := range p {
			processed[reqID] = struct{}{}
		}
		processedMu.Unlock()
	}
	lsn.pruneConfirmedRequestCounts()
}

// MaybeSubtractReservedLink figures out how much LINK is reserved for other VRF requests that
// have not been fully confirmed yet on-chain, and subtracts that from the given startBalance,
// and returns that value if there are no errors.
func (lsn *listenerV2) MaybeSubtractReservedLink(ctx context.Context, startBalance *big.Int, chainID *big.Int, subID *big.Int, vrfVersion vrfcommon.Version) (*big.Int, error) {
	var metaField string
	if vrfVersion == vrfcommon.V2Plus {
		metaField = txMetaGlobalSubId
	} else if vrfVersion == vrfcommon.V2 {
		metaField = txMetaFieldSubId
	} else {
		return nil, errors.Errorf("unsupported vrf version %s", vrfVersion)
	}

	txes, err := lsn.chain.TxManager().FindTxesByMetaFieldAndStates(ctx, metaField, subID.String(), reserveEthLinkQueryStates, chainID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("TXM FindTxesByMetaFieldAndStates failed: %w", err)
	}

	reservedLinkSum := big.NewInt(0)
	// Aggregate non-null MaxLink from all txes returned
	for _, tx := range txes {
		var meta *txmgrtypes.TxMeta[common.Address, common.Hash]
		meta, err = tx.GetMeta()
		if err != nil {
			return nil, fmt.Errorf("GetMeta for Tx failed: %w", err)
		}
		if meta != nil && meta.MaxLink != nil {
			txMaxLink, success := new(big.Int).SetString(*meta.MaxLink, 10)
			if !success {
				return nil, fmt.Errorf("converting reserved LINK %s", *meta.MaxLink)
			}

			reservedLinkSum.Add(reservedLinkSum, txMaxLink)
		}
	}

	return new(big.Int).Sub(startBalance, reservedLinkSum), nil
}

// MaybeSubtractReservedEth figures out how much ether is reserved for other VRF requests that
// have not been fully confirmed yet on-chain, and subtracts that from the given startBalance,
// and returns that value if there are no errors.
func (lsn *listenerV2) MaybeSubtractReservedEth(ctx context.Context, startBalance *big.Int, chainID *big.Int, subID *big.Int, vrfVersion vrfcommon.Version) (*big.Int, error) {
	var metaField string
	if vrfVersion == vrfcommon.V2Plus {
		metaField = txMetaGlobalSubId
	} else if vrfVersion == vrfcommon.V2 {
		// native payment is not supported for v2, so returning 0 reserved ETH
		return big.NewInt(0), nil
	} else {
		return nil, errors.Errorf("unsupported vrf version %s", vrfVersion)
	}
	txes, err := lsn.chain.TxManager().FindTxesByMetaFieldAndStates(ctx, metaField, subID.String(), reserveEthLinkQueryStates, chainID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("TXM FindTxesByMetaFieldAndStates failed: %w", err)
	}

	reservedEthSum := big.NewInt(0)
	// Aggregate non-null MaxEth from all txes returned
	for _, tx := range txes {
		var meta *txmgrtypes.TxMeta[common.Address, common.Hash]
		meta, err = tx.GetMeta()
		if err != nil {
			return nil, fmt.Errorf("GetMeta for Tx failed: %w", err)
		}
		if meta != nil && meta.MaxEth != nil {
			txMaxEth, success := new(big.Int).SetString(*meta.MaxEth, 10)
			if !success {
				return nil, fmt.Errorf("converting reserved ETH %s", *meta.MaxEth)
			}

			reservedEthSum.Add(reservedEthSum, txMaxEth)
		}
	}

	if startBalance != nil {
		return new(big.Int).Sub(startBalance, reservedEthSum), nil
	}
	return big.NewInt(0), nil
}

func (lsn *listenerV2) processRequestsPerSubBatchHelper(
	ctx context.Context,
	subID *big.Int,
	startBalance *big.Int,
	startBalanceNoReserved *big.Int,
	reqs []pendingRequest,
	subIsActive bool,
	nativePayment bool,
) (processed map[string]struct{}) {
	start := time.Now()
	processed = make(map[string]struct{})

	// Base the max gas for a batch on the max gas limit for a single callback.
	// Since the max gas limit for a single callback is usually quite large already,
	// we probably don't want to exceed it too much so that we can reliably get
	// batch fulfillments included, while also making sure that the biggest gas guzzler
	// callbacks are included.
	config, err := lsn.coordinator.GetConfig(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		lsn.l.Errorw("Couldn't get config from coordinator", "err", err)
		return processed
	}

	// Add very conservative upper bound estimate on verification costs.
	batchMaxGas := config.MaxGasLimit() + 400_000

	l := lsn.l.With(
		"subID", subID,
		"eligibleSubReqs", len(reqs),
		"startBalance", startBalance.String(),
		"startBalanceNoReserved", startBalanceNoReserved.String(),
		"batchMaxGas", batchMaxGas,
		"subIsActive", subIsActive,
		"nativePayment", nativePayment,
	)

	defer func() {
		l.Infow("Finished processing for sub",
			"endBalance", startBalanceNoReserved.String(),
			"totalProcessed", len(processed),
			"totalUnique", uniqueReqs(reqs),
			"time", time.Since(start).String())
	}()

	l.Infow("Processing requests for subscription with batching")

	ready, expired := lsn.getReadyAndExpired(l, reqs)
	for _, reqID := range expired {
		processed[reqID] = struct{}{}
	}

	// Process requests in chunks in order to kick off as many jobs
	// as configured in parallel. Then we can combine into fulfillment
	// batches afterwards.
	for chunkStart := 0; chunkStart < len(ready); chunkStart += int(lsn.job.VRFSpec.ChunkSize) {
		chunkEnd := chunkStart + int(lsn.job.VRFSpec.ChunkSize)
		if chunkEnd > len(ready) {
			chunkEnd = len(ready)
		}
		chunk := ready[chunkStart:chunkEnd]

		var unfulfilled []pendingRequest
		alreadyFulfilled, err := lsn.checkReqsFulfilled(ctx, l, chunk)
		if errors.Is(err, context.Canceled) {
			l.Infow("Context canceled, stopping request processing", "err", err)
			return processed
		} else if err != nil {
			l.Errorw("Error checking for already fulfilled requests, proceeding anyway", "err", err)
		}
		for i, a := range alreadyFulfilled {
			if a {
				processed[chunk[i].req.RequestID().String()] = struct{}{}
			} else {
				unfulfilled = append(unfulfilled, chunk[i])
			}
		}

		// All fromAddresses passed to the VRFv2 job have the same KeySpecific-MaxPrice value.
		fromAddresses := lsn.fromAddresses()
		maxGasPriceWei := lsn.feeCfg.PriceMaxKey(fromAddresses[0])

		// Cases:
		// 1. Never simulated: in this case, we want to observe the time until simulated
		// on the utcTimestamp field of the pending request.
		// 2. Simulated before: in this case, lastTry will be set to a non-zero time value,
		// in which case we'd want to use that as a relative point from when we last tried
		// the request.
		observeRequestSimDuration(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, lsn.coordinator.Version(), unfulfilled)

		pipelines := lsn.runPipelines(ctx, l, maxGasPriceWei, unfulfilled)
		batches := newBatchFulfillments(batchMaxGas, lsn.coordinator.Version())
		outOfBalance := false
		for _, p := range pipelines {
			ll := l.With("reqID", p.req.req.RequestID().String(),
				"txHash", p.req.req.Raw().TxHash,
				"maxGasPrice", maxGasPriceWei.String(),
				"fundsNeeded", p.fundsNeeded.String(),
				"maxFee", p.maxFee.String(),
				"gasLimit", p.gasLimit,
				"attempts", p.req.attempts,
				"remainingBalance", startBalanceNoReserved.String(),
				"consumerAddress", p.req.req.Sender(),
				"blockNumber", p.req.req.Raw().BlockNumber,
				"blockHash", p.req.req.Raw().BlockHash,
			)
			fromAddresses := lsn.fromAddresses()
			fromAddress, err := lsn.gethks.GetRoundRobinAddress(ctx, lsn.chainID, fromAddresses...)
			if err != nil {
				l.Errorw("Couldn't get next from address", "err", err)
				continue
			}
			ll = ll.With("fromAddress", fromAddress)

			if p.err != nil {
				if errors.Is(p.err, errBlockhashNotInStore{}) {
					// Running the blockhash store feeder in backwards mode will be required to
					// resolve this.
					ll.Criticalw("Pipeline error", "err", p.err)
				} else if errors.Is(p.err, errProofVerificationFailed{}) {
					// This occurs when the proof reverts in the simulation
					// This is almost always (if not always) due to a proof generated with an out-of-date
					// blockhash
					// we can simply mark as processed and move on, since we will eventually
					// process the request with the right blockhash
					ll.Infow("proof reverted in simulation, likely stale blockhash")
					processed[p.req.req.RequestID().String()] = struct{}{}
				} else {
					ll.Errorw("Pipeline error", "err", p.err)
					if !subIsActive {
						ll.Warnw("Force-fulfilling a request with insufficient funds on a cancelled sub")
						etx, err := lsn.enqueueForceFulfillment(ctx, p, fromAddress)
						if err != nil {
							ll.Errorw("Error enqueuing force-fulfillment, re-queueing request", "err", err)
							continue
						}
						ll.Infow("Successfully enqueued force-fulfillment", "ethTxID", etx.ID)
						processed[p.req.req.RequestID().String()] = struct{}{}

						// Need to put a continue here, otherwise the next if statement will be hit
						// and we'd break out of the loop prematurely.
						// If a sub is canceled, we want to force-fulfill ALL of it's pending requests
						// before saying we're done with it.
						continue
					}

					if startBalanceNoReserved.Cmp(p.fundsNeeded) < 0 && errors.Is(p.err, errPossiblyInsufficientFunds{}) {
						ll.Infow("Insufficient balance to fulfill a request based on estimate, breaking", "err", p.err)
						outOfBalance = true

						// break out of this inner loop to process the currently constructed batch
						break
					}

					// Ensure consumer is valid, otherwise drop the request.
					if !lsn.isConsumerValidAfterFinalityDepthElapsed(ctx, p.req) {
						lsn.l.Infow(
							"Dropping request that was made by an invalid consumer.",
							"consumerAddress", p.req.req.Sender(),
							"reqID", p.req.req.RequestID(),
							"blockNumber", p.req.req.Raw().BlockNumber,
							"blockHash", p.req.req.Raw().BlockHash,
						)
						processed[p.req.req.RequestID().String()] = struct{}{}
						continue
					}
				}
				continue
			}

			if startBalanceNoReserved.Cmp(p.maxFee) < 0 {
				// Insufficient funds, have to wait for a user top up.
				// Break out of the loop now and process what we are able to process
				// in the constructed batches.
				ll.Infow("Insufficient balance to fulfill a request, breaking")
				break
			}

			batches.addRun(p, fromAddress)

			startBalanceNoReserved.Sub(startBalanceNoReserved, p.maxFee)
		}

		var processedRequestIDs []string
		for _, batch := range batches.fulfillments {
			l.Debugw("Processing batch", "batchSize", len(batch.proofs))
			p := lsn.processBatch(l, subID, startBalanceNoReserved, batchMaxGas, batch, batch.fromAddress)
			processedRequestIDs = append(processedRequestIDs, p...)
		}

		for _, reqID := range processedRequestIDs {
			processed[reqID] = struct{}{}
		}

		// outOfBalance is set to true if the current sub we are processing
		// has run out of funds to process any remaining requests. After enqueueing
		// this constructed batch, we break out of this outer loop in order to
		// avoid unnecessarily processing the remaining requests.
		if outOfBalance {
			break
		}
	}

	return
}

func (lsn *listenerV2) processRequestsPerSubBatch(
	ctx context.Context,
	subID *big.Int,
	startLinkBalance *big.Int,
	startEthBalance *big.Int,
	reqs []pendingRequest,
	subIsActive bool,
) map[string]struct{} {
	var processed = make(map[string]struct{})
	startBalanceNoReserveLink, err := lsn.MaybeSubtractReservedLink(
		ctx, startLinkBalance, lsn.chainID, subID, lsn.coordinator.Version())
	if err != nil {
		lsn.l.Errorw("Couldn't get reserved LINK for subscription", "sub", reqs[0].req.SubID(), "err", err)
		return processed
	}
	startBalanceNoReserveEth, err := lsn.MaybeSubtractReservedEth(
		ctx, startEthBalance, lsn.chainID, subID, lsn.coordinator.Version())
	if err != nil {
		lsn.l.Errorw("Couldn't get reserved ether for subscription", "sub", reqs[0].req.SubID(), "err", err)
		return processed
	}

	// Split the requests into native and LINK requests.
	var (
		nativeRequests []pendingRequest
		linkRequests   []pendingRequest
	)
	for _, req := range reqs {
		if req.req.NativePayment() {
			nativeRequests = append(nativeRequests, req)
		} else {
			linkRequests = append(linkRequests, req)
		}
	}
	// process the native and link requests in parallel
	var wg sync.WaitGroup
	var nativeProcessed, linkProcessed map[string]struct{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		nativeProcessed = lsn.processRequestsPerSubBatchHelper(ctx, subID, startEthBalance, startBalanceNoReserveEth, nativeRequests, subIsActive, true)
	}()
	go func() {
		defer wg.Done()
		linkProcessed = lsn.processRequestsPerSubBatchHelper(ctx, subID, startLinkBalance, startBalanceNoReserveLink, linkRequests, subIsActive, false)
	}()
	wg.Wait()
	// combine the processed link and native requests into the processed map
	for k, v := range nativeProcessed {
		processed[k] = v
	}
	for k, v := range linkProcessed {
		processed[k] = v
	}

	return processed
}

// enqueueForceFulfillment enqueues a forced fulfillment through the
// VRFOwner contract. It estimates gas again on the transaction due
// to the extra steps taken within VRFOwner.fulfillRandomWords.
func (lsn *listenerV2) enqueueForceFulfillment(
	ctx context.Context,
	p vrfPipelineResult,
	fromAddress common.Address,
) (etx txmgr.Tx, err error) {
	if lsn.job.VRFSpec.VRFOwnerAddress == nil {
		err = errors.New("vrf owner address not set in job spec, recreate job and provide it to force-fulfill")
		return
	}

	if p.payload == "" {
		// should probably never happen
		// a critical log will be logged if this is the case in simulateFulfillment
		err = errors.New("empty payload in vrfPipelineResult")
		return
	}

	// fulfill the request through the VRF owner
	lsn.l.Infow("VRFOwner.fulfillRandomWords vs. VRFCoordinatorV2.fulfillRandomWords",
		"vrf_owner.fulfillRandomWords", hexutil.Encode(vrfOwnerABI.Methods["fulfillRandomWords"].ID),
		"vrf_coordinator_v2.fulfillRandomWords", hexutil.Encode(coordinatorV2ABI.Methods["fulfillRandomWords"].ID),
	)

	vrfOwnerAddress1 := lsn.vrfOwner.Address()
	vrfOwnerAddressSpec := lsn.job.VRFSpec.VRFOwnerAddress.Address()
	lsn.l.Infow("addresses diff", "wrapper_address", vrfOwnerAddress1, "spec_address", vrfOwnerAddressSpec)

	lsn.l.Infow("fulfillRandomWords payload", "proof", p.proof, "commitment", p.reqCommitment.Get(), "payload", p.payload)
	txData := hexutil.MustDecode(p.payload)
	if err != nil {
		err = fmt.Errorf("abi pack VRFOwner.fulfillRandomWords: %w", err)
		return
	}
	estimateGasLimit, err := lsn.chain.Client().EstimateGas(ctx, ethereum.CallMsg{
		From: fromAddress,
		To:   &vrfOwnerAddressSpec,
		Data: txData,
	})
	if err != nil {
		err = fmt.Errorf("failed to estimate gas on VRFOwner.fulfillRandomWords: %w", err)
		return
	}

	lsn.l.Infow("Estimated gas limit on force fulfillment",
		"estimateGasLimit", estimateGasLimit, "pipelineGasLimit", p.gasLimit)
	if estimateGasLimit < p.gasLimit {
		estimateGasLimit = p.gasLimit
	}

	requestID := common.BytesToHash(p.req.req.RequestID().Bytes())
	subID := p.req.req.SubID()
	requestTxHash := p.req.req.Raw().TxHash
	return lsn.chain.TxManager().CreateTransaction(ctx, txmgr.TxRequest{
		FromAddress:    fromAddress,
		ToAddress:      lsn.vrfOwner.Address(),
		EncodedPayload: txData,
		FeeLimit:       estimateGasLimit,
		Strategy:       txmgrcommon.NewSendEveryStrategy(),
		Meta: &txmgr.TxMeta{
			RequestID:     &requestID,
			SubID:         ptr(subID.Uint64()),
			RequestTxHash: &requestTxHash,
			// No max link since simulation failed
		},
	})
}

// For an errored pipeline run, wait until the finality depth of the chain to have elapsed,
// then check if the failing request is being called by an invalid sender. Return false if this is the case,
// otherwise true.
func (lsn *listenerV2) isConsumerValidAfterFinalityDepthElapsed(ctx context.Context, req pendingRequest) bool {
	latestHead := lsn.getLatestHead()
	if latestHead-req.req.Raw().BlockNumber > uint64(lsn.cfg.FinalityDepth()) {
		code, err := lsn.chain.Client().CodeAt(ctx, req.req.Sender(), big.NewInt(int64(latestHead)))
		if err != nil {
			lsn.l.Warnw("Failed to fetch contract code", "err", err)
			return true // error fetching code, give the benefit of doubt to the consumer
		}
		if len(code) == 0 {
			return false // invalid consumer
		}
	}

	return true // valid consumer, or finality depth has not elapsed
}

// processRequestsPerSubHelper processes a set of pending requests for the provided sub id.
// It returns a set of request IDs that were processed.
// Note that the provided startBalanceNoReserve is the balance of the subscription
// minus any pending requests that have already been processed and not yet fulfilled onchain.
func (lsn *listenerV2) processRequestsPerSubHelper(
	ctx context.Context,
	subID *big.Int,
	startBalance *big.Int,
	startBalanceNoReserved *big.Int,
	reqs []pendingRequest,
	subIsActive bool,
	nativePayment bool,
) (processed map[string]struct{}) {
	start := time.Now()
	processed = make(map[string]struct{})

	l := lsn.l.With(
		"subID", subID,
		"eligibleSubReqs", len(reqs),
		"startBalance", startBalance.String(),
		"startBalanceNoReserved", startBalanceNoReserved.String(),
		"subIsActive", subIsActive,
		"nativePayment", nativePayment,
	)

	defer func() {
		l.Infow("Finished processing for sub",
			"endBalance", startBalanceNoReserved.String(),
			"totalProcessed", len(processed),
			"totalUnique", uniqueReqs(reqs),
			"time", time.Since(start).String())
	}()

	l.Infow("Processing requests for subscription")

	ready, expired := lsn.getReadyAndExpired(l, reqs)
	for _, reqID := range expired {
		processed[reqID] = struct{}{}
	}

	// Process requests in chunks
	for chunkStart := 0; chunkStart < len(ready); chunkStart += int(lsn.job.VRFSpec.ChunkSize) {
		chunkEnd := chunkStart + int(lsn.job.VRFSpec.ChunkSize)
		if chunkEnd > len(ready) {
			chunkEnd = len(ready)
		}
		chunk := ready[chunkStart:chunkEnd]

		var unfulfilled []pendingRequest
		alreadyFulfilled, err := lsn.checkReqsFulfilled(ctx, l, chunk)
		if errors.Is(err, context.Canceled) {
			l.Infow("Context canceled, stopping request processing", "err", err)
			return processed
		} else if err != nil {
			l.Errorw("Error checking for already fulfilled requests, proceeding anyway", "err", err)
		}
		for i, a := range alreadyFulfilled {
			if a {
				processed[chunk[i].req.RequestID().String()] = struct{}{}
			} else {
				unfulfilled = append(unfulfilled, chunk[i])
			}
		}

		// All fromAddresses passed to the VRFv2 job have the same KeySpecific-MaxPrice value.
		fromAddresses := lsn.fromAddresses()
		maxGasPriceWei := lsn.feeCfg.PriceMaxKey(fromAddresses[0])
		observeRequestSimDuration(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, lsn.coordinator.Version(), unfulfilled)
		pipelines := lsn.runPipelines(ctx, l, maxGasPriceWei, unfulfilled)
		for _, p := range pipelines {
			ll := l.With("reqID", p.req.req.RequestID().String(),
				"txHash", p.req.req.Raw().TxHash,
				"maxGasPrice", maxGasPriceWei.String(),
				"fundsNeeded", p.fundsNeeded.String(),
				"maxFee", p.maxFee.String(),
				"gasLimit", p.gasLimit,
				"attempts", p.req.attempts,
				"remainingBalance", startBalanceNoReserved.String(),
				"consumerAddress", p.req.req.Sender(),
				"blockNumber", p.req.req.Raw().BlockNumber,
				"blockHash", p.req.req.Raw().BlockHash,
			)
			fromAddress, err := lsn.gethks.GetRoundRobinAddress(ctx, lsn.chainID, fromAddresses...)
			if err != nil {
				l.Errorw("Couldn't get next from address", "err", err)
				continue
			}
			ll = ll.With("fromAddress", fromAddress)

			if p.err != nil {
				if errors.Is(p.err, errBlockhashNotInStore{}) {
					// Running the blockhash store feeder in backwards mode will be required to
					// resolve this.
					ll.Criticalw("Pipeline error", "err", p.err)
				} else if errors.Is(p.err, errProofVerificationFailed{}) {
					// This occurs when the proof reverts in the simulation
					// This is almost always (if not always) due to a proof generated with an out-of-date
					// blockhash
					// we can simply mark as processed and move on, since we will eventually
					// process the request with the right blockhash
					ll.Infow("proof reverted in simulation, likely stale blockhash")
					processed[p.req.req.RequestID().String()] = struct{}{}
				} else {
					ll.Errorw("Pipeline error", "err", p.err)

					if !subIsActive {
						lsn.l.Warnw("Force-fulfilling a request with insufficient funds on a cancelled sub")
						etx, err2 := lsn.enqueueForceFulfillment(ctx, p, fromAddress)
						if err2 != nil {
							ll.Errorw("Error enqueuing force-fulfillment, re-queueing request", "err", err2)
							continue
						}
						ll.Infow("Enqueued force-fulfillment", "ethTxID", etx.ID)
						processed[p.req.req.RequestID().String()] = struct{}{}

						// Need to put a continue here, otherwise the next if statement will be hit
						// and we'd break out of the loop prematurely.
						// If a sub is canceled, we want to force-fulfill ALL of it's pending requests
						// before saying we're done with it.
						continue
					}

					if startBalanceNoReserved.Cmp(p.fundsNeeded) < 0 {
						ll.Infow("Insufficient balance to fulfill a request based on estimate, returning", "err", p.err)
						return processed
					}

					// Ensure consumer is valid, otherwise drop the request.
					if !lsn.isConsumerValidAfterFinalityDepthElapsed(ctx, p.req) {
						lsn.l.Infow(
							"Dropping request that was made by an invalid consumer.",
							"consumerAddress", p.req.req.Sender(),
							"reqID", p.req.req.RequestID(),
							"blockNumber", p.req.req.Raw().BlockNumber,
							"blockHash", p.req.req.Raw().BlockHash,
						)
						processed[p.req.req.RequestID().String()] = struct{}{}
						continue
					}
				}
				continue
			}

			if startBalanceNoReserved.Cmp(p.maxFee) < 0 {
				// Insufficient funds, have to wait for a user top up. Leave it unprocessed for now
				ll.Infow("Insufficient balance to fulfill a request, returning")
				return processed
			}

			ll.Infow("Enqueuing fulfillment")
			var transaction txmgr.Tx
			err = sqlutil.TransactDataSource(ctx, lsn.ds, nil, func(tx sqlutil.DataSource) error {
				if err = lsn.pipelineRunner.InsertFinishedRun(ctx, tx, p.run, true); err != nil {
					return err
				}

				var maxLink, maxEth *string
				tmp := p.maxFee.String()
				if p.reqCommitment.NativePayment() {
					maxEth = &tmp
				} else {
					maxLink = &tmp
				}
				var (
					txMetaSubID       *uint64
					txMetaGlobalSubID *string
				)
				if lsn.coordinator.Version() == vrfcommon.V2Plus {
					txMetaGlobalSubID = ptr(p.req.req.SubID().String())
				} else if lsn.coordinator.Version() == vrfcommon.V2 {
					txMetaSubID = ptr(p.req.req.SubID().Uint64())
				}
				requestID := common.BytesToHash(p.req.req.RequestID().Bytes())
				coordinatorAddress := lsn.coordinator.Address()
				requestTxHash := p.req.req.Raw().TxHash
				transaction, err = lsn.chain.TxManager().CreateTransaction(ctx, txmgr.TxRequest{
					FromAddress:    fromAddress,
					ToAddress:      lsn.coordinator.Address(),
					EncodedPayload: hexutil.MustDecode(p.payload),
					FeeLimit:       p.gasLimit,
					Meta: &txmgr.TxMeta{
						RequestID:     &requestID,
						MaxLink:       maxLink,
						MaxEth:        maxEth,
						SubID:         txMetaSubID,
						GlobalSubID:   txMetaGlobalSubID,
						RequestTxHash: &requestTxHash,
					},
					Strategy: txmgrcommon.NewSendEveryStrategy(),
					Checker: txmgr.TransmitCheckerSpec{
						CheckerType:           lsn.transmitCheckerType(),
						VRFCoordinatorAddress: &coordinatorAddress,
						VRFRequestBlockNumber: new(big.Int).SetUint64(p.req.req.Raw().BlockNumber),
					},
				})
				return err
			})
			if err != nil {
				ll.Errorw("Error enqueuing fulfillment, requeuing request", "err", err)
				continue
			}
			ll.Infow("Enqueued fulfillment", "ethTxID", transaction.GetID())

			// If we successfully enqueued for the txm, subtract that balance
			// And loop to attempt to enqueue another fulfillment
			startBalanceNoReserved.Sub(startBalanceNoReserved, p.maxFee)
			processed[p.req.req.RequestID().String()] = struct{}{}
			vrfcommon.IncProcessedReqs(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, lsn.coordinator.Version())
		}
	}

	return
}

func (lsn *listenerV2) transmitCheckerType() txmgrtypes.TransmitCheckerType {
	if lsn.coordinator.Version() == vrfcommon.V2 {
		return txmgr.TransmitCheckerTypeVRFV2
	}
	return txmgr.TransmitCheckerTypeVRFV2Plus
}

func (lsn *listenerV2) processRequestsPerSub(
	ctx context.Context,
	subID *big.Int,
	startLinkBalance *big.Int,
	startEthBalance *big.Int,
	reqs []pendingRequest,
	subIsActive bool,
) map[string]struct{} {
	if lsn.job.VRFSpec.BatchFulfillmentEnabled && lsn.batchCoordinator != nil {
		return lsn.processRequestsPerSubBatch(ctx, subID, startLinkBalance, startEthBalance, reqs, subIsActive)
	}

	var processed = make(map[string]struct{})
	chainId := lsn.chain.Client().ConfiguredChainID()
	startBalanceNoReserveLink, err := lsn.MaybeSubtractReservedLink(
		ctx, startLinkBalance, chainId, subID, lsn.coordinator.Version())
	if err != nil {
		lsn.l.Errorw("Couldn't get reserved LINK for subscription", "sub", reqs[0].req.SubID(), "err", err)
		return processed
	}
	startBalanceNoReserveEth, err := lsn.MaybeSubtractReservedEth(
		ctx, startEthBalance, lsn.chainID, subID, lsn.coordinator.Version())
	if err != nil {
		lsn.l.Errorw("Couldn't get reserved ETH for subscription", "sub", reqs[0].req.SubID(), "err", err)
		return processed
	}

	// Split the requests into native and LINK requests.
	var (
		nativeRequests []pendingRequest
		linkRequests   []pendingRequest
	)
	for _, req := range reqs {
		if req.req.NativePayment() {
			if !lsn.inflightCache.Contains(req.req.Raw()) {
				nativeRequests = append(nativeRequests, req)
			} else {
				lsn.l.Debugw("Skipping native request because it is already inflight",
					"reqID", req.req.RequestID())
			}
		} else {
			if !lsn.inflightCache.Contains(req.req.Raw()) {
				linkRequests = append(linkRequests, req)
			} else {
				lsn.l.Debugw("Skipping link request because it is already inflight",
					"reqID", req.req.RequestID())
			}
		}
	}
	// process the native and link requests in parallel
	var (
		wg                             sync.WaitGroup
		nativeProcessed, linkProcessed map[string]struct{}
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		nativeProcessed = lsn.processRequestsPerSubHelper(
			ctx,
			subID,
			startEthBalance,
			startBalanceNoReserveEth,
			nativeRequests,
			subIsActive,
			true)
	}()
	go func() {
		defer wg.Done()
		linkProcessed = lsn.processRequestsPerSubHelper(
			ctx,
			subID,
			startLinkBalance,
			startBalanceNoReserveLink,
			linkRequests,
			subIsActive,
			false)
	}()
	wg.Wait()
	// combine the native and link processed requests into the processed map
	for k, v := range nativeProcessed {
		processed[k] = v
	}
	for k, v := range linkProcessed {
		processed[k] = v
	}

	return processed
}

func (lsn *listenerV2) requestCommitmentPayload(requestID *big.Int) (payload []byte, err error) {
	if lsn.coordinator.Version() == vrfcommon.V2Plus {
		return coordinatorV2PlusABI.Pack("s_requestCommitments", requestID)
	} else if lsn.coordinator.Version() == vrfcommon.V2 {
		return coordinatorV2ABI.Pack("getCommitment", requestID)
	}
	return nil, errors.Errorf("unsupported coordinator version: %s", lsn.coordinator.Version())
}

// checkReqsFulfilled returns a bool slice the same size of the given reqs slice
// where each slice element indicates whether that request was already fulfilled
// or not.
func (lsn *listenerV2) checkReqsFulfilled(ctx context.Context, l logger.Logger, reqs []pendingRequest) ([]bool, error) {
	var (
		start     = time.Now()
		calls     = make([]rpc.BatchElem, len(reqs))
		fulfilled = make([]bool, len(reqs))
	)

	for i, req := range reqs {
		payload, err := lsn.requestCommitmentPayload(req.req.RequestID())
		if err != nil {
			// This shouldn't happen
			return fulfilled, fmt.Errorf("creating getCommitment payload: %w", err)
		}

		reqBlockNumber := new(big.Int).SetUint64(req.req.Raw().BlockNumber)

		// Subtract 5 since the newest block likely isn't indexed yet and will cause "header not
		// found" errors.
		currBlock := new(big.Int).SetUint64(lsn.getLatestHead() - 5)
		m := bigmath.Max(reqBlockNumber, currBlock)

		var result string
		calls[i] = rpc.BatchElem{
			Method: "eth_call",
			Args: []interface{}{
				map[string]interface{}{
					"to":   lsn.coordinator.Address(),
					"data": hexutil.Bytes(payload),
				},
				// The block at which we want to make the call
				hexutil.EncodeBig(m),
			},
			Result: &result,
		}
	}

	err := lsn.chain.Client().BatchCallContext(ctx, calls)
	if err != nil {
		return fulfilled, fmt.Errorf("making batch call: %w", err)
	}

	var errs error
	for i, call := range calls {
		if call.Error != nil {
			errs = multierr.Append(errs, fmt.Errorf("checking request %s with hash %s: %w",
				reqs[i].req.RequestID().String(), reqs[i].req.Raw().TxHash.String(), call.Error))
			continue
		}

		rString, ok := call.Result.(*string)
		if !ok {
			errs = multierr.Append(errs,
				fmt.Errorf("unexpected result %+v on request %s with hash %s",
					call.Result, reqs[i].req.RequestID().String(), reqs[i].req.Raw().TxHash.String()))
			continue
		}
		result, err := hexutil.Decode(*rString)
		if err != nil {
			errs = multierr.Append(errs,
				fmt.Errorf("decoding batch call result %+v %s request %s with hash %s: %w",
					call.Result, *rString, reqs[i].req.RequestID().String(), reqs[i].req.Raw().TxHash.String(), err))
			continue
		}

		if utils.IsEmpty(result) {
			l.Infow("Request already fulfilled",
				"reqID", reqs[i].req.RequestID().String(),
				"attempts", reqs[i].attempts,
				"txHash", reqs[i].req.Raw().TxHash)
			fulfilled[i] = true
		}
	}

	l.Debugw("Done checking fulfillment status",
		"numChecked", len(reqs), "time", time.Since(start).String())
	return fulfilled, errs
}

func (lsn *listenerV2) runPipelines(
	ctx context.Context,
	l logger.Logger,
	maxGasPriceWei *assets.Wei,
	reqs []pendingRequest,
) []vrfPipelineResult {
	var (
		start   = time.Now()
		results = make([]vrfPipelineResult, len(reqs))
		wg      = sync.WaitGroup{}
	)

	for i, req := range reqs {
		wg.Add(1)
		go func(i int, req pendingRequest) {
			defer wg.Done()
			ll := l.With("reqID", req.req.RequestID().String())
			results[i] = lsn.simulateFulfillment(ctx, maxGasPriceWei, req, ll)
		}(i, req)
	}
	wg.Wait()

	l.Debugw("Finished running pipelines",
		"count", len(reqs), "time", time.Since(start).String())
	return results
}

func (lsn *listenerV2) estimateFee(
	ctx context.Context,
	req RandomWordsRequested,
	maxGasPriceWei *assets.Wei,
) (*big.Int, error) {
	// NativePayment() returns true if and only if the version is V2+ and the
	// request was made in ETH.
	if req.NativePayment() {
		return EstimateFeeWei(req.CallbackGasLimit(), maxGasPriceWei.ToInt())
	}

	// In the event we are using LINK we need to estimate the fee in juels
	// Don't use up too much time to get this info, it's not critical for operating vrf.
	callCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	roundData, err := lsn.aggregator.LatestRoundData(&bind.CallOpts{Context: callCtx})
	if err != nil {
		return nil, fmt.Errorf("get aggregator latestAnswer: %w", err)
	}

	return EstimateFeeJuels(
		req.CallbackGasLimit(),
		maxGasPriceWei.ToInt(),
		roundData.Answer,
	)
}

// Here we use the pipeline to parse the log, generate a vrf response
// then simulate the transaction at the max gas price to determine its maximum link cost.
func (lsn *listenerV2) simulateFulfillment(
	ctx context.Context,
	maxGasPriceWei *assets.Wei,
	req pendingRequest,
	lg logger.Logger,
) vrfPipelineResult {
	var (
		res = vrfPipelineResult{req: req}
		err error
	)
	// estimate how much funds are needed so that we can log it if the simulation fails.
	res.fundsNeeded, err = lsn.estimateFee(ctx, req.req, maxGasPriceWei)
	if err != nil {
		// not critical, just log and continue
		lg.Warnw("unable to estimate funds needed for request, continuing anyway",
			"err", err)
		res.fundsNeeded = big.NewInt(0)
	}

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"databaseID":    lsn.job.ID,
			"externalJobID": lsn.job.ExternalJobID,
			"name":          lsn.job.Name.ValueOrZero(),
			"publicKey":     lsn.job.VRFSpec.PublicKey[:],
			"maxGasPrice":   maxGasPriceWei.ToInt().String(),
			"evmChainID":    lsn.job.VRFSpec.EVMChainID.String(),
		},
		"jobRun": map[string]interface{}{
			"logBlockHash":   req.req.Raw().BlockHash.Bytes(),
			"logBlockNumber": req.req.Raw().BlockNumber,
			"logTxHash":      req.req.Raw().TxHash,
			"logTopics":      req.req.Raw().Topics,
			"logData":        req.req.Raw().Data,
		},
	})
	var trrs pipeline.TaskRunResults
	res.run, trrs, err = lsn.pipelineRunner.ExecuteRun(ctx, *lsn.job.PipelineSpec, vars)
	if err != nil {
		res.err = fmt.Errorf("executing run: %w", err)
		return res
	}
	// The call task will fail if there are insufficient funds
	if res.run.AllErrors.HasError() {
		res.err = errors.WithStack(res.run.AllErrors.ToError())

		if strings.Contains(res.err.Error(), "blockhash not found in store") {
			res.err = multierr.Combine(res.err, errBlockhashNotInStore{})
		} else if isProofVerificationError(res.err.Error()) {
			res.err = multierr.Combine(res.err, errProofVerificationFailed{})
		} else if strings.Contains(res.err.Error(), "execution reverted") {
			// Even if the simulation fails, we want to get the
			// txData for the fulfillRandomWords call, in case
			// we need to force fulfill.
			for _, trr := range trrs {
				if trr.Task.Type() == pipeline.TaskTypeVRFV2 {
					if trr.Result.Error != nil {
						// error in VRF proof generation
						// this means that we won't be able to force-fulfill in the event of a
						// canceled sub and active requests.
						// since this would be an extraordinary situation,
						// we can log loudly here.
						lg.Criticalw("failed to generate VRF proof", "err", trr.Result.Error)
						break
					}

					// extract the abi-encoded tx data to fulfillRandomWords from the VRF task.
					// that's all we need in the event of a force-fulfillment.
					m := trr.Result.Value.(map[string]any)
					res.payload = m["output"].(string)
					res.proof = FromV2Proof(m["proof"].(vrf_coordinator_v2.VRFProof))
					res.reqCommitment = NewRequestCommitment(m["requestCommitment"])
				}
			}
			res.err = multierr.Combine(res.err, errPossiblyInsufficientFunds{})
		}

		return res
	}
	finalResult := trrs.FinalResult()
	if len(finalResult.Values) != 1 {
		res.err = errors.Errorf("unexpected number of outputs, expected 1, was %d", len(finalResult.Values))
		return res
	}

	// Run succeeded, we expect a byte array representing the billing amount
	b, ok := finalResult.Values[0].([]uint8)
	if !ok {
		res.err = errors.New("expected []uint8 final result")
		return res
	}

	res.maxFee, err = hex.ParseBig(hexutil.Encode(b)[2:])
	if err != nil {
		res.err = err
		return res
	}

	for _, trr := range trrs {
		if trr.Task.Type() == pipeline.TaskTypeVRFV2 {
			m := trr.Result.Value.(map[string]interface{})
			res.payload = m["output"].(string)
			res.proof = FromV2Proof(m["proof"].(vrf_coordinator_v2.VRFProof))
			res.reqCommitment = NewRequestCommitment(m["requestCommitment"])
		}

		if trr.Task.Type() == pipeline.TaskTypeVRFV2Plus {
			m := trr.Result.Value.(map[string]interface{})
			res.payload = m["output"].(string)
			res.proof = FromV2PlusProof(m["proof"].(vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalProof))
			res.reqCommitment = NewRequestCommitment(m["requestCommitment"])
		}

		if trr.Task.Type() == pipeline.TaskTypeEstimateGasLimit {
			res.gasLimit = trr.Result.Value.(uint64)
		}
	}
	return res
}

func (lsn *listenerV2) fromAddresses() []common.Address {
	var addresses []common.Address
	for _, a := range lsn.job.VRFSpec.FromAddresses {
		addresses = append(addresses, a.Address())
	}
	return addresses
}
