package ocr3

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/mitchellh/mapstructure"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/datafeeds"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/mercury"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const (
	ocrCapabilityID = "offchain_reporting"
)

var info = capabilities.MustNewCapabilityInfo(
	ocrCapabilityID,
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
	lggr   logger.Logger

	requestTimeout time.Duration
	clock          clockwork.Clock

	aggregators map[string]types.Aggregator

	encoderFactory EncoderFactory
	encoders       map[string]types.Encoder

	transmitCh chan *response
	newTimerCh chan *request
}

var _ capabilityIface = (*capability)(nil)

func newCapability(s *store, clock clockwork.Clock, requestTimeout time.Duration, encoderFactory EncoderFactory, lggr logger.Logger) *capability {
	o := &capability{
		CapabilityInfo: info,
		store:          s,
		clock:          clock,
		requestTimeout: requestTimeout,
		stopCh:         make(chan struct{}),
		lggr:           logger.Named(lggr, "OCR3CapabilityClient"),
		encoderFactory: encoderFactory,
		aggregators:    map[string]types.Aggregator{},
		encoders:       map[string]types.Encoder{},

		transmitCh: make(chan *response),
		newTimerCh: make(chan *request),
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

type workflowConfig struct {
	AggregationMethod string      `mapstructure:"aggregation_method"`
	AggregationConfig *values.Map `mapstructure:"aggregation_config"`
	Encoder           string      `mapstructure:"encoder"`
	EncoderConfig     *values.Map `mapstructure:"encoder_config"`
}

func newWorkflowConfig() *workflowConfig {
	return &workflowConfig{
		EncoderConfig:     values.EmptyMap(),
		AggregationConfig: values.EmptyMap(),
	}
}

func (o *capability) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	c := newWorkflowConfig()
	err := request.Config.UnwrapTo(c)
	if err != nil {
		return err
	}

	switch c.AggregationMethod {
	case "data_feeds_2_0":
		mc := mercury.NewCodec()
		agg, err := datafeeds.NewDataFeedsAggregator(*c.AggregationConfig, mc, o.lggr)
		if err != nil {
			return err
		}

		o.aggregators[request.Metadata.WorkflowID] = agg

		encoder, err := o.encoderFactory(c.EncoderConfig)
		if err != nil {
			return err
		}
		o.encoders[request.Metadata.WorkflowID] = encoder
	default:
		return fmt.Errorf("aggregator %s not supported", c.AggregationMethod)
	}

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

	o.newTimerCh <- r
	return nil
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

func (o *capability) handleTransmitMsg(ctx context.Context, resp *response) {
	req, wasPresent := o.store.evict(ctx, resp.WorkflowExecutionID)
	if !wasPresent {
		return
	}

	select {
	case <-req.RequestCtx.Done():
		// This should only happen if the client has closed the upstream context.
		// In this case, the request is cancelled and we shouldn't transmit.
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
		resp := &response{
			WorkflowExecutionID: r.WorkflowExecutionID,
			CapabilityResponse: capabilities.CapabilityResponse{
				Err:   fmt.Errorf("timeout exceeded: could not process request before expiry %s", r.WorkflowExecutionID),
				Value: nil,
			},
		}

		o.transmitCh <- resp
	}
}

func (o *capability) transmitResponse(ctx context.Context, resp *response) error {
	o.transmitCh <- resp
	return nil
}

func (o *capability) unmarshalRequest(ctx context.Context, r capabilities.CapabilityRequest, callback chan<- capabilities.CapabilityResponse) (*request, error) {
	valuesMap, err := r.Inputs.Unwrap()
	if err != nil {
		return nil, err
	}

	expiresAt := o.clock.Now().Add(o.requestTimeout)
	req := &request{
		RequestCtx:          context.Background(), // TODO: set correct context
		CallbackCh:          callback,
		WorkflowExecutionID: r.Metadata.WorkflowExecutionID,
		WorkflowID:          r.Metadata.WorkflowID,
		Observations:        r.Inputs.Underlying["observations"],
		ExpiresAt:           expiresAt,
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
