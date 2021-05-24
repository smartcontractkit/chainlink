package pipeline_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
)

func TestVars_GetSet(t *testing.T) {
	t.Parallel()

	t.Run("gets the values at keypaths that exist", func(t *testing.T) {
		t.Parallel()

		vars := pipeline.Vars{
			"foo": []interface{}{1, "bar", false},
			"bar": 321,
		}

		got, err := vars.Get("foo.1")
		require.NoError(t, err)
		require.Equal(t, "bar", got)

		got, err = vars.Get("bar")
		require.NoError(t, err)
		require.Equal(t, 321, got)
	})

	t.Run("errors when getting the values at keypaths that don't exist", func(t *testing.T) {
		t.Parallel()

		vars := pipeline.Vars{
			"foo": []interface{}{1, "bar", false},
			"bar": 321,
		}

		_, err := vars.Get("foo.blah")
		require.Error(t, err)
	})

	t.Run("sets values at simple keypaths", func(t *testing.T) {
		t.Parallel()

		vars := pipeline.Vars{}

		err := vars.Set("foo", 123)
		require.NoError(t, err)

		err = vars.Set("bar", []interface{}{"a", "b"})
		require.NoError(t, err)

		expected := pipeline.Vars{
			"foo": 123,
			"bar": []interface{}{"a", "b"},
		}
		require.Equal(t, expected, vars)
	})

	t.Run("sets values in slices that exist", func(t *testing.T) {
		t.Parallel()

		vars := pipeline.Vars{
			"foo": []interface{}{1, "bar", false},
			"bar": 321,
		}
		err := vars.Set("foo.1", 123)
		require.NoError(t, err)

		expected := pipeline.Vars{
			"foo": []interface{}{1, 123, false},
			"bar": 321,
		}
		require.Equal(t, expected, vars)
	})

	t.Run("sets values in maps that exist", func(t *testing.T) {
		t.Parallel()

		vars := pipeline.Vars{
			"foo": map[string]interface{}{
				"bar": "hello",
			},
			"bar": 321,
		}
		err := vars.Set("foo.chain", "link")
		require.NoError(t, err)

		expected := pipeline.Vars{
			"foo": map[string]interface{}{
				"bar":   "hello",
				"chain": "link",
			},
			"bar": 321,
		}
		require.Equal(t, expected, vars)
	})

	t.Run("sets values in nested maps", func(t *testing.T) {
		t.Parallel()

		vars := pipeline.Vars{
			"foo": map[string]interface{}{
				"bar": map[string]interface{}{
					"chain": "link",
				},
			},
			"bar": 321,
		}
		err := vars.Set("foo.bar.sergey", 123)
		require.NoError(t, err)

		expected := pipeline.Vars{
			"foo": map[string]interface{}{
				"bar": map[string]interface{}{
					"chain":  "link",
					"sergey": 123,
				},
			},
			"bar": 321,
		}
		require.Equal(t, expected, vars)
	})

	t.Run("sets values in nested maps that don't exist by creating them", func(t *testing.T) {
		t.Parallel()

		vars := pipeline.Vars{
			"foo": map[string]interface{}{},
			"bar": 321,
		}
		err := vars.Set("foo.bar.sergey.chain", "link")
		require.NoError(t, err)

		expected := pipeline.Vars{
			"foo": map[string]interface{}{
				"bar": map[string]interface{}{
					"sergey": map[string]interface{}{
						"chain": "link",
					},
				},
			},
			"bar": 321,
		}
		require.Equal(t, expected, vars)
	})

	t.Run("errors when setting values in nested slices outside of their current size", func(t *testing.T) {
		t.Parallel()

		vars := pipeline.Vars{
			"foo": []interface{}{1, 2},
			"bar": 321,
		}
		err := vars.Set("foo.2", "link")
		require.Error(t, err)
	})
}

