package vrf

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	heaps "github.com/theodesp/go-heaps"
	"github.com/theodesp/go-heaps/pairing"
	"go.uber.org/multierr"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/batch_vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
	bigmath "github.com/smartcontractkit/chainlink/core/utils/big_math"
)

var (
	_                     log.Listener   = &listenerV2{}
	_                     job.ServiceCtx = &listenerV2{}
	coordinatorV2ABI                     = evmtypes.MustGetABI(vrf_coordinator_v2.VRFCoordinatorV2ABI)
	batchCoordinatorV2ABI                = evmtypes.MustGetABI(batch_vrf_coordinator_v2.BatchVRFCoordinatorV2ABI)
)

const (
	// GasAfterPaymentCalculation is the gas used after computing the payment
	GasAfterPaymentCalculation = 21000 + // base cost of the transaction
		100 + 5000 + // warm subscription balance read and update. See https://eips.ethereum.org/EIPS/eip-2929
		2*2100 + 20000 - // cold read oracle address and oracle balance and first time oracle balance update, note first time will be 20k, but 5k subsequently
		4800 + // request delete refund (refunds happen after execution), note pre-london fork was 15k. See https://eips.ethereum.org/EIPS/eip-3529
		6685 // Positive static costs of argument encoding etc. note that it varies by +/- x*12 for every x bytes of non-zero data in the proof.

	// BatchFulfillmentIterationGasCost is the cost of a single iteration of the batch coordinator's
	// loop. This is used to determine the gas allowance for a batch fulfillment call.
	BatchFulfillmentIterationGasCost = 52_000

	// backoffFactor is the factor by which to increase the delay each time a request fails.
	backoffFactor = 1.3
)

type errPossiblyInsufficientFunds struct{}

func (errPossiblyInsufficientFunds) Error() string {
	return "Simulation errored, possibly insufficient funds. Request will remain unprocessed until funds are available"
}

type errBlockhashNotInStore struct{}

func (errBlockhashNotInStore) Error() string {
	return "Blockhash not in store"
}

func newListenerV2(
	cfg Config,
	l logger.Logger,
	ethClient evmclient.Client,
	chainID *big.Int,
	logBroadcaster log.Broadcaster,
	q pg.Q,
	coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
	batchCoordinator batch_vrf_coordinator_v2.BatchVRFCoordinatorV2Interface,
	aggregator *aggregator_v3_interface.AggregatorV3Interface,
	txm txmgr.TxManager,
	pipelineRunner pipeline.Runner,
	gethks keystore.Eth,
	job job.Job,
	mailMon *utils.MailboxMonitor,
	reqLogs *utils.Mailbox[log.Broadcast],
	reqAdded func(),
	respCount map[string]uint64,
	headBroadcaster httypes.HeadBroadcasterRegistry,
	deduper *logDeduper,
) *listenerV2 {
	return &listenerV2{
		cfg:                cfg,
		l:                  logger.Sugared(l),
		ethClient:          ethClient,
		chainID:            chainID,
		logBroadcaster:     logBroadcaster,
		txm:                txm,
		mailMon:            mailMon,
		coordinator:        coordinator,
		batchCoordinator:   batchCoordinator,
		pipelineRunner:     pipelineRunner,
		job:                job,
		q:                  q,
		gethks:             gethks,
		reqLogs:            reqLogs,
		chStop:             make(chan struct{}),
		reqAdded:           reqAdded,
		respCount:          respCount,
		blockNumberToReqID: pairing.New(),
		headBroadcaster:    headBroadcaster,
		latestHeadMu:       sync.RWMutex{},
		wg:                 &sync.WaitGroup{},
		aggregator:         aggregator,
		deduper:            deduper,
	}
}

type pendingRequest struct {
	confirmedAtBlock uint64
	req              *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested
	lb               log.Broadcast
	utcTimestamp     time.Time

	// used for exponential backoff when retrying
	attempts int
	lastTry  time.Time
}

type vrfPipelineResult struct {
	err           error
	maxLink       *big.Int
	juelsNeeded   *big.Int
	run           pipeline.Run
	payload       string
	gasLimit      uint32
	req           pendingRequest
	proof         vrf_coordinator_v2.VRFProof
	reqCommitment vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment
}

type listenerV2 struct {
	utils.StartStopOnce
	cfg            Config
	l              logger.SugaredLogger
	ethClient      evmclient.Client
	chainID        *big.Int
	logBroadcaster log.Broadcaster
	txm            txmgr.TxManager
	mailMon        *utils.MailboxMonitor

	coordinator      vrf_coordinator_v2.VRFCoordinatorV2Interface
	batchCoordinator batch_vrf_coordinator_v2.BatchVRFCoordinatorV2Interface

	pipelineRunner pipeline.Runner
	job            job.Job
	q              pg.Q
	gethks         keystore.Eth
	reqLogs        *utils.Mailbox[log.Broadcast]
	chStop         chan struct{}
	// We can keep these pending logs in memory because we
	// only mark them confirmed once we send a corresponding fulfillment transaction.
	// So on node restart in the middle of processing, the lb will resend them.
	reqsMu   sync.Mutex // Both the log listener and the request handler write to reqs
	reqs     []pendingRequest
	reqAdded func() // A simple debug helper

	// Data structures for reorg attack protection
	// We want a map so we can do an O(1) count update every fulfillment log we get.
	respCountMu sync.Mutex
	respCount   map[string]uint64
	// This auxiliary heap is used when we need to purge the
	// respCount map - we repeatedly want to remove the minimum log.
	// You could use a sorted list if the completed logs arrive in order, but they may not.
	blockNumberToReqID *pairing.PairHeap

	// head tracking data structures
	headBroadcaster  httypes.HeadBroadcasterRegistry
	latestHeadMu     sync.RWMutex
	latestHeadNumber uint64

	// Wait group to wait on all goroutines to shut down.
	wg *sync.WaitGroup

	// aggregator client to get link/eth feed prices from chain.
	aggregator aggregator_v3_interface.AggregatorV3InterfaceInterface

	// deduper prevents processing duplicate requests from the log broadcaster.
	deduper *logDeduper
}

