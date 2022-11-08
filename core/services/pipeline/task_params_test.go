package pipeline_test

import (
	"math"
	"math/big"
	"net/url"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
)

func TestStringParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	var nilObjectParam *pipeline.ObjectParam

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		// valid
		{"string", "foo bar baz", pipeline.StringParam("foo bar baz"), nil},
		{"[]byte", []byte("foo bar baz"), pipeline.StringParam("foo bar baz"), nil},
		{"*object", mustNewObjectParam(t, `boz bar bap`), pipeline.StringParam("boz bar bap"), nil},
		{"object", *mustNewObjectParam(t, `boz bar bap`), pipeline.StringParam("boz bar bap"), nil},
		// invalid
		{"int", 12345, pipeline.StringParam(""), pipeline.ErrBadInput},
		{"nil", nil, pipeline.StringParam(""), pipeline.ErrBadInput},
		{"nil ObjectParam", nilObjectParam, pipeline.StringParam(""), pipeline.ErrBadInput},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var p pipeline.StringParam
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
}

func TestStringSliceParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	expected := pipeline.StringSliceParam{"foo", "bar", "baz"}

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		{"json", `[ "foo", "bar", "baz" ]`, expected, nil},
		{"[]string", []string{"foo", "bar", "baz"}, expected, nil},
		{"[]interface{} with strings", []interface{}{"foo", "bar", "baz"}, expected, nil},
		{"[]interface{} with []byte", []interface{}{[]byte("foo"), []byte("bar"), []byte("baz")}, expected, nil},
		{"SliceParam", pipeline.SliceParam([]interface{}{"foo", "bar", "baz"}), expected, nil},

		{"nil", nil, pipeline.StringSliceParam(nil), nil},

		{"bad json", `[ "foo", 1, false ]`, nil, pipeline.ErrBadInput},
		{"[]interface{} with bad types", []interface{}{123, true}, nil, pipeline.ErrBadInput},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var p pipeline.StringSliceParam
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			if test.expected != nil {
				require.Equal(t, test.expected, p)
			}
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
		{"int", 12345, pipeline.BytesParam(nil), pipeline.ErrBadInput},
		{"hex-invalid", "0xh", pipeline.BytesParam("0xh"), nil},
		{"valid-hex", hexutil.MustDecode("0xd3184d"), pipeline.BytesParam(hexutil.MustDecode("0xd3184d")), nil},
		{"*object", mustNewObjectParam(t, `boz bar bap`), pipeline.BytesParam("boz bar bap"), nil},
		{"object", *mustNewObjectParam(t, `boz bar bap`), pipeline.BytesParam("boz bar bap"), nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var p pipeline.BytesParam
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
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

		{"42-char string with 0x but wrong characters", "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadzzzz", nil, pipeline.ErrBadInput},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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
		{"nil", nil, pipeline.AddressSliceParam(nil), nil},

		{"bad json", `[ "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef" "0xcafebabecafebabecafebabecafebabecafebabe" ]`, nil, pipeline.ErrBadInput},
		{"[]interface{} with bad types", []interface{}{123, true}, nil, pipeline.ErrBadInput},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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
		// positive
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
		{"float64", float64(123), pipeline.Uint64Param(123), nil},
		// negative
		{"bool", true, pipeline.Uint64Param(0), pipeline.ErrBadInput},
		{"negative int", int(-123), pipeline.Uint64Param(0), pipeline.ErrBadInput},
		{"negative int8", int8(-123), pipeline.Uint64Param(0), pipeline.ErrBadInput},
		{"negative int16", int16(-123), pipeline.Uint64Param(0), pipeline.ErrBadInput},
		{"negative int32", int32(-123), pipeline.Uint64Param(0), pipeline.ErrBadInput},
		{"negative int64", int64(-123), pipeline.Uint64Param(0), pipeline.ErrBadInput},
		{"negative float64", float64(-123), pipeline.Uint64Param(0), pipeline.ErrBadInput},
		{"out of bounds float64", math.MaxFloat64, pipeline.Uint64Param(0), pipeline.ErrBadInput},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var p pipeline.Uint64Param
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			if test.err == nil {
				require.Equal(t, test.expected, p)
			}
		})
	}
}

