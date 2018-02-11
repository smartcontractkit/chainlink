package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestHttpNotAUrlError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		adapter adapters.Adapter
	}{
		{"HttpGet", &adapters.HttpGet{URL: cltest.MustParseWebURL("NotAURL")}},
		{"HttpPost", &adapters.HttpGet{URL: cltest.MustParseWebURL("NotAURL")}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.adapter.Perform(models.RunResult{}, nil)
			assert.Equal(t, models.JSON{}, result.Output)
			assert.NotNil(t, result.Error)
		})
	}
}

func TestHttpGetAdapterPerform(t *testing.T) {
	t.Parallel()

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
		{"server error", 400, "", false, true, `Invalid request`},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			input := models.RunResultWithValue("unused")
			mock, cleanup := cltest.NewHTTPMockServer(t, test.status, "GET", test.response,
				func(body string) { assert.Equal(t, ``, body) })
			defer cleanup()

			hga := adapters.HttpGet{URL: cltest.MustParseWebURL(mock.URL)}
			result := hga.Perform(input, store)

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
	t.Parallel()

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
		{"server error", 500, "", false, true, `big error`},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			input := models.RunResultWithValue("modern")
			wantedBody := `{"value":"modern"}`
			mock, cleanup := cltest.NewHTTPMockServer(t, test.status, "POST", test.response,
				func(body string) { assert.Equal(t, wantedBody, body) })
			defer cleanup()

			hpa := adapters.HttpPost{URL: cltest.MustParseWebURL(mock.URL)}
			result := hpa.Perform(input, store)

			val, err := result.Get("value")
			assert.Nil(t, err)
			assert.Equal(t, test.want, val.String())
			assert.Equal(t, test.wantExists, val.Exists())
			assert.Equal(t, test.wantErrored, result.HasError())
			assert.Equal(t, false, result.Pending)
		})
	}
}
