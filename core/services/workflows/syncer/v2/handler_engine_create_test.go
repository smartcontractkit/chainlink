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
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
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
	triggers.TriggerCoordinator
	registerErr error
	registered  []triggers.Subscriber
}

func (c *recordingCoordinator) RegisterTriggers(_ context.Context, subscriber triggers.Subscriber, _ triggers.RegistrationParams) ([]string, error) {
	c.registered = append(c.registered, subscriber)
	return nil, c.registerErr
}

// nodeConfig is the node state that decides the routing policy.
type nodeConfig struct {
	flagOpen        bool
	withCoordinator bool
	sharded         bool
}

// engineCreateFixture is an eventHandler whose engine factory mirrors the real
// one's contract: it reads the routing policy itself and returns the matching
// engine, leaving tryEngineCreate to route on the engine's IsCoordinated.
type engineCreateFixture struct {
	h                 *eventHandler
	legacyEngine      v2.WorkflowEngine
	coordinatedEngine *fakeCoordinatedEngine
	factoryCalls      atomic.Int32
	factoryErr        error
	coordinator       *recordingCoordinator
}

func newEngineCreateFixture(t *testing.T, cfg nodeConfig) *engineCreateFixture {
	t.Helper()
	lggr := logger.Test(t)
	registry := capreg.NewRegistry(lggr)
	registry.SetRegistryMetadata(&capreg.TestRegistryMetadata{})

	f := &engineCreateFixture{
		legacyEngine:      &mockEngine{},
		coordinatedEngine: &fakeCoordinatedEngine{},
	}
	// A sharded node gets the ShardFailoverManager wrapper, as the real factory
	// builds it.
	if cfg.sharded {
		f.legacyEngine = newTestShardFailoverManager(t, &mockEngine{})
	}

	f.h = &eventHandler{
		lggr:           lggr,
		capRegistry:    registry,
		engineRegistry: NewEngineRegistry(),
		featureFlags:   &v2.EngineFeatureFlags{CoordinatedEngine: limits.NewGateLimiter(cfg.flagOpen)},
		tracer:         noop.NewTracerProvider().Tracer(""),
		engineFactory: func(ctx context.Context, wfid, _ string, _ types.WorkflowName, _ string, _, _ []byte, _ string, initDone chan<- error) (v2.WorkflowEngine, error) {
			f.factoryCalls.Add(1)
			if f.factoryErr != nil {
				return nil, f.factoryErr
			}
			initDone <- nil
			if f.h.useCoordinatedEngine(ctx, wfid) {
				return f.coordinatedEngine, nil
			}
			return f.legacyEngine, nil
		},
	}
	if cfg.withCoordinator {
		f.coordinator = &recordingCoordinator{}
		f.h.triggerCoordinator = f.coordinator
	}
	if cfg.sharded {
		f.h.shardingEnabled = true
		f.h.dispatcher = newFakeDispatcher()
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

	t.Run("flag open with a coordinator takes the coordinated path", func(t *testing.T) {
		t.Parallel()
		f := newEngineCreateFixture(t, nodeConfig{flagOpen: true, withCoordinator: true})
		spec, wid := newTestWorkflowSpec(t)

		require.NoError(t, f.h.tryEngineCreate(t.Context(), spec, "TestSource"))

		require.Equal(t, int32(1), f.factoryCalls.Load())
		entry, ok := f.h.engineRegistry.Get(wid)
		require.True(t, ok)
		require.True(t, entry.Coordinated())
		require.NotEmpty(t, entry.ReconcileKey)
		require.Equal(t, []triggers.Subscriber{f.coordinatedEngine}, f.coordinator.registered)
	})

	t.Run("flag closed takes the legacy path", func(t *testing.T) {
		t.Parallel()
		f := newEngineCreateFixture(t, nodeConfig{flagOpen: false, withCoordinator: true})
		spec, wid := newTestWorkflowSpec(t)

		require.NoError(t, f.h.tryEngineCreate(t.Context(), spec, "TestSource"))

		require.Equal(t, int32(1), f.factoryCalls.Load())
		entry, ok := f.h.engineRegistry.Get(wid)
		require.True(t, ok)
		require.False(t, entry.Coordinated())
		require.Empty(t, f.coordinator.registered)
	})

	t.Run("flag open without a coordinator takes the legacy path", func(t *testing.T) {
		t.Parallel()
		f := newEngineCreateFixture(t, nodeConfig{flagOpen: true, withCoordinator: false})
		spec, wid := newTestWorkflowSpec(t)

		require.NoError(t, f.h.tryEngineCreate(t.Context(), spec, "TestSource"))

		entry, ok := f.h.engineRegistry.Get(wid)
		require.True(t, ok)
		require.False(t, entry.Coordinated())
	})

	t.Run("a sharded node stays on the legacy path with the flag open", func(t *testing.T) {
		t.Parallel()
		f := newEngineCreateFixture(t, nodeConfig{flagOpen: true, withCoordinator: true, sharded: true})
		spec, wid := newTestWorkflowSpec(t)

		require.NoError(t, f.h.tryEngineCreate(t.Context(), spec, "TestSource"))

		entry, ok := f.h.engineRegistry.Get(wid)
		require.True(t, ok)
		require.IsType(t, &ShardFailoverManager{}, entry.Service)
		require.False(t, entry.Coordinated(), "a sharded engine must never classify as coordinated")
		require.Empty(t, f.coordinator.registered, "the coordinator must not be used on a sharded node")
		require.Equal(t, engineCounts{legacy: 1}, countEngines(f.h.engineRegistry.GetAll()))
	})

	t.Run("RegisterTriggers failure removes and closes the coordinated engine", func(t *testing.T) {
		t.Parallel()
		f := newEngineCreateFixture(t, nodeConfig{flagOpen: true, withCoordinator: true})
		f.coordinator.registerErr = errors.New("boom")
		spec, wid := newTestWorkflowSpec(t)

		err := f.h.tryEngineCreate(t.Context(), spec, "TestSource")

		require.ErrorContains(t, err, "failed to register triggers via coordinator")
		require.False(t, f.h.engineRegistry.Contains(wid))
		require.True(t, f.coordinatedEngine.closed.Load())
	})

	t.Run("a factory failure is reported and registers nothing", func(t *testing.T) {
		t.Parallel()
		f := newEngineCreateFixture(t, nodeConfig{withCoordinator: true})
		f.factoryErr = errors.New("boom")
		spec, wid := newTestWorkflowSpec(t)

		err := f.h.tryEngineCreate(t.Context(), spec, "TestSource")

		require.ErrorContains(t, err, "failed to create workflow engine")
		require.False(t, f.h.engineRegistry.Contains(wid))
	})

	t.Run("invalid spec is rejected before any engine is built", func(t *testing.T) {
		t.Parallel()
		f := newEngineCreateFixture(t, nodeConfig{withCoordinator: true})
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
