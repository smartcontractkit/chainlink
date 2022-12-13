package vrf

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	heaps "github.com/theodesp/go-heaps"
	"github.com/theodesp/go-heaps/pairing"

	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/recovery"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
	bigmath "github.com/smartcontractkit/chainlink/core/utils/big_math"
)

var (
	_ log.Listener   = &listenerV1{}
	_ job.ServiceCtx = &listenerV1{}
)

const callbacksTimeout = 10 * time.Second

type request struct {
	confirmedAtBlock uint64
	req              *solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest
	lb               log.Broadcast
	utcTimestamp     time.Time
}

type listenerV1 struct {
	utils.StartStopOnce

	cfg             Config
	l               logger.SugaredLogger
	logBroadcaster  log.Broadcaster
	coordinator     *solidity_vrf_coordinator_interface.VRFCoordinator
	pipelineRunner  pipeline.Runner
	job             job.Job
	q               pg.Q
	headBroadcaster httypes.HeadBroadcasterRegistry
	txm             txmgr.TxManager
	gethks          GethKeyStore
	mailMon         *utils.MailboxMonitor
	reqLogs         *utils.Mailbox[log.Broadcast]
	chStop          chan struct{}
	waitOnStop      chan struct{}
	newHead         chan struct{}
	latestHead      uint64
	latestHeadMu    sync.RWMutex
	// We can keep these pending logs in memory because we
	// only mark them confirmed once we send a corresponding fulfillment transaction.
	// So on node restart in the middle of processing, the lb will resend them.
	reqsMu   sync.Mutex // Both goroutines write to reqs
	reqs     []request
	reqAdded func() // A simple debug helper

	// Data structures for reorg attack protection
	// We want a map so we can do an O(1) count update every fulfillment log we get.
	respCountMu sync.Mutex
	respCount   map[[32]byte]uint64
	// This auxiliary heap is to used when we need to purge the
	// respCount map - we repeatedly want remove the minimum log.
	// You could use a sorted list if the completed logs arrive in order, but they may not.
	blockNumberToReqID *pairing.PairHeap

	// deduper prevents processing duplicate requests from the log broadcaster.
	deduper *logDeduper
}

// Note that we have 2 seconds to do this processing
func (lsn *listenerV1) OnNewLongestChain(_ context.Context, head *evmtypes.Head) {
	lsn.setLatestHead(head)
	select {
	case lsn.newHead <- struct{}{}:
	default:
	}
}

func (lsn *listenerV1) setLatestHead(h *evmtypes.Head) {
	lsn.latestHeadMu.Lock()
	defer lsn.latestHeadMu.Unlock()
	num := uint64(h.Number)
	if num > lsn.latestHead {
		lsn.latestHead = num
	}
}

func (lsn *listenerV1) getLatestHead() uint64 {
	lsn.latestHeadMu.RLock()
	defer lsn.latestHeadMu.RUnlock()
	return lsn.latestHead
}

// Start complies with job.Service
func (lsn *listenerV1) Start(context.Context) error {
	return lsn.StartOnce("VRFListener", func() error {
		spec := job.LoadEnvConfigVarsVRF(lsn.cfg, *lsn.job.VRFSpec)

		unsubscribeLogs := lsn.logBroadcaster.Register(lsn, log.ListenerOpts{
			Contract: lsn.coordinator.Address(),
			ParseLog: lsn.coordinator.ParseLog,
			LogsWithTopics: map[common.Hash][][]log.Topic{
				solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest{}.Topic(): {
					{
						log.Topic(lsn.job.ExternalIDEncodeStringToTopic()),
						log.Topic(lsn.job.ExternalIDEncodeBytesToTopic()),
					},
				},
				solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequestFulfilled{}.Topic(): {},
			},
			// If we set this to MinIncomingConfirmations, since both the log broadcaster and head broadcaster get heads
			// at the same time from the head tracker whether we process the log at MinIncomingConfirmations or
			// MinIncomingConfirmations+1 would depend on the order in which their OnNewLongestChain callbacks got
			// called.
			// We listen one block early so that the log can be stored in pendingRequests to avoid this.
			MinIncomingConfirmations: spec.MinIncomingConfirmations - 1,
			ReplayStartedCallback:    lsn.ReplayStartedCallback,
		})
		// Subscribe to the head broadcaster for handling
		// per request conf requirements.
		latestHead, unsubscribeHeadBroadcaster := lsn.headBroadcaster.Subscribe(lsn)
		if latestHead != nil {
			lsn.setLatestHead(latestHead)
		}
		go lsn.runLogListener([]func(){unsubscribeLogs}, spec.MinIncomingConfirmations)
		go lsn.runHeadListener(unsubscribeHeadBroadcaster)

		lsn.mailMon.Monitor(lsn.reqLogs, "VRFListener", "RequestLogs", fmt.Sprint(lsn.job.ID))
		return nil
	})
}

