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
	require.NoError(t, er.Add(workflowID1, srv))
	ok = er.Contains(workflowID1)
	require.True(t, ok)

	// add another item
	// this verifies that keys are unique
	require.NoError(t, er.Add(workflowID2, srv))
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
	require.NoError(t, er.Add(workflowID1, srv))

	// pop all
	es = er.PopAll()
	require.Len(t, es, 2)
}

type fakeService struct{}

func (f fakeService) Start(ctx context.Context) error { return nil }

func (f fakeService) Close() error { return nil }

func (f fakeService) Ready() error { return nil }

func (f fakeService) HealthReport() map[string]error { return map[string]error{} }

func (f fakeService) Name() string { return "" }
