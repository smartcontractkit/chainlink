package consensus

import (
	"context"
	"fmt"
	"sync"

	"github.com/jonboulle/clockwork"
	"github.com/mitchellh/mapstructure"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

var info = capabilities.MustNewCapabilityInfo(
	"ocr3",
	capabilities.CapabilityTypeConsensus,
	"OCR3 consensus exposed as a capability.",
	"v1.0.0",
)

type capability struct {
	services.StateMachine
	capabilities.CapabilityInfo
	store  *store
	stopCh services.StopChan
	wg     sync.WaitGroup

	clock clockwork.Clock

	newExpiryWorkerCh chan *request
}

func newCapability(s *store, clock clockwork.Clock) *capability {
	o := &capability{
		CapabilityInfo:    info,
		store:             s,
		newExpiryWorkerCh: make(chan *request),
		clock:             clock,
		stopCh:            make(chan struct{}),
	}
	return o
}

func (o *capability) Start(ctx context.Context) error {
	return o.StartOnce("OCR3Capability", func() error {
		o.wg.Add(1)
		go o.loop()
		return nil
	})
}

func (o *capability) Close() error {
	return o.StopOnce("OCR3Capability", func() error {
		close(o.stopCh)
		o.wg.Wait()
		return nil
	})
}

func (o *capability) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (o *capability) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}

func (o *capability) Execute(ctx context.Context, callback chan<- capabilities.CapabilityResponse, request capabilities.CapabilityRequest) error {
	// Receives and stores an observation to do consensus on
	// Receives an aggregation method; at this point the method has been validated
	// Returns the consensus result over a channel
	r, err := o.unmarshalRequest(ctx, request, callback)
	if err != nil {
		return err
	}

	err = o.store.add(ctx, r)
	if err != nil {
		return err
	}

	o.newExpiryWorkerCh <- r
	return nil
}

func (o *capability) loop() {
	ctx, cancel := o.stopCh.NewCtx()
	defer cancel()
	defer o.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case r := <-o.newExpiryWorkerCh:
			o.wg.Add(1)
			go o.expiryWorker(ctx, r)
		}
	}
}

func (o *capability) expiryWorker(ctx context.Context, r *request) {
	defer o.wg.Done()

	d := r.ExpiresAt.Sub(o.clock.Now())
	tr := o.clock.NewTimer(d)
	defer tr.Stop()

	select {
	case <-ctx.Done():
		return
	case <-tr.Chan():
		wasPresent := o.store.evict(ctx, r.RequestID)
		if !wasPresent {
			// the item was already evicted,
			// we'll assume it was processed successfully
			// and return
			return
		}

		timeoutResp := capabilities.CapabilityResponse{
			Err: fmt.Errorf("timeout exceeded: could not process request before expiry %+v", r.RequestID),
		}

		select {
		case <-r.RequestCtx.Done():
		case r.CallbackCh <- timeoutResp:
			close(r.CallbackCh)
		}
	}
}

//nolint:unused
func (o *capability) response(ctx context.Context, resp response) error {
	req, err := o.store.get(ctx, resp.RequestID)
	if err != nil {
		return err
	}

	r := capabilities.CapabilityResponse{
		Value: resp.Value,
		Err:   resp.Err,
	}

	select {
	case <-req.RequestCtx.Done():
		return fmt.Errorf("request canceled: not propagating response %+v to caller", resp)
	case req.CallbackCh <- r:
		close(req.CallbackCh)
		o.store.evict(ctx, resp.RequestID)
		return nil
	}
}

func (o *capability) unmarshalRequest(ctx context.Context, r capabilities.CapabilityRequest, callback chan<- capabilities.CapabilityResponse) (*request, error) {
	valuesMap, err := r.Inputs.Unwrap()
	if err != nil {
		return nil, err
	}

	req := &request{
		RequestCtx:   ctx,
		CallbackCh:   callback,
		RequestID:    r.Metadata.WorkflowExecutionID,
		Observations: r.Inputs.Underlying["observations"],
	}
	err = mapstructure.Decode(valuesMap, req)
	if err != nil {
		return nil, err
	}

	configMap, err := r.Config.Unwrap()
	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(configMap, req)
	if err != nil {
		return nil, err
	}

	return req, err
}
