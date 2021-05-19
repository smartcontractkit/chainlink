package pipeline

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"
	"go.uber.org/multierr"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink/core/logger"
)

type HTTPTask struct {
	BaseTask                       `mapstructure:",squash"`
	Method                         string
	URL                            string
	RequestData                    string `json:"requestData"`
	AllowUnrestrictedNetworkAccess string

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

func (t *HTTPTask) SetDefaults(inputValues map[string]string, g TaskDAG, self TaskDAGNode) error {
	return nil
}

func (t *HTTPTask) Run(ctx context.Context, vars Vars, _ JSONSerializable, inputs []Result) Result {
	_, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: err}
	}

	var (
		method                         StringParam
		url                            URLParam
		requestData                    MapParam
		allowUnrestrictedNetworkAccess MaybeBoolParam
	)
	err = multierr.Combine(
		vars.ResolveValue(&method, From(NonemptyString(t.Method))),
		vars.ResolveValue(&url, From(NonemptyString(t.URL))),
		vars.ResolveValue(&requestData, From(VariableExpr(t.RequestData), NonemptyString(t.RequestData), Input(inputs, 0))),
		vars.ResolveValue(&allowUnrestrictedNetworkAccess, From(t.AllowUnrestrictedNetworkAccess)),
	)
	if err != nil {
		return Result{Error: err}
	}

	var bodyReader io.Reader
	if requestData != nil {
		bodyBytes, err := json.Marshal(requestData)
		if err != nil {
			return Result{Error: errors.Wrap(err, "failed to encode request body as JSON")}
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	request, err := http.NewRequest(string(method), url.String(), bodyReader)
	if err != nil {
		return Result{Error: errors.Wrap(err, "failed to create http.Request")}
	}
	request.Header.Set("Content-Type", "application/json")

	config := utils.HTTPRequestConfig{
		Timeout:                        t.config.DefaultHTTPTimeout().Duration(),
		MaxAttempts:                    t.config.DefaultMaxHTTPAttempts(),
		SizeLimit:                      t.config.DefaultHTTPLimit(),
		AllowUnrestrictedNetworkAccess: t.allowUnrestrictedNetworkAccess(allowUnrestrictedNetworkAccess),
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
	promHTTPFetchTime.WithLabelValues(t.DotID()).Set(float64(elapsed))
	promHTTPResponseBodySize.WithLabelValues(t.DotID()).Set(float64(len(responseBytes)))

	if statusCode >= 400 {
		maybeErr := bestEffortExtractError(responseBytes)
		return Result{Error: errors.Errorf("got error from %s: (status code %v) %s", url.String(), statusCode, maybeErr)}
	}

	logger.Debugw("HTTP task got response",
		"response", string(responseBytes),
		"url", url.String(),
		"dotID", t.DotID(),
	)
	// NOTE: We always stringify the response since this is required for all current jobs.
	// If a binary response is required we might consider adding an adapter
	// flag such as  "BinaryMode: true" which passes through raw binary as the
	// value instead.
	return Result{Value: string(responseBytes)}
}

func (t *HTTPTask) allowUnrestrictedNetworkAccess(mb MaybeBoolParam) bool {
	b, isSet := mb.Bool()
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
