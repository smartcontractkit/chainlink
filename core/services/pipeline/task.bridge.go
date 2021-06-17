package pipeline

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type BridgeTask struct {
	BaseTask `mapstructure:",squash"`

	Name              string `json:"name"`
	RequestData       string `json:"requestData"`
	IncludeInputAtKey string `json:"includeInputAtKey"`

	tx     *gorm.DB
	config Config
}

var _ Task = (*BridgeTask)(nil)

func (t *BridgeTask) Type() TaskType {
	return TaskTypeBridge
}

func (t *BridgeTask) Run(ctx context.Context, vars Vars, meta JSONSerializable, inputs []Result) Result {
	inputValues, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: err}
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
	switch v := meta.Val.(type) {
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
		logger.Warnw(`The "includeInputAtKey" parameter on Bridge tasks is deprecated. Please migrate to variable interpolation syntax as soon as possible (see CHANGELOG).`,
			"task", t.DotID(),
		)
		if len(inputValues) > 0 {
			requestData[string(includeInputAtKey)] = inputValues[0]
		}
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

	responseBytes, elapsed, err := makeHTTPRequest(ctx, "POST", URLParam(url), requestData, allowUnrestrictedNetworkAccess, t.config)
	if err != nil {
		return Result{Error: err}
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
	task := models.TaskType(name)

	bridge, err := FindBridge(t.tx, task)
	if err != nil {
		return URLParam{}, err
	}
	return URLParam(bridge.URL), nil
}

func withMeta(request MapParam, meta MapParam) MapParam {
	output := make(MapParam)
	for k, v := range request {
		output[k] = v
	}
	output["meta"] = meta
	return output
}
