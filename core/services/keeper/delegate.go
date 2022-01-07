package keeper

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

// To make sure Delegate struct implements job.Delegate interface
var _ job.Delegate = (*Delegate)(nil)

type transmitter interface {
	CreateEthTransaction(newTx bulletprooftxmanager.NewTx, qopts ...pg.QOpt) (etx bulletprooftxmanager.EthTx, err error)
}

type Delegate struct {
	logger   logger.Logger
	db       *sqlx.DB
	jrm      job.ORM
	pr       pipeline.Runner
	chainSet evm.ChainSet
}

// NewDelegate is the constructor of Delegate
func NewDelegate(
	db *sqlx.DB,
	jrm job.ORM,
	pr pipeline.Runner,
	logger logger.Logger,
	chainSet evm.ChainSet,
) *Delegate {
	return &Delegate{
		logger:   logger,
		db:       db,
		jrm:      jrm,
		pr:       pr,
		chainSet: chainSet,
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
	chain, err := d.chainSet.Get(spec.KeeperSpec.EVMChainID.ToInt())
	if err != nil {
		return nil, err
	}

	contractAddress := spec.KeeperSpec.ContractAddress
	contract, err := keeper_registry_wrapper.NewKeeperRegistry(
		contractAddress.Address(),
		chain.Client(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create keeper registry contract wrapper")
	}
	strategy := bulletprooftxmanager.NewQueueingTxStrategy(spec.ExternalJobID, chain.Config().KeeperDefaultTransactionQueueDepth(), false)

	orm := NewORM(d.db, d.logger, chain.TxManager(), chain.Config(), strategy)

	svcLogger := d.logger.With(
		"jobID", spec.ID,
		"registryAddress", contractAddress.Hex(),
	)

	minIncomingConfirmations := chain.Config().MinIncomingConfirmations()
	if spec.KeeperSpec.MinIncomingConfirmations != nil {
		minIncomingConfirmations = *spec.KeeperSpec.MinIncomingConfirmations
	}

	registrySynchronizer := NewRegistrySynchronizer(RegistrySynchronizerOptions{
		Job:                      spec,
		Contract:                 contract,
		ORM:                      orm,
		JRM:                      d.jrm,
		LogBroadcaster:           chain.LogBroadcaster(),
		SyncInterval:             chain.Config().KeeperRegistrySyncInterval(),
		MinIncomingConfirmations: minIncomingConfirmations,
		Logger:                   svcLogger,
		SyncUpkeepQueueSize:      chain.Config().KeeperRegistrySyncUpkeepQueueSize(),
	})
	upkeepExecuter := NewUpkeepExecuter(
		spec,
		orm,
		d.pr,
		chain.Client(),
		chain.HeadBroadcaster(),
		chain.TxManager().GetGasEstimator(),
		svcLogger,
		chain.Config(),
	)

	return []job.Service{
		registrySynchronizer,
		upkeepExecuter,
	}, nil
}