// Start starts listenerV2.
func (lsn *listenerV2) Start(ctx context.Context) error {
	return lsn.StartOnce("VRFListenerV2", func() error {
		// Check gas limit configuration
		confCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		conf, err := lsn.coordinator.GetConfig(&bind.CallOpts{Context: confCtx})
		gasLimit := lsn.cfg.EvmGasLimitDefault()
		if lsn.cfg.EvmGasLimitVRFJobType() != nil {
			gasLimit = *lsn.cfg.EvmGasLimitVRFJobType()
		}
		if err != nil {
			lsn.l.Criticalw("Error getting coordinator config for gas limit check, starting anyway.", "err", err)
		} else if conf.MaxGasLimit+(GasProofVerification*2) > uint32(gasLimit) {
			lsn.l.Criticalw("Node gas limit setting may not be high enough to fulfill all requests; it should be increased. Starting anyway.",
				"currentGasLimit", gasLimit,
				"neededGasLimit", conf.MaxGasLimit+(GasProofVerification*2),
				"callbackGasLimit", conf.MaxGasLimit,
				"proofVerificationGas", GasProofVerification)
		}

		spec := job.LoadEnvConfigVarsVRF(lsn.cfg, *lsn.job.VRFSpec)

		unsubscribeLogs := lsn.logBroadcaster.Register(lsn, log.ListenerOpts{
			Contract: lsn.coordinator.Address(),
			ParseLog: lsn.coordinator.ParseLog,
			LogsWithTopics: map[common.Hash][][]log.Topic{
				vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{}.Topic(): {
					{
						log.Topic(spec.PublicKey.MustHash()),
					},
				},
			},
			// Specify a min incoming confirmations of 1 so that we can receive a request log
			// right away. We set the real number of confirmations on a per-request basis in
			// the getConfirmedAt method.
			MinIncomingConfirmations: 1,
			ReplayStartedCallback:    lsn.ReplayStartedCallback,
		})

		latestHead, unsubscribeHeadBroadcaster := lsn.headBroadcaster.Subscribe(lsn)
		if latestHead != nil {
			lsn.setLatestHead(latestHead)
		}

		// Log listener gathers request logs
		lsn.wg.Add(1)
		go func() {
			lsn.runLogListener([]func(){unsubscribeLogs, unsubscribeHeadBroadcaster}, spec.MinIncomingConfirmations, lsn.wg)
		}()

		// Request handler periodically computes a set of logs which can be fulfilled.
		lsn.wg.Add(1)
		go func() {
			lsn.runRequestHandler(spec.PollPeriod, lsn.wg)
		}()

		lsn.mailMon.Monitor(lsn.reqLogs, "VRFListenerV2", "RequestLogs", fmt.Sprint(lsn.job.ID))
		return nil
	})
}

func (lsn *listenerV2) setLatestHead(head *evmtypes.Head) {
	lsn.latestHeadMu.Lock()
	defer lsn.latestHeadMu.Unlock()
	num := uint64(head.Number)
	if num > lsn.latestHeadNumber {
		lsn.latestHeadNumber = num
	}
}

// OnNewLongestChain is called by the head broadcaster when a new head is available.
func (lsn *listenerV2) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
	lsn.setLatestHead(head)
}

func (lsn *listenerV2) getLatestHead() uint64 {
	lsn.latestHeadMu.RLock()
	defer lsn.latestHeadMu.RUnlock()
	return uint64(lsn.latestHeadNumber)
}

