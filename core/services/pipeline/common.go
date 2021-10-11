package pipeline

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gopkg.in/guregu/null.v4"
)

//go:generate mockery --name Config --output ./mocks/ --case=underscore

type (
	Task interface {
		Type() TaskType
		ID() int
		DotID() string
		Run(ctx context.Context, vars Vars, inputs []Result) Result
		Base() *BaseTask
		Outputs() []Task
		Inputs() []Task
		OutputIndex() int32
		TaskTimeout() (time.Duration, bool)
	}

	Config interface {
		BridgeResponseURL() *url.URL
		DatabaseMaximumTxDuration() time.Duration
		DatabaseURL() url.URL
		DefaultHTTPLimit() int64
		DefaultHTTPTimeout() models.Duration
		DefaultMaxHTTPAttempts() uint
		DefaultHTTPAllowUnrestrictedNetworkAccess() bool
		EthGasLimitDefault() uint64
		EthMaxQueuedTransactions() uint64
		TriggerFallbackDBPollInterval() time.Duration
		JobPipelineMaxRunDuration() time.Duration
		JobPipelineReaperInterval() time.Duration
		JobPipelineReaperThreshold() time.Duration
	}
)

var (
	ErrWrongInputCardinality = errors.New("wrong number of task inputs")
	ErrBadInput              = errors.New("bad input for task")
	ErrInputTaskErrored      = errors.New("input task errored")
	ErrParameterEmpty        = errors.New("parameter is empty")
	ErrTooManyErrors         = errors.New("too many errors")
	ErrTimeout               = errors.New("timeout")
	ErrTaskRunFailed         = errors.New("task run failed")
)

const (
	InputTaskKey = "input"
)

// Result is the result of a TaskRun
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
	if result.Error != nil {
		errString = null.StringFrom(result.Error.Error())
	}
	return errString
}

// FinalResult is the result of a Run
type FinalResult struct {
	Values []interface{}
	Errors []error
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
	ID         uuid.UUID
	Task       Task
	TaskRun    TaskRun
	Result     Result
	CreatedAt  time.Time
	FinishedAt null.Time
}

func (result *TaskRunResult) IsPending() bool {
	return !result.FinishedAt.Valid && result.Result == Result{}
}

func (result *TaskRunResult) IsTerminal() bool {
	return len(result.Task.Outputs()) == 0
}

// TaskRunResults represents a collection of results for all task runs for one pipeline run
type TaskRunResults []TaskRunResult

// FinalResult pulls the FinalResult for the pipeline_run from the task runs
// It needs to respect the output index of each task
func (trrs TaskRunResults) FinalResult() FinalResult {
	var found bool
	var fr FinalResult
	sort.Slice(trrs, func(i, j int) bool {
		return trrs[i].Task.OutputIndex() < trrs[j].Task.OutputIndex()
	})
	for _, trr := range trrs {
		if trr.IsTerminal() {
			fr.Values = append(fr.Values, trr.Result.Value)
			fr.Errors = append(fr.Errors, trr.Result.Error)
			found = true
		}
	}

	if !found {
		logger.Errorw("expected at least one task to be final", "tasks", trrs)
		panic("expected at least one task to be final")
	}
	return fr
}

type RunWithResults struct {
	Run            Run
	TaskRunResults TaskRunResults
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

func (js *JSONSerializable) Empty() bool {
	return js == nil || js.Null
}

type TaskType string

func (t TaskType) String() string {
	return string(t)
}

const (
	TaskTypeHTTP            TaskType = "http"
	TaskTypeBridge          TaskType = "bridge"
	TaskTypeMean            TaskType = "mean"
	TaskTypeMedian          TaskType = "median"
	TaskTypeMode            TaskType = "mode"
	TaskTypeSum             TaskType = "sum"
	TaskTypeMultiply        TaskType = "multiply"
	TaskTypeDivide          TaskType = "divide"
	TaskTypeJSONParse       TaskType = "jsonparse"
	TaskTypeCBORParse       TaskType = "cborparse"
	TaskTypeAny             TaskType = "any"
	TaskTypeVRF             TaskType = "vrf"
	TaskTypeETHCall         TaskType = "ethcall"
	TaskTypeETHTx           TaskType = "ethtx"
	TaskTypeETHABIEncode    TaskType = "ethabiencode"
	TaskTypeETHABIDecode    TaskType = "ethabidecode"
	TaskTypeETHABIDecodeLog TaskType = "ethabidecodelog"
	TaskTypeMerge           TaskType = "merge"

	// Testing only.
	TaskTypePanic TaskType = "panic"
)

var (
	stringType  = reflect.TypeOf("")
	bytesType   = reflect.TypeOf([]byte(nil))
	bytes20Type = reflect.TypeOf([20]byte{})
	int32Type   = reflect.TypeOf(int32(0))
)

func UnmarshalTaskFromMap(taskType TaskType, taskMap interface{}, ID int, dotID string) (_ Task, err error) {
	defer utils.WrapIfError(&err, "UnmarshalTaskFromMap")

	switch taskMap.(type) {
	default:
		return nil, errors.Errorf("UnmarshalTaskFromMap only accepts a map[string]interface{} or a map[string]string. Got %v (%#v) of type %T", taskMap, taskMap, taskMap)
	case map[string]interface{}, map[string]string:
	}

	taskType = TaskType(strings.ToLower(string(taskType)))

	var task Task
	switch taskType {
	case TaskTypePanic:
		task = &PanicTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeHTTP:
		task = &HTTPTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeBridge:
		task = &BridgeTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeMean:
		task = &MeanTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeMedian:
		task = &MedianTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeMode:
		task = &ModeTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeSum:
		task = &SumTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeAny:
		task = &AnyTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeJSONParse:
		task = &JSONParseTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeMultiply:
		task = &MultiplyTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeDivide:
		task = &DivideTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeVRF:
		task = &VRFTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeETHCall:
		task = &ETHCallTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeETHTx:
		task = &ETHTxTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeETHABIEncode:
		task = &ETHABIEncodeTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeETHABIDecode:
		task = &ETHABIDecodeTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeETHABIDecodeLog:
		task = &ETHABIDecodeLogTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeCBORParse:
		task = &CBORParseTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeMerge:
		task = &MergeTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	default:
		return nil, errors.Errorf(`unknown task type: "%v"`, taskType)
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: task,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
				switch from {
				case stringType:
					switch to {
					case int32Type:
						i, err2 := strconv.ParseInt(data.(string), 10, 32)
						return int32(i), err2
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

func CheckInputs(inputs []Result, minLen, maxLen, maxErrors int) ([]interface{}, error) {
	if minLen >= 0 && len(inputs) < minLen {
		return nil, errors.Wrapf(ErrWrongInputCardinality, "min: %v max: %v (got %v)", minLen, maxLen, len(inputs))
	} else if maxLen >= 0 && len(inputs) > maxLen {
		return nil, errors.Wrapf(ErrWrongInputCardinality, "min: %v max: %v (got %v)", minLen, maxLen, len(inputs))
	}
	var vals []interface{}
	var errs int
	for _, input := range inputs {
		if input.Error != nil {
			errs++
			continue
		}
		vals = append(vals, input.Value)
	}
	if maxErrors >= 0 && errs > maxErrors {
		return nil, ErrTooManyErrors
	}
	return vals, nil
}
