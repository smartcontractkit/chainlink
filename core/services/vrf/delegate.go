package vrf

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	heaps "github.com/theodesp/go-heaps"
	"github.com/theodesp/go-heaps/pairing"

	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"

	"gopkg.in/guregu/null.v4"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

type Delegate struct {
	cfg  Config
	db   *gorm.DB
	txm  bulletprooftxmanager.TxManager
	pr   pipeline.Runner
	porm pipeline.ORM
	ks   *keystore.Master
	ec   eth.Client
	hb   httypes.HeadBroadcasterRegistry
	lb   log.Broadcaster
}

//go:generate mockery --name GethKeyStore --output mocks/ --case=underscore

type GethKeyStore interface {
	GetRoundRobinAddress(addresses ...common.Address) (common.Address, error)
}

type Config interface {
	MinIncomingConfirmations() uint32
	EthGasLimitDefault() uint64
}

func NewDelegate(
	db *gorm.DB,
	txm bulletprooftxmanager.TxManager,
	ks *keystore.Master,
	pr pipeline.Runner,
	porm pipeline.ORM,
	lb log.Broadcaster,
	headBroadcaster httypes.HeadBroadcasterRegistry,
	ec eth.Client,
	cfg Config) *Delegate {
	return &Delegate{
		cfg:  cfg,
		db:   db,
		txm:  txm,
		ks:   ks,
		pr:   pr,
		porm: porm,
		hb:   headBroadcaster,
		lb:   lb,
		ec:   ec,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.VRF
}

func (d *Delegate) AfterJobCreated(spec job.Job)  {}
func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) ServicesForSpec(jb job.Job) ([]job.Service, error) {
	if jb.VRFSpec == nil {
		return nil, errors.Errorf("vrf.Delegate expects a *job.VRFSpec to be present, got %+v", jb)
	}
	coordinator, err := solidity_vrf_coordinator_interface.NewVRFCoordinator(jb.VRFSpec.CoordinatorAddress.Address(), d.ec)
	if err != nil {
		return nil, err
	}
	abi := eth.MustGetABI(solidity_vrf_coordinator_interface.VRFCoordinatorABI)
	l := logger.CreateLogger(logger.Default.SugaredLogger.With(
		"jobID", jb.ID,
		"externalJobID", jb.ExternalJobID,
		"coordinatorAddress", jb.VRFSpec.CoordinatorAddress,
	))

	vorm := keystore.NewVRFORM(d.db)

	logListener := &listener{
		cfg:                d.cfg,
		l:                  *l,
		headBroadcaster:    d.hb,
		logBroadcaster:     d.lb,
		db:                 d.db,
		txm:                d.txm,
		abi:                abi,
		coordinator:        coordinator,
		pipelineRunner:     d.pr,
		vorm:               vorm,
		vrfks:              d.ks.VRF(),
		gethks:             d.ks.Eth(),
		pipelineORM:        d.porm,
		job:                jb,
		reqLogs:            utils.NewMailbox(1000),
		chStop:             make(chan struct{}),
		waitOnStop:         make(chan struct{}),
		newHead:            make(chan struct{}, 1),
		respCount:          getStartingResponseCounts(d.db, l),
		blockNumberToReqID: pairing.New(),
		reqAdded:           func() {},
	}
	return []job.Service{logListener}, nil
}

func getStartingResponseCounts(db *gorm.DB, l *logger.Logger) map[[32]byte]uint64 {
	respCounts := make(map[[32]byte]uint64)
	var counts []struct {
		RequestID string
		Count     int
	}
	// Allow any state, not just confirmed, on purpose.
	// We assume once a ethtx is queued it will go through.
	err := db.Raw(`SELECT meta->'RequestID' AS request_id, count(meta->'RequestID') as count 
			FROM eth_txes 
			WHERE meta->'RequestID' IS NOT NULL 
		    GROUP BY meta->'RequestID'`).Scan(&counts).Error
	if err != nil {
		// Continue with an empty map, do not block job on this.
		l.Errorw("VRFListener: unable to read previous fulfillments", "err", err)
		return respCounts
	}
	for _, c := range counts {
		// Remove the quotes from the json
		req := strings.Replace(c.RequestID, `"`, ``, 2)
		// Remove the 0x prefix
		b, err := hex.DecodeString(req[2:])
		if err != nil {
			l.Errorw("VRFListener: unable to read fulfillment", "err", err, "reqID", c.RequestID)
			continue
		}
		var reqID [32]byte
		copy(reqID[:], b[:])
		respCounts[reqID] = uint64(c.Count)
	}
	return respCounts
}

var (
	_ log.Listener = &listener{}
	_ job.Service  = &listener{}
)

type request struct {
	confirmedAtBlock uint64
	req              *solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest
	lb               log.Broadcast
}

type listener struct {
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
	vrfks           *keystore.VRF
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

func (lsn *listener) Connect(head *models.Head) error {
	return nil
}

// Note that we have 2 seconds to do this processing
func (lsn *listener) OnNewLongestChain(_ context.Context, head models.Head) {
	lsn.setLatestHead(head)
	select {
	case lsn.newHead <- struct{}{}:
	default:
	}
}

func (lsn *listener) setLatestHead(h models.Head) {
	lsn.latestHeadMu.Lock()
	defer lsn.latestHeadMu.Unlock()
	num := uint64(h.Number)
	if num > lsn.latestHead {
		lsn.latestHead = num
	}
}

func (lsn *listener) getLatestHead() uint64 {
	lsn.latestHeadMu.RLock()
	defer lsn.latestHeadMu.RUnlock()
	return lsn.latestHead
}

// Start complies with job.Service
func (lsn *listener) Start() error {
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
		go gracefulpanic.WrapRecover(func() {
			lsn.runLogListener([]func(){unsubscribeLogs}, minConfs)
		})
		go gracefulpanic.WrapRecover(func() {
			lsn.runHeadListener(unsubscribeHeadBroadcaster)
		})
		return nil
	})
}

