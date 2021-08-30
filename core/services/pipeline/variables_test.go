package pipeline_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
)

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
		require.Equal(t, pipeline.ErrKeypathTooDeep, errors.Cause(err))
	})

	t.Run("errors when getting a value at a keypath where the first part is not a map/slice", func(t *testing.T) {
		vars := pipeline.NewVarsFrom(map[string]interface{}{
			"foo": 123,
		})
		_, err := vars.Get("foo.bar")
		require.Equal(t, pipeline.ErrKeypathNotFound, errors.Cause(err))
	})

	t.Run("errors when getting a value at a keypath with more than 2 components", func(t *testing.T) {
		vars := pipeline.NewVarsFrom(map[string]interface{}{
			"foo": 123,
		})
		_, err := vars.Get("foo.bar.baz")
		require.Equal(t, pipeline.ErrKeypathTooDeep, errors.Cause(err))
	})
}

func TestResolveValue(t *testing.T) {
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

		expectedValue := 123
		expectedErr := errors.New("some issue")

		param := new(mocks.PipelineParamUnmarshaler)
		param.On("UnmarshalPipelineParam", expectedValue).Return(expectedErr)

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

func TestGetters_VarExpr(t *testing.T) {
	t.Parallel()

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": 42,
		},
	})

	tests := []struct {
		expr  string
		value interface{}
		err   error
	}{
		{"$(foo.bar)", 42, nil},
		{" $(foo.bar)", 42, nil},
		{"$(foo.bar) ", 42, nil},
		{"$( foo.bar)", 42, nil},
		{"$(foo.bar )", 42, nil},
		{"$( foo.bar )", 42, nil},
		{" $( foo.bar )", 42, nil},
		{"$()", nil, pipeline.ErrVarsRoot},
		{"$(foo.bar", nil, pipeline.ErrParameterEmpty},
		{"$foo.bar)", nil, pipeline.ErrParameterEmpty},
		{"(foo.bar)", nil, pipeline.ErrParameterEmpty},
		{"foo.bar", nil, pipeline.ErrParameterEmpty},
	}

	for _, test := range tests {
		test := test
		t.Run(test.expr, func(t *testing.T) {
			val, err := pipeline.VarExpr(test.expr, vars)()
			require.Equal(t, test.value, val)
			require.Equal(t, test.err, errors.Cause(err))
		})
	}
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

func TestKeypath(t *testing.T) {
	t.Run("can be constructed from a period-delimited string with 2 or fewer parts", func(t *testing.T) {
		kp, err := pipeline.NewKeypathFromString("")
		require.NoError(t, err)
		require.Equal(t, pipeline.Keypath{nil, nil}, kp)

		kp, err = pipeline.NewKeypathFromString("foo")
		require.NoError(t, err)
		require.Equal(t, pipeline.Keypath{[]byte("foo"), nil}, kp)

		kp, err = pipeline.NewKeypathFromString("foo.bar")
		require.NoError(t, err)
		require.Equal(t, pipeline.Keypath{[]byte("foo"), []byte("bar")}, kp)
	})

	t.Run("errors if constructor is passed more than 2 parts", func(t *testing.T) {
		_, err := pipeline.NewKeypathFromString("foo.bar.baz")
		require.Equal(t, pipeline.ErrKeypathTooDeep, errors.Cause(err))
	})

	t.Run("accurately reports its NumParts", func(t *testing.T) {
		kp, err := pipeline.NewKeypathFromString("")
		require.NoError(t, err)
		require.Equal(t, pipeline.Keypath{nil, nil}, kp)
		require.Equal(t, 0, kp.NumParts())

		kp, err = pipeline.NewKeypathFromString("foo")
		require.NoError(t, err)
		require.Equal(t, pipeline.Keypath{[]byte("foo"), nil}, kp)
		require.Equal(t, 1, kp.NumParts())

		kp, err = pipeline.NewKeypathFromString("foo.bar")
		require.NoError(t, err)
		require.Equal(t, pipeline.Keypath{[]byte("foo"), []byte("bar")}, kp)
		require.Equal(t, 2, kp.NumParts())
	})

	t.Run("stringifies correctly", func(t *testing.T) {
		kp, err := pipeline.NewKeypathFromString("")
		require.NoError(t, err)
		require.Equal(t, pipeline.Keypath{nil, nil}, kp)
		require.Equal(t, "(empty)", kp.String())

		kp, err = pipeline.NewKeypathFromString("foo")
		require.NoError(t, err)
		require.Equal(t, pipeline.Keypath{[]byte("foo"), nil}, kp)
		require.Equal(t, "foo", kp.String())

		kp, err = pipeline.NewKeypathFromString("foo.bar")
		require.NoError(t, err)
		require.Equal(t, pipeline.Keypath{[]byte("foo"), []byte("bar")}, kp)
		require.Equal(t, "foo.bar", kp.String())
	})
}
