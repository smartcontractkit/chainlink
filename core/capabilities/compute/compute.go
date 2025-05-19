package compute

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/metering"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	wasmpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/pb"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/validation"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
)

const (
	CapabilityIDCompute = "custom-compute@1.0.0"

	binaryKey       = "binary"
	configKey       = "config"
	maxMemoryMBsKey = "maxMemoryMBs"
	timeoutKey      = "timeout"
	tickIntervalKey = "tickInterval"
)

var (
	computeWASMInit = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "compute_wasm_module_init",
		Help: "how long it takes to initialize a WASM module",
		Buckets: []float64{
			float64(50 * time.Millisecond),
			float64(100 * time.Millisecond),
			float64(200 * time.Millisecond),
			float64(500 * time.Millisecond),
			float64(1 * time.Second),
			float64(2 * time.Second),
			float64(4 * time.Second),
			float64(8 * time.Second),
		},
	}, []string{"workflowID", "stepRef"})
	computeWASMExec = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "compute_wasm_module_exec",
		Help: "how long it takes to execute a request from a WASM module",
		Buckets: []float64{
			float64(50 * time.Millisecond),
			float64(100 * time.Millisecond),
			float64(200 * time.Millisecond),
			float64(500 * time.Millisecond),
			float64(1 * time.Second),
			float64(2 * time.Second),
			float64(4 * time.Second),
			float64(8 * time.Second),
		},
	}, []string{"workflowID", "stepRef"})
)

var _ capabilities.ActionCapability = (*Compute)(nil)

type FetcherFn func(ctx context.Context, req *host.FetchRequest) (*host.FetchResponse, error)

type FetcherFactory interface {
	NewFetcher(log logger.Logger, emitter custmsg.MessageEmitter) FetcherFn
}

type Compute struct {
	stopCh services.StopChan
	log    logger.Logger

	// emitter is used to emit messages from the WASM module to a configured collector.
	emitter  custmsg.MessageEmitter
	registry coretypes.CapabilitiesRegistry
	modules  *moduleCache

	// transformer is used to transform a values.Map into a ParsedConfig struct on each execution
	// of a request.
	transformer *transformer

	fetcherFactory FetcherFactory

	numWorkers           int
	maxResponseSizeBytes uint64
	queue                chan request
	wg                   sync.WaitGroup
}

func (c *Compute) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (c *Compute) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}

func generateID(binary []byte) string {
	id := sha256.Sum256(binary)
	return hex.EncodeToString(id[:])
}

func (c *Compute) Execute(ctx context.Context, request capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	ch, err := c.enqueueRequest(ctx, request)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	select {
	case <-c.stopCh:
		return capabilities.CapabilityResponse{}, errors.New("service shutting down, aborting request")
	case <-ctx.Done():
		return capabilities.CapabilityResponse{}, fmt.Errorf("request cancelled by upstream: %w", ctx.Err())
	case resp := <-ch:
		return resp.resp, resp.err
	}
}

type request struct {
	ch  chan response
	req capabilities.CapabilityRequest
	ctx func() context.Context
}

type response struct {
	resp capabilities.CapabilityResponse
	err  error
}

func (c *Compute) enqueueRequest(ctx context.Context, req capabilities.CapabilityRequest) (<-chan response, error) {
	ch := make(chan response)
	r := request{
		ch:  ch,
		req: req,
		ctx: func() context.Context { return ctx },
	}
	select {
	case <-c.stopCh:
		return nil, errors.New("service shutting down, aborting request")
	case <-ctx.Done():
		return nil, fmt.Errorf("could not enqueue request: %w", ctx.Err())
	case c.queue <- r:
		return ch, nil
	}
}

func (c *Compute) execute(ctx context.Context, respCh chan response, req capabilities.CapabilityRequest) {
	copiedReq, cfg, err := c.transformer.Transform(req)
	if err != nil {
		respCh <- response{err: fmt.Errorf("invalid request: could not transform config: %w", err)}
		return
	}

	id := generateID(cfg.Binary)

	m, ok := c.modules.get(id)
	if !ok {
		mod, innerErr := c.initModule(id, cfg.ModuleConfig, cfg.Binary, copiedReq.Metadata)
		if innerErr != nil {
			respCh <- response{err: innerErr}
			return
		}

		m = mod
	}

	resp, err := c.executeWithModule(ctx, m.module, cfg.Config, copiedReq)
	select {
	case <-c.stopCh:
	case <-ctx.Done():
	case respCh <- response{resp: resp, err: err}:
	}
}

