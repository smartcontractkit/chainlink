package models_test

import (
	"encoding/hex"
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
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

func TestJSON_Keys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want []string
	}{
		{"empty object", "{}", []string{}},
		{"ordered", `{"a":1,"b":1,"c":1}`, []string{"a", "b", "c"}},
		{"unordered", `{"c":1,"a":1,"b":1}`, []string{"a", "b", "c"}},
		{"duplicates", `{"a":1,"a":1,"b":1}`, []string{"a", "b"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			j := cltest.JSONFromString(test.in)
			assert.Equal(t, test.want, j.Keys())
		})
	}
}

func TestJSON_CBOR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   models.JSON
		want string
	}{
		{"empty object", models.JSON{}, "a0"},
		{"hello world",
			cltest.JSONFromString(`{"path":["recent","usd"],"url":"https://etherprice.com/api"}`),
			`a264706174688266726563656e74637573646375726c781a68747470733a2f2f657468657270726963652e636f6d2f617069`},
		{"complex object",
			cltest.JSONFromString(`{"a":{"1":[{"b":"free"},{"c":"more"},{"d":["less", {"nesting":{"4":"life"}}]}]}}`),
			`a16161a1613183a161626466726565a16163646d6f7265a1616482646c657373a1676e657374696e67a16134646c696665`},
		{"unordered keys",
			cltest.JSONFromString(`{"b":0,"a":0}`),
			`a26161fb00000000000000006162fb0000000000000000`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cbor, err := test.in.CBOR()
			assert.Nil(t, err)

			cborHex := hex.EncodeToString(cbor)
			assert.Equal(t, test.want, cborHex)
		})
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
