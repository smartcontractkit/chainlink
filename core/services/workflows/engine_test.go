package workflows

import (
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	coreCap "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	pocCapabilities "github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"
)

const hardcodedWorkflow = `
triggers:
  - type: "mercury-trigger"
    config:
      feedIds:
        - "0x1111111111111111111100000000000000000000000000000000000000000000"
        - "0x2222222222222222222200000000000000000000000000000000000000000000"
        - "0x3333333333333333333300000000000000000000000000000000000000000000"

consensus:
  - type: "offchain_reporting"
    ref: "evm_median"
    inputs:
      observations:
        - "$(trigger.outputs)"
    config:
      aggregation_method: "data_feeds_2_0"
      aggregation_config:
        "0x1111111111111111111100000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: 3600
        "0x2222222222222222222200000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: 3600
        "0x3333333333333333333300000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: 3600
      encoder: "EVM"
      encoder_config:
        abi: "mercury_reports bytes[]"

targets:
  - type: "write_polygon-testnet-mumbai"
    inputs:
      report: "$(evm_median.outputs.report)"
    config:
      address: "0x3F3554832c636721F1fD1822Ccca0354576741Ef"
      params: ["$(report)"]
      abi: "receive(report bytes)"
  - type: "write_ethereum-testnet-sepolia"
    inputs:
      report: "$(evm_median.outputs.report)"
    config:
      address: "0x54e220867af6683aE6DcBF535B4f952cB5116510"
      params: ["$(report)"]
      abi: "receive(report bytes)"
`

// newTestEngineFromSpec creates a new engine with some test defaults.
func newTestEngineFromSpec(t *testing.T, reg *coreCap.Registry, spec string) (eng *Engine, initFailed chan struct{}) {
	builder := yamlBuilder{Spec: spec}
	return newTestEngine(t, reg, builder)
}

func newTestEngine(t *testing.T, reg *coreCap.Registry, builder specBuilder) (eng *Engine, initFailed chan struct{}) {
	peerID := p2ptypes.PeerID{}
	initFailed = make(chan struct{})
	cfg := Config{
		Lggr:        logger.TestLogger(t),
		Registry:    reg,
		SpecBuilder: builder,
		DONInfo:     nil,
		PeerID:      func() *p2ptypes.PeerID { return &peerID },
		maxRetries:  1,
		retryMs:     100,
		afterInit: func(success bool) {
			if !success {
				close(initFailed)
			}
		},
	}
	eng, err := NewEngine(cfg)
	require.NoError(t, err)
	return eng, initFailed
}

// getExecutionId returns the execution id of the workflow that is
// currently being executed by the engine.
//
// If the engine fails to initialize, the test will fail rather
// than blocking indefinitely.
func getExecutionId(t *testing.T, eng *Engine, initFailed <-chan struct{}) string {
	var eid string
	select {
	case <-initFailed:
		t.FailNow()
	case eid = <-eng.xxxExecutionFinished:
	}
	return eid
}

type mockCapability struct {
	capabilities.CapabilityInfo
	capabilities.CallbackExecutable
	response  chan capabilities.CapabilityResponse
	transform func(capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error)
}

func newMockCapability(info capabilities.CapabilityInfo, transform func(capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error)) *mockCapability {
	return &mockCapability{
		transform:      transform,
		CapabilityInfo: info,
		response:       make(chan capabilities.CapabilityResponse, 10),
	}
}

func (m *mockCapability) Execute(ctx context.Context, req capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	cr, err := m.transform(req)
	if err != nil {
		return nil, err
	}

	ch := make(chan capabilities.CapabilityResponse, 10)

	m.response <- cr
	ch <- cr
	close(ch)
	return ch, nil
}

func (m *mockCapability) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (m *mockCapability) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}

type mockTriggerCapability struct {
	capabilities.CapabilityInfo
	triggerEvent capabilities.CapabilityResponse
	ch           chan capabilities.CapabilityResponse
}

var _ capabilities.TriggerCapability = (*mockTriggerCapability)(nil)

