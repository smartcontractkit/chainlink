package adapters_test

import (
	"encoding/json"
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
				URL:         cltest.WebURL(t, mock.URL),
				Headers:     test.headers,
				QueryParams: test.queryParams,
			}
			assert.Equal(t, test.queryParams, hga.QueryParams)

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
				URL:         cltest.WebURL(t, mock.URL),
				Headers:     test.headers,
				QueryParams: test.queryParams,
			}
			assert.Equal(t, test.queryParams, hpa.QueryParams)

			result := hpa.Perform(input, leanStore())

			val := result.Result()
			assert.Equal(t, test.want, val.String())
			assert.Equal(t, true, val.Exists())
			assert.Equal(t, test.wantErrored, result.HasError())
			assert.Equal(t, false, result.Status.PendingBridge())
		})
	}
}

func TestQueryParameters_Success(t *testing.T) {
	t.Parallel()

	baseUrl := "http://example.com"

	cases := []struct {
		name           string
		queryParams    string
		startingUrl    string
		expectedParams adapters.QueryParameters
		expectedURL    string
	}{
		{"empty",
			`""`,
			baseUrl,
			adapters.QueryParameters{},
			baseUrl,
		},
		{
			"array of params",
			`["firstKey","firstVal","secondKey","secondVal"]`,
			baseUrl,
			adapters.QueryParameters{
				"firstKey":  []string{"firstVal"},
				"secondKey": []string{"secondVal"},
			},
			"http://example.com?firstKey=firstVal&secondKey=secondVal",
		},
		{
			"string of params",
			`"firstKey=firstVal&secondKey=secondVal"`,
			baseUrl,
			adapters.QueryParameters{
				"firstKey":  []string{"firstVal"},
				"secondKey": []string{"secondVal"},
			},
			"http://example.com?firstKey=firstVal&secondKey=secondVal",
		},
		{
			"string has question mark",
			`"?firstKey=firstVal&secondKey=secondVal"`,
			baseUrl,
			adapters.QueryParameters{
				"firstKey":  []string{"firstVal"},
				"secondKey": []string{"secondVal"},
			},
			"http://example.com?firstKey=firstVal&secondKey=secondVal",
		},
		{
			"starting URL has existing params",
			`"?firstKey=firstVal&secondKey=secondVal"`,
			"http://example.com?firstKey=hardVal",
			adapters.QueryParameters{
				"firstKey":  []string{"firstVal"},
				"secondKey": []string{"secondVal"},
			},
			"http://example.com?firstKey=hardVal&secondKey=secondVal",
		},
	}

	for _, tt := range cases {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			qp := adapters.QueryParameters{}
			err := json.Unmarshal([]byte(test.queryParams), &qp)
			hga := adapters.HTTPGet{
				URL:         cltest.WebURL(t, test.startingUrl),
				QueryParams: qp,
			}
			hpa := adapters.HTTPPost{
				URL:         cltest.WebURL(t, test.startingUrl),
				QueryParams: qp,
			}
			requestGET, _ := hga.GetRequest()
			assert.Equal(t, test.expectedURL, requestGET.URL.String())
			assert.Equal(t, test.expectedParams, hga.QueryParams)
			requestPOST, _ := hpa.GetRequest("")
			assert.Equal(t, test.expectedURL, requestPOST.URL.String())
			assert.Equal(t, test.expectedParams, hpa.QueryParams)
			assert.Nil(t, err)
		})
	}
}

func TestQueryParameters_Error(t *testing.T) {
	t.Parallel()

	baseUrl := "http://example.com"

	cases := []struct {
		name           string
		queryParams    string
		startingUrl    string
		expectedParams adapters.QueryParameters
		expectedURL    string
	}{
		{
			"odd number of params",
			`["firstKey","firstVal","secondKey","secondVal","bad"]`,
			baseUrl,
			adapters.QueryParameters{},
			baseUrl,
		},
		{
			"bad format of string",
			`"firstKey=firstVal&secondKey=secondVal&bad"`,
			baseUrl,
			adapters.QueryParameters{},
			baseUrl,
		},
	}

	for _, tt := range cases {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			qp := adapters.QueryParameters{}
			err := json.Unmarshal([]byte(test.queryParams), &qp)
			hga := adapters.HTTPGet{
				URL:         cltest.WebURL(t, test.startingUrl),
				QueryParams: qp,
			}
			hpa := adapters.HTTPPost{
				URL:         cltest.WebURL(t, test.startingUrl),
				QueryParams: qp,
			}
			requestGET, _ := hga.GetRequest()
			assert.Equal(t, test.expectedURL, requestGET.URL.String())
			assert.Equal(t, test.expectedParams, hga.QueryParams)
			requestPOST, _ := hpa.GetRequest("")
			assert.Equal(t, test.expectedURL, requestPOST.URL.String())
			assert.Equal(t, test.expectedParams, hpa.QueryParams)
			assert.NotNil(t, err)
		})
	}
}

