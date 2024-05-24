package pipeline_test

import (
	"math/big"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestGetters_VarExpr(t *testing.T) {
	t.Parallel()

	vars := createTestVars()

	tests := []struct {
		expr   string
		result interface{}
		err    error
	}{
		// no errors
		{"$(foo.bar)", "value", nil},
		{" $(foo.bar)", "value", nil},
		{"$(foo.bar) ", "value", nil},
		{"$( foo.bar)", "value", nil},
		{"$(foo.bar )", "value", nil},
		{"$( foo.bar )", "value", nil},
		{" $( foo.bar )", "value", nil},
		// errors
		{"  ", nil, pipeline.ErrParameterEmpty},
		{"$()", nil, pipeline.ErrParameterEmpty},
		{"$", nil, pipeline.ErrParameterEmpty},
		{"$(", nil, pipeline.ErrParameterEmpty},
		{"$)", nil, pipeline.ErrParameterEmpty},
		{"()", nil, pipeline.ErrParameterEmpty},
		{"$(foo.bar", nil, pipeline.ErrParameterEmpty},
		{"$foo.bar)", nil, pipeline.ErrParameterEmpty},
		{"(foo.bar)", nil, pipeline.ErrParameterEmpty},
		{"foo.bar", nil, pipeline.ErrParameterEmpty},
		{"$(err)", nil, pipeline.ErrTooManyErrors},
	}

	for _, test := range tests {
		t.Run(test.expr, func(t *testing.T) {
			getter := pipeline.VarExpr(test.expr, vars)
			v, err := getter()
			if test.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, test.result, v)
			} else {
				assert.Equal(t, test.err, errors.Cause(err))
			}
		})
	}
}

func TestGetters_JSONWithVarExprs(t *testing.T) {
	t.Parallel()

	vars := createTestVars()

	errVal, err := vars.Get("err")
	require.NoError(t, err)

	big, ok := new(big.Int).SetString("314159265358979323846264338327950288419716939937510582097494459", 10)
	require.True(t, ok)

	tests := []struct {
		json        string
		field       string
		result      interface{}
		err         error
		allowErrors bool
	}{
		// no errors
		{`{ "x": $(zet) }`, "x", 123, nil, false},
		{`{ "x": $( zet ) }`, "x", 123, nil, false},
		{`{ "x": { "y": $(zet) } }`, "x", map[string]interface{}{"y": 123}, nil, false},
		{`{ "z": "foo" }`, "z", "foo", nil, false},
		{`{ "a": $(arr.1) }`, "a", 200, nil, false},
		{`{}`, "", map[string]interface{}{}, nil, false},
		{`{ "e": $(err) }`, "e", errVal, nil, true},
		{`null`, "", nil, nil, false},
		{`{ "x": 314159265358979323846264338327950288419716939937510582097494459 }`, "x", big, nil, false},
		{`{ "x": 3141592653589 }`, "x", int64(3141592653589), nil, false},
		{`{ "x": 18446744073709551615 }`, "x", uint64(18446744073709551615), nil, false},
		{`{ "x": 3141592653589.567 }`, "x", float64(3141592653589.567), nil, false},
		// errors
		{`  `, "", nil, pipeline.ErrParameterEmpty, false},
		{`{ "x": $(missing) }`, "x", nil, pipeline.ErrKeypathNotFound, false},
		{`{ "x": "$(zet)" }`, "x", "$(zet)", pipeline.ErrBadInput, false},
		{`{ "$(foo.bar)": $(zet) }`, "value", 123, pipeline.ErrBadInput, false},
		{`{ "x": { "__chainlink_key_path__": 0 } }`, "", nil, pipeline.ErrBadInput, false},
		{`{ "e": $(err)`, "e", nil, pipeline.ErrBadInput, false},
	}

	for _, test := range tests {
		t.Run(test.json, func(t *testing.T) {
			getter := pipeline.JSONWithVarExprs(test.json, vars, test.allowErrors)
			v, err := getter()
			if test.err != nil {
				assert.Equal(t, test.err, errors.Cause(err))
			} else {
				m, is := v.(map[string]interface{})
				if is && test.field != "" {
					assert.Equal(t, test.result, m[test.field])
				} else {
					assert.Equal(t, test.result, v)
				}
			}
		})
	}
}

