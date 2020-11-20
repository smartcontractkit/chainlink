package pipeline

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type (
	MaybeBool string
)

const (
	MaybeBoolTrue  = MaybeBool("true")
	MaybeBoolFalse = MaybeBool("false")
	MaybeBoolNull  = MaybeBool("")
)

func MaybeBoolFromString(s string) (MaybeBool, error) {
	switch s {
	case "true":
		return MaybeBoolTrue, nil
	case "false":
		return MaybeBoolFalse, nil
	case "":
		return MaybeBoolNull, nil
	default:
		return "", errors.Errorf("unknown value for bool: %s", s)
	}
}

func (m MaybeBool) Bool() (b bool, isSet bool) {
	switch m {
	case MaybeBoolTrue:
		return true, true
	case MaybeBoolFalse:
		return false, true
	default:
		return false, false
	}
}

type HTTPTask struct {
	BaseTask                       `mapstructure:",squash"`
	Method                         string
	URL                            models.WebURL
	RequestData                    HttpRequestData `json:"requestData"`
	AllowUnrestrictedNetworkAccess MaybeBool

	config Config
}

type PossibleErrorResponses struct {
	Error        string `json:"error"`
	ErrorMessage string `json:"errorMessage"`
}

var _ Task = (*HTTPTask)(nil)

func (t *HTTPTask) Type() TaskType {
	return TaskTypeHTTP
}

func (t *HTTPTask) Run(ctx context.Context, taskRun TaskRun, inputs []Result) Result {
	if len(inputs) > 0 {
		return Result{Error: errors.Wrapf(ErrWrongInputCardinality, "HTTPTask requires 0 inputs")}
	}

	var bodyReader io.Reader
	if t.RequestData != nil {
		bodyBytes, err := json.Marshal(t.RequestData)
		if err != nil {
			return Result{Error: errors.Wrap(err, "failed to encode request body as JSON")}
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	request, err := http.NewRequest(t.Method, t.URL.String(), bodyReader)
	if err != nil {
		return Result{Error: errors.Wrap(err, "failed to create http.Request")}
	}
	request.Header.Set("Content-Type", "application/json")

	config := utils.HTTPRequestConfig{
		Timeout:                        t.config.DefaultHTTPTimeout().Duration(),
		MaxAttempts:                    t.config.DefaultMaxHTTPAttempts(),
		SizeLimit:                      t.config.DefaultHTTPLimit(),
		AllowUnrestrictedNetworkAccess: t.allowUnrestrictedNetworkAccess(),
	}

	httpRequest := utils.HTTPRequest{
		Request: request,
		Config:  config,
	}

	responseBytes, statusCode, err := httpRequest.SendRequest(ctx)
	if err != nil {
		return Result{Error: errors.Wrapf(err, "error making http request")}
	}

	if statusCode >= 400 {
		maybeErr := bestEffortExtractError(responseBytes)
		return Result{Error: errors.Errorf("got error from %s: (status code %v) %s", t.URL.String(), statusCode, maybeErr)}
	}

	logger.Debugw("HTTP task got response",
		"response", string(responseBytes),
		"url", t.URL.String(),
	)
	return Result{Value: responseBytes}
}

func (t *HTTPTask) allowUnrestrictedNetworkAccess() bool {
	b, isSet := t.AllowUnrestrictedNetworkAccess.Bool()
	if isSet {
		return b
	}
	return t.config.DefaultHTTPAllowUnrestrictedNetworkAccess()
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
