package pipeline

import (
	"context"
	"errors"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	pkgerrors "github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	cutils "github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	cnull "github.com/smartcontractkit/chainlink/v2/core/null"
)

const (
	BlockHeaderFeederJobType       string = "blockheaderfeeder"
	BlockhashStoreJobType          string = "blockhashstore"
	BootstrapJobType               string = "bootstrap"
	CronJobType                    string = "cron"
	DirectRequestJobType           string = "directrequest"
	FluxMonitorJobType             string = "fluxmonitor"
	GatewayJobType                 string = "gateway"
	KeeperJobType                  string = "keeper"
	LegacyGasStationServerJobType  string = "legacygasstationserver"
	LegacyGasStationSidecarJobType string = "legacygasstationsidecar"
	OffchainReporting2JobType      string = "offchainreporting2"
	OffchainReportingJobType       string = "offchainreporting"
	StreamJobType                  string = "stream"
	VRFJobType                     string = "vrf"
	WebhookJobType                 string = "webhook"
	WorkflowJobType                string = "workflow"
)

//go:generate mockery --quiet --name Config --output ./mocks/ --case=underscore

type (
	Task interface {
		Type() TaskType
		ID() int
		DotID() string
		Run(ctx context.Context, lggr logger.Logger, vars Vars, inputs []Result) (Result, RunInfo)
		Base() *BaseTask
		Outputs() []Task
		Inputs() []TaskDependency
		OutputIndex() int32
		TaskTimeout() (time.Duration, bool)
		TaskRetries() uint32
		TaskMinBackoff() time.Duration
		TaskMaxBackoff() time.Duration
	}

	Config interface {
		DefaultHTTPLimit() int64
		DefaultHTTPTimeout() commonconfig.Duration
		MaxRunDuration() time.Duration
		ReaperInterval() time.Duration
		ReaperThreshold() time.Duration
		VerboseLogging() bool
	}

	BridgeConfig interface {
		BridgeResponseURL() *url.URL
		BridgeCacheTTL() time.Duration
	}
)

// Wraps the input Task for the given dependent task along with a bool variable PropagateResult,
// which Indicates whether result of InputTask should be propagated to its dependent task.
// If the edge between these tasks was an implicit edge, then results are not propagated. This is because
// some tasks cannot handle an input from an edge which wasn't specified in the spec.
type TaskDependency struct {
	PropagateResult bool
	InputTask       Task
}

var (
	ErrWrongInputCardinality = errors.New("wrong number of task inputs")
	ErrBadInput              = errors.New("bad input for task")
	ErrInputTaskErrored      = errors.New("input task errored")
	ErrParameterEmpty        = errors.New("parameter is empty")
	ErrIndexOutOfRange       = errors.New("index out of range")
	ErrTooManyErrors         = errors.New("too many errors")
	ErrTimeout               = errors.New("timeout")
	ErrTaskRunFailed         = errors.New("task run failed")
	ErrCancelled             = errors.New("task run cancelled (fail early)")
)

const (
	InputTaskKey = "input"
)

// RunInfo contains additional information about the finished TaskRun
type RunInfo struct {
	IsRetryable bool
	IsPending   bool
}

// retryableMeta should be returned if the error is non-deterministic; i.e. a
// repeated attempt sometime later _might_ succeed where the current attempt
// failed
func retryableRunInfo() RunInfo {
	return RunInfo{IsRetryable: true}
}

func pendingRunInfo() RunInfo {
	return RunInfo{IsPending: true}
}

func isRetryableHTTPError(statusCode int, err error) bool {
	if statusCode >= 400 && statusCode < 500 {
		// Client errors are not likely to succeed by resubmitting the exact same information again
		return false
	} else if statusCode >= 500 {
		// Remote errors _might_ work on a retry
		return true
	}
	return err != nil
}

// Result is the result of a TaskRun
type Result struct {
	Value interface{}
	Error error
}

// OutputDB dumps a single result output for a pipeline_run or pipeline_task_run
func (result Result) OutputDB() jsonserializable.JSONSerializable {
	return jsonserializable.JSONSerializable{Val: result.Value, Valid: !(result.Value == nil || (reflect.ValueOf(result.Value).Kind() == reflect.Ptr && reflect.ValueOf(result.Value).IsNil()))}
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
	Values      []interface{}
	AllErrors   []error
	FatalErrors []error
}

// HasFatalErrors returns true if the final result has any errors
func (result FinalResult) HasFatalErrors() bool {
	for _, err := range result.FatalErrors {
		if err != nil {
			return true
		}
	}
	return false
}

// HasErrors returns true if the final result has any errors
func (result FinalResult) HasErrors() bool {
	for _, err := range result.AllErrors {
		if err != nil {
			return true
		}
	}
	return false
}

