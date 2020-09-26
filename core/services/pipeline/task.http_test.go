package pipeline_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"

	"github.com/stretchr/testify/require"
)

func TestHTTPTask_NotAUrlError(t *testing.T) {
	t.Parallel()

	task := &pipeline.HTTPTask{URL: cltest.WebURL(t, "NotAURL")}
	task.HelperSetConfig(orm.NewConfig())
	result := task.Run(nil)
	require.Error(t, result.Error)
	require.Nil(t, result.Value)
}

func TestHTTPTask_GET(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		status      int
		want        string
		wantErrored bool
		response    string
		headers     pipeline.Header
		queryParams pipeline.QueryParameters
	}{
		{"success", http.StatusOK, "results!", false, `results!`, nil, nil},
		{"success but error in body", http.StatusOK, `{"error": "results!"}`, false, `{"error": "results!"}`, nil, nil},
		{"success with HTML", http.StatusOK, `<html>results!</html>`, false, `<html>results!</html>`, nil, nil},
		{"success with headers", http.StatusOK, "results!", false, `results!`,
			pipeline.Header{
				"Key1": []string{"value"},
				"Key2": []string{"value", "value"},
			}, nil},
		{"not found", http.StatusBadRequest, "", true, `<html>so bad</html>`, nil, nil},
		{"server error", http.StatusBadRequest, "", true, `Invalid request`, nil, nil},
		{"success with params", http.StatusOK, "results!", false, `results!`, nil,
			pipeline.QueryParameters{
				"Key1": []string{"value0"},
				"Key2": []string{"value1", "value2"},
			}},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			mock, cleanup := cltest.NewHTTPMockServer(t, test.status, "GET", test.response,
				func(header http.Header, body string) {
					require.Equal(t, ``, body)
					for key, values := range test.headers {
						require.Equal(t, values, header[key])
					}
				})
			defer cleanup()

			task := pipeline.HTTPTask{
				URL:                            cltest.WebURL(t, mock.URL),
				Method:                         "GET",
				Headers:                        test.headers,
				QueryParams:                    test.queryParams,
				AllowUnrestrictedNetworkAccess: true,
			}
			task.HelperSetConfig(orm.NewConfig())

			result := task.Run(nil)

			if test.wantErrored {
				require.Error(t, result.Error)
			} else {
				require.NoError(t, result.Error)
				require.Equal(t, test.want, string(result.Value.([]byte)))
			}
		})
	}
}

func TestHTTPTask_TooLarge(t *testing.T) {
	cfg := orm.NewConfig()
	cfg.Set("DEFAULT_HTTP_LIMIT", "1")
	cfg.Set("MAX_HTTP_ATTEMPTS", "3")

	tests := []struct {
		verb    string
		factory func(models.WebURL) *pipeline.HTTPTask
	}{
		{"GET", func(url models.WebURL) *pipeline.HTTPTask {
			return &pipeline.HTTPTask{Method: "GET", URL: url, AllowUnrestrictedNetworkAccess: true}
		}},
		{"POST", func(url models.WebURL) *pipeline.HTTPTask {
			return &pipeline.HTTPTask{Method: "POST", URL: url, AllowUnrestrictedNetworkAccess: true}
		}},
	}
	for _, test := range tests {
		t.Run(test.verb, func(t *testing.T) {
			largePayload := "12"
			mock, cleanup := cltest.NewHTTPMockServer(t, http.StatusOK, test.verb, largePayload)
			defer cleanup()

			task := test.factory(cltest.WebURL(t, mock.URL))
			task.HelperSetConfig(cfg)
			result := task.Run(nil)

			require.Error(t, result.Error)
			require.Contains(t, result.Error.Error(), "HTTP response too large")
			require.Nil(t, result.Value)
		})
	}
}

func TestHTTPTask_WithRestrictedIP(t *testing.T) {
	cfg := orm.NewConfig()

	tests := []struct {
		verb    string
		factory func(models.WebURL) *pipeline.HTTPTask
	}{
		{"GET", func(url models.WebURL) *pipeline.HTTPTask {
			return &pipeline.HTTPTask{Method: "GET", URL: url, AllowUnrestrictedNetworkAccess: false}
		}},
		{"POST", func(url models.WebURL) *pipeline.HTTPTask {
			return &pipeline.HTTPTask{Method: "POST", URL: url, AllowUnrestrictedNetworkAccess: false}
		}},
	}
	for _, test := range tests {
		t.Run(test.verb, func(t *testing.T) {
			payload := ""
			mock, _ := cltest.NewHTTPMockServer(t, http.StatusOK, test.verb, payload)
			defer mock.Close()

			task := test.factory(cltest.WebURL(t, mock.URL))
			task.HelperSetConfig(cfg)
			result := task.Run(nil)

			require.Error(t, result.Error)
			require.Contains(t, result.Error.Error(), "disallowed IP")
			require.Nil(t, result.Value)
		})
	}
}

func stringRef(str string) *string {
	return &str
}

