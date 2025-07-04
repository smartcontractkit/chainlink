package crecli

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDataSource_AddObjectAndGetData(t *testing.T) {
	ds := NewDataSource()

	// Handler that returns the param as bytes
	handler := func(ctx context.Context, param string) ([]byte, error) {
		return []byte("data:" + param), nil
	}

	ds.AddObject("obj1", handler, "param1")

	// Should retrieve correct data
	data, err := ds.GetData(context.Background(), "obj1")
	require.NoError(t, err)
	require.Equal(t, []byte("data:param1"), data)

	// Add another object with a different handler
	handler2 := func(ctx context.Context, param string) ([]byte, error) {
		return []byte(param + "-suffix"), nil
	}
	ds.AddObject("obj2", handler2, "foo")
	data2, err2 := ds.GetData(context.Background(), "obj2")
	require.NoError(t, err2)
	require.Equal(t, []byte("foo-suffix"), data2)
}

func TestDataSource_GetData_NoHandle(t *testing.T) {
	ds := NewDataSource()
	_, err := ds.GetData(context.Background(), "missing")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no handle found for object missing")
}

func TestDataSource_GetData_NoParam(t *testing.T) {
	ds := NewDataSource()
	ds.handles["obj"] = func(ctx context.Context, param string) ([]byte, error) {
		return nil, errors.New("should not be called")
	}
	// params map is empty, so param for "obj" is missing
	_, err := ds.GetData(context.Background(), "obj")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no parameter found for object obj")
}
