package pipeline

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

//
// Return types:
//     string
//
type BridgeTask struct {
	BaseTask `mapstructure:",squash"`

	Name              string `json:"name"`
	RequestData       string `json:"requestData"`
	IncludeInputAtKey string `json:"includeInputAtKey"`
	Async             string `json:"async"`

	queryer    pg.Queryer
	config     Config
	httpClient *http.Client
}

var _ Task = (*BridgeTask)(nil)

var zeroURL = new(url.URL)

func (t *BridgeTask) Type() TaskType {
	return TaskTypeBridge
}

func (t *BridgeTask) Run(ctx context.Context, lggr logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	inputValues, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
	}

	var (
		name              StringParam
		requestData       MapParam
		includeInputAtKey StringParam
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&name, From(NonemptyString(t.Name))), "name"),
		errors.Wrap(ResolveParam(&requestData, From(VarExpr(t.RequestData, vars), JSONWithVarExprs(t.RequestData, vars, false), nil)), "requestData"),
		errors.Wrap(ResolveParam(&includeInputAtKey, From(t.IncludeInputAtKey)), "includeInputAtKey"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	url, err := t.getBridgeURLFromName(name)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	var metaMap MapParam

	meta, _ := vars.Get("jobRun.meta")
	switch v := meta.(type) {
	case map[string]interface{}:
		metaMap = MapParam(v)
	case nil:
	default:
		lggr.Warnw(`"meta" field on task run is malformed, discarding`,
			"task", t.DotID(),
			"meta", meta,
		)
	}

	requestData = withRunInfo(requestData, metaMap)
	if t.IncludeInputAtKey != "" {
		if len(inputValues) > 0 {
			requestData[string(includeInputAtKey)] = inputValues[0]
		}
	}

	if t.Async == "true" {
		responseURL := t.config.BridgeResponseURL()
		if *responseURL != *zeroURL {
			responseURL.Path = path.Join(responseURL.Path, "/v2/resume/", t.uuid.String())
		}
		requestData["responseURL"] = responseURL.String()
	}

	requestDataJSON, err := json.Marshal(requestData)
	if err != nil {
		return Result{Error: err}, runInfo
	}
	lggr.Debugw("Bridge task: sending request",
		"requestData", string(requestDataJSON),
		"url", url.String(),
	)

	requestCtx, cancel := httpRequestCtx(ctx, t, t.config)
	defer cancel()

	responseBytes, statusCode, headers, elapsed, err := makeHTTPRequest(requestCtx, lggr, "POST", URLParam(url), requestData, t.httpClient, t.config.DefaultHTTPLimit())
	if err != nil {
		return Result{Error: err}, RunInfo{IsRetryable: isRetryableHTTPError(statusCode, err)}
	}

	if t.Async == "true" {
		// Look for a `pending` flag. This check is case-insensitive because http.Header normalizes header names
		if _, ok := headers["X-Chainlink-Pending"]; ok {
			return result, pendingRunInfo()
		}

		var response struct {
			Pending bool `json:"pending"`
		}
		if err := json.Unmarshal(responseBytes, &response); err == nil && response.Pending {
			return Result{}, pendingRunInfo()
		}
	}

	// NOTE: We always stringify the response since this is required for all current jobs.
	// If a binary response is required we might consider adding an adapter
	// flag such as  "BinaryMode: true" which passes through raw binary as the
	// value instead.
	result = Result{Value: string(responseBytes)}

	promHTTPFetchTime.WithLabelValues(t.DotID()).Set(float64(elapsed))
	promHTTPResponseBodySize.WithLabelValues(t.DotID()).Set(float64(len(responseBytes)))

	lggr.Debugw("Bridge task: fetched answer",
		"answer", result.Value,
		"url", url.String(),
		"dotID", t.DotID(),
	)
	return result, runInfo
}

func (t BridgeTask) getBridgeURLFromName(name StringParam) (URLParam, error) {
	var bt bridges.BridgeType
	err := t.queryer.Get(&bt, "SELECT * FROM bridge_types WHERE name = $1", string(name))
	if err != nil {
		return URLParam{}, errors.Wrapf(err, "could not find bridge with name '%s'", name)
	}
	return URLParam(bt.URL), nil
}

func withRunInfo(request MapParam, meta MapParam) MapParam {
	output := make(MapParam)
	for k, v := range request {
		output[k] = v
	}
	if meta != nil {
		output["meta"] = meta
	}
	return output
}
