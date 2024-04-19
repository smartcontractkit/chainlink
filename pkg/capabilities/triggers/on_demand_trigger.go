package triggers

import (
	"context"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var info = capabilities.MustNewCapabilityInfo(
	"on-demand-trigger",
	capabilities.CapabilityTypeTrigger,
	"An example on-demand trigger.",
	"v1.0.0",
)

type workflowID string

type onDemandTriggerConfig struct{}

type OnDemand struct {
	log logger.Logger
	capabilities.Validator[onDemandTriggerConfig, any, capabilities.CapabilityResponse]
	capabilities.CapabilityInfo
	chans map[workflowID]chan<- capabilities.CapabilityResponse
	mu    sync.Mutex
}

var _ capabilities.TriggerCapability = (*OnDemand)(nil)

// Somewhat useless usage of validator, since
// the types we're validating against are essentially `any`.
//
// Leaving it in for completeness sake.
//
// Consider looking at MercuryTrigger for a more complete example.
var onDemandValidator = capabilities.NewValidator[onDemandTriggerConfig, any, capabilities.CapabilityResponse](capabilities.ValidatorArgs{Info: info})

// NewOnDemand creates a new on-demand trigger. The sendChannelBufferSize should be sized to ensure each client has sufficient
// time to process events, once this buffer is full new events for the client will be dropped.
func NewOnDemand(log logger.Logger) *OnDemand {
	return &OnDemand{
		log:            log,
		CapabilityInfo: info,
		Validator:      onDemandValidator,
		chans:          map[workflowID]chan<- capabilities.CapabilityResponse{},
	}
}

func (o *OnDemand) FanOutEvent(ctx context.Context, response capabilities.CapabilityResponse) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	for workFlowID, ch := range o.chans {
		select {
		case ch <- response:
		default:
			o.log.Warn("dropping event for workflowId %s due to slow consumer", workFlowID)
		}
	}
	return nil
}

// SendEvent sends an event to a specific workflowId. If the workflowId is not registered an error is returned.
func (o *OnDemand) SendEvent(ctx context.Context, wid string, event capabilities.CapabilityResponse) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	ch, ok := o.chans[workflowID(wid)]
	if !ok {
		return fmt.Errorf("no registration for %s", wid)
	}

	select {
	case ch <- event:
	default:
		return fmt.Errorf("event for workflowId %s not sent as send buffer is full", wid)
	}

	return nil
}

func (o *OnDemand) RegisterTrigger(ctx context.Context, req capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	wid := req.Metadata.WorkflowID
	o.mu.Lock()
	defer o.mu.Unlock()

	ch := make(chan capabilities.CapabilityResponse, defaultSendChannelBufferSize)
	o.chans[workflowID(wid)] = ch
	return ch, nil
}

func (o *OnDemand) UnregisterTrigger(ctx context.Context, req capabilities.CapabilityRequest) error {
	wid := req.Metadata.WorkflowID

	o.mu.Lock()
	defer o.mu.Unlock()

	ch, ok := o.chans[workflowID(wid)]
	if ok {
		close(ch)
	}
	delete(o.chans, workflowID(wid))
	return nil
}
