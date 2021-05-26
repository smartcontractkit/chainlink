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
}

func TestVars_ResolveValue(t *testing.T) {
	t.Parallel()

	t.Run("calls getters in order until the first one that returns without ErrParameterEmpty", func(t *testing.T) {
		t.Parallel()

		param := new(mocks.PipelineParamUnmarshaler)
		param.On("UnmarshalPipelineParam", mock.Anything).Return(nil)

		called := []int{}
		getters := []pipeline.GetterFunc{
			func() (interface{}, error) {
				called = append(called, 0)
				return nil, errors.Wrap(pipeline.ErrParameterEmpty, "make sure it still notices when wrapped")
			},
			func() (interface{}, error) {
				called = append(called, 1)
				return 123, nil
			},
			func() (interface{}, error) {
				called = append(called, 2)
				return 123, nil
			},
		}

		err := pipeline.ResolveParam(param, getters)
		require.NoError(t, err)
		require.Equal(t, []int{0, 1}, called)

		param.AssertExpectations(t)
	})

	t.Run("returns any GetterFunc error that isn't ErrParameterEmpty", func(t *testing.T) {
		t.Parallel()

		param := new(mocks.PipelineParamUnmarshaler)
		called := []int{}
		expectedErr := errors.New("some other issue")

		getters := []pipeline.GetterFunc{
			func() (interface{}, error) {
				called = append(called, 0)
				return nil, expectedErr
			},
			func() (interface{}, error) {
				called = append(called, 1)
				return 123, nil
			},
			func() (interface{}, error) {
				called = append(called, 2)
				return 123, nil
			},
		}

		err := pipeline.ResolveParam(param, getters)
		require.Equal(t, expectedErr, err)
		require.Equal(t, []int{0}, called)
	})

	t.Run("calls UnmarshalPipelineParam with the value obtained from the GetterFuncs", func(t *testing.T) {
		t.Parallel()

		expectedValue := 123

		param := new(mocks.PipelineParamUnmarshaler)
		param.On("UnmarshalPipelineParam", expectedValue).Return(nil)

		getters := []pipeline.GetterFunc{
			func() (interface{}, error) {
				return expectedValue, nil
			},
		}

		err := pipeline.ResolveParam(param, getters)
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
			func() (interface{}, error) {
				return expectedValue, nil
			},
		}

		err := pipeline.ResolveParam(param, getters)
		require.Equal(t, expectedErr, err)

		param.AssertExpectations(t)
	})
}

func TestGetters_NonemptyString(t *testing.T) {
	t.Parallel()

	t.Run("returns any non-empty string", func(t *testing.T) {
		t.Parallel()
		val, err := pipeline.NonemptyString("foo bar")()
		require.NoError(t, err)
		require.Equal(t, "foo bar", val)
	})

	t.Run("returns ErrParameterEmpty when given an empty string (including only spaces)", func(t *testing.T) {
		t.Parallel()
		_, err := pipeline.NonemptyString("")()
		require.Equal(t, pipeline.ErrParameterEmpty, errors.Cause(err))
		_, err = pipeline.NonemptyString(" ")()
		require.Equal(t, pipeline.ErrParameterEmpty, errors.Cause(err))
	})
}

func TestGetters_Input(t *testing.T) {
	t.Parallel()

	t.Run("returns the requested input's Value and Error if they exist", func(t *testing.T) {
		t.Parallel()
		expectedVal := "bar"
		expectedErr := errors.New("some err")
		val, err := pipeline.Input([]pipeline.Result{{Value: "foo"}, {Value: expectedVal, Error: expectedErr}, {Value: "baz"}}, 1)()
		require.Equal(t, expectedVal, val)
		require.Equal(t, expectedErr, err)
	})

	t.Run("returns ErrParameterEmpty if the specified input does not exist", func(t *testing.T) {
		t.Parallel()
		_, err := pipeline.Input([]pipeline.Result{{Value: "foo"}}, 1)()
		require.Equal(t, pipeline.ErrParameterEmpty, errors.Cause(err))
	})
}

func TestGetters_Inputs(t *testing.T) {
	t.Parallel()

	theErr := errors.New("some issue")

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
				{Error: theErr},
				{Value: "baz"},
			},
			[]interface{}{"foo", theErr, "baz"}, nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run("returns all of the inputs' Values if the provided inputs meet the spec", func(t *testing.T) {
			t.Parallel()

			val, err := pipeline.Inputs(test.inputs)()
			require.Equal(t, test.expectedErr, errors.Cause(err))
			require.Equal(t, test.expected, val)
		})
	}
}
