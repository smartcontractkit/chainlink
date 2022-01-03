package pipeline

import (
	"context"
	"encoding/json"

	"go.uber.org/multierr"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink/core/logger"
)

//
// Return types:
//     string
//
type HTTPTask struct {
	BaseTask                       `mapstructure:",squash"`
	Method                         string
	URL                            string
	RequestData                    string `json:"requestData"`
	AllowUnrestrictedNetworkAccess string

	config Config
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

func (t *HTTPTask) Run(ctx context.Context, lggr logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	_, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
	}

	var (
		method                         StringParam
		url                            URLParam
		requestData                    MapParam
		allowUnrestrictedNetworkAccess BoolParam
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&method, From(NonemptyString(t.Method), "GET")), "method"),
		errors.Wrap(ResolveParam(&url, From(VarExpr(t.URL, vars), NonemptyString(t.URL))), "url"),
		errors.Wrap(ResolveParam(&requestData, From(VarExpr(t.RequestData, vars), JSONWithVarExprs(t.RequestData, vars, false), nil)), "requestData"),
		errors.Wrap(ResolveParam(&allowUnrestrictedNetworkAccess, From(NonemptyString(t.AllowUnrestrictedNetworkAccess), !variableRegexp.MatchString(t.URL))), "allowUnrestrictedNetworkAccess"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	requestDataJSON, err := json.Marshal(requestData)
	if err != nil {
		return Result{Error: err}, runInfo
	}
	lggr.Debugw("HTTP task: sending request",
		"requestData", string(requestDataJSON),
		"url", url.String(),
		"method", method,
		"allowUnrestrictedNetworkAccess", allowUnrestrictedNetworkAccess,
	)

	responseBytes, statusCode, _, elapsed, err := makeHTTPRequest(ctx, lggr, method, url, requestData, allowUnrestrictedNetworkAccess, t.config)
	if err != nil {
		if errors.Cause(err) == ErrDisallowedIP {
			err = errors.Wrap(err, "connections to local resources are disabled by default, if you are sure this is safe, you can enable on a per-task basis by setting allowUnrestrictedNetworkAccess=true in the pipeline task spec")
		}
		return Result{Error: err}, RunInfo{IsRetryable: isRetryableHTTPError(statusCode, err)}
	}

	lggr.Debugw("HTTP task got response",
		"response", string(responseBytes),
		"url", url.String(),
		"dotID", t.DotID(),
	)

	promHTTPFetchTime.WithLabelValues(t.DotID()).Set(float64(elapsed))
	promHTTPResponseBodySize.WithLabelValues(t.DotID()).Set(float64(len(responseBytes)))

	// NOTE: We always stringify the response since this is required for all current jobs.
	// If a binary response is required we might consider adding an adapter
	// flag such as  "BinaryMode: true" which passes through raw binary as the
	// value instead.
	return Result{Value: string(responseBytes)}, runInfo
}
