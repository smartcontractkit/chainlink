package workflows

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	gcmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"

	coreCap "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/compute"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
)

const (
	testWorkflowID    = "<workflow-id>"
	testWorkflowOwner = "testowner"
	testWorkflowName  = "testworkflow"
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
      cre_step_timeout: 610
`

type testHooks struct {
	initFailed        chan struct{}
	initSuccessful    chan struct{}
	executionFinished chan string
	rateLimited       chan string
}

func newTestDBStore(t *testing.T, clock clockwork.Clock) store.Store {
	// Taken from https://github.com/smartcontractkit/chainlink/blob/d736d9e0838983a021677bc608556b3994f46690/core/services/job/orm.go#L412
	// We need to insert this row so that we dont get foreign key constraint errors
	// based on the workflow_id
	db := pgtest.NewSqlxDB(t)
	sql := `INSERT INTO workflow_specs (workflow, workflow_id, workflow_owner, workflow_name, created_at, updated_at)
	VALUES (:workflow, :workflow_id, :workflow_owner, :workflow_name, NOW(), NOW())
	RETURNING id;`
	var wfSpec job.WorkflowSpec
	wfSpec.Workflow = simpleWorkflow
	wfSpec.WorkflowID = testWorkflowID
	wfSpec.WorkflowOwner = testWorkflowOwner
	wfSpec.WorkflowName = testWorkflowName
	_, err := db.NamedExec(sql, wfSpec)
	require.NoError(t, err)

	return store.NewDBStore(db, logger.TestLogger(t), clock)
}

type testConfigProvider struct {
	localNode           func(ctx context.Context) (capabilities.Node, error)
	configForCapability func(ctx context.Context, capabilityID string, donID uint32) (registrysyncer.CapabilityConfiguration, error)
}

func (t testConfigProvider) LocalNode(ctx context.Context) (capabilities.Node, error) {
	if t.localNode != nil {
		return t.localNode(ctx)
	}

	peerID := p2ptypes.PeerID{}
	return capabilities.Node{
		WorkflowDON: capabilities.DON{
			ID: 1,
		},
		PeerID: &peerID,
	}, nil
}

func (t testConfigProvider) ConfigForCapability(ctx context.Context, capabilityID string, donID uint32) (registrysyncer.CapabilityConfiguration, error) {
	if t.configForCapability != nil {
		return t.configForCapability(ctx, capabilityID, donID)
	}

	return registrysyncer.CapabilityConfiguration{}, nil
}

func newTestEngineWithYAMLSpec(t *testing.T, reg *coreCap.Registry, spec string, opts ...func(c *Config)) (*Engine, *testHooks) {
	sdkSpec, err := (&job.WorkflowSpec{
		Workflow: spec,
		SpecType: job.YamlSpec,
	}).SDKSpec(testutils.Context(t))
	require.NoError(t, err)

	return newTestEngine(t, reg, sdkSpec, opts...)
}

type mockSecretsFetcher struct{}

func (s mockSecretsFetcher) SecretsFor(ctx context.Context, workflowOwner, hexWorkflowName, decodedWorkflowName, workflowID string) (map[string]string, error) {
	return map[string]string{}, nil
}

// newTestEngine creates a new engine with some test defaults.
func newTestEngine(t *testing.T, reg *coreCap.Registry, sdkSpec sdk.WorkflowSpec, opts ...func(c *Config)) (*Engine, *testHooks) {
	initFailed := make(chan struct{})
	initSuccessful := make(chan struct{})
	executionFinished := make(chan string, 100)
	rateLimited := make(chan string)
	clock := clockwork.NewFakeClock()
	rl, err := ratelimiter.NewRateLimiter(ratelimiter.Config{
		GlobalRPS:      1000.0,
		GlobalBurst:    1000,
		PerSenderRPS:   100.0,
		PerSenderBurst: 100,
	})
	require.NoError(t, err)

	reg.SetLocalRegistry(&testConfigProvider{})
	cfg := Config{
		WorkflowID:    testWorkflowID,
		WorkflowOwner: testWorkflowOwner,
		WorkflowName: defaultName{
			name: testWorkflowName,
		},
		Lggr:       logger.TestLogger(t),
		Registry:   reg,
		Workflow:   sdkSpec,
		maxRetries: 1,
		retryMs:    100,
		afterInit: func(success bool) {
			if success {
				close(initSuccessful)
			} else {
				close(initFailed)
			}
		},
		onExecutionFinished: func(weid string) {
			executionFinished <- weid
		},
		onRateLimit: func(weid string) {
			rateLimited <- weid
		},
		SecretsFetcher: mockSecretsFetcher{},
		clock:          clock,
		RateLimiter:    rl,
	}
	for _, o := range opts {
		o(&cfg)
	}
	// We use the cfg clock incase they override it
	if cfg.Store == nil {
		cfg.Store = newTestDBStore(t, cfg.clock)
	}
	eng, err := NewEngine(testutils.Context(t), cfg)
	require.NoError(t, err)
	return eng, &testHooks{initSuccessful: initSuccessful, initFailed: initFailed, executionFinished: executionFinished, rateLimited: rateLimited}
}

// getExecutionId returns the execution id of the workflow that is
// currently being executed by the engine.
//
// If the engine fails to initialize, the test will fail rather
// than blocking indefinitely.
func getExecutionID(t *testing.T, eng *Engine, hooks *testHooks) string {
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
	capabilities.Executable
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

func (m *mockCapability) Execute(ctx context.Context, req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	cr, err := m.transform(req)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	m.response <- cr
	return cr, nil
}

func (m *mockCapability) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (m *mockCapability) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}

type mockTriggerCapability struct {
	capabilities.CapabilityInfo
	triggerEvent               *capabilities.TriggerResponse
	ch                         chan capabilities.TriggerResponse
	registerTriggerCallCounter map[string]int
}

var _ capabilities.TriggerCapability = (*mockTriggerCapability)(nil)

func (m *mockTriggerCapability) RegisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	m.registerTriggerCallCounter[req.TriggerID]++
	if m.triggerEvent != nil {
		m.ch <- *m.triggerEvent
	}
	return m.ch, nil
}

func (m *mockTriggerCapability) UnregisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) error {
	if m.registerTriggerCallCounter[req.TriggerID] == 0 {
		return errors.New("failed to unregister a non-registered trigger")
	}
	return nil
}

func TestEngineWithHardcodedWorkflow(t *testing.T) {
	ctx := testutils.Context(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))

	trigger, cr := mockTrigger(t)

	require.NoError(t, reg.Add(ctx, trigger))
	require.NoError(t, reg.Add(ctx, mockConsensus("")))
	target1 := mockTarget("")
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

	eng, testHooks := newTestEngineWithYAMLSpec(
		t,
		reg,
		hardcodedWorkflow,
	)

	servicetest.Run(t, eng)

	eid := getExecutionID(t, eng, testHooks)
	resp1 := <-target1.response
	assert.Equal(t, cr.Event.Outputs, resp1.Value)

	resp2 := <-target2.response
	assert.Equal(t, cr.Event.Outputs, resp2.Value)

	state, err := eng.executionStates.Get(ctx, eid)
	require.NoError(t, err)

	assert.Equal(t, state.Status, store.StatusCompleted)
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

func mockTrigger(t *testing.T) (capabilities.TriggerCapability, capabilities.TriggerResponse) {
	mt := &mockTriggerCapability{
		CapabilityInfo: capabilities.MustNewCapabilityInfo(
			"mercury-trigger@1.0.0",
			capabilities.CapabilityTypeTrigger,
			"issues a trigger when a mercury report is received.",
		),
		ch:                         make(chan capabilities.TriggerResponse, 10),
		registerTriggerCallCounter: make(map[string]int),
	}
	resp, err := values.NewMap(map[string]any{
		"123": decimal.NewFromFloat(1.00),
		"456": decimal.NewFromFloat(1.25),
		"789": decimal.NewFromFloat(1.50),
	})
	require.NoError(t, err)
	tr := capabilities.TriggerResponse{
		Event: capabilities.TriggerEvent{
			TriggerType: mt.ID,
			ID:          time.Now().UTC().Format(time.RFC3339),
			Outputs:     resp,
		},
	}
	mt.triggerEvent = &tr
	return mt, tr
}

func mockNoopTrigger(t *testing.T) capabilities.TriggerCapability {
	mt := &mockTriggerCapability{
		CapabilityInfo: capabilities.MustNewCapabilityInfo(
			"mercury-trigger@1.0.0",
			capabilities.CapabilityTypeTrigger,
			"issues a trigger when a mercury report is received.",
		),
		ch:                         make(chan capabilities.TriggerResponse, 10),
		registerTriggerCallCounter: make(map[string]int),
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

func mockConsensusWithEarlyTermination(id string) *mockCapability {
	if len(id) == 0 {
		id = "offchain_reporting@1.0.0"
	}
	return newMockCapability(
		capabilities.MustNewCapabilityInfo(
			id,
			capabilities.CapabilityTypeConsensus,
			"an ocr3 consensus capability",
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			return capabilities.CapabilityResponse{},
				// copy error object to make sure message comparison works as expected
				errors.New(capabilities.ErrStopExecution.Error())
		},
	)
}

func mockConsensus(id string) *mockCapability {
	if len(id) == 0 {
		id = "offchain_reporting@1.0.0"
	}
	return newMockCapability(
		capabilities.MustNewCapabilityInfo(
			id,
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

func mockTarget(id string) *mockCapability {
	if len(id) == 0 {
		id = "write_polygon-testnet-mumbai@1.0.0"
	}
	return newMockCapability(
		capabilities.MustNewCapabilityInfo(
			id,
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

func TestEngine_RateLimit(t *testing.T) {
	t.Run("per user rate limit", func(t *testing.T) {
		ctx := testutils.Context(t)
		reg := coreCap.NewRegistry(logger.TestLogger(t))

		trigger, _ := mockTrigger(t)
		require.NoError(t, reg.Add(ctx, trigger))
		require.NoError(t, reg.Add(ctx, mockConsensus("")))
		target1 := mockTarget("")
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

		setRateLimiter := func(c *Config) {
			rl, err := ratelimiter.NewRateLimiter(ratelimiter.Config{
				GlobalRPS:      1000.0,
				GlobalBurst:    1000,
				PerSenderRPS:   1.0,
				PerSenderBurst: 1,
			})
			require.NoError(t, err)
			c.RateLimiter = rl
		}

		eng, testHooks := newTestEngineWithYAMLSpec(
			t,
			reg,
			hardcodedWorkflow,
			setRateLimiter,
		)

		// Call RateLimiter once as owner, so next execution gets blocked by per user limit
		senderAllow, globalAllow := eng.ratelimiter.Allow(testWorkflowOwner)
		require.True(t, senderAllow)
		require.True(t, globalAllow)
		servicetest.Run(t, eng)

		select {
		case <-testHooks.rateLimited:
		case <-ctx.Done():
			t.FailNow()
		}
	})

	t.Run("global rate limit", func(t *testing.T) {
		ctx := testutils.Context(t)
		reg := coreCap.NewRegistry(logger.TestLogger(t))

		trigger, _ := mockTrigger(t)
		require.NoError(t, reg.Add(ctx, trigger))
		require.NoError(t, reg.Add(ctx, mockConsensus("")))
		target1 := mockTarget("")
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

		setRateLimiter := func(c *Config) {
			rl, err := ratelimiter.NewRateLimiter(ratelimiter.Config{
				GlobalRPS:      1.0,
				GlobalBurst:    1,
				PerSenderRPS:   100.0,
				PerSenderBurst: 100,
			})
			require.NoError(t, err)
			c.RateLimiter = rl
		}

		eng, testHooks := newTestEngineWithYAMLSpec(
			t,
			reg,
			hardcodedWorkflow,
			setRateLimiter,
		)

		// Call RateLimiter once as other owner, so next execution gets blocked by global limit
		senderAllow, globalAllow := eng.ratelimiter.Allow("some other owner")
		require.True(t, senderAllow)
		require.True(t, globalAllow)
		servicetest.Run(t, eng)

		select {
		case <-testHooks.rateLimited:
		case <-ctx.Done():
			t.FailNow()
		}
	})
}

func TestEngine_ErrorsTheWorkflowIfAStepErrors(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))

	trigger, _ := mockTrigger(t)

	require.NoError(t, reg.Add(ctx, trigger))
	require.NoError(t, reg.Add(ctx, mockFailingConsensus()))
	require.NoError(t, reg.Add(ctx, mockTarget("write_polygon-testnet-mumbai@1.0.0")))

	eng, hooks := newTestEngineWithYAMLSpec(t, reg, simpleWorkflow)

	servicetest.Run(t, eng)

	eid := getExecutionID(t, eng, hooks)
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
	require.NoError(t, reg.Add(ctx, mockConsensusWithEarlyTermination("")))
	require.NoError(t, reg.Add(ctx, mockTarget("")))

	eng, hooks := newTestEngineWithYAMLSpec(t, reg, simpleWorkflow)
	servicetest.Run(t, eng)

	eid := getExecutionID(t, eng, hooks)
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

func mockAction(t *testing.T) (*mockCapability, values.Value) {
	outputs, err := values.NewMap(map[string]any{"output": "foo"})
	require.NoError(t, err)
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

	trigger, tr := mockTrigger(t)

	require.NoError(t, reg.Add(ctx, trigger))
	require.NoError(t, reg.Add(ctx, mockConsensus("")))
	require.NoError(t, reg.Add(ctx, mockTarget("")))

	action, out := mockAction(t)
	require.NoError(t, reg.Add(ctx, action))

	eng, hooks := newTestEngineWithYAMLSpec(t, reg, multiStepWorkflow)
	servicetest.Run(t, eng)

	eid := getExecutionID(t, eng, hooks)
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

	require.NoError(t, err)
	uo, err := values.Unwrap(tr.Event.Outputs)
	require.NoError(t, err)
	assert.Equal(t, obs.([]any)[0].(map[string]any), uo)

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
	require.NoError(t, reg.Add(ctx, mockConsensus("")))
	require.NoError(t, reg.Add(ctx, mockTarget("")))

	action, _ := mockAction(t)
	require.NoError(t, reg.Add(ctx, action))
	dbstore := newTestDBStore(t, clockwork.NewFakeClock())
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
		WorkflowID:  testWorkflowID,
		ExecutionID: "<execution-ID>",
		Status:      store.StatusStarted,
	}
	_, err = dbstore.Add(ctx, ec)
	require.NoError(t, err)

	eng, hooks := newTestEngineWithYAMLSpec(
		t,
		reg,
		multiStepWorkflow,
		func(c *Config) { c.Store = dbstore },
	)
	servicetest.Run(t, eng)

	eid := getExecutionID(t, eng, hooks)
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
	require.NoError(t, reg.Add(ctx, mockConsensus("")))
	require.NoError(t, reg.Add(ctx, mockTarget("")))

	action, _ := mockAction(t)
	require.NoError(t, reg.Add(ctx, action))

	clock := clockwork.NewFakeClock()
	dbstore := newTestDBStore(t, clock)
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
		WorkflowID:  testWorkflowID,
		ExecutionID: "<execution-ID>",
		Status:      store.StatusStarted,
	}
	_, err = dbstore.Add(ctx, ec)
	require.NoError(t, err)

	eng, hooks := newTestEngineWithYAMLSpec(
		t,
		reg,
		multiStepWorkflow,
		func(c *Config) {
			c.Store = dbstore
			c.clock = clock
		},
	)
	clock.Advance(15 * time.Minute)
	servicetest.Run(t, eng)

	_ = getExecutionID(t, eng, hooks)
	gotEx, err := dbstore.Get(ctx, "<execution-ID>")
	require.NoError(t, err)
	assert.Equal(t, store.StatusTimeout, gotEx.Status)
}

const (
	delayedWorkflow = `
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
      deltaStage: 2s
      schedule: allAtOnce
`
)

func TestEngine_WrapsTargets(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))

	trigger, _ := mockTrigger(t)

	require.NoError(t, reg.Add(ctx, trigger))
	require.NoError(t, reg.Add(ctx, mockConsensus("")))
	require.NoError(t, reg.Add(ctx, mockTarget("")))

	clock := clockwork.NewFakeClock()
	dbstore := newTestDBStore(t, clock)

	eng, hooks := newTestEngineWithYAMLSpec(
		t,
		reg,
		delayedWorkflow,
		func(c *Config) {
			c.Store = dbstore
			c.clock = clock
		},
	)
	servicetest.Run(t, eng)

	<-hooks.initSuccessful

	err := eng.workflow.walkDo(workflows.KeywordTrigger, func(s *step) error {
		if s.Ref == workflows.KeywordTrigger {
			return nil
		}

		info, err2 := s.capability.Info(ctx)
		require.NoError(t, err2)

		if info.CapabilityType == capabilities.CapabilityTypeTarget {
			assert.Equal(t, "*transmission.LocalTargetCapability", fmt.Sprintf("%T", s.capability))
		} else {
			assert.NotEqual(t, "*transmission.LocalTargetCapability", fmt.Sprintf("%T", s.capability))
		}

		return nil
	})
	require.NoError(t, err)
}

func TestEngine_GetsNodeInfoDuringInitialization(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))

	trigger, _ := mockTrigger(t)

	require.NoError(t, reg.Add(ctx, trigger))
	require.NoError(t, reg.Add(ctx, mockConsensus("")))
	require.NoError(t, reg.Add(ctx, mockTarget("")))

	clock := clockwork.NewFakeClock()
	dbstore := newTestDBStore(t, clock)

	var peerID p2ptypes.PeerID
	node := capabilities.Node{
		PeerID: &peerID,
		WorkflowDON: capabilities.DON{
			ID: 1,
		},
	}
	retryCount := 0

	reg.SetLocalRegistry(testConfigProvider{
		localNode: func(ctx context.Context) (capabilities.Node, error) {
			n := capabilities.Node{}
			err := errors.New("peer not initialized")
			if retryCount > 0 {
				n = node
				err = nil
			}
			retryCount++
			return n, err
		},
	})
	eng, hooks := newTestEngineWithYAMLSpec(
		t,
		reg,
		delayedWorkflow,
		func(c *Config) {
			c.Store = dbstore
			c.clock = clock
			c.maxRetries = 2
			c.retryMs = 0
		},
	)
	servicetest.Run(t, eng)

	<-hooks.initSuccessful

	assert.Equal(t, node, eng.localNode)
}

const passthroughInterpolationWorkflow = `
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
  - id: "write_ethereum-testnet-sepolia@1.0.0"
    inputs: "$(evm_median.outputs)"
    config:
      address: "0x54e220867af6683aE6DcBF535B4f952cB5116510"
      params: ["$(report)"]
      abi: "receive(report bytes)"
`

func TestEngine_PassthroughInterpolation(t *testing.T) {
	ctx := testutils.Context(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))

	trigger, _ := mockTrigger(t)

	require.NoError(t, reg.Add(ctx, trigger))
	require.NoError(t, reg.Add(ctx, mockConsensus("")))
	writeID := "write_ethereum-testnet-sepolia@1.0.0"
	target := newMockCapability(
		capabilities.MustNewCapabilityInfo(
			writeID,
			capabilities.CapabilityTypeTarget,
			"a write capability targeting ethereum sepolia testnet",
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			return capabilities.CapabilityResponse{
				Value: req.Inputs,
			}, nil
		},
	)
	require.NoError(t, reg.Add(ctx, target))

	eng, testHooks := newTestEngineWithYAMLSpec(
		t,
		reg,
		passthroughInterpolationWorkflow,
	)

	servicetest.Run(t, eng)

	eid := getExecutionID(t, eng, testHooks)

	state, err := eng.executionStates.Get(ctx, eid)
	require.NoError(t, err)

	assert.Equal(t, state.Status, store.StatusCompleted)

	// There is passthrough interpolation between the consensus and target steps,
	// so the input of one should be the output of the other, exactly.
	gotInputs, err := values.Unwrap(state.Steps[writeID].Inputs)
	require.NoError(t, err)

	gotOutputs, err := values.Unwrap(state.Steps["evm_median"].Outputs.Value)
	require.NoError(t, err)
	assert.Equal(t, gotInputs, gotOutputs)
}

func TestEngine_Error(t *testing.T) {
	err := errors.New("some error")
	tests := []struct {
		name   string
		labels map[string]string
		err    error
		reason string
		want   string
	}{
		{
			name:   "Error with error and reason",
			labels: map[string]string{platform.KeyWorkflowID: "my-workflow-id"},
			err:    err,
			reason: "some reason",
			want:   "workflowID my-workflow-id: some reason: some error",
		},
		{
			name:   "Error with error and no reason",
			labels: map[string]string{platform.KeyWorkflowExecutionID: "dd3708ac7d8dd6fa4fae0fb87b73f318a4da2526c123e159b72435e3b2fe8751"},
			err:    err,
			want:   "workflowExecutionID dd3708ac7d8dd6fa4fae0fb87b73f318a4da2526c123e159b72435e3b2fe8751: some error",
		},
		{
			name:   "Error with no error and reason",
			labels: map[string]string{platform.KeyCapabilityID: "streams-trigger:network_eth@1.0.0"},
			reason: "some reason",
			want:   "capabilityID streams-trigger:network_eth@1.0.0: some reason",
		},
		{
			name:   "Error with no error and no reason",
			labels: map[string]string{platform.KeyTriggerID: "wf_123_trigger_456"},
			want:   "triggerID wf_123_trigger_456: ",
		},
		{
			name:   "Error with no labels",
			labels: map[string]string{},
			err:    err,
			reason: "some reason",
			want:   "some reason: some error",
		},
		{
			name: "Multiple labels",
			labels: map[string]string{
				platform.KeyWorkflowID:          "my-workflow-id",
				platform.KeyWorkflowExecutionID: "dd3708ac7d8dd6fa4fae0fb87b73f318a4da2526c123e159b72435e3b2fe8751",
				platform.KeyCapabilityID:        "streams-trigger:network_eth@1.0.0",
			},
			err:    err,
			reason: "some reason",
			want:   "workflowID my-workflow-id: workflowExecutionID dd3708ac7d8dd6fa4fae0fb87b73f318a4da2526c123e159b72435e3b2fe8751: capabilityID streams-trigger:network_eth@1.0.0: some reason: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &workflowError{
				labels: tt.labels,
				err:    tt.err,
				reason: tt.reason,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("err string mismatch\ngot = %v\nwant = %v", got, tt.want)
			}
		})
	}
}

func TestEngine_MergesWorkflowConfigAndCRConfig(t *testing.T) {
	var (
		ctx            = testutils.Context(t)
		writeID        = "write_polygon-testnet-mumbai@1.0.0"
		gotConfig      = values.EmptyMap()
		wantConfigKeys = []string{"deltaStage", "schedule", "address", "params", "abi"}
	)

	giveRegistryConfig, err := values.WrapMap(map[string]any{
		"deltaStage": "1s",
		"schedule":   "allAtOnce",
	})
	assert.NoError(t, err, "failed to wrap map of registry config")

	// Mock the capabilities of the simple workflow.
	reg := coreCap.NewRegistry(logger.TestLogger(t))
	trigger, _ := mockTrigger(t)
	consensus := mockConsensus("")
	target := newMockCapability(
		// Create a remote capability so we don't use the local transmission protocol.
		capabilities.MustNewRemoteCapabilityInfo(
			writeID,
			capabilities.CapabilityTypeTarget,
			"a write capability targeting polygon testnet",
			&capabilities.DON{ID: 1},
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			// Replace the empty config with the write target config.
			gotConfig = req.Config

			return capabilities.CapabilityResponse{
				Value: req.Inputs,
			}, nil
		},
	)

	require.NoError(t, reg.Add(ctx, trigger))
	require.NoError(t, reg.Add(ctx, consensus))
	require.NoError(t, reg.Add(ctx, target))

	eng, testHooks := newTestEngineWithYAMLSpec(
		t,
		reg,
		simpleWorkflow,
	)
	reg.SetLocalRegistry(testConfigProvider{
		configForCapability: func(ctx context.Context, capabilityID string, donID uint32) (registrysyncer.CapabilityConfiguration, error) {
			if capabilityID != writeID {
				return registrysyncer.CapabilityConfiguration{}, nil
			}

			var cb []byte
			cb, err = proto.Marshal(&capabilitiespb.CapabilityConfig{
				DefaultConfig: values.ProtoMap(giveRegistryConfig),
			})
			return registrysyncer.CapabilityConfiguration{
				Config: cb,
			}, err
		},
	})

	servicetest.Run(t, eng)

	eid := getExecutionID(t, eng, testHooks)

	state, err := eng.executionStates.Get(ctx, eid)
	require.NoError(t, err)

	assert.Equal(t, state.Status, store.StatusCompleted)

	// Assert that the config from the CR is merged with the default config from the registry.
	m, err := values.Unwrap(gotConfig)
	require.NoError(t, err)
	assert.Equal(t, m.(map[string]any)["deltaStage"], "1s")
	assert.Equal(t, m.(map[string]any)["schedule"], "allAtOnce")

	for _, k := range wantConfigKeys {
		assert.Contains(t, m.(map[string]any), k)
	}
}

const customComputeWorkflow = `
triggers:
  - id: "mercury-trigger@1.0.0"
    config:
      feedlist:
        - "0x1111111111111111111100000000000000000000000000000000000000000000" # ETHUSD
        - "0x2222222222222222222200000000000000000000000000000000000000000000" # LINKUSD
        - "0x3333333333333333333300000000000000000000000000000000000000000000" # BTCUSD

actions:
  - id: custom-compute@1.0.0
    ref: custom-compute
    config:
      maxMemoryMBs: 128
      tickInterval: 100ms
      timeout: 300ms
    inputs:
      action:
        - $(trigger.outputs)

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
  - id: "write_ethereum-testnet-sepolia@1.0.0"
    inputs: "$(evm_median.outputs)"
    config:
      address: "0x54e220867af6683aE6DcBF535B4f952cB5116510"
      params: ["$(report)"]
      abi: "receive(report bytes)"
`

// TestEngine_MergesWorkflowConfigAndCRConfig_CRConfigPrecedence tests that the engine merges the
// workflow config with the CR config correctly, with the CR config taking precedence.
func TestEngine_MergesWorkflowConfigAndCRConfig_CRConfigPrecedence(t *testing.T) {
	var (
		ctx              = testutils.Context(t)
		actionID         = "custom-compute@1.0.0"
		giveTimeout      = 300 * time.Millisecond
		giveTickInterval = 100 * time.Millisecond
		registryConfig   = map[string]any{
			"maxMemoryMBs": int64(64),
			"timeout":      giveTimeout.String(),
			"tickInterval": giveTickInterval.String(),
		}
		gotConfig = values.EmptyMap()
	)

	giveRegistryConfig, err := values.WrapMap(registryConfig)
	assert.NoError(t, err, "failed to wrap map of registry config")

	// Mock the capabilities of the simple workflow.
	reg := coreCap.NewRegistry(logger.TestLogger(t))
	trigger, _ := mockTrigger(t)
	target := mockTarget("write_ethereum-testnet-sepolia@1.0.0")
	action := newMockCapability(
		// Create a remote capability so we don't use the local transmission protocol.
		capabilities.MustNewRemoteCapabilityInfo(
			actionID,
			capabilities.CapabilityTypeAction,
			"a custom compute action with custom config",
			&capabilities.DON{ID: 1},
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			// Replace the empty config with the write target config.
			gotConfig = req.Config

			return capabilities.CapabilityResponse{
				Value: req.Inputs,
			}, nil
		},
	)

	consensus := mockConsensus("")

	require.NoError(t, reg.Add(ctx, trigger))
	require.NoError(t, reg.Add(ctx, action))
	require.NoError(t, reg.Add(ctx, target))
	require.NoError(t, reg.Add(ctx, consensus))

	eng, testHooks := newTestEngineWithYAMLSpec(
		t,
		reg,
		customComputeWorkflow,
	)
	reg.SetLocalRegistry(testConfigProvider{
		configForCapability: func(ctx context.Context, capabilityID string, donID uint32) (registrysyncer.CapabilityConfiguration, error) {
			if capabilityID != actionID {
				return registrysyncer.CapabilityConfiguration{}, nil
			}

			var cb []byte
			cb, err = proto.Marshal(&capabilitiespb.CapabilityConfig{
				RestrictedConfig: values.ProtoMap(giveRegistryConfig),
				RestrictedKeys:   []string{"maxMemoryMBs", "tickInterval", "timeout"},
			})
			return registrysyncer.CapabilityConfiguration{
				Config: cb,
			}, err
		},
	})

	servicetest.Run(t, eng)

	eid := getExecutionID(t, eng, testHooks)

	state, err := eng.executionStates.Get(ctx, eid)
	require.NoError(t, err)

	assert.Equal(t, state.Status, store.StatusCompleted)

	// Assert that the config from the CR is merged with the default config from the registry. With
	// the CR config taking precedence.
	m, err := values.Unwrap(gotConfig)
	require.NoError(t, err)
	assert.Equalf(t, registryConfig["maxMemoryMBs"], m.(map[string]any)["maxMemoryMBs"], "maxMemoryMBs should be %d", registryConfig["maxMemoryMBs"])
	assert.Equalf(t, registryConfig["timeout"], m.(map[string]any)["timeout"], "timeout should be %s", registryConfig["timeout"])
	assert.Equalf(t, registryConfig["tickInterval"], m.(map[string]any)["tickInterval"], "tickInterval should be %s", registryConfig["tickInterval"])
}

func TestEngine_HandlesNilConfigOnchain(t *testing.T) {
	ctx := testutils.Context(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))

	trigger, _ := mockTrigger(t)

	require.NoError(t, reg.Add(ctx, trigger))
	require.NoError(t, reg.Add(ctx, mockConsensus("")))
	writeID := "write_polygon-testnet-mumbai@1.0.0"

	gotConfig := values.EmptyMap()
	target := newMockCapability(
		// Create a remote capability so we don't use the local transmission protocol.
		capabilities.MustNewRemoteCapabilityInfo(
			writeID,
			capabilities.CapabilityTypeTarget,
			"a write capability targeting polygon testnet",
			&capabilities.DON{ID: 1},
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			gotConfig = req.Config

			return capabilities.CapabilityResponse{
				Value: req.Inputs,
			}, nil
		},
	)
	require.NoError(t, reg.Add(ctx, target))

	eng, testHooks := newTestEngineWithYAMLSpec(
		t,
		reg,
		simpleWorkflow,
	)
	reg.SetLocalRegistry(testConfigProvider{})

	servicetest.Run(t, eng)

	eid := getExecutionID(t, eng, testHooks)

	state, err := eng.executionStates.Get(ctx, eid)
	require.NoError(t, err)

	assert.Equal(t, state.Status, store.StatusCompleted)

	m, err := values.Unwrap(gotConfig)
	require.NoError(t, err)
	// The write target config contains three keys
	assert.Len(t, m.(map[string]any), 3)
}

func TestEngine_MultiBranchExecution(t *testing.T) {
	// This workflow describes 2 branches in the workflow graph.
	// A -> B -> C
	// A -> D -> E
	workflowSpec := `
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
  - id: "early_exit_offchain_reporting@1.0.0"
    ref: "evm_median_early_exit"
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
  - id: "write_polygon-testnet-early-exit@1.0.0"
    inputs:
      report: "$(evm_median_early_exit.outputs.report)"
    config:
      address: "0x3F3554832c636721F1fD1822Ccca0354576741Ef"
      params: ["$(report)"]
      abi: "receive(report bytes)"
`
	ctx := testutils.Context(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))

	trigger, _ := mockTrigger(t)
	require.NoError(t, reg.Add(ctx, trigger))
	require.NoError(t, reg.Add(ctx, mockConsensus("")))
	require.NoError(t, reg.Add(ctx, mockConsensusWithEarlyTermination("early_exit_offchain_reporting@1.0.0")))
	require.NoError(t, reg.Add(ctx, mockTarget("")))
	require.NoError(t, reg.Add(ctx, mockTarget("write_polygon-testnet-early-exit@1.0.0")))

	eng, hooks := newTestEngineWithYAMLSpec(t, reg, workflowSpec)
	servicetest.Run(t, eng)

	eid := getExecutionID(t, eng, hooks)
	state, err := eng.executionStates.Get(ctx, eid)
	require.NoError(t, err)

	assert.Equal(t, store.StatusCompletedEarlyExit, state.Status)
}

func basicTestTrigger(t *testing.T) *mockTriggerCapability {
	mt := &mockTriggerCapability{
		CapabilityInfo: capabilities.MustNewCapabilityInfo(
			"basic-test-trigger@1.0.0",
			capabilities.CapabilityTypeTrigger,
			"basic test trigger",
		),
		ch:                         make(chan capabilities.TriggerResponse, 10),
		registerTriggerCallCounter: make(map[string]int),
	}

	resp, err := values.NewMap(map[string]any{
		"cool_output": "foo",
	})
	require.NoError(t, err)
	tr := capabilities.TriggerResponse{
		Event: capabilities.TriggerEvent{
			TriggerType: mt.ID,
			ID:          time.Now().UTC().Format(time.RFC3339),
			Outputs:     resp,
		},
	}
	mt.triggerEvent = &tr
	return mt
}

func TestEngine_WithCustomComputeStep(t *testing.T) {
	cmd := "core/services/workflows/test/wasm/cmd"
	binary := "test/wasm/cmd/testmodule.wasm"

	ctx := testutils.Context(t)
	log := logger.TestLogger(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))
	cfg := compute.Config{
		ServiceConfig: webapi.ServiceConfig{
			RateLimiter: common.RateLimiterConfig{
				GlobalRPS:      100.0,
				GlobalBurst:    100,
				PerSenderRPS:   100.0,
				PerSenderBurst: 100,
			},
		},
	}

	connector := gcmocks.NewGatewayConnector(t)
	connector.EXPECT().GatewayIDs().Return([]string{"gateway1"})
	handler, err := webapi.NewOutgoingConnectorHandler(
		connector,
		cfg.ServiceConfig,
		ghcapabilities.MethodComputeAction, log)
	require.NoError(t, err)

	idGeneratorFn := func() string { return "validRequestID" }
	fetcher, err := compute.NewOutgoingConnectorFetcherFactory(handler, idGeneratorFn)
	require.NoError(t, err)
	compute, err := compute.NewAction(cfg, log, reg, fetcher)
	require.NoError(t, err)
	require.NoError(t, compute.Start(ctx))
	defer compute.Close()

	trigger := basicTestTrigger(t)
	require.NoError(t, reg.Add(ctx, trigger))

	binaryB := wasmtest.CreateTestBinary(cmd, binary, true, t)

	spec, err := host.GetWorkflowSpec(
		ctx,
		&host.ModuleConfig{Logger: log},
		binaryB,
		nil, // config
	)
	require.NoError(t, err)
	eng, testHooks := newTestEngine(
		t,
		reg,
		*spec,
		func(c *Config) {
			c.Binary = binaryB
			c.Config = nil
		},
	)
	reg.SetLocalRegistry(testConfigProvider{})

	servicetest.Run(t, eng)

	eid := getExecutionID(t, eng, testHooks)

	state, err := eng.executionStates.Get(ctx, eid)
	require.NoError(t, err)

	assert.Equal(t, state.Status, store.StatusCompleted)
	res, ok := state.ResultForStep("compute")
	assert.True(t, ok)
	assert.True(t, res.Outputs.(*values.Map).Underlying["Value"].(*values.Bool).Underlying)
}

func TestEngine_CustomComputePropagatesBreaks(t *testing.T) {
	cmd := "core/services/workflows/test/break/cmd"
	binary := "test/wasm/break/testmodule.wasm"

	ctx := testutils.Context(t)
	log := logger.TestLogger(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))
	cfg := compute.Config{
		ServiceConfig: webapi.ServiceConfig{
			RateLimiter: common.RateLimiterConfig{
				GlobalRPS:      100.0,
				GlobalBurst:    100,
				PerSenderRPS:   100.0,
				PerSenderBurst: 100,
			},
		},
	}
	connector := gcmocks.NewGatewayConnector(t)
	connector.EXPECT().GatewayIDs().Return([]string{"gateway1"})
	handler, err := webapi.NewOutgoingConnectorHandler(
		connector,
		cfg.ServiceConfig,
		ghcapabilities.MethodComputeAction, log)
	require.NoError(t, err)

	idGeneratorFn := func() string { return "validRequestID" }
	fetcher, err := compute.NewOutgoingConnectorFetcherFactory(handler, idGeneratorFn)
	require.NoError(t, err)
	compute, err := compute.NewAction(cfg, log, reg, fetcher)
	require.NoError(t, err)
	require.NoError(t, compute.Start(ctx))
	defer compute.Close()

	trigger := basicTestTrigger(t)
	require.NoError(t, reg.Add(ctx, trigger))

	binaryB := wasmtest.CreateTestBinary(cmd, binary, true, t)

	spec, err := host.GetWorkflowSpec(
		ctx,
		&host.ModuleConfig{Logger: log},
		binaryB,
		nil, // config
	)
	require.NoError(t, err)
	eng, testHooks := newTestEngine(
		t,
		reg,
		*spec,
		func(c *Config) {
			c.Binary = binaryB
			c.Config = nil
		},
	)
	reg.SetLocalRegistry(testConfigProvider{})

	servicetest.Run(t, eng)

	eid := getExecutionID(t, eng, testHooks)

	state, err := eng.executionStates.Get(ctx, eid)
	require.NoError(t, err)

	assert.Equal(t, store.StatusCompletedEarlyExit, state.Status)
}

const secretsWorkflow = `
triggers:
  - id: "mercury-trigger@1.0.0"
    config:
      feedlist:
        - "0x1111111111111111111100000000000000000000000000000000000000000000" # ETHUSD
        - "0x2222222222222222222200000000000000000000000000000000000000000000" # LINKUSD
        - "0x3333333333333333333300000000000000000000000000000000000000000000" # BTCUSD

actions:
  - id: custom-compute@1.0.0
    ref: custom-compute
    config:
      fidelityToken: $(ENV.secrets.fidelity)
    inputs:
      action:
        - $(trigger.outputs)

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
  - id: "write_ethereum-testnet-sepolia@1.0.0"
    inputs: "$(evm_median.outputs)"
    config:
      address: "0x54e220867af6683aE6DcBF535B4f952cB5116510"
      params: ["$(report)"]
      abi: "receive(report bytes)"
`

type mockFetcher struct {
	retval map[string]string
	retErr error
}

func (m *mockFetcher) SecretsFor(ctx context.Context, workflowOwner, hexWorkflowName, decodedWorkflowName, workflowID string) (map[string]string, error) {
	return m.retval, m.retErr
}

func TestEngine_FetchesSecrets(t *testing.T) {
	ctx := testutils.Context(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))

	trigger, _ := mockTrigger(t)
	require.NoError(t, reg.Add(ctx, trigger))

	require.NoError(t, reg.Add(ctx, mockConsensus("")))

	target := mockTarget("write_ethereum-testnet-sepolia@1.0.0")
	require.NoError(t, reg.Add(ctx, target))

	var gotConfig *values.Map
	action := newMockCapability(
		// Create a remote capability so we don't use the local transmission protocol.
		capabilities.MustNewRemoteCapabilityInfo(
			"custom-compute@1.0.0",
			capabilities.CapabilityTypeAction,
			"a custom compute action with custom config",
			&capabilities.DON{ID: 1},
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			// Replace the empty config with the write target config.
			gotConfig = req.Config

			return capabilities.CapabilityResponse{
				Value: req.Inputs,
			}, nil
		},
	)
	require.NoError(t, reg.Add(ctx, action))

	t.Run("successfully fetches secrets", func(t *testing.T) {
		eng, testHooks := newTestEngineWithYAMLSpec(
			t,
			reg,
			secretsWorkflow,
			func(c *Config) {
				c.SecretsFetcher = &mockFetcher{
					retval: map[string]string{
						"fidelity": "aFidelitySecret",
					},
				}
			},
		)

		servicetest.Run(t, eng)

		eid := getExecutionID(t, eng, testHooks)

		state, err := eng.executionStates.Get(ctx, eid)
		require.NoError(t, err)

		assert.Equal(t, store.StatusCompleted, state.Status)

		expected := map[string]any{
			"fidelityToken": "aFidelitySecret",
		}
		expm, err := values.Wrap(expected)
		require.NoError(t, err)
		assert.Equal(t, gotConfig, expm)
	})
}

func TestEngine_CloseHappensOnlyIfWorkflowHasBeenRegistered(t *testing.T) {
	ctx := testutils.Context(t)
	reg := coreCap.NewRegistry(logger.TestLogger(t))

	trigger, _ := mockTrigger(t)

	require.NoError(t, reg.Add(ctx, trigger))

	require.NoError(t, reg.Add(ctx, mockConsensus("")))

	target := mockTarget("write_ethereum-testnet-sepolia@1.0.0")
	require.NoError(t, reg.Add(ctx, target))

	action := newMockCapability(
		// Create a remote capability so we don't use the local transmission protocol.
		capabilities.MustNewRemoteCapabilityInfo(
			"custom-compute@1.0.0",
			capabilities.CapabilityTypeAction,
			"a custom compute action with custom config",
			&capabilities.DON{ID: 1},
		),
		func(req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
			return capabilities.CapabilityResponse{
				Value: req.Inputs,
			}, nil
		},
	)
	require.NoError(t, reg.Add(ctx, action))

	eng, testHooks := newTestEngineWithYAMLSpec(
		t,
		reg,
		secretsWorkflow,
		func(c *Config) {
			c.SecretsFetcher = &mockFetcher{
				retval: map[string]string{},
				retErr: errors.New("failed to fetch secrets XXX"),
			}
		},
	)

	err := eng.Start(ctx)
	require.NoError(t, err)

	// simulate WorkflowUpdatedEvent that calls tryEngineCleanup
	<-testHooks.initFailed
	err = eng.Close()
	require.NoError(t, err)
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name             string
		baseConfig       map[string]any
		expectedConfig   map[string]any
		capabilityConfig capabilities.CapabilityConfiguration
	}{
		{
			name: "no remote config",
			baseConfig: map[string]any{
				"foo": "bar",
			},
			expectedConfig: map[string]any{
				"foo": "bar",
			},
			capabilityConfig: capabilities.CapabilityConfiguration{},
		},
		{
			name: "user provides restricted config",
			baseConfig: map[string]any{
				"restrictedXXX": "restrictedYYY",
				"foo":           "bar",
			},
			expectedConfig: map[string]any{
				"foo": "bar",
			},
			capabilityConfig: capabilities.CapabilityConfiguration{
				RestrictedKeys: []string{"restrictedXXX"},
			},
		},
		{
			name: "user provides restricted config; capability contains restricted",
			baseConfig: map[string]any{
				"restrictedXXX": "restrictedYYY",
				"foo":           "bar",
			},
			expectedConfig: map[string]any{
				"foo":           "bar",
				"restrictedXXX": "restrictedXXXSetRemotely",
			},
			capabilityConfig: capabilities.CapabilityConfiguration{
				RestrictedKeys: []string{"restrictedXXX"},
				RestrictedConfig: &values.Map{
					Underlying: map[string]values.Value{
						"restrictedXXX": values.NewString("restrictedXXXSetRemotely"),
					},
				},
			},
		},
		{
			name: "default overridden by what user provides",
			baseConfig: map[string]any{
				"restrictedXXX": "restrictedYYY",
				"foo":           "bar",
				"baz":           "overridden",
			},
			expectedConfig: map[string]any{
				"foo": "bar",
				"baz": "overridden",
			},
			capabilityConfig: capabilities.CapabilityConfiguration{
				RestrictedKeys: []string{"restrictedXXX"},
				DefaultConfig: &values.Map{
					Underlying: map[string]values.Value{
						"baz": values.NewString("qux"),
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(st *testing.T) {
			bc, err := values.NewMap(tc.baseConfig)
			require.NoError(t, err)
			got := merge(bc, tc.capabilityConfig)
			gotMap, err := got.Unwrap()
			require.NoError(t, err)
			assert.Equal(t, tc.expectedConfig, gotMap)
		})
	}
}