// Returns all the confirmed logs from
// the pending queue by subscription
func (lsn *listenerV2) getAndRemoveConfirmedLogsBySub(latestHead uint64) map[uint64][]pendingRequest {
	lsn.reqsMu.Lock()
	defer lsn.reqsMu.Unlock()
	updateQueueSize(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, v2, uniqueReqs(lsn.reqs))
	var toProcess = make(map[uint64][]pendingRequest)
	var toKeep []pendingRequest
	for i := 0; i < len(lsn.reqs); i++ {
		if r := lsn.reqs[i]; lsn.ready(r, latestHead) {
			toProcess[r.req.SubId] = append(toProcess[r.req.SubId], r)
		} else {
			toKeep = append(toKeep, lsn.reqs[i])
		}
	}
	lsn.reqs = toKeep
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
	lsn.respCountMu.Lock()
	defer lsn.respCountMu.Unlock()
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
// Note we have to consider the pending reqs already in the txm as already "spent" link,
// using a max link consumed in their metadata.
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
func (lsn *listenerV2) processPendingVRFRequests(ctx context.Context) {
	confirmed := lsn.getAndRemoveConfirmedLogsBySub(lsn.getLatestHead())
	processed := make(map[string]struct{})
	start := time.Now()

	// Add any unprocessed requests back to lsn.reqs after request processing is complete.
	defer func() {
		var toKeep []pendingRequest
		for _, subReqs := range confirmed {
			for _, req := range subReqs {
				if _, ok := processed[req.req.RequestId.String()]; !ok {
					req.attempts++
					req.lastTry = time.Now().UTC()
					toKeep = append(toKeep, req)
					if lsn.job.VRFSpec.BackoffInitialDelay != 0 {
						lsn.l.Infow("Request failed, next retry will be delayed.",
							"reqID", req.req.RequestId.String(),
							"subID", req.req.SubId,
							"attempts", req.attempts,
							"lastTry", req.lastTry.String(),
							"nextTry", nextTry(
								req.attempts,
								lsn.job.VRFSpec.BackoffInitialDelay,
								lsn.job.VRFSpec.BackoffMaxDelay,
								req.lastTry))
					}
				} else {
					lsn.markLogAsConsumed(req.lb)
				}
			}
		}
		// There could be logs accumulated to this slice while request processor is running,
		// so we merged the new ones with the ones that need to be requeued.
		lsn.reqsMu.Lock()
		lsn.reqs = append(lsn.reqs, toKeep...)
		lsn.l.Infow("Finished processing pending requests",
			"totalProcessed", len(processed),
			"totalFailed", len(toKeep),
			"total", len(lsn.reqs),
			"time", time.Since(start).String())
		lsn.reqsMu.Unlock() // unlock here since len(lsn.reqs) is a read, to avoid a data race.
	}()

	// Get subscription balance. Note that outside of this request handler, this can only decrease while there
	// are no pending requests
	if len(confirmed) == 0 {
		lsn.l.Infow("No pending requests ready for processing")
		return
	}
	for subID, reqs := range confirmed {
		sub, err := lsn.coordinator.GetSubscription(&bind.CallOpts{
			Context: ctx,
		}, subID)

		if err != nil {
			if strings.Contains(err.Error(), "execution reverted") {
				lsn.l.Warnw("Subscription not found", "subID", subID, "err", err)
				for _, req := range reqs {
					lsn.l.Infow("Skipping requests without valid subscription", "subID", subID, "reqID", req.req.RequestId)
					processed[req.req.RequestId.String()] = struct{}{}
				}
			} else {
				lsn.l.Errorw("Unable to read subscription balance", "subID", subID, "err", err)
			}
			continue
		}

		// Sort requests in ascending order by CallbackGasLimit
		// so that we process the "cheapest" requests for each subscription
		// first. This allows us to break out of the processing loop as early as possible
		// in the event that a subscription is too underfunded to have it's
		// requests processed.
		slices.SortFunc(reqs, func(a, b pendingRequest) bool {
			return a.req.CallbackGasLimit < b.req.CallbackGasLimit
		})

		startBalance := sub.Balance
		p := lsn.processRequestsPerSub(ctx, subID, startBalance, reqs)
		for reqID := range p {
			processed[reqID] = struct{}{}
		}
	}
	lsn.pruneConfirmedRequestCounts()
}

// MaybeSubtractReservedLink figures out how much LINK is reserved for other VRF requests that
// have not been fully confirmed yet on-chain, and subtracts that from the given startBalance,
// and returns that value if there are no errors.
func MaybeSubtractReservedLink(q pg.Q, startBalance *big.Int, chainID, subID uint64) (*big.Int, error) {
	var reservedLink string
	err := q.Get(&reservedLink, `SELECT SUM(CAST(meta->>'MaxLink' AS NUMERIC(78, 0)))
				   FROM eth_txes
				   WHERE meta->>'MaxLink' IS NOT NULL
				   AND evm_chain_id = $1
				   AND CAST(meta->>'SubId' AS NUMERIC) = $2
				   AND state IN ('unconfirmed', 'unstarted', 'in_progress')
				   GROUP BY meta->>'SubId'`, chainID, subID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "getting reserved LINK")
	}

	if reservedLink != "" {
		reservedLinkInt, success := big.NewInt(0).SetString(reservedLink, 10)
		if !success {
			return nil, fmt.Errorf("converting reserved LINK %s", reservedLink)
		}

		return new(big.Int).Sub(startBalance, reservedLinkInt), nil
	}

	return new(big.Int).Set(startBalance), nil
}

type fulfilledReqV2 struct {
	blockNumber uint64
	reqID       string
}

func (a fulfilledReqV2) Compare(b heaps.Item) int {
	a1 := a
	a2 := b.(fulfilledReqV2)
	switch {
	case a1.blockNumber > a2.blockNumber:
		return 1
	case a1.blockNumber < a2.blockNumber:
		return -1
	default:
		return 0
	}
}

