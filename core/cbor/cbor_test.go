package cbor

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func Test_ParseCBOR(t *testing.T) {
	t.Parallel()

	address, err := utils.TryParseHex("0x8bd112d3f8f92e41c861939545ad387307af9703")
	require.NoError(t, err)

	tests := []struct {
		name        string
		in          string
		want        interface{}
		wantErrored bool
	}{
		{
			"hello world",
			`0xbf6375726c781a68747470733a2f2f657468657270726963652e636f6d2f61706964706174689f66726563656e7463757364ffff`,
			jsonMustUnmarshal(t, `{"path":["recent","usd"],"url":"https://etherprice.com/api"}`),
			false,
		},
		{
			"trailing empty bytes",
			`0xbf6375726c781a68747470733a2f2f657468657270726963652e636f6d2f61706964706174689f66726563656e7463757364ffff000000`,
			jsonMustUnmarshal(t, `{"path":["recent","usd"],"url":"https://etherprice.com/api"}`),
			false,
		},
		{
			"nested maps",
			`0xbf657461736b739f6868747470706f7374ff66706172616d73bf636d73676f68656c6c6f5f636861696e6c696e6b6375726c75687474703a2f2f6c6f63616c686f73743a36363930ffff`,
			jsonMustUnmarshal(t, `{"params":{"msg":"hello_chainlink","url":"http://localhost:6690"},"tasks":["httppost"]}`),
			false,
		},
		{
			"missing initial start map marker",
			`0x636B65796576616C7565ff`,
			jsonMustUnmarshal(t, `{"key":"value"}`),
			false,
		},
		{
			"with address encoded",
			`0x6d72656d6f7465436861696e4964186a6e6c69627261727956657273696f6e016f636f6e747261637441646472657373548bd112d3f8f92e41c861939545ad387307af97036d636f6e6669726d6174696f6e730a68626c6f636b4e756d69307831336261626264`,
			map[string]interface{}{
				"blockNum":        "0x13babbd",
				"confirmations":   uint64(10),
				"contractAddress": address,
				"libraryVersion":  uint64(1),
				"remoteChainId":   uint64(106),
			},
			false,
		},
		{
			"bignums",
			"0x" +
				"bf" + // map(*)
				"67" + // text(7)
				"6269676e756d73" + // "bignums"
				"9f" + // array(*)
				"c2" + // tag(2) == unsigned bignum
				"5820" + // bytes(32)
				"0000000000000000000000000000000000000000000000010000000000000000" +
				// int(18446744073709551616)
				"c2" + // tag(2) == unsigned bignum
				"5820" + // bytes(32)
				"4000000000000000000000000000000000000000000000000000000000000000" +
				// int(28948022309329048855892746252171976963317496166410141009864396001978282409984)
				"c3" + // tag(3) == signed bignum
				"5820" + // bytes(32)
				"0000000000000000000000000000000000000000000000010000000000000000" +
				// int(18446744073709551616)
				"c3" + // tag(3) == signed bignum
				"5820" + // bytes(32)
				"3fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
				// int(28948022309329048855892746252171976963317496166410141009864396001978282409983)
				"ff" + // primitive(*)
				"ff", // primitive(*)
			map[string]interface{}{
				"bignums": []interface{}{
					testutils.MustParseBigInt(t, "18446744073709551616"),
					testutils.MustParseBigInt(t, "28948022309329048855892746252171976963317496166410141009864396001978282409984"),
					testutils.MustParseBigInt(t, "-18446744073709551617"),
					testutils.MustParseBigInt(t, "-28948022309329048855892746252171976963317496166410141009864396001978282409984"),
				},
			},
			false,
		},
		{
			"bignums",
			"0x" +
				"67" + // text(7)
				"6269676e756d73" + // "bignums"
				"9f" + // array(*)
				"c2" + // tag(2) == unsigned bignum
				"5820" + // bytes(32)
				"0000000000000000000000000000000000000000000000010000000000000000" +
				// int(18446744073709551616)
				"c2" + // tag(2) == unsigned bignum
				"5820" + // bytes(32)
				"4000000000000000000000000000000000000000000000000000000000000000" +
				// int(28948022309329048855892746252171976963317496166410141009864396001978282409984)
				"c3" + // tag(3) == signed bignum
				"5820" + // bytes(32)
				"0000000000000000000000000000000000000000000000010000000000000000" +
				// int(18446744073709551616)
				"c3" + // tag(3) == signed bignum
				"5820" + // bytes(32)
				"3fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
				// int(28948022309329048855892746252171976963317496166410141009864396001978282409983)
				"ff", // primitive(*)
			map[string]interface{}{
				"bignums": []interface{}{
					testutils.MustParseBigInt(t, "18446744073709551616"),
					testutils.MustParseBigInt(t, "28948022309329048855892746252171976963317496166410141009864396001978282409984"),
					testutils.MustParseBigInt(t, "-18446744073709551617"),
					testutils.MustParseBigInt(t, "-28948022309329048855892746252171976963317496166410141009864396001978282409984"),
				},
			},
			false,
		},
		{"empty object", `0xa0`, jsonMustUnmarshal(t, `{}`), false},
		{"empty string", `0x`, jsonMustUnmarshal(t, `{}`), false},
		{"invalid CBOR", `0xff`, jsonMustUnmarshal(t, `{}`), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := hexutil.Decode(test.in)
			assert.NoError(t, err)

			json, err := ParseDietCBOR(b)
			if test.wantErrored {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, json)
			}
		})
	}
}

