package vrf

import (
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
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
	lb   log.Broadcaster
	hb   httypes.HeadBroadcasterRegistry
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
	hb httypes.HeadBroadcasterRegistry,
	ec eth.Client,
	cfg Config) *Delegate {
	return &Delegate{
		cfg:  cfg,
		db:   db,
		txm:  txm,
		ks:   ks,
		pr:   pr,
		porm: porm,
		lb:   lb,
		ec:   ec,
		hb:   hb,
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

	vorm := keystore.NewVRFORM(d.db)
	return []job.Service{
		&listenerV1{
			cfg:            d.cfg,
			l:              *l,
			logBroadcaster: d.lb,
			db:             d.db,
			abi:            abi,
			coordinator:    coordinator,
			txm:            d.txm,
			pipelineRunner: d.pr,
			vorm:           vorm,
			vrfks:          d.ks.VRF(),
			gethks:         d.ks.Eth(),
			pipelineORM:    d.porm,
			job:            jb,
			mbLogs:         utils.NewMailbox(1000),
			chStop:         make(chan struct{}),
			waitOnStop:     make(chan struct{}),
		},
		&listenerV2{
			cfg:             d.cfg,
			l:               *l,
			ethClient:       d.ec,
			logBroadcaster:  d.lb,
			headBroadcaster: d.hb,
			db:              d.db,
			abi:             abiV2,
			coordinator:     coordinatorV2,
			txm:             d.txm,
			pipelineRunner:  d.pr,
			vorm:            vorm,
			vrfks:           d.ks.VRF(),
			gethks:          d.ks.Eth(),
			pipelineORM:     d.porm,
			job:             jb,
			mbLogs:          utils.NewMailbox(1000),
			chStop:          make(chan struct{}),
			waitOnStop:      make(chan struct{}),
		},
	}, nil
}

var (
	_ log.Listener = &listenerV1{}
	_ job.Service  = &listenerV1{}
	_ log.Listener = &listenerV2{}
	_ job.Service  = &listenerV2{}
)