func TestExtendedPath_Success(t *testing.T) {
	t.Parallel()

	baseUrl := "http://example.com"

	cases := []struct {
		name         string
		startingUrl  string
		path         string
		expectedPath adapters.ExtendedPath
		expectedURL  string
	}{
		{
			"empty",
			baseUrl,
			`""`,
			adapters.ExtendedPath{""},
			baseUrl,
		},
		{
			"two paths",
			baseUrl,
			`"one/two"`,
			adapters.ExtendedPath{
				"one",
				"two",
			},
			"http://example.com/one/two",
		},
		{
			"existing path no trailing slash",
			"http://example.com/one",
			`"two/three"`,
			adapters.ExtendedPath{
				"two",
				"three",
			},
			"http://example.com/one/two/three",
		},
		{
			"existing path with trailing slash",
			"http://example.com/one/",
			`"two/three"`,
			adapters.ExtendedPath{
				"two",
				"three",
			},
			"http://example.com/one/two/three",
		},
		{
			"input as arrays",
			baseUrl,
			`["one","two"]`,
			adapters.ExtendedPath{
				"one",
				"two",
			},
			"http://example.com/one/two",
		},
		{
			"input begins with slash",
			baseUrl,
			`"/one/two"`,
			adapters.ExtendedPath{
				"",
				"one",
				"two",
			},
			"http://example.com/one/two",
		},
		{
			"input ends with slash",
			baseUrl,
			`"one/two/"`,
			adapters.ExtendedPath{
				"one",
				"two",
				"",
			},
			"http://example.com/one/two",
		},
	}

	for _, tt := range cases {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			ep := adapters.ExtendedPath{}
			err := json.Unmarshal([]byte(test.path), &ep)
			hga := adapters.HTTPGet{
				URL:          cltest.WebURL(t, test.startingUrl),
				ExtendedPath: ep,
			}
			hpa := adapters.HTTPPost{
				URL:          cltest.WebURL(t, test.startingUrl),
				ExtendedPath: ep,
			}
			requestGET, _ := hga.GetRequest()
			assert.Equal(t, test.expectedURL, requestGET.URL.String())
			assert.Equal(t, test.expectedPath, hga.ExtendedPath)
			requestPOST, _ := hpa.GetRequest("")
			assert.Equal(t, test.expectedURL, requestPOST.URL.String())
			assert.Equal(t, test.expectedPath, hpa.ExtendedPath)
			assert.Nil(t, err)
		})
	}
}

func TestExtendedPath_Error(t *testing.T) {
	t.Parallel()

	baseUrl := "http://example.com"

	cases := []struct {
		name         string
		startingUrl  string
		path         string
		expectedPath adapters.ExtendedPath
		expectedURL  string
	}{
		{
			"bad array input",
			baseUrl,
			`["one","two"`,
			adapters.ExtendedPath{},
			baseUrl,
		},
	}

	for _, tt := range cases {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			ep := adapters.ExtendedPath{}
			err := json.Unmarshal([]byte(test.path), &ep)
			hga := adapters.HTTPGet{
				URL:          cltest.WebURL(t, test.startingUrl),
				ExtendedPath: ep,
			}
			hpa := adapters.HTTPPost{
				URL:          cltest.WebURL(t, test.startingUrl),
				ExtendedPath: ep,
			}
			requestGET, _ := hga.GetRequest()
			assert.Equal(t, test.expectedURL, requestGET.URL.String())
			assert.Equal(t, test.expectedPath, hga.ExtendedPath)
			requestPOST, _ := hpa.GetRequest("")
			assert.Equal(t, test.expectedURL, requestPOST.URL.String())
			assert.Equal(t, test.expectedPath, hpa.ExtendedPath)
			assert.NotNil(t, err)

		})
	}
}

func TestHTTP_BuildingURL(t *testing.T) {
	t.Parallel()

	baseUrl := "http://example.com"

	cases := []struct {
		name        string
		startingUrl string
		path        string
		queryParams string
		expectedURL string
	}{
		{
			"one of each",
			baseUrl,
			`"one"`,
			`"firstKey=firstVal"`,
			"http://example.com/one?firstKey=firstVal",
		},
		{
			"query params no path",
			baseUrl,
			`""`,
			`"firstKey=firstVal"`,
			"http://example.com?firstKey=firstVal",
		},
		{
			"path no query params",
			baseUrl,
			`"one"`,
			`""`,
			"http://example.com/one",
		},
		{
			"many of each",
			baseUrl,
			`"one/two/three"`,
			`"firstKey=firstVal&secondKey=secondVal"`,
			"http://example.com/one/two/three?firstKey=firstVal&secondKey=secondVal",
		},
		{
			"existing path",
			"http://example.com/one",
			`"two"`,
			`"firstKey=firstVal"`,
			"http://example.com/one/two?firstKey=firstVal",
		},
		{
			"existing query param",
			"http://example.com?firstKey=firstVal",
			`"one"`,
			`"secondKey=secondVal"`,
			"http://example.com/one?firstKey=firstVal&secondKey=secondVal",
		},
		{
			"existing path and query param",
			"http://example.com/one?firstKey=firstVal",
			`"two"`,
			`"secondKey=secondVal"`,
			"http://example.com/one/two?firstKey=firstVal&secondKey=secondVal",
		},
	}

	for _, tt := range cases {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			ep := adapters.ExtendedPath{}
			qp := adapters.QueryParameters{}
			err := json.Unmarshal([]byte(test.path), &ep)
			err = json.Unmarshal([]byte(test.queryParams), &qp)
			hga := adapters.HTTPGet{
				URL:          cltest.WebURL(t, test.startingUrl),
				QueryParams:  qp,
				ExtendedPath: ep,
			}
			hpa := adapters.HTTPPost{
				URL:          cltest.WebURL(t, test.startingUrl),
				QueryParams:  qp,
				ExtendedPath: ep,
			}
			requestGET, _ := hga.GetRequest()
			assert.Equal(t, test.expectedURL, requestGET.URL.String())
			requestPOST, _ := hpa.GetRequest("")
			assert.Equal(t, test.expectedURL, requestPOST.URL.String())
			assert.Nil(t, err)
		})
	}
}