func (lsn *listenerV2) processRequestsPerSubBatch(
	ctx context.Context,
	subID uint64,
	startBalance *big.Int,
	reqs []pendingRequest,
) map[string]struct{} {
	start := time.Now()
	var processed = make(map[string]struct{})
	startBalanceNoReserveLink, err := MaybeSubtractReservedLink(
		lsn.q, startBalance, lsn.chainID.Uint64(), subID)
	if err != nil {
		lsn.l.Errorw("Couldn't get reserved LINK for subscription", "sub", reqs[0].req.SubId, "err", err)
		return processed
	}

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
	batchMaxGas := uint32(config.MaxGasLimit + 400_000)

	l := lsn.l.With(
		"subID", reqs[0].req.SubId,
		"eligibleSubReqs", len(reqs),
		"startBalance", startBalance.String(),
		"startBalanceNoReservedLink", startBalanceNoReserveLink.String(),
		"batchMaxGas", batchMaxGas,
	)

	defer func() {
		l.Infow("Finished processing for sub",
			"endBalance", startBalanceNoReserveLink.String(),
			"totalProcessed", len(processed),
			"totalUnique", uniqueReqs(reqs),
			"time", time.Since(start).String())
	}()

	l.Infow("Processing requests for subscription with batching")

	// Check for already consumed or expired reqs
	unconsumed, processedReqs := lsn.getUnconsumed(l, reqs)
	for _, reqID := range processedReqs {
		processed[reqID] = struct{}{}
	}

	// Process requests in chunks in order to kick off as many jobs
	// as configured in parallel. Then we can combine into fulfillment
	// batches afterwards.
	for chunkStart := 0; chunkStart < len(unconsumed); chunkStart += int(lsn.job.VRFSpec.ChunkSize) {
		chunkEnd := chunkStart + int(lsn.job.VRFSpec.ChunkSize)
		if chunkEnd > len(unconsumed) {
			chunkEnd = len(unconsumed)
		}
		chunk := unconsumed[chunkStart:chunkEnd]

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
				processed[chunk[i].req.RequestId.String()] = struct{}{}
			} else {
				unfulfilled = append(unfulfilled, chunk[i])
			}
		}

		// All fromAddresses passed to the VRFv2 job have the same KeySpecificMaxGasPriceWei value.
		fromAddresses := lsn.fromAddresses()
		maxGasPriceWei := lsn.cfg.KeySpecificMaxGasPriceWei(fromAddresses[0])

		// Cases:
		// 1. Never simulated: in this case, we want to observe the time until simulated
		// on the utcTimestamp field of the pending request.
		// 2. Simulated before: in this case, lastTry will be set to a non-zero time value,
		// in which case we'd want to use that as a relative point from when we last tried
		// the request.
		observeRequestSimDuration(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, v2, unfulfilled)

		pipelines := lsn.runPipelines(ctx, l, maxGasPriceWei, unfulfilled)
		batches := newBatchFulfillments(batchMaxGas)
		outOfBalance := false
		for _, p := range pipelines {
			ll := l.With("reqID", p.req.req.RequestId.String(),
				"txHash", p.req.req.Raw.TxHash,
				"maxGasPrice", maxGasPriceWei.String(),
				"juelsNeeded", p.juelsNeeded.String(),
				"maxLink", p.maxLink.String(),
				"gasLimit", p.gasLimit,
				"attempts", p.req.attempts,
				"remainingBalance", startBalanceNoReserveLink.String(),
			)

			if p.err != nil {
				if startBalanceNoReserveLink.Cmp(p.juelsNeeded) < 0 && errors.Is(p.err, errPossiblyInsufficientFunds{}) {
					ll.Infow("Insufficient link balance to fulfill a request based on estimate, breaking", "err", p.err)
					outOfBalance = true

					// break out of this inner loop to process the currently constructed batch
					break
				}

				if errors.Is(p.err, errBlockhashNotInStore{}) {
					// Running the blockhash store feeder in backwards mode will be required to
					// resolve this.
					ll.Criticalw("Pipeline error", "err", p.err)
				} else {
					ll.Errorw("Pipeline error", "err", p.err)
					// Ensure consumer is valid, otherwise drop the request.
					if !lsn.isConsumerValidAfterFinalityDepthElapsed(ctx, p.req) {
						lsn.l.Infow(
							"Dropping request that was made by an invalid consumer.",
							"consumerAddress", p.req.req.Sender,
							"reqID", p.req.req.RequestId,
							"blockNumber", p.req.req.Raw.BlockNumber,
							"blockHash", p.req.req.Raw.BlockHash,
						)
						lsn.markLogAsConsumed(p.req.lb)
					}
				}
				continue
			}

			if startBalanceNoReserveLink.Cmp(p.maxLink) < 0 {
				// Insufficient funds, have to wait for a user top up.
				// Break out of the loop now and process what we are able to process
				// in the constructed batches.
				ll.Infow("Insufficient link balance to fulfill a request, breaking")
				break
			}

			batches.addRun(p)

			startBalanceNoReserveLink.Sub(startBalanceNoReserveLink, p.maxLink)
		}

		var processedRequestIDs []string
		for _, batch := range batches.fulfillments {
			l.Debugw("Processing batch", "batchSize", len(batch.proofs))
			p := lsn.processBatch(l, subID, startBalanceNoReserveLink, batchMaxGas, batch)
			processedRequestIDs = append(processedRequestIDs, p...)
		}

		for _, reqID := range processedRequestIDs {
			processed[reqID] = struct{}{}
		}

		// outOfBalance is set to true if the current sub we are processing
		// has ran out of funds to process any remaining requests. After enqueueing
		// this constructed batch, we break out of this outer loop in order to
		// avoid unnecessarily processing the remaining requests.
		if outOfBalance {
			break
		}
	}

	return processed
}