// Removes and returns all the confirmed logs from
// the pending queue.
func (lsn *listenerV1) extractConfirmedLogs() []request {
	lsn.reqsMu.Lock()
	defer lsn.reqsMu.Unlock()
	updateQueueSize(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, v1, len(lsn.reqs))
	var toProcess, toKeep []request
	for i := 0; i < len(lsn.reqs); i++ {
		if lsn.reqs[i].confirmedAtBlock <= lsn.getLatestHead() {
			toProcess = append(toProcess, lsn.reqs[i])
		} else {
			toKeep = append(toKeep, lsn.reqs[i])
		}
	}
	lsn.reqs = toKeep
	return toProcess
}

type fulfilledReq struct {
	blockNumber uint64
	reqID       [32]byte
}

func (a fulfilledReq) Compare(b heaps.Item) int {
	a1 := a
	a2 := b.(fulfilledReq)
	switch {
	case a1.blockNumber > a2.blockNumber:
		return 1
	case a1.blockNumber < a2.blockNumber:
		return -1
	default:
		return 0
	}
}

// Remove all entries 10000 blocks or older
// to avoid a memory leak.
func (lsn *listenerV1) pruneConfirmedRequestCounts() {
	lsn.respCountMu.Lock()
	defer lsn.respCountMu.Unlock()
	min := lsn.blockNumberToReqID.FindMin()
	for min != nil {
		m := min.(fulfilledReq)
		if m.blockNumber > (lsn.getLatestHead() - 10000) {
			break
		}
		delete(lsn.respCount, m.reqID)
		lsn.blockNumberToReqID.DeleteMin()
		min = lsn.blockNumberToReqID.FindMin()
	}
}

// Listen for new heads
func (lsn *listenerV1) runHeadListener(unsubscribe func()) {
	ctx, cancel := utils.ContextFromChan(lsn.chStop)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			unsubscribe()
			lsn.waitOnStop <- struct{}{}
			return
		case <-lsn.newHead:
			recovery.WrapRecover(lsn.l, func() {
				toProcess := lsn.extractConfirmedLogs()
				var toRetry []request
				for _, r := range toProcess {
					if success := lsn.ProcessRequest(ctx, r); !success {
						toRetry = append(toRetry, r)
					}
				}
				lsn.reqsMu.Lock()
				defer lsn.reqsMu.Unlock()
				lsn.reqs = append(lsn.reqs, toRetry...)
				lsn.pruneConfirmedRequestCounts()
			})
		}
	}
}

func (lsn *listenerV1) runLogListener(unsubscribes []func(), minConfs uint32) {
	lsn.l.Infow("Listening for run requests",
		"gasLimit", lsn.cfg.EvmGasLimitDefault(),
		"minConfs", minConfs)
	for {
		select {
		case <-lsn.chStop:
			for _, f := range unsubscribes {
				f()
			}
			lsn.waitOnStop <- struct{}{}
			return
		case <-lsn.reqLogs.Notify():
			// Process all the logs in the queue if one is added
			for {
				lb, exists := lsn.reqLogs.Retrieve()
				if !exists {
					break
				}
				recovery.WrapRecover(lsn.l, func() {
					lsn.handleLog(lb, minConfs)
				})
			}
		}
	}
}

