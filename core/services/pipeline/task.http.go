package pipeline

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
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

var (
	promHTTPFetchTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pipeline_task_http_fetch_time",
		Help: "Time taken to fully execute the HTTP request",
	},
		[]string{"pipeline_task_spec_id"},
	)
	promHTTPResponseBodySize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pipeline_task_http_response_body_size",
		Help: "Size (in bytes) of the HTTP response body",
	},
		[]string{"pipeline_task_spec_id"},
	)
)

func (t *HTTPTask) Type() TaskType {
	return TaskTypeHTTP
}

func (t *HTTPTask) SetDefaults(inputValues map[string]string, g TaskDAG, self taskDAGNode) error {
	return nil
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

	start := time.Now()
	responseBytes, statusCode, err := httpRequest.SendRequest(ctx)
	if err != nil {
		if ctx.Err() != nil {
			return Result{Error: errors.New("http request timed out or interrupted")}
		}
		return Result{Error: errors.Wrapf(err, "error making http request")}
	}
	elapsed := time.Since(start)
	promHTTPFetchTime.WithLabelValues(fmt.Sprintf("%d", taskRun.PipelineTaskSpecID)).Set(float64(elapsed))
	promHTTPResponseBodySize.WithLabelValues(fmt.Sprintf("%d", taskRun.PipelineTaskSpecID)).Set(float64(len(responseBytes)))

	if statusCode >= 400 {
		maybeErr := bestEffortExtractError(responseBytes)
		return Result{Error: errors.Errorf("got error from %s: (status code %v) %s", t.URL.String(), statusCode, maybeErr)}
	}

	logger.Debugw("HTTP task got response",
		"response", string(responseBytes),
		"url", t.URL.String(),
		"pipelineTaskSpecID", taskRun.PipelineTaskSpecID,
	)
	// NOTE: We always stringify the response since this is required for all current jobs.
	// If a binary response is required we might consider adding an adapter
	// flag such as  "BinaryMode: true" which passes through raw binary as the
	// value instead.
	return Result{Value: string(responseBytes)}
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
