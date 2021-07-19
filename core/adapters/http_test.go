package adapters_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func leanStore() *store.Store {
	return &store.Store{Config: config.NewConfig()}
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.adapter.Perform(models.RunInput{}, store, nil)
			assert.True(t, result.HasError())
			assert.Empty(t, result.Data())
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
		{"success", http.StatusOK, "results!", false, `results!`, nil, nil},
		{"success but error in body", http.StatusOK, `{"error": "results!"}`, false, `{"error": "results!"}`, nil, nil},
		{"success with HTML", http.StatusOK, `<html>results!</html>`, false, `<html>results!</html>`, nil, nil},
		{"success with headers", http.StatusOK, "results!", false, `results!`,
			http.Header{
				"Key1": []string{"value"},
				"Key2": []string{"value", "value"},
			}, nil},
		{"not found", http.StatusBadRequest, "", true, `<html>so bad</html>`, nil, nil},
		{"server error", http.StatusBadRequest, "", true, `Invalid request`, nil, nil},
		{"success with params", http.StatusOK, "results!", false, `results!`, nil,
			adapters.QueryParameters{
				"Key1": []string{"value0"},
				"Key2": []string{"value1", "value2"},
			}},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			input := cltest.NewRunInputWithResult("inputValue")
			mock, cleanup := cltest.NewHTTPMockServer(t, test.status, "GET", test.response,
				func(header http.Header, body string) {
					assert.Equal(t, ``, body)
					for key, values := range test.headers {
						assert.Equal(t, values, header[key])
					}
				})
			defer cleanup()

			hga := adapters.HTTPGet{
				URL:                            cltest.WebURL(t, mock.URL),
				Headers:                        test.headers,
				QueryParams:                    test.queryParams,
				AllowUnrestrictedNetworkAccess: true,
			}
			assert.Equal(t, test.queryParams, hga.QueryParams)

			result := hga.Perform(input, store, nil)

			if test.wantErrored {
				require.Error(t, result.Error())
			} else {
				require.NoError(t, result.Error())
				assert.Equal(t, test.want, result.Result().String())
			}
			assert.Equal(t, false, result.Status().PendingBridge())
		})
	}
}

func TestHTTP_TooLarge(t *testing.T) {
	cfg := config.NewConfig()
	cfg.Set("DEFAULT_HTTP_LIMIT", "1")
	cfg.Set("MAX_HTTP_ATTEMPTS", "3")

	store := &store.Store{Config: cfg}

	tests := []struct {
		verb    string
		factory func(models.WebURL) adapters.BaseAdapter
	}{
		{"GET", func(url models.WebURL) adapters.BaseAdapter {
			return &adapters.HTTPGet{URL: url, AllowUnrestrictedNetworkAccess: true}
		}},
		{"POST", func(url models.WebURL) adapters.BaseAdapter {
			return &adapters.HTTPPost{URL: url, AllowUnrestrictedNetworkAccess: true}
		}},
	}
	for _, test := range tests {
		t.Run(test.verb, func(t *testing.T) {
			input := cltest.NewRunInputWithResult("inputValue")
			largePayload := "12"
			mock, cleanup := cltest.NewHTTPMockServer(t, http.StatusOK, test.verb, largePayload)
			defer cleanup()

			hga := test.factory(cltest.WebURL(t, mock.URL))
			result := hga.Perform(input, store, nil)

			require.Error(t, result.Error())
			assert.Contains(t, result.Error().Error(), "HTTP response too large")
			assert.Equal(t, "", result.Result().String())
		})
	}
}

func TestHTTP_PerformWithRestrictedIP(t *testing.T) {
	cfg := config.NewConfig()
	store := &store.Store{Config: cfg}

	tests := []struct {
		verb    string
		factory func(models.WebURL) adapters.BaseAdapter
	}{
		{"GET", func(url models.WebURL) adapters.BaseAdapter {
			return &adapters.HTTPGet{URL: url, AllowUnrestrictedNetworkAccess: false}
		}},
		{"POST", func(url models.WebURL) adapters.BaseAdapter {
			return &adapters.HTTPPost{URL: url, AllowUnrestrictedNetworkAccess: false}
		}},
	}
	for _, test := range tests {
		t.Run(test.verb, func(t *testing.T) {
			input := cltest.NewRunInputWithResult("inputValue")
			payload := ""
			mock, _ := cltest.NewHTTPMockServer(t, http.StatusOK, test.verb, payload)
			defer mock.Close()

			h := test.factory(cltest.WebURL(t, mock.URL))
			result := h.Perform(input, store, nil)

			require.Error(t, result.Error())
			assert.Contains(t, result.Error().Error(), "disallowed IP")
			assert.Equal(t, "", result.Result().String())
		})
	}
}

