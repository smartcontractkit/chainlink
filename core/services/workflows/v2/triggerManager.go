package v2

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

type triggerManager struct {
}

func (m *triggerManager) RegisterTrigger(ctx context.Context, trigger capabilities.TriggerCapability, request capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {

	return trigger.RegisterTrigger(ctx, request)
}

func (m *triggerManager) UnregisterTrigger(ctx context.Context, trigger capabilities.TriggerCapability, request capabilities.TriggerRegistrationRequest) error {
	return trigger.UnregisterTrigger(ctx, request)
}