func (result FinalResult) CombinedError() error {
	if !result.HasErrors() {
		return nil
	}
	return errors.Join(result.AllErrors...)
}

// SingularResult returns a single result if the FinalResult only has one set of outputs/errors
func (result FinalResult) SingularResult() (Result, error) {
	if len(result.FatalErrors) != 1 || len(result.Values) != 1 {
		return Result{}, pkgerrors.Errorf("cannot cast FinalResult to singular result; it does not have exactly 1 error and exactly 1 output: %#v", result)
	}
	return Result{Error: result.FatalErrors[0], Value: result.Values[0]}, nil
}

// TaskRunResult describes the result of a task run, suitable for database
// update or insert.
// ID might be zero if the TaskRun has not been inserted yet
// TaskSpecID will always be non-zero
type TaskRunResult struct {
	ID         uuid.UUID
	Task       Task    `json:"-"`
	TaskRun    TaskRun `json:"-"`
	Result     Result
	Attempts   uint
	CreatedAt  time.Time
	FinishedAt null.Time
	// runInfo is never persisted
	runInfo RunInfo
}

func (result *TaskRunResult) IsPending() bool {
	return !result.FinishedAt.Valid && result.Result == Result{}
}

func (result *TaskRunResult) IsTerminal() bool {
	return len(result.Task.Outputs()) == 0
}

// TaskRunResults represents a collection of results for all task runs for one pipeline run
type TaskRunResults []TaskRunResult

// GetTaskRunResultsFinishedAt returns latest finishedAt time from TaskRunResults.
func (trrs TaskRunResults) GetTaskRunResultsFinishedAt() time.Time {
	var finishedTime time.Time
	for _, trr := range trrs {
		if trr.FinishedAt.Valid && trr.FinishedAt.Time.After(finishedTime) {
			finishedTime = trr.FinishedAt.Time
		}
	}
	return finishedTime
}

// FinalResult pulls the FinalResult for the pipeline_run from the task runs
// It needs to respect the output index of each task
func (trrs TaskRunResults) FinalResult(l logger.Logger) FinalResult {
	var found bool
	var fr FinalResult
	sort.Slice(trrs, func(i, j int) bool {
		return trrs[i].Task.OutputIndex() < trrs[j].Task.OutputIndex()
	})
	for _, trr := range trrs {
		fr.AllErrors = append(fr.AllErrors, trr.Result.Error)
		if trr.IsTerminal() {
			fr.Values = append(fr.Values, trr.Result.Value)
			fr.FatalErrors = append(fr.FatalErrors, trr.Result.Error)
			found = true
		}
	}

	if !found {
		l.Panicw("Expected at least one task to be final", "tasks", trrs)
	}
	return fr
}

// Terminals returns all terminal task run results
func (trrs TaskRunResults) Terminals() (terminals []TaskRunResult) {
	for _, trr := range trrs {
		if trr.IsTerminal() {
			terminals = append(terminals, trr)
		}
	}
	return
}

// GetNextTaskOf returns the task with the next id or nil if it does not exist
func (trrs *TaskRunResults) GetNextTaskOf(task TaskRunResult) *TaskRunResult {
	nextID := task.Task.Base().id + 1

	for _, trr := range *trrs {
		if trr.Task.Base().id == nextID {
			return &trr
		}
	}

	return nil
}

type TaskType string

func (t TaskType) String() string {
	return string(t)
}

const (
	TaskTypeAny              TaskType = "any"
	TaskTypeBase64Decode     TaskType = "base64decode"
	TaskTypeBase64Encode     TaskType = "base64encode"
	TaskTypeBridge           TaskType = "bridge"
	TaskTypeCBORParse        TaskType = "cborparse"
	TaskTypeConditional      TaskType = "conditional"
	TaskTypeDivide           TaskType = "divide"
	TaskTypeETHABIDecode     TaskType = "ethabidecode"
	TaskTypeETHABIDecodeLog  TaskType = "ethabidecodelog"
	TaskTypeETHABIEncode     TaskType = "ethabiencode"
	TaskTypeETHABIEncode2    TaskType = "ethabiencode2"
	TaskTypeETHCall          TaskType = "ethcall"
	TaskTypeETHTx            TaskType = "ethtx"
	TaskTypeEstimateGasLimit TaskType = "estimategaslimit"
	TaskTypeHTTP             TaskType = "http"
	TaskTypeHexDecode        TaskType = "hexdecode"
	TaskTypeHexEncode        TaskType = "hexencode"
	TaskTypeJSONParse        TaskType = "jsonparse"
	TaskTypeLength           TaskType = "length"
	TaskTypeLessThan         TaskType = "lessthan"
	TaskTypeLookup           TaskType = "lookup"
	TaskTypeLowercase        TaskType = "lowercase"
	TaskTypeMean             TaskType = "mean"
	TaskTypeMedian           TaskType = "median"
	TaskTypeMerge            TaskType = "merge"
	TaskTypeMode             TaskType = "mode"
	TaskTypeMultiply         TaskType = "multiply"
	TaskTypeSum              TaskType = "sum"
	TaskTypeUppercase        TaskType = "uppercase"
	TaskTypeVRF              TaskType = "vrf"
	TaskTypeVRFV2            TaskType = "vrfv2"
	TaskTypeVRFV2Plus        TaskType = "vrfv2plus"

	// Testing only.
	TaskTypePanic TaskType = "panic"
	TaskTypeMemo  TaskType = "memo"
	TaskTypeFail  TaskType = "fail"
)

