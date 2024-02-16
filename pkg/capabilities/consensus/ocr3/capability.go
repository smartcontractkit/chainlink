package ocr3

import (
	"context"
	"fmt"
	"sync"

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
	lggr   logger.Logger

	clock clockwork.Clock

	newExpiryWorkerCh chan *request

	aggregators map[string]types.Aggregator

	encoderFactory EncoderFactory
	encoders       map[string]types.Encoder
}

var _ capabilityIface = (*capability)(nil)

func newCapability(s *store, clock clockwork.Clock, encoderFactory EncoderFactory, lggr logger.Logger) *capability {
	o := &capability{
		CapabilityInfo:    info,
		store:             s,
		newExpiryWorkerCh: make(chan *request),
		clock:             clock,
		stopCh:            make(chan struct{}),
		lggr:              lggr,
		encoderFactory:    encoderFactory,
		aggregators:       map[string]types.Aggregator{},
		encoders:          map[string]types.Encoder{},
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

type workflowConfig struct {
	AggregationMethod       string `mapstructure:"aggregation_method"`
	AggregationMethodConfig map[string]any
	Encoder                 string
	EncoderConfig           map[string]any
}

func (o *capability) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	confMap, err := request.Config.Unwrap()
	if err != nil {
		return err
	}

	// TODO: values lib should a wrapped version of decode
	// which can handle passthrough translations of maps to values.Map.
	// This will avoid the need to translate/untranslate
	c := &workflowConfig{}
	err = mapstructure.Decode(confMap, c)
	if err != nil {
		return err
	}

	switch c.AggregationMethod {
	case "data_feeds_2_0":
		cm, err := values.NewMap(c.AggregationMethodConfig)
		if err != nil {
			return err
		}

		mc := mercury.NewCodec()
		agg, err := datafeeds.NewDataFeedsAggregator(*cm, mc, o.lggr)
		if err != nil {
			return err
		}

		o.aggregators[request.Metadata.WorkflowID] = agg

		em, err := values.NewMap(c.EncoderConfig)
		if err != nil {
			return err
		}

		encoder, err := o.encoderFactory(em)
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
		wasPresent := o.store.evict(ctx, r.WorkflowExecutionID)
		if !wasPresent {
			// the item was already evicted,
			// we'll assume it was processed successfully
			// and return
			return
		}

		timeoutResp := capabilities.CapabilityResponse{
			Err: fmt.Errorf("timeout exceeded: could not process request before expiry %+v", r.WorkflowExecutionID),
		}

		select {
		case <-r.RequestCtx.Done():
		case r.CallbackCh <- timeoutResp:
			close(r.CallbackCh)
		}
	}
}

func (o *capability) transmitResponse(ctx context.Context, resp response) error {
	req, err := o.store.get(ctx, resp.WorkflowExecutionID)
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
		o.store.evict(ctx, resp.WorkflowExecutionID)
		return nil
	}
}

func (o *capability) unmarshalRequest(ctx context.Context, r capabilities.CapabilityRequest, callback chan<- capabilities.CapabilityResponse) (*request, error) {
	valuesMap, err := r.Inputs.Unwrap()
	if err != nil {
		return nil, err
	}

	req := &request{
		RequestCtx:          ctx,
		CallbackCh:          callback,
		WorkflowExecutionID: r.Metadata.WorkflowExecutionID,
		WorkflowID:          r.Metadata.WorkflowID,
		Observations:        r.Inputs.Underlying["observations"],
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
