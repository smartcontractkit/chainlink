package vrf

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/theodesp/go-heaps/pairing"
	"gorm.io/gorm"
)

var (
	_ log.Listener = &listenerV3{}
	_ job.Service  = &listenerV3{}
)

type listenerV3 struct {
	utils.StartStopOnce
	cfg             Config
	l               logger.Logger
	abi             abi.ABI
	ethClient       eth.Client
	logBroadcaster  log.Broadcaster
	txm             bulletprooftxmanager.TxManager
	headBroadcaster httypes.HeadBroadcasterRegistry
	coordinator     *vrf_coordinator_v2.VRFCoordinatorV2
	pipelineRunner  pipeline.Runner
	pipelineORM     pipeline.ORM
	vorm            keystore.VRFORM
	job             job.Job
	db              *gorm.DB
	vrfks           keystore.VRF
	gethks          keystore.Eth
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
	reqs     []pendingRequest
	reqAdded func() // A simple debug helper

	// Data structures for reorg attack protection
	// We want a map so we can do an O(1) count update every fulfillment log we get.
	respCountMu sync.Mutex
	respCount   map[string]uint64
	// This auxiliary heap is to used when we need to purge the
	// respCount map - we repeatedly want remove the minimum log.
	// You could use a sorted list if the completed logs arrive in order, but they may not.
	blockNumberToReqID *pairing.PairHeap
}

func (lsn *listenerV3) Start() error {
	return lsn.StartOnce("VRFListenerV2", func() error {
		// Take the larger of the global vs specific.
		// Note that the v2 vrf requests specify their own confirmation requirements.
		// We wait for max(minConfs, request required confs) to be safe.
		minConfs := lsn.cfg.MinIncomingConfirmations()
		if lsn.job.VRFSpec.Confirmations > lsn.cfg.MinIncomingConfirmations() {
			minConfs = lsn.job.VRFSpec.Confirmations
		}
		unsubscribeLogs := lsn.logBroadcaster.Register(lsn, log.ListenerOpts{
			Contract: lsn.coordinator.Address(),
			ParseLog: lsn.coordinator.ParseLog,
			LogsWithTopics: map[common.Hash][][]log.Topic{
				vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{}.Topic(): {
					{
						log.Topic(lsn.job.VRFSpec.PublicKey.MustHash()),
					},
				},
			},
			// Do not specify min confirmations, as it varies from request to request.
		})

		// Log listener gathers request logs
		go gracefulpanic.WrapRecover(func() {
			lsn.runLogListener([]func(){unsubscribeLogs}, minConfs)
		})
		// Request handler periodically computes a set of logs which can be fulfilled.
		go gracefulpanic.WrapRecover(func() {
			lsn.runRequestHandler()
		})
		return nil
	})
}

func (lsn *listenerV3) Connect(head *eth.Head) error {
	lsn.latestHead = uint64(head.Number)
	return nil
}

// Returns all the confirmed logs from
// the pending queue.
func (lsn *listenerV3) getConfirmedLogs(latestHead uint64) []pendingRequest {
	lsn.reqsMu.Lock()
	defer lsn.reqsMu.Unlock()
	var toProcess []pendingRequest
	for i := 0; i < len(lsn.reqs); i++ {
		if lsn.reqs[i].confirmedAtBlock <= latestHead {
			toProcess = append(toProcess, lsn.reqs[i])
		}
	}
	return toProcess
}

