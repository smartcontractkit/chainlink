package pipeline_test

import (
	"net/url"
	"testing"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestStringParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		{"string", "foo bar baz", pipeline.StringParam("foo bar baz"), nil},
		{"[]byte", []byte("foo bar baz"), pipeline.StringParam("foo bar baz"), nil},
		{"int", 12345, pipeline.StringParam(""), pipeline.ErrBadInput},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var p pipeline.StringParam
			err := p.UnmarshalPipelineParam(test.input, nil)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
}

func TestBytesParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		{"string", "foo bar baz", pipeline.BytesParam("foo bar baz"), nil},
		{"[]byte", []byte("foo bar baz"), pipeline.BytesParam("foo bar baz"), nil},
		{"int", 12345, pipeline.BytesParam(""), pipeline.ErrBadInput},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var p pipeline.BytesParam
			err := p.UnmarshalPipelineParam(test.input, nil)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
}

func TestUint64Param_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		{"string", "123", pipeline.Uint64Param(123), nil},
		{"int", int(123), pipeline.Uint64Param(123), nil},
		{"int8", int8(123), pipeline.Uint64Param(123), nil},
		{"int16", int16(123), pipeline.Uint64Param(123), nil},
		{"int32", int32(123), pipeline.Uint64Param(123), nil},
		{"int64", int64(123), pipeline.Uint64Param(123), nil},
		{"uint", uint(123), pipeline.Uint64Param(123), nil},
		{"uint8", uint8(123), pipeline.Uint64Param(123), nil},
		{"uint16", uint16(123), pipeline.Uint64Param(123), nil},
		{"uint32", uint32(123), pipeline.Uint64Param(123), nil},
		{"uint64", uint64(123), pipeline.Uint64Param(123), nil},
		{"bool", true, pipeline.Uint64Param(0), pipeline.ErrBadInput},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var p pipeline.Uint64Param
			err := p.UnmarshalPipelineParam(test.input, nil)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
}

func TestBoolParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		{"string true", "true", pipeline.BoolParam(true), nil},
		{"string false", "false", pipeline.BoolParam(false), nil},
		{"bool true", true, pipeline.BoolParam(true), nil},
		{"bool false", false, pipeline.BoolParam(false), nil},
		{"int", int8(123), pipeline.BoolParam(false), pipeline.ErrBadInput},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var p pipeline.BoolParam
			err := p.UnmarshalPipelineParam(test.input, nil)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
}

func TestMaybeBoolParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         interface{}
		expected      interface{}
		expectedIsSet bool
		err           error
	}{
		{"string true", "true", pipeline.MaybeBoolTrue, true, nil},
		{"string false", "false", pipeline.MaybeBoolFalse, true, nil},
		{"string empty", "", pipeline.MaybeBoolNull, false, nil},
		{"bool true", true, pipeline.MaybeBoolTrue, true, nil},
		{"bool false", false, pipeline.MaybeBoolFalse, true, nil},
		{"int", int8(123), pipeline.MaybeBoolNull, false, pipeline.ErrBadInput},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var p pipeline.MaybeBoolParam
			err := p.UnmarshalPipelineParam(test.input, nil)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)

			_, isSet := p.Bool()
			require.Equal(t, test.expectedIsSet, isSet)
		})
	}
}

func TestDecimalParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	d := decimal.NewFromFloat(123.45)
	dNull := decimal.Decimal{}

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		{"string", "123.45", pipeline.DecimalParam(d), nil},
		{"float32", float32(123.45), pipeline.DecimalParam(d), nil},
		{"float64", float64(123.45), pipeline.DecimalParam(d), nil},
		{"bool", false, pipeline.DecimalParam(dNull), pipeline.ErrBadInput},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var p pipeline.DecimalParam
			err := p.UnmarshalPipelineParam(test.input, nil)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
}

func TestURLParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	good, err := url.ParseRequestURI("https://chain.link/foo?bar=sergey")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		{"good", "https://chain.link/foo?bar=sergey", pipeline.URLParam(*good), nil},
		{"bad", "asdlkfjlskdfj", pipeline.URLParam(url.URL{}), pipeline.ErrBadInput},
		{"bool", true, pipeline.URLParam(url.URL{}), pipeline.ErrBadInput},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var p pipeline.URLParam
			err := p.UnmarshalPipelineParam(test.input, nil)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
}

func TestMapParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	vars := pipeline.Vars{
		"foo": map[string]interface{}{
			"abc": "def",
		},
		"bar": "123",
	}

	inputStr := `
    {
        "chain": $(foo),
        "link": {
            $(bar): "satoshi",
            "sergey": $(foo.abc)
        }
    }`

	inputMap := map[string]interface{}{
		"chain": map[string]interface{}{
			"abc": "def",
		},
		"link": map[string]interface{}{
			"sergey": "def",
			"123":    "satoshi",
		},
	}

	expected := pipeline.MapParam{
		"chain": map[string]interface{}{
			"abc": "def",
		},
		"link": map[string]interface{}{
			"sergey": "def",
			"123":    "satoshi",
		},
	}

	var got1 pipeline.MapParam
	err := got1.UnmarshalPipelineParam(inputStr, vars)
	require.NoError(t, err)
	require.Equal(t, expected, got1)

	var got2 pipeline.MapParam
	err = got2.UnmarshalPipelineParam(inputMap, vars)
	require.NoError(t, err)
	require.Equal(t, expected, got2)
}

func TestSliceParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		{"[]interface{}", []interface{}{1, 2, 3}, pipeline.SliceParam([]interface{}{1, 2, 3}), nil},
		{"[]interface{} with var", []interface{}{1, 2, "$(foo)"}, pipeline.SliceParam([]interface{}{1, 2, "$(foo)"}), nil},
		{"[]byte", []byte(`[1, 2, 3]`), pipeline.SliceParam([]interface{}{float64(1), float64(2), float64(3)}), nil},
		{"[]byte with var", []byte(`[1, 2, $(foo)]`), pipeline.SliceParam([]interface{}{float64(1), float64(2), "42"}), nil},
		{"string", `[1, 2, 3]`, pipeline.SliceParam([]interface{}{float64(1), float64(2), float64(3)}), nil},
		{"string with var", `[1, 2, $(foo)]`, pipeline.SliceParam([]interface{}{float64(1), float64(2), "42"}), nil},
		{"bool", true, pipeline.SliceParam(nil), pipeline.ErrBadInput},
	}

	vars := pipeline.Vars{"foo": "42"}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var p pipeline.SliceParam
			err := p.UnmarshalPipelineParam(test.input, vars)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
}

func TestSliceParam_FilterErrors(t *testing.T) {
	t.Parallel()

	s := pipeline.SliceParam{"foo", errors.New("bar"), "baz"}
	vals, n := s.FilterErrors()
	require.Equal(t, 1, n)
	require.Equal(t, pipeline.SliceParam{"foo", "baz"}, vals)
}

func TestDecimalSliceParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	expected := pipeline.DecimalSliceParam{*mustDecimal(t, "1.1"), *mustDecimal(t, "2.2"), *mustDecimal(t, "3.3")}

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		{"[]interface{}", []interface{}{1.1, "2.2", *mustDecimal(t, "3.3")}, expected, nil},
		{"string", `[1.1, "2.2", 3.3]`, expected, nil},
		{"[]byte", `[1.1, "2.2", 3.3]`, expected, nil},
		{"[]interface{} with error", `[1.1, true, "abc"]`, pipeline.DecimalSliceParam(nil), pipeline.ErrBadInput},
		{"bool", true, pipeline.DecimalSliceParam(nil), pipeline.ErrBadInput},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var p pipeline.DecimalSliceParam
			err := p.UnmarshalPipelineParam(test.input, nil)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
}

func TestStringSliceParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	expected := pipeline.StringSliceParam{"1.1", "2.2", "3.3", "sergey"}

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		{"[]interface{}", []interface{}{"1.1", "2.2", "3.3", "sergey"}, expected, nil},
		{"string", `1.1,2.2,3.3,sergey`, expected, nil},
		{"[]byte", []byte(`1.1,2.2,3.3,sergey`), expected, nil},
		{"bool", true, pipeline.StringSliceParam(nil), pipeline.ErrBadInput},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var p pipeline.StringSliceParam
			err := p.UnmarshalPipelineParam(test.input, nil)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
}
