package workflows

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	coreCap "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
)

const hardcodedWorkflow = `
triggers:
  - id: "mercury-trigger@1.0.0"
    config:
      feedIds:
        - "0x1111111111111111111100000000000000000000000000000000000000000000"
        - "0x2222222222222222222200000000000000000000000000000000000000000000"
        - "0x3333333333333333333300000000000000000000000000000000000000000000"

consensus:
  - id: "offchain_reporting@1.0.0"
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
  - id: "write_polygon-testnet-mumbai@1.0.0"
    inputs:
      report: "$(evm_median.outputs.report)"
    config:
      address: "0x3F3554832c636721F1fD1822Ccca0354576741Ef"
      params: ["$(report)"]
      abi: "receive(report bytes)"
  - id: "write_ethereum-testnet-sepolia@1.0.0"
    inputs:
      report: "$(evm_median.outputs.report)"
    config:
      address: "0x54e220867af6683aE6DcBF535B4f952cB5116510"
      params: ["$(report)"]
      abi: "receive(report bytes)"
`

type testHooks struct {
	initFailed        chan struct{}
	executionFinished chan string
}

// newTestEngine creates a new engine with some test defaults.
func newTestEngine(t *testing.T, reg *coreCap.Registry, spec string, opts ...func(c *Config)) (*Engine, *testHooks) {
	peerID := p2ptypes.PeerID{}
	initFailed := make(chan struct{})
	executionFinished := make(chan string, 100)
	cfg := Config{
		Lggr:     logger.TestLogger(t),
		Registry: reg,
		Spec:     spec,
		DONInfo: &capabilities.DON{
			ID: "00010203",
		},
		PeerID:     func() *p2ptypes.PeerID { return &peerID },
		maxRetries: 1,
		retryMs:    100,
		afterInit: func(success bool) {
			if !success {
				close(initFailed)
			}
		},
		onExecutionFinished: func(weid string) {
			executionFinished <- weid
		},
		clock: clockwork.NewFakeClock(),
	}
	for _, o := range opts {
		o(&cfg)
	}
	eng, err := NewEngine(cfg)
	require.NoError(t, err)
	return eng, &testHooks{initFailed: initFailed, executionFinished: executionFinished}
}

// getExecutionId returns the execution id of the workflow that is
// currently being executed by the engine.
//
// If the engine fails to initialize, the test will fail rather
// than blocking indefinitely.
func getExecutionId(t *testing.T, eng *Engine, hooks *testHooks) string {
	var eid string
	select {
	case <-hooks.initFailed:
		t.FailNow()
	case eid = <-hooks.executionFinished:
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
	triggerEvent *capabilities.CapabilityResponse
	ch           chan capabilities.CapabilityResponse
}

var _ capabilities.TriggerCapability = (*mockTriggerCapability)(nil)

func (m *mockTriggerCapability) RegisterTrigger(ctx context.Context, req capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	if m.triggerEvent != nil {
		m.ch <- *m.triggerEvent
	}
	return m.ch, nil
}

func (m *mockTriggerCapability) UnregisterTrigger(ctx context.Context, req capabilities.CapabilityRequest) error {
	return nil
}

func TestEngineWithHardcodedWorkflow(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		store store.Store
	}{
		{
			name:  "db-engine",
			store: store.NewDBStore(pgtest.NewSqlxDB(t), clockwork.NewFakeClock()),
		},
		{
			name:  "in-memory-engine",
			store: store.NewInMemoryStore(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testutils.Context(t)
			reg := coreCap.NewRegistry(logger.TestLogger(t))

			trigger, cr := mockTrigger(t)

			require.NoError(t, reg.Add(ctx, trigger))
			require.NoError(t, reg.Add(ctx, mockConsensus()))
			target1 := mockTarget()
			require.NoError(t, reg.Add(ctx, target1))

			target2 := newMockCapability(
				capabilities.MustNewCapabilityInfo(
					"write_ethereum-testnet-sepolia@1.0.0",
					capabilities.CapabilityTypeTarget,
					"a write capability targeting ethereum sepolia testnet",
				),
				func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
					m := req.Inputs.Underlying["report"].(*values.Map)
					return capabilities.CapabilityResponse{
						Value: m,
					}, nil
				},
			)
			require.NoError(t, reg.Add(ctx, target2))

			eng, testHooks := newTestEngine(
				t,
				reg,
				hardcodedWorkflow,
				func(c *Config) { c.Store = tc.store },
			)

			err := eng.Start(ctx)
			require.NoError(t, err)
			defer eng.Close()

			eid := getExecutionId(t, eng, testHooks)
			assert.Equal(t, cr, <-target1.response)
			assert.Equal(t, cr, <-target2.response)

			state, err := eng.executionStates.Get(ctx, eid)
			require.NoError(t, err)

			assert.Equal(t, state.Status, store.StatusCompleted)
		})
	}
}

