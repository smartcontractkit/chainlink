package starknet

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/juno/utils"
)

type Backoff func(wait time.Duration) time.Duration

type RejectedTransaction struct {
	FailureReason string `json:"transaction_failure_reason"`
}

type TransactionFailureReason struct {
	Code         string `json:"code"`
	ErrorMessage string `json:"error_message"`
}

type FeederClient struct {
	url        string
	client     *http.Client
	backoff    Backoff
	maxRetries int
	maxWait    time.Duration
	minWait    time.Duration
	log        utils.SimpleLogger
}

func (c *FeederClient) WithBackoff(b Backoff) *FeederClient {
	c.backoff = b
	return c
}

func (c *FeederClient) WithMaxRetries(num int) *FeederClient {
	c.maxRetries = num
	return c
}

func (c *FeederClient) WithMaxWait(d time.Duration) *FeederClient {
	c.maxWait = d
	return c
}

func (c *FeederClient) WithMinWait(d time.Duration) *FeederClient {
	c.minWait = d
	return c
}

func (c *FeederClient) WithLogger(log utils.SimpleLogger) *FeederClient {
	c.log = log
	return c
}

func ExponentialBackoff(wait time.Duration) time.Duration {
	return wait * 2
}

func NopBackoff(d time.Duration) time.Duration {
	return 0
}

func NewFeederClient(clientURL string) *FeederClient {
	return &FeederClient{
		url:        clientURL,
		client:     http.DefaultClient,
		backoff:    ExponentialBackoff,
		maxRetries: 5,
		maxWait:    10 * time.Second,
		minWait:    time.Second,
		log:        utils.NewNopZapLogger(),
	}
}

// buildQueryString builds the query url with encoded parameters
func (c *FeederClient) buildQueryString(endpoint string, args map[string]string) string {
	base, err := url.Parse(c.url)
	if err != nil {
		panic("Malformed feeder base URL")
	}

	base.Path += endpoint

	params := url.Values{}
	for k, v := range args {
		params.Add(k, v)
	}
	base.RawQuery = params.Encode()

	return base.String()
}

// get performs a "GET" http request with the given URL and returns the response body
func (c *FeederClient) get(ctx context.Context, queryURL string) (io.ReadCloser, error) {
	var res *http.Response
	var err error
	wait := time.Duration(0)
	for i := 0; i <= c.maxRetries; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(wait):
			var req *http.Request
			req, err = http.NewRequestWithContext(ctx, "GET", queryURL, http.NoBody)
			if err != nil {
				return nil, err
			}

			res, err = c.client.Do(req)
			if err == nil {
				if res.StatusCode == http.StatusOK {
					return res.Body, nil
				}

				if res.StatusCode != http.StatusOK {
					err = errors.New(res.Status)
				}
				res.Body.Close()
			}

			if wait < c.minWait {
				wait = c.minWait
			}
			wait = c.backoff(wait)
			if wait > c.maxWait {
				wait = c.maxWait
			}
			c.log.Warnw("failed query to feeder, retrying...", "retryAfter", wait.String())
		}
	}
	return nil, err
}

func (c *FeederClient) TransactionFailure(ctx context.Context, transactionHash *felt.Felt) (*TransactionFailureReason, error) {
	queryURL := c.buildQueryString("get_transaction", map[string]string{
		"transactionHash": transactionHash.String(),
	})

	body, err := c.get(ctx, queryURL)

	if err != nil {
		return nil, err
	}
	defer body.Close()

	txStatus := new(TransactionFailureReason)
	if err = json.NewDecoder(body).Decode(txStatus); err != nil {
		return nil, err
	}
	return txStatus, nil
}

// Only responds on valid /get_transaction?transactionHash=<TRANSACTION_HASH> requests. fails otherwise
func NewTestFeederServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryMap, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil || !strings.HasSuffix(r.URL.Path, "get_transaction") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		queryArg := "transactionHash"

		_, found := queryMap[queryArg]
		if !found {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defaultErr := []byte(`{
			"code": "SOME_ERROR",
			"error_message": "some error was encountered"
		}`)

		_, err = w.Write(defaultErr)
		if err != nil {
			panic(err)
		}
	}))
}

// NewTestFeederClient returns a client and a function to close a test server.
func NewTestFeederClient(t *testing.T) *FeederClient {
	srv := NewTestFeederServer()
	t.Cleanup(srv.Close)

	c := NewFeederClient(srv.URL).WithBackoff(NopBackoff).WithMaxRetries(0)
	c.client = &http.Client{
		Transport: &http.Transport{
			// On macOS tests often fail with the following error:
			//
			// "Get "http://127.0.0.1:xxxx/get_{feeder gateway method}?{arg}={value}": dial tcp 127.0.0.1:xxxx:
			//    connect: can't assign requested address"
			//
			// This error makes running local tests, in quick succession, difficult because we have to wait for the OS to release ports.
			// Sometimes the sync tests will hang because sync process will keep making requests if there was some error.
			// This problem is further exacerbated by having parallel tests.
			//
			// Increasing test client's idle conns allows for large concurrent requests to be made from a single test client.
			MaxIdleConnsPerHost: 1000,
		},
	}
	return c
}
