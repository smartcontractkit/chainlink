package pipeline

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	clhttp "github.com/smartcontractkit/chainlink/core/utils/http"
)

func makeHTTPRequest(
	ctx context.Context,
	lggr logger.Logger,
	method StringParam,
	url URLParam,
	requestData MapParam,
	client *http.Client,
	httpLimit int64,
) ([]byte, int, http.Header, time.Duration, error) {

	var bodyReader io.Reader
	if requestData != nil {
		bodyBytes, err := json.Marshal(requestData)
		if err != nil {
			return nil, 0, nil, 0, errors.Wrap(err, "failed to encode request body as JSON")
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	request, err := http.NewRequestWithContext(ctx, string(method), url.String(), bodyReader)
	if err != nil {
		return nil, 0, nil, 0, errors.Wrap(err, "failed to create http.Request")
	}
	request.Header.Set("Content-Type", "application/json")

	httpRequest := clhttp.HTTPRequest{
		Client:  client,
		Request: request,
		Config:  clhttp.HTTPRequestConfig{SizeLimit: httpLimit},
		Logger:  lggr.Named("HTTPRequest"),
	}

	start := time.Now()
	responseBytes, statusCode, headers, err := httpRequest.SendRequest()
	if ctx.Err() != nil {
		return nil, 0, nil, 0, errors.New("http request timed out or interrupted")
	}
	if err != nil {
		return nil, 0, nil, 0, errors.Wrapf(err, "error making http request")
	}
	elapsed := time.Since(start) // TODO: return elapsed from utils/http

	if statusCode >= 400 {
		maybeErr := bestEffortExtractError(responseBytes)
		return nil, statusCode, headers, 0, errors.Errorf("got error from %s: (status code %v) %s", url.String(), statusCode, maybeErr)
	}
	return responseBytes, statusCode, headers, elapsed, nil
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
