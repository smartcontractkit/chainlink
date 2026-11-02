package v2

import (
	"context"
	"encoding/hex"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"

	capreg "github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
)

// fakeCoordinatedEngine is a v2.WorkflowEngine that reports as coordinated.
// Only the lifecycle methods tryEngineCreate calls, and the drain state the
// registry's engine counts read through a ShardFailoverManager, are implemented.
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

func (f *fakeCoordinatedEngine) DrainStartedAt() (time.Time, bool) { return time.Time{}, false }

// recordingCoordinator records the subscribers handed to RegisterTriggers.
type recordingCoordinator struct {
	triggers.Coordinator
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
// engine, leaving tryEngineCreate to route on the engine's IsCoordinated. On a
// sharded node it wraps whichever engine it picks in a ShardFailoverManager,
// as the real factory does.
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

	f.h = &eventHandler{
		lggr:           lggr,
		capRegistry:    registry,
		engineRegistry: NewEngineRegistry(),
		featureFlags:   &v2.EngineFeatureFlags{CoordinatedEngine: limits.NewGateLimiter(cfg.flagOpen)},
		tracer:         noop.NewTracerProvider().Tracer(""),
		engineFactory: func(ctx context.Context, wfid, _ string, _ types.WorkflowName, _ string, _, _ []byte, _ string, _ []byte, initDone chan<- error) (v2.WorkflowEngine, error) {
			f.factoryCalls.Add(1)
			if f.factoryErr != nil {
				return nil, f.factoryErr
			}
			initDone <- nil
			engine := f.legacyEngine
			if f.h.useCoordinatedEngine(ctx, wfid) {
				engine = f.coordinatedEngine
			}
			if cfg.sharded {
				return newTestShardFailoverManager(t, engine), nil
			}
			return engine, nil
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
	return newNamedTestWorkflowSpec(t, "engine-create-test")
}

// newNamedTestWorkflowSpec is newTestWorkflowSpec for a chosen name, so a test can
// run several workflows with distinct IDs.
func newNamedTestWorkflowSpec(t *testing.T, name string) (*job.WorkflowSpec, types.WorkflowID) {
	t.Helper()
	owner := []byte{0x01, 0x02, 0x03}
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

	t.Run("a sharded node takes the coordinated path with the flag open", func(t *testing.T) {
		t.Parallel()
		f := newEngineCreateFixture(t, nodeConfig{flagOpen: true, withCoordinator: true, sharded: true})
		spec, wid := newTestWorkflowSpec(t)

		require.NoError(t, f.h.tryEngineCreate(t.Context(), spec, "TestSource"))

		entry, ok := f.h.engineRegistry.Get(wid)
		require.True(t, ok)
		manager, ok := entry.Service.(*ShardFailoverManager)
		require.True(t, ok, "a sharded node registers the ShardFailoverManager, got %T", entry.Service)
		require.Same(t, f.coordinatedEngine, manager.engine, "the manager wraps the coordinated engine")
		require.True(t, entry.Coordinated())
		// The coordinator must subscribe and deliver through the manager, so
		// every event passes its shard admission check.
		require.Equal(t, []triggers.Subscriber{manager}, f.coordinator.registered)
		require.Equal(t, engineCounts{coordinated: 1}, countEngines(f.h.engineRegistry.GetAll()))
	})

	t.Run("a sharded node takes the legacy path with the flag closed", func(t *testing.T) {
		t.Parallel()
		f := newEngineCreateFixture(t, nodeConfig{flagOpen: false, withCoordinator: true, sharded: true})
		spec, wid := newTestWorkflowSpec(t)

		require.NoError(t, f.h.tryEngineCreate(t.Context(), spec, "TestSource"))

		entry, ok := f.h.engineRegistry.Get(wid)
		require.True(t, ok)
		manager, ok := entry.Service.(*ShardFailoverManager)
		require.True(t, ok, "a sharded node registers the ShardFailoverManager, got %T", entry.Service)
		require.Same(t, f.legacyEngine, manager.engine, "the manager wraps the legacy engine")
		require.False(t, entry.Coordinated())
		require.Empty(t, f.coordinator.registered)
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

	t.Run("sharded node with the flag open", func(t *testing.T) {
		t.Parallel()
		require.True(t, newHandler(true, true, true).useCoordinatedEngine(t.Context(), "wfid"))
	})

	t.Run("sharded node with the flag closed", func(t *testing.T) {
		t.Parallel()
		require.False(t, newHandler(false, true, true).useCoordinatedEngine(t.Context(), "wfid"))
	})

	t.Run("nil feature flags", func(t *testing.T) {
		t.Parallel()
		h := newHandler(true, true, false)
		h.featureFlags = nil
		require.False(t, h.useCoordinatedEngine(t.Context(), "wfid"))
	})

	t.Run("the deprecated global flag open overrides a closed per-workflow flag", func(t *testing.T) {
		t.Parallel()
		h := newHandler(false, true, false)
		h.featureFlags.CoordinatedEngineGlobal = limits.NewGateLimiter(true)
		require.True(t, h.useCoordinatedEngine(t.Context(), "wfid"))
	})
}

// Test_useCoordinatedEngine_Settings resolves the flags from real settings, as a node
// does: PerWorkflow.CoordinatedEngineEnabled can be set for one workflow, owner or org,
// the most specific scope wins, and the deprecated global key still enables every workflow.
func Test_useCoordinatedEngine_Settings(t *testing.T) {
	t.Parallel()

	const (
		orgOn        = "org-on"
		ownerOn      = "0000000000000000000000000000000000000a11"
		ownerDefault = "0000000000000000000000000000000000000b22"
		wfOn         = "1111111111111111111111111111111111111111111111111111111111111111"
		wfOff        = "2222222222222222222222222222222222222222222222222222222222222222"
		wfDefault    = "3333333333333333333333333333333333333333333333333333333333333333"
	)
	perWorkflow := func(enabled string) string {
		return `{"PerWorkflow":{"CoordinatedEngineEnabled":"` + enabled + `"}}`
	}
	settingsJSON := `{
		"global":` + perWorkflow("false") + `,
		"org":{"` + orgOn + `":` + perWorkflow("true") + `},
		"owner":{"` + ownerOn + `":` + perWorkflow("true") + `},
		"workflow":{"` + wfOn + `":` + perWorkflow("true") + `,"` + wfOff + `":` + perWorkflow("false") + `}
	}`

	newHandler := func(t *testing.T, settingsJSON string) *eventHandler {
		getter, err := settings.NewJSONGetter([]byte(settingsJSON))
		require.NoError(t, err)
		flags, err := v2.NewFeatureFlags(limits.Factory{Logger: logger.Test(t), Settings: getter}, nil)
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, flags.Close()) })
		return &eventHandler{lggr: logger.Test(t), featureFlags: flags, triggerCoordinator: &recordingCoordinator{}}
	}
	use := func(t *testing.T, h *eventHandler, cre contexts.CRE) bool {
		return h.useCoordinatedEngine(contexts.WithCRE(t.Context(), cre), cre.Workflow)
	}

	t.Run("most specific scope wins", func(t *testing.T) {
		t.Parallel()
		h := newHandler(t, settingsJSON)
		for _, tc := range []struct {
			name string
			cre  contexts.CRE
			want bool
		}{
			{"no override uses the global default", contexts.CRE{Org: "org-other", Owner: ownerDefault, Workflow: wfDefault}, false},
			{"workflow override", contexts.CRE{Org: "org-other", Owner: ownerDefault, Workflow: wfOn}, true},
			{"owner override", contexts.CRE{Org: "org-other", Owner: ownerOn, Workflow: wfDefault}, true},
			{"org override", contexts.CRE{Org: orgOn, Owner: ownerDefault, Workflow: wfDefault}, true},
			{"workflow override beats its org and owner", contexts.CRE{Org: orgOn, Owner: ownerOn, Workflow: wfOff}, false},
			{"missing org skips only the org scope", contexts.CRE{Owner: ownerOn, Workflow: wfDefault}, true},
		} {
			require.Equal(t, tc.want, use(t, h, tc.cre), tc.name)
		}
	})

	t.Run("missing workflow tenant takes the legacy path", func(t *testing.T) {
		t.Parallel()
		h := newHandler(t, `{"global":`+perWorkflow("true")+`}`)
		require.False(t, h.useCoordinatedEngine(t.Context(), wfDefault), "no tenant in ctx")
		require.True(t, use(t, h, contexts.CRE{Owner: ownerDefault, Workflow: wfDefault}), "same settings, with the tenant")
	})

	t.Run("deprecated global key still enables every workflow", func(t *testing.T) {
		t.Parallel()
		h := newHandler(t, `{"global":{"CoordinatedEngineEnabled":"true","PerWorkflow":{"CoordinatedEngineEnabled":"false"}},
			"workflow":{"`+wfOff+`":`+perWorkflow("false")+`}}`)
		require.True(t, use(t, h, contexts.CRE{Owner: ownerDefault, Workflow: wfDefault}))
		require.True(t, use(t, h, contexts.CRE{Owner: ownerDefault, Workflow: wfOff}),
			"the deprecated key is not overridable per workflow; move to PerWorkflow to opt one out")
	})
}

// switchableGetter is a settings.Getter whose JSON a test can replace, as an operator
// changes CRE settings on a running node.
type switchableGetter struct {
	mu     sync.Mutex
	getter settings.Getter
}

func newSwitchableGetter(t *testing.T, settingsJSON string) *switchableGetter {
	g := &switchableGetter{}
	g.set(t, settingsJSON)
	return g
}

func (g *switchableGetter) set(t *testing.T, settingsJSON string) {
	t.Helper()
	getter, err := settings.NewJSONGetter([]byte(settingsJSON))
	require.NoError(t, err)
	g.mu.Lock()
	defer g.mu.Unlock()
	g.getter = getter
}

func (g *switchableGetter) GetScoped(ctx context.Context, scope settings.Scope, key string) (string, error) {
	g.mu.Lock()
	getter := g.getter
	g.mu.Unlock()
	return getter.GetScoped(ctx, scope, key)
}

// unregisteringCoordinator is a recordingCoordinator that also accepts the
// UnregisterTriggers call engine cleanup makes for a coordinated workflow.
type unregisteringCoordinator struct {
	recordingCoordinator
	unregistered []string
}

func (c *unregisteringCoordinator) UnregisterTriggers(_ context.Context, workflowID string) error {
	c.unregistered = append(c.unregistered, workflowID)
	return nil
}

// Test_CoordinatedEngineFlag_PerWorkflowSwitch checks that the per-workflow
// CoordinatedEngineEnabled flag moves one workflow between the legacy and the
// coordinated engine while another workflow keeps running. The flag is read only when
// an engine is created, so a change applies when the workflow's engine is next
// recreated, and never touches the other workflow's engine.
func Test_CoordinatedEngineFlag_PerWorkflowSwitch(t *testing.T) {
	t.Parallel()

	f := newEngineCreateFixture(t, nodeConfig{withCoordinator: true})
	coordinator := &unregisteringCoordinator{}
	f.h.triggerCoordinator = coordinator

	const nameA, nameB = "flag-switch-a", "flag-switch-b"
	_, widA := newNamedTestWorkflowSpec(t, nameA)
	_, widB := newNamedTestWorkflowSpec(t, nameB)
	require.NotEqual(t, widA, widB)

	flagFor := func(workflowID, enabled string) string {
		return `{"global":{"PerWorkflow":{"CoordinatedEngineEnabled":"false"}},
			"workflow":{"` + workflowID + `":{"PerWorkflow":{"CoordinatedEngineEnabled":"` + enabled + `"}}}}`
	}
	getter := newSwitchableGetter(t, flagFor(widA.Hex(), "false"))
	flags, err := v2.NewFeatureFlags(limits.Factory{Logger: logger.Test(t), Settings: getter}, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, flags.Close()) })
	f.h.featureFlags = flags

	// A fresh engine per creation, so a test can tell one workflow's engine from another's.
	f.h.engineFactory = func(ctx context.Context, wfid, _ string, _ types.WorkflowName, _ string, _, _ []byte, _ string, _ []byte, initDone chan<- error) (v2.WorkflowEngine, error) {
		f.factoryCalls.Add(1)
		initDone <- nil
		if f.h.useCoordinatedEngine(ctx, wfid) {
			return &fakeCoordinatedEngine{}, nil
		}
		return &mockEngine{}, nil
	}

	// ctx carries the tenant, as Handle sets it before any engine is created.
	tenantCtx := func(spec *job.WorkflowSpec) context.Context {
		return contexts.WithCRE(t.Context(), contexts.CRE{Owner: spec.WorkflowOwner, Workflow: spec.WorkflowID})
	}
	// create loads a fresh spec each time, as the syncer does: engine creation consumes
	// the spec's binary.
	create := func(name string) {
		spec, _ := newNamedTestWorkflowSpec(t, name)
		require.NoError(t, f.h.tryEngineCreate(tenantCtx(spec), spec, "TestSource"))
	}
	// recreate is the replace path of workflowRegisteredEvent: clean up, then create.
	recreate := func(name string) {
		spec, wid := newNamedTestWorkflowSpec(t, name)
		require.NoError(t, f.h.tryEngineCleanup(tenantCtx(spec), wid))
		create(name)
	}
	entry := func(wid types.WorkflowID) ServiceWithMetadata {
		e, ok := f.h.engineRegistry.Get(wid)
		require.True(t, ok)
		return e
	}

	create(nameA)
	create(nameB)
	require.False(t, entry(widA).Coordinated(), "A starts on legacy")
	require.False(t, entry(widB).Coordinated(), "B starts on legacy")
	engineB := entry(widB).Service

	getter.set(t, flagFor(widA.Hex(), "true"))
	require.False(t, entry(widA).Coordinated(), "a running engine keeps its path until it is recreated")

	recreate(nameA)
	require.True(t, entry(widA).Coordinated(), "A is coordinated once recreated")
	require.Len(t, coordinator.registered, 1, "only A registers with the coordinator")
	require.Same(t, engineB, entry(widB).Service, "B's engine is untouched")
	require.False(t, entry(widB).Coordinated())

	getter.set(t, flagFor(widA.Hex(), "false"))
	coordinatedA := entry(widA).Service.(*fakeCoordinatedEngine)
	recreate(nameA)
	require.False(t, entry(widA).Coordinated(), "A is back on legacy once recreated")
	require.Equal(t, []string{widA.Hex()}, coordinator.unregistered, "the coordinated engine's triggers are unregistered")
	require.True(t, coordinatedA.closed.Load(), "the coordinated engine is closed")
	require.Same(t, engineB, entry(widB).Service, "B's engine is still untouched")

	require.Equal(t, engineCounts{legacy: 2}, countEngines(f.h.engineRegistry.GetAll()))
	require.Equal(t, int32(4), f.factoryCalls.Load(), "one creation per create or recreate")
}
