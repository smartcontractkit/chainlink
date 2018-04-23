package models_test

import (
	"encoding/json"
	"net/url"
	"reflect"
	"testing"
	"time"

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
		{"hello world",
			`0xbf6375726c781a68747470733a2f2f657468657270726963652e636f6d2f61706964706174689f66726563656e7463757364ffff`,
			cltest.JSONFromString(`{"path":["recent","usd"],"url":"https://etherprice.com/api"}`),
			false},
		{"trailing empty bytes",
			`0xbf6375726c781a68747470733a2f2f657468657270726963652e636f6d2f61706964706174689f66726563656e7463757364ffff000000`,
			cltest.JSONFromString(`{"path":["recent","usd"],"url":"https://etherprice.com/api"}`),
			false},
		{"empty object", `a0`, cltest.JSONFromString(`{}`), false},
		{"empty string", ``, models.JSON{}, true},
		{"invalid CBOR", `ff`, models.JSON{}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := utils.HexToBytes(test.in)
			assert.Nil(t, err)

			json, err := models.ParseCBOR(b)
			assert.Equal(t, test.want, json)
			assert.Equal(t, test.wantErrored, (err != nil))
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
			assert.Nil(t, err)

			var decoded interface{}
			cbor := codec.NewDecoderBytes(encoded, new(codec.CborHandle))
			assert.Nil(t, cbor.Decode(&decoded))

			decoded = coerceInterfaceMapToStringMap(decoded)

			assert.True(t, reflect.DeepEqual(test.in.Value(), decoded))
		})
	}
}

func coerceInterfaceMapToStringMap(in interface{}) interface{} {
	switch typed := in.(type) {
	case map[interface{}]interface{}:
		m := map[string]interface{}{}
		for k, v := range typed {
			m[k.(string)] = coerceInterfaceMapToStringMap(v)
		}
		return m
	case []interface{}:
		r := make([]interface{}, len(typed))
		for i, v := range typed {
			r[i] = coerceInterfaceMapToStringMap(v)
		}
		return r
	default:
		return in
	}
}

func TestWebURL_UnmarshalJSON_Error(t *testing.T) {
	t.Parallel()
	j := []byte(`"NotAUrl"`)
	wurl := &models.WebURL{}
	err := json.Unmarshal(j, wurl)
	assert.NotNil(t, err)
}

func TestWebURL_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	j := []byte(`"http://www.duckduckgo.com"`)
	wurl := &models.WebURL{}
	err := json.Unmarshal(j, wurl)
	assert.Nil(t, err)
}

func TestWebURL_MarshalJSON(t *testing.T) {
	t.Parallel()

	str := "http://www.duckduckgo.com"
	parsed, err := url.ParseRequestURI(str)
	assert.Nil(t, err)
	wurl := &models.WebURL{URL: parsed}
	b, err := json.Marshal(wurl)
	assert.Nil(t, err)
	assert.Equal(t, `"`+str+`"`, string(b))
}

func TestWebURL_String_HasURL(t *testing.T) {
	t.Parallel()

	u, _ := url.Parse("http://www.duckduckgo.com")
	w := models.WebURL{
		URL: u,
	}

	assert.Equal(t, "http://www.duckduckgo.com", w.String())
}

func TestWebURL_String_HasNilURL(t *testing.T) {
	t.Parallel()

	w := models.WebURL{}

	assert.Equal(t, "", w.String())
}

func TestTimeDurationFromNow(t *testing.T) {
	t.Parallel()
	future := models.Time{Time: time.Now().Add(time.Second)}
	duration := future.DurationFromNow()
	assert.True(t, 0 < duration)
}