func (lsn *listenerV1) handleLog(lb log.Broadcast, minConfs uint32) {
	lggr := lsn.l.With(
		"log", lb.String(),
		"decodedLog", lb.DecodedLog(),
		"blockNumber", lb.RawLog().BlockNumber,
		"blockHash", lb.RawLog().BlockHash,
		"txHash", lb.RawLog().TxHash,
	)

	lggr.Infow("Log received")
	if v, ok := lb.DecodedLog().(*solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequestFulfilled); ok {
		lggr.Debugw("Got fulfillment log",
			"requestID", hex.EncodeToString(v.RequestId[:]))
		if !lsn.shouldProcessLog(lb) {
			return
		}
		lsn.respCountMu.Lock()
		lsn.respCount[v.RequestId]++
		lsn.blockNumberToReqID.Insert(fulfilledReq{
			blockNumber: v.Raw.BlockNumber,
			reqID:       v.RequestId,
		})
		lsn.respCountMu.Unlock()
		lsn.markLogAsConsumed(lb)
		return
	}

	req, err := lsn.coordinator.ParseRandomnessRequest(lb.RawLog())
	if err != nil {
		lggr.Errorw("Failed to parse RandomnessRequest log", "err", err)
		if !lsn.shouldProcessLog(lb) {
			return
		}
		lsn.markLogAsConsumed(lb)
		return
	}

	confirmedAt := lsn.getConfirmedAt(req, minConfs)
	lsn.reqsMu.Lock()
	lsn.reqs = append(lsn.reqs, request{
		confirmedAtBlock: confirmedAt,
		req:              req,
		lb:               lb,
		utcTimestamp:     time.Now().UTC(),
	})
	lsn.reqAdded()
	lsn.reqsMu.Unlock()
	lggr.Infow("Enqueued randomness request",
		"requestID", hex.EncodeToString(req.RequestID[:]),
		"requestJobID", hex.EncodeToString(req.JobID[:]),
		"keyHash", hex.EncodeToString(req.KeyHash[:]),
		"fee", req.Fee,
		"sender", req.Sender.Hex(),
		"txHash", lb.RawLog().TxHash)
}

func (lsn *listenerV1) shouldProcessLog(lb log.Broadcast) bool {
	consumed, err := lsn.logBroadcaster.WasAlreadyConsumed(lb)
	if err != nil {
		lsn.l.Errorw("Could not determine if log was already consumed", "error", err, "txHash", lb.RawLog().TxHash)
		// Do not process, let lb resend it as a retry mechanism.
		return false
	}
	return !consumed
}

func (lsn *listenerV1) markLogAsConsumed(lb log.Broadcast) {
	err := lsn.logBroadcaster.MarkConsumed(lb)
	lsn.l.ErrorIf(err, fmt.Sprintf("Unable to mark log %v as consumed", lb.String()))
}

func (lsn *listenerV1) getConfirmedAt(req *solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest, minConfs uint32) uint64 {
	lsn.respCountMu.Lock()
	defer lsn.respCountMu.Unlock()
	newConfs := uint64(minConfs) * (1 << lsn.respCount[req.RequestID])
	// We cap this at 200 because solidity only supports the most recent 256 blocks
	// in the contract so if it was older than that, fulfillments would start failing
	// without the blockhash store feeder. We use 200 to give the node plenty of time
	// to fulfill even on fast chains.
	if newConfs > 200 {
		newConfs = 200
	}
	if lsn.respCount[req.RequestID] > 0 {
		lsn.l.Warnw("Duplicate request found after fulfillment, doubling incoming confirmations",
			"txHash", req.Raw.TxHash,
			"blockNumber", req.Raw.BlockNumber,
			"blockHash", req.Raw.BlockHash,
			"requestID", hex.EncodeToString(req.RequestID[:]),
			"newConfs", newConfs)
		incDupeReqs(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, v1)
	}
	return req.Raw.BlockNumber + newConfs
}

