package models

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func Test_ParseCBOR(t *testing.T) {

	tests := []struct {
		name        string
		in          string
		want        JSON
		wantErrored bool
	}{
		{
			"hello world",
			`0xbf6375726c781a68747470733a2f2f657468657270726963652e636f6d2f61706964706174689f66726563656e7463757364ffff`,
			jsonMustUnmarshal(`{"path":["recent","usd"],"url":"https://etherprice.com/api"}`),
			false,
		},
		{
			"trailing empty bytes",
			`0xbf6375726c781a68747470733a2f2f657468657270726963652e636f6d2f61706964706174689f66726563656e7463757364ffff000000`,
			jsonMustUnmarshal(`{"path":["recent","usd"],"url":"https://etherprice.com/api"}`),
			false,
		},
		{
			"nested maps",
			`0xbf657461736b739f6868747470706f7374ff66706172616d73bf636d73676f68656c6c6f5f636861696e6c696e6b6375726c75687474703a2f2f6c6f63616c686f73743a36363930ffff`,
			jsonMustUnmarshal(`{"params":{"msg":"hello_chainlink","url":"http://localhost:6690"},"tasks":["httppost"]}`),
			false,
		},
		{
			"missing initial start map marker",
			`0x636B65796576616C7565ff`,
			jsonMustUnmarshal(`{"key":"value"}`),
			false,
		},
		{"empty object", `0xa0`, jsonMustUnmarshal(`{}`), false},
		{"empty string", `0x`, JSON{}, false},
		{"invalid CBOR", `0xff`, JSON{}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := hexutil.Decode(test.in)
			assert.NoError(t, err)

			json, err := ParseCBOR(b)
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
			"empty",
			[]byte{},
			[]byte{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, autoAddMapDelimiters(test.in))
		})
	}
}

func jsonMustUnmarshal(in string) JSON {
	var j JSON
	err := json.Unmarshal([]byte(in), &j)
	if err != nil {
		log.Panicf("Failed to unmarshal '%s'", in)
	}
	return j
}
