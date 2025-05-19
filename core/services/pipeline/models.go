package pipeline

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type Spec struct {
	ID                int32
	DotDagSource      string          `json:"dotDagSource"`
	CreatedAt         time.Time       `json:"-"`
	MaxTaskDuration   models.Interval `json:"-"`
	GasLimit          *uint32         `json:"-"`
	ForwardingAllowed bool            `json:"-"`

	JobID   int32  `json:"-"`
	JobName string `json:"-"`
	JobType string `json:"-"`

	Pipeline *Pipeline `json:"-" db:"-"` // This may be nil, or may be populated manually as a cache. There is no locking on this, so be careful
}

func (s *Spec) GetOrParsePipeline() (*Pipeline, error) {
	if s.Pipeline != nil {
		return s.Pipeline, nil
	}
	return s.ParsePipeline()
}

func (s *Spec) ParsePipeline() (*Pipeline, error) {
	return Parse(s.DotDagSource)
}

type Run struct {
	ID             int64                             `json:"-"`
	JobID          int32                             `json:"-"`
	PipelineSpecID int32                             `json:"-"`
	PruningKey     int32                             `json:"-"` // This currently refers to the upstream job ID
	PipelineSpec   Spec                              `json:"pipelineSpec"`
	Meta           jsonserializable.JSONSerializable `json:"meta"`
	// The errors are only ever strings
	// DB example: [null, null, "my error"]
	AllErrors   RunErrors                         `json:"all_errors"`
	FatalErrors RunErrors                         `json:"fatal_errors"`
	Inputs      jsonserializable.JSONSerializable `json:"inputs"`
	// Its expected that Output.Val is of type []interface{}.
	// DB example: [1234, {"a": 10}, null]
	Outputs          jsonserializable.JSONSerializable `json:"outputs"`
	CreatedAt        time.Time                         `json:"createdAt"`
	FinishedAt       null.Time                         `json:"finishedAt"`
	PipelineTaskRuns []TaskRun                         `json:"taskRuns"`
	State            RunStatus                         `json:"state"`

	Pending bool
	// FailSilently is used to signal that a task with the failEarly flag has failed, and we want to not put this in the db
	FailSilently bool
}

func (r Run) GetID() string {
	return strconv.FormatInt(r.ID, 10)
}

func (r *Run) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	r.ID = ID
	return nil
}

func (r Run) HasFatalErrors() bool {
	for _, err := range r.FatalErrors {
		if !err.IsZero() {
			return true
		}
	}
	return false
}

func (r Run) HasErrors() bool {
	for _, err := range r.AllErrors {
		if !err.IsZero() {
			return true
		}
	}
	return false
}

// Status determines the status of the run.
func (r *Run) Status() RunStatus {
	if r.HasFatalErrors() {
		return RunStatusErrored
	} else if r.FinishedAt.Valid {
		return RunStatusCompleted
	}

	return RunStatusRunning
}

func (r *Run) ByDotID(id string) *TaskRun {
	for i, run := range r.PipelineTaskRuns {
		if run.DotID == id {
			return &r.PipelineTaskRuns[i]
		}
	}
	return nil
}

func (r *Run) StringOutputs() ([]*string, error) {
	// The UI expects all outputs to be strings.
	var outputs []*string
	// Note for async jobs, Outputs can be nil/invalid
	if r.Outputs.Valid {
		outs, ok := r.Outputs.Val.([]interface{})
		if !ok {
			return nil, fmt.Errorf("unable to process output type %T", r.Outputs.Val)
		}

		if r.Outputs.Valid && r.Outputs.Val != nil {
			for _, out := range outs {
				switch v := out.(type) {
				case string:
					s := v
					outputs = append(outputs, &s)
				case map[string]interface{}:
					b, _ := json.Marshal(v)
					bs := string(b)
					outputs = append(outputs, &bs)
				case decimal.Decimal:
					s := v.String()
					outputs = append(outputs, &s)
				case *decimal.Decimal:
					s := v.String()
					outputs = append(outputs, &s)
				case big.Int:
					s := v.String()
					outputs = append(outputs, &s)
				case *big.Int:
					s := v.String()
					outputs = append(outputs, &s)
				case int8, uint8, int16, uint16, int32, uint32, int64, uint64:
					s := fmt.Sprintf("%v", v)
					outputs = append(outputs, &s)
				case float64:
					s := strconv.FormatFloat(v, 'f', -1, 64)
					outputs = append(outputs, &s)
				case nil:
					outputs = append(outputs, nil)
				default:
					return nil, fmt.Errorf("unable to process output type %T", out)
				}
			}
		}
	}

	return outputs, nil
}

