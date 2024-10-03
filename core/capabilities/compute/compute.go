package compute

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	wasmpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/pb"
)

const (
	CapabilityIDCompute = "custom_compute@1.0.0"

	binaryKey = "binary"
	configKey = "config"
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

type Compute struct {
	log      logger.Logger
	registry coretypes.CapabilitiesRegistry
	modules  *moduleCache
}

func (c *Compute) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (c *Compute) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}

func generateID(binary []byte) string {
	id := sha256.Sum256(binary)
	return fmt.Sprintf("%x", id)
}

func copyRequest(req capabilities.CapabilityRequest) capabilities.CapabilityRequest {
	return capabilities.CapabilityRequest{
		Metadata: req.Metadata,
		Inputs:   req.Inputs.CopyMap(),
		Config:   req.Config.CopyMap(),
	}
}

func (c *Compute) Execute(ctx context.Context, request capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	copied := copyRequest(request)

	binary, err := c.popBytesValue(copied.Config, binaryKey)
	if err != nil {
		return capabilities.CapabilityResponse{}, fmt.Errorf("invalid request: %w", err)
	}

	config, err := c.popBytesValue(copied.Config, configKey)
	if err != nil {
		return capabilities.CapabilityResponse{}, fmt.Errorf("invalid request: %w", err)
	}

	id := generateID(binary)

	m, ok := c.modules.get(id)
	if !ok {
		mod, err := c.initModule(id, binary, request.Metadata.WorkflowID, request.Metadata.ReferenceID)
		if err != nil {
			return capabilities.CapabilityResponse{}, err
		}

		m = mod
	}

	return c.executeWithModule(m.module, config, request)
}

func (c *Compute) initModule(id string, binary []byte, workflowID, referenceID string) (*module, error) {
	initStart := time.Now()
	mod, err := host.NewModule(&host.ModuleConfig{Logger: c.log}, binary)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate WASM module: %w", err)
	}

	mod.Start()

	initDuration := time.Since(initStart)
	computeWASMInit.WithLabelValues(workflowID, referenceID).Observe(float64(initDuration))

	m := &module{module: mod}
	c.modules.add(id, m)
	return m, nil
}

func (c *Compute) popBytesValue(m *values.Map, key string) ([]byte, error) {
	v, ok := m.Underlying[key]
	if !ok {
		return nil, fmt.Errorf("could not find %q in map", key)
	}

	vb, ok := v.(*values.Bytes)
	if !ok {
		return nil, fmt.Errorf("value is not bytes: %q", key)
	}

	delete(m.Underlying, key)
	return vb.Underlying, nil
}

func (c *Compute) executeWithModule(module *host.Module, config []byte, req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
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
	resp, err := module.Run(wasmReq)
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

	computeWASMExec.WithLabelValues(
		req.Metadata.WorkflowID,
		req.Metadata.ReferenceID,
	).Observe(float64(time.Since(executeStart)))

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
	return c.registry.Add(ctx, c)
}

func (c *Compute) Close() error {
	c.modules.close()
	return nil
}

func NewAction(log logger.Logger, registry coretypes.CapabilitiesRegistry) *Compute {
	compute := &Compute{
		log:      logger.Named(log, "CustomCompute"),
		registry: registry,
		modules:  newModuleCache(clockwork.NewRealClock(), 1*time.Minute, 10*time.Minute, 3),
	}
	return compute
}
