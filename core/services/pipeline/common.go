package pipeline

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/store/models"
)

type (
	Type string

	Spec interface {
		JobID() *models.ID
		Type() Type
		Tasks() []Task
	}

	Service interface {
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