func (c *Compute) initModule(id string, cfg *host.ModuleConfig, binary []byte, requestMetadata capabilities.RequestMetadata) (*module, error) {
	initStart := time.Now()

	cfg.Fetch = c.fetcherFactory.NewFetcher(c.log, c.emitter)

	cfg.MaxResponseSizeBytes = c.maxResponseSizeBytes
	mod, err := host.NewModule(cfg, binary)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate WASM module: %w", err)
	}

	mod.Start()

	initDuration := time.Since(initStart)
	computeWASMInit.WithLabelValues(requestMetadata.WorkflowID, requestMetadata.ReferenceID).Observe(float64(initDuration))

	m := &module{module: mod}
	err = c.modules.add(id, m)
	if err != nil {
		c.log.Warnf("failed to add module to cache: %s", err.Error())
	}
	return m, nil
}

func (c *Compute) executeWithModule(ctx context.Context, module host.ModuleV1, config []byte, req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	executeStart := time.Now()
	capReq := capabilitiespb.CapabilityRequestToProto(req)

	wasmReq := &wasmpb.Request{
		Id:     uuid.New().String(),
		Config: config,
		Message: &wasmpb.Request_ComputeRequest{
			ComputeRequest: &wasmpb.ComputeRequest{
				Request: capReq,
			},
		},
	}
	resp, err := module.Run(ctx, wasmReq)
	if err != nil {
		return capabilities.CapabilityResponse{}, fmt.Errorf("error running module: %w", err)
	}

	cresppb := resp.GetComputeResponse().GetResponse()
	if cresppb == nil {
		return capabilities.CapabilityResponse{}, errors.New("got nil compute response")
	}

	cresp, err := capabilitiespb.CapabilityResponseFromProto(cresppb)
	if err != nil {
		return capabilities.CapabilityResponse{}, fmt.Errorf("could not convert response proto into response: %w", err)
	}

	executionTime := time.Since(executeStart)
	computeWASMExec.WithLabelValues(
		req.Metadata.WorkflowID,
		req.Metadata.ReferenceID,
	).Observe(float64(executionTime))

	// Add execution time to response metadata
	cresp.Metadata.Metering = append(cresp.Metadata.Metering, capabilities.MeteringNodeDetail{
		SpendUnit:  metering.ComputeUnit.Name,
		SpendValue: strconv.FormatInt(int64(executionTime.Round(time.Second)/(1_000_000_000)), 10),
	})

	return cresp, nil
}

func (c *Compute) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.NewCapabilityInfo(
		CapabilityIDCompute,
		capabilities.CapabilityTypeAction,
		"WASM custom compute capability",
	)
}

func (c *Compute) Start(ctx context.Context) error {
	c.modules.start()

	c.wg.Add(c.numWorkers)
	for i := 0; i < c.numWorkers; i++ {
		go func() {
			innerCtx, cancel := c.stopCh.NewCtx()
			defer cancel()

			defer c.wg.Done()
			c.worker(innerCtx)
		}()
	}
	return c.registry.Add(ctx, c)
}

func (c *Compute) worker(ctx context.Context) {
	for {
		select {
		case <-c.stopCh:
			return
		case req := <-c.queue:
			c.execute(req.ctx(), req.ch, req.req)
		}
	}
}

func (c *Compute) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c.modules.close()
	close(c.stopCh)

	err := c.registry.Remove(ctx, CapabilityIDCompute)
	if err != nil {
		return err
	}

	c.wg.Wait()

	return nil
}

type outgoingConnectorFetcherFactory struct {
	outgoingConnectorHandler *webapi.OutgoingConnectorHandler
	idGenerator              func() string
	metrics                  *computeMetricsLabeler
}

func NewOutgoingConnectorFetcherFactory(
	outgoingConnectorHandler *webapi.OutgoingConnectorHandler,
	idGenerator func() string,
) (FetcherFactory, error) {
	metricsLabeler, err := newComputeMetricsLabeler(metrics.NewLabeler().With("capability", CapabilityIDCompute))
	if err != nil {
		return nil, fmt.Errorf("failed to create compute metrics labeler: %w", err)
	}

	factory := &outgoingConnectorFetcherFactory{
		outgoingConnectorHandler: outgoingConnectorHandler,
		idGenerator:              idGenerator,
		metrics:                  metricsLabeler,
	}

	return factory, nil
}

