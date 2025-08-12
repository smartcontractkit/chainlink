package mocks

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// HTTPRoundTripper is a mock implementation of http.RoundTripper for testing
type HTTPRoundTripper struct {
	Response     *http.Response
	ResponseBody string
	Error        error
	ExpectedURL  string
}

func (m *HTTPRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	if m.ExpectedURL != "" && req.URL.String() != m.ExpectedURL {
		return nil, fmt.Errorf("unexpected URL: got %q, expected %q", req.URL.String(), m.ExpectedURL)
	}

	response := *m.Response
	response.Body = io.NopCloser(strings.NewReader(m.ResponseBody))
	response.Request = req
	return &response, nil
}

// NewMockHTTPClient creates a new HTTP client with mock transport for testing
func NewMockHTTPClient(responseBody string, statusCode int) *http.Client {
	return &http.Client{
		Transport: &HTTPRoundTripper{
			Response: &http.Response{
				StatusCode: statusCode,
				Header:     make(http.Header),
				Body:       http.NoBody,
			},
			ResponseBody: responseBody,
		},
	}
}

// NewMockHTTPClientWithExpectedURL creates a new HTTP client that validates the request URL
func NewMockHTTPClientWithExpectedURL(responseBody string, statusCode int, expectedURL string) *http.Client {
	return &http.Client{
		Transport: &HTTPRoundTripper{
			Response: &http.Response{
				StatusCode: statusCode,
				Header:     make(http.Header),
				Body:       http.NoBody,
			},
			ResponseBody: responseBody,
			ExpectedURL:  expectedURL,
		},
	}
}
