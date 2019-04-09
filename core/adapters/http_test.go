package adapters_test

import (
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
)

func TestHttpAdapters_NotAUrlError(t *testing.T) {
	tests := []struct {
		name    string
		adapter adapters.BaseAdapter
	}{
		{"HTTPGet", &adapters.HTTPGet{URL: cltest.WebURL("NotAURL")}},
		{"HTTPPost", &adapters.HTTPPost{URL: cltest.WebURL("NotAURL")}},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := test.adapter.Perform(models.RunResult{}, nil)
			assert.Equal(t, models.JSON{}, result.Data)
			assert.True(t, result.HasError())
		})
	}
}

func TestHttpGet_Perform(t *testing.T) {
	cases := []struct {
		name        string
		status      int
		want        string
		wantErrored bool
		response    string
		headers     http.Header
	}{
		{"success", 200, "results!", false, `results!`, nil},
		{"success but error in body", 200, `{"error": "results!"}`, false, `{"error": "results!"}`, nil},
		{"success with HTML", 200, `<html>results!</html>`, false, `<html>results!</html>`, nil},
		{"success with headers", 200, "results!", false, `results!`,
			http.Header{
				"Key1": []string{"value"},
				"Key2": []string{"value", "value"},
			}},
		{"not found", 400, "inputValue", true, `<html>so bad</html>`, nil},
		{"server error", 400, "inputValue", true, `Invalid request`, nil},
	}

	for _, tt := range cases {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := cltest.RunResultWithResult("inputValue")
			mock, cleanup := cltest.NewHTTPMockServer(t, test.status, "GET", test.response,
				func(header http.Header, body string) {
					assert.Equal(t, ``, body)
					for key, values := range test.headers {
						assert.Equal(t, values, header[key])
					}
				})
			defer cleanup()

			hga := adapters.HTTPGet{URL: cltest.WebURL(mock.URL), Headers: test.headers}
			result := hga.Perform(input, nil)

			val, err := result.ResultString()
			assert.NoError(t, err)
			assert.Equal(t, test.want, val)
			assert.Equal(t, test.wantErrored, result.HasError())
			assert.Equal(t, false, result.Status.PendingBridge())
		})
	}
}

func TestHttpPost_Perform(t *testing.T) {
	cases := []struct {
		name        string
		status      int
		want        string
		wantErrored bool
		response    string
		headers     http.Header
	}{
		{"success", 200, "results!", false, `results!`, nil},
		{"success but error in body", 200, `{"error": "results!"}`, false, `{"error": "results!"}`, nil},
		{"success with HTML", 200, `<html>results!</html>`, false, `<html>results!</html>`, nil},
		{"success with headers", 200, "results!", false, `results!`,
			http.Header{
				"Key1": []string{"value"},
				"Key2": []string{"value", "value"},
			}},
		{"not found", 400, "inputVal", true, `<html>so bad</html>`, nil},
		{"server error", 500, "inputVal", true, `big error`, nil},
	}

	for _, tt := range cases {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := cltest.RunResultWithResult("inputVal")
			wantedBody := `{"result":"inputVal"}`
			mock, cleanup := cltest.NewHTTPMockServer(t, test.status, "POST", test.response,
				func(header http.Header, body string) {
					assert.Equal(t, wantedBody, body)
					for key, values := range test.headers {
						assert.Equal(t, values, header[key])
					}
				})
			defer cleanup()

			hpa := adapters.HTTPPost{URL: cltest.WebURL(mock.URL), Headers: test.headers}
			result := hpa.Perform(input, nil)

			val := result.Result()
			assert.Equal(t, test.want, val.String())
			assert.Equal(t, true, val.Exists())
			assert.Equal(t, test.wantErrored, result.HasError())
			assert.Equal(t, false, result.Status.PendingBridge())
		})
	}
}
