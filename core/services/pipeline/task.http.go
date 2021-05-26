package pipeline

import (
	"context"

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
	RequestData                    string `json:"requestData" pipeline:"@expand_vars"`
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

func (t *HTTPTask) Run(ctx context.Context, _ JSONSerializable, inputs []Result) Result {
	var (
		method                         StringParam
		url                            URLParam
		requestData                    MapParam
		allowUnrestrictedNetworkAccess BoolParam
	)
	err := multierr.Combine(
		errors.Wrap(ResolveParam(&method, From(NonemptyString(t.Method), "GET")), "method"),
		errors.Wrap(ResolveParam(&url, From(NonemptyString(t.URL))), "url"),
		errors.Wrap(ResolveParam(&requestData, From(NonemptyString(t.RequestData), Input(inputs, 0), nil)), "requestData"),
		errors.Wrap(ResolveParam(&allowUnrestrictedNetworkAccess, From(NonemptyString(t.AllowUnrestrictedNetworkAccess), t.config.DefaultHTTPAllowUnrestrictedNetworkAccess())), "allowUnrestrictedNetworkAccess"),
	)
	if err != nil {
		return Result{Error: err}
	}

	responseBytes, elapsed, err := makeHTTPRequest(ctx, method, url, requestData, allowUnrestrictedNetworkAccess, t.config)
	if err != nil {
		return Result{Error: err}
	}

	logger.Debugw("HTTP task got response",
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
	return Result{Value: string(responseBytes)}
}