// Removes and returns all the confirmed logs from
// the pending queue.
func (lsn *listenerV3) extractConfirmedLogs() []pendingRequest {
	lsn.reqsMu.Lock()
	defer lsn.reqsMu.Unlock()
	var toProcess, toKeep []pendingRequest
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

func (lsn *listenerV3) getLatestHead() uint64 {
	lsn.latestHeadMu.RLock()
	defer lsn.latestHeadMu.RUnlock()
	return lsn.latestHead
}

// Remove all entries 10000 blocks or older
// to avoid a memory leak.
func (lsn *listenerV3) pruneConfirmedRequestCounts() {
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

// Every tick, we want to determine a set of logs that are confirmed
// and the subscription has sufficient balance to fulfill,
// given a eth call with the max gas price.
// Note we have to consider the pending reqs already in the bptxm as already "spent" link,
// using a max link consumed in their metadata.
// A user will need a minBalance capable of fulfilling a single req at the max gas price or nothing will happen.
// This is acceptable as users can choose different keyhashes which have different max gas prices.
func (lsn *listenerV3) runRequestHandler() {
	// TODO: Probably would have to be a configuration parameter per job so chains could have faster ones
	tick := time.NewTicker(2 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-lsn.chStop:
			lsn.waitOnStop <- struct{}{}
			return
		case <-tick.C:
			latestHead, err := lsn.ethClient.HeaderByNumber(context.Background(), nil)
			lsn.l.Infow("VRFListenerV2", "head", latestHead.Number, "err", err, "confs", lsn.job.VRFSpec.Confirmations)
			if err != nil {
				continue
			}
			confirmed := lsn.getConfirmedLogs(latestHead.Number.Uint64())
			// TODO: Group these by subscription, for now assume all the same sub
			// TODO: also probably want to order these by request time so we service oldest first
			// Get subscription balance. Note that outside of this request handler, this can only decrease while there
			// are no pending requests
			// The reqConf period is how long we have between doing the sub read
			// and enqueuing txes for the bptxm to ensure there's no race.
			if len(confirmed) == 0 {
				continue
			}
			sub, err := lsn.coordinator.GetSubscription(nil, confirmed[0].req.SubId)
			if err != nil {
				// TODO: log error
				continue
			}
			keys, err := lsn.gethks.SendingKeys()
			if err != nil {
				// TODO: log error
				continue
			}
			fromAddress := keys[0].Address
			if lsn.job.VRFSpec.FromAddress != nil {
				fromAddress = *lsn.job.VRFSpec.FromAddress
			}
			maxGasPrice := lsn.cfg.KeySpecificMaxGasPriceWei(fromAddress.Address())
			startBalance := sub.Balance
			// TODO: test this
			var reservedLink string
			err = lsn.db.Raw(`SELECT SUM(CAST(meta->>'MaxLink' AS NUMERIC(78, 0))) 
				FROM eth_txes
				WHERE meta->>'MaxLink' IS NOT NULL
				GROUP BY from_address = ?`, fromAddress).Scan(&reservedLink).Error
			if err != nil {
				lsn.l.Errorw("VRFListenerV2", "err", err)
				continue
			}
			lsn.l.Infow("VRFListenerV2", "reserved link", reservedLink)

			var processed []pendingRequest
			for _, req := range confirmed {
				vars := pipeline.NewVarsFrom(map[string]interface{}{
					"jobSpec": map[string]interface{}{
						"databaseID":    lsn.job.ID,
						"externalJobID": lsn.job.ExternalJobID,
						"name":          lsn.job.Name.ValueOrZero(),
						"publicKey":     lsn.job.VRFSpec.PublicKey[:],
						"maxGasPrice":   maxGasPrice.String(),
					},
					"jobRun": map[string]interface{}{
						"logBlockHash":   req.req.Raw.BlockHash[:],
						"logBlockNumber": req.req.Raw.BlockNumber,
						"logTxHash":      req.req.Raw.TxHash,
						"logTopics":      req.req.Raw.Topics,
						"logData":        req.req.Raw.Data,
					},
				})
				run, trrs, err := lsn.pipelineRunner.ExecuteRun(context.Background(), *lsn.job.PipelineSpec, vars, lsn.l)
				if err != nil {
					logger.Errorw("VRFListenerV2: failed executing run", "err", err)
					continue
				}
				// The call task will fail if there are insufficient funds
				if run.Errors.HasError() {
					logger.Errorw("VRFListenerV2: run errored", "err", err, "max gas price", maxGasPrice)
					continue
				}
				if len(trrs.FinalResult().Values) != 1 {
					logger.Errorw("VRFListenerV2: unexpected number of outputs", "err", err)
					continue
				}
				// Run succeeded, we expect a byte array representing the billing amount
				b, ok := trrs.FinalResult().Values[0].([]uint8)
				if !ok {
					logger.Errorw("VRFListenerV2: unexpected type", "err", err)
					continue
				}
				bi := utils.HexToBig(hexutil.Encode(b)[2:])
				lsn.l.Infow("VRFListenerV2: run output",
					"out", bi.String(),
					"start", startBalance)
				if startBalance.Cmp(bi) > 0 {
					// We have enough balance to service it, lets enqueue for bptxm
					err = postgres.NewGormTransactionManager(lsn.db).Transact(func(ctx context.Context) error {
						tx := postgres.TxFromContext(ctx, lsn.db)
						_, err := lsn.pipelineRunner.InsertFinishedRun(postgres.UnwrapGorm(tx), run, true)
						if err != nil {
							return err
						}
						if err = lsn.logBroadcaster.MarkConsumed(tx, req.lb); err != nil {
							return err
						}
						var (
							payload  string
							gaslimit uint64
						)
						for _, trr := range trrs {
							if trr.Task.Type() == pipeline.TaskTypeVRFV2 {
								m := trr.Result.Value.(map[string]interface{})
								payload = m["output"].(string)
							}
							if trr.Task.Type() == pipeline.TaskTypeEstimateGasLimit {
								gaslimit = trr.Result.Value.(uint64)
							}
						}
						_, err = lsn.txm.CreateEthTransaction(tx, bulletprooftxmanager.NewTx{
							FromAddress:    fromAddress.Address(),
							ToAddress:      lsn.coordinator.Address(),
							EncodedPayload: hexutil.MustDecode(payload),
							GasLimit:       gaslimit,
							Meta: &bulletprooftxmanager.EthTxMeta{
								RequestID: common.BytesToHash(req.req.RequestId.Bytes()),
								MaxLink:   bi.String(),
							},
							MinConfirmations: null.Uint32From(uint32(lsn.cfg.MinRequiredOutgoingConfirmations())),
							Strategy:         bulletprooftxmanager.NewSendEveryStrategy(false), // We already simd
						})
						// TODO: maybe save the eth tx id somewhere to link it
						return err
					})
					if err != nil {
						// TODO: log error
						continue
					}
					// If we successfully enqueued for the bptxm, subtract that balance
					startBalance = startBalance.Sub(startBalance, bi)
					processed = append(processed, req)
				} else {
					// We don't, leave it pending for now.
					// Its possible that a subsequent req uses less and we can fulfill it.
					continue
				}
			}
			lsn.reqsMu.Lock()
			var toKeep []pendingRequest
			for _, req := range lsn.reqs {
				for _, confirmed := range processed {
					if confirmed.req.RequestId.String() != req.req.RequestId.String() {
						toKeep = append(toKeep, confirmed)
					}
				}
			}
			lsn.reqs = toKeep
			lsn.reqsMu.Unlock()
		}
	}

}

// Listen for new heads
func (lsn *listenerV3) runHeadListener(unsubscribe func()) {
	for {
		select {
		case <-lsn.chStop:
			unsubscribe()
			lsn.waitOnStop <- struct{}{}
			return
		case <-lsn.newHead:
			toProcess := lsn.extractConfirmedLogs()
			// From a given set of logs,
			for _, r := range toProcess {
				lsn.ProcessV2VRFRequest(r.req, r.lb)
			}
			lsn.pruneConfirmedRequestCounts()
		}
	}
}

func (lsn *listenerV3) runLogListener(unsubscribes []func(), minConfs uint32) {
	lsn.l.Infow("VRFListenerV2: listening for run requests",
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
					panic(fmt.Sprintf("VRFListenerV2: invariant violated, expected log.Broadcast got %T", i))
				}
				lsn.handleLog(lb, minConfs)
			}
		}
	}
}

