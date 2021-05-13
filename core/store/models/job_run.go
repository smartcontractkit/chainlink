package models

import (
	"fmt"
	"math/big"
	"time"

	"gorm.io/gorm"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/assets"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	null "gopkg.in/guregu/null.v4"
)

var (
	promTotalRunUpdates = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "run_status_update_total",
		Help: "The total number of status updates for Job Runs",
	},
		[]string{"job_spec_id", "from_status", "status"},
	)
)

// JobRun tracks the status of a job by holding its TaskRuns and the
// Result of each Run.
type JobRun struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;not null"`
	JobSpecID      JobID          `json:"jobId" gorm:"type:uuid"`
	Result         RunResult      `json:"result" gorm:"foreignkey:ResultID"`
	ResultID       clnull.Int64   `json:"-"`
	RunRequest     RunRequest     `json:"-" gorm:"foreignkey:RunRequestID"`
	RunRequestID   clnull.Int64   `json:"-"`
	Status         RunStatus      `json:"status" gorm:"default:'unstarted'"`
	TaskRuns       []TaskRun      `json:"taskRuns" gorm:"foreignKey:JobRunID"`
	CreatedAt      time.Time      `json:"createdAt"`
	FinishedAt     null.Time      `json:"finishedAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	Initiator      Initiator      `json:"initiator" gorm:"foreignkey:InitiatorID;->"`
	InitiatorID    int64          `json:"-"`
	CreationHeight *utils.Big     `json:"creationHeight"`
	ObservedHeight *utils.Big     `json:"observedHeight"`
	DeletedAt      gorm.DeletedAt `json:"-"`
	Payment        *assets.Link   `json:"payment,omitempty"`
}

// MakeJobRun returns a new JobRun copy
func MakeJobRun(job *JobSpec, now time.Time, initiator *Initiator, currentHeight *big.Int, runRequest *RunRequest) JobRun {
	run := JobRun{
		ID:          uuid.NewV4(),
		JobSpecID:   job.ID,
		CreatedAt:   now,
		UpdatedAt:   now,
		Initiator:   *initiator,
		InitiatorID: initiator.ID,
		TaskRuns:    make([]TaskRun, len(job.Tasks)),
		RunRequest:  *runRequest,
		Payment:     runRequest.Payment,
	}
	if currentHeight != nil {
		run.CreationHeight = utils.NewBig(currentHeight)
		run.ObservedHeight = utils.NewBig(currentHeight)
	}
	for i, task := range job.Tasks {
		run.TaskRuns[i] = TaskRun{
			ID:         uuid.NewV4(),
			JobRunID:   run.ID,
			TaskSpec:   task,
			TaskSpecID: task.ID,
			Status:     RunStatusUnstarted,
		}
	}
	run.SetStatus(RunStatusInProgress)
	return run
}

// GetID returns the ID of this structure for jsonapi serialization.
func (jr JobRun) GetID() string {
	return jr.ID.String()
}

// GetName returns the pluralized "type" of this structure for jsonapi serialization.
func (jr JobRun) GetName() string {
	return "runs"
}

// SetStatus updates run status.
func (jr *JobRun) SetStatus(status RunStatus) {
	oldStatus := jr.Status
	jr.Status = status
	if jr.Status.Completed() && jr.TasksRemain() {
		jr.Status = RunStatusInProgress
	} else if jr.Status.Finished() {
		jr.FinishedAt = null.TimeFrom(time.Now())
	}
	promTotalRunUpdates.WithLabelValues(jr.JobSpecID.String(), string(oldStatus), string(status)).Inc()
}

// GetStatus returns the JobRun's RunStatus
func (jr *JobRun) GetStatus() RunStatus {
	return jr.Status
}

// SetID is used to set the ID of this structure when deserializing from jsonapi documents.
func (jr *JobRun) SetID(value string) error {
	return jr.ID.UnmarshalText([]byte(value))
}

// ForLogger formats the JobRun for a common formatting in the log.
func (jr JobRun) ForLogger(kvs ...interface{}) []interface{} {
	output := []interface{}{
		"jobID", jr.JobSpecID.String(),
		"runID", jr.ID.String(),
		"status", jr.Status,
	}

	if jr.CreationHeight != nil {
		output = append(output, "creation_height", jr.CreationHeight.ToInt())
	}

	if jr.ObservedHeight != nil {
		output = append(output, "observed_height", jr.ObservedHeight.ToInt())
	}

	if jr.HasError() {
		output = append(output, "job_error", jr.ErrorString())
	}

	if jr.Status.Completed() {
		output = append(output, "link_earned", jr.Payment)
	} else {
		output = append(output, "input_amount", jr.Payment)
	}

	if jr.RunRequest.RequestID != nil {
		output = append(output, "requestID", jr.RunRequest.RequestID)
	}

	if jr.RunRequest.TxHash != nil {
		output = append(output, "txHash", jr.RunRequest.TxHash)
	}

	if jr.RunRequest.BlockHash != nil {
		output = append(output, "blockHash", jr.RunRequest.BlockHash)
	}

	return append(kvs, output...)
}

// HasError returns true if this JobRun has errored
func (jr JobRun) HasError() bool {
	return jr.Status.Errored()
}

// NextTaskRunIndex returns the position of the next unfinished task
func (jr *JobRun) NextTaskRunIndex() (int, bool) {
	for index, tr := range jr.TaskRuns {
		if tr.Status.Pending() || tr.Status == RunStatusUnstarted || tr.Status == RunStatusInProgress {
			return index, true
		}
	}
	return 0, false
}

// NextTaskRun returns the next immediate TaskRun in the list
// of unfinished TaskRuns.
func (jr *JobRun) NextTaskRun() *TaskRun {
	nextTaskIndex, runnable := jr.NextTaskRunIndex()
	if runnable {
		return &jr.TaskRuns[nextTaskIndex]
	}
	return nil
}

// PreviousTaskRun returns the last task to be processed, if it exists
func (jr *JobRun) PreviousTaskRun() *TaskRun {
	index, runnable := jr.NextTaskRunIndex()
	if runnable && index > 0 {
		return &jr.TaskRuns[index-1]
	}
	return nil
}

// TasksRemain returns true if there are unfinished tasks left for this job run
func (jr *JobRun) TasksRemain() bool {
	_, runnable := jr.NextTaskRunIndex()
	return runnable
}

// SetError sets this job run to failed and saves the error message
func (jr *JobRun) SetError(err error) {
	jr.Result.ErrorMessage = null.StringFrom(err.Error())
	jr.SetStatus(RunStatusErrored)
}

// Cancel sets this run as cancelled, it should no longer be processed.
func (jr *JobRun) Cancel() {
	currentTaskRun := jr.NextTaskRun()
	if currentTaskRun != nil {
		currentTaskRun.Status = RunStatusCancelled
	}
	jr.SetStatus(RunStatusCancelled)
}

// ApplyOutput updates the JobRun's Result and Status
func (jr *JobRun) ApplyOutput(result RunOutput) {
	if result.HasError() {
		jr.SetError(result.Error())
		return
	}
	jr.Result.Data = result.Data()
	jr.SetStatus(result.Status())
}

// ApplyBridgeRunResult saves the input from a BridgeAdapter
func (jr *JobRun) ApplyBridgeRunResult(result BridgeRunResult) {
	if result.HasError() {
		jr.SetError(result.GetError())
	}
	jr.Result.Data = result.Data
	jr.SetStatus(result.Status)
}

// ErrorString returns the error as a string if present, otherwise "".
func (jr *JobRun) ErrorString() string {
	return jr.Result.ErrorMessage.ValueOrZero()
}

// RunRequest stores the fields used to initiate the parent job run.
type RunRequest struct {
	ID            int64 `gorm:"primary_key"`
	RequestID     *common.Hash
	TxHash        *common.Hash
	BlockHash     *common.Hash
	Requester     *common.Address
	CreatedAt     time.Time
	Payment       *assets.Link
	RequestParams JSON `gorm:"type:jsonb;default:'{}'"`
}

// NewRunRequest returns a new RunRequest instance.
func NewRunRequest(requestParams JSON) *RunRequest {
	return &RunRequest{CreatedAt: time.Now(), RequestParams: requestParams}
}

// TaskRun stores the Task and represents the status of the
// Task to be ran.
type TaskRun struct {
	ID                               uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;not null"`
	JobRunID                         uuid.UUID     `json:"-" gorm:"type:uuid"`
	Result                           RunResult     `json:"result"`
	ResultID                         clnull.Int64  `json:"-"`
	Status                           RunStatus     `json:"status" gorm:"default:'unstarted'"`
	TaskSpec                         TaskSpec      `json:"task" gorm:"->"`
	TaskSpecID                       int64         `json:"-"`
	MinRequiredIncomingConfirmations clnull.Uint32 `json:"minimumConfirmations" gorm:"column:minimum_confirmations"`
	ObservedIncomingConfirmations    clnull.Uint32 `json:"confirmations" gorm:"column:confirmations"`
	CreatedAt                        time.Time     `json:"-"`
	UpdatedAt                        time.Time     `json:"-"`
}

