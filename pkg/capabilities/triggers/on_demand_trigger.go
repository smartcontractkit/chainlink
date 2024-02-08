package triggers

import (
	"context"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

var info = capabilities.MustNewCapabilityInfo(
	"on-demand-trigger",
	capabilities.CapabilityTypeTrigger,
	"An example on-demand trigger.",
	"v1.0.0",
)

type workflowID string

type OnDemand struct {
	capabilities.CapabilityInfo
	chans map[workflowID]chan<- capabilities.CapabilityResponse
	mu    sync.Mutex
}

var _ capabilities.TriggerCapability = (*OnDemand)(nil)

func NewOnDemand() *OnDemand {
	return &OnDemand{
		CapabilityInfo: info,
		chans:          map[workflowID]chan<- capabilities.CapabilityResponse{},
	}
}

func (o *OnDemand) FanOutEvent(ctx context.Context, event capabilities.CapabilityResponse) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, ch := range o.chans {
		ch <- event
	}
	return nil
}

// Execute executes the on-demand trigger.
func (o *OnDemand) SendEvent(ctx context.Context, wid string, event capabilities.CapabilityResponse) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	ch, ok := o.chans[workflowID(wid)]
	if !ok {
		return fmt.Errorf("no registration for %s", wid)
	}

	ch <- event
	return nil
}

func (o *OnDemand) RegisterTrigger(ctx context.Context, callback chan<- capabilities.CapabilityResponse, req capabilities.CapabilityRequest) error {
	wid := req.Metadata.WorkflowID

	o.mu.Lock()
	defer o.mu.Unlock()

	o.chans[workflowID(wid)] = callback
	return nil
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
