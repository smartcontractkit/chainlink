package pipeline

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --name JobSpec --output ./mocks/ --case=underscore
//go:generate mockery --name JobService --output ./mocks/ --case=underscore
//go:generate mockery --name Task --output ./mocks/ --case=underscore

type (
	JobType string

	JobSpec interface {
		JobID() *models.ID
		JobType() JobType
		TaskDAG() TaskDAG
	}

	JobService interface {
		Start() error
		Stop() error
	}

	Task interface {
		Run(inputs []Result) Result
		OutputTasks() []Task
		SetOutputTasks(tasks []Task)
	}

	Result struct {
		Value interface{}
		Error error
	}
)

var (
	ErrWrongInputCardinality = errors.New("wrong number of task inputs")
)

type TaskType string

const (
	TaskTypeHTTP      TaskType = "http"
	TaskTypeBridge    TaskType = "bridge"
	TaskTypeMedian    TaskType = "median"
	TaskTypeMultiply  TaskType = "multiply"
	TaskTypeJSONParse TaskType = "jsonparse"
)

func NewTaskByType(taskType TaskType) (Task, error) {
	switch taskType {
	case TaskTypeHTTP:
		return &HTTPTask{}, nil
	case TaskTypeBridge:
		return &BridgeTask{}, nil
	case TaskTypeMedian:
		return &MedianTask{}, nil
	case TaskTypeJSONParse:
		return &JSONParseTask{}, nil
	case TaskTypeMultiply:
		return &MultiplyTask{}, nil
	default:
		return nil, errors.Errorf(`unknown task type: "%v"`, taskType)
	}
}

type BaseTask struct {
	outputTasks []Task `json:"-"`
}

func (t BaseTask) OutputTasks() []Task                { return t.outputTasks }
func (t *BaseTask) SetOutputTasks(outputTasks []Task) { t.outputTasks = outputTasks }

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
