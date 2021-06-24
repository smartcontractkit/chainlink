package keeper

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"gorm.io/gorm"
)

type transmitter interface {
	CreateEthTransaction(db *gorm.DB, fromAddress, toAddress common.Address, payload []byte, gasLimit uint64, meta interface{}) (etx models.EthTx, err error)
}

type Delegate struct {
	config          orm.ConfigReader
	orm             ORM
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
	config *orm.Config,
) *Delegate {
	return &Delegate{
		config:          config,
		orm:             NewORM(db, txm, config),
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

func (Delegate) OnJobCreated(spec job.Job) {}
func (Delegate) OnJobDeleted(spec job.Job) {}

func (d *Delegate) ServicesForSpec(spec job.Job) (services []job.Service, err error) {
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

	registrySynchronizer := NewRegistrySynchronizer(
		spec,
		contract,
		d.orm,
		d.jrm,
		d.logBroadcaster,
		d.config.KeeperRegistrySyncInterval(),
		d.config.KeeperMinimumRequiredConfirmations(),
	)
	upkeepExecuter := NewUpkeepExecuter(
		spec,
		d.orm,
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
