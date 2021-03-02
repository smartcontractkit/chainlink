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

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
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
		DatabaseURL() url.URL
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

// Result is the result of a TaskRun
type Result struct {
	Value interface{}
	Error error
}

// FinalResult is the result of a Run
// TODO: Get rid of FinalErrors and use FinalResult instead
// https://www.pivotaltracker.com/story/show/176557536
type FinalResult struct {
	Values []interface{}
	Errors []error
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

// OutputsDB dumps a result output for a pipeline_run
func (result FinalResult) OutputsDB() JSONSerializable {
	return JSONSerializable{Val: result.Values, Null: false}
}

// ErrorsDB dumps a result error for a pipeline_run
func (result FinalResult) ErrorsDB() JSONSerializable {
	errStrs := make([]null.String, len(result.Errors))
	for i, err := range result.Errors {
		if err == nil {
			errStrs[i] = null.String{}
		} else {
			errStrs[i] = null.StringFrom(err.Error())
		}
	}

	return JSONSerializable{Val: errStrs, Null: false}
}

// HasErrors returns true if the final result has any errors
func (result FinalResult) HasErrors() bool {
	for _, err := range result.Errors {
		if err != nil {
			return true
		}
	}
	return false
}

// SingularResult returns a single result if the FinalResult only has one set of outputs/errors
func (result FinalResult) SingularResult() (Result, error) {
	if len(result.Errors) != 1 || len(result.Values) != 1 {
		return Result{}, errors.Errorf("cannot cast FinalResult to singular result; it does not have exactly 1 error and exactly 1 output: %#v", result)
	}
	return Result{Error: result.Errors[0], Value: result.Values[0]}, nil
}

// TaskRunResult describes the result of a task run, suitable for database
// update or insert.
// ID might be zero if the TaskRun has not been inserted yet
// TaskSpecID will always be non-zero
type TaskRunResult struct {
	ID         int64
	TaskSpecID int32
	Result     Result
	FinishedAt time.Time
	IsTerminal bool
}

// TaskRunResults represents a collection of results for all task runs for one pipeline run
type TaskRunResults []TaskRunResult

// FinalResult pulls the FinalResult for the pipeline_run from the task runs
func (trrs TaskRunResults) FinalResult() (result FinalResult) {
	var found bool
	for _, trr := range trrs {
		if trr.IsTerminal {
			// FIXME: This is a mess because of the special `__result__` task.
			// It gets much simpler and will change when the magical
			// "__result__" type is removed.
			// https://www.pivotaltracker.com/story/show/176557536
			values, is := trr.Result.Value.([]interface{})
			if !is {
				panic("expected terminal task run result to have multiple values")
			}
			result.Values = append(result.Values, values...)

			finalErrs, is := trr.Result.Error.(FinalErrors)
			if !is {
				panic("expected terminal task run result to be FinalErrors")
			}
			errs := make([]error, len(finalErrs))
			for i, finalErr := range finalErrs {
				if finalErr.IsZero() {
					errs[i] = nil
				} else {
					errs[i] = errors.New(finalErr.ValueOrZero())
				}
			}
			result.Errors = append(result.Errors, errs...)
			found = true
		}
	}

	if !found {
		logger.Errorw("expected at least one task to be final", "tasks", trrs)
		panic("expected at least one task to be final")
	}
	return
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
	TaskTypeAny       TaskType = "any"
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
	case TaskTypeAny:
		task = &AnyTask{BaseTask: BaseTask{dotID: dotID}}
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
