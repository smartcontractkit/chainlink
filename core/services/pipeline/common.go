package pipeline

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name Task --output ./mocks/ --case=underscore
//go:generate mockery --name Config --output ./mocks/ --case=underscore

type (
	Task interface {
		Type() TaskType
		DotID() string
		Run(ctx context.Context, taskRun TaskRun, inputs []Result) Result
		OutputTask() Task
		SetOutputTask(task Task)
		OutputIndex() int32
	}

	Result struct {
		Value interface{}
		Error error
	}

	Config interface {
		BridgeResponseURL() *url.URL
		DatabaseMaximumTxDuration() time.Duration
		DatabaseURL() string
		DefaultHTTPLimit() int64
		DefaultHTTPTimeout() models.Duration
		DefaultMaxHTTPAttempts() uint
		DefaultHTTPAllowUnrestrictedNetworkAccess() bool
		TriggerFallbackDBPollInterval() time.Duration
		JobPipelineMaxTaskDuration() time.Duration
		JobPipelineParallelism() uint8
		JobPipelineReaperInterval() time.Duration
		JobPipelineReaperThreshold() time.Duration
	}
)

var (
	ErrWrongInputCardinality = errors.New("wrong number of task inputs")
	ErrBadInput              = errors.New("bad input for task")
)

type BaseTask struct {
	outputTask Task
	dotID      string `mapstructure:"-"`
	Index      int32  `mapstructure:"index" json:"-" `
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
	switch x := js.Val.(type) {
	case []byte:
		return json.Marshal(string(x))
	default:
		return json.Marshal(js.Val)
	}
}

func (js *JSONSerializable) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.Errorf("JSONSerializable#Scan received a value of type %T", value)
	}
	if js == nil {
		*js = JSONSerializable{}
	}
	return js.UnmarshalJSON(bytes)
}

func (js JSONSerializable) Value() (driver.Value, error) {
	return js.MarshalJSON()
}

type TaskType string

const (
	TaskTypeHTTP      TaskType = "http"
	TaskTypeBridge    TaskType = "bridge"
	TaskTypeMedian    TaskType = "median"
	TaskTypeMultiply  TaskType = "multiply"
	TaskTypeJSONParse TaskType = "jsonparse"
	TaskTypeResult    TaskType = "result"
)

const ResultTaskDotID = "__result__"

func UnmarshalTaskFromMap(taskType TaskType, taskMap interface{}, dotID string, config Config, txdb *gorm.DB) (_ Task, err error) {
	defer utils.WrapIfError(&err, "UnmarshalTaskFromMap")

	switch taskMap.(type) {
	default:
		return nil, errors.New("UnmarshalTaskFromMap only accepts a map[string]interface{} or a map[string]string")
	case map[string]interface{}, map[string]string:
	}

	taskType = TaskType(strings.ToLower(string(taskType)))

	var task Task
	switch taskType {
	case TaskTypeHTTP:
		task = &HTTPTask{config: config, BaseTask: BaseTask{dotID: dotID}}
	case TaskTypeBridge:
		task = &BridgeTask{config: config, txdb: txdb, BaseTask: BaseTask{dotID: dotID}}
	case TaskTypeMedian:
		task = &MedianTask{BaseTask: BaseTask{dotID: dotID}}
	case TaskTypeJSONParse:
		task = &JSONParseTask{BaseTask: BaseTask{dotID: dotID}}
	case TaskTypeMultiply:
		task = &MultiplyTask{BaseTask: BaseTask{dotID: dotID}}
	case TaskTypeResult:
		task = &ResultTask{BaseTask: BaseTask{dotID: ResultTaskDotID}}
	default:
		return nil, errors.Errorf(`unknown task type: "%v"`, taskType)
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: task,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToSliceHookFunc(","),
			func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
				switch f {
				case reflect.TypeOf(""):
					switch t {
					case reflect.TypeOf(models.WebURL{}):
						u, err2 := url.Parse(data.(string))
						if err2 != nil {
							return nil, err2
						}
						return models.WebURL(*u), nil

					case reflect.TypeOf(HttpRequestData{}):
						var m map[string]interface{}
						err2 := json.Unmarshal([]byte(data.(string)), &m)
						return HttpRequestData(m), err2

					case reflect.TypeOf(decimal.Decimal{}):
						return decimal.NewFromString(data.(string))

					case reflect.TypeOf(int32(0)):
						i, err2 := strconv.ParseInt(data.(string), 10, 32)
						return int32(i), err2
					case reflect.TypeOf(true):
						b, err2 := strconv.ParseBool(data.(string))
						return b, err2
					}
				}
				return data, nil
			},
		),
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

func WrapResultIfError(result *Result, msg string, args ...interface{}) {
	if result.Error != nil {
		logger.Errorf(msg+": %+v", append(args, result.Error)...)
		result.Error = errors.Wrapf(result.Error, msg, args...)
	}
}

type HttpRequestData map[string]interface{}

func (h *HttpRequestData) Scan(value interface{}) error { return json.Unmarshal(value.([]byte), h) }
func (h HttpRequestData) Value() (driver.Value, error)  { return json.Marshal(h) }
func (h HttpRequestData) AsMap() map[string]interface{} { return h }
