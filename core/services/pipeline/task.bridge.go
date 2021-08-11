package pipeline

import (
	"context"
	"encoding/json"
	"net/url"
	"path"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
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

	db     *gorm.DB
	config Config
}

var _ Task = (*BridgeTask)(nil)

var zeroURL = new(url.URL)

func (t *BridgeTask) Type() TaskType {
	return TaskTypeBridge
}

var ErrPending = errors.New("pending")

func (t *BridgeTask) Run(ctx context.Context, vars Vars, inputs []Result) Result {
	inputValues, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}
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
		return Result{Error: err}
	}

	url, err := t.getBridgeURLFromName(name)
	if err != nil {
		return Result{Error: err}
	}

	var metaMap MapParam

	meta, _ := vars.Get("jobRun.meta")
	switch v := meta.(type) {
	case map[string]interface{}:
		metaMap = MapParam(v)
	case nil:
	default:
		logger.Warnw(`"meta" field on task run is malformed, discarding`,
			"task", t.DotID(),
			"meta", meta,
		)
	}

	requestData = withMeta(requestData, metaMap)
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

	// URL is "safe" because it comes from the node's own database
	// Some node operators may run external adapters on their own hardware
	allowUnrestrictedNetworkAccess := BoolParam(true)

	requestDataJSON, err := json.Marshal(requestData)
	if err != nil {
		return Result{Error: err}
	}
	logger.Debugw("Bridge task: sending request",
		"requestData", string(requestDataJSON),
		"url", url.String(),
	)

	responseBytes, headers, elapsed, err := makeHTTPRequest(ctx, "POST", URLParam(url), requestData, allowUnrestrictedNetworkAccess, t.config)
	if err != nil {
		return Result{Error: err}
	}

	if t.Async == "true" {
		// Look for a `pending` flag. This check is case-insensitive because http.Header normalizes header names
		if _, ok := headers["X-Chainlink-Pending"]; ok {
			return Result{Error: ErrPending}
		}

		var response struct {
			Pending bool `json:"pending"`
		}
		if err := json.Unmarshal(responseBytes, &response); err == nil && response.Pending {
			return Result{Error: ErrPending}
		}
	}

	// NOTE: We always stringify the response since this is required for all current jobs.
	// If a binary response is required we might consider adding an adapter
	// flag such as  "BinaryMode: true" which passes through raw binary as the
	// value instead.
	result := Result{Value: string(responseBytes)}

	promHTTPFetchTime.WithLabelValues(t.DotID()).Set(float64(elapsed))
	promHTTPResponseBodySize.WithLabelValues(t.DotID()).Set(float64(len(responseBytes)))

	logger.Debugw("Bridge task: fetched answer",
		"answer", result.Value,
		"url", url.String(),
		"dotID", t.DotID(),
	)
	return result
}

func (t BridgeTask) getBridgeURLFromName(name StringParam) (URLParam, error) {
	var bt models.BridgeType
	err := t.db.First(&bt, "name = ?", string(name)).Error
	if err != nil {
		return URLParam{}, errors.Wrapf(err, "could not find bridge with name '%s'", name)
	}
	return URLParam(bt.URL), nil
}

func withMeta(request MapParam, meta MapParam) MapParam {
	output := make(MapParam)
	for k, v := range request {
		output[k] = v
	}
	if meta != nil {
		output["meta"] = meta
	}
	return output
}