func TestMaybeUint64Param_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		// positive
		{"string", "123", pipeline.NewMaybeUint64Param(123, true), nil},
		{"int", int(123), pipeline.NewMaybeUint64Param(123, true), nil},
		{"int8", int8(123), pipeline.NewMaybeUint64Param(123, true), nil},
		{"int16", int16(123), pipeline.NewMaybeUint64Param(123, true), nil},
		{"int32", int32(123), pipeline.NewMaybeUint64Param(123, true), nil},
		{"int64", int64(123), pipeline.NewMaybeUint64Param(123, true), nil},
		{"uint", uint(123), pipeline.NewMaybeUint64Param(123, true), nil},
		{"uint8", uint8(123), pipeline.NewMaybeUint64Param(123, true), nil},
		{"uint16", uint16(123), pipeline.NewMaybeUint64Param(123, true), nil},
		{"uint32", uint32(123), pipeline.NewMaybeUint64Param(123, true), nil},
		{"uint64", uint64(123), pipeline.NewMaybeUint64Param(123, true), nil},
		{"float64", float64(123), pipeline.NewMaybeUint64Param(123, true), nil},
		{"empty string", "", pipeline.NewMaybeUint64Param(0, false), nil},
		// negative
		{"bool", true, pipeline.NewMaybeUint64Param(0, false), pipeline.ErrBadInput},
		{"negative int", int(-123), pipeline.NewMaybeUint64Param(0, false), pipeline.ErrBadInput},
		{"negative int8", int8(-123), pipeline.NewMaybeUint64Param(0, false), pipeline.ErrBadInput},
		{"negative int16", int16(-123), pipeline.NewMaybeUint64Param(0, false), pipeline.ErrBadInput},
		{"negative int32", int32(-123), pipeline.NewMaybeUint64Param(0, false), pipeline.ErrBadInput},
		{"negative int64", int64(-123), pipeline.NewMaybeUint64Param(0, false), pipeline.ErrBadInput},
		{"negative float64", float64(-123), pipeline.NewMaybeUint64Param(0, false), pipeline.ErrBadInput},
		{"out of bounds float64", math.MaxFloat64, pipeline.NewMaybeUint64Param(0, false), pipeline.ErrBadInput},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var p pipeline.MaybeUint64Param
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			if err == nil {
				require.Equal(t, test.expected, p)
			}
		})
	}
}

func TestMaybeBigIntParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	fromInt := func(n int64) pipeline.MaybeBigIntParam {
		return pipeline.NewMaybeBigIntParam(big.NewInt(n))
	}

	intDecimal := *mustDecimal(t, "123")
	floatDecimal := *mustDecimal(t, "123.45")

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		// positive
		{"string", "123", fromInt(123), nil},
		{"empty string", "", pipeline.NewMaybeBigIntParam(nil), nil},
		{"nil", nil, pipeline.NewMaybeBigIntParam(nil), nil},
		{"*big.Int", big.NewInt(123), fromInt(123), nil},
		{"int", int(123), fromInt(123), nil},
		{"int8", int8(123), fromInt(123), nil},
		{"int16", int16(123), fromInt(123), nil},
		{"int32", int32(123), fromInt(123), nil},
		{"int64", int64(123), fromInt(123), nil},
		{"uint", uint(123), fromInt(123), nil},
		{"uint8", uint8(123), fromInt(123), nil},
		{"uint16", uint16(123), fromInt(123), nil},
		{"uint32", uint32(123), fromInt(123), nil},
		{"uint64", uint64(123), fromInt(123), nil},
		{"float64", float64(123), fromInt(123), nil},
		{"float64", float64(-123), fromInt(-123), nil},
		{"decimal.Decimal", intDecimal, fromInt(123), nil},
		{"*decimal.Decimal", &intDecimal, fromInt(123), nil},
		// negative
		{"bool", true, pipeline.NewMaybeBigIntParam(nil), pipeline.ErrBadInput},
		{"negative out of bound float64", -math.MaxFloat64, pipeline.NewMaybeBigIntParam(nil), pipeline.ErrBadInput},
		{"positive out of bound float64", math.MaxFloat64, pipeline.NewMaybeBigIntParam(nil), pipeline.ErrBadInput},
		{"non-integer decimal.Decimal", floatDecimal, pipeline.NewMaybeBigIntParam(nil), pipeline.ErrBadInput},
		{"non-integer *decimal.Decimal", &floatDecimal, pipeline.NewMaybeBigIntParam(nil), pipeline.ErrBadInput},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var p pipeline.MaybeBigIntParam
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			if test.err == nil {
				require.Equal(t, test.expected, p)
			}
		})
	}
}

