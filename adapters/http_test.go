package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestHttpNotAUrlError(t *testing.T) {
	tests := []struct {
		name    string
		adapter adapters.Adapter
	}{
		{"HttpGet", &adapters.HttpGet{URL: cltest.MustParseWebURL("NotAURL")}},
		{"HttpPost", &adapters.HttpGet{URL: cltest.MustParseWebURL("NotAURL")}},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := test.adapter.Perform(models.RunResult{}, nil)
			assert.Equal(t, models.JSON{}, result.Data)
			assert.NotNil(t, result.Error)
		})
	}
}

func TestHttpGetAdapterPerform(t *testing.T) {
	cases := []struct {
		name        string
		status      int
		want        string
		wantExists  bool
		wantErrored bool
		response    string
	}{
		{"success", 200, "so good", true, false, `so good`},
		{"success but error in body", 200, `{"error": "so good"}`, true, false, `{"error": "so good"}`},
		{"success with HTML", 200, `<html>so good</html>`, true, false, `<html>so good</html>`},
		{"not found", 400, "", false, true, `<html>so bad</html>`},
		{"server error", 400, "", false, true, `Invalid request`},
	}

	for _, tt := range cases {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := models.RunResultWithValue("unused")
			mock, cleanup := cltest.NewHTTPMockServer(t, test.status, "GET", test.response,
				func(body string) { assert.Equal(t, ``, body) })
			defer cleanup()

			hga := adapters.HttpGet{URL: cltest.MustParseWebURL(mock.URL)}
			result := hga.Perform(input, nil)

			val, err := result.Get("value")
			assert.Nil(t, err)
			assert.Equal(t, test.want, val.String())
			assert.Equal(t, test.wantExists, val.Exists())
			assert.Equal(t, test.wantErrored, result.HasError())
			assert.Equal(t, false, result.Pending)
		})
	}
}

func TestHttpPostAdapterPerform(t *testing.T) {
	cases := []struct {
		name        string
		status      int
		want        string
		wantExists  bool
		wantErrored bool
		response    string
	}{
		{"success", 200, "so meta", true, false, `so meta`},
		{"success but error in body", 200, `{"error": "so meta"}`, true, false, `{"error": "so meta"}`},
		{"success with HTML", 200, `<html>so meta</html>`, true, false, `<html>so meta</html>`},
		{"not found", 400, "", false, true, `<html>so bad</html>`},
		{"server error", 500, "", false, true, `big error`},
	}

	for _, tt := range cases {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := models.RunResultWithValue("modern")
			wantedBody := `{"value":"modern"}`
			mock, cleanup := cltest.NewHTTPMockServer(t, test.status, "POST", test.response,
				func(body string) { assert.Equal(t, wantedBody, body) })
			defer cleanup()

			hpa := adapters.HttpPost{URL: cltest.MustParseWebURL(mock.URL)}
			result := hpa.Perform(input, nil)

			val, err := result.Get("value")
			assert.Nil(t, err)
			assert.Equal(t, test.want, val.String())
			assert.Equal(t, test.wantExists, val.Exists())
			assert.Equal(t, test.wantErrored, result.HasError())
			assert.Equal(t, false, result.Pending)
		})
	}
}