func (f *outgoingConnectorFetcherFactory) NewFetcher(log logger.Logger, emitter custmsg.MessageEmitter) FetcherFn {
	return func(ctx context.Context, req *host.FetchRequest) (*host.FetchResponse, error) {
		if err := validation.ValidateWorkflowOrExecutionID(req.Metadata.WorkflowID); err != nil {
			return nil, fmt.Errorf("workflow ID %q is invalid: %w", req.Metadata.WorkflowID, err)
		}
		if err := validation.ValidateWorkflowOrExecutionID(req.Metadata.WorkflowExecutionID); err != nil {
			return nil, fmt.Errorf("workflow execution ID %q is invalid: %w", req.Metadata.WorkflowExecutionID, err)
		}

		cma := emitter.With(
			platform.KeyWorkflowID, req.Metadata.WorkflowID,
			platform.KeyWorkflowName, req.Metadata.DecodedWorkflowName,
			platform.KeyWorkflowOwner, req.Metadata.WorkflowOwner,
			platform.KeyWorkflowExecutionID, req.Metadata.WorkflowExecutionID,
			timestampKey, time.Now().UTC().Format(time.RFC3339Nano),
		)

		messageID := strings.Join([]string{
			req.Metadata.WorkflowExecutionID,
			ghcapabilities.MethodComputeAction,
			f.idGenerator(),
		}, "/")

		resp, err := f.outgoingConnectorHandler.HandleSingleNodeRequest(ctx, messageID, ghcapabilities.Request{
			URL:        req.URL,
			Method:     req.Method,
			Headers:    req.Headers,
			Body:       req.Body,
			TimeoutMs:  req.TimeoutMs,
			WorkflowID: req.Metadata.WorkflowID,
		})
		if err != nil {
			return nil, err
		}

		log.Debugw("received gateway response", "donID", resp.Body.DonId, "msgID", resp.Body.MessageId, "receiver", resp.Body.Receiver, "sender", resp.Body.Sender)
		var response host.FetchResponse
		err = json.Unmarshal(resp.Body.Payload, &response)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal fetch response: %w", err)
		}

		f.metrics.with(
			"status", strconv.FormatUint(uint64(response.StatusCode), 10),
			platform.KeyWorkflowID, req.Metadata.WorkflowID,
			platform.KeyWorkflowName, req.Metadata.WorkflowName,
			platform.KeyWorkflowOwner, req.Metadata.WorkflowOwner,
		).incrementHTTPRequestCounter(ctx)

		// Only log if the response is not in the 200 range
		if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
			msg := fmt.Sprintf("compute fetch request failed with status code %d", response.StatusCode)
			err = cma.Emit(ctx, msg)
			if err != nil {
				log.Errorf("failed to send custom message with msg: %s, err: %v", msg, err)
			}
		}

		return &response, nil
	}
}

const (
	defaultNumWorkers                = 3
	defaultMaxMemoryMBs              = 128
	defaultMaxTickInterval           = 100 * time.Millisecond
	defaultMaxTimeout                = 10 * time.Second
	defaultMaxCompressedBinarySize   = 20 * 1024 * 1024  // 20 MB
	defaultMaxDecompressedBinarySize = 100 * 1024 * 1024 // 100 MB
	defaultMaxResponseSizeBytes      = 5 * 1024 * 1024   // 5 MB
)

type Config struct {
	webapi.ServiceConfig
	NumWorkers                int
	MaxMemoryMBs              uint64
	MaxTimeout                time.Duration
	MaxTickInterval           time.Duration
	MaxCompressedBinarySize   uint64
	MaxDecompressedBinarySize uint64
	MaxResponseSizeBytes      uint64
}

func (c *Config) ApplyDefaults() {
	if c.NumWorkers == 0 {
		c.NumWorkers = defaultNumWorkers
	}
	if c.MaxMemoryMBs == 0 {
		c.MaxMemoryMBs = defaultMaxMemoryMBs
	}
	if c.MaxTimeout == 0 {
		c.MaxTimeout = defaultMaxTimeout
	}
	if c.MaxTickInterval == 0 {
		c.MaxTickInterval = defaultMaxTickInterval
	}
	if c.MaxCompressedBinarySize == 0 {
		c.MaxCompressedBinarySize = uint64(defaultMaxCompressedBinarySize)
	}
	if c.MaxDecompressedBinarySize == 0 {
		c.MaxDecompressedBinarySize = uint64(defaultMaxDecompressedBinarySize)
	}
	if c.MaxResponseSizeBytes == 0 {
		c.MaxResponseSizeBytes = uint64(defaultMaxResponseSizeBytes)
	}
}

func NewAction(
	config Config,
	log logger.Logger,
	registry coretypes.CapabilitiesRegistry,
	fetcherFactory FetcherFactory,
	opts ...func(*Compute),
) (*Compute, error) {
	config.ApplyDefaults()

	var (
		lggr    = logger.Named(log, "CustomCompute")
		labeler = custmsg.NewLabeler()
		compute = &Compute{
			stopCh:               make(services.StopChan),
			log:                  lggr,
			emitter:              labeler,
			registry:             registry,
			modules:              newModuleCache(clockwork.NewRealClock(), 1*time.Minute, 10*time.Minute, 3),
			transformer:          NewTransformer(lggr, labeler, config),
			fetcherFactory:       fetcherFactory,
			queue:                make(chan request),
			numWorkers:           config.NumWorkers,
			maxResponseSizeBytes: config.MaxResponseSizeBytes,
		}
	)

	for _, opt := range opts {
		opt(compute)
	}

	return compute, nil
}
