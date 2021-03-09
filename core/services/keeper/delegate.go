package keeper

import (
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"gorm.io/gorm"
)

type Delegate struct {
	keeperORM    KeeperORM
	ethClient    eth.Client
	syncInterval time.Duration
}

func NewDelegate(db *gorm.DB, ethClient eth.Client, config *orm.Config) *Delegate {
	keeperORM := NewORM(db)
	return &Delegate{
		keeperORM:    keeperORM,
		ethClient:    ethClient,
		syncInterval: config.KeeperRegistrySyncInterval(),
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
		return nil, err
	}

	registrySynchronizer := NewRegistrySynchronizer(spec, contract, d.keeperORM, d.syncInterval)

	return []job.Service{
		registrySynchronizer,
	}, nil
}
