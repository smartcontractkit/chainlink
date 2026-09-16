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

// engineCtors maps each engine implementation name to its constructor. It is the
// table that proves ExecutionEngine is a faithful, behaviorally-identical copy of
// Engine: any test parameterized over it exercises both types identically. When
// CRE-6176 stripping begins, tests that encode behavior being moved to the
// TriggerDispatcher will start failing here — informatively — telling the owner
// exactly which assertions are affected.
var engineCtors = map[string]func(*v2.EngineConfig) (v2.WorkflowEngine, error){
	"Engine":          func(c *v2.EngineConfig) (v2.WorkflowEngine, error) { return v2.NewEngine(c) },
	"ExecutionEngine": func(c *v2.EngineConfig) (v2.WorkflowEngine, error) { return v2.NewExecutionEngine(c) },
}

// TestEngine_SatisfiesWorkflowEngine is a construction smoke test parameterized
// over both constructors. It proves both types satisfy the WorkflowEngine
// interface and can be constructed from the same valid config.
func TestEngine_SatisfiesWorkflowEngine(t *testing.T) {
	t.Parallel()
	for name, ctor := range engineCtors {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			cfg := defaultTestConfig(t, nil)
			withStubbedLocalNode(t, cfg)

			engine, err := ctor(cfg)
			require.NoError(t, err)
			require.NotNil(t, engine)

			// The interface is satisfied (compile-time assertions live in the
			// source files); this is the runtime construction smoke test.
			var _ v2.WorkflowEngine = engine
			// Both engines report the identical service name (see plan §2e): the
			// name feeds health-report keys, so it must not differ between them.
			require.Equal(t, "WorkflowEngine.WorkflowEngineV2", engine.Name())
		})
	}
}

// TestExecutionEngine_SatisfiesInterfaces asserts the interface contracts the
// copy is expected to hold today. The Acknowledger assertion is expected to be
// deleted during CRE-6176 stripping — its removal is the signal that AC 4 landed.
func TestExecutionEngine_SatisfiesInterfaces(t *testing.T) {
	t.Parallel()
	cfg := defaultTestConfig(t, nil)
	withStubbedLocalNode(t, cfg)

	engine, err := v2.NewExecutionEngine(cfg)
	require.NoError(t, err)
	require.NotNil(t, engine)

	var _ v2.WorkflowEngine = engine
	var _ v2.EventSink = engine
	var _ v2.Acknowledger = engine // self-injects today; removed when ACK moves to the dispatcher
}
