package pipeline_test

import (
	"net/url"
	"testing"

	"github.com/ethereum/go-ethereum/common"

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
		{"object", pipeline.MustNewObjectParam(`boz bar bap`), pipeline.StringParam("boz bar bap"), nil},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var p pipeline.StringParam
			err := p.UnmarshalPipelineParam(test.input)
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
		{"string", "0x11AAFF", pipeline.BytesParam([]byte{0x11, 0xAA, 0xFF}), nil},
		{"[]byte", []byte("foo bar baz"), pipeline.BytesParam("foo bar baz"), nil},
		{"int", 12345, pipeline.BytesParam(nil), pipeline.ErrBadInput},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var p pipeline.BytesParam
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
}

func TestBytesParam_MarshalJSON(t *testing.T) {
	t.Parallel()

	bp := pipeline.BytesParam([]byte{0x11, 0xAA, 0xFF})
	json, err := bp.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `"0x11aaff"`, string(json))
}

func TestAddressParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	var addr pipeline.AddressParam
	copy(addr[:], []byte("deadbeefdeadbeefdead"))

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		{"20-char string", "deadbeefdeadbeefdead", addr, nil},
		{"21-char string", "deadbeefdeadbeefdeadb", nil, pipeline.ErrBadInput},
		{"19-char string", "deadbeefdeadbeefdea", nil, pipeline.ErrBadInput},
		{"20-char []byte", []byte("deadbeefdeadbeefdead"), addr, nil},
		{"21-char []byte", []byte("deadbeefdeadbeefdeadb"), nil, pipeline.ErrBadInput},
		{"19-char []byte", []byte("deadbeefdeadbeefdea"), nil, pipeline.ErrBadInput},

		{"42-char string with 0x", "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef", pipeline.AddressParam(common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")), nil},
		{"41-char string with 0x", "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbee", nil, pipeline.ErrBadInput},
		{"43-char string with 0x", "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefd", nil, pipeline.ErrBadInput},
		{"42-char string without 0x", "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefde", nil, pipeline.ErrBadInput},
		{"40-char string without 0x", "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef", nil, pipeline.ErrBadInput},

		{"42-char []byte with 0x", []byte("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"), pipeline.AddressParam(common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")), nil},
		{"41-char []byte with 0x", []byte("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbee"), nil, pipeline.ErrBadInput},
		{"43-char []byte with 0x", []byte("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefd"), nil, pipeline.ErrBadInput},
		{"42-char []byte without 0x", []byte("deadbeefdeadbeefdeadbeefdeadbeefdeadbeefde"), nil, pipeline.ErrBadInput},
		{"40-char []byte without 0x", []byte("deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"), nil, pipeline.ErrBadInput},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var p pipeline.AddressParam
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			if test.expected != nil {
				require.Equal(t, test.expected, p)
			}
		})
	}
}

func TestAddressSliceParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	addr1 := common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	addr2 := common.HexToAddress("0xcafebabecafebabecafebabecafebabecafebabe")
	expected := pipeline.AddressSliceParam{addr1, addr2}

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		{"json", `[ "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef", "0xcafebabecafebabecafebabecafebabecafebabe" ]`, expected, nil},
		{"[]common.Address", []common.Address{addr1, addr2}, expected, nil},
		{"[]interface{} with common.Address", []interface{}{addr1, addr2}, expected, nil},
		{"[]interface{} with strings", []interface{}{addr1.String(), addr2.String()}, expected, nil},
		{"[]interface{} with []byte", []interface{}{[]byte(addr1.String()), []byte(addr2.String())}, expected, nil},
		{"[]interface{} with []byte", []interface{}{[]byte(addr1.String()), []byte(addr2.String())}, expected, nil},
		{"nil", nil, pipeline.AddressSliceParam(nil), nil},

		{"bad json", `[ "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef" "0xcafebabecafebabecafebabecafebabecafebabe" ]`, nil, pipeline.ErrBadInput},
		{"[]interface{} with bad types", []interface{}{123, true}, nil, pipeline.ErrBadInput},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var p pipeline.AddressSliceParam
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			if test.expected != nil {
				require.Equal(t, test.expected, p)
			}
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
			err := p.UnmarshalPipelineParam(test.input)
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
		{"object", pipeline.MustNewObjectParam(true), pipeline.BoolParam(true), nil},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var p pipeline.BoolParam
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
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
		{"object", pipeline.MustNewObjectParam(123.45), pipeline.DecimalParam(d), nil},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var p pipeline.DecimalParam
			err := p.UnmarshalPipelineParam(test.input)
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
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
}

func TestMapParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	inputStr := `
    {
        "chain": {"abc": "def"},
        "link": {
            "123": "satoshi",
            "sergey": "def"
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
	err := got1.UnmarshalPipelineParam(inputStr)
	require.NoError(t, err)
	require.Equal(t, expected, got1)

	var got2 pipeline.MapParam
	err = got2.UnmarshalPipelineParam(inputMap)
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
		{"[]byte", []byte(`[1, 2, 3]`), pipeline.SliceParam([]interface{}{float64(1), float64(2), float64(3)}), nil},
		{"string", `[1, 2, 3]`, pipeline.SliceParam([]interface{}{float64(1), float64(2), float64(3)}), nil},
		{"bool", true, pipeline.SliceParam(nil), pipeline.ErrBadInput},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var p pipeline.SliceParam
			err := p.UnmarshalPipelineParam(test.input)
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
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
}

func TestJSONPathParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	expected := pipeline.JSONPathParam{"1.1", "2.2", "3.3", "sergey"}

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		{"[]interface{}", []interface{}{"1.1", "2.2", "3.3", "sergey"}, expected, nil},
		{"string", `1.1,2.2,3.3,sergey`, expected, nil},
		{"[]byte", []byte(`1.1,2.2,3.3,sergey`), expected, nil},
		{"bool", true, pipeline.JSONPathParam(nil), pipeline.ErrBadInput},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var p pipeline.JSONPathParam
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
}