// For an errored pipeline run, wait until the finality depth of the chain to have elapsed,
// then check if the failing request is being called by an invalid sender. Return false if this is the case,
// otherwise true.
func (lsn *listenerV2) isConsumerValidAfterFinalityDepthElapsed(ctx context.Context, req pendingRequest) bool {
	latestHead := lsn.getLatestHead()
	if latestHead-req.req.Raw.BlockNumber > uint64(lsn.cfg.EvmFinalityDepth()) {
		code, err := lsn.ethClient.CodeAt(ctx, req.req.Sender, big.NewInt(int64(latestHead)))
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

func (lsn *listenerV2) processRequestsPerSub(
	ctx context.Context,
	subID uint64,
	startBalance *big.Int,
	reqs []pendingRequest,
) map[string]struct{} {
	if lsn.job.VRFSpec.BatchFulfillmentEnabled && lsn.batchCoordinator != nil {
		return lsn.processRequestsPerSubBatch(ctx, subID, startBalance, reqs)
	}

	start := time.Now()
	var processed = make(map[string]struct{})
	startBalanceNoReserveLink, err := MaybeSubtractReservedLink(
		lsn.q, startBalance, lsn.ethClient.ChainID().Uint64(), subID)
	if err != nil {
		lsn.l.Errorw("Couldn't get reserved LINK for subscription", "sub", reqs[0].req.SubId, "err", err)
		return processed
	}

	l := lsn.l.With(
		"subID", reqs[0].req.SubId,
		"eligibleSubReqs", len(reqs),
		"startBalance", startBalance.String(),
		"startBalanceNoReservedLink", startBalanceNoReserveLink.String(),
	)

	defer func() {
		l.Infow("Finished processing for sub",
			"endBalance", startBalanceNoReserveLink.String(),
			"totalProcessed", len(processed),
			"totalUnique", uniqueReqs(reqs),
			"time", time.Since(start).String())
	}()

	l.Infow("Processing requests for subscription")

	// Check for already consumed or expired reqs
	unconsumed, processedReqs := lsn.getUnconsumed(l, reqs)
	for _, reqID := range processedReqs {
		processed[reqID] = struct{}{}
	}

	// Process requests in chunks
	for chunkStart := 0; chunkStart < len(unconsumed); chunkStart += int(lsn.job.VRFSpec.ChunkSize) {
		chunkEnd := chunkStart + int(lsn.job.VRFSpec.ChunkSize)
		if chunkEnd > len(unconsumed) {
			chunkEnd = len(unconsumed)
		}
		chunk := unconsumed[chunkStart:chunkEnd]

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
				processed[chunk[i].req.RequestId.String()] = struct{}{}
			} else {
				unfulfilled = append(unfulfilled, chunk[i])
			}
		}

		// All fromAddresses passed to the VRFv2 job have the same KeySpecificMaxGasPriceWei value.
		fromAddresses := lsn.fromAddresses()
		maxGasPriceWei := lsn.cfg.KeySpecificMaxGasPriceWei(fromAddresses[0])
		observeRequestSimDuration(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, v2, unfulfilled)
		pipelines := lsn.runPipelines(ctx, l, maxGasPriceWei, unfulfilled)
		for _, p := range pipelines {
			ll := l.With("reqID", p.req.req.RequestId.String(),
				"txHash", p.req.req.Raw.TxHash,
				"maxGasPrice", maxGasPriceWei.String(),
				"juelsNeeded", p.juelsNeeded.String(),
				"maxLink", p.maxLink.String(),
				"gasLimit", p.gasLimit,
				"attempts", p.req.attempts,
				"remainingBalance", startBalanceNoReserveLink.String(),
			)

			if p.err != nil {
				if startBalanceNoReserveLink.Cmp(p.juelsNeeded) < 0 && errors.Is(p.err, errPossiblyInsufficientFunds{}) {
					ll.Infow("Insufficient link balance to fulfill a request based on estimate, returning", "err", p.err)
					return processed
				}

				if errors.Is(p.err, errBlockhashNotInStore{}) {
					// Running the blockhash store feeder in backwards mode will be required to
					// resolve this.
					ll.Criticalw("Pipeline error", "err", p.err)
				} else {
					ll.Errorw("Pipeline error", "err", p.err)

					// Ensure consumer is valid, otherwise drop the request.
					if !lsn.isConsumerValidAfterFinalityDepthElapsed(ctx, p.req) {
						lsn.l.Infow(
							"Dropping request that was made by an invalid consumer.",
							"consumerAddress", p.req.req.Sender,
							"reqID", p.req.req.RequestId,
							"blockNumber", p.req.req.Raw.BlockNumber,
							"blockHash", p.req.req.Raw.BlockHash,
						)
						lsn.markLogAsConsumed(p.req.lb)
					}
				}
				continue
			}

			if startBalanceNoReserveLink.Cmp(p.maxLink) < 0 {
				// Insufficient funds, have to wait for a user top up. Leave it unprocessed for now
				ll.Infow("Insufficient link balance to fulfill a request, returning")
				return processed
			}

			fromAddress, err := lsn.gethks.GetRoundRobinAddress(lsn.chainID, fromAddresses...)
			if err != nil {
				l.Errorw("Couldn't get next from address", "err", err)
				continue
			}
			ll = ll.With("fromAddress", fromAddress)

			ll.Infow("Enqueuing fulfillment")
			var ethTX txmgr.EthTx
			err = lsn.q.Transaction(func(tx pg.Queryer) error {
				if err = lsn.pipelineRunner.InsertFinishedRun(&p.run, true, pg.WithQueryer(tx)); err != nil {
					return err
				}
				if err = lsn.logBroadcaster.MarkConsumed(p.req.lb, pg.WithQueryer(tx)); err != nil {
					return err
				}

				maxLinkString := p.maxLink.String()
				requestID := common.BytesToHash(p.req.req.RequestId.Bytes())
				coordinatorAddress := lsn.coordinator.Address()
				ethTX, err = lsn.txm.CreateEthTransaction(txmgr.NewTx{
					FromAddress:    fromAddress,
					ToAddress:      lsn.coordinator.Address(),
					EncodedPayload: hexutil.MustDecode(p.payload),
					GasLimit:       p.gasLimit,
					Meta: &txmgr.EthTxMeta{
						RequestID:     &requestID,
						MaxLink:       &maxLinkString,
						SubID:         &p.req.req.SubId,
						RequestTxHash: &p.req.req.Raw.TxHash,
					},
					Strategy: txmgr.NewSendEveryStrategy(),
					Checker: txmgr.TransmitCheckerSpec{
						CheckerType:           txmgr.TransmitCheckerTypeVRFV2,
						VRFCoordinatorAddress: &coordinatorAddress,
						VRFRequestBlockNumber: new(big.Int).SetUint64(p.req.req.Raw.BlockNumber),
					},
				}, pg.WithQueryer(tx), pg.WithParentCtx(ctx))
				return err
			})
			if err != nil {
				ll.Errorw("Error enqueuing fulfillment, requeuing request", "err", err)
				continue
			}
			ll.Infow("Enqueued fulfillment", "ethTxID", ethTX.ID)

			// If we successfully enqueued for the txm, subtract that balance
			// And loop to attempt to enqueue another fulfillment
			startBalanceNoReserveLink.Sub(startBalanceNoReserveLink, p.maxLink)
			processed[p.req.req.RequestId.String()] = struct{}{}
			incProcessedReqs(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, v2)
		}
	}

	return processed
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
		payload, err := coordinatorV2ABI.Pack("getCommitment", req.req.RequestId)
		if err != nil {
			// This shouldn't happen
			return fulfilled, errors.Wrap(err, "creating getCommitment payload")
		}

		reqBlockNumber := new(big.Int).SetUint64(req.req.Raw.BlockNumber)

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

	err := lsn.ethClient.BatchCallContext(ctx, calls)
	if err != nil {
		return fulfilled, errors.Wrap(err, "making batch call")
	}

	var errs error
	for i, call := range calls {
		if call.Error != nil {
			errs = multierr.Append(errs, fmt.Errorf("checking request %s with hash %s: %w",
				reqs[i].req.RequestId.String(), reqs[i].req.Raw.TxHash.String(), call.Error))
			continue
		}

		rString, ok := call.Result.(*string)
		if !ok {
			errs = multierr.Append(errs,
				fmt.Errorf("unexpected result %+v on request %s with hash %s",
					call.Result, reqs[i].req.RequestId.String(), reqs[i].req.Raw.TxHash.String()))
			continue
		}
		result, err := hexutil.Decode(*rString)
		if err != nil {
			errs = multierr.Append(errs,
				fmt.Errorf("decoding batch call result %+v %s request %s with hash %s: %w",
					call.Result, *rString, reqs[i].req.RequestId.String(), reqs[i].req.Raw.TxHash.String(), err))
			continue
		}

		if utils.IsEmpty(result) {
			l.Infow("Request already fulfilled",
				"reqID", reqs[i].req.RequestId.String(),
				"attempts", reqs[i].attempts,
				"txHash", reqs[i].req.Raw.TxHash)
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
			results[i] = lsn.simulateFulfillment(ctx, maxGasPriceWei, req, l)
		}(i, req)
	}
	wg.Wait()

	l.Debugw("Finished running pipelines",
		"count", len(reqs), "time", time.Since(start).String())
	return results
}

