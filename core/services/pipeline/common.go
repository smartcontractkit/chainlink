package pipeline

import (
	"encoding/json"
	"net/url"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --name Task --output ./mocks/ --case=underscore
//go:generate mockery --name Config --output ./mocks/ --case=underscore

type (
	Task interface {
		Type() TaskType
		Run(inputs []Result) Result
		OutputTask() Task
		SetOutputTask(task Task)
	}

	Result struct {
		Value interface{}
		Error error
	}

	Config interface {
		DefaultHTTPTimeout() models.Duration
		DefaultMaxHTTPAttempts() uint
		DefaultHTTPLimit() int64
	}
)

var (
	ErrWrongInputCardinality = errors.New("wrong number of task inputs")
)

type BaseTask struct {
	outputTask Task `json:"-"`
}

func (t BaseTask) OutputTask() Task               { return t.outputTask }
func (t *BaseTask) SetOutputTask(outputTask Task) { t.outputTask = outputTask }

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

type TaskType string

const (
	TaskTypeHTTP      TaskType = "http"
	TaskTypeBridge    TaskType = "bridge"
	TaskTypeMedian    TaskType = "median"
	TaskTypeMultiply  TaskType = "multiply"
	TaskTypeJSONParse TaskType = "jsonparse"
)

func UnmarshalTask(taskType TaskType, taskMap interface{}, orm ORM, config Config) (Task, error) {
	switch taskMap.(type) {
	default:
		return nil, errors.New("UnmarshalTask only accepts a map[string]interface{} or a map[string]string")
	case map[string]interface{}, map[string]string:
	}

	var task Task
	switch taskType {
	case TaskTypeHTTP:
		task = &HTTPTask{config: config}
	case TaskTypeBridge:
		task = &BridgeTask{orm: orm, config: config}
	case TaskTypeMedian:
		task = &MedianTask{}
	case TaskTypeJSONParse:
		task = &JSONParseTask{}
	case TaskTypeMultiply:
		task = &MultiplyTask{}
	default:
		return nil, errors.Errorf(`unknown task type: "%v"`, taskType)
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToSliceHookFunc(","),
			func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
				switch f {
				case reflect.TypeOf(""):
					switch t {
					case reflect.TypeOf(models.WebURL{}):
						u, err := url.Parse(data.(string))
						if err != nil {
							return nil, err
						}
						return models.WebURL(*u), nil

					case reflect.TypeOf(HttpRequestData{}):
						var m map[string]interface{}
						err := json.Unmarshal([]byte(data.(string)), &m)
						return HttpRequestData(m), err

					case reflect.TypeOf(decimal.Decimal{}):
						return decimal.NewFromString(data.(string))
					}
				}
				return data, nil
			},
		),
		Result: task,
	})
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(taskMap)
	if err != nil {
		return nil, err
	}
	return task, nil
}
