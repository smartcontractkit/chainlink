package v2_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
)

// withStubbedLocalNode replaces the config's capabilities registry with a mock
// whose LocalNode is stubbed, which NewEngine/NewExecutionEngine require.
func withStubbedLocalNode(t *testing.T, cfg *v2.EngineConfig) {
	t.Helper()
	reg := regmocks.NewCapabilitiesRegistry(t)
	reg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)
	cfg.CapRegistry = reg
}

// withNoopAcknowledger sets a no-op TriggerAcknowledger, which NewExecutionEngine
// requires. (Engine self-injects when nil, so it must not be set for Engine.)
func withNoopAcknowledger(cfg *v2.EngineConfig) {
	cfg.TriggerAcknowledger = noopAcknowledger{}
}

// engineImpl pairs an engine constructor with whether it is the execution-only
// variant (which the test coordinator must register/ACK on behalf of).
type engineImpl struct {
	ctor          engineCtor
	executionOnly bool
}

// engineImpls maps each engine implementation name to its constructor metadata.
// Tests parameterized over it exercise both types through the shared
// WorkflowEngine interface. Engine owns trigger registration; ExecutionEngine is
// execution-only and depends on an injected acknowledger (the test coordinator).
var engineImpls = map[string]engineImpl{
	"Engine": {ctor: func(c *v2.EngineConfig) (v2.WorkflowEngine, error) { return v2.NewEngine(c) }, executionOnly: false},
	"ExecutionEngine": {ctor: func(c *v2.EngineConfig) (v2.WorkflowEngine, error) {
		return v2.NewExecutionEngine(c)
	}, executionOnly: true},
}

// forEachEngineImpl runs fn once per engine implementation, as a named subtest. fn
// receives a config and must construct the engine via newCoordinatedEngine with
// the supplied impl.
func forEachEngineImpl(t *testing.T, fn func(t *testing.T, impl engineImpl)) {
	t.Helper()
	for name, impl := range engineImpls {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fn(t, impl)
		})
	}
}

// TestEngine_SatisfiesWorkflowEngine is a construction smoke test parameterized
// over both constructors. It proves both types satisfy the WorkflowEngine
// interface and can be constructed from the same valid config.
func TestEngine_SatisfiesWorkflowEngine(t *testing.T) {
	t.Parallel()
	for name, impl := range engineImpls {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			cfg := defaultTestConfig(t, nil)
			withStubbedLocalNode(t, cfg)
			if impl.executionOnly {
				withNoopAcknowledger(cfg)
			}

			engine, err := impl.ctor(cfg)
			require.NoError(t, err)
			require.NotNil(t, engine)
			require.Equal(t, "WorkflowEngine.WorkflowEngineV2", engine.Name())
		})
	}
}

// TestExecutionEngine_SatisfiesInterfaces asserts the interface contracts the
// execution-only engine is expected to hold.
func TestExecutionEngine_SatisfiesInterfaces(t *testing.T) {
	t.Parallel()
	cfg := defaultTestConfig(t, nil)
	withStubbedLocalNode(t, cfg)
	withNoopAcknowledger(cfg)

	engine, err := v2.NewExecutionEngine(cfg)
	require.NoError(t, err)
	require.NotNil(t, engine)

	var _ v2.WorkflowEngine = engine
	var _ v2.EventSink = engine
}