func (r *Run) StringFatalErrors() []*string {
	var fatalErrors []*string

	for _, err := range r.FatalErrors {
		if err.Valid {
			s := err.String
			fatalErrors = append(fatalErrors, &s)
		} else {
			fatalErrors = append(fatalErrors, nil)
		}
	}

	return fatalErrors
}

func (r *Run) StringAllErrors() []*string {
	var allErrors []*string

	for _, err := range r.AllErrors {
		if err.Valid {
			s := err.String
			allErrors = append(allErrors, &s)
		} else {
			allErrors = append(allErrors, nil)
		}
	}

	return allErrors
}

type RunErrors []null.String

func (re *RunErrors) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.Errorf("RunErrors#Scan received a value of type %T", value)
	}
	return json.Unmarshal(bytes, re)
}

func (re RunErrors) Value() (driver.Value, error) {
	if len(re) == 0 {
		return nil, nil
	}
	return json.Marshal(re)
}

func (re RunErrors) HasError() bool {
	for _, e := range re {
		if !e.IsZero() {
			return true
		}
	}
	return false
}

// ToError coalesces all non-nil errors into a single error object.
// This is useful for logging.
func (re RunErrors) ToError() error {
	toErr := func(ns null.String) error {
		if !ns.IsZero() {
			return errors.New(ns.String)
		}
		return nil
	}
	errs := []error{}
	for _, e := range re {
		errs = append(errs, toErr(e))
	}
	return multierr.Combine(errs...)
}

type ResumeRequest struct {
	Error null.String     `json:"error"`
	Value json.RawMessage `json:"value"`
}

func (rr ResumeRequest) ToResult() (Result, error) {
	var res Result
	if rr.Error.Valid && rr.Value == nil {
		res.Error = errors.New(rr.Error.ValueOrZero())
		return res, nil
	}
	if !rr.Error.Valid && rr.Value != nil {
		res.Value = []byte(rr.Value)
		return res, nil
	}
	return Result{}, errors.New("must provide only one of either 'value' or 'error' key")
}

type TaskRun struct {
	ID            uuid.UUID                         `json:"id"`
	Type          TaskType                          `json:"type"`
	PipelineRun   Run                               `json:"-"`
	PipelineRunID int64                             `json:"-"`
	Output        jsonserializable.JSONSerializable `json:"output"`
	Error         null.String                       `json:"error"`
	CreatedAt     time.Time                         `json:"createdAt"`
	FinishedAt    null.Time                         `json:"finishedAt"`
	Index         int32                             `json:"index"`
	DotID         string                            `json:"dotId"`

	// Used internally for sorting completed results
	task Task
}

func (tr TaskRun) GetID() string {
	return fmt.Sprintf("%v", tr.ID)
}

func (tr *TaskRun) SetID(value string) error {
	ID, err := uuid.Parse(value)
	if err != nil {
		return err
	}
	tr.ID = ID
	return nil
}

func (tr TaskRun) GetDotID() string {
	return tr.DotID
}

func (tr TaskRun) Result() Result {
	var result Result
	if !tr.Error.IsZero() {
		result.Error = errors.New(tr.Error.ValueOrZero())
	} else if tr.Output.Valid && tr.Output.Val != nil {
		result.Value = tr.Output.Val
	}
	return result
}

func (tr *TaskRun) IsPending() bool {
	return !tr.FinishedAt.Valid && tr.Output.Empty() && tr.Error.IsZero()
}

// RunStatus represents the status of a run
type RunStatus string

const (
	// RunStatusUnknown is the when the run status cannot be determined.
	RunStatusUnknown RunStatus = "unknown"
	// RunStatusRunning is used for when a run is actively being executed.
	RunStatusRunning RunStatus = "running"
	// RunStatusSuspended is used when a run is paused and awaiting further results.
	RunStatusSuspended RunStatus = "suspended"
	// RunStatusErrored is used for when a run has errored and will not complete.
	RunStatusErrored RunStatus = "errored"
	// RunStatusCompleted is used for when a run has successfully completed execution.
	RunStatusCompleted RunStatus = "completed"
)

// Completed returns true if the status is RunStatusCompleted.
func (s RunStatus) Completed() bool {
	return s == RunStatusCompleted
}

// Errored returns true if the status is RunStatusErrored.
func (s RunStatus) Errored() bool {
	return s == RunStatusErrored
}

// Finished returns true if the status is final and can't be changed.
func (s RunStatus) Finished() bool {
	return s.Completed() || s.Errored()
}