func (m *mockTriggerCapability) RegisterTrigger(ctx context.Context, req capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	m.ch <- m.triggerEvent
	return m.ch, nil
}

func (m *mockTriggerCapability) UnregisterTrigger(ctx context.Context, req capabilities.CapabilityRequest) error {
	return nil
}

func TestEngineWithHardcodedWorkflow(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))

	trigger, cr := mockTrigger(t)

	require.NoError(t, reg.Add(ctx, trigger))
	require.NoError(t, reg.Add(ctx, mockConsensus()))
	target1 := mockTarget()
	require.NoError(t, reg.Add(ctx, target1))

	target2 := newMockCapability(
		capabilities.MustNewCapabilityInfo(
			"write_ethereum-testnet-sepolia",
			capabilities.CapabilityTypeTarget,
			"a write capability targeting ethereum sepolia testnet",
			"v1.0.0",
			nil,
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			m := req.Inputs.Underlying["report"].(*values.Map)
			return capabilities.CapabilityResponse{
				Value: m,
			}, nil
		},
	)
	require.NoError(t, reg.Add(ctx, target2))

	eng, initFailed := newTestEngineFromSpec(t, reg, hardcodedWorkflow)

	err := eng.Start(ctx)
	require.NoError(t, err)
	defer eng.Close()

	eid := getExecutionId(t, eng, initFailed)
	assert.Equal(t, cr, <-target1.response)
	assert.Equal(t, cr, <-target2.response)

	state, err := eng.executionStates.get(ctx, eid)
	require.NoError(t, err)

	assert.Equal(t, state.status, statusCompleted)
}

const (
	simpleWorkflow = `
triggers:
  - type: "mercury-trigger"
    config:
      feedlist:
        - "0x1111111111111111111100000000000000000000000000000000000000000000" # ETHUSD
        - "0x2222222222222222222200000000000000000000000000000000000000000000" # LINKUSD
        - "0x3333333333333333333300000000000000000000000000000000000000000000" # BTCUSD
        
consensus:
  - type: "offchain_reporting"
    ref: "evm_median"
    inputs:
      observations:
        - "$(trigger.outputs)"
    config:
      aggregation_method: "data_feeds_2_0"
      aggregation_config:
        "0x1111111111111111111100000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: "30m"
        "0x2222222222222222222200000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: "30m"
        "0x3333333333333333333300000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: "30m"
      encoder: "EVM"
      encoder_config:
        abi: "mercury_reports bytes[]"

targets:
  - type: "write_polygon-testnet-mumbai"
    inputs:
      report: "$(evm_median.outputs.report)"
    config:
      address: "0x3F3554832c636721F1fD1822Ccca0354576741Ef"
      params: ["$(report)"]
      abi: "receive(report bytes)"
`
)

func mockTrigger(t *testing.T) (capabilities.TriggerCapability, capabilities.CapabilityResponse) {
	mt := &mockTriggerCapability{
		CapabilityInfo: capabilities.MustNewCapabilityInfo(
			"mercury-trigger",
			capabilities.CapabilityTypeTrigger,
			"issues a trigger when a mercury report is received.",
			"v1.0.0",
			nil,
		),
		ch: make(chan capabilities.CapabilityResponse, 10),
	}
	resp, err := values.NewMap(map[string]any{
		"Values": map[string]any{
			"123": decimal.NewFromFloat(1.00),
			"456": decimal.NewFromFloat(1.25),
			"789": decimal.NewFromFloat(1.50),
		},
		"Decimals": map[string]any{
			"123": map[string]int{"Places": 19},
			"456": map[string]int{"Places": 19},
			"789": map[string]int{"Places": 8},
		},
		"Metadata": map[string]any{"TriggerRef": "Mercury"},
	})
	require.NoError(t, err)
	cr := capabilities.CapabilityResponse{
		Value: resp,
	}
	mt.triggerEvent = cr
	return mt, cr
}

func mockFailingConsensus() *mockCapability {
	return newMockCapability(
		capabilities.MustNewCapabilityInfo(
			"offchain_reporting",
			capabilities.CapabilityTypeConsensus,
			"an ocr3 consensus capability",
			"v3.0.0",
			nil,
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			return capabilities.CapabilityResponse{}, errors.New("fatal consensus error")
		},
	)
}

