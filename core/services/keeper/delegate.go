package keeper

import (
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_contract"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

type Delegate struct {
	keeperORM KeeperORM
	ethClient eth.Client
}

func NewDelegate(orm *orm.ORM, ethClient eth.Client) *Delegate {
	keeperORM := NewORM(orm)
	return &Delegate{
		keeperORM: keeperORM,
		ethClient: ethClient,
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
	contract, err := keeper_registry_contract.NewKeeperRegistryContract(
		contractAddress.Address(),
		d.ethClient,
	)
	if err != nil {
		return nil, err
	}

	registrySynchronizer := NewRegistrySynchronizer(spec, contract, d.keeperORM, 10*time.Second) // TODO - RYAN

	return []job.Service{
		registrySynchronizer,
	}, nil
}
