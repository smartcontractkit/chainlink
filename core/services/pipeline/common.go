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
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
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
		ID() int64
		DotID() string
		Run(ctx context.Context, vars Vars, meta JSONSerializable, inputs []Result) Result
		Base() *BaseTask
		Outputs() []Task
		Inputs() []Task
		OutputTask() Task
		SetOutputTask(task Task)
		OutputIndex() int32
		TaskTimeout() (time.Duration, bool)
		NumPredecessors() int
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
		JobPipelineMaxRunDuration() time.Duration
		JobPipelineReaperInterval() time.Duration
		JobPipelineReaperThreshold() time.Duration
	}
)

var (
	ErrWrongInputCardinality = errors.New("wrong number of task inputs")
	ErrBadInput              = errors.New("bad input for task")
	ErrParameterEmpty        = errors.New("parameter is empty")
	ErrTooManyErrors         = errors.New("too many errors")
)

const (
	InputTaskKey = "input"
)

// Bundled tx and txmutex for multiple goroutines inside the same transaction.
// This mutex is necessary to work to avoid
// concurrent database calls inside the same transaction to fail.
// With the pq driver: `pq: unexpected Parse response 'C'`
// With the pgx driver: `conn busy`.
type SafeTx struct {
	tx   *gorm.DB
	txMu *sync.Mutex
}

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

// OutputsDB dumps a result output for a pipeline_run
func (result FinalResult) OutputsDB() JSONSerializable {
	return JSONSerializable{Val: result.Values, Null: false}
}

// ErrorsDB dumps a result error for a pipeline_run
func (result FinalResult) ErrorsDB() RunErrors {
	errStrs := make([]null.String, len(result.Errors))
	for i, err := range result.Errors {
		if err == nil {
			errStrs[i] = null.String{}
		} else {
			errStrs[i] = null.StringFrom(err.Error())
		}
	}

	return errStrs
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
	Task       Task
	TaskRun    TaskRun
	Result     Result
	CreatedAt  time.Time
	FinishedAt time.Time
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

type TaskType string

func (t TaskType) String() string {
	return string(t)
}

const (
	TaskTypeHTTP      TaskType = "http"
	TaskTypeBridge    TaskType = "bridge"
	TaskTypeMedian    TaskType = "median"
	TaskTypeMultiply  TaskType = "multiply"
	TaskTypeJSONParse TaskType = "jsonparse"
	TaskTypeAny       TaskType = "any"
	TaskTypeVRF       TaskType = "vrf"

	// Testing only.
	TaskTypePanic TaskType = "panic"
)

var (
	stringType = reflect.TypeOf("")
	int32Type  = reflect.TypeOf(int32(0))
)

func UnmarshalTaskFromMap(taskType TaskType, taskMap interface{}, ID int64, dotID string, config Config, txdb *gorm.DB, txdbMutex *sync.Mutex, numPredecessors int) (_ Task, err error) {
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
		task = &PanicTask{BaseTask: BaseTask{id: ID, dotID: dotID, numPredecessors: numPredecessors}}
	case TaskTypeHTTP:
		task = &HTTPTask{config: config, BaseTask: BaseTask{id: ID, dotID: dotID, numPredecessors: numPredecessors}}
	case TaskTypeBridge:
		task = &BridgeTask{config: config, safeTx: SafeTx{txdb, txdbMutex}, BaseTask: BaseTask{id: ID, dotID: dotID, numPredecessors: numPredecessors}}
	case TaskTypeMedian:
		task = &MedianTask{BaseTask: BaseTask{id: ID, dotID: dotID, numPredecessors: numPredecessors}}
	case TaskTypeAny:
		task = &AnyTask{BaseTask: BaseTask{id: ID, dotID: dotID, numPredecessors: numPredecessors}}
	case TaskTypeJSONParse:
		task = &JSONParseTask{BaseTask: BaseTask{id: ID, dotID: dotID, numPredecessors: numPredecessors}}
	case TaskTypeMultiply:
		task = &MultiplyTask{BaseTask: BaseTask{id: ID, dotID: dotID, numPredecessors: numPredecessors}}
	case TaskTypeVRF:
		task = &VRFTask{BaseTask: BaseTask{id: ID, dotID: dotID, numPredecessors: numPredecessors}}
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
