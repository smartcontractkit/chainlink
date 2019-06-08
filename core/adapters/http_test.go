package adapters_test

import (
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
)

func leanStore() *store.Store {
	return &store.Store{Config: store.NewConfig()}
}

func TestHttpAdapters_NotAUrlError(t *testing.T) {
	t.Parallel()

	store := leanStore()
	tests := []struct {
		name    string
		adapter adapters.BaseAdapter
	}{
		{"HTTPGet", &adapters.HTTPGet{URL: cltest.WebURL(t, "NotAURL")}},
		{"HTTPPost", &adapters.HTTPPost{URL: cltest.WebURL(t, "NotAURL")}},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			result := test.adapter.Perform(models.RunResult{}, store)
			assert.Equal(t, models.JSON{}, result.Data)
			assert.True(t, result.HasError())
		})
	}
}

func TestHTTPGet_Perform(t *testing.T) {
	t.Parallel()

	store := leanStore()
	cases := []struct {
		name        string
		status      int
		want        string
		wantErrored bool
		response    string
		headers     http.Header
		queryParams adapters.QueryParameters
	}{
		{"success", 200, "results!", false, `results!`, nil, nil},
		{"success but error in body", 200, `{"error": "results!"}`, false, `{"error": "results!"}`, nil, nil},
		{"success with HTML", 200, `<html>results!</html>`, false, `<html>results!</html>`, nil, nil},
		{"success with headers", 200, "results!", false, `results!`,
			http.Header{
				"Key1": []string{"value"},
				"Key2": []string{"value", "value"},
			}, nil},
		{"not found", 400, "inputValue", true, `<html>so bad</html>`, nil, nil},
		{"server error", 400, "inputValue", true, `Invalid request`, nil, nil},
		{"success with params", 200, "results!", false, `results!`, nil,
			adapters.QueryParameters{
				"Key1": []string{"value0"},
				"Key2": []string{"value1", "value2"},
			}},
	}

	for _, tt := range cases {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			input := cltest.RunResultWithResult("inputValue")
			mock, cleanup := cltest.NewHTTPMockServer(t, test.status, "GET", test.response,
				func(header http.Header, body string) {
					assert.Equal(t, ``, body)
					for key, values := range test.headers {
						assert.Equal(t, values, header[key])
					}
				})
			defer cleanup()

			hga := adapters.HTTPGet{
				URL: cltest.WebURL(t, mock.URL),
				Headers: test.headers,
				QueryParams: test.queryParams,
			}
			for key, _ := range hga.QueryParams {
				assert.Equal(t, test.queryParams[key], hga.QueryParams[key])
			}

			result := hga.Perform(input, store)

			val, err := result.ResultString()
			assert.NoError(t, err)
			assert.Equal(t, test.want, val)
			assert.Equal(t, test.wantErrored, result.HasError())
			assert.Equal(t, false, result.Status.PendingBridge())
		})
	}
}

func TestHTTP_TooLarge(t *testing.T) {
	cfg := store.NewConfig()
	cfg.Set("DEFAULT_HTTP_LIMIT", "1")
	store := &store.Store{Config: cfg}

	tests := []struct {
		verb    string
		factory func(models.WebURL) adapters.BaseAdapter
	}{
		{"GET", func(url models.WebURL) adapters.BaseAdapter { return &adapters.HTTPGet{URL: url} }},
		{"POST", func(url models.WebURL) adapters.BaseAdapter { return &adapters.HTTPPost{URL: url} }},
	}
	for _, test := range tests {
		t.Run(test.verb, func(t *testing.T) {
			input := cltest.RunResultWithResult("inputValue")
			largePayload := "12"
			mock, cleanup := cltest.NewHTTPMockServer(t, 200, test.verb, largePayload)
			defer cleanup()

			hga := test.factory(cltest.WebURL(t, mock.URL))
			result := hga.Perform(input, store)

			assert.Equal(t, true, result.HasError())
			assert.Equal(t, "HTTP request too large, must be less than 1 bytes", result.Error())
			assert.Equal(t, "inputValue", result.Result().String())
		})
	}
}

func TestHttpPost_Perform(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		status      int
		want        string
		wantErrored bool
		response    string
		headers     http.Header
		queryParams adapters.QueryParameters
	}{
		{"success", 200, "results!", false, `results!`, nil, nil},
		{"success but error in body", 200, `{"error": "results!"}`, false, `{"error": "results!"}`, nil, nil},
		{"success with HTML", 200, `<html>results!</html>`, false, `<html>results!</html>`, nil, nil},
		{"success with headers", 200, "results!", false, `results!`,
			http.Header{
				"Key1": []string{"value"},
				"Key2": []string{"value", "value"},
			}, nil},
		{"not found", 400, "inputVal", true, `<html>so bad</html>`, nil, nil},
		{"server error", 500, "inputVal", true, `big error`, nil, nil},
		{"success with params", 200, "results!", false, `results!`, nil,
			adapters.QueryParameters{
				"Key1": []string{"value"},
				"Key2": []string{"value", "value"},
			}},
	}

	for _, tt := range cases {
		test := tt
		t.Run(test.name, func(t *testing.T) {
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

			hpa := adapters.HTTPPost{
				URL: cltest.WebURL(t, mock.URL),
				Headers: test.headers,
				QueryParams: test.queryParams,
			}
			for key, _ := range hpa.QueryParams {
				assert.Equal(t, test.queryParams[key], hpa.QueryParams[key])
			}

			result := hpa.Perform(input, leanStore())

			val := result.Result()
			assert.Equal(t, test.want, val.String())
			assert.Equal(t, true, val.Exists())
			assert.Equal(t, test.wantErrored, result.HasError())
			assert.Equal(t, false, result.Status.PendingBridge())
		})
	}
}