func (lsn *listenerV3) shouldProcessLog(lb log.Broadcast) bool {
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	consumed, err := lsn.logBroadcaster.WasAlreadyConsumed(lsn.db.WithContext(ctx), lb)
	if err != nil {
		lsn.l.Errorw("VRFListenerV2: could not determine if log was already consumed", "error", err, "txHash", lb.RawLog().TxHash)
		// Do not process, let lb resend it as a retry mechanism.
		return false
	}
	return !consumed
}

func (lsn *listenerV3) getConfirmedAt(req *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested, minConfs uint32) uint64 {
	lsn.respCountMu.Lock()
	defer lsn.respCountMu.Unlock()
	newConfs := uint64(minConfs) * (1 << lsn.respCount[req.RequestId.String()])
	// We cap this at 200 because solidity only supports the most recent 256 blocks
	// in the contract so if it was older than that, fulfillments would start failing
	// without the blockhash store feeder. We use 200 to give the node plenty of time
	// to fulfill even on fast chains.
	if newConfs > 200 {
		newConfs = 200
	}
	if lsn.respCount[req.RequestId.String()] > 0 {
		lsn.l.Warnw("VRFListenerV2: duplicate request found after fulfillment, doubling incoming confirmations",
			"txHash", req.Raw.TxHash,
			"blockNumber", req.Raw.BlockNumber,
			"blockHash", req.Raw.BlockHash,
			"reqID", req.RequestId.String(),
			"newConfs", newConfs)
	}
	return req.Raw.BlockNumber + newConfs
}