func (lsn *listenerV2) estimateFeeJuels(
	ctx context.Context,
	req *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested,
	maxGasPriceWei *assets.Wei,
) (*big.Int, error) {
	// Don't use up too much time to get this info, it's not critical for operating vrf.
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	roundData, err := lsn.aggregator.LatestRoundData(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, errors.Wrap(err, "get aggregator latestAnswer")
	}

	juelsNeeded, err := EstimateFeeJuels(
		req.CallbackGasLimit,
		maxGasPriceWei.ToInt(),
		roundData.Answer,
	)
	if err != nil {
		return nil, errors.Wrap(err, "estimate fee juels")
	}
	return juelsNeeded, nil
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
	// estimate how much juels are needed so that we can log it if the simulation fails.
	res.juelsNeeded, err = lsn.estimateFeeJuels(ctx, req.req, maxGasPriceWei)
	if err != nil {
		// not critical, just log and continue
		lg.Warnw("unable to estimate juels needed for request, continuing anyway",
			"reqID", req.req.RequestId,
			"err", err)
		res.juelsNeeded = big.NewInt(0)
	}

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"databaseID":    lsn.job.ID,
			"externalJobID": lsn.job.ExternalJobID,
			"name":          lsn.job.Name.ValueOrZero(),
			"publicKey":     lsn.job.VRFSpec.PublicKey[:],
			"maxGasPrice":   maxGasPriceWei.ToInt().String(),
		},
		"jobRun": map[string]interface{}{
			"logBlockHash":   req.req.Raw.BlockHash[:],
			"logBlockNumber": req.req.Raw.BlockNumber,
			"logTxHash":      req.req.Raw.TxHash,
			"logTopics":      req.req.Raw.Topics,
			"logData":        req.req.Raw.Data,
		},
	})
	var trrs pipeline.TaskRunResults
	res.run, trrs, err = lsn.pipelineRunner.ExecuteRun(ctx, *lsn.job.PipelineSpec, vars, lg)
	if err != nil {
		res.err = errors.Wrap(err, "executing run")
		return res
	}
	// The call task will fail if there are insufficient funds
	if res.run.AllErrors.HasError() {
		res.err = errors.WithStack(res.run.AllErrors.ToError())

		if strings.Contains(res.err.Error(), "blockhash not found in store") {
			res.err = multierr.Combine(res.err, errBlockhashNotInStore{})
		} else if strings.Contains(res.err.Error(), "execution reverted") {
			res.err = multierr.Combine(res.err, errPossiblyInsufficientFunds{})
		}

		return res
	}
	if len(trrs.FinalResult(lg).Values) != 1 {
		res.err = errors.Errorf("unexpected number of outputs, expected 1, was %d", len(trrs.FinalResult(lg).Values))
		return res
	}

	// Run succeeded, we expect a byte array representing the billing amount
	b, ok := trrs.FinalResult(lg).Values[0].([]uint8)
	if !ok {
		res.err = errors.New("expected []uint8 final result")
		return res
	}
	res.maxLink = utils.HexToBig(hexutil.Encode(b)[2:])
	for _, trr := range trrs {
		if trr.Task.Type() == pipeline.TaskTypeVRFV2 {
			m := trr.Result.Value.(map[string]interface{})
			res.payload = m["output"].(string)
			res.proof = m["proof"].(vrf_coordinator_v2.VRFProof)
			res.reqCommitment = m["requestCommitment"].(vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment)
		}

		if trr.Task.Type() == pipeline.TaskTypeEstimateGasLimit {
			res.gasLimit = trr.Result.Value.(uint32)
		}
	}
	return res
}

