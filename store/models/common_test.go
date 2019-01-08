package models_test

import (
	"encoding/json"
	"math/big"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	"github.com/ugorji/go/codec"
)

func Test_ParseCBOR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		in          string
		want        models.JSON
		wantErrored bool
	}{
		{
			"hello world",
			`0xbf6375726c781a68747470733a2f2f657468657270726963652e636f6d2f61706964706174689f66726563656e7463757364ffff`,
			cltest.JSONFromString(`{"path":["recent","usd"],"url":"https://etherprice.com/api"}`),
			false,
		},
		{
			"trailing empty bytes",
			`0xbf6375726c781a68747470733a2f2f657468657270726963652e636f6d2f61706964706174689f66726563656e7463757364ffff000000`,
			cltest.JSONFromString(`{"path":["recent","usd"],"url":"https://etherprice.com/api"}`),
			false,
		},
		{
			"nested maps",
			`0xbf657461736b739f6868747470706f7374ff66706172616d73bf636d73676f68656c6c6f5f636861696e6c696e6b6375726c75687474703a2f2f6c6f63616c686f73743a36363930ffff`,
			cltest.JSONFromString(`{"params":{"msg":"hello_chainlink","url":"http://localhost:6690"},"tasks":["httppost"]}`),
			false,
		},
		{
			"missing initial start map marker",
			`0x636B65796576616C7565ff`,
			cltest.JSONFromString(`{"key":"value"}`),
			false,
		},
		{
			"missing trailing end map marker",
			`0xbf636B65796576616C7565`,
			cltest.JSONFromString(`{"key":"value"}`),
			false,
		},
		{
			"missing both start and end map marker",
			`0x636B65796576616C7565`,
			cltest.JSONFromString(`{"key":"value"}`),
			false,
		},
		{"empty object", `0xa0`, cltest.JSONFromString(`{}`), false},
		{"empty string", `0x`, models.JSON{}, true},
		{"invalid CBOR", `0xff`, models.JSON{}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := hexutil.Decode(test.in)
			assert.NoError(t, err)

			json, err := models.ParseCBOR(b)
			if test.wantErrored {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, json)
			}
		})
	}
}

func TestJSON_Merge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       string
		want        string
		wantErrored bool
	}{
		{"new field", `{"extra":"fields"}`,
			`{"value":"OLD","other":1,"extra":"fields"}`, false},
		{"overwritting fields", `{"value":["new","new"],"extra":2}`,
			`{"value":["new","new"],"other":1,"extra":2}`, false},
		{"nested JSON", `{"extra":{"fields": ["more", 1]}}`,
			`{"value":"OLD","other":1,"extra":{"fields":["more",1]}}`, false},
		{"empty JSON", `{}`,
			`{"value":"OLD","other":1}`, false},
		{"null values", `{"value":null}`,
			`{"value":null,"other":1}`, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			orig := `{"value":"OLD","other":1}`
			j1 := cltest.JSONFromString(orig)
			j2 := cltest.JSONFromString(test.input)

			merged, err := j1.Merge(j2)
			assert.Equal(t, test.wantErrored, (err != nil))
			assert.JSONEq(t, test.want, merged.String())
			assert.JSONEq(t, orig, j1.String())
		})
	}
}

