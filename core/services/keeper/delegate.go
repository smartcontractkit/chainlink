package keeper

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
)

type Delegate struct {
}

func NewDelegate() *Delegate {
	return &Delegate{}
}

func (d *Delegate) JobType() job.Type {
	return job.Keeper
}

func (d *Delegate) ServicesForSpec(spec job.Job) (services []job.Service, err error) {
	if spec.KeeperSpec == nil {
		return nil, errors.Errorf("Delegate expects a *job.KeeperSpec to be present, got %v", spec)
	}
	// TODO
	return nil, nil
}
