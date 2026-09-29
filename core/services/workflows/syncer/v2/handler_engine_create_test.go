package v2

import (
	"context"
	"encoding/hex"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"

	capreg "github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

// fakeCoordinatedEngine is a v2.WorkflowEngine that reports as coordinated.
// Only the lifecycle methods tryEngineCreate calls are implemented.
type fakeCoordinatedEngine struct {
	v2.WorkflowEngine
	closed atomic.Bool
}

func (f *fakeCoordinatedEngine) Start(context.Context) error { return nil }

func (f *fakeCoordinatedEngine) Close() error {
	f.closed.Store(true)
	return nil
}

func (f *fakeCoordinatedEngine) IsCoordinated() bool { return true }

// recordingCoordinator records the subscribers handed to RegisterTriggers.
type recordingCoordinator struct {
	TriggerCoordinator
	registerErr error
	registered  []v2.Subscriber
}

func (c *recordingCoordinator) RegisterTriggers(_ context.Context, subscriber v2.Subscriber, _ RegistrationParams) ([]string, error) {
	c.registered = append(c.registered, subscriber)
	return nil, c.registerErr
}

// engineCreateFixture is an eventHandler whose engine factory returns a
// preset engine. The factory owns the flag read in production, so these tests
// cover the dispatch.
type engineCreateFixture struct {
	h            *eventHandler
	engine       v2.WorkflowEngine
	factoryCalls atomic.Int32
	factoryErr   error
	coordinator  *recordingCoordinator
}

// newEngineCreateFixture wires a handler whose factory hands back engine.
// Passing a fakeCoordinatedEngine exercises the coordinated path, a mockEngine
// the legacy one.
func newEngineCreateFixture(t *testing.T, engine v2.WorkflowEngine, withCoordinator bool) *engineCreateFixture {
	t.Helper()
	lggr := logger.Test(t)
	registry := capreg.NewRegistry(lggr)
	registry.SetRegistryMetadata(&capreg.TestRegistryMetadata{})

	f := &engineCreateFixture{engine: engine}
	f.h = &eventHandler{
		lggr:           lggr,
		capRegistry:    registry,
		engineRegistry: NewEngineRegistry(),
		tracer:         noop.NewTracerProvider().Tracer(""),
		engineFactory: func(_ context.Context, _, _ string, _ types.WorkflowName, _ string, _, _ []byte, _ string, initDone chan<- error) (v2.WorkflowEngine, error) {
			f.factoryCalls.Add(1)
			if f.factoryErr != nil {
				return nil, f.factoryErr
			}
			initDone <- nil
			return f.engine, nil
		},
	}
	if withCoordinator {
		f.coordinator = &recordingCoordinator{}
		f.h.triggerCoordinator = f.coordinator
	}
	return f
}

// newTestWorkflowSpec returns a spec whose workflow ID matches its artifacts,
// so it passes prepareEngineInputs.
func newTestWorkflowSpec(t *testing.T) (*job.WorkflowSpec, types.WorkflowID) {
	t.Helper()
	owner := []byte{0x01, 0x02, 0x03}
	name := "engine-create-test"
	binary := []byte("binary")
	config := []byte("config")
	id, err := pkgworkflows.GenerateWorkflowID(owner, name, binary, config, "")
	require.NoError(t, err)
	wid := types.WorkflowID(id)
	return &job.WorkflowSpec{
		Workflow:      hex.EncodeToString(binary),
		Config:        string(config),
		WorkflowID:    wid.Hex(),
		WorkflowOwner: hex.EncodeToString(owner),
		WorkflowName:  name,
	}, wid
}

func Test_tryEngineCreate_routing(t *testing.T) {
	t.Parallel()

	t.Run("a coordinated engine is registered and handed to the coordinator", func(t *testing.T) {
		t.Parallel()
		engine := &fakeCoordinatedEngine{}
		f := newEngineCreateFixture(t, engine, true)
		spec, wid := newTestWorkflowSpec(t)

		require.NoError(t, f.h.tryEngineCreate(t.Context(), spec, "TestSource"))

		require.Equal(t, int32(1), f.factoryCalls.Load())
		entry, ok := f.h.engineRegistry.Get(wid)
		require.True(t, ok)
		require.True(t, entry.Coordinated())
		require.NotEmpty(t, entry.ReconcileKey)
		require.Equal(t, []v2.Subscriber{engine}, f.coordinator.registered)
	})

	t.Run("a legacy engine is registered without touching the coordinator", func(t *testing.T) {
		t.Parallel()
		f := newEngineCreateFixture(t, &mockEngine{}, true)
		spec, wid := newTestWorkflowSpec(t)

		require.NoError(t, f.h.tryEngineCreate(t.Context(), spec, "TestSource"))

		require.Equal(t, int32(1), f.factoryCalls.Load())
		entry, ok := f.h.engineRegistry.Get(wid)
		require.True(t, ok)
		require.False(t, entry.Coordinated())
		require.Empty(t, f.coordinator.registered)
	})

	t.Run("RegisterTriggers failure removes and closes the coordinated engine", func(t *testing.T) {
		t.Parallel()
		engine := &fakeCoordinatedEngine{}
		f := newEngineCreateFixture(t, engine, true)
		f.coordinator.registerErr = errors.New("boom")
		spec, wid := newTestWorkflowSpec(t)

		err := f.h.tryEngineCreate(t.Context(), spec, "TestSource")

		require.ErrorContains(t, err, "failed to register triggers via coordinator")
		require.False(t, f.h.engineRegistry.Contains(wid))
		require.True(t, engine.closed.Load())
	})

	t.Run("a factory failure is reported and registers nothing", func(t *testing.T) {
		t.Parallel()
		f := newEngineCreateFixture(t, &mockEngine{}, true)
		f.factoryErr = errors.New("boom")
		spec, wid := newTestWorkflowSpec(t)

		err := f.h.tryEngineCreate(t.Context(), spec, "TestSource")

		require.ErrorContains(t, err, "failed to create workflow engine")
		require.False(t, f.h.engineRegistry.Contains(wid))
	})

	t.Run("invalid spec is rejected before any engine is built", func(t *testing.T) {
		t.Parallel()
		f := newEngineCreateFixture(t, &mockEngine{}, true)
		spec, _ := newTestWorkflowSpec(t)
		spec.WorkflowID = types.WorkflowID{0xff}.Hex()

		err := f.h.tryEngineCreate(t.Context(), spec, "TestSource")

		require.ErrorContains(t, err, "workflowID mismatch")
		require.Equal(t, int32(0), f.factoryCalls.Load())
	})
}

// Test_useCoordinatedEngine covers the routing policy the engine factory reads.
func Test_useCoordinatedEngine(t *testing.T) {
	t.Parallel()

	newHandler := func(flagOpen, withCoordinator, sharded bool) *eventHandler {
		h := &eventHandler{
			lggr:         logger.Test(t),
			featureFlags: &v2.EngineFeatureFlags{CoordinatedEngine: limits.NewGateLimiter(flagOpen)},
		}
		if withCoordinator {
			h.triggerCoordinator = &recordingCoordinator{}
		}
		if sharded {
			h.shardingEnabled = true
			h.dispatcher = newFakeDispatcher()
		}
		return h
	}

	t.Run("flag open with a coordinator", func(t *testing.T) {
		t.Parallel()
		require.True(t, newHandler(true, true, false).useCoordinatedEngine(t.Context(), "wfid"))
	})

	t.Run("flag closed", func(t *testing.T) {
		t.Parallel()
		require.False(t, newHandler(false, true, false).useCoordinatedEngine(t.Context(), "wfid"))
	})

	t.Run("no coordinator wired", func(t *testing.T) {
		t.Parallel()
		require.False(t, newHandler(true, false, false).useCoordinatedEngine(t.Context(), "wfid"))
	})

	t.Run("sharded nodes stay on the legacy engine", func(t *testing.T) {
		t.Parallel()
		require.False(t, newHandler(true, true, true).useCoordinatedEngine(t.Context(), "wfid"))
	})

	t.Run("nil feature flags", func(t *testing.T) {
		t.Parallel()
		h := newHandler(true, true, false)
		h.featureFlags = nil
		require.False(t, h.useCoordinatedEngine(t.Context(), "wfid"))
	})
}
