package vrf

import (
	"encoding/hex"
	"strings"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"

	"github.com/theodesp/go-heaps/pairing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

type Delegate struct {
	db   *gorm.DB
	pr   pipeline.Runner
	porm pipeline.ORM
	ks   *keystore.Master
	cc   evm.ChainCollection
}

//go:generate mockery --name GethKeyStore --output mocks/ --case=underscore

type GethKeyStore interface {
	GetRoundRobinAddress(addresses ...common.Address) (common.Address, error)
}

type Config interface {
	MinIncomingConfirmations() uint32
	EvmGasLimitDefault() uint64
}

func NewDelegate(
	db *gorm.DB,
	ks *keystore.Master,
	pr pipeline.Runner,
	porm pipeline.ORM,
	chainCollection evm.ChainCollection) *Delegate {
	return &Delegate{
		db:   db,
		ks:   ks,
		pr:   pr,
		porm: porm,
		cc:   chainCollection,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.VRF
}

func (d *Delegate) AfterJobCreated(spec job.Job)  {}
func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) ServicesForSpec(jb job.Job) ([]job.Service, error) {
	if jb.VRFSpec == nil || jb.PipelineSpec == nil {
		return nil, errors.Errorf("vrf.Delegate expects a VRFSpec and PipelineSpec to be present, got %+v", jb)
	}
	pl, err := jb.PipelineSpec.Pipeline()
	if err != nil {
		return nil, err
	}
	// TODO: Fix with https://app.clubhouse.io/chainlinklabs/story/14615/add-ability-to-set-chain-id-in-all-pipeline-tasks-that-interact-with-evm
	chain, err := d.cc.Default()
	if err != nil {
		return nil, err
	}
	coordinator, err := solidity_vrf_coordinator_interface.NewVRFCoordinator(jb.VRFSpec.CoordinatorAddress.Address(), chain.Client())
	if err != nil {
		return nil, err
	}
	coordinatorV2, err := vrf_coordinator_v2.NewVRFCoordinatorV2(jb.VRFSpec.CoordinatorAddress.Address(), chain.Client())
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

	vorm := keystore.NewVRFORM(d.db)
	for _, task := range pl.Tasks {
		if _, ok := task.(*pipeline.VRFTaskV2); ok {
			return []job.Service{&listenerV2{
				cfg:             chain.Config(),
				l:               *l,
				ethClient:       chain.Client(),
				logBroadcaster:  chain.LogBroadcaster(),
				headBroadcaster: chain.HeadBroadcaster(),
				db:              d.db,
				abi:             abiV2,
				coordinator:     coordinatorV2,
				txm:             chain.TxManager(),
				pipelineRunner:  d.pr,
				vorm:            vorm,
				vrfks:           d.ks.VRF(),
				gethks:          d.ks.Eth(),
				pipelineORM:     d.porm,
				job:             jb,
				mbLogs:          utils.NewMailbox(1000),
				chStop:          make(chan struct{}),
				waitOnStop:      make(chan struct{}),
			}}, nil
		}
		if _, ok := task.(*pipeline.VRFTask); ok {
			return []job.Service{&listenerV1{
				cfg:                chain.Config(),
				l:                  *l,
				headBroadcaster:    chain.HeadBroadcaster(),
				logBroadcaster:     chain.LogBroadcaster(),
				db:                 d.db,
				txm:                chain.TxManager(),
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
			}}, nil
		}
	}
	return nil, errors.New("invalid job spec expected a vrf task")
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
