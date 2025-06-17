package syncer

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

func TestEngineRegistry(t *testing.T) {
	var srv services.Service = &fakeService{}

	const id1 = "foo"
	owner := []byte{1, 2, 3, 4, 5}
	name := "my-workflow"
	workflowID := types.WorkflowID([32]byte{0, 1, 2, 3, 4})
	er := NewEngineRegistry()
	require.False(t, er.Contains(EngineRegistryKey{Owner: owner, Name: name}))

	e, ok := er.Get(EngineRegistryKey{Owner: owner, Name: name})
	require.False(t, ok)
	require.Nil(t, e.Service)
	require.Equal(t, ServiceWithMetadata{}, e)

	e, err := er.Pop(EngineRegistryKey{Owner: owner, Name: name})
	require.ErrorIs(t, err, errNotFound)
	require.Nil(t, e.Service)
	require.Equal(t, ServiceWithMetadata{}, e)

	// add
	require.NoError(t, er.Add(EngineRegistryKey{Owner: owner, Name: name}, srv, workflowID))
	require.True(t, er.Contains(EngineRegistryKey{Owner: owner, Name: name}))

	// add another item
	// this verifies that keys are unique
	name2 := "my-workflow-2"
	require.NoError(t, er.Add(EngineRegistryKey{Owner: owner, Name: name2}, srv, workflowID))
	require.True(t, er.Contains(EngineRegistryKey{Owner: owner, Name: name}))

	// get
	e, ok = er.Get(EngineRegistryKey{Owner: owner, Name: name})
	require.True(t, ok)
	require.Equal(t, srv, e.Service)
	require.Equal(t, workflowID, e.WorkflowID)
	require.Equal(t, owner, e.WorkflowOwner)
	require.Equal(t, name, e.WorkflowName)

	// get all
	es := er.GetAll()
	require.Len(t, es, 2)

	// remove
	e, err = er.Pop(EngineRegistryKey{Owner: owner, Name: name})
	require.NoError(t, err)
	require.Equal(t, srv, e.Service)
	require.Equal(t, workflowID, e.WorkflowID)
	require.Equal(t, owner, e.WorkflowOwner)
	require.Equal(t, name, e.WorkflowName)
	require.False(t, er.Contains(EngineRegistryKey{Owner: owner, Name: name}))

	// re-add
	require.NoError(t, er.Add(EngineRegistryKey{Owner: owner, Name: name}, srv, workflowID))

	// pop all
	es = er.PopAll()
	require.Len(t, es, 2)
}

func TestEngineRegistry_keyFor(t *testing.T) {
	owner := []byte("owner")
	k := EngineRegistryKey{Owner: owner, Name: "name"}
	assert.Equal(t, k.keyFor(), fmt.Sprintf("%x-name", owner))

	k = EngineRegistryKey{Owner: owner, Name: "name2"}
	assert.Equal(t, k.keyFor(), fmt.Sprintf("%x-name2", owner))
}

type fakeService struct{}

func (f fakeService) Start(ctx context.Context) error { return nil }

func (f fakeService) Close() error { return nil }

func (f fakeService) Ready() error { return nil }

func (f fakeService) HealthReport() map[string]error { return map[string]error{} }

func (f fakeService) Name() string { return "" }