func stringRef(str string) *string {
	return &str
}

func TestHttpPost_Perform(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		status      int
		want        string
		wantBody    string
		wantErrored bool
		response    string
		headers     http.Header
		queryParams adapters.QueryParameters
		body        *string
	}{
		{
			"success", http.StatusOK, "results!",
			`{"result":"inputVal"}`,
			false,
			`results!`,
			nil,
			nil,
			nil,
		},
		{
			"success but error in body",
			http.StatusOK,
			`{"error": "results!"}`,
			`{"result":"inputVal"}`,
			false,
			`{"error": "results!"}`,
			nil,
			nil,
			nil,
		},
		{
			"success with HTML",
			http.StatusOK,
			`<html>results!</html>`,
			`{"result":"inputVal"}`,
			false,
			`<html>results!</html>`,
			nil,
			nil,
			nil,
		},
		{
			"success with headers",
			http.StatusOK, "results!", `{"result":"inputVal"}`, false, `results!`,
			http.Header{
				"Key1": []string{"value"},
				"Key2": []string{"value", "value"},
			},
			nil,
			nil,
		},
		{
			"not found",
			http.StatusBadRequest,
			"",
			`{"result":"inputVal"}`,
			true,
			`<html>so bad</html>`,
			nil,
			nil,
			nil,
		},
		{
			"server error",
			http.StatusInternalServerError,
			"",
			`{"result":"inputVal"}`,
			true,
			`big error`,
			nil,
			nil,
			nil,
		},
		{
			"success with params",
			http.StatusOK,
			"results!",
			`{"result":"inputVal"}`,
			false,
			`results!`,
			nil,
			adapters.QueryParameters{
				"Key1": []string{"value"},
				"Key2": []string{"value", "value"},
			},
			nil,
		},
		{
			"success with body",
			http.StatusOK,
			"results!",
			`{"Key1":"value","Key2":"value"}`,
			false,
			`results!`,
			nil,
			nil,
			stringRef(`{"Key1":"value","Key2":"value"}`),
		},
		{
			"success with body",
			http.StatusOK,
			"results!",
			"",
			false,
			`results!`,
			nil,
			nil,
			stringRef(""),
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			input := cltest.NewRunInputWithResult("inputVal")
			mock, cleanup := cltest.NewHTTPMockServer(t, test.status, "POST", test.response,
				func(header http.Header, body string) {
					assert.Equal(t, test.wantBody, body)
					for key, values := range test.headers {
						assert.Equal(t, values, header[key])
					}
				})
			defer cleanup()

			hpa := adapters.HTTPPost{
				URL:                            cltest.WebURL(t, mock.URL),
				Headers:                        test.headers,
				QueryParams:                    test.queryParams,
				Body:                           test.body,
				AllowUnrestrictedNetworkAccess: true,
			}
			assert.Equal(t, test.queryParams, hpa.QueryParams)

			result := hpa.Perform(input, leanStore(), nil)

			val := result.Result()
			assert.Equal(t, test.want, val.String())
			assert.NotEqual(t, test.wantErrored, val.Exists())
			require.Equal(t, test.wantErrored, result.HasError())
			assert.Equal(t, false, result.Status().PendingBridge())
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

	for _, test := range cases {
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
		{
			"invalid type",
			`{"firstKey": "firstVal", "secondKey": "secondVal"}`,
			baseUrl,
			adapters.QueryParameters{},
			baseUrl,
		},
		{
			"invalid json",
			"invalid",
			baseUrl,
			adapters.QueryParameters{},
			baseUrl,
		},
	}

	for _, test := range cases {
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

	for _, test := range cases {
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

	for _, test := range cases {
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
			"subdirectory with trailing slash",
			"http://example.com/subdir/",
			`""`,
			`"?firstKey=firstVal"`,
			"http://example.com/subdir/?firstKey=firstVal",
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

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			ep := adapters.ExtendedPath{}
			qp := adapters.QueryParameters{}
			err := json.Unmarshal([]byte(test.path), &ep)
			assert.NoError(t, err, "failed to unmarshal path: %s to adapter.", test.path)
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

func TestHTTP_JSONDeserializationDoesNotSetAllowUnrestrictedNetworkAccess(t *testing.T) {
	hga := adapters.HTTPGet{}
	err := json.Unmarshal([]byte(`{"allowUnrestrictedNetworkAccess": true}`), &hga)
	require.NoError(t, err)
	assert.False(t, hga.AllowUnrestrictedNetworkAccess)

	hpa := adapters.HTTPPost{}
	err = json.Unmarshal([]byte(`{"allowUnrestrictedNetworkAccess": true}`), &hpa)
	require.NoError(t, err)
	assert.False(t, hpa.AllowUnrestrictedNetworkAccess)
}

func TestHTTP_RetryPolicy(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	t.Parallel()
	str := leanStore()
	input := cltest.NewRunInputWithResult("testRetryPolicy")

	t.Run("don't retry if the response status is", func(t *testing.T) {
		for _, statusCode := range []int{200, 300, 400} {
			t.Run(strconv.Itoa(statusCode), func(t *testing.T) {
				t.Parallel()
				counter := uint32(0)
				srv := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						atomic.AddUint32(&counter, 1)
						w.WriteHeader(statusCode)
					}))
				defer srv.Close()
				hga := makeHTTPGetAdapter(t, srv)
				_ = hga.Perform(input, str, nil)
				if atomic.LoadUint32(&counter) != 1 {
					t.Fatalf("expected retry count to be 1 for status %d but is %d", statusCode, counter)
				}
			})
		}
	})
	t.Run("retry if the response is 5xx", func(t *testing.T) {
		t.Parallel()
		counter := uint32(0)
		srv := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				atomic.AddUint32(&counter, 1)
				if counter <= 2 {
					w.WriteHeader(500)
					return
				}
				w.WriteHeader(200)
			}))
		defer srv.Close()
		hga := makeHTTPGetAdapter(t, srv)
		_ = hga.Perform(input, str, nil)
		if atomic.LoadUint32(&counter) != 3 {
			t.Fatalf("expected adapter to make 3 call, when the first 2 are 500s, instead it made %d calls", counter)
		}
	})
	t.Run("don't retry if response body is too large", func(t *testing.T) {
		t.Parallel()
		counter := uint32(0)
		srv := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				atomic.AddUint32(&counter, 1)
				w.WriteHeader(200)
				largeBody := fillBlob(str.Config.DefaultHTTPLimit() + 10)
				w.Write(largeBody)
			}))
		defer srv.Close()
		hga := makeHTTPGetAdapter(t, srv)
		_ = hga.Perform(input, str, nil)
		if atomic.LoadUint32(&counter) != 1 {
			t.Fatalf("expected adapter to give up when it receives a large response but instead it tried %d times", counter)
		}
	})
	t.Run("retry maxAttempts times then give up", func(t *testing.T) {
		t.Parallel()
		var counter uint32 = 0
		srv := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				atomic.AddUint32(&counter, 1)
				w.WriteHeader(500)
			}))
		defer srv.Close()
		hga := makeHTTPGetAdapter(t, srv)
		_ = hga.Perform(input, str, nil)
		expected := str.Config.DefaultMaxHTTPAttempts()
		if atomic.LoadUint32(&counter) != uint32(expected) {
			t.Fatalf("expected adapter to give up after %d attempts but instead it tried %d times", expected, counter)
		}
	})
	t.Run("retry if the server is broken", func(t *testing.T) {
		t.Parallel()
		var counter uint32 = 0
		srv := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				atomic.AddUint32(&counter, 1)
				hj, ok := w.(http.Hijacker)
				if !ok {
					t.Fatalf("Unable to hijack response writer!")
				}
				conn, _, err := hj.Hijack()
				if err != nil {
					require.NoError(t, err)
				}
				conn.Close()
			}))
		defer srv.Close()
		hga := makeHTTPGetAdapter(t, srv)
		_ = hga.Perform(input, str, nil)
		expected := uint32(str.Config.DefaultMaxHTTPAttempts())
		if atomic.LoadUint32(&counter) != expected {
			t.Fatalf("expected adapter to try %d times but got %d when the server is broken", expected, counter)
		}
	})
}

// Helpers

func makeHTTPGetAdapter(t *testing.T, server *httptest.Server) *adapters.HTTPGet {
	return &adapters.HTTPGet{
		URL:                            cltest.WebURL(t, server.URL),
		AllowUnrestrictedNetworkAccess: true,
	}
}

func fillBlob(size int64) []byte {
	body := make([]byte, size)
	var i int64
	for i = 0; i < size; i++ {
		body[i] = 'x'
	}
	return body
}