// String returns info on the TaskRun as "ID,Type,Status,Result".
func (tr TaskRun) String() string {
	return fmt.Sprintf("TaskRun(%v,%v,%v,%v)", tr.ID.String(), tr.TaskSpec.Type, tr.Status, tr.Result)
}

// SetError sets this task run to failed and saves the error message
func (tr *TaskRun) SetError(err error) {
	tr.Result.ErrorMessage = null.StringFrom(err.Error())
	tr.Status = RunStatusErrored
}

// ApplyBridgeRunResult updates the TaskRun's Result and Status
func (tr *TaskRun) ApplyBridgeRunResult(result BridgeRunResult) {
	if result.HasError() {
		tr.SetError(result.GetError())
	}
	tr.Result.Data = result.Data
	tr.Status = result.Status
}

// ApplyOutput updates the TaskRun's Result and Status
func (tr *TaskRun) ApplyOutput(result RunOutput) {
	if result.HasError() {
		tr.SetError(result.Error())
		return
	}
	tr.Result.Data = result.Data()
	tr.Status = result.Status()
}

// RunResult keeps track of the outcome of a TaskRun or JobRun. It stores the
// Data and ErrorMessage.
type RunResult struct {
	ID           int64       `json:"-" gorm:"primary_key;auto_increment"`
	Data         JSON        `json:"data"`
	ErrorMessage null.String `json:"error"`
	CreatedAt    time.Time   `json:"-"`
	UpdatedAt    time.Time   `json:"-"`
}