func TestJSON_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		json        string
		wantErrored bool
	}{
		{"basic", `{"number": 100, "string": "100", "bool": true}`, false},
		{"invalid JSON", `{`, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var j models.JSON
			err := json.Unmarshal([]byte(test.json), &j)
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func TestJSON_ParseJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		in          string
		want        models.JSON
		wantErrored bool
	}{
		{"basic", `{"num": 100}`, cltest.JSONFromString(`{"num": 100}`), false},
		{"empty string", ``, cltest.JSONFromString(`{}`), false},
		{"invalid JSON", `{`, models.JSON{}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			json, err := models.ParseJSON([]byte(test.in))
			assert.Equal(t, test.want, json)
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func TestJSON_Add(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		value   interface{}
		errored bool
		want    string
	}{
		{"adding string", "b", "2", false, `{"a":"1","b":"2"}`},
		{"adding int", "b", 2, false, `{"a":"1","b":2}`},
		{"overriding", "a", "2", false, `{"a":"2"}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			json := cltest.JSONFromString(`{"a":"1"}`)

			json, err := json.Add(test.key, test.value)
			assert.Equal(t, test.errored, (err != nil))
			assert.Equal(t, test.want, json.String())
		})
	}
}

func TestJSON_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		key  string
		want string
	}{
		{"remove existing key", "b", `{"a":"1"}`},
		{"remove non-existing key", "c", `{"a":"1","b":2}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			json := cltest.JSONFromString(`{"a":"1","b":2}`)

			json, err := json.Delete(test.key)

			assert.NoError(t, err)
			assert.Equal(t, test.want, json.String())
		})
	}
}

func TestJSON_CBOR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   models.JSON
	}{
		{"empty object", models.JSON{}},
		{"array", cltest.JSONFromString(`[1,2,3,4]`)},
		{
			"hello world",
			cltest.JSONFromString(`{"path":["recent","usd"],"url":"https://etherprice.com/api"}`),
		},
		{
			"complex object",
			cltest.JSONFromString(`{"a":{"1":[{"b":"free"},{"c":"more"},{"d":["less", {"nesting":{"4":"life"}}]}]}}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encoded, err := test.in.CBOR()
			assert.NoError(t, err)

			var decoded interface{}
			cbor := codec.NewDecoderBytes(encoded, new(codec.CborHandle))
			assert.NoError(t, cbor.Decode(&decoded))

			decoded, err = utils.CoerceInterfaceMapToStringMap(decoded)
			assert.NoError(t, err)
			assert.True(t, reflect.DeepEqual(test.in.Value(), decoded))
		})
	}
}

func TestWebURL_UnmarshalJSON_Error(t *testing.T) {
	t.Parallel()
	j := []byte(`"NotAUrl"`)
	wurl := &models.WebURL{}
	err := json.Unmarshal(j, wurl)
	assert.Error(t, err)
}

func TestWebURL_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	j := []byte(`"http://www.duckduckgo.com"`)
	wurl := &models.WebURL{}
	err := json.Unmarshal(j, wurl)
	assert.NoError(t, err)
}

func TestWebURL_MarshalJSON(t *testing.T) {
	t.Parallel()

	str := "http://www.duckduckgo.com"
	parsed, err := url.ParseRequestURI(str)
	assert.NoError(t, err)
	wurl := models.WebURL(*parsed)
	b, err := json.Marshal(wurl)
	assert.NoError(t, err)
	assert.Equal(t, `"`+str+`"`, string(b))
}

func TestWebURL_String_HasURL(t *testing.T) {
	t.Parallel()

	u, _ := url.Parse("http://www.duckduckgo.com")
	w := models.WebURL(*u)

	assert.Equal(t, "http://www.duckduckgo.com", w.String())
}

func TestWebURL_String_HasNilURL(t *testing.T) {
	t.Parallel()

	w := models.WebURL{}

	assert.Equal(t, "", w.String())
}

func TestTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		errored bool
	}{
		{"unix string", `"1529445491"`, time.Unix(1529445491, 0).UTC(), false},
		{"unix int", `1529445491`, time.Unix(1529445491, 0).UTC(), false},
		{"iso8601 time", `"2018-06-19T22:17:19Z"`, time.Unix(1529446639, 0).UTC(), false},
		{"iso8601 date", `"2018-06-19"`, time.Unix(1529366400, 0).UTC(), false},
		{"iso8601 year", `"2018"`, time.Unix(1514764800, 0).UTC(), false},
		{"invalid string", `"1000h"`, time.Now(), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual models.Time
			err := json.Unmarshal([]byte(test.input), &actual)
			if test.errored {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, actual.Time)
			}
		})
	}
}

func TestTime_DurationFromNow(t *testing.T) {
	t.Parallel()
	future := models.Time{Time: time.Now().Add(time.Second)}
	duration := future.DurationFromNow()
	assert.True(t, 0 < duration)
}

func TestInt_UnmarshalText(t *testing.T) {
	t.Parallel()

	i := &models.Int{}
	tests := []struct {
		name      string
		input     string
		wantError bool
		want      *big.Int
	}{
		{"number", `1234`, false, big.NewInt(1234)},
		{"string", `"1234"`, false, big.NewInt(1234)},
		{"hex number", `0x1234`, false, big.NewInt(4660)},
		{"hex string", `"0x1234"`, false, big.NewInt(4660)},
		{"single quoted", `'1234'`, false, big.NewInt(1234)},
		{"quoted word", `"word"`, true, big.NewInt(0)},
		{"word", `word`, true, big.NewInt(0)},
		{"empty", ``, true, big.NewInt(0)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := i.UnmarshalText([]byte(test.input))
			cltest.AssertError(t, test.wantError, err)
			assert.Equal(t, test.want, i.ToBig())
		})
	}
}
