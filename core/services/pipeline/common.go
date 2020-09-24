package pipeline

import (
	"database/sql/driver"
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"

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
		DotID() string
		Run(inputs []Result) Result
		OutputTask() Task
		SetOutputTask(task Task)
		OutputIndex() int32
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
	outputTask Task   `json:"-"`
	dotID      string `json:"-" mapstructure:"-"`
	Index      int32  `json:"-" mapstructure:"index"`
}

func (t BaseTask) DotID() string                  { return t.dotID }
func (t BaseTask) OutputIndex() int32             { return t.Index }
func (t BaseTask) OutputTask() Task               { return t.outputTask }
func (t *BaseTask) SetOutputTask(outputTask Task) { t.outputTask = outputTask }

type JSONSerializable struct {
	Val interface{}
}

func (js *JSONSerializable) UnmarshalJSON(bs []byte) error {
	if js == nil {
		*js = JSONSerializable{}
	}
	return json.Unmarshal(bs, &js.Val)
}

func (js JSONSerializable) MarshalJSON() ([]byte, error) {
	return json.Marshal(js.Val)
}

func (js *JSONSerializable) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.Errorf("JSONSerializable#Scan received a value of type %T", value)
	}
	if js == nil {
		*js = JSONSerializable{}
	}
	return json.Unmarshal(bytes, &js.Val)
}

func (js JSONSerializable) Value() (driver.Value, error) {
	return json.Marshal(js.Val)
}

type TaskType string

const (
	TaskTypeHTTP      TaskType = "http"
	TaskTypeBridge    TaskType = "bridge"
	TaskTypeMedian    TaskType = "median"
	TaskTypeMultiply  TaskType = "multiply"
	TaskTypeJSONParse TaskType = "jsonparse"
)

func UnmarshalTaskFromMap(taskType TaskType, taskMap interface{}, dotID string, orm ORM, config Config) (Task, error) {
	switch taskMap.(type) {
	default:
		return nil, errors.New("UnmarshalTaskFromMap only accepts a map[string]interface{} or a map[string]string")
	case map[string]interface{}, map[string]string:
	}

	var task Task
	switch taskType {
	case TaskTypeHTTP:
		task = &HTTPTask{config: config, BaseTask: BaseTask{dotID: dotID}}
	case TaskTypeBridge:
		task = &BridgeTask{orm: orm, config: config, BaseTask: BaseTask{dotID: dotID}}
	case TaskTypeMedian:
		task = &MedianTask{BaseTask: BaseTask{dotID: dotID}}
	case TaskTypeJSONParse:
		task = &JSONParseTask{BaseTask: BaseTask{dotID: dotID}}
	case TaskTypeMultiply:
		task = &MultiplyTask{BaseTask: BaseTask{dotID: dotID}}
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

					case reflect.TypeOf(int32(0)):
						i, err := strconv.Atoi(data.(string))
						return int32(i), err
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