func TestGetters_Input(t *testing.T) {
	t.Parallel()

	t.Run("returns the requested input's Value and Error if they exist", func(t *testing.T) {
		expectedVal := "bar"
		expectedErr := errors.New("some err")
		val, err := pipeline.Input([]pipeline.Result{{Value: "foo"}, {Value: expectedVal, Error: expectedErr}, {Value: "baz"}}, 1)()
		assert.Equal(t, expectedVal, val)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("returns ErrIndexOutOfRange if the specified index is out of range", func(t *testing.T) {
		_, err := pipeline.Input([]pipeline.Result{{Value: "foo"}}, 1)()
		assert.Equal(t, pipeline.ErrIndexOutOfRange, errors.Cause(err))
		_, err = pipeline.Input([]pipeline.Result{{Value: "foo"}}, -1)()
		assert.Equal(t, pipeline.ErrIndexOutOfRange, errors.Cause(err))
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
			[]interface{}{"foo", theErr, "baz"},
			nil,
		},
		{
			"returns nil array",
			[]pipeline.Result{},
			nil,
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val, err := pipeline.Inputs(test.inputs)()
			assert.Equal(t, test.expectedErr, errors.Cause(err))
			assert.Equal(t, test.expected, val)
		})
	}
}

func TestGetters_NonemptyString(t *testing.T) {
	t.Parallel()

	t.Run("returns any non-empty string", func(t *testing.T) {
		val, err := pipeline.NonemptyString("foo bar")()
		assert.NoError(t, err)
		assert.Equal(t, "foo bar", val)
	})

	t.Run("returns ErrParameterEmpty when given an empty string (including only spaces)", func(t *testing.T) {
		_, err := pipeline.NonemptyString("")()
		assert.Equal(t, pipeline.ErrParameterEmpty, errors.Cause(err))
		_, err = pipeline.NonemptyString(" ")()
		assert.Equal(t, pipeline.ErrParameterEmpty, errors.Cause(err))
	})
}

func TestGetters_ValidDurationInSeconds(t *testing.T) {
	t.Parallel()

	t.Run("returns duration in seconds", func(t *testing.T) {
		val, err := pipeline.ValidDurationInSeconds("10s")()
		assert.NoError(t, err)
		assert.Equal(t, 10, val)

		val, err = pipeline.ValidDurationInSeconds("1m")()
		assert.NoError(t, err)
		assert.Equal(t, 60, val)

		val, err = pipeline.ValidDurationInSeconds("1h")()
		assert.NoError(t, err)
		assert.Equal(t, 3600, val)
	})

	t.Run("returns ErrParameterEmpty when given an empty string (including only spaces)", func(t *testing.T) {
		_, err := pipeline.ValidDurationInSeconds("")()
		assert.Equal(t, pipeline.ErrParameterEmpty, errors.Cause(err))
		_, err = pipeline.ValidDurationInSeconds(" ")()
		assert.Equal(t, pipeline.ErrParameterEmpty, errors.Cause(err))
	})

	t.Run("returns duration errors when given invalid durations", func(t *testing.T) {
		_, err := pipeline.ValidDurationInSeconds("1b")()
		assert.Contains(t, err.Error(), "unknown unit")
		_, err = pipeline.ValidDurationInSeconds("5")()
		assert.Contains(t, err.Error(), "missing unit")
		_, err = pipeline.ValidDurationInSeconds("!m")()
		assert.Contains(t, err.Error(), "invalid duration")
	})
}

func TestGetters_From(t *testing.T) {
	t.Parallel()

	t.Run("no inputs", func(t *testing.T) {
		getters := pipeline.From()
		assert.Empty(t, getters)
	})

	var fooGetter1 pipeline.GetterFunc = func() (interface{}, error) {
		return "foo", nil
	}
	var fooGetter2 pipeline.GetterFunc = func() (interface{}, error) {
		return "foo", nil
	}

	tests := []struct {
		name     string
		input    []interface{}
		expected string
	}{
		{
			"only getters",
			[]interface{}{fooGetter1, fooGetter2},
			"foo",
		},
		{
			"mix of getters and values",
			[]interface{}{fooGetter1, "foo"},
			"foo",
		},
		{
			"only values",
			[]interface{}{"foo", "foo"},
			"foo",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			getters := pipeline.From(test.input...)
			assert.Len(t, getters, 2)

			for _, getter := range getters {
				val, err := getter()
				assert.NoError(t, err)
				assert.Equal(t, test.expected, val)
			}
		})
	}
}

func createTestVars() pipeline.Vars {
	return pipeline.NewVarsFrom(map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "value",
		},
		"zet": 123,
		"arr": []interface{}{
			100, 200, 300,
		},
		"err": errors.New("some error"),
	})
}
