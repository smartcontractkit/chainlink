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
	"github.com/smartcontractkit/chainlink-common/pkg/services"
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

// engineCreateFixture is an eventHandler with counting fakes for both engine
// factories, so a test can tell which creation path tryEngineCreate took.
type engineCreateFixture struct {
	h                 *eventHandler
	legacyEngine      *mockEngine
	coordinatedEngine *fakeCoordinatedEngine
	legacyCalls       atomic.Int32
	coordinatedCalls  atomic.Int32
	coordinator       *recordingCoordinator
}

func newEngineCreateFixture(t *testing.T, withCoordinator, flagOpen bool) *engineCreateFixture {
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
		featureFlags:   &v2.EngineFeatureFlags{CoordinatedEngine: limits.NewGateLimiter(flagOpen)},
		tracer:         noop.NewTracerProvider().Tracer(""),
		legacyEngineFactory: func(_ context.Context, _, _ string, _ types.WorkflowName, _ string, _, _ []byte, _ string, initDone chan<- error) (services.Service, error) {
			f.legacyCalls.Add(1)
			initDone <- nil
			return f.legacyEngine, nil
		},
		coordinatedEngineFactory: func(_ context.Context, _, _ string, _ types.WorkflowName, _ string, _, _ []byte, _ string, initDone chan<- error) (v2.WorkflowEngine, error) {
			f.coordinatedCalls.Add(1)
			initDone <- nil
			return f.coordinatedEngine, nil
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

	t.Run("flag open with a coordinator takes the coordinated path", func(t *testing.T) {
		t.Parallel()
		f := newEngineCreateFixture(t, true, true)
		spec, wid := newTestWorkflowSpec(t)

		require.NoError(t, f.h.tryEngineCreate(t.Context(), spec, "TestSource"))

		require.Equal(t, int32(1), f.coordinatedCalls.Load())
		require.Equal(t, int32(0), f.legacyCalls.Load())
		entry, ok := f.h.engineRegistry.Get(wid)
		require.True(t, ok)
		require.True(t, entry.Coordinated)
		require.NotEmpty(t, entry.ReconcileKey)
		require.Equal(t, []v2.Subscriber{f.coordinatedEngine}, f.coordinator.registered)
	})

	t.Run("flag closed takes the legacy path", func(t *testing.T) {
		t.Parallel()
		f := newEngineCreateFixture(t, true, false)
		spec, wid := newTestWorkflowSpec(t)

		require.NoError(t, f.h.tryEngineCreate(t.Context(), spec, "TestSource"))

		require.Equal(t, int32(0), f.coordinatedCalls.Load())
		require.Equal(t, int32(1), f.legacyCalls.Load())
		entry, ok := f.h.engineRegistry.Get(wid)
		require.True(t, ok)
		require.False(t, entry.Coordinated)
		require.Empty(t, f.coordinator.registered)
	})

	t.Run("flag open without a coordinator takes the legacy path", func(t *testing.T) {
		t.Parallel()
		f := newEngineCreateFixture(t, false, true)
		spec, wid := newTestWorkflowSpec(t)

		require.NoError(t, f.h.tryEngineCreate(t.Context(), spec, "TestSource"))

		require.Equal(t, int32(0), f.coordinatedCalls.Load())
		require.Equal(t, int32(1), f.legacyCalls.Load())
		entry, ok := f.h.engineRegistry.Get(wid)
		require.True(t, ok)
		require.False(t, entry.Coordinated)
	})

	t.Run("RegisterTriggers failure removes and closes the coordinated engine", func(t *testing.T) {
		t.Parallel()
		f := newEngineCreateFixture(t, true, true)
		f.coordinator.registerErr = errors.New("boom")
		spec, wid := newTestWorkflowSpec(t)

		err := f.h.tryEngineCreate(t.Context(), spec, "TestSource")

		require.ErrorContains(t, err, "failed to register triggers via coordinator")
		require.False(t, f.h.engineRegistry.Contains(wid))
		require.True(t, f.coordinatedEngine.closed.Load())
	})

	t.Run("invalid spec is rejected before any engine is built", func(t *testing.T) {
		t.Parallel()
		f := newEngineCreateFixture(t, true, true)
		spec, _ := newTestWorkflowSpec(t)
		spec.WorkflowID = types.WorkflowID{0xff}.Hex()

		err := f.h.tryEngineCreate(t.Context(), spec, "TestSource")

		require.ErrorContains(t, err, "workflowID mismatch")
		require.Equal(t, int32(0), f.coordinatedCalls.Load())
		require.Equal(t, int32(0), f.legacyCalls.Load())
	})
}
