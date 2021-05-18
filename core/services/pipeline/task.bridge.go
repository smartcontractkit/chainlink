package pipeline

import (
	"context"
	"fmt"
	"net/url"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"go.uber.org/multierr"
)

type BridgeTask struct {
	BaseTask `mapstructure:",squash"`

	Name              string `json:"name"`
	RequestData       string `json:"requestData"`
	IncludeInputAtKey string `json:"includeInputAtKey"`

	safeTx SafeTx
	config Config
}

var _ Task = (*BridgeTask)(nil)

func (t *BridgeTask) Type() TaskType {
	return TaskTypeBridge
}

func (t *BridgeTask) SetDefaults(inputValues map[string]string, g TaskDAG, self TaskDAGNode) error {
	return nil
}

func (t *BridgeTask) Run(ctx context.Context, vars Vars, meta JSONSerializable, inputs []Result) (result Result) {
	inputValues, err := CheckInputs(inputs, 0, 1, 0)
	if err != nil {
		return Result{Error: err}
	}

	var (
		name              StringParam
		requestData       MapParam
		includeInputAtKey StringParam
	)
	err = multierr.Combine(
		vars.ResolveValue(&name, From(NonemptyString(t.Name))),
		vars.ResolveValue(&requestData, From(NonemptyString(t.RequestData))),
		vars.ResolveValue(&includeInputAtKey, From(NonemptyString(t.IncludeInputAtKey))),
	)
	if err != nil {
		return Result{Error: err}
	}

	url, err := t.getBridgeURLFromName(name)
	if err != nil {
		return Result{Error: err}
	}
	fmt.Println(url)

	var metaMap map[string]interface{}
	switch v := meta.Val.(type) {
	case map[string]interface{}:
		metaMap = v
	case nil:
	default:
		logger.Warnw(`"meta" field on task run is malformed, discarding`,
			"task", t.DotID(),
			"meta", meta,
		)
	}

	requestData = withMeta(requestData, metaMap)
	if t.IncludeInputAtKey != "" && len(inputValues) > 0 {
		requestData[string(includeInputAtKey)] = inputValues[0]
	}

	// result = (&HTTPTask{
	// 	URL:         models.WebURL(url),
	// 	Method:      "POST",
	// 	RequestData: requestData,
	// 	// URL is "safe" because it comes from the node's own database
	// 	// Some node operators may run external adapters on their own hardware
	// 	AllowUnrestrictedNetworkAccess: MaybeBoolTrue,
	// 	config:                         t.config,
	// }).Run(ctx, meta, nil)
	// if result.Error != nil {
	// 	return result
	// }
	// logger.Debugw("Bridge task: fetched answer",
	// 	"answer", result.Value,
	// 	"url", url.String(),
	// )
	return result
}

func (t BridgeTask) getBridgeURLFromName(name StringParam) (url.URL, error) {
	task := models.TaskType(name)

	if t.safeTx.txMu != nil {
		t.safeTx.txMu.Lock()
		defer t.safeTx.txMu.Unlock()
	}

	bridge, err := FindBridge(t.safeTx.tx, task)
	if err != nil {
		return url.URL{}, err
	}
	bridgeURL := url.URL(bridge.URL)
	return bridgeURL, nil
}

func withMeta(request MapParam, meta MapParam) MapParam {
	output := make(MapParam)
	for k, v := range request {
		output[k] = v
	}
	output["meta"] = meta
	return output
}
