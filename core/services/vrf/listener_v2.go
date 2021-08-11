package vrf

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

const (
	// Gas to be used
	GasAfterPaymentCalculation = 5000 + // subID balance update
		2100 + // cold subscription balance read
		20000 + // first time oracle balance update, note first time will be 20k, but 5k subsequently
		2*2100 - // cold read oracle address and oracle balance
		4800 + // request delete refund, note pre-london fork was 15k
		21000 + // base cost of the transaction
		7748 // Static costs of argument encoding etc. note that it varies by +/- x*12 for every x bytes of non-zero data in the proof.
)

var (
	_ log.Listener = &listenerV2{}
	_ job.Service  = &listenerV2{}
)

type pendingRequest struct {
	confirmedAtBlock uint64
	req              *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested
	lb               log.Broadcast
}

type listenerV2 struct {
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
	vrfks           *keystore.VRF
	gethks          *keystore.Eth
	mbLogs          *utils.Mailbox
	chStop          chan struct{}
	waitOnStop      chan struct{}
	latestHead      uint64
	// We can keep these pending logs in memory because we
	// only mark them confirmed once we send a corresponding fulfillment transaction.
	// So on node restart in the middle of processing, the lb will resend them.
	pendingLogs []pendingRequest
}

func (lsn *listenerV2) Start() error {
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

		// Subscribe to the head broadcaster for handling
		// per request conf requirements.
		_, unsubscribeHeadBroadcaster := lsn.headBroadcaster.Subscribe(lsn)

		go gracefulpanic.WrapRecover(func() {
			lsn.run([]func(){unsubscribeLogs, unsubscribeHeadBroadcaster}, minConfs)
		})
		return nil
	})
}

func (lsn *listenerV2) Connect(head *models.Head) error {
	lsn.latestHead = uint64(head.Number)
	return nil
}

func (lsn *listenerV2) OnNewLongestChain(ctx context.Context, head models.Head) {
	// Check if any v2 logs are ready for processing.
	lsn.latestHead = uint64(head.Number)
	var remainingLogs []pendingRequest
	for _, pl := range lsn.pendingLogs {
		if pl.confirmedAtBlock <= lsn.latestHead {
			// Note below makes API calls and opens a database transaction
			// TODO: Batch these requests in a follow up.
			lsn.ProcessV2VRFRequest(pl.req, pl.lb)
		} else {
			remainingLogs = append(remainingLogs, pl)
		}
	}
	lsn.pendingLogs = remainingLogs
}

func (lsn *listenerV2) run(unsubscribeLogs []func(), minConfs uint32) {
	lsn.l.Infow("VRFListenerV2: listening for run requests",
		"minConfs", minConfs)
	for {
		select {
		case <-lsn.chStop:
			for _, us := range unsubscribeLogs {
				us()
			}
			lsn.waitOnStop <- struct{}{}
			return
		case <-lsn.mbLogs.Notify():
			// Process all the logs in the queue if one is added
			for {
				i, exists := lsn.mbLogs.Retrieve()
				if !exists {
					break
				}
				lb, ok := i.(log.Broadcast)
				if !ok {
					panic(fmt.Sprintf("VRFListenerV2: invariant violated, expected log.Broadcast got %T", i))
				}
				alreadyConsumed, err := lsn.logBroadcaster.WasAlreadyConsumed(lsn.db, lb)
				if err != nil {
					lsn.l.Errorw("VRFListenerV2: could not determine if log was already consumed", "error", err, "txHash", lb.RawLog().TxHash)
					continue
				} else if alreadyConsumed {
					continue
				}
				req, err := lsn.coordinator.ParseRandomWordsRequested(lb.RawLog())
				if err != nil {
					lsn.l.Errorw("VRFListenerV2: failed to parse log", "err", err, "txHash", lb.RawLog().TxHash)
					lsn.markLogAsConsumed(lb)
					return
				}
				lsn.pendingLogs = append(lsn.pendingLogs, pendingRequest{
					confirmedAtBlock: req.Raw.BlockNumber + uint64(req.MinimumRequestConfirmations),
					req:              req,
					lb:               lb,
				})
			}
		}
	}
}

func (lsn *listenerV2) markLogAsConsumed(lb log.Broadcast) {
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	err := lsn.logBroadcaster.MarkConsumed(lsn.db.WithContext(ctx), lb)
	lsn.l.ErrorIf(errors.Wrapf(err, "VRFListenerV2: unable to mark log %v as consumed", lb.String()))
}

func (lsn *listenerV2) ProcessV2VRFRequest(req *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested, lb log.Broadcast) {
	// Check if the vrf req has already been fulfilled
	callback, err := lsn.coordinator.GetCommitment(nil, req.PreSeedAndRequestId)
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
		"reqID", req.PreSeedAndRequestId.String(),
		"keyHash", hex.EncodeToString(req.KeyHash[:]),
		"txHash", req.Raw.TxHash,
		"blockNumber", req.Raw.BlockNumber,
		"blockHash", req.Raw.BlockHash,
		"seed", req.PreSeedAndRequestId)

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
			logger.Errorw("VRFListener: failed mark consumed", "err", err)
		}
		return nil
	}); err != nil {
		logger.Errorw("VRFListener: failed executing run", "err", err)
	}
}

// Close complies with job.Service
func (lsn *listenerV2) Close() error {
	return lsn.StopOnce("VRFListenerV2", func() error {
		close(lsn.chStop)
		<-lsn.waitOnStop
		return nil
	})
}

func (lsn *listenerV2) HandleLog(lb log.Broadcast) {
	wasOverCapacity := lsn.mbLogs.Deliver(lb)
	if wasOverCapacity {
		logger.Error("VRFListenerV2: log mailbox is over capacity - dropped the oldest log")
	}
}

// Job complies with log.Listener
func (lsn *listenerV2) JobID() int32 {
	return lsn.job.ID
}