func mockConsensus() *mockCapability {
	return newMockCapability(
		capabilities.MustNewCapabilityInfo(
			"offchain_reporting",
			capabilities.CapabilityTypeConsensus,
			"an ocr3 consensus capability",
			"v3.0.0",
			nil,
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			obs := req.Inputs.Underlying["observations"]
			report := obs.(*values.List)
			rm := map[string]any{
				"report": report.Underlying[0],
			}
			rv, err := values.NewMap(rm)
			if err != nil {
				return capabilities.CapabilityResponse{}, err
			}

			return capabilities.CapabilityResponse{
				Value: rv,
			}, nil
		},
	)
}

func mockTarget() *mockCapability {
	return newMockCapability(
		capabilities.MustNewCapabilityInfo(
			"write_polygon-testnet-mumbai",
			capabilities.CapabilityTypeTarget,
			"a write capability targeting polygon mumbai testnet",
			"v1.0.0",
			nil,
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			m := req.Inputs.Underlying["report"].(*values.Map)
			return capabilities.CapabilityResponse{
				Value: m,
			}, nil
		},
	)
}

func TestEngine_ErrorsTheWorkflowIfAStepErrors(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))

	trigger, _ := mockTrigger(t)

	require.NoError(t, reg.Add(ctx, trigger))
	require.NoError(t, reg.Add(ctx, mockFailingConsensus()))
	require.NoError(t, reg.Add(ctx, mockTarget()))

	eng, initFailed := newTestEngineFromSpec(t, reg, simpleWorkflow)

	err := eng.Start(ctx)
	require.NoError(t, err)
	defer eng.Close()

	eid := getExecutionId(t, eng, initFailed)
	state, err := eng.executionStates.get(ctx, eid)
	require.NoError(t, err)

	assert.Equal(t, state.status, statusErrored)
	// evm_median is the ref of our failing consensus step
	assert.Equal(t, state.steps["evm_median"].status, statusErrored)
}

const (
	multiStepWorkflow = `
triggers:
  - type: "mercury-trigger"
    config:
      feedlist:
        - "0x1111111111111111111100000000000000000000000000000000000000000000" # ETHUSD
        - "0x2222222222222222222200000000000000000000000000000000000000000000" # LINKUSD
        - "0x3333333333333333333300000000000000000000000000000000000000000000" # BTCUSD

actions:
  - type: "read_chain_action"
    ref: "read_chain_action"
    inputs:
      action:
        - "$(trigger.outputs)"
        
consensus:
  - type: "offchain_reporting"
    ref: "evm_median"
    inputs:
      observations:
        - "$(trigger.outputs)"
        - "$(read_chain_action.outputs)"
    config:
      aggregation_method: "data_feeds_2_0"
      aggregation_config:
        "0x1111111111111111111100000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: "30m"
        "0x2222222222222222222200000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: "30m"
        "0x3333333333333333333300000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: "30m"
      encoder: "EVM"
      encoder_config:
        abi: "mercury_reports bytes[]"

targets:
  - type: "write_polygon-testnet-mumbai"
    inputs:
      report: "$(evm_median.outputs.report)"
    config:
      address: "0x3F3554832c636721F1fD1822Ccca0354576741Ef"
      params: ["$(report)"]
      abi: "receive(report bytes)"
`
)

