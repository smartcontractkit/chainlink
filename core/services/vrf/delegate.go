package vrf

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
	"time"

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
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
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

func (d *Delegate) OnJobCreated(spec job.Job) {}
func (d *Delegate) OnJobDeleted(spec job.Job) {}

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
		cfg:             d.cfg,
		l:               *l,
		headBroadcaster: d.hb,
		logBroadcaster:  d.lb,
		db:              d.db,
		txm:             d.txm,
		abi:             abi,
		coordinator:     coordinator,
		pipelineRunner:  d.pr,
		vorm:            vorm,
		vrfks:           d.ks.VRF(),
		gethks:          d.ks.Eth(),
		pipelineORM:     d.porm,
		job:             jb,
		reqLogs:         utils.NewMailbox(1000),
		chStop:          make(chan struct{}),
		waitOnStop:      make(chan struct{}),
		newHead:         make(chan struct{}, 1),
		respCount:       make(map[[32]byte]uint64),
	}
	return []job.Service{logListener}, nil
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
	respCount       map[[32]byte]uint64
	chStop          chan struct{}
	waitOnStop      chan struct{}
	newHead         chan struct{}
	// We can keep these pending logs in memory because we
	// only mark them confirmed once we send a corresponding fulfillment transaction.
	// So on node restart in the middle of processing, the lb will resend them.
	latestHead uint64
	reqsMu     sync.Mutex
	reqs       []request
	reqAdded   func()
}

func (lsn *listener) Connect(head *models.Head) error {
	lsn.latestHead = uint64(head.Number)
	return nil
}

// Note that we have 2 seconds to do this processing
func (lsn *listener) OnNewLongestChain(ctx context.Context, head models.Head) {
	lsn.latestHead = uint64(head.Number)
	select {
	case lsn.newHead <- struct{}{}:
	default:
	}
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
		unsubscribeHeadBroadcaster := lsn.headBroadcaster.Subscribe(lsn)
		go gracefulpanic.WrapRecover(func() {
			lsn.runLogListener([]func(){unsubscribeLogs}, minConfs)
		})
		go gracefulpanic.WrapRecover(func() {
			lsn.runHeadListener(unsubscribeHeadBroadcaster)
		})
		return nil
	})
}

