package v2

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

func TestEngineRegistry(t *testing.T) {
	t.Parallel()
	workflowID1 := types.WorkflowID([32]byte{0, 1, 2, 3, 4})
	workflowID2 := types.WorkflowID([32]byte{0, 1, 2, 3, 4, 5})

	var srv services.Service = &fakeService{}

	er := NewEngineRegistry()
	ok := er.Contains(workflowID1)
	require.False(t, ok)

	e, ok := er.Get(workflowID1)
	require.False(t, ok)
	require.Nil(t, e.Service)

	e, err := er.Pop(workflowID1)
	require.ErrorIs(t, err, ErrNotFound)
	require.Nil(t, e.Service)

	// add
	require.NoError(t, er.Add(workflowID1, "TestSource", srv))
	ok = er.Contains(workflowID1)
	require.True(t, ok)

	// add another item
	// this verifies that keys are unique
	require.NoError(t, er.Add(workflowID2, "TestSource", srv))
	ok = er.Contains(workflowID2)
	require.True(t, ok)

	// get
	e, ok = er.Get(workflowID1)
	require.True(t, ok)
	require.Equal(t, srv, e.Service)

	// get all
	es := er.GetAll()
	require.Len(t, es, 2)

	// remove
	e, err = er.Pop(workflowID1)
	require.NoError(t, err)
	require.Equal(t, srv, e.Service)
	ok = er.Contains(workflowID1)
	require.False(t, ok)

	// re-add
	require.NoError(t, er.Add(workflowID1, "TestSource", srv))

	// pop all
	es = er.PopAll()
	require.Len(t, es, 2)
}

func TestEngineRegistry_SourceTracking(t *testing.T) {
	t.Parallel()
	er := NewEngineRegistry()

	wfID1 := types.WorkflowID([32]byte{1})
	wfID2 := types.WorkflowID([32]byte{2})
	wfID3 := types.WorkflowID([32]byte{3})

	// Add engines from different sources
	require.NoError(t, er.Add(wfID1, ContractWorkflowSourceName, &fakeService{}))
	require.NoError(t, er.Add(wfID2, ContractWorkflowSourceName, &fakeService{}))
	require.NoError(t, er.Add(wfID3, GRPCWorkflowSourceName, &fakeService{}))

	// GetBySource filters correctly
	contractEngines := er.GetBySource(ContractWorkflowSourceName)
	require.Len(t, contractEngines, 2)

	grpcEngines := er.GetBySource(GRPCWorkflowSourceName)
	require.Len(t, grpcEngines, 1)

	// Unknown source returns empty
	unknownEngines := er.GetBySource("UnknownSource")
	require.Empty(t, unknownEngines)
}

func TestEngineRegistry_SourceInMetadata(t *testing.T) {
	t.Parallel()
	er := NewEngineRegistry()
	wfID := types.WorkflowID([32]byte{1})

	require.NoError(t, er.Add(wfID, "TestSource", &fakeService{}))

	engine, ok := er.Get(wfID)
	require.True(t, ok)
	require.Equal(t, "TestSource", engine.Source)
}

func TestEngineRegistry_GetAllIncludesSource(t *testing.T) {
	t.Parallel()
	er := NewEngineRegistry()

	wfID1 := types.WorkflowID([32]byte{1})
	wfID2 := types.WorkflowID([32]byte{2})

	require.NoError(t, er.Add(wfID1, ContractWorkflowSourceName, &fakeService{}))
	require.NoError(t, er.Add(wfID2, GRPCWorkflowSourceName, &fakeService{}))

	engines := er.GetAll()
	require.Len(t, engines, 2)

	// Verify each engine has its source
	sources := make(map[string]bool)
	for _, e := range engines {
		sources[e.Source] = true
	}
	require.True(t, sources[ContractWorkflowSourceName])
	require.True(t, sources[GRPCWorkflowSourceName])
}

func TestEngineRegistry_PopReturnsSource(t *testing.T) {
	t.Parallel()
	er := NewEngineRegistry()
	wfID := types.WorkflowID([32]byte{1})

	require.NoError(t, er.Add(wfID, ContractWorkflowSourceName, &fakeService{}))

	engine, err := er.Pop(wfID)
	require.NoError(t, err)
	require.Equal(t, ContractWorkflowSourceName, engine.Source)
}

func TestEngineRegistry_PopAllReturnsSource(t *testing.T) {
	t.Parallel()
	er := NewEngineRegistry()

	wfID1 := types.WorkflowID([32]byte{1})
	wfID2 := types.WorkflowID([32]byte{2})

	require.NoError(t, er.Add(wfID1, ContractWorkflowSourceName, &fakeService{}))
	require.NoError(t, er.Add(wfID2, GRPCWorkflowSourceName, &fakeService{}))

	engines := er.PopAll()
	require.Len(t, engines, 2)

	// Verify sources are preserved
	sources := make(map[string]bool)
	for _, e := range engines {
		sources[e.Source] = true
	}
	require.True(t, sources[ContractWorkflowSourceName])
	require.True(t, sources[GRPCWorkflowSourceName])
}

type fakeService struct{}

func (f fakeService) Start(ctx context.Context) error { return nil }

func (f fakeService) Close() error { return nil }

func (f fakeService) Ready() error { return nil }

