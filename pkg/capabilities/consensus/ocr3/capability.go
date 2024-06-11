package ocr3

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const (
	ocrCapabilityID = "offchain_reporting@1.0.0"

	methodStartRequest = "start_request"
	methodSendResponse = "send_response"
	methodHeader       = "method"
	transmissionHeader = "transmission"
	terminateHeader    = "terminate"
)

var info = capabilities.MustNewCapabilityInfo(
	ocrCapabilityID,
	capabilities.CapabilityTypeConsensus,
	"OCR3 consensus exposed as a capability.",
)

type capability struct {
	services.StateMachine
	capabilities.CapabilityInfo
	capabilities.Validator[config, inputs, outputs]

	store  *store
	stopCh services.StopChan
	wg     sync.WaitGroup
	lggr   logger.Logger

	requestTimeout time.Duration
	clock          clockwork.Clock

	aggregatorFactory types.AggregatorFactory
	aggregators       map[string]types.Aggregator

	encoderFactory types.EncoderFactory
	encoders       map[string]types.Encoder

	transmitCh chan *outputs
	newTimerCh chan *request

	callbackChannelBufferSize int
}

var _ capabilityIface = (*capability)(nil)
var _ capabilities.ConsensusCapability = (*capability)(nil)
var ocr3CapabilityValidator = capabilities.NewValidator[config, inputs, outputs](capabilities.ValidatorArgs{Info: info})

func newCapability(s *store, clock clockwork.Clock, requestTimeout time.Duration, aggregatorFactory types.AggregatorFactory, encoderFactory types.EncoderFactory, lggr logger.Logger,
	callbackChannelBufferSize int) *capability {
	o := &capability{
		CapabilityInfo:    info,
		Validator:         ocr3CapabilityValidator,
		store:             s,
		clock:             clock,
		requestTimeout:    requestTimeout,
		stopCh:            make(chan struct{}),
		lggr:              logger.Named(lggr, "OCR3CapabilityClient"),
		aggregatorFactory: aggregatorFactory,
		aggregators:       map[string]types.Aggregator{},
		encoderFactory:    encoderFactory,
		encoders:          map[string]types.Encoder{},

		transmitCh: make(chan *outputs),
		newTimerCh: make(chan *request),

		callbackChannelBufferSize: callbackChannelBufferSize,
	}
	return o
}

