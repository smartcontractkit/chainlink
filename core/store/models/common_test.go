package models_test

import (
	"encoding/json"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSON_Merge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		original  string
		input     string
		want      string
		wantError bool
	}{
		{
			"new field",
			`{"value":"OLD","other":1}`,
			`{"extra":"fields"}`,
			`{"value":"OLD","other":1,"extra":"fields"}`,
			false,
		},
		{
			"overwritting fields",
			`{"value":"OLD","other":1}`,
			`{"value":["new","new"],"extra":2}`,
			`{"value":["new","new"],"other":1,"extra":2}`,
			false,
		},
		{
			"nested JSON",
			`{"value":"OLD","other":1}`,
			`{"extra":{"fields": ["more", 1]}}`,
			`{"value":"OLD","other":1,"extra":{"fields":["more",1]}}`,
			false,
		},
		{
			"empty JSON",
			`{"value":"OLD","other":1}`,
			`{}`,
			`{"value":"OLD","other":1}`,
			false,
		},
		{
			"null values",
			`{"value":"OLD","other":1}`,
			`{"value":null}`,
			`{"value":null,"other":1}`,
			false,
		},
		{
			"string",
			`"string"`,
			`{}`,
			"",
			true,
		},
		{
			"array",
			`["a1"]`,
			`{"value": null}`,
			"",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			j1 := cltest.JSONFromString(t, test.original)
			j2 := cltest.JSONFromString(t, test.input)

			merged, err := models.Merge(j1, j2)
			if test.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.JSONEq(t, test.want, merged.String())
				assert.JSONEq(t, test.original, j1.String())
			}
		})
	}
}

func TestJSON_MergeNull(t *testing.T) {
	merged, err := models.Merge(models.JSON{}, models.JSON{})
	require.NoError(t, err)
	assert.Equal(t, `{}`, merged.String())
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
		{"basic", `{"num": 100}`, cltest.JSONFromString(t, `{"num": 100}`), false},
		{"empty string", ``, models.JSON{}, false},
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
		{"escaped quote", "a", `"2"`, false, `{"a":"\"2\""}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			json := cltest.JSONFromString(t, `{"a":"1"}`)

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
			json := cltest.JSONFromString(t, `{"a":"1","b":2}`)

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
		{"array", cltest.JSONFromString(t, `[1,2,3,4]`)},
		{
			"hello world",
			cltest.JSONFromString(t, `{"path":["recent","usd"],"url":"https://etherprice.com/api"}`),
		},
		{
			"complex object",
			cltest.JSONFromString(t, `{"a":{"1":[{"b":"free"},{"c":"more"},{"d":["less", {"nesting":{"4":"life"}}]}]}}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encoded, err := test.in.CBOR()
			assert.NoError(t, err)

			var decoded interface{}
			err = cbor.Unmarshal(encoded, &decoded)

			assert.NoError(t, err)

			decoded, err = utils.CoerceInterfaceMapToStringMap(decoded)
			assert.NoError(t, err)
			assert.True(t, reflect.DeepEqual(test.in.Result.Value(), decoded))
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

func TestAnyTime_UnmarshalJSON_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  time.Time
	}{
		{"unix string", `"1529445491"`, time.Unix(1529445491, 0).UTC()},
		{"unix int", `1529445491`, time.Unix(1529445491, 0).UTC()},
		{"iso8601 time", `"2018-06-19T22:17:19Z"`, time.Unix(1529446639, 0).UTC()},
		{"iso8601 date", `"2018-06-19"`, time.Unix(1529366400, 0).UTC()},
		{"iso8601 year", `"2018"`, time.Unix(1514764800, 0).UTC()},
		{"null", `null`, time.Time{}},
		{"empty", `""`, time.Time{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual models.AnyTime
			err := json.Unmarshal([]byte(test.input), &actual)
			require.NoError(t, err)
			assert.Equal(t, test.want, actual.Time)
		})
	}
}

func TestAnyTime_UnmarshalJSON_Error(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"invalid string", `"1000h"`},
		{"float", `"1000.123"`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual models.AnyTime
			err := json.Unmarshal([]byte(test.input), &actual)
			assert.Error(t, err)
		})
	}
}

func TestAnyTime_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input models.AnyTime
		want  string
	}{
		{"valid", models.NewAnyTime(time.Unix(1529446639, 0).UTC()), `"2018-06-19T22:17:19Z"`},
		{"invalid", models.AnyTime{}, `null`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := json.Marshal(&test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.want, string(b))
		})
	}
}

func TestDuration_MarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input models.Duration
		want  string
	}{
		{"zero", models.MustMakeDuration(0), `"0s"`},
		{"one second", models.MustMakeDuration(time.Second), `"1s"`},
		{"one minute", models.MustMakeDuration(time.Minute), `"1m0s"`},
		{"one hour", models.MustMakeDuration(time.Hour), `"1h0m0s"`},
		{"one hour thirty minutes", models.MustMakeDuration(time.Hour + 30*time.Minute), `"1h30m0s"`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := json.Marshal(&test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.want, string(b))
		})
	}
}

func TestCron_UnmarshalJSON_Success(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{"valid 5-field cron", `"CRON_TZ=UTC 0 0/5 * * *"`},
		{"valid 6-field cron", `"CRON_TZ=UTC 30 0 0/5 * * *"`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual models.Cron
			err := json.Unmarshal([]byte(test.input), &actual)
			assert.NoError(t, err)
		})
	}
}

func TestCron_UnmarshalJSON_Invalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantError string
	}{
		{"5-field cron without time zone", `"0 0/5 * * *"`, "Cron: specs must specify a time zone using CRON_TZ, e.g. 'CRON_TZ=UTC 5 * * * *'"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual models.Cron
			err := json.Unmarshal([]byte(test.input), &actual)
			assert.EqualError(t, err, test.wantError)
		})
	}
}
