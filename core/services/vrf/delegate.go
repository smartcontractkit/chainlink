package vrf

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"

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
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

type Delegate struct {
	cfg    Config
	vorm   ORM
	db     *gorm.DB
	pr     pipeline.Runner
	porm   pipeline.ORM
	vrfks  *VRFKeyStore
	gethks GethKeyStore
	ec     eth.Client
	lb     log.Broadcaster
}

//go:generate mockery --name GethKeyStore --output mocks/ --case=underscore

type GethKeyStore interface {
	GetRoundRobinAddress(addresses ...common.Address) (common.Address, error)
}

type Config interface {
	MinIncomingConfirmations() uint32
	EthGasLimitDefault() uint64
	EthMaxQueuedTransactions() uint64
}

func NewDelegate(
	db *gorm.DB,
	vorm ORM,
	gethks GethKeyStore,
	vrfks *VRFKeyStore,
	pr pipeline.Runner,
	porm pipeline.ORM,
	lb log.Broadcaster,
	ec eth.Client,
	cfg Config) *Delegate {
	return &Delegate{
		cfg:    cfg,
		db:     db,
		vrfks:  vrfks,
		gethks: gethks,
		vorm:   vorm,
		pr:     pr,
		porm:   porm,
		lb:     lb,
		ec:     ec,
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
	coordinatorV2, err := vrf_coordinator_v2.NewVRFCoordinatorV2(jb.VRFSpec.CoordinatorAddress.Address(), d.ec)
	if err != nil {
		return nil, err
	}
	abi := eth.MustGetABI(solidity_vrf_coordinator_interface.VRFCoordinatorABI)
	abiV2 := eth.MustGetABI(vrf_coordinator_v2.VRFCoordinatorV2ABI)
	l := logger.CreateLogger(logger.Default.SugaredLogger.With(
		"jobID", jb.ID,
		"externalJobID", jb.ExternalJobID,
		"coordinatorAddress", jb.VRFSpec.CoordinatorAddress,
	))

	logListener := &listener{
		cfg:            d.cfg,
		l:              *l,
		logBroadcaster: d.lb,
		db:             d.db,
		abi:            abi,
		abiV2:          abiV2,
		coordinator:    coordinator,
		coordinatorV2:  coordinatorV2,
		pipelineRunner: d.pr,
		vorm:           d.vorm,
		vrfks:          d.vrfks,
		gethks:         d.gethks,
		pipelineORM:    d.porm,
		job:            jb,
		mbLogs:         utils.NewMailbox(1000),
		chStop:         make(chan struct{}),
		waitOnStop:     make(chan struct{}),
	}
	return []job.Service{logListener}, nil
}

var (
	_ log.Listener = &listener{}
	_ job.Service  = &listener{}
)

type listener struct {
	cfg            Config
	l              logger.Logger
	abi            abi.ABI
	abiV2          abi.ABI
	logBroadcaster log.Broadcaster
	coordinator    *solidity_vrf_coordinator_interface.VRFCoordinator
	coordinatorV2  *vrf_coordinator_v2.VRFCoordinatorV2
	pipelineRunner pipeline.Runner
	pipelineORM    pipeline.ORM
	vorm           ORM
	job            job.Job
	db             *gorm.DB
	vrfks          *VRFKeyStore
	gethks         GethKeyStore
	mbLogs         *utils.Mailbox
	chStop         chan struct{}
	waitOnStop     chan struct{}
	utils.StartStopOnce
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
		//unsubscribeLogs := lsn.logBroadcaster.Register(lsn, log.ListenerOpts{
		//	Contract: lsn.coordinator,
		//	LogsWithTopics: map[common.Hash][][]log.Topic{
		//		solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest{}.Topic(): {
		//			{
		//				log.Topic(lsn.job.ExternalIDToTopicHash()),
		//			},
		//		},
		//	},
		//	NumConfirmations: uint64(minConfs),
		//})
		unsubscribeLogsV2 := lsn.logBroadcaster.Register(lsn, log.ListenerOpts{
			Contract: lsn.coordinatorV2,
			LogsWithTopics: map[common.Hash][][]log.Topic{
				vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{}.Topic(): {
					{
						log.Topic(lsn.job.VRFSpec.PublicKey.MustHash()),
					},
				},
			},
		})
		go gracefulpanic.WrapRecover(func() {
			lsn.run([]func(){unsubscribeLogsV2}, minConfs)
		})
		return nil
	})
}

