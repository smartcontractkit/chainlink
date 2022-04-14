package keeper

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	registry1_1 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry1_2 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

type RegistryVersion int32

const (
	RegistryVersion_1_0 RegistryVersion = iota
	RegistryVersion_1_1
	RegistryVersion_1_2
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

func (Delegate) AfterJobCreated(spec job.Job) {}

func (Delegate) BeforeJobDeleted(spec job.Job) {}

// ServicesForSpec satisfies the job.Delegate interface.
func (d *Delegate) ServicesForSpec(spec job.Job) (services []job.ServiceCtx, err error) {
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
	registryAddress := spec.KeeperSpec.ContractAddress
	contract1_1, err := registry1_1.NewKeeperRegistry(
		registryAddress.Address(),
		chain.Client(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create keeper registry 1_1 contract wrapper")
	}
	contract1_2, err := registry1_2.NewKeeperRegistry(
		registryAddress.Address(),
		chain.Client(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create keeper registry 1_2 contract wrapper")
	}

	strategy := txmgr.NewQueueingTxStrategy(spec.ExternalJobID, chain.Config().KeeperDefaultTransactionQueueDepth())
	orm := NewORM(d.db, d.logger, chain.Config(), strategy)
	svcLogger := d.logger.With(
		"jobID", spec.ID,
		"registryAddress", registryAddress.Hex(),
	)

	registryVersion, err := getRegistryVersion(contract1_1)
	if err != nil {
		return nil, errors.Wrap(err, "unable to determine version of keeper registry contract")
	}
	svcLogger.Debug("Registry version is: ", *registryVersion)

	minIncomingConfirmations := chain.Config().MinIncomingConfirmations()
	if spec.KeeperSpec.MinIncomingConfirmations != nil {
		minIncomingConfirmations = *spec.KeeperSpec.MinIncomingConfirmations
	}

	registrySynchronizer := NewRegistrySynchronizer(RegistrySynchronizerOptions{
		Job:                      spec,
		Contract1_1:              contract1_1,
		Contract1_2:              contract1_2,
		Version:                  *registryVersion,
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

	return []job.ServiceCtx{
		registrySynchronizer,
		upkeepExecuter,
	}, nil
}

func getRegistryVersion(contract1_1 *registry1_1.KeeperRegistry) (*RegistryVersion, error) {
	// Use registry 1_1 wrapper to get version information
	typeAndVersion, err := contract1_1.TypeAndVersion(nil)
	if err != nil {
		// Version 1.0 does not support typeAndVersion interface, hence gives an error on this call
		version := RegistryVersion_1_0
		return &version, nil
	}
	switch typeAndVersion {
	case "KeeperRegistry 1.1.0":
		version := RegistryVersion_1_1
		return &version, nil
	case "KeeperRegistry 1.2.0":
		version := RegistryVersion_1_2
		return &version, nil
	default:
		return nil, errors.Errorf("Registry version %s not supported", typeAndVersion)
	}
}
