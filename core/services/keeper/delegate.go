package keeper

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"gorm.io/gorm"
)

type transmitter interface {
	CreateEthTransaction(db *gorm.DB, fromAddress, toAddress common.Address, payload []byte, gasLimit uint64, meta interface{}, strategy bulletprooftxmanager.TxStrategy) (etx bulletprooftxmanager.EthTx, err error)
}

type Delegate struct {
	config          Config
	db              *gorm.DB
	txm             transmitter
	jrm             job.ORM
	pr              pipeline.Runner
	ethClient       eth.Client
	headBroadcaster httypes.HeadBroadcaster
	logBroadcaster  log.Broadcaster
}

var _ job.Delegate = (*Delegate)(nil)

func NewDelegate(
	db *gorm.DB,
	txm transmitter,
	jrm job.ORM,
	pr pipeline.Runner,
	ethClient eth.Client,
	headBroadcaster httypes.HeadBroadcaster,
	logBroadcaster log.Broadcaster,
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
	}
}

func (d *Delegate) JobType() job.Type {
	return job.Keeper
}

func (Delegate) AfterJobCreated(spec job.Job)  {}
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

	registrySynchronizer := NewRegistrySynchronizer(
		spec,
		contract,
		orm,
		d.jrm,
		d.logBroadcaster,
		d.config.KeeperRegistrySyncInterval(),
		d.config.KeeperMinimumRequiredConfirmations(),
	)
	upkeepExecuter := NewUpkeepExecuter(
		spec,
		orm,
		d.pr,
		d.ethClient,
		d.headBroadcaster,
		d.config,
	)

	return []job.Service{
		registrySynchronizer,
		upkeepExecuter,
	}, nil
}
