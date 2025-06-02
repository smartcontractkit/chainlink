package v1

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	heaps "github.com/theodesp/go-heaps"
	"github.com/theodesp/go-heaps/pairing"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mathutil"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/log"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/recovery"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	_ log.Listener   = &Listener{}
	_ job.ServiceCtx = &Listener{}
)

const callbacksTimeout = 10 * time.Second

type request struct {
	confirmedAtBlock uint64
	req              *solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest
	lb               log.Broadcast
	utcTimestamp     time.Time
}

type Listener struct {
	services.StateMachine

	Cfg            vrfcommon.Config
	FeeCfg         vrfcommon.FeeConfig
	L              logger.SugaredLogger
	Coordinator    *solidity_vrf_coordinator_interface.VRFCoordinator
	PipelineRunner pipeline.Runner
	Job            job.Job
	GethKs         vrfcommon.GethKeyStore
	MailMon        *mailbox.Monitor
	ReqLogs        *mailbox.Mailbox[log.Broadcast]
	ChStop         services.StopChan
	WaitOnStop     chan struct{}
	NewHead        chan struct{}
	LatestHead     uint64
	LatestHeadMu   sync.RWMutex
	Chain          legacyevm.Chain

	// We can keep these pending logs in memory because we
	// only mark them confirmed once we send a corresponding fulfillment transaction.
	// So on node restart in the middle of processing, the lb will resend them.
	ReqsMu   sync.Mutex // Both goroutines write to Reqs
	Reqs     []request
	ReqAdded func() // A simple debug helper

	// Data structures for reorg attack protection
	// We want a map so we can do an O(1) count update every fulfillment log we get.
	RespCountMu   sync.Mutex
	ResponseCount map[[32]byte]uint64
	// This auxiliary heap is to used when we need to purge the
	// ResponseCount map - we repeatedly want remove the minimum log.
	// You could use a sorted list if the completed logs arrive in order, but they may not.
	BlockNumberToReqID *pairing.PairHeap

	// Deduper prevents processing duplicate requests from the log broadcaster.
	Deduper *vrfcommon.LogDeduper
}

// Note that we have 2 seconds to do this processing
func (lsn *Listener) OnNewLongestChain(_ context.Context, head *evmtypes.Head) {
	lsn.setLatestHead(head)
	select {
	case lsn.NewHead <- struct{}{}:
	default:
	}
}

func (lsn *Listener) setLatestHead(h *evmtypes.Head) {
	lsn.LatestHeadMu.Lock()
	defer lsn.LatestHeadMu.Unlock()
	num := uint64(h.Number)
	if num > lsn.LatestHead {
		lsn.LatestHead = num
	}
}

func (lsn *Listener) getLatestHead() uint64 {
	lsn.LatestHeadMu.RLock()
	defer lsn.LatestHeadMu.RUnlock()
	return lsn.LatestHead
}