func (lsn *listenerV2) runRequestHandler(pollPeriod time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	tick := time.NewTicker(pollPeriod)
	defer tick.Stop()
	ctx, cancel := utils.ContextFromChan(lsn.chStop)
	defer cancel()
	for {
		select {
		case <-lsn.chStop:
			return
		case <-tick.C:
			lsn.processPendingVRFRequests(ctx)
		}
	}
}

func (lsn *listenerV2) runLogListener(unsubscribes []func(), minConfs uint32, wg *sync.WaitGroup) {
	defer wg.Done()
	lsn.l.Infow("Listening for run requests",
		"minConfs", minConfs)
	for {
		select {
		case <-lsn.chStop:
			for _, f := range unsubscribes {
				f()
			}
			return
		case <-lsn.reqLogs.Notify():
			// Process all the logs in the queue if one is added
			for {
				lb, exists := lsn.reqLogs.Retrieve()
				if !exists {
					break
				}
				lsn.handleLog(lb, minConfs)
			}
		}
	}
}

func (lsn *listenerV2) getConfirmedAt(req *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested, nodeMinConfs uint32) uint64 {
	lsn.respCountMu.Lock()
	defer lsn.respCountMu.Unlock()
	// Take the max(nodeMinConfs, requestedConfs + requestedConfsDelay).
	// Add the requested confs delay if provided in the jobspec so that we avoid an edge case
	// where the primary and backup VRF v2 nodes submit a proof at the same time.
	minConfs := nodeMinConfs
	if uint32(req.MinimumRequestConfirmations)+uint32(lsn.job.VRFSpec.RequestedConfsDelay) > nodeMinConfs {
		minConfs = uint32(req.MinimumRequestConfirmations) + uint32(lsn.job.VRFSpec.RequestedConfsDelay)
	}
	newConfs := uint64(minConfs) * (1 << lsn.respCount[req.RequestId.String()])
	// We cap this at 200 because solidity only supports the most recent 256 blocks
	// in the contract so if it was older than that, fulfillments would start failing
	// without the blockhash store feeder. We use 200 to give the node plenty of time
	// to fulfill even on fast chains.
	if newConfs > 200 {
		newConfs = 200
	}
	if lsn.respCount[req.RequestId.String()] > 0 {
		lsn.l.Warnw("Duplicate request found after fulfillment, doubling incoming confirmations",
			"txHash", req.Raw.TxHash,
			"blockNumber", req.Raw.BlockNumber,
			"blockHash", req.Raw.BlockHash,
			"reqID", req.RequestId.String(),
			"newConfs", newConfs)
		incDupeReqs(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, v2)
	}
	return req.Raw.BlockNumber + newConfs
}