func (lsn *listener) run(unsubscribeLogs []func(), minConfs uint32) {
	lsn.l.Infow("VRFListener: listening for run requests",
		"maxUnconfirmed", lsn.cfg.EthMaxQueuedTransactions(),
		"gasLimit", lsn.cfg.EthGasLimitDefault(),
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
					panic(fmt.Sprintf("VRFListener: invariant violated, expected log.Broadcast got %T", i))
				}
				alreadyConsumed, err := lsn.logBroadcaster.WasAlreadyConsumed(lsn.db, lb)
				if err != nil {
					lsn.l.Errorw("VRFListener: could not determine if log was already consumed", "error", err, "txHash", lb.RawLog().TxHash)
					continue
				} else if alreadyConsumed {
					continue
				}
				// Try v1 first then v2
				switch lb.DecodedLog().(type) {
				case *solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest:
					lsn.ProcessV1VRFRequest(lb)
				case *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested:
					lsn.ProcessV2VRFRequest(lb)
				}
			}
		}
	}
}

func (lsn *listener) ProcessV2VRFRequest(lb log.Broadcast) {
	req, err := lsn.coordinatorV2.ParseRandomWordsRequested(lb.RawLog())
	if err != nil {
		lsn.l.Errorw("VRFListenerV2: failed to parse log", "err", err, "txHash", lb.RawLog().TxHash)
		lsn.l.ErrorIf(lsn.logBroadcaster.MarkConsumed(lsn.db, lb), "failed to mark consumed")
		return
	}

	// Check if the vrf req has already been fulfilled
	callback, err := lsn.coordinatorV2.SCallbacks(nil, req.PreSeed)
	if err != nil {
		lsn.l.Errorw("VRFListenerV2: unable to check if already fulfilled, processing anyways", "err", err, "txHash", req.Raw.TxHash)
	} else if utils.IsEmpty(callback.SeedAndBlockNum[:]) {
		// If seedAndBlockNumber is zero then the response has been fulfilled
		// and we should skip it
		lsn.l.Infow("VRFListenerV2: request already fulfilled", "txHash", req.Raw.TxHash, "subID", req.SubId, "callback", callback)
		lsn.l.ErrorIf(lsn.logBroadcaster.MarkConsumed(lsn.db, lb), "failed to mark consumed")
		return
	}

	s := time.Now()
	vrfCoordinatorPayload, _, err := lsn.ProcessLogV2(req, lb)
	f := time.Now()
	err = postgres.GormTransactionWithDefaultContext(lsn.db, func(tx *gorm.DB) error {
		if err == nil {
			// No errors processing the log, submit a transaction
			var etx models.EthTx
			var from common.Address
			from, err = lsn.gethks.GetRoundRobinAddress()
			if err != nil {
				return err
			}
			etx, err = bulletprooftxmanager.CreateEthTransaction(tx,
				from,
				lsn.coordinator.Address(),
				vrfCoordinatorPayload,
				lsn.cfg.EthGasLimitDefault(),
				lsn.cfg.EthMaxQueuedTransactions(),
				&models.EthTxMetaV2{
					JobID: lsn.job.ID,
					//RequestID:     req.PreSeed,
					RequestTxHash: lb.RawLog().TxHash,
				})
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

func (lsn *listener) ProcessV1VRFRequest(lb log.Broadcast) {
	req, err := lsn.coordinator.ParseRandomnessRequest(lb.RawLog())
	if err != nil {
		lsn.l.Errorw("VRFListener: failed to parse log", "err", err, "txHash", lb.RawLog().TxHash)
		lsn.l.ErrorIf(lsn.logBroadcaster.MarkConsumed(lsn.db, lb), "failed to mark consumed")
		return
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
		return
	}

	s := time.Now()
	vrfCoordinatorPayload, req, err := lsn.ProcessLog(req, lb)
	f := time.Now()
	err = postgres.GormTransactionWithDefaultContext(lsn.db, func(tx *gorm.DB) error {
		if err == nil {
			// No errors processing the log, submit a transaction
			var etx models.EthTx
			var from common.Address
			from, err = lsn.gethks.GetRoundRobinAddress()
			if err != nil {
				return err
			}
			etx, err = bulletprooftxmanager.CreateEthTransaction(tx,
				from,
				lsn.coordinator.Address(),
				vrfCoordinatorPayload,
				lsn.cfg.EthGasLimitDefault(),
				lsn.cfg.EthMaxQueuedTransactions(),
				&models.EthTxMetaV2{
					JobID:         lsn.job.ID,
					RequestID:     req.RequestID,
					RequestTxHash: lb.RawLog().TxHash,
				})
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

func (lsn *listener) ProcessLogV2(req *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested, lb log.Broadcast) ([]byte, *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested, error) {
	lsn.l.Infow("VRFListenerV2: received log request",
		"log", lb.String(),
		"reqID", req.PreSeed.String(),
		"keyHash", hex.EncodeToString(req.KeyHash[:]),
		"txHash", req.Raw.TxHash,
		"blockNumber", req.Raw.BlockNumber,
		"seed", req.PreSeed.String())
	// Validate the key against the spec
	kh, err := lsn.job.VRFSpec.PublicKey.Hash()
	if err != nil {
		return nil, nil, err
	}
	if !bytes.Equal(req.KeyHash[:], kh[:]) {
		return nil, nil, errors.New(fmt.Sprintf("invalid key hash %v expected %v", hex.EncodeToString(req.KeyHash[:]), hex.EncodeToString(kh[:])))
	}
	// uint256(keccak256(abi.encode(keyHash, msg.sender, nonce)))
	preSeed, err := BigToSeed(req.PreSeed)
	if err != nil {
		return nil, nil, errors.New("unable to parse preseed")
	}
	seed := PreSeedData{
		PreSeed:   preSeed,
		BlockHash: req.Raw.BlockHash,
		BlockNum:  req.Raw.BlockNumber,
	}
	lsn.l.Infow("generating proof", "pk", lsn.job.VRFSpec.PublicKey.String(), "seed", preSeed, "blockHash", req.Raw.BlockHash.String())
	solidityProof, err := lsn.vrfks.GenerateProof(lsn.job.VRFSpec.PublicKey, seed)
	if err != nil {
		lsn.l.Errorw("VRFListener: error generating proof", "err", err)
		return nil, nil, err
	}

	vrfCoordinatorArgs, err := lsn.abiV2.Methods["fulfillRandomWords"].Inputs.PackValues(
		[]interface{}{
			solidityProof[:], // geth expects slice, even if arg is constant-length
		})
	if err != nil {
		lsn.l.Errorw("VRFListener: error building fulfill args", "err", err)
		return nil, nil, err
	}

	return append(lsn.abiV2.Methods["fulfillRandomWords"].ID, vrfCoordinatorArgs...), nil, nil
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

	solidityProof, err := lsn.vrfks.GenerateProof(inputs.pk, inputs.seed)
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
		return inputs, errors.New(fmt.Sprintf("invalid key hash %v expected %v", hex.EncodeToString(request.KeyHash[:]), hex.EncodeToString(kh[:])))
	}
	preSeed, err := BigToSeed(request.Seed)
	if err != nil {
		return inputs, errors.New("unable to parse preseed")
	}
	expectedJobID := jb.ExternalIDToTopicHash()
	if !bytes.Equal(expectedJobID[:], request.JobID[:]) {
		return inputs, errors.New(fmt.Sprintf("request jobID %v doesn't match expected %v", request.JobID[:], jb.ExternalIDToTopicHash().Bytes()))
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
		<-lsn.waitOnStop
		return nil
	})
}

func (lsn *listener) HandleLog(lb log.Broadcast) {
	wasOverCapacity := lsn.mbLogs.Deliver(lb)
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
