package test

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

const ConfigTOML = `[Foo]
Bar = "Baz"
`

var (

	//CapabilitiesRegistry
	GetID          = "get-id"
	GetTriggerID   = "get-trigger-id"
	GetActionID    = "get-action-id"
	GetConsensusID = "get-consensus-id"
	GetTargetID    = "get-target-id"
	CapabilityInfo = capabilities.CapabilityInfo{
		ID:             "capability-info-id@1.0.0",
		CapabilityType: 2,
		Description:    "capability-info-description",
	}
)

var _ capabilities.BaseCapability = (*baseCapability)(nil)

type baseCapability struct {
}

func (e baseCapability) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return CapabilityInfo, nil
}
