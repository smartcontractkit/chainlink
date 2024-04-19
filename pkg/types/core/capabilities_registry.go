package core

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

//go:generate mockery --quiet --name CapabilitiesRegistry --output ./mocks/ --case=underscore
type CapabilitiesRegistry interface {
	Get(ctx context.Context, ID string) (capabilities.BaseCapability, error)
	GetTrigger(ctx context.Context, ID string) (capabilities.TriggerCapability, error)
	GetAction(ctx context.Context, ID string) (capabilities.ActionCapability, error)
	GetConsensus(ctx context.Context, ID string) (capabilities.ConsensusCapability, error)
	GetTarget(ctx context.Context, ID string) (capabilities.TargetCapability, error)
	List(ctx context.Context) ([]capabilities.BaseCapability, error)
	Add(ctx context.Context, c capabilities.BaseCapability) error
}