func TestMaybeInt32Param_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		{"string", "123", pipeline.NewMaybeInt32Param(123, true), nil},
		{"int", int(123), pipeline.NewMaybeInt32Param(123, true), nil},
		{"int8", int8(123), pipeline.NewMaybeInt32Param(123, true), nil},
		{"int16", int16(123), pipeline.NewMaybeInt32Param(123, true), nil},
		{"int32", int32(123), pipeline.NewMaybeInt32Param(123, true), nil},
		{"int64", int64(123), pipeline.NewMaybeInt32Param(123, true), nil},
		{"uint", uint(123), pipeline.NewMaybeInt32Param(123, true), nil},
		{"uint8", uint8(123), pipeline.NewMaybeInt32Param(123, true), nil},
		{"uint16", uint16(123), pipeline.NewMaybeInt32Param(123, true), nil},
		{"uint32", uint32(123), pipeline.NewMaybeInt32Param(123, true), nil},
		{"uint64", uint64(123), pipeline.NewMaybeInt32Param(123, true), nil},
		{"float64", float64(123), pipeline.NewMaybeInt32Param(123, true), nil},
		{"bool", true, pipeline.NewMaybeInt32Param(0, false), pipeline.ErrBadInput},
		{"empty string", "", pipeline.NewMaybeInt32Param(0, false), nil},
		{"string overflow", "100000000000", pipeline.NewMaybeInt32Param(0, false), pipeline.ErrBadInput},
		{"int64 overflow", int64(123 << 32), pipeline.NewMaybeInt32Param(0, false), pipeline.ErrBadInput},
		{"negative int64 overflow", -int64(123 << 32), pipeline.NewMaybeInt32Param(0, false), pipeline.ErrBadInput},
		{"uint64 overflow", uint64(123 << 32), pipeline.NewMaybeInt32Param(0, false), pipeline.ErrBadInput},
		{"float overflow", float64(123 << 32), pipeline.NewMaybeInt32Param(0, false), pipeline.ErrBadInput},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var p pipeline.MaybeInt32Param
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
		{"*object", mustNewObjectParam(t, true), pipeline.BoolParam(true), nil},
		{"object", *mustNewObjectParam(t, true), pipeline.BoolParam(true), nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var p pipeline.BoolParam
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
}

func TestDecimalParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	var nilObjectParam *pipeline.ObjectParam
	d := decimal.NewFromFloat(123.45)
	dNull := decimal.Decimal{}

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		// valid
		{"string", "123.45", pipeline.DecimalParam(d), nil},
		{"float32", float32(123.45), pipeline.DecimalParam(d), nil},
		{"float64", float64(123.45), pipeline.DecimalParam(d), nil},
		{"object", mustNewObjectParam(t, 123.45), pipeline.DecimalParam(d), nil},
		// invalid
		{"bool", false, pipeline.DecimalParam(dNull), pipeline.ErrBadInput},
		{"nil", nil, pipeline.DecimalParam(dNull), pipeline.ErrBadInput},
		{"nil ObjectParam", nilObjectParam, pipeline.DecimalParam(dNull), pipeline.ErrBadInput},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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
		t.Run(test.name, func(t *testing.T) {
			var p pipeline.URLParam
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
}

func TestMapParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	var nilObjectParam *pipeline.ObjectParam

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

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		// valid
		{"from string", inputStr, expected, nil},
		{"from []byte", []byte(inputStr), expected, nil},
		{"from map", inputMap, expected, nil},
		{"from nil", nil, pipeline.MapParam(nil), nil},
		{"from *object", mustNewObjectParam(t, inputMap), expected, nil},
		{"from object", *mustNewObjectParam(t, inputMap), expected, nil},
		// invalid
		{"wrong type", 123, pipeline.MapParam(nil), pipeline.ErrBadInput},
		{"nil ObjectParam", nilObjectParam, pipeline.MapParam(nil), pipeline.ErrBadInput},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var p pipeline.MapParam
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
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
		{"nil", nil, pipeline.SliceParam(nil), nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var p pipeline.SliceParam
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
}

func TestHashSliceParam_UnmarshalPipelineParam(t *testing.T) {
	t.Parallel()

	hash1 := common.HexToHash("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	hash2 := common.HexToHash("0xcafebabecafebabecafebabecafebabecafebabedeadbeefdeadbeefdeadbeef")
	expected := pipeline.HashSliceParam{hash1, hash2}

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		err      error
	}{
		{"json", `[ "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef", "0xcafebabecafebabecafebabecafebabecafebabedeadbeefdeadbeefdeadbeef" ]`, expected, nil},
		{"[]common.Hash", []common.Hash{hash1, hash2}, expected, nil},
		{"[]interface{} with common.Hash", []interface{}{hash1, hash2}, expected, nil},
		{"[]interface{} with strings", []interface{}{hash1.String(), hash2.String()}, expected, nil},
		{"[]interface{} with []byte", []interface{}{[]byte(hash1.String()), []byte(hash2.String())}, expected, nil},
		{"nil", nil, pipeline.HashSliceParam(nil), nil},
		{"bad json", `[ "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef" "0xcafebabecafebabecafebabecafebabecafebabedeadbeefdeadbeefdeadbeef" ]`, nil, pipeline.ErrBadInput},
		{"[]interface{} with bad types", []interface{}{123, true}, nil, pipeline.ErrBadInput},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var p pipeline.HashSliceParam
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			if test.expected != nil {
				require.Equal(t, test.expected, p)
			}
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
	decimalsSlice := []decimal.Decimal{*mustDecimal(t, "1.1"), *mustDecimal(t, "2.2"), *mustDecimal(t, "3.3")}

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
		{"nil", nil, pipeline.DecimalSliceParam(nil), nil},
		{"[]decimal.Decimal", decimalsSlice, expected, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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
		{"nil", nil, pipeline.JSONPathParam(nil), nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var p pipeline.JSONPathParam
			err := p.UnmarshalPipelineParam(test.input)
			require.Equal(t, test.err, errors.Cause(err))
			require.Equal(t, test.expected, p)
		})
	}
}

func TestResolveValue(t *testing.T) {
	t.Parallel()

	t.Run("calls getters in order until the first one that returns without ErrParameterEmpty", func(t *testing.T) {
		param := mocks.NewPipelineParamUnmarshaler(t)
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
	})

	t.Run("returns any GetterFunc error that isn't ErrParameterEmpty", func(t *testing.T) {
		param := mocks.NewPipelineParamUnmarshaler(t)
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
		expectedValue := 123

		param := mocks.NewPipelineParamUnmarshaler(t)
		param.On("UnmarshalPipelineParam", expectedValue).Return(nil)

		getters := []pipeline.GetterFunc{
			func() (interface{}, error) {
				return expectedValue, nil
			},
		}

		err := pipeline.ResolveParam(param, getters)
		require.NoError(t, err)
	})

	t.Run("returns any error returned by UnmarshalPipelineParam", func(t *testing.T) {
		expectedValue := 123
		expectedErr := errors.New("some issue")

		param := mocks.NewPipelineParamUnmarshaler(t)
		param.On("UnmarshalPipelineParam", expectedValue).Return(expectedErr)

		getters := []pipeline.GetterFunc{
			func() (interface{}, error) {
				return expectedValue, nil
			},
		}

		err := pipeline.ResolveParam(param, getters)
		require.Equal(t, expectedErr, err)
	})
}
