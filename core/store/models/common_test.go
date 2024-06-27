package models_test

import (
	"encoding/json"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"

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

func TestNewInterval(t *testing.T) {
	t.Parallel()

	duration := 33 * time.Second
	interval := models.NewInterval(duration)

	require.Equal(t, duration, interval.Duration())
}

func TestSha256Hash_MarshalJSON_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	hash := models.MustSha256HashFromHex("f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5")
	json, err := hash.MarshalJSON()
	require.NoError(t, err)
	require.NotEmpty(t, json)

	var newHash models.Sha256Hash
	err = newHash.UnmarshalJSON(json)
	require.NoError(t, err)

	require.Equal(t, hash, newHash)
}

func TestSha256Hash_Sha256HashFromHex(t *testing.T) {
	t.Parallel()

	_, err := models.Sha256HashFromHex("abczzz")
	require.Error(t, err)

	_, err = models.Sha256HashFromHex("f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5")
	require.NoError(t, err)

	_, err = models.Sha256HashFromHex("f5bf259689b26f1374e6")
	require.NoError(t, err)
}

func TestSha256Hash_String(t *testing.T) {
	t.Parallel()

	hash := models.MustSha256HashFromHex("f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5")
	assert.Equal(t, "f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5", hash.String())
}

func TestSha256Hash_Scan_Value(t *testing.T) {
	t.Parallel()

	hash := models.MustSha256HashFromHex("f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5")
	val, err := hash.Value()
	require.NoError(t, err)

	var newHash models.Sha256Hash
	err = newHash.Scan(val)
	require.NoError(t, err)

	require.Equal(t, hash, newHash)
}

func TestAddressCollection_Scan_Value(t *testing.T) {
	t.Parallel()

	ac := models.AddressCollection{
		common.HexToAddress(strings.Repeat("AA", 20)),
		common.HexToAddress(strings.Repeat("BB", 20)),
	}

	val, err := ac.Value()
	require.NoError(t, err)

	var acNew models.AddressCollection
	err = acNew.Scan(val)
	require.NoError(t, err)

	require.Equal(t, ac, acNew)
}

func TestAddressCollection_ToStrings(t *testing.T) {
	t.Parallel()

	hex1 := "0xaAaAaAaaAaAaAaaAaAAAAAAAAaaaAaAaAaaAaaAa"
	hex2 := "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"

	ac := models.AddressCollection{
		common.HexToAddress(hex1),
		common.HexToAddress(hex2),
	}

	acStrings := ac.ToStrings()
	require.Len(t, acStrings, 2)
	require.Equal(t, hex1, acStrings[0])
	require.Equal(t, hex2, acStrings[1])
}

func TestInterval_IsZero(t *testing.T) {
	t.Parallel()

	i := models.NewInterval(0)
	require.NotNil(t, i)
	require.True(t, i.IsZero())

	i = models.NewInterval(1)
	require.NotNil(t, i)
	require.False(t, i.IsZero())
}

func TestInterval_Scan_Value(t *testing.T) {
	t.Parallel()

	i := models.NewInterval(100)
	require.NotNil(t, i)

	val, err := i.Value()
	require.NoError(t, err)

	iNew := models.NewInterval(0)
	err = iNew.Scan(val)
	require.NoError(t, err)

	require.Equal(t, i, iNew)
}

func TestInterval_MarshalText_UnmarshalText(t *testing.T) {
	t.Parallel()

	i := models.NewInterval(100)
	require.NotNil(t, i)

	txt, err := i.MarshalText()
	require.NoError(t, err)

	iNew := models.NewInterval(0)
	err = iNew.UnmarshalText(txt)
	require.NoError(t, err)

	require.Equal(t, i, iNew)
}

func TestDuration_Scan_Value(t *testing.T) {
	t.Parallel()

	d := models.MustMakeDuration(100)
	require.NotNil(t, d)

	val, err := d.Value()
	require.NoError(t, err)

	dNew := models.MustMakeDuration(0)
	err = dNew.Scan(val)
	require.NoError(t, err)

	require.Equal(t, d, dNew)
}

func TestDuration_MarshalJSON_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	d := models.MustMakeDuration(100)
	require.NotNil(t, d)

	json, err := d.MarshalJSON()
	require.NoError(t, err)

	dNew := models.MustMakeDuration(0)
	err = dNew.UnmarshalJSON(json)
	require.NoError(t, err)

	require.Equal(t, d, dNew)
}

func TestDuration_MakeDurationFromString(t *testing.T) {
	t.Parallel()

	d, err := models.ParseDuration("1s")
	require.NoError(t, err)
	require.Equal(t, 1*time.Second, d.Duration())

	_, err = models.ParseDuration("xyz")
	require.Error(t, err)
}

func TestWebURL_Scan_Value(t *testing.T) {
	t.Parallel()

	u, err := url.Parse("https://chain.link")
	require.NoError(t, err)

	w := models.WebURL(*u)

	val, err := w.Value()
	require.NoError(t, err)

	var wNew models.WebURL
	err = wNew.Scan(val)
	require.NoError(t, err)

	require.Equal(t, w, wNew)
}

func TestJSON_Scan_Value(t *testing.T) {
	t.Parallel()

	js, err := models.ParseJSON([]byte(`{"foo":123}`))
	require.NoError(t, err)

	val, err := js.Value()
	require.NoError(t, err)

	var jsNew models.JSON
	err = jsNew.Scan(val)
	require.NoError(t, err)

	require.Equal(t, js, jsNew)
}

func TestJSON_Bytes(t *testing.T) {
	t.Parallel()

	jsBytes := []byte(`{"foo":123}`)

	js, err := models.ParseJSON(jsBytes)
	require.NoError(t, err)

	require.Equal(t, jsBytes, js.Bytes())
}

func TestJSON_MarshalJSON(t *testing.T) {
	t.Parallel()

	jsBytes := []byte(`{"foo":123}`)

	js, err := models.ParseJSON(jsBytes)
	require.NoError(t, err)

	bs, err := js.MarshalJSON()
	require.NoError(t, err)

	require.Equal(t, jsBytes, bs)
}

func TestJSON_UnmarshalTOML(t *testing.T) {
	t.Parallel()

	jsBytes := []byte(`{"foo":123}`)

	var js models.JSON
	err := js.UnmarshalTOML(jsBytes)
	require.NoError(t, err)
	require.Equal(t, jsBytes, js.Bytes())

	err = js.UnmarshalTOML(string(jsBytes))
	require.NoError(t, err)
	require.Equal(t, jsBytes, js.Bytes())
}
