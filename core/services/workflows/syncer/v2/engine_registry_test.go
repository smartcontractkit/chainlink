package v2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

func TestEngineRegistry(t *testing.T) {
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
	er := NewEngineRegistry()
	wfID := types.WorkflowID([32]byte{1})

	require.NoError(t, er.Add(wfID, "TestSource", &fakeService{}))

	engine, ok := er.Get(wfID)
	require.True(t, ok)
	require.Equal(t, "TestSource", engine.Source)
}

func TestEngineRegistry_GetAllIncludesSource(t *testing.T) {
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
	er := NewEngineRegistry()
	wfID := types.WorkflowID([32]byte{1})

	require.NoError(t, er.Add(wfID, ContractWorkflowSourceName, &fakeService{}))

	engine, err := er.Pop(wfID)
	require.NoError(t, err)
	require.Equal(t, ContractWorkflowSourceName, engine.Source)
}

func TestEngineRegistry_PopAllReturnsSource(t *testing.T) {
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