func TestVars_ResolveValue(t *testing.T) {
	t.Parallel()

	t.Run("calls getters in order until the first one that returns without ErrParameterEmpty", func(t *testing.T) {
		t.Parallel()

		vars := pipeline.Vars{}

		param := new(mocks.PipelineParamUnmarshaler)
		param.On("UnmarshalPipelineParam", mock.Anything, vars).Return(nil)

		called := []int{}
		getters := []pipeline.GetterFunc{
			func(_ pipeline.Vars) (interface{}, error) {
				called = append(called, 0)
				return nil, errors.Wrap(pipeline.ErrParameterEmpty, "make sure it still notices when wrapped")
			},
			func(_ pipeline.Vars) (interface{}, error) {
				called = append(called, 1)
				return 123, nil
			},
			func(_ pipeline.Vars) (interface{}, error) {
				called = append(called, 2)
				return 123, nil
			},
		}

		err := vars.ResolveValue(param, getters)
		require.NoError(t, err)
		require.Equal(t, []int{0, 1}, called)

		param.AssertExpectations(t)
	})

	t.Run("returns any GetterFunc error that isn't ErrParameterEmpty", func(t *testing.T) {
		t.Parallel()

		vars := pipeline.Vars{}
		param := new(mocks.PipelineParamUnmarshaler)
		called := []int{}
		expectedErr := errors.New("some other issue")

		getters := []pipeline.GetterFunc{
			func(_ pipeline.Vars) (interface{}, error) {
				called = append(called, 0)
				return nil, expectedErr
			},
			func(_ pipeline.Vars) (interface{}, error) {
				called = append(called, 1)
				return 123, nil
			},
			func(_ pipeline.Vars) (interface{}, error) {
				called = append(called, 2)
				return 123, nil
			},
		}

		err := vars.ResolveValue(param, getters)
		require.Equal(t, expectedErr, err)
		require.Equal(t, []int{0}, called)
	})

	t.Run("calls UnmarshalPipelineParam with the value obtained from the GetterFuncs", func(t *testing.T) {
		t.Parallel()

		vars := pipeline.Vars{}
		expectedValue := 123

		param := new(mocks.PipelineParamUnmarshaler)
		param.On("UnmarshalPipelineParam", expectedValue, vars).Return(nil)

		getters := []pipeline.GetterFunc{
			func(_ pipeline.Vars) (interface{}, error) {
				return expectedValue, nil
			},
		}

		err := vars.ResolveValue(param, getters)
		require.NoError(t, err)

		param.AssertExpectations(t)
	})

	t.Run("returns any error returned by UnmarshalPipelineParam", func(t *testing.T) {
		t.Parallel()

		vars := pipeline.Vars{}
		expectedValue := 123
		expectedErr := errors.New("some issue")

		param := new(mocks.PipelineParamUnmarshaler)
		param.On("UnmarshalPipelineParam", expectedValue, vars).Return(expectedErr)

		getters := []pipeline.GetterFunc{
			func(_ pipeline.Vars) (interface{}, error) {
				return expectedValue, nil
			},
		}

		err := vars.ResolveValue(param, getters)
		require.Equal(t, expectedErr, err)

		param.AssertExpectations(t)
	})
}

func TestGetters_VariableExpr(t *testing.T) {
	t.Parallel()

	vars := pipeline.Vars{
		"foo": map[string]interface{}{
			"bar": []interface{}{0, 1, 42},
		},
	}

	tests := []struct {
		expr    string
		keypath interface{}
		err     error
	}{
		{"$(foo.bar.2)", 42, nil},
		{" $(foo.bar.2)", 42, nil},
		{"$(foo.bar.2) ", 42, nil},
		{"$( foo.bar.2)", 42, nil},
		{"$(foo.bar.2 )", 42, nil},
		{"$( foo.bar.2 )", 42, nil},
		{" $( foo.bar.2 )", 42, nil},
		{"$()", (map[string]interface{})(vars), nil},
		{"$(foo.bar.2", nil, pipeline.ErrParameterEmpty},
		{"$foo.bar.2)", nil, pipeline.ErrParameterEmpty},
		{"(foo.bar.2)", nil, pipeline.ErrParameterEmpty},
		{"foo.bar.2", nil, pipeline.ErrParameterEmpty},
	}

	for _, test := range tests {
		test := test
		t.Run(test.expr, func(t *testing.T) {
			val, err := pipeline.VariableExpr(test.expr)(vars)
			require.Equal(t, test.keypath, val)
			require.Equal(t, test.err, errors.Cause(err))
		})
	}
}

func TestGetters_NonemptyString(t *testing.T) {
	t.Parallel()

	t.Run("returns any non-empty string", func(t *testing.T) {
		t.Parallel()
		val, err := pipeline.NonemptyString("foo bar")(nil)
		require.NoError(t, err)
		require.Equal(t, "foo bar", val)
	})

	t.Run("returns ErrParameterEmpty when given an empty string (including only spaces)", func(t *testing.T) {
		t.Parallel()
		_, err := pipeline.NonemptyString("")(nil)
		require.Equal(t, pipeline.ErrParameterEmpty, errors.Cause(err))
		_, err = pipeline.NonemptyString(" ")(nil)
		require.Equal(t, pipeline.ErrParameterEmpty, errors.Cause(err))
	})
}

func TestGetters_Input(t *testing.T) {
	t.Parallel()

	t.Run("returns the requested input's Value and Error if they exist", func(t *testing.T) {
		t.Parallel()
		expectedVal := "bar"
		expectedErr := errors.New("some err")
		val, err := pipeline.Input([]pipeline.Result{{Value: "foo"}, {Value: expectedVal, Error: expectedErr}, {Value: "baz"}}, 1)(nil)
		require.Equal(t, expectedVal, val)
		require.Equal(t, expectedErr, err)
	})

	t.Run("returns ErrParameterEmpty if the specified input does not exist", func(t *testing.T) {
		t.Parallel()
		_, err := pipeline.Input([]pipeline.Result{{Value: "foo"}}, 1)(nil)
		require.Equal(t, pipeline.ErrParameterEmpty, errors.Cause(err))
	})
}

func TestGetters_Inputs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		inputs      []pipeline.Result
		expected    []interface{}
		expectedErr error
	}{
		{
			"returns the values and errors",
			[]pipeline.Result{
				{Value: "foo"},
				{Error: errors.New("some issue")},
				{Value: "baz"},
			},
			[]interface{}{"foo", errors.New("some issue"), "baz"}, nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run("returns all of the inputs' Values if the provided inputs meet the spec", func(t *testing.T) {
			t.Parallel()

			val, err := pipeline.Inputs(test.inputs)(nil)
			require.Equal(t, test.expectedErr, errors.Cause(err))
			require.Equal(t, test.expected, val)
		})
	}
}
