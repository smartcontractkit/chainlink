package job

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type (
	JobType string

	JobSpec interface {
		JobID() *models.ID
		JobType() JobType
	}

	JobService interface {
		Start() error
		Stop() error
	}

	Task interface {
		Run(inputs []Result) (interface{}, error)
		InputTasks() []Task
		SetInputTasks(tasks []Task)
	}
)

var (
	ErrWrongInputCardinality = errors.New("wrong number of task inputs")
)

type TaskType string

var (
	TaskTypeHttpFetcher          TaskType = "http"
	TaskTypeBridgeFetcher        TaskType = "bridge"
	TaskTypeMedianFetcher        TaskType = "median"
	TaskTypeMultiplyTransformer  TaskType = "multiply"
	TaskTypeJSONParseTransformer TaskType = "jsonparse"
)

func UnmarshalTaskJSON(bs []byte) (_ Task, err error) {
	var header struct {
		Type TaskType `json:"type"`
	}
	err = json.Unmarshal(bs, &header)
	if err != nil {
		return nil, err
	}

	var task Task
	switch header.Type {
	case TaskTypeBridgeFetcher:
		var bridgeFetcher BridgeFetcher
		err = json.Unmarshal(bs, &bridgeFetcher)
		if err != nil {
			return nil, err
		}
		task = &bridgeFetcher

	case TaskTypeHttpFetcher:
		var httpFetcher HttpFetcher
		err = json.Unmarshal(bs, &httpFetcher)
		if err != nil {
			return nil, err
		}
		task = &httpFetcher

	case TaskTypeMedianFetcher:
		var medianFetcher MedianFetcher
		err = json.Unmarshal(bs, &medianFetcher)
		if err != nil {
			return nil, err
		}
		task = &medianFetcher

	case TaskTypeJSONParseTransformer:
		var jsonTransformer JSONParseTransformer
		err = json.Unmarshal(bs, &jsonTransformer)
		if err != nil {
			return err
		}
		task = &jsonTransformer

	case TaskTypeMultiplyTransformer:
		var multiplyTransformer MultiplyTransformer
		err = json.Unmarshal(bs, &multiplyTransformer)
		if err != nil {
			return err
		}
		task = &multiplyTransformer

	default:
		return nil, errors.New("unknown fetcher type")
	}

	return fetcher, nil
}

type BaseTask struct {
	ID         uint64 `gorm:"primary_key;auto_increment"`
	inputTasks []Task `json:"-" gorm:"-"`
}

func (t BaseTask) InputTasks() []Task               { return t.inputTasks }
func (t *BaseTask) SetInputTasks(inputTasks []Task) { t.inputTasks = inputTasks }

type JSONSerializable struct {
	Value interface{}
}

func (js *JSONSerializable) UnmarshalJSON(bs []byte) error {
	if js == nil {
		*js = JSONSerializable{}
	}
	return json.Unmarshal(bs, &js.Value)
}

func (js JSONSerializable) MarshalJSON() ([]byte, error) {
	return json.Marshal(js.Value)
}
