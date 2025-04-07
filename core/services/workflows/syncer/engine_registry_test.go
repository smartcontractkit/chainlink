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
	er := NewEngineRegistry()
	require.False(t, er.Contains(id1))

	e, err := er.Get(id1)
	require.ErrorIs(t, err, errNotFound)
	require.Nil(t, e)

	e, err = er.Pop(id1)
	require.ErrorIs(t, err, errNotFound)
	require.Nil(t, e)

	// add
	require.NoError(t, er.Add(id1, srv))
	require.True(t, er.Contains(id1))

	e, err = er.Get(id1)
	require.NoError(t, err)
	require.Equal(t, srv, e)

	// remove
	e, err = er.Pop(id1)
	require.NoError(t, err)
	require.Equal(t, srv, e)
	require.False(t, er.Contains(id1))

	// re-add
	require.NoError(t, er.Add(id1, srv))

	es := er.PopAll()
	require.Len(t, es, 1)
	require.Equal(t, srv, es[0])
}

type fakeService struct{}

func (f fakeService) Start(ctx context.Context) error { return nil }

func (f fakeService) Close() error { return nil }

func (f fakeService) Ready() error { return nil }

func (f fakeService) HealthReport() map[string]error { return map[string]error{} }

func (f fakeService) Name() string { return "" }