// Removes and returns all the confirmed logs from
// the pending queue.
func (lsn *listener) extractConfirmedLogs() []request {
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
func (lsn *listener) pruneConfirmedRequestCounts() {
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
func (lsn *listener) runHeadListener(unsubscribe func()) {
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

func (lsn *listener) runLogListener(unsubscribes []func(), minConfs uint32) {
	lsn.l.Infow("VRFListener: listening for run requests",
		"gasLimit", lsn.cfg.EthGasLimitDefault(),
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

func (lsn *listener) handleLog(lb log.Broadcast, minConfs uint32) {
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
		lsn.l.Errorw("VRFListener: failed to parse log", "err", err, "txHash", lb.RawLog().TxHash)
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

func (lsn *listener) shouldProcessLog(lb log.Broadcast) bool {
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()
	consumed, err := lsn.logBroadcaster.WasAlreadyConsumed(lsn.db.WithContext(ctx), lb)
	if err != nil {
		lsn.l.Errorw("VRFListener: could not determine if log was already consumed", "error", err, "txHash", lb.RawLog().TxHash)
		// Do not process, let lb resend it as a retry mechanism.
		return false
	}
	return !consumed
}

func (lsn *listener) markLogAsConsumed(lb log.Broadcast) {
	ctx, cancel := postgres.DefaultQueryCtx()
	defer cancel()

	err := lsn.logBroadcaster.MarkConsumed(lsn.db.WithContext(ctx), lb)
	lsn.l.ErrorIf(errors.Wrapf(err, "VRFListener: unable to mark log %v as consumed", lb.String()))
}

func (lsn *listener) getConfirmedAt(req *solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest, minConfs uint32) uint64 {
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
		lsn.l.Warnw("VRFListener: duplicate request found after fulfillment, doubling incoming confirmations",
			"txHash", req.Raw.TxHash,
			"blockNumber", req.Raw.BlockNumber,
			"blockHash", req.Raw.BlockHash,
			"reqID", hex.EncodeToString(req.RequestID[:]),
			"newConfs", newConfs)
	}
	return req.Raw.BlockNumber + uint64(minConfs)*(1<<lsn.respCount[req.RequestID])
}

func (lsn *listener) ProcessRequest(req *solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest, lb log.Broadcast) {
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
		lsn.l.Errorw("VRFListener: unable to check if already fulfilled, processing anyways", "err", err, "txHash", req.Raw.TxHash)
	} else if utils.IsEmpty(callback.SeedAndBlockNum[:]) {
		// If seedAndBlockNumber is zero then the response has been fulfilled
		// and we should skip it
		lsn.l.Infow("VRFListener: request already fulfilled", "txHash", req.Raw.TxHash, "reqID", req.RequestID)
		lsn.markLogAsConsumed(lb)
		return
	}

	s := time.Now()
	lsn.l.Infow("VRFListener: received log request",
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
	run, trrs, err := lsn.pipelineRunner.ExecuteRun(context.Background(), *lsn.job.PipelineSpec, vars, lsn.l)
	if err != nil {
		logger.Errorw("VRFListener: failed executing run", "err", err)
	}
	f := time.Now()
	err = postgres.GormTransactionWithDefaultContext(lsn.db, func(tx *gorm.DB) error {
		_, err = lsn.pipelineRunner.InsertFinishedRun(tx, pipeline.Run{
			State:          pipeline.RunStatusCompleted,
			PipelineSpecID: run.PipelineSpecID,
			Errors:         run.Errors,
			Outputs:        run.Outputs,
			Meta:           run.Meta,
			CreatedAt:      s,
			FinishedAt:     null.TimeFrom(f),
		}, trrs, true)
		if err != nil {
			return errors.Wrap(err, "VRFListener: failed to insert finished run")
		}
		// Always mark consumed regardless of whether the proof failed or not.
		return lsn.logBroadcaster.MarkConsumed(tx, lb)
	})
	if err != nil {
		lsn.l.Errorw("VRFListener failed to save run", "err", err)
	}
}

// Close complies with job.Service
func (lsn *listener) Close() error {
	return lsn.StopOnce("VRFListener", func() error {
		close(lsn.chStop)
		<-lsn.waitOnStop // Log listener
		<-lsn.waitOnStop // Head listener
		return nil
	})
}

func (lsn *listener) HandleLog(lb log.Broadcast) {
	wasOverCapacity := lsn.reqLogs.Deliver(lb)
	if wasOverCapacity {
		logger.Error("VRFListener: l mailbox is over capacity - dropped the oldest l")
	}
}

// JobID complies with log.Listener
func (*listener) JobID() models.JobID {
	return models.NilJobID
}

// Job complies with log.Listener
func (lsn *listener) JobIDV2() int32 {
	return lsn.job.ID
}

// IsV2Job complies with log.Listener
func (*listener) IsV2Job() bool {
	return true
}
