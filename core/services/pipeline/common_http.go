package pipeline

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/utils"
)

func makeHTTPRequest(
	ctx context.Context,
	method StringParam,
	url URLParam,
	requestData MapParam,
	allowUnrestrictedNetworkAccess BoolParam,
	cfg Config,
) ([]byte, http.Header, time.Duration, error) {

	var bodyReader io.Reader
	if requestData != nil {
		bodyBytes, err := json.Marshal(requestData)
		if err != nil {
			return nil, nil, 0, errors.Wrap(err, "failed to encode request body as JSON")
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, cfg.DefaultHTTPTimeout().Duration())
	defer cancel()

	request, err := http.NewRequestWithContext(timeoutCtx, string(method), url.String(), bodyReader)
	if err != nil {
		return nil, nil, 0, errors.Wrap(err, "failed to create http.Request")
	}
	request.Header.Set("Content-Type", "application/json")

	httpRequest := utils.HTTPRequest{
		Request: request,
		Config: utils.HTTPRequestConfig{
			SizeLimit:                      cfg.DefaultHTTPLimit(),
			AllowUnrestrictedNetworkAccess: bool(allowUnrestrictedNetworkAccess),
		},
	}

	start := time.Now()
	responseBytes, statusCode, headers, err := httpRequest.SendRequest()
	if ctx.Err() != nil {
		return nil, nil, 0, errors.New("http request timed out or interrupted")
	}
	if err != nil {
		return nil, nil, 0, errors.Wrapf(err, "error making http request")
	}
	elapsed := time.Since(start) // TODO: return elapsed from utils/http

	if statusCode >= 400 {
		maybeErr := bestEffortExtractError(responseBytes)
		return nil, headers, 0, errors.Errorf("got error from %s: (status code %v) %s", url.String(), statusCode, maybeErr)
	}
	return responseBytes, headers, elapsed, nil
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
