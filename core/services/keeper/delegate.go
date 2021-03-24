package keeper

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"gorm.io/gorm"
)

type Delegate struct {
	config          *orm.Config
	db              *gorm.DB
	ethClient       eth.Client
	headBroadcaster *services.HeadBroadcaster
	logBroadcaster  log.Broadcaster
}

func NewDelegate(
	db *gorm.DB,
	ethClient eth.Client,
	headBroadcaster *services.HeadBroadcaster,
	logBroadcaster log.Broadcaster,
	config *orm.Config,
) *Delegate {
	return &Delegate{
		config:          config,
		db:              db,
		ethClient:       ethClient,
		headBroadcaster: headBroadcaster,
		logBroadcaster:  logBroadcaster,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.Keeper
}

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
		d.db,
		d.headBroadcaster,
		d.logBroadcaster,
		d.config.KeeperRegistrySyncInterval(),
		d.config.KeeperMinimumRequiredConfirmations(),
	)
	upkeepExecutor := NewUpkeepExecutor(
		spec,
		d.db,
		d.ethClient,
		d.headBroadcaster,
		d.config.KeeperMaximumGracePeriod(),
	)

	return []job.Service{
		registrySynchronizer,
		upkeepExecutor,
	}, nil
}
