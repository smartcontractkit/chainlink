package vrf

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	heaps "github.com/theodesp/go-heaps"
	"github.com/theodesp/go-heaps/pairing"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	_ log.Listener = &listenerV1{}
	_ job.Service  = &listenerV1{}
)

type request struct {
	confirmedAtBlock uint64
	req              *solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest
	lb               log.Broadcast
}

type listenerV1 struct {
	utils.StartStopOnce

	cfg             Config
	l               logger.Logger
	abi             abi.ABI
	logBroadcaster  log.Broadcaster
	coordinator     *solidity_vrf_coordinator_interface.VRFCoordinator
	pipelineRunner  pipeline.Runner
	pipelineORM     pipeline.ORM
	vorm            keystore.VRFORM
	job             job.Job
	db              *gorm.DB
	headBroadcaster httypes.HeadBroadcasterRegistry
	txm             bulletprooftxmanager.TxManager
	vrfks           keystore.VRF
	gethks          GethKeyStore
	reqLogs         *utils.Mailbox
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
}

// Note that we have 2 seconds to do this processing
func (lsn *listenerV1) OnNewLongestChain(_ context.Context, head eth.Head) {
	lsn.setLatestHead(head)
	select {
	case lsn.newHead <- struct{}{}:
	default:
	}
}

func (lsn *listenerV1) setLatestHead(h eth.Head) {
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
func (lsn *listenerV1) Start() error {
	return lsn.StartOnce("VRFListener", func() error {
		// Take the larger of the global vs specific
		// Note that runtime changes to incoming confirmations require a job delete/add
		// because we need to resubscribe to the lb with the new min.
		minConfs := lsn.cfg.MinIncomingConfirmations()
		if lsn.job.VRFSpec.Confirmations > lsn.cfg.MinIncomingConfirmations() {
			minConfs = lsn.job.VRFSpec.Confirmations
		}
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
			// If we set this to minConfs, since both the log broadcaster and head broadcaster get heads
			// at the same time from the head tracker whether we process the log at minConfs or minConfs+1
			// would depend on the order in which their OnNewLongestChain callbacks got called.
			// We listen one block early so that the log can be stored in pendingRequests
			// to avoid this.
			NumConfirmations: uint64(minConfs - 1),
		})
		// Subscribe to the head broadcaster for handling
		// per request conf requirements.
		latestHead, unsubscribeHeadBroadcaster := lsn.headBroadcaster.Subscribe(lsn)
		if latestHead != nil {
			lsn.setLatestHead(*latestHead)
		}
		go gracefulpanic.WrapRecover(lsn.l, func() {
			lsn.runLogListener([]func(){unsubscribeLogs}, minConfs)
		})
		go gracefulpanic.WrapRecover(lsn.l, func() {
			lsn.runHeadListener(unsubscribeHeadBroadcaster)
		})
		return nil
	})
}

// Removes and returns all the confirmed logs from
// the pending queue.
func (lsn *listenerV1) extractConfirmedLogs() []request {
	lsn.reqsMu.Lock()
	defer lsn.reqsMu.Unlock()
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
	for {
		select {
		case <-lsn.chStop:
			unsubscribe()
			lsn.waitOnStop <- struct{}{}
			return
		case <-lsn.newHead:
			toProcess := lsn.extractConfirmedLogs()
			for _, r := range toProcess {
				lsn.ProcessRequest(r.req, r.lb)
			}
			lsn.pruneConfirmedRequestCounts()
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
				i, exists := lsn.reqLogs.Retrieve()
				if !exists {
					break
				}
				lb, ok := i.(log.Broadcast)
				if !ok {
					panic(fmt.Sprintf("VRFListener: invariant violated, expected log.Broadcast got %T", i))
				}
				lsn.handleLog(lb, minConfs)
			}
		}
	}
}