func TestHTTPTask_POST(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		status      int
		want        string
		wantBody    string
		wantErrored bool
		response    string
		headers     pipeline.Header
		queryParams pipeline.QueryParameters
		body        pipeline.HttpRequestData
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
			pipeline.Header{
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
			pipeline.QueryParameters{
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
			pipeline.HttpRequestData{"Key1": "value", "Key2": "value"},
		},
		{
			"success with body",
			http.StatusOK,
			"results!",
			"{}",
			false,
			`results!`,
			nil,
			nil,
			pipeline.HttpRequestData{},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			mock, cleanup := cltest.NewHTTPMockServer(t, test.status, "POST", test.response,
				func(header http.Header, body string) {
					require.Equal(t, test.wantBody, body)
					for key, values := range test.headers {
						require.Equal(t, values, header[key])
					}
				})
			defer cleanup()

			task := pipeline.HTTPTask{
				Method:                         "POST",
				URL:                            cltest.WebURL(t, mock.URL),
				Headers:                        test.headers,
				QueryParams:                    test.queryParams,
				RequestData:                    test.body,
				AllowUnrestrictedNetworkAccess: true,
			}
			task.HelperSetConfig(orm.NewConfig())

			result := task.Run([]pipeline.Result{{Value: []byte(`{"result":"inputVal"}`)}})

			if !test.wantErrored {
				require.Equal(t, test.want, string(result.Value.([]byte)))
				require.NoError(t, result.Error)
			} else {
				require.Error(t, result.Error)
				require.Nil(t, result.Value)
			}
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
		expectedParams pipeline.QueryParameters
		expectedURL    string
	}{
		{"empty",
			`""`,
			baseUrl,
			pipeline.QueryParameters{},
			baseUrl,
		},
		{
			"array of params",
			`["firstKey","firstVal","secondKey","secondVal"]`,
			baseUrl,
			pipeline.QueryParameters{
				"firstKey":  []string{"firstVal"},
				"secondKey": []string{"secondVal"},
			},
			"http://example.com?firstKey=firstVal&secondKey=secondVal",
		},
		{
			"string of params",
			`"firstKey=firstVal&secondKey=secondVal"`,
			baseUrl,
			pipeline.QueryParameters{
				"firstKey":  []string{"firstVal"},
				"secondKey": []string{"secondVal"},
			},
			"http://example.com?firstKey=firstVal&secondKey=secondVal",
		},
		{
			"string has question mark",
			`"?firstKey=firstVal&secondKey=secondVal"`,
			baseUrl,
			pipeline.QueryParameters{
				"firstKey":  []string{"firstVal"},
				"secondKey": []string{"secondVal"},
			},
			"http://example.com?firstKey=firstVal&secondKey=secondVal",
		},
		{
			"starting URL has existing params",
			`"?firstKey=firstVal&secondKey=secondVal"`,
			"http://example.com?firstKey=hardVal",
			pipeline.QueryParameters{
				"firstKey":  []string{"firstVal"},
				"secondKey": []string{"secondVal"},
			},
			"http://example.com?firstKey=hardVal&secondKey=secondVal",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			qp := pipeline.QueryParameters{}
			err := json.Unmarshal([]byte(test.queryParams), &qp)
			require.NoError(t, err)
			task := pipeline.HTTPTask{
				URL:         cltest.WebURL(t, test.startingUrl),
				QueryParams: qp,
			}
			request, _ := task.Request("", nil)
			require.Equal(t, test.expectedURL, request.URL.String())
			require.Equal(t, test.expectedParams, task.QueryParams)
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
		expectedParams pipeline.QueryParameters
		expectedURL    string
	}{
		{
			"odd number of params",
			`["firstKey","firstVal","secondKey","secondVal","bad"]`,
			baseUrl,
			pipeline.QueryParameters{},
			baseUrl,
		},
		{
			"bad format of string",
			`"firstKey=firstVal&secondKey=secondVal&bad"`,
			baseUrl,
			pipeline.QueryParameters{},
			baseUrl,
		},
		{
			"invalid type",
			`{"firstKey": "firstVal", "secondKey": "secondVal"}`,
			baseUrl,
			pipeline.QueryParameters{},
			baseUrl,
		},
		{
			"invalid json",
			"invalid",
			baseUrl,
			pipeline.QueryParameters{},
			baseUrl,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			qp := pipeline.QueryParameters{}
			err := json.Unmarshal([]byte(test.queryParams), &qp)
			require.Error(t, err)
			task := pipeline.HTTPTask{
				URL:         cltest.WebURL(t, test.startingUrl),
				QueryParams: qp,
			}
			request, _ := task.Request("", nil)
			require.Equal(t, test.expectedURL, request.URL.String())
			require.Equal(t, test.expectedParams, task.QueryParams)
		})
	}
}

func TestHTTPTask_ExtendedPath_Success(t *testing.T) {
	t.Parallel()

	baseUrl := "http://example.com"

	cases := []struct {
		name         string
		startingUrl  string
		path         string
		expectedPath pipeline.ExtendedPath
		expectedURL  string
	}{
		{
			"empty",
			baseUrl,
			`""`,
			pipeline.ExtendedPath{""},
			baseUrl,
		},
		{
			"two paths",
			baseUrl,
			`"one/two"`,
			pipeline.ExtendedPath{
				"one",
				"two",
			},
			"http://example.com/one/two",
		},
		{
			"existing path no trailing slash",
			"http://example.com/one",
			`"two/three"`,
			pipeline.ExtendedPath{
				"two",
				"three",
			},
			"http://example.com/one/two/three",
		},
		{
			"existing path with trailing slash",
			"http://example.com/one/",
			`"two/three"`,
			pipeline.ExtendedPath{
				"two",
				"three",
			},
			"http://example.com/one/two/three",
		},
		{
			"input as arrays",
			baseUrl,
			`["one","two"]`,
			pipeline.ExtendedPath{
				"one",
				"two",
			},
			"http://example.com/one/two",
		},
		{
			"input begins with slash",
			baseUrl,
			`"/one/two"`,
			pipeline.ExtendedPath{
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
			pipeline.ExtendedPath{
				"one",
				"two",
				"",
			},
			"http://example.com/one/two",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			ep := pipeline.ExtendedPath{}
			err := json.Unmarshal([]byte(test.path), &ep)
			task := pipeline.HTTPTask{
				URL:          cltest.WebURL(t, test.startingUrl),
				ExtendedPath: ep,
			}
			task.HelperSetConfig(orm.NewConfig())
			requestGET, _ := task.Request("", nil)
			require.Equal(t, test.expectedURL, requestGET.URL.String())
			require.Equal(t, test.expectedPath, task.ExtendedPath)
			require.Nil(t, err)
		})
	}
}

func TestHTTPTask_ExtendedPath_Error(t *testing.T) {
	t.Parallel()

	baseUrl := "http://example.com"

	cases := []struct {
		name         string
		startingUrl  string
		path         string
		expectedPath pipeline.ExtendedPath
		expectedURL  string
	}{
		{
			"bad array input",
			baseUrl,
			`["one","two"`,
			pipeline.ExtendedPath{},
			baseUrl,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			ep := pipeline.ExtendedPath{}
			err := json.Unmarshal([]byte(test.path), &ep)
			task := pipeline.HTTPTask{
				URL:          cltest.WebURL(t, test.startingUrl),
				ExtendedPath: ep,
			}
			task.HelperSetConfig(orm.NewConfig())
			requestGET, _ := task.Request("", nil)
			require.Equal(t, test.expectedURL, requestGET.URL.String())
			require.Equal(t, test.expectedPath, task.ExtendedPath)
			require.NotNil(t, err)

		})
	}
}