func (o *capability) Start(ctx context.Context) error {
	return o.StartOnce("OCR3Capability", func() error {
		o.wg.Add(1)
		go o.worker()
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

func (o *capability) Name() string { return o.lggr.Name() }

func (o *capability) HealthReport() map[string]error {
	return map[string]error{o.Name(): o.Healthy()}
}

func (o *capability) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	c, err := o.ValidateConfig(request.Config)
	if err != nil {
		return err
	}

	agg, err := o.aggregatorFactory(c.AggregationMethod, *c.AggregationConfig, o.lggr)
	if err != nil {
		return err
	}
	o.aggregators[request.Metadata.WorkflowID] = agg

	encoder, err := o.encoderFactory(c.EncoderConfig)
	if err != nil {
		return err
	}
	o.encoders[request.Metadata.WorkflowID] = encoder
	return nil
}

func (o *capability) getAggregator(workflowID string) (types.Aggregator, error) {
	agg, ok := o.aggregators[workflowID]
	if !ok {
		return nil, fmt.Errorf("no aggregator found for workflowID %s", workflowID)
	}

	return agg, nil
}

func (o *capability) getEncoder(workflowID string) (types.Encoder, error) {
	enc, ok := o.encoders[workflowID]
	if !ok {
		return nil, fmt.Errorf("no aggregator found for workflowID %s", workflowID)
	}

	return enc, nil
}

func (o *capability) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	delete(o.aggregators, request.Metadata.WorkflowID)
	delete(o.encoders, request.Metadata.WorkflowID)
	return nil
}

// Execute enqueues a new consensus request, passing it to the reporting plugin as needed.
// IMPORTANT: OCR3 only exposes signatures via the contractTransmitter, which is located
// in a separate process to the reporting plugin LOOPP. However, only the reporting plugin
// LOOPP is able to transmit responses back to the workflow engine. As a workaround to this, we've implemented a custom contract transmitter which fetches this capability from the
// registry and calls Execute with the response, setting "method = `methodSendResponse`".
func (o *capability) Execute(ctx context.Context, r capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	m := struct {
		Method       string
		Transmission map[string]any
		Terminate    bool
	}{
		Method: methodStartRequest,
	}
	err := r.Inputs.UnwrapTo(&m)
	if err != nil {
		o.lggr.Warnf("could not unwrap method from CapabilityRequest, using default: %w", err)
	}

	switch m.Method {
	case methodSendResponse:
		inputs, err := values.NewMap(m.Transmission)
		if err != nil {
			return nil, fmt.Errorf("failed to create map for response inputs: %w", err)
		}
		o.lggr.Debugw("Execute - sending response", "workflowExecutionID", r.Metadata.WorkflowExecutionID, "inputs", inputs)
		var responseErr error
		if m.Terminate {
			o.lggr.Debugw("Execute - terminating execution", "workflowExecutionID", r.Metadata.WorkflowExecutionID)
			responseErr = capabilities.ErrStopExecution
		}
		out := &outputs{
			WorkflowExecutionID: r.Metadata.WorkflowExecutionID,
			CapabilityResponse: capabilities.CapabilityResponse{
				Value: inputs,
				Err:   responseErr,
			},
		}
		o.transmitCh <- out

		// Return a dummy response back to the caller
		// This allows the transmitter to block on a response before
		// returning from Transmit()
		// TODO(cedric): our current stream-based implementation for the Execute
		// returns immediately without waiting for the server-side to complete. This
		// breaks the API since callers can no longer rely on a non-error response
		// from Execute() serving as an acknowledgement that the request in being handled.
		callbackCh := make(chan capabilities.CapabilityResponse, 1)
		callbackCh <- capabilities.CapabilityResponse{}
		close(callbackCh)
		return callbackCh, nil
	case methodStartRequest:
		// Receives and stores an observation to do consensus on
		// Receives an aggregation method; at this point the method has been validated
		// Returns the consensus result over a channel
		inputs, err := o.ValidateInputs(r.Inputs)
		if err != nil {
			return nil, err
		}

		config, err := o.ValidateConfig(r.Config)
		if err != nil {
			return nil, err
		}

		return o.queueRequestForProcessing(ctx, r.Metadata, inputs, config)
	}

	return nil, fmt.Errorf("unknown method: %s", m.Method)
}

// queueRequestForProcessing queues a request for processing by the worker
// goroutine by adding the request to its store.
//
// When a request is queued, a timer is started to ensure that the request does not exceed its expiry time.
func (o *capability) queueRequestForProcessing(
	ctx context.Context,
	metadata capabilities.RequestMetadata,
	i *inputs,
	c *config,
) (<-chan capabilities.CapabilityResponse, error) {
	callbackCh := make(chan capabilities.CapabilityResponse, o.callbackChannelBufferSize)

	// Use the capability-level request timeout unless the request's config specifies
	// its own timeout, in which case we'll use that instead. This allows the workflow spec
	// to configure more granular timeouts depending on the circumstances.
	requestTimeout := o.requestTimeout
	if c.RequestTimeoutMS != 0 {
		requestTimeout = time.Duration(c.RequestTimeoutMS) * time.Millisecond
	}

	r := &request{
		StopCh:              make(chan struct{}),
		CallbackCh:          callbackCh,
		WorkflowExecutionID: metadata.WorkflowExecutionID,
		WorkflowID:          metadata.WorkflowID,
		WorkflowOwner:       metadata.WorkflowOwner,
		WorkflowName:        metadata.WorkflowName,
		ReportID:            c.ReportID,
		WorkflowDonID:       metadata.WorkflowDonID,
		Observations:        i.Observations,
		ExpiresAt:           o.clock.Now().Add(requestTimeout),
	}

	o.lggr.Debugw("Execute - adding to store", "workflowID", r.WorkflowID, "workflowExecutionID", r.WorkflowExecutionID, "observations", r.Observations)

	err := o.store.add(ctx, r)
	if err != nil {
		return nil, err
	}

	o.newTimerCh <- r
	return callbackCh, nil
}

func (o *capability) worker() {
	ctx, cancel := o.stopCh.NewCtx()
	defer cancel()
	defer o.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case r := <-o.newTimerCh:
			o.wg.Add(1)
			go o.expiryTimer(ctx, r)
		case resp := <-o.transmitCh:
			o.handleTransmitMsg(ctx, resp)
		}
	}
}

func (o *capability) handleTransmitMsg(ctx context.Context, resp *outputs) {
	req, wasPresent := o.store.evict(ctx, resp.WorkflowExecutionID)
	if !wasPresent {
		return
	}

	select {
	// This should only happen if the client has closed the upstream context.
	// In this case, the request is cancelled and we shouldn't transmit.
	case <-req.StopCh:
	case req.CallbackCh <- resp.CapabilityResponse:
		close(req.CallbackCh)
	}
}

func (o *capability) expiryTimer(ctx context.Context, r *request) {
	defer o.wg.Done()

	d := r.ExpiresAt.Sub(o.clock.Now())
	tr := o.clock.NewTimer(d)
	defer tr.Stop()

	select {
	case <-ctx.Done():
		return
	case <-tr.Chan():
		resp := &outputs{
			WorkflowExecutionID: r.WorkflowExecutionID,
			CapabilityResponse: capabilities.CapabilityResponse{
				Err:   fmt.Errorf("timeout exceeded: could not process request before expiry %s", r.WorkflowExecutionID),
				Value: nil,
			},
		}
		o.lggr.Debugw("expiryTimer - expired request", "workflowExecutionID", r.WorkflowExecutionID)
		o.transmitCh <- resp
	}
}
