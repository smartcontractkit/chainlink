package pipeline_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestVars_Set(t *testing.T) {
	t.Parallel()

	vars := pipeline.NewVarsFrom(nil)

	err := vars.Set("xyz", "foo")
	require.NoError(t, err)
	v, err := vars.Get("xyz")
	require.NoError(t, err)
	require.Equal(t, "foo", v)

	err = vars.Set("  ", "foo")
	require.ErrorIs(t, err, pipeline.ErrVarsRoot)

	err = vars.Set("x.y", "foo")
	require.ErrorIs(t, err, pipeline.ErrVarsSetNested)
}

func TestVars_Get(t *testing.T) {
	t.Parallel()

	t.Run("gets the values at keypaths that exist", func(t *testing.T) {
		t.Parallel()

		vars := pipeline.NewVarsFrom(map[string]interface{}{
			"foo": []interface{}{1, "bar", false},
			"bar": 321,
		})

		got, err := vars.Get("foo.1")
		require.NoError(t, err)
		require.Equal(t, "bar", got)

		got, err = vars.Get("bar")
		require.NoError(t, err)
		require.Equal(t, 321, got)
	})

	t.Run("errors when getting the values at keypaths that don't exist", func(t *testing.T) {
		t.Parallel()

		vars := pipeline.NewVarsFrom(map[string]interface{}{
			"foo": []interface{}{1, "bar", false},
			"bar": 321,
		})

		_, err := vars.Get("foo.blah")
		require.Equal(t, pipeline.ErrKeypathNotFound, errors.Cause(err))
	})

	t.Run("errors when asked for a keypath with more than 2 parts", func(t *testing.T) {
		t.Parallel()

		vars := pipeline.NewVarsFrom(map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": map[string]interface{}{
					"chainlink": 123,
				},
			},
		})
		_, err := vars.Get("foo.bar.chainlink")
		require.Equal(t, pipeline.ErrWrongKeypath, errors.Cause(err))
	})

	t.Run("errors when getting a value at a keypath where the first part is not a map/slice", func(t *testing.T) {
		t.Parallel()

		vars := pipeline.NewVarsFrom(map[string]interface{}{
			"foo": 123,
		})
		_, err := vars.Get("foo.bar")
		require.Equal(t, pipeline.ErrKeypathNotFound, errors.Cause(err))
	})

	t.Run("errors when getting a value at a keypath with more than 2 components", func(t *testing.T) {
		t.Parallel()

		vars := pipeline.NewVarsFrom(map[string]interface{}{
			"foo": 123,
		})
		_, err := vars.Get("foo.bar.baz")
		require.Equal(t, pipeline.ErrWrongKeypath, errors.Cause(err))
	})

	t.Run("index out of range", func(t *testing.T) {
		t.Parallel()

		vars := pipeline.NewVarsFrom(map[string]interface{}{
			"foo": []interface{}{1, "bar", false},
		})

		_, err := vars.Get("foo.4")
		require.ErrorIs(t, err, pipeline.ErrIndexOutOfRange)

		_, err = vars.Get("foo.-1")
		require.ErrorIs(t, err, pipeline.ErrIndexOutOfRange)
	})
}

func TestVars_Copy(t *testing.T) {
	t.Parallel()

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"nested": map[string]interface{}{
			"foo": "zet",
		},
		"bar": 321,
	})

	varsCopy := vars.Copy()
	require.Equal(t, vars, varsCopy)
}