func (f fakeService) HealthReport() map[string]error { return map[string]error{} }

func (f fakeService) Name() string { return "" }

// fakeDrainableService reports a fixed drain state; it is otherwise a fakeService.
type fakeDrainableService struct {
	fakeService
	draining bool
}

func (f *fakeDrainableService) Drain() bool { return false }

func (f *fakeDrainableService) ActiveExecutions() int32 { return 0 }

func (f *fakeDrainableService) DrainStartedAt() (time.Time, bool) {
	if !f.draining {
		return time.Time{}, false
	}
	return time.Unix(1, 0), true
}

// fakeCoordinatedDrainableEngine is a v2.WorkflowEngine that reports as
// coordinated and drainable. Only the methods ServiceWithMetadata.Coordinated
// and the DrainableService assertion touch are implemented.
type fakeCoordinatedDrainableEngine struct {
	v2.WorkflowEngine
	draining bool
}

func (f *fakeCoordinatedDrainableEngine) IsCoordinated() bool { return true }

func (f *fakeCoordinatedDrainableEngine) Drain() bool { return false }

func (f *fakeCoordinatedDrainableEngine) ActiveExecutions() int32 { return 0 }

func (f *fakeCoordinatedDrainableEngine) DrainStartedAt() (time.Time, bool) {
	if !f.draining {
		return time.Time{}, false
	}
	return time.Unix(1, 0), true
}

func TestEngineRegistry_Coordinated(t *testing.T) {
	t.Parallel()
	legacyID := types.WorkflowID([32]byte{1})
	coordinatedID := types.WorkflowID([32]byte{2})

	er := NewEngineRegistry()
	require.NoError(t, er.AddWithReconcileKey(legacyID, "TestSource", "legacy-key", &fakeService{}))
	require.NoError(t, er.AddWithReconcileKey(coordinatedID, "TestSource", "coordinated-key", &fakeCoordinatedDrainableEngine{}))
	require.ErrorIs(t, er.AddWithReconcileKey(coordinatedID, "TestSource", "coordinated-key", &fakeCoordinatedDrainableEngine{}), ErrAlreadyExists)

	entry, ok := er.Get(coordinatedID)
	require.True(t, ok)
	require.True(t, entry.Coordinated())
	require.Equal(t, "coordinated-key", entry.ReconcileKey)

	entry, ok = er.Get(legacyID)
	require.True(t, ok)
	require.False(t, entry.Coordinated())

	byID := map[types.WorkflowID]bool{}
	for _, e := range er.GetAll() {
		byID[e.WorkflowID] = e.Coordinated()
	}
	require.Equal(t, map[types.WorkflowID]bool{legacyID: false, coordinatedID: true}, byID)

	for _, e := range er.GetBySource("TestSource") {
		require.Equal(t, e.WorkflowID == coordinatedID, e.Coordinated())
	}

	popped, err := er.Pop(coordinatedID)
	require.NoError(t, err)
	require.True(t, popped.Coordinated())

	for _, e := range er.PopAll() {
		require.False(t, e.Coordinated())
	}
}

func TestCountEngines(t *testing.T) {
	t.Parallel()

	t.Run("empty registry", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, engineCounts{}, countEngines(nil))
	})

	t.Run("mixed engine types", func(t *testing.T) {
		t.Parallel()
		engines := []ServiceWithMetadata{
			{Service: &fakeCoordinatedDrainableEngine{}},
			{Service: &fakeCoordinatedDrainableEngine{draining: true}},
			{Service: &fakeDrainableService{}},
			{Service: &fakeDrainableService{draining: true}},
			{Service: &fakeService{}},
			{Service: newTestShardFailoverManager(t, &mockEngine{})},
		}
		require.Equal(t, engineCounts{draining: 2, coordinated: 2, legacy: 4}, countEngines(engines))
	})

	t.Run("a sharded wrapper follows the engine it wraps", func(t *testing.T) {
		t.Parallel()
		engines := []ServiceWithMetadata{
			{Service: newTestShardFailoverManager(t, &fakeCoordinatedDrainableEngine{})},
		}
		require.Equal(t, engineCounts{coordinated: 1}, countEngines(engines))
	})
}

func TestTriggerEngineRegistry(t *testing.T) {
	t.Parallel()
	wfID := types.WorkflowID([32]byte{9})

	tests := []struct {
		name            string
		registered      services.Service // nil registers nothing
		wantOK          bool
		wantCoordinated bool
	}{
		{
			name:       "unregistered workflow",
			registered: nil,
			wantOK:     false,
		},
		{
			name:       "service that is not a RegisteredEngine",
			registered: &fakeService{},
			wantOK:     false,
		},
		{
			name:            "coordinated engine",
			registered:      &fakeCoordinatedDrainableEngine{},
			wantOK:          true,
			wantCoordinated: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			er := NewEngineRegistry()
			if tt.registered != nil {
				require.NoError(t, er.Add(wfID, "src", tt.registered))
			}

			got, ok := NewTriggerEngineRegistry(er).Get(wfID)
			require.Equal(t, tt.wantOK, ok)
			if !tt.wantOK {
				require.Nil(t, got)
				return
			}
			require.Equal(t, tt.wantCoordinated, got.IsCoordinated())
		})
	}
}