// Listen for new heads
func (lsn *listener) runHeadListener(unsubscribe func()) {
	for {
		select {
		case <-lsn.chStop:
			unsubscribe()
			lsn.waitOnStop <- struct{}{}
		case <-lsn.newHead:
			var toProcess []request
			lsn.reqsMu.Lock()
			sort.Slice(lsn.reqs, func(i, j int) bool {
				return lsn.reqs[i].confirmedAtBlock < lsn.reqs[j].confirmedAtBlock
			})
			i := sort.Search(len(lsn.reqs), func(i int) bool {
				return lsn.reqs[i].confirmedAtBlock <= lsn.latestHead
			})
			if i < len(lsn.reqs) && lsn.reqs[i].confirmedAtBlock <= lsn.latestHead {
				toProcess = append(toProcess, lsn.reqs[:i+1]...)
				lsn.reqs = lsn.reqs[i+1:]
			}
			lsn.reqsMu.Unlock()
			for _, r := range toProcess {
				lsn.ProcessRequest(r.req, r.lb)
			}
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
				alreadyConsumed, err := lsn.logBroadcaster.WasAlreadyConsumed(lsn.db, lb)
				if err != nil {
					lsn.l.Errorw("VRFListener: could not determine if log was already consumed", "error", err, "txHash", lb.RawLog().TxHash)
					continue
				} else if alreadyConsumed {
					continue
				}
				if v, ok := lb.DecodedLog().(*solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequestFulfilled); ok {
					lsn.respCount[v.RequestId]++
					lsn.l.ErrorIf(lsn.logBroadcaster.MarkConsumed(lsn.db, lb), "failed to mark consumed")
					continue
				}
				req, err := lsn.coordinator.ParseRandomnessRequest(lb.RawLog())
				if err != nil {
					lsn.l.Errorw("VRFListener: failed to parse log", "err", err, "txHash", lb.RawLog().TxHash)
					lsn.l.ErrorIf(lsn.logBroadcaster.MarkConsumed(lsn.db, lb), "failed to mark consumed")
					continue
				}

				// Check if the vrf req has already been fulfilled
				callback, err := lsn.coordinator.Callbacks(nil, req.RequestID)
				if err != nil {
					lsn.l.Errorw("VRFListener: unable to check if already fulfilled, processing anyways", "err", err, "txHash", req.Raw.TxHash)
				} else if utils.IsEmpty(callback.SeedAndBlockNum[:]) {
					// If seedAndBlockNumber is zero then the response has been fulfilled
					// and we should skip it
					lsn.l.Infow("VRFListener: request already fulfilled", "txHash", req.Raw.TxHash)
					lsn.l.ErrorIf(lsn.logBroadcaster.MarkConsumed(lsn.db, lb), "failed to mark consumed")
					continue
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
		}
	}
}

func (lsn *listener) getConfirmedAt(req *solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest, minConfs uint32) uint64 {
	newConfs := uint64(minConfs) * (1 << lsn.respCount[req.RequestID])
	if newConfs > 200 {
		newConfs = 200
	}
	if lsn.respCount[req.RequestID] > 0 {
		lsn.l.Warn("VRFListener: duplicate request found after fulfillment, doubling incoming confirmations",
			"reqID", req.RequestID,
			"newConfs", newConfs)
	}
	return req.Raw.BlockNumber + uint64(minConfs)*(1<<lsn.respCount[req.RequestID])
}

func (lsn *listener) ProcessRequest(req *solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest, lb log.Broadcast) {
	s := time.Now()
	vrfCoordinatorPayload, req, err := lsn.ProcessLog(req, lb)
	f := time.Now()
	err = postgres.GormTransactionWithDefaultContext(lsn.db, func(tx *gorm.DB) error {
		if err == nil {
			// No errors processing the log, submit a transaction
			var etx bulletprooftxmanager.EthTx
			var from common.Address
			from, err = lsn.gethks.GetRoundRobinAddress()
			if err != nil {
				return err
			}
			etx, err = lsn.txm.CreateEthTransaction(tx,
				from,
				lsn.coordinator.Address(),
				vrfCoordinatorPayload,
				lsn.cfg.EthGasLimitDefault(),
				&models.EthTxMetaV2{
					JobID:         lsn.job.ID,
					RequestID:     req.RequestID,
					RequestTxHash: lb.RawLog().TxHash,
				},
				bulletprooftxmanager.SendEveryStrategy{},
			)
			if err != nil {
				return err
			}
			// TODO: Once we have eth tasks supported, we can use the pipeline directly
			// and be able to save errored proof generations. Until then only save
			// successful runs and log errors.
			_, err = lsn.pipelineRunner.InsertFinishedRun(tx, pipeline.Run{
				PipelineSpecID: lsn.job.PipelineSpecID,
				Errors:         []null.String{{}},
				Outputs: pipeline.JSONSerializable{
					Val: []interface{}{fmt.Sprintf("queued tx from %v to %v txdata %v",
						etx.FromAddress,
						etx.ToAddress,
						hex.EncodeToString(etx.EncodedPayload))},
				},
				Meta: pipeline.JSONSerializable{
					Val: map[string]interface{}{"eth_tx_id": etx.ID},
				},
				CreatedAt:  s,
				FinishedAt: &f,
			}, nil, false)
			if err != nil {
				return errors.Wrap(err, "VRFListener: failed to insert finished run")
			}
		}
		// Always mark consumed regardless of whether the proof failed or not.
		err = lsn.logBroadcaster.MarkConsumed(tx, lb)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		lsn.l.Errorw("VRFListener failed to save run", "err", err)
	}
}

func (lsn *listener) ProcessLog(req *solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest, lb log.Broadcast) ([]byte, *solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest, error) {
	lsn.l.Infow("VRFListener: received log request",
		"log", lb.String(),
		"reqID", hex.EncodeToString(req.RequestID[:]),
		"keyHash", hex.EncodeToString(req.KeyHash[:]),
		"txHash", req.Raw.TxHash,
		"blockNumber", req.Raw.BlockNumber,
		"seed", req.Seed,
		"fee", req.Fee)
	// Validate the key against the spec
	inputs, err := GetVRFInputs(lsn.job, req)
	if err != nil {
		lsn.l.Errorw("VRFListener: invalid log", "err", err)
		return nil, req, err
	}

	solidityProof, err := GenerateProofResponse(lsn.vrfks, inputs.pk, inputs.seed)
	if err != nil {
		lsn.l.Errorw("VRFListener: error generating proof", "err", err)
		return nil, req, err
	}

	vrfCoordinatorArgs, err := models.VRFFulfillMethod().Inputs.PackValues(
		[]interface{}{
			solidityProof[:], // geth expects slice, even if arg is constant-length
		})
	if err != nil {
		lsn.l.Errorw("VRFListener: error building fulfill args", "err", err)
		return nil, req, err
	}

	return append(lsn.abi.Methods["fulfillRandomnessRequest"].ID, vrfCoordinatorArgs...), req, nil
}

type VRFInputs struct {
	pk   secp256k1.PublicKey
	seed PreSeedData
}

// Check the key hash against the spec's pubkey
func GetVRFInputs(jb job.Job, request *solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest) (VRFInputs, error) {
	var inputs VRFInputs
	kh, err := jb.VRFSpec.PublicKey.Hash()
	if err != nil {
		return inputs, err
	}
	if !bytes.Equal(request.KeyHash[:], kh[:]) {
		return inputs, fmt.Errorf("invalid key hash %v expected %v", hex.EncodeToString(request.KeyHash[:]), hex.EncodeToString(kh[:]))
	}
	preSeed, err := BigToSeed(request.Seed)
	if err != nil {
		return inputs, errors.New("unable to parse preseed")
	}
	expectedJobID := jb.ExternalIDEncodeStringToTopic()
	if !bytes.Equal(expectedJobID[:], request.JobID[:]) {
		return inputs, fmt.Errorf("request jobID %v doesn't match expected %v", request.JobID[:], jb.ExternalIDEncodeStringToTopic().Bytes())
	}
	return VRFInputs{
		pk: jb.VRFSpec.PublicKey,
		seed: PreSeedData{
			PreSeed:   preSeed,
			BlockHash: request.Raw.BlockHash,
			BlockNum:  request.Raw.BlockNumber,
		},
	}, nil
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
