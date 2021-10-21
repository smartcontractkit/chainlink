package keeper

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

// To make sure Delegate struct implements job.Delegate interface
var _ job.Delegate = (*Delegate)(nil)

type transmitter interface {
	CreateEthTransaction(db *gorm.DB, newTx bulletprooftxmanager.NewTx) (etx bulletprooftxmanager.EthTx, err error)
}

type Delegate struct {
	config          Config
	logger          *logger.Logger
	db              *gorm.DB
	txm             transmitter
	jrm             job.ORM
	pr              pipeline.Runner
	ethClient       eth.Client
	headBroadcaster httypes.HeadBroadcaster
	logBroadcaster  log.Broadcaster
}

// NewDelegate is the constructor of Delegate
func NewDelegate(
	db *gorm.DB,
	txm transmitter,
	jrm job.ORM,
	pr pipeline.Runner,
	ethClient eth.Client,
	headBroadcaster httypes.HeadBroadcaster,
	logBroadcaster log.Broadcaster,
	logger *logger.Logger,
	config Config,
) *Delegate {
	return &Delegate{
		config:          config,
		db:              db,
		txm:             txm,
		jrm:             jrm,
		pr:              pr,
		ethClient:       ethClient,
		headBroadcaster: headBroadcaster,
		logBroadcaster:  logBroadcaster,
		logger:          logger,
	}
}

// JobType returns job type
func (d *Delegate) JobType() job.Type {
	return job.Keeper
}

func (Delegate) AfterJobCreated(spec job.Job) {}

func (Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) ServicesForSpec(spec job.Job) (services []job.Service, err error) {
	// TODO: we need to fill these out manually, find a better fix
	spec.PipelineSpec.JobName = spec.Name.ValueOrZero()
	spec.PipelineSpec.JobID = spec.ID

	if spec.KeeperSpec == nil {
		return nil, errors.Errorf("Delegate expects a *job.KeeperSpec to be present, got %v", spec)
	}

	contractAddress := spec.KeeperSpec.ContractAddress
	contract, err := keeper_registry_wrapper.NewKeeperRegistry(
		contractAddress.Address(),
		d.ethClient,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create keeper registry contract wrapper")
	}
	strategy := bulletprooftxmanager.NewQueueingTxStrategy(spec.ExternalJobID, d.config.KeeperDefaultTransactionQueueDepth())

	orm := NewORM(d.db, d.txm, d.config, strategy)

	svcLogger := d.logger.With(
		"jobID", spec.ID,
		"registryAddress", contractAddress.Hex(),
	)

	registrySynchronizer := NewRegistrySynchronizer(
		spec,
		contract,
		orm,
		d.jrm,
		d.logBroadcaster,
		d.config.KeeperRegistrySyncInterval(),
		d.config.KeeperMinimumRequiredConfirmations(),
		svcLogger.Named("RegistrySynchronizer"),
	)
	upkeepExecuter := NewUpkeepExecuter(
		spec,
		orm,
		d.pr,
		d.ethClient,
		d.headBroadcaster,
		svcLogger.Named("UpkeepExecuter"),
		d.config,
	)

	return []job.Service{
		registrySynchronizer,
		upkeepExecuter,
	}, nil
}