func TestHTTPTask_BuildingURL(t *testing.T) {
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

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			ep := pipeline.ExtendedPath{}
			qp := pipeline.QueryParameters{}
			err := json.Unmarshal([]byte(test.path), &ep)
			require.NoError(t, err, "failed to unmarshal path: %s to adapter.", test.path)
			err = json.Unmarshal([]byte(test.queryParams), &qp)
			task := pipeline.HTTPTask{
				URL:          cltest.WebURL(t, test.startingUrl),
				QueryParams:  qp,
				ExtendedPath: ep,
			}
			requestGET, _ := task.Request("", nil)
			require.Equal(t, test.expectedURL, requestGET.URL.String())
			require.Nil(t, err)
		})
	}
}

func TestHTTPTask_JSONDeserializationDoesNotSetAllowUnrestrictedNetworkAccess(t *testing.T) {
	var task pipeline.HTTPTask
	err := json.Unmarshal([]byte(`{"allowUnrestrictedNetworkAccess": true}`), &task)
	require.NoError(t, err)
	require.False(t, task.AllowUnrestrictedNetworkAccess)
}

func TestHTTPTask_RetryPolicy(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	t.Parallel()
	str := leanStore()

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
				task := makeHTTPGetAdapter(t, srv, str)
				_ = task.Run(nil)
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
		task := makeHTTPGetAdapter(t, srv, str)
		_ = task.Run(nil)
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
		task := makeHTTPGetAdapter(t, srv, str)
		_ = task.Run(nil)
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
		task := makeHTTPGetAdapter(t, srv, str)
		_ = task.Run(nil)
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
		task := makeHTTPGetAdapter(t, srv, str)
		_ = task.Run(nil)
		expected := uint32(str.Config.DefaultMaxHTTPAttempts())
		if atomic.LoadUint32(&counter) != expected {
			t.Fatalf("expected adapter to try %d times but got %d when the server is broken", expected, counter)
		}
	})
}

func leanStore() *store.Store {
	return &store.Store{Config: orm.NewConfig()}
}

func makeHTTPGetAdapter(t *testing.T, server *httptest.Server, str *store.Store) *pipeline.HTTPTask {
	task := &pipeline.HTTPTask{
		URL:                            cltest.WebURL(t, server.URL),
		Method:                         "GET",
		AllowUnrestrictedNetworkAccess: true,
	}
	task.HelperSetConfig(str.Config)
	return task
}

func fillBlob(size int64) []byte {
	body := make([]byte, size)
	var i int64
	for i = 0; i < size; i++ {
		body[i] = 'x'
	}
	return body
}
