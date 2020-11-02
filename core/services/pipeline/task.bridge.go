package pipeline

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type BridgeTask struct {
	BaseTask `mapstructure:",squash"`

	Name        string          `json:"name"`
	RequestData HttpRequestData `json:"requestData"`

	txdb   *gorm.DB
	config Config
}

var _ Task = (*BridgeTask)(nil)

func (t *BridgeTask) Type() TaskType {
	return TaskTypeBridge
}

func (t *BridgeTask) Run(ctx context.Context, taskRun TaskRun, inputs []Result) (result Result) {
	if len(inputs) > 0 {
		return Result{Error: errors.Wrapf(ErrWrongInputCardinality, "BridgeTask requires 0 inputs")}
	}

	url, err := t.getBridgeURLFromName()
	if err != nil {
		return Result{Error: err}
	}

	var meta map[string]interface{}
	switch v := taskRun.PipelineRun.Meta.Val.(type) {
	case map[string]interface{}:
		meta = v
	case nil:
	default:
		logger.Warnw(`"meta" field on task run is malformed, discarding`,
			"jobID", taskRun.PipelineRun.PipelineSpecID,
			"taskRunID", taskRun.ID,
			"task", taskRun.PipelineTaskSpec.DotID,
			"meta", taskRun.PipelineRun.Meta.Val,
		)
	}

	result = (&HTTPTask{
		URL:         models.WebURL(url),
		Method:      "POST",
		RequestData: withIDAndMeta(t.RequestData, taskRun.PipelineRunID, meta),
		config:      t.config,
	}).Run(ctx, taskRun, inputs)
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
	bridge, err := FindBridge(t.txdb, task)
	if err != nil {
		return url.URL{}, err
	}
	bridgeURL := url.URL(bridge.URL)
	return bridgeURL, nil
}

func withIDAndMeta(request HttpRequestData, runID int64, meta HttpRequestData) HttpRequestData {
	output := make(HttpRequestData)
	for k, v := range request {
		output[k] = v
	}
	output["id"] = fmt.Sprintf("%d", runID)
	output["meta"] = meta
	return output
}
