package syncer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

func TestEngineRegistry(t *testing.T) {
	var srv services.Service = &fakeService{}

	const id1 = "foo"
	owner := []byte{1, 2, 3, 4, 5}
	name := "my-workflow"
	workflowID := [32]byte{0, 1, 2, 3, 4}
	er := NewEngineRegistry()
	require.False(t, er.Contains(EngineRegistryKey{Owner: owner, Name: name}))

	e, err := er.Get(EngineRegistryKey{Owner: owner, Name: name})
	require.ErrorIs(t, err, errNotFound)
	require.Nil(t, e.Service)
	require.Equal(t, ServiceWithMetadata{}, e)

	e, err = er.Pop(EngineRegistryKey{Owner: owner, Name: name})
	require.ErrorIs(t, err, errNotFound)
	require.Nil(t, e.Service)
	require.Equal(t, ServiceWithMetadata{}, e)

	// add
	require.NoError(t, er.Add(EngineRegistryKey{Owner: owner, Name: name}, srv, workflowID))
	require.True(t, er.Contains(EngineRegistryKey{Owner: owner, Name: name}))

	// get
	e, err = er.Get(EngineRegistryKey{Owner: owner, Name: name})
	require.NoError(t, err)
	require.Equal(t, srv, e.Service)
	require.Equal(t, workflowID, e.workflowID)
	require.Equal(t, owner, e.workflowOwner)
	require.Equal(t, name, e.workflowName)

	// get all
	es := er.GetAll()
	require.Len(t, es, 1)
	require.Equal(t, srv, es[0].Service)
	require.Equal(t, es[0].workflowID, e.workflowID)
	require.Equal(t, es[0].workflowOwner, e.workflowOwner)
	require.Equal(t, es[0].workflowName, e.workflowName)

	// remove
	e, err = er.Pop(EngineRegistryKey{Owner: owner, Name: name})
	require.NoError(t, err)
	require.Equal(t, srv, e.Service)
	require.Equal(t, workflowID, e.workflowID)
	require.Equal(t, owner, e.workflowOwner)
	require.Equal(t, name, e.workflowName)
	require.False(t, er.Contains(EngineRegistryKey{Owner: owner, Name: name}))

	// re-add
	require.NoError(t, er.Add(EngineRegistryKey{Owner: owner, Name: name}, srv, workflowID))

	// pop all
	es = er.PopAll()
	require.Len(t, es, 1)
	require.Equal(t, srv, es[0].Service)
	require.Equal(t, es[0].workflowID, e.workflowID)
	require.Equal(t, es[0].workflowOwner, e.workflowOwner)
	require.Equal(t, es[0].workflowName, e.workflowName)
}

type fakeService struct{}

func (f fakeService) Start(ctx context.Context) error { return nil }

func (f fakeService) Close() error { return nil }

func (f fakeService) Ready() error { return nil }

func (f fakeService) HealthReport() map[string]error { return map[string]error{} }

func (f fakeService) Name() string { return "" }
