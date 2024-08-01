package core

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

type CapabilitiesRegistry interface {
	LocalNode(ctx context.Context) (capabilities.Node, error)
	ConfigForCapability(ctx context.Context, capabilityID string, donID uint32) (capabilities.CapabilityConfiguration, error)

	Get(ctx context.Context, ID string) (capabilities.BaseCapability, error)
	GetTrigger(ctx context.Context, ID string) (capabilities.TriggerCapability, error)
	GetAction(ctx context.Context, ID string) (capabilities.ActionCapability, error)
	GetConsensus(ctx context.Context, ID string) (capabilities.ConsensusCapability, error)
	GetTarget(ctx context.Context, ID string) (capabilities.TargetCapability, error)
	List(ctx context.Context) ([]capabilities.BaseCapability, error)
	Add(ctx context.Context, c capabilities.BaseCapability) error
}
