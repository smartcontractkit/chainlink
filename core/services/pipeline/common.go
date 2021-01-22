package pipeline

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gopkg.in/guregu/null.v4"
)

//go:generate mockery --name Config --output ./mocks/ --case=underscore

type (
	Task interface {
		Type() TaskType
		DotID() string
		Run(ctx context.Context, taskRun TaskRun, inputs []Result) Result
		OutputTask() Task
		SetOutputTask(task Task)
		OutputIndex() int32
		TaskTimeout() (time.Duration, bool)
		SetDefaults(inputValues map[string]string, g TaskDAG, self taskDAGNode) error
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
		JobPipelineMaxRunDuration() time.Duration
		JobPipelineParallelism() uint8
		JobPipelineReaperInterval() time.Duration
		JobPipelineReaperThreshold() time.Duration
	}
)

var (
	ErrWrongInputCardinality = errors.New("wrong number of task inputs")
	ErrBadInput              = errors.New("bad input for task")
)

type Result struct {
	Value interface{}
	Error error
}

// OutputDB dumps a single result output for a pipeline_run or pipeline_task_run
func (result Result) OutputDB() JSONSerializable {
	return JSONSerializable{Val: result.Value, Null: result.Value == nil}
}

// ErrorDB dumps a single result error for a pipeline_task_run
func (result Result) ErrorDB() null.String {
	var errString null.String
	if finalErrors, is := result.Error.(FinalErrors); is {
		errString = null.StringFrom(finalErrors.Error())
	} else if result.Error != nil {
		errString = null.StringFrom(result.Error.Error())
	}
	return errString
}

// ErrorsDB dumps a result error for a pipeline_run
func (result Result) ErrorsDB() JSONSerializable {
	var val interface{}
	if finalErrors, is := result.Error.(FinalErrors); is {
		val = finalErrors
	} else if result.Error != nil {
		val = result.Error.Error()
	} else {
		val = nil
	}
	return JSONSerializable{Val: val, Null: false}
}

// TaskRunResult describes the result of a task run, suitable for database
// update or insert
type TaskRunResult struct {
	ID         int64
	Result     Result
	FinishedAt time.Time
	IsFinal    bool
}

type BaseTask struct {
	outputTask Task
	dotID      string        `mapstructure:"-"`
	Index      int32         `mapstructure:"index" json:"-" `
	Timeout    time.Duration `mapstructure:"timeout"`
}

func (t BaseTask) DotID() string                  { return t.dotID }
func (t BaseTask) OutputIndex() int32             { return t.Index }
func (t BaseTask) OutputTask() Task               { return t.outputTask }
func (t *BaseTask) SetOutputTask(outputTask Task) { t.outputTask = outputTask }
func (t BaseTask) TaskTimeout() (time.Duration, bool) {
	if t.Timeout == time.Duration(0) {
		return time.Duration(0), false
	}
	return t.Timeout, true
}

type JSONSerializable struct {
	Val  interface{}
	Null bool
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
	if value == nil {
		*js = JSONSerializable{Null: true}
		return nil
	}
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
	if js.Null {
		return nil, nil
	}
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

func UnmarshalTaskFromMap(taskType TaskType, taskMap interface{}, dotID string, config Config, txdb *gorm.DB, txdbMutex *sync.Mutex) (_ Task, err error) {
	defer utils.WrapIfError(&err, "UnmarshalTaskFromMap")

	switch taskMap.(type) {
	default:
		return nil, errors.Errorf("UnmarshalTaskFromMap only accepts a map[string]interface{} or a map[string]string. Got %v (%#v) of type %T", taskMap, taskMap, taskMap)
	case map[string]interface{}, map[string]string:
	}

	taskType = TaskType(strings.ToLower(string(taskType)))

	var task Task
	switch taskType {
	case TaskTypeHTTP:
		task = &HTTPTask{config: config, BaseTask: BaseTask{dotID: dotID}}
	case TaskTypeBridge:
		task = &BridgeTask{config: config, txdb: txdb, txdbMutex: txdbMutex, BaseTask: BaseTask{dotID: dotID}}
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
			mapstructure.StringToTimeDurationHookFunc(),
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
					case reflect.TypeOf(uint32(0)):
						i, err2 := strconv.ParseInt(data.(string), 10, 32)
						return uint32(i), err2
					case reflect.TypeOf(int64(0)):
						i, err2 := strconv.ParseInt(data.(string), 10, 64)
						return uint32(i), err2
					case reflect.TypeOf(uint64(0)):
						i, err2 := strconv.ParseInt(data.(string), 10, 64)
						return uint64(i), err2
					case reflect.TypeOf(true):
						b, err2 := strconv.ParseBool(data.(string))
						return b, err2
					case reflect.TypeOf(MaybeBool("")):
						return MaybeBoolFromString(data.(string))
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