// ProcessRequest attempts to process the VRF request. Returns true if successful, false otherwise.
func (lsn *listenerV1) ProcessRequest(ctx context.Context, req request) bool {
	// This check to see if the log was consumed needs to be in the same
	// goroutine as the mark consumed to avoid processing duplicates.
	if !lsn.shouldProcessLog(req.lb) {
		return true
	}

	lggr := lsn.l.With(
		"log", req.lb.String(),
		"requestID", hex.EncodeToString(req.req.RequestID[:]),
		"txHash", req.req.Raw.TxHash,
		"keyHash", hex.EncodeToString(req.req.KeyHash[:]),
		"jobID", hex.EncodeToString(req.req.JobID[:]),
		"sender", req.req.Sender.Hex(),
		"blockNumber", req.req.Raw.BlockNumber,
		"blockHash", req.req.Raw.BlockHash,
		"seed", req.req.Seed,
		"fee", req.req.Fee,
	)

	// Check if the vrf req has already been fulfilled
	// Note we have to do this after the log has been confirmed.
	// If not, the following problematic (example) scenario can arise:
	// 1. Request log comes in block 100
	// 2. Fulfill the request in block 110
	// 3. Reorg both request and fulfillment, now request lives at
	// block 101 and fulfillment lives at block 115
	// 4. The eth node sees the request reorg and tells us about it. We do our fulfillment
	// check and the node says its already fulfilled (hasn't seen the fulfillment reorged yet),
	// so we don't process the request.
	// Subtract 5 since the newest block likely isn't indexed yet and will cause "header not
	// found" errors.
	currBlock := new(big.Int).SetUint64(lsn.getLatestHead() - 5)
	m := bigmath.Max(req.confirmedAtBlock, currBlock)
	ctx, cancel := context.WithTimeout(ctx, callbacksTimeout)
	defer cancel()
	callback, err := lsn.coordinator.Callbacks(&bind.CallOpts{
		BlockNumber: m,
		Context:     ctx,
	}, req.req.RequestID)
	if err != nil {
		lggr.Errorw("Unable to check if already fulfilled, processing anyways", "err", err)
	} else if utils.IsEmpty(callback.SeedAndBlockNum[:]) {
		// If seedAndBlockNumber is zero then the response has been fulfilled
		// and we should skip it
		lggr.Infow("Request already fulfilled")
		lsn.markLogAsConsumed(req.lb)
		return true
	}

	// Check if we can ignore the request due to its age.
	if time.Now().UTC().Sub(req.utcTimestamp) >= lsn.job.VRFSpec.RequestTimeout {
		lggr.Infow("Request too old, dropping it")
		lsn.markLogAsConsumed(req.lb)
		return true
	}

	lggr.Infow("Processing log request")

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"databaseID":    lsn.job.ID,
			"externalJobID": lsn.job.ExternalJobID,
			"name":          lsn.job.Name.ValueOrZero(),
			"publicKey":     lsn.job.VRFSpec.PublicKey[:],
			"from":          lsn.fromAddresses(),
		},
		"jobRun": map[string]interface{}{
			"logBlockHash":   req.req.Raw.BlockHash[:],
			"logBlockNumber": req.req.Raw.BlockNumber,
			"logTxHash":      req.req.Raw.TxHash,
			"logTopics":      req.req.Raw.Topics,
			"logData":        req.req.Raw.Data,
		},
	})

	run := pipeline.NewRun(*lsn.job.PipelineSpec, vars)
	// The VRF pipeline has no async tasks, so we don't need to check for `incomplete`
	if _, err = lsn.pipelineRunner.Run(ctx, &run, lggr, true, func(tx pg.Queryer) error {
		// Always mark consumed regardless of whether the proof failed or not.
		if err = lsn.logBroadcaster.MarkConsumed(req.lb, pg.WithQueryer(tx)); err != nil {
			lggr.Errorw("Failed mark consumed", "err", err)
		}
		return nil
	}); err != nil {
		lggr.Errorw("Failed to execute VRFV1 pipeline run",
			"err", err)
		return false
	}

	// At this point the pipeline runner has completed the run of the pipeline,
	// but it may have errored out.
	if run.HasErrors() || run.HasFatalErrors() {
		lggr.Error("VRFV1 pipeline run failed with errors",
			"runErrors", run.AllErrors.ToError(),
			"runFatalErrors", run.FatalErrors.ToError(),
		)
		return false
	}

	// At this point, the pipeline run executed successfully, and we mark
	// the request as processed.
	lggr.Infow("Executed VRFV1 fulfillment run")
	incProcessedReqs(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, v1)
	return true
}

// Close complies with job.Service
func (lsn *listenerV1) Close() error {
	return lsn.StopOnce("VRFListener", func() error {
		close(lsn.chStop)
		<-lsn.waitOnStop // Log listenerV1
		<-lsn.waitOnStop // Head listenerV1
		return lsn.reqLogs.Close()
	})
}

func (lsn *listenerV1) HandleLog(lb log.Broadcast) {
	if !lsn.deduper.shouldDeliver(lb.RawLog()) {
		lsn.l.Tracew("skipping duplicate log broadcast", "log", lb.RawLog())
		return
	}

	wasOverCapacity := lsn.reqLogs.Deliver(lb)
	if wasOverCapacity {
		lsn.l.Error("log mailbox is over capacity - dropped the oldest log")
		incDroppedReqs(lsn.job.Name.ValueOrZero(), lsn.job.ExternalJobID, v1, reasonMailboxSize)
	}
}

func (lsn *listenerV1) fromAddresses() []common.Address {
	var addresses []common.Address
	for _, a := range lsn.job.VRFSpec.FromAddresses {
		addresses = append(addresses, a.Address())
	}
	return addresses
}

// Job complies with log.Listener
func (lsn *listenerV1) JobID() int32 {
	return lsn.job.ID
}

// ReplayStartedCallback is called by the log broadcaster when a replay is about to start.
func (lsn *listenerV1) ReplayStartedCallback() {
	// Clear the log deduper cache so that we don't incorrectly ignore logs that have been sent that
	// are already in the cache.
	lsn.deduper.clear()
}
