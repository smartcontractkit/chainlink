package keeper

import (
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"gorm.io/gorm"
)

type Delegate struct {
	ethClient       eth.Client
	headBroadcaster *services.HeadBroadcaster
	db              *gorm.DB
	syncInterval    time.Duration
}

func NewDelegate(db *gorm.DB, ethClient eth.Client, headBroadcaster *services.HeadBroadcaster, config *orm.Config) *Delegate {
	return &Delegate{
		ethClient:       ethClient,
		headBroadcaster: headBroadcaster,
		db:              db,
		syncInterval:    config.KeeperRegistrySyncInterval(),
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

	registrySynchronizer := NewRegistrySynchronizer(spec, contract, d.db, d.syncInterval)
	upkeepExecutor := NewUpkeepExecutor(spec, d.db, d.ethClient, d.headBroadcaster)

	return []job.Service{
		registrySynchronizer,
		upkeepExecutor,
	}, nil
}
