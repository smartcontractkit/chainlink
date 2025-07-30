package mocks

import (
	"io"
	"net/http"
	"strings"
)

// HTTPRoundTripper is a mock implementation of http.RoundTripper for testing
type HTTPRoundTripper struct {
	Response     *http.Response
	ResponseBody string
	Error        error
}

func (m *HTTPRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.Error != nil {
		return nil, m.Error
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