const (
	simpleWorkflow = `
triggers:
  - id: "mercury-trigger@1.0.0"
    config:
      feedlist:
        - "0x1111111111111111111100000000000000000000000000000000000000000000" # ETHUSD
        - "0x2222222222222222222200000000000000000000000000000000000000000000" # LINKUSD
        - "0x3333333333333333333300000000000000000000000000000000000000000000" # BTCUSD
        
consensus:
  - id: "offchain_reporting@1.0.0"
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
  - id: "write_polygon-testnet-mumbai@1.0.0"
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
			"mercury-trigger@1.0.0",
			capabilities.CapabilityTypeTrigger,
			"issues a trigger when a mercury report is received.",
		),
		ch: make(chan capabilities.CapabilityResponse, 10),
	}
	resp, err := values.NewMap(map[string]any{
		"123": decimal.NewFromFloat(1.00),
		"456": decimal.NewFromFloat(1.25),
		"789": decimal.NewFromFloat(1.50),
	})
	require.NoError(t, err)
	cr := capabilities.CapabilityResponse{
		Value: resp,
	}
	mt.triggerEvent = &cr
	return mt, cr
}

func mockNoopTrigger(t *testing.T) capabilities.TriggerCapability {
	mt := &mockTriggerCapability{
		CapabilityInfo: capabilities.MustNewCapabilityInfo(
			"mercury-trigger@1.0.0",
			capabilities.CapabilityTypeTrigger,
			"issues a trigger when a mercury report is received.",
		),
		ch: make(chan capabilities.CapabilityResponse, 10),
	}
	return mt
}

func mockFailingConsensus() *mockCapability {
	return newMockCapability(
		capabilities.MustNewCapabilityInfo(
			"offchain_reporting@1.0.0",
			capabilities.CapabilityTypeConsensus,
			"an ocr3 consensus capability",
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			return capabilities.CapabilityResponse{}, errors.New("fatal consensus error")
		},
	)
}

func mockConsensusWithEarlyTermination() *mockCapability {
	return newMockCapability(
		capabilities.MustNewCapabilityInfo(
			"offchain_reporting@1.0.0",
			capabilities.CapabilityTypeConsensus,
			"an ocr3 consensus capability",
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			return capabilities.CapabilityResponse{
				Err: capabilities.ErrStopExecution,
			}, nil
		},
	)
}

func mockConsensus() *mockCapability {
	return newMockCapability(
		capabilities.MustNewCapabilityInfo(
			"offchain_reporting@1.0.0",
			capabilities.CapabilityTypeConsensus,
			"an ocr3 consensus capability",
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
			"write_polygon-testnet-mumbai@1.0.0",
			capabilities.CapabilityTypeTarget,
			"a write capability targeting polygon mumbai testnet",
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

	eng, hooks := newTestEngine(t, reg, simpleWorkflow)

	err := eng.Start(ctx)
	require.NoError(t, err)
	defer eng.Close()

	eid := getExecutionId(t, eng, hooks)
	state, err := eng.executionStates.Get(ctx, eid)
	require.NoError(t, err)

	assert.Equal(t, state.Status, store.StatusErrored)
	// evm_median is the ref of our failing consensus step
	assert.Equal(t, state.Steps["evm_median"].Status, store.StatusErrored)
}

func TestEngine_GracefulEarlyTermination(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))

	trigger, _ := mockTrigger(t)

	require.NoError(t, reg.Add(ctx, trigger))
	require.NoError(t, reg.Add(ctx, mockConsensusWithEarlyTermination()))
	require.NoError(t, reg.Add(ctx, mockTarget()))

	eng, hooks := newTestEngine(t, reg, simpleWorkflow)

	err := eng.Start(ctx)
	require.NoError(t, err)
	defer eng.Close()

	eid := getExecutionId(t, eng, hooks)
	state, err := eng.executionStates.Get(ctx, eid)
	require.NoError(t, err)

	assert.Equal(t, state.Status, store.StatusCompletedEarlyExit)
	assert.Nil(t, state.Steps["write_polygon-testnet-mumbai"])
}

const (
	multiStepWorkflow = `
triggers:
  - id: "mercury-trigger@1.0.0"
    config:
      feedlist:
        - "0x1111111111111111111100000000000000000000000000000000000000000000" # ETHUSD
        - "0x2222222222222222222200000000000000000000000000000000000000000000" # LINKUSD
        - "0x3333333333333333333300000000000000000000000000000000000000000000" # BTCUSD

actions:
  - id: "read_chain_action@1.0.0"
    ref: "read_chain_action"
    config: {}
    inputs:
      action:
        - "$(trigger.outputs)"
        
consensus:
  - id: "offchain_reporting@1.0.0"
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
  - id: "write_polygon-testnet-mumbai@1.0.0"
    inputs:
      report: "$(evm_median.outputs.report)"
    config:
      address: "0x3F3554832c636721F1fD1822Ccca0354576741Ef"
      params: ["$(report)"]
      abi: "receive(report bytes)"
`
)

func mockAction() (*mockCapability, values.Value) {
	outputs := values.NewString("output")
	return newMockCapability(
		capabilities.MustNewCapabilityInfo(
			"read_chain_action@1.0.0",
			capabilities.CapabilityTypeAction,
			"a read chain action",
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

	eng, hooks := newTestEngine(t, reg, multiStepWorkflow)
	err := eng.Start(ctx)
	require.NoError(t, err)
	defer eng.Close()

	eid := getExecutionId(t, eng, hooks)
	state, err := eng.executionStates.Get(ctx, eid)
	require.NoError(t, err)

	assert.Equal(t, state.Status, store.StatusCompleted)

	// The inputs to the consensus step should
	// be the outputs of the two dependents.
	inputs := state.Steps["evm_median"].Inputs
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

func TestEngine_ResumesPendingExecutions(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))

	trigger := mockNoopTrigger(t)
	resp, err := values.NewMap(map[string]any{
		"123": decimal.NewFromFloat(1.00),
		"456": decimal.NewFromFloat(1.25),
		"789": decimal.NewFromFloat(1.50),
	})
	require.NoError(t, err)

	require.NoError(t, reg.Add(ctx, trigger))
	require.NoError(t, reg.Add(ctx, mockConsensus()))
	require.NoError(t, reg.Add(ctx, mockTarget()))

	action, _ := mockAction()
	require.NoError(t, reg.Add(ctx, action))

	dbstore := store.NewDBStore(pgtest.NewSqlxDB(t), clockwork.NewFakeClock())
	ec := &store.WorkflowExecution{
		Steps: map[string]*store.WorkflowExecutionStep{
			workflows.KeywordTrigger: {
				Outputs: store.StepOutput{
					Value: resp,
				},
				Status:      store.StatusCompleted,
				ExecutionID: "<execution-ID>",
				Ref:         workflows.KeywordTrigger,
			},
		},
		WorkflowID:  "",
		ExecutionID: "<execution-ID>",
		Status:      store.StatusStarted,
	}
	err = dbstore.Add(ctx, ec)
	require.NoError(t, err)

	eng, hooks := newTestEngine(
		t,
		reg,
		multiStepWorkflow,
		func(c *Config) { c.Store = dbstore },
	)
	err = eng.Start(ctx)
	require.NoError(t, err)

	eid := getExecutionId(t, eng, hooks)
	gotEx, err := dbstore.Get(ctx, eid)
	require.NoError(t, err)
	assert.Equal(t, store.StatusCompleted, gotEx.Status)
}

func TestEngine_TimesOutOldExecutions(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))

	trigger := mockNoopTrigger(t)
	resp, err := values.NewMap(map[string]any{
		"123": decimal.NewFromFloat(1.00),
		"456": decimal.NewFromFloat(1.25),
		"789": decimal.NewFromFloat(1.50),
	})
	require.NoError(t, err)

	require.NoError(t, reg.Add(ctx, trigger))
	require.NoError(t, reg.Add(ctx, mockConsensus()))
	require.NoError(t, reg.Add(ctx, mockTarget()))

	action, _ := mockAction()
	require.NoError(t, reg.Add(ctx, action))

	clock := clockwork.NewFakeClock()
	dbstore := store.NewDBStore(pgtest.NewSqlxDB(t), clock)
	ec := &store.WorkflowExecution{
		Steps: map[string]*store.WorkflowExecutionStep{
			workflows.KeywordTrigger: {
				Outputs: store.StepOutput{
					Value: resp,
				},
				Status:      store.StatusCompleted,
				ExecutionID: "<execution-ID>",
				Ref:         workflows.KeywordTrigger,
			},
		},
		WorkflowID:  "",
		ExecutionID: "<execution-ID>",
		Status:      store.StatusStarted,
	}
	err = dbstore.Add(ctx, ec)
	require.NoError(t, err)

	eng, hooks := newTestEngine(
		t,
		reg,
		multiStepWorkflow,
		func(c *Config) {
			c.Store = dbstore
			c.clock = clock
		},
	)
	clock.Advance(15 * time.Minute)
	err = eng.Start(ctx)
	require.NoError(t, err)

	_ = getExecutionId(t, eng, hooks)
	gotEx, err := dbstore.Get(ctx, "<execution-ID>")
	require.NoError(t, err)
	assert.Equal(t, store.StatusTimeout, gotEx.Status)
}