// Start complies with job.Service
func (lsn *Listener) Start(ctx context.Context) error {
	return lsn.StartOnce("VRFListener", func() error {
		spec := job.LoadDefaultVRFPollPeriod(*lsn.Job.VRFSpec)

		unsubscribeLogs := lsn.Chain.LogBroadcaster().Register(lsn, log.ListenerOpts{
			Contract: lsn.Coordinator.Address(),
			ParseLog: lsn.Coordinator.ParseLog,
			LogsWithTopics: map[common.Hash][][]log.Topic{
				solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest{}.Topic(): {
					{
						log.Topic(lsn.Job.ExternalIDEncodeStringToTopic()),
						log.Topic(lsn.Job.ExternalIDEncodeBytesToTopic()),
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
		latestHead, unsubscribeHeadBroadcaster := lsn.Chain.HeadBroadcaster().Subscribe(lsn)
		if latestHead != nil {
			lsn.setLatestHead(latestHead)
		}

		// Populate the response count map
		lsn.RespCountMu.Lock()
		defer lsn.RespCountMu.Unlock()
		respCount, err := lsn.GetStartingResponseCountsV1(ctx)
		if err != nil {
			return err
		}
		lsn.ResponseCount = respCount
		go lsn.RunLogListener([]func(){unsubscribeLogs}, spec.MinIncomingConfirmations)
		go lsn.RunHeadListener(unsubscribeHeadBroadcaster)

		lsn.MailMon.Monitor(lsn.ReqLogs, "VRFListener", "RequestLogs", strconv.Itoa(int(lsn.Job.ID)))
		return nil
	})
}

func (lsn *Listener) GetStartingResponseCountsV1(ctx context.Context) (respCount map[[32]byte]uint64, err error) {
	respCounts := make(map[[32]byte]uint64)
	var latestBlockNum *big.Int
	// Retry client call for LatestBlockHeight if fails
	// Want to avoid failing startup due to potential faulty RPC call
	err = retry.Do(func() error {
		latestBlockNum, err = lsn.Chain.Client().LatestBlockHeight(ctx)
		return err
	}, retry.Attempts(10), retry.Delay(500*time.Millisecond))
	if err != nil {
		return nil, err
	}
	if latestBlockNum == nil {
		return nil, errors.New("LatestBlockHeight return nil block num")
	}
	confirmedBlockNum := latestBlockNum.Int64() - int64(lsn.Chain.Config().EVM().FinalityDepth())
	// Only check as far back as the evm finality depth for completed transactions.
	var counts []vrfcommon.RespCountEntry
	counts, err = vrfcommon.GetRespCounts(ctx, lsn.Chain.TxManager(), lsn.Chain.Client().ConfiguredChainID(), confirmedBlockNum)
	if err != nil {
		// Continue with an empty map, do not block job on this.
		lsn.L.Errorw("Unable to read previous confirmed fulfillments", "err", err)
		return respCounts, nil
	}

	for _, c := range counts {
		// Remove the quotes from the json
		req := strings.Replace(c.RequestID, `"`, ``, 2)
		// Remove the 0x prefix
		b, err := hex.DecodeString(req[2:])
		if err != nil {
			lsn.L.Errorw("Unable to read fulfillment", "err", err, "reqID", c.RequestID)
			continue
		}
		var reqID [32]byte
		copy(reqID[:], b)
		respCounts[reqID] = uint64(c.Count)
	}

	return respCounts, nil
}

// Removes and returns all the confirmed logs from
// the pending queue.
func (lsn *Listener) extractConfirmedLogs() []request {
	lsn.ReqsMu.Lock()
	defer lsn.ReqsMu.Unlock()
	vrfcommon.UpdateQueueSize(lsn.Job.Name.ValueOrZero(), lsn.Job.ExternalJobID, vrfcommon.V1, len(lsn.Reqs))
	var toProcess, toKeep []request
	for i := 0; i < len(lsn.Reqs); i++ {
		if lsn.Reqs[i].confirmedAtBlock <= lsn.getLatestHead() {
			toProcess = append(toProcess, lsn.Reqs[i])
		} else {
			toKeep = append(toKeep, lsn.Reqs[i])
		}
	}
	lsn.Reqs = toKeep
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
func (lsn *Listener) pruneConfirmedRequestCounts() {
	lsn.RespCountMu.Lock()
	defer lsn.RespCountMu.Unlock()
	min := lsn.BlockNumberToReqID.FindMin()
	for min != nil {
		m := min.(fulfilledReq)
		if m.blockNumber > (lsn.getLatestHead() - 10000) {
			break
		}
		delete(lsn.ResponseCount, m.reqID)
		lsn.BlockNumberToReqID.DeleteMin()
		min = lsn.BlockNumberToReqID.FindMin()
	}
}

// Listen for new heads
func (lsn *Listener) RunHeadListener(unsubscribe func()) {
	ctx, cancel := lsn.ChStop.NewCtx()
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			unsubscribe()
			lsn.WaitOnStop <- struct{}{}
			return
		case <-lsn.NewHead:
			recovery.WrapRecover(lsn.L, func() {
				toProcess := lsn.extractConfirmedLogs()
				var toRetry []request
				for _, r := range toProcess {
					if success := lsn.ProcessRequest(ctx, r); !success {
						toRetry = append(toRetry, r)
					}
				}
				lsn.ReqsMu.Lock()
				defer lsn.ReqsMu.Unlock()
				lsn.Reqs = append(lsn.Reqs, toRetry...)
				lsn.pruneConfirmedRequestCounts()
			})
		}
	}
}

func (lsn *Listener) RunLogListener(unsubscribes []func(), minConfs uint32) {
	ctx, cancel := lsn.ChStop.NewCtx()
	defer cancel()
	lsn.L.Infow("Listening for run requests",
		"gasLimit", lsn.FeeCfg.LimitDefault(),
		"minConfs", minConfs)
	for {
		select {
		case <-lsn.ChStop:
			for _, f := range unsubscribes {
				f()
			}
			lsn.WaitOnStop <- struct{}{}
			return
		case <-lsn.ReqLogs.Notify():
			// Process all the logs in the queue if one is added
			for {
				lb, exists := lsn.ReqLogs.Retrieve()
				if !exists {
					break
				}
				recovery.WrapRecover(lsn.L, func() {
					lsn.handleLog(ctx, lb, minConfs)
				})
			}
		}
	}
}

func (lsn *Listener) handleLog(ctx context.Context, lb log.Broadcast, minConfs uint32) {
	lggr := lsn.L.With(
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
		if !lsn.shouldProcessLog(ctx, lb) {
			return
		}
		lsn.RespCountMu.Lock()
		lsn.ResponseCount[v.RequestId]++
		lsn.BlockNumberToReqID.Insert(fulfilledReq{
			blockNumber: v.Raw.BlockNumber,
			reqID:       v.RequestId,
		})
		lsn.RespCountMu.Unlock()
		lsn.markLogAsConsumed(ctx, lb)
		return
	}

	req, err := lsn.Coordinator.ParseRandomnessRequest(lb.RawLog())
	if err != nil {
		lggr.Errorw("Failed to parse RandomnessRequest log", "err", err)
		if !lsn.shouldProcessLog(ctx, lb) {
			return
		}
		lsn.markLogAsConsumed(ctx, lb)
		return
	}

	confirmedAt := lsn.getConfirmedAt(req, minConfs)
	lsn.ReqsMu.Lock()
	lsn.Reqs = append(lsn.Reqs, request{
		confirmedAtBlock: confirmedAt,
		req:              req,
		lb:               lb,
		utcTimestamp:     time.Now().UTC(),
	})
	lsn.ReqAdded()
	lsn.ReqsMu.Unlock()
	lggr.Infow("Enqueued randomness request",
		"requestID", hex.EncodeToString(req.RequestID[:]),
		"requestJobID", hex.EncodeToString(req.JobID[:]),
		"keyHash", hex.EncodeToString(req.KeyHash[:]),
		"fee", req.Fee,
		"sender", req.Sender.Hex(),
		"txHash", lb.RawLog().TxHash)
}

func (lsn *Listener) shouldProcessLog(ctx context.Context, lb log.Broadcast) bool {
	consumed, err := lsn.Chain.LogBroadcaster().WasAlreadyConsumed(ctx, lb)
	if err != nil {
		lsn.L.Errorw("Could not determine if log was already consumed", "err", err, "txHash", lb.RawLog().TxHash)
		// Do not process, let lb resend it as a retry mechanism.
		return false
	}
	return !consumed
}

func (lsn *Listener) markLogAsConsumed(ctx context.Context, lb log.Broadcast) {
	err := lsn.Chain.LogBroadcaster().MarkConsumed(ctx, nil, lb)
	lsn.L.ErrorIf(err, fmt.Sprintf("Unable to mark log %v as consumed", lb.String()))
}

func (lsn *Listener) getConfirmedAt(req *solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest, minConfs uint32) uint64 {
	lsn.RespCountMu.Lock()
	defer lsn.RespCountMu.Unlock()
	newConfs := uint64(minConfs) * (1 << lsn.ResponseCount[req.RequestID])
	// We cap this at 200 because solidity only supports the most recent 256 blocks
	// in the contract so if it was older than that, fulfillments would start failing
	// without the blockhash store feeder. We use 200 to give the node plenty of time
	// to fulfill even on fast chains.
	if newConfs > 200 {
		newConfs = 200
	}
	if lsn.ResponseCount[req.RequestID] > 0 {
		lsn.L.Warnw("Duplicate request found after fulfillment, doubling incoming confirmations",
			"txHash", req.Raw.TxHash,
			"blockNumber", req.Raw.BlockNumber,
			"blockHash", req.Raw.BlockHash,
			"requestID", hex.EncodeToString(req.RequestID[:]),
			"newConfs", newConfs)
		vrfcommon.IncDupeReqs(lsn.Job.Name.ValueOrZero(), lsn.Job.ExternalJobID, vrfcommon.V1)
	}
	return req.Raw.BlockNumber + newConfs
}

// ProcessRequest attempts to process the VRF request. Returns true if successful, false otherwise.
func (lsn *Listener) ProcessRequest(ctx context.Context, req request) bool {
	// This check to see if the log was consumed needs to be in the same
	// goroutine as the mark consumed to avoid processing duplicates.
	if !lsn.shouldProcessLog(ctx, req.lb) {
		return true
	}

	lggr := lsn.L.With(
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
	m := mathutil.Max(req.confirmedAtBlock, lsn.getLatestHead()-5)
	ctx, cancel := context.WithTimeout(ctx, callbacksTimeout)
	defer cancel()
	callback, err := lsn.Coordinator.Callbacks(&bind.CallOpts{
		BlockNumber: big.NewInt(int64(m)),
		Context:     ctx,
	}, req.req.RequestID)
	if err != nil {
		lggr.Errorw("Unable to check if already fulfilled, processing anyways", "err", err)
	} else if utils.IsEmpty(callback.SeedAndBlockNum[:]) {
		// If seedAndBlockNumber is zero then the response has been fulfilled
		// and we should skip it
		lggr.Infow("Request already fulfilled")
		lsn.markLogAsConsumed(ctx, req.lb)
		return true
	}

	// Check if we can ignore the request due to its age.
	if time.Now().UTC().Sub(req.utcTimestamp) >= lsn.Job.VRFSpec.RequestTimeout {
		lggr.Infow("Request too old, dropping it")
		lsn.markLogAsConsumed(ctx, req.lb)
		return true
	}

	lggr.Infow("Processing log request")

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"databaseID":    lsn.Job.ID,
			"externalJobID": lsn.Job.ExternalJobID,
			"name":          lsn.Job.Name.ValueOrZero(),
			"publicKey":     lsn.Job.VRFSpec.PublicKey[:],
			"from":          lsn.fromAddresses(),
			"evmChainID":    lsn.Job.VRFSpec.EVMChainID.String(),
		},
		"jobRun": map[string]interface{}{
			"logBlockHash":   req.req.Raw.BlockHash[:],
			"logBlockNumber": req.req.Raw.BlockNumber,
			"logTxHash":      req.req.Raw.TxHash,
			"logTopics":      req.req.Raw.Topics,
			"logData":        req.req.Raw.Data,
		},
	})

	run := pipeline.NewRun(*lsn.Job.PipelineSpec, vars)
	// The VRF pipeline has no async tasks, so we don't need to check for `incomplete`
	if _, err = lsn.PipelineRunner.Run(ctx, run, true, func(tx sqlutil.DataSource) error {
		// Always mark consumed regardless of whether the proof failed or not.
		if err = lsn.Chain.LogBroadcaster().MarkConsumed(ctx, tx, req.lb); err != nil {
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
	vrfcommon.IncProcessedReqs(lsn.Job.Name.ValueOrZero(), lsn.Job.ExternalJobID, vrfcommon.V1)
	return true
}

// Close complies with job.Service
func (lsn *Listener) Close() error {
	return lsn.StopOnce("VRFListener", func() error {
		close(lsn.ChStop)
		<-lsn.WaitOnStop // Log Listener
		<-lsn.WaitOnStop // Head Listener
		return lsn.ReqLogs.Close()
	})
}

func (lsn *Listener) HandleLog(ctx context.Context, lb log.Broadcast) {
	if !lsn.Deduper.ShouldDeliver(lb.RawLog()) {
		lsn.L.Tracew("skipping duplicate log broadcast", "log", lb.RawLog())
		return
	}

	wasOverCapacity := lsn.ReqLogs.Deliver(lb)
	if wasOverCapacity {
		lsn.L.Error("log mailbox is over capacity - dropped the oldest log")
		vrfcommon.IncDroppedReqs(lsn.Job.Name.ValueOrZero(), lsn.Job.ExternalJobID, vrfcommon.V1, vrfcommon.ReasonMailboxSize)
	}
}

func (lsn *Listener) fromAddresses() []common.Address {
	var addresses []common.Address
	for _, a := range lsn.Job.VRFSpec.FromAddresses {
		addresses = append(addresses, a.Address())
	}
	return addresses
}

// Job complies with log.Listener
func (lsn *Listener) JobID() int32 {
	return lsn.Job.ID
}

// ReplayStartedCallback is called by the log broadcaster when a replay is about to start.
func (lsn *Listener) ReplayStartedCallback() {
	// Clear the log Deduper cache so that we don't incorrectly ignore logs that have been sent that
	// are already in the cache.
	lsn.Deduper.Clear()
}