func Test_autoAddMapDelimiters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []byte
		want []byte
	}{
		{
			"map(0)",
			hexutil.MustDecode("0xA0"),
			hexutil.MustDecode("0xA0"),
		},
		{
			`map(1) {"key":"value"}`,
			hexutil.MustDecode("0xA1636B65796576616C7565"),
			hexutil.MustDecode("0xA1636B65796576616C7565"),
		},
		{
			"array(0)",
			hexutil.MustDecode("0x80"),
			hexutil.MustDecode("0x80"),
		},
		{
			`map(*) {"key":"value"}`,
			hexutil.MustDecode("0xbf636B65796576616C7565ff"),
			hexutil.MustDecode("0xbf636B65796576616C7565ff"),
		},
		{
			`map(*) {"key":"value"} missing open delimiter`,
			hexutil.MustDecode("0x636B65796576616C7565ff"),
			hexutil.MustDecode("0xbf636B65796576616C7565ffff"),
		},
		{
			`map(*) {"key":"value"} missing closing delimiter`,
			hexutil.MustDecode("0xbf636B65796576616C7565"),
			hexutil.MustDecode("0xbf636B65796576616C7565"),
		},
		{
			`map(*) {"key":"value"} missing both delimiters`,
			hexutil.MustDecode("0x636B65796576616C7565"),
			hexutil.MustDecode("0xbf636B65796576616C7565ff"),
		},
		{
			"empty input adds delimiters",
			[]byte{},
			[]byte{0xbf, 0xff},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, autoAddMapDelimiters(test.in))
		})
	}
}

func jsonMustUnmarshal(t *testing.T, in string) interface{} {
	var j interface{}
	err := json.Unmarshal([]byte(in), &j)
	require.NoError(t, err)
	return j
}

func TestCoerceInterfaceMapToStringMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{"empty map", map[interface{}]interface{}{}, map[string]interface{}{}},
		{"simple map", map[interface{}]interface{}{"key": "value"}, map[string]interface{}{"key": "value"}},
		{"int map", map[int]interface{}{1: "value"}, map[int]interface{}{1: "value"}},
		{
			"nested string map map",
			map[string]interface{}{"key": map[interface{}]interface{}{"nk": "nv"}},
			map[string]interface{}{"key": map[string]interface{}{"nk": "nv"}},
		},
		{
			"nested map map",
			map[interface{}]interface{}{"key": map[interface{}]interface{}{"nk": "nv"}},
			map[string]interface{}{"key": map[string]interface{}{"nk": "nv"}},
		},
		{
			"nested map array",
			map[interface{}]interface{}{"key": []interface{}{1, "value"}},
			map[string]interface{}{"key": []interface{}{1, "value"}},
		},
		{"empty array", []interface{}{}, []interface{}{}},
		{"simple array", []interface{}{1, "value"}, []interface{}{1, "value"}},
		{
			"nested array map",
			[]interface{}{map[interface{}]interface{}{"key": map[interface{}]interface{}{"nk": "nv"}}},
			[]interface{}{map[string]interface{}{"key": map[string]interface{}{"nk": "nv"}}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decoded, err := CoerceInterfaceMapToStringMap(test.input)
			require.NoError(t, err)
			assert.True(t, reflect.DeepEqual(test.want, decoded))
		})
	}
}

func TestCoerceInterfaceMapToStringMap_BadInputs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input interface{}
	}{
		{"error map", map[interface{}]interface{}{1: "value"}},
		{"error array", []interface{}{map[interface{}]interface{}{1: "value"}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := CoerceInterfaceMapToStringMap(test.input)
			assert.Error(t, err)
		})
	}
}

func TestJSON_CBOR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
	}{
		{"empty object", jsonMustUnmarshal(t, `{}`)},
		{"array", jsonMustUnmarshal(t, `[1,2,3,4]`)},
		{
			"basic object",
			jsonMustUnmarshal(t, `{"path":["recent","usd"],"url":"https://etherprice.com/api"}`),
		},
		{
			"complex object",
			jsonMustUnmarshal(t, `{"a":{"1":[{"b":"free"},{"c":"more"},{"d":["less", {"nesting":{"4":"life"}}]}]}}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encoded := mustMarshal(t, test.in)

			var decoded interface{}
			err := cbor.Unmarshal(encoded, &decoded)
			require.NoError(t, err)

			decoded, err = CoerceInterfaceMapToStringMap(decoded)
			require.NoError(t, err)
			assert.True(t, reflect.DeepEqual(test.in, decoded))
		})
	}
}

// mustMarshal returns a bytes array of the JSON map or array encoded to CBOR.
func mustMarshal(t *testing.T, j interface{}) []byte {
	switch v := j.(type) {
	case map[string]interface{}, []interface{}, nil:
		b, err := cbor.Marshal(v)
		if err != nil {
			t.Fatalf("failed to marshal CBOR: %v", err)
		}
		return b
	default:
		t.Fatalf("unable to coerce JSON to CBOR for type %T", v)
		return nil
	}
}
