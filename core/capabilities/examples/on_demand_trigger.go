package examples

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
)

// OnDemandTrigger is an example on-demand trigger.
type OnDemandTrigger struct {
	capabilities.CapabilityInfo
	// map between Workflow Execution IDs and callback channels
}

// NewOnDemandTrigger returns a new OnDemandTrigger.
func NewOnDemandTrigger() (*OnDemandTrigger, error) {
	onDemandTriggerInfo, err := capabilities.NewCapabilityInfo(
		"on-demand-trigger",
		capabilities.CapabilityTypeTrigger,
		"An example on-demand trigger.",
		"v1.0.0",
	)

	if err != nil {
		return nil, err
	}

	return &OnDemandTrigger{
		CapabilityInfo: onDemandTriggerInfo,
	}, nil
}

// TODO SendMultipleEvents
// func (o *OnDemandTrigger) FanOutEvent(ctx context.Context, event capabilities.CapabilityResponse) error {
// 	o.eventsChannel <- event
// 	return nil
// }

// Execute executes the on-demand trigger.
func (o *OnDemandTrigger) SendEvent(ctx context.Context, event capabilities.CapabilityResponse) error {
	// TODO: Add a "workflow execution id" to be able to differentiate between recipient channels
	return nil
}

func (o *OnDemandTrigger) RegisterTrigger(ctx context.Context, callback chan capabilities.CapabilityResponse, inputs values.Map) error {
	// Validate inputs.
	// Get 1) a workflow execution id from the input so we can register a trigger
	// and 2)
	// Send any new events from the eventChannel to the newly register callback channel

	return nil
}

func (o *OnDemandTrigger) UnregisterTrigger(ctx context.Context, inputs values.Map) error {
	return nil
}
