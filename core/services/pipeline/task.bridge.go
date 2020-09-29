package pipeline

import (
	"net/url"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type BridgeTask struct {
	BaseTask `mapstructure:",squash"`

	Name        string          `json:"name"`
	RequestData HttpRequestData `json:"requestData"`

	orm    ORM
	config Config
}

var _ Task = (*BridgeTask)(nil)

func (t *BridgeTask) Type() TaskType {
	return TaskTypeBridge
}

func (t *BridgeTask) Run(inputs []Result) (result Result) {
	if len(inputs) > 0 {
		return Result{Error: errors.Wrapf(ErrWrongInputCardinality, "BridgeTask requires 0 inputs")}
	}

	url, err := t.getBridgeURLFromName()
	if err != nil {
		return Result{Error: err}
	}

	// add an arbitrary "id" field to the request json
	// this is done in order to keep request payloads consistent in format
	// between flux monitor polling requests and http/bridge adapters
	requestData := withIDAndMeta(t.RequestData, meta)

	result = (&HTTPTask{
		URL:         models.WebURL(url),
		Method:      "POST",
		RequestData: requestData,
		config:      t.config,
	}).Run(inputs)
	if result.Error != nil {
		return result
	}
	logger.Debugw("Bridge task: fetched answer",
		"answer", string(result.Value.([]byte)),
		"url", url.String(),
	)
	return result
}

func (t BridgeTask) getBridgeURLFromName() (url.URL, error) {
	task := models.TaskType(t.Name)
	bridge, err := t.orm.FindBridge(task)
	if err != nil {
		return url.URL{}, err
	}
	bridgeURL := url.URL(bridge.URL)
	return bridgeURL, nil
}

func withIDAndMeta(request, meta map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	for k, v := range request {
		output[k] = v
	}
	output["id"] = models.NewID()
	output["meta"] = meta
	return output
}
