package vrf

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"time"

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

type Config struct {
	minIncomingConfs   uint32
	params             utils.ScryptParams
	gasLimit           uint64
	maxUnconfirmedTxes uint64
}

func NewConfig(minIncomingConfs uint32, params utils.ScryptParams, gasLimit uint64, maxUnconfirmedTxes uint64) Config {
	return Config{
		minIncomingConfs:   minIncomingConfs,
		params:             params,
		gasLimit:           gasLimit,
		maxUnconfirmedTxes: maxUnconfirmedTxes,
	}
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

func (d *Delegate) ServicesForSpec(jb job.Job) ([]job.Service, error) {
	if jb.VRFSpec == nil {
		return nil, errors.Errorf("vrf.Delegate expects a *job.VRFSpec to be present, got %+v", jb)
	}
	coordinator, err := solidity_vrf_coordinator_interface.NewVRFCoordinator(jb.VRFSpec.CoordinatorAddress.Address(), d.ec)
	if err != nil {
		return nil, err
	}

	// Take the larger of the global vs specific
	if jb.VRFSpec.Confirmations > d.cfg.minIncomingConfs {
		d.cfg.minIncomingConfs = jb.VRFSpec.Confirmations
	}

	logListener := &listener{
		cfg:            d.cfg,
		logBroadcaster: d.lb,
		db:             d.db,
		coordinator:    coordinator,
		pipelineRunner: d.pr,
		vorm:           d.vorm,
		vrfks:          d.vrfks,
		gethks:         d.gethks,
		pipelineORM:    d.porm,
		job:            jb,
		mbLogs:         utils.NewMailbox(1000),
		chStop:         make(chan struct{}),
	}
	return []job.Service{logListener}, nil
}

var (
	_ log.Listener = &listener{}
	_ job.Service  = &listener{}
)

type listener struct {
	cfg              Config
	logBroadcaster   log.Broadcaster
	coordinator      *solidity_vrf_coordinator_interface.VRFCoordinator
	pipelineRunner   pipeline.Runner
	pipelineORM      pipeline.ORM
	vorm             ORM
	job              job.Job
	db               *gorm.DB
	vrfks            *VRFKeyStore
	gethks           GethKeyStore
	mbLogs           *utils.Mailbox
	minConfirmations uint32
	chStop           chan struct{}
	utils.StartStopOnce
}

// Start complies with job.Service
func (l *listener) Start() error {
	cs := NewContractSubmitter()
	return l.StartOnce("VRFListener", func() error {
		unsubscribeLogs := l.logBroadcaster.Register(l, log.ListenerOpts{
			Contract: l.coordinator,
			LogsWithTopics: map[common.Hash][][]log.Topic{
				solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest{}.Topic(): {
					{
						log.Topic(l.job.ExternalIDToTopicHash()),
					},
				},
			},
			NumConfirmations: uint64(l.minConfirmations),
		})
		go gracefulpanic.WrapRecover(func() {
			l.run(unsubscribeLogs, cs)
		})
		return nil
	})
}

func (l *listener) run(unsubscribeLogs func(), submitter ContractSubmitter) {
	logger.Infow("VRFListener: listening for run requests", "jobTopic", l.job.ExternalIDToTopicHash())
	for {
		select {
		case <-l.chStop:
			unsubscribeLogs()
			return
		case <-l.mbLogs.Notify():
			// Process all the logs in the queue if one is added
			for {
				i, exists := l.mbLogs.Retrieve()
				if !exists {
					break
				}
				lb, ok := i.(log.Broadcast)
				if !ok {
					panic(fmt.Sprintf("VRFListener: invariant violated, expected log.Broadcast got %T", i))
				}
				alreadyConsumed, err := l.logBroadcaster.WasAlreadyConsumed(l.db, lb)
				if err != nil {
					logger.Errorw("VRFListener: could not determine if log was already consumed", "error", err)
					continue
				} else if alreadyConsumed {
					continue
				}
				s := time.Now()
				vrfCoordinatorPayload, req, err := l.ProcessLog(lb)
				f := time.Now()
				err = postgres.GormTransactionWithDefaultContext(l.db, func(tx *gorm.DB) error {
					if err == nil {
						// No errors processing the log, submit a transaction
						var etx *models.EthTx
						var from common.Address
						from, err = l.gethks.GetRoundRobinAddress() // TODO TX
						if err != nil {
							return err
						}
						etx, err = submitter.CreateEthTransaction(
							tx, models.EthTxMetaV2{
								JobID:         l.job.ID,
								RequestID:     req.RequestID,
								RequestTxHash: lb.RawLog().TxHash,
							},
							from, l.coordinator.Address(),
							vrfCoordinatorPayload,
							l.cfg.gasLimit, l.cfg.maxUnconfirmedTxes,
						)
						if err != nil {
							return err
						}
						// TODO: Once we have eth tasks supported, we can use the pipeline directly
						// and be able to save errored proof generations. Until then only save
						// successful runs and log errors.
						_, err = l.pipelineRunner.InsertFinishedRun(tx, pipeline.Run{
							PipelineSpecID: l.job.PipelineSpecID,
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
					err = l.logBroadcaster.MarkConsumed(tx, lb)
					if err != nil {
						return err
					}
					return nil
				})
				if err != nil {
					logger.Errorw("VRFListener failed to save run", "err", err)
				}
			}
		}
	}
}

func (l *listener) ProcessLog(lb log.Broadcast) ([]byte, *solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest, error) {
	req, err := l.coordinator.ParseRandomnessRequest(lb.RawLog())
	if err != nil {
		logger.Errorw("VRFListener: failed to parse log", "err", err)
		return nil, req, err
	}
	// Validate the key against the spec
	inputs, err := GetVRFInputs(l.job, req)
	if err != nil {
		logger.Errorw("VRFListener: invalid log", "err", err)
		return nil, req, err
	}

	solidityProof, err := l.vrfks.GenerateProof(inputs.pk, inputs.seed)
	if err != nil {
		logger.Errorw("VRFListener: error generating proof", "err", err)
		return nil, req, err
	}

	vrfCoordinatorArgs, err := models.VRFFulfillMethod().Inputs.PackValues(
		[]interface{}{
			solidityProof[:], // geth expects slice, even if arg is constant-length
		})
	if err != nil {
		logger.Errorw("VRFListener: error building fulfill args", "err", err)
		return nil, req, err
	}
	return vrfCoordinatorArgs, req, nil
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
func (l *listener) Close() error {
	return l.StopOnce("VRFListener", func() error {
		close(l.chStop)
		return nil
	})
}

func (l *listener) HandleLog(lb log.Broadcast) {
	wasOverCapacity := l.mbLogs.Deliver(lb)
	if wasOverCapacity {
		logger.Error("VRFListener: log mailbox is over capacity - dropped the oldest log")
	}
}

// JobID complies with log.Listener
func (*listener) JobID() models.JobID {
	return models.NilJobID
}

// Job complies with log.Listener
func (l *listener) JobIDV2() int32 {
	return l.job.ID
}

// IsV2Job complies with log.Listener
func (*listener) IsV2Job() bool {
	return true
}
