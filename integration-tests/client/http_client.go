package client

import (
	"fmt"
	"net/http"

	resty "github.com/go-resty/resty/v2"
)

// APIClient handles basic request sending logic and cookie handling
type APIClient struct {
	RC     *resty.Client
	Header http.Header
}

type ReqParams struct {
	QueryParams map[string]string
	PathParams  map[string]string
}

type SetReqArgs func(r *ReqParams)

func WithPathParam(p map[string]string) SetReqArgs {
	return func(r *ReqParams) {
		r.PathParams = p
	}
}

func WithQueryParam(q map[string]string) SetReqArgs {
	return func(r *ReqParams) {
		r.QueryParams = q
	}
}

// NewAPIClient returns new basic resty client configured with an base URL
func NewAPIClient(baseURL string) *APIClient {
	rc := resty.New()
	rc.SetBaseURL(baseURL)
	return &APIClient{
		RC: rc,
	}
}

func (c *APIClient) WithHeader(header http.Header) *APIClient {
	c.Header = header
	return c
}

func (c *APIClient) WithRetryCount(retries int) *APIClient {
	c.RC.SetRetryCount(retries)
	return c
}

func (c *APIClient) Request(method,
	endpoint string,
	body interface{},
	obj interface{},
	expectedStatusCode int,
	args ...SetReqArgs,
) (*resty.Response, error) {
	req := c.RC.R()
	req.Method = method
	req.URL = endpoint
	rArgs := &ReqParams{}
	for _, f := range args {
		f(rArgs)
	}
	if rArgs.PathParams != nil {
		req.SetPathParams(rArgs.PathParams)
	}
	if rArgs.QueryParams != nil {
		req.SetQueryParams(rArgs.QueryParams)
	}
	resp, err := req.
		SetHeaderMultiValues(c.Header).
		SetBody(body).
		SetResult(&obj).
		Send()
	if err != nil {
		return nil, err
	} else if resp.IsError() {
		return resp, fmt.Errorf(
			"unexpected response code, got %d",
			resp.StatusCode(),
		)
	} else if resp.StatusCode() != expectedStatusCode {
		return resp, fmt.Errorf(
			"unexpected response code, got %d, expected 200\nURL: %s\nresponse received: %s",
			resp.StatusCode(),
			resp.Request.URL,
			resp.String(),
		)
	}
	return resp, err
}