func (lsn *listenerV2) handleLog(lb log.Broadcast, minConfs uint32) {
	if v, ok := lb.DecodedLog().(*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled); ok {
		lsn.l.Debugw("Received fulfilled log", "reqID", v.RequestId, "success", v.Success)
		consumed, err := lsn.logBroadcaster.WasAlreadyConsumed(lb)
		if err != nil {
			lsn.l.Errorw("Could not determine if log was already consumed", "error", err, "txHash", lb.RawLog().TxHash)
			return
		} else if consumed {
			return
		}
		lsn.respCountMu.Lock()
		lsn.respCount[v.RequestId.String()]++
		lsn.respCountMu.Unlock()
		lsn.blockNumberToReqID.Insert(fulfilledReqV2{
			blockNumber: v.Raw.BlockNumber,
			reqID:       v.RequestId.String(),
		})
		lsn.markLogAsConsumed(lb)
		return
	}

	req, err := lsn.coordinator.ParseRandomWordsRequested(lb.RawLog())
	if err != nil {
		lsn.l.Errorw("Failed to parse log", "err", err, "txHash", lb.RawLog().TxHash)
		consumed, err := lsn.logBroadcaster.WasAlreadyConsumed(lb)
		if err != nil {
			lsn.l.Errorw("Could not determine if log was already consumed", "error", err, "txHash", lb.RawLog().TxHash)
			return
		} else if consumed {
			return
		}
		lsn.markLogAsConsumed(lb)
		return
	}

	confirmedAt := lsn.getConfirmedAt(req, minConfs)
	lsn.l.Infow("VRFListenerV2: Received log request", "reqID", req.RequestId, "confirmedAt", confirmedAt, "subID", req.SubId, "sender", req.Sender)
	lsn.reqsMu.Lock()
	lsn.reqs = append(lsn.reqs, pendingRequest{
		confirmedAtBlock: confirmedAt,
		req:              req,
		lb:               lb,
		utcTimestamp:     time.Now().UTC(),
	})
	lsn.reqAdded()
	lsn.reqsMu.Unlock()
}

func (lsn *listenerV2) markLogAsConsumed(lb log.Broadcast) {
	err := lsn.logBroadcaster.MarkConsumed(lb)
	lsn.l.ErrorIf(err, fmt.Sprintf("Unable to mark log %v as consumed", lb.String()))
}

// Close complies with job.Service
func (lsn *listenerV2) Close() error {
	return lsn.StopOnce("VRFListenerV2", func() error {
		close(lsn.chStop)
		// wait on the request handler, log listener, and head listener to stop
		lsn.wg.Wait()
		return lsn.reqLogs.Close()
	})
}

func (lsn *listenerV2) HandleLog(lb log.Broadcast) {
	if !lsn.deduper.shouldDeliver(lb.RawLog()) {
		lsn.l.Tracew("skipping duplicate log broadcast", "log", lb.RawLog())
		return
	}

	wasOverCapacity := lsn.reqLogs.Deliver(lb)
	if wasOverCapacity {
		lsn.l.Error("Log mailbox is over capacity - dropped the oldest log")
		incDroppedReqs(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, v2, reasonMailboxSize)
	}
}

// JobID complies with log.Listener
func (lsn *listenerV2) JobID() int32 {
	return lsn.job.ID
}

// ReplayStartedCallback is called by the log broadcaster when a replay is about to start.
func (lsn *listenerV2) ReplayStartedCallback() {
	// Clear the log deduper cache so that we don't incorrectly ignore logs that have been sent that
	// are already in the cache.
	lsn.deduper.clear()
}

func (lsn *listenerV2) fromAddresses() []common.Address {
	var addresses []common.Address
	for _, a := range lsn.job.VRFSpec.FromAddresses {
		addresses = append(addresses, a.Address())
	}
	return addresses
}

func uniqueReqs(reqs []pendingRequest) int {
	s := map[string]struct{}{}
	for _, r := range reqs {
		s[r.req.RequestId.String()] = struct{}{}
	}
	return len(s)
}

// GasProofVerification is an upper limit on the gas used for verifying the VRF proof on-chain.
// It can be used to estimate the amount of LINK needed to fulfill a request.
const GasProofVerification uint32 = 200_000

// EstimateFeeJuels estimates the amount of link needed to fulfill a request
// given the callback gas limit, the gas price, and the wei per unit link.
// An error is returned if the wei per unit link provided is zero.
func EstimateFeeJuels(callbackGasLimit uint32, maxGasPriceWei, weiPerUnitLink *big.Int) (*big.Int, error) {
	if weiPerUnitLink.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("wei per unit link is zero")
	}
	maxGasUsed := big.NewInt(int64(callbackGasLimit + GasProofVerification))
	costWei := maxGasUsed.Mul(maxGasUsed, maxGasPriceWei)
	// Multiply by 1e18 first so that we don't lose a ton of digits due to truncation when we divide
	// by weiPerUnitLink
	numerator := costWei.Mul(costWei, big.NewInt(1e18))
	costJuels := numerator.Quo(numerator, weiPerUnitLink)
	return costJuels, nil
}
