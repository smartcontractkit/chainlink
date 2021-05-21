package vrf

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
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
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type Delegate struct {
	minIncomingConfs uint32
	vorm             ORM
	db               *gorm.DB
	pr               pipeline.Runner
	porm             pipeline.ORM
	ks               *VRFKeyStore
	ec               eth.Client
	lb               log.Broadcaster
}

func NewDelegate(minIncomingConfs uint32, params utils.ScryptParams, db *gorm.DB, pr pipeline.Runner, porm pipeline.ORM, lb log.Broadcaster, ec eth.Client) *Delegate {
	vorm := NewORM(db)
	ks := NewVRFKeyStore(vorm, params)
	return &Delegate{
		minIncomingConfs: minIncomingConfs,
		db:               db,
		ks:               ks,
		vorm:             vorm,
		pr:               pr,
		porm:             porm,
		lb:               lb,
		ec:               ec,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.VRF
}

func (d *Delegate) ServicesForSpec(jb job.Job) ([]job.Service, error) {
	if jb.VRFSpec == nil {
		return nil, errors.Errorf("vrf.Delegate expects a *job.VRFSpec to be present, got %v", jb)
	}
	coordinator, err := solidity_vrf_coordinator_interface.NewVRFCoordinator(jb.VRFSpec.CoordinatorAddress.Address(), d.ec)
	if err != nil {
		return nil, err
	}

	// Take the larger of the global vs specific
	minConfirmations := d.minIncomingConfs
	if jb.VRFSpec.Confirmations > d.minIncomingConfs {
		minConfirmations = jb.VRFSpec.Confirmations
	}

	logListener := &listener{
		logBroadcaster:   d.lb,
		db:               d.db,
		coordinator:      coordinator,
		pipelineRunner:   d.pr,
		ks:               d.ks,
		pipelineORM:      d.porm,
		job:              jb,
		mbLogs:           utils.NewMailbox(1000),
		minConfirmations: minConfirmations,
		chStop:           make(chan struct{}),
	}
	return []job.Service{logListener}, nil
}

var (
	_ log.Listener = &listener{}
	_ job.Service  = &listener{}
)

type listener struct {
	logBroadcaster   log.Broadcaster
	coordinator      *solidity_vrf_coordinator_interface.VRFCoordinator
	pipelineRunner   pipeline.Runner
	pipelineORM      pipeline.ORM
	vorm             ORM
	job              job.Job
	db               *gorm.DB
	ks               *VRFKeyStore
	mbLogs           *utils.Mailbox
	minConfirmations uint32
	chStop           chan struct{}
	utils.StartStopOnce
}

// Start complies with job.Service
func (l *listener) Start() error {
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
		go func() {
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
						lb, ok  := i.(log.Broadcast)
						if !ok {
							panic("blah")
						}
						alreadyConsumed, err := l.logBroadcaster.WasAlreadyConsumed(l.db, lb)
						if err != nil {
							logger.Errorw("VRFListener: could not determine if log was already consumed", "error", err)
							continue
						} else if alreadyConsumed {
							continue
						}

						s := time.Now()
						req, err := l.coordinator.ParseRandomnessRequest(lb.RawLog())
						if err != nil {
							logger.Error("VRFListener: invalid log")
							continue
						}
						// Validate the key against the spec
						inputs, err := GetVRFInputs(l.job, req)
						if err != nil {
							logger.Error("VRFListener: invalid log")
							continue
						}
						var re pipeline.RunErrors
						var output pipeline.JSONSerializable
						solidityProof, errGeneratingProof := l.ks.GenerateProof(inputs.pk, inputs.seed)
						if errGeneratingProof != nil {
							logger.Errorw("VRFListener: error generating proof", "err", errGeneratingProof)
							re = append(re, null.StringFrom(errGeneratingProof.Error()))
							output = pipeline.JSONSerializable{Null: true}
						} else {
							re = append(re, null.String{})
							output = pipeline.JSONSerializable{Val: solidityProof}
						}

						vrfCoordinatorArgs, err := models.VRFFulfillMethod().Inputs.PackValues(
							[]interface{}{
								solidityProof[:], // geth expects slice, even if arg is constant-length
							})
						if err != nil {
							//TODO
							continue
						}
						f := time.Now()
						err = postgres.GormTransactionWithDefaultContext(l.db, func(tx *gorm.DB) error {
							var etx *models.EthTx
							if errGeneratingProof != nil {
								etx, err = l.vorm.CreateEthTransaction(tx, common.Address{}, common.Address{}, vrfCoordinatorArgs, 0, 0)
								if err != nil {
									return err
								}
								_, err = l.pipelineRunner.InsertFinishedRun(tx, pipeline.Run{
									PipelineSpecID: l.job.PipelineSpecID,
									Meta: pipeline.JSONSerializable{
										Val: map[string]interface{}{"eth_tx_id": etx.ID},
									},
									Errors:     re,
									Outputs:    output,
									CreatedAt:  s,
									FinishedAt: &f,
								}, nil, false)
								if err != nil {
									return errors.Wrap(err, "VRFListener: failed to insert finished run")
								}
							} else {
								// Do not submit an eth tx, insert failure
								_, err = l.pipelineRunner.InsertFinishedRun(tx, pipeline.Run{
									PipelineSpecID: l.job.PipelineSpecID,
									Meta:           pipeline.JSONSerializable{Null: true},
									Errors:         re,
									Outputs:        output,
									CreatedAt:      s,
									FinishedAt:     &f,
								}, nil, false)
								if err != nil {
									return errors.Wrap(err, "VRFListener: failed to insert finished run")
								}
							}

							err = l.logBroadcaster.MarkConsumed(tx, lb)
							if err != nil {
								return err
							}
							return nil
						})
						if err != nil {
							logger.Error("VRFListener failed to save run", "err", err)
						}
					}
				}
			}
		}()
		return nil
	})
}

type VRFInputs struct {
	pk   secp256k1.PublicKey
	seed PreSeedData
}

func GetVRFInputs(jb job.Job, request *solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest) (VRFInputs, error) {
	// Check the key hash against the spec's pubkey, seed fields are not empty etc. etc.
	return VRFInputs{}, nil
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