var (
	stringType     = reflect.TypeOf("")
	bytesType      = reflect.TypeOf([]byte(nil))
	bytes20Type    = reflect.TypeOf([20]byte{})
	int32Type      = reflect.TypeOf(int32(0))
	nullUint32Type = reflect.TypeOf(cnull.Uint32{})
)

func UnmarshalTaskFromMap(taskType TaskType, taskMap interface{}, ID int, dotID string) (_ Task, err error) {
	defer cutils.WrapIfError(&err, "UnmarshalTaskFromMap")

	switch taskMap.(type) {
	default:
		return nil, pkgerrors.Errorf("UnmarshalTaskFromMap only accepts a map[string]interface{} or a map[string]string. Got %v (%#v) of type %T", taskMap, taskMap, taskMap)
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
	case TaskTypeMemo:
		task = &MemoTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeMultiply:
		task = &MultiplyTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeDivide:
		task = &DivideTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeVRF:
		task = &VRFTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeVRFV2:
		task = &VRFTaskV2{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeVRFV2Plus:
		task = &VRFTaskV2Plus{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeEstimateGasLimit:
		task = &EstimateGasLimitTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeETHCall:
		task = &ETHCallTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeETHTx:
		task = &ETHTxTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeETHABIEncode:
		task = &ETHABIEncodeTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeETHABIEncode2:
		task = &ETHABIEncodeTask2{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeETHABIDecode:
		task = &ETHABIDecodeTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeETHABIDecodeLog:
		task = &ETHABIDecodeLogTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeCBORParse:
		task = &CBORParseTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeFail:
		task = &FailTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeMerge:
		task = &MergeTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeLength:
		task = &LengthTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeLessThan:
		task = &LessThanTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeLookup:
		task = &LookupTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeLowercase:
		task = &LowercaseTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeUppercase:
		task = &UppercaseTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeConditional:
		task = &ConditionalTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeHexDecode:
		task = &HexDecodeTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeHexEncode:
		task = &HexEncodeTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeBase64Decode:
		task = &Base64DecodeTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	case TaskTypeBase64Encode:
		task = &Base64EncodeTask{BaseTask: BaseTask{id: ID, dotID: dotID}}
	default:
		return nil, pkgerrors.Errorf(`unknown task type: "%v"`, taskType)
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           task,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
				if from != stringType {
					return data, nil
				}
				switch to {
				case nullUint32Type:
					i, err2 := strconv.ParseUint(data.(string), 10, 32)
					return cnull.Uint32From(uint32(i)), err2
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
		return nil, pkgerrors.Wrapf(ErrWrongInputCardinality, "min: %v max: %v (got %v)", minLen, maxLen, len(inputs))
	} else if maxLen >= 0 && len(inputs) > maxLen {
		return nil, pkgerrors.Wrapf(ErrWrongInputCardinality, "min: %v max: %v (got %v)", minLen, maxLen, len(inputs))
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

var ErrInvalidEVMChainID = errors.New("invalid EVM chain ID")

func SelectGasLimit(ge config.GasEstimator, jobType string, specGasLimit *uint32) uint64 {
	if specGasLimit != nil {
		return uint64(*specGasLimit)
	}

	jt := ge.LimitJobType()
	var jobTypeGasLimit *uint32
	switch jobType {
	case DirectRequestJobType:
		jobTypeGasLimit = jt.DR()
	case FluxMonitorJobType:
		jobTypeGasLimit = jt.FM()
	case OffchainReportingJobType:
		jobTypeGasLimit = jt.OCR()
	case OffchainReporting2JobType:
		jobTypeGasLimit = jt.OCR2()
	case KeeperJobType:
		jobTypeGasLimit = jt.Keeper()
	case VRFJobType:
		jobTypeGasLimit = jt.VRF()
	}

	if jobTypeGasLimit != nil {
		return uint64(*jobTypeGasLimit)
	}
	return ge.LimitDefault()
}

func selectBlock(block string) (string, error) {
	if block == "" {
		return "latest", nil
	}
	block = strings.ToLower(block)
	if block == "pending" || block == "latest" {
		return block, nil
	}
	return "", pkgerrors.Errorf("unsupported block param: %s", block)
}
