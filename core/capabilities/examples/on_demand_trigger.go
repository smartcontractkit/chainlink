package examples

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/mitchellh/mapstructure"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
)

// OnDemandTrigger is an example on-demand trigger.
type OnDemandTrigger struct {
	capabilities.CapabilityInfo
	chans map[string]chan capabilities.CapabilityResponse
	mu    sync.Mutex
}

// NewOnDemandTrigger returns a new OnDemandTrigger.
func NewOnDemandTrigger() *OnDemandTrigger {
	onDemandTriggerInfo, err := capabilities.NewCapabilityInfo(
		"on-demand-trigger",
		capabilities.CapabilityTypeTrigger,
		"An example on-demand trigger.",
		"v1.0.0",
	)

	if err != nil {
		panic(err)
	}

	return &OnDemandTrigger{
		CapabilityInfo: onDemandTriggerInfo,
		chans:          map[string]chan capabilities.CapabilityResponse{},
	}
}

func (o *OnDemandTrigger) FanOutEvent(ctx context.Context, event capabilities.CapabilityResponse) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, ch := range o.chans {
		ch <- event
	}
	return nil
}

// Execute executes the on-demand trigger.
func (o *OnDemandTrigger) SendEvent(ctx context.Context, wid string, event capabilities.CapabilityResponse) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	ch, ok := o.chans[wid]
	if !ok {
		return fmt.Errorf("no registration for %s", wid)
	}

	ch <- event
	return nil
}

type triggerRequest struct {
	WorkflowExecutionID string `mapstructure:"weid"`
}

func (o *OnDemandTrigger) Validate(inputs *values.Map) (*triggerRequest, error) {
	i := &triggerRequest{}
	m, err := inputs.Unwrap()
	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(m, i)
	if err != nil {
		return nil, err
	}

	fmt.Print(i, m, inputs)
	if i.WorkflowExecutionID == "" {
		return nil, errors.New("must provide workflow execution id")
	}

	return i, err

}

func (o *OnDemandTrigger) RegisterTrigger(ctx context.Context, callback chan capabilities.CapabilityResponse, inputs *values.Map) error {
	i, err := o.Validate(inputs)
	if err != nil {
		return err
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	o.chans[i.WorkflowExecutionID] = callback
	return nil
}

func (o *OnDemandTrigger) UnregisterTrigger(ctx context.Context, inputs *values.Map) error {
	i, err := o.Validate(inputs)
	if err != nil {
		return err
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	ch, ok := o.chans[i.WorkflowExecutionID]
	if ok {
		close(ch)
	}
	delete(o.chans, i.WorkflowExecutionID)
	return nil
}
