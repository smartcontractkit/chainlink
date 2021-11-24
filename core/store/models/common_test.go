package models_test

import (
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"

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