func (lsn *listenerV1) handleLog(lb log.Broadcast, minConfs uint32) {
	lsn.l.Infow("Log received", lb.String(), lb.DecodedLog())
	if v, ok := lb.DecodedLog().(*solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequestFulfilled); ok {
		if !lsn.shouldProcessLog(lb) {
			return
		}
		lsn.respCountMu.Lock()
		lsn.respCount[v.RequestId]++
		lsn.respCountMu.Unlock()
		lsn.blockNumberToReqID.Insert(fulfilledReq{
			blockNumber: v.Raw.BlockNumber,
			reqID:       v.RequestId,
		})
		lsn.markLogAsConsumed(lb)
		return
	}

	req, err := lsn.coordinator.ParseRandomnessRequest(lb.RawLog())
	if err != nil {
		lsn.l.Errorw("Failed to parse log", "err", err, "txHash", lb.RawLog().TxHash)
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
	})
	lsn.reqAdded()
	lsn.reqsMu.Unlock()
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
			"reqID", hex.EncodeToString(req.RequestID[:]),
			"newConfs", newConfs)
	}
	return req.Raw.BlockNumber + newConfs
}

func (lsn *listenerV1) ProcessRequest(req *solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest, lb log.Broadcast) {
	// This check to see if the log was consumed needs to be in the same
	// goroutine as the mark consumed to avoid processing duplicates.
	if !lsn.shouldProcessLog(lb) {
		return
	}

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
	callback, err := lsn.coordinator.Callbacks(nil, req.RequestID)
	if err != nil {
		lsn.l.Errorw("Unable to check if already fulfilled, processing anyways", "err", err, "txHash", req.Raw.TxHash)
	} else if utils.IsEmpty(callback.SeedAndBlockNum[:]) {
		// If seedAndBlockNumber is zero then the response has been fulfilled
		// and we should skip it
		lsn.l.Infow("Request already fulfilled", "txHash", req.Raw.TxHash, "reqID", req.RequestID)
		lsn.markLogAsConsumed(lb)
		return
	}

	lsn.l.Infow("Received log request",
		"log", lb.String(),
		"reqID", hex.EncodeToString(req.RequestID[:]),
		"keyHash", hex.EncodeToString(req.KeyHash[:]),
		"txHash", req.Raw.TxHash,
		"blockNumber", req.Raw.BlockNumber,
		"blockHash", req.Raw.BlockHash,
		"seed", req.Seed,
		"fee", req.Fee)

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"databaseID":    lsn.job.ID,
			"externalJobID": lsn.job.ExternalJobID,
			"name":          lsn.job.Name.ValueOrZero(),
			"publicKey":     lsn.job.VRFSpec.PublicKey[:],
		},
		"jobRun": map[string]interface{}{
			"logBlockHash":   req.Raw.BlockHash[:],
			"logBlockNumber": req.Raw.BlockNumber,
			"logTxHash":      req.Raw.TxHash,
			"logTopics":      req.Raw.Topics,
			"logData":        req.Raw.Data,
		},
	})

	run := pipeline.NewRun(*lsn.job.PipelineSpec, vars)
	if _, err = lsn.pipelineRunner.Run(context.Background(), &run, lsn.l, true, func(tx postgres.Queryer) error {
		// Always mark consumed regardless of whether the proof failed or not.
		if err = lsn.logBroadcaster.MarkConsumed(lb, postgres.WithQueryer(tx)); err != nil {
			lsn.l.Errorw("Failed mark consumed", "err", err)
		}
		return nil
	}); err != nil {
		lsn.l.Errorw("Failed executing run", "err", err)
	}
}

// Close complies with job.Service
func (lsn *listenerV1) Close() error {
	return lsn.StopOnce("VRFListener", func() error {
		close(lsn.chStop)
		<-lsn.waitOnStop // Log listenerV1
		<-lsn.waitOnStop // Head listenerV1
		return nil
	})
}

func (lsn *listenerV1) HandleLog(lb log.Broadcast) {
	wasOverCapacity := lsn.reqLogs.Deliver(lb)
	if wasOverCapacity {
		lsn.l.Error("l mailbox is over capacity - dropped the oldest l")
	}
}

// Job complies with log.Listener
func (lsn *listenerV1) JobID() int32 {
	return lsn.job.ID
}