func mockAction() (*mockCapability, values.Value) {
	outputs, _ := values.NewMap(map[string]any{"Read": "output"})
	return newMockCapability(
		capabilities.MustNewCapabilityInfo(
			"read_chain_action",
			capabilities.CapabilityTypeAction,
			"a read chain action",
			"v1.0.0",
			nil,
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {

			return capabilities.CapabilityResponse{
				Value: outputs,
			}, nil
		},
	), outputs
}

func TestEngine_MultiStepDependencies(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))

	trigger, cr := mockTrigger(t)

	require.NoError(t, reg.Add(ctx, trigger))
	require.NoError(t, reg.Add(ctx, mockConsensus()))
	require.NoError(t, reg.Add(ctx, mockTarget()))

	action, out := mockAction()
	require.NoError(t, reg.Add(ctx, action))

	eng, initFailed := newTestEngineFromSpec(t, reg, multiStepWorkflow)
	err := eng.Start(ctx)
	require.NoError(t, err)
	defer eng.Close()

	eid := getExecutionId(t, eng, initFailed)
	state, err := eng.executionStates.get(ctx, eid)
	require.NoError(t, err)

	assert.Equal(t, state.status, statusCompleted)

	// The inputs to the consensus step should
	// be the outputs of the two dependents.
	inputs := state.steps["evm_median"].inputs
	unw, err := values.Unwrap(inputs)
	require.NoError(t, err)

	obs := unw.(map[string]any)["observations"]
	assert.Len(t, obs, 2)

	tunw, err := values.Unwrap(cr.Value)
	require.NoError(t, err)
	assert.Equal(t, obs.([]any)[0], tunw)

	o, err := values.Unwrap(out)
	require.NoError(t, err)
	assert.Equal(t, obs.([]any)[1], o)
}

// Note that eventually we want to allow encoding in the consensus to be customizable as well.
// This is out of scope for this POC.
const multiStepWorkflowCodeConfig = `
type_map:
  - write_chain: write_polygon-testnet-mumbai
config:
  - mercury-trigger:
      feedlist:
        - '0x1111111111111111111100000000000000000000000000000000000000000000'
        - '0x2222222222222222222200000000000000000000000000000000000000000000'
  - evm_median: null
    aggregation_method: data_feeds_2_0
    aggregation_config:
      '0x1111111111111111111100000000000000000000000000000000000000000000':
        deviation: '0.001'
        heartbeat: 30m
      '0x2222222222222222222200000000000000000000000000000000000000000000':
        deviation: '0.001'
        heartbeat: 30m
      '0x3333333333333333333300000000000000000000000000000000000000000000':
        deviation: '0.001'
        heartbeat: 30m
    encoder: EVM
    encoder_config:
      abi: 'mercury_reports bytes[]'
  - write_chain:
      address: '0x3F3554832c636721F1fD1822Ccca0354576741Ef'
      params:
        - $(report)
      abi: receive(report bytes)
`

func TestEngine_MultiStepDependenciesCode(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))

	// This should be done automatically when they run for with code
	require.NoError(t, reg.Add(context.Background(), &localCodeCapability{
		Workflow:       nil,
		CapabilityType: capabilities.CapabilityTypeAction,
		Id:             pocCapabilities.LocalCodeActionCapability,
	}))
	require.NoError(t, reg.Add(context.Background(), &localCodeCapability{
		Workflow:       nil,
		CapabilityType: capabilities.CapabilityTypeConsensus,
		Id:             pocCapabilities.LocalCodeConsensusCapability,
	}))

	trigger, cr := mockTrigger(t)

	require.NoError(t, reg.Add(ctx, trigger))
	require.NoError(t, reg.Add(ctx, mockTarget()))

	require.NoError(t, reg.Add(ctx, mockConsensus()))
	action, out := mockAction()
	require.NoError(t, reg.Add(ctx, action))

	eng, initFailed := newTestEngine(t, reg, multiStepWorkflow)
	err := eng.Start(ctx)
	require.NoError(t, err)
	defer eng.Close()

	eid := getExecutionId(t, eng, initFailed)
	state, err := eng.executionStates.get(ctx, eid)
	require.NoError(t, err)

	assert.Equal(t, state.status, statusCompleted)

	// The inputs to the consensus step should
	// be the outputs of the two dependents.
	inputs := state.steps["evm_median"].inputs
	unw, err := values.Unwrap(inputs)
	require.NoError(t, err)

	obs := unw.(map[string]any)["observations"]
	assert.Len(t, obs, 2)

	tunw, err := values.Unwrap(cr.Value)
	require.NoError(t, err)
	assert.Equal(t, obs.([]any)[0], tunw)

	o, err := values.Unwrap(out)
	require.NoError(t, err)
	assert.Equal(t, obs.([]any)[1], o)
}
