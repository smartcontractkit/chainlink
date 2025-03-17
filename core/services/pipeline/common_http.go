package pipeline

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	clhttp "github.com/smartcontractkit/chainlink/v2/core/utils/http"
)

func makeHTTPRequest(
	ctx context.Context,
	lggr logger.Logger,
	method StringParam,
	url URLParam,
	reqHeaders []string,
	requestData MapParam,
	client *http.Client,
	httpLimit int64,
) (responseBytes []byte, statusCode int, respHeaders http.Header, start, finish time.Time, err error) {
	var bodyReader io.Reader
	if requestData != nil {
		var bodyBytes []byte
		bodyBytes, err = json.Marshal(requestData)
		if err != nil {
			err = errors.Wrap(err, "failed to encode request body as JSON")
			return
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	var request *http.Request
	request, err = http.NewRequestWithContext(ctx, string(method), url.String(), bodyReader)
	if err != nil {
		err = errors.Wrap(err, "failed to create http.Request")
		return
	}
	request.Header.Set("Content-Type", "application/json")
	if len(reqHeaders)%2 != 0 {
		panic("headers must have an even number of elements")
	}
	for i := 0; i+1 < len(reqHeaders); i += 2 {
		request.Header.Set(reqHeaders[i], reqHeaders[i+1])
	}

	httpRequest := clhttp.HTTPRequest{
		Client:  client,
		Request: request,
		Config:  clhttp.HTTPRequestConfig{SizeLimit: httpLimit},
		Logger:  logger.Sugared(lggr).Named("HTTPRequest"),
	}

	start = time.Now()
	responseBytes, statusCode, respHeaders, err = httpRequest.SendRequest()
	finish = time.Now()
	if ctx.Err() != nil {
		err = errors.New("http request timed out or interrupted")
		return
	}
	if err != nil {
		err = errors.Wrapf(err, "error making http request")
		return
	}

	if statusCode >= 400 {
		err = errors.Errorf("got error from %s: (status code %v) %s", url.String(), statusCode, bestEffortExtractError(responseBytes))
	}
	return
}

type PossibleErrorResponses struct {
	Error        string `json:"error"`
	ErrorMessage string `json:"errorMessage"`
}

func bestEffortExtractError(responseBytes []byte) string {
	var resp PossibleErrorResponses
	err := json.Unmarshal(responseBytes, &resp)
	if err != nil {
		return ""
	}
	if resp.Error != "" {
		return resp.Error
	} else if resp.ErrorMessage != "" {
		return resp.ErrorMessage
	}
	return string(responseBytes)
}

func httpRequestCtx(ctx context.Context, t Task, cfg Config) (requestCtx context.Context, cancel context.CancelFunc) {
	// Only set the default timeout if the task timeout is missing; task
	// timeout if present will have already been set on the context at a higher
	// level. If task timeout is explicitly set to zero, we must not override
	// with the default http timeout here (since it has been explicitly
	// disabled).
	//
	// DefaultHTTPTimeout is not used if set to 0.
	if _, isSet := t.TaskTimeout(); !isSet && cfg.DefaultHTTPTimeout().Duration() > 0 {
		requestCtx, cancel = context.WithTimeout(ctx, cfg.DefaultHTTPTimeout().Duration())
	} else {
		requestCtx = ctx
		cancel = func() {}
	}
	return
}

// statusCodeGroup maps to course status code group (e.g. 2xx, 4xx, 5xx) to reduce metric cardinality.
func statusCodeGroup(status int) string {
	switch {
	case status >= 100 && status < 200:
		return "1xx"
	case status >= 200 && status < 300:
		return "2xx"
	case status >= 300 && status < 400:
		return "3xx"
	case status >= 400 && status < 500:
		return "4xx"
	case status >= 500 && status < 600:
		return "5xx"
	default:
		return "unknown"
	}
}
