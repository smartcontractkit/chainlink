package keeper

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

// To make sure Delegate struct implements job.Delegate interface
var _ job.Delegate = (*Delegate)(nil)

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

func (d *Delegate) BeforeJobCreated(spec job.Job) {}
func (d *Delegate) AfterJobCreated(spec job.Job)  {}
func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

// ServicesForSpec satisfies the job.Delegate interface.
func (d *Delegate) ServicesForSpec(spec job.Job) (services []job.ServiceCtx, err error) {
	if spec.KeeperSpec == nil {
		return nil, errors.Errorf("Delegate expects a *job.KeeperSpec to be present, got %v", spec)
	}
	chain, err := d.chainSet.Get(spec.KeeperSpec.EVMChainID.ToInt())
	if err != nil {
		return nil, err
	}
	registryAddress := spec.KeeperSpec.ContractAddress

	cfg := chain.Config()
	strategy := txmgr.NewQueueingTxStrategy(spec.ExternalJobID, cfg.KeeperDefaultTransactionQueueDepth(), cfg.DatabaseDefaultQueryTimeout())
	orm := NewORM(d.db, d.logger, chain.Config(), strategy)
	svcLogger := d.logger.With(
		"jobID", spec.ID,
		"registryAddress", registryAddress.Hex(),
	)

	registryWrapper, err := NewRegistryWrapper(registryAddress, chain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "unable to create keeper registry wrapper")
	}
	svcLogger.Info("Registry version is: ", registryWrapper.Version)

	minIncomingConfirmations := chain.Config().MinIncomingConfirmations()
	if spec.KeeperSpec.MinIncomingConfirmations != nil {
		minIncomingConfirmations = *spec.KeeperSpec.MinIncomingConfirmations
	}

	// effectiveKeeperAddress is the keeper address registered on the registry. This is by default the EOA account on the node.
	// In the case of forwarding, the keeper address is the forwarder contract deployed onchain between EOA and Registry.
	effectiveKeeperAddress := spec.KeeperSpec.FromAddress.Address()
	if spec.ForwardingAllowed {
		fwdrAddress, fwderr := chain.TxManager().GetForwarderForEOA(spec.KeeperSpec.FromAddress.Address())
		if fwderr == nil {
			effectiveKeeperAddress = fwdrAddress
		} else {
			svcLogger.Warnw("Skipping forwarding for job, will fallback to default behavior", "job", spec.Name, "err", fwderr)
		}
	}

	registrySynchronizer := NewRegistrySynchronizer(RegistrySynchronizerOptions{
		Job:                      spec,
		RegistryWrapper:          *registryWrapper,
		ORM:                      orm,
		JRM:                      d.jrm,
		LogBroadcaster:           chain.LogBroadcaster(),
		SyncInterval:             chain.Config().KeeperRegistrySyncInterval(),
		MinIncomingConfirmations: minIncomingConfirmations,
		Logger:                   svcLogger,
		SyncUpkeepQueueSize:      chain.Config().KeeperRegistrySyncUpkeepQueueSize(),
		EffectiveKeeperAddress:   effectiveKeeperAddress,
		newTurnEnabled:           chain.Config().KeeperTurnFlagEnabled(),
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
		effectiveKeeperAddress,
	)

	return []job.ServiceCtx{
		registrySynchronizer,
		upkeepExecuter,
	}, nil
}