func (lsn *listenerV3) handleLog(lb log.Broadcast, minConfs uint32) {
	if v, ok := lb.DecodedLog().(*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled); ok {
		lsn.l.Infow("Received fulfilled log", "reqID", v.RequestId, "success", v.Success)
		if !lsn.shouldProcessLog(lb) {
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
		lsn.l.Errorw("VRFListenerV2: failed to parse log", "err", err, "txHash", lb.RawLog().TxHash)
		if !lsn.shouldProcessLog(lb) {
			return
		}
		lsn.markLogAsConsumed(lb)
		return
	}

	confirmedAt := lsn.getConfirmedAt(req, minConfs)
	lsn.reqsMu.Lock()
	lsn.reqs = append(lsn.reqs, pendingRequest{
		confirmedAtBlock: confirmedAt,
		req:              req,
		lb:               lb,
	})
	lsn.reqAdded()
	lsn.reqsMu.Unlock()
}

func (lsn *listenerV3) markLogAsConsumed(lb log.Broadcast) {
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	err := lsn.logBroadcaster.MarkConsumed(lsn.db.WithContext(ctx), lb)
	lsn.l.ErrorIf(errors.Wrapf(err, "VRFListenerV2: unable to mark log %v as consumed", lb.String()))
}

func (lsn *listenerV3) ProcessV2VRFRequest(req *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested, lb log.Broadcast) {
	// Check if the vrf req has already been fulfilled
	callback, err := lsn.coordinator.GetCommitment(nil, req.RequestId)
	if err != nil {
		lsn.l.Errorw("VRFListenerV2: unable to check if already fulfilled, processing anyways", "err", err, "txHash", req.Raw.TxHash)
	} else if utils.IsEmpty(callback[:]) {
		// If seedAndBlockNumber is zero then the response has been fulfilled
		// and we should skip it
		lsn.l.Infow("VRFListenerV2: request already fulfilled", "txHash", req.Raw.TxHash, "subID", req.SubId, "callback", callback)
		lsn.markLogAsConsumed(lb)
		return
	}

	lsn.l.Infow("VRFListenerV2: received log request",
		"log", lb.String(),
		"reqID", req.RequestId.String(),
		"txHash", req.Raw.TxHash,
		"blockNumber", req.Raw.BlockNumber,
		"blockHash", req.Raw.BlockHash,
		"seed", req.PreSeed)

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
	if _, err = lsn.pipelineRunner.Run(context.Background(), &run, lsn.l, true, func(tx *gorm.DB) error {
		// Always mark consumed regardless of whether the proof failed or not.
		if err = lsn.logBroadcaster.MarkConsumed(tx, lb); err != nil {
			logger.Errorw("VRFListenerV2: failed mark consumed", "err", err)
		}
		return nil
	}); err != nil {
		logger.Errorw("VRFListenerV2: failed executing run", "err", err)
	}
}

// Close complies with job.Service
func (lsn *listenerV3) Close() error {
	return lsn.StopOnce("VRFListenerV2", func() error {
		close(lsn.chStop)
		<-lsn.waitOnStop
		return nil
	})
}

func (lsn *listenerV3) HandleLog(lb log.Broadcast) {
	wasOverCapacity := lsn.reqLogs.Deliver(lb)
	if wasOverCapacity {
		logger.Error("VRFListenerV2: log mailbox is over capacity - dropped the oldest log")
	}
}

// Job complies with log.Listener
func (lsn *listenerV3) JobID() int32 {
	return lsn.job.ID
}
